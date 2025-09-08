package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type accrualRepository struct {
	q *db.Queries
	publisher ports.EventPublisher
}

func NewAccrualRepository(q *db.Queries, pub ports.EventPublisher) ports.AccrualRepository {
	return &accrualRepository{
		q:         q,
		publisher: pub,
	}
}

func (r *accrualRepository) Create(ctx context.Context, a db.Accrual) (db.Accrual, error) {
	arg := db.CreateAccrualParams{
		Description: a.Description,
		Amount:      a.Amount,
		AccrualDate: a.AccrualDate,
		AccountID:   a.AccountID,
		CreatedBy:   a.CreatedBy,
		UpdatedBy:   a.UpdatedBy,
	}
	res, err := r.q.CreateAccrual(ctx, arg)
	if err != nil {
		return db.Accrual{}, err
	}
	domainAccrual := toDomainAccrual(res)


	return domainAccrual, nil
}

func (r *accrualRepository) Get(ctx context.Context, id uuid.UUID) (db.Accrual, error) {
	res, err := r.q.GetAccrualById(ctx, id)
	if err != nil {
		return db.Accrual{}, err
	}
	return toDomainAccrual(res), nil
}

func (r *accrualRepository) Update(ctx context.Context, a db.Accrual) (db.Accrual, error) {
	arg := db.UpdateAccrualParams{
		ID:          a.ID,
		Description: a.Description,
		Amount:      a.Amount,
		AccrualDate: a.AccrualDate,
		AccountID:   a.AccountID,
		UpdatedBy:   a.UpdatedBy,
	}
	res, err := r.q.UpdateAccrual(ctx, arg)
	if err != nil {
		return db.Accrual{}, err
	}
	domainAccrual := toDomainAccrual(res)

	return domainAccrual, nil
}

func (r *accrualRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteAccrual(ctx, id); err != nil {
		return err
	}

	return nil
}

func (r *accrualRepository) List(ctx context.Context, limit, offset int32) ([]db.Accrual, error) {
	rows, err := r.q.ListAccruals(ctx, db.ListAccrualsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	result := make([]db.Accrual, 0, len(rows))
	result = append(result, rows...)
	return result, nil
}

// --- mapping helpers ---
func toDomainAccrual(a db.Accrual) db.Accrual {
	return db.Accrual{
		ID:          a.ID,
		Description: a.Description,
		Amount:      a.Amount,
		AccrualDate: a.AccrualDate,
		AccountID:   a.AccountID,
		CreatedAt:   a.CreatedAt,
		CreatedBy:   a.CreatedBy,
		UpdatedAt:   a.UpdatedAt,
		UpdatedBy:   a.UpdatedBy,
		Revision:    a.Revision,
	}
}

type AllocationRuleRepository struct {
	q         *db.Queries
	publisher ports.EventPublisher
}

func NewAllocationRuleRepository(q *db.Queries, pub ports.EventPublisher) *AllocationRuleRepository {
	return &AllocationRuleRepository{
		q:         q,
		publisher: pub,
	}
}

func (r *AllocationRuleRepository) Create(ctx context.Context, rule db.AllocationRule) (db.AllocationRule, error) {
	arg := db.CreateAllocationRuleParams{
		Name:                rule.Name,
		Basis:               rule.Basis,
		SourceAccountID:     rule.SourceAccountID,
		TargetCostCenterIds: rule.TargetCostCenterIds,
		Formula:             rule.Formula,
		CreatedBy:           rule.CreatedBy,
		UpdatedBy:           rule.UpdatedBy,
	}
	res, err := r.q.CreateAllocationRule(ctx, arg)
	if err != nil {
		return db.AllocationRule{}, err
	}


	return res, nil
}

func (r *AllocationRuleRepository) Get(ctx context.Context, id uuid.UUID) (db.AllocationRule, error) {
	return r.q.GetAllocationRule(ctx, id)
}

func (r *AllocationRuleRepository) List(ctx context.Context, limit, offset int32) ([]db.AllocationRule, error) {
	return r.q.ListAllocationRules(ctx, db.ListAllocationRulesParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *AllocationRuleRepository) Update(ctx context.Context, rule db.AllocationRule) (db.AllocationRule, error) {
	arg := db.UpdateAllocationRuleParams{
		ID:                  rule.ID,
		Name:                rule.Name,
		Basis:               rule.Basis,
		SourceAccountID:     rule.SourceAccountID,
		TargetCostCenterIds: rule.TargetCostCenterIds,
		Formula:             rule.Formula,
		UpdatedBy:           rule.UpdatedBy,
	}
	res, err := r.q.UpdateAllocationRule(ctx, arg)
	if err != nil {
		return db.AllocationRule{}, err
	}

	return res, nil
}

func (r *AllocationRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteAllocationRule(ctx, id); err != nil {
		return err
	}
	return nil
}
