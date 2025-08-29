package grpc_server

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/google/uuid"
	money "google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccrualHandler struct {
	pb.UnimplementedAccrualServiceServer
	svc *services.AccrualService
}

func NewAccrualHandler(svc *services.AccrualService) *AccrualHandler {
	return &AccrualHandler{svc: svc}
}

func (h *AccrualHandler) CreateAccrual(ctx context.Context, req *pb.CreateAccrualRequest) (*pb.Accrual, error) {
	id, err := uuid.Parse(req.Accrual.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}
	acc, err := h.svc.Create(ctx, domain.Accrual{
		ID:          id,
		Description: stringPtr(req.Accrual.Description),
		Amount:      moneyToString(req.Accrual.Amount),
		AccrualDate: req.Accrual.AccrualDate.AsTime(),
		AccountID:   req.Accrual.AccountId,
		UpdatedBy:   getUserFromContext(ctx),
	})
	if err != nil {
		return nil, err

	}
	return toPbAccrual(acc), nil
}

func (h *AccrualHandler) GetAccrualById(ctx context.Context, req *pb.GetAccrualByIdRequest) (*pb.Accrual, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}
	acc, err := h.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toPbAccrual(acc), nil
}

func (h *AccrualHandler) UpdateAccrual(ctx context.Context, req *pb.UpdateAccrualRequest) (*pb.Accrual, error) {
	id, err := uuid.Parse(req.Accrual.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}
	acc, err := h.svc.Update(ctx, domain.Accrual{
		ID:          id,
		Description: &req.Accrual.Description,
		Amount:      moneyToString(req.Accrual.Amount),
		AccrualDate: req.Accrual.AccrualDate.AsTime(),
		AccountID:   req.Accrual.AccountId,
		UpdatedBy:   getUserFromContext(ctx),
	})
	if err != nil {
		return nil, err
	}
	return toPbAccrual(acc), nil
}

func (h *AccrualHandler) DeleteAccrual(ctx context.Context, req *pb.DeleteAccrualRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *AccrualHandler) ListAccruals(ctx context.Context, req *pb.ListAccrualsRequest) (*pb.ListAccrualsResponse, error) {
	// default page size
	limit := int32(50)
	if req.GetPage().GetPageSize() > 0 {
		limit = req.Page.PageSize
	}

	// decode page_token into offset
	var offset int32
	if req.GetPage().GetPageToken() != "" {
		o, err := strconv.Atoi(req.Page.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page_token")
		}
		offset = int32(o)
	}

	// query service
	accruals, err := h.svc.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// map to proto
	pbAccruals := make([]*pb.Accrual, 0, len(accruals))
	for _, a := range accruals {
		pbAccruals = append(pbAccruals, toPbAccrual(a))
	}

	totalCount := int64(len(accruals)) // TODO: ideally from DB

	// compute next_page_token
	nextToken := ""
	if int64(offset)+int64(limit) < totalCount {
		nextToken = strconv.Itoa(int(offset + limit))
	}

	return &pb.ListAccrualsResponse{
		Accruals: pbAccruals,
		Page: &pb.PageResponse{
			NextPageToken: nextToken,
			TotalSize:     totalCount,
		},
	}, nil
}

// helper mapper
func toPbAccrual(a domain.Accrual) *pb.Accrual {
	return &pb.Accrual{
		Id:          a.ID.String(),
		Description: derefString(a.Description),
		Amount:      stringToMoney(a.Amount, "USD"), // convert float64 → *money.Money
		AccrualDate: timestamppb.New(a.AccrualDate),
		AccountId:   a.AccountID,
	}
}

func floatToMoney(amount float64, currency string) *money.Money {
	units := int64(amount)
	nanos := int32((amount - float64(units)) * 1e9)
	return &money.Money{
		CurrencyCode: currency,
		Units:        units,
		Nanos:        nanos,
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringPtr(s string) *string {
	return &s
}

func stringToMoney(amountStr string, currency string) *money.Money {
	f, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil
	}
	return floatToMoney(f, currency)
}

// *money.Money → string
func moneyToString(m *money.Money) string {
	if m == nil {
		return "0"
	}
	return fmt.Sprintf("%d.%09d", m.Units, m.Nanos) // crude formatting
}

func getUserFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get("user-id")
	if len(values) > 0 {
		return values[0]
	}
	return ""
}
