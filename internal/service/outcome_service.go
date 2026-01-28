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

type OutcomeServiceInterface interface {
	Create(ctx context.Context, name string, amount int, categoryId int, createdAt *time.Time) (*domain.Outcome, error)
	GetAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.Outcome, error)
	GetById(ctx context.Context, id int) (*domain.Outcome, error)
	Patch(ctx context.Context, id int, name string, amount int, categoryId int, createdAt *time.Time) (*domain.Outcome, error)
	DeleteById(ctx context.Context, id int) error
	GetSum(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.CategorySum, error)
	GetTotal(ctx context.Context, from *time.Time, to *time.Time) (int, error)
	GetSeries(ctx context.Context, from *time.Time, to *time.Time) ([]domain.MonthlySeries, error)
	GetTotalSeries(ctx context.Context, from *time.Time, to *time.Time) ([]domain.MonthlyTotalSeries, error)
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

func (s *OutcomeService) GetAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.Outcome, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	if categoryId != 0 {
		_, err := s.categoryRepo.FindById(ctx, categoryId)
		if err != nil {
			return nil, &domain.InvalidEntityError{
				UnderlyingCause: errors.New("invalid category"),
			}
		}
	}

	outcomes, err := s.repo.FindAll(ctx, from, to, categoryId)
	if err != nil {
		return nil, err
	}

	return outcomes, nil
}

func (s *OutcomeService) GetById(ctx context.Context, id int) (*domain.Outcome, error) {
	if id <= 0 {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid id"),
		}
	}

	outcome, err := s.repo.FindById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &domain.EntityNotFoundError{
				UnderlyingCause: err,
			}
		}
		return nil, err
	}

	return outcome, nil
}

func (s *OutcomeService) Patch(ctx context.Context, id int, name string, amount int, categoryId int, createdAt *time.Time) (*domain.Outcome, error) {
	outcome, err := s.repo.FindById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &domain.EntityNotFoundError{
				UnderlyingCause: err,
			}
		}
		return nil, err
	}

	o := &domain.Outcome{
		ID: outcome.ID,
	}

	if name != "" {
		o.Name = name
	} else {
		o.Name = outcome.Name
	}

	if amount != 0 {
		o.Amount = amount
	} else {
		o.Amount = outcome.Amount
	}

	if categoryId != 0 {
		_, err := s.categoryRepo.FindById(ctx, categoryId)
		if err != nil {
			return nil, &domain.InvalidEntityError{
				UnderlyingCause: errors.New("invalid category"),
			}
		}
		o.CategoryId = categoryId
	} else {
		o.CategoryId = outcome.CategoryId
	}

	if createdAt != nil {
		o.CreatedAt = createdAt
	} else {
		o.CreatedAt = outcome.CreatedAt
	}

	errUpdt := s.repo.Update(ctx, o)
	if errUpdt != nil {
		return nil, errUpdt
	}

	return o, nil
}

func (s *OutcomeService) DeleteById(ctx context.Context, id int) error {
	if id <= 0 {
		return &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid id"),
		}
	}

	return s.repo.DeleteById(ctx, id)
}

func (s *OutcomeService) GetSum(ctx context.Context, from *time.Time, to *time.Time, categoryId int) ([]domain.CategorySum, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	if categoryId != 0 {
		_, err := s.categoryRepo.FindById(ctx, categoryId)
		if err != nil {
			return nil, &domain.InvalidEntityError{
				UnderlyingCause: errors.New("invalid category"),
			}
		}
	}

	return s.repo.GetSumByCategory(ctx, from, to, categoryId)
}

func (s *OutcomeService) GetTotal(ctx context.Context, from *time.Time, to *time.Time) (int, error) {
	if from != nil && to != nil && from.After(*to) {
		return 0, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	return s.repo.GetTotalSum(ctx, from, to)
}

func (s *OutcomeService) GetSeries(ctx context.Context, from *time.Time, to *time.Time) ([]domain.MonthlySeries, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	return s.repo.GetMonthlySeries(ctx, from, to)
}

func (s *OutcomeService) GetTotalSeries(ctx context.Context, from *time.Time, to *time.Time) ([]domain.MonthlyTotalSeries, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	return s.repo.GetMonthlyTotalSeries(ctx, from, to)
}
