package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

var (
	ErrInvalidName = errors.New("invalid allocation rule name")
)

type AccrualService struct {
	repo ports.AccrualRepository
}

func NewAccrualService(r ports.AccrualRepository) *AccrualService {
	return &AccrualService{repo: r}
}

func (s *AccrualService) Create(ctx context.Context, a db.Accrual) (db.Accrual, error) {
	return s.repo.Create(ctx, a)
}

func (s *AccrualService) Get(ctx context.Context, id uuid.UUID) (db.Accrual, error) {
	return s.repo.Get(ctx, id)
}

func (s *AccrualService) Update(ctx context.Context, a db.Accrual) (db.Accrual, error) {
	return s.repo.Update(ctx, a)
}

func (s *AccrualService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *AccrualService) List(ctx context.Context, limit, offset int32) ([]db.Accrual, error) {
	return s.repo.List(ctx, limit, offset)
}

type allocationService struct {
	repo ports.AllocationRuleRepository
}

func NewAllocationService(r ports.AllocationRuleRepository) ports.AllocationService {
	return &allocationService{repo: r}
}

func (s *allocationService) CreateRule(ctx context.Context, r db.AllocationRule) (db.AllocationRule, error) {
	if r.Name == "" {
		return db.AllocationRule{}, ErrInvalidName
	}
	return s.repo.Create(ctx, r)
}

func (s *allocationService) GetRule(ctx context.Context, id uuid.UUID) (db.AllocationRule, error) {
	return s.repo.Get(ctx, id)
}

func (s *allocationService) ListRules(ctx context.Context, limit, offset int32) ([]db.AllocationRule, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *allocationService) UpdateRule(ctx context.Context, r db.AllocationRule) (db.AllocationRule, error) {
	return s.repo.Update(ctx, r)
}

func (s *allocationService) DeleteRule(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *allocationService) ApplyRule(ctx context.Context, ruleID uuid.UUID) error {
	return nil
}