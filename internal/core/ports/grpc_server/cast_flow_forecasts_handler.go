package grpc_server

import (
	"context"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

// CashFlowGRPCServer implements the gRPC server for cash flow forecasts
type CashFlowGRPCServer struct {
	pb.UnimplementedCashFlowServiceServer
	service ports.CashFlowService
}

// NewCashFlowGRPCServer creates a new gRPC server instance
func NewCashFlowGRPCServer(s ports.CashFlowService) *CashFlowGRPCServer {
	return &CashFlowGRPCServer{service: s}
}

// GenerateForecast handles GenerateForecast gRPC request
func (s *CashFlowGRPCServer) GenerateForecast(ctx context.Context, req *pb.CashFlowForecastRequest) (*pb.CashFlowForecastResponse, error) {
    result, err := s.service.GenerateForecastFromPeriod(ctx, req.GetPeriod())
    if err != nil {
        return nil, err
    }
    return &pb.CashFlowForecastResponse{ForecastDetails: result}, nil
}



func (s *CashFlowGRPCServer) GetForecast(ctx context.Context, req *pb.CashFlowForecastRequest) (*pb.CashFlowForecastResponse, error) {
	result, err := s.service.GetForecastFromPeriod(ctx, req.GetPeriod())
	if err != nil {
		return nil, err
	}

	return &pb.CashFlowForecastResponse{
		ForecastDetails: result,
	}, nil
}


// ListForecasts handles ListForecasts gRPC request
func (s *CashFlowGRPCServer) ListForecasts(ctx context.Context, req *pb.CashFlowForecastRequest) (*pb.CashFlowForecastResponse, error) {
	result, err := s.service.ListForecastsFromPeriod(ctx, req.GetPeriod())
	if err != nil {
		return nil, err
	}

	return &pb.CashFlowForecastResponse{
		ForecastDetails: result,
	}, nil
}

