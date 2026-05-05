package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() *Config {
	return &Config{
		Port: getEnv("NOTIFICATION_SERVICE_PORT", "8005"),

		DBHost:     getEnv("NOTIFICATION_SERVICE_DB_HOST", "localhost"),
		DBPort:     getEnv("NOTIFICATION_SERVICE_DB_PORT", "5432"),
		DBUser:     getEnv("NOTIFICATION_SERVICE_DB_USER", "postgres"),
		DBPassword: getEnv("NOTIFICATION_SERVICE_DB_PASSWORD", "postgres"),
		DBName:     getEnv("NOTIFICATION_SERVICE_DB_NAME", "notification_db"),
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
