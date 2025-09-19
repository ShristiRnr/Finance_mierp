package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
)

type InvoiceService struct {
	repo      ports.InvoiceRepository
	eventSvc  *FinanceEventService
	publisher    ports.EventPublisher
}

// Updated constructor to include publisher
func NewInvoiceService(repo ports.InvoiceRepository,eventSvc *FinanceEventService, publisher ports.EventPublisher) *InvoiceService {
	return &InvoiceService{
		repo:      repo,
		eventSvc: eventSvc,
		publisher:    publisher,
	}
}

// ---------- Invoice CRUD ----------
func (s *InvoiceService) CreateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error) {
	// Step 1: Create invoice in DB
	i, err := s.repo.CreateInvoice(ctx, inv)
	if err != nil {
		return i, err
	}

	// Step 2: Build event from invoice row
	event := db.FinanceInvoiceCreatedEvent{
		InvoiceID:      i.ID,
		InvoiceNumber:  i.InvoiceNumber,
		InvoiceDate:    i.InvoiceDate,
		Total:          i.GrandTotal,
		OrganizationID: i.OrganizationID,
	}

	// Step 3: Record + publish event
	if _, err := s.eventSvc.RecordInvoiceCreated(ctx, event); err != nil {
		fmt.Printf("Error recording invoice.created event: %v\n", err)
	}

	return i, nil
}

func (s *InvoiceService) GetInvoice(ctx context.Context, id uuid.UUID) (db.Invoice, error) {
	return s.repo.GetInvoice(ctx, id)
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error) {
	i, err := s.repo.UpdateInvoice(ctx, inv)
	if err != nil {
		return i, err
	}

	if err := s.publisher.PublishInvoiceUpdated(ctx, &i); err != nil {
		fmt.Printf("Kafka publish error (invoice.updated): %v\n", err)
	}
	return i, nil
}

func (s *InvoiceService) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteInvoice(ctx, id); err != nil {
		return err
	}

	if err := s.publisher.PublishInvoiceDeleted(ctx, id.String()); err != nil {
		fmt.Printf("Kafka publish error (invoice.deleted): %v\n", err)
	}
	return nil
}

func (s *InvoiceService) ListInvoices(ctx context.Context, limit, offset int32) ([]db.Invoice, error) {
	return s.repo.ListInvoices(ctx, limit, offset)
}

func (s *InvoiceService) SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error) {
	return s.repo.SearchInvoices(ctx, query, limit, offset)
}

// ---------- Invoice Items ----------
func (s *InvoiceService) CreateInvoiceItem(ctx context.Context, item db.InvoiceItem) (db.InvoiceItem, error) {
	it, err := s.repo.CreateInvoiceItem(ctx, item)
	if err != nil {
		return it, err
	}

	if err := s.publisher.PublishInvoiceItemCreated(ctx, &it); err != nil {
		fmt.Printf("Kafka publish error (invoice.item.created): %v\n", err)
	}
	return it, nil
}

func (s *InvoiceService) ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]db.InvoiceItem, error) {
	return s.repo.ListInvoiceItems(ctx, invoiceID)
}

// ---------- Invoice Taxes & Discounts ----------
func (s *InvoiceService) AddInvoiceTax(ctx context.Context, tax db.InvoiceTax) (db.InvoiceTax, error) {
	t, err := s.repo.AddInvoiceTax(ctx, tax)
	if err != nil {
		return t, err
	}

	if err := s.publisher.PublishInvoiceTaxAdded(ctx, &t); err != nil {
		fmt.Printf("Kafka publish error (invoice.tax.added): %v\n", err)
	}
	return t, nil
}

func (s *InvoiceService) AddInvoiceDiscount(ctx context.Context, disc db.InvoiceDiscount) (db.InvoiceDiscount, error) {
	d, err := s.repo.AddInvoiceDiscount(ctx, disc)
	if err != nil {
		return d, err
	}

	if err := s.publisher.PublishInvoiceDiscountAdded(ctx, &d); err != nil {
		fmt.Printf("Kafka publish error (invoice.discount.added): %v\n", err)
	}
	return d, nil
}
