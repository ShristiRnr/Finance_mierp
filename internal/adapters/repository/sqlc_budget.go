package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
)

type BudgetRepository struct {
	queries *db.Queries
}

func NewBudgetRepo(queries *sql.DB) *BudgetRepository {
    return &BudgetRepository{
		queries: db.New(queries),
	}
}

// --------------------- Budgets ---------------------

func (r *BudgetRepository) Create(ctx context.Context, b *domain.Budget) (*domain.Budget, error) {
	row, err := r.queries.CreateBudget(ctx, db.CreateBudgetParams{
		Name:        b.Name,
		TotalAmount: b.TotalAmount,
		Column3:     nil,
		CreatedBy:   sql.NullString{String: b.CreatedBy, Valid: true},
		UpdatedBy:   sql.NullString{String: b.UpdatedBy, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return mapBudget(row), nil
}

func (r *BudgetRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Budget, error) {
	row, err := r.queries.GetBudget(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapBudget(row), nil
}

func (r *BudgetRepository) List(ctx context.Context, limit, offset int32) ([]*domain.Budget, error) {
	rows, err := r.queries.ListBudgets(ctx, db.ListBudgetsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Budget, len(rows))
	for i, row := range rows {
		result[i] = mapBudget(row)
	}
	return result, nil
}

func (r *BudgetRepository) Update(ctx context.Context, b *domain.Budget) (*domain.Budget, error) {
	row, err := r.queries.UpdateBudget(ctx, db.UpdateBudgetParams{
		ID:          b.ID,
		Name:        b.Name,
		TotalAmount: b.TotalAmount,
		Status:      b.Status,
		UpdatedBy:   sql.NullString{String: b.UpdatedBy, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return mapBudget(row), nil
}

func (r *BudgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteBudget(ctx, id)
}

// --------------------- Budget Allocations ---------------------

func (r *BudgetRepository) Allocate(ctx context.Context, ba *domain.BudgetAllocation) (*domain.BudgetAllocation, error) {
	row, err := r.queries.AllocateBudget(ctx, db.AllocateBudgetParams{
		BudgetID:        ba.BudgetID,
		DepartmentID:    ba.DepartmentID,
		AllocatedAmount: ba.AllocatedAmount,
		Column4:         nil, // default spent amount
		CreatedBy:       sql.NullString{String: ba.CreatedBy, Valid: true},
		UpdatedBy:       sql.NullString{String: ba.UpdatedBy, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return mapBudgetAllocation(row), nil
}

func (r *BudgetRepository) GetAllocation(ctx context.Context, id uuid.UUID) (*domain.BudgetAllocation, error) {
	row, err := r.queries.GetBudgetAllocation(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapBudgetAllocation(row), nil
}

func (r *BudgetRepository) ListAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]*domain.BudgetAllocation, error) {
	rows, err := r.queries.ListBudgetAllocations(ctx, db.ListBudgetAllocationsParams{
		BudgetID: budgetID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain.BudgetAllocation, len(rows))
	for i, row := range rows {
		result[i] = mapBudgetAllocation(row)
	}
	return result, nil
}

func (r *BudgetRepository) UpdateAllocation(ctx context.Context, ba *domain.BudgetAllocation) (*domain.BudgetAllocation, error) {
	row, err := r.queries.UpdateBudgetAllocation(ctx, db.UpdateBudgetAllocationParams{
		ID:              ba.ID,
		DepartmentID:    ba.DepartmentID,
		AllocatedAmount: ba.AllocatedAmount,
		SpentAmount:     sql.NullString{String: ba.SpentAmount, Valid: true},
		UpdatedBy:       sql.NullString{String: ba.UpdatedBy, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return mapBudgetAllocation(row), nil
}

func (r *BudgetRepository) DeleteAllocation(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteBudgetAllocation(ctx, id)
}

// --------------------- Budget Comparison ---------------------

func (r *BudgetRepository) GetBudgetComparison(ctx context.Context, id uuid.UUID) (*domain.BudgetComparisonReport, error) {
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

	return &domain.BudgetComparisonReport{
		BudgetID:       row.BudgetID,
		TotalBudget:    row.TotalBudget,
		TotalAllocated: totalAllocated,
		TotalSpent:     totalSpent,
		RemainingBudget:      row.RemainingBudget,
	}, nil
}

// --------------------- Helpers ---------------------

func mapBudget(row db.Budget) *domain.Budget {
	var revision int32
	if row.Revision.Valid {
		revision = row.Revision.Int32
	}

	return &domain.Budget{
		ID:          row.ID,
		Name:        row.Name,
		TotalAmount: row.TotalAmount,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt.Time,
		CreatedBy:   row.CreatedBy.String,
		UpdatedAt:   row.UpdatedAt.Time,
		UpdatedBy:   row.UpdatedBy.String,
		Revision:    revision,
	}
}

func mapBudgetAllocation(row db.BudgetAllocation) *domain.BudgetAllocation {
	var revision int32
	if row.Revision.Valid {
		revision = row.Revision.Int32
	}

	var spentAmount, remainingAmount string
	if row.SpentAmount.Valid {
		spentAmount = row.SpentAmount.String
	}
	if row.RemainingAmount.Valid {
		remainingAmount = row.RemainingAmount.String
	}

	return &domain.BudgetAllocation{
		ID:              row.ID,
		BudgetID:        row.BudgetID,
		DepartmentID:    row.DepartmentID,
		AllocatedAmount: row.AllocatedAmount,
		SpentAmount:     spentAmount,
		RemainingAmount: remainingAmount,
		CreatedAt:       row.CreatedAt.Time,
		CreatedBy:       row.CreatedBy.String,
		UpdatedAt:       row.UpdatedAt.Time,
		UpdatedBy:       row.UpdatedBy.String,
		Revision:        revision,
	}
}
