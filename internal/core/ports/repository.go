package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
)

// AccountRepository handles chart of accounts
type AccountRepository interface {
	Create(ctx context.Context, a domain.Account) (domain.Account, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Account, error)
	List(ctx context.Context, limit, offset int32) ([]domain.Account, error)
	Update(ctx context.Context, a domain.Account) (domain.Account, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// JournalEntryRepository handles journal entries
type JournalEntryRepository interface {
	Create(ctx context.Context, j domain.JournalEntry) (domain.JournalEntry, error)
	Get(ctx context.Context, id uuid.UUID) (domain.JournalEntry, error)
	List(ctx context.Context, limit, offset int32) ([]domain.JournalEntry, error)
	Update(ctx context.Context, j domain.JournalEntry) (domain.JournalEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// JournalLineRepository handles lines inside journal entries
type JournalLineRepository interface {
	Add(ctx context.Context, l domain.JournalLine) (domain.JournalLine, error)
	ListByEntry(ctx context.Context, entryID uuid.UUID) ([]domain.JournalLine, error)
}
