package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type GstRepo struct {
	q *db.Queries
}

func NewGstRepo(q *db.Queries) ports.GstRepository {
	return &GstRepo{q: q}
}

// -------------------- helpers --------------------

func f2s(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ptrF2NullStr(p *float64) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: f2s(*p), Valid: true}
}

func ptrT2NullTime(p *time.Time) sql.NullTime {
	if p == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *p, Valid: true}
}

func ptrS2NullStr(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func nullStrToPtrFloat(ns sql.NullString) *float64 {
	if !ns.Valid {
		return nil
	}
	v, err := strconv.ParseFloat(ns.String, 64)
	if err != nil {
		return nil
	}
	return &v
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	t := nt.Time
	return &t
}

func nullStrToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

// -------------------- mappers --------------------

func mapGstBreakup(r db.GstBreakup) domain.GstBreakup {
	taxable, _ := strconv.ParseFloat(r.TaxableAmount, 64)
	return domain.GstBreakup{
		ID:            r.ID,
		InvoiceID:     r.InvoiceID,
		TaxableAmount: taxable,
		CGST:          nullStrToPtrFloat(r.Cgst),
		SGST:          nullStrToPtrFloat(r.Sgst),
		IGST:          nullStrToPtrFloat(r.Igst),
		TotalGST:      nullStrToPtrFloat(r.TotalGst),
		CreatedAt:     r.CreatedAt.Time,
		CreatedBy:     r.CreatedBy.String,
		Revision:      r.Revision.Int32,
	}
}

func mapGstRegime(r db.GstRegime) domain.GstRegime {
	var rc *bool
	if r.ReverseCharge.Valid {
		v := r.ReverseCharge.Bool
		rc = &v
	}
	return domain.GstRegime{
		ID:            r.ID,
		InvoiceID:     r.InvoiceID,
		GSTIN:         r.Gstin,
		PlaceOfSupply: r.PlaceOfSupply,
		ReverseCharge: rc,
		CreatedAt:     r.CreatedAt.Time,
		CreatedBy:     r.CreatedBy.String,
		Revision:      r.Revision.Int32,
	}
}

func mapGstDocStatus(r db.GstDocStatus) domain.GstDocStatus {
	return domain.GstDocStatus{
		ID:             r.ID,
		InvoiceID:      r.InvoiceID,
		EinvoiceStatus: nullStrToPtr(r.EinvoiceStatus),
		IRN:            nullStrToPtr(r.Irn),
		AckNo:          nullStrToPtr(r.AckNo),
		AckDate:        nullTimeToPtr(r.AckDate),
		EwayStatus:     nullStrToPtr(r.EwayStatus),
		EwayBillNo:     nullStrToPtr(r.EwayBillNo),
		EwayValidUpto:  nullTimeToPtr(r.EwayValidUpto),
		LastError:      nullStrToPtr(r.LastError),
		LastSyncedAt:   nullTimeToPtr(r.LastSyncedAt),
		CreatedAt:      r.CreatedAt.Time,
		CreatedBy:      r.CreatedBy.String,
		Revision:       r.Revision.Int32,
	}
}

// -------------------- repo methods --------------------

// Breakup
func (r *GstRepo) AddGstBreakup(ctx context.Context, invoiceID uuid.UUID, taxableAmount float64, cgst, sgst, igst, totalGst *float64) (domain.GstBreakup, error) {
	dbRow, err := r.q.AddGstBreakup(ctx, db.AddGstBreakupParams{
		InvoiceID:     invoiceID,
		TaxableAmount: f2s(taxableAmount),
		Cgst:          ptrF2NullStr(cgst),
		Sgst:          ptrF2NullStr(sgst),
		Igst:          ptrF2NullStr(igst),
		TotalGst:      ptrF2NullStr(totalGst),
	})
	if err != nil {
		return domain.GstBreakup{}, err
	}
	return mapGstBreakup(dbRow), nil
}

func (r *GstRepo) GetGstBreakup(ctx context.Context, invoiceID uuid.UUID) (domain.GstBreakup, error) {
	dbRow, err := r.q.GetGstBreakup(ctx, invoiceID)
	if err != nil {
		return domain.GstBreakup{}, err
	}
	return mapGstBreakup(dbRow), nil
}

// Regime
func (r *GstRepo) AddGstRegime(ctx context.Context, invoiceID uuid.UUID, gstin, placeOfSupply string, reverseCharge *bool) (domain.GstRegime, error) {
	var rc sql.NullBool
	if reverseCharge != nil {
		rc = sql.NullBool{Bool: *reverseCharge, Valid: true}
	}
	dbRow, err := r.q.AddGstRegime(ctx, db.AddGstRegimeParams{
		InvoiceID:     invoiceID,
		Gstin:         gstin,
		PlaceOfSupply: placeOfSupply,
		ReverseCharge: rc,
	})
	if err != nil {
		return domain.GstRegime{}, err
	}
	return mapGstRegime(dbRow), nil
}

func (r *GstRepo) GetGstRegime(ctx context.Context, invoiceID uuid.UUID) (domain.GstRegime, error) {
	dbRow, err := r.q.GetGstRegime(ctx, invoiceID)
	if err != nil {
		return domain.GstRegime{}, err
	}
	return mapGstRegime(dbRow), nil
}

// Doc Status
func (r *GstRepo) AddGstDocStatus(
	ctx context.Context,
	invoiceID uuid.UUID,
	einvoiceStatus, irn, ackNo *string,
	ackDate *time.Time,
	ewayStatus, ewayBillNo *string,
	ewayValidUpto *time.Time,
	lastError *string,
	lastSyncedAt *time.Time,
) (domain.GstDocStatus, error) {
	dbRow, err := r.q.AddGstDocStatus(ctx, db.AddGstDocStatusParams{
		InvoiceID:     invoiceID,
		EinvoiceStatus: ptrS2NullStr(einvoiceStatus),
		Irn:            ptrS2NullStr(irn),
		AckNo:          ptrS2NullStr(ackNo),
		AckDate:        ptrT2NullTime(ackDate),
		EwayStatus:     ptrS2NullStr(ewayStatus),
		EwayBillNo:     ptrS2NullStr(ewayBillNo),
		EwayValidUpto:  ptrT2NullTime(ewayValidUpto),
		LastError:      ptrS2NullStr(lastError),
		LastSyncedAt:   ptrT2NullTime(lastSyncedAt),
	})
	if err != nil {
		return domain.GstDocStatus{}, err
	}
	return mapGstDocStatus(dbRow), nil
}

func (r *GstRepo) GetGstDocStatus(ctx context.Context, invoiceID uuid.UUID) (domain.GstDocStatus, error) {
	dbRow, err := r.q.GetGstDocStatus(ctx, invoiceID)
	if err != nil {
		return domain.GstDocStatus{}, err
	}
	return mapGstDocStatus(dbRow), nil
}
