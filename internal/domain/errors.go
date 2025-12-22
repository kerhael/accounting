package domain

import "fmt"

type ErrorResponse struct {
	Message string `json:"message"`
}

type InvalidEntityError struct {
	UnderlyingCause error
}

func (e *InvalidEntityError) Error() string {
	return fmt.Sprintf("invalid entity: %v", e.UnderlyingCause)
}

func (e *InvalidEntityError) Unwrap() error {
	return e.UnderlyingCause
}

type EntityNotFoundError struct {
	UnderlyingCause error
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("entity not found: %v", e.UnderlyingCause)
}

func (e *EntityNotFoundError) Unwrap() error {
	return e.UnderlyingCause
}
