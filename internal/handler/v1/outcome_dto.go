package v1

import (
	"time"
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

type OutcomeResponse struct {
	Name       string     `json:"name"`       // Name of the expense
	CreatedAt  *time.Time `json:"createdAt"`  // Date of the expense (ex: "2026-01-01T00:00:00Z")
	Amount     int        `json:"amount"`     // Amount in cents (ex: 1999 for 19.99€)
	CategoryId int        `json:"categoryId"` // ID of the associated category
	ID         int        `json:"id"`         // ID of the expense
}

type PatchOutcomeRequest struct {
	Name       *string    `json:"name"`       // Name of the expense (optional)
	CreatedAt  *time.Time `json:"createdAt"`  // Date of the expense (optional, ex: "2026-01-01T00:00:00Z")
	Amount     *int       `json:"amount"`     // Amount in cents (optional, ex: 1999 for 19.99€)
	CategoryId *int       `json:"categoryId"` // ID of the associated category (optional)
}

type CategorySumResponse struct {
	CategoryId int `json:"categoryId"` // Category ID
	Total      int `json:"total"`      // Total amount in cents for this category
}

type SumOutcomeResponse []CategorySumResponse

type TotalOutcomeResponse struct {
	Total int `json:"total"` // Total amount in cents
}

type MonthlySeries struct {
	Month      string      `json:"month"`      // Month in YYYY-MM format
	Categories map[int]int `json:"categories"` // Map of categoryId to total amount
}

type SeriesOutcomeResponse []MonthlySeries

type MonthlyTotalSeries struct {
	Month string `json:"month"` // Month in YYYY-MM format
	Total int    `json:"total"` // Total amount
}

type TotalSeriesOutcomeResponse []MonthlyTotalSeries
