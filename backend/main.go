package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API IS RUNNING")
	})

	mux.HandleFunc("POST /api/preview", previewHandler)

	log.Println("API running at http://localhost:8080")

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Print(err)
	}
}
