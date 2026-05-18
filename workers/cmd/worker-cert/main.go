package main

import (
	"log"
	"os"

	"temporal-progress-mock/workers/internal/activities"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	address := env("TEMPORAL_ADDRESS", "localhost:7233")
	namespace := env("TEMPORAL_NAMESPACE", "default")
	certTaskQueue := env("CERT_TASK_QUEUE", "cert-task-queue")

	c, err := client.Dial(client.Options{
		HostPort:  address,
		Namespace: namespace,
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, certTaskQueue, worker.Options{})

	w.RegisterActivityWithOptions(activities.RequestCertificate, activity.RegisterOptions{Name: "RequestCertificate"})
	w.RegisterActivityWithOptions(activities.DeployApplication, activity.RegisterOptions{Name: "DeployApplication"})
	w.RegisterActivityWithOptions(activities.VerifyApplication, activity.RegisterOptions{Name: "VerifyApplication"})

	log.Printf("starting cert/deploy activity worker on task queue %s", certTaskQueue)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("unable to start worker", err)
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
