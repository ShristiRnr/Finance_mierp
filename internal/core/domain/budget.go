package domain

import (
	"time"

	"github.com/google/uuid"
)

// Budget represents a financial budget for a specific period
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

// BudgetAllocation represents the allocation of a budget to a specific department
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

// BudgetComparisonReport represents a report comparing budget planned vs actual income.
type BudgetComparisonReport struct {
	BudgetID       uuid.UUID
	TotalBudget    string
	TotalAllocated string
	TotalSpent     string
	RemainingBudget int32
}

// GetBudgetComparisonReportRow represents a row in the budget comparison report.
type GetBudgetComparisonReportRow struct {
	BudgetID        uuid.UUID
	TotalBudget     string
	TotalAllocated  interface{}
	TotalSpent      interface{}
	RemainingBudget int32
}
