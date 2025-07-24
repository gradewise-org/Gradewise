package main

import (
	"gradewise/backend/internal/api"
	"gradewise/backend/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection with graceful fallback
	dbConfig := database.NewConfig()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Println("Continuing without database connection...")
	}

	server := api.NewServer()
	r := gin.Default()

	// Health and monitoring endpoints
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.GET("/health", func(c *gin.Context) {
		response := gin.H{
			"status": "ok",
			"time":   gin.H{"timestamp": time.Now().Format(time.RFC3339)},
		}
		
		// Add database health check
		if db != nil {
			if err := db.HealthCheck(); err != nil {
				response["database"] = gin.H{
					"status": "error",
					"error":  err.Error(),
				}
			} else {
				response["database"] = gin.H{
					"status": "connected",
				}
			}
		} else {
			response["database"] = gin.H{
				"status": "not_connected",
			}
		}
		
		c.JSON(http.StatusOK, response)
	})

	// Database connectivity test endpoint
	r.GET("/db-test", func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Database not connected",
			})
			return
		}

		// Test database operations
		status, checkedAt, err := db.GetHealthCheckStatus()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get health check status",
				"details": err.Error(),
			})
			return
		}

		// Update health check
		if err := db.UpdateHealthCheck("api_test"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update health check",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Database test successful",
			"last_status": status,
			"last_checked": checkedAt.Format(time.RFC3339),
			"current_time": time.Now().Format(time.RFC3339),
		})
	})

	// Register the API handlers
	api.RegisterHandlers(r, server)

	s :=&http.Server{
		Handler: r,
		Addr: "0.0.0.0:8080",
	}

	log.Fatalln(s.ListenAndServe())
}