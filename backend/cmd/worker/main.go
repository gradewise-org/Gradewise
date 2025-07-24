package main

import (
	"gradewise/backend/internal/temporal"
	"gradewise/backend/internal/temporal/activity"
	"gradewise/backend/internal/temporal/workflow"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the Temporal client
	c, err := client.Dial(client.Options{
		HostPort: temporal.GetTemporalAddress(),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create the Temporal worker
	w := worker.New(c, temporal.TaskQueueName, worker.Options{})

	// Inject local greeting (through instantiation) into the GreetingActivities struct
	activities := &activity.GreetingActivities{
		LocalGreeting: "おはよ",
	}

	// Register Workflows, Activities, etc.
	w.RegisterActivity(activities)
	w.RegisterWorkflow(workflow.InteractWithMe)

	// Start the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
