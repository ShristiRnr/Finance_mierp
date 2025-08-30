package services

import (
	"context"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
)

// creditDebitNoteService implements the CreditDebitNoteService port.
type creditDebitNoteService struct {
	repo ports.CreditDebitNoteRepository
}

// NewCreditDebitNoteService creates a new service instance.
func NewCreditDebitNoteService(repo ports.CreditDebitNoteRepository) ports.CreditDebitNoteService {
	return &creditDebitNoteService{repo: repo}
}

func (s *creditDebitNoteService) Create(ctx context.Context, note domain.CreditDebitNote) (domain.CreditDebitNote, error) {
	// Business logic would go here (e.g., validation, calculations).
	// For now, it's a direct pass-through to the repository.
	return s.repo.Create(ctx, note)
}

func (s *creditDebitNoteService) Get(ctx context.Context, id uuid.UUID) (domain.CreditDebitNote, error) {
	return s.repo.Get(ctx, id)
}

func (s *creditDebitNoteService) List(ctx context.Context, limit, offset int32) ([]domain.CreditDebitNote, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *creditDebitNoteService) Update(ctx context.Context, note domain.CreditDebitNote) (domain.CreditDebitNote, error) {
	return s.repo.Update(ctx, note)
}

func (s *creditDebitNoteService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}