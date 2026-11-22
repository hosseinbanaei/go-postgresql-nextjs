package main

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
)

func readAuthJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || contentType != "application/json" {
		authError(w, http.StatusUnsupportedMediaType, "send application/json")
		return false
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		authError(w, http.StatusBadRequest, "could not read requestt body")
		return false
	}

	err = json.Unmarshal(body, destination)
	if err != nil {
		authError(w, http.StatusBadRequest, "invalid json")
		return false
	}

	return true
}

func authJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		log.Printf("write auth response: %v", err)
	}
}

func authError(w http.ResponseWriter, status int, message string) {
	authJSON(w, status, map[string]string{"error": message})
}
