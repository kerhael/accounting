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
	FindById(ctx context.Context, id int) (*domain.Outcome, error)
	Update(ctx context.Context, o *domain.Outcome) error
	DeleteById(ctx context.Context, id int) error
	GetSumByCategory(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.CategorySum, error)
	GetTotalSum(ctx context.Context, from *time.Time, to *time.Time) (int, error)
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

func (r *PostgresOutcomeRepository) FindById(ctx context.Context, id int) (*domain.Outcome, error) {
	var c domain.Outcome

	query := `
		SELECT id, name, amount, category_id, created_at FROM outcomes
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Amount, &c.CategoryId, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *PostgresOutcomeRepository) Update(ctx context.Context, o *domain.Outcome) error {
	query := `UPDATE outcomes SET name = $1, amount = $2, category_id = $3, created_at = $4 WHERE id = $5`

	_, err := r.db.Exec(ctx, query, o.Name, o.Amount, o.CategoryId, o.CreatedAt, o.ID)
	return err
}

func (r *PostgresOutcomeRepository) DeleteById(ctx context.Context, id int) error {
	query := `
		DELETE FROM outcomes
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PostgresOutcomeRepository) GetSumByCategory(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.CategorySum, error) {
	query := `
		SELECT c.id as category_id, COALESCE(SUM(o.amount), 0) as total
		FROM categories c
		LEFT JOIN outcomes o ON c.id = o.category_id`
	args := []interface{}{}
	argCount := 0

	if from != nil || to != nil {
		query += ` AND (`
		conditionsAdded := false

		if from != nil {
			argCount++
			query += `o.created_at >= $` + strconv.Itoa(argCount)
			args = append(args, *from)
			conditionsAdded = true
		}

		if to != nil {
			if conditionsAdded {
				query += ` AND `
			}
			argCount++
			query += `o.created_at <= $` + strconv.Itoa(argCount)
			args = append(args, *to)
		} else if from != nil {
			query += ` AND o.created_at <= NOW()`
		}

		query += `)`
	}

	if categoryId != 0 {
		argCount++
		query += ` WHERE c.id = $` + strconv.Itoa(argCount)
		args = append(args, categoryId)
	}

	query += ` GROUP BY c.id ORDER BY c.id`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sums []domain.CategorySum
	for rows.Next() {
		var s domain.CategorySum
		if err := rows.Scan(&s.CategoryId, &s.Total); err != nil {
			return nil, err
		}
		sums = append(sums, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sums, nil
}

func (r *PostgresOutcomeRepository) GetTotalSum(ctx context.Context, from *time.Time, to *time.Time) (int, error) {
	query := `SELECT SUM(amount) as total FROM outcomes WHERE 1 = 1`
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

	var total int
	err := r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
