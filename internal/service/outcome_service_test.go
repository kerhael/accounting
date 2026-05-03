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

	userId := 123
	category := &domain.Category{
		ID:     1,
		Label:  "Food",
		UserId: userId,
	}
	mockCategoryRepo.On("FindById", ctx, category.ID, userId).Return(category, nil)

	name := "Restaurant"
	amount := 1999
	categoryId := category.ID
	createdAt := time.Now()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Outcome")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.Outcome)
		arg.ID = 1
	})

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt, userId)

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

	userId := 123
	category := &domain.Category{
		ID:     1,
		Label:  "Food",
		UserId: userId,
	}
	mockCategoryRepo.On("FindById", ctx, category.ID, userId).Return(category, nil)

	name := ""
	amount := 100
	categoryId := category.ID
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt, userId)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateOutcome_InvalidName_Whitespace(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	category := &domain.Category{
		ID:     1,
		Label:  "Food",
		UserId: userId,
	}
	mockCategoryRepo.On("FindById", ctx, category.ID, userId).Return(category, nil)

	name := "   "
	amount := 100
	categoryId := category.ID
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt, userId)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateOutcome_InvalidAmount_Zero(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	category := &domain.Category{
		ID:     1,
		Label:  "Food",
		UserId: userId,
	}
	mockCategoryRepo.On("FindById", ctx, category.ID, userId).Return(category, nil)

	name := "Restaurant"
	amount := 0
	categoryId := category.ID
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt, userId)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateOutcome_InvalidAmount_Negative(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	category := &domain.Category{
		ID:     1,
		Label:  "Food",
		UserId: userId,
	}
	mockCategoryRepo.On("FindById", ctx, category.ID, userId).Return(category, nil)

	name := "Restaurant"
	amount := -1
	categoryId := category.ID
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt, userId)

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

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt, 123)

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
	userId := 123
	mockCategoryRepo.On("FindById", ctx, categoryId, userId).Return((*domain.Category)(nil), errors.New("not found"))

	name := "Restaurant"
	amount := 1999
	createdAt := time.Now()

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt, userId)

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

	userId := 123
	category := &domain.Category{
		ID:     1,
		Label:  "Food",
		UserId: userId,
	}
	mockCategoryRepo.On("FindById", ctx, category.ID, userId).Return(category, nil)

	name := "Restaurant"
	amount := 1999
	categoryId := category.ID
	var createdAt *time.Time = nil

	outcome, err := service.Create(ctx, name, amount, categoryId, createdAt, userId)

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

	userId := 123
	category := &domain.Category{
		ID:     1,
		Label:  "Food",
		UserId: userId,
	}
	mockCategoryRepo.On("FindById", ctx, category.ID, userId).Return(category, nil)

	name := "Restaurant"
	amount := 1999
	categoryId := category.ID
	createdAt := time.Now()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Outcome")).Return(errors.New("repo error"))

	outcome, err := service.Create(ctx, name, amount, categoryId, &createdAt, userId)

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

	userId := 123
	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: 1,
			CreatedAt:  &time.Time{},
			UserId:     userId,
		},
		{
			ID:         2,
			Name:       "Groceries",
			Amount:     5000,
			CategoryId: 2,
			CreatedAt:  &time.Time{},
			UserId:     userId,
		},
	}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId, 20, 0).Return(expectedOutcomes, nil)
	mockRepo.On("CountAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return(2, nil)

	outcomes, total, err := service.GetAll(ctx, nil, nil, 0, userId, 20, 0)

	assert.NoError(t, err)
	assert.NotNil(t, outcomes)
	assert.Len(t, outcomes, 2)
	assert.Equal(t, 2, total)
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

	outcomes, total, err := service.GetAll(ctx, &from, &to, 0, 123, 20, 0)

	assert.Error(t, err)
	assert.Nil(t, outcomes)
	assert.Equal(t, 0, total)
	assert.IsType(t, &domain.InvalidDateError{}, err)

	// Repository should not be called since validation happens first
	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything, mock.Anything, mock.Anything, 123)
}

func TestGetAllOutcomes_CategoryNotFound(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	categoryId := 1
	userId := 123
	mockCategoryRepo.On("FindById", ctx, categoryId, userId).Return((*domain.Category)(nil), errors.New("not found"))

	outcomes, total, err := service.GetAll(ctx, nil, nil, categoryId, userId, 20, 0)

	assert.Error(t, err)
	assert.Nil(t, outcomes)
	assert.Equal(t, 0, total)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	// Repository should not be called since validation happens first
	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything, mock.Anything, mock.Anything, userId)
}

func TestGetAllOutcomes_EmptyList(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	expectedOutcomes := []domain.Outcome{}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, 123, 20, 0).Return(expectedOutcomes, nil)
	mockRepo.On("CountAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, 123).Return(0, nil)

	outcomes, total, err := service.GetAll(ctx, nil, nil, 0, 123, 20, 0)

	assert.NoError(t, err)
	assert.NotNil(t, outcomes)
	assert.Len(t, outcomes, 0)
	assert.Empty(t, outcomes)
	assert.Equal(t, 0, total)

	mockRepo.AssertExpectations(t)
}

func TestGetAllOutcomes_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId, 20, 0).Return([]domain.Outcome(nil), errors.New("repo error"))

	outcomes, total, err := service.GetAll(ctx, nil, nil, 0, userId, 20, 0)

	assert.Error(t, err)
	assert.Nil(t, outcomes)
	assert.Equal(t, 0, total)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetAllOutcomes_CountError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	expectedOutcomes := []domain.Outcome{
		{
			ID:         1,
			Name:       "Restaurant",
			Amount:     1999,
			CategoryId: 1,
			CreatedAt:  &time.Time{},
			UserId:     userId,
		},
	}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId, 20, 0).Return(expectedOutcomes, nil)
	mockRepo.On("CountAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, userId).Return(0, errors.New("count error"))

	outcomes, total, err := service.GetAll(ctx, nil, nil, 0, userId, 20, 0)

	assert.Error(t, err)
	assert.Nil(t, outcomes)
	assert.Equal(t, 0, total)
	assert.Equal(t, "count error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetById_Success(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	expectedOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Restaurant",
		Amount:     1999,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
		UserId:     userId,
	}
	mockRepo.On("FindById", ctx, 1, userId).Return(expectedOutcome, nil)

	outcome, err := service.GetById(ctx, 1, userId)

	assert.NoError(t, err)
	assert.NotNil(t, outcome)
	assert.Equal(t, expectedOutcome.ID, outcome.ID)
	assert.Equal(t, expectedOutcome.Name, outcome.Name)
	assert.Equal(t, expectedOutcome.Amount, outcome.Amount)
	assert.Equal(t, expectedOutcome.CategoryId, outcome.CategoryId)
	assert.Equal(t, expectedOutcome.CreatedAt, outcome.CreatedAt)
	assert.Equal(t, expectedOutcome.UserId, outcome.UserId)

	mockRepo.AssertExpectations(t)
}

func TestGetById_InvalidId_Zero(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	outcome, err := service.GetById(ctx, 0, 123)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "FindById", mock.Anything, mock.Anything, 123)
}

func TestGetById_InvalidId_Negative(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	outcome, err := service.GetById(ctx, -1, 123)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "FindById", mock.Anything, mock.Anything, 123)
}

func TestGetById_NotFound(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, 999, 123).Return((*domain.Outcome)(nil), pgx.ErrNoRows)

	outcome, err := service.GetById(ctx, 999, 123)

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
	mockRepo.On("FindById", ctx, 1, 123).Return((*domain.Outcome)(nil), repoErr)

	outcome, err := service.GetById(ctx, 1, 123)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestPatchById_Success_NameOnly(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	existingOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Old Name",
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
		UserId:     userId,
	}
	mockRepo.On("FindById", ctx, 1, userId).Return(existingOutcome, nil)

	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Outcome")).Return(nil).Run(func(args mock.Arguments) {
		updated := args.Get(1).(*domain.Outcome)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, 1000, updated.Amount)
		assert.Equal(t, 1, updated.CategoryId)
		assert.Equal(t, existingOutcome.CreatedAt, updated.CreatedAt)
	})

	outcome, err := service.PatchById(ctx, 1, "New Name", 0, 0, nil, userId)

	assert.NoError(t, err)
	assert.NotNil(t, outcome)
	assert.Equal(t, 1, outcome.ID)
	assert.Equal(t, "New Name", outcome.Name)
	assert.Equal(t, 1000, outcome.Amount)
	assert.Equal(t, 1, outcome.CategoryId)
	assert.Equal(t, existingOutcome.CreatedAt, outcome.CreatedAt)
	assert.Equal(t, existingOutcome.UserId, outcome.UserId)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestPatchById_Success_AllFields(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	existingOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Old Name",
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
		UserId:     userId,
	}
	mockRepo.On("FindById", ctx, 1, userId).Return(existingOutcome, nil)

	newCategory := &domain.Category{
		ID:     2,
		Label:  "Transport",
		UserId: userId,
	}
	mockCategoryRepo.On("FindById", ctx, 2, userId).Return(newCategory, nil)

	newCreatedAt := time.Now()
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Outcome")).Return(nil).Run(func(args mock.Arguments) {
		updated := args.Get(1).(*domain.Outcome)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, 2000, updated.Amount)
		assert.Equal(t, 2, updated.CategoryId)
		assert.Equal(t, &newCreatedAt, updated.CreatedAt)
		assert.Equal(t, userId, updated.UserId)
	})

	outcome, err := service.PatchById(ctx, 1, "New Name", 2000, 2, &newCreatedAt, userId)

	assert.NoError(t, err)
	assert.NotNil(t, outcome)
	assert.Equal(t, "New Name", outcome.Name)
	assert.Equal(t, 2000, outcome.Amount)
	assert.Equal(t, 2, outcome.CategoryId)
	assert.Equal(t, &newCreatedAt, outcome.CreatedAt)
	assert.Equal(t, userId, outcome.UserId)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestPatchById_InvalidCategory(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	existingOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Old Name",
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
		UserId:     userId,
	}
	mockRepo.On("FindById", ctx, 1, userId).Return(existingOutcome, nil)

	mockCategoryRepo.On("FindById", ctx, 999, userId).Return((*domain.Category)(nil), errors.New("not found"))

	outcome, err := service.PatchById(ctx, 1, "", 0, 999, nil, userId)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestPatchById_NotFound(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	mockRepo.On("FindById", ctx, 999, userId).Return((*domain.Outcome)(nil), pgx.ErrNoRows)

	outcome, err := service.PatchById(ctx, 999, "New Name", 0, 0, nil, userId)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.IsType(t, &domain.EntityNotFoundError{}, err)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestPatchById_UpdateError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	existingOutcome := &domain.Outcome{
		ID:         1,
		Name:       "Old Name",
		Amount:     1000,
		CategoryId: 1,
		CreatedAt:  &time.Time{},
		UserId:     userId,
	}
	mockRepo.On("FindById", ctx, 1, userId).Return(existingOutcome, nil)

	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Outcome")).Return(errors.New("update error"))

	outcome, err := service.PatchById(ctx, 1, "New Name", 0, 0, nil, userId)

	assert.Error(t, err)
	assert.Nil(t, outcome)
	assert.Equal(t, "update error", err.Error())

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestOutcomeDeleteById_Success(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	mockRepo.On("DeleteById", ctx, 1, userId).Return(nil)

	err := service.DeleteById(ctx, 1, userId)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestOutcomeDeleteById_InvalidId_Zero(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	err := service.DeleteById(ctx, 0, 123)

	assert.Error(t, err)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "DeleteById", mock.Anything, mock.Anything, 123)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestOutcomeDeleteById_InvalidId_Negative(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	err := service.DeleteById(ctx, -1, 123)

	assert.Error(t, err)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "DeleteById", mock.Anything, mock.Anything, 123)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestOutcomeDeleteById_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	repoErr := errors.New("repo error")
	mockRepo.On("DeleteById", ctx, 1, userId).Return(repoErr)

	err := service.DeleteById(ctx, 1, userId)

	assert.Error(t, err)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetSum_Success_NoFilters(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	categorySums := []domain.CategorySum{
		{CategoryId: 1, Total: 3000},
		{CategoryId: 2, Total: 1500},
	}
	mockRepo.On("GetSumByCategory", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, 123).Return(categorySums, nil)

	result, err := service.GetSum(ctx, nil, nil, 0, 123)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, result[0].CategoryId)
	assert.Equal(t, 3000, result[0].Total)
	assert.Equal(t, 2, result[1].CategoryId)
	assert.Equal(t, 1500, result[1].Total)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetSum_Success_WithFilters(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	category := &domain.Category{ID: 1, Label: "Food", UserId: userId}
	mockCategoryRepo.On("FindById", ctx, 1, userId).Return(category, nil)

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	categorySums := []domain.CategorySum{
		{CategoryId: 1, Total: 3000},
	}
	mockRepo.On("GetSumByCategory", ctx, &from, &to, 1, userId).Return(categorySums, nil)

	result, err := service.GetSum(ctx, &from, &to, 1, userId)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, result[0].CategoryId)
	assert.Equal(t, 3000, result[0].Total)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestGetSum_InvalidDates(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	to := time.Now()
	from := to.Add(24 * time.Hour)

	result, err := service.GetSum(ctx, &from, &to, 0, 123)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &domain.InvalidDateError{}, err)

	mockRepo.AssertNotCalled(t, "GetSumByCategory")
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetSum_InvalidCategory(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	mockCategoryRepo.On("FindById", ctx, 999, 123).Return((*domain.Category)(nil), errors.New("not found"))

	result, err := service.GetSum(ctx, nil, nil, 999, 123)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "GetSumByCategory")
	mockCategoryRepo.AssertExpectations(t)
}

func TestGetSum_EmptyList(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	categorySums := []domain.CategorySum{}
	mockRepo.On("GetSumByCategory", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, 123).Return(categorySums, nil)

	result, err := service.GetSum(ctx, nil, nil, 0, 123)

	assert.NoError(t, err)
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetSum_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	mockRepo.On("GetSumByCategory", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 0, 123).Return([]domain.CategorySum(nil), errors.New("repo error"))

	result, err := service.GetSum(ctx, nil, nil, 0, 123)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetTotal_Success_NoFilters(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	expectedTotal := 4500
	mockRepo.On("GetTotalSum", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 123).Return(expectedTotal, nil)

	result, err := service.GetTotal(ctx, nil, nil, 123)

	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, result)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetTotal_Success_WithFilters(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedTotal := 3000
	userId := 123
	mockRepo.On("GetTotalSum", ctx, &from, &to, userId).Return(expectedTotal, nil)

	result, err := service.GetTotal(ctx, &from, &to, userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, result)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetTotal_InvalidDates(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	to := time.Now()
	from := to.Add(24 * time.Hour)
	userId := 123

	result, err := service.GetTotal(ctx, &from, &to, userId)

	assert.Error(t, err)
	assert.Equal(t, 0, result)
	assert.IsType(t, &domain.InvalidDateError{}, err)

	mockRepo.AssertNotCalled(t, "GetTotalSum")
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetTotal_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	mockRepo.On("GetTotalSum", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), 123).Return(0, errors.New("repo error"))

	result, err := service.GetTotal(ctx, nil, nil, 123)

	assert.Error(t, err)
	assert.Equal(t, 0, result)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetSeries_Success_NoFilters(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	expectedSeries := []domain.MonthlySeries{
		{
			Month: "2025-07",
			Categories: map[int]int{
				1: 3000,
				2: 1500,
				3: 0, // All categories included even with 0
			},
		},
		{
			Month: "2025-08",
			Categories: map[int]int{
				1: 2500,
				2: 0,
				3: 500,
			},
		},
	}
	mockRepo.On("GetMonthlySeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(expectedSeries, nil)

	result, err := service.GetSeries(ctx, nil, nil, userId)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "2025-07", result[0].Month)
	assert.Equal(t, map[int]int{1: 3000, 2: 1500, 3: 0}, result[0].Categories)
	assert.Equal(t, "2025-08", result[1].Month)
	assert.Equal(t, map[int]int{1: 2500, 2: 0, 3: 500}, result[1].Categories)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetSeries_Success_WithFilters(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedSeries := []domain.MonthlySeries{
		{
			Month: "2025-01",
			Categories: map[int]int{
				1: 3000,
			},
		},
	}
	mockRepo.On("GetMonthlySeries", ctx, &from, &to, userId).Return(expectedSeries, nil)

	result, err := service.GetSeries(ctx, &from, &to, userId)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "2025-01", result[0].Month)
	assert.Equal(t, map[int]int{1: 3000}, result[0].Categories)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetSeries_InvalidDates(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	to := time.Now()
	from := to.Add(24 * time.Hour)

	result, err := service.GetSeries(ctx, &from, &to, 123)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &domain.InvalidDateError{}, err)

	mockRepo.AssertNotCalled(t, "GetMonthlySeries")
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetSeries_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	mockRepo.On("GetMonthlySeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(nil, errors.New("repo error"))

	result, err := service.GetSeries(ctx, nil, nil, userId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetTotalSeries_Success_NoFilters(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	expectedSeries := []domain.MonthlyTotalSeries{
		{
			Month: "2025-07",
			Total: 3000,
		},
		{
			Month: "2025-08",
			Total: 2500,
		},
	}
	mockRepo.On("GetMonthlyTotalSeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(expectedSeries, nil)

	result, err := service.GetTotalSeries(ctx, nil, nil, userId)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "2025-07", result[0].Month)
	assert.Equal(t, 3000, result[0].Total)
	assert.Equal(t, "2025-08", result[1].Month)
	assert.Equal(t, 2500, result[1].Total)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetTotalSeries_Success_WithFilters(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedSeries := []domain.MonthlyTotalSeries{
		{
			Month: "2025-01",
			Total: 3000,
		},
	}
	mockRepo.On("GetMonthlyTotalSeries", ctx, &from, &to, userId).Return(expectedSeries, nil)

	result, err := service.GetTotalSeries(ctx, &from, &to, userId)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "2025-01", result[0].Month)
	assert.Equal(t, 3000, result[0].Total)

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetTotalSeries_InvalidDates(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	to := time.Now()
	from := to.Add(24 * time.Hour)

	result, err := service.GetTotalSeries(ctx, &from, &to, 123)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &domain.InvalidDateError{}, err)

	mockRepo.AssertNotCalled(t, "GetMonthlyTotalSeries")
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}

func TestGetTotalSeries_RepoError(t *testing.T) {
	mockRepo := new(mocks.OutcomeRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	service := NewOutcomeService(mockRepo, mockCategoryRepo)
	ctx := context.Background()

	userId := 123
	mockRepo.On("GetMonthlyTotalSeries", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time"), userId).Return(nil, errors.New("repo error"))

	result, err := service.GetTotalSeries(ctx, nil, nil, userId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertNotCalled(t, "FindById")
}
