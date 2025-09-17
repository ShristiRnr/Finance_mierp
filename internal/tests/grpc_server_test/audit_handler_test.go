package grpc_server_test

import (
	"context"
	"testing"
	"time"
	"database/sql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/google/uuid"

	grpcserver "github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
)

// --- Mocks ---

type MockAuditService struct {
	mock.Mock
}

func (m *MockAuditService) Record(ctx context.Context, event *db.AuditEvent) (*db.AuditEvent, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(*db.AuditEvent), args.Error(1)
}

func (m *MockAuditService) List(ctx context.Context, page db.Pagination) ([]db.AuditEvent, error) {
	args := m.Called(ctx, page)
	return args.Get(0).([]db.AuditEvent), args.Error(1)
}

func (m *MockAuditService) GetByID(ctx context.Context, id string) (*db.AuditEvent, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.AuditEvent), args.Error(1)
}

func (m *MockAuditService) Filter(ctx context.Context, filter db.FilterParams, page db.Pagination) ([]db.AuditEvent, error) {
	args := m.Called(ctx, filter, page)
	return args.Get(0).([]db.AuditEvent), args.Error(1)
}

// --- Tests ---

func TestRecordAuditEvent(t *testing.T) {
	mockSvc := new(MockAuditService)
	handler := grpcserver.NewAuditHandler(mockSvc, nil)

	inputEvent := &pb.AuditEvent{
		UserId:       "user1",
		Action:       "create",
		Timestamp:    timestamppb.New(time.Now()),
		Details:      "details",
		ResourceType: "type",
		ResourceId:   "res1",
	}

	domainEvent := &db.AuditEvent{
		UserID:       inputEvent.UserId,
		Action:       inputEvent.Action,
		Timestamp:    inputEvent.Timestamp.AsTime(),
		Details:      sql.NullString{String: inputEvent.Details, Valid: true},
		ResourceType: sql.NullString{String: inputEvent.ResourceType, Valid: true},
		ResourceID:   sql.NullString{String: inputEvent.ResourceId, Valid: true},
	}

	mockSvc.On("Record", mock.Anything, domainEvent).Return(domainEvent, nil)

	resp, err := handler.RecordAuditEvent(context.Background(), &pb.RecordAuditEventRequest{Event: inputEvent})
	assert.NoError(t, err)
	assert.Equal(t, inputEvent.UserId, resp.UserId)
	assert.Equal(t, inputEvent.Action, resp.Action)
}

func TestGetAuditEventById(t *testing.T) {
	mockSvc := new(MockAuditService)
	handler := grpcserver.NewAuditHandler(mockSvc, nil)

	eventID := "123"
	domainEvent := &db.AuditEvent{
		ID:     uuid.New(),
		UserID: "user1",
		Action: "update",
	}

	mockSvc.On("GetByID", mock.Anything, eventID).Return(domainEvent, nil)

	resp, err := handler.GetAuditEventById(context.Background(), &pb.GetAuditEventByIdRequest{Id: eventID})
	assert.NoError(t, err)
	assert.Equal(t, domainEvent.UserID, resp.UserId)
}

func TestListAuditEvents(t *testing.T) {
	mockSvc := new(MockAuditService)
	handler := grpcserver.NewAuditHandler(mockSvc, nil)

	events := []db.AuditEvent{
		{ID: uuid.New(), UserID: "user1", Action: "create"},
		{ID: uuid.New(), UserID: "user2", Action: "delete"},
	}

	mockSvc.On("List", mock.Anything, db.Pagination{Offset: 0, Limit: 50}).Return(events, nil)

	resp, err := handler.ListAuditEvents(context.Background(), &pb.ListAuditEventsRequest{
		Page: &pb.PageRequest{PageSize: 50},
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Events, len(events))
	assert.Equal(t, events[0].UserID, resp.Events[0].UserId)
}

func TestFilterAuditEvents(t *testing.T) {
	mockSvc := new(MockAuditService)
	handler := grpcserver.NewAuditHandler(mockSvc, nil)

	filter := db.FilterParams{
		UserID:       ptrString("user1"),
		ResourceType: ptrString("type"),
	}
	events := []db.AuditEvent{
		{ID: uuid.New(), UserID: "user1", Action: "update"},
	}

	mockSvc.On("Filter", mock.Anything, filter, db.Pagination{Offset: 0, Limit: 50}).Return(events, nil)

	resp, err := handler.FilterAuditEvents(context.Background(), &pb.FilterAuditEventsRequest{
		UserId:       "user1",
		ResourceType: "type",
		Page:         &pb.PageRequest{PageSize: 50},
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Events, len(events))
	assert.Equal(t, events[0].UserID, resp.Events[0].UserId)
}

// --- Helpers ---

func ptrString(s string) *string {
	return &s
}
