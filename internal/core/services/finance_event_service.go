package services

import (
	"context"
	"fmt"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type FinanceEventService struct {
	repo      ports.FinanceEventRepository
	publisher ports.EventPublisher
}

// Updated constructor to include publisher
func NewFinanceEventService(r ports.FinanceEventRepository, publisher ports.EventPublisher) *FinanceEventService {
	return &FinanceEventService{
		repo:      r,
		publisher: publisher,
	}
}

// -------------------- Invoice Created --------------------
func (s *FinanceEventService) RecordInvoiceCreated(ctx context.Context, e db.FinanceInvoiceCreatedEvent) (db.FinanceInvoiceCreatedEvent, error) {
	// Persist event in DB
	event, err := s.repo.InsertInvoiceCreated(ctx, e)
	if err != nil {
		return event, err
	}

	// Publish to Kafka
	if err := s.publisher.PublishFinanceInvoiceCreated(ctx, &event); err != nil {
		fmt.Printf("Kafka publish error (invoice.created): %v\n", err)
	}

	return event, nil
}


func (s *FinanceEventService) GetInvoiceCreatedEvents(ctx context.Context, orgID string, limit, offset int32) ([]db.FinanceInvoiceCreatedEvent, error) {
	return s.repo.ListInvoiceCreated(ctx, orgID, limit, offset)
}

// -------------------- Payment Received --------------------
func (s *FinanceEventService) RecordPaymentReceived(ctx context.Context, e db.FinancePaymentReceivedEvent) (db.FinancePaymentReceivedEvent, error) {
	ev, err := s.repo.InsertPaymentReceived(ctx, e)
	if err != nil {
		return ev, err
	}

	if err := s.publisher.PublishFinancePaymentReceived(ctx, &ev); err != nil {
		fmt.Printf("Kafka publish error (payment.received): %v\n", err)
	}
	return ev, nil
}

func (s *FinanceEventService) GetPaymentReceivedEvents(ctx context.Context, orgID string, limit, offset int32) ([]db.FinancePaymentReceivedEvent, error) {
	return s.repo.ListPaymentReceived(ctx, orgID, limit, offset)
}

// -------------------- Inventory Cost Posted --------------------
func (s *FinanceEventService) RecordInventoryCostPosted(ctx context.Context, e db.InventoryCostPostedEvent) (db.InventoryCostPostedEvent, error) {
	ev, err := s.repo.InsertInventoryCostPosted(ctx, e)
	if err != nil {
		return ev, err
	}

	if err := s.publisher.PublishInventoryCostPosted(ctx, &ev); err != nil {
		fmt.Printf("Kafka publish error (inventory.cost.posted): %v\n", err)
	}
	return ev, nil
}

func (s *FinanceEventService) GetInventoryCostPostedEvents(ctx context.Context, orgID string, limit, offset int32) ([]db.InventoryCostPostedEvent, error) {
	return s.repo.ListInventoryCostPosted(ctx, orgID, limit, offset)
}

// -------------------- Payroll Posted --------------------
func (s *FinanceEventService) RecordPayrollPosted(ctx context.Context, e db.PayrollPostedEvent) (db.PayrollPostedEvent, error) {
	ev, err := s.repo.InsertPayrollPosted(ctx, e)
	if err != nil {
		return ev, err
	}

	if err := s.publisher.PublishPayrollPosted(ctx, &ev); err != nil {
		fmt.Printf("Kafka publish error (payroll.posted): %v\n", err)
	}
	return ev, nil
}

func (s *FinanceEventService) GetPayrollPostedEvents(ctx context.Context, orgID string, limit, offset int32) ([]db.PayrollPostedEvent, error) {
	return s.repo.ListPayrollPosted(ctx, orgID, limit, offset)
}

// -------------------- Vendor Bill Approved --------------------
func (s *FinanceEventService) RecordVendorBillApproved(ctx context.Context, e db.VendorBillApprovedEvent) (db.VendorBillApprovedEvent, error) {
	ev, err := s.repo.InsertVendorBillApproved(ctx, e)
	if err != nil {
		return ev, err
	}

	if err := s.publisher.PublishVendorBillApproved(ctx, &ev); err != nil {
		fmt.Printf("Kafka publish error (vendor.bill.approved): %v\n", err)
	}
	return ev, nil
}

func (s *FinanceEventService) GetVendorBillApprovedEvents(ctx context.Context, orgID string, limit, offset int32) ([]db.VendorBillApprovedEvent, error) {
	return s.repo.ListVendorBillApproved(ctx, orgID, limit, offset)
}
