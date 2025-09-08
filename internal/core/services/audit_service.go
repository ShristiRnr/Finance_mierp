package services

import (
	"context"
	"fmt"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

// AuditService provides application logic for managing audit events.
type AuditService struct {
	repo ports.AuditRepository
}

// NewAuditService creates a new AuditService.
func NewAuditService(repo ports.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Record(ctx context.Context, event *db.AuditEvent) (*db.AuditEvent, error) {
	if event.UserID == "" || event.Action == "" {
		return nil, fmt.Errorf("user_id and action are required")
	}
	return s.repo.RecordAuditEvent(ctx, event)
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