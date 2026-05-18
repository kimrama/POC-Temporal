package activities

import (
	"context"
	"fmt"
	"time"

	"temporal-progress-mock/workers/internal/shared"

	"go.temporal.io/sdk/activity"
)

func CreateNamespace(ctx context.Context, input shared.ProvisionInput) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "creating namespace")
	time.Sleep(2 * time.Second)

	return shared.ActivityResult{
		Message: fmt.Sprintf("Namespace %s created on cluster %s", input.Namespace, input.Cluster),
		Value:   input.Namespace,
	}, nil
}

func CreateSecret(ctx context.Context, input shared.ProvisionInput) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "creating secret")
	time.Sleep(2 * time.Second)

	return shared.ActivityResult{
		Message: fmt.Sprintf("Secret for %s created", input.AppName),
		Value:   fmt.Sprintf("%s-secret", input.AppName),
	}, nil
}

func DeleteSecret(ctx context.Context, input shared.ProvisionInput) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "deleting secret")
	time.Sleep(2 * time.Second)

	return shared.ActivityResult{
		Message: fmt.Sprintf("Secret for %s deleted", input.AppName),
	}, nil
}
