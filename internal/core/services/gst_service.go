package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type GstService struct {
	repo      ports.GstRepository
	publisher ports.EventPublisher
}

func NewGstService(repo ports.GstRepository, publisher ports.EventPublisher) *GstService {
	return &GstService{
		repo:      repo,
		publisher: publisher,
	}
}

// ---------- Breakup ----------
func (s *GstService) AddGstBreakup(ctx context.Context, invoiceID string, taxableAmount float64, cgst, sgst, igst, totalGst *float64) (db.GstBreakup, error) {
	id, err := uuid.Parse(invoiceID)
	if err != nil {
		return db.GstBreakup{}, err
	}

	breakup, err := s.repo.AddGstBreakup(ctx, id, taxableAmount, cgst, sgst, igst, totalGst)
	if err != nil {
		return breakup, err
	}

	// Publish GST breakup event
	if err := s.publisher.PublishGstBreakupAdded(ctx, &breakup); err != nil {
		fmt.Printf("Kafka publish error (gst.breakup.added): %v\n", err)
	}

	return breakup, nil
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

	regime, err := s.repo.AddGstRegime(ctx, id, gstin, placeOfSupply, reverseCharge)
	if err != nil {
		return regime, err
	}

	// Publish GST regime event
	if err := s.publisher.PublishGstRegimeAdded(ctx, &regime); err != nil {
		fmt.Printf("Kafka publish error (gst.regime.added): %v\n", err)
	}

	return regime, nil
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

	status, err := s.repo.AddGstDocStatus(ctx, id,
		einvoiceStatus, irn, ackNo, ackDate,
		ewayStatus, ewayBillNo, ewayValidUpto,
		lastError, lastSyncedAt,
	)
	if err != nil {
		return status, err
	}

	// Publish GST doc status event
	if err := s.publisher.PublishGstDocStatusAdded(ctx, &status); err != nil {
		fmt.Printf("Kafka publish error (gst.doc.status.added): %v\n", err)
	}

	return status, nil
}

func (s *GstService) GetGstDocStatus(ctx context.Context, invoiceID string) (db.GstDocStatus, error) {
	id, err := uuid.Parse(invoiceID)
	if err != nil {
		return db.GstDocStatus{}, err
	}
	return s.repo.GetGstDocStatus(ctx, id)
}
