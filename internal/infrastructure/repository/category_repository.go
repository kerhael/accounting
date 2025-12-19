package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerhael/accounting/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
}

type PostgresCategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{db: db}
}

func (r *PostgresCategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	query := `
		INSERT INTO categories (label)
		VALUES ($1)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, c.Label).Scan(&c.ID)
}
