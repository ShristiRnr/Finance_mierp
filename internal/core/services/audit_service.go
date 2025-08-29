package services

import (
	"context"
	"fmt"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
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

func (s *AuditService) Record(ctx context.Context, event *domain.AuditEvent) (*domain.AuditEvent, error) {
	// In a real application, you might add validation or business logic here.
	if event.UserID == "" || event.Action == "" {
		return nil, fmt.Errorf("user_id and action are required")
	}
	return s.repo.RecordAuditEvent(ctx, event)
}

func (s *AuditService) GetByID(ctx context.Context, id string) (*domain.AuditEvent, error) {
	return s.repo.GetAuditEventByID(ctx, id)
}

func (s *AuditService) List(ctx context.Context, page domain.Pagination) ([]domain.AuditEvent, error) {
	return s.repo.ListAuditEvents(ctx, page)
}

func (s *AuditService) Filter(ctx context.Context, filter domain.FilterParams, page domain.Pagination) ([]domain.AuditEvent, error) {
	return s.repo.FilterAuditEvents(ctx, filter, page)
}