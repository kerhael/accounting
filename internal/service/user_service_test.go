package service

import (
	"context"
	"errors"
	"testing"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/infrastructure/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).
		Return(nil).
		Run(func(args mock.Arguments) {
			u := args.Get(1).(*domain.User)
			u.ID = 1
		})

	user, err := svc.Create(ctx, "John", "Doe", "john@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "john@example.com", user.Email)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, "password123", user.PasswordHash)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_NormalizesEmail(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := svc.Create(ctx, "Jane", "Doe", "  JANE@EXAMPLE.COM  ", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "jane@example.com", user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_EmptyFirstName(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.Create(ctx, "", "Doe", "john@example.com", "password123")

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))
	assert.Equal(t, "firstName is required", invalidErr.UnderlyingCause.Error())

	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_WhitespaceFirstName(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.Create(ctx, "   ", "Doe", "john@example.com", "password123")

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))

	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_EmptyLastName(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.Create(ctx, "John", "", "john@example.com", "password123")

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))
	assert.Equal(t, "lastName is required", invalidErr.UnderlyingCause.Error())

	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_WhitespaceLastName(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.Create(ctx, "John", "   ", "john@example.com", "password123")

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))

	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_EmptyEmail(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.Create(ctx, "John", "Doe", "", "password123")

	assert.Nil(t, user)
	assert.Error(t, err)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_InvalidEmailFormat(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.Create(ctx, "John", "Doe", "not-an-email", "password123")

	assert.Nil(t, user)
	assert.Error(t, err)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_RepoError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()
	repoErr := errors.New("db failure")

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(repoErr)

	user, err := svc.Create(ctx, "John", "Doe", "john@example.com", "password123")

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "db failure", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindByEmail_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()
	expectedUser := &domain.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	mockRepo.On("FindByEmail", ctx, "john@example.com").Return(expectedUser, nil)

	user, err := svc.FindByEmail(ctx, "john@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.FirstName, user.FirstName)
	assert.Equal(t, expectedUser.LastName, user.LastName)
	assert.Equal(t, expectedUser.Email, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindByEmail_NormalizesEmail(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()
	expectedUser := &domain.User{
		ID:        1,
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@example.com",
	}

	mockRepo.On("FindByEmail", ctx, "jane@example.com").Return(expectedUser, nil)

	user, err := svc.FindByEmail(ctx, "  JANE@EXAMPLE.COM  ")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "jane@example.com", user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindByEmail_EmptyEmail(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.FindByEmail(ctx, "")

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))

	mockRepo.AssertNotCalled(t, "FindByEmail")
}

func TestUserService_FindByEmail_WhitespaceEmail(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.FindByEmail(ctx, "   ")

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))

	mockRepo.AssertNotCalled(t, "FindByEmail")
}

func TestUserService_FindByEmail_InvalidEmailFormat(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.FindByEmail(ctx, "not-an-email")

	assert.Nil(t, user)
	assert.Error(t, err)

	mockRepo.AssertNotCalled(t, "FindByEmail")
}

func TestUserService_FindByEmail_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()
	repoErr := errors.New("user not found")

	mockRepo.On("FindByEmail", ctx, "nonexistent@example.com").Return((*domain.User)(nil), repoErr)

	user, err := svc.FindByEmail(ctx, "nonexistent@example.com")

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindById_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()
	expectedUser := &domain.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	mockRepo.On("FindById", ctx, 1).Return(expectedUser, nil)

	user, err := svc.FindById(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.FirstName, user.FirstName)
	assert.Equal(t, expectedUser.LastName, user.LastName)
	assert.Equal(t, expectedUser.Email, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_FindById_InvalidID(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.FindById(ctx, 0)

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))
	assert.Equal(t, "invalid id", invalidErr.UnderlyingCause.Error())

	mockRepo.AssertNotCalled(t, "FindById")
}

func TestUserService_FindById_NegativeID(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.FindById(ctx, -1)

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))

	mockRepo.AssertNotCalled(t, "FindById")
}

func TestUserService_FindById_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()
	repoErr := errors.New("user not found")

	mockRepo.On("FindById", ctx, 999).Return((*domain.User)(nil), repoErr)

	user, err := svc.FindById(ctx, 999)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}
