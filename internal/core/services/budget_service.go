package services

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

// BudgetService with Kafka publisher
type BudgetService struct {
	repo      ports.BudgetRepository
	publisher ports.EventPublisher
}

func NewBudgetService(repo ports.BudgetRepository, pub ports.EventPublisher) *BudgetService {
	return &BudgetService{
		repo:      repo,
		publisher: pub,
	}
}

// ------------------- Budgets -------------------
func (s *BudgetService) CreateBudget(ctx context.Context, b *db.Budget) (*db.Budget, error) {
	created, err := s.repo.Create(ctx, b)
	if err != nil {
		return nil, err
	}

	// Publish event asynchronously
	if s.publisher != nil {
		go func(budget *db.Budget) {
			if err := s.publisher.PublishBudgetCreated(context.Background(), budget); err != nil {
				log.Printf("failed to publish budget created event: %v", err)
			}
		}(created)
	}

	return created, nil
}

func (s *BudgetService) GetBudget(ctx context.Context, id uuid.UUID) (*db.Budget, error) {
	return s.repo.Get(ctx, id)
}

func (s *BudgetService) ListBudgets(ctx context.Context, limit, offset int32) ([]db.Budget, error) {
	budgets, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]db.Budget, len(budgets))
	copy(result, budgets)
	return result, nil
}

func (s *BudgetService) UpdateBudget(ctx context.Context, b *db.Budget) (*db.Budget, error) {
	updated, err := s.repo.Update(ctx, b)
	if err != nil {
		return nil, err
	}

	// Publish event asynchronously
	if s.publisher != nil {
		go func(budget *db.Budget) {
			if err := s.publisher.PublishBudgetUpdated(context.Background(), budget); err != nil {
				log.Printf("failed to publish budget updated event: %v", err)
			}
		}(updated)
	}

	return updated, nil
}

func (s *BudgetService) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish deletion event asynchronously
	if s.publisher != nil {
		go func() {
			if err := s.publisher.PublishBudgetDeleted(context.Background(), id.String()); err != nil {
				log.Printf("failed to publish budget deleted event: %v", err)
			}
		}()
	}

	return nil
}

// ------------------- Budget Allocations -------------------
func (s *BudgetService) AllocateBudget(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	alloc, err := s.repo.Allocate(ctx, ba)
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		go func(alloc *db.BudgetAllocation) {
			if err := s.publisher.PublishBudgetAllocated(context.Background(), alloc); err != nil {
				log.Printf("failed to publish budget allocation event: %v", err)
			}
		}(alloc)
	}

	return alloc, nil
}

func (s *BudgetService) GetBudgetAllocation(ctx context.Context, id uuid.UUID) (*db.BudgetAllocation, error) {
	return s.repo.GetAllocation(ctx, id)
}

func (s *BudgetService) ListBudgetAllocations(ctx context.Context, budgetID uuid.UUID, limit, offset int32) ([]db.BudgetAllocation, error) {
	allocs, err := s.repo.ListAllocations(ctx, budgetID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]db.BudgetAllocation, len(allocs))
	copy(result, allocs)
	return result, nil
}

func (s *BudgetService) UpdateBudgetAllocation(ctx context.Context, ba *db.BudgetAllocation) (*db.BudgetAllocation, error) {
	updated, err := s.repo.UpdateAllocation(ctx, ba)
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		go func(alloc *db.BudgetAllocation) {
			if err := s.publisher.PublishBudgetAllocationUpdated(context.Background(), alloc); err != nil {
				log.Printf("failed to publish budget allocation updated event: %v", err)
			}
		}(updated)
	}

	return updated, nil
}

func (s *BudgetService) DeleteBudgetAllocation(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteAllocation(ctx, id); err != nil {
		return err
	}

	if s.publisher != nil {
		go func() {
			if err := s.publisher.PublishBudgetAllocationDeleted(context.Background(), id.String()); err != nil {
				log.Printf("failed to publish budget allocation deleted event: %v", err)
			}
		}()
	}

	return nil
}

// ------------------- Budget Comparison -------------------
func (s *BudgetService) GetBudgetComparisonReport(ctx context.Context, id uuid.UUID) (*db.GetBudgetComparisonReportRow, error) {
	row, err := s.repo.GetBudgetComparisonReport(ctx, id)
	if err != nil {
		return nil, err
	}
	return row, nil
}