package v1

import (
	"time"
)

type CreateIncomeRequest struct {
	Name      string    `json:"name"`      // Name of the income
	CreatedAt time.Time `json:"createdAt"` // Date of the income (ex: "2026-01-01T00:00:00Z")
	Amount    int       `json:"amount"`    // Amount in cents (ex: 1999 for 19.99€)
}

type GetAllIncomeRequest struct {
	From time.Time `json:"from"` // Start date (optional)
	To   time.Time `json:"to"`   // End date (optional)
}

type IncomeResponse struct {
	Name      string     `json:"name"`      // Name of the income
	CreatedAt *time.Time `json:"createdAt"` // Date of the income (ex: "2026-01-01T00:00:00Z")
	Amount    int        `json:"amount"`    // Amount in cents (ex: 1999 for 19.99€)
	ID        int        `json:"id"`        // ID of the income
}

type PatchIncomeRequest struct {
	Name      *string    `json:"name"`      // Name of the income (optional)
	CreatedAt *time.Time `json:"createdAt"` // Date of the income (optional, ex: "2026-01-01T00:00:00Z")
	Amount    *int       `json:"amount"`    // Amount in cents (optional, ex: 1999 for 19.99€)
}
