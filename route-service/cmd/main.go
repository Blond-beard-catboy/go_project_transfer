package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"go_project_transfer/route-service/internal/clients"
	"go_project_transfer/route-service/internal/config"
	"go_project_transfer/route-service/internal/db"
	"go_project_transfer/route-service/internal/handlers"
	"go_project_transfer/route-service/internal/repository"
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

	routeRepo := repository.NewRouteRepository(database)
	cargoClient := clients.NewCargoClient(cfg.CargoServiceURL)
	routeHandler := handlers.NewRouteHandler(routeRepo, cargoClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/routes", routeHandler.CreateRoute)
	r.Get("/api/routes/{id}", routeHandler.GetRoute)
	r.Post("/api/routes/{id}/cargo/{cargo_id}", routeHandler.AddCargoToRoute)
	// можно добавить другие эндпоинты по мере необходимости

	log.Printf("Route Service starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
