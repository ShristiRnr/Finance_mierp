package repository

import (
	"context"
	"database/sql"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db" // sqlc generated package
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
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
	accounts, err := r.q.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.Account, 0, len(accounts))
	for _, acc := range accounts {
		result = append(result, domain.Account{
			ID:     acc.ID,
			Code:   acc.Code,
			Name:   acc.Name,
			Type:   acc.Type,
			Status: acc.Status,
		})
	}
	return result, nil
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
		CreatedAt:          a.CreatedAt.Time,
		CreatedBy:          a.CreatedBy.String,
		UpdatedAt:          a.UpdatedAt.Time,
		UpdatedBy:          a.UpdatedBy.String,
		Revision:           a.Revision,
	}
}

type journalRepository struct {
	q *db.Queries
}

func NewJournalRepository(q *db.Queries) ports.JournalRepository {
	return &journalRepository{q: q}
}

func (r *journalRepository) Create(ctx context.Context, j domain.JournalEntry) (domain.JournalEntry, error) {
	arg := db.CreateJournalEntryParams{
		JournalDate: j.JournalDate,
		Reference:   sql.NullString{String: derefString(j.Reference), Valid: j.Reference != nil},
		Memo:        sql.NullString{String: derefString(j.Memo), Valid: j.Memo != nil},
		SourceType:  sql.NullString{String: derefString(j.SourceType), Valid: j.SourceType != nil},
		SourceID:    sql.NullString{String: derefString(j.SourceID), Valid: j.SourceID != nil},
		CreatedBy:   sql.NullString{String: j.CreatedBy, Valid: j.CreatedBy != ""},
		UpdatedBy:   sql.NullString{String: j.UpdatedBy, Valid: j.UpdatedBy != ""},
	}
	res, err := r.q.CreateJournalEntry(ctx, arg)
	if err != nil {
		return domain.JournalEntry{}, err
	}
	return toDomainJournal(res), nil
}

func (r *journalRepository) Get(ctx context.Context, id uuid.UUID) (domain.JournalEntry, error) {
	res, err := r.q.GetJournalEntry(ctx, id)
	if err != nil {
		return domain.JournalEntry{}, err
	}
	return toDomainJournal(res), nil
}

func (r *journalRepository) List(ctx context.Context, limit, offset int32) ([]domain.JournalEntry, error) {
	rows, err := r.q.ListJournalEntries(ctx, db.ListJournalEntriesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.JournalEntry, 0, len(rows))
	for _, rj := range rows {
		result = append(result, toDomainJournal(rj))
	}
	return result, nil
}

func (r *journalRepository) Update(ctx context.Context, j domain.JournalEntry) (domain.JournalEntry, error) {
	arg := db.UpdateJournalEntryParams{
		ID:         j.ID,
		JournalDate: j.JournalDate,
		Reference:   sql.NullString{String: derefString(j.Reference), Valid: j.Reference != nil},
		Memo:        sql.NullString{String: derefString(j.Memo), Valid: j.Memo != nil},
		SourceType:  sql.NullString{String: derefString(j.SourceType), Valid: j.SourceType != nil},
		SourceID:    sql.NullString{String: derefString(j.SourceID), Valid: j.SourceID != nil},
		UpdatedBy:   sql.NullString{String: j.UpdatedBy, Valid: j.UpdatedBy != ""},
	}
	res, err := r.q.UpdateJournalEntry(ctx, arg)
	if err != nil {
		return domain.JournalEntry{}, err
	}
	return toDomainJournal(res), nil
}

func (r *journalRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteJournalEntry(ctx, id)
}

// mapper
func toDomainJournal(j db.JournalEntry) domain.JournalEntry {
	return domain.JournalEntry{
		ID:         j.ID,
		JournalDate: j.JournalDate,
		Reference:   nullStringToPtr(j.Reference),
		Memo:        nullStringToPtr(j.Memo),
		SourceType:  nullStringToPtr(j.SourceType),
		SourceID:    nullStringToPtr(j.SourceID),
		CreatedAt:   j.CreatedAt.Time,
		CreatedBy:   j.CreatedBy.String,
		UpdatedAt:   j.UpdatedAt.Time,
		UpdatedBy:   j.UpdatedBy.String,
		Revision:    j.Revision,
	}
}

type ledgerRepository struct {
	q *db.Queries
}

func NewLedgerRepository(q *db.Queries) ports.LedgerRepository {
	return &ledgerRepository{q: q}
}

func (r *ledgerRepository) List(ctx context.Context, limit, offset int32) ([]domain.LedgerEntry, error) {
	rows, err := r.q.ListLedgerEntries(ctx, db.ListLedgerEntriesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.LedgerEntry, 0, len(rows))
	for _, l := range rows {
		result = append(result, domain.LedgerEntry{
			EntryID:   l.ID,
			AccountID: l.AccountID,
			Side:      l.Side,
			Amount:    l.Amount,
			PostedAt:  l.TransactionDate,
		})
	}
	return result, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}