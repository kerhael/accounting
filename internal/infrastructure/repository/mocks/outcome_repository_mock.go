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

func (m *OutcomeRepository) FindAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.Outcome, error) {
	args := m.Called(ctx, from, to, categoryId)

	var outcomes []domain.Outcome
	if args.Get(0) != nil {
		outcomes = args.Get(0).([]domain.Outcome)
	}

	return outcomes, args.Error(1)
}

func (m *OutcomeRepository) FindById(ctx context.Context, id int) (*domain.Outcome, error) {
	args := m.Called(ctx, id)

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

func (m *OutcomeRepository) DeleteById(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *OutcomeRepository) GetSumByCategory(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.CategorySum, error) {
	args := m.Called(ctx, from, to, categoryId)

	var sums []domain.CategorySum
	if args.Get(0) != nil {
		sums = args.Get(0).([]domain.CategorySum)
	}

	return sums, args.Error(1)
}
