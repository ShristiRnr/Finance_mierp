package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type BudgetService struct {
	repo ports.BudgetRepository
}

func NewBudgetService(repo ports.BudgetRepository) *BudgetService {
	return &BudgetService{repo: repo}
}

// Budgets
func (s *BudgetService) CreateBudget(ctx context.Context, b *domain.Budget) (*domain.Budget, error) {
	return s.repo.Create(ctx, b)
}

func (s *BudgetService) GetBudget(ctx context.Context, id uuid.UUID) (*domain.Budget, error) {
	return s.repo.Get(ctx, id)
}

func (s *BudgetService) ListBudgets(ctx context.Context, limit, offset int32) ([]*domain.Budget, error) {
	budgets, err := s.repo.List(ctx, limit, offset) // returns []domain.Budget
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Budget, len(budgets))
	for i := range budgets {
		result[i] = &budgets[i]
	}
	return result, nil
}

func (s *BudgetService) UpdateBudget(ctx context.Context, b *domain.Budget) (*domain.Budget, error) {
	return s.repo.Update(ctx, b)
}

func (s *BudgetService) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Budget Allocations
func (s *BudgetService) AllocateBudget(ctx context.Context, ba *domain.BudgetAllocation) (*domain.BudgetAllocation, error) {
	return s.repo.Allocate(ctx, ba)
}

func (s *BudgetService) GetBudgetAllocation(ctx context.Context, id uuid.UUID) (*domain.BudgetAllocation, error) {
	return s.repo.GetAllocation(ctx, id)
}

func (s *BudgetService) ListBudgetAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]*domain.BudgetAllocation, error) {
	allocs, err := s.repo.ListAllocations(ctx, budgetID, limit, offset) // returns []domain.BudgetAllocation
	if err != nil {
		return nil, err
	}

	result := make([]*domain.BudgetAllocation, len(allocs))
	for i := range allocs {
		result[i] = &allocs[i]
	}
	return result, nil
}
func (s *BudgetService) UpdateBudgetAllocation(ctx context.Context, ba *domain.BudgetAllocation) (*domain.BudgetAllocation, error) {
	return s.repo.UpdateAllocation(ctx, ba)
}

func (s *BudgetService) DeleteBudgetAllocation(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteAllocation(ctx, id)
}

// Budget Comparison
func (s *BudgetService) GetBudgetComparison(ctx context.Context, id uuid.UUID) (*domain.BudgetComparisonReport, error) {
	return s.repo.GetBudgetComparison(ctx, id)
}
