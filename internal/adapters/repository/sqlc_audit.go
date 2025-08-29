package repository

import (
	"context"
	"database/sql"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"// Your generated sqlc package
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/google/uuid"

)
// AuditRepository implements the repository.AuditRepository interface using sqlc.
type AuditRepository struct {
	q *db.Queries
}

func NewAuditRepository(database *sql.DB) *AuditRepository {
	return &AuditRepository{
		q: db.New(database),
	}
}

// Helper to convert from sqlc model to domain model
func toDomain(e db.AuditEvent) domain.AuditEvent {
	return domain.AuditEvent{
		ID:           e.ID.String(),
		UserID:       e.UserID,
		Action:       e.Action,
		Timestamp:    e.Timestamp,
		Details:      e.Details.String,
		ResourceType: e.ResourceType.String,
		ResourceID:   e.ResourceID.String,
	}
}

func (r *AuditRepository) RecordAuditEvent(ctx context.Context, event *domain.AuditEvent) (*domain.AuditEvent, error) {
	params := db.RecordAuditEventParams{
		UserID:       event.UserID,
		Action:       event.Action,
		Timestamp:    event.Timestamp,
		Details:      sql.NullString{String: event.Details, Valid: event.Details != ""},
		ResourceType: sql.NullString{String: event.ResourceType, Valid: event.ResourceType != ""},
		ResourceID:   sql.NullString{String: event.ResourceID, Valid: event.ResourceID != ""},
	}
	recordedEvent, err := r.q.RecordAuditEvent(ctx, params)
	if err != nil {
		return nil, err
	}
	res := toDomain(recordedEvent)
	return &res, nil
}

func (r *AuditRepository) GetAuditEventByID(ctx context.Context, id string) (*domain.AuditEvent, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err // Invalid UUID format
	}
	event, err := r.q.GetAuditEventById(ctx, uuid)
	if err != nil {
		return nil, err
	}
	res := toDomain(event)
	return &res, nil
}

func (r *AuditRepository) ListAuditEvents(ctx context.Context, page domain.Pagination) ([]domain.AuditEvent, error) {
	params := db.ListAuditEventsParams{
		Limit:  page.Limit,
		Offset: page.Offset,
	}
	events, err := r.q.ListAuditEvents(ctx, params)
	if err != nil {
		return nil, err
	}
	domainEvents := make([]domain.AuditEvent, len(events))
	for i, e := range events {
		domainEvents[i] = toDomain(e)
	}
	return domainEvents, nil
}

func (r *AuditRepository) FilterAuditEvents(ctx context.Context, filter domain.FilterParams, page domain.Pagination) ([]domain.AuditEvent, error) {
	params := db.FilterAuditEventsParams{
		Limit:  page.Limit,
		Offset: page.Offset,
	}
	if filter.UserID != nil {
		params.Column1 = *filter.UserID
	}
	if filter.Action != nil {
		params.Column2 = *filter.Action
	}
	if filter.ResourceType != nil {
		params.Column3 = *filter.ResourceType
	}
	if filter.ResourceID != nil {
		params.Column4 = *filter.ResourceID
	}
	if filter.FromDate != nil {
		params.Column5 = *filter.FromDate
	}
	if filter.ToDate != nil {
		params.Column6 = *filter.ToDate
	}
	
	events, err := r.q.FilterAuditEvents(ctx, params)
	if err != nil {
		return nil, err
	}

	domainEvents := make([]domain.AuditEvent, len(events))
	for i, e := range events {
		domainEvents[i] = toDomain(e)
	}
	return domainEvents, nil
}