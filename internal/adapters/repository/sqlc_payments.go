package repository

import (
	"context"
	"time"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
)

// ================================
// BankAccount Repository
// ================================
type BankAccountRepository struct {
	queries *db.Queries
}

func NewBankAccountRepository(q *db.Queries) *BankAccountRepository {
	return &BankAccountRepository{queries: q}
}

func (r *BankAccountRepository) CreateBankAccount(ctx context.Context, ba domain.BankAccount) (domain.BankAccount, error) {
	// Example: assign ID & timestamps
	ba.ID = uuid.New()
	ba.CreatedAt = time.Now()
	ba.UpdatedAt = time.Now()
	ba.Revision = 1

	// TODO: use r.queries to insert into DB
	return ba, nil
}

func (r *BankAccountRepository) GetBankAccount(ctx context.Context, id uuid.UUID) (domain.BankAccount, error) {
	// TODO: fetch from DB
	return domain.BankAccount{ID: id}, nil
}

func (r *BankAccountRepository) UpdateBankAccount(ctx context.Context, ba domain.BankAccount) (domain.BankAccount, error) {
	ba.UpdatedAt = time.Now()
	ba.Revision += 1
	// TODO: update in DB
	return ba, nil
}

func (r *BankAccountRepository) DeleteBankAccount(ctx context.Context, id uuid.UUID) error {
	// TODO: delete from DB
	return nil
}

func (r *BankAccountRepository) ListBankAccounts(ctx context.Context, limit, offset int32) ([]domain.BankAccount, error) {
	// TODO: fetch list from DB
	return []domain.BankAccount{}, nil
}

// ================================
// PaymentDue Repository
// ================================
type PaymentDueRepository struct {
	queries *db.Queries
}

func NewPaymentDueRepository(q *db.Queries) *PaymentDueRepository {
	return &PaymentDueRepository{queries: q}
}

func (r *PaymentDueRepository) CreatePaymentDue(ctx context.Context, pd domain.PaymentDue) (domain.PaymentDue, error) {
	pd.ID = uuid.New()
	pd.CreatedAt = time.Now()
	pd.UpdatedAt = time.Now()
	pd.Revision = 1
	// TODO: insert into DB
	return pd, nil
}

func (r *PaymentDueRepository) GetPaymentDue(ctx context.Context, id uuid.UUID) (domain.PaymentDue, error) {
	// TODO: fetch from DB
	return domain.PaymentDue{ID: id}, nil
}

func (r *PaymentDueRepository) UpdatePaymentDue(ctx context.Context, pd domain.PaymentDue) (domain.PaymentDue, error) {
	pd.UpdatedAt = time.Now()
	pd.Revision += 1
	// TODO: update in DB
	return pd, nil
}

func (r *PaymentDueRepository) DeletePaymentDue(ctx context.Context, id uuid.UUID) error {
	// TODO: delete from DB
	return nil
}

func (r *PaymentDueRepository) ListPaymentDues(ctx context.Context, limit, offset int32) ([]domain.PaymentDue, error) {
	// TODO: fetch list from DB
	return []domain.PaymentDue{}, nil
}

func (r *PaymentDueRepository) MarkPaymentAsPaid(ctx context.Context, id uuid.UUID, updatedBy string) (domain.PaymentDue, error) {
	// TODO: mark as paid in DB
	pd := domain.PaymentDue{
		ID:        id,
		Status:    "Paid",
		UpdatedAt: time.Now(),
		UpdatedBy: sql.NullString{String: updatedBy, Valid: true},
	}
	return pd, nil
}

// ================================
// BankTransaction Repository
// ================================
type BankTransactionRepository struct {
	queries *db.Queries
}

func NewBankTransactionRepository(q *db.Queries) *BankTransactionRepository {
	return &BankTransactionRepository{queries: q}
}

func (r *BankTransactionRepository) ImportBankTransaction(ctx context.Context, tx domain.BankTransaction) (domain.BankTransaction, error) {
	tx.ID = uuid.New()
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	tx.Revision = 1
	// TODO: insert into DB
	return tx, nil
}

func (r *BankTransactionRepository) ListBankTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int32) ([]domain.BankTransaction, error) {
	// TODO: fetch from DB
	return []domain.BankTransaction{}, nil
}

func (r *BankTransactionRepository) ReconcileTransaction(ctx context.Context, tx domain.BankTransaction) (domain.BankTransaction, error) {
	tx.Reconciled = sql.NullBool{Bool: true, Valid: true}
	tx.UpdatedAt = time.Now()
	tx.Revision += 1
	// TODO: update in DB
	return tx, nil
}
