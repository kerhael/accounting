package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerhael/accounting/internal/infrastructure/repository"
	"github.com/kerhael/accounting/internal/service"
)

type Handlers struct {
	Health   *HealthHandler
	Category *CategoryHandler
}

func NewHandlers(db *pgxpool.Pool) *Handlers {
	healthRepo := repository.NewHealthRepository(db)
	healthService := service.NewHealthService(healthRepo)
	healthHandler := NewHealthHandler(healthService)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := NewCategoryHandler(categoryService)

	return &Handlers{
		Health:   healthHandler,
		Category: categoryHandler,
	}
}
