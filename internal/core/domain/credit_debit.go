package domain

import (
	"time"

	"github.com/google/uuid"
)

type NoteType string

const (
	NoteTypeCredit NoteType = "CREDIT"
	NoteTypeDebit  NoteType = "DEBIT"
)

type ExternalRef struct {
	ID        uuid.UUID
	System    string
	RefID     string
	CreatedAt time.Time
}

// CreditDebitNote represents a credit or debit note for an invoice
type CreditDebitNote struct {
	ID        uuid.UUID
	InvoiceID uuid.UUID
	Type      NoteType
	Amount    string // mapping done from google.type.Money
	Reason    string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	Revision  int32
	Refs      []ExternalRef
}
