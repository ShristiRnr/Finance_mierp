package services

import (
	"context"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
)

type creditDebitNoteService struct {
	repo      ports.CreditDebitNoteRepository
	publisher ports.EventPublisher
}

// NewCreditDebitNoteService creates a new service instance with Kafka publisher.
func NewCreditDebitNoteService(repo ports.CreditDebitNoteRepository, publisher ports.EventPublisher) ports.CreditDebitNoteService {
	return &creditDebitNoteService{
		repo:      repo,
		publisher: publisher,
	}
}

// Create a new credit/debit note and publish event
// Create a new credit/debit note and publish event
func (s *creditDebitNoteService) Create(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error) {
	created, err := s.repo.Create(ctx, note)
	if err != nil {
		return db.CreditDebitNote{}, err
	}

	// Publish typed event
	if err := s.publisher.PublishCreditDebitNoteCreated(ctx, &created); err != nil {
		println("Kafka publish error:", err.Error())
	}

	return created, nil
}

// Get by ID (no event needed)
func (s *creditDebitNoteService) Get(ctx context.Context, id uuid.UUID) (db.CreditDebitNote, error) {
	return s.repo.Get(ctx, id)
}

// List (no event needed)
func (s *creditDebitNoteService) List(ctx context.Context, limit, offset int32) ([]db.CreditDebitNote, error) {
	return s.repo.List(ctx, limit, offset)
}

// Update a credit/debit note and publish event
func (s *creditDebitNoteService) Update(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error) {
	updated, err := s.repo.Update(ctx, note)
	if err != nil {
		return db.CreditDebitNote{}, err
	}

	if err := s.publisher.PublishCreditDebitNoteUpdated(ctx, &updated); err != nil {
		println("Kafka publish error:", err.Error())
	}

	return updated, nil
}

// Delete a credit/debit note and publish event
func (s *creditDebitNoteService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.publisher.PublishCreditDebitNoteDeleted(ctx, id.String()); err != nil {
		println("Kafka publish error:", err.Error())
	}

	return nil
}
