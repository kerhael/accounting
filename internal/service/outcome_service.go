package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/infrastructure/repository"
)

type OutcomeServiceInterface interface {
	Create(ctx context.Context, name string, amount int, categoryId int, createdAt *time.Time) (*domain.Outcome, error)
}

type OutcomeService struct {
	repo         repository.OutcomeRepository
	categoryRepo repository.CategoryRepository
}

func NewOutcomeService(repo repository.OutcomeRepository, categoryRepo repository.CategoryRepository) *OutcomeService {
	return &OutcomeService{repo: repo, categoryRepo: categoryRepo}
}

func (s *OutcomeService) Create(ctx context.Context, name string, amount int, categoryId int, createdAt *time.Time) (*domain.Outcome, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid name"),
		}
	}

	if amount <= 0 {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid amount"),
		}
	}

	if categoryId == 0 {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid category"),
		}
	}
	_, err := s.categoryRepo.FindById(ctx, categoryId)
	if err != nil {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid category"),
		}
	}

	if createdAt == nil {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid creation date"),
		}
	}

	outcome := &domain.Outcome{
		Name:       name,
		CreatedAt:  createdAt,
		Amount:     amount,
		CategoryId: categoryId,
	}

	if err := s.repo.Create(ctx, outcome); err != nil {
		return nil, err
	}

	return outcome, nil
}
