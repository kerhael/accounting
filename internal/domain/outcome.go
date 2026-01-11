package domain

import "time"

type Outcome struct {
	Name       string     `json:"name"`       // Name of the expense
	CreatedAt  *time.Time `json:"createdAt"`  // Date of the expense (ex: "2026-01-01T00:00:00Z")
	Amount     int        `json:"amount"`     // Amount in cents (ex: 1999 for 19.99€)
	CategoryId int        `json:"categoryId"` // ID of the associated category
	ID         int        `json:"id"`         // ID of the expense
}

type CategorySum struct {
	CategoryId int `json:"categoryId"` // Category ID
	Total      int `json:"total"`      // Total amount in cents for this category
}

type MonthlySeries struct {
	Month      string      `json:"month"`      // Month in YYYY-MM format
	Categories map[int]int `json:"categories"` // Map of categoryId to total amount
}

type MonthlyTotalSeries struct {
	Month string `json:"month"` // Month in YYYY-MM format
	Total int    `json:"total"` // Total amount
}
