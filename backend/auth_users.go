package main

import (
	"errors"
	"net/mail"
	"strings"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func cleanEmail(value string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(value))
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", errors.New("enter a valid email address")
	}

	if address.Address != email {
		return "", errors.New("enter a valid email address")
	}

	return email, nil
}
