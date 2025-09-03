package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
)


// =========================================== BankAccount Repository ========================================

type BankAccountRepository struct {
	queries *db.Queries
}

func NewBankAccountRepository(q *db.Queries) *BankAccountRepository {
	return &BankAccountRepository{queries: q}
}

func (r *BankAccountRepository) CreateBankAccount(ctx context.Context, ba domain.BankAccount) (domain.BankAccount, error) {
	params := db.CreateBankAccountParams{
		Name:            ba.Name,
		AccountNumber:   ba.AccountNumber,
		IfscOrSwift:     ba.IfscOrSwift,
		LedgerAccountID: ba.LedgerAccountID,
		CreatedBy:       ba.CreatedBy,
		UpdatedBy:       ba.UpdatedBy,
	}
	b, err := r.queries.CreateBankAccount(ctx, params)
	if err != nil {
		return domain.BankAccount{}, err
	}
	return dbBankAccountToDomain(b), nil
}

func (r *BankAccountRepository) GetBankAccount(ctx context.Context, id uuid.UUID) (domain.BankAccount, error) {
	b, err := r.queries.GetBankAccount(ctx, id)
	if err != nil {
		return domain.BankAccount{}, err
	}
	return dbBankAccountToDomain(b), nil
}

func (r *BankAccountRepository) UpdateBankAccount(ctx context.Context, ba domain.BankAccount) (domain.BankAccount, error) {
	params := db.UpdateBankAccountParams{
		ID:              ba.ID,
		Name:            ba.Name,
		AccountNumber:   ba.AccountNumber,
		IfscOrSwift:     ba.IfscOrSwift,
		LedgerAccountID: ba.LedgerAccountID,
		UpdatedBy:       ba.UpdatedBy,
	}
	b, err := r.queries.UpdateBankAccount(ctx, params)
	if err != nil {
		return domain.BankAccount{}, err
	}
	return dbBankAccountToDomain(b), nil
}

func (r *BankAccountRepository) DeleteBankAccount(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteBankAccount(ctx, id)
}

func (r *BankAccountRepository) ListBankAccounts(ctx context.Context, limit, offset int32) ([]domain.BankAccount, error) {
	params := db.ListBankAccountsParams{Limit: limit, Offset: offset}
	list, err := r.queries.ListBankAccounts(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]domain.BankAccount, len(list))
	for i, b := range list {
		result[i] = dbBankAccountToDomain(b)
	}
	return result, nil
}

//============================================ PaymentDue Repository ========================================

type PaymentDueRepository struct {
	queries *db.Queries
}

func NewPaymentDueRepository(q *db.Queries) *PaymentDueRepository {
	return &PaymentDueRepository{queries: q}
}

func (r *PaymentDueRepository) CreatePaymentDue(ctx context.Context, pd domain.PaymentDue) (domain.PaymentDue, error) {
	params := db.CreatePaymentDueParams{
		InvoiceID: pd.InvoiceID,
		AmountDue: pd.AmountDue,
		DueDate:   pd.DueDate,
		Status:    pd.Status,
		CreatedBy: pd.CreatedBy,
		UpdatedBy: pd.UpdatedBy,
	}
	b, err := r.queries.CreatePaymentDue(ctx, params)
	if err != nil {
		return domain.PaymentDue{}, err
	}
	return dbPaymentDueToDomain(b), nil
}

func (r *PaymentDueRepository) GetPaymentDue(ctx context.Context, id uuid.UUID) (domain.PaymentDue, error) {
	b, err := r.queries.GetPaymentDue(ctx, id)
	if err != nil {
		return domain.PaymentDue{}, err
	}
	return dbPaymentDueToDomain(b), nil
}

func (r *PaymentDueRepository) UpdatePaymentDue(ctx context.Context, pd domain.PaymentDue) (domain.PaymentDue, error) {
	params := db.UpdatePaymentDueParams{
		ID:        pd.ID,
		InvoiceID: pd.InvoiceID,
		AmountDue: pd.AmountDue,
		DueDate:   pd.DueDate,
		Status:    pd.Status,
		UpdatedBy: pd.UpdatedBy,
	}
	b, err := r.queries.UpdatePaymentDue(ctx, params)
	if err != nil {
		return domain.PaymentDue{}, err
	}
	return dbPaymentDueToDomain(b), nil
}

func (r *PaymentDueRepository) DeletePaymentDue(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeletePaymentDue(ctx, id)
}

func (r *PaymentDueRepository) ListPaymentDues(ctx context.Context, limit, offset int32) ([]domain.PaymentDue, error) {
	params := db.ListPaymentDuesParams{Limit: limit, Offset: offset}
	list, err := r.queries.ListPaymentDues(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]domain.PaymentDue, len(list))
	for i, b := range list {
		result[i] = dbPaymentDueToDomain(b)
	}
	return result, nil
}

func (r *PaymentDueRepository) MarkPaymentAsPaid(ctx context.Context, id uuid.UUID, updatedBy string) (domain.PaymentDue, error) {
	params := db.MarkPaymentAsPaidParams{
		ID:        id,
		UpdatedBy: sql.NullString{String: updatedBy, Valid: true},
	}
	b, err := r.queries.MarkPaymentAsPaid(ctx, params)
	if err != nil {
		return domain.PaymentDue{}, err
	}
	return dbPaymentDueToDomain(b), nil
}


// ======================================== BankTransaction Repository ========================================

type BankTransactionRepository struct {
	queries *db.Queries
}

func NewBankTransactionRepository(q *db.Queries) *BankTransactionRepository {
	return &BankTransactionRepository{queries: q}
}

func (r *BankTransactionRepository) ImportBankTransaction(ctx context.Context, tx domain.BankTransaction) (domain.BankTransaction, error) {
	params := db.ImportBankTransactionParams{
		BankAccountID:        tx.BankAccountID,
		Amount:               tx.Amount,
		TransactionDate:      tx.TransactionDate,
		Description:          tx.Description,
		Reference:            tx.Reference,
		Reconciled:           tx.Reconciled,
		MatchedReferenceType: tx.MatchedReferenceType,
		MatchedReferenceID:   tx.MatchedReferenceID,
		CreatedBy:            tx.CreatedBy,
		UpdatedBy:            tx.UpdatedBy,
	}
	b, err := r.queries.ImportBankTransaction(ctx, params)
	if err != nil {
		return domain.BankTransaction{}, err
	}
	return dbBankTransactionToDomain(b), nil
}

func (r *BankTransactionRepository) ListBankTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int32) ([]domain.BankTransaction, error) {
	params := db.ListBankTransactionsParams{
		BankAccountID: bankAccountID,
		Limit:         limit,
		Offset:        offset,
	}
	list, err := r.queries.ListBankTransactions(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]domain.BankTransaction, len(list))
	for i, b := range list {
		result[i] = dbBankTransactionToDomain(b)
	}
	return result, nil
}

func (r *BankTransactionRepository) ReconcileTransaction(ctx context.Context, tx domain.BankTransaction) (domain.BankTransaction, error) {
	params := db.ReconcileTransactionParams{
		ID:                   tx.ID,
		MatchedReferenceType: tx.MatchedReferenceType,
		MatchedReferenceID:   tx.MatchedReferenceID,
		UpdatedBy:            tx.UpdatedBy,
	}
	b, err := r.queries.ReconcileTransaction(ctx, params)
	if err != nil {
		return domain.BankTransaction{}, err
	}
	return dbBankTransactionToDomain(b), nil
}


// =================================================== Helpers ==========================================

func dbBankAccountToDomain(b db.BankAccount) domain.BankAccount {
	return domain.BankAccount{
		ID:              b.ID,
		Name:            b.Name,
		AccountNumber:   b.AccountNumber,
		IfscOrSwift:     b.IfscOrSwift,
		LedgerAccountID: b.LedgerAccountID,
		CreatedAt:       b.CreatedAt.Time,
		CreatedBy:       b.CreatedBy,
		UpdatedAt:       b.UpdatedAt.Time,
		UpdatedBy:       b.UpdatedBy,
		Revision:        b.Revision.Int32,
	}
}

func dbPaymentDueToDomain(b db.PaymentDue) domain.PaymentDue {
	return domain.PaymentDue{
		ID:        b.ID,
		InvoiceID: b.InvoiceID,
		AmountDue: b.AmountDue,
		DueDate:   b.DueDate,
		Status:    b.Status,
		CreatedAt:       b.CreatedAt.Time,
		CreatedBy:       b.CreatedBy,
		UpdatedAt:       b.UpdatedAt.Time,
		UpdatedBy:       b.UpdatedBy,
		Revision:        b.Revision.Int32,
	}
}

func dbBankTransactionToDomain(b db.BankTransaction) domain.BankTransaction {
	return domain.BankTransaction{
		ID:                  b.ID,
		BankAccountID:       b.BankAccountID,
		Amount:              b.Amount,
		TransactionDate:     b.TransactionDate,
		Description:         b.Description,
		Reference:           b.Reference,
		Reconciled:          b.Reconciled,
		MatchedReferenceType: b.MatchedReferenceType,
		MatchedReferenceID:   b.MatchedReferenceID,
		CreatedAt:           b.CreatedAt.Time,
		CreatedBy:           b.CreatedBy,
		UpdatedAt:           b.UpdatedAt.Time,
		UpdatedBy:           b.UpdatedBy,
		Revision:            b.Revision.Int32,
	}
}
