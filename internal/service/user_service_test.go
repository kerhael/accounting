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

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))

	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_InvalidEmailFormat(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	svc := NewUserService(mockRepo)

	ctx := context.Background()

	user, err := svc.Create(ctx, "John", "Doe", "not-an-email", "password123")

	assert.Nil(t, user)
	assert.Error(t, err)

	var invalidErr *domain.InvalidEntityError
	assert.True(t, errors.As(err, &invalidErr))

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
