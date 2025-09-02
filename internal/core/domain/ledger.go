package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)


// Account = Chart of Account
type Account struct {
	ID                 uuid.UUID
	Code               string
	Name               string
	Type               string
	ParentID           *uuid.UUID
	Status             string
	AllowManualJournal bool
	CreatedAt          time.Time
	CreatedBy          string
	UpdatedAt          time.Time
	UpdatedBy          string
	Revision           sql.NullInt32
}

// JournalEntry = Transaction header
type JournalEntry struct {
	ID         uuid.UUID
	JournalDate time.Time
	Reference   *string
	Memo        *string
	SourceType  *string
	SourceID    *string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	Revision    sql.NullInt32
}

// JournalLine = Debit/Credit detail
type JournalLine struct {
	AccountID    uuid.UUID
	Side         string // "DEBIT" or "CREDIT"
	Amount       string // decimal.Decimal if using shopspring/decimal
	CostCenterID *string
	Description  *string
	CreatedAt    time.Time
}

// Ledger Entry (read-only projection)
type LedgerEntry struct {
	EntryID   uuid.UUID
	AccountID uuid.UUID
	Side      string
	Amount    string
	PostedAt  time.Time
}