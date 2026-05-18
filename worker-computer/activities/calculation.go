package activities

import (
	"context"
	"fmt"
	"os"
	"time"

	"temporal-progress-mock/workers/shared"

	"go.temporal.io/sdk/activity"
)

const VALUE_FILE_PATH = "value.txt"

func SetInitialValues(ctx context.Context, input shared.ComputeInput) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "setting initial values")
	time.Sleep(3 * time.Second)

	valueFile, err := os.OpenFile(VALUE_FILE_PATH, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return shared.ActivityResult{}, err
	}
	defer valueFile.Close()

	_, err = valueFile.WriteString(fmt.Sprintf("%d", input.InitialValue))
	if err != nil {
		return shared.ActivityResult{
			Message: "failed to write initial value to file",
		}, err
	}

	return shared.ActivityResult{
		Message: fmt.Sprintf("Initial values set for compute = %d", input.InitialValue),
		Value:   fmt.Sprintf("%d-tls", input.InitialValue),
	}, nil

}
