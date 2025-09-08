// repository/exchange_rate_repo.go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type ExchangeRateRepo struct {
    q *db.Queries
}

func NewExchangeRateRepo(q *db.Queries) ports.ExchangeRateRepository {
    return &ExchangeRateRepo{q: q}
}

func (r *ExchangeRateRepo) Create(ctx context.Context, rate db.ExchangeRate) (db.ExchangeRate, error) {
    arg := db.CreateExchangeRateParams{
        BaseCurrency:  rate.BaseCurrency,
        QuoteCurrency: rate.QuoteCurrency,
        Rate:          rate.Rate,
        AsOf:          rate.AsOf,
        CreatedBy:     rate.CreatedBy,
        UpdatedBy:     rate.UpdatedBy,
        Revision:      rate.Revision,
    }
    rec, err := r.q.CreateExchangeRate(ctx, arg)
    if err != nil {
        return db.ExchangeRate{}, err
    }
    return mapDBToDomain(rec), nil
}

func (r *ExchangeRateRepo) Get(ctx context.Context, id uuid.UUID) (db.ExchangeRate, error) {
    rec, err := r.q.GetExchangeRate(ctx, id)
    if err != nil {
        return db.ExchangeRate{}, err
    }
    return mapDBToDomain(rec), nil
}

func (r *ExchangeRateRepo) Update(ctx context.Context, rate db.ExchangeRate) (db.ExchangeRate, error) {
    arg := db.UpdateExchangeRateParams{
        ID:            rate.ID,
        BaseCurrency:  rate.BaseCurrency,
        QuoteCurrency: rate.QuoteCurrency,
        Rate:          rate.Rate,
        AsOf:          rate.AsOf,
        UpdatedBy:     rate.UpdatedBy,
        Revision:      rate.Revision,
    }
    rec, err := r.q.UpdateExchangeRate(ctx, arg)
    if err != nil {
        return db.ExchangeRate{}, err
    }
    return mapDBToDomain(rec), nil
}

func (r *ExchangeRateRepo) Delete(ctx context.Context, id uuid.UUID) error {
    return r.q.DeleteExchangeRate(ctx, id)
}

func (r *ExchangeRateRepo) List(ctx context.Context, base, quote *string, limit, offset int32) ([]db.ExchangeRate, error) {
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
    result := make([]db.ExchangeRate, len(recs))
    for i, rdb := range recs {
        result[i] = mapDBToDomain(rdb)
    }
    return result, nil
}

func (r *ExchangeRateRepo) GetLatest(ctx context.Context, base, quote string, asOf time.Time) (db.ExchangeRate, error) {
    rec, err := r.q.GetLatestRate(ctx, db.GetLatestRateParams{
        BaseCurrency: base,
        QuoteCurrency: quote,
        AsOf: asOf,
    })
    if err != nil {
        return db.ExchangeRate{}, err
    }
    return mapDBToDomain(rec), nil
}

func safePtr(ptr *string) string {
    if ptr == nil {
        return ""
    }
    return *ptr
}

func mapDBToDomain(rec db.ExchangeRate) db.ExchangeRate {
    return db.ExchangeRate{
        ID:            rec.ID,
        BaseCurrency:  rec.BaseCurrency,
        QuoteCurrency: rec.QuoteCurrency,
        Rate:          rec.Rate,
        AsOf:          rec.AsOf,
        CreatedAt:     rec.CreatedAt,
        UpdatedAt:     rec.UpdatedAt,
        CreatedBy:     rec.CreatedBy,
        UpdatedBy:     rec.UpdatedBy,
        Revision:      rec.Revision,
    }
}

