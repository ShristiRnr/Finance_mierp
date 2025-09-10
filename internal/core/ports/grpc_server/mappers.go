package grpc_server

import (
	"context"
	"strconv"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	money "google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc/metadata"
)

//
// ---------- Account Mapping ----------
//

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
		return nil
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
func toPbLedgerEntry(le db.LedgerEntry) *pb.LedgerEntry {
	return &pb.LedgerEntry{
		Id:              le.EntryID.String(),
		AccountId:       le.AccountID.String(),
		Side:            toPbLedgerSide(le.Side),
		Amount:          parseMoney(le.Amount),
		TransactionDate: timestamppb.New(le.TransactionDate),
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

func toPbAccrual(a db.Accrual) *pb.Accrual {
	return &pb.Accrual{
		Id:          a.ID.String(),
		Description: a.Description.String,
		Amount:      stringToMoney(a.Amount, "USD"), // convert float64 → *money.Money
		AccrualDate: timestamppb.New(a.AccrualDate),
		AccountId:   a.AccountID,
	}
}

func floatToMoney(amount float64, currency string) *money.Money {
	units := int64(amount)
	nanos := int32((amount - float64(units)) * 1e9)
	return &money.Money{
		CurrencyCode: currency,
		Units:        units,
		Nanos:        nanos,
	}
}

func stringToMoney(amountStr string, currency string) *money.Money {
	f, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil
	}
	return floatToMoney(f, currency)
}


func moneyToString(m *money.Money) string {
	if m == nil {
		return "0"
	}
	return fmt.Sprintf("%d.%09d", m.Units, m.Nanos) // crude formatting
}

func getUserFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get("user-id")
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func mapDomainNoteTypeToProto(t db.Type) pb.NoteType {
	switch t {
	case "CREDIT":
		return pb.NoteType_NOTE_TYPE_CREDIT
	case "DEBIT":
		return pb.NoteType_NOTE_TYPE_DEBIT
	default:
		return pb.NoteType_NOTE_TYPE_UNSPECIFIED
	}
}


func mapDomainAmountToProto(amount string) *money.Money {
	if amount == "" {
		return nil
	}
	dec, err := decimal.NewFromString(amount) // from shopspring/decimal
	if err != nil {
		return nil
	}
	units := dec.IntPart()
	nanos := (dec.Sub(decimal.NewFromInt(units))).Mul(decimal.New(1_000_000_000, 0)).IntPart()

	return &money.Money{
		CurrencyCode: "USD", // pick or inject dynamically
		Units:        units,
		Nanos:        int32(nanos),
	}
}


// mapDomainToProtoCreditDebitNote converts a domain CreditDebitNote to a protobuf message.
func mapDomainToProtoCreditDebitNote(note db.CreditDebitNote) *pb.CreditDebitNote {
	return &pb.CreditDebitNote{
		Id:        note.ID.String(),
		InvoiceId: note.InvoiceID.String(),
		Type:      mapDomainNoteTypeToProto(note.Type),
		Amount:    mapDomainAmountToProto(note.Amount),
		Reason:    note.Reason.String,
	}
}

func mapDomainToProtoExchangeRate(rate db.ExchangeRate) *pb.ExchangeRate {
    f64 := 0.0
    if rate.Rate != "" {
        if parsed, err := strconv.ParseFloat(rate.Rate, 64); err == nil {
            f64 = parsed
        }
    }

    return &pb.ExchangeRate{
        Id:            rate.ID.String(),
        BaseCurrency:  rate.BaseCurrency,
        QuoteCurrency: rate.QuoteCurrency,
        Rate:          f64,
        AsOf:          timestamppb.New(rate.AsOf),
    }
}

func stringPtr(s string) *string {
	return &s
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}