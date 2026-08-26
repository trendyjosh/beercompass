package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

// Run database migrations
func RunMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("RunMigrations: creating migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("RunMigrations: running migrations: %w", err)
	}

	return nil
}
