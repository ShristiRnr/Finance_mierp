package services

import (
	"context"
	"fmt"

	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
)

type CashFlowService struct {
	repo      ports.CashFlowForecastRepository
	publisher ports.EventPublisher
}

func NewCashFlowService(repo ports.CashFlowForecastRepository, publisher ports.EventPublisher) *CashFlowService {
	return &CashFlowService{
		repo:      repo,
		publisher: publisher,
	}
}

// Accepts a ReportPeriod and returns forecast details as string
func (s *CashFlowService) GenerateForecastFromPeriod(ctx context.Context, period *pb.ReportPeriod) (string, error) {
	if period == nil {
		return "", fmt.Errorf("period is required")
	}

	// Simulate generating forecast
	result := fmt.Sprintf("Generated forecast for period: %v", period)

	// Publish to Kafka
	err := s.publisher.Publish(ctx, "cash_flow_forecasts", "forecast.generated", []byte(result))
	if err != nil {
		// Log but don't block the main flow
		fmt.Printf("Kafka publish error: %v\n", err)
	}

	return result, nil
}

func (s *CashFlowService) GetForecastFromPeriod(ctx context.Context, period *pb.ReportPeriod) (string, error) {
	if period == nil {
		return "", fmt.Errorf("period is required")
	}

	// Simulate fetching forecast
	result := fmt.Sprintf("Fetched forecast for period: %v", period)

	// Publish to Kafka
	err := s.publisher.Publish(ctx, "cash_flow_forecasts", "forecast.fetched", []byte(result))
	if err != nil {
		fmt.Printf("Kafka publish error: %v\n", err)
	}

	return result, nil
}

func (s *CashFlowService) ListForecastsFromPeriod(ctx context.Context, period *pb.ReportPeriod) (string, error) {
	if period == nil {
		return "", fmt.Errorf("period is required")
	}

	// Simulate listing forecasts
	result := fmt.Sprintf("Listed forecasts for period: %v", period)

	// Publish to Kafka
	err := s.publisher.Publish(ctx, "cash_flow_forecasts", "forecast.listed", []byte(result))
	if err != nil {
		fmt.Printf("Kafka publish error: %v\n", err)
	}

	return result, nil
}

