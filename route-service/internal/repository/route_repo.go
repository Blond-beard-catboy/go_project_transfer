package repository

import (
	"database/sql"
	"time"

	"go_project_transfer/route-service/internal/models"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

// CreateRoute создаёт пустой маршрут
func (r *RouteRepository) CreateRoute(status string) (*models.Route, error) {
	query := `INSERT INTO routes (status, created_at) VALUES ($1, $2) RETURNING id, order_id, status, created_at`
	var route models.Route
	err := r.db.QueryRow(query, status, time.Now()).Scan(&route.ID, &route.OrderID, &route.Status, &route.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &route, nil
}

// GetRouteByID получает маршрут вместе с точками
func (r *RouteRepository) GetRouteByID(id int) (*models.Route, []models.RoutePoint, error) {
	// сначала маршрут
	routeQuery := `SELECT id, order_id, status, created_at FROM routes WHERE id = $1`
	var route models.Route
	err := r.db.QueryRow(routeQuery, id).Scan(&route.ID, &route.OrderID, &route.Status, &route.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	// точки
	pointsQuery := `SELECT id, route_id, type, cargo_id, address, planned_time, actual_time, status 
                    FROM route_points WHERE route_id = $1 ORDER BY id`
	rows, err := r.db.Query(pointsQuery, id)
	if err != nil {
		return &route, nil, err
	}
	defer rows.Close()

	var points []models.RoutePoint
	for rows.Next() {
		var p models.RoutePoint
		err := rows.Scan(&p.ID, &p.RouteID, &p.Type, &p.CargoID, &p.Address, &p.PlannedTime, &p.ActualTime, &p.Status)
		if err != nil {
			continue
		}
		points = append(points, p)
	}
	return &route, points, nil
}

// AddPoint добавляет точку к маршруту
func (r *RouteRepository) AddPoint(routeID int, pointType, address string, cargoID *int, plannedTime *time.Time) (*models.RoutePoint, error) {
	query := `INSERT INTO route_points (route_id, type, cargo_id, address, planned_time, status)
              VALUES ($1, $2, $3, $4, $5, 'pending') RETURNING id, route_id, type, cargo_id, address, planned_time, actual_time, status`
	var p models.RoutePoint
	err := r.db.QueryRow(query, routeID, pointType, cargoID, address, plannedTime).Scan(
		&p.ID, &p.RouteID, &p.Type, &p.CargoID, &p.Address, &p.PlannedTime, &p.ActualTime, &p.Status)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdatePointStatus отмечает точку выполненной
func (r *RouteRepository) UpdatePointStatus(pointID int, status string, actualTime time.Time) error {
	query := `UPDATE route_points SET status = $1, actual_time = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, actualTime, pointID)
	return err
}

// UpdateRouteStatus обновляет статус маршрута (если нужно)
func (r *RouteRepository) UpdateRouteStatus(routeID int, status string) error {
	query := `UPDATE routes SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, routeID)
	return err
}
