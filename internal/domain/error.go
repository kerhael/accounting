package domain

import "fmt"

type InvalidDateError struct {
	UnderlyingCause error
}

func (e *InvalidDateError) Error() string {
	return fmt.Sprintf("invalid date range: %v", e.UnderlyingCause)
}

func (e *InvalidDateError) Unwrap() error {
	return e.UnderlyingCause
}

type InvalidEntityError struct {
	UnderlyingCause error
}

func (e *InvalidEntityError) Error() string {
	return fmt.Sprintf("invalid entity data: %v", e.UnderlyingCause)
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
