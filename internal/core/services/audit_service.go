package services

import (
	"context"
	"fmt"
	"log"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

// AuditService provides application logic for managing audit events.
type AuditService struct {
	repo ports.AuditRepository
	publisher ports.EventPublisher
}

// NewAuditService creates a new AuditService.
func NewAuditService(repo ports.AuditRepository, publisher ports.EventPublisher) *AuditService {
	return &AuditService{repo: repo, publisher:  publisher}
}


func (s *AuditService) Record(ctx context.Context, event *db.AuditEvent) (*db.AuditEvent, error) {
	if event == nil {
		return nil, fmt.Errorf("event cannot be nil")
	}
	if event.UserID == "" || event.Action == "" {
		return nil, fmt.Errorf("user_id and action are required")
	}

	recorded, err := s.repo.RecordAuditEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	// Publish Kafka event (async, best effort)
	if s.publisher != nil {
		go func(e *db.AuditEvent) {
			if err := s.publisher.PublishAuditRecorded(context.Background(), e); err != nil {
				log.Printf("failed to publish audit event: %v", err)
			}
		}(recorded)
	}

	return recorded, nil
}


func (s *AuditService) GetByID(ctx context.Context, id string) (*db.AuditEvent, error) {
	return s.repo.GetAuditEventByID(ctx, id)
}

func (s *AuditService) List(ctx context.Context, page db.Pagination) ([]db.AuditEvent, error) {
	return s.repo.ListAuditEvents(ctx, page)
}

func (s *AuditService) Filter(ctx context.Context, filter db.FilterParams, page db.Pagination) ([]db.AuditEvent, error) {
	return s.repo.FilterAuditEvents(ctx, filter, page)
}