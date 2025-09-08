package ports

import (
    "context"

    "github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

// EventPublisher defines an abstraction for publishing domain events.
type EventPublisher interface {
    Publish(ctx context.Context, topic, key string, payload []byte) error
    // Account events
    PublishAccountCreated(ctx context.Context, a db.Account) error
    PublishAccountUpdated(ctx context.Context, a db.Account) error
    PublishAccountDeleted(ctx context.Context, id string) error

    // Journal events
    PublishJournalCreated(ctx context.Context, j db.JournalEntry) error
    PublishJournalUpdated(ctx context.Context, j db.JournalEntry) error
    PublishJournalDeleted(ctx context.Context, id string) error

    // Accrual events
    PublishAccrualCreated(ctx context.Context, a db.Accrual) error
    PublishAccrualUpdated(ctx context.Context, a db.Accrual) error
    PublishAccrualDeleted(ctx context.Context, id string) error

    // AllocationRule events
    PublishAllocationRuleCreated(ctx context.Context, r db.AllocationRule) error
    PublishAllocationRuleUpdated(ctx context.Context, r db.AllocationRule) error
    PublishAllocationRuleDeleted(ctx context.Context, id string) error
}
