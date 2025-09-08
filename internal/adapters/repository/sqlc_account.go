package repository

import (
	"context"
	"database/sql"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
)

type AccountSQLCRepository struct {
	queries *db.Queries
	publisher ports.EventPublisher
}

var _ ports.AccountRepository = (*AccountSQLCRepository)(nil)
func NewAccountRepository(q *db.Queries, pub ports.EventPublisher) ports.AccountRepository {
	return &AccountSQLCRepository{
		queries:  q,
		publisher: pub,
	}
}

func (r *AccountSQLCRepository) Create(ctx context.Context, a db.Account) (db.Account, error) {
	params := db.CreateAccountParams{
		Code:               a.Code,
		Name:               a.Name,
		Type:               a.Type,
		Status:             a.Status,
		AllowManualJournal: a.AllowManualJournal,
	}
	acc, err := r.queries.CreateAccount(ctx, params)
	if err != nil {
		return db.Account{}, err
	}
	domainAcc := toDomainAccount(acc)
	// Publish Kafka event
	_ = r.publisher.PublishAccountCreated(ctx, domainAcc)

	return domainAcc, nil
}

func (r *AccountSQLCRepository) Get(ctx context.Context, id uuid.UUID) (db.Account, error) {
	acc, err := r.queries.GetAccount(ctx, id)
	if err != nil {
		return db.Account{}, err
	}
	return toDomainAccount(acc), nil
}

func (r *AccountSQLCRepository) Update(ctx context.Context, a db.Account) (db.Account, error) {
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
		return db.Account{}, err
	}

	domainAcc := toDomainAccount(dbAcc)
	// Publish Kafka event
	_ = r.publisher.PublishAccountUpdated(ctx, domainAcc)

	return domainAcc, nil
}

func (r *AccountSQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteAccount(ctx, id); err != nil {
		return err
	}
	// Publish Kafka event
	_ = r.publisher.PublishAccountDeleted(ctx, id.String())
	return nil
}

func (r *AccountSQLCRepository) List(ctx context.Context, limit, offset int32) ([]db.Account, error) {
	dbAccs, err := r.queries.ListAccounts(ctx, db.ListAccountsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	accounts := make([]db.Account, len(dbAccs))
	copy(accounts, dbAccs)
	return accounts, nil
}

// ---------------- Journals ----------------

type JournalSQLCRepository struct {
	db        *sql.DB
	queries   *db.Queries
	publisher ports.EventPublisher
}

func NewJournalRepository(db *sql.DB, q *db.Queries, pub ports.EventPublisher) ports.JournalRepository {
	return &JournalSQLCRepository{
		db:        db,
		queries:   q,
		publisher: pub,
	}
}

// Create header + lines in a transaction
func (r *JournalSQLCRepository) Create(ctx context.Context, j db.JournalEntry) (db.JournalEntry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return db.JournalEntry{}, err
	}
	defer tx.Rollback()

	txQueries := db.New(tx)
	entry, err := txQueries.CreateJournalEntry(ctx, db.CreateJournalEntryParams{
		JournalDate: j.JournalDate,
		Reference:   j.Reference,
		Memo:        j.Memo,
		SourceType:  j.SourceType,
		SourceID:    j.SourceID,
		CreatedBy:   j.CreatedBy,
		UpdatedBy:   j.UpdatedBy,
	})
	if err != nil {
		return db.JournalEntry{}, err
	}
	if err := txQueries.DeleteJournalLinesByEntry(ctx, j.ID); err != nil {
		return db.JournalEntry{}, err
	}

	lines := make([]db.JournalLine, len(j.Lines))
	for i, line := range j.Lines {
		dbLine, err := txQueries.CreateJournalLine(ctx, db.CreateJournalLineParams{
			JournalID:    j.ID,
			AccountID:    line.AccountID,
			Side:         line.Side,
			Amount:       line.Amount,
			CostCenterID: line.CostCenterID,
			Description:  line.Description,
		})
		if err != nil {
			return db.JournalEntry{}, err
		}
		lines[i] = dbLine
	}

	if err := tx.Commit(); err != nil {
		return db.JournalEntry{}, err
	}

	domainEntry := toDomainJournalWithLines(entry, lines)
	// Publish Kafka event
	_ = r.publisher.PublishJournalCreated(ctx, domainEntry)

	return domainEntry, nil
}

func (r *JournalSQLCRepository) Get(ctx context.Context, id uuid.UUID) (db.JournalEntry, error) {
	entry, err := r.queries.GetJournalEntry(ctx, id)
	if err != nil {
		return db.JournalEntry{}, err
	}
	lines, err := r.queries.ListJournalLinesByEntry(ctx, id)
	if err != nil {
		return db.JournalEntry{}, err
	}
	domainLines := make([]db.JournalLine, len(lines))
	copy(domainLines, lines)
	return toDomainJournalWithLines(entry, domainLines), nil
}

// Update header, handle lines separately (delete + recreate)
func (r *JournalSQLCRepository) Update(ctx context.Context, j db.JournalEntry) (db.JournalEntry, error) {
	_, err := r.queries.UpdateJournalEntry(ctx, db.UpdateJournalEntryParams{
		ID:          j.ID,
		JournalDate: j.JournalDate,
		Reference:   j.Reference,
		Memo:        j.Memo,
		SourceType:  j.SourceType,
		SourceID:    j.SourceID,
	})
	if err != nil {
		return db.JournalEntry{}, err
	}

	if err := r.queries.DeleteJournalLinesByEntry(ctx, j.ID); err != nil {
		return db.JournalEntry{}, err
	}
	lines := make([]db.JournalLine, len(j.Lines))
	for i, line := range j.Lines {
		dbLine, err := r.queries.CreateJournalLine(ctx, db.CreateJournalLineParams{
			JournalID:    j.ID,
			AccountID:    line.AccountID,
			Side:         line.Side,
			Amount:       line.Amount,
			CostCenterID: line.CostCenterID,
			Description:  line.Description,
		})
		if err != nil {
			return db.JournalEntry{}, err
		}
		lines[i] = toDomainJournalLine(dbLine)
	}

	entry, err := r.queries.GetJournalEntry(ctx, j.ID)
	if err != nil {
		return db.JournalEntry{}, err
	}
	domainEntry := toDomainJournalWithLines(entry, lines)
	// Publish Kafka event
	_ = r.publisher.PublishJournalUpdated(ctx, domainEntry)

	return domainEntry, nil
}

func (r *JournalSQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteJournalEntry(ctx, id); err != nil {
		return err
	}
	_ = r.publisher.PublishJournalDeleted(ctx, id.String())
	return nil
}

func (r *JournalSQLCRepository) List(ctx context.Context, limit, offset int32) ([]db.JournalEntry, error) {
	dbJEs, err := r.queries.ListJournalEntries(ctx, db.ListJournalEntriesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	journals := make([]db.JournalEntry, len(dbJEs))
	copy(journals, dbJEs)
	return journals, nil
}

// ---------------- Ledger ----------------

type LedgerSQLCRepository struct {
	queries *db.Queries
}

func NewLedgerRepository(q *db.Queries) ports.LedgerRepository {
	return &LedgerSQLCRepository{queries: q}
}

func (r *LedgerSQLCRepository) List(ctx context.Context, limit, offset int32) ([]db.LedgerEntry, error) {
	rows, err := r.queries.ListLedgerEntries(ctx, db.ListLedgerEntriesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	entries := make([]db.LedgerEntry, len(rows))
	copy(entries, rows)
	return entries, nil
}

// ---------------- mapping helpers ----------------

func toDomainAccount(a db.Account) db.Account {
	return db.Account{
		ID:                 a.ID,
		Code:               a.Code,
		Name:               a.Name,
		Type:               a.Type,
		Status:             a.Status,
		AllowManualJournal: a.AllowManualJournal,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

func toDomainJournal(j db.JournalEntry) db.JournalEntry {
	return db.JournalEntry{
		ID:          j.ID,
		JournalDate: j.JournalDate,
		Reference:   j.Reference,
		Memo:        j.Memo,
		SourceType:  j.SourceType,
		SourceID:    j.SourceID,
		CreatedAt:   j.CreatedAt,
		CreatedBy:   j.CreatedBy,
		UpdatedAt:   j.UpdatedAt,
		UpdatedBy:   j.UpdatedBy,
		Revision:    j.Revision,
	}
}

func toDomainJournalWithLines(j db.JournalEntry, lines []db.JournalLine) db.JournalEntry {
	je := toDomainJournal(j)
	je.Lines = lines
	return je
}

func toDomainJournalLine(l db.JournalLine) db.JournalLine {
	return db.JournalLine{
		AccountID:    l.AccountID,
		Side:         l.Side,
		Amount:       l.Amount,
		CostCenterID: l.CostCenterID,
		Description:  l.Description,
		CreatedAt:    l.CreatedAt,
	}
}
