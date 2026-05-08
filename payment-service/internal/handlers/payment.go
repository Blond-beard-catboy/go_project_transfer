package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_project_transfer/payment-service/internal/models"
	"go_project_transfer/payment-service/internal/repository"
)

type PaymentHandler struct {
	repo *repository.PaymentRepository
}

func NewPaymentHandler(repo *repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{repo: repo}
}

type createPaymentRequest struct {
	OrderID int     `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	payment := &models.Payment{
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  "pending",
	}
	if err := h.repo.Create(payment); err != nil {
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/payments/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid payment id", http.StatusBadRequest)
		return
	}
	payment, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	payments, err := h.repo.List()
	if err != nil {
		http.Error(w, "Failed to fetch payments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (h *PaymentHandler) Pay(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/payments/"):]
	idStr = idStr[:len(idStr)-len("/pay")]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid payment id", http.StatusBadRequest)
		return
	}

	payment, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}
	if payment.Status != "pending" {
		http.Error(w, "Payment already processed", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if err := h.repo.UpdateStatus(id, "paid", &now); err != nil {
		http.Error(w, "Failed to update payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "paid"})
}
