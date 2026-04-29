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

	CargoServiceURL string
}

func Load() *Config {
	return &Config{
		Port: getEnv("ROUTE_SERVICE_PORT", "8003"),

		DBHost:     getEnv("ROUTE_SERVICE_DB_HOST", "localhost"),
		DBPort:     getEnv("ROUTE_SERVICE_DB_PORT", "5432"),
		DBUser:     getEnv("ROUTE_SERVICE_DB_USER", "postgres"),
		DBPassword: getEnv("ROUTE_SERVICE_DB_PASSWORD", "postgres"),
		DBName:     getEnv("ROUTE_SERVICE_DB_NAME", "route_db"),

		CargoServiceURL: getEnv("CARGO_SERVICE_URL", "http://localhost:8001"),
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
