package models

import "time"

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Subject   *string   `json:"subject,omitempty"`
	Body      string    `json:"body"`
	Status    string    `json:"status"` // sent, pending, failed
	CreatedAt time.Time `json:"created_at"`
}
