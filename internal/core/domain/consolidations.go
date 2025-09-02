package domain

import (
	"time"

	"github.com/google/uuid"
)

// Consolidation represents a financial consolidation for a specific period
type Consolidation struct {
	ID         uuid.UUID
	EntityIds  []string
	PeriodStart time.Time
	PeriodEnd   time.Time
	Report      string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	Revision    int32
}
