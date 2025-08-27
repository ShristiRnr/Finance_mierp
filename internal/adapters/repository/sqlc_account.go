package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db" // sqlc generated package
)

type accountRepository struct {
	q *db.Queries
}

func NewAccountRepository(q *db.Queries) ports.AccountRepository {
	return &accountRepository{q: q}
}

func (r *accountRepository) Create(ctx context.Context, a domain.Account) (domain.Account, error) {
	arg := db.CreateAccountParams{
		Code:     a.Code,
		Name:     a.Name,
		Type:     a.Type,
		ParentID: uuid.NullUUID{UUID: *a.ParentID, Valid: a.ParentID != nil},
		Status:   a.Status,
		AllowManualJournal: sql.NullBool{Bool: a.AllowManualJournal, Valid: true},
		CreatedBy: sql.NullString{String: a.CreatedBy, Valid: a.CreatedBy != ""},
		UpdatedBy: sql.NullString{String: a.UpdatedBy, Valid: a.UpdatedBy != ""},
	}
	res, err := r.q.CreateAccount(ctx, arg)
	if err != nil {
		return domain.Account{}, err
	}
	return toDomainAccount(res), nil
}

func (r *accountRepository) Get(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	res, err := r.q.GetAccount(ctx, id)
	if err != nil {
		return domain.Account{}, err
	}
	return toDomainAccount(res), nil
}

func (r *accountRepository) List(ctx context.Context, limit, offset int32) ([]domain.Account, error) {
	res, err := r.q.ListAccounts(ctx, db.ListAccountsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	var accounts []domain.Account
	for _, a := range res {
		accounts = append(accounts, toDomainAccount(a))
	}
	return accounts, nil
}

func (r *accountRepository) Update(ctx context.Context, a domain.Account) (domain.Account, error) {
	arg := db.UpdateAccountParams{
		ID:       a.ID,
		Code:     a.Code,
		Name:     a.Name,
		Type:     a.Type,
		ParentID: uuid.NullUUID{UUID: *a.ParentID, Valid: a.ParentID != nil},
		Status:   a.Status,
		AllowManualJournal: sql.NullBool{Bool: a.AllowManualJournal, Valid: true},
		UpdatedBy: sql.NullString{String: a.UpdatedBy, Valid: a.UpdatedBy != ""},
	}
	res, err := r.q.UpdateAccount(ctx, arg)
	if err != nil {
		return domain.Account{}, err
	}
	return toDomainAccount(res), nil
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteAccount(ctx, id)
}

// helper mapper
func toDomainAccount(a db.Account) domain.Account {
	var parentID *uuid.UUID
	if a.ParentID.Valid {
		parentID = &a.ParentID.UUID
	}
	return domain.Account{
		ID:                 a.ID,
		Code:               a.Code,
		Name:               a.Name,
		Type:               a.Type,
		ParentID:           parentID,
		Status:             a.Status,
		AllowManualJournal: a.AllowManualJournal.Bool,
		CreatedAt:          a.CreatedAt,
		CreatedBy:          a.CreatedBy.String,
		UpdatedAt:          a.UpdatedAt,
		UpdatedBy:          a.UpdatedBy.String,
		Revision:           a.Revision,
	}
}
