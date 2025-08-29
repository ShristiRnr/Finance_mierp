// internal/core/services/consolidation_service.go
package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type ConsolidationService struct {
	repo ports.ConsolidationRepository
}

func NewConsolidationService(repo ports.ConsolidationRepository) *ConsolidationService {
	return &ConsolidationService{repo: repo}
}

func (s *ConsolidationService) Create(ctx context.Context, c domain.Consolidation) (domain.Consolidation, error) {
	return s.repo.Create(ctx, c)
}

func (s *ConsolidationService) Get(ctx context.Context, id uuid.UUID) (domain.Consolidation, error) {
	return s.repo.Get(ctx, id)
}

func (s *ConsolidationService) List(ctx context.Context, entityIds []string, start, end time.Time, limit, offset int32) ([]domain.Consolidation, error) {
	return s.repo.List(ctx, entityIds, start, end, limit, offset)
}

func (s *ConsolidationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
