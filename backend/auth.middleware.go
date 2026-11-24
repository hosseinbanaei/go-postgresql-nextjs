package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type authenticatedHandler func(http.ResponseWriter, *http.Request, User)

func requireAuth(db *sql.DB, next authenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || len(cookie.Value) != 64 {
			authError(w, http.StatusUnauthorized, "please log in")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var user User

		err = db.QueryRowContext(ctx, `
			SELECT users.id, users.name, users.email
			FROM sessions
			JOIN users ON users.id = sessions.user_id
			WHERE sessions.token_hash = $1
			AND sessions.expires_at > NOW()
		`, hashSessionToken(cookie.Value)).Scan(&user.ID, &user.Name, &user.Email)

		if errors.Is(err, sql.ErrNoRows) {
			clearSessionCookie(w)
			authError(w, http.StatusInternalServerError, "could not check session")
			return
		}
		if err != nil {
			log.Printf("check session: %v", err)
			authError(w, http.StatusUnauthorized, "could not check session")
			return
		}

		next(w, r, user)
	}
}

func meHandler(w http.ResponseWriter, _ *http.Request, user User) {
	authJSON(w, http.StatusOK, user)
}

func trustedBrowserRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		frontend := os.Getenv("FRONTEND_ORIGIN")
		apiOrigin := "http://localhost:8080"

		if origin != "" {
			if origin != apiOrigin && (frontend == "" || origin != frontend) {
				authError(w, http.StatusForbidden, "origin is not allowed")
				return
			}
		} else if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
			authError(w, http.StatusForbidden, "cross-site request is not allowed")
			return
		}

		next(w, r)
	}
}
