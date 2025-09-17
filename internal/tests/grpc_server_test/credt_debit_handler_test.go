package grpc_server_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"database/sql"
	
	"github.com/stretchr/testify/mock"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	money "google.golang.org/genproto/googleapis/type/money"
)

// --- Mock Service ---
type MockCreditDebitNoteService struct {
	mock.Mock
}

func (m *MockCreditDebitNoteService) Create(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error) {
	args := m.Called(ctx, note)
	return args.Get(0).(db.CreditDebitNote), args.Error(1)
}
func (m *MockCreditDebitNoteService) Get(ctx context.Context, id uuid.UUID) (db.CreditDebitNote, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.CreditDebitNote), args.Error(1)
}
func (m *MockCreditDebitNoteService) List(ctx context.Context, limit, offset int32) ([]db.CreditDebitNote, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.CreditDebitNote), args.Error(1)
}
func (m *MockCreditDebitNoteService) Update(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error) {
	args := m.Called(ctx, note)
	return args.Get(0).(db.CreditDebitNote), args.Error(1)
}
func (m *MockCreditDebitNoteService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func makeTestNote() db.CreditDebitNote {
	now := time.Now()
	return db.CreditDebitNote{
		ID:        uuid.New(),
		InvoiceID: uuid.New(),
		Type:      "credit",
		Amount:    "500.00",
		Reason:    toNullString("discount"),
		CreatedBy: toNullString("tester"),
		CreatedAt: sql.NullTime{Time: now, Valid: true},
	}
}

// --- Tests ---

func TestCreateCreditDebitNote_Success(t *testing.T) {
	mockSvc := new(MockCreditDebitNoteService)
	handler := grpc_server.NewGRPCServer(mockSvc)

	req := &pb.CreateCreditDebitNoteRequest{
		Note: &pb.CreditDebitNote{
			InvoiceId: uuid.New().String(),
			Type:      pb.NoteType_NOTE_TYPE_CREDIT,
			Amount:    &money.Money{CurrencyCode: "INR", Units: 500},
			Audit:     &pb.AuditFields{CreatedBy: "tester"},
		},
	}

	expected := makeTestNote()
	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("db.CreditDebitNote")).
		Return(expected, nil).Once()

	resp, err := handler.CreateCreditDebitNote(context.Background(), req)
	assert.NoError(t, err)

	// ✅ Convert resp.Amount to string before comparing
	respAmount := fmt.Sprintf("%.2f", float64(resp.Amount.Units)+float64(resp.Amount.Nanos)/1e9)
	assert.Equal(t, expected.Amount, respAmount)

	mockSvc.AssertExpectations(t)
}


func TestCreateCreditDebitNote_InvalidInvoiceID(t *testing.T) {
	mockSvc := new(MockCreditDebitNoteService)
	handler := grpc_server.NewGRPCServer(mockSvc)

	req := &pb.CreateCreditDebitNoteRequest{
		Note: &pb.CreditDebitNote{
			InvoiceId: "not-a-uuid",
		},
	}

	resp, err := handler.CreateCreditDebitNote(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestGetCreditDebitNote_Success(t *testing.T) {
	mockSvc := new(MockCreditDebitNoteService)
	handler := grpc_server.NewGRPCServer(mockSvc)

	note := makeTestNote()
	mockSvc.On("Get", mock.Anything, note.ID).Return(note, nil)

	resp, err := handler.GetCreditDebitNote(context.Background(), &pb.GetCreditDebitNoteRequest{Id: note.ID.String()})
	assert.NoError(t, err)
	assert.Equal(t, note.ID.String(), resp.Id)
}

func TestGetCreditDebitNote_NotFound(t *testing.T) {
	mockSvc := new(MockCreditDebitNoteService)
	handler := grpc_server.NewGRPCServer(mockSvc)

	id := uuid.New()
	mockSvc.On("Get", mock.Anything, id).Return(db.CreditDebitNote{}, errors.New("not found"))

	resp, err := handler.GetCreditDebitNote(context.Background(), &pb.GetCreditDebitNoteRequest{Id: id.String()})
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestListCreditDebitNotes_Success(t *testing.T) {
	mockSvc := new(MockCreditDebitNoteService)
	handler := grpc_server.NewGRPCServer(mockSvc)

	notes := []db.CreditDebitNote{makeTestNote(), makeTestNote()}
	mockSvc.On("List", mock.Anything, int32(50), int32(0)).Return(notes, nil)

	resp, err := handler.ListCreditDebitNotes(context.Background(), &pb.ListCreditDebitNotesRequest{})
	assert.NoError(t, err)
	assert.Len(t, resp.Notes, 2)
}

func TestUpdateCreditDebitNote_Success(t *testing.T) {
	mockSvc := new(MockCreditDebitNoteService)
	handler := grpc_server.NewGRPCServer(mockSvc)

	note := makeTestNote()
	mockSvc.On("Update", mock.Anything, mock.AnythingOfType("db.CreditDebitNote")).Return(note, nil)

	req := &pb.UpdateCreditDebitNoteRequest{
		Note: &pb.CreditDebitNote{
			Id:        note.ID.String(),
			InvoiceId: note.InvoiceID.String(),
			Type:       pb.NoteType_NOTE_TYPE_CREDIT,
			Amount:    &money.Money{CurrencyCode: "INR", Units: 500},
			Reason:    "discount",
			Audit:     &pb.AuditFields{UpdatedBy: "tester"},
		},
	}

	resp, err := handler.UpdateCreditDebitNote(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, note.ID.String(), resp.Id)
}

func TestDeleteCreditDebitNote_Success(t *testing.T) {
	mockSvc := new(MockCreditDebitNoteService)
	handler := grpc_server.NewGRPCServer(mockSvc)

	id := uuid.New()
	mockSvc.On("Delete", mock.Anything, id).Return(nil)

	resp, err := handler.DeleteCreditDebitNote(context.Background(), &pb.DeleteCreditDebitNoteRequest{Id: id.String()})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func toNullString(s string) sql.NullString {
    if s == "" {
        return sql.NullString{Valid: false}
    }
    return sql.NullString{String: s, Valid: true}
}