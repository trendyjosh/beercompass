package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (Config, error) {
	pw := os.Getenv("POSTGRES_PASSWORD")
	if pw == "" {
		return Config{}, fmt.Errorf("config: POSTGRES_PASSWORD is required")
	}

	host := getenvDefault("POSTGRES_HOST", "localhost")
	port := getenvDefault("POSTGRES_PORT", "5432")
	user := getenvDefault("POSTGRES_USER", "osm")
	name := getenvDefault("POSTGRES_DB", "osm")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, url.QueryEscape(pw), host, port, name)

	return Config{
		DatabaseURL: dbURL,
		Port:        getenvDefault("PORT", "8080"),
	}, nil
}

func getenvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
