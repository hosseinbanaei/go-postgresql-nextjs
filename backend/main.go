package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
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

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API IS RUNNING")
	})

	mux.HandleFunc("POST /api/preview", previewHandler)

	server := &http.Server{
		Addr:              "localhost:8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("API running at http://localhost:8080")

	return server.ListenAndServe()
}
