package services

import (
	"context"
	"fmt"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

// -------------------- Bank Accounts --------------------
type BankService struct {
	repo      ports.BankAccountRepository
	publisher ports.EventPublisher
}

func NewBankService(repo ports.BankAccountRepository, publisher ports.EventPublisher) *BankService {
	return &BankService{repo: repo, publisher: publisher}
}

func (s *BankService) CreateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error) {
	acc, err := s.repo.CreateBankAccount(ctx, ba)
	if err != nil {
		return acc, err
	}

	// Publish event with correct fields
	_ = s.publisher.PublishAccountCreated(ctx, &db.Account{
		ID:   acc.ID,   // uuid.UUID
		Name: acc.Name, // field from db.BankAccount
	})

	return acc, nil
}

func (s *BankService) UpdateBankAccount(ctx context.Context, ba db.BankAccount) (db.BankAccount, error) {
	acc, err := s.repo.UpdateBankAccount(ctx, ba)
	if err != nil {
		return acc, err
	}

	_ = s.publisher.PublishAccountUpdated(ctx, &db.Account{
		ID:   acc.ID,   // uuid.UUID
		Name: acc.Name, // correct field
	})

	return acc, nil
}

func (s *BankService) DeleteBankAccount(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteBankAccount(ctx, id); err != nil {
		return err
	}

	_ = s.publisher.PublishAccountDeleted(ctx, id.String()) // string is okay for delete

	return nil
}

func (s *BankService) GetBankAccount(ctx context.Context, id uuid.UUID) (db.BankAccount, error) {
	return s.repo.GetBankAccount(ctx, id)
}

func (s *BankService) ListBankAccounts(ctx context.Context, limit, offset int32) ([]db.BankAccount, error) {
	return s.repo.ListBankAccounts(ctx, limit, offset)
}

// -------------------- Payment Due --------------------
type PaymentDueService struct {
	repo      ports.PaymentDueRepository
	publisher ports.EventPublisher
}

func NewPaymentDueService(repo ports.PaymentDueRepository, publisher ports.EventPublisher) *PaymentDueService {
	return &PaymentDueService{repo: repo, publisher: publisher}
}

func (s *PaymentDueService) CreatePaymentDue(ctx context.Context, pd db.PaymentDue) (db.PaymentDue, error) {
	pay, err := s.repo.CreatePaymentDue(ctx, pd)
	if err != nil {
		return pay, err
	}

	// Publish audit using db.AuditEvent (SQL struct)
	_ = s.publisher.PublishAuditRecorded(ctx, &db.AuditEvent{
		ID:           uuid.New(),                          // uuid.UUID
		UserID:       pay.CreatedBy.String,      // assuming CreatedBy is a string UUID
		Action:       "payment.due.created",
		Timestamp:    time.Now(),
		Details:      sql.NullString{String: fmt.Sprintf("%+v", pay), Valid: true},
		ResourceType: sql.NullString{String: "PaymentDue", Valid: true},
		ResourceID:   sql.NullString{String: pay.ID.String(), Valid: true},
	})

	return pay, nil
}


func (s *PaymentDueService) UpdatePaymentDue(ctx context.Context, pd db.PaymentDue) (db.PaymentDue, error) {
	pay, err := s.repo.UpdatePaymentDue(ctx, pd)
	if err != nil {
		return pay, err
	}

	// Publish audit using db.AuditEvent (SQL struct)
	_ = s.publisher.PublishAuditRecorded(ctx, &db.AuditEvent{
		ID:           uuid.New(),
		UserID:       pay.UpdatedBy.String, // make sure this field exists on PaymentDue
		Action:       "payment.due.updated",
		Timestamp:    time.Now(),
		Details:      sql.NullString{String: fmt.Sprintf("%+v", pay), Valid: true},
		ResourceType: sql.NullString{String: "PaymentDue", Valid: true},
		ResourceID:   sql.NullString{String: pay.ID.String(), Valid: true},
	})

	return pay, nil
}


func (s *PaymentDueService) DeletePaymentDue(ctx context.Context, id uuid.UUID, deletedBy string) error {
	if err := s.repo.DeletePaymentDue(ctx, id); err != nil {
		return err
	}

	_ = s.publisher.PublishAuditRecorded(ctx, &db.AuditEvent{
		ID:           uuid.New(),
		UserID:       deletedBy, // the user who deleted
		Action:       "payment.due.deleted",
		Timestamp:    time.Now(),
		Details:      sql.NullString{String: fmt.Sprintf("Deleted PaymentDue ID: %s", id.String()), Valid: true},
		ResourceType: sql.NullString{String: "PaymentDue", Valid: true},
		ResourceID:   sql.NullString{String: id.String(), Valid: true},
	})

	return nil
}

func (s *PaymentDueService) MarkPaymentAsPaid(ctx context.Context, id uuid.UUID, updatedBy string) (db.PaymentDue, error) {
	pay, err := s.repo.MarkPaymentAsPaid(ctx, id, updatedBy)
	if err != nil {
		return pay, err
	}

	_ = s.publisher.PublishAuditRecorded(ctx, &db.AuditEvent{
		ID:           uuid.New(),
		UserID:       updatedBy,
		Action:       "payment.due.paid",
		Timestamp:    time.Now(),
		Details:      sql.NullString{String: fmt.Sprintf("%+v", pay), Valid: true},
		ResourceType: sql.NullString{String: "PaymentDue", Valid: true},
		ResourceID:   sql.NullString{String: pay.ID.String(), Valid: true},
	})

	return pay, nil
}


func (s *PaymentDueService) GetPaymentDue(ctx context.Context, id uuid.UUID) (db.PaymentDue, error) {
	return s.repo.GetPaymentDue(ctx, id)
}

func (s *PaymentDueService) ListPaymentDues(ctx context.Context, limit, offset int32) ([]db.PaymentDue, error) {
	return s.repo.ListPaymentDues(ctx, limit, offset)
}

// -------------------- Bank Transactions --------------------
type BankTransactionService struct {
	repo      ports.BankTransactionRepository
	publisher ports.EventPublisher
}

func NewBankTransactionService(repo ports.BankTransactionRepository, publisher ports.EventPublisher) *BankTransactionService {
	return &BankTransactionService{repo: repo, publisher: publisher}
}

func (s *BankTransactionService) ImportBankTransaction(ctx context.Context, tx db.BankTransaction, userID string) (db.BankTransaction, error) {
	trx, err := s.repo.ImportBankTransaction(ctx, tx)
	if err != nil {
		return trx, err
	}

	_ = s.publisher.PublishAuditRecorded(ctx, &db.AuditEvent{
		ID:           uuid.New(),
		UserID:       userID,
		Action:       "bank.transaction.imported",
		Timestamp:    time.Now(),
		Details:      sql.NullString{String: fmt.Sprintf("%+v", trx), Valid: true},
		ResourceType: sql.NullString{String: "BankTransaction", Valid: true},
		ResourceID:   sql.NullString{String: trx.ID.String(), Valid: true},
	})

	return trx, nil
}

func (s *BankTransactionService) ReconcileTransaction(ctx context.Context, tx db.BankTransaction, userID string) (db.BankTransaction, error) {
	trx, err := s.repo.ReconcileTransaction(ctx, tx)
	if err != nil {
		return trx, err
	}

	_ = s.publisher.PublishAuditRecorded(ctx, &db.AuditEvent{
		ID:           uuid.New(),
		UserID:       userID,
		Action:       "bank.transaction.reconciled",
		Timestamp:    time.Now(),
		Details:      sql.NullString{String: fmt.Sprintf("%+v", trx), Valid: true},
		ResourceType: sql.NullString{String: "BankTransaction", Valid: true},
		ResourceID:   sql.NullString{String: trx.ID.String(), Valid: true},
	})

	return trx, nil
}


func (s *BankTransactionService) ListBankTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int32) ([]db.BankTransaction, error) {
	return s.repo.ListBankTransactions(ctx, bankAccountID, limit, offset)
}
