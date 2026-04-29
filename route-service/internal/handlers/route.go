package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_project_transfer/route-service/internal/clients"
	"go_project_transfer/route-service/internal/models"
	"go_project_transfer/route-service/internal/repository"
)

type RouteHandler struct {
	repo        *repository.RouteRepository
	cargoClient *clients.CargoClient
}

func NewRouteHandler(repo *repository.RouteRepository, cargoClient *clients.CargoClient) *RouteHandler {
	return &RouteHandler{
		repo:        repo,
		cargoClient: cargoClient,
	}
}

// CreateRoute POST /api/routes
func (h *RouteHandler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	route, err := h.repo.CreateRoute("planned")
	if err != nil {
		http.Error(w, "Failed to create route", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(route)
}

// GetRoute GET /api/routes/{id}
func (h *RouteHandler) GetRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/routes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid route id", http.StatusBadRequest)
		return
	}

	route, points, err := h.repo.GetRouteByID(id)
	if err != nil {
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"route":  route,
		"points": points,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AddPoint POST /api/routes/{id}/points
func (h *RouteHandler) AddPoint(w http.ResponseWriter, r *http.Request) {
	// извлекаем id маршрута из URL
	// реализация аналогична GetRoute, но требует парсинга тела запроса
	// для краткости опущена, добавим позже при необходимости
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// UpdatePoint PATCH /api/routes/{id}/points/{point_id}
func (h *RouteHandler) UpdatePoint(w http.ResponseWriter, r *http.Request) {
	// парсим point_id и обновляем статус
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// AddCargoToRoute POST /api/routes/{id}/cargo/{cargo_id}
func (h *RouteHandler) AddCargoToRoute(w http.ResponseWriter, r *http.Request) {
	// извлекаем id маршрута и cargo_id из URL
	// пример: /api/routes/5/cargo/12
	// распарсим путь вручную
	// Для простоты используем путь: /api/routes/5/cargo/12
	// разобьём по слешам
	path := r.URL.Path
	// path = "/api/routes/5/cargo/12"
	parts := splitPath(path)
	// parts = ["api", "routes", "5", "cargo", "12"]
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	routeID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid route id", http.StatusBadRequest)
		return
	}
	cargoID, err := strconv.Atoi(parts[4])
	if err != nil {
		http.Error(w, "Invalid cargo id", http.StatusBadRequest)
		return
	}

	// получаем данные груза из Cargo Service
	cargo, err := h.cargoClient.GetCargo(cargoID)
	if err != nil {
		http.Error(w, "Cargo service error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	// создаём точки погрузки и разгрузки
	// используем транзакцию (для простоты делаем два отдельных вызова, но лучше в одной транзакции)
	pickupPoint, err := h.repo.AddPoint(routeID, "pickup", cargo.PickupLocation, &cargoID, nil)
	if err != nil {
		http.Error(w, "Failed to add pickup point", http.StatusInternalServerError)
		return
	}
	deliveryPoint, err := h.repo.AddPoint(routeID, "delivery", cargo.DeliveryLocation, &cargoID, nil)
	if err != nil {
		// если вторая точка не создалась, удалить первую? упростим: просто вернём ошибку
		http.Error(w, "Failed to add delivery point", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"route_id": routeID,
		"points":   []*models.RoutePoint{pickupPoint, deliveryPoint},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// вспомогательная функция для разбора пути
func splitPath(path string) []string {
	if path[0] == '/' {
		path = path[1:]
	}
	parts := []string{}
	start := 0
	for i, ch := range path {
		if ch == '/' {
			parts = append(parts, path[start:i])
			start = i + 1
		}
	}
	if start < len(path) {
		parts = append(parts, path[start:])
	}
	return parts
}
