package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type ExpenseService struct {
	repo ports.ExpenseRepository
}

func NewExpenseService(r ports.ExpenseRepository) ports.ExpenseService {
	return &ExpenseService{repo: r}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, exp domain.Expense) (domain.Expense, error) {
	return s.repo.Create(ctx, exp)
}

func (s *ExpenseService) GetExpense(ctx context.Context, id uuid.UUID) (domain.Expense, error) {
	return s.repo.Get(ctx, id)
}

func (s *ExpenseService) ListExpenses(ctx context.Context, limit, offset int32) ([]domain.Expense, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, exp domain.Expense) (domain.Expense, error) {
	return s.repo.Update(ctx, exp)
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

type CostCenterService struct {
	repo ports.CostCenterRepository
}

func NewCostCenterService(r ports.CostCenterRepository) ports.CostCenterService {
	return &CostCenterService{repo: r}
}

func (s *CostCenterService) CreateCostCenter(ctx context.Context, cc domain.CostCenter) (domain.CostCenter, error) {
	return s.repo.Create(ctx, cc)
}

func (s *CostCenterService) GetCostCenter(ctx context.Context, id uuid.UUID) (domain.CostCenter, error) {
	return s.repo.Get(ctx, id)
}

func (s *CostCenterService) ListCostCenters(ctx context.Context, limit, offset int32) ([]domain.CostCenter, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *CostCenterService) UpdateCostCenter(ctx context.Context, cc domain.CostCenter) (domain.CostCenter, error) {
	return s.repo.Update(ctx, cc)
}

func (s *CostCenterService) DeleteCostCenter(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}


type CostAllocationService struct {
	repo ports.CostAllocationRepository
}

func NewCostAllocationService(r ports.CostAllocationRepository) ports.CostAllocationService {
	return &CostAllocationService{repo: r}
}

func (s *CostAllocationService) AllocateCost(ctx context.Context, ca domain.CostAllocation) (domain.CostAllocation, error) {
	return s.repo.Allocate(ctx, ca)
}

func (s *CostAllocationService) ListAllocations(ctx context.Context, limit, offset int32) ([]domain.CostAllocation, error) {
	return s.repo.List(ctx, limit, offset)
}