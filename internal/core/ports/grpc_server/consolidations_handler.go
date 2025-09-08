package grpc_server

import (
	"context"
	"strconv"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConsolidationHandler struct {
	service *services.ConsolidationService
	pb.UnimplementedConsolidationServiceServer
}

func NewConsolidationHandler(svc *services.ConsolidationService) *ConsolidationHandler {
	return &ConsolidationHandler{service: svc}
}

func (h *ConsolidationHandler) CreateConsolidation(ctx context.Context, req *pb.CreateConsolidationRequest) (*pb.Consolidation, error) {
	period := req.GetConsolidation().GetPeriod()
c := db.Consolidation{
    EntityIds:   req.GetConsolidation().GetEntityIds(),
    PeriodStart: period.GetStartDate().AsTime(),
    PeriodEnd:   period.GetEndDate().AsTime(),
    Report:      req.GetConsolidation().GetReport(),
}

	created, err := h.service.Create(ctx, c)
	if err != nil {
		return nil, err
	}

	return &pb.Consolidation{
		Id:        created.ID.String(),
		EntityIds: created.EntityIds,
		Period: &pb.ReportPeriod{
			StartDate: timestamppb.New(created.PeriodStart),
			EndDate:   timestamppb.New(created.PeriodEnd),
		},
		Report: created.Report,
	}, nil
}

func (h *ConsolidationHandler) GetConsolidation(ctx context.Context, req *pb.GetConsolidationRequest) (*pb.Consolidation, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	c, err := h.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.Consolidation{
		Id:        c.ID.String(),
		EntityIds: c.EntityIds,
		Period: &pb.ReportPeriod{
			StartDate: timestamppb.New(c.PeriodStart),
			EndDate:   timestamppb.New(c.PeriodEnd),
		},
		Report: c.Report,
	}, nil
}

func (h *ConsolidationHandler) ListConsolidations(ctx context.Context, req *pb.ListConsolidationsRequest) (*pb.ListConsolidationsResponse, error) {

    // Pagination
    limit := req.GetPage().GetPageSize()
    var offset int32
    if req.GetPage().GetPageToken() != "" {
        o, err := strconv.Atoi(req.GetPage().GetPageToken())
        if err != nil {
            return nil, err
        }
        offset = int32(o)
    }

    // Call service
    items, err := h.service.List(ctx, req.GetEntityIds(), req.GetPeriod().GetStartDate().AsTime(), req.GetPeriod().GetEndDate().AsTime(), limit, offset)
    if err != nil {
        return nil, err
    }

    // Map to protobuf
    pbItems := make([]*pb.Consolidation, len(items))
    for i, c := range items {
        pbItems[i] = &pb.Consolidation{
            Id:        c.ID.String(),
            EntityIds: c.EntityIds,
            Period: &pb.ReportPeriod{
                StartDate: timestamppb.New(c.PeriodStart),
                EndDate:   timestamppb.New(c.PeriodEnd),
            },
            Report: c.Report,
        }
    }

    return &pb.ListConsolidationsResponse{
		Consolidations: pbItems,
		Page: &pb.PageResponse{
			NextPageToken: strconv.Itoa(int(offset + limit)),
			TotalSize:     int64(len(pbItems)),  
		},
	}, nil
}

func (h *ConsolidationHandler) DeleteConsolidation(ctx context.Context, req *pb.DeleteConsolidationRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := h.service.Delete(ctx, id); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
