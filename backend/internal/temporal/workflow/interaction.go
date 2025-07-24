package workflow

import (
	"fmt"
	"gradewise/backend/internal/temporal/activity"
	"strings"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func InteractWithMe(ctx workflow.Context, name string) (string, error) {
	// Define the activity options, including the retry policy
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second, //amount of time that must elapse before the first retry occurs
			MaximumInterval:    time.Minute, //maximum interval between retries
			BackoffCoefficient: 2,           //how much the retry interval increases
			// MaximumAttempts: 5, // Uncomment this if you want to limit attempts
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Run our activities
	var chat_log []string

	var acts *activity.GreetingActivities

	var r string
	err := workflow.ExecuteActivity(ctx, acts.Greet, name).Get(ctx, &r)

	if err != nil {
		return "", fmt.Errorf("Failed to greet: %s", err)
	}
	chat_log = append(chat_log, r)

	err = workflow.ExecuteActivity(ctx, acts.SayGoodbye).Get(ctx, &r)
	if err != nil {
		return "", fmt.Errorf("Failed to say goodbye: %s", err)
	}
	chat_log = append(chat_log, r)

	return strings.Join(chat_log, "\n"), nil
}
