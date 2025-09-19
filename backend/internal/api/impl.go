package api

import (
	"gradewise/backend/internal/temporal"
	"gradewise/backend/internal/temporal/workflow"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

// Ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Server)(nil)

type Server struct {
	Temporal client.Client
}

func NewServer(t client.Client) *Server {
	return &Server{Temporal: t}
}

func (s *Server) RegisterFaculty(c *gin.Context) {
	var req FacultyRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Error{Code: "BAD_REQUEST", Message: "invalid JSON body"})
		return
	}

	// TODO: hash password, insert into Postgres, create verification code, send email... basically do everyting useful lol
	c.JSON(http.StatusCreated, FacultyRegistrationResponse{
		FacultyId: uuid.MustParse(uuid.NewString()), // matches openapi_types.UUID
		Email:     req.Email,
		Status:    PENDINGVERIFICATION,
	})
}

func (s *Server) GreetUser(c *gin.Context, params GreetUserParams) {
	if s.Temporal == nil {
		c.String(http.StatusServiceUnavailable, "Temporal unavailable")
		return
	}

	ctx := c.Request.Context()

	we, err := s.Temporal.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{TaskQueue: temporal.TaskQueueName},
		workflow.InteractWithMe,
		params.Name,
	)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var result string
	if err := s.Temporal.GetWorkflow(ctx, we.GetID(), we.GetRunID()).Get(ctx, &result); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}
