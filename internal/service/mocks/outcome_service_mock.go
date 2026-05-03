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

func (m *OutcomeService) Create(ctx context.Context, name string, amount int, categoryId int, createdAt *time.Time, userId int) (*domain.Outcome, error) {
	args := m.Called(ctx, name, amount, categoryId, createdAt, userId)
	if outcome, ok := args.Get(0).(*domain.Outcome); ok {
		return outcome, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OutcomeService) GetAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int, userId int, limit int, offset int) ([]domain.Outcome, int, error) {
	args := m.Called(ctx, from, to, categoryId, userId, limit, offset)

	var outcomes []domain.Outcome
	if args.Get(0) != nil {
		outcomes = args.Get(0).([]domain.Outcome)
	}

	var total int
	if args.Get(1) != nil {
		total = args.Get(1).(int)
	}

	return outcomes, total, args.Error(2)
}

func (m *OutcomeService) GetById(ctx context.Context, id int, userId int) (*domain.Outcome, error) {
	args := m.Called(ctx, id, userId)
	if outcome, ok := args.Get(0).(*domain.Outcome); ok {
		return outcome, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OutcomeService) PatchById(ctx context.Context, id int, name string, amount int, categoryId int, createdAt *time.Time, userId int) (*domain.Outcome, error) {
	args := m.Called(ctx, id, name, amount, categoryId, createdAt, userId)
	if outcome, ok := args.Get(0).(*domain.Outcome); ok {
		return outcome, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OutcomeService) DeleteById(ctx context.Context, id int, userId int) error {
	args := m.Called(ctx, id, userId)
	return args.Error(0)
}

func (m *OutcomeService) GetSum(ctx context.Context, from *time.Time, to *time.Time, categoryId int, userId int) ([]domain.CategorySum, error) {
	args := m.Called(ctx, from, to, categoryId, userId)

	var sums []domain.CategorySum
	if args.Get(0) != nil {
		sums = args.Get(0).([]domain.CategorySum)
	}

	return sums, args.Error(1)
}

func (m *OutcomeService) GetTotal(ctx context.Context, from *time.Time, to *time.Time, userId int) (int, error) {
	args := m.Called(ctx, from, to, userId)

	var total int
	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *OutcomeService) GetSeries(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.MonthlySeries, error) {
	args := m.Called(ctx, from, to, userId)

	var series []domain.MonthlySeries
	if args.Get(0) != nil {
		series = args.Get(0).([]domain.MonthlySeries)
	}

	return series, args.Error(1)
}

func (m *OutcomeService) GetTotalSeries(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.MonthlyTotalSeries, error) {
	args := m.Called(ctx, from, to, userId)

	var series []domain.MonthlyTotalSeries
	if args.Get(0) != nil {
		series = args.Get(0).([]domain.MonthlyTotalSeries)
	}

	return series, args.Error(1)
}
