package services_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Repository ---
type MockCreditDebitNoteRepo struct {
	mock.Mock
}

func (m *MockCreditDebitNoteRepo) Create(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error) {
	args := m.Called(ctx, note)
	return args.Get(0).(db.CreditDebitNote), args.Error(1)
}

func (m *MockCreditDebitNoteRepo) Get(ctx context.Context, id uuid.UUID) (db.CreditDebitNote, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.CreditDebitNote), args.Error(1)
}

func (m *MockCreditDebitNoteRepo) List(ctx context.Context, limit, offset int32) ([]db.CreditDebitNote, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.CreditDebitNote), args.Error(1)
}

func (m *MockCreditDebitNoteRepo) Update(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error) {
	args := m.Called(ctx, note)
	return args.Get(0).(db.CreditDebitNote), args.Error(1)
}

func (m *MockCreditDebitNoteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Mock Publisher ---
type MockePublisher struct {
	mock.Mock
}

func (m *MockePublisher) PublishCreditDebitNoteCreated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}
func (m *MockePublisher) PublishCreditDebitNoteUpdated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}
func (m *MockePublisher) PublishCreditDebitNoteDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Helper ---
func makeTestNote() db.CreditDebitNote {
	return db.CreditDebitNote{
		ID:        uuid.New(),
		InvoiceID: uuid.New(),
		Type:      "CREDIT",
		Amount:    "500.00",
		Reason:    sql.NullString{String: "test", Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true}, 
		CreatedBy: sql.NullString{String: "tester", Valid: true},
	}
}

func TestCreditDebitNoteService_Create_Success(t *testing.T) {
	mockRepo := new(MockCreditDebitNoteRepo)
	mockPub := new(MockPublisher)
	svc := services.NewCreditDebitNoteService(mockRepo, mockPub)

	input := makeTestNote()
	expected := input

	mockRepo.On("Create", mock.Anything, input).Return(expected, nil).Once()
	mockPub.On("PublishCreditDebitNoteCreated", mock.Anything, &expected).Return(nil).Once()

	got, err := svc.Create(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, got.ID)

	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestCreditDebitNoteService_Create_RepoError(t *testing.T) {
	mockRepo := new(MockCreditDebitNoteRepo)
	mockPub := new(MockPublisher)
	svc := services.NewCreditDebitNoteService(mockRepo, mockPub)

	input := makeTestNote()

	mockRepo.On("Create", mock.Anything, input).Return(db.CreditDebitNote{}, errors.New("db error")).Once()

	_, err := svc.Create(context.Background(), input)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCreditDebitNoteService_Get(t *testing.T) {
	mockRepo := new(MockCreditDebitNoteRepo)
	mockPub := new(MockPublisher)
	svc := services.NewCreditDebitNoteService(mockRepo, mockPub)

	note := makeTestNote()
	mockRepo.On("Get", mock.Anything, note.ID).Return(note, nil).Once()

	got, err := svc.Get(context.Background(), note.ID)
	assert.NoError(t, err)
	assert.Equal(t, note.ID, got.ID)

	mockRepo.AssertExpectations(t)
}

func TestCreditDebitNoteService_List(t *testing.T) {
	mockRepo := new(MockCreditDebitNoteRepo)
	mockPub := new(MockPublisher)
	svc := services.NewCreditDebitNoteService(mockRepo, mockPub)

	notes := []db.CreditDebitNote{makeTestNote(), makeTestNote()}
	mockRepo.On("List", mock.Anything, int32(10), int32(0)).Return(notes, nil).Once()

	got, err := svc.List(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, got, 2)

	mockRepo.AssertExpectations(t)
}

func TestCreditDebitNoteService_Update(t *testing.T) {
	mockRepo := new(MockCreditDebitNoteRepo)
	mockPub := new(MockPublisher)
	svc := services.NewCreditDebitNoteService(mockRepo, mockPub)

	input := makeTestNote()
	updated := input
	updated.Reason = sql.NullString{String: "updated reason", Valid: true}

	mockRepo.On("Update", mock.Anything, input).Return(updated, nil).Once()
	mockPub.On("PublishCreditDebitNoteUpdated", mock.Anything, &updated).Return(nil).Once()

	got, err := svc.Update(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, "updated reason", got.Reason.String)

	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestCreditDebitNoteService_Delete(t *testing.T) {
	mockRepo := new(MockCreditDebitNoteRepo)
	mockPub := new(MockPublisher)
	svc := services.NewCreditDebitNoteService(mockRepo, mockPub)

	id := uuid.New()

	mockRepo.On("Delete", mock.Anything, id).Return(nil).Once()
	mockPub.On("PublishCreditDebitNoteDeleted", mock.Anything, id.String()).Return(nil).Once()

	err := svc.Delete(context.Background(), id)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}
