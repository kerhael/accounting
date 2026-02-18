package security

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/kerhael/accounting/internal/domain"
)

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateEmail(email string) error {
	email = NormalizeEmail(email)

	if email == "" {
		return &domain.InvalidEntityError{
			UnderlyingCause: errors.New("email is required"),
		}
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid email format"),
		}
	}

	return nil
}
