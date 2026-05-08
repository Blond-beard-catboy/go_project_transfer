package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go_project_transfer/order-service/internal/clients"
	"go_project_transfer/order-service/internal/config"
	"go_project_transfer/order-service/internal/models"
	"go_project_transfer/order-service/internal/repository"
)

type OrderHandler struct {
	repo          *repository.OrderRepository
	cargoClient   *clients.CargoClient
	routeClient   *clients.RouteClient
	paymentClient *clients.PaymentClient
	cfg           *config.Config
}

func NewOrderHandler(
	repo *repository.OrderRepository,
	cargoClient *clients.CargoClient,
	routeClient *clients.RouteClient,
	cfg *config.Config,
	paymentClient *clients.PaymentClient,
) *OrderHandler {
	return &OrderHandler{
		repo:          repo,
		cargoClient:   cargoClient,
		routeClient:   routeClient,
		cfg:           cfg,
		paymentClient: paymentClient,
	}
}

type createOrderRequest struct {
	CargoID    int `json:"cargo_id"`
	CustomerID int `json:"customer_id"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 1. Проверить груз в Cargo Service
	_, err := h.cargoClient.GetCargo(req.CargoID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cargo service error: %v", err), http.StatusServiceUnavailable)
		return
	}

	// 2. Создать маршрут в Route Service
	routeID, err := h.routeClient.CreateRoute()
	if err != nil {
		http.Error(w, fmt.Sprintf("Route service error: %v", err), http.StatusServiceUnavailable)
		return
	}

	// 3. Добавить груз в маршрут
	if err := h.routeClient.AddCargoToRoute(routeID, req.CargoID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add cargo to route: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Сохранить заказ в БД
	order := &models.Order{
		CargoID:    req.CargoID,
		CustomerID: req.CustomerID,
		RouteID:    routeID,
		Status:     "new",
	}
	if err := h.repo.Create(order); err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/orders/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order id", http.StatusBadRequest)
		return
	}
	order, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID заказа из URL (последний сегмент)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-2] // предпоследний элемент – это id
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order id", http.StatusBadRequest)
		return
	}

	// Получаем заказ
	order, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	if order.Status != "new" {
		http.Error(w, "Order cannot be confirmed", http.StatusBadRequest)
		return
	}

	// Генерируем PDF (заглушка)
	contractFile := fmt.Sprintf("/contracts/order_%d.pdf", order.ID)
	if err := h.repo.UpdateStatus(order.ID, "confirmed", &contractFile); err != nil {
		http.Error(w, "Failed to confirm order", http.StatusInternalServerError)
		return
	}

	// Отправляем уведомление (если конфиг содержит URL)
	if h.cfg != nil && h.cfg.NotificationServiceURL != "" {
		notifClient := clients.NewNotificationClient(h.cfg.NotificationServiceURL)
		subject := "Order confirmed"
		body := fmt.Sprintf("Your order #%d has been confirmed", order.ID)
		if err := notifClient.SendNotification(order.CustomerID, subject, body); err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	}

	// Добавляем создание платежа
	if h.paymentClient != nil {
		amount := 1000.0 // или рассчитайте из веса груза, например
		if err := h.paymentClient.CreatePayment(order.ID, amount); err != nil {
			log.Printf("Failed to create payment: %v", err)
		}
	}

	// Ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":        "confirmed",
		"contract_file": contractFile,
	})
}
