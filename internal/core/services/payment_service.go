package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type BankService struct {
	repo ports.BankAccountRepository
}

func NewBankService(repo ports.BankAccountRepository) *BankService {
	return &BankService{repo: repo}
}

func (s *BankService) CreateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error) {
	return s.repo.CreateBankAccount(ctx, ba)
}

func (s *BankService) GetBankAccount(ctx context.Context, id uuid.UUID) (db.BankAccount, error) {
	return s.repo.GetBankAccount(ctx, id)
}

func (s *BankService) UpdateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error) {
	return s.repo.UpdateBankAccount(ctx, ba)
}

func (s *BankService) DeleteBankAccount(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteBankAccount(ctx, id)
}

func (s *BankService) ListBankAccounts(ctx context.Context, limit, offset int32) ([]db.BankAccount, error) {
	return s.repo.ListBankAccounts(ctx, limit, offset)
}

type PaymentDueService struct {
	repo ports.PaymentDueRepository
}

func NewPaymentDueService(repo ports.PaymentDueRepository) *PaymentDueService {
	return &PaymentDueService{repo: repo}
}

func (s *PaymentDueService) CreatePaymentDue(ctx context.Context, pd db.PaymentDue) (db.PaymentDue, error) {
	return s.repo.CreatePaymentDue(ctx, pd)
}

func (s *PaymentDueService) GetPaymentDue(ctx context.Context, id uuid.UUID) (db.PaymentDue, error) {
	return s.repo.GetPaymentDue(ctx, id)
}

func (s *PaymentDueService) UpdatePaymentDue(ctx context.Context, pd db.PaymentDue) (db.PaymentDue, error) {
	return s.repo.UpdatePaymentDue(ctx, pd)
}

func (s *PaymentDueService) DeletePaymentDue(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePaymentDue(ctx, id)
}

func (s *PaymentDueService) ListPaymentDues(ctx context.Context, limit, offset int32) ([]db.PaymentDue, error) {
	return s.repo.ListPaymentDues(ctx, limit, offset)
}

func (s *PaymentDueService) MarkPaymentAsPaid(ctx context.Context, id uuid.UUID, updatedBy string) (db.PaymentDue, error) {
	return s.repo.MarkPaymentAsPaid(ctx, id, updatedBy)
}

type BankTransactionService struct {
	repo ports.BankTransactionRepository
}

func NewBankTransactionService(repo ports.BankTransactionRepository) *BankTransactionService {
	return &BankTransactionService{repo: repo}
}

func (s *BankTransactionService) ImportBankTransaction(ctx context.Context, tx db.BankTransaction) (db.BankTransaction, error) {
	return s.repo.ImportBankTransaction(ctx, tx)
}

func (s *BankTransactionService) ListBankTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int32) ([]db.BankTransaction, error) {
	return s.repo.ListBankTransactions(ctx, bankAccountID, limit, offset)
}

func (s *BankTransactionService) ReconcileTransaction(ctx context.Context, tx db.BankTransaction) (db.BankTransaction, error) {
	return s.repo.ReconcileTransaction(ctx, tx)
}