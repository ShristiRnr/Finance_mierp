package services_test

import (
	"context"
	"testing"
	"time"
	"database/sql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"

)

// ------------------ Mock Definitions ------------------
type MockGstsRepository struct {
	mock.Mock
}

func (m *MockGstsRepository) AddGstBreakup(ctx context.Context, invoiceID uuid.UUID, taxableAmount float64, cgst, sgst, igst, totalGst *float64) (db.GstBreakup, error) {
	args := m.Called(ctx, invoiceID, taxableAmount, cgst, sgst, igst, totalGst)
	return args.Get(0).(db.GstBreakup), args.Error(1)
}

func (m *MockGstsRepository) GetGstBreakup(ctx context.Context, invoiceID uuid.UUID) (db.GstBreakup, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstBreakup), args.Error(1)
}

func (m *MockGstsRepository) AddGstRegime(ctx context.Context, invoiceID uuid.UUID, gstin, placeOfSupply string, reverseCharge *bool) (db.GstRegime, error) {
	args := m.Called(ctx, invoiceID, gstin, placeOfSupply, reverseCharge)
	return args.Get(0).(db.GstRegime), args.Error(1)
}

func (m *MockGstsRepository) GetGstRegime(ctx context.Context, invoiceID uuid.UUID) (db.GstRegime, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstRegime), args.Error(1)
}

func (m *MockGstsRepository) AddGstDocStatus(ctx context.Context, invoiceID uuid.UUID, einvoiceStatus, irn, ackNo *string, ackDate *time.Time, ewayStatus, ewayBillNo *string, ewayValidUpto *time.Time, lastError *string, lastSyncedAt *time.Time) (db.GstDocStatus, error) {
	args := m.Called(ctx, invoiceID, einvoiceStatus, irn, ackNo, ackDate, ewayStatus, ewayBillNo, ewayValidUpto, lastError, lastSyncedAt)
	return args.Get(0).(db.GstDocStatus), args.Error(1)
}

func (m *MockGstsRepository) GetGstDocStatus(ctx context.Context, invoiceID uuid.UUID) (db.GstDocStatus, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstDocStatus), args.Error(1)
}

type MockaPublisher struct {
	mock.Mock
}

// Generic publish
func (m *MockaPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	args := m.Called(ctx, topic, key, payload)
	return args.Error(0)
}

func (m *MockaPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ---------------- Account events ----------------
func (m *MockaPublisher) PublishAccountCreated(ctx context.Context, a *db.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockaPublisher) PublishAccountUpdated(ctx context.Context, a *db.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockaPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Journal events ----------------
func (m *MockaPublisher) PublishJournalCreated(ctx context.Context, j *db.JournalEntry) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockaPublisher) PublishJournalUpdated(ctx context.Context, j *db.JournalEntry) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockaPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Accrual events ----------------
func (m *MockaPublisher) PublishAccrualCreated(ctx context.Context, a *db.Accrual) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockaPublisher) PublishAccrualUpdated(ctx context.Context, a *db.Accrual) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockaPublisher) PublishAccrualDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- AllocationRule events ----------------
func (m *MockaPublisher) PublishAllocationRuleCreated(ctx context.Context, r *db.AllocationRule) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockaPublisher) PublishAllocationRuleUpdated(ctx context.Context, r *db.AllocationRule) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockaPublisher) PublishAllocationRuleDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Audit events ----------------
func (m *MockaPublisher) PublishAuditRecorded(ctx context.Context, event *db.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ---------------- Budget events ----------------
func (m *MockaPublisher) PublishBudgetCreated(ctx context.Context, b *db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockaPublisher) PublishBudgetUpdated(ctx context.Context, b *db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockaPublisher) PublishBudgetDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockaPublisher) PublishBudgetAllocated(ctx context.Context, ba *db.BudgetAllocation) error {
	args := m.Called(ctx, ba)
	return args.Error(0)
}

func (m *MockaPublisher) PublishBudgetAllocationUpdated(ctx context.Context, ba *db.BudgetAllocation) error {
	args := m.Called(ctx, ba)
	return args.Error(0)
}

func (m *MockaPublisher) PublishBudgetAllocationDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Cash Flow Forecast events ----------------
func (m *MockaPublisher) PublishCashFlowForecastGenerated(ctx context.Context, forecast *db.CashFlowForecast) error {
	args := m.Called(ctx, forecast)
	return args.Error(0)
}

func (m *MockaPublisher) PublishCashFlowForecastFetched(ctx context.Context, forecast *db.CashFlowForecast) error {
	args := m.Called(ctx, forecast)
	return args.Error(0)
}

func (m *MockaPublisher) PublishCashFlowForecastListed(ctx context.Context, forecasts []db.CashFlowForecast) error {
	args := m.Called(ctx, forecasts)
	return args.Error(0)
}

// ---------------- Consolidation events ----------------
func (m *MockaPublisher) PublishConsolidationCreated(ctx context.Context, c *db.Consolidation) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockaPublisher) PublishConsolidationDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Credit/Debit Note events ----------------
func (m *MockaPublisher) PublishCreditDebitNoteCreated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockaPublisher) PublishCreditDebitNoteUpdated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockaPublisher) PublishCreditDebitNoteDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Exchange Rate events ----------------
func (m *MockaPublisher) PublishExchangeRateCreated(ctx context.Context, rate *db.ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockaPublisher) PublishExchangeRateUpdated(ctx context.Context, rate *db.ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockaPublisher) PublishExchangeRateDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Expense events ----------------
func (m *MockaPublisher) PublishExpenseCreated(ctx context.Context, exp *db.Expense) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockaPublisher) PublishExpenseUpdated(ctx context.Context, exp *db.Expense) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockaPublisher) PublishExpenseDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- CostCenter events ----------------
func (m *MockaPublisher) PublishCostCenterCreated(ctx context.Context, cc *db.CostCenter) error {
	args := m.Called(ctx, cc)
	return args.Error(0)
}

func (m *MockaPublisher) PublishCostCenterUpdated(ctx context.Context, cc *db.CostCenter) error {
	args := m.Called(ctx, cc)
	return args.Error(0)
}

func (m *MockaPublisher) PublishCostCenterDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- CostAllocation events ----------------
func (m *MockaPublisher) PublishCostAllocationAllocated(ctx context.Context, ca *db.CostAllocation) error {
	args := m.Called(ctx, ca)
	return args.Error(0)
}

func (m *MockaPublisher) PublishCostAllocationListed(ctx context.Context, allocs []db.CostAllocation) error {
	args := m.Called(ctx, allocs)
	return args.Error(0)
}

// ---------------- Finance domain events ----------------
func (m *MockaPublisher) PublishFinanceInvoiceCreated(ctx context.Context, ev *db.FinanceInvoiceCreatedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockaPublisher) PublishFinancePaymentReceived(ctx context.Context, ev *db.FinancePaymentReceivedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockaPublisher) PublishInventoryCostPosted(ctx context.Context, ev *db.InventoryCostPostedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockaPublisher) PublishPayrollPosted(ctx context.Context, ev *db.PayrollPostedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockaPublisher) PublishVendorBillApproved(ctx context.Context, ev *db.VendorBillApprovedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

// ---------------- GST events ----------------
func (m *MockaPublisher) PublishGstBreakupAdded(ctx context.Context, breakup *db.GstBreakup) error {
	args := m.Called(ctx, breakup)
	return args.Error(0)
}

func (m *MockaPublisher) PublishGstRegimeAdded(ctx context.Context, regime *db.GstRegime) error {
	args := m.Called(ctx, regime)
	return args.Error(0)
}

func (m *MockaPublisher) PublishGstDocStatusAdded(ctx context.Context, status *db.GstDocStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

// ---------------- Invoice events ----------------
func (m *MockaPublisher) PublishInvoiceUpdated(ctx context.Context, inv *db.Invoice) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockaPublisher) PublishInvoiceDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockaPublisher) PublishInvoiceItemCreated(ctx context.Context, item *db.InvoiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockaPublisher) PublishInvoiceTaxAdded(ctx context.Context, tax *db.InvoiceTax) error {
	args := m.Called(ctx, tax)
	return args.Error(0)
}

func (m *MockaPublisher) PublishInvoiceDiscountAdded(ctx context.Context, disc *db.InvoiceDiscount) error {
	args := m.Called(ctx, disc)
	return args.Error(0)
}

// ---------------- Bank Transaction events ----------------
func (m *MockaPublisher) PublishBankTransactionImported(ctx context.Context, ev *db.BankAccount) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockaPublisher) PublishBankTransactionReconciled(ctx context.Context, ev *db.BankAccount) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}
// ------------------ Test Cases ------------------
func TestGstService_AddAndGetBreakup(t *testing.T) {
	mockRepo := new(MockGstsRepository)
	mockPub := new(MockaPublisher)
	service := services.NewGstService(mockRepo, mockPub)

	ctx := context.Background()
	invoiceID := uuid.New()
	cgst := 10.0
	sgst := 10.0
	igst := 0.0
	totalGst := 20.0

	expected := db.GstBreakup{
		ID:           uuid.New(),
		InvoiceID:    invoiceID,
		TaxableAmount: "100.0",
		Cgst:         sql.NullString{String: "10", Valid: true},
		Sgst:         sql.NullString{String: "10", Valid: true},
		Igst:         sql.NullString{String: "0", Valid: true},
		TotalGst:     sql.NullString{String: "20", Valid: true},
	}

	mockRepo.On("AddGstBreakup", ctx, invoiceID, 100.0, &cgst, &sgst, &igst, &totalGst).Return(expected, nil)
	mockPub.On("PublishGstBreakupAdded", ctx, &expected).Return(nil)
	mockRepo.On("GetGstBreakup", ctx, invoiceID).Return(expected, nil)

	// AddGstBreakup
	result, err := service.AddGstBreakup(ctx, invoiceID.String(), 100.0, &cgst, &sgst, &igst, &totalGst)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)

	// GetGstBreakup
	result, err = service.GetGstBreakup(ctx, invoiceID.String())
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestGstService_AddAndGetRegime(t *testing.T) {
	mockRepo := new(MockGstsRepository)
	mockPub := new(MockaPublisher)
	service := services.NewGstService(mockRepo, mockPub)

	ctx := context.Background()
	invoiceID := uuid.New()
	gstin := "27AAAAA0000A1Z5"
	place := "Maharashtra"
	reverse := true

	expected := db.GstRegime{
		ID:           uuid.New(),
		InvoiceID:    invoiceID,
		Gstin:        gstin,
		PlaceOfSupply: place,
		ReverseCharge: sql.NullBool{Bool: reverse, Valid: true},
	}

	mockRepo.On("AddGstRegime", ctx, invoiceID, gstin, place, &reverse).Return(expected, nil)
	mockPub.On("PublishGstRegimeAdded", ctx, &expected).Return(nil)
	mockRepo.On("GetGstRegime", ctx, invoiceID).Return(expected, nil)

	// AddGstRegime
	result, err := service.AddGstRegime(ctx, invoiceID.String(), gstin, place, &reverse)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)

	// GetGstRegime
	result, err = service.GetGstRegime(ctx, invoiceID.String())
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestGstService_AddAndGetDocStatus(t *testing.T) {
	mockRepo := new(MockGstsRepository)
	mockPub := new(MockaPublisher)
	service := services.NewGstService(mockRepo, mockPub)

	ctx := context.Background()
	invoiceID := uuid.New()
	now := time.Now()
	einvoiceStatus := "SUCCESS"
	irn := "IRN123"
	ackNo := "ACK456"
	ewayStatus := "VALID"
	ewayBillNo := "EB123"

	expected := db.GstDocStatus{
		ID:            uuid.New(),
		InvoiceID:     invoiceID,
		EinvoiceStatus: sql.NullString{String: einvoiceStatus, Valid: true},
		Irn:           sql.NullString{String: irn, Valid: true},
		AckNo:         sql.NullString{String: ackNo, Valid: true},
		AckDate:       sql.NullTime{Time: now, Valid: true},
		EwayStatus:    sql.NullString{String: ewayStatus, Valid: true},
		EwayBillNo:    sql.NullString{String: ewayBillNo, Valid: true},
		EwayValidUpto: sql.NullTime{Time: now, Valid: true},
	}

	mockRepo.On("AddGstDocStatus", ctx, invoiceID,
		&einvoiceStatus, &irn, &ackNo, &now,
		&ewayStatus, &ewayBillNo, &now,
		(*string)(nil), (*time.Time)(nil),
	).Return(expected, nil)
	mockPub.On("PublishGstDocStatusAdded", ctx, &expected).Return(nil)
	mockRepo.On("GetGstDocStatus", ctx, invoiceID).Return(expected, nil)

	// AddGstDocStatus
	result, err := service.AddGstDocStatus(ctx, invoiceID.String(), &einvoiceStatus, &irn, &ackNo, &now, &ewayStatus, &ewayBillNo, &now, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)

	// GetGstDocStatus
	result, err = service.GetGstDocStatus(ctx, invoiceID.String())
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}
