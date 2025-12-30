package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerhael/accounting/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
	FindAll(ctx context.Context) ([]domain.Category, error)
	FindById(ctx context.Context, id int) (*domain.Category, error)
	DeleteById(ctx context.Context, id int) error
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

func (r *PostgresCategoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	query := `SELECT id, label FROM categories ORDER BY label`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Label); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) FindById(ctx context.Context, id int) (*domain.Category, error) {
	var c domain.Category

	query := `
		SELECT id, label FROM categories
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Label)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *PostgresCategoryRepository) DeleteById(ctx context.Context, id int) error {
	query := `
		DELETE FROM categories
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
