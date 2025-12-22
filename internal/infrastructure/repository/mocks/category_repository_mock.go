package mocks

import (
	"context"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/stretchr/testify/mock"
)

type CategoryRepository struct {
	mock.Mock
}

func (m *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *CategoryRepository) FindById(ctx context.Context, id int) (*domain.Category, error) {
	args := m.Called(ctx, id)

	var category *domain.Category
	if args.Get(0) != nil {
		category = args.Get(0).(*domain.Category)
	}

	return category, args.Error(1)
}

func (m *CategoryRepository) DeleteById(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
