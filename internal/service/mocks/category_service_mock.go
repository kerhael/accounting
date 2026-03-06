package mocks

import (
	"context"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/stretchr/testify/mock"
)

type CategoryService struct {
	mock.Mock
}

func (m *CategoryService) Create(ctx context.Context, label string, userId int) (*domain.Category, error) {
	args := m.Called(ctx, label, userId)
	if cat, ok := args.Get(0).(*domain.Category); ok {
		return cat, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryService) GetAll(ctx context.Context, userId int) ([]domain.Category, error) {
	args := m.Called(ctx, userId)
	var categories []domain.Category
	if args.Get(0) != nil {
		categories = args.Get(0).([]domain.Category)
	}

	return categories, args.Error(1)
}

func (m *CategoryService) GetById(ctx context.Context, id int, userId int) (*domain.Category, error) {
	args := m.Called(ctx, id, userId)
	if cat, ok := args.Get(0).(*domain.Category); ok {
		return cat, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryService) DeleteById(ctx context.Context, id int, userId int) error {
	args := m.Called(ctx, id, userId)
	return args.Error(0)
}
