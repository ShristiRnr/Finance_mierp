package domain

import (
	"time"

	"github.com/google/uuid"
	"database/sql"
)

// =====================
// Bank Accounts
// =====================
type BankAccount struct {
	ID              uuid.UUID
	Name            string
	AccountNumber   string
	IfscOrSwift     string
	LedgerAccountID sql.NullString
	CreatedAt       time.Time
	CreatedBy       sql.NullString
	UpdatedAt       time.Time
	UpdatedBy       sql.NullString
	Revision        int32
}

// =====================
// Bank Transactions
// =====================
type BankTransaction struct {
	ID                 uuid.UUID
	BankAccountID      uuid.UUID
	Amount             string
	TransactionDate    time.Time
	Description        sql.NullString
	Reference          sql.NullString
	Reconciled         sql.NullBool
	MatchedReferenceType sql.NullString
	MatchedReferenceID   uuid.NullUUID
	CreatedAt          time.Time
	CreatedBy          sql.NullString
	UpdatedAt          time.Time
	UpdatedBy          sql.NullString
	Revision           int32
}

// =====================
// Payment Dues
// =====================
type PaymentDue struct {
	ID        uuid.UUID
	InvoiceID uuid.UUID
	AmountDue string
	DueDate   time.Time
	Status    string
	CreatedAt time.Time
	CreatedBy sql.NullString
	UpdatedAt time.Time
	UpdatedBy sql.NullString
	Revision  int32
}
