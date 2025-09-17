package grpc_server_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

// ---------- Helpers ----------
func floatPtr(f float64) *float64    { return &f }
func boolPtr(b bool) *bool           { return &b }

type MockGstRepository struct {
	mock.Mock
}

func (m *MockGstRepository) AddGstBreakup(ctx context.Context, invoiceID uuid.UUID, taxable float64, cgst, sgst, igst, total *float64) (db.GstBreakup, error) {
	args := m.Called(ctx, invoiceID, taxable, cgst, sgst, igst, total)
	return args.Get(0).(db.GstBreakup), args.Error(1)
}

func (m *MockGstRepository) GetGstBreakup(ctx context.Context, invoiceID uuid.UUID) (db.GstBreakup, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstBreakup), args.Error(1)
}

func (m *MockGstRepository) AddGstRegime(ctx context.Context, invoiceID uuid.UUID, gstin, placeOfSupply string, reverseCharge *bool) (db.GstRegime, error) {
	args := m.Called(ctx, invoiceID, gstin, placeOfSupply, reverseCharge)
	return args.Get(0).(db.GstRegime), args.Error(1)
}

func (m *MockGstRepository) GetGstRegime(ctx context.Context, invoiceID uuid.UUID) (db.GstRegime, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstRegime), args.Error(1)
}

func (m *MockGstRepository) AddGstDocStatus(ctx context.Context, invoiceID uuid.UUID, einvoiceStatus, irn, ackNo *string, ackDate *time.Time,
	ewayStatus, ewayBillNo *string, ewayValidUpto *time.Time, lastError *string, lastSyncedAt *time.Time) (db.GstDocStatus, error) {
	args := m.Called(ctx, invoiceID, einvoiceStatus, irn, ackNo, ackDate, ewayStatus, ewayBillNo, ewayValidUpto, lastError, lastSyncedAt)
	return args.Get(0).(db.GstDocStatus), args.Error(1)
}

func (m *MockGstRepository) GetGstDocStatus(ctx context.Context, invoiceID uuid.UUID) (db.GstDocStatus, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(db.GstDocStatus), args.Error(1)
}

// ---------- Mocks ----------
// MockPublisher implements ports.EventPublisher for testing
type MockPublisher struct {
	mock.Mock
}

// Generic publish
func (m *MockPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	args := m.Called(ctx, topic, key, payload)
	return args.Error(0)
}

func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ---------------- Account events ----------------
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

// ---------------- Journal events ----------------
func (m *MockPublisher) PublishJournalCreated(ctx context.Context, j *db.JournalEntry) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockPublisher) PublishJournalUpdated(ctx context.Context, j *db.JournalEntry) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Accrual events ----------------
func (m *MockPublisher) PublishAccrualCreated(ctx context.Context, a *db.Accrual) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockPublisher) PublishAccrualUpdated(ctx context.Context, a *db.Accrual) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockPublisher) PublishAccrualDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- AllocationRule events ----------------
func (m *MockPublisher) PublishAllocationRuleCreated(ctx context.Context, r *db.AllocationRule) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockPublisher) PublishAllocationRuleUpdated(ctx context.Context, r *db.AllocationRule) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockPublisher) PublishAllocationRuleDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Audit events ----------------
func (m *MockPublisher) PublishAuditRecorded(ctx context.Context, event *db.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ---------------- Budget events ----------------
func (m *MockPublisher) PublishBudgetCreated(ctx context.Context, b *db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockPublisher) PublishBudgetUpdated(ctx context.Context, b *db.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockPublisher) PublishBudgetDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPublisher) PublishBudgetAllocated(ctx context.Context, ba *db.BudgetAllocation) error {
	args := m.Called(ctx, ba)
	return args.Error(0)
}

func (m *MockPublisher) PublishBudgetAllocationUpdated(ctx context.Context, ba *db.BudgetAllocation) error {
	args := m.Called(ctx, ba)
	return args.Error(0)
}

func (m *MockPublisher) PublishBudgetAllocationDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Cash Flow Forecast events ----------------
func (m *MockPublisher) PublishCashFlowForecastGenerated(ctx context.Context, forecast *db.CashFlowForecast) error {
	args := m.Called(ctx, forecast)
	return args.Error(0)
}

func (m *MockPublisher) PublishCashFlowForecastFetched(ctx context.Context, forecast *db.CashFlowForecast) error {
	args := m.Called(ctx, forecast)
	return args.Error(0)
}

func (m *MockPublisher) PublishCashFlowForecastListed(ctx context.Context, forecasts []db.CashFlowForecast) error {
	args := m.Called(ctx, forecasts)
	return args.Error(0)
}

// ---------------- Consolidation events ----------------
func (m *MockPublisher) PublishConsolidationCreated(ctx context.Context, c *db.Consolidation) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockPublisher) PublishConsolidationDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Credit/Debit Note events ----------------
func (m *MockPublisher) PublishCreditDebitNoteCreated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockPublisher) PublishCreditDebitNoteUpdated(ctx context.Context, note *db.CreditDebitNote) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockPublisher) PublishCreditDebitNoteDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Exchange Rate events ----------------
func (m *MockPublisher) PublishExchangeRateCreated(ctx context.Context, rate *db.ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockPublisher) PublishExchangeRateUpdated(ctx context.Context, rate *db.ExchangeRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockPublisher) PublishExchangeRateDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- Expense events ----------------
func (m *MockPublisher) PublishExpenseCreated(ctx context.Context, exp *db.Expense) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockPublisher) PublishExpenseUpdated(ctx context.Context, exp *db.Expense) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockPublisher) PublishExpenseDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- CostCenter events ----------------
func (m *MockPublisher) PublishCostCenterCreated(ctx context.Context, cc *db.CostCenter) error {
	args := m.Called(ctx, cc)
	return args.Error(0)
}

func (m *MockPublisher) PublishCostCenterUpdated(ctx context.Context, cc *db.CostCenter) error {
	args := m.Called(ctx, cc)
	return args.Error(0)
}

func (m *MockPublisher) PublishCostCenterDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------------- CostAllocation events ----------------
func (m *MockPublisher) PublishCostAllocationAllocated(ctx context.Context, ca *db.CostAllocation) error {
	args := m.Called(ctx, ca)
	return args.Error(0)
}

func (m *MockPublisher) PublishCostAllocationListed(ctx context.Context, allocs []db.CostAllocation) error {
	args := m.Called(ctx, allocs)
	return args.Error(0)
}

// ---------------- Finance domain events ----------------
func (m *MockPublisher) PublishFinanceInvoiceCreated(ctx context.Context, ev *db.FinanceInvoiceCreatedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockPublisher) PublishFinancePaymentReceived(ctx context.Context, ev *db.FinancePaymentReceivedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockPublisher) PublishInventoryCostPosted(ctx context.Context, ev *db.InventoryCostPostedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockPublisher) PublishPayrollPosted(ctx context.Context, ev *db.PayrollPostedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockPublisher) PublishVendorBillApproved(ctx context.Context, ev *db.VendorBillApprovedEvent) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

// ---------------- GST events ----------------
func (m *MockPublisher) PublishGstBreakupAdded(ctx context.Context, breakup *db.GstBreakup) error {
	args := m.Called(ctx, breakup)
	return args.Error(0)
}

func (m *MockPublisher) PublishGstRegimeAdded(ctx context.Context, regime *db.GstRegime) error {
	args := m.Called(ctx, regime)
	return args.Error(0)
}

func (m *MockPublisher) PublishGstDocStatusAdded(ctx context.Context, status *db.GstDocStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

// ---------------- Invoice events ----------------
func (m *MockPublisher) PublishInvoiceUpdated(ctx context.Context, inv *db.Invoice) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockPublisher) PublishInvoiceDeleted(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPublisher) PublishInvoiceItemCreated(ctx context.Context, item *db.InvoiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockPublisher) PublishInvoiceTaxAdded(ctx context.Context, tax *db.InvoiceTax) error {
	args := m.Called(ctx, tax)
	return args.Error(0)
}

func (m *MockPublisher) PublishInvoiceDiscountAdded(ctx context.Context, disc *db.InvoiceDiscount) error {
	args := m.Called(ctx, disc)
	return args.Error(0)
}

// ---------------- Bank Transaction events ----------------
func (m *MockPublisher) PublishBankTransactionImported(ctx context.Context, ev *db.BankAccount) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func (m *MockPublisher) PublishBankTransactionReconciled(ctx context.Context, ev *db.BankAccount) error {
	args := m.Called(ctx, ev)
	return args.Error(0)
}

func setupTestHandler(t *testing.T) (*grpc_server.GstHandler, *MockGstRepository, *MockPublisher) {
    _ = t // ignore unused parameter warning

    mockRepo := new(MockGstRepository)
    mockPub := new(MockPublisher)
    svc := services.NewGstService(mockRepo, mockPub)
    handler := grpc_server.NewGstHandler(svc)
    return handler, mockRepo, mockPub
}

// ---------- Tests ----------

func TestGstHandler_AddAndGetBreakup(t *testing.T) {
	handler, mockRepo, mockPub := setupTestHandler(t)
	invoiceID := uuid.New()

	exp := db.GstBreakup{
		ID:            uuid.New(),
		InvoiceID:     invoiceID,
		TaxableAmount: "1000",
		Cgst:          sql.NullString{String: "50", Valid: true},
		Sgst:          sql.NullString{String: "50", Valid: true},
		Igst:          sql.NullString{String: "0", Valid: true},
		TotalGst:      sql.NullString{String: "100", Valid: true},
	}

	mockRepo.On("AddGstBreakup", mock.Anything, invoiceID, 1000.0, floatPtr(50), floatPtr(50), floatPtr(0), floatPtr(100)).Return(exp, nil)
	mockRepo.On("GetGstBreakup", mock.Anything, invoiceID).Return(exp, nil)
	mockPub.On("PublishGstBreakupAdded", mock.Anything, mock.AnythingOfType("*db.GstBreakup")).Return(nil)

	// --- Add POST ---
	body, _ := json.Marshal(map[string]interface{}{
		"invoice_id":     invoiceID.String(),
		"taxable_amount": 1000,
		"cgst":           50,
		"sgst":           50,
		"igst":           0,
		"total_gst":      100,
	})
	req := httptest.NewRequest(http.MethodPost, "/gst/breakup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.AddGstBreakup(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var resp db.GstBreakup
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "1000", resp.TaxableAmount)
	assert.Equal(t, "50", resp.Cgst.String)

	// --- GET ---
	reqGet := httptest.NewRequest(http.MethodGet, "/gst/breakup/"+invoiceID.String(), nil)
	rw := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("invoice_id", invoiceID.String())
	reqGet = reqGet.WithContext(context.WithValue(reqGet.Context(), chi.RouteCtxKey, chiCtx))
	handler.GetGstBreakup(rw, reqGet)
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)

	var getResp db.GstBreakup
	_ = json.NewDecoder(rw.Body).Decode(&getResp)
	assert.Equal(t, "1000", getResp.TaxableAmount)
}

func TestGstHandler_AddAndGetRegime(t *testing.T) {
	handler, mockRepo, mockPub := setupTestHandler(t)
	invoiceID := uuid.New()

	exp := db.GstRegime{
		ID:            uuid.New(),
		InvoiceID:     invoiceID,
		Gstin:         "27ABCDE1234F2Z5",
		PlaceOfSupply: "MH",
		ReverseCharge: sql.NullBool{Bool: false, Valid: true},
	}

	mockRepo.On("AddGstRegime", mock.Anything, invoiceID, exp.Gstin, exp.PlaceOfSupply, boolPtr(false)).Return(exp, nil)
	mockRepo.On("GetGstRegime", mock.Anything, invoiceID).Return(exp, nil)
	mockPub.On("PublishGstRegimeAdded", mock.Anything, mock.AnythingOfType("*db.GstRegime")).Return(nil)

	body, _ := json.Marshal(map[string]interface{}{
		"invoice_id":      invoiceID.String(),
		"gstin":           exp.Gstin,
		"place_of_supply": exp.PlaceOfSupply,
		"reverse_charge":  false,
	})
	req := httptest.NewRequest(http.MethodPost, "/gst/regime", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.AddGstRegime(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var resp db.GstRegime
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, exp.Gstin, resp.Gstin)

	// --- GET ---
	reqGet := httptest.NewRequest(http.MethodGet, "/gst/regime/"+invoiceID.String(), nil)
	rw := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("invoice_id", invoiceID.String())
	reqGet = reqGet.WithContext(context.WithValue(reqGet.Context(), chi.RouteCtxKey, chiCtx))
	handler.GetGstRegime(rw, reqGet)
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
}

func TestGstHandler_AddAndGetDocStatus(t *testing.T) {
	handler, mockRepo, mockPub := setupTestHandler(t)
	invoiceID := uuid.New() // Keep as UUID
	now := time.Now()

	exp := db.GstDocStatus{
		ID:             uuid.New(),
		InvoiceID:      invoiceID,
		EinvoiceStatus: sql.NullString{String: "SUCCESS", Valid: true},
		Irn:            sql.NullString{String: "IRN123", Valid: true},
		AckNo:          sql.NullString{String: "ACK456", Valid: true},
		AckDate:        sql.NullTime{Time: now, Valid: true},
		EwayStatus:     sql.NullString{String: "VALID", Valid: true},
		EwayBillNo:     sql.NullString{String: "EB123", Valid: true},
		EwayValidUpto:  sql.NullTime{Time: now, Valid: true},
	}

	// --- Mock repository call with correct types ---
	mockRepo.On(
		"AddGstDocStatus",
		mock.Anything, // context
		invoiceID,     // UUID
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*time.Time"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("*time.Time"),
		(*string)(nil),
		(*time.Time)(nil),
	).Return(exp, nil)

	mockRepo.On("GetGstDocStatus", mock.Anything, invoiceID).Return(exp, nil)
	mockPub.On(
		"PublishGstDocStatusAdded",
		mock.Anything,
		mock.AnythingOfType("*db.GstDocStatus"),
	).Return(nil)

	// --- POST request ---
	body, _ := json.Marshal(map[string]interface{}{
		"invoice_id":      invoiceID.String(),
		"einvoice_status": "SUCCESS",
		"irn":             "IRN123",
		"ack_no":          "ACK456",
		"ack_date":        now.Format(time.RFC3339),
		"eway_status":     "VALID",
		"eway_bill_no":    "EB123",
		"eway_valid_upto": now.Format(time.RFC3339),
	})

	req := httptest.NewRequest(http.MethodPost, "/gst/doc-status", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.AddGstDocStatus(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var resp db.GstDocStatus
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "SUCCESS", resp.EinvoiceStatus.String)

	// --- GET request ---
	reqGet := httptest.NewRequest(http.MethodGet, "/gst/doc-status/"+invoiceID.String(), nil)
	rw := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("invoice_id", invoiceID.String())
	reqGet = reqGet.WithContext(context.WithValue(reqGet.Context(), chi.RouteCtxKey, chiCtx))
	handler.GetGstDocStatus(rw, reqGet)

	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
}