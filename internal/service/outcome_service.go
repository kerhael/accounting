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
	Create(ctx context.Context, name string, amount int, categoryId int, createdAt *time.Time, userId int) (*domain.Outcome, error)
	GetAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int, userId int, limit int, offset int) ([]domain.Outcome, int, error)
	GetById(ctx context.Context, id int, userId int) (*domain.Outcome, error)
	PatchById(ctx context.Context, id int, name string, amount int, categoryId int, createdAt *time.Time, userId int) (*domain.Outcome, error)
	DeleteById(ctx context.Context, id int, userId int) error
	GetSum(ctx context.Context, from *time.Time, to *time.Time, categoryId int, userId int) ([]domain.CategorySum, error)
	GetTotal(ctx context.Context, from *time.Time, to *time.Time, userId int) (int, error)
	GetSeries(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.MonthlySeries, error)
	GetTotalSeries(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.MonthlyTotalSeries, error)
}

type OutcomeService struct {
	repo         repository.OutcomeRepository
	categoryRepo repository.CategoryRepository
}

func NewOutcomeService(repo repository.OutcomeRepository, categoryRepo repository.CategoryRepository) *OutcomeService {
	return &OutcomeService{repo: repo, categoryRepo: categoryRepo}
}

func (s *OutcomeService) Create(ctx context.Context, name string, amount int, categoryId int, createdAt *time.Time, userId int) (*domain.Outcome, error) {
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
	_, err := s.categoryRepo.FindById(ctx, categoryId, userId)
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
		UserId:     userId,
	}

	if err := s.repo.Create(ctx, outcome); err != nil {
		return nil, err
	}

	return outcome, nil
}

func (s *OutcomeService) GetAll(ctx context.Context, from *time.Time, to *time.Time, categoryId int, userId int, limit int, offset int) ([]domain.Outcome, int, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, 0, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	if categoryId != 0 {
		_, err := s.categoryRepo.FindById(ctx, categoryId, userId)
		if err != nil {
			return nil, 0, &domain.InvalidEntityError{
				UnderlyingCause: errors.New("invalid category"),
			}
		}
	}

	outcomes, err := s.repo.FindAll(ctx, from, to, categoryId, userId, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountAll(ctx, from, to, categoryId, userId)
	if err != nil {
		return nil, 0, err
	}

	return outcomes, total, nil
}

func (s *OutcomeService) GetById(ctx context.Context, id int, userId int) (*domain.Outcome, error) {
	if id <= 0 {
		return nil, &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid id"),
		}
	}

	outcome, err := s.repo.FindById(ctx, id, userId)
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

func (s *OutcomeService) PatchById(ctx context.Context, id int, name string, amount int, categoryId int, createdAt *time.Time, userId int) (*domain.Outcome, error) {
	outcome, err := s.repo.FindById(ctx, id, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &domain.EntityNotFoundError{
				UnderlyingCause: err,
			}
		}
		return nil, err
	}

	o := &domain.Outcome{
		ID:     outcome.ID,
		UserId: outcome.UserId,
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
		_, err := s.categoryRepo.FindById(ctx, categoryId, userId)
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

func (s *OutcomeService) DeleteById(ctx context.Context, id int, userId int) error {
	if id <= 0 {
		return &domain.InvalidEntityError{
			UnderlyingCause: errors.New("invalid id"),
		}
	}

	return s.repo.DeleteById(ctx, id, userId)
}

func (s *OutcomeService) GetSum(ctx context.Context, from *time.Time, to *time.Time, categoryId int, userId int) ([]domain.CategorySum, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	if categoryId != 0 {
		_, err := s.categoryRepo.FindById(ctx, categoryId, userId)
		if err != nil {
			return nil, &domain.InvalidEntityError{
				UnderlyingCause: errors.New("invalid category"),
			}
		}
	}

	return s.repo.GetSumByCategory(ctx, from, to, categoryId, userId)
}

func (s *OutcomeService) GetTotal(ctx context.Context, from *time.Time, to *time.Time, userId int) (int, error) {
	if from != nil && to != nil && from.After(*to) {
		return 0, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	return s.repo.GetTotalSum(ctx, from, to, userId)
}

func (s *OutcomeService) GetSeries(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.MonthlySeries, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	return s.repo.GetMonthlySeries(ctx, from, to, userId)
}

func (s *OutcomeService) GetTotalSeries(ctx context.Context, from *time.Time, to *time.Time, userId int) ([]domain.MonthlyTotalSeries, error) {
	if from != nil && to != nil && from.After(*to) {
		return nil, &domain.InvalidDateError{
			UnderlyingCause: errors.New("start date must be before end date"),
		}
	}

	return s.repo.GetMonthlyTotalSeries(ctx, from, to, userId)
}
