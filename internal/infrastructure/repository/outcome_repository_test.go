package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPostgresOutcomeRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	ctx := context.Background()

	outcome := &domain.Outcome{
		Name:       "Test Outcome",
		Amount:     1000,
		CategoryId: 1,
		UserId:     123,
	}

	rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO outcomes").
		WithArgs("Test Outcome", 1000, 1, pgxmock.AnyArg(), 123).
		WillReturnRows(rows)

	err = repo.Create(ctx, outcome)

	assert.NoError(t, err)
	assert.Equal(t, 1, outcome.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOutcomeRepository_FindAll(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	now := time.Now()
	rows := pgxmock.NewRows(
		[]string{"id", "name", "amount", "category_id", "created_at", "user_id"},
	).
		AddRow(1, "Rent", 1000, 1, &now, 123).
		AddRow(2, "Food", 200, 2, &now, 123)

	mock.ExpectQuery("SELECT (.+) FROM outcomes").
		WithArgs(123).
		WillReturnRows(rows)

	outcomes, err := repo.FindAll(context.Background(), nil, nil, 0, 123)

	assert.NoError(t, err)
	assert.Len(t, outcomes, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOutcomeRepository_FindById(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	rows := pgxmock.NewRows(
		[]string{"id", "name", "amount", "category_id", "created_at", "user_id"},
	).AddRow(1, "Rent", 1000, 1, time.Now(), 123)

	mock.ExpectQuery("SELECT (.+) FROM outcomes").
		WithArgs(1, 123).
		WillReturnRows(rows)

	outcome, err := repo.FindById(context.Background(), 1, 123)

	assert.NoError(t, err)
	assert.Equal(t, "Rent", outcome.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOutcomeRepository_Update(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	outcome := &domain.Outcome{
		ID:         1,
		Name:       "Updated",
		Amount:     3000,
		CategoryId: 2,
		UserId:     123,
	}

	mock.ExpectExec("UPDATE outcomes").
		WithArgs("Updated", 3000, 2, pgxmock.AnyArg(), 1, 123).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(context.Background(), outcome)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOutcomeRepository_DeleteById(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	mock.ExpectExec("DELETE FROM outcomes").
		WithArgs(1, 123).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteById(context.Background(), 1, 123)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
