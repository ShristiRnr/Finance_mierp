package services

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
    "github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type GstService struct {
    repo ports.GstRepository
}

func NewGstService(repo ports.GstRepository) *GstService {
    return &GstService{repo: repo}
}

// ---------- Breakup ----------
func (s *GstService) AddGstBreakup(ctx context.Context, invoiceID string, taxableAmount float64, cgst, sgst, igst, totalGst *float64) (db.GstBreakup, error) {
    id, err := uuid.Parse(invoiceID)
    if err != nil {
        return db.GstBreakup{}, err
    }
    return s.repo.AddGstBreakup(ctx, id, taxableAmount, cgst, sgst, igst, totalGst)
}

func (s *GstService) GetGstBreakup(ctx context.Context, invoiceID string) (db.GstBreakup, error) {
    id, err := uuid.Parse(invoiceID)
    if err != nil {
        return db.GstBreakup{}, err
    }
    return s.repo.GetGstBreakup(ctx, id)
}

// ---------- Regime ----------
func (s *GstService) AddGstRegime(ctx context.Context, invoiceID, gstin, placeOfSupply string, reverseCharge *bool) (db.GstRegime, error) {
    id, err := uuid.Parse(invoiceID)
    if err != nil {
        return db.GstRegime{}, err
    }
    return s.repo.AddGstRegime(ctx, id, gstin, placeOfSupply, reverseCharge)
}

func (s *GstService) GetGstRegime(ctx context.Context, invoiceID string) (db.GstRegime, error) {
    id, err := uuid.Parse(invoiceID)
    if err != nil {
        return db.GstRegime{}, err
    }
    return s.repo.GetGstRegime(ctx, id)
}

// ---------- Doc Status ----------
func (s *GstService) AddGstDocStatus(
    ctx context.Context,
    invoiceID string,
    einvoiceStatus, irn, ackNo *string,
    ackDate *time.Time,
    ewayStatus, ewayBillNo *string,
    ewayValidUpto *time.Time,
    lastError *string,
    lastSyncedAt *time.Time,
) (db.GstDocStatus, error) {
    id, err := uuid.Parse(invoiceID)
    if err != nil {
        return db.GstDocStatus{}, err
    }
    return s.repo.AddGstDocStatus(ctx, id,
        einvoiceStatus, irn, ackNo, ackDate,
        ewayStatus, ewayBillNo, ewayValidUpto,
        lastError, lastSyncedAt,
    )
}

func (s *GstService) GetGstDocStatus(ctx context.Context, invoiceID string) (db.GstDocStatus, error) {
    id, err := uuid.Parse(invoiceID)
    if err != nil {
        return db.GstDocStatus{}, err
    }
    return s.repo.GetGstDocStatus(ctx, id)
}
