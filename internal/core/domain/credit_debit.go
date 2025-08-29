package domain

import (
	"time"

	"github.com/google/uuid"
)

// NoteType could be an enum in proto, mapped to string here
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

type CreditDebitNote struct {
	ID        uuid.UUID
	InvoiceID uuid.UUID
	Type      NoteType
	Amount    string // map from google.type.Money later
	Reason    string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	Revision  int32
	Refs      []ExternalRef
}
