package workflows

import (
	"fmt"
	"time"

	"temporal-progress-mock/workers/internal/shared"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	K8STaskQueue  = "k8s-task-queue"
	CertTaskQueue = "cert-task-queue"
)

func ProvisionWorkflow(ctx workflow.Context, input shared.ProvisionInput) (string, error) {
	state := shared.ProgressState{
		Status: "running",
		Steps: []shared.WorkflowStep{
			{Key: "create_namespace", Label: "Create namespace", Status: shared.StepPending},
			{Key: "create_secret", Label: "Create secret", Status: shared.StepPending},
			{Key: "request_certificate", Label: "Request certificate", Status: shared.StepPending},
			{Key: "deploy_application", Label: "Deploy application", Status: shared.StepPending},
			{Key: "verify_application", Label: "Verify application", Status: shared.StepPending},
		},
	}

	info := workflow.GetInfo(ctx)
	state.WorkflowID = info.WorkflowExecution.ID

	if err := workflow.SetQueryHandler(ctx, "get_progress", func() (shared.ProgressState, error) {
		return state, nil
	}); err != nil {
		return "", err
	}

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    5 * time.Second,
		MaximumAttempts:    3,
	}

	runActivity := func(taskQueue string, stepKey string, activityName string) error {
		setStep(ctx, &state, stepKey, shared.StepRunning, fmt.Sprintf("Running %s", activityName))

		activityOptions := workflow.ActivityOptions{
			TaskQueue:              taskQueue,
			StartToCloseTimeout:    30 * time.Second,
			ScheduleToCloseTimeout: 2 * time.Minute,
			HeartbeatTimeout:       10 * time.Second,
			RetryPolicy:            retryPolicy,
		}

		activityCtx := workflow.WithActivityOptions(ctx, activityOptions)

		var result shared.ActivityResult
		err := workflow.ExecuteActivity(activityCtx, activityName, input).Get(activityCtx, &result)
		if err != nil {
			setStep(ctx, &state, stepKey, shared.StepFailed, err.Error())
			state.Status = "failed"
			state.Error = err.Error()
			return err
		}

		setStep(ctx, &state, stepKey, shared.StepCompleted, result.Message)
		return nil
	}

	if err := runActivity(K8STaskQueue, "create_namespace", "CreateNamespace"); err != nil {
		return "", err
	}

	if err := runActivity(K8STaskQueue, "create_secret", "CreateSecret"); err != nil {
		return "", err
	}

	if err := runActivity(CertTaskQueue, "request_certificate", "RequestCertificate"); err != nil {
		_ = runCompensation(ctx, &state, input)
		return "", err
	}

	if err := runActivity(CertTaskQueue, "deploy_application", "DeployApplication"); err != nil {
		_ = runCompensation(ctx, &state, input)
		return "", err
	}

	if err := runActivity(CertTaskQueue, "verify_application", "VerifyApplication"); err != nil {
		_ = runCompensation(ctx, &state, input)
		return "", err
	}

	state.Status = "completed"
	state.CurrentStep = ""
	return fmt.Sprintf("Provisioned %s in namespace %s on cluster %s", input.AppName, input.Namespace, input.Cluster), nil
}

func setStep(ctx workflow.Context, state *shared.ProgressState, key string, status shared.StepStatus, message string) {
	now := workflow.Now(ctx).Format(time.RFC3339)
	state.CurrentStep = key

	for i := range state.Steps {
		if state.Steps[i].Key != key {
			continue
		}

		state.Steps[i].Status = status
		state.Steps[i].Message = message

		if status == shared.StepRunning {
			state.Steps[i].StartedAt = now
		}

		if status == shared.StepCompleted || status == shared.StepFailed || status == shared.StepRollback {
			state.Steps[i].CompletedAt = now
		}

		return
	}
}

func runCompensation(ctx workflow.Context, state *shared.ProgressState, input shared.ProvisionInput) error {
	state.CurrentStep = "rollback"
	state.Steps = append(state.Steps,
		shared.WorkflowStep{
			Key:     "rollback_secret",
			Label:   "Rollback secret",
			Status:  shared.StepRunning,
			Message: "Deleting secret",
		},
	)

	activityOptions := workflow.ActivityOptions{
		TaskQueue:              K8STaskQueue,
		StartToCloseTimeout:    30 * time.Second,
		ScheduleToCloseTimeout: time.Minute,
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOptions)

	var result shared.ActivityResult
	err := workflow.ExecuteActivity(activityCtx, "DeleteSecret", input).Get(activityCtx, &result)
	if err != nil {
		state.Error = fmt.Sprintf("rollback failed: %s", err.Error())
		return err
	}

	last := len(state.Steps) - 1
	state.Steps[last].Status = shared.StepRollback
	state.Steps[last].Message = result.Message
	state.Steps[last].CompletedAt = workflow.Now(ctx).Format(time.RFC3339)

	return nil
}
