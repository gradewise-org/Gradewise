package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DevAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		if method == http.MethodPost && path == "/api/faculty/register" {
			c.Next()
			return
		}

		if path == "/health" || path == "/api/health" {
			c.Next()
			return
		}

		if strings.HasPrefix(path, "/lti") {
			c.Next()
			return
		}

		id := c.GetHeader("X-Faculty-ID")
		if id == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "X-Faculty-ID required (dev)"})
			return
		}
		u, err := uuid.Parse(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "invalid X-Faculty-ID"})
			return
		}
		c.Set("facultyID", u)
		c.Next()
	}
}
