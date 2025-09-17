package grpc_server

import (
	"context"
	"strconv"
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


func toPbLedgerSide(s string) pb.LedgerSide {
    switch s {
    case "DEBIT","DR", "Debit":
        return pb.LedgerSide_LEDGER_SIDE_DEBIT
    case "CREDIT","CR", "Credit":
        return pb.LedgerSide_LEDGER_SIDE_CREDIT
    default:
        return pb.LedgerSide_LEDGER_SIDE_UNSPECIFIED
    }
}

//
// ---------- Ledger Entry Mapping ----------
//
func toPbLedgerEntry(e *db.LedgerEntry) *pb.LedgerEntry {
	return &pb.LedgerEntry{
		Id:        e.EntryID.String(),
		AccountId: e.AccountID.String(),
		Side:      toPbLedgerSide(e.Side),
		Amount: &money.Money{
			CurrencyCode: "USD",
			Units:        parseAmount(e.Amount),
			Nanos:        0,
		},
		TransactionDate: timestamppb.New(e.TransactionDate),
	}
}

//
// ---------- Helpers ----------
//

func toPbAccrual(a db.Accrual) *pb.Accrual {
	return &pb.Accrual{
		Id:          a.ID.String(),
		Description: a.Description.String,
		Amount:      stringToMoney(a.Amount, "USD"), // convert float64 → *money.Money
		AccrualDate: timestamppb.New(a.AccrualDate),
		AccountId:   a.AccountID,
	}
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

func mapDomainNoteTypeToProto(t string) pb.NoteType {
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

func parseAmount(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}