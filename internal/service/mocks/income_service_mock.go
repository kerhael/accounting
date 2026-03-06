package mocks

import (
	"context"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/stretchr/testify/mock"
)

type IncomeService struct {
	mock.Mock
}

func (m *IncomeService) Create(ctx context.Context, name string, amount int, createdAt *time.Time, userId int) (*domain.Income, error) {
	args := m.Called(ctx, name, amount, createdAt, userId)
	if income, ok := args.Get(0).(*domain.Income); ok {
		return income, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *IncomeService) GetAll(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.Income, error) {
	args := m.Called(ctx, from, to, userId)

	var incomes []domain.Income
	if args.Get(0) != nil {
		incomes = args.Get(0).([]domain.Income)
	}

	return incomes, args.Error(1)
}

func (m *IncomeService) GetById(ctx context.Context, id int, userId int) (*domain.Income, error) {
	args := m.Called(ctx, id, userId)
	if income, ok := args.Get(0).(*domain.Income); ok {
		return income, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *IncomeService) PatchById(ctx context.Context, id int, name string, amount int, createdAt *time.Time, userId int) (*domain.Income, error) {
	args := m.Called(ctx, id, name, amount, createdAt, userId)
	if income, ok := args.Get(0).(*domain.Income); ok {
		return income, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *IncomeService) DeleteById(ctx context.Context, id int, userId int) error {
	args := m.Called(ctx, id, userId)
	return args.Error(0)
}
