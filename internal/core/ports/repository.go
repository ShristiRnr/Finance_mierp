package ports

import (
	"context"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/google/uuid"
)

// Accounts
type AccountRepository interface {
	Create(ctx context.Context, a domain.Account) (domain.Account, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Account, error)
	Update(ctx context.Context, a domain.Account) (domain.Account, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]domain.Account, error)
}

// Journals
type JournalRepository interface {
	Create(ctx context.Context, j domain.JournalEntry) (domain.JournalEntry, error)
	Get(ctx context.Context, id uuid.UUID) (domain.JournalEntry, error)
	Update(ctx context.Context, j domain.JournalEntry) (domain.JournalEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]domain.JournalEntry, error)
}

// Ledger (read-only projection)
type LedgerRepository interface {
	List(ctx context.Context, limit, offset int32) ([]domain.LedgerEntry, error)
}

// Accruals
type AccrualRepository interface {
	Create(ctx context.Context, a domain.Accrual) (domain.Accrual, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Accrual, error)
	Update(ctx context.Context, a domain.Accrual) (domain.Accrual, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]domain.Accrual, error)
}

// Allocation
type AllocationRuleRepository interface {
	Create(ctx context.Context, r domain.AllocationRule) (domain.AllocationRule, error)
	Get(ctx context.Context, id uuid.UUID) (domain.AllocationRule, error)
	Update(ctx context.Context, r domain.AllocationRule) (domain.AllocationRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]domain.AllocationRule, error)
}

// AllocationService defines the business methods for allocation rules
type AllocationService interface {
	CreateRule(ctx context.Context, r domain.AllocationRule) (domain.AllocationRule, error)
	GetRule(ctx context.Context, id uuid.UUID) (domain.AllocationRule, error)
	ListRules(ctx context.Context, limit, offset int32) ([]domain.AllocationRule, error)
	UpdateRule(ctx context.Context, r domain.AllocationRule) (domain.AllocationRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error
	ApplyRule(ctx context.Context, ruleID uuid.UUID) error
}

// AuditRepository is the port for database interactions related to audit events.
type AuditRepository interface {
	RecordAuditEvent(ctx context.Context, event *domain.AuditEvent) (*domain.AuditEvent, error)
	ListAuditEvents(ctx context.Context, page domain.Pagination) ([]domain.AuditEvent, error)
	GetAuditEventByID(ctx context.Context, id string) (*domain.AuditEvent, error)
	FilterAuditEvents(ctx context.Context, filter domain.FilterParams, page domain.Pagination) ([]domain.AuditEvent, error)
}

// BudgetRepository is the port for database interactions related to budgets.
type BudgetRepository interface {
	Create(ctx context.Context, b *domain.Budget) (*domain.Budget, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Budget, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Budget, error)
	Update(ctx context.Context, b *domain.Budget) (*domain.Budget, error)
	Delete(ctx context.Context, id uuid.UUID) error

	Allocate(ctx context.Context, ba *domain.BudgetAllocation) (*domain.BudgetAllocation, error)
	GetAllocation(ctx context.Context, id uuid.UUID) (*domain.BudgetAllocation, error)
	ListAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]domain.BudgetAllocation, error)
	UpdateAllocation(ctx context.Context, ba *domain.BudgetAllocation) (*domain.BudgetAllocation, error)
	DeleteAllocation(ctx context.Context, id uuid.UUID) error

	GetBudgetComparison(ctx context.Context, id uuid.UUID) (*domain.BudgetComparisonReport, error)
}

type CashFlowForecastRepository interface {
	Generate(ctx context.Context, cf *domain.CashFlowForecast) (*domain.CashFlowForecast, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.CashFlowForecast, error)
	List(ctx context.Context, organizationID string, limit, offset int32) ([]*domain.CashFlowForecast, error)
}

type ConsolidationRepository interface {
	Create(ctx context.Context, c domain.Consolidation) (domain.Consolidation, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Consolidation, error)
	List(ctx context.Context, entityIds []string, start, end time.Time, limit, offset int32) ([]domain.Consolidation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreditDebitNoteRepository interface {
	Create(ctx context.Context, note domain.CreditDebitNote) (domain.CreditDebitNote, error)
	Get(ctx context.Context, id uuid.UUID) (domain.CreditDebitNote, error)
	List(ctx context.Context, limit, offset int32) ([]domain.CreditDebitNote, error)
	Update(ctx context.Context, note domain.CreditDebitNote) (domain.CreditDebitNote, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreditDebitNoteService interface {
	Create(ctx context.Context, note domain.CreditDebitNote) (domain.CreditDebitNote, error)
	Get(ctx context.Context, id uuid.UUID) (domain.CreditDebitNote, error)
	List(ctx context.Context, limit, offset int32) ([]domain.CreditDebitNote, error)
	Update(ctx context.Context, note domain.CreditDebitNote) (domain.CreditDebitNote, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ExchangeRateRepository interface {
	Create(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error)
	Get(ctx context.Context, id uuid.UUID) (domain.ExchangeRate, error)
	Update(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, base, quote *string, limit, offset int32) ([]domain.ExchangeRate, error)
	GetLatest(ctx context.Context, base, quote string, asOf time.Time) (domain.ExchangeRate, error)
}

type ExchangeRateService interface {
	Create(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error)
	Get(ctx context.Context, id uuid.UUID) (domain.ExchangeRate, error)
	Update(ctx context.Context, rate domain.ExchangeRate) (domain.ExchangeRate, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, base, quote *string, limit, offset int32) ([]domain.ExchangeRate, error)
	GetLatest(ctx context.Context, base, quote string, asOf time.Time) (domain.ExchangeRate, error)
}

type ExpenseRepository interface {
	Create(ctx context.Context, expense domain.Expense) (domain.Expense, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Expense, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Expense, error)
	Update(ctx context.Context, expense domain.Expense) (domain.Expense, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CostCenterRepository interface {
	Create(ctx context.Context, cc domain.CostCenter) (domain.CostCenter, error)
	Get(ctx context.Context, id uuid.UUID) (domain.CostCenter, error)
	List(ctx context.Context, limit, offset int32) ([]domain.CostCenter, error)
	Update(ctx context.Context, cc domain.CostCenter) (domain.CostCenter, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CostAllocationRepository interface {
	Allocate(ctx context.Context, ca domain.CostAllocation) (domain.CostAllocation, error)
	List(ctx context.Context, limit, offset int32) ([]domain.CostAllocation, error)
}

type ExpenseService interface {
	CreateExpense(ctx context.Context, exp domain.Expense) (domain.Expense, error)
	GetExpense(ctx context.Context, id uuid.UUID) (domain.Expense, error)
	ListExpenses(ctx context.Context, limit, offset int32) ([]domain.Expense, error)
	UpdateExpense(ctx context.Context, exp domain.Expense) (domain.Expense, error)
	DeleteExpense(ctx context.Context, id uuid.UUID) error
}

type CostCenterService interface {
	CreateCostCenter(ctx context.Context, cc domain.CostCenter) (domain.CostCenter, error)
	GetCostCenter(ctx context.Context, id uuid.UUID) (domain.CostCenter, error)
	ListCostCenters(ctx context.Context, limit, offset int32) ([]domain.CostCenter, error)
	UpdateCostCenter(ctx context.Context, cc domain.CostCenter) (domain.CostCenter, error)
	DeleteCostCenter(ctx context.Context, id uuid.UUID) error
}

type CostAllocationService interface {
	AllocateCost(ctx context.Context, ca domain.CostAllocation) (domain.CostAllocation, error)
	ListAllocations(ctx context.Context, limit, offset int32) ([]domain.CostAllocation, error)
}

type FinanceEventRepository interface {
	InsertInvoiceCreated(ctx context.Context, e domain.FinanceInvoiceCreatedEvent) (domain.FinanceInvoiceCreatedEvent, error)
	ListInvoiceCreated(ctx context.Context, orgID string, limit, offset int32) ([]domain.FinanceInvoiceCreatedEvent, error)

	InsertPaymentReceived(ctx context.Context, e domain.FinancePaymentReceivedEvent) (domain.FinancePaymentReceivedEvent, error)
	ListPaymentReceived(ctx context.Context, orgID string, limit, offset int32) ([]domain.FinancePaymentReceivedEvent, error)

	InsertInventoryCostPosted(ctx context.Context, e domain.InventoryCostPostedEvent) (domain.InventoryCostPostedEvent, error)
	ListInventoryCostPosted(ctx context.Context, orgID string, limit, offset int32) ([]domain.InventoryCostPostedEvent, error)

	InsertPayrollPosted(ctx context.Context, e domain.PayrollPostedEvent) (domain.PayrollPostedEvent, error)
	ListPayrollPosted(ctx context.Context, orgID string, limit, offset int32) ([]domain.PayrollPostedEvent, error)

	InsertVendorBillApproved(ctx context.Context, e domain.VendorBillApprovedEvent) (domain.VendorBillApprovedEvent, error)
	ListVendorBillApproved(ctx context.Context, orgID string, limit, offset int32) ([]domain.VendorBillApprovedEvent, error)
}

type FinancialReportsRepository interface {
	// Profit & Loss
	GenerateProfitLossReport(ctx context.Context, report domain.ProfitLossReport) (domain.ProfitLossReport, error)
	GetProfitLossReport(ctx context.Context, id uuid.UUID) (domain.ProfitLossReport, error)
	ListProfitLossReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.ProfitLossReport, error)

	// Balance Sheet
	GenerateBalanceSheetReport(ctx context.Context, report domain.BalanceSheetReport) (domain.BalanceSheetReport, error)
	GetBalanceSheetReport(ctx context.Context, id uuid.UUID) (domain.BalanceSheetReport, error)
	ListBalanceSheetReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.BalanceSheetReport, error)

	// Trial Balance
	CreateTrialBalanceReport(ctx context.Context, report domain.TrialBalanceReport) (domain.TrialBalanceReport, error)
	AddTrialBalanceEntry(ctx context.Context, entry domain.TrialBalanceEntry) (domain.TrialBalanceEntry, error)
	GetTrialBalanceReport(ctx context.Context, id uuid.UUID) (domain.TrialBalanceReport, error)
	ListTrialBalanceReports(ctx context.Context, orgID string, limit, offset int32) ([]domain.TrialBalanceReport, error)
	ListTrialBalanceEntries(ctx context.Context, reportID uuid.UUID) ([]domain.TrialBalanceEntry, error)

	// Compliance
	GenerateComplianceReport(ctx context.Context, report domain.ComplianceReport) (domain.ComplianceReport, error)
	GetComplianceReport(ctx context.Context, id uuid.UUID) (domain.ComplianceReport, error)
	ListComplianceReports(ctx context.Context, orgID, jurisdiction string, limit, offset int32) ([]domain.ComplianceReport, error)
}

type GstRepository interface {
	// Breakup
	AddGstBreakup(ctx context.Context, invoiceID uuid.UUID, taxableAmount float64, cgst, sgst, igst, totalGst *float64) (domain.GstBreakup, error)
	GetGstBreakup(ctx context.Context, invoiceID uuid.UUID) (domain.GstBreakup, error)

	// Regime
	AddGstRegime(ctx context.Context, invoiceID uuid.UUID, gstin, placeOfSupply string, reverseCharge *bool) (domain.GstRegime, error)
	GetGstRegime(ctx context.Context, invoiceID uuid.UUID) (domain.GstRegime, error)

	// Doc Status
	AddGstDocStatus(ctx context.Context, invoiceID uuid.UUID, einvoiceStatus, irn, ackNo *string, ackDate *time.Time, ewayStatus, ewayBillNo *string, ewayValidUpto *time.Time,
		lastError *string, lastSyncedAt *time.Time) (domain.GstDocStatus, error)

	GetGstDocStatus(ctx context.Context, invoiceID uuid.UUID) (domain.GstDocStatus, error)
}

type InvoiceRepository interface {
    CreateInvoice(ctx context.Context, inv domain.Invoice) (domain.Invoice, error)
    UpdateInvoice(ctx context.Context, inv domain.Invoice) (domain.Invoice, error)
    DeleteInvoice(ctx context.Context, id uuid.UUID) error
    GetInvoice(ctx context.Context, id uuid.UUID) (domain.Invoice, error)
    ListInvoices(ctx context.Context, limit, offset int32) ([]domain.Invoice, error)
    SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]domain.Invoice, error)

    CreateInvoiceItem(ctx context.Context, item domain.InvoiceItem) (domain.InvoiceItem, error)
    ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]domain.InvoiceItem, error)

    AddInvoiceTax(ctx context.Context, tax domain.InvoiceTax) (domain.InvoiceTax, error)
    AddInvoiceDiscount(ctx context.Context, discount domain.InvoiceDiscount) (domain.InvoiceDiscount, error)
}

type BankAccountRepository interface {
    CreateBankAccount(ctx context.Context, ba domain.BankAccount) (domain.BankAccount, error)
    GetBankAccount(ctx context.Context, id uuid.UUID) (domain.BankAccount, error)
    UpdateBankAccount(ctx context.Context, ba domain.BankAccount) (domain.BankAccount, error)
    DeleteBankAccount(ctx context.Context, id uuid.UUID) error
    ListBankAccounts(ctx context.Context, limit, offset int32) ([]domain.BankAccount, error)
}


type PaymentDueRepository interface {
    CreatePaymentDue(ctx context.Context, pd domain.PaymentDue) (domain.PaymentDue, error)
    GetPaymentDue(ctx context.Context, id uuid.UUID) (domain.PaymentDue, error)
    UpdatePaymentDue(ctx context.Context, pd domain.PaymentDue) (domain.PaymentDue, error)
    DeletePaymentDue(ctx context.Context, id uuid.UUID) error
    ListPaymentDues(ctx context.Context, limit, offset int32) ([]domain.PaymentDue, error)
    MarkPaymentAsPaid(ctx context.Context, id uuid.UUID, updatedBy string) (domain.PaymentDue, error)
}

type BankTransactionRepository interface {
    ImportBankTransaction(ctx context.Context, tx domain.BankTransaction) (domain.BankTransaction, error)
    ListBankTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int32) ([]domain.BankTransaction, error)
    ReconcileTransaction(ctx context.Context, tx domain.BankTransaction) (domain.BankTransaction, error)
}
