package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type InvoiceService struct {
	repo ports.InvoiceRepository
}

func NewInvoiceService(repo ports.InvoiceRepository) *InvoiceService {
	return &InvoiceService{repo: repo}
}

// ---------- Invoice CRUD ----------
func (s *InvoiceService) CreateInvoice(ctx context.Context, inv domain.Invoice) (domain.Invoice, error) {
	return s.repo.CreateInvoice(ctx, inv)
}

func (s *InvoiceService) GetInvoice(ctx context.Context, id uuid.UUID) (domain.Invoice, error) {
	return s.repo.GetInvoice(ctx, id)
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, inv domain.Invoice) (domain.Invoice, error) {
	return s.repo.UpdateInvoice(ctx, inv)
}

func (s *InvoiceService) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteInvoice(ctx, id)
}

func (s *InvoiceService) ListInvoices(ctx context.Context, limit, offset int32) ([]domain.Invoice, error) {
	return s.repo.ListInvoices(ctx, limit, offset)
}

func (s *InvoiceService) SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]domain.Invoice, error) {
	return s.repo.SearchInvoices(ctx, query, limit, offset)
}

// ---------- Invoice Items ----------
func (s *InvoiceService) CreateInvoiceItem(ctx context.Context, item domain.InvoiceItem) (domain.InvoiceItem, error) {
	return s.repo.CreateInvoiceItem(ctx, item)
}

func (s *InvoiceService) ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]domain.InvoiceItem, error) {
	return s.repo.ListInvoiceItems(ctx, invoiceID)
}

// ---------- Invoice Taxes & Discounts ----------
func (s *InvoiceService) AddInvoiceTax(ctx context.Context, tax domain.InvoiceTax) (domain.InvoiceTax, error) {
	return s.repo.AddInvoiceTax(ctx, tax)
}

func (s *InvoiceService) AddInvoiceDiscount(ctx context.Context, disc domain.InvoiceDiscount) (domain.InvoiceDiscount, error) {
	return s.repo.AddInvoiceDiscount(ctx, disc)
}