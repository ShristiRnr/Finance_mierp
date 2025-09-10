// internal/core/services/consolidation_service.go
package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type ConsolidationService struct {
	repo      ports.ConsolidationRepository
	publisher ports.EventPublisher
}

func NewConsolidationService(repo ports.ConsolidationRepository, publisher ports.EventPublisher) *ConsolidationService {
	return &ConsolidationService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *ConsolidationService) Create(ctx context.Context, c db.Consolidation) (db.Consolidation, error) {
	created, err := s.repo.Create(ctx, c)
	if err != nil {
		return db.Consolidation{}, err
	}

	// Publish event
	if s.publisher != nil {
		if err := s.publisher.PublishConsolidationCreated(ctx, &created); err != nil {
			// Log error but don't block main flow
			// Could use a logger here
			println("Kafka publish error (CreateConsolidation):", err.Error())
		}
	}

	return created, nil
}

func (s *ConsolidationService) Get(ctx context.Context, id uuid.UUID) (db.Consolidation, error) {
	return s.repo.Get(ctx, id)
}

func (s *ConsolidationService) List(ctx context.Context, entityIds []string, start, end time.Time, limit, offset int32) ([]db.Consolidation, error) {
	return s.repo.List(ctx, entityIds, start, end, limit, offset)
}

func (s *ConsolidationService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.publisher != nil {
		if err := s.publisher.PublishConsolidationDeleted(ctx, id.String()); err != nil {
			println("Kafka publish error (DeleteConsolidation):", err.Error())
		}
	}

	return nil
}