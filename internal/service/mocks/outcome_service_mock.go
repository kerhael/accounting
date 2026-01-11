package mocks

import (
	"context"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/stretchr/testify/mock"
)

type OutcomeService struct {
	mock.Mock
}

func (m *OutcomeService) Create(ctx context.Context, name string, amount int, categoryId int, createdAt *time.Time) (*domain.Outcome, error) {
	args := m.Called(ctx, name, amount, categoryId, createdAt)
	if outcome, ok := args.Get(0).(*domain.Outcome); ok {
		return outcome, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OutcomeService) GetAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.Outcome, error) {
	args := m.Called(ctx, from, to, categoryId)

	var outcomes []domain.Outcome
	if args.Get(0) != nil {
		outcomes = args.Get(0).([]domain.Outcome)
	}

	return outcomes, args.Error(1)
}

func (m *OutcomeService) GetById(ctx context.Context, id int) (*domain.Outcome, error) {
	args := m.Called(ctx, id)
	if outcome, ok := args.Get(0).(*domain.Outcome); ok {
		return outcome, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OutcomeService) Patch(ctx context.Context, id int, name string, amount int, categoryId int, createdAt *time.Time) (*domain.Outcome, error) {
	args := m.Called(ctx, id, name, amount, categoryId, createdAt)
	if outcome, ok := args.Get(0).(*domain.Outcome); ok {
		return outcome, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OutcomeService) DeleteById(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *OutcomeService) GetSum(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.CategorySum, error) {
	args := m.Called(ctx, from, to, categoryId)

	var sums []domain.CategorySum
	if args.Get(0) != nil {
		sums = args.Get(0).([]domain.CategorySum)
	}

	return sums, args.Error(1)
}

func (m *OutcomeService) GetTotal(ctx context.Context, from *time.Time, to *time.Time) (int, error) {
	args := m.Called(ctx, from, to)

	var total int
	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *OutcomeService) GetSeries(ctx context.Context, from *time.Time, to *time.Time) ([]domain.MonthlySeries, error) {
	args := m.Called(ctx, from, to)

	var series []domain.MonthlySeries
	if args.Get(0) != nil {
		series = args.Get(0).([]domain.MonthlySeries)
	}

	return series, args.Error(1)
}

func (m *OutcomeService) GetTotalSeries(ctx context.Context, from *time.Time, to *time.Time) ([]domain.MonthlyTotalSeries, error) {
	args := m.Called(ctx, from, to)

	var series []domain.MonthlyTotalSeries
	if args.Get(0) != nil {
		series = args.Get(0).([]domain.MonthlyTotalSeries)
	}

	return series, args.Error(1)
}
