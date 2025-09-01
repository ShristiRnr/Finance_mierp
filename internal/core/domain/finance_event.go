package domain

import (
	"time"

	"github.com/google/uuid"
)

type FinanceInvoiceCreatedEvent struct {
	ID             uuid.UUID
	InvoiceID      uuid.UUID
	InvoiceNumber  string
	InvoiceDate    time.Time
	PartyRefID     uuid.NullUUID
	Total          string
	OrganizationID string
	CreatedAt      time.Time
}

type FinancePaymentReceivedEvent struct {
	ID             uuid.UUID
	PaymentDueID   uuid.UUID
	InvoiceID      uuid.UUID
	AmountPaid     string
	PaidAt         time.Time
	Reference      *string
	OrganizationID string
	CreatedAt      time.Time
}

type InventoryCostPostedEvent struct {
	ID             uuid.UUID
	ReferenceType  string
	ReferenceID    uuid.UUID
	Amount         string
	CostCenterID   *string
	OrganizationID string
	CreatedAt      time.Time
}

type PayrollPostedEvent struct {
	ID             uuid.UUID
	PayrollRunID   uuid.UUID
	TotalGross     string
	TotalNet       string
	RunDate        time.Time
	OrganizationID string
	CreatedAt      time.Time
}

type VendorBillApprovedEvent struct {
	ID             uuid.UUID
	VendorBillID   uuid.UUID
	Amount         string
	ApprovedAt     time.Time
	OrganizationID string
	CreatedAt      time.Time
}
