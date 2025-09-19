package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	db "github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
)

func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *db.Queries) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	queries := db.New(sqlDB)
	return sqlDB, mock, queries
}

// ==================== PAYMENT DUE TESTS ====================

func TestPaymentDueRepository_CRUD(t *testing.T) {
	sqlDB, mock, queries := setup(t)
	defer sqlDB.Close()

	repo := repository.NewPaymentDueRepository(queries)

	ctx := context.Background()
	id := uuid.New()
	invoiceID := uuid.New()
	dueDate := time.Now().AddDate(0, 0, 30)

	// --- Create ---
	mock.ExpectQuery(`INSERT INTO payment_dues`).
		WithArgs(invoiceID, int64(1000), dueDate, "PENDING", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	created, err := repo.CreatePaymentDue(ctx, db.PaymentDue{
		InvoiceID: invoiceID,
		AmountDue: "1000",
		DueDate:   dueDate,
		Status:    "PENDING",
		CreatedBy: sql.NullString{String: "tester", Valid: true},
		UpdatedBy: sql.NullString{String: "tester", Valid: true},
	})
	require.NoError(t, err)
	require.Equal(t, id, created.ID)


	// --- Get ---
	mock.ExpectQuery(`SELECT .* FROM payment_dues`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "invoice_id", "amount_due", "due_date"}).
			AddRow(id, invoiceID, int64(1000), dueDate))

	got, err := repo.GetPaymentDue(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, got.ID)

	// --- Update ---
	mock.ExpectQuery(`UPDATE payment_dues`).
		WithArgs(id, invoiceID, int64(2000), dueDate, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	_, err = repo.UpdatePaymentDue(ctx, db.PaymentDue{
		ID:        id,
		InvoiceID: invoiceID,
		AmountDue: "2000",
		DueDate:   dueDate,
	})
	require.NoError(t, err)

	// --- Delete ---
	mock.ExpectExec(`DELETE FROM payment_dues`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeletePaymentDue(ctx, id)
	require.NoError(t, err)

	// --- List ---
	mock.ExpectQuery(`SELECT .* FROM payment_dues`).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "invoice_id", "amount_due", "due_date"}).
			AddRow(id, invoiceID, int64(2000), dueDate))

	list, err := repo.ListPaymentDues(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ==================== BANK ACCOUNT TESTS ====================

func TestBankAccountRepository_CRUD(t *testing.T) {
	sqlDB, mock, queries := setup(t)
	defer sqlDB.Close()

	repo := repository.NewBankAccountRepository(queries)

	ctx := context.Background()
	id := uuid.New()

	// --- Create ---
	mock.ExpectQuery(`INSERT INTO bank_accounts`).
		WithArgs("HDFC", "1234567890", "IFSC001", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	acc, err := repo.CreateBankAccount(ctx, db.BankAccount{
		Name:          "HDFC",
		AccountNumber: "1234567890",
		IfscOrSwift:   "IFSC001",
	})
	require.NoError(t, err)
	require.Equal(t, id, acc.ID)

	// --- Get ---
	mock.ExpectQuery(`SELECT .* FROM bank_accounts`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_number", "ifsc_or_swift"}).
			AddRow(id, "HDFC", "1234567890", "IFSC001"))

	got, err := repo.GetBankAccount(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, got.ID)

	// --- Update ---
	mock.ExpectQuery(`UPDATE bank_accounts`).
		WithArgs(id, "SBI", "1234567890", "IFSC001", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	_, err = repo.UpdateBankAccount(ctx, db.BankAccount{
		ID:            id,
		Name:          "SBI",
		AccountNumber: "1234567890",
		IfscOrSwift:   "IFSC001",
	})
	require.NoError(t, err)

	// --- Delete ---
	mock.ExpectExec(`DELETE FROM bank_accounts`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteBankAccount(ctx, id)
	require.NoError(t, err)

	// --- List ---
	mock.ExpectQuery(`SELECT .* FROM bank_accounts`).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "account_number", "ifsc_or_swift"}).
			AddRow(id, "SBI", "1234567890", "IFSC001"))

	list, err := repo.ListBankAccounts(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)

	require.NoError(t, mock.ExpectationsWereMet())
}

// ==================== BANK TRANSACTION TESTS ====================

func TestBankTransactionRepository_Import(t *testing.T) {
	sqlDB, mock, queries := setup(t)
	defer sqlDB.Close()

	repo := repository.NewBankTransactionRepository(queries)

	ctx := context.Background()
	id := uuid.New()
	bankAccID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`INSERT INTO bank_transactions`).
		WithArgs(bankAccID, "2000", now, sql.NullString{}, sql.NullString{}, sql.NullBool{}, sql.NullString{}, sql.NullString{}, sql.NullString{}, sql.NullString{}).
		WillReturnRows(sqlmock.NewRows([]string{"id", "bank_account_id", "amount", "transaction_date"}).
			AddRow(id, bankAccID, "2000", now))

	tx, err := repo.ImportBankTransaction(ctx, db.BankTransaction{
		BankAccountID:   bankAccID,
		Amount:          "2000",
		TransactionDate: now,
	})
	require.NoError(t, err)
	require.Equal(t, "2000", tx.Amount)
	require.Equal(t, bankAccID, tx.BankAccountID)

	require.NoError(t, mock.ExpectationsWereMet())
}
