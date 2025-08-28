package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
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

type AccrualRepository interface {
	Create(ctx context.Context, a domain.Accrual) (domain.Accrual, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Accrual, error)
	Update(ctx context.Context, a domain.Accrual) (domain.Accrual, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]domain.Accrual, error)

	AddExternalRef(ctx context.Context, ref domain.AccrualExternalRef) (domain.AccrualExternalRef, error)
	ListExternalRefs(ctx context.Context, accrualID uuid.UUID) ([]domain.AccrualExternalRef, error)
}

type AllocationRuleRepository interface {
	Create(ctx context.Context, r domain.AllocationRule) (domain.AllocationRule, error)
	Get(ctx context.Context, id uuid.UUID) (domain.AllocationRule, error)
	Update(ctx context.Context, r domain.AllocationRule) (domain.AllocationRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]domain.AllocationRule, error)
}