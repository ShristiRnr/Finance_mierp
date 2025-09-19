package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
)

func TestBudgetRepository_CRUD(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	repo := repository.NewBudgetRepository(dbConn)
	ctx := context.Background()

	// ---------- Common data ----------
	budgetID := uuid.New()
	allocationID := uuid.New()
	now := time.Now()

	budget := &db.Budget{
		ID:          budgetID,
		Name:        "Marketing Budget",
		TotalAmount: "10000",
		Status:      "DRAFT",
		CreatedBy:   sql.NullString{String: "admin", Valid: true},
		UpdatedBy:   sql.NullString{String: "admin", Valid: true},
		CreatedAt:   sql.NullTime{Time: now, Valid: true},
		UpdatedAt:   sql.NullTime{Time: now, Valid: true},
	}

	allocation := &db.BudgetAllocation{
		ID:              allocationID,
		BudgetID:        budgetID,
		DepartmentID:    "Sales",
		AllocatedAmount: "5000",
		SpentAmount:     sql.NullString{String: "2000", Valid: true},
		RemainingAmount: sql.NullString{String: "3000", Valid: true},
		CreatedBy:       sql.NullString{String: "admin", Valid: true},
		UpdatedBy:       sql.NullString{String: "admin", Valid: true},
		CreatedAt:       sql.NullTime{Time: now, Valid: true},
		UpdatedAt:       sql.NullTime{Time: now, Valid: true},
	}

	// ================== Budget: Create ==================
	mock.ExpectQuery(`INSERT INTO budgets`).
		WithArgs(budget.Name, budget.TotalAmount, sql.NullString{}, budget.CreatedBy, budget.UpdatedBy).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "total_amount", "status", "created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			budgetID, budget.Name, budget.TotalAmount, budget.Status,
			budget.CreatedAt, budget.CreatedBy, budget.UpdatedAt, budget.UpdatedBy, 1,
		))

	created, err := repo.Create(ctx, budget)
	require.NoError(t, err)
	require.Equal(t, budgetID, created.ID)

	// ================== Budget: Get ==================
	mock.ExpectQuery(`SELECT .* FROM budgets WHERE id = \$1`).
		WithArgs(budgetID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "total_amount", "status", "created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			budgetID, budget.Name, budget.TotalAmount, budget.Status,
			budget.CreatedAt, budget.CreatedBy, budget.UpdatedAt, budget.UpdatedBy, 1,
		))

	got, err := repo.Get(ctx, budgetID)
	require.NoError(t, err)
	require.Equal(t, budgetID, got.ID)

	// ================== Budget: List ==================
	mock.ExpectQuery(`SELECT .* FROM budgets ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "total_amount", "status", "created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			budgetID, budget.Name, budget.TotalAmount, budget.Status,
			budget.CreatedAt, budget.CreatedBy, budget.UpdatedAt, budget.UpdatedBy, 1,
		))

	list, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, budgetID, list[0].ID)

	// ================== Budget: Update ==================
	budget.Status = "APPROVED"
	mock.ExpectQuery(`UPDATE budgets`).
		WithArgs(budgetID, budget.Name, budget.TotalAmount, budget.Status, budget.UpdatedBy).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "total_amount", "status", "created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			budgetID, budget.Name, budget.TotalAmount, budget.Status,
			budget.CreatedAt, budget.CreatedBy, budget.UpdatedAt, budget.UpdatedBy, 2,
		))

	updated, err := repo.Update(ctx, budget)
	require.NoError(t, err)
	require.Equal(t, "APPROVED", updated.Status)

	// ================== Budget: Delete ==================
	mock.ExpectExec(`DELETE FROM budgets WHERE id = \$1`).
		WithArgs(budgetID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(ctx, budgetID)
	require.NoError(t, err)

	// ================== Allocation: Allocate ==================
	mock.ExpectQuery(`INSERT INTO budget_allocations`).
		WithArgs(allocation.BudgetID, allocation.DepartmentID, allocation.AllocatedAmount, sql.NullString{}, allocation.CreatedBy, allocation.UpdatedBy).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "budget_id", "department_id", "allocated_amount", "spent_amount", "remaining_amount", "created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			allocationID, allocation.BudgetID, allocation.DepartmentID, allocation.AllocatedAmount, allocation.SpentAmount,
			allocation.RemainingAmount, allocation.CreatedAt, allocation.CreatedBy, allocation.UpdatedAt, allocation.UpdatedBy, 1,
		))

	alloc, err := repo.Allocate(ctx, allocation)
	require.NoError(t, err)
	require.Equal(t, allocationID, alloc.ID)

	// ================== Allocation: Get ==================
	mock.ExpectQuery(`SELECT .* FROM budget_allocations WHERE id = \$1`).
		WithArgs(allocationID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "budget_id", "department_id", "allocated_amount", "spent_amount", "remaining_amount",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			allocationID, allocation.BudgetID, allocation.DepartmentID, allocation.AllocatedAmount, allocation.SpentAmount,
			allocation.RemainingAmount, allocation.CreatedAt, allocation.CreatedBy, allocation.UpdatedAt, allocation.UpdatedBy, 1,
		))

	gotAlloc, err := repo.GetAllocation(ctx, allocationID)
	require.NoError(t, err)
	require.Equal(t, allocationID, gotAlloc.ID)

	// ================== Allocation: List ==================
	mock.ExpectQuery(`SELECT .* FROM budget_allocations WHERE budget_id = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(budgetID, int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "budget_id", "department_id", "allocated_amount", "spent_amount", "remaining_amount",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			allocationID, allocation.BudgetID, allocation.DepartmentID, allocation.AllocatedAmount, allocation.SpentAmount,
			allocation.RemainingAmount, allocation.CreatedAt, allocation.CreatedBy, allocation.UpdatedAt, allocation.UpdatedBy, 1,
		))

	listAlloc, err := repo.ListAllocations(ctx, budgetID, 10, 0)
	require.NoError(t, err)
	require.Len(t, listAlloc, 1)
	require.Equal(t, allocationID, listAlloc[0].ID)

	// ================== Allocation: Update ==================
	allocation.SpentAmount = sql.NullString{String: "2500", Valid: true}
	mock.ExpectQuery(`UPDATE budget_allocations`).
		WithArgs(allocation.ID, allocation.DepartmentID, allocation.AllocatedAmount, allocation.SpentAmount, allocation.UpdatedBy).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "budget_id", "department_id", "allocated_amount", "spent_amount", "remaining_amount",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			allocationID, allocation.BudgetID, allocation.DepartmentID, allocation.AllocatedAmount, allocation.SpentAmount,
			allocation.RemainingAmount, allocation.CreatedAt, allocation.CreatedBy, allocation.UpdatedAt, allocation.UpdatedBy, 2,
		))

	updatedAlloc, err := repo.UpdateAllocation(ctx, allocation)
	require.NoError(t, err)
	require.Equal(t, "2500", updatedAlloc.SpentAmount.String)

	// ================== Allocation: Delete ==================
	mock.ExpectExec(`DELETE FROM budget_allocations WHERE id = \$1`).
		WithArgs(allocationID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteAllocation(ctx, allocationID)
	require.NoError(t, err)

	// ================== Budget Comparison Report ==================
	mock.ExpectQuery(`SELECT .* FROM budgets b LEFT JOIN budget_allocations`).
		WithArgs(budgetID).
		WillReturnRows(sqlmock.NewRows([]string{
			"budget_id", "total_budget", "total_allocated", "total_spent", "remaining_budget",
		}).AddRow(
			budgetID, budget.TotalAmount, "5000", "2000", 5000,
		))

	report, err := repo.GetBudgetComparisonReport(ctx, budgetID)
	require.NoError(t, err)
	require.Equal(t, budgetID, report.BudgetID)

	require.NoError(t, mock.ExpectationsWereMet())
}