package repository

import (
	"database/sql"
	"time"

	"go_project_transfer/notification-service/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notif *models.Notification) error {
	query := `INSERT INTO notifications (user_id, type, subject, body, status, created_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(query, notif.UserID, notif.Type, notif.Subject, notif.Body, notif.Status, time.Now()).Scan(&notif.ID)
	return err
}

func (r *NotificationRepository) List(userID int, role string) ([]models.Notification, error) {
	var rows *sql.Rows
	var err error
	if role == "admin" {
		rows, err = r.db.Query(`SELECT id, user_id, type, subject, body, status, created_at FROM notifications ORDER BY created_at DESC`)
	} else {
		rows, err = r.db.Query(`SELECT id, user_id, type, subject, body, status, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Subject, &n.Body, &n.Status, &n.CreatedAt)
		if err != nil {
			continue
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}
