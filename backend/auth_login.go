package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input LoginInput

		if !readAuthJSON(w, r, &input) {
			return
		}

		email, err := cleanEmail(input.Email)
		if err != nil || len(input.Password) < 8 || len(input.Password) > 72 {
			authError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var (
			user       User
			storedHash string
		)

		err = db.QueryRowContext(ctx, `
			SELECT id, name, email, password_hash
			FROM users WHERE email = $1
		`, email).Scan(&user.ID, &user.Name, &user.Email, &storedHash)

		if errors.Is(err, sql.ErrNoRows) {
			authError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		if err != nil {
			log.Printf("find login user: %v", err)
			authError(w, http.StatusUnauthorized, "could not log in")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(input.Password))
		if err != nil {
			authError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		token, expiresAt, err := createSession(ctx, db, user.ID)
		if err != nil {
			log.Printf("create session: %v", err)
			authError(w, http.StatusInternalServerError, "could not log in")
			return
		}

		setSessionCookie(w, token, expiresAt)
		authJSON(w, http.StatusOK, user)
	}
}
