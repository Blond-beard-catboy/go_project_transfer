package config

import (
	"os"
)

type Config struct {
	Port            string
	UserServiceURL  string
	CargoServiceURL string
	RouteServiceURL string
	OrderServiceURL string
	JWTSecret       string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("GATEWAY_PORT", "8002"),
		UserServiceURL:  getEnv("USER_SERVICE_URL", "http://localhost:8000"),
		CargoServiceURL: getEnv("CARGO_SERVICE_URL", "http://localhost:8001"),
		RouteServiceURL: getEnv("ROUTE_SERVICE_URL", "http://localhost:8003"),
		OrderServiceURL: getEnv("ORDER_SERVICE_URL", "http://localhost:8004"),
		JWTSecret:       getEnv("JWT_SECRET", "super-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
