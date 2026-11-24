package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
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

		input.Name = strings.TrimSpace(input.Name)
		if input.Name == "" || len(input.Name) > 100 {
			authError(w, http.StatusBadRequest, "name must be 1 to 100 bytes")
			return
		}

		email, err := cleanEmail(input.Email)
		if err != nil {
			authError(w, http.StatusBadRequest, err.Error())
		}
		if len(input.Password) < 8 || len(input.Password) > 72 {
			authError(w, http.StatusBadRequest, "password must be 8 to 72 bytes")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("hash password: %v", err)
			authError(w, http.StatusInternalServerError, "could not create account")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var user User

		err = db.QueryRowContext(ctx, `
			INSERT INTO users (name, email, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id, name, email
		`, input.Name, email, string(hash)).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			var databaseErr *pgconn.PgError
			if errors.As(err, &databaseErr) && databaseErr.Code == "23505" {
				authError(w, http.StatusConflict, "email is already in use")
			}
			log.Printf("insert user: %v", err)
			authError(w, http.StatusInternalServerError, "could not create account")
			return
		}

		authJSON(w, http.StatusCreated, user)
	}
}
