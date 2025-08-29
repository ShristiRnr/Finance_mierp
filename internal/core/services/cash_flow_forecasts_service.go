package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type CashFlowService struct {
	repo ports.CashFlowForecastRepository
}

func NewCashFlowService(repo ports.CashFlowForecastRepository) *CashFlowService {
	return &CashFlowService{repo: repo}
}

func (s *CashFlowService) GenerateForecast(ctx context.Context, cf *domain.CashFlowForecast) (*domain.CashFlowForecast, error) {
	return s.repo.Generate(ctx, cf)
}

func (s *CashFlowService) GetForecast(ctx context.Context, id uuid.UUID) (*domain.CashFlowForecast, error) {
	return s.repo.Get(ctx, id)
}

func (s *CashFlowService) ListForecasts(ctx context.Context, organizationID string, limit, offset int32) ([]*domain.CashFlowForecast, error) {
	return s.repo.List(ctx, organizationID, limit, offset)
}
