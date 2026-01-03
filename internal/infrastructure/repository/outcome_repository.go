package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerhael/accounting/internal/domain"
)

type OutcomeRepository interface {
	Create(ctx context.Context, c *domain.Outcome) error
}

type PostgresOutcomeRepository struct {
	db *pgxpool.Pool
}

func NewOutcomeRepository(db *pgxpool.Pool) *PostgresOutcomeRepository {
	return &PostgresOutcomeRepository{db: db}
}

func (r *PostgresOutcomeRepository) Create(ctx context.Context, o *domain.Outcome) error {
	query := `
		INSERT INTO outcomes (name, amount, category_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, o.Name, o.Amount, o.CategoryId, &o.CreatedAt).Scan(&o.ID)
}
