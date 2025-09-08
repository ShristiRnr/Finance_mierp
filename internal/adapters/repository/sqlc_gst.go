package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
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

// ========================================= mappers =========================================

func mapGstBreakup(r db.GstBreakup) db.GstBreakup {
	return db.GstBreakup{
		ID:            r.ID,
		InvoiceID:     r.InvoiceID,
		TaxableAmount: r.TaxableAmount,
		Cgst:          r.Cgst,
		Sgst:          r.Sgst,
		Igst:         r.Igst,
		TotalGst:    r.TotalGst,
		CreatedAt:    r.CreatedAt,
		CreatedBy:    r.CreatedBy,
		Revision:     r.Revision,
	}
}

func mapGstRegime(r db.GstRegime) db.GstRegime {
	return db.GstRegime{
		ID:            r.ID,
		InvoiceID:     r.InvoiceID,
		Gstin:         r.Gstin,
		PlaceOfSupply: r.PlaceOfSupply,
		ReverseCharge: r.ReverseCharge,
		CreatedAt:     r.CreatedAt,
		CreatedBy:     r.CreatedBy,
		Revision:      r.Revision,
	}
}

func mapGstDocStatus(r db.GstDocStatus) db.GstDocStatus {
	return db.GstDocStatus{
		ID:             r.ID,
		InvoiceID:      r.InvoiceID,
		EinvoiceStatus: r.EinvoiceStatus,
		Irn:            r.Irn,
		AckNo:          r.AckNo,
		AckDate:        r.AckDate,
		EwayStatus:     r.EwayStatus,
		EwayBillNo:     r.EwayBillNo,
		EwayValidUpto:  r.EwayValidUpto,
		LastError:      r.LastError,
		LastSyncedAt:   r.LastSyncedAt,
		CreatedAt:      r.CreatedAt,
		CreatedBy:      r.CreatedBy,
		Revision:       r.Revision,
	}
}

// ================================================ repo methods ================================================

// Breakup
func (r *GstRepo) AddGstBreakup(ctx context.Context, invoiceID uuid.UUID, taxableAmount float64, cgst, sgst, igst, totalGst *float64) (db.GstBreakup, error) {
	dbRow, err := r.q.AddGstBreakup(ctx, db.AddGstBreakupParams{
		InvoiceID:     invoiceID,
		TaxableAmount: f2s(taxableAmount),
		Cgst:          ptrF2NullStr(cgst),
		Sgst:          ptrF2NullStr(sgst),
		Igst:          ptrF2NullStr(igst),
		TotalGst:      ptrF2NullStr(totalGst),
	})
	if err != nil {
		return db.GstBreakup{}, err
	}
	return mapGstBreakup(dbRow), nil
}

func (r *GstRepo) GetGstBreakup(ctx context.Context, invoiceID uuid.UUID) (db.GstBreakup, error) {
	dbRow, err := r.q.GetGstBreakup(ctx, invoiceID)
	if err != nil {
		return db.GstBreakup{}, err
	}
	return mapGstBreakup(dbRow), nil
}

// Regime
func (r *GstRepo) AddGstRegime(ctx context.Context, invoiceID uuid.UUID, gstin, placeOfSupply string, reverseCharge *bool) (db.GstRegime, error) {
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
		return db.GstRegime{}, err
	}
	return mapGstRegime(dbRow), nil
}

func (r *GstRepo) GetGstRegime(ctx context.Context, invoiceID uuid.UUID) (db.GstRegime, error) {
	dbRow, err := r.q.GetGstRegime(ctx, invoiceID)
	if err != nil {
		return db.GstRegime{}, err
	}
	return mapGstRegime(dbRow), nil
}

// ================================================ Doc Status ================================================

func (r *GstRepo) AddGstDocStatus(
	ctx context.Context,
	invoiceID uuid.UUID,
	einvoiceStatus, irn, ackNo *string,
	ackDate *time.Time,
	ewayStatus, ewayBillNo *string,
	ewayValidUpto *time.Time,
	lastError *string,
	lastSyncedAt *time.Time,
) (db.GstDocStatus, error) {
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
		return db.GstDocStatus{}, err
	}
	return mapGstDocStatus(dbRow), nil
}

func (r *GstRepo) GetGstDocStatus(ctx context.Context, invoiceID uuid.UUID) (db.GstDocStatus, error) {
	dbRow, err := r.q.GetGstDocStatus(ctx, invoiceID)
	if err != nil {
		return db.GstDocStatus{}, err
	}
	return mapGstDocStatus(dbRow), nil
}
