package config

import (
	"os"

	"profconnect-api/internal/adapter/database"
)

// LoadConfig reads environment variables and returns database config
func LoadConfig() database.Config {
	return database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5434"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "profconnect_local"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv retrieves an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
