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
		Port: getEnv("PAYMENT_SERVICE_PORT", "8007"),

		DBHost:     getEnv("PAYMENT_SERVICE_DB_HOST", "localhost"),
		DBPort:     getEnv("PAYMENT_SERVICE_DB_PORT", "5432"),
		DBUser:     getEnv("PAYMENT_SERVICE_DB_USER", "postgres"),
		DBPassword: getEnv("PAYMENT_SERVICE_DB_PASSWORD", "postgres"),
		DBName:     getEnv("PAYMENT_SERVICE_DB_NAME", "payment_db"),
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
