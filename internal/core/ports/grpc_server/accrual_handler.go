package grpc_server

import (
	"context"
	"strconv"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
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
	acc, err := h.svc.Create(ctx, domain.Accrual{
		ID:          uuid.New(),
		Description: &req.Description,
		Amount:      req.Amount,
		AccrualDate: req.AccrualDate.AsTime(),
		AccountID:   req.AccountId,
		CreatedBy:   req.UserId,
		UpdatedBy:   req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return toPbAccrual(acc), nil
}

func (h *AccrualHandler) GetAccrual(ctx context.Context, req *pb.GetAccrualRequest) (*pb.Accrual, error) {
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
		Amount:      req.Accrual.Amount,
		AccrualDate: req.Accrual.AccrualDate.AsTime(),
		AccountID:   req.Accrual.AccountId,
		UpdatedBy:   req.UserId,
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
		Amount:      a.Amount,
		AccrualDate: timestamppb.New(a.AccrualDate),
		AccountId:   a.AccountID,
	}
}
