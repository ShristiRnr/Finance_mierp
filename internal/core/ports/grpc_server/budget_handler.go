package grpc_server

import (
	"context"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
)

type BudgetHandler struct {
	service *services.BudgetService
}

func NewBudgetHandler(service *services.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

// ---------------- Budgets ----------------

func (h *BudgetHandler) CreateBudget(ctx context.Context, b *domain.Budget) (*domain.Budget, error) {
	return h.service.CreateBudget(ctx, b)
}

func (h *BudgetHandler) GetBudget(ctx context.Context, id uuid.UUID) (*domain.Budget, error) {
	return h.service.GetBudget(ctx, id)
}

func (h *BudgetHandler) ListBudget(ctx context.Context, limit, offset int32) ([]*domain.Budget, error) {
	return h.service.ListBudgets(ctx, limit, offset)
}

func (h *BudgetHandler) UpdateBudget(ctx context.Context, b *domain.Budget) (*domain.Budget, error) {
	return h.service.UpdateBudget(ctx, b)
}

func (h *BudgetHandler) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	return h.service.DeleteBudget(ctx, id)
}

// ---------------- Budget Allocations ----------------

func (h *BudgetHandler) AllocateBudget(ctx context.Context, ba *domain.BudgetAllocation) (*domain.BudgetAllocation, error) {
	return h.service.AllocateBudget(ctx, ba)
}

func (h *BudgetHandler) GetBudgetAllocation(ctx context.Context, id uuid.UUID) (*domain.BudgetAllocation, error) {
	return h.service.GetBudgetAllocation(ctx, id)
}

func (h *BudgetHandler) ListBudgetAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]*domain.BudgetAllocation, error) {
	return h.service.ListBudgetAllocations(ctx, budgetID, limit, offset)
}

func (h *BudgetHandler) UpdateBudgetAllocation(ctx context.Context, ba *domain.BudgetAllocation) (*domain.BudgetAllocation, error) {
	return h.service.UpdateBudgetAllocation(ctx, ba)
}

func (h *BudgetHandler) DeleteBudgetAllocation(ctx context.Context, id uuid.UUID) error {
	return h.service.DeleteBudgetAllocation(ctx, id)
}

// ---------------- Budget Comparison ----------------

func (h *BudgetHandler) GetBudgetComparison(ctx context.Context, id uuid.UUID) (*domain.BudgetComparisonReport, error) {
	return h.service.GetBudgetComparison(ctx, id)
}
