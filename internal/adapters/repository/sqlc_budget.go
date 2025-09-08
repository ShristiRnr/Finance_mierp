package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

type BudgetRepository struct {
	queries *db.Queries
}

func NewBudgetRepo(queries *sql.DB) *BudgetRepository {
    return &BudgetRepository{
		queries: db.New(queries),
	}
}

// ======================================== Budgets ========================================

func (r *BudgetRepository) Create(ctx context.Context, b *db.Budget) (*db.Budget, error) {
	row, err := r.queries.CreateBudget(ctx, db.CreateBudgetParams{
		Name:        b.Name,
		TotalAmount: b.TotalAmount,
		Column3:     nil,
		CreatedBy:   b.CreatedBy,
		UpdatedBy:   b.UpdatedBy,
	})
	if err != nil {
		return nil, err
	}

	return mapBudget(row), nil
}

func (r *BudgetRepository) Get(ctx context.Context, id uuid.UUID) (*db.Budget, error) {
	row, err := r.queries.GetBudget(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapBudget(row), nil
}

func (r *BudgetRepository) List(ctx context.Context, limit, offset int32) ([]*db.Budget, error) {
	rows, err := r.queries.ListBudgets(ctx, db.ListBudgetsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*db.Budget, len(rows))
	for i, row := range rows {
		result[i] = mapBudget(row)
	}
	return result, nil
}

func (r *BudgetRepository) Update(ctx context.Context, b *db.Budget) (*db.Budget, error) {
	row, err := r.queries.UpdateBudget(ctx, db.UpdateBudgetParams{
		ID:          b.ID,
		Name:        b.Name,
		TotalAmount: b.TotalAmount,
		Status:      b.Status,
		UpdatedBy:   b.UpdatedBy,
	})
	if err != nil {
		return nil, err
	}
	return mapBudget(row), nil
}

func (r *BudgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteBudget(ctx, id)
}

// ===================================== Budget Allocations =======================================

func (r *BudgetRepository) Allocate(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	row, err := r.queries.AllocateBudget(ctx, db.AllocateBudgetParams{
		BudgetID:        ba.BudgetID,
		DepartmentID:    ba.DepartmentID,
		AllocatedAmount: ba.AllocatedAmount,
		Column4:         nil, // default spent amount
		CreatedBy:       ba.CreatedBy,
		UpdatedBy:       ba.UpdatedBy,
	})
	if err != nil {
		return nil, err
	}
	return mapBudgetAllocation(row), nil
}

func (r *BudgetRepository) GetAllocation(ctx context.Context, id uuid.UUID) (*db.BudgetAllocation, error) {
	row, err := r.queries.GetBudgetAllocation(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapBudgetAllocation(row), nil
}

func (r *BudgetRepository) ListAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]*db.BudgetAllocation, error) {
	rows, err := r.queries.ListBudgetAllocations(ctx, db.ListBudgetAllocationsParams{
		BudgetID: budgetID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*db.BudgetAllocation, len(rows))
	for i, row := range rows {
		result[i] = mapBudgetAllocation(row)
	}
	return result, nil
}

func (r *BudgetRepository) UpdateAllocation(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	row, err := r.queries.UpdateBudgetAllocation(ctx, db.UpdateBudgetAllocationParams{
		ID:              ba.ID,
		DepartmentID:    ba.DepartmentID,
		AllocatedAmount: ba.AllocatedAmount,
		SpentAmount:     ba.SpentAmount,
		UpdatedBy:       ba.UpdatedBy,
	})
	if err != nil {
		return nil, err
	}
	return mapBudgetAllocation(row), nil
}

func (r *BudgetRepository) DeleteAllocation(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteBudgetAllocation(ctx, id)
}

// ======================================== Budget Comparison =======================================

func (r *BudgetRepository) GetBudgetComparison(ctx context.Context, id uuid.UUID) (*db.GetBudgetComparisonReportRow, error) {
	row, err := r.queries.GetBudgetComparisonReport(ctx, id)
	if err != nil {
		return nil, err
	}

	totalAllocated := ""
	totalSpent := ""

	if row.TotalAllocated != nil {
		totalAllocated = fmt.Sprintf("%v", row.TotalAllocated)
	}
	if row.TotalSpent != nil {
		totalSpent = fmt.Sprintf("%v", row.TotalSpent)
	}

	return &db.GetBudgetComparisonReportRow{
		BudgetID:       row.BudgetID,
		TotalBudget:    row.TotalBudget,
		TotalAllocated: totalAllocated,
		TotalSpent:     totalSpent,
		RemainingBudget:      row.RemainingBudget,
	}, nil
}

// --------------------- Helpers ---------------------

func mapBudget(row db.Budget) *db.Budget {

	return &db.Budget{
		ID:          row.ID,
		Name:        row.Name,
		TotalAmount: row.TotalAmount,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt,
		CreatedBy:   row.CreatedBy,
		UpdatedAt:   row.UpdatedAt,
		UpdatedBy:   row.UpdatedBy,
		Revision:    row.Revision,
	}
}

func mapBudgetAllocation(row db.BudgetAllocation) *db.BudgetAllocation {

	return &db.BudgetAllocation{
		ID:              row.ID,
		BudgetID:        row.BudgetID,
		DepartmentID:    row.DepartmentID,
		AllocatedAmount: row.AllocatedAmount,
		SpentAmount:     row.SpentAmount,
		RemainingAmount: row.RemainingAmount,
		CreatedAt:       row.CreatedAt,
		CreatedBy:       row.CreatedBy,
		UpdatedAt:       row.UpdatedAt,
		UpdatedBy:       row.UpdatedBy,
		Revision:        row.Revision,
	}
}
