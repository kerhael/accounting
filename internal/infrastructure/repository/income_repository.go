package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/kerhael/accounting/internal/domain"
)

type IncomeRepository interface {
	Create(ctx context.Context, c *domain.Income) error
	FindAll(ctx context.Context, from *time.Time, to *time.Time, userId int, limit int, offset int) ([]domain.Income, error)
	CountAll(ctx context.Context, from *time.Time, to *time.Time, userId int) (int, error)
	FindById(ctx context.Context, id int, userId int) (*domain.Income, error)
	Update(ctx context.Context, o *domain.Income) error
	DeleteById(ctx context.Context, id int, userId int) error
}

type PostgresIncomeRepository struct {
	db DB
}

func NewIncomeRepository(db DB) *PostgresIncomeRepository {
	return &PostgresIncomeRepository{db: db}
}

func (r *PostgresIncomeRepository) Create(ctx context.Context, i *domain.Income) error {
	query := `
		INSERT INTO incomes (name, amount, created_at, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, i.Name, i.Amount, &i.CreatedAt, i.UserId).Scan(&i.ID)
}

func (r *PostgresIncomeRepository) FindAll(ctx context.Context, from *time.Time, to *time.Time, userId int, limit int, offset int) ([]domain.Income, error) {
	query := `SELECT id, name, amount, created_at, user_id FROM incomes WHERE user_id = $1`
	args := []any{userId}
	argCount := 1

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

	query += ` ORDER BY created_at DESC, id DESC`
	argCount++
	query += ` LIMIT $` + strconv.Itoa(argCount)
	args = append(args, limit)
	argCount++
	query += ` OFFSET $` + strconv.Itoa(argCount)
	args = append(args, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []domain.Income
	for rows.Next() {
		var i domain.Income
		if err := rows.Scan(&i.ID, &i.Name, &i.Amount, &i.CreatedAt, &i.UserId); err != nil {
			return nil, err
		}
		incomes = append(incomes, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incomes, nil
}

func (r *PostgresIncomeRepository) CountAll(ctx context.Context, from *time.Time, to *time.Time, userId int) (int, error) {
	query := `SELECT COUNT(*) FROM incomes WHERE user_id = $1`
	args := []any{userId}
	argCount := 1

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

func (r *PostgresIncomeRepository) FindById(ctx context.Context, id int, userId int) (*domain.Income, error) {
	var i domain.Income

	query := `
		SELECT id, name, amount, created_at, user_id FROM incomes
		WHERE id = $1 AND user_id = $2
	`

	err := r.db.QueryRow(ctx, query, id, userId).Scan(&i.ID, &i.Name, &i.Amount, &i.CreatedAt, &i.UserId)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (r *PostgresIncomeRepository) Update(ctx context.Context, i *domain.Income) error {
	query := `UPDATE incomes SET name = $1, amount = $2, created_at = $3 WHERE id = $4 AND user_id = $5`

	_, err := r.db.Exec(ctx, query, i.Name, i.Amount, i.CreatedAt, i.ID, i.UserId)
	return err
}

func (r *PostgresIncomeRepository) DeleteById(ctx context.Context, id int, userId int) error {
	query := `
		DELETE FROM incomes
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.db.Exec(ctx, query, id, userId)
	return err
}
