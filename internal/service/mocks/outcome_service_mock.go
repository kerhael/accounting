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
	if cat, ok := args.Get(0).(*domain.Outcome); ok {
		return cat, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OutcomeService) GetAll(ctx context.Context, from *time.Time, to *time.Time) ([]domain.Outcome, error) {
	args := m.Called(ctx, from, to)

	var outcomes []domain.Outcome
	if args.Get(0) != nil {
		outcomes = args.Get(0).([]domain.Outcome)
	}

	return outcomes, args.Error(1)
}
