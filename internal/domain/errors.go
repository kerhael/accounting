package domain

import "fmt"

type InvalidEntityError struct {
	UnderlyingCause error
}

func (e *InvalidEntityError) Error() string {
	return fmt.Sprintf("invalid entity: %v", e.UnderlyingCause)
}

func (e *InvalidEntityError) Unwrap() error {
	return e.UnderlyingCause
}
