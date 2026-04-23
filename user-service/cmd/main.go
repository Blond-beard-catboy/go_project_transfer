package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"go_project_transfer/user-service/internal/config"
	"go_project_transfer/user-service/internal/db"
	"go_project_transfer/user-service/internal/handlers"
	"go_project_transfer/user-service/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)
	authHandler := handlers.NewAuthHandler(userRepo, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	log.Printf("User Service starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
