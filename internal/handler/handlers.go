package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/kerhael/accounting/internal/handler/v1"
	"github.com/kerhael/accounting/internal/infrastructure/repository"
	"github.com/kerhael/accounting/internal/service"
)

type HandlersV1 struct {
	Health   *v1.HealthHandler
	Category *v1.CategoryHandler
}

type Handlers struct {
	V1 *HandlersV1
}

func NewHandlers(db *pgxpool.Pool) *Handlers {
	healthRepo := repository.NewHealthRepository(db)
	healthService := service.NewHealthService(healthRepo)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)

	return &Handlers{
		V1: &HandlersV1{
			Health:   v1.NewHealthHandler(healthService),
			Category: v1.NewCategoryHandler(categoryService),
		},
	}
}
