// repository/exchange_rate_repo.go
package repository

import (
    "context"
    "database/sql"
    "time"

    "github.com/google/uuid"
    "github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type ExchangeRateRepo struct {
    q *db.Queries
}

func NewExchangeRateRepo(q *db.Queries) ports.ExchangeRateRepository {
    return &ExchangeRateRepo{q: q}
}

func (r *ExchangeRateRepo) Create(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error) {
    arg := db.CreateExchangeRateParams{
        BaseCurrency:  rate.BaseCurrency,
        QuoteCurrency: rate.QuoteCurrency,
        Rate:          rate.Rate,
        AsOf:          rate.AsOf,
        CreatedBy:     sql.NullString{String: safeString(rate.CreatedBy), Valid: rate.CreatedBy != nil},
        UpdatedBy:     sql.NullString{String: safeString(rate.UpdatedBy), Valid: rate.UpdatedBy != nil},
        Revision:      sql.NullString{String: safeString(rate.Revision), Valid: rate.Revision != nil},
    }
    rec, err := r.q.CreateExchangeRate(ctx, arg)
    if err != nil {
        return domain.ExchangeRate{}, err
    }
    return mapDBToDomain(rec), nil
}

func (r *ExchangeRateRepo) Get(ctx context.Context, id uuid.UUID) (domain.ExchangeRate, error) {
    rec, err := r.q.GetExchangeRate(ctx, id)
    if err != nil {
        return domain.ExchangeRate{}, err
    }
    return mapDBToDomain(rec), nil
}

func (r *ExchangeRateRepo) Update(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error) {
    arg := db.UpdateExchangeRateParams{
        ID:            rate.ID,
        BaseCurrency:  rate.BaseCurrency,
        QuoteCurrency: rate.QuoteCurrency,
        Rate:          rate.Rate,
        AsOf:          rate.AsOf,
        UpdatedBy:     sql.NullString{String: safeString(rate.UpdatedBy), Valid: rate.UpdatedBy != nil},
        Revision:      sql.NullString{String: safeString(rate.Revision), Valid: rate.Revision != nil},
    }
    rec, err := r.q.UpdateExchangeRate(ctx, arg)
    if err != nil {
        return domain.ExchangeRate{}, err
    }
    return mapDBToDomain(rec), nil
}

func (r *ExchangeRateRepo) Delete(ctx context.Context, id uuid.UUID) error {
    return r.q.DeleteExchangeRate(ctx, id)
}

func (r *ExchangeRateRepo) List(ctx context.Context, base, quote *string, limit, offset int32) ([]domain.ExchangeRate, error) {
    arg := db.ListExchangeRatesParams{
        Column1: safePtr(base),
        Column2: safePtr(quote),
        Limit:   limit,
        Offset:  offset,
    }
    recs, err := r.q.ListExchangeRates(ctx, arg)
    if err != nil {
        return nil, err
    }
    result := make([]domain.ExchangeRate, len(recs))
    for i, rdb := range recs {
        result[i] = mapDBToDomain(rdb)
    }
    return result, nil
}

func (r *ExchangeRateRepo) GetLatest(ctx context.Context, base, quote string, asOf time.Time) (domain.ExchangeRate, error) {
    rec, err := r.q.GetLatestRate(ctx, db.GetLatestRateParams{
        BaseCurrency: base,
        QuoteCurrency: quote,
        AsOf: asOf,
    })
    if err != nil {
        return domain.ExchangeRate{}, err
    }
    return mapDBToDomain(rec), nil
}

func safeString(ptr *string) string {
    if ptr == nil {
        return ""
    }
    return *ptr
}
func safePtr(ptr *string) string {
    if ptr == nil {
        return ""
    }
    return *ptr
}

func mapDBToDomain(rec db.ExchangeRate) domain.ExchangeRate {
    return domain.ExchangeRate{
        ID:            rec.ID,
        BaseCurrency:  rec.BaseCurrency,
        QuoteCurrency: rec.QuoteCurrency,
        Rate:          rec.Rate,
        AsOf:          rec.AsOf,
        CreatedAt:     rec.CreatedAt.Time,
        UpdatedAt:     rec.UpdatedAt.Time,
        CreatedBy:     nullableToPtr(rec.CreatedBy),
        UpdatedBy:     nullableToPtr(rec.UpdatedBy),
        Revision:      nullableToPtr(rec.Revision),
    }
}

func nullableToPtr(ns sql.NullString) *string {
    if ns.Valid {
        return &ns.String
    }
    return nil
}
