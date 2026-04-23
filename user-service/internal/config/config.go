package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret          string
	JWTExpirationHours int64
}

func Load() *Config {
	return &Config{
		Port: getEnv("USER_SERVICE_PORT", "8000"),

		DBHost:     getEnv("USER_SERVICE_DB_HOST", "localhost"),
		DBPort:     getEnv("USER_SERVICE_DB_PORT", "5432"),
		DBUser:     getEnv("USER_SERVICE_DB_USER", "postgres"),
		DBPassword: getEnv("USER_SERVICE_DB_PASSWORD", "postgres"),
		DBName:     getEnv("USER_SERVICE_DB_NAME", "user_db"),

		JWTSecret:          getEnv("JWT_SECRET", "super-secret-key"),
		JWTExpirationHours: getEnvAsInt64("JWT_EXPIRATION_HOURS", 24),
	}
}

func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}
