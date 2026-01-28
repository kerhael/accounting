package domain

import "time"

type Income struct {
	Name      string
	CreatedAt *time.Time
	Amount    int
	ID        int
}
