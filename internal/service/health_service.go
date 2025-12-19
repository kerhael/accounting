package service

import (
	"context"

	"github.com/kerhael/accounting/internal/infrastructure/repository"
)

type HealthService struct {
	repo repository.HealthRepository
}

func NewHealthService(repo repository.HealthRepository) *HealthService {
	return &HealthService{repo: repo}
}

func (s *HealthService) Check(ctx context.Context) error {
	return s.repo.Check(ctx)
}
