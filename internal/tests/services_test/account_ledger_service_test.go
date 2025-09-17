package services_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

// ---------------- Mock Repositories ----------------
type MockAccountRepo struct {
	mock.Mock
}

func (m *MockAccountRepo) Create(ctx context.Context, a *db.Account) (*db.Account, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*db.Account), args.Error(1)
}

func (m *MockAccountRepo) Get(ctx context.Context, id uuid.UUID) (*db.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.Account), args.Error(1)
}

func (m *MockAccountRepo) Update(ctx context.Context, a *db.Account) (*db.Account, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*db.Account), args.Error(1)
}

func (m *MockAccountRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAccountRepo) List(ctx context.Context, limit, offset int32) ([]*db.Account, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*db.Account), args.Error(1)
}


type MockJournalRepo struct{ mock.Mock }

func (m *MockJournalRepo) Create(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error) {
	args := m.Called(ctx, j)
	return args.Get(0).(*db.JournalEntry), args.Error(1)
}
func (m *MockJournalRepo) Get(ctx context.Context, id uuid.UUID) (*db.JournalEntry, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.JournalEntry), args.Error(1)
}
func (m *MockJournalRepo) Update(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error) {
	args := m.Called(ctx, j)
	return args.Get(0).(*db.JournalEntry), args.Error(1)
}
func (m *MockJournalRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockJournalRepo) List(ctx context.Context, limit, offset int32) ([]*db.JournalEntry, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*db.JournalEntry), args.Error(1)
}

type MockLedgerRepo struct{ mock.Mock }

func (m *MockLedgerRepo) List(ctx context.Context, limit, offset int32) ([]*db.LedgerEntry, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*db.LedgerEntry), args.Error(1)
}

// ---------------- Mock Publisher ----------------
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	args := m.Called(ctx, topic, key, payload)
	return args.Error(0)
}

func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPublisher) PublishAccountCreated(ctx context.Context, a *db.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockPublisher) PublishAccountUpdated(ctx context.Context, a *db.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPublisher) PublishJournalCreated(ctx context.Context, j *db.JournalEntry) error {
	return m.Called(ctx, j).Error(0)
}
func (m *MockPublisher) PublishJournalUpdated(ctx context.Context, j *db.JournalEntry) error {
	return m.Called(ctx, j).Error(0)
}
func (m *MockPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- Accrual events ----------------
func (m *MockPublisher) PublishAccrualCreated(ctx context.Context, a *db.Accrual) error {
	return m.Called(ctx, a).Error(0)
}
func (m *MockPublisher) PublishAccrualUpdated(ctx context.Context, a *db.Accrual) error {
	return m.Called(ctx, a).Error(0)
}
func (m *MockPublisher) PublishAccrualDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- AllocationRule events ----------------
func (m *MockPublisher) PublishAllocationRuleCreated(ctx context.Context, r *db.AllocationRule) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockPublisher) PublishAllocationRuleUpdated(ctx context.Context, r *db.AllocationRule) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockPublisher) PublishAllocationRuleDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- Audit events ----------------
func (m *MockPublisher) PublishAuditRecorded(ctx context.Context, event *db.AuditEvent) error {
	return m.Called(ctx, event).Error(0)
}

// ---------------- Budget events ----------------
func (m *MockPublisher) PublishBudgetCreated(ctx context.Context, b *db.Budget) error {
	return m.Called(ctx, b).Error(0)
}
func (m *MockPublisher) PublishBudgetUpdated(ctx context.Context, b *db.Budget) error {
	return m.Called(ctx, b).Error(0)
}
func (m *MockPublisher) PublishBudgetDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockPublisher) PublishBudgetAllocated(ctx context.Context, ba *db.BudgetAllocation) error {
	return m.Called(ctx, ba).Error(0)
}
func (m *MockPublisher) PublishBudgetAllocationUpdated(ctx context.Context, ba *db.BudgetAllocation) error {
	return m.Called(ctx, ba).Error(0)
}
func (m *MockPublisher) PublishBudgetAllocationDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- Cash Flow Forecast events ----------------
func (m *MockPublisher) PublishCashFlowForecastGenerated(ctx context.Context, forecast *db.CashFlowForecast) error {
	return m.Called(ctx, forecast).Error(0)
}
func (m *MockPublisher) PublishCashFlowForecastFetched(ctx context.Context, forecast *db.CashFlowForecast) error {
	return m.Called(ctx, forecast).Error(0)
}
func (m *MockPublisher) PublishCashFlowForecastListed(ctx context.Context, forecasts []db.CashFlowForecast) error {
	return m.Called(ctx, forecasts).Error(0)
}

// ---------------- Consolidation events ----------------
func (m *MockPublisher) PublishConsolidationCreated(ctx context.Context, c *db.Consolidation) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockPublisher) PublishConsolidationDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- Credit/Debit Note events ----------------
func (m *MockPublisher) PublishCreditDebitNoteCreated(ctx context.Context, note *db.CreditDebitNote) error {
	return m.Called(ctx, note).Error(0)
}
func (m *MockPublisher) PublishCreditDebitNoteUpdated(ctx context.Context, note *db.CreditDebitNote) error {
	return m.Called(ctx, note).Error(0)
}
func (m *MockPublisher) PublishCreditDebitNoteDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- Exchange Rate events ----------------
func (m *MockPublisher) PublishExchangeRateCreated(ctx context.Context, rate *db.ExchangeRate) error {
	return m.Called(ctx, rate).Error(0)
}
func (m *MockPublisher) PublishExchangeRateUpdated(ctx context.Context, rate *db.ExchangeRate) error {
	return m.Called(ctx, rate).Error(0)
}
func (m *MockPublisher) PublishExchangeRateDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- Expense events ----------------
func (m *MockPublisher) PublishExpenseCreated(ctx context.Context, exp *db.Expense) error {
	return m.Called(ctx, exp).Error(0)
}
func (m *MockPublisher) PublishExpenseUpdated(ctx context.Context, exp *db.Expense) error {
	return m.Called(ctx, exp).Error(0)
}
func (m *MockPublisher) PublishExpenseDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- CostCenter events ----------------
func (m *MockPublisher) PublishCostCenterCreated(ctx context.Context, cc *db.CostCenter) error {
	return m.Called(ctx, cc).Error(0)
}
func (m *MockPublisher) PublishCostCenterUpdated(ctx context.Context, cc *db.CostCenter) error {
	return m.Called(ctx, cc).Error(0)
}
func (m *MockPublisher) PublishCostCenterDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ---------------- CostAllocation events ----------------
func (m *MockPublisher) PublishCostAllocationAllocated(ctx context.Context, ca *db.CostAllocation) error {
	return m.Called(ctx, ca).Error(0)
}
func (m *MockPublisher) PublishCostAllocationListed(ctx context.Context, allocs []db.CostAllocation) error {
	return m.Called(ctx, allocs).Error(0)
}

// ---------------- Finance domain events ----------------
func (m *MockPublisher) PublishFinanceInvoiceCreated(ctx context.Context, ev *db.FinanceInvoiceCreatedEvent) error {
	return m.Called(ctx, ev).Error(0)
}
func (m *MockPublisher) PublishFinancePaymentReceived(ctx context.Context, ev *db.FinancePaymentReceivedEvent) error {
	return m.Called(ctx, ev).Error(0)
}
func (m *MockPublisher) PublishInventoryCostPosted(ctx context.Context, ev *db.InventoryCostPostedEvent) error {
	return m.Called(ctx, ev).Error(0)
}
func (m *MockPublisher) PublishPayrollPosted(ctx context.Context, ev *db.PayrollPostedEvent) error {
	return m.Called(ctx, ev).Error(0)
}
func (m *MockPublisher) PublishVendorBillApproved(ctx context.Context, ev *db.VendorBillApprovedEvent) error {
	return m.Called(ctx, ev).Error(0)
}

// ---------------- GST events ----------------
func (m *MockPublisher) PublishGstBreakupAdded(ctx context.Context, breakup *db.GstBreakup) error {
	return m.Called(ctx, breakup).Error(0)
}
func (m *MockPublisher) PublishGstRegimeAdded(ctx context.Context, regime *db.GstRegime) error {
	return m.Called(ctx, regime).Error(0)
}
func (m *MockPublisher) PublishGstDocStatusAdded(ctx context.Context, status *db.GstDocStatus) error {
	return m.Called(ctx, status).Error(0)
}

// ---------------- Invoice events ----------------
func (m *MockPublisher) PublishInvoiceUpdated(ctx context.Context, inv *db.Invoice) error {
	return m.Called(ctx, inv).Error(0)
}
func (m *MockPublisher) PublishInvoiceDeleted(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockPublisher) PublishInvoiceItemCreated(ctx context.Context, item *db.InvoiceItem) error {
	return m.Called(ctx, item).Error(0)
}
func (m *MockPublisher) PublishInvoiceTaxAdded(ctx context.Context, tax *db.InvoiceTax) error {
	return m.Called(ctx, tax).Error(0)
}
func (m *MockPublisher) PublishInvoiceDiscountAdded(ctx context.Context, disc *db.InvoiceDiscount) error {
	return m.Called(ctx, disc).Error(0)
}

// ---------------- Bank Transaction events ----------------
func (m *MockPublisher) PublishBankTransactionImported(ctx context.Context, ev *db.BankAccount) error {
	return m.Called(ctx, ev).Error(0)
}
func (m *MockPublisher) PublishBankTransactionReconciled(ctx context.Context, ev *db.BankAccount) error {
	return m.Called(ctx, ev).Error(0)
}

// ---------------- Tests ----------------
func TestServices_CRUD(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	accID := uuid.New()
	jID := uuid.New()

	// ---------------- AccountService Test ----------------
	acc := &db.Account{
		ID:                 accID,
		Code:               "1001",
		Name:               "Cash",
		Type:               "ASSET",
		Status:             "ACTIVE",
		AllowManualJournal: true,
		CreatedAt:          now,
		CreatedBy:          sql.NullString{String: "admin", Valid: true},
		UpdatedAt:          now,
		UpdatedBy:          sql.NullString{String: "admin", Valid: true},
	}

	accRepo := new(MockAccountRepo)
	pub := new(MockPublisher)
	accService := services.NewAccountService(accRepo, pub)

	accRepo.On("Create", ctx, acc).Return(acc, nil)
	pub.On("PublishAccountCreated", mock.Anything, acc).Return(nil)
	createdAcc, err := accService.Create(ctx, acc)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, createdAcc.ID)

	accRepo.On("Get", ctx, accID).Return(acc, nil)
	gotAcc, err := accService.Get(ctx, accID)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, gotAcc.ID)

	acc.Name = "Updated Cash"
	accRepo.On("Update", ctx, acc).Return(acc, nil)
	pub.On("PublishAccountUpdated", mock.Anything, acc).Return(nil)
	updatedAcc, err := accService.Update(ctx, acc)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Cash", updatedAcc.Name)

	accRepo.On("List", ctx, int32(10), int32(0)).Return([]*db.Account{acc}, nil)
	listAcc, err := accService.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, listAcc, 1)

	accRepo.On("Delete", ctx, accID).Return(nil)
	pub.On("PublishAccountDeleted", mock.Anything, accID.String()).Return(nil)
	err = accService.Delete(ctx, accID)
	assert.NoError(t, err)

	// ---------------- JournalService Test ----------------
	journal := &db.JournalEntry{
		ID:         jID,
		JournalDate: now,
		Reference:  sql.NullString{String: "Ref001", Valid: true},
		Memo:       sql.NullString{String: "Memo", Valid: true},
		CreatedAt:  now, 
		CreatedBy:  sql.NullString{String: "admin", Valid: true},
		UpdatedAt:  now, 
		UpdatedBy:  sql.NullString{String: "admin", Valid: true},
	}

	jRepo := new(MockJournalRepo)
	jService := services.NewJournalService(jRepo, pub)

	jRepo.On("Create", ctx, journal).Return(journal, nil)
	pub.On("PublishJournalCreated", mock.Anything, journal).Return(nil)
	createdJournal, err := jService.Create(ctx, journal)
	assert.NoError(t, err)
	assert.Equal(t, jID, createdJournal.ID)

	jRepo.On("Get", ctx, jID).Return(journal, nil)
	gotJournal, err := jService.Get(ctx, jID)
	assert.NoError(t, err)
	assert.Equal(t, jID, gotJournal.ID)

	journal.Memo = sql.NullString{String: "Updated Memo", Valid: true}
	jRepo.On("Update", ctx, journal).Return(journal, nil)
	pub.On("PublishJournalUpdated", mock.Anything, journal).Return(nil)
	updatedJournal, err := jService.Update(ctx, journal)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Memo", updatedJournal.Memo.String)

	jRepo.On("List", ctx, int32(10), int32(0)).Return([]*db.JournalEntry{journal}, nil)
	listJournal, err := jService.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, listJournal, 1)

	jRepo.On("Delete", ctx, jID).Return(nil)
	pub.On("PublishJournalDeleted", mock.Anything, jID.String()).Return(nil)
	err = jService.Delete(ctx, jID)
	assert.NoError(t, err)

	// ---------------- LedgerService Test ----------------
	ledgerEntry := &db.LedgerEntry{
		EntryID:        uuid.New(),
		AccountID:      accID,
		Side:           "DEBIT",
		Amount:         "100.00",
		TransactionDate: now,
	}
	lRepo := new(MockLedgerRepo)
	lService := services.NewLedgerService(lRepo)

	lRepo.On("List", ctx, int32(10), int32(0)).Return([]*db.LedgerEntry{ledgerEntry}, nil)
	listLedger, err := lService.List(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, listLedger, 1)
	assert.Equal(t, ledgerEntry.EntryID, listLedger[0].EntryID)
}
