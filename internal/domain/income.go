package domain

import "time"

type Income struct {
	Name      string     `json:"name"`      // Name of the income
	CreatedAt *time.Time `json:"createdAt"` // Date of the income (ex: "2026-01-01T00:00:00Z")
	Amount    int        `json:"amount"`    // Amount in cents (ex: 1999 for 19.99€)
	ID        int        `json:"id"`        // ID of the income
}
