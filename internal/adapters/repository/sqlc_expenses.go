package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

type ExpenseRepo struct {
	q *db.Queries
}

func NewExpenseRepo(q *db.Queries) ports.ExpenseRepository {
	return &ExpenseRepo{q: q}
}

func (r *ExpenseRepo) Create(ctx context.Context, exp db.Expense) (db.Expense, error) {
	arg := db.CreateExpenseParams{
		Category:     exp.Category,
		Amount:       exp.Amount,
		ExpenseDate:  exp.ExpenseDate,
		CostCenterID: exp.CostCenterID,
		CreatedBy:    exp.CreatedBy,
		UpdatedBy:    exp.UpdatedBy,
	}
	row, err := r.q.CreateExpense(ctx, arg)
	if err != nil {
		return db.Expense{}, err
	}
	return mapDbExpense(row), nil
}

func (r *ExpenseRepo) Get(ctx context.Context, id uuid.UUID) (db.Expense, error) {
	row, err := r.q.GetExpense(ctx, id)
	if err != nil {
		return db.Expense{}, err
	}
	return mapDbExpense(row), nil
}

func (r *ExpenseRepo) List(ctx context.Context, limit, offset int32) ([]db.Expense, error) {
	rows, err := r.q.ListExpenses(ctx, db.ListExpensesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	exps := make([]db.Expense, 0, len(rows))
	for _, row := range rows {
		exps = append(exps, mapDbExpense(row))
	}
	return exps, nil
}

func (r *ExpenseRepo) Update(ctx context.Context, exp db.Expense) (db.Expense, error) {
	arg := db.UpdateExpenseParams{
		ID:           exp.ID,
		Category:     exp.Category,
		Amount:       exp.Amount,
		ExpenseDate:  exp.ExpenseDate,
		CostCenterID: exp.CostCenterID,
		UpdatedBy:    exp.UpdatedBy,
	}
	row, err := r.q.UpdateExpense(ctx, arg)
	if err != nil {
		return db.Expense{}, err
	}
	return mapDbExpense(row), nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteExpense(ctx, id)
}

type CostCenterRepo struct {
	q *db.Queries
}

func NewCostCenterRepo(q *db.Queries) ports.CostCenterRepository {
	return &CostCenterRepo{q: q}
}

func (r *CostCenterRepo) Create(ctx context.Context, cc db.CostCenter) (db.CostCenter, error) {
	arg := db.CreateCostCenterParams{
		Name:        cc.Name,
		Description: cc.Description,
		CreatedBy:   cc.CreatedBy,
		UpdatedBy:   cc.UpdatedBy,
	}
	row, err := r.q.CreateCostCenter(ctx, arg)
	if err != nil {
		return db.CostCenter{}, err
	}
	return mapDbCostCenter(row), nil
}

func (r *CostCenterRepo) Get(ctx context.Context, id uuid.UUID) (db.CostCenter, error) {
	row, err := r.q.GetCostCenter(ctx, id)
	if err != nil {
		return db.CostCenter{}, err
	}
	return mapDbCostCenter(row), nil
}

func (r *CostCenterRepo) List(ctx context.Context, limit, offset int32) ([]db.CostCenter, error) {
	rows, err := r.q.ListCostCenters(ctx, db.ListCostCentersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	centers := make([]db.CostCenter, 0, len(rows))
	for _, row := range rows {
		centers = append(centers, mapDbCostCenter(row))
	}
	return centers, nil
}

func (r *CostCenterRepo) Update(ctx context.Context, cc db.CostCenter) (db.CostCenter, error) {
	arg := db.UpdateCostCenterParams{
		ID:          cc.ID,
		Name:        cc.Name,
		Description: cc.Description,
		UpdatedBy:   cc.UpdatedBy,
	}
	row, err := r.q.UpdateCostCenter(ctx, arg)
	if err != nil {
		return db.CostCenter{}, err
	}
	return mapDbCostCenter(row), nil
}

func (r *CostCenterRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteCostCenter(ctx, id)
}

type CostAllocationRepo struct {
	q *db.Queries
}

func NewCostAllocationRepo(q *db.Queries) ports.CostAllocationRepository {
	return &CostAllocationRepo{q: q}
}

func (r *CostAllocationRepo) Allocate(ctx context.Context, ca db.CostAllocation) (db.CostAllocation, error) {
	arg := db.AllocateCostParams{
		CostCenterID:  ca.CostCenterID,
		Amount:        ca.Amount,
		ReferenceType: ca.ReferenceType,
		ReferenceID:   ca.ReferenceID,
		CreatedBy:    ca.CreatedBy,
		UpdatedBy:     ca.UpdatedBy,
	}
	row, err := r.q.AllocateCost(ctx, arg)
	if err != nil {
		return db.CostAllocation{}, err
	}
	return mapDbCostAllocation(row), nil
}

func (r *CostAllocationRepo) List(ctx context.Context, limit, offset int32) ([]db.CostAllocation, error) {
	rows, err := r.q.ListCostAllocations(ctx, db.ListCostAllocationsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	allocs := make([]db.CostAllocation, 0, len(rows))
	for _, row := range rows {
		allocs = append(allocs, mapDbCostAllocation(row))
	}
	return allocs, nil
}

func mapDbExpense(row db.Expense) db.Expense {
	return db.Expense{
		ID:           row.ID,
		Category:     row.Category,
		Amount:       row.Amount,
		ExpenseDate:  row.ExpenseDate,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		CreatedBy:    row.CreatedBy,
		UpdatedBy:    row.UpdatedBy,
	}
}

func mapDbCostCenter(row db.CostCenter) db.CostCenter {
	return db.CostCenter{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		CreatedBy:   row.CreatedBy,
		UpdatedBy:   row.UpdatedBy,
	}
}

func mapDbCostAllocation(row db.CostAllocation) db.CostAllocation {
	return db.CostAllocation{
		ID:            row.ID,
		CostCenterID:  row.CostCenterID,
		Amount:        row.Amount,
		ReferenceType: row.ReferenceType,
		ReferenceID:   row.ReferenceID,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
		CreatedBy:     row.CreatedBy,
		UpdatedBy:     row.UpdatedBy,
	}
}