package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Accrual struct {
	ID          uuid.UUID
	Description *string
	Amount      string
	AccrualDate time.Time
	AccountID   string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	Revision    sql.NullInt32
}

type AccrualExternalRef struct {
	ID        uuid.UUID
	AccrualID uuid.UUID
	System    string
	RefID     string
	CreatedAt time.Time
}

type AllocationRule struct {
	ID                  uuid.UUID
	Name                string
	Basis               string
	SourceAccountID     string
	TargetCostCenterIds []string
	Formula             *string
	CreatedAt           time.Time
	CreatedBy           string
	UpdatedAt           time.Time
	UpdatedBy           string
	Revision            sql.NullInt32
}
