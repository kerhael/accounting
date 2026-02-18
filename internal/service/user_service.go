package service

import (
	"context"
	"errors"
	"strings"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/infrastructure/repository"
	"github.com/kerhael/accounting/pkg/security"
)

type UserServiceInterface interface {
	Create(ctx context.Context, firstName string, lastName string, email string, password string) (*domain.User, error)
}

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, firstName string, lastName string, email string, password string) (*domain.User, error) {
	firstName = strings.TrimSpace(firstName)
	if firstName == "" {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("firstName is required"),
		}
	}
	lastName = strings.TrimSpace(lastName)
	if lastName == "" {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("lastName is required"),
		}
	}
	email = security.NormalizeEmail(email)
	err := security.ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
