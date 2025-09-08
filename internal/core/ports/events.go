package ports

import (
    "context"

    "github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

// EventPublisher defines an abstraction for publishing domain events.
type EventPublisher interface {
    PublishAccountCreated(ctx context.Context, a db.Account) error
	PublishAccountUpdated(ctx context.Context, a db.Account) error
    PublishAccountDeleted(ctx context.Context, id string) error

    PublishJournalCreated(ctx context.Context, j db.JournalEntry) error
	PublishJournalUpdated(ctx context.Context, j db.JournalEntry) error
    PublishJournalDeleted(ctx context.Context, id string) error
}