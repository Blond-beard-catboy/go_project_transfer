package db

import (
	"database/sql"
	"log"

	"go_project_transfer/notification-service/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDBConnString())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Notification Service database connected")
	return db, nil
}
