package grpc_server

import (
	"context"
	"strconv"
	"database/sql"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

type AccrualHandler struct {
	pb.UnimplementedAccrualServiceServer
	svc *services.AccrualService
	publisher ports.EventPublisher
}

func NewAccrualHandler(svc *services.AccrualService,  pub ports.EventPublisher) *AccrualHandler {
	return &AccrualHandler{svc: svc, publisher: pub}
}

func (h *AccrualHandler) CreateAccrual(ctx context.Context, req *pb.CreateAccrualRequest) (*pb.Accrual, error) {
	id, err := uuid.Parse(req.Accrual.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}
	acc, err := h.svc.Create(ctx, db.Accrual{
		ID:          id,
		Description: toNullString(req.Accrual.Description),
		Amount:      moneyToString(req.Accrual.Amount),
		AccrualDate: req.Accrual.AccrualDate.AsTime(),
		AccountID:   req.Accrual.AccountId,
		UpdatedBy:   toNullString(getUserFromContext(ctx)),
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
	acc, err := h.svc.Update(ctx, db.Accrual{
		ID:          id,
		Description: toNullString(req.Accrual.Description),
		Amount:      moneyToString(req.Accrual.Amount),
		AccrualDate: req.Accrual.AccrualDate.AsTime(),
		AccountID:   req.Accrual.AccountId,
		UpdatedBy:   toNullString(getUserFromContext(ctx)),
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

type AllocationHandler struct {
	pb.UnimplementedAllocationAutomationServiceServer
	svc       ports.AllocationService
	publisher ports.EventPublisher // Kafka publisher
}

func NewAllocationHandler(svc ports.AllocationService, pub ports.EventPublisher) *AllocationHandler {
	return &AllocationHandler{svc: svc, publisher: pub}
}

func (h *AllocationHandler) CreateRule(ctx context.Context, req *pb.CreateAllocationRuleRequest) (*pb.AllocationRule, error) {
	rule, err := h.svc.CreateRule(ctx, db.AllocationRule{
		Name:                req.Rule.Name,
		Basis:               req.Rule.Basis,
		SourceAccountID:     req.Rule.SourceAccountId,
		TargetCostCenterIds: req.Rule.TargetCostCenterIds,
		Formula:             toNullString(req.Rule.Formula),
		CreatedBy:           toNullString(getUserFromContext(ctx)),
		UpdatedBy:           toNullString(getUserFromContext(ctx)),
	})
	if err != nil {
		return nil, err
	}
	return toPbAllocationRule(rule), nil
}

func (h *AllocationHandler) UpdateRule(ctx context.Context, req *pb.UpdateAllocationRuleRequest) (*pb.AllocationRule, error) {
	id, err := uuid.Parse(req.Rule.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	rule, err := h.svc.UpdateRule(ctx, db.AllocationRule{
		ID:                  id,
		Name:                req.Rule.Name,
		Basis:               req.Rule.Basis,
		SourceAccountID:     req.Rule.SourceAccountId,
		TargetCostCenterIds: req.Rule.TargetCostCenterIds,
		Formula:             toNullString(req.Rule.Formula),
		UpdatedBy:           toNullString(getUserFromContext(ctx)),
	})
	if err != nil {
		return nil, err
	}

	return toPbAllocationRule(rule), nil
}

func (h *AllocationHandler) DeleteRule(ctx context.Context, req *pb.DeleteAllocationRuleRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	if err := h.svc.DeleteRule(ctx, id); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func toPbAllocationRule(r db.AllocationRule) *pb.AllocationRule {
	return &pb.AllocationRule{
		Id:                  r.ID.String(),
		Name:                r.Name,
		Basis:               r.Basis,
		SourceAccountId:     r.SourceAccountID,
		TargetCostCenterIds: r.TargetCostCenterIds,
		Formula:             fromNullString(r.Formula),
	}
}

func fromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
