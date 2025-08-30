package grpc_server

import (
	"context"
	"strconv"
	"strings"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ExchangeRateHandler struct {
	svc ports.ExchangeRateService
}

func NewExchangeRateHandler(svc ports.ExchangeRateService) *ExchangeRateHandler {
	return &ExchangeRateHandler{svc: svc}
}

func (h *ExchangeRateHandler) CreateExchangeRate(ctx context.Context, req *pb.CreateExchangeRateRequest) (*pb.ExchangeRate, error) {
    // convert float64 → string
    rateStr := strconv.FormatFloat(req.Rate.Rate, 'f', -1, 64)

    rate := domain.ExchangeRate{
		ID:            uuid.New(),
        BaseCurrency:  req.Rate.BaseCurrency,
        QuoteCurrency: req.Rate.QuoteCurrency,
        Rate:          rateStr,                 // ✅ now string
        AsOf:          req.Rate.AsOf.AsTime(),
        CreatedBy:     &req.Meta.AuthSubject,
    }

    created, err := h.svc.Create(ctx, rate)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create rate: %v", err)
    }
    return mapDomainToProtoExchangeRate(created), nil
}


func (h *ExchangeRateHandler) GetExchangeRate(ctx context.Context, req *pb.GetExchangeRateRequest) (*pb.ExchangeRate, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	rate, err := h.svc.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not found: %v", err)
	}
	return mapDomainToProtoExchangeRate(rate), nil
}

func (h *ExchangeRateHandler) UpdateExchangeRate(ctx context.Context, req *pb.UpdateExchangeRateRequest) (*pb.ExchangeRate, error) {
    id, err := uuid.Parse(req.Rate.Id)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
    }

    // Convert float64 → string
    rateStr := strconv.FormatFloat(req.Rate.Rate, 'f', -1, 64)

    rate := domain.ExchangeRate{
        ID:            id,
        BaseCurrency:  req.Rate.BaseCurrency,
        QuoteCurrency: req.Rate.QuoteCurrency,
        Rate:          rateStr,
        AsOf:          req.Rate.AsOf.AsTime(),
        UpdatedBy:     &req.Meta.AuthSubject,
    }

    updated, err := h.svc.Update(ctx, rate)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to update: %v", err)
    }
    return mapDomainToProtoExchangeRate(updated), nil
}


func (h *ExchangeRateHandler) DeleteExchangeRate(ctx context.Context, req *pb.DeleteExchangeRateRequest) (emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		return emptypb.Empty{}, status.Errorf(codes.Internal, "failed to delete: %v", err)
	}
	return emptypb.Empty{}, nil
}

func (h *ExchangeRateHandler) ListExchangeRates(ctx context.Context, req *pb.ListExchangeRatesRequest) (*pb.ListExchangeRatesResponse, error) {
    // default page size
    limit := int32(50)
    if req.Page.PageSize > 0 {
        limit = req.Page.PageSize
    }

    // decode page_token into offset
    var offset int32
    if req.Page.PageToken != "" {
        if o, err := strconv.Atoi(req.Page.PageToken); err == nil {
            offset = int32(o)
        }
    }

    // parse filter if you want (optional)
    var base, quote *string
    if req.Page.Filter != "" {
        // Example: parse "base_currency=USD AND quote_currency=EUR"
        filters := strings.Split(req.Page.Filter, "AND")
        for _, f := range filters {
            f = strings.TrimSpace(f)
            if strings.HasPrefix(f, "base_currency=") {
                val := strings.TrimPrefix(f, "base_currency=")
                val = strings.Trim(val, "'\" ")
                base = &val
            } else if strings.HasPrefix(f, "quote_currency=") {
                val := strings.TrimPrefix(f, "quote_currency=")
                val = strings.Trim(val, "'\" ")
                quote = &val
            }
        }
    }

    rates, err := h.svc.List(ctx, base, quote, limit, offset)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to list: %v", err)
    }

    var pbRates []*pb.ExchangeRate
    for _, r := range rates {
        pbRates = append(pbRates, mapDomainToProtoExchangeRate(r))
    }

    // generate next_page_token
    nextPageToken := strconv.Itoa(int(offset + int32(len(rates))))

    return &pb.ListExchangeRatesResponse{
        Rates: pbRates,
        Page: &pb.PageResponse{
            NextPageToken: nextPageToken,
            TotalSize:     int64(len(pbRates)), // optionally query total size
        },
    }, nil
}

