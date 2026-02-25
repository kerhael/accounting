package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerhael/accounting/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, c *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindById(ctx context.Context, id int) (*domain.User, error)
	DeleteById(ctx context.Context, id int) error
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, u.FirstName, u.LastName, u.Email, u.PasswordHash).Scan(&u.ID)
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, first_name, last_name, email, password_hash, created_at FROM users WHERE email = $1 AND deleted_at IS NULL`

	row := r.db.QueryRow(ctx, query, email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindById(ctx context.Context, id int) (*domain.User, error) {
	var u domain.User

	query := `SELECT id, first_name, last_name, email, password_hash, created_at FROM users WHERE id = $1  AND deleted_at IS NULL`

	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresUserRepository) DeleteById(ctx context.Context, id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
