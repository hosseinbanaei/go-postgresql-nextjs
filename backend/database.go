package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func openDB(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return db, nil
}
