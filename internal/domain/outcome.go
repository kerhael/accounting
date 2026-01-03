package domain

import "time"

type Outcome struct {
	Name       string     `json:"name"`
	CreatedAt  *time.Time `json:"createdAt"`
	Amount     int        `json:"amount"`
	CategoryId int        `json:"categoryId"`
	ID         int        `json:"id"`
}
