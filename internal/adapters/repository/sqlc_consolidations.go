package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db" // sqlc generated package
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

// ConsolidationRepo implements ports.ConsolidationRepository
type ConsolidationRepo struct {
	queries *db.Queries
}

// NewConsolidationRepo returns a new instance of ConsolidationRepo.
func NewConsolidationRepo(q *db.Queries) ports.ConsolidationRepository {
	return &ConsolidationRepo{queries: q}
}

// Create inserts a new consolidation record into the database.
func (r *ConsolidationRepo) Create(ctx context.Context, c db.Consolidation) (db.Consolidation, error) {
	params := db.CreateConsolidationParams{
		EntityIds:   c.EntityIds,
		PeriodStart: c.PeriodStart,
		PeriodEnd:   c.PeriodEnd,
		Report:      c.Report,
	}
	sqlcCon, err := r.queries.CreateConsolidation(ctx, params)
	if err != nil {
		return db.Consolidation{}, err
	}
	return mapSQLCToDomain(sqlcCon), nil
}

// Get retrieves a consolidation by its ID.
func (r *ConsolidationRepo) Get(ctx context.Context, id uuid.UUID) (db.Consolidation, error) {
	sqlcCon, err := r.queries.GetConsolidation(ctx, id)
	if err != nil {
		return db.Consolidation{}, err
	}
	return mapSQLCToDomain(sqlcCon), nil
}

// List retrieves consolidations for given entity IDs and period range with pagination.
func (r *ConsolidationRepo) List(ctx context.Context, entityIds []string, start, end time.Time, limit, offset int32) ([]db.Consolidation, error) {
	params := db.ListConsolidationsParams{
		PeriodStart: start,
		PeriodEnd:   end,
		Limit:       limit,
		Offset:      offset,
	}
	sqlcItems, err := r.queries.ListConsolidations(ctx, params)
	if err != nil {
		return nil, err
	}

	domainItems := make([]db.Consolidation, len(sqlcItems))
	for i, c := range sqlcItems {
		domainItems[i] = mapSQLCToDomain(c)
	}
	return domainItems, nil
}

// Delete removes a consolidation by its ID.
func (r *ConsolidationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteConsolidation(ctx, id)
}

func mapSQLCToDomain(c db.Consolidation) db.Consolidation {

	return db.Consolidation{
		ID:          c.ID,
		PeriodStart: c.PeriodStart,
		PeriodEnd:   c.PeriodEnd,
		Report:      c.Report,
		CreatedAt:   c.CreatedAt,
		CreatedBy:   c.CreatedBy,
		UpdatedAt:   c.UpdatedAt,
		UpdatedBy:   c.UpdatedBy,
		Revision:    c.Revision,
	}
}