package models

import "time"

type Order struct {
	ID           int       `json:"id"`
	CargoID      int       `json:"cargo_id"`
	CustomerID   int       `json:"customer_id"`
	DriverID     *int      `json:"driver_id,omitempty"`
	RouteID      int       `json:"route_id"`
	Status       string    `json:"status"` // new, confirmed, in_progress, completed, cancelled
	ContractFile *string   `json:"contract_file,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// Ответы от внешних сервисов
type CargoData struct {
	ID               int     `json:"id"`
	PickupLocation   string  `json:"pickup_location"`
	DeliveryLocation string  `json:"delivery_location"`
	Title            string  `json:"title"`
	Weight           float64 `json:"weight"`
}
