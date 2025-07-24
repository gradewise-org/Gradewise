package api

import (
	"context"
	"gradewise/backend/internal/temporal"
	"gradewise/backend/internal/temporal/workflow"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

// Ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (Server) GreetUser(c *gin.Context, params GreetUserParams) {
	// Create Temporal client
	t, err := client.Dial(client.Options{
		HostPort: temporal.GetTemporalAddress(),
	})
	if err != nil {
		log.Println("Unable to create client", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer t.Close()

	we, err := t.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			// ID: "",
			TaskQueue: temporal.TaskQueueName,
		},
		workflow.InteractWithMe,
		params.Name,
	)
	if err != nil {
		log.Println("Unable to execute workflow", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var workflowResult string
	// err = we.Get(context.Background(), &workflowResult)
	desc, err := t.DescribeWorkflowExecution(context.Background(), we.GetID(), we.GetRunID())
	if err != nil {
		log.Println("Unable to describe workflow execution")
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("Status = %d", desc.WorkflowExecutionInfo.Status)

	err = t.GetWorkflow(context.Background(), we.GetID(), we.GetRunID()).Get(context.Background(), &workflowResult)
	if err != nil {
		log.Println("Unable to get workflow result", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	desc, err = t.DescribeWorkflowExecution(context.Background(), we.GetID(), we.GetRunID())
	if err != nil {
		log.Println("Unable to describe workflow execution")
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("Status = %d", desc.WorkflowExecutionInfo.Status)

	c.JSON(http.StatusOK, workflowResult)
}
