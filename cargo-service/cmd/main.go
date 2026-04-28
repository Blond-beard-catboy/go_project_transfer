package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"go_project_transfer/cargo-service/internal/config"
	"go_project_transfer/cargo-service/internal/db"
	"go_project_transfer/cargo-service/internal/handlers"
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

	cargoHandler := handlers.NewCargoHandler(database)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/cargo", cargoHandler.Create)
	r.Get("/api/cargo", cargoHandler.List)

	log.Printf("Cargo Service starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
