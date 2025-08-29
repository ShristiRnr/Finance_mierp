package grpc_server

import (
	"strconv"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	money "google.golang.org/genproto/googleapis/type/money"
)

//
// ---------- Account Mapping ----------
//
func toDomainAccountType(t pb.AccountType) string {
	return pb.AccountType_name[int32(t)] // direct proto string ("ACCOUNT_ASSET")
}

func toPbAccountType(s string) pb.AccountType {
	if val, ok := pb.AccountType_value[s]; ok {
		return pb.AccountType(val)
	}
	return pb.AccountType_ACCOUNT_TYPE_UNSPECIFIED
}

func toDomainStatus(s pb.AccountStatus) string {
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

func toPbStatus(s string) pb.AccountStatus {
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


func toPbLedgerSide(s string) pb.LedgerSide {
    switch s {
    case "DEBIT":
        return pb.LedgerSide_LEDGER_SIDE_DEBIT
    case "CREDIT":
        return pb.LedgerSide_LEDGER_SIDE_CREDIT
    default:
        return pb.LedgerSide_LEDGER_SIDE_UNSPECIFIED
    }
}

func parseMoney(amount string) *money.Money {
	if amount == "" {
		return nil
	}
	val, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return nil // or handle error properly
	}

	units := int64(val)
	nanos := int32((val - float64(units)) * 1e9)

	return &money.Money{
		CurrencyCode: "USD",
		Units:        units,
		Nanos:        nanos,
	}
}

//
// ---------- Ledger Entry Mapping ----------
//
func toPbLedgerEntry(le domain.LedgerEntry) *pb.LedgerEntry {
	return &pb.LedgerEntry{
		Id:              le.EntryID.String(),
		AccountId:       le.AccountID.String(),
		Side:            toPbLedgerSide(le.Side),
		Amount:          parseMoney(le.Amount),
		TransactionDate: timestamppb.New(le.PostedAt),

		// Not available in domain yet
		Description:   "",
		CostCenterId:  "",
		ReferenceType: "",
		ReferenceId:   "",
		Audit:         nil,
		ExternalRefs:  nil,
	}
}

//
// ---------- Helpers ----------
//
func strOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
