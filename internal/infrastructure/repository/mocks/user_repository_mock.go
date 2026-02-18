package mocks

import (
	"context"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(ctx context.Context, c *domain.User) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)

	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}

	return user, args.Error(1)
}

func (m *UserRepository) FindById(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)

	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}

	return user, args.Error(1)
}

func (m *UserRepository) DeleteById(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
