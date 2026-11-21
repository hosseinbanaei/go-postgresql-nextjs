package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type PreviewInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func previewHandler(w http.ResponseWriter, r *http.Request) {
	var input PreviewInput

	err := json.NewDecoder(r.Body).Decode(&input)

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(input)
}
