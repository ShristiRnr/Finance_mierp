// service/exchange_rate_service.go
package services

import (
    "context"
    "time"
    "fmt"

    "github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
)

type ExchangeRateSvc struct {
    repo      ports.ExchangeRateRepository
    publisher ports.EventPublisher
}

func NewExchangeRateService(repo ports.ExchangeRateRepository, publisher ports.EventPublisher) ports.ExchangeRateService {
    return &ExchangeRateSvc{
        repo:      repo,
        publisher: publisher,
    }
}

func (s *ExchangeRateSvc) Create(ctx context.Context, rate db.ExchangeRate) (db.ExchangeRate, error) {
	r, err := s.repo.Create(ctx, rate)
	if err != nil {
		return r, err
	}

	// Publish typed event
	if err := s.publisher.PublishExchangeRateCreated(ctx, &r); err != nil {
		fmt.Printf("Kafka publish error (create): %v\n", err)
	}

	return r, nil
}


func (s *ExchangeRateSvc) Get(ctx context.Context, id uuid.UUID) (db.ExchangeRate, error) {
    return s.repo.Get(ctx, id)
}

func (s *ExchangeRateSvc) Update(ctx context.Context, rate db.ExchangeRate) (db.ExchangeRate, error) {
	r, err := s.repo.Update(ctx, rate)
	if err != nil {
		return r, err
	}

	// Publish typed event
	if err := s.publisher.PublishExchangeRateUpdated(ctx, &r); err != nil {
		fmt.Printf("Kafka publish error (update): %v\n", err)
	}

	return r, nil
}

func (s *ExchangeRateSvc) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Publish typed event
	if err := s.publisher.PublishExchangeRateDeleted(ctx, id.String()); err != nil {
		fmt.Printf("Kafka publish error (delete): %v\n", err)
	}

	return nil
}

func (s *ExchangeRateSvc) List(ctx context.Context, base, quote *string, limit, offset int32) ([]db.ExchangeRate, error) {
    return s.repo.List(ctx, base, quote, limit, offset)
}

func (s *ExchangeRateSvc) GetLatest(ctx context.Context, base, quote string, asOf time.Time) (db.ExchangeRate, error) {
    return s.repo.GetLatest(ctx, base, quote, asOf)
}