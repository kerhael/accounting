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
		WithArgs(123, 20, 0).
		WillReturnRows(rows)

	outcomes, err := repo.FindAll(context.Background(), nil, nil, 0, 123, 20, 0)

	assert.NoError(t, err)
	assert.Len(t, outcomes, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOutcomeRepository_CountAll(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	rows := pgxmock.NewRows([]string{"count"}).AddRow(2)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(123).
		WillReturnRows(rows)

	total, err := repo.CountAll(context.Background(), nil, nil, 0, 123)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)

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

func TestPostgresOutcomeRepository_GetSumByCategory(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	rows := pgxmock.NewRows([]string{"category_id", "total"}).
		AddRow(1, 1000).
		AddRow(2, 2000)

	mock.ExpectQuery("SELECT (.+) FROM categories").
		WithArgs(123).
		WillReturnRows(rows)

	sums, err := repo.GetSumByCategory(context.Background(), nil, nil, 0, 123)

	assert.NoError(t, err)
	assert.Len(t, sums, 2)
	assert.Equal(t, 1, sums[0].CategoryId)
	assert.Equal(t, 1000, sums[0].Total)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOutcomeRepository_GetTotalSum(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	rows := pgxmock.NewRows([]string{"total"}).AddRow(3000)

	mock.ExpectQuery("SELECT (.+) FROM outcomes").
		WithArgs(123).
		WillReturnRows(rows)

	total, err := repo.GetTotalSum(context.Background(), nil, nil, 123)

	assert.NoError(t, err)
	assert.Equal(t, 3000, total)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOutcomeRepository_GetMonthlySeries(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)

	rows := pgxmock.NewRows([]string{"month", "category_id", "total"}).
		AddRow("2023-01", 1, 1000).
		AddRow("2023-01", 2, 2000).
		AddRow("2023-02", 1, 0).
		AddRow("2023-02", 2, 0).
		AddRow("2023-03", 1, 3000).
		AddRow("2023-03", 2, 4000).
		AddRow("2023-04", 1, 0).
		AddRow("2023-04", 2, 0)

	mock.ExpectQuery("WITH months AS").
		WithArgs(from, to, 123).
		WillReturnRows(rows)

	series, err := repo.GetMonthlySeries(context.Background(), &from, &to, 123)

	assert.NoError(t, err)
	assert.Len(t, series, 4)
	assert.Equal(t, "2023-01", series[0].Month)
	assert.Equal(t, 2, len(series[0].Categories))
	assert.Equal(t, 1000, series[0].Categories[1])
	assert.Equal(t, 2000, series[0].Categories[2])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresOutcomeRepository_GetMonthlyTotalSeries(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewOutcomeRepository(mock)

	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)

	rows := pgxmock.NewRows([]string{"month", "total"}).
		AddRow("2023-01", 3000).
		AddRow("2023-02", 0).
		AddRow("2023-03", 7000).
		AddRow("2023-04", 0)

	mock.ExpectQuery("WITH months AS").
		WithArgs(from, to, 123).
		WillReturnRows(rows)

	series, err := repo.GetMonthlyTotalSeries(context.Background(), &from, &to, 123)

	assert.NoError(t, err)
	assert.Len(t, series, 4)
	assert.Equal(t, "2023-01", series[0].Month)
	assert.Equal(t, 3000, series[0].Total)

	assert.NoError(t, mock.ExpectationsWereMet())
}
