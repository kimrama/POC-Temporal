package main

import (
	"log"
	"os"

	"temporal-progress-mock/workers/activities"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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

	computerWorker.RegisterActivityWithOptions(activities.SetInitialValues, activity.RegisterOptions{Name: "SetInitialValues"})

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
