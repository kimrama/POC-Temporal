package workflows

import (
	"fmt"
	"time"

	"temporal-progress-mock/workers/shared"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	ComputerTaskQueue = "computer-task-queue"
)

func ComputeWorkflow(ctx workflow.Context, input shared.ComputeInput) (string, error) {
	state := shared.ProgressState{
		Status: "running",
		Steps: []shared.WorkflowStep{
			{Key: "initial_values", Label: "Set initial values", Status: shared.StepPending},
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

	if err := runActivity(ComputerTaskQueue, "initial_values", "SetInitialValues"); err != nil {
		return "", err
	}
	state.Status = "completed"
	state.CurrentStep = ""
	return "Compute workflow completed successfully", nil
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

// func runCompensation(ctx workflow.Context, state *shared.ProgressState, input shared.ComputeInput) error {
// 	state.CurrentStep = "rollback"
// 	state.Steps = append(state.Steps,
// 		shared.WorkflowStep{
// 			Key:     "rollback_secret",
// 			Label:   "Rollback secret",
// 			Status:  shared.StepRunning,
// 			Message: "Deleting secret",
// 		},
// 	)

// 	activityOptions := workflow.ActivityOptions{
// 		TaskQueue:              K8STaskQueue,
// 		StartToCloseTimeout:    30 * time.Second,
// 		ScheduleToCloseTimeout: time.Minute,
// 	}
// 	activityCtx := workflow.WithActivityOptions(ctx, activityOptions)

// 	var result shared.ActivityResult
// 	err := workflow.ExecuteActivity(activityCtx, "DeleteSecret", input).Get(activityCtx, &result)
// 	if err != nil {
// 		state.Error = fmt.Sprintf("rollback failed: %s", err.Error())
// 		return err
// 	}

// 	last := len(state.Steps) - 1
// 	state.Steps[last].Status = shared.StepRollback
// 	state.Steps[last].Message = result.Message
// 	state.Steps[last].CompletedAt = workflow.Now(ctx).Format(time.RFC3339)

// 	return nil
// }
