package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
)

func newTestRepo(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *repository.CreditDebitNoteRepo) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)

	queries := db.New(dbConn)
	repo := repository.NewCreditDebitNoteRepo(queries).(*repository.CreditDebitNoteRepo)

	return dbConn, mock, repo
}

func TestCreditDebitNoteRepo_Create(t *testing.T) {
	dbConn, mock, repo := newTestRepo(t)
	defer dbConn.Close()

	ctx := context.Background()
	id := uuid.New()
	invoiceID := uuid.New()
	now := time.Now()

	note := db.CreditDebitNote{
		InvoiceID: invoiceID,
		Type:      "credit",
		Amount:    "1000.00",
		Reason:    sql.NullString{String: "discount", Valid: true},
		CreatedBy: sql.NullString{String: "user1", Valid: true},
		UpdatedBy: sql.NullString{String: "user1", Valid: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "invoice_id", "type", "amount", "reason",
		"created_at", "created_by", "updated_at", "updated_by", "revision",
	}).AddRow(id, invoiceID, "credit", "1000.00", "discount", now, "user1", now, "user1", 1)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO credit_debit_notes`)).
		WithArgs(note.InvoiceID, note.Type, note.Amount, note.Reason, note.CreatedBy, note.UpdatedBy).
		WillReturnRows(rows)

	created, err := repo.Create(ctx, note)
	require.NoError(t, err)
	require.Equal(t, id, created.ID)
	require.Equal(t, "credit", created.Type)
	require.Equal(t, "1000.00", created.Amount)
	require.Equal(t, "discount", created.Reason.String)
}

func TestCreditDebitNoteRepo_Get(t *testing.T) {
	dbConn, mock, repo := newTestRepo(t)
	defer dbConn.Close()

	ctx := context.Background()
	id := uuid.New()
	invoiceID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "invoice_id", "type", "amount", "reason",
		"created_at", "created_by", "updated_at", "updated_by", "revision",
	}).AddRow(id, invoiceID, "debit", "500.00", "adjustment", now, "user1", now, "user2", 1)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, invoice_id, type, amount, reason, created_at, created_by, updated_at, updated_by, revision 
 FROM credit_debit_notes WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	got, err := repo.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "debit", got.Type)
	require.Equal(t, "500.00", got.Amount)
	require.Equal(t, "adjustment", got.Reason.String)
}

func TestCreditDebitNoteRepo_List(t *testing.T) {
	dbConn, mock, repo := newTestRepo(t)
	defer dbConn.Close()

	ctx := context.Background()
	id1 := uuid.New()
	id2 := uuid.New()
	invoiceID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "invoice_id", "type", "amount", "reason",
		"created_at", "created_by", "updated_at", "updated_by", "revision",
	}).AddRow(id1, invoiceID, "credit", "100.00", "promo", now, "u1", now, "u1", 1).
		AddRow(id2, invoiceID, "debit", "200.00", "adjustment", now, "u2", now, "u2", 1)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, invoice_id, type, amount, reason, created_at, created_by, updated_at, updated_by, revision 
 FROM credit_debit_notes
 ORDER BY created_at DESC
 LIMIT $1 OFFSET $2`)).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(rows)

	list, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, "credit", list[0].Type)
	require.Equal(t, "debit", list[1].Type)
}

func TestCreditDebitNoteRepo_Update(t *testing.T) {
	dbConn, mock, repo := newTestRepo(t)
	defer dbConn.Close()

	ctx := context.Background()
	id := uuid.New()
	invoiceID := uuid.New()
	now := time.Now()

	note := db.CreditDebitNote{
		ID:        id,
		InvoiceID: invoiceID,
		Type:      "credit",
		Amount:    "1500.00",
		Reason:    sql.NullString{String: "correction", Valid: true},
		UpdatedBy: sql.NullString{String: "user2", Valid: true},
	}

	rows := sqlmock.NewRows([]string{
		"id", "invoice_id", "type", "amount", "reason",
		"created_at", "created_by", "updated_at", "updated_by", "revision",
	}).AddRow(id, invoiceID, "credit", "1500.00", "correction", now, "u1", now, "user2", 2)

	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE credit_debit_notes
 SET invoice_id = $2, type = $3, amount = $4, reason = $5, updated_by = $6, revision = revision + 1, updated_at = now()
 WHERE id = $1
 RETURNING id, invoice_id, type, amount, reason, created_at, created_by, updated_at, updated_by, revision`)).
		WithArgs(note.ID, note.InvoiceID, note.Type, note.Amount, note.Reason, note.UpdatedBy).
		WillReturnRows(rows)

	updated, err := repo.Update(ctx, note)
	require.NoError(t, err)
	require.Equal(t, "1500.00", updated.Amount)
	require.Equal(t, "correction", updated.Reason.String)
}

func TestCreditDebitNoteRepo_Delete(t *testing.T) {
	dbConn, mock, repo := newTestRepo(t)
	defer dbConn.Close()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM credit_debit_notes WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(ctx, id)
	require.NoError(t, err)
}
