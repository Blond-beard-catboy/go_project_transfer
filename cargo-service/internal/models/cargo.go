package models

import "time"

type Cargo struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description,omitempty"`
	Weight           float64   `json:"weight"`
	PickupLocation   string    `json:"pickup_location"`
	DeliveryLocation string    `json:"delivery_location"`
	PickupDate       time.Time `json:"pickup_date"`
	DeliveryDate     time.Time `json:"delivery_date"`
	OwnerID          int       `json:"owner_id"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}
