package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/infrastructure/repository"
)

type CategoryServiceInterface interface {
	Create(ctx context.Context, label string) (*domain.Category, error)
	GetById(ctx context.Context, id int) (*domain.Category, error)
	DeleteById(ctx context.Context, id int) error
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

func (s *CategoryService) GetById(ctx context.Context, id int) (*domain.Category, error) {
	if id <= 0 {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid id"),
		}
	}

	category, err := s.repo.FindById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &domain.EntityNotFoundError{
				UnderlyingCause: err,
			}
		}
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) DeleteById(ctx context.Context, id int) error {
	if id <= 0 {
		return &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid id"),
		}
	}

	return s.repo.DeleteById(ctx, id)
}
