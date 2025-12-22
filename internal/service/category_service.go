package service

import (
	"context"
	"errors"
	"strings"

	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/infrastructure/repository"
)

type CategoryServiceInterface interface {
	Create(ctx context.Context, label string) (*domain.Category, error)
}

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, label string) (*domain.Category, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("label is required"),
		}
	}

	category := &domain.Category{
		Label: label,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}
