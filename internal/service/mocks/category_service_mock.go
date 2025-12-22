package mocks

import (
	"context"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/stretchr/testify/mock"
)

type CategoryService struct {
	mock.Mock
}

func (m *CategoryService) Create(ctx context.Context, label string) (*domain.Category, error) {
	args := m.Called(ctx, label)
	if cat, ok := args.Get(0).(*domain.Category); ok {
		return cat, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryService) GetById(ctx context.Context, id int) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if cat, ok := args.Get(0).(*domain.Category); ok {
		return cat, args.Error(1)
	}
	return nil, args.Error(1)
}
