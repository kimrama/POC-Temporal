package workflows

import (
	"temporal-progress-mock/shared"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"time"
)

const (
	ComputerTaskQueue = "computer-task-queue"
)

func ComputeWorkflow(ctx workflow.Context, input shared.ComputeInput) (string, error) {

	state := shared.ProgressState{
		Status: "running",
		Steps: []shared.ComputeStep{
			{Key: "set_initial_values", Label: "Set Initial Values", Status: shared.StepPending},
			{Key: "plus_one", Label: "Add 1", Status: shared.StepPending},
			{Key: "times_two", Label: "Multiply by 2", Status: shared.StepPending},
		},
	}

	info := workflow.GetInfo(ctx)
	state.WorkflowID = info.WorkflowExecution.ID

	if err := workflow.SetQueryHandler(ctx, "get_progress", func() (shared.ProgressState, error) {
		return state, nil
	}); err != nil {
		return "", err
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 3},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	compensations := []shared.Compensation{}
	addCompensation := func(name string, do interface{}) {
		compensations = append(compensations, shared.Compensation{Name: name, Do: do})
	}

	rollback := func() {
		for i := len(compensations) - 1; i >= 0; i-- {
			compensation := compensations[i]

			err := workflow.ExecuteActivity(ctx, compensation.Do).Get(ctx, nil)
			if err != nil {
				workflow.GetLogger(ctx).Error("compensation failed", "compensation", compensation.Name, "error", err)
			} else {
				workflow.GetLogger(ctx).Info("compensation succeeded", "compensation", compensation.Name)
			}
		}
	}

	// Step 1: Set initial values

	setStep(ctx, &state, "set_initial_values", shared.StepRunning, "Setting initial values")
	err := workflow.ExecuteActivity(ctx, "SetInitialValues", input).Get(ctx, nil)
	if err != nil {
		setStep(ctx, &state, "set_initial_values", shared.StepFailed, err.Error())
		return "", err
	}
	setStep(ctx, &state, "set_initial_values", shared.StepCompleted, "Initial values set")

	addCompensation("ResetValue", "ResetValue")

	// Step 2: Add 1

	setStep(ctx, &state, "plus_one", shared.StepRunning, "Adding 1")
	err = workflow.ExecuteActivity(ctx, "PlusOne").Get(ctx, nil)
	if err != nil {
		setStep(ctx, &state, "plus_one", shared.StepFailed, err.Error())
		rollback()
		return "", err
	}
	setStep(ctx, &state, "plus_one", shared.StepCompleted, "Added 1")
	addCompensation("PlusOne", "MinusOne")

	// Step 3: Multiply by 2
	setStep(ctx, &state, "times_two", shared.StepRunning, "Multiplying by 2")
	err = workflow.ExecuteActivity(ctx, "TimesTwo").Get(ctx, nil)
	if err != nil {
		setStep(ctx, &state, "times_two", shared.StepFailed, err.Error())
		rollback()
		return "", err
	}
	setStep(ctx, &state, "times_two", shared.StepCompleted, "Multiplied by 2")
	addCompensation("TimesTwo", "DivideByTwo")

	return "Computation completed successfully", nil
}

func setStep(ctx workflow.Context, state *shared.ProgressState, stepKey string, status shared.StepStatus, message string) {
	now := workflow.Now(ctx).Format(time.RFC3339)
	state.CurrentStep = stepKey

	if status == shared.StepFailed {
		state.Status = "failed"
		state.Error = message
	}

	for i := range state.Steps {
		if state.Steps[i].Key != stepKey {
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
