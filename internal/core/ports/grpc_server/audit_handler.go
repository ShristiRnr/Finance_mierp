package grpc_server

import (
	"context"
	"strconv"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the gRPC server for the AuditTrailService.
type AuditHandler struct {
	pb.UnimplementedAuditTrailServiceServer
	service *services.AuditService
}

func NewAuditHandler(service *services.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

// Helper to convert from domain model to proto message
func toProto(e *domain.AuditEvent) *pb.AuditEvent {
	return &pb.AuditEvent{
		Id:           e.ID,
		UserId:       e.UserID,
		Action:       e.Action,
		Timestamp:    timestamppb.New(e.Timestamp),
		Details:      e.Details,
		ResourceType: e.ResourceType,
		ResourceId:   e.ResourceID,
	}
}

func (s *AuditHandler) RecordAuditEvent(ctx context.Context, req *pb.RecordAuditEventRequest) (*pb.AuditEvent, error) {
	event := req.GetEvent()
	domainEvent := &domain.AuditEvent{
		UserID:       event.GetUserId(),
		Action:       event.GetAction(),
		Timestamp:    event.GetTimestamp().AsTime(),
		Details:      event.GetDetails(),
		ResourceType: event.GetResourceType(),
		ResourceID:   event.GetResourceId(),
	}

	recordedEvent, err := s.service.Record(ctx, domainEvent)
	if err != nil {
		return nil, err // Convert to gRPC error status
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
	events, err := s.service.List(ctx, domain.Pagination{Offset: offset, Limit: limit})
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
	filter := domain.FilterParams{}
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
	events, err := s.service.Filter(ctx, filter, domain.Pagination{Offset: offset, Limit: limit})
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