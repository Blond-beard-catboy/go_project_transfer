package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"go_project_transfer/order-service/internal/clients"
	"go_project_transfer/order-service/internal/config"
	"go_project_transfer/order-service/internal/db"
	"go_project_transfer/order-service/internal/handlers"
	"go_project_transfer/order-service/internal/repository"
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

	orderRepo := repository.NewOrderRepository(database)
	cargoClient := clients.NewCargoClient(cfg.CargoServiceURL)
	routeClient := clients.NewRouteClient(cfg.RouteServiceURL)
	orderHandler := handlers.NewOrderHandler(orderRepo, cargoClient, routeClient, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/orders", orderHandler.CreateOrder)
	r.Get("/api/orders/{id}", orderHandler.GetOrder)
	r.Post("/api/orders/{id}/confirm", orderHandler.ConfirmOrder)

	log.Printf("Order Service starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
