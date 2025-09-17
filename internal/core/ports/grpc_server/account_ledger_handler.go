package grpc_server

import (
	"context"
	"strconv"
	"database/sql"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LedgerHandler struct {
	pb.UnimplementedLedgerServiceServer
	accountSvc ports.AccountService
	journalSvc ports.JournalService
	ledgerSvc  ports.LedgerService
}

func NewLedgerHandler(acc ports.AccountService, j ports.JournalService, l ports.LedgerService) *LedgerHandler {
	return &LedgerHandler{
		accountSvc: acc,
		journalSvc: j,
		ledgerSvc:  l,
	}
}

// ---------------- Accounts ----------------
func (h *LedgerHandler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	acc, err := h.accountSvc.Create(ctx, &db.Account{
		ID:                 uuid.New(),
		Code:               req.Account.Code,
		Name:               req.Account.Name,
		Type:               fromPbAccountType(req.Account.Type),
		Status:             fromPbAccountStatus(req.Account.Status),
		AllowManualJournal: req.Account.AllowManualJournal,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create account: %v", err)
	}
	return toPbAccount(acc), nil
}

func (h *LedgerHandler) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	acc, err := h.accountSvc.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get account: %v", err)
	}
	return toPbAccount(acc), nil
}

func (h *LedgerHandler) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.Account, error) {
	id, err := uuid.Parse(req.Account.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	acc, err := h.accountSvc.Update(ctx, &db.Account{
		ID:                 id,
		Code:               req.Account.Code,
		Name:               req.Account.Name,
		Type:               fromPbAccountType(req.Account.Type),
		Status:             fromPbAccountStatus(req.Account.Status),
		AllowManualJournal: req.Account.AllowManualJournal,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update account: %v", err)
	}
	return toPbAccount(acc), nil
}

func (h *LedgerHandler) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	if err := h.accountSvc.Delete(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete account: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (h *LedgerHandler) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	limit := int32(100)
	if req.GetPage().GetPageSize() > 0 {
		limit = req.Page.PageSize
	}
	offset := int32(0)
	if req.GetPage().GetPageToken() != "" {
		o, err := strconv.Atoi(req.Page.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page_token")
		}
		offset = int32(o)
	}

	accounts, err := h.accountSvc.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list accounts: %v", err)
	}

	pbAccounts := make([]*pb.Account, len(accounts))
	for i, a := range accounts {
		pbAccounts[i] = toPbAccount(a)
	}

	totalCount := int64(len(accounts)) // TODO: should be DB total count
	nextToken := ""
	if int64(offset)+int64(limit) < totalCount {
		nextToken = strconv.Itoa(int(offset + limit))
	}

	return &pb.ListAccountsResponse{
		Accounts: pbAccounts,
		Page: &pb.PageResponse{
			NextPageToken: nextToken,
			TotalSize:     totalCount,
		},
	}, nil
}

// ---------------- Journals ----------------
func (h *LedgerHandler) CreateJournalEntry(ctx context.Context, req *pb.CreateJournalEntryRequest) (*pb.JournalEntry, error) {
	j, err := h.journalSvc.Create(ctx, &db.JournalEntry{
		ID:          uuid.New(),
		JournalDate: req.Entry.JournalDate.AsTime(),
		Memo:        toNullString(req.Entry.Memo),
		SourceType:  toNullString(req.Entry.SourceType),
		SourceID:    toNullString(req.Entry.SourceId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create journal: %v", err)
	}
	return toPbJournalEntry(j), nil
}

func (h *LedgerHandler) GetJournalEntry(ctx context.Context, req *pb.GetJournalEntryRequest) (*pb.JournalEntry, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	j, err := h.journalSvc.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get journal: %v", err)
	}
	return toPbJournalEntry(j), nil
}

func (h *LedgerHandler) UpdateJournalEntry(ctx context.Context, req *pb.UpdateJournalEntryRequest) (*pb.JournalEntry, error) {
	id, err := uuid.Parse(req.Entry.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	j, err := h.journalSvc.Update(ctx, &db.JournalEntry{
		ID:          id,
		JournalDate: req.Entry.JournalDate.AsTime(),
		Memo:        toNullString(req.Entry.Memo),
		SourceType:  toNullString(req.Entry.SourceType),
		SourceID:    toNullString(req.Entry.SourceId),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update journal: %v", err)
	}
	return toPbJournalEntry(j), nil
}

func (h *LedgerHandler) DeleteJournalEntry(ctx context.Context, req *pb.DeleteJournalEntryRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	if err := h.journalSvc.Delete(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "delete journal: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (h *LedgerHandler) ListJournalEntries(ctx context.Context, req *pb.ListJournalEntriesRequest) (*pb.ListJournalEntriesResponse, error) {
	limit := int32(50)
	if req.GetPage().GetPageSize() > 0 {
		limit = req.Page.PageSize
	}
	offset := int32(0)
	if req.GetPage().GetPageToken() != "" {
		o, err := strconv.Atoi(req.Page.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page_token")
		}
		offset = int32(o)
	}

	entries, err := h.journalSvc.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list journal entries: %v", err)
	}

	pbEntries := make([]*pb.JournalEntry, len(entries))
	for i, e := range entries {
		pbEntries[i] = toPbJournalEntry(e)
	}

	totalCount := int64(len(entries)) // TODO: DB count
	nextToken := ""
	if int64(offset)+int64(limit) < totalCount {
		nextToken = strconv.Itoa(int(offset + limit))
	}

	return &pb.ListJournalEntriesResponse{
		Entries: pbEntries,
		Page: &pb.PageResponse{
			NextPageToken: nextToken,
			TotalSize:     totalCount,
		},
	}, nil
}

// ---------------- Ledger ----------------
func (h *LedgerHandler) ListLedgerEntries(ctx context.Context, req *pb.ListLedgerEntriesRequest) (*pb.ListLedgerEntriesResponse, error) {
	limit := int32(100)
	if req.GetPage().GetPageSize() > 0 {
		limit = req.Page.PageSize
	}
	offset := int32(0)
	if req.GetPage().GetPageToken() != "" {
		o, err := strconv.Atoi(req.Page.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page_token")
		}
		offset = int32(o)
	}

	entries, err := h.ledgerSvc.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list ledger entries: %v", err)
	}

	pbEntries := make([]*pb.LedgerEntry, len(entries))
	for i, e := range entries {
		pbEntries[i] = toPbLedgerEntry(e)
	}

	totalCount := int64(len(entries)) // TODO: DB count
	nextToken := ""
	if int64(offset)+int64(limit) < totalCount {
		nextToken = strconv.Itoa(int(offset + limit))
	}

	return &pb.ListLedgerEntriesResponse{
		Entries: pbEntries,
		Page: &pb.PageResponse{
			NextPageToken: nextToken,
			TotalSize:     totalCount,
		},
	}, nil
}

// ---------------- Helper Conversions ----------------
func toPbAccount(a *db.Account) *pb.Account {
	return &pb.Account{
		Id:                 a.ID.String(),
		Code:               a.Code,
		Name:               a.Name,
		Type:               toPbAccountType(a.Type),
		Status:             toPbAccountStatus(a.Status),
		AllowManualJournal: a.AllowManualJournal,
	}
}

func toPbJournalEntry(j *db.JournalEntry) *pb.JournalEntry {
	return &pb.JournalEntry{
		Id:          j.ID.String(),
		JournalDate: timestamppb.New(j.JournalDate),
		Memo:        j.Memo.String,
		SourceType:  j.SourceType.String,
		SourceId:    j.SourceID.String,
	}
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// ---------------- Enum Conversions ----------------

func fromPbAccountType(t pb.AccountType) string {
	switch t {
	case pb.AccountType_ACCOUNT_ASSET:
		return "ASSET"
	case pb.AccountType_ACCOUNT_LIABILITY:
		return "LIABILITY"
	case pb.AccountType_ACCOUNT_EQUITY:
		return "EQUITY"
	case pb.AccountType_ACCOUNT_REVENUE:
		return "REVENUE"
	case pb.AccountType_ACCOUNT_EXPENSE:
		return "EXPENSE"
	default:
		return "UNSPECIFIED"
	}
}

func fromPbAccountStatus(s pb.AccountStatus) string {
	switch s {
	case pb.AccountStatus_ACCOUNT_ACTIVE:
		return "ACTIVE"
	case pb.AccountStatus_ACCOUNT_INACTIVE:
		return "INACTIVE"
	case pb.AccountStatus_ACCOUNT_ARCHIVED:
		return "ARCHIVED"
	default:
		return "UNSPECIFIED"
	}
}

func toPbAccountType(s string) pb.AccountType {
	switch s {
	case "ASSET":
		return pb.AccountType_ACCOUNT_ASSET
	case "LIABILITY":
		return pb.AccountType_ACCOUNT_LIABILITY
	case "EQUITY":
		return pb.AccountType_ACCOUNT_EQUITY
	case "REVENUE":
		return pb.AccountType_ACCOUNT_REVENUE
	case "EXPENSE":
		return pb.AccountType_ACCOUNT_EXPENSE
	default:
		return pb.AccountType_ACCOUNT_TYPE_UNSPECIFIED
	}
}

func toPbAccountStatus(s string) pb.AccountStatus {
	switch s {
	case "ACTIVE":
		return pb.AccountStatus_ACCOUNT_ACTIVE
	case "INACTIVE":
		return pb.AccountStatus_ACCOUNT_INACTIVE
	case "ARCHIVED":
		return pb.AccountStatus_ACCOUNT_ARCHIVED
	default:
		return pb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
	}
}
