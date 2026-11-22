package main

import (
	"database/sql"
	"net/http"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input RegisterInput
		if !readAuthJSON(w, r, &input) {
			return
		}
	}
}
