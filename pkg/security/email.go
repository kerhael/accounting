package security

import (
	"errors"
	"net/mail"
	"strings"
)

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateEmail(email string) error {
	email = NormalizeEmail(email)

	if email == "" {
		return errors.New("email is required")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	return nil
}
