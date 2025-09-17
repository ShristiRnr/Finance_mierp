package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

// AuditService provides application logic for managing audit events.
type AuditService struct {
	repo ports.AuditRepository
	publisher ports.EventPublisher
	async     bool
}

// NewAuditService creates a new AuditService.
func NewAuditService(repo ports.AuditRepository, publisher ports.EventPublisher) *AuditService {
	return &AuditService{repo: repo, publisher:  publisher}
}


func (s *AuditService) Record(ctx context.Context, event *db.AuditEvent) (*db.AuditEvent, error) {
    recorded, err := s.repo.RecordAuditEvent(ctx, event)
    if err != nil {
        return nil, err
    }

    if s.publisher != nil {
        if s.async {
            go s.publisher.PublishAuditRecorded(context.Background(), recorded)
        } else {
            _ = s.publisher.PublishAuditRecorded(ctx, recorded)
        }
    }
    return recorded, nil
}


func (s *AuditService) GetByID(ctx context.Context, id string) (*db.AuditEvent, error) {
    // Convert string ID to uuid.UUID
    uid, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid UUID: %w", err)
    }

    // Pass the uuid.UUID to repo
    return s.repo.GetAuditEventByID(ctx, uid)
}


func (s *AuditService) List(ctx context.Context, page db.Pagination) ([]db.AuditEvent, error) {
	return s.repo.ListAuditEvents(ctx, page)
}

func (s *AuditService) Filter(ctx context.Context, filter db.FilterParams, page db.Pagination) ([]db.AuditEvent, error) {
	return s.repo.FilterAuditEvents(ctx, filter, page)
}