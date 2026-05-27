package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr           string
	ShutdownGrace  time.Duration
	DatabaseURL    string
}

func Load() (Config, error) {
	cfg := Config{
		Addr:          getenv("PAYMENTS_ADDR", ":8080"),
		DatabaseURL:   getenv("PAYMENTS_DATABASE_URL", "postgres://localhost/payments?sslmode=disable"),
		ShutdownGrace: 5 * time.Second,
	}

	if raw := os.Getenv("PAYMENTS_SHUTDOWN_GRACE_SECONDS"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse PAYMENTS_SHUTDOWN_GRACE_SECONDS: %w", err)
		}
		cfg.ShutdownGrace = time.Duration(n) * time.Second
	}
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
