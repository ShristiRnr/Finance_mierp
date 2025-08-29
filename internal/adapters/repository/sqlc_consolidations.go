package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db" // sqlc generated package
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
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
func (r *ConsolidationRepo) Create(ctx context.Context, c domain.Consolidation) (domain.Consolidation, error) {
	params := db.CreateConsolidationParams{
		EntityIds:   c.EntityIds,
		PeriodStart: c.PeriodStart,
		PeriodEnd:   c.PeriodEnd,
		Report:      c.Report,
	}
	sqlcCon, err := r.queries.CreateConsolidation(ctx, params)
	if err != nil {
		return domain.Consolidation{}, err
	}
	return mapSQLCToDomain(sqlcCon), nil
}

// Get retrieves a consolidation by its ID.
func (r *ConsolidationRepo) Get(ctx context.Context, id uuid.UUID) (domain.Consolidation, error) {
	sqlcCon, err := r.queries.GetConsolidation(ctx, id)
	if err != nil {
		return domain.Consolidation{}, err
	}
	return mapSQLCToDomain(sqlcCon), nil
}

// List retrieves consolidations for given entity IDs and period range with pagination.
func (r *ConsolidationRepo) List(ctx context.Context, entityIds []string, start, end time.Time, limit, offset int32) ([]domain.Consolidation, error) {
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

	domainItems := make([]domain.Consolidation, len(sqlcItems))
	for i, c := range sqlcItems {
		domainItems[i] = mapSQLCToDomain(c)
	}
	return domainItems, nil
}

// Delete removes a consolidation by its ID.
func (r *ConsolidationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteConsolidation(ctx, id)
}

// mapSQLCToDomain converts the SQLC Consolidation type to the domain Consolidation type.
func mapSQLCToDomain(c db.Consolidation) domain.Consolidation {
	var createdAt time.Time
	if c.CreatedAt.Valid {
		createdAt = c.CreatedAt.Time
	}

	var updatedAt time.Time
	if c.UpdatedAt.Valid {
		updatedAt = c.UpdatedAt.Time
	}

	var createdBy string
	if c.CreatedBy.Valid {
		createdBy = c.CreatedBy.String
	}

	var updatedBy string
	if c.UpdatedBy.Valid {
		updatedBy = c.UpdatedBy.String
	}

	var revision int32
	if c.Revision.Valid {
		revision = c.Revision.Int32
	}

	return domain.Consolidation{
		ID:          c.ID,
		PeriodStart: c.PeriodStart,
		PeriodEnd:   c.PeriodEnd,
		Report:      c.Report,
		CreatedAt:   createdAt,
		CreatedBy:   createdBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   updatedBy,
		Revision:    revision,
	}
}