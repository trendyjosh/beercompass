package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/trendyjosh/beercompass/internal/config"
	"github.com/trendyjosh/beercompass/internal/db"
	"github.com/trendyjosh/beercompass/internal/importer"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %w", err)
	}

	// Connect to DB
	pool, err := db.ConnectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Download OSM Data
	fmt.Println("Starting download...")
	if err := importer.DownloadOsm(); err != nil {
		log.Fatalf("failed to download OSM data: %v", err)
	}
	fmt.Println("Download complete.")

	// Parse OSM Data
	fmt.Println("Parsing OSM data...")
	pubs, err := importer.ParseOsm()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Data parsing complete.")

	// Store in DB
	fmt.Println("Storing pubs in DB...")
	if err := db.StorePubs(ctx, pool, pubs); err != nil {
		log.Fatalf("failed to store pubs: %v", err)
	}
	fmt.Println("DB store complete.")
}
