package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/infrastructure/repository"
)

type IncomeServiceInterface interface {
	Create(ctx context.Context, name string, amount int, createdAt *time.Time, userId int) (*domain.Income, error)
	GetAll(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.Income, error)
	GetById(ctx context.Context, id int, userId int) (*domain.Income, error)
	PatchById(ctx context.Context, id int, name string, amount int, createdAt *time.Time, userId int) (*domain.Income, error)
	DeleteById(ctx context.Context, id int, userId int) error
}

type IncomeService struct {
	repo repository.IncomeRepository
}

func NewIncomeService(repo repository.IncomeRepository) *IncomeService {
	return &IncomeService{repo: repo}
}

func (s *IncomeService) Create(ctx context.Context, name string, amount int, createdAt *time.Time, userId int) (*domain.Income, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("name cannot be empty"),
		}
	}

	if amount <= 0 {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("amount must be greater than zero"),
		}
	}

	if createdAt == nil {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("creation date is required"),
		}
	}

	income := &domain.Income{
		Name:      name,
		CreatedAt: createdAt,
		Amount:    amount,
		UserId:    userId,
	}

	if err := s.repo.Create(ctx, income); err != nil {
		return nil, err
	}

	return income, nil
}

func (s *IncomeService) GetAll(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.Income, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	incomes, err := s.repo.FindAll(ctx, from, to, userId)
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func (s *IncomeService) GetById(ctx context.Context, id int, userId int) (*domain.Income, error) {
	if id <= 0 {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid id"),
		}
	}

	income, err := s.repo.FindById(ctx, id, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &domain.EntityNotFoundError{
				UnderlyingCause: err,
			}
		}
		return nil, err
	}

	return income, nil
}

func (s *IncomeService) PatchById(ctx context.Context, id int, name string, amount int, createdAt *time.Time, userId int) (*domain.Income, error) {
	income, err := s.repo.FindById(ctx, id, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &domain.EntityNotFoundError{
				UnderlyingCause: err,
			}
		}
		return nil, err
	}

	i := &domain.Income{
		ID:     income.ID,
		UserId: income.UserId,
	}

	if name != "" {
		i.Name = name
	} else {
		i.Name = income.Name
	}

	if amount != 0 {
		i.Amount = amount
	} else {
		i.Amount = income.Amount
	}

	if createdAt != nil {
		i.CreatedAt = createdAt
	} else {
		i.CreatedAt = income.CreatedAt
	}

	errUpdt := s.repo.Update(ctx, i)
	if errUpdt != nil {
		return nil, errUpdt
	}

	return i, nil
}

func (s *IncomeService) DeleteById(ctx context.Context, id int, userId int) error {
	if id <= 0 {
		return &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid id"),
		}
	}

	return s.repo.DeleteById(ctx, id, userId)
}
