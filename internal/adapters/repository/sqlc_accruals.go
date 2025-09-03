package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type accrualRepository struct {
	q *db.Queries
}

func NewAccrualRepository(q *db.Queries) ports.AccrualRepository {
	return &accrualRepository{q: q}
}

func (r *accrualRepository) Create(ctx context.Context, a domain.Accrual) (domain.Accrual, error) {
	arg := db.CreateAccrualParams{
		Description: sql.NullString{String: derefString(a.Description), Valid: a.Description != nil},
		Amount:      a.Amount,
		AccrualDate: a.AccrualDate,
		AccountID:   a.AccountID,
		CreatedBy:   sql.NullString{String: a.CreatedBy, Valid: a.CreatedBy != ""},
		UpdatedBy:   sql.NullString{String: a.UpdatedBy, Valid: a.UpdatedBy != ""},
	}
	res, err := r.q.CreateAccrual(ctx, arg)
	if err != nil {
		return domain.Accrual{}, err
	}
	return toDomainAccrual(res), nil
}

func (r *accrualRepository) Get(ctx context.Context, id uuid.UUID) (domain.Accrual, error) {
	res, err := r.q.GetAccrualById(ctx, id)
	if err != nil {
		return domain.Accrual{}, err
	}
	return toDomainAccrual(res), nil
}

func (r *accrualRepository) Update(ctx context.Context, a domain.Accrual) (domain.Accrual, error) {
	arg := db.UpdateAccrualParams{
		ID:          a.ID,
		Description: sql.NullString{String: derefString(a.Description), Valid: a.Description != nil},
		Amount:      a.Amount,
		AccrualDate: a.AccrualDate,
		AccountID:   a.AccountID,
		UpdatedBy:   sql.NullString{String: a.UpdatedBy, Valid: a.UpdatedBy != ""},
	}
	res, err := r.q.UpdateAccrual(ctx, arg)
	if err != nil {
		return domain.Accrual{}, err
	}
	return toDomainAccrual(res), nil
}

func (r *accrualRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteAccrual(ctx, id)
}

func (r *accrualRepository) List(ctx context.Context, limit, offset int32) ([]domain.Accrual, error) {
	rows, err := r.q.ListAccruals(ctx, db.ListAccrualsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	result := make([]domain.Accrual, 0, len(rows))
	for _, row := range rows {
		result = append(result, toDomainAccrual(row))
	}
	return result, nil
}

func (r *accrualRepository) ListExternalRefs(ctx context.Context, accrualID uuid.UUID) ([]domain.AccrualExternalRef, error) {
	rows, err := r.q.ListAccrualExternalRefs(ctx, accrualID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.AccrualExternalRef, 0, len(rows))
	for _, row := range rows {
		result = append(result, toDomainAccrualRef(row))
	}
	return result, nil
}

func (r *accrualRepository) AddExternalRef(ctx context.Context, ref domain.AccrualExternalRef) (domain.AccrualExternalRef, error) {
	arg := db.AddAccrualExternalRefParams{
		AccrualID: ref.AccrualID,
		System:    ref.System,
		RefID:     ref.RefID,
	}
	res, err := r.q.AddAccrualExternalRef(ctx, arg)
	if err != nil {
		return domain.AccrualExternalRef{}, err
	}
	return toDomainAccrualRef(res), nil
}

// --- mapping helpers ---
func toDomainAccrual(a db.Accrual) domain.Accrual {
	return domain.Accrual{
		ID:          a.ID,
		Description: nullStringToPtr(a.Description),
		Amount:      a.Amount,
		AccrualDate: a.AccrualDate,
		AccountID:   a.AccountID,
		CreatedAt:   a.CreatedAt.Time,
		CreatedBy:   a.CreatedBy.String,
		UpdatedAt:   a.UpdatedAt.Time,
		UpdatedBy:   a.UpdatedBy.String,
		Revision:    a.Revision,
	}
}

func toDomainAccrualRef(r db.AccrualExternalRef) domain.AccrualExternalRef {
	return domain.AccrualExternalRef{
		ID:        r.ID,
		AccrualID: r.AccrualID,
		System:    r.System,
		RefID:     r.RefID,
		CreatedAt: r.CreatedAt.Time,
	}
}
