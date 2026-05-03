package repository

import (
	"database/sql"
	"time"

	"go_project_transfer/order-service/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	query := `INSERT INTO orders (cargo_id, customer_id, driver_id, route_id, status, contract_file, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(query, order.CargoID, order.CustomerID, order.DriverID,
		order.RouteID, order.Status, order.ContractFile, time.Now()).Scan(&order.ID)
	return err
}

func (r *OrderRepository) GetByID(id int) (*models.Order, error) {
	query := `SELECT id, cargo_id, customer_id, driver_id, route_id, status, contract_file, created_at
              FROM orders WHERE id = $1`
	var order models.Order
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.CargoID, &order.CustomerID,
		&order.DriverID, &order.RouteID, &order.Status, &order.ContractFile, &order.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateStatus(id int, status string, contractFile *string) error {
	query := `UPDATE orders SET status = $1, contract_file = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, contractFile, id)
	return err
}
