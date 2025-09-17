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

// -------------------- MOCKS --------------------

type MockBudgetRepo struct {
	mock.Mock
}

func (m *MockBudgetRepo) Create(ctx context.Context, b *db.Budget) (*db.Budget, error) {
	args := m.Called(ctx, b)
	return args.Get(0).(*db.Budget), args.Error(1)
}
func (m *MockBudgetRepo) Get(ctx context.Context, id uuid.UUID) (*db.Budget, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.Budget), args.Error(1)
}
func (m *MockBudgetRepo) List(ctx context.Context, limit, offset int32) ([]db.Budget, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]db.Budget), args.Error(1)
}
func (m *MockBudgetRepo) Update(ctx context.Context, b *db.Budget) (*db.Budget, error) {
	args := m.Called(ctx, b)
	return args.Get(0).(*db.Budget), args.Error(1)
}
func (m *MockBudgetRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockBudgetRepo) Allocate(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	args := m.Called(ctx, ba)
	return args.Get(0).(*db.BudgetAllocation), args.Error(1)
}
func (m *MockBudgetRepo) GetAllocation(ctx context.Context, id uuid.UUID) (*db.BudgetAllocation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.BudgetAllocation), args.Error(1)
}
func (m *MockBudgetRepo) ListAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]db.BudgetAllocation, error) {
	args := m.Called(ctx, budgetID, limit, offset)
	return args.Get(0).([]db.BudgetAllocation), args.Error(1)
}
func (m *MockBudgetRepo) UpdateAllocation(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	args := m.Called(ctx, ba)
	return args.Get(0).(*db.BudgetAllocation), args.Error(1)
}
func (m *MockBudgetRepo) DeleteAllocation(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockBudgetRepo) GetBudgetComparisonReport(ctx context.Context, id uuid.UUID) (*db.GetBudgetComparisonReportRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*db.GetBudgetComparisonReportRow), args.Error(1)
}

// -------------------- PUBLISHER MOCK --------------------

type MocksPublisher struct {
	mock.Mock
}

// ---------------- Generic ----------------
func (m *MocksPublisher) Publish(ctx context.Context, topic string, key string, value []byte) error {
	args := m.Called(ctx, topic, key, value)
	return args.Error(0)
}

func (m *MocksPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ---------------- Account events ----------------
func (m *MocksPublisher) PublishAccountCreated(ctx context.Context, a *db.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MocksPublisher) PublishAccountUpdated(ctx context.Context, a *db.Account) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MocksPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Journal events ----------------
func (m *MocksPublisher) PublishJournalCreated(ctx context.Context, j *db.JournalEntry) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MocksPublisher) PublishJournalUpdated(ctx context.Context, j *db.JournalEntry) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MocksPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Accrual events ----------------
func (m *MocksPublisher) PublishAccrualCreated(ctx context.Context, a *db.Accrual) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MocksPublisher) PublishAccrualUpdated(ctx context.Context, a *db.Accrual) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MocksPublisher) PublishAccrualDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- AllocationRule events ----------------
func (m *MocksPublisher) PublishAllocationRuleCreated(ctx context.Context, r *db.AllocationRule) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MocksPublisher) PublishAllocationRuleUpdated(ctx context.Context, r *db.AllocationRule) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MocksPublisher) PublishAllocationRuleDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Audit events ----------------
func (m *MocksPublisher) PublishAuditRecorded(ctx context.Context, event *db.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ---------------- Budget events ----------------
func (m *MocksPublisher) PublishBudgetCreated(ctx context.Context, b *db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MocksPublisher) PublishBudgetUpdated(ctx context.Context, b *db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MocksPublisher) PublishBudgetDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MocksPublisher) PublishBudgetAllocated(ctx context.Context, ba *db.BudgetAllocation) error {
	args := m.Called(ctx, ba)
	return args.Error(0)
}

func (m *MocksPublisher) PublishBudgetAllocationUpdated(ctx context.Context, ba *db.BudgetAllocation) error {
	args := m.Called(ctx, ba)
	return args.Error(0)
}

func (m *MocksPublisher) PublishBudgetAllocationDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Cash Flow Forecast events ----------------
func (m *MocksPublisher) PublishCashFlowForecastGenerated(ctx context.Context, forecast *db.CashFlowForecast) error {
	args := m.Called(ctx, forecast)
	return args.Error(0)
}

func (m *MocksPublisher) PublishCashFlowForecastFetched(ctx context.Context, forecast *db.CashFlowForecast) error {
	args := m.Called(ctx, forecast)
	return args.Error(0)
}

func (m *MocksPublisher) PublishCashFlowForecastListed(ctx context.Context, forecasts []db.CashFlowForecast) error {
	args := m.Called(ctx, forecasts)
	return args.Error(0)
}

// ---------------- Consolidation events ----------------
func (m *MocksPublisher) PublishConsolidationCreated(ctx context.Context, c *db.Consolidation) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MocksPublisher) PublishConsolidationDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Credit/Debit Note events ----------------
func (m *MocksPublisher) PublishCreditDebitNoteCreated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MocksPublisher) PublishCreditDebitNoteUpdated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MocksPublisher) PublishCreditDebitNoteDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Exchange Rate events ----------------
func (m *MocksPublisher) PublishExchangeRateCreated(ctx context.Context, rate *db.ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MocksPublisher) PublishExchangeRateUpdated(ctx context.Context, rate *db.ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MocksPublisher) PublishExchangeRateDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Expense events ----------------
func (m *MocksPublisher) PublishExpenseCreated(ctx context.Context, exp *db.Expense) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MocksPublisher) PublishExpenseUpdated(ctx context.Context, exp *db.Expense) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MocksPublisher) PublishExpenseDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- CostCenter events ----------------
func (m *MocksPublisher) PublishCostCenterCreated(ctx context.Context, cc *db.CostCenter) error {
	args := m.Called(ctx, cc)
	return args.Error(0)
}

func (m *MocksPublisher) PublishCostCenterUpdated(ctx context.Context, cc *db.CostCenter) error {
	args := m.Called(ctx, cc)
	return args.Error(0)
}

func (m *MocksPublisher) PublishCostCenterDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- CostAllocation events ----------------
func (m *MocksPublisher) PublishCostAllocationAllocated(ctx context.Context, ca *db.CostAllocation) error {
	args := m.Called(ctx, ca)
	return args.Error(0)
}

func (m *MocksPublisher) PublishCostAllocationListed(ctx context.Context, allocs []db.CostAllocation) error {
	args := m.Called(ctx, allocs)
	return args.Error(0)
}

// ---------------- Finance domain events ----------------
func (m *MocksPublisher) PublishFinanceInvoiceCreated(ctx context.Context, ev *db.FinanceInvoiceCreatedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MocksPublisher) PublishFinancePaymentReceived(ctx context.Context, ev *db.FinancePaymentReceivedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MocksPublisher) PublishInventoryCostPosted(ctx context.Context, ev *db.InventoryCostPostedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MocksPublisher) PublishPayrollPosted(ctx context.Context, ev *db.PayrollPostedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MocksPublisher) PublishVendorBillApproved(ctx context.Context, ev *db.VendorBillApprovedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

// ---------------- GST events ----------------
func (m *MocksPublisher) PublishGstBreakupAdded(ctx context.Context, breakup *db.GstBreakup) error {
	args := m.Called(ctx, breakup)
	return args.Error(0)
}

func (m *MocksPublisher) PublishGstRegimeAdded(ctx context.Context, regime *db.GstRegime) error {
	args := m.Called(ctx, regime)
	return args.Error(0)
}

func (m *MocksPublisher) PublishGstDocStatusAdded(ctx context.Context, status *db.GstDocStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

// ---------------- Invoice events ----------------
func (m *MocksPublisher) PublishInvoiceUpdated(ctx context.Context, inv *db.Invoice) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MocksPublisher) PublishInvoiceDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MocksPublisher) PublishInvoiceItemCreated(ctx context.Context, item *db.InvoiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MocksPublisher) PublishInvoiceTaxAdded(ctx context.Context, tax *db.InvoiceTax) error {
	args := m.Called(ctx, tax)
	return args.Error(0)
}

func (m *MocksPublisher) PublishInvoiceDiscountAdded(ctx context.Context, disc *db.InvoiceDiscount) error {
	args := m.Called(ctx, disc)
	return args.Error(0)
}

// ---------------- Bank Transaction events ----------------
func (m *MocksPublisher) PublishBankTransactionImported(ctx context.Context, ev *db.BankAccount) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MocksPublisher) PublishBankTransactionReconciled(ctx context.Context, ev *db.BankAccount) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

// -------------------- TESTS --------------------

func TestBudgetService_CreateBudget(t *testing.T) {
    mockRepo := new(MockBudgetRepo)
    mockPub := new(MockPublisher)
    service := services.NewBudgetService(mockRepo, mockPub)

    ctx := context.Background()
    input := &db.Budget{Name: "New Budget", TotalAmount: "1000"}
    created := &db.Budget{ID: uuid.New(), Name: "New Budget", TotalAmount: "1000"}

    // Repo expectation
    mockRepo.On("Create", mock.Anything, input).Return(created, nil)

    // Publisher expectation (async publish → signal done)
    done := make(chan struct{})
    mockPub.On("PublishBudgetCreated", mock.Anything, created).
        Run(func(args mock.Arguments) { close(done) }).
        Return(nil)

    // Act
    got, err := service.CreateBudget(ctx, input)

    assert.NoError(t, err)
    assert.Equal(t, created, got)

    // Wait for async goroutine to fire
    <-done

    mockRepo.AssertExpectations(t)
    mockPub.AssertExpectations(t)
}


func TestBudgetService_GetBudget(t *testing.T) {
	mockRepo := new(MockBudgetRepo)
	service := services.NewBudgetService(mockRepo, nil)

	id := uuid.New()
	expected := &db.Budget{ID: id, Name: "IT Budget"}
	mockRepo.On("Get", mock.Anything, id).Return(expected, nil)

	got, err := service.GetBudget(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	mockRepo.AssertExpectations(t)
}

func TestBudgetService_ListBudgets(t *testing.T) {
	mockRepo := new(MockBudgetRepo)
	service := services.NewBudgetService(mockRepo, nil)

	expected := []db.Budget{{Name: "B1"}, {Name: "B2"}}
	mockRepo.On("List", mock.Anything, int32(5), int32(0)).Return(expected, nil)

	got, err := service.ListBudgets(context.Background(), 5, 0)

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	mockRepo.AssertExpectations(t)
}

func TestBudgetService_UpdateBudget(t *testing.T) {
    mockRepo := new(MockBudgetRepo)
    mockPub := new(MockPublisher)
    service := services.NewBudgetService(mockRepo, mockPub)

    ctx := context.Background()
    input := &db.Budget{ID: uuid.New(), Name: "Old Name", TotalAmount: "1000"}
    updated := &db.Budget{ID: input.ID, Name: "Updated", TotalAmount: "5000"}

    // Repo expectation
    mockRepo.On("Update", mock.Anything, input).Return(updated, nil)

    // Publisher expectation (async → use done channel)
    done := make(chan struct{})
    mockPub.On("PublishBudgetUpdated", mock.Anything, updated).
        Run(func(args mock.Arguments) { close(done) }).
        Return(nil)

    // Act
    got, err := service.UpdateBudget(ctx, input)

    assert.NoError(t, err)
    assert.Equal(t, updated, got)

    // Wait for async goroutine to fire
    <-done

    mockRepo.AssertExpectations(t)
    mockPub.AssertExpectations(t)
}


func TestBudgetService_DeleteBudget(t *testing.T) {
    mockRepo := new(MockBudgetRepo)
    mockPub := new(MockPublisher)
    service := services.NewBudgetService(mockRepo, mockPub)

    id := uuid.New()
    ctx := context.Background()

    // Repo expectation
    mockRepo.On("Delete", mock.Anything, id).Return(nil)

    // Publisher expectation (async) → signal done when called
    done := make(chan struct{})
    mockPub.On("PublishBudgetDeleted", mock.Anything, id.String()).
        Run(func(args mock.Arguments) { close(done) }).
        Return(nil)

    // Act
    err := service.DeleteBudget(ctx, id)
    assert.NoError(t, err)

    // Wait for async goroutine
    <-done

    mockRepo.AssertExpectations(t)
    mockPub.AssertExpectations(t)
}


func TestBudgetService_AllocateBudget(t *testing.T) {
    mockRepo := new(MockBudgetRepo)
    mockPub := new(MockPublisher)
    service := services.NewBudgetService(mockRepo, mockPub)

    ctx := context.Background()
    allocInput := &db.BudgetAllocation{
        ID:             uuid.New(),
        BudgetID:       uuid.New(),
        DepartmentID:   uuid.New().String(),
        AllocatedAmount:"1000",
    }
    allocResult := &db.BudgetAllocation{
        ID:             allocInput.ID,
        BudgetID:       allocInput.BudgetID,
        DepartmentID:   allocInput.DepartmentID,
        AllocatedAmount:"1000",
        SpentAmount:    sql.NullString{String:"200",Valid:true},
        RemainingAmount:sql.NullString{String:"800",Valid:true},
    }

    // Repo expectation
    mockRepo.On("Allocate", mock.Anything, allocInput).Return(allocResult, nil)

    // Publisher expectation (async → use done channel)
    done := make(chan struct{})
    mockPub.On("PublishBudgetAllocated", mock.Anything, allocResult).
        Run(func(args mock.Arguments) { close(done) }).
        Return(nil)

    // Act
    got, err := service.AllocateBudget(ctx, allocInput)

    assert.NoError(t, err)
    assert.Equal(t, allocResult, got)

    // Wait for async goroutine
    <-done

    mockRepo.AssertExpectations(t)
    mockPub.AssertExpectations(t)
}

func TestBudgetService_GetBudgetComparisonReport(t *testing.T) {
	mockRepo := new(MockBudgetRepo)
	service := services.NewBudgetService(mockRepo, nil)

	id := uuid.New()
	row := &db.GetBudgetComparisonReportRow{
		BudgetID:        id,
		TotalBudget:     "10000",
		TotalAllocated:  "6000",
		TotalSpent:      "4000",
		RemainingBudget: 6000,
	}

	mockRepo.On("GetBudgetComparisonReport", mock.Anything, id).Return(row, nil)

	got, err := service.GetBudgetComparisonReport(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, row, got)
	mockRepo.AssertExpectations(t)
}

func TestBudgetService_RepoError(t *testing.T) {
	mockRepo := new(MockBudgetRepo)
	service := services.NewBudgetService(mockRepo, nil)

	budget := &db.Budget{Name: "Fail"}
	mockRepo.On("Create", mock.Anything, budget).Return(&db.Budget{}, errors.New("db error"))

	_, err := service.CreateBudget(context.Background(), budget)
	assert.Error(t, err)
}
