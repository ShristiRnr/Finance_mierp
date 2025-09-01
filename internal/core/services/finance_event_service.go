package services

import (
	"context"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type FinanceEventService struct {
	repo ports.FinanceEventRepository
}

func NewFinanceEventService(r ports.FinanceEventRepository) *FinanceEventService {
	return &FinanceEventService{repo: r}
}

func (s *FinanceEventService) RecordInvoiceCreated(ctx context.Context, e domain.FinanceInvoiceCreatedEvent) (domain.FinanceInvoiceCreatedEvent, error) {
	return s.repo.InsertInvoiceCreated(ctx, e)
}

func (s *FinanceEventService) GetInvoiceCreatedEvents(ctx context.Context, orgID string, limit, offset int32) ([]domain.FinanceInvoiceCreatedEvent, error) {
	return s.repo.ListInvoiceCreated(ctx, orgID, limit, offset)
}

func (s *FinanceEventService) RecordPaymentReceived(ctx context.Context, e domain.FinancePaymentReceivedEvent) (domain.FinancePaymentReceivedEvent, error) {
	return s.repo.InsertPaymentReceived(ctx, e)
}

func (s *FinanceEventService) GetPaymentReceivedEvents(ctx context.Context, orgID string, limit, offset int32) ([]domain.FinancePaymentReceivedEvent, error) {
	return s.repo.ListPaymentReceived(ctx, orgID, limit, offset)
}

func (s *FinanceEventService) RecordInventoryCostPosted(ctx context.Context, e domain.InventoryCostPostedEvent) (domain.InventoryCostPostedEvent, error) {
	return s.repo.InsertInventoryCostPosted(ctx, e)
}

func (s *FinanceEventService) GetInventoryCostPostedEvents(ctx context.Context, orgID string, limit, offset int32) ([]domain.InventoryCostPostedEvent, error) {
	return s.repo.ListInventoryCostPosted(ctx, orgID, limit, offset)
}

func (s *FinanceEventService) RecordPayrollPosted(ctx context.Context, e domain.PayrollPostedEvent) (domain.PayrollPostedEvent, error) {
	return s.repo.InsertPayrollPosted(ctx, e)
}

func (s *FinanceEventService) GetPayrollPostedEvents(ctx context.Context, orgID string, limit, offset int32) ([]domain.PayrollPostedEvent, error) {
	return s.repo.ListPayrollPosted(ctx, orgID, limit, offset)
}

func (s *FinanceEventService) RecordVendorBillApproved(ctx context.Context, e domain.VendorBillApprovedEvent) (domain.VendorBillApprovedEvent, error) {
	return s.repo.InsertVendorBillApproved(ctx, e)
}

func (s *FinanceEventService) GetVendorBillApprovedEvents(ctx context.Context, orgID string, limit, offset int32) ([]domain.VendorBillApprovedEvent, error) {
	return s.repo.ListVendorBillApproved(ctx, orgID, limit, offset)
}
