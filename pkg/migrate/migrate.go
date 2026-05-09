package migrate

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run выполняет миграции для заданной базы данных.
// migrationsPath должен указывать на папку с миграциями, например "file://./migrations"
// databaseURL - строка подключения к PostgreSQL (например, "postgres://user:pass@localhost:5432/db?sslmode=disable")
func Run(migrationsPath, databaseURL string) error {
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}
	log.Printf("Migrations applied successfully for %s", databaseURL)
	return nil
}
