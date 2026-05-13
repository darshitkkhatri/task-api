package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// server
	Port         string
	ReadTimeout  int
	WriteTimeout int

	// database
	DBPath string

	// auth
	JWTSecret      string
	TokenExpiryHrs int

	// environment
	Env string // "development" | "production"
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		ReadTimeout:    getEnvInt("READ_TIMEOUT", 5),
		WriteTimeout:   getEnvInt("WRITE_TIMEOUT", 10),
		DBPath:         getEnv("DB_PATH", "tasks.db"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		TokenExpiryHrs: getEnvInt("TOKEN_EXPIRY_HRS", 24),
		Env:            getEnv("ENV", "development"),
	}

	// validate required fields
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	return nil
}

func (c *Config) IsProd() bool {
	return c.Env == "production"
}

// helpers
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}
