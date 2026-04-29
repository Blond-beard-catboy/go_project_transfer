package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"go_project_transfer/api-gateway/internal/config"
	myauth "go_project_transfer/api-gateway/internal/middleware"
	"go_project_transfer/api-gateway/internal/proxy"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	userProxy := proxy.NewProxy(cfg.UserServiceURL)
	cargoProxy := proxy.NewProxy(cfg.CargoServiceURL)
	routeProxy := proxy.NewProxy(cfg.RouteServiceURL)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(myauth.JWTAuthMiddleware(cfg.JWTSecret))

	// Добавляем заголовки X-User-ID и X-User-Role для проксируемых запросов
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if userID, ok := r.Context().Value(myauth.UserIDKey).(int); ok {
				r.Header.Set("X-User-ID", strconv.Itoa(userID))
			}
			if role, ok := r.Context().Value(myauth.UserRoleKey).(string); ok {
				r.Header.Set("X-User-Role", role)
			}
			next.ServeHTTP(w, r)
		})
	})

	// Публичные маршруты (регистрация, логин)
	r.HandleFunc("/api/register", userProxy.ServeHTTP)
	r.HandleFunc("/api/login", userProxy.ServeHTTP)

	// Прокси для user-service (защищённые)
	r.Handle("/api/users/*", http.StripPrefix("/api/users", userProxy))

	// Прокси для cargo-service (защищённые)
	r.Handle("/api/cargo", cargoProxy)
	r.Handle("/api/cargo/*", cargoProxy)

	// Прокси для route-service (защищённые)
	r.Handle("/api/routes", routeProxy)
	r.Handle("/api/routes/*", routeProxy)

	log.Printf("API Gateway starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
