package main

import (
	"log"
	"os"

	"temporal-progress-mock/workers/internal/activities"
	"temporal-progress-mock/workers/internal/workflows"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	address := env("TEMPORAL_ADDRESS", "localhost:7233")
	namespace := env("TEMPORAL_NAMESPACE", "default")
	workflowTaskQueue := env("WORKFLOW_TASK_QUEUE", "provision-task-queue")
	k8sTaskQueue := env("K8S_TASK_QUEUE", "k8s-task-queue")

	c, err := client.Dial(client.Options{
		HostPort:  address,
		Namespace: namespace,
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	workflowWorker := worker.New(c, workflowTaskQueue, worker.Options{})
	workflowWorker.RegisterWorkflowWithOptions(workflows.ProvisionWorkflow, workflow.RegisterOptions{
		Name: "ProvisionWorkflow",
	})

	k8sWorker := worker.New(c, k8sTaskQueue, worker.Options{})
	k8sWorker.RegisterActivityWithOptions(activities.CreateNamespace, activity.RegisterOptions{Name: "CreateNamespace"})
	k8sWorker.RegisterActivityWithOptions(activities.CreateSecret, activity.RegisterOptions{Name: "CreateSecret"})
	k8sWorker.RegisterActivityWithOptions(activities.DeleteSecret, activity.RegisterOptions{Name: "DeleteSecret"})

	errCh := make(chan error, 2)

	go func() {
		log.Printf("starting workflow worker on task queue %s", workflowTaskQueue)
		errCh <- workflowWorker.Run(worker.InterruptCh())
	}()

	go func() {
		log.Printf("starting k8s activity worker on task queue %s", k8sTaskQueue)
		errCh <- k8sWorker.Run(worker.InterruptCh())
	}()

	if err := <-errCh; err != nil {
		log.Fatalln("worker stopped with error", err)
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
