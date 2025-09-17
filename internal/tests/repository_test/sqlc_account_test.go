package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
	"github.com/google/uuid"
)

// ---------------- Account Repository ----------------

func TestAccountRepository_Create_Get_Update_Delete_List(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	queries := db.New(dbConn)
	repo := repository.NewAccountRepository(queries, nil)

	ctx := context.Background()
	id := uuid.New()
	now := time.Now()

	// ---------- CREATE ----------
	mock.ExpectQuery(`(?s)INSERT INTO accounts`).
		WithArgs("A100", "Cash", "Asset", sqlmock.AnyArg(), "Active", true,
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "code", "name", "type", "parent_id", "status", "allow_manual_journal",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			id, "A100", "Cash", "Asset", nil, "Active", true,
			now, sql.NullString{String: "tester", Valid: true},
			now, sql.NullString{String: "tester", Valid: true},
			1,
		))

	acc := &db.Account{
		Code:               "A100",
		Name:               "Cash",
		Type:               "Asset",
		Status:             "Active",
		AllowManualJournal: true,
	}
	created, err := repo.Create(ctx, acc)
	require.NoError(t, err)
	require.Equal(t, "Cash", created.Name)

	// ---------- GET ----------
	mock.ExpectQuery(`(?s)SELECT (.+) FROM accounts WHERE id =`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "code", "name", "type", "parent_id", "status", "allow_manual_journal",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			id, "A100", "Cash", "Asset", nil, "Active", true,
			now, sql.NullString{String: "tester", Valid: true},
			now, sql.NullString{String: "tester", Valid: true},
			1,
		))

	got, err := repo.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, created.ID, got.ID)

	// ---------- UPDATE ----------
	mock.ExpectQuery(`(?s)UPDATE accounts`).
		WithArgs(
			id,               // $1 id
			"A100",           // $2 code
			"Cash-Updated",   // $3 name
			"Asset",          // $4 type
			nil,              // $5 parent_id
			"Active",         // $6 status
			true,             // $7 allow_manual_journal
			sqlmock.AnyArg(), // $8 updated_by
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "code", "name", "type", "parent_id", "status", "allow_manual_journal",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			id, "A100", "Cash-Updated", "Asset", nil, "Active", true,
			now, sql.NullString{String: "tester", Valid: true},
			time.Now(), sql.NullString{String: "tester", Valid: true},
			2,
		))

	updated, err := repo.Update(ctx, &db.Account{
		ID:                 id,
		Code:               "A100",
		Name:               "Cash-Updated",
		Type:               "Asset",
		Status:             "Active",
		AllowManualJournal: true,
	})
	require.NoError(t, err)
	require.Equal(t, "Cash-Updated", updated.Name)

	// ---------- LIST ----------
	mock.ExpectQuery(`(?s)SELECT (.+) FROM accounts ORDER BY code LIMIT`).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "code", "name", "type", "parent_id", "status", "allow_manual_journal",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).
			AddRow(id, "A100", "Cash-Updated", "Asset", nil, "Active", true,
				now, sql.NullString{String: "tester", Valid: true},
				now, sql.NullString{String: "tester", Valid: true}, 2),
		)

	list, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "Cash-Updated", list[0].Name)

	// ---------- DELETE ----------
	mock.ExpectExec(`(?s)DELETE FROM accounts WHERE id =`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ---------------- Journal Repository ----------------

func TestJournalRepository_Create_Get_Update_Delete_List(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	queries := db.New(dbConn)
	repo := repository.NewJournalRepository(dbConn, queries, nil)

	ctx := context.Background()
	journalID := uuid.New()
	lineID := uuid.New()
	accountID := uuid.New()
	now := time.Now()

	// ---------- CREATE ----------
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO journal_entries`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "journal_date", "reference", "memo", "source_type", "source_id",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			journalID, now,
			sql.NullString{String: "REF-1", Valid: true},
			sql.NullString{String: "Test", Valid: true},
			sql.NullString{String: "System", Valid: true},
			sql.NullString{String: "SRC123", Valid: true},
			now, sql.NullString{String: "tester", Valid: true},
			now, sql.NullString{String: "tester", Valid: true},
			1,
		))

	mock.ExpectExec(`DELETE FROM journal_lines WHERE journal_id = \$1`).
		WithArgs(journalID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`INSERT INTO journal_lines`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "journal_id", "account_id", "side", "amount", "cost_center_id", "description", "created_at",
		}).AddRow(
			lineID, journalID, accountID, "Debit", "100", nil,
			sql.NullString{String: "Line 1", Valid: true}, now,
		))

	mock.ExpectCommit()

	journal := &db.JournalEntry{
		JournalDate: now,
		Reference:   sql.NullString{String: "REF-1", Valid: true},
		Memo:        sql.NullString{String: "Test", Valid: true},
		SourceType:  sql.NullString{String: "System", Valid: true},
		SourceID:    sql.NullString{String: "SRC123", Valid: true},
		CreatedBy:   sql.NullString{String: "tester", Valid: true},
		UpdatedBy:   sql.NullString{String: "tester", Valid: true},
		Lines: []db.JournalLine{{
			AccountID:   accountID,
			Side:        "Debit",
			Amount:      "100",
			Description: sql.NullString{String: "Line 1", Valid: true},
		}},
	}

	created, err := repo.Create(ctx, journal)
	require.NoError(t, err)
	require.Equal(t, "REF-1", created.Reference.String)

	// ---------- GET ----------
	mock.ExpectQuery(`SELECT .* FROM journal_entries WHERE id = \$1`).
		WithArgs(journalID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "journal_date", "reference", "memo", "source_type", "source_id",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			journalID, now,
			sql.NullString{String: "REF-1", Valid: true},
			sql.NullString{String: "Test", Valid: true},
			sql.NullString{String: "System", Valid: true},
			sql.NullString{String: "SRC123", Valid: true},
			now, sql.NullString{String: "tester", Valid: true},
			now, sql.NullString{String: "tester", Valid: true},
			1,
		))

	mock.ExpectQuery(`SELECT .* FROM journal_lines WHERE journal_id = \$1`).
		WithArgs(journalID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "journal_id", "account_id", "side", "amount", "cost_center_id", "description", "created_at",
		}).AddRow(
			lineID, journalID, accountID, "Debit", "100", nil,
			sql.NullString{String: "Line 1", Valid: true}, now,
		))

	got, err := repo.Get(ctx, journalID)
	require.NoError(t, err)
	require.Equal(t, created.ID, got.ID)

	// ---------- UPDATE ----------
	mock.ExpectBegin()
	mock.ExpectQuery(`UPDATE journal_entries`).
		WithArgs(
			journalID,
			sqlmock.AnyArg(), // journal_date can vary
			sql.NullString{String: "REF-2", Valid: true},
			sql.NullString{String: "Updated", Valid: true},
			sql.NullString{String: "System", Valid: true},
			sql.NullString{String: "SRC123", Valid: true},
			sql.NullString{String: "tester", Valid: true},
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "journal_date", "reference", "memo", "source_type", "source_id",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			journalID, now,
			sql.NullString{String: "REF-2", Valid: true},
			sql.NullString{String: "Updated", Valid: true},
			sql.NullString{String: "System", Valid: true},
			sql.NullString{String: "SRC123", Valid: true},
			now, sql.NullString{String: "tester", Valid: true},
			time.Now(), sql.NullString{String: "tester", Valid: true},
			2,
		))

	mock.ExpectExec(`DELETE FROM journal_lines WHERE journal_id = \$1`).
		WithArgs(journalID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`INSERT INTO journal_lines`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "journal_id", "account_id", "side", "amount", "cost_center_id", "description", "created_at",
		}).AddRow(
			uuid.New(), journalID, accountID, "Credit", "200", nil,
			sql.NullString{String: "Line 2", Valid: true}, now,
		))

	mock.ExpectCommit()

	updated, err := repo.Update(ctx, &db.JournalEntry{
		ID:          journalID,
		JournalDate: now,
		Reference:   sql.NullString{String: "REF-2", Valid: true},
		Memo:        sql.NullString{String: "Updated", Valid: true},
		SourceType:  sql.NullString{String: "System", Valid: true},
		SourceID:    sql.NullString{String: "SRC123", Valid: true},
		UpdatedBy:   sql.NullString{String: "tester", Valid: true},
		Lines: []db.JournalLine{{
			AccountID:   accountID,
			Side:        "Credit",
			Amount:      "200",
			Description: sql.NullString{String: "Line 2", Valid: true},
		}},
	})
	require.NoError(t, err)
	require.Equal(t, "REF-2", updated.Reference.String)

	// ---------- LIST ----------
	mock.ExpectQuery(`SELECT .* FROM journal_entries ORDER BY journal_date.*LIMIT \$1 OFFSET \$2`).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "journal_date", "reference", "memo", "source_type", "source_id",
			"created_at", "created_by", "updated_at", "updated_by", "revision",
		}).AddRow(
			journalID, now,
			sql.NullString{String: "REF-2", Valid: true},
			sql.NullString{String: "Updated", Valid: true},
			sql.NullString{String: "System", Valid: true},
			sql.NullString{String: "SRC123", Valid: true},
			now, sql.NullString{String: "tester", Valid: true},
			now, sql.NullString{String: "tester", Valid: true}, 2,
		))

	list, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "REF-2", list[0].Reference.String)

	// ---------- DELETE ----------
	mock.ExpectExec(`DELETE FROM journal_entries WHERE id = \$1`).
		WithArgs(journalID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(ctx, journalID)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}


// ---------------- Ledger Repository ----------------

func TestLedgerRepository_List(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	queries := db.New(dbConn)
	repo := repository.NewLedgerRepository(queries)

	ctx := context.Background()
	now := time.Now()

	// mock select (must match sqlc return columns)
	mock.ExpectQuery(`(?s)SELECT (.+) FROM ledger_entries`).
		WillReturnRows(sqlmock.NewRows([]string{
			"entry_id", "account_id", "side", "amount", "transaction_date",
		}).AddRow(
			uuid.New(), uuid.New(), "Debit", "500", now,
		))

	list, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, list)

	require.NoError(t, mock.ExpectationsWereMet())
}
