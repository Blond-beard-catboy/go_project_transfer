package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go_project_transfer/notification-service/internal/models"
	"go_project_transfer/notification-service/internal/repository"
)

type NotificationHandler struct {
	repo *repository.NotificationRepository
}

func NewNotificationHandler(repo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

type createNotificationRequest struct {
	UserID  int    `json:"user_id"`
	Type    string `json:"type"`
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body"`
}

func (h *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	notif := &models.Notification{
		UserID: req.UserID,
		Type:   req.Type,
		Body:   req.Body,
		Status: "sent", // сразу помечаем как отправленное
	}
	if req.Subject != "" {
		notif.Subject = &req.Subject
	}

	if err := h.repo.Create(notif); err != nil {
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	// Эмуляция отправки – просто пишем в лог
	subj := ""
	if notif.Subject != nil {
		subj = *notif.Subject
	}
	log.Printf("📨 Sending notification to user %d: %s - %s", notif.UserID, subj, notif.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notif)
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	// Заголовки X-User-ID и X-User-Role должны быть добавлены Gateway (или переданы вручную при тесте)
	userIDHeader := r.Header.Get("X-User-ID")
	roleHeader := r.Header.Get("X-User-Role")
	// Для простоты примем значения из заголовков
	userID := 0
	role := "driver"
	if userIDHeader != "" {
		// парсинг ID (в реальности нужен strconv)
	}
	if roleHeader != "" {
		role = roleHeader
	}

	notifs, err := h.repo.List(userID, role)
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}
