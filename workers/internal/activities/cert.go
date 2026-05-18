package activities

import (
	"context"
	"fmt"
	"time"

	"temporal-progress-mock/workers/internal/shared"

	"go.temporal.io/sdk/activity"
)

func RequestCertificate(ctx context.Context, input shared.ProvisionInput) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "requesting certificate")
	time.Sleep(3 * time.Second)

	return shared.ActivityResult{
		Message: fmt.Sprintf("Certificate issued for %s.%s.svc", input.AppName, input.Namespace),
		Value:   fmt.Sprintf("%s-tls", input.AppName),
	}, nil
}

func DeployApplication(ctx context.Context, input shared.ProvisionInput) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "applying manifests")
	time.Sleep(3 * time.Second)

	activity.RecordHeartbeat(ctx, "waiting for pods")
	time.Sleep(3 * time.Second)

	return shared.ActivityResult{
		Message: fmt.Sprintf("Application %s deployed", input.AppName),
		Value:   input.AppName,
	}, nil
}

func VerifyApplication(ctx context.Context, input shared.ProvisionInput) (shared.ActivityResult, error) {
	activity.RecordHeartbeat(ctx, "checking service health")
	time.Sleep(2 * time.Second)

	return shared.ActivityResult{
		Message: fmt.Sprintf("Application %s is healthy", input.AppName),
		Value:   "healthy",
	}, nil
}
