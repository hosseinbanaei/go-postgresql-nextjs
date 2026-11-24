package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return errors.New("DATABASE_URL is required")
	}

	db, err := openDB(url)
	if err != nil {
		return err
	}

	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler(db))
	mux.HandleFunc("POST /api/preview", previewHandler)
	mux.HandleFunc("POST /api/auth/register", registerHandler(db))
	mux.HandleFunc("POST /api/auth/login", loginHandler(db))
	mux.HandleFunc("POST /api/auth/me", requireAuth(db, meHandler))

	server := &http.Server{
		Addr:              "localhost:8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("API running at http://localhost:8080")

	return server.ListenAndServe()
}
