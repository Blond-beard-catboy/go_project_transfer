package models

import "time"

type Route struct {
	ID        int       `json:"id"`
	OrderID   *int      `json:"order_id,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type RoutePoint struct {
	ID          int        `json:"id"`
	RouteID     int        `json:"route_id"`
	Type        string     `json:"type"` // pickup, delivery, service
	CargoID     *int       `json:"cargo_id,omitempty"`
	Address     string     `json:"address"`
	PlannedTime *time.Time `json:"planned_time,omitempty"`
	ActualTime  *time.Time `json:"actual_time,omitempty"`
	Status      string     `json:"status"` // pending, done
}

// Структура для получения данных груза из Cargo Service
type CargoData struct {
	ID               int    `json:"id"`
	PickupLocation   string `json:"pickup_location"`
	DeliveryLocation string `json:"delivery_location"`
	// другие поля не нужны
}
