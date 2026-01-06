package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerhael/accounting/internal/domain"
)

type OutcomeRepository interface {
	Create(ctx context.Context, c *domain.Outcome) error
	FindAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.Outcome, error)
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

func (r *PostgresOutcomeRepository) FindAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.Outcome, error) {
	query := `SELECT id, name, amount, category_id, created_at FROM outcomes WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if from != nil {
		argCount++
		query += ` AND created_at >= $` + strconv.Itoa(argCount)
		args = append(args, *from)
	}

	if to != nil {
		argCount++
		query += ` AND created_at <= $` + strconv.Itoa(argCount)
		args = append(args, *to)
	} else {
		query += ` AND created_at <= NOW()`
	}

	if categoryId != 0 {
		argCount++
		query += ` AND category_id = $` + strconv.Itoa(argCount)
		args = append(args, categoryId)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outcomes []domain.Outcome
	for rows.Next() {
		var o domain.Outcome
		if err := rows.Scan(&o.ID, &o.Name, &o.Amount, &o.CategoryId, &o.CreatedAt); err != nil {
			return nil, err
		}
		outcomes = append(outcomes, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return outcomes, nil
}
