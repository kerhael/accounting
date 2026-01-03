package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/infrastructure/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOutcome_Success(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	category := &domain.Category{
		ID:    1,
		Label: "Food",
	}
	mockCategoryRepo.On("FindById", ctx, category.ID).Return(category, nil)

	name := "Restaurant"
	amount := 1999
	categoryId := category.ID
	createdAt := time.Now()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Outcome")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.Outcome)
		arg.ID = 1
	})

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt)

	assert.NoError(t, err)
	assert.NotNil(t, outcome)
	assert.Equal(t, 1, outcome.ID)
	assert.Equal(t, name, outcome.Name)
	assert.Equal(t, amount, outcome.Amount)
	assert.Equal(t, categoryId, outcome.CategoryId)
	assert.Equal(t, createdAt, *outcome.CreatedAt)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestCreateOutcome_InvalidName(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	category := &domain.Category{
		ID:    1,
		Label: "Food",
	}
	mockCategoryRepo.On("FindById", ctx, category.ID).Return(category, nil)

	name := ""
	amount := 100
	categoryId := category.ID
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateOutcome_InvalidName_Whitespace(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	category := &domain.Category{
		ID:    1,
		Label: "Food",
	}
	mockCategoryRepo.On("FindById", ctx, category.ID).Return(category, nil)

	name := "   "
	amount := 100
	categoryId := category.ID
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateOutcome_InvalidAmount_Zero(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	category := &domain.Category{
		ID:    1,
		Label: "Food",
	}
	mockCategoryRepo.On("FindById", ctx, category.ID).Return(category, nil)

	name := "Restaurant"
	amount := 0
	categoryId := category.ID
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateOutcome_InvalidAmount_Negative(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	category := &domain.Category{
		ID:    1,
		Label: "Food",
	}
	mockCategoryRepo.On("FindById", ctx, category.ID).Return(category, nil)

	name := "Restaurant"
	amount := -1
	categoryId := category.ID
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateOutcome_InvalidCategoryId(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	name := "Restaurant"
	amount := 1999
	categoryId := 0
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateOutcome_CategoryNotFound(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	categoryId := 1
	mockCategoryRepo.On("FindById", ctx, categoryId).Return((*domain.Category)(nil), errors.New("not found"))

	name := "Restaurant"
	amount := 1999
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCreateOutcome_InvalidCreatedAt(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	category := &domain.Category{
		ID:    1,
		Label: "Food",
	}
	mockCategoryRepo.On("FindById", ctx, category.ID).Return(category, nil)

	name := "Restaurant"
	amount := 1999
	categoryId := category.ID
	var createdAt *time.Time = nil

	outcome, err := service.Create(ctx, name, amount, categoryId, createdAt)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCreateOutcome_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	category := &domain.Category{
		ID:    1,
		Label: "Food",
	}
	mockCategoryRepo.On("FindById", ctx, category.ID).Return(category, nil)

	name := "Restaurant"
	amount := 1999
	categoryId := category.ID
	createdAt := time.Now()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Outcome")).Return(errors.New("repo error"))

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}
