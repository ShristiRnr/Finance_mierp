
package services

import (
	"context"
	"errors"
	"log"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
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
	publisher ports.EventPublisher
}

func NewAccountService(repo ports.AccountRepository, pub ports.EventPublisher) *AccountService {
	return &AccountService{repo: repo, publisher: pub}
}

func (s *AccountService) Create(ctx context.Context, a *db.Account) (*db.Account, error) {
	createdVal, err := s.repo.Create(ctx, a) // repo expects value type
	if err != nil {
		return nil, err
	}
	if s.publisher != nil {
		go func(acct *db.Account) {
			if err := s.publisher.PublishAccountCreated(context.Background(), acct); err != nil {
				log.Printf("failed to publish account created: %v", err)
			}
		}(createdVal)
	}

	return createdVal, nil
}

func (s *AccountService) Get(ctx context.Context, id uuid.UUID) (*db.Account, error) {
	acc, err := s.repo.Get(ctx, id) // returns db.Account
	if err != nil {
		return nil, err
	}
	return acc, nil
}


// AccountService
func (s *AccountService) Update(ctx context.Context, a *db.Account) (*db.Account, error) {
	updatedVal, err := s.repo.Update(ctx, a)
	if err != nil {
		return nil, err
	}


	if s.publisher != nil {
		go func(acct *db.Account) {
			if err := s.publisher.PublishAccountUpdated(context.Background(), acct); err != nil {
				log.Printf("failed to publish account updated: %v", err)
			}
		}(updatedVal)
	}

	return updatedVal, nil
}

func (s *AccountService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.publisher != nil {
		go func() {
			if err := s.publisher.PublishAccountDeleted(context.Background(), id.String()); err != nil {
				log.Printf("failed to publish account deleted: %v", err)
			}
		}()
	}

	return nil
}


func (s *AccountService) List(ctx context.Context, limit, offset int32) ([]*db.Account, error) {
	accs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// convert slice of values to slice of pointers
	result := make([]*db.Account, len(accs))
	copy(result, accs)


	return result, nil
}

// JournalService
type JournalService struct {
	repo ports.JournalRepository
	publisher ports.EventPublisher
}

func NewJournalService(repo ports.JournalRepository, pub ports.EventPublisher) *JournalService {
	return &JournalService{repo: repo, publisher: pub}
}

func (s *JournalService) Create(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error) {
	createdVal, err := s.repo.Create(ctx, j)
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		go func(jrnl *db.JournalEntry) {
			if err := s.publisher.PublishJournalCreated(context.Background(), jrnl); err != nil {
				log.Printf("failed to publish journal created: %v", err)
			}
		}(createdVal)
	}

	return createdVal, nil
}

func (s *JournalService) Get(ctx context.Context, id uuid.UUID) (*db.JournalEntry, error) {
	return s.repo.Get(ctx, id)
}

func (s *JournalService) Update(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error) {
    // Convert pointer to value for repo
    updatedVal, err := s.repo.Update(ctx, j)
    if err != nil {
        return nil, err
    }

    // Pointer for return and publishing
    if s.publisher != nil {
        go func(jrnl *db.JournalEntry) {
            if err := s.publisher.PublishJournalUpdated(context.Background(), jrnl); err != nil {
                log.Printf("failed to publish journal updated: %v", err)
            }
        }(updatedVal)
    }

    return updatedVal, nil
}




func (s *JournalService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.publisher != nil {
		go func() {
			if err := s.publisher.PublishJournalDeleted(context.Background(), id.String()); err != nil {
				log.Printf("failed to publish journal deleted: %v", err)
			}
		}()
	}

	return nil
}

func (s *JournalService) List(ctx context.Context, limit, offset int32) ([]*db.JournalEntry, error) {
	return s.repo.List(ctx, limit, offset)
}

// LedgerService
type LedgerService struct {
	repo ports.LedgerRepository
}

func NewLedgerService(repo ports.LedgerRepository) *LedgerService {
	return &LedgerService{repo: repo}
}

func (s *LedgerService) List(ctx context.Context, limit, offset int32) ([]*db.LedgerEntry, error) {
	return s.repo.List(ctx, limit, offset)
}