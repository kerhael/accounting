package mocks

import (
	"context"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/stretchr/testify/mock"
)

type IncomeRepository struct {
	mock.Mock
}

func (m *IncomeRepository) Create(ctx context.Context, o *domain.Income) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *IncomeRepository) FindAll(ctx context.Context, from *time.Time, to *time.Time) ([]domain.Income, error) {
	args := m.Called(ctx, from, to)

	var incomes []domain.Income
	if args.Get(0) != nil {
		incomes = args.Get(0).([]domain.Income)
	}

	return incomes, args.Error(1)
}

func (m *IncomeRepository) FindById(ctx context.Context, id int) (*domain.Income, error) {
	args := m.Called(ctx, id)

	var income *domain.Income
	if args.Get(0) != nil {
		income = args.Get(0).(*domain.Income)
	}

	return income, args.Error(1)
}

func (m *IncomeRepository) Update(ctx context.Context, o *domain.Income) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *IncomeRepository) DeleteById(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
