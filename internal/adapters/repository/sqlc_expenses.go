package repository

import (
	"context"

	"database/sql"
	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

type ExpenseRepo struct {
	q *db.Queries
}

func NewExpenseRepo(q *db.Queries) ports.ExpenseRepository {
	return &ExpenseRepo{q: q}
}

func (r *ExpenseRepo) Create(ctx context.Context, exp domain.Expense) (domain.Expense, error) {
	arg := db.CreateExpenseParams{
		Category:     exp.Category,
		Amount:       exp.Amount,
		ExpenseDate:  exp.ExpenseDate,
		CostCenterID: uuidToNullUUID(uuid.UUID(exp.CostCenter.NodeID())),
		CreatedBy:    strToNullString(exp.CreatedBy),
		UpdatedBy:    strToNullString(exp.UpdatedBy),
	}
	row, err := r.q.CreateExpense(ctx, arg)
	if err != nil {
		return domain.Expense{}, err
	}
	return mapDbExpense(row), nil
}

func (r *ExpenseRepo) Get(ctx context.Context, id uuid.UUID) (domain.Expense, error) {
	row, err := r.q.GetExpense(ctx, id)
	if err != nil {
		return domain.Expense{}, err
	}
	return mapDbExpense(row), nil
}

func (r *ExpenseRepo) List(ctx context.Context, limit, offset int32) ([]domain.Expense, error) {
	rows, err := r.q.ListExpenses(ctx, db.ListExpensesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	exps := make([]domain.Expense, 0, len(rows))
	for _, row := range rows {
		exps = append(exps, mapDbExpense(row))
	}
	return exps, nil
}

func (r *ExpenseRepo) Update(ctx context.Context, exp domain.Expense) (domain.Expense, error) {
	arg := db.UpdateExpenseParams{
		ID:           exp.ID,
		Category:     exp.Category,
		Amount:       exp.Amount,
		ExpenseDate:  exp.ExpenseDate,
		CostCenterID: uuidToNullUUID(uuid.UUID(exp.CostCenter.NodeID())),
		UpdatedBy:    strToNullString(exp.UpdatedBy),
	}
	row, err := r.q.UpdateExpense(ctx, arg)
	if err != nil {
		return domain.Expense{}, err
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

func (r *CostCenterRepo) Create(ctx context.Context, cc domain.CostCenter) (domain.CostCenter, error) {
	arg := db.CreateCostCenterParams{
		Name:        cc.Name,
		Description: strToNullString(cc.Description),
		CreatedBy:   strToNullString(cc.CreatedBy),
		UpdatedBy:   strToNullString(cc.UpdatedBy),
	}
	row, err := r.q.CreateCostCenter(ctx, arg)
	if err != nil {
		return domain.CostCenter{}, err
	}
	return mapDbCostCenter(row), nil
}

func (r *CostCenterRepo) Get(ctx context.Context, id uuid.UUID) (domain.CostCenter, error) {
	row, err := r.q.GetCostCenter(ctx, id)
	if err != nil {
		return domain.CostCenter{}, err
	}
	return mapDbCostCenter(row), nil
}

func (r *CostCenterRepo) List(ctx context.Context, limit, offset int32) ([]domain.CostCenter, error) {
	rows, err := r.q.ListCostCenters(ctx, db.ListCostCentersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	centers := make([]domain.CostCenter, 0, len(rows))
	for _, row := range rows {
		centers = append(centers, mapDbCostCenter(row))
	}
	return centers, nil
}

func (r *CostCenterRepo) Update(ctx context.Context, cc domain.CostCenter) (domain.CostCenter, error) {
	arg := db.UpdateCostCenterParams{
		ID:          cc.ID,
		Name:        cc.Name,
		Description: strToNullString(cc.Description),
		UpdatedBy:   strToNullString(cc.UpdatedBy),
	}
	row, err := r.q.UpdateCostCenter(ctx, arg)
	if err != nil {
		return domain.CostCenter{}, err
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

func (r *CostAllocationRepo) Allocate(ctx context.Context, ca domain.CostAllocation) (domain.CostAllocation, error) {
	arg := db.AllocateCostParams{
		CostCenterID:  ca.CostCenterID,
		Amount:        ca.Amount,
		ReferenceType: ca.ReferenceType,
		ReferenceID:   ca.ReferenceID,
		CreatedBy:     strToNullString(ca.CreatedBy),
		UpdatedBy:     strToNullString(ca.UpdatedBy),
	}
	row, err := r.q.AllocateCost(ctx, arg)
	if err != nil {
		return domain.CostAllocation{}, err
	}
	return mapDbCostAllocation(row), nil
}

func (r *CostAllocationRepo) List(ctx context.Context, limit, offset int32) ([]domain.CostAllocation, error) {
	rows, err := r.q.ListCostAllocations(ctx, db.ListCostAllocationsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	allocs := make([]domain.CostAllocation, 0, len(rows))
	for _, row := range rows {
		allocs = append(allocs, mapDbCostAllocation(row))
	}
	return allocs, nil
}

func mapDbExpense(row db.Expense) domain.Expense {
	return domain.Expense{
		ID:           row.ID,
		Category:     row.Category,
		Amount:       row.Amount,
		ExpenseDate:  row.ExpenseDate,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		CreatedBy:    nullStringToPtr(row.CreatedBy),
		UpdatedBy:    nullStringToPtr(row.UpdatedBy),
	}
}

func mapDbCostCenter(row db.CostCenter) domain.CostCenter {
	return domain.CostCenter{
		ID:          row.ID,
		Name:        row.Name,
		Description: nullStringToPtr(row.Description),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
		CreatedBy:   nullStringToPtr(row.CreatedBy),
		UpdatedBy:   nullStringToPtr(row.UpdatedBy),
	}
}

func mapDbCostAllocation(row db.CostAllocation) domain.CostAllocation {
	return domain.CostAllocation{
		ID:            row.ID,
		CostCenterID:  row.CostCenterID,
		Amount:        row.Amount,
		ReferenceType: row.ReferenceType,
		ReferenceID:   row.ReferenceID,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		CreatedBy:     nullStringToPtr(row.CreatedBy),
		UpdatedBy:     nullStringToPtr(row.UpdatedBy),
	}
}

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}


func uuidToNullUUID(id uuid.UUID) uuid.NullUUID {
	if id == uuid.Nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: id, Valid: true}
}

func strToNullString(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}