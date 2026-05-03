package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPostgresIncomeRepository_Create(t *testing.T) {

	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewIncomeRepository(mock)

	ctx := context.Background()

	income := &domain.Income{
		Name:   "Test Income",
		Amount: 1000,
		UserId: 123,
	}

	rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO incomes").
		WithArgs("Test Income", 1000, pgxmock.AnyArg(), 123).
		WillReturnRows(rows)

	err = repo.Create(ctx, income)

	assert.NoError(t, err)
	assert.Equal(t, 1, income.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresIncomeRepository_FindAll(t *testing.T) {

	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewIncomeRepository(mock)

	now := time.Now()
	rows := pgxmock.NewRows(
		[]string{"id", "name", "amount", "created_at", "user_id"},
	).
		AddRow(1, "Salary", 2000, &now, 123).
		AddRow(2, "Freelance", 500, &now, 123)

	mock.ExpectQuery("SELECT (.+) FROM incomes").
		WithArgs(123, 20, 0).
		WillReturnRows(rows)

	incomes, err := repo.FindAll(context.Background(), nil, nil, 123, 20, 0)

	assert.NoError(t, err)
	assert.Len(t, incomes, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresIncomeRepository_CountAll(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewIncomeRepository(mock)

	rows := pgxmock.NewRows([]string{"count"}).AddRow(2)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(123).
		WillReturnRows(rows)

	total, err := repo.CountAll(context.Background(), nil, nil, 123)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresIncomeRepository_FindById(t *testing.T) {

	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewIncomeRepository(mock)

	now := time.Now()
	rows := pgxmock.NewRows(
		[]string{"id", "name", "amount", "created_at", "user_id"},
	).AddRow(1, "Salary", 2000, &now, 123)

	mock.ExpectQuery("SELECT (.+) FROM incomes").
		WithArgs(1, 123).
		WillReturnRows(rows)

	income, err := repo.FindById(context.Background(), 1, 123)

	assert.NoError(t, err)
	assert.Equal(t, "Salary", income.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresIncomeRepository_Update(t *testing.T) {

	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewIncomeRepository(mock)

	income := &domain.Income{
		ID:     1,
		Name:   "Updated",
		Amount: 3000,
		UserId: 123,
	}

	mock.ExpectExec("UPDATE incomes").
		WithArgs("Updated", 3000, pgxmock.AnyArg(), 1, 123).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(context.Background(), income)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresIncomeRepository_DeleteById(t *testing.T) {

	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewIncomeRepository(mock)

	mock.ExpectExec("DELETE FROM incomes").
		WithArgs(1, 123).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteById(context.Background(), 1, 123)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
