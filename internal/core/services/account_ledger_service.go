package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type AccountService struct {
	repo ports.AccountRepository
}

func NewAccountService(repo ports.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, a domain.Account) (domain.Account, error) {
	// Example validation rule
	if a.Code == "" || a.Name == "" {
		return domain.Account{}, ErrInvalidInput
	}
	return s.repo.Create(ctx, a)
}

func (s *AccountService) GetAccount(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	return s.repo.Get(ctx, id)
}

func (s *AccountService) ListAccounts(ctx context.Context, limit, offset int32) ([]domain.Account, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *AccountService) UpdateAccount(ctx context.Context, a domain.Account) (domain.Account, error) {
	return s.repo.Update(ctx, a)
}

func (s *AccountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
