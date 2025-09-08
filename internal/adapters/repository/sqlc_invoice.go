package repository

import (
	"context"

	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/google/uuid"
)

// ---------- InvoiceRepo ----------
type InvoiceRepo struct {
	q *db.Queries
}

// NewInvoiceRepo creates a new InvoiceRepo instance
func NewInvoiceRepo(q *db.Queries) ports.InvoiceRepository {
	return &InvoiceRepo{q: q}
}

// ======================================= Invoice =========================================

func (r *InvoiceRepo) CreateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error) {
	dbRow, err := r.q.CreateInvoice(ctx, db.CreateInvoiceParams{
		InvoiceNumber:        inv.InvoiceNumber,
		Type:                 inv.Type,
		InvoiceDate:          inv.InvoiceDate,
		DueDate:              inv.DueDate,
		DeliveryDate:         inv.DeliveryDate,
		OrganizationID:       inv.OrganizationID,
		PoNumber:             inv.PoNumber,
		EwayNumberLegacy:     inv.EwayNumberLegacy,
		StatusNote:           inv.StatusNote,
		Status:               inv.Status,
		PaymentReference:     inv.PaymentReference,
		ChallanNumber:        inv.ChallanNumber,
		ChallanDate:          inv.ChallanDate,
		LrNumber:             inv.LrNumber,
		TransporterName:      inv.TransporterName,
		TransporterID:        inv.TransporterID,
		VehicleNumber:        inv.VehicleNumber,
		AgainstInvoiceNumber: inv.AgainstInvoiceNumber,
		AgainstInvoiceDate:   inv.AgainstInvoiceDate,
		Subtotal:             inv.Subtotal,
		GstCgst:              inv.GstCgst,
		GstSgst:              inv.GstSgst,
		GstIgst:              inv.GstIgst,
		GstRate:              inv.GstRate,
		GrandTotal:           inv.GrandTotal,
		CreatedBy:            inv.CreatedBy,
		UpdatedBy:            inv.UpdatedBy,
		Revision:             inv.Revision,
	})
	if err != nil {
		return db.Invoice{}, err
	}
	return mapInvoice(dbRow), nil
}

func (r *InvoiceRepo) GetInvoice(ctx context.Context, id uuid.UUID) (db.Invoice, error) {
	dbRow, err := r.q.GetInvoice(ctx, id)
	if err != nil {
		return db.Invoice{}, err
	}
	return mapInvoice(dbRow), nil
}

func (r *InvoiceRepo) ListInvoices(ctx context.Context, limit, offset int32) ([]db.Invoice, error) {
	dbRows, err := r.q.ListInvoices(ctx, db.ListInvoicesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	invoices := make([]db.Invoice, len(dbRows))
	for i, row := range dbRows {
		invoices[i] = mapInvoice(row)
	}
	return invoices, nil
}

func (r *InvoiceRepo) UpdateInvoice(ctx context.Context, inv db.Invoice) (db.Invoice, error) {
	dbRow, err := r.q.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:                   inv.ID,
		InvoiceNumber:        inv.InvoiceNumber,
		Type:                 inv.Type,
		InvoiceDate:          inv.InvoiceDate,
		DueDate:              inv.DueDate,
		DeliveryDate:         inv.DeliveryDate,
		OrganizationID:       inv.OrganizationID,
		PoNumber:             inv.PoNumber,
		EwayNumberLegacy:     inv.EwayNumberLegacy,
		StatusNote:           inv.StatusNote,
		Status:               inv.Status,
		PaymentReference:     inv.PaymentReference,
		ChallanNumber:        inv.ChallanNumber,
		ChallanDate:          inv.ChallanDate,
		LrNumber:             inv.LrNumber,
		TransporterName:      inv.TransporterName,
		TransporterID:        inv.TransporterID,
		VehicleNumber:        inv.VehicleNumber,
		AgainstInvoiceNumber: inv.AgainstInvoiceNumber,
		AgainstInvoiceDate:   inv.AgainstInvoiceDate,
		Subtotal:             inv.Subtotal,
		GstCgst:              inv.GstCgst,
		GstSgst:              inv.GstSgst,
		GstIgst:              inv.GstIgst,
		GstRate:              inv.GstRate,
		GrandTotal:           inv.GrandTotal,
		UpdatedBy:            inv.UpdatedBy,
		Revision:             inv.Revision,
	})
	if err != nil {
		return db.Invoice{}, err
	}
	return mapInvoice(dbRow), nil
}

func (r *InvoiceRepo) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteInvoice(ctx, id)
}

func (r *InvoiceRepo) SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]db.Invoice, error) {
	dbRows, err := r.q.SearchInvoices(ctx, db.SearchInvoicesParams{
		Column1: ptrS2NullStr(&query),
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	invoices := make([]db.Invoice, len(dbRows))
	for i, row := range dbRows {
		invoices[i] = mapInvoice(row)
	}
	return invoices, nil
}

// =========================================== InvoiceItem ===========================================

func (r *InvoiceRepo) CreateInvoiceItem(ctx context.Context, item db.InvoiceItem) (db.InvoiceItem, error) {
	dbRow, err := r.q.CreateInvoiceItem(ctx, db.CreateInvoiceItemParams{
		InvoiceID:    item.InvoiceID,
		Name:         item.Name,
		Description:  item.Description,
		Hsn:          item.Hsn,
		Quantity:     item.Quantity,
		UnitPrice:    item.UnitPrice,
		LineSubtotal: item.LineSubtotal,
		LineTotal:    item.LineTotal,
		CostCenterID: item.CostCenterID,
	})
	if err != nil {
		return db.InvoiceItem{}, err
	}
	return mapInvoiceItem(dbRow), nil
}

func (r *InvoiceRepo) ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]db.InvoiceItem, error) {
	dbRows, err := r.q.ListInvoiceItems(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	items := make([]db.InvoiceItem, len(dbRows))
	for i, row := range dbRows {
		items[i] = mapInvoiceItem(row)
	}
	return items, nil
}

// ---------- InvoiceTax ----------
func (r *InvoiceRepo) AddInvoiceTax(ctx context.Context, tax db.InvoiceTax) (db.InvoiceTax, error) {
	dbRow, err := r.q.AddInvoiceTax(ctx, db.AddInvoiceTaxParams{
		InvoiceID: tax.InvoiceID,
		Name:      tax.Name,
		Rate:      tax.Rate,
		Amount:    tax.Amount,
	})
	if err != nil {
		return db.InvoiceTax{}, err
	}
	return mapInvoiceTax(dbRow), nil
}

// =========================================== InvoiceDiscount ===========================================

func (r *InvoiceRepo) AddInvoiceDiscount(ctx context.Context, disc db.InvoiceDiscount) (db.InvoiceDiscount, error) {
	dbRow, err := r.q.AddInvoiceDiscount(ctx, db.AddInvoiceDiscountParams{
		InvoiceID:   disc.InvoiceID,
		Description: disc.Description,
		Amount:      disc.Amount,
	})
	if err != nil {
		return db.InvoiceDiscount{}, err
	}
	return mapInvoiceDiscount(dbRow), nil
}

// ============================================ Mapping functions =======================================

func mapInvoice(i db.Invoice) db.Invoice {
	return db.Invoice{
		ID:                   i.ID,
		InvoiceNumber:        i.InvoiceNumber,
		Type:                 i.Type,
		InvoiceDate:          i.InvoiceDate,
		DueDate:              i.DueDate,
		DeliveryDate:         i.DeliveryDate,
		OrganizationID:       i.OrganizationID,
		PoNumber:             i.PoNumber,
		EwayNumberLegacy:     i.EwayNumberLegacy,
		StatusNote:           i.StatusNote,
		Status:               i.Status,
		PaymentReference:     i.PaymentReference,
		ChallanNumber:        i.ChallanNumber,
		ChallanDate:          i.ChallanDate,
		LrNumber:             i.LrNumber,
		TransporterName:      i.TransporterName,
		TransporterID:        i.TransporterID,
		VehicleNumber:        i.VehicleNumber,
		AgainstInvoiceNumber: i.AgainstInvoiceNumber,
		AgainstInvoiceDate:   i.AgainstInvoiceDate,
		Subtotal:             i.Subtotal,
		GstCgst:              i.GstCgst,
		GstSgst:              i.GstSgst,
		GstIgst:              i.GstIgst,
		GstRate:              i.GstRate,
		GrandTotal:           i.GrandTotal,
		CreatedBy:            i.CreatedBy,
		UpdatedBy:            i.UpdatedBy,
		Revision:             i.Revision,
	}
}

func mapInvoiceItem(i db.InvoiceItem) db.InvoiceItem {
	return db.InvoiceItem{
		ID:          i.ID,
		InvoiceID:   i.InvoiceID,
		Name:        i.Name,
		Description: i.Description,
		Hsn:         i.Hsn,
		Quantity:    i.Quantity,
		UnitPrice:   i.UnitPrice,
		LineSubtotal:i.LineSubtotal,
		LineTotal:   i.LineTotal,
		CostCenterID:i.CostCenterID,
	}
}

func mapInvoiceTax(i db.InvoiceTax) db.InvoiceTax {
	return db.InvoiceTax{
		ID:        i.ID,
		InvoiceID: i.InvoiceID,
		Name:      i.Name,
		Rate:      i.Rate,
		Amount:    i.Amount,
	}
}

func mapInvoiceDiscount(i db.InvoiceDiscount) db.InvoiceDiscount {
	return db.InvoiceDiscount{
		ID:          i.ID,
		InvoiceID:   i.InvoiceID,
		Description: i.Description,
		Amount:      i.Amount,
	}
}

