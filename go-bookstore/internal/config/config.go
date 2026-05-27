// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds runtime configuration for the bookstore service.
type Config struct {
	Addr           string
	ShutdownGrace  time.Duration
	DatabaseURL    string
	MaxPageSize    int
	DefaultPageLen int
	PaymentAPIURL  string
	AuditFilePath  string
}

// Load reads configuration from environment variables, falling back to defaults.
func Load() (Config, error) {
	cfg := Config{
		Addr:           getenv("BOOKSTORE_ADDR", ":8080"),
		DatabaseURL:    getenv("BOOKSTORE_DATABASE_URL", "postgres://localhost/bookstore?sslmode=disable"),
		PaymentAPIURL:  getenv("BOOKSTORE_PAYMENT_API_URL", "https://payments.example.com"),
		AuditFilePath:  getenv("BOOKSTORE_AUDIT_FILE", "audit.log"),
		ShutdownGrace:  5 * time.Second,
		MaxPageSize:    100,
		DefaultPageLen: 25,
	}

	if raw := os.Getenv("BOOKSTORE_SHUTDOWN_GRACE_SECONDS"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse BOOKSTORE_SHUTDOWN_GRACE_SECONDS: %w", err)
		}
		cfg.ShutdownGrace = time.Duration(n) * time.Second
	}

	if raw := os.Getenv("BOOKSTORE_MAX_PAGE_SIZE"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse BOOKSTORE_MAX_PAGE_SIZE: %w", err)
		}
		if n <= 0 {
			return Config{}, fmt.Errorf("BOOKSTORE_MAX_PAGE_SIZE must be positive, got %d", n)
		}
		cfg.MaxPageSize = n
	}

	if raw := os.Getenv("BOOKSTORE_DEFAULT_PAGE_LEN"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse BOOKSTORE_DEFAULT_PAGE_LEN: %w", err)
		}
		if n <= 0 || n > cfg.MaxPageSize {
			return Config{}, fmt.Errorf("BOOKSTORE_DEFAULT_PAGE_LEN must be in (0, %d], got %d", cfg.MaxPageSize, n)
		}
		cfg.DefaultPageLen = n
	}

	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
