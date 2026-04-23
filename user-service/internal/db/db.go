package db

import (
	"database/sql"
	"log"

	"go_project_transfer/user-service/internal/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.GetDBConnString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("User Service database connected")
	return db, nil
}
