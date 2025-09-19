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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ----------- Mocks -----------

type MockInvoiceRepo struct {
	mock.Mock
}

func (m *MockInvoiceRepo) CreateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error) {
	args := m.Called(ctx, inv)
	return args.Get(0).(db.Invoice), args.Error(1)
}

func (m *MockInvoiceRepo) GetInvoice(ctx context.Context, id uuid.UUID) (db.Invoice, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Invoice), args.Error(1)
}

func (m *MockInvoiceRepo) UpdateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error) {
	args := m.Called(ctx, inv)
	return args.Get(0).(db.Invoice), args.Error(1)
}

func (m *MockInvoiceRepo) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInvoiceRepo) ListInvoices(ctx context.Context, limit, offset int32) ([]db.Invoice, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.Invoice), args.Error(1)
}

func (m *MockInvoiceRepo) SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error) {
	args := m.Called(ctx, query, limit, offset)
	return args.Get(0).([]db.Invoice), args.Error(1)
}

func (m *MockInvoiceRepo) CreateInvoiceItem(ctx context.Context, item db.InvoiceItem) (db.InvoiceItem, error) {
	args := m.Called(ctx, item)
	return args.Get(0).(db.InvoiceItem), args.Error(1)
}

func (m *MockInvoiceRepo) ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]db.InvoiceItem, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).([]db.InvoiceItem), args.Error(1)
}

func (m *MockInvoiceRepo) AddInvoiceTax(ctx context.Context, tax db.InvoiceTax) (db.InvoiceTax, error) {
	args := m.Called(ctx, tax)
	return args.Get(0).(db.InvoiceTax), args.Error(1)
}

func (m *MockInvoiceRepo) AddInvoiceDiscount(ctx context.Context, disc db.InvoiceDiscount) (db.InvoiceDiscount, error) {
	args := m.Called(ctx, disc)
	return args.Get(0).(db.InvoiceDiscount), args.Error(1)
}

// ---- Mock FinanceEventService ----
type MockFinanceEventService struct {
	mock.Mock
	*services.FinanceEventService
}

func (m *MockFinanceEventService) RecordInvoiceCreated(ctx context.Context, event db.FinanceInvoiceCreatedEvent) (db.FinanceInvoiceCreatedEvent, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(db.FinanceInvoiceCreatedEvent), args.Error(1)
}

func (m *MockFinanceEventService) RecordInvoiceUpdated(ctx context.Context, event db.Invoice) (db.Invoice, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(db.Invoice), args.Error(1)
}

func (m *MockFinanceEventService) RecordInvoiceDeleted(ctx context.Context, invoiceID uuid.UUID) error {
	args := m.Called(ctx, invoiceID)
	return args.Error(0)
}

func (m *MockFinanceEventService) RecordBudgetCreated(ctx context.Context, b db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockFinanceEventService) RecordBudgetUpdated(ctx context.Context, b db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockFinanceEventService) RecordBudgetDeleted(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}


// ---- Mock Publisher ----

// MockPublisher implements ports.EventPublisher fully for testing
type MockmPublisher struct {
	mock.Mock
}

// Generic publish
func (m *MockmPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	args := m.Called(ctx, topic, key, payload)
	return args.Error(0)
}

func (m *MockmPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ---------------- Account events ----------------
func (m *MockmPublisher) PublishAccountCreated(ctx context.Context, a *db.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockmPublisher) PublishAccountUpdated(ctx context.Context, a *db.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockmPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Journal events ----------------
func (m *MockmPublisher) PublishJournalCreated(ctx context.Context, j *db.JournalEntry) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockmPublisher) PublishJournalUpdated(ctx context.Context, j *db.JournalEntry) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockmPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Accrual events ----------------
func (m *MockmPublisher) PublishAccrualCreated(ctx context.Context, a *db.Accrual) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockmPublisher) PublishAccrualUpdated(ctx context.Context, a *db.Accrual) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockmPublisher) PublishAccrualDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- AllocationRule events ----------------
func (m *MockmPublisher) PublishAllocationRuleCreated(ctx context.Context, r *db.AllocationRule) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockmPublisher) PublishAllocationRuleUpdated(ctx context.Context, r *db.AllocationRule) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockmPublisher) PublishAllocationRuleDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Audit events ----------------
func (m *MockmPublisher) PublishAuditRecorded(ctx context.Context, event *db.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ---------------- Budget events ----------------
func (m *MockmPublisher) PublishBudgetCreated(ctx context.Context, b *db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockmPublisher) PublishBudgetUpdated(ctx context.Context, b *db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockmPublisher) PublishBudgetDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockmPublisher) PublishBudgetAllocated(ctx context.Context, ba *db.BudgetAllocation) error {
	args := m.Called(ctx, ba)
	return args.Error(0)
}

func (m *MockmPublisher) PublishBudgetAllocationUpdated(ctx context.Context, ba *db.BudgetAllocation) error {
	args := m.Called(ctx, ba)
	return args.Error(0)
}

func (m *MockmPublisher) PublishBudgetAllocationDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Cash Flow Forecast events ----------------
func (m *MockmPublisher) PublishCashFlowForecastGenerated(ctx context.Context, forecast *db.CashFlowForecast) error {
	args := m.Called(ctx, forecast)
	return args.Error(0)
}

func (m *MockmPublisher) PublishCashFlowForecastFetched(ctx context.Context, forecast *db.CashFlowForecast) error {
	args := m.Called(ctx, forecast)
	return args.Error(0)
}

func (m *MockmPublisher) PublishCashFlowForecastListed(ctx context.Context, forecasts []db.CashFlowForecast) error {
	args := m.Called(ctx, forecasts)
	return args.Error(0)
}

// ---------------- Consolidation events ----------------
func (m *MockmPublisher) PublishConsolidationCreated(ctx context.Context, c *db.Consolidation) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockmPublisher) PublishConsolidationDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Credit/Debit Note events ----------------
func (m *MockmPublisher) PublishCreditDebitNoteCreated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockmPublisher) PublishCreditDebitNoteUpdated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockmPublisher) PublishCreditDebitNoteDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Exchange Rate events ----------------
func (m *MockmPublisher) PublishExchangeRateCreated(ctx context.Context, rate *db.ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockmPublisher) PublishExchangeRateUpdated(ctx context.Context, rate *db.ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockmPublisher) PublishExchangeRateDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Expense events ----------------
func (m *MockmPublisher) PublishExpenseCreated(ctx context.Context, exp *db.Expense) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockmPublisher) PublishExpenseUpdated(ctx context.Context, exp *db.Expense) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockmPublisher) PublishExpenseDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- CostCenter events ----------------
func (m *MockmPublisher) PublishCostCenterCreated(ctx context.Context, cc *db.CostCenter) error {
	args := m.Called(ctx, cc)
	return args.Error(0)
}

func (m *MockmPublisher) PublishCostCenterUpdated(ctx context.Context, cc *db.CostCenter) error {
	args := m.Called(ctx, cc)
	return args.Error(0)
}

func (m *MockmPublisher) PublishCostCenterDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- CostAllocation events ----------------
func (m *MockmPublisher) PublishCostAllocationAllocated(ctx context.Context, ca *db.CostAllocation) error {
	args := m.Called(ctx, ca)
	return args.Error(0)
}

func (m *MockmPublisher) PublishCostAllocationListed(ctx context.Context, allocs []db.CostAllocation) error {
	args := m.Called(ctx, allocs)
	return args.Error(0)
}

// ---------------- Finance domain events ----------------
func (m *MockmPublisher) PublishFinanceInvoiceCreated(ctx context.Context, ev *db.FinanceInvoiceCreatedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockmPublisher) PublishFinancePaymentReceived(ctx context.Context, ev *db.FinancePaymentReceivedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockmPublisher) PublishInventoryCostPosted(ctx context.Context, ev *db.InventoryCostPostedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockmPublisher) PublishPayrollPosted(ctx context.Context, ev *db.PayrollPostedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockmPublisher) PublishVendorBillApproved(ctx context.Context, ev *db.VendorBillApprovedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

// ---------------- GST events ----------------
func (m *MockmPublisher) PublishGstBreakupAdded(ctx context.Context, breakup *db.GstBreakup) error {
	args := m.Called(ctx, breakup)
	return args.Error(0)
}

func (m *MockmPublisher) PublishGstRegimeAdded(ctx context.Context, regime *db.GstRegime) error {
	args := m.Called(ctx, regime)
	return args.Error(0)
}

func (m *MockmPublisher) PublishGstDocStatusAdded(ctx context.Context, status *db.GstDocStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

// ---------------- Invoice events ----------------
func (m *MockmPublisher) PublishInvoiceUpdated(ctx context.Context, inv *db.Invoice) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockmPublisher) PublishInvoiceDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockmPublisher) PublishInvoiceItemCreated(ctx context.Context, item *db.InvoiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockmPublisher) PublishInvoiceTaxAdded(ctx context.Context, tax *db.InvoiceTax) error {
	args := m.Called(ctx, tax)
	return args.Error(0)
}

func (m *MockmPublisher) PublishInvoiceDiscountAdded(ctx context.Context, disc *db.InvoiceDiscount) error {
	args := m.Called(ctx, disc)
	return args.Error(0)
}

// ---------------- Bank Transaction events ----------------
func (m *MockmPublisher) PublishBankTransactionImported(ctx context.Context, ev *db.BankAccount) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockmPublisher) PublishBankTransactionReconciled(ctx context.Context, ev *db.BankAccount) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}


// ----------- Tests -----------

func TestInvoiceService_CreateInvoice(t *testing.T) {
	ctx := context.Background()

	repo := new(MockInvoiceRepo)
	eventSvc := new(MockFinanceEventService)
	pub := new(MockPublisher)

	service := services.NewInvoiceService(repo, nil, pub)

	inv := db.Invoice{
		ID:             uuid.New(),
		InvoiceNumber:  "INV-001",
		InvoiceDate:    time.Now(),
		GrandTotal:     "1000",
		OrganizationID: "ORG-123",
	}

	event := db.FinanceInvoiceCreatedEvent{
		InvoiceID:      inv.ID,
		InvoiceNumber:  inv.InvoiceNumber,
		InvoiceDate:    inv.InvoiceDate,
		Total:          inv.GrandTotal,
		OrganizationID: inv.OrganizationID,
	}

	repo.On("CreateInvoice", ctx, inv).Return(inv, nil)
	eventSvc.On("RecordInvoiceCreated", ctx, event).Return(event, nil)

	res, err := service.CreateInvoice(ctx, inv)
	require.NoError(t, err)
	require.Equal(t, inv.ID, res.ID)

	repo.AssertCalled(t, "CreateInvoice", ctx, inv)
	eventSvc.AssertCalled(t, "RecordInvoiceCreated", ctx, event)
}


func TestInvoiceService_UpdateInvoice(t *testing.T) {
	ctx := context.Background()
	repo := new(MockInvoiceRepo)
	pub := new(MockmPublisher)
	svc := services.NewInvoiceService(repo, nil, pub)

	inv := db.Invoice{ID: uuid.New(), InvoiceNumber: "INV-123"}
	updated := inv
	updated.Status = "Paid"

	repo.On("UpdateInvoice", ctx, inv).Return(updated, nil)
	pub.On("PublishInvoiceUpdated", ctx, &updated).Return(nil)

	got, err := svc.UpdateInvoice(ctx, inv)
	require.NoError(t, err)
	require.Equal(t, updated, got)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestInvoiceService_DeleteInvoice(t *testing.T) {
	ctx := context.Background()
	repo := new(MockInvoiceRepo)
	pub := new(MockmPublisher)
	svc := services.NewInvoiceService(repo, nil, pub)

	id := uuid.New()

	repo.On("DeleteInvoice", ctx, id).Return(nil)
	pub.On("PublishInvoiceDeleted", ctx, id.String()).Return(nil)

	err := svc.DeleteInvoice(ctx, id)
	require.NoError(t, err)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestInvoiceService_CreateInvoiceItem(t *testing.T) {
	ctx := context.Background()
	repo := new(MockInvoiceRepo)
	pub := new(MockmPublisher)
	svc := services.NewInvoiceService(repo, nil, pub)

	item := db.InvoiceItem{ID: uuid.New(), InvoiceID: uuid.New(), Name: "Item1"}
	repo.On("CreateInvoiceItem", ctx, item).Return(item, nil)
	pub.On("PublishInvoiceItemCreated", ctx, &item).Return(nil)

	got, err := svc.CreateInvoiceItem(ctx, item)
	require.NoError(t, err)
	require.Equal(t, item, got)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestInvoiceService_AddInvoiceTax(t *testing.T) {
	ctx := context.Background()
	repo := new(MockInvoiceRepo)
	pub := new(MockmPublisher)
	svc := services.NewInvoiceService(repo, nil, pub)

	tax := db.InvoiceTax{ID: uuid.New(), InvoiceID: uuid.New(), Name: "GST", Rate: "18%", Amount: "180"}
	repo.On("AddInvoiceTax", ctx, tax).Return(tax, nil)
	pub.On("PublishInvoiceTaxAdded", ctx, &tax).Return(nil)

	got, err := svc.AddInvoiceTax(ctx, tax)
	require.NoError(t, err)
	require.Equal(t, tax, got)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestInvoiceService_AddInvoiceDiscount(t *testing.T) {
	ctx := context.Background()
	repo := new(MockInvoiceRepo)
	pub := new(MockmPublisher)
	svc := services.NewInvoiceService(repo, nil, pub)

	disc := db.InvoiceDiscount{ID: uuid.New(), InvoiceID: uuid.New(), Description: sql.NullString{String: "Promo", Valid: true}, Amount: "100"}
	repo.On("AddInvoiceDiscount", ctx, disc).Return(disc, nil)
	pub.On("PublishInvoiceDiscountAdded", ctx, &disc).Return(nil)

	got, err := svc.AddInvoiceDiscount(ctx, disc)
	require.NoError(t, err)
	require.Equal(t, disc, got)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestInvoiceService_GetInvoice_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockInvoiceRepo)
	svc := services.NewInvoiceService(repo, nil, nil)

	id := uuid.New()
	repo.On("GetInvoice", ctx, id).Return(db.Invoice{}, errors.New("not found"))

	_, err := svc.GetInvoice(ctx, id)
	require.Error(t, err)
	repo.AssertExpectations(t)
}
