package activities

import (
	"context"
	"fmt"
	"os"
	"time"

	"temporal-progress-mock/shared"

	"go.temporal.io/sdk/activity"
)

const VALUE_FILE_PATH = "value.txt"

func ResetValue(ctx context.Context) error {
	return os.WriteFile(VALUE_FILE_PATH, []byte("0"), 0644)
}

func SetInitialValues(ctx context.Context, input shared.ComputeInput) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "setting initial values")
	time.Sleep(5 * time.Second)

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
		Value:   fmt.Sprintf("%d", input.InitialValue),
	}, nil

}
func MinusOne(ctx context.Context) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "subtracting 1 from value")
	time.Sleep(5 * time.Second)

	valueFile, err := os.OpenFile(VALUE_FILE_PATH, os.O_RDWR, 0644)
	if err != nil {
		return shared.ActivityResult{}, err
	}
	defer valueFile.Close()

	var currentValue int
	_, err = fmt.Fscanf(valueFile, "%d", &currentValue)
	if err != nil {
		return shared.ActivityResult{
			Message: "failed to read current value from file",
		}, err
	}

	newValue := currentValue - 1

	_, err = valueFile.Seek(0, 0)
	if err != nil {
		return shared.ActivityResult{}, err
	}

	_, err = valueFile.WriteString(fmt.Sprintf("%d", newValue))
	if err != nil {
		return shared.ActivityResult{
			Message: "failed to write new value to file",
		}, err
	}

	return shared.ActivityResult{
		Message: fmt.Sprintf("Subtracted 1 from value, new value = %d", newValue),
		Value:   fmt.Sprintf("%d", newValue),
	}, nil
}

func PlusOne(ctx context.Context) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "adding 1 to value")
	time.Sleep(5 * time.Second)

	valueFile, err := os.OpenFile(VALUE_FILE_PATH, os.O_RDWR, 0644)
	if err != nil {
		return shared.ActivityResult{}, err
	}
	defer valueFile.Close()

	var currentValue int
	_, err = fmt.Fscanf(valueFile, "%d", &currentValue)
	if err != nil {
		return shared.ActivityResult{
			Message: "failed to read current value from file",
		}, err
	}

	newValue := currentValue + 1

	_, err = valueFile.Seek(0, 0)
	if err != nil {
		return shared.ActivityResult{}, err
	}

	_, err = valueFile.WriteString(fmt.Sprintf("%d", newValue))
	if err != nil {
		return shared.ActivityResult{
			Message: "failed to write new value to file",
		}, err
	}

	return shared.ActivityResult{
		Message: fmt.Sprintf("Added 1 to value, new value = %d", newValue),
		Value:   fmt.Sprintf("%d", newValue),
	}, nil
}

func DivideByTwo(ctx context.Context) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "dividing value by 2")
	time.Sleep(5 * time.Second)

	valueFile, err := os.OpenFile(VALUE_FILE_PATH, os.O_RDWR, 0644)
	if err != nil {
		return shared.ActivityResult{}, err
	}
	defer valueFile.Close()

	var currentValue int
	_, err = fmt.Fscanf(valueFile, "%d", &currentValue)
	if err != nil {
		return shared.ActivityResult{
			Message: "failed to read current value from file",
		}, err
	}

	newValue := currentValue / 2

	_, err = valueFile.Seek(0, 0)
	if err != nil {
		return shared.ActivityResult{}, err
	}

	_, err = valueFile.WriteString(fmt.Sprintf("%d", newValue))
	if err != nil {
		return shared.ActivityResult{
			Message: "failed to write new value to file",
		}, err
	}

	return shared.ActivityResult{
		Message: fmt.Sprintf("Divided value by 2, new value = %d", newValue),
		Value:   fmt.Sprintf("%d", newValue),
	}, nil
}
func TimesTwo(ctx context.Context) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "multiplying value by 2")
	time.Sleep(5 * time.Second)

	return shared.ActivityResult{
		Message: "simulating failure in TimesTwo activity",
	}, fmt.Errorf("simulated error in TimesTwo")

	// valueFile, err := os.OpenFile(VALUE_FILE_PATH, os.O_RDWR, 0644)
	// if err != nil {
	// 	return shared.ActivityResult{}, err
	// }
	// defer valueFile.Close()

	// var currentValue int
	// _, err = fmt.Fscanf(valueFile, "%d", &currentValue)
	// if err != nil {
	// 	return shared.ActivityResult{
	// 		Message: "failed to read current value from file",
	// 	}, err
	// }

	// newValue := currentValue * 2

	// _, err = valueFile.Seek(0, 0)
	// if err != nil {
	// 	return shared.ActivityResult{}, err
	// }

	// _, err = valueFile.WriteString(fmt.Sprintf("%d", newValue))
	// if err != nil {
	// 	return shared.ActivityResult{
	// 		Message: "failed to write new value to file",
	// 	}, err
	// }

	// return shared.ActivityResult{
	// 	Message: fmt.Sprintf("Multiplied value by 2, new value = %d", newValue),
	// 	Value:   fmt.Sprintf("%d", newValue),
	// }, nil
}
