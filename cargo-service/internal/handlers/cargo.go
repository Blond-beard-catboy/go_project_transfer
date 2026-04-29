package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_project_transfer/cargo-service/internal/models"

	"github.com/go-chi/chi/v5"
)

type CargoHandler struct {
	db *sql.DB
}

func NewCargoHandler(db *sql.DB) *CargoHandler {
	return &CargoHandler{db: db}
}

type createCargoRequest struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Weight           float64   `json:"weight"`
	PickupLocation   string    `json:"pickup_location"`
	DeliveryLocation string    `json:"delivery_location"`
	PickupDate       time.Time `json:"pickup_date"`
	DeliveryDate     time.Time `json:"delivery_date"`
	OwnerID          int       `json:"owner_id"`
}

func (h *CargoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCargoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO cargos (title, description, weight, pickup_location, delivery_location, pickup_date, delivery_date, owner_id, status, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	var id int
	err := h.db.QueryRow(query, req.Title, req.Description, req.Weight, req.PickupLocation, req.DeliveryLocation,
		req.PickupDate, req.DeliveryDate, req.OwnerID, "pending", time.Now()).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create cargo", http.StatusInternalServerError)
		return
	}

	cargo := models.Cargo{
		ID:               id,
		Title:            req.Title,
		Description:      req.Description,
		Weight:           req.Weight,
		PickupLocation:   req.PickupLocation,
		DeliveryLocation: req.DeliveryLocation,
		PickupDate:       req.PickupDate,
		DeliveryDate:     req.DeliveryDate,
		OwnerID:          req.OwnerID,
		Status:           "pending",
		CreatedAt:        time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cargo)
}

func (h *CargoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid cargo id", http.StatusBadRequest)
		return
	}

	var cargo models.Cargo
	query := `SELECT id, title, description, weight, pickup_location, delivery_location, pickup_date, delivery_date, owner_id, status, created_at
              FROM cargos WHERE id = $1`
	row := h.db.QueryRow(query, id)
	err = row.Scan(&cargo.ID, &cargo.Title, &cargo.Description, &cargo.Weight,
		&cargo.PickupLocation, &cargo.DeliveryLocation, &cargo.PickupDate,
		&cargo.DeliveryDate, &cargo.OwnerID, &cargo.Status, &cargo.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Cargo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cargo)
}

func (h *CargoHandler) List(w http.ResponseWriter, r *http.Request) {
	ownerIDStr := r.URL.Query().Get("owner_id")
	if ownerIDStr == "" {
		http.Error(w, "owner_id required", http.StatusBadRequest)
		return
	}
	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		http.Error(w, "Invalid owner_id", http.StatusBadRequest)
		return
	}

	rows, err := h.db.Query(`SELECT id, title, description, weight, pickup_location, delivery_location, pickup_date, delivery_date, owner_id, status, created_at FROM cargos WHERE owner_id = $1`, ownerID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cargos []models.Cargo
	for rows.Next() {
		var c models.Cargo
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Weight, &c.PickupLocation, &c.DeliveryLocation,
			&c.PickupDate, &c.DeliveryDate, &c.OwnerID, &c.Status, &c.CreatedAt)
		if err != nil {
			continue
		}
		cargos = append(cargos, c)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cargos)
}
