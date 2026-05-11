package config

import (
	"os"
)

type Config struct {
	Port                   string
	UserServiceURL         string
	CargoServiceURL        string
	RouteServiceURL        string
	OrderServiceURL        string
	NotificationServiceURL string
	PaymentServiceURL      string
	JWTSecret              string
}

func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "8005"),
		UserServiceURL:         getEnv("USER_SERVICE_URL", "http://localhost:8000"),
		CargoServiceURL:        getEnv("CARGO_SERVICE_URL", "http://localhost:8001"),
		RouteServiceURL:        getEnv("ROUTE_SERVICE_URL", "http://localhost:8002"),
		OrderServiceURL:        getEnv("ORDER_SERVICE_URL", "http://localhost:8003"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8006"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "http://localhost:8007"),
		JWTSecret:              getEnv("JWT_SECRET", "super-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
