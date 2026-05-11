package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"go_project_transfer/payment-service/internal/config"
	"go_project_transfer/payment-service/internal/db"
	"go_project_transfer/payment-service/internal/handlers"
	"go_project_transfer/payment-service/internal/repository"
	"go_project_transfer/pkg/migrate"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()
	dbURL := cfg.GetDBConnString()

	if err := migrate.Run("file://./migrations", dbURL); err != nil {
		log.Fatal("Migration failed:", err)

	}

	database, err := db.Connect(cfg)

	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer database.Close()

	repo := repository.NewPaymentRepository(database)
	handler := handlers.NewPaymentHandler(repo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/payments", handler.Create)
	r.Get("/api/payments", handler.List)
	r.Get("/api/payments/{id}", handler.GetByID)
	r.Patch("/api/payments/{id}/pay", handler.Pay)

	log.Printf("Payment Service starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
