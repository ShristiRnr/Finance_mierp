package services_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

// =====================
// Mocks
// =====================
type MockBankRepo struct{ mock.Mock }
type MockPaymentRepo struct{ mock.Mock }
type MockTxnRepo struct{ mock.Mock }
type MockpPublisher struct{ mock.Mock }

func (m *MockBankRepo) CreateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error) {
	args := m.Called(ctx, ba)
	return args.Get(0).(db.BankAccount), args.Error(1)
}
func (m *MockBankRepo) UpdateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error) {
	args := m.Called(ctx, ba)
	return args.Get(0).(db.BankAccount), args.Error(1)
}
func (m *MockBankRepo) DeleteBankAccount(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockBankRepo) GetBankAccount(ctx context.Context, id uuid.UUID) (db.BankAccount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.BankAccount), args.Error(1)
}
func (m *MockBankRepo) ListBankAccounts(ctx context.Context, limit, offset int32) ([]db.BankAccount, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.BankAccount), args.Error(1)
}

// PaymentDue Repo
func (m *MockPaymentRepo) CreatePaymentDue(ctx context.Context, pd db.PaymentDue) (db.PaymentDue, error) {
	args := m.Called(ctx, pd)
	return args.Get(0).(db.PaymentDue), args.Error(1)
}
func (m *MockPaymentRepo) UpdatePaymentDue(ctx context.Context, pd db.PaymentDue) (db.PaymentDue, error) {
	args := m.Called(ctx, pd)
	return args.Get(0).(db.PaymentDue), args.Error(1)
}
func (m *MockPaymentRepo) DeletePaymentDue(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockPaymentRepo) MarkPaymentAsPaid(ctx context.Context, id uuid.UUID, updatedBy string) (db.PaymentDue, error) {
	args := m.Called(ctx, id, updatedBy)
	return args.Get(0).(db.PaymentDue), args.Error(1)
}
func (m *MockPaymentRepo) GetPaymentDue(ctx context.Context, id uuid.UUID) (db.PaymentDue, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.PaymentDue), args.Error(1)
}
func (m *MockPaymentRepo) ListPaymentDues(ctx context.Context, limit, offset int32) ([]db.PaymentDue, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.PaymentDue), args.Error(1)
}

// BankTransaction Repo
func (m *MockTxnRepo) ImportBankTransaction(ctx context.Context, tx db.BankTransaction) (db.BankTransaction, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).(db.BankTransaction), args.Error(1)
}
func (m *MockTxnRepo) ReconcileTransaction(ctx context.Context, tx db.BankTransaction) (db.BankTransaction, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).(db.BankTransaction), args.Error(1)
}
func (m *MockTxnRepo) ListBankTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int32) ([]db.BankTransaction, error) {
	args := m.Called(ctx, bankAccountID, limit, offset)
	return args.Get(0).([]db.BankTransaction), args.Error(1)
}

// Event Publisher
func (m *MockpPublisher) PublishAccountCreated(ctx context.Context, acc *db.Account) error {
	return m.Called(ctx, acc).Error(0)
}
func (m *MockpPublisher) PublishAccountUpdated(ctx context.Context, acc *db.Account) error {
	return m.Called(ctx, acc).Error(0)
}
func (m *MockpPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockpPublisher) PublishAuditRecorded(ctx context.Context, e *db.AuditEvent) error {
	return m.Called(ctx, e).Error(0)
}

// =====================
// Tests
// =====================
func TestBankService_CreateBankAccount(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBankRepo)
	mockPub := new(MockPublisher)
	svc := services.NewBankService(mockRepo, mockPub)

	acc := db.BankAccount{ID: uuid.New(), Name: "HDFC"}
	mockRepo.On("CreateBankAccount", ctx, acc).Return(acc, nil)
	mockPub.On("PublishAccountCreated", ctx, mock.Anything).Return(nil)

	res, err := svc.CreateBankAccount(ctx, acc)
	assert.NoError(t, err)
	assert.Equal(t, "HDFC", res.Name)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestPaymentDueService_CreatePaymentDue(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPaymentRepo)
	mockPub := new(MockPublisher)
	svc := services.NewPaymentDueService(mockRepo, mockPub)

	pd := db.PaymentDue{ID: uuid.New(), Status: "PENDING", CreatedBy: sql.NullString{String: "user1", Valid: true}}
	mockRepo.On("CreatePaymentDue", ctx, pd).Return(pd, nil)
	mockPub.On("PublishAuditRecorded", ctx, mock.Anything).Return(nil)

	res, err := svc.CreatePaymentDue(ctx, pd)
	assert.NoError(t, err)
	assert.Equal(t, "PENDING", res.Status)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestPaymentDueService_MarkPaymentAsPaid(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPaymentRepo)
	mockPub := new(MockPublisher)
	svc := services.NewPaymentDueService(mockRepo, mockPub)

	id := uuid.New()
	pd := db.PaymentDue{ID: id, Status: "PAID"}
	mockRepo.On("MarkPaymentAsPaid", ctx, id, "admin").Return(pd, nil)
	mockPub.On("PublishAuditRecorded", ctx, mock.Anything).Return(nil)

	res, err := svc.MarkPaymentAsPaid(ctx, id, "admin")
	assert.NoError(t, err)
	assert.Equal(t, "PAID", res.Status)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestBankTransactionService_ImportBankTransaction(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTxnRepo)
	mockPub := new(MockPublisher)
	svc := services.NewBankTransactionService(mockRepo, mockPub)

	tx := db.BankTransaction{ID: uuid.New(), Amount: "1000"}
	mockRepo.On("ImportBankTransaction", ctx, tx).Return(tx, nil)
	mockPub.On("PublishAuditRecorded", ctx, mock.Anything).Return(nil)

	res, err := svc.ImportBankTransaction(ctx, tx, "admin")
	assert.NoError(t, err)
	assert.Equal(t, "1000", res.Amount)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestBankTransactionService_ReconcileTransaction_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTxnRepo)
	mockPub := new(MockPublisher)
	svc := services.NewBankTransactionService(mockRepo, mockPub)

	tx := db.BankTransaction{ID: uuid.New()}
	mockRepo.On("ReconcileTransaction", ctx, tx).Return(db.BankTransaction{}, errors.New("db error"))

	_, err := svc.ReconcileTransaction(ctx, tx, "admin")
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBankService_UpdateBankAccount(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBankRepo)
	mockPub := new(MockPublisher)
	svc := services.NewBankService(mockRepo, mockPub)

	acc := db.BankAccount{ID: uuid.New(), Name: "ICICI"}
	mockRepo.On("UpdateBankAccount", ctx, acc).Return(acc, nil)
	mockPub.On("PublishAccountUpdated", ctx, mock.Anything).Return(nil)

	res, err := svc.UpdateBankAccount(ctx, acc)
	assert.NoError(t, err)
	assert.Equal(t, "ICICI", res.Name)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestBankService_GetBankAccount(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBankRepo)
	mockPub := new(MockPublisher)
	svc := services.NewBankService(mockRepo, mockPub)

	id := uuid.New()
	acc := db.BankAccount{ID: id, Name: "SBI"}
	mockRepo.On("GetBankAccount", ctx, id).Return(acc, nil)

	res, err := svc.GetBankAccount(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "SBI", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestBankService_DeleteBankAccount(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBankRepo)
	mockPub := new(MockPublisher)
	svc := services.NewBankService(mockRepo, mockPub)

	id := uuid.New()
	mockRepo.On("DeleteBankAccount", ctx, id).Return(nil)
	mockPub.On("PublishAccountDeleted", ctx, id.String()).Return(nil)

	err := svc.DeleteBankAccount(ctx, id)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestBankService_ListBankAccounts(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBankRepo)
	mockPub := new(MockPublisher)
	svc := services.NewBankService(mockRepo, mockPub)

	accs := []db.BankAccount{
		{ID: uuid.New(), Name: "Axis"},
		{ID: uuid.New(), Name: "Kotak"},
	}
	mockRepo.On("ListBankAccounts", ctx, int32(10), int32(0)).Return(accs, nil)

	res, err := svc.ListBankAccounts(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Axis", res[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestPaymentDueService_UpdatePaymentDue(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPaymentRepo)
	mockPub := new(MockPublisher)
	svc := services.NewPaymentDueService(mockRepo, mockPub)

	pd := db.PaymentDue{ID: uuid.New(), Status: "UPDATED"}
	mockRepo.On("UpdatePaymentDue", ctx, pd).Return(pd, nil)
	mockPub.On("PublishAuditRecorded", ctx, mock.Anything).Return(nil)

	res, err := svc.UpdatePaymentDue(ctx, pd)
	assert.NoError(t, err)
	assert.Equal(t, "UPDATED", res.Status)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestPaymentDueService_GetPaymentDue(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPaymentRepo)
	mockPub := new(MockPublisher)
	svc := services.NewPaymentDueService(mockRepo, mockPub)

	id := uuid.New()
	pd := db.PaymentDue{ID: id, Status: "PENDING"}
	mockRepo.On("GetPaymentDue", ctx, id).Return(pd, nil)

	res, err := svc.GetPaymentDue(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "PENDING", res.Status)
	mockRepo.AssertExpectations(t)
}

func TestPaymentDueService_DeletePaymentDue(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPaymentRepo)
	mockPub := new(MockPublisher)
	svc := services.NewPaymentDueService(mockRepo, mockPub)

	id := uuid.New()
	updatedBy:= "admin"
	mockRepo.On("DeletePaymentDue", ctx, id).Return(nil)
	mockPub.On("PublishAuditRecorded", ctx, mock.Anything).Return(nil)

	err := svc.DeletePaymentDue(ctx, id, updatedBy)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestPaymentDueService_ListPaymentDues(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPaymentRepo)
	mockPub := new(MockPublisher)
	svc := services.NewPaymentDueService(mockRepo, mockPub)

	pds := []db.PaymentDue{
		{ID: uuid.New(), Status: "PENDING"},
		{ID: uuid.New(), Status: "PAID"},
	}
	mockRepo.On("ListPaymentDues", ctx, int32(5), int32(0)).Return(pds, nil)

	res, err := svc.ListPaymentDues(ctx, 5, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	mockRepo.AssertExpectations(t)
}
