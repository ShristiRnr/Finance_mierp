package ports

import (
	"context"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

// EventPublisher defines an abstraction for publishing domain events.
type EventPublisher interface {
	// Generic publish
	Publish(ctx context.Context, topic, key string, payload []byte) error
	Close() error

	// ---------------- Account events ----------------
	PublishAccountCreated(ctx context.Context, a *db.Account) error
	PublishAccountUpdated(ctx context.Context, a *db.Account) error
	PublishAccountDeleted(ctx context.Context, id string) error

	// ---------------- Journal events ----------------
	PublishJournalCreated(ctx context.Context, j *db.JournalEntry) error
	PublishJournalUpdated(ctx context.Context, j *db.JournalEntry) error
	PublishJournalDeleted(ctx context.Context, id string) error

	// ---------------- Accrual events ----------------
	PublishAccrualCreated(ctx context.Context, a *db.Accrual) error
	PublishAccrualUpdated(ctx context.Context, a *db.Accrual) error
	PublishAccrualDeleted(ctx context.Context, id string) error

	// ---------------- AllocationRule events ----------------
	PublishAllocationRuleCreated(ctx context.Context, r *db.AllocationRule) error
	PublishAllocationRuleUpdated(ctx context.Context, r *db.AllocationRule) error
	PublishAllocationRuleDeleted(ctx context.Context, id string) error

	// ---------------- Audit events ----------------
	PublishAuditRecorded(ctx context.Context, event *db.AuditEvent) error

	// ---------------- Budget events ----------------
	PublishBudgetCreated(ctx context.Context, b *db.Budget) error
	PublishBudgetUpdated(ctx context.Context, b *db.Budget) error
	PublishBudgetDeleted(ctx context.Context, id string) error
	PublishBudgetAllocated(ctx context.Context, ba *db.BudgetAllocation) error
	PublishBudgetAllocationUpdated(ctx context.Context, ba *db.BudgetAllocation) error
	PublishBudgetAllocationDeleted(ctx context.Context, id string) error

	// ---------------- Cash Flow Forecast events ----------------
	PublishCashFlowForecastGenerated(ctx context.Context, forecast *db.CashFlowForecast) error
	PublishCashFlowForecastFetched(ctx context.Context, forecast *db.CashFlowForecast) error
	PublishCashFlowForecastListed(ctx context.Context, forecasts []db.CashFlowForecast) error

	// ---------------- Consolidation events ----------------
	PublishConsolidationCreated(ctx context.Context, c *db.Consolidation) error
	PublishConsolidationDeleted(ctx context.Context, id string) error

	// ---------------- Credit/Debit Note events ----------------
	PublishCreditDebitNoteCreated(ctx context.Context, note *db.CreditDebitNote) error
	PublishCreditDebitNoteUpdated(ctx context.Context, note *db.CreditDebitNote) error
	PublishCreditDebitNoteDeleted(ctx context.Context, id string) error

	// ---------------- Exchange Rate events ----------------
	PublishExchangeRateCreated(ctx context.Context, rate *db.ExchangeRate) error
	PublishExchangeRateUpdated(ctx context.Context, rate *db.ExchangeRate) error
	PublishExchangeRateDeleted(ctx context.Context, id string) error

	// ---------------- Expense events ----------------
	PublishExpenseCreated(ctx context.Context, exp *db.Expense) error
	PublishExpenseUpdated(ctx context.Context, exp *db.Expense) error
	PublishExpenseDeleted(ctx context.Context, id string) error

	// ---------------- CostCenter events ----------------
	PublishCostCenterCreated(ctx context.Context, cc *db.CostCenter) error
	PublishCostCenterUpdated(ctx context.Context, cc *db.CostCenter) error
	PublishCostCenterDeleted(ctx context.Context, id string) error

	// ---------------- CostAllocation events ----------------
	PublishCostAllocationAllocated(ctx context.Context, ca *db.CostAllocation) error
	PublishCostAllocationListed(ctx context.Context, allocs []db.CostAllocation) error

	// ---------------- Finance domain events ----------------
	PublishFinanceInvoiceCreated(ctx context.Context, ev *db.FinanceInvoiceCreatedEvent) error
	PublishFinancePaymentReceived(ctx context.Context, ev *db.FinancePaymentReceivedEvent) error
	PublishInventoryCostPosted(ctx context.Context, ev *db.InventoryCostPostedEvent) error
	PublishPayrollPosted(ctx context.Context, ev *db.PayrollPostedEvent) error
	PublishVendorBillApproved(ctx context.Context, ev *db.VendorBillApprovedEvent) error

	// ---------------- GST events ----------------
	PublishGstBreakupAdded(ctx context.Context, breakup *db.GstBreakup) error
	PublishGstRegimeAdded(ctx context.Context, regime *db.GstRegime) error
	PublishGstDocStatusAdded(ctx context.Context, status *db.GstDocStatus) error

	// ---------------- Invoice events ----------------
	PublishInvoiceUpdated(ctx context.Context, inv *db.Invoice) error
	PublishInvoiceDeleted(ctx context.Context, id string) error
	PublishInvoiceItemCreated(ctx context.Context, item *db.InvoiceItem) error
	PublishInvoiceTaxAdded(ctx context.Context, tax *db.InvoiceTax) error
	PublishInvoiceDiscountAdded(ctx context.Context, disc *db.InvoiceDiscount) error

	// ---------------- Bank Transaction events ----------------
	PublishBankTransactionImported(ctx context.Context, ev *db.BankAccount) error
	PublishBankTransactionReconciled(ctx context.Context, ev *db.BankAccount) error
	
}
