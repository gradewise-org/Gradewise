package api

import (
	"errors"
	"gradewise/backend/internal/database"
	"gradewise/backend/internal/temporal"
	"gradewise/backend/internal/temporal/workflow"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.temporal.io/sdk/client"
	"golang.org/x/crypto/bcrypt"
)

// Ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Server)(nil)

type Server struct {
	Temporal client.Client
	DB       *database.DB
}

func NewServer(t client.Client, db *database.DB) *Server {
	return &Server{Temporal: t, DB: db}
}

func (s *Server) RegisterFaculty(c *gin.Context) {
	if s.DB == nil {
		c.JSON(http.StatusServiceUnavailable, Error{Code: "DB_UNAVAILABLE", Message: "database not available"})
		return
	}

	var req FacultyRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Error{Code: "BAD_REQUEST", Message: "invalid JSON body"})
		return
	}

	// TODO: Make improvements to validation once we figure out what we want for validation (special chars, certain amount of chars, etc.)
	if len(req.Password) < 10 {
		c.JSON(http.StatusBadRequest, Error{Code: "BAD_REQUEST", Message: "password must be at least 8 characters"})
		return
	}
	email := strings.ToLower(string(req.Email)) // normalize

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Code: "HASH_ERROR", Message: "failed to hash password"})
		return
	}

	ctx := c.Request.Context()

	// Insert faculty only (no verification logic yet)
	var facultyID uuid.UUID
	if err := s.DB.QueryRowContext(ctx, `
		INSERT INTO faculty (email, pass_hash, full_name, institution, status)
		VALUES ($1, $2, $3, $4, 'PENDING_VERIFICATION')
		RETURNING id
	`, email, string(hash), req.FullName, req.Institution).Scan(&facultyID); err != nil {

		// Map unique violation (email) -> 409 Conflict
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, Error{Code: "EMAIL_EXISTS", Message: "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, Error{Code: "DB_INSERT_ERROR", Message: "failed to create faculty"})
		return
	}

	// TODO: verification flow (disabled for now)
	// - generate one-time code
	// - upsert into faculty_verifications(email, code, expires_at)
	// - dispatch email via Temporal/worker
	// Example (commented):
	// code := "<generate-6-digit>"
	// expires := time.Now().Add(15 * time.Minute)
	// _, _ = s.DB.ExecContext(ctx, `
	//   INSERT INTO faculty_verifications (email, code, expires_at)
	//   VALUES ($1,$2,$3)
	//   ON CONFLICT (email) DO UPDATE
	//   SET code=EXCLUDED.code, expires_at=EXCLUDED.expires_at, created_at=now()
	// `, email, code, expires)
	// _ = s.Temporal.ExecuteWorkflow(ctx, client.StartWorkflowOptions{TaskQueue: temporal.TaskQueueName}, workflow.SendVerificationEmail, email, code)

	c.JSON(http.StatusCreated, FacultyRegistrationResponse{
		FacultyId: facultyID,
		Email:     openapi_types.Email(email),
		Status:    PENDINGVERIFICATION, // keep pending until we add verify endpoint
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
