// Package database provides PostgreSQL connection management and health checks
package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB wraps sql.DB with additional methods
type DB struct {
	*sql.DB
}

// Config holds database connection configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfig creates a database configuration from environment variables
// Defaults to k8s service names for local development
func NewConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "gradewise_user"),
		Password: getEnv("DB_PASSWORD", "gradewise_pass"),
		DBName:   getEnv("DB_NAME", "gradewise_api"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewConnection establishes a PostgreSQL connection with connection pooling
func NewConnection(config *Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for production readiness
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection is working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) HealthCheck() error {
	var result int
	if err := db.QueryRow("SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	return nil
}

func (db *DB) GetHealthCheckStatus() (string, time.Time, error) {
	var status string
	var checkedAt time.Time
	
	query := "SELECT status, checked_at FROM health_check ORDER BY checked_at DESC LIMIT 1"
	err := db.QueryRow(query).Scan(&status, &checkedAt)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get health check status: %w", err)
	}
	
	return status, checkedAt, nil
}

func (db *DB) UpdateHealthCheck(status string) error {
	query := "INSERT INTO health_check (status) VALUES ($1)"
	_, err := db.Exec(query, status)
	if err != nil {
		return fmt.Errorf("failed to update health check: %w", err)
	}
	return nil
} 