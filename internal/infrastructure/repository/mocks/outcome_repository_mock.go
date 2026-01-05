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

func (m *OutcomeRepository) FindAll(ctx context.Context, from *time.Time, to *time.Time) ([]domain.Outcome, error) {
	args := m.Called(ctx, from, to)

	var outcomes []domain.Outcome
	if args.Get(0) != nil {
		outcomes = args.Get(0).([]domain.Outcome)
	}

	return outcomes, args.Error(1)
}
