package main

import (
	"log"

	"github.com/trendyjosh/beercompass/internal/config"
	"github.com/trendyjosh/beercompass/internal/db"
)

// Run database migrations and log outcome.
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("migrate: loading config: %v", err)
	}

	if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrate: running migrations: %v", err)
	}

	log.Println("migrate: migrations applied successfully")
}
