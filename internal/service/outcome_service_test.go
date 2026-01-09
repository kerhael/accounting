package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
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

func TestGetAllOutcomes_Success(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: 1,
			CreatedAt:  &time.Time{},
		},
		{
			ID:         2,
			Name:       "Groceries",
			Amount:     5000,
			CategoryId: 2,
			CreatedAt:  &time.Time{},
		},
	}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0).Return(expectedOutcomes, nil)

	outcomes, err := service.GetAll(ctx, nil, nil, 0)

	assert.NoError(t, err)
	assert.NotNil(t, outcomes)
	assert.Len(t, outcomes, 2)
	assert.Equal(t, expectedOutcomes[0].ID, outcomes[0].ID)
	assert.Equal(t, expectedOutcomes[0].Name, outcomes[0].Name)
	assert.Equal(t, expectedOutcomes[1].ID, outcomes[1].ID)
	assert.Equal(t, expectedOutcomes[1].Name, outcomes[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetAllOutcomes_InvalidDates(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	to := time.Now()
	from := to.Add(24 * time.Hour)

	outcomes, err := service.GetAll(ctx, &from, &to, 0)

	assert.Error(t, err)
	assert.Nil(t, outcomes)
	assert.IsType(t, &domain.InvalidDateError{}, err)

	// Repository should not be called since validation happens first
	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetAllOutcomes_CategoryNotFound(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	categoryId := 1
	mockCategoryRepo.On("FindById", ctx, categoryId).Return((*domain.Category)(nil), errors.New("not found"))

	outcomes, err := service.GetAll(ctx, nil, nil, categoryId)

	assert.Error(t, err)
	assert.Nil(t, outcomes)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	// Repository should not be called since validation happens first
	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetAllOutcomes_EmptyList(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	expectedOutcomes := []domain.Outcome{}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0).Return(expectedOutcomes, nil)

	outcomes, err := service.GetAll(ctx, nil, nil, 0)

	assert.NoError(t, err)
	assert.NotNil(t, outcomes)
	assert.Len(t, outcomes, 0)
	assert.Empty(t, outcomes)

	mockRepo.AssertExpectations(t)
}

func TestGetAllOutcomes_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0).Return([]domain.Outcome(nil), errors.New("repo error"))

	outcomes, err := service.GetAll(ctx, nil, nil, 0)

	assert.Error(t, err)
	assert.Nil(t, outcomes)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetById_Success(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	expectedOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(expectedOutcome, nil)

	outcome, err := service.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, outcome)
	assert.Equal(t, expectedOutcome.ID, outcome.ID)
	assert.Equal(t, expectedOutcome.Name, outcome.Name)
	assert.Equal(t, expectedOutcome.Amount, outcome.Amount)
	assert.Equal(t, expectedOutcome.CategoryId, outcome.CategoryId)
	assert.Equal(t, expectedOutcome.CreatedAt, outcome.CreatedAt)

	mockRepo.AssertExpectations(t)
}

func TestGetById_InvalidId_Zero(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	outcome, err := service.GetById(ctx, 0)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "FindById", mock.Anything, mock.Anything)
}

func TestGetById_InvalidId_Negative(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	outcome, err := service.GetById(ctx, -1)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "FindById", mock.Anything, mock.Anything)
}

func TestGetById_NotFound(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, 999).Return((*domain.Outcome)(nil), pgx.ErrNoRows)

	outcome, err := service.GetById(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.EntityNotFoundError{}, err)

	mockRepo.AssertExpectations(t)
}

func TestGetById_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	repoErr := errors.New("repo error")
	mockRepo.On("FindById", ctx, 1).Return((*domain.Outcome)(nil), repoErr)

	outcome, err := service.GetById(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestPatchOutcome_Success_NameOnly(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	existingOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Old Name",
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(existingOutcome, nil)

	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Outcome")).Return(nil).Run(func(args mock.Arguments) {
		updated := args.Get(1).(*domain.Outcome)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, 1000, updated.Amount)
		assert.Equal(t, 1, updated.CategoryId)
		assert.Equal(t, existingOutcome.CreatedAt, updated.CreatedAt)
	})

	outcome, err := service.Patch(ctx, 1, "New Name", 0, 0, nil)

	assert.NoError(t, err)
	assert.NotNil(t, outcome)
	assert.Equal(t, 1, outcome.ID)
	assert.Equal(t, "New Name", outcome.Name)
	assert.Equal(t, 1000, outcome.Amount)
	assert.Equal(t, 1, outcome.CategoryId)
	assert.Equal(t, existingOutcome.CreatedAt, outcome.CreatedAt)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestPatchOutcome_Success_AllFields(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	existingOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Old Name",
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(existingOutcome, nil)

	newCategory := &domain.Category{
		ID:    2,
		Label: "Transport",
	}
	mockCategoryRepo.On("FindById", ctx, 2).Return(newCategory, nil)

	newCreatedAt := time.Now()
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Outcome")).Return(nil).Run(func(args mock.Arguments) {
		updated := args.Get(1).(*domain.Outcome)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, 2000, updated.Amount)
		assert.Equal(t, 2, updated.CategoryId)
		assert.Equal(t, &newCreatedAt, updated.CreatedAt)
	})

	outcome, err := service.Patch(ctx, 1, "New Name", 2000, 2, &newCreatedAt)

	assert.NoError(t, err)
	assert.NotNil(t, outcome)
	assert.Equal(t, "New Name", outcome.Name)
	assert.Equal(t, 2000, outcome.Amount)
	assert.Equal(t, 2, outcome.CategoryId)
	assert.Equal(t, &newCreatedAt, outcome.CreatedAt)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestPatchOutcome_InvalidCategory(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	existingOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Old Name",
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(existingOutcome, nil)

	mockCategoryRepo.On("FindById", ctx, 999).Return((*domain.Category)(nil), errors.New("not found"))

	outcome, err := service.Patch(ctx, 1, "", 0, 999, nil)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestPatchOutcome_NotFound(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, 999).Return((*domain.Outcome)(nil), pgx.ErrNoRows)

	outcome, err := service.Patch(ctx, 999, "New Name", 0, 0, nil)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.EntityNotFoundError{}, err)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestPatchOutcome_UpdateError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	existingOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Old Name",
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(existingOutcome, nil)

	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Outcome")).Return(errors.New("update error"))

	outcome, err := service.Patch(ctx, 1, "New Name", 0, 0, nil)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.Equal(t, "update error", err.Error())

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}
