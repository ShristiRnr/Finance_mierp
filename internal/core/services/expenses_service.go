package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type ExpenseService struct {
	repo      ports.ExpenseRepository
	publisher ports.EventPublisher
}

func NewExpenseService(r ports.ExpenseRepository, publisher ports.EventPublisher) ports.ExpenseService {
	return &ExpenseService{repo: r, publisher: publisher}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, exp db.Expense) (db.Expense, error) {
	e, err := s.repo.Create(ctx, exp)
	if err != nil {
		return e, err
	}

	if err := s.publisher.PublishExpenseCreated(ctx, &e); err != nil {
		fmt.Printf("Kafka publish error (expense.created): %v\n", err)
	}
	return e, nil
}

func (s *ExpenseService) GetExpense(ctx context.Context, id uuid.UUID) (db.Expense, error) {
	return s.repo.Get(ctx, id)
}

func (s *ExpenseService) ListExpenses(ctx context.Context, limit, offset int32) ([]db.Expense, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, exp db.Expense) (db.Expense, error) {
	e, err := s.repo.Update(ctx, exp)
	if err != nil {
		return e, err
	}

	if err := s.publisher.PublishExpenseUpdated(ctx, &e); err != nil {
		fmt.Printf("Kafka publish error (expense.updated): %v\n", err)
	}
	return e, nil
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.publisher.PublishExpenseDeleted(ctx, id.String()); err != nil {
		fmt.Printf("Kafka publish error (expense.deleted): %v\n", err)
	}
	return nil
}

type CostCenterService struct {
	repo      ports.CostCenterRepository
	publisher ports.EventPublisher
}

func NewCostCenterService(r ports.CostCenterRepository, publisher ports.EventPublisher) ports.CostCenterService {
	return &CostCenterService{repo: r, publisher: publisher}
}

func (s *CostCenterService) CreateCostCenter(ctx context.Context, cc db.CostCenter) (db.CostCenter, error) {
	c, err := s.repo.Create(ctx, cc)
	if err != nil {
		return c, err
	}

	if err := s.publisher.PublishCostCenterCreated(ctx, &c); err != nil {
		fmt.Printf("Kafka publish error (cost_center.created): %v\n", err)
	}
	return c, nil
}

func (s *CostCenterService) GetCostCenter(ctx context.Context, id uuid.UUID) (db.CostCenter, error) {
	return s.repo.Get(ctx, id)
}

func (s *CostCenterService) ListCostCenters(ctx context.Context, limit, offset int32) ([]db.CostCenter, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *CostCenterService) UpdateCostCenter(ctx context.Context, cc db.CostCenter) (db.CostCenter, error) {
	c, err := s.repo.Update(ctx, cc)
	if err != nil {
		return c, err
	}

	if err := s.publisher.PublishCostCenterUpdated(ctx, &c); err != nil {
		fmt.Printf("Kafka publish error (cost_center.updated): %v\n", err)
	}
	return c, nil
}

func (s *CostCenterService) DeleteCostCenter(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.publisher.PublishCostCenterDeleted(ctx, id.String()); err != nil {
		fmt.Printf("Kafka publish error (cost_center.deleted): %v\n", err)
	}
	return nil
}

type CostAllocationService struct {
	repo      ports.CostAllocationRepository
	publisher ports.EventPublisher
}

func (s *CostAllocationService) AllocateCost(ctx context.Context, ca db.CostAllocation) (db.CostAllocation, error) {
	c, err := s.repo.Allocate(ctx, ca)
	if err != nil {
		return c, err
	}

	if err := s.publisher.PublishCostAllocationAllocated(ctx, &c); err != nil {
		fmt.Printf("Kafka publish error (cost_allocation.allocated): %v\n", err)
	}
	return c, nil
}

func (s *CostAllocationService) ListAllocations(ctx context.Context, limit, offset int32) ([]db.CostAllocation, error) {
	allocs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// optional: publish event for listing
	if err := s.publisher.PublishCostAllocationListed(ctx, allocs); err != nil {
		fmt.Printf("Kafka publish error (cost_allocation.listed): %v\n", err)
	}

	return allocs, nil
}
