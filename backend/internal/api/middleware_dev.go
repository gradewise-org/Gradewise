package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DevAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
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
