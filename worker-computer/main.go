package main

import (
	"log"
	"os"

	"temporal-progress-mock/activities"
	"temporal-progress-mock/workflows"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	address := env("TEMPORAL_ADDRESS", "localhost:7233")
	namespace := env("TEMPORAL_NAMESPACE", "default")
	ComputerTaskQueue := env("COMPUTER_TASK_QUEUE", "computer-task-queue")

	c, err := client.Dial(client.Options{
		HostPort:  address,
		Namespace: namespace,
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	computerWorker := worker.New(c, ComputerTaskQueue, worker.Options{})

	computerWorker.RegisterWorkflowWithOptions(workflows.ComputeWorkflow, workflow.RegisterOptions{
		Name: "ComputeWorkflow",
	})

	computerWorker.RegisterActivityWithOptions(activities.SetInitialValues, activity.RegisterOptions{Name: "SetInitialValues"})
	computerWorker.RegisterActivityWithOptions(activities.ResetValue, activity.RegisterOptions{Name: "ResetValue"})
	computerWorker.RegisterActivityWithOptions(activities.PlusOne, activity.RegisterOptions{Name: "PlusOne"})
	computerWorker.RegisterActivityWithOptions(activities.TimesTwo, activity.RegisterOptions{Name: "TimesTwo"})
	computerWorker.RegisterActivityWithOptions(activities.MinusOne, activity.RegisterOptions{Name: "MinusOne"})
	computerWorker.RegisterActivityWithOptions(activities.DivideByTwo, activity.RegisterOptions{Name: "DivideByTwo"})
	log.Printf("starting computer activity worker on task queue %s", ComputerTaskQueue)

	if err := computerWorker.Run(worker.InterruptCh()); err != nil {
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
