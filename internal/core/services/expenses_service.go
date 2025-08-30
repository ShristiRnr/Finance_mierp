// service/exchange_rate_service.go
package services

import (
    "context"
    "time"

    "github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
)

type ExchangeRateSvc struct {
    repo ports.ExchangeRateRepository
}

func NewExchangeRateService(repo ports.ExchangeRateRepository) ports.ExchangeRateService {
    return &ExchangeRateSvc{repo: repo}
}

func (s *ExchangeRateSvc) Create(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error) {
    return s.repo.Create(ctx, rate)
}

func (s *ExchangeRateSvc) Get(ctx context.Context, id uuid.UUID) (domain.ExchangeRate, error) {
    return s.repo.Get(ctx, id)
}

func (s *ExchangeRateSvc) Update(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error) {
    return s.repo.Update(ctx, rate)
}

func (s *ExchangeRateSvc) Delete(ctx context.Context, id uuid.UUID) error {
    return s.repo.Delete(ctx, id)
}

func (s *ExchangeRateSvc) List(ctx context.Context, base, quote *string, limit, offset int32) ([]domain.ExchangeRate, error) {
    return s.repo.List(ctx, base, quote, limit, offset)
}

func (s *ExchangeRateSvc) GetLatest(ctx context.Context, base, quote string, asOf time.Time) (domain.ExchangeRate, error) {
    return s.repo.GetLatest(ctx, base, quote, asOf)
}
