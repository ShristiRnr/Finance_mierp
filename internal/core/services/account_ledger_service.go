package services

import (
	"context"
	"errors"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
)

// AccountService
type AccountService struct {
	repo ports.AccountRepository
}

func NewAccountService(repo ports.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Create(ctx context.Context, a domain.Account) (domain.Account, error) {
	if a.Code == "" || a.Name == "" {
		return domain.Account{}, ErrInvalidInput
	}
	return s.repo.Create(ctx, a)
}

func (s *AccountService) Get(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	return s.repo.Get(ctx, id)
}

func (s *AccountService) Update(ctx context.Context, a domain.Account) (domain.Account, error) {
	return s.repo.Update(ctx, a)
}

func (s *AccountService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *AccountService) List(ctx context.Context, limit, offset int32) ([]domain.Account, error) {
	return s.repo.List(ctx, limit, offset)
}

// JournalService
type JournalService struct {
	repo ports.JournalRepository
}

func NewJournalService(repo ports.JournalRepository) *JournalService {
	return &JournalService{repo: repo}
}

func (s *JournalService) Create(ctx context.Context, j domain.JournalEntry) (domain.JournalEntry, error) {
	if j.JournalDate.IsZero() {
		return domain.JournalEntry{}, ErrInvalidInput
	}
	return s.repo.Create(ctx, j)
}

func (s *JournalService) Get(ctx context.Context, id uuid.UUID) (domain.JournalEntry, error) {
	return s.repo.Get(ctx, id)
}

func (s *JournalService) Update(ctx context.Context, j domain.JournalEntry) (domain.JournalEntry, error) {
	return s.repo.Update(ctx, j)
}

func (s *JournalService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *JournalService) List(ctx context.Context, limit, offset int32) ([]domain.JournalEntry, error) {
	return s.repo.List(ctx, limit, offset)
}

// LedgerService
type LedgerService struct {
	repo ports.LedgerRepository
}

func NewLedgerService(repo ports.LedgerRepository) *LedgerService {
	return &LedgerService{repo: repo}
}

func (s *LedgerService) List(ctx context.Context, limit, offset int32) ([]domain.LedgerEntry, error) {
	return s.repo.List(ctx, limit, offset)
}
