package domain

import (
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
	Revision           int32
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
	Revision    int32
}

// JournalLine = Debit/Credit detail
type JournalLine struct {
	ID           uuid.UUID
	EntryID      uuid.UUID
	AccountID    uuid.UUID
	Side         string // "DEBIT" or "CREDIT"
	Amount       string // could be decimal.Decimal if using shopspring/decimal
	CostCenterID *string
	Description  *string
	CreatedAt    time.Time
}
