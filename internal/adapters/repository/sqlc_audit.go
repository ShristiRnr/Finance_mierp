package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"// Your generated sqlc package
	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)
// AuditRepository implements the repository.AuditRepository interface using sqlc.
type AuditRepository struct {
	queries *db.Queries
	publisher ports.EventPublisher
}

func NewAuditRepository(database *sql.DB, pub ports.EventPublisher) *AuditRepository {
	return &AuditRepository{
		queries: db.New(database),
		publisher: pub,
	}
}

// Helper to convert from sqlc model to domain model
func toDomain(e db.AuditEvent) db.AuditEvent {
	return db.AuditEvent{
		ID:           e.ID,
		UserID:       e.UserID,
		Action:       e.Action,
		Timestamp:    e.Timestamp,
		Details:      e.Details,
		ResourceType: e.ResourceType,
		ResourceID:   e.ResourceID,
	}
}

func (r *AuditRepository) RecordAuditEvent(ctx context.Context, event *db.AuditEvent) (*db.AuditEvent, error) {
	params := db.RecordAuditEventParams{
		UserID:       event.UserID,
		Action:       event.Action,
		Timestamp:    event.Timestamp,
		Details:      event.Details,
		ResourceType: event.ResourceType,
		ResourceID:   event.ResourceID,
	}
	recordedEvent, err := r.queries.RecordAuditEvent(ctx, params)
	if err != nil {
		return nil, err
	}

	res := toDomain(recordedEvent)

	return &res, nil
}

func (r *AuditRepository) GetAuditEventByID(ctx context.Context, id uuid.UUID) (*db.AuditEvent, error) {
	event, err := r.queries.GetAuditEventById(ctx, id) // use the parameter
	if err != nil {
		return nil, err
	}
	res := toDomain(event)
	return &res, nil
}


func (r *AuditRepository) ListAuditEvents(ctx context.Context, page db.Pagination) ([]db.AuditEvent, error) {
	params := db.ListAuditEventsParams{
		Limit:  int32(page.Limit),
		Offset: int32(page.Offset),
	}
	events, err := r.queries.ListAuditEvents(ctx, params)
	if err != nil {
		return nil, err
	}
	domainEvents := make([]db.AuditEvent, len(events))
	for i, e := range events {
		domainEvents[i] = toDomain(e)
	}
	return domainEvents, nil
}

func (r *AuditRepository) FilterAuditEvents(ctx context.Context, filter db.FilterParams, page db.Pagination) ([]db.AuditEvent, error) {
	params := db.FilterAuditEventsParams{
		Limit:  int32(page.Limit),
		Offset: int32(page.Offset),
	}

	// Set defaults for strings
	params.Column1 = ""
	params.Column2 = ""
	params.Column3 = ""
	params.Column4 = ""
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

	// Set defaults for dates
	params.Column5 = time.Time{}
	params.Column6 = time.Time{}
	if filter.FromDate != nil {
		params.Column5 = *filter.FromDate
	}
	if filter.ToDate != nil {
		params.Column6 = *filter.ToDate
	}

	events, err := r.queries.FilterAuditEvents(ctx, params)
	if err != nil {
		return nil, err
	}

	domainEvents := make([]db.AuditEvent, len(events))
	for i, e := range events {
		domainEvents[i] = toDomain(e)
	}
	return domainEvents, nil
}
