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

	CargoServiceURL        string
	RouteServiceURL        string
	NotificationServiceURL string
	PaymentServiceURL      string
}

func Load() *Config {
	return &Config{
		Port: getEnv("ORDER_SERVICE_PORT", "8004"),

		DBHost:     getEnv("ORDER_SERVICE_DB_HOST", "localhost"),
		DBPort:     getEnv("ORDER_SERVICE_DB_PORT", "5432"),
		DBUser:     getEnv("ORDER_SERVICE_DB_USER", "postgres"),
		DBPassword: getEnv("ORDER_SERVICE_DB_PASSWORD", "postgres"),
		DBName:     getEnv("ORDER_SERVICE_DB_NAME", "order_db"),

		CargoServiceURL:        getEnv("CARGO_SERVICE_URL", "http://localhost:8001"),
		RouteServiceURL:        getEnv("ROUTE_SERVICE_URL", "http://localhost:8003"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8005"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "http://localhost:8006"),
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
