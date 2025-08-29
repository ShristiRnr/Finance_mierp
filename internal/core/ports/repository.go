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

	AddExternalRef(ctx context.Context, ref domain.AccrualExternalRef) (domain.AccrualExternalRef, error)
	ListExternalRefs(ctx context.Context, accrualID uuid.UUID) ([]domain.AccrualExternalRef, error)
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
