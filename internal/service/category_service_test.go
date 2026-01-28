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

func TestDeleteById_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	id := 1

	mockRepo.On("DeleteById", ctx, id).Return(nil)

	err := service.DeleteById(ctx, id)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteById_InvalidId(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	invalidId := -1

	err := service.DeleteById(ctx, invalidId)

	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))
	assert.Equal(t, "invalid id", invalidErr.UnderlyingCause.Error())
}

func TestDeleteById_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	id := 1
	repoErr := errors.New("foreign key constraint violation")

	mockRepo.On("DeleteById", ctx, id).Return(repoErr)

	err := service.DeleteById(ctx, id)

	assert.Error(t, err)
	assert.Equal(t, repoErr.Error(), err.Error())

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

	assert.Equal(t, "invalid entity data: invalid id", err.Error())

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

func TestGetAllCategories_Success(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	category1 := domain.Category{
		ID:    1,
		Label: "Books",
	}
	category2 := domain.Category{
		ID:    2,
		Label: "Food",
	}

	mockRepo.On("FindAll", ctx).Return([]domain.Category{category1, category2}, nil)

	categories, err := service.GetAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, []domain.Category{category1, category2})

	assert.Equal(t, categories[0].ID, category1.ID)
	assert.Equal(t, categories[1].ID, category2.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCategories_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	service := NewCategoryService(mockRepo)

	ctx := context.Background()
	repoErr := errors.New("database connection failed")
	mockRepo.On("FindAll", ctx).Return(nil, repoErr)

	categories, err := service.GetAll(ctx)

	assert.Nil(t, categories)
	assert.Error(t, err)
	assert.Equal(t, repoErr.Error(), err.Error())

	mockRepo.AssertExpectations(t)
}
