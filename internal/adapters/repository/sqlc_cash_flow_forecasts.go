package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type CashFlowForecastRepository struct {
	queries *db.Queries
}

func NewCashFlowForecastRepo(dbConn *sql.DB) ports.CashFlowForecastRepository {
	return &CashFlowForecastRepository{
		queries: db.New(dbConn),
	}
}

func (r *CashFlowForecastRepository) Generate(ctx context.Context, cf *db.CashFlowForecast) (*db.CashFlowForecast, error) {
	row, err := r.queries.GenerateCashFlowForecast(ctx, db.GenerateCashFlowForecastParams{
		OrganizationID:  cf.OrganizationID,
		PeriodStart:     cf.PeriodStart,
		PeriodEnd:       cf.PeriodEnd,
		ForecastDetails: cf.ForecastDetails,
	})
	if err != nil {
		return nil, err
	}
	return mapCashFlowForecast(row), nil
}

func (r *CashFlowForecastRepository) Get(ctx context.Context, id uuid.UUID) (*db.CashFlowForecast, error) {
	row, err := r.queries.GetCashFlowForecast(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapCashFlowForecast(row), nil
}

func (r *CashFlowForecastRepository) List(ctx context.Context, organizationID string, limit, offset int32) ([]*db.CashFlowForecast, error) {
	rows, err := r.queries.ListCashFlowForecasts(ctx, db.ListCashFlowForecastsParams{
		OrganizationID: organizationID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*db.CashFlowForecast, len(rows))
	for i, row := range rows {
		result[i] = mapCashFlowForecast(row)
	}
	return result, nil
}

// ===================================================== Helpers ========================================================

func mapCashFlowForecast(row db.CashFlowForecast) *db.CashFlowForecast {

	return &db.CashFlowForecast{
		ID:              row.ID,
		OrganizationID:  row.OrganizationID,
		PeriodStart:     row.PeriodStart,
		PeriodEnd:       row.PeriodEnd,
		ForecastDetails: row.ForecastDetails,
		CreatedAt:       row.CreatedAt,
		CreatedBy:       row.CreatedBy,
		UpdatedAt:       row.UpdatedAt,
		UpdatedBy:       row.UpdatedBy,
		Revision:        row.Revision,
	}
}
