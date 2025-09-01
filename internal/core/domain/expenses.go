package domain

import(
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID
	Category    string
	Amount      string // consider decimal.Decimal for money
	ExpenseDate time.Time
	CostCenter  *uuid.UUID
	CreatedAt   time.Time
	CreatedBy   *string
	UpdatedAt   time.Time
	UpdatedBy   *string
	Revision    int32
}

type ExpenseExternalRef struct {
	ID        uuid.UUID
	ExpenseID uuid.UUID
	System    string
	RefID     string
	CreatedAt time.Time
}

type CostCenter struct {
	ID          uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time
	CreatedBy   *string
	UpdatedAt   time.Time
	UpdatedBy   *string
	Revision    int32
}

type CostAllocation struct {
	ID           uuid.UUID
	CostCenterID uuid.UUID
	Amount       string // again consider decimal
	ReferenceType string
	ReferenceID   string
	CreatedAt    time.Time
	CreatedBy    *string
	UpdatedAt    time.Time
	UpdatedBy    *string
	Revision     int32
}