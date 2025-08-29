package domain

import (
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	ID          uuid.UUID
	Name        string
	TotalAmount string
	Status      string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	Revision    int32
}

type BudgetAllocation struct {
	ID             uuid.UUID
	BudgetID       uuid.UUID
	DepartmentID   string
	AllocatedAmount string
	SpentAmount     string
	RemainingAmount string
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedAt      time.Time
	UpdatedBy      string
	Revision       int32
}

type BudgetComparisonReport struct {
	BudgetID       uuid.UUID
	TotalBudget    string
	TotalAllocated string
	TotalSpent     string
	RemainingBudget int32
}

type GetBudgetComparisonReportRow struct {
	BudgetID        uuid.UUID
	TotalBudget     string
	TotalAllocated  interface{}
	TotalSpent      interface{}
	RemainingBudget int32
}
