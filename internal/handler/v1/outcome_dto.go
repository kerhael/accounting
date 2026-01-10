package v1

import (
	"time"

	"github.com/kerhael/accounting/internal/domain"
)

type CreateOutcomeRequest struct {
	Name       string    `json:"name"`       // Name of the expense
	CreatedAt  time.Time `json:"createdAt"`  // Date of the expense (ex: "2026-01-01T00:00:00Z")
	Amount     int       `json:"amount"`     // Amount in cents (ex: 1999 for 19.99€)
	CategoryId int       `json:"categoryId"` // ID of the associated category
}

type GetAllOutcomeRequest struct {
	From       time.Time `json:"from"`       // Start date (optional)
	To         time.Time `json:"to"`         // End date (optional)
	CategoryId int       `json:"categoryId"` // ID of a category (optional)
}

type OutcomeResponse domain.Outcome

type PatchOutcomeRequest struct {
	Name       *string    `json:"name"`       // Name of the expense (optional)
	CreatedAt  *time.Time `json:"createdAt"`  // Date of the expense (optional, ex: "2026-01-01T00:00:00Z")
	Amount     *int       `json:"amount"`     // Amount in cents (optional, ex: 1999 for 19.99€)
	CategoryId *int       `json:"categoryId"` // ID of the associated category (optional)
}

type CategorySumResponse domain.CategorySum

type SumOutcomeResponse []CategorySumResponse

type TotalOutcomeResponse struct {
	Total int `json:"total"` // Total amount in cents
}
