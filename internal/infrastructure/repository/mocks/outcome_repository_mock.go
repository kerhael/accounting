package mocks

import (
	"context"

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
