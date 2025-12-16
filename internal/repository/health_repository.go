package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthRepository interface {
	Check(ctx context.Context) error
}

type PostgresHealthRepository struct {
	db *pgxpool.Pool
}

func NewPostgresHealthRepository(db *pgxpool.Pool) *PostgresHealthRepository {
	return &PostgresHealthRepository{db: db}
}

func (r *PostgresHealthRepository) Check(ctx context.Context) error {
	return r.db.Ping(ctx)
}
