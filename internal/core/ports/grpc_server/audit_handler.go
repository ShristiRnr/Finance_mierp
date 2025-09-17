package grpc_server

import (
	"context"
	"strconv"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the gRPC server for the AuditTrailService.
type AuditHandler struct {
	pb.UnimplementedAuditTrailServiceServer
	service ports.AuditService
	producer ports.EventPublisher
}

func NewAuditHandler(service ports.AuditService,producer ports.EventPublisher) *AuditHandler {
	return &AuditHandler{service: service,
	producer: producer,}
}

// Helper to convert from domain model to proto message
func toProto(e *db.AuditEvent) *pb.AuditEvent {
	return &pb.AuditEvent{
		Id:           e.ID.String(),
		UserId:       e.UserID,
		Action:       e.Action,
		Timestamp:    timestamppb.New(e.Timestamp),
		Details:      e.Details.String,
		ResourceType: e.ResourceType.String,
		ResourceId:   e.ResourceID.String,
	}
}

func (s *AuditHandler) RecordAuditEvent(ctx context.Context, req *pb.RecordAuditEventRequest) (*pb.AuditEvent, error) {
	event := req.GetEvent()
	domainEvent := &db.AuditEvent{
		UserID:       event.GetUserId(),
		Action:       event.GetAction(),
		Timestamp:    event.GetTimestamp().AsTime(),
		Details:      toNullString(event.GetDetails()),
		ResourceType: toNullString(event.GetResourceType()),
		ResourceID:   toNullString(event.GetResourceId()),
	}

	// Save + publish (delegated to service)
	recordedEvent, err := s.service.Record(ctx, domainEvent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to record audit: %v", err)
	}

	return toProto(recordedEvent), nil
}


func (s *AuditHandler) GetAuditEventById(ctx context.Context, req *pb.GetAuditEventByIdRequest) (*pb.AuditEvent, error) {
	event, err := s.service.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return toProto(event), nil
}

func (s *AuditHandler) ListAuditEvents(ctx context.Context, req *pb.ListAuditEventsRequest) (*pb.ListAuditEventsResponse, error) {
	// Default page size if not provided
	limit := int32(50)
	if req.GetPage().GetPageSize() > 0 {
		limit = req.Page.PageSize
	}

	// Decode page_token into offset
	var offset int32
	if req.GetPage().GetPageToken() != "" {
		o, err := strconv.Atoi(req.Page.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page_token")
		}
		offset = int32(o)
	}

	// Get events + total count from service
	events, err := s.service.List(ctx, db.Pagination{ Offset: int(offset), Limit: int(limit) })
	if err != nil {
		return nil, err
	}

	// Convert to proto
	protoEvents := make([]*pb.AuditEvent, 0, len(events))
	for _, e := range events {
		event := e
		protoEvents = append(protoEvents, toProto(&event))
	}

	// Compute next page token

	totalCount := int64(len(events))
	var nextToken string
	if offset+limit < int32(totalCount) {
		nextToken = strconv.Itoa(int(offset + limit))
	}

	return &pb.ListAuditEventsResponse{
		Events: protoEvents,
		Page: &pb.PageResponse{
			NextPageToken: nextToken,
			TotalSize:     totalCount, // cast because COUNT() is int64
		},
	}, nil
}



func (s *AuditHandler) FilterAuditEvents(ctx context.Context, req *pb.FilterAuditEventsRequest) (*pb.FilterAuditEventsResponse, error) {
	filter := db.FilterParams{}
	if req.UserId != "" {
		filter.UserID = &req.UserId
	}
	if req.Action != "" {
		filter.Action = &req.Action
	}
	if req.ResourceType != "" {
		filter.ResourceType = &req.ResourceType
	}
	if req.ResourceId != "" {
		filter.ResourceID = &req.ResourceId
	}
	if req.FromDate.IsValid() {
		fromDate := req.FromDate.AsTime()
		filter.FromDate = &fromDate
	}
	if req.ToDate.IsValid() {
		toDate := req.ToDate.AsTime()
		filter.ToDate = &toDate
	}

	limit := int32(50)
	if req.GetPage().GetPageSize() > 0 {
		limit = req.Page.PageSize
	}

	// Decode page_token into offset
	var offset int32
	if req.GetPage().GetPageToken() != "" {
		o, err := strconv.Atoi(req.Page.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page_token")
		}
		offset = int32(o)
	}

	// Call the service
	events, err := s.service.Filter(ctx, filter, db.Pagination{Limit: int(limit), Offset: int(offset)})
	if err != nil {
		return nil, err
	}

	// Convert domain events to proto
	protoEvents := make([]*pb.AuditEvent, len(events))
	for i, e := range events {
		protoEvents[i] = toProto(&e)
	}

	totalCount := int64(len(events))
	var nextToken string
	if int64(offset+limit) < totalCount {
		nextToken = strconv.Itoa(int(offset + limit))
	}

	return &pb.FilterAuditEventsResponse{
		Events: protoEvents,
		Page: &pb.PageResponse{
			NextPageToken: nextToken,
			TotalSize:     totalCount,
		},
	}, nil
}