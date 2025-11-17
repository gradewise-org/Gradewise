package api

import (
	"database/sql"
	"errors"
	"gradewise/backend/internal/database"
	"gradewise/backend/internal/temporal"
	"gradewise/backend/internal/temporal/workflow"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		VALUES ($1, $2, $3, $4, 'VERIFIED')
		RETURNING id
	`, email, string(hash), req.FullName, req.Institution).Scan(&facultyID); err != nil {

		var existingStatus string
		if err2 := s.DB.QueryRowContext(ctx, `
				SELECT id, status
				FROM faculty
				WHERE email = $1
			`, email).Scan(&facultyID, &existingStatus); err2 != nil {
			// If we somehow fail here, treat as server error.
			c.JSON(http.StatusInternalServerError, Error{
				Code:    "DB_LOOKUP_ERROR",
				Message: "email already exists but failed to lookup faculty",
			})
			return
		}

		// Return 200 OK with the existing faculty record
		c.JSON(http.StatusOK, FacultyRegistrationResponse{
			FacultyId: facultyID,
			Email:     openapi_types.Email(email),
			Status:    FacultyRegistrationResponseStatus(existingStatus),
		})
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
		Status:    VERIFIED, // keep pending until we add verify endpoint
	})
}

func (s *Server) CreateCourse(c *gin.Context) {
	if s.DB == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	v, ok := c.Get("facultyID")
	if !ok {
		c.JSON(http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "missing faculty identity"})
		return
	}
	facultyID := v.(uuid.UUID)
	var status string
	if err := s.DB.QueryRowContext(c, `SELECT status FROM faculty WHERE id = $1`, facultyID).Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "faculty not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, Error{Code: "DB_ERROR", Message: "failed to load faculty"})
		return
	}
	if status != "VERIFIED" {
		c.JSON(http.StatusForbidden, Error{Code: "NOT_VERIFIED", Message: "account not verified"})
		return
	}

	var req CourseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Error{Code: "BAD_REQUEST", Message: "invalid JSON body"})
		return
	}

	req.Code = strings.TrimSpace(req.Code)
	req.Term = strings.TrimSpace(req.Term)
	req.Title = strings.TrimSpace(req.Title)
	if req.Code == "" || req.Term == "" || req.Title == "" {
		c.JSON(http.StatusBadRequest, Error{Code: "BAD_REQUEST", Message: "code, term, and title are required"})
		return
	}

	var courseID uuid.UUID
	err := s.DB.QueryRowContext(
		c,
		`INSERT INTO courses
           (code, term, title, section, start_date, end_date, creator_faculty_id)
         VALUES ($1,$2,$3,$4,$5,$6,$7)
         RETURNING id`,
		req.Code,
		req.Term,
		req.Title,
		nullableString(req.Section),
		nullableDate(req.StartDate),
		nullableDate(req.EndDate),
		facultyID,
	).Scan(&courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Code: "DB_INSERT_ERROR", Message: "failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, CourseCreateResponse{CourseId: &courseID})
}

func (s *Server) ListCourses(c *gin.Context, params ListCoursesParams) {
	if s.DB == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	ctx := c.Request.Context()

	// `mine` defaults to true
	mine := true
	if params.Mine != nil {
		mine = *params.Mine
	}

	term := strings.TrimSpace(stringOrEmpty(params.Term))
	code := strings.TrimSpace(stringOrEmpty(params.Code))

	// If mine=true, we need a facultyID in context
	var facultyID uuid.UUID
	var haveFaculty bool
	if mine {
		if v, ok := c.Get("facultyID"); ok {
			if id, ok2 := v.(uuid.UUID); ok2 {
				facultyID = id
				haveFaculty = true
			}
		}
		if !haveFaculty {
			c.JSON(http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "missing faculty identity"})
			return
		}
	}

	var (
		rows *sql.Rows
		err  error
	)

	if mine {
		// Courses I own OR where I am a grader
		rows, err = s.DB.QueryContext(ctx, `
		  SELECT c.id, c.code, c.term, c.title, c.section,
		         c.start_date, c.end_date,
		         c.creator_faculty_id, c.created_at
		    FROM courses c
		   WHERE (c.creator_faculty_id = $1
		          OR EXISTS (
		               SELECT 1 FROM course_graders g
		                WHERE g.course_id = c.id
		                  AND g.faculty_id = $1
		          ))
		     AND ($2 = '' OR c.term = $2)
		     AND ($3 = '' OR c.code ILIKE $3)
		   ORDER BY c.created_at DESC
		   LIMIT 200
		`, facultyID, term, likeParam(code))
	} else {
		// All courses (optionally filtered)
		rows, err = s.DB.QueryContext(ctx, `
		  SELECT c.id, c.code, c.term, c.title, c.section,
		         c.start_date, c.end_date,
		         c.creator_faculty_id, c.created_at
		    FROM courses c
		   WHERE ($1 = '' OR c.term = $1)
		     AND ($2 = '' OR c.code ILIKE $2)
		   ORDER BY c.created_at DESC
		   LIMIT 200
		`, term, likeParam(code))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Code: "DB_QUERY_ERROR", Message: "failed to list courses"})
		return
	}
	defer rows.Close()

	var out []Course

	for rows.Next() {
		var (
			id      uuid.UUID
			codeV   string
			termV   string
			title   string
			section sql.NullString
			start   sql.NullTime
			end     sql.NullTime
			creator uuid.UUID
			created time.Time
		)

		if err := rows.Scan(&id, &codeV, &termV, &title, &section, &start, &end, &creator, &created); err != nil {
			c.JSON(http.StatusInternalServerError, Error{Code: "DB_SCAN_ERROR", Message: "failed to read courses"})
			return
		}

		var sectionPtr *string
		if section.Valid {
			s := section.String
			sectionPtr = &s
		}

		var startPtr *openapi_types.Date
		if start.Valid {
			startPtr = &openapi_types.Date{Time: start.Time}
		}

		var endPtr *openapi_types.Date
		if end.Valid {
			endPtr = &openapi_types.Date{Time: end.Time}
		}

		out = append(out, Course{
			Id:               &id,
			Code:             &codeV,
			Term:             &termV,
			Title:            &title,
			Section:          sectionPtr,
			StartDate:        startPtr,
			EndDate:          endPtr,
			CreatorFacultyId: &creator,
			CreatedAt:        &created,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, Error{Code: "DB_SCAN_ERROR", Message: "error iterating courses"})
		return
	}

	c.JSON(http.StatusOK, out)
}

func stringOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func likeParam(s string) string {
	if s == "" {
		return ""
	}
	return "%" + s + "%"
}

func nullableString(s *string) any {
	if s == nil {
		return nil
	}
	v := strings.TrimSpace(*s)
	if v == "" {
		return nil
	}
	return v
}

func nullableDate(d *openapi_types.Date) any {
	if d == nil {
		return nil
	}
	return d.Time
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
