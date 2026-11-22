package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var result int

		if err := db.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
			log.Printf("health query: %v", err)
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"status":   "ok",
			"database": "connected",
		}); err != nil {
			log.Printf("encode health: %v", err)
		}
	}
}
