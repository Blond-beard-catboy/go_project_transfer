package repository

import (
	"database/sql"
	"errors"
	"time"

	"go_project_transfer/user-service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, hashed_password, role, full_name, created_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, user.Email, user.HashedPassword, user.Role, user.FullName, time.Now()).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, hashed_password, role, full_name, created_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.HashedPassword, &user.Role, &user.FullName, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
