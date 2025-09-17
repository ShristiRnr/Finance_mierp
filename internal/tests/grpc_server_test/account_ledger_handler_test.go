package grpc_server_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	grpcserver "github.com/ShristiRnr/Finance_mierp/internal/core/ports/grpc_server"
)

//
// ------------------ Mock Services ------------------
//

type MockAccountService struct{ mock.Mock }

func (m *MockAccountService) Create(ctx context.Context, a *db.Account) (*db.Account, error) {
	args := m.Called(ctx, a)
	if res := args.Get(0); res != nil {
		return res.(*db.Account), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAccountService) Get(ctx context.Context, id uuid.UUID) (*db.Account, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*db.Account), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAccountService) Update(ctx context.Context, a *db.Account) (*db.Account, error) {
	args := m.Called(ctx, a)
	if res := args.Get(0); res != nil {
		return res.(*db.Account), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAccountService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockAccountService) List(ctx context.Context, limit, offset int32) ([]*db.Account, error) {
	args := m.Called(ctx, limit, offset)
	if res := args.Get(0); res != nil {
		return res.([]*db.Account), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockJournalService struct{ mock.Mock }

func (m *MockJournalService) Create(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error) {
	args := m.Called(ctx, j)
	if res := args.Get(0); res != nil {
		return res.(*db.JournalEntry), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockJournalService) Get(ctx context.Context, id uuid.UUID) (*db.JournalEntry, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*db.JournalEntry), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockJournalService) Update(ctx context.Context, j *db.JournalEntry) (*db.JournalEntry, error) {
	args := m.Called(ctx, j)
	if res := args.Get(0); res != nil {
		return res.(*db.JournalEntry), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockJournalService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockJournalService) List(ctx context.Context, limit, offset int32) ([]*db.JournalEntry, error) {
	args := m.Called(ctx, limit, offset)
	if res := args.Get(0); res != nil {
		return res.([]*db.JournalEntry), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockLedgerService struct{ mock.Mock }

func (m *MockLedgerService) List(ctx context.Context, limit, offset int32) ([]*db.LedgerEntry, error) {
	args := m.Called(ctx, limit, offset)
	if res := args.Get(0); res != nil {
		return res.([]*db.LedgerEntry), args.Error(1)
	}
	return nil, args.Error(1)
}

// --------------------- Tests ---------------------

func TestLedgerHandler_Accounts(t *testing.T) {
	ctx := context.Background()
	mockAcc := new(MockAccountService)
	h := grpcserver.NewLedgerHandler(mockAcc, nil, nil)

	// --- CreateAccount success ---
	req := &pb.CreateAccountRequest{Account: &pb.Account{
		Code: "1001", Name: "Cash", Type: pb.AccountType_ACCOUNT_ASSET,
		Status: pb.AccountStatus_ACCOUNT_ACTIVE, AllowManualJournal: true,
	}}
	expected := &db.Account{ID: uuid.New(), Code: "1001", Name: "Cash", Type: "ASSET", Status: "ACTIVE"}
	mockAcc.On("Create", mock.Anything, mock.AnythingOfType("*db.Account")).Return(expected, nil)

	resp, err := h.CreateAccount(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "Cash", resp.Name)
	mockAcc.AssertExpectations(t)

	// --- GetAccount error (invalid id) ---
	_, err = h.GetAccount(ctx, &pb.GetAccountRequest{Id: "bad-uuid"})
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	// --- UpdateAccount error (service error) ---
	accID := uuid.New()
	mockAcc.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))
	_, err = h.UpdateAccount(ctx, &pb.UpdateAccountRequest{Account: &pb.Account{Id: accID.String()}})
	st, _ = status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())

	// --- DeleteAccount success ---
	mockAcc.On("Delete", mock.Anything, accID).Return(nil)
	_, err = h.DeleteAccount(ctx, &pb.DeleteAccountRequest{Id: accID.String()})
	assert.NoError(t, err)

	// --- ListAccounts ---
	mockAcc.On("List", mock.Anything, int32(100), int32(0)).
		Return([]*db.Account{expected}, nil)
	listResp, err := h.ListAccounts(ctx, &pb.ListAccountsRequest{})
	assert.NoError(t, err)
	assert.Len(t, listResp.Accounts, 1)
}

func TestLedgerHandler_Journals(t *testing.T) {
	ctx := context.Background()
	mockJnl := new(MockJournalService)
	h := grpcserver.NewLedgerHandler(nil, mockJnl, nil)

	// --- CreateJournalEntry ---
	expected := &db.JournalEntry{
		ID:          uuid.New(),
		JournalDate: time.Now(),
		Memo:        sql.NullString{String: "test", Valid: true},
	}
	mockJnl.On("Create", mock.Anything, mock.Anything).Return(expected, nil)
	resp, err := h.CreateJournalEntry(ctx, &pb.CreateJournalEntryRequest{
		Entry: &pb.JournalEntry{JournalDate: timestamppb.New(time.Now()), Memo: "test"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "test", resp.Memo)

	// --- GetJournalEntry invalid id ---
	_, err = h.GetJournalEntry(ctx, &pb.GetJournalEntryRequest{Id: "bad"})
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	// --- UpdateJournalEntry error ---
	id := uuid.New()
	mockJnl.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))
	_, err = h.UpdateJournalEntry(ctx, &pb.UpdateJournalEntryRequest{
		Entry: &pb.JournalEntry{Id: id.String()},
	})
	st, _ = status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())

	// --- DeleteJournalEntry success ---
	mockJnl.On("Delete", mock.Anything, id).Return(nil)
	_, err = h.DeleteJournalEntry(ctx, &pb.DeleteJournalEntryRequest{Id: id.String()})
	assert.NoError(t, err)

	// --- ListJournalEntries ---
	mockJnl.On("List", mock.Anything, int32(50), int32(0)).
		Return([]*db.JournalEntry{expected}, nil)
	listResp, err := h.ListJournalEntries(ctx, &pb.ListJournalEntriesRequest{})
	assert.NoError(t, err)
	assert.Len(t, listResp.Entries, 1)
}

func TestLedgerHandler_Ledger(t *testing.T) {
	ctx := context.Background()
	mockLgr := new(MockLedgerService)
	h := grpcserver.NewLedgerHandler(nil, nil, mockLgr)

	entry := &db.LedgerEntry{
		EntryID:         uuid.New(),
		AccountID:       uuid.New(),
		Side:            "Debit", // or "Credit"
		Amount:          "1000",
		TransactionDate: time.Now(),
	}

	mockLgr.On("List", mock.Anything, int32(100), int32(0)).
		Return([]*db.LedgerEntry{entry}, nil)

	resp, err := h.ListLedgerEntries(ctx, &pb.ListLedgerEntriesRequest{})
	assert.NoError(t, err)
	assert.Len(t, resp.Entries, 1)

	entryResp := resp.Entries[0]

	// --- Compare IDs ---
	assert.Equal(t, entry.EntryID.String(), entryResp.Id)
	assert.Equal(t, entry.AccountID.String(), entryResp.AccountId)

	// --- Convert enum to string for comparison ---
	sideStr := ""
	switch entryResp.Side {
	case pb.LedgerSide_LEDGER_SIDE_DEBIT:
		sideStr = "Debit"
	case pb.LedgerSide_LEDGER_SIDE_CREDIT:
		sideStr = "Credit"
	default:
		sideStr = "Unspecified"
	}
	assert.Equal(t, entry.Side, sideStr)

	// --- Convert money.Money to string for comparison ---
	amountStr := fmt.Sprintf("%d", entryResp.Amount.Units)
	assert.Equal(t, entry.Amount, amountStr)
}

