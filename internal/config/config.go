package config

import (
	"os"
)

type Config struct {
	Host              string
	Port              string
	User              string
	Password          string
	DBName            string
	SSLMode           string
	JWTSecret         string
	SeedAdminName     string
	SeedAdminEmail    string
	SeedAdminPassword string
}

// LoadConfig reads environment variables and returns database config
func LoadConfig() Config {
	return Config{
		Host:              getEnv("DB_HOST", "localhost"),
		Port:              getEnv("DB_PORT", "5434"),
		User:              getEnv("DB_USER", "postgres"),
		Password:          getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "profconnect_local"),
		SSLMode:           getEnv("DB_SSLMODE", "disable"),
		JWTSecret:         getEnv("JWT_SECRET", "your_jwt_secret_key"),
		SeedAdminName:     getEnv("SEED_ADMIN_NAME", ""),
		SeedAdminEmail:    getEnv("SEED_ADMIN_EMAIL", ""),
		SeedAdminPassword: getEnv("SEED_ADMIN_PASSWORD", ""),
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
