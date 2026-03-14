package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Open creates a configured postgres connection pool.
func Open(ctx context.Context, dsn string, maxOpen, maxIdle int) (*sql.DB, error) {
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	database.SetMaxOpenConns(maxOpen)
	database.SetMaxIdleConns(maxIdle)
	database.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := database.PingContext(pingCtx); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return database, nil
}
