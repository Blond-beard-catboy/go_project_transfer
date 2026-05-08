package repository

import (
	"database/sql"
	"time"

	"go_project_transfer/payment-service/internal/models"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(p *models.Payment) error {
	query := `INSERT INTO payments (order_id, amount, status, due_date, created_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, p.OrderID, p.Amount, p.Status, p.DueDate, time.Now()).Scan(&p.ID)
	return err
}

func (r *PaymentRepository) GetByID(id int) (*models.Payment, error) {
	query := `SELECT id, order_id, amount, status, due_date, paid_at, created_at FROM payments WHERE id = $1`
	var p models.Payment
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.DueDate, &p.PaidAt, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepository) List() ([]models.Payment, error) {
	rows, err := r.db.Query(`SELECT id, order_id, amount, status, due_date, paid_at, created_at FROM payments ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.DueDate, &p.PaidAt, &p.CreatedAt); err != nil {
			continue
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *PaymentRepository) UpdateStatus(id int, status string, paidAt *time.Time) error {
	query := `UPDATE payments SET status = $1, paid_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, paidAt, id)
	return err
}
