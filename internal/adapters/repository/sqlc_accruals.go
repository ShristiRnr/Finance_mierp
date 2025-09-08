package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type accrualRepository struct {
	q *db.Queries
}

func NewAccrualRepository(q *db.Queries) ports.AccrualRepository {
	return &accrualRepository{q: q}
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
	return toDomainAccrual(res), nil
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
	return toDomainAccrual(res), nil
}

func (r *accrualRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteAccrual(ctx, id)
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
