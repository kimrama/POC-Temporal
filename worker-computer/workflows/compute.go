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

	if err := workflow.SetQueryHandler(ctx, "progress", func() (shared.ProgressState, error) {
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
	err := workflow.ExecuteActivity(ctx, "SetInitialValues", input).Get(ctx, nil)
	if err != nil {
		return "", err
	}

	addCompensation("ResetValue", "ResetValue")

	// Step 2: Add 1
	err = workflow.ExecuteActivity(ctx, "PlusOne").Get(ctx, nil)
	if err != nil {
		rollback()
		return "", err
	}

	addCompensation("PlusOne", "MinusOne")

	// Step 3: Multiply by 2
	err = workflow.ExecuteActivity(ctx, "TimesTwo").Get(ctx, nil)
	if err != nil {
		rollback()
		return "", err
	}

	addCompensation("TimesTwo", "DivideByTwo")

	return "Computation completed successfully", nil
}
