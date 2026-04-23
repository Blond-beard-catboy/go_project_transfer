package models

import "time"

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Role           string    `json:"role"`
	FullName       string    `json:"full_name"`
	CreatedAt      time.Time `json:"created_at"`
}
