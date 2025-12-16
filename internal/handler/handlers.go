package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerhael/accounting/internal/repository"
)

type Handlers struct {
	Health *HealthHandler
}

func NewHandlers(db *pgxpool.Pool) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(
			repository.NewPostgresHealthRepository(db),
		),
	}
}
