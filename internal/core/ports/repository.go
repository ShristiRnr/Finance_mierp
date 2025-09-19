package ports

import (
	"context"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/google/uuid"
	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
)

// Accounts
type AccountRepository interface {
	Create(ctx context.Context, a *db.Account) (*db.Account, error)
	Get(ctx context.Context, id uuid.UUID) (*db.Account, error)
	Update(ctx context.Context, a *db.Account) (*db.Account, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]*db.Account, error)
}

// Journals
type JournalRepository interface {
	Create(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error)
	Get(ctx context.Context, id uuid.UUID) (*db.JournalEntry, error)
	Update(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]*db.JournalEntry, error)
}

// Ledger (read-only projection)
type LedgerRepository interface {
	List(ctx context.Context, limit, offset int32) ([]*db.LedgerEntry, error)
}

// Accounts
type AccountService interface {
    Create(ctx context.Context, a *db.Account) (*db.Account, error)
    Get(ctx context.Context, id uuid.UUID) (*db.Account, error)
    Update(ctx context.Context, a *db.Account) (*db.Account, error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, limit, offset int32) ([]*db.Account, error)
}

// Journals
type JournalService interface {
    Create(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error)
    Get(ctx context.Context, id uuid.UUID) (*db.JournalEntry, error)
    Update(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, limit, offset int32) ([]*db.JournalEntry, error)
}

// Ledger (read-only projection)
type LedgerService interface {
    List(ctx context.Context, limit, offset int32) ([]*db.LedgerEntry, error)
}

// Accruals
type AccrualRepository interface {
	Create(ctx context.Context, a db.Accrual) (db.Accrual, error)
	Get(ctx context.Context, id uuid.UUID) (db.Accrual, error)
	Update(ctx context.Context, a db.Accrual) (db.Accrual, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]db.Accrual, error)
}

type AccrualService interface {
	Create(ctx context.Context, a db.Accrual) (db.Accrual, error)
	Get(ctx context.Context, id uuid.UUID) (db.Accrual, error)
	Update(ctx context.Context, a db.Accrual) (db.Accrual, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]db.Accrual, error)
}

// Allocation
type AllocationRuleRepository interface {
	Create(ctx context.Context, r db.AllocationRule) (db.AllocationRule, error)
	Get(ctx context.Context, id uuid.UUID) (db.AllocationRule, error)
	Update(ctx context.Context, r db.AllocationRule) (db.AllocationRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]db.AllocationRule, error)
}

// AllocationService defines the business methods for allocation rules
type AllocationService interface {
	CreateRule(ctx context.Context, r db.AllocationRule) (db.AllocationRule, error)
	GetRule(ctx context.Context, id uuid.UUID) (db.AllocationRule, error)
	ListRules(ctx context.Context, limit, offset int32) ([]db.AllocationRule, error)
	UpdateRule(ctx context.Context, r db.AllocationRule) (db.AllocationRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error
}

// AuditRepository is the port for database interactions related to audit events.
type AuditRepository interface {
	RecordAuditEvent(ctx context.Context, event *db.AuditEvent) (*db.AuditEvent, error)
	ListAuditEvents(ctx context.Context, page db.Pagination) ([]db.AuditEvent, error)
	GetAuditEventByID(ctx context.Context, id uuid.UUID) (*db.AuditEvent, error)
	FilterAuditEvents(ctx context.Context, filter db.FilterParams, page db.Pagination) ([]db.AuditEvent, error)
}

type AuditService interface {
	Record(ctx context.Context, event *db.AuditEvent) (*db.AuditEvent, error)
	List(ctx context.Context, page db.Pagination) ([]db.AuditEvent, error)
	GetByID(ctx context.Context, id string) (*db.AuditEvent, error)
	Filter(ctx context.Context, filter db.FilterParams, page db.Pagination) ([]db.AuditEvent, error)
}

// BudgetRepository is the port for database interactions related to budgets.
type BudgetRepository interface {
	Create(ctx context.Context, b *db.Budget) (*db.Budget, error)
	Get(ctx context.Context, id uuid.UUID) (*db.Budget, error)
	List(ctx context.Context, limit, offset int32) ([]db.Budget, error)
	Update(ctx context.Context, b *db.Budget) (*db.Budget, error)
	Delete(ctx context.Context, id uuid.UUID) error

	Allocate(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error)
	GetAllocation(ctx context.Context, id uuid.UUID) (*db.BudgetAllocation, error)
	ListAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]db.BudgetAllocation, error)
	UpdateAllocation(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error)
	DeleteAllocation(ctx context.Context, id uuid.UUID) error

	GetBudgetComparisonReport(ctx context.Context, id uuid.UUID) (*db.GetBudgetComparisonReportRow, error)
}

type BudgetService interface {
	// --- Budget Operations ---
	CreateBudget(ctx context.Context, b *db.Budget) (*db.Budget, error)
	GetBudget(ctx context.Context, id uuid.UUID) (*db.Budget, error)
	ListBudgets(ctx context.Context, limit, offset int32) ([]db.Budget, error)
	UpdateBudget(ctx context.Context, b *db.Budget) (*db.Budget, error)
	DeleteBudget(ctx context.Context, id uuid.UUID) error

	// --- Budget Allocation Operations ---
	AllocateBudget(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error)
	GetBudgetAllocation(ctx context.Context, id uuid.UUID) (*db.BudgetAllocation, error)
	ListBudgetAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]db.BudgetAllocation, error)
	UpdateBudgetAllocation(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error)
	DeleteBudgetAllocation(ctx context.Context, id uuid.UUID) error

	// --- Reports ---
	GetBudgetComparisonReport(ctx context.Context, id uuid.UUID) (*db.GetBudgetComparisonReportRow, error)
}

type CashFlowForecastRepository interface {
	Generate(ctx context.Context, cf *db.CashFlowForecast) (*db.CashFlowForecast, error)
	Get(ctx context.Context, id uuid.UUID) (*db.CashFlowForecast, error)
	List(ctx context.Context, organizationID string, limit, offset int32) ([]*db.CashFlowForecast, error)
}

type CashFlowService interface {
    GenerateForecastFromPeriod(ctx context.Context, period *pb.ReportPeriod) (string, error)
    GetForecastFromPeriod(ctx context.Context, period *pb.ReportPeriod) (string, error)
    ListForecastsFromPeriod(ctx context.Context, period *pb.ReportPeriod) (string, error)
}


type ConsolidationRepository interface {
	Create(ctx context.Context, c db.Consolidation) (db.Consolidation, error)
	Get(ctx context.Context, id uuid.UUID) (db.Consolidation, error)
	List(ctx context.Context, entityIds []string, start, end time.Time, limit, offset int32) ([]db.Consolidation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreditDebitNoteRepository interface {
	Create(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error)
	Get(ctx context.Context, id uuid.UUID) (db.CreditDebitNote, error)
	List(ctx context.Context, limit, offset int32) ([]db.CreditDebitNote, error)
	Update(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreditDebitNoteService interface {
	Create(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error)
	Get(ctx context.Context, id uuid.UUID) (db.CreditDebitNote, error)
	List(ctx context.Context, limit, offset int32) ([]db.CreditDebitNote, error)
	Update(ctx context.Context, note db.CreditDebitNote) (db.CreditDebitNote, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ExchangeRateRepository interface {
	Create(ctx context.Context, rate db.ExchangeRate) (db.ExchangeRate, error)
	Get(ctx context.Context, id uuid.UUID) (db.ExchangeRate, error)
	Update(ctx context.Context, rate db.ExchangeRate) (db.ExchangeRate, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, base, quote *string, limit, offset int32) ([]db.ExchangeRate, error)
	GetLatest(ctx context.Context, base, quote string, asOf time.Time) (db.ExchangeRate, error)
}

type ExchangeRateService interface {
	Create(ctx context.Context, rate db.ExchangeRate) (db.ExchangeRate, error)
	Get(ctx context.Context, id uuid.UUID) (db.ExchangeRate, error)
	Update(ctx context.Context, rate db.ExchangeRate) (db.ExchangeRate, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, base, quote *string, limit, offset int32) ([]db.ExchangeRate, error)
	GetLatest(ctx context.Context, base, quote string, asOf time.Time) (db.ExchangeRate, error)
}

type ExpenseRepository interface {
	Create(ctx context.Context, expense db.Expense) (db.Expense, error)
	Get(ctx context.Context, id uuid.UUID) (db.Expense, error)
	List(ctx context.Context, limit, offset int32) ([]db.Expense, error)
	Update(ctx context.Context, expense db.Expense) (db.Expense, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CostCenterRepository interface {
	Create(ctx context.Context, cc db.CostCenter) (db.CostCenter, error)
	Get(ctx context.Context, id uuid.UUID) (db.CostCenter, error)
	List(ctx context.Context, limit, offset int32) ([]db.CostCenter, error)
	Update(ctx context.Context, cc db.CostCenter) (db.CostCenter, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CostAllocationRepository interface {
	Allocate(ctx context.Context, ca db.CostAllocation) (db.CostAllocation, error)
	List(ctx context.Context, limit, offset int32) ([]db.CostAllocation, error)
}

type ExpenseService interface {
	CreateExpense(ctx context.Context, exp db.Expense) (db.Expense, error)
	GetExpense(ctx context.Context, id uuid.UUID) (db.Expense, error)
	ListExpenses(ctx context.Context, limit, offset int32) ([]db.Expense, error)
	UpdateExpense(ctx context.Context, exp db.Expense) (db.Expense, error)
	DeleteExpense(ctx context.Context, id uuid.UUID) error
}

type CostCenterService interface {
	CreateCostCenter(ctx context.Context, cc db.CostCenter) (db.CostCenter, error)
	GetCostCenter(ctx context.Context, id uuid.UUID) (db.CostCenter, error)
	ListCostCenters(ctx context.Context, limit, offset int32) ([]db.CostCenter, error)
	UpdateCostCenter(ctx context.Context, cc db.CostCenter) (db.CostCenter, error)
	DeleteCostCenter(ctx context.Context, id uuid.UUID) error
}

type CostAllocationService interface {
	AllocateCost(ctx context.Context, ca db.CostAllocation) (db.CostAllocation, error)
	ListAllocations(ctx context.Context, limit, offset int32) ([]db.CostAllocation, error)
}

type FinanceEventRepository interface {
	InsertInvoiceCreated(ctx context.Context, e db.FinanceInvoiceCreatedEvent) (db.FinanceInvoiceCreatedEvent, error)
	ListInvoiceCreated(ctx context.Context, orgID string, limit, offset int32) ([]db.FinanceInvoiceCreatedEvent, error)

	InsertPaymentReceived(ctx context.Context, e db.FinancePaymentReceivedEvent) (db.FinancePaymentReceivedEvent, error)
	ListPaymentReceived(ctx context.Context, orgID string, limit, offset int32) ([]db.FinancePaymentReceivedEvent, error)

	InsertInventoryCostPosted(ctx context.Context, e db.InventoryCostPostedEvent) (db.InventoryCostPostedEvent, error)
	ListInventoryCostPosted(ctx context.Context, orgID string, limit, offset int32) ([]db.InventoryCostPostedEvent, error)

	InsertPayrollPosted(ctx context.Context, e db.PayrollPostedEvent) (db.PayrollPostedEvent, error)
	ListPayrollPosted(ctx context.Context, orgID string, limit, offset int32) ([]db.PayrollPostedEvent, error)

	InsertVendorBillApproved(ctx context.Context, e db.VendorBillApprovedEvent) (db.VendorBillApprovedEvent, error)
	ListVendorBillApproved(ctx context.Context, orgID string, limit, offset int32) ([]db.VendorBillApprovedEvent, error)
}

type FinancialReportsRepository interface {
	// Profit & Loss
	GenerateProfitLossReport(ctx context.Context, report db.ProfitLossReport) (db.ProfitLossReport, error)
	GetProfitLossReport(ctx context.Context, id uuid.UUID) (db.ProfitLossReport, error)
	ListProfitLossReports(ctx context.Context, orgID string, limit, offset int32) ([]db.ProfitLossReport, error)

	// Balance Sheet
	GenerateBalanceSheetReport(ctx context.Context, report db.BalanceSheetReport) (db.BalanceSheetReport, error)
	GetBalanceSheetReport(ctx context.Context, id uuid.UUID) (db.BalanceSheetReport, error)
	ListBalanceSheetReports(ctx context.Context, orgID string, limit, offset int32) ([]db.BalanceSheetReport, error)

	// Trial Balance
	CreateTrialBalanceReport(ctx context.Context, report db.TrialBalanceReport) (db.TrialBalanceReport, error)
	AddTrialBalanceEntry(ctx context.Context, entry db.TrialBalanceEntry) (db.TrialBalanceEntry, error)
	GetTrialBalanceReport(ctx context.Context, id uuid.UUID) (db.TrialBalanceReport, error)
	ListTrialBalanceReports(ctx context.Context, orgID string, limit, offset int32) ([]db.TrialBalanceReport, error)
	ListTrialBalanceEntries(ctx context.Context, reportID uuid.UUID) ([]db.TrialBalanceEntry, error)

	// Compliance
	GenerateComplianceReport(ctx context.Context, report db.ComplianceReport) (db.ComplianceReport, error)
	GetComplianceReport(ctx context.Context, id uuid.UUID) (db.ComplianceReport, error)
	ListComplianceReports(ctx context.Context, orgID, jurisdiction string, limit, offset int32) ([]db.ComplianceReport, error)
}

type GstRepository interface {
	// Breakup
	AddGstBreakup(ctx context.Context, invoiceID uuid.UUID, taxableAmount float64, cgst, sgst, igst, totalGst *float64) (db.GstBreakup, error)
	GetGstBreakup(ctx context.Context, invoiceID uuid.UUID) (db.GstBreakup, error)

	// Regime
	AddGstRegime(ctx context.Context, invoiceID uuid.UUID, gstin, placeOfSupply string, reverseCharge *bool) (db.GstRegime, error)
	GetGstRegime(ctx context.Context, invoiceID uuid.UUID) (db.GstRegime, error)

	// Doc Status
	AddGstDocStatus(ctx context.Context, invoiceID uuid.UUID, einvoiceStatus, irn, ackNo *string, ackDate *time.Time, ewayStatus, ewayBillNo *string, ewayValidUpto *time.Time,
		lastError *string, lastSyncedAt *time.Time) (db.GstDocStatus, error)

	GetGstDocStatus(ctx context.Context, invoiceID uuid.UUID) (db.GstDocStatus, error)
}

type GstServiceInterface interface {
	AddGstBreakup(ctx context.Context, invoiceID string, taxableAmount float64, cgst, sgst, igst, totalGst *float64) (db.GstBreakup, error)
	GetGstBreakup(ctx context.Context, invoiceID string) (db.GstBreakup, error)

	AddGstRegime(ctx context.Context, invoiceID, gstin, placeOfSupply string, reverseCharge *bool) (db.GstRegime, error)
	GetGstRegime(ctx context.Context, invoiceID string) (db.GstRegime, error)

	AddGstDocStatus(ctx context.Context, invoiceID string, einvoiceStatus, irn, ackNo *string, ackDate *time.Time, ewayStatus, ewayBillNo *string, ewayValidUpto *time.Time, lastError *string, lastSyncedAt *time.Time) (db.GstDocStatus, error)
	GetGstDocStatus(ctx context.Context, invoiceID string) (db.GstDocStatus, error)
}

type InvoiceRepository interface {
    CreateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error)
    UpdateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error)
    DeleteInvoice(ctx context.Context, id uuid.UUID) error
    GetInvoice(ctx context.Context, id uuid.UUID) (db.Invoice, error)
    ListInvoices(ctx context.Context, limit, offset int32) ([]db.Invoice, error)
    SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error)

    CreateInvoiceItem(ctx context.Context, item db.InvoiceItem) (db.InvoiceItem, error)
    ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]db.InvoiceItem, error)

    AddInvoiceTax(ctx context.Context, tax db.InvoiceTax) (db.InvoiceTax, error)
    AddInvoiceDiscount(ctx context.Context, discount db.InvoiceDiscount) (db.InvoiceDiscount, error)
}

type InvoiceServiceInterface interface {
    CreateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error)
    GetInvoice(ctx context.Context, id uuid.UUID) (db.Invoice, error)
    UpdateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error)
    DeleteInvoice(ctx context.Context, id uuid.UUID) error
    ListInvoices(ctx context.Context, limit, offset int32) ([]db.Invoice, error)
    SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error)
}

type BankAccountRepository interface {
    CreateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error)
    GetBankAccount(ctx context.Context, id uuid.UUID) (db.BankAccount, error)
    UpdateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error)
    DeleteBankAccount(ctx context.Context, id uuid.UUID) error
    ListBankAccounts(ctx context.Context, limit, offset int32) ([]db.BankAccount, error)
}


type PaymentDueRepository interface {
    CreatePaymentDue(ctx context.Context, pd db.PaymentDue) (db.PaymentDue, error)
    GetPaymentDue(ctx context.Context, id uuid.UUID) (db.PaymentDue, error)
    UpdatePaymentDue(ctx context.Context, pd db.PaymentDue) (db.PaymentDue, error)
    DeletePaymentDue(ctx context.Context, id uuid.UUID) error
    ListPaymentDues(ctx context.Context, limit, offset int32) ([]db.PaymentDue, error)
    MarkPaymentAsPaid(ctx context.Context, id uuid.UUID, updatedBy string) (db.PaymentDue, error)
}

type BankTransactionRepository interface {
    ImportBankTransaction(ctx context.Context, tx db.BankTransaction) (db.BankTransaction, error)
    ListBankTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int32) ([]db.BankTransaction, error)
    ReconcileTransaction(ctx context.Context, tx db.BankTransaction) (db.BankTransaction, error)
}

type BankService interface {
    CreateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error)
    GetBankAccount(ctx context.Context, id uuid.UUID) (db.BankAccount, error)
    UpdateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error)
    DeleteBankAccount(ctx context.Context, id uuid.UUID) error
    ListBankAccounts(ctx context.Context, limit, offset int32) ([]db.BankAccount, error)
}


type FinanceEventServiceInterface interface {
    RecordInvoiceCreated(ctx context.Context, event db.FinanceInvoiceCreatedEvent) (db.FinanceInvoiceCreatedEvent, error)
    RecordInvoiceUpdated(ctx context.Context, event db.Invoice) (db.Invoice, error)
    RecordInvoiceDeleted(ctx context.Context, invoiceID uuid.UUID) error
    RecordBudgetCreated(ctx context.Context, b db.Budget) error
    RecordBudgetUpdated(ctx context.Context, b db.Budget) error
    RecordBudgetDeleted(ctx context.Context, id uuid.UUID) error
}
