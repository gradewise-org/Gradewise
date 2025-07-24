// Package temporal provides shared configuration for Temporal workflows
package temporal

import (
	"os"
)

// TaskQueueName is the task queue used by all gradewise workflows and activities
const TaskQueueName = "gradewise-backend"

// GetTemporalAddress returns the Temporal server address
// Uses environment variable or defaults to k8s service name
func GetTemporalAddress() string {
	if addr := os.Getenv("TEMPORAL_HOST_PORT"); addr != "" {
		return addr
	}
	// Default to k8s service name and port
	return "temporal-server:7233"
}
