package domain

import "time"

type Outcome struct {
	Name       string
	CreatedAt  *time.Time
	Amount     int
	CategoryId int
	ID         int
}

type CategorySum struct {
	CategoryId int
	Total      int
}

type MonthlySeries struct {
	Month      string
	Categories map[int]int
}

type MonthlyTotalSeries struct {
	Month string
	Total int
}
