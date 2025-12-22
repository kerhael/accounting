package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/infrastructure/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCategory_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	label := "Food"

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Category")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.Category)
		arg.ID = 1
	})

	category, err := service.Create(ctx, label)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, 1, category.ID)
	assert.Equal(t, label, category.Label)

	mockRepo.AssertExpectations(t)
}

func TestCreateCategory_InvalidLabel(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	category, err := service.Create(ctx, "  ")

	assert.Nil(t, category)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))
	assert.Equal(t, "label is required", invalidErr.UnderlyingCause.Error())
}

func TestCreateCategory_RepoError(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	label := "Travel"

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Category")).Return(errors.New("db failure"))

	category, err := service.Create(ctx, label)

	assert.Nil(t, category)
	assert.Error(t, err)

	assert.Equal(t, "db failure", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetCategoryById_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	category := &domain.Category{
		ID:    1,
		Label: "Books",
	}

	mockRepo.On("FindById", ctx, category.ID).Return(category, nil)

	c, err := service.GetById(ctx, category.ID)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, c.ID, category.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetCategoryById_InvalidId(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()

	category, err := service.GetById(ctx, -10)

	assert.Nil(t, category)
	assert.Error(t, err)

	assert.Equal(t, "invalid entity: invalid id", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetCategoryById_NotFound(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	id := 2

	mockRepo.On("FindById", ctx, id).Return(nil, pgx.ErrNoRows)

	category, err := service.GetById(ctx, id)

	assert.Nil(t, category)

	var notFoundErr *domain.EntityNotFoundError
	assert.True(t, errors.As(err, &notFoundErr), "expected EntityNotFoundError")

	mockRepo.AssertExpectations(t)
}
