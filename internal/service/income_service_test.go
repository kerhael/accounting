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

func TestCreateIncome_Success(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	name := "Restaurant"
	amount := 1999
	createdAt := time.Now()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Income")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.Income)
		arg.ID = 1
	})

	income, err := service.Create(ctx, name, amount, &createdAt)

	assert.NoError(t, err)
	assert.NotNil(t, income)
	assert.Equal(t, 1, income.ID)
	assert.Equal(t, name, income.Name)
	assert.Equal(t, amount, income.Amount)
	assert.Equal(t, createdAt, *income.CreatedAt)

	mockRepo.AssertExpectations(t)
}

func TestCreateIncome_InvalidName(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	name := ""
	amount := 100
	createdAt := time.Now()

	income, err := service.Create(ctx, name, amount, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateIncome_InvalidName_Whitespace(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	name := "   "
	amount := 100
	createdAt := time.Now()

	income, err := service.Create(ctx, name, amount, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateIncome_InvalidAmount_Zero(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	name := "Restaurant"
	amount := 0
	createdAt := time.Now()

	income, err := service.Create(ctx, name, amount, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateIncome_InvalidAmount_Negative(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	name := "Restaurant"
	amount := -1
	createdAt := time.Now()

	income, err := service.Create(ctx, name, amount, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateIncome_InvalidCreatedAt(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	name := "Restaurant"
	amount := 1999
	var createdAt *time.Time = nil

	income, err := service.Create(ctx, name, amount, createdAt)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.InvalidEntityError{}, err)
}

func TestCreateIncome_RepoError(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	name := "Restaurant"
	amount := 1999
	createdAt := time.Now()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Income")).Return(errors.New("repo error"))

	income, err := service.Create(ctx, name, amount, &createdAt)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetAllIncomes_Success(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	expectedIncomes := []domain.Income{
		{
			ID:        1,
			Name:      "Restaurant",
			Amount:    1999,
			CreatedAt: &time.Time{},
		},
		{
			ID:        2,
			Name:      "Groceries",
			Amount:    5000,
			CreatedAt: &time.Time{},
		},
	}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return(expectedIncomes, nil)

	incomes, err := service.GetAll(ctx, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, incomes)
	assert.Len(t, incomes, 2)
	assert.Equal(t, expectedIncomes[0].ID, incomes[0].ID)
	assert.Equal(t, expectedIncomes[0].Name, incomes[0].Name)
	assert.Equal(t, expectedIncomes[1].ID, incomes[1].ID)
	assert.Equal(t, expectedIncomes[1].Name, incomes[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetAllIncomes_InvalidDates(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	to := time.Now()
	from := to.Add(24 * time.Hour)

	incomes, err := service.GetAll(ctx, &from, &to)

	assert.Error(t, err)
	assert.Nil(t, incomes)
	assert.IsType(t, &domain.InvalidDateError{}, err)

	// Repository should not be called since validation happens first
	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetAllIncomes_EmptyList(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	expectedIncomes := []domain.Income{}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return(expectedIncomes, nil)

	incomes, err := service.GetAll(ctx, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, incomes)
	assert.Len(t, incomes, 0)
	assert.Empty(t, incomes)

	mockRepo.AssertExpectations(t)
}

func TestGetAllIncomes_RepoError(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx, mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return([]domain.Income(nil), errors.New("repo error"))

	incomes, err := service.GetAll(ctx, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, incomes)
	assert.Equal(t, "repo error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetIncomeById_Success(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	expectedIncome := &domain.Income{
		ID:        1,
		Name:      "Restaurant",
		Amount:    1999,
		CreatedAt: &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(expectedIncome, nil)

	income, err := service.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, income)
	assert.Equal(t, expectedIncome.ID, income.ID)
	assert.Equal(t, expectedIncome.Name, income.Name)
	assert.Equal(t, expectedIncome.Amount, income.Amount)
	assert.Equal(t, expectedIncome.CreatedAt, income.CreatedAt)

	mockRepo.AssertExpectations(t)
}

func TestGetIncomeById_InvalidId_Zero(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	income, err := service.GetById(ctx, 0)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "FindById", mock.Anything, mock.Anything)
}

func TestGetIncomeById_InvalidId_Negative(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	income, err := service.GetById(ctx, -1)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "FindById", mock.Anything, mock.Anything)
}

func TestGetIncomeById_NotFound(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, 999).Return((*domain.Income)(nil), pgx.ErrNoRows)

	income, err := service.GetById(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.EntityNotFoundError{}, err)

	mockRepo.AssertExpectations(t)
}

func TestGetIncomeById_RepoError(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	repoErr := errors.New("repo error")
	mockRepo.On("FindById", ctx, 1).Return((*domain.Income)(nil), repoErr)

	income, err := service.GetById(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestPatchIncome_Success_NameOnly(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	existingIncome := &domain.Income{
		ID:        1,
		Name:      "Old Name",
		Amount:    1000,
		CreatedAt: &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(existingIncome, nil)

	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Income")).Return(nil).Run(func(args mock.Arguments) {
		updated := args.Get(1).(*domain.Income)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, 1000, updated.Amount)
		assert.Equal(t, existingIncome.CreatedAt, updated.CreatedAt)
	})

	income, err := service.Patch(ctx, 1, "New Name", 0, nil)

	assert.NoError(t, err)
	assert.NotNil(t, income)
	assert.Equal(t, 1, income.ID)
	assert.Equal(t, "New Name", income.Name)
	assert.Equal(t, 1000, income.Amount)
	assert.Equal(t, existingIncome.CreatedAt, income.CreatedAt)

	mockRepo.AssertExpectations(t)
}

func TestPatchIncome_Success_AllFields(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	existingIncome := &domain.Income{
		ID:        1,
		Name:      "Old Name",
		Amount:    1000,
		CreatedAt: &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(existingIncome, nil)

	newCreatedAt := time.Now()
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Income")).Return(nil).Run(func(args mock.Arguments) {
		updated := args.Get(1).(*domain.Income)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, 2000, updated.Amount)
		assert.Equal(t, &newCreatedAt, updated.CreatedAt)
	})

	income, err := service.Patch(ctx, 1, "New Name", 2000, &newCreatedAt)

	assert.NoError(t, err)
	assert.NotNil(t, income)
	assert.Equal(t, "New Name", income.Name)
	assert.Equal(t, 2000, income.Amount)
	assert.Equal(t, &newCreatedAt, income.CreatedAt)

	mockRepo.AssertExpectations(t)
}

func TestPatchIncome_NotFound(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindById", ctx, 999).Return((*domain.Income)(nil), pgx.ErrNoRows)

	income, err := service.Patch(ctx, 999, "New Name", 0, nil)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.IsType(t, &domain.EntityNotFoundError{}, err)

	mockRepo.AssertExpectations(t)
}

func TestPatchIncome_UpdateError(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	existingIncome := &domain.Income{
		ID:        1,
		Name:      "Old Name",
		Amount:    1000,
		CreatedAt: &time.Time{},
	}
	mockRepo.On("FindById", ctx, 1).Return(existingIncome, nil)

	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Income")).Return(errors.New("update error"))

	income, err := service.Patch(ctx, 1, "New Name", 0, nil)

	assert.Error(t, err)
	assert.Nil(t, income)
	assert.Equal(t, "update error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestIncomeDeleteById_Success(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	mockRepo.On("DeleteById", ctx, 1).Return(nil)

	err := service.DeleteById(ctx, 1)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestIncomeDeleteById_InvalidId_Zero(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	err := service.DeleteById(ctx, 0)

	assert.Error(t, err)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "DeleteById", mock.Anything, mock.Anything)
}

func TestIncomeDeleteById_InvalidId_Negative(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	err := service.DeleteById(ctx, -1)

	assert.Error(t, err)
	assert.IsType(t, &domain.InvalidEntityError{}, err)

	mockRepo.AssertNotCalled(t, "DeleteById", mock.Anything, mock.Anything)
}

func TestIncomeDeleteById_RepoError(t *testing.T) {
	mockRepo := new(mocks.IncomeRepository)
	service := NewIncomeService(mockRepo)
	ctx := context.Background()

	repoErr := errors.New("repo error")
	mockRepo.On("DeleteById", ctx, 1).Return(repoErr)

	err := service.DeleteById(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}
