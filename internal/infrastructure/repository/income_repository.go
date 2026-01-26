package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerhael/accounting/internal/domain"
)

type IncomeRepository interface {
	Create(ctx context.Context, c *domain.Income) error
	FindAll(ctx context.Context, from *time.Time, to *time.Time) ([]domain.Income, error)
	FindById(ctx context.Context, id int) (*domain.Income, error)
	Update(ctx context.Context, o *domain.Income) error
	DeleteById(ctx context.Context, id int) error
}

type PostgresIncomeRepository struct {
	db *pgxpool.Pool
}

func NewIncomeRepository(db *pgxpool.Pool) *PostgresIncomeRepository {
	return &PostgresIncomeRepository{db: db}
}

func (r *PostgresIncomeRepository) Create(ctx context.Context, o *domain.Income) error {
	query := `
		INSERT INTO incomes (name, amount, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, o.Name, o.Amount, &o.CreatedAt).Scan(&o.ID)
}

func (r *PostgresIncomeRepository) FindAll(ctx context.Context, from *time.Time, to *time.Time) ([]domain.Income, error) {
	query := `SELECT id, name, amount, created_at FROM incomes WHERE 1=1`
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

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []domain.Income
	for rows.Next() {
		var o domain.Income
		if err := rows.Scan(&o.ID, &o.Name, &o.Amount, &o.CreatedAt); err != nil {
			return nil, err
		}
		incomes = append(incomes, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incomes, nil
}

func (r *PostgresIncomeRepository) FindById(ctx context.Context, id int) (*domain.Income, error) {
	var c domain.Income

	query := `
		SELECT id, name, amount, created_at FROM incomes
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Amount, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *PostgresIncomeRepository) Update(ctx context.Context, o *domain.Income) error {
	query := `UPDATE incomes SET name = $1, amount = $2, created_at = $3 WHERE id = $4`

	_, err := r.db.Exec(ctx, query, o.Name, o.Amount, o.CreatedAt, o.ID)
	return err
}

func (r *PostgresIncomeRepository) DeleteById(ctx context.Context, id int) error {
	query := `
		DELETE FROM incomes
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
