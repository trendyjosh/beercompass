package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trendyjosh/beercompass/internal/importer"
)

// Establish a connection pool to the database
func ConnectDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("ConnectDB: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ConnectDB: pinging database: %w", err)
	}

	return db, nil
}

// Store extracted pubs in the database
func StorePubs(ctx context.Context, pool *pgxpool.Pool, pubs []importer.Pub) error {
	// ST_MakePoint stores lng before lat (mathmatics vector notation)
	const query = `
        INSERT INTO pubs (osm_id, name, location)
        VALUES (
			$1,
			$2,
			ST_SetSRID(
				ST_MakePoint($3, $4),
				4326
			)
		)
        ON CONFLICT (osm_id) DO UPDATE
            SET name     = EXCLUDED.name,
                location = EXCLUDED.location`

	for _, pub := range pubs {
		_, err := pool.Exec(ctx, query, pub.OsmID, pub.Name, pub.Lon, pub.Lat)
		if err != nil {
			return fmt.Errorf("storePubs: upserting pub %d: %w", pub.OsmID, err)
		}
	}

	return nil
}
