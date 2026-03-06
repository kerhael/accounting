package mocks

import (
	"context"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/stretchr/testify/mock"
)

type OutcomeRepository struct {
	mock.Mock
}

func (m *OutcomeRepository) Create(ctx context.Context, o *domain.Outcome) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *OutcomeRepository) FindAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int, userId int) ([]domain.Outcome, error) {
	args := m.Called(ctx, from, to, categoryId, userId)

	var outcomes []domain.Outcome
	if args.Get(0) != nil {
		outcomes = args.Get(0).([]domain.Outcome)
	}

	return outcomes, args.Error(1)
}

func (m *OutcomeRepository) FindById(ctx context.Context, id int, userId int) (*domain.Outcome, error) {
	args := m.Called(ctx, id, userId)

	var outcome *domain.Outcome
	if args.Get(0) != nil {
		outcome = args.Get(0).(*domain.Outcome)
	}

	return outcome, args.Error(1)
}

func (m *OutcomeRepository) Update(ctx context.Context, o *domain.Outcome) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *OutcomeRepository) DeleteById(ctx context.Context, id int, userId int) error {
	args := m.Called(ctx, id, userId)
	return args.Error(0)
}

func (m *OutcomeRepository) GetSumByCategory(ctx context.Context, from *time.Time, to *time.Time, categoryId int, userId int) ([]domain.CategorySum, error) {
	args := m.Called(ctx, from, to, categoryId, userId)

	var sums []domain.CategorySum
	if args.Get(0) != nil {
		sums = args.Get(0).([]domain.CategorySum)
	}

	return sums, args.Error(1)
}

func (m *OutcomeRepository) GetTotalSum(ctx context.Context, from *time.Time, to *time.Time, userId int) (int, error) {
	args := m.Called(ctx, from, to, userId)

	var total int
	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *OutcomeRepository) GetMonthlySeries(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.MonthlySeries, error) {
	args := m.Called(ctx, from, to, userId)

	var series []domain.MonthlySeries
	if args.Get(0) != nil {
		series = args.Get(0).([]domain.MonthlySeries)
	}

	return series, args.Error(1)
}

func (m *OutcomeRepository) GetMonthlyTotalSeries(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.MonthlyTotalSeries, error) {
	args := m.Called(ctx, from, to, userId)

	var series []domain.MonthlyTotalSeries
	if args.Get(0) != nil {
		series = args.Get(0).([]domain.MonthlyTotalSeries)
	}

	return series, args.Error(1)
}
