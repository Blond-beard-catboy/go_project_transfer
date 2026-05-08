package models

import "time"

type Payment struct {
	ID        int        `json:"id"`
	OrderID   int        `json:"order_id"`
	Amount    float64    `json:"amount"`
	Status    string     `json:"status"` // pending, paid, failed, cancelled
	DueDate   *time.Time `json:"due_date,omitempty"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
