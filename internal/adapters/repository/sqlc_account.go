package repository

import (
	"context"
	"database/sql"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
)

// ---------------- Account Repository ----------------
type AccountSQLCRepository struct {
	queries   *db.Queries
	publisher ports.EventPublisher
}

var _ ports.AccountRepository = (*AccountSQLCRepository)(nil)

func NewAccountRepository(q *db.Queries, pub ports.EventPublisher) ports.AccountRepository {
	return &AccountSQLCRepository{
		queries:   q,
		publisher: pub,
	}
}

func (r *AccountSQLCRepository) Create(ctx context.Context, a *db.Account) (*db.Account, error) {
	params := db.CreateAccountParams{
		Code:               a.Code,
		Name:               a.Name,
		Type:               a.Type,
		Status:             a.Status,
		AllowManualJournal: a.AllowManualJournal,
	}
	acc, err := r.queries.CreateAccount(ctx, params)
	if err != nil {
		return nil, err
	}
	domainAcc := &db.Account{
		ID:                 acc.ID,
		Code:               acc.Code,
		Name:               acc.Name,
		Type:               acc.Type,
		Status:             acc.Status,
		AllowManualJournal: acc.AllowManualJournal,
		CreatedAt:          acc.CreatedAt,
		UpdatedAt:          acc.UpdatedAt,
	}
	return domainAcc, nil
}

func (r *AccountSQLCRepository) Get(ctx context.Context, id uuid.UUID) (*db.Account, error) {
	acc, err := r.queries.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	return &db.Account{
		ID:                 acc.ID,
		Code:               acc.Code,
		Name:               acc.Name,
		Type:               acc.Type,
		Status:             acc.Status,
		AllowManualJournal: acc.AllowManualJournal,
		CreatedAt:          acc.CreatedAt,
		UpdatedAt:          acc.UpdatedAt,
	}, nil
}

func (r *AccountSQLCRepository) Update(ctx context.Context, a *db.Account) (*db.Account, error) {
	arg := db.UpdateAccountParams{
		ID:                 a.ID,
		Code:               a.Code,
		Name:               a.Name,
		Type:               a.Type,
		Status:             a.Status,
		AllowManualJournal: a.AllowManualJournal,
	}
	dbAcc, err := r.queries.UpdateAccount(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &db.Account{
		ID:                 dbAcc.ID,
		Code:               dbAcc.Code,
		Name:               dbAcc.Name,
		Type:               dbAcc.Type,
		Status:             dbAcc.Status,
		AllowManualJournal: dbAcc.AllowManualJournal,
		CreatedAt:          dbAcc.CreatedAt,
		UpdatedAt:          dbAcc.UpdatedAt,
	}, nil
}

func (r *AccountSQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteAccount(ctx, id)
}

func (r *AccountSQLCRepository) List(ctx context.Context, limit, offset int32) ([]*db.Account, error) {
	dbAccs, err := r.queries.ListAccounts(ctx, db.ListAccountsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	accounts := make([]*db.Account, len(dbAccs))
	for i := range dbAccs {
		accounts[i] = &dbAccs[i]
	}
	return accounts, nil
}

// ---------------- Journal Repository ----------------

type JournalSQLCRepository struct {
	db        *sql.DB
	queries   *db.Queries
	publisher ports.EventPublisher
}

var _ ports.JournalRepository = (*JournalSQLCRepository)(nil)

func NewJournalRepository(db *sql.DB, q *db.Queries, pub ports.EventPublisher) ports.JournalRepository {
	return &JournalSQLCRepository{
		db:        db,
		queries:   q,
		publisher: pub,
	}
}

func (r *JournalSQLCRepository) Create(ctx context.Context, journal *db.JournalEntry) (*db.JournalEntry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	// Insert entry
	created, err := qtx.CreateJournalEntry(ctx, db.CreateJournalEntryParams{
		JournalDate: journal.JournalDate,
		Reference:   journal.Reference,
		Memo:        journal.Memo,
		SourceType:  journal.SourceType,
		SourceID:    journal.SourceID,
		CreatedBy:   journal.CreatedBy,
		UpdatedBy:   journal.UpdatedBy,
	})
	if err != nil {
		return nil, err
	}

	if err := qtx.DeleteJournalLinesByEntry(ctx, created.ID); err != nil {
		return nil, err
	}

	// Insert lines
	for _, l := range journal.Lines {
		_, err := qtx.CreateJournalLine(ctx, db.CreateJournalLineParams{
			JournalID:   created.ID,
			AccountID:   l.AccountID,
			Side:        l.Side,
			Amount:      l.Amount,
			Description: l.Description,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *JournalSQLCRepository) Get(ctx context.Context, id uuid.UUID) (*db.JournalEntry, error) {
	entry, err := r.queries.GetJournalEntry(ctx, id)
	if err != nil {
		return nil, err
	}
	lines, err := r.queries.ListJournalLinesByEntry(ctx, id)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines
	return &entry, nil
}

func (r *JournalSQLCRepository) Update(ctx context.Context, journal *db.JournalEntry) (*db.JournalEntry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	// Update journal entry
	updated, err := qtx.UpdateJournalEntry(ctx, db.UpdateJournalEntryParams{
		ID:          journal.ID, 
		JournalDate: journal.JournalDate,
		Reference:   journal.Reference,
		Memo:        journal.Memo,
		SourceType:  journal.SourceType,
		SourceID:    journal.SourceID,
		UpdatedBy:   journal.UpdatedBy,// ✅ important!
	})
	if err != nil {
		return nil, err
	}

	// ✅ Correct journal.ID passed here
	if err := qtx.DeleteJournalLinesByEntry(ctx, journal.ID); err != nil {
		return nil, err
	}

	// Re-insert lines
	for _, l := range journal.Lines {
		_, err := qtx.CreateJournalLine(ctx, db.CreateJournalLineParams{
			JournalID:   journal.ID,
			AccountID:   l.AccountID,
			Side:        l.Side,
			Amount:      l.Amount,
			Description: l.Description,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *JournalSQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteJournalEntry(ctx, id)
}

func (r *JournalSQLCRepository) List(ctx context.Context, limit, offset int32) ([]*db.JournalEntry, error) {
	dbJEs, err := r.queries.ListJournalEntries(ctx, db.ListJournalEntriesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	journals := make([]*db.JournalEntry, len(dbJEs))
	for i := range dbJEs {
		journals[i] = &dbJEs[i]
	}
	return journals, nil
}

// ---------------- Ledger Repository ----------------

type LedgerSQLCRepository struct {
	queries *db.Queries
}

var _ ports.LedgerRepository = (*LedgerSQLCRepository)(nil)

func NewLedgerRepository(q *db.Queries) ports.LedgerRepository {
	return &LedgerSQLCRepository{queries: q}
}

func (r *LedgerSQLCRepository) List(ctx context.Context, limit, offset int32) ([]*db.LedgerEntry, error) {
	rows, err := r.queries.ListLedgerEntries(ctx, db.ListLedgerEntriesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	entries := make([]*db.LedgerEntry, len(rows))
	for i := range rows {
		entries[i] = &rows[i]
	}
	return entries, nil
}
