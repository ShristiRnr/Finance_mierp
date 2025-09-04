package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
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

func (r *InvoiceRepo) CreateInvoice(ctx context.Context, inv domain.Invoice) (domain.Invoice, error) {
	dbRow, err := r.q.CreateInvoice(ctx, db.CreateInvoiceParams{
		InvoiceNumber:        inv.InvoiceNumber,
		Type:                 inv.Type,
		InvoiceDate:          inv.InvoiceDate,
		DueDate:              ptrT2NullTime(inv.DueDate),
		DeliveryDate:         ptrT2NullTime(inv.DeliveryDate),
		OrganizationID:       inv.OrganizationID,
		PoNumber:             ptrS2NullStr(inv.PoNumber),
		EwayNumberLegacy:     ptrS2NullStr(inv.EwayNumberLegacy),
		StatusNote:           ptrS2NullStr(inv.StatusNote),
		Status:               inv.Status,
		PaymentReference:     ptrS2NullStr(inv.PaymentReference),
		ChallanNumber:        ptrS2NullStr(inv.ChallanNumber),
		ChallanDate:          ptrT2NullTime(inv.ChallanDate),
		LrNumber:             ptrS2NullStr(inv.LrNumber),
		TransporterName:      ptrS2NullStr(inv.TransporterName),
		TransporterID:        ptrS2NullStr(inv.TransporterID),
		VehicleNumber:        ptrS2NullStr(inv.VehicleNumber),
		AgainstInvoiceNumber: ptrS2NullStr(inv.AgainstInvoiceNumber),
		AgainstInvoiceDate:   ptrT2NullTime(inv.AgainstInvoiceDate),
		Subtotal:             inv.Subtotal,
		GstCgst:              ptrS2NullStr(inv.GstCgst),
		GstSgst:              ptrS2NullStr(inv.GstSgst),
		GstIgst:              ptrS2NullStr(inv.GstIgst),
		GstRate:              ptrS2NullStr(inv.GstRate),
		GrandTotal:           inv.GrandTotal,
		CreatedBy:            ptrS2NullStr(inv.CreatedBy),
		UpdatedBy:            ptrS2NullStr(inv.UpdatedBy),
		Revision:             i32ToNullInt(inv.Revision),
	})
	if err != nil {
		return domain.Invoice{}, err
	}
	return mapInvoice(dbRow), nil
}

func (r *InvoiceRepo) GetInvoice(ctx context.Context, id uuid.UUID) (domain.Invoice, error) {
	dbRow, err := r.q.GetInvoice(ctx, id)
	if err != nil {
		return domain.Invoice{}, err
	}
	return mapInvoice(dbRow), nil
}

func (r *InvoiceRepo) ListInvoices(ctx context.Context, limit, offset int32) ([]domain.Invoice, error) {
	dbRows, err := r.q.ListInvoices(ctx, db.ListInvoicesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	invoices := make([]domain.Invoice, len(dbRows))
	for i, row := range dbRows {
		invoices[i] = mapInvoice(row)
	}
	return invoices, nil
}

func (r *InvoiceRepo) UpdateInvoice(ctx context.Context, inv domain.Invoice) (domain.Invoice, error) {
	dbRow, err := r.q.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:                   inv.ID,
		InvoiceNumber:        inv.InvoiceNumber,
		Type:                 inv.Type,
		InvoiceDate:          inv.InvoiceDate,
		DueDate:              ptrT2NullTime(inv.DueDate),
		DeliveryDate:         ptrT2NullTime(inv.DeliveryDate),
		OrganizationID:       inv.OrganizationID,
		PoNumber:             ptrS2NullStr(inv.PoNumber),
		EwayNumberLegacy:     ptrS2NullStr(inv.EwayNumberLegacy),
		StatusNote:           ptrS2NullStr(inv.StatusNote),
		Status:               inv.Status,
		PaymentReference:     ptrS2NullStr(inv.PaymentReference),
		ChallanNumber:        ptrS2NullStr(inv.ChallanNumber),
		ChallanDate:          ptrT2NullTime(inv.ChallanDate),
		LrNumber:             ptrS2NullStr(inv.LrNumber),
		TransporterName:      ptrS2NullStr(inv.TransporterName),
		TransporterID:        ptrS2NullStr(inv.TransporterID),
		VehicleNumber:        ptrS2NullStr(inv.VehicleNumber),
		AgainstInvoiceNumber: ptrS2NullStr(inv.AgainstInvoiceNumber),
		AgainstInvoiceDate:   ptrT2NullTime(inv.AgainstInvoiceDate),
		Subtotal:             inv.Subtotal,
		GstCgst:              ptrS2NullStr(inv.GstCgst),
		GstSgst:              ptrS2NullStr(inv.GstSgst),
		GstIgst:              ptrS2NullStr(inv.GstIgst),
		GstRate:              ptrS2NullStr(inv.GstRate),
		GrandTotal:           inv.GrandTotal,
		UpdatedBy:            ptrS2NullStr(inv.UpdatedBy),
		Revision:             i32ToNullInt(inv.Revision),
	})
	if err != nil {
		return domain.Invoice{}, err
	}
	return mapInvoice(dbRow), nil
}

func (r *InvoiceRepo) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteInvoice(ctx, id)
}

func (r *InvoiceRepo) SearchInvoices(ctx context.Context, query string, limit, offset int32) ([]domain.Invoice, error) {
	dbRows, err := r.q.SearchInvoices(ctx, db.SearchInvoicesParams{
		Column1: ptrS2NullStr(&query),
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	invoices := make([]domain.Invoice, len(dbRows))
	for i, row := range dbRows {
		invoices[i] = mapInvoice(row)
	}
	return invoices, nil
}

// =========================================== InvoiceItem ===========================================

func (r *InvoiceRepo) CreateInvoiceItem(ctx context.Context, item domain.InvoiceItem) (domain.InvoiceItem, error) {
	dbRow, err := r.q.CreateInvoiceItem(ctx, db.CreateInvoiceItemParams{
		InvoiceID:    item.InvoiceID,
		Name:         item.Name,
		Description:  ptrS2NullStr(item.Description),
		Hsn:          ptrS2NullStr(item.Hsn),
		Quantity:     item.Quantity,
		UnitPrice:    item.UnitPrice,
		LineSubtotal: item.LineSubtotal,
		LineTotal:    item.LineTotal,
		CostCenterID: ptrS2NullStr(item.CostCenterID),
	})
	if err != nil {
		return domain.InvoiceItem{}, err
	}
	return mapInvoiceItem(dbRow), nil
}

func (r *InvoiceRepo) ListInvoiceItems(ctx context.Context, invoiceID uuid.UUID) ([]domain.InvoiceItem, error) {
	dbRows, err := r.q.ListInvoiceItems(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	items := make([]domain.InvoiceItem, len(dbRows))
	for i, row := range dbRows {
		items[i] = mapInvoiceItem(row)
	}
	return items, nil
}

// ---------- InvoiceTax ----------
func (r *InvoiceRepo) AddInvoiceTax(ctx context.Context, tax domain.InvoiceTax) (domain.InvoiceTax, error) {
	dbRow, err := r.q.AddInvoiceTax(ctx, db.AddInvoiceTaxParams{
		InvoiceID: tax.InvoiceID,
		Name:      tax.Name,
		Rate:      tax.Rate,
		Amount:    tax.Amount,
	})
	if err != nil {
		return domain.InvoiceTax{}, err
	}
	return mapInvoiceTax(dbRow), nil
}

// =========================================== InvoiceDiscount ===========================================

func (r *InvoiceRepo) AddInvoiceDiscount(ctx context.Context, disc domain.InvoiceDiscount) (domain.InvoiceDiscount, error) {
	dbRow, err := r.q.AddInvoiceDiscount(ctx, db.AddInvoiceDiscountParams{
		InvoiceID:   disc.InvoiceID,
		Description: ptrS2NullStr(disc.Description),
		Amount:      disc.Amount,
	})
	if err != nil {
		return domain.InvoiceDiscount{}, err
	}
	return mapInvoiceDiscount(dbRow), nil
}

// ============================================ Mapping functions =======================================

func mapInvoice(i db.Invoice) domain.Invoice {
	return domain.Invoice{
		ID:                   i.ID,
		InvoiceNumber:        i.InvoiceNumber,
		Type:                 i.Type,
		InvoiceDate:          i.InvoiceDate,
		DueDate:              nullTime2Ptr(i.DueDate),
		DeliveryDate:         nullTime2Ptr(i.DeliveryDate),
		OrganizationID:       i.OrganizationID,
		PoNumber:             nullStr2Ptr(i.PoNumber),
		EwayNumberLegacy:     nullStr2Ptr(i.EwayNumberLegacy),
		StatusNote:           nullStr2Ptr(i.StatusNote),
		Status:               i.Status,
		PaymentReference:     nullStr2Ptr(i.PaymentReference),
		ChallanNumber:        nullStr2Ptr(i.ChallanNumber),
		ChallanDate:          nullTime2Ptr(i.ChallanDate),
		LrNumber:             nullStr2Ptr(i.LrNumber),
		TransporterName:      nullStr2Ptr(i.TransporterName),
		TransporterID:        nullStr2Ptr(i.TransporterID),
		VehicleNumber:        nullStr2Ptr(i.VehicleNumber),
		AgainstInvoiceNumber: nullStr2Ptr(i.AgainstInvoiceNumber),
		AgainstInvoiceDate:   nullTime2Ptr(i.AgainstInvoiceDate),
		Subtotal:             i.Subtotal,
		GstCgst:              nullStr2Ptr(i.GstCgst),
		GstSgst:              nullStr2Ptr(i.GstSgst),
		GstIgst:              nullStr2Ptr(i.GstIgst),
		GstRate:              nullStr2Ptr(i.GstRate),
		GrandTotal:           i.GrandTotal,
		CreatedBy:            nullStr2Ptr(i.CreatedBy),
		UpdatedBy:            nullStr2Ptr(i.UpdatedBy),
		Revision:             i.Revision.Int32,
	}
}

func mapInvoiceItem(i db.InvoiceItem) domain.InvoiceItem {
	return domain.InvoiceItem{
		ID:          i.ID,
		InvoiceID:   i.InvoiceID,
		Name:        i.Name,
		Description: nullStr2Ptr(i.Description),
		Hsn:         nullStr2Ptr(i.Hsn),
		Quantity:    i.Quantity,
		UnitPrice:   i.UnitPrice,
		LineSubtotal:i.LineSubtotal,
		LineTotal:   i.LineTotal,
		CostCenterID:nullStr2Ptr(i.CostCenterID),
	}
}

func mapInvoiceTax(i db.InvoiceTax) domain.InvoiceTax {
	return domain.InvoiceTax{
		ID:        i.ID,
		InvoiceID: i.InvoiceID,
		Name:      i.Name,
		Rate:      i.Rate,
		Amount:    i.Amount,
	}
}

func mapInvoiceDiscount(i db.InvoiceDiscount) domain.InvoiceDiscount {
	return domain.InvoiceDiscount{
		ID:          i.ID,
		InvoiceID:   i.InvoiceID,
		Description: nullStr2Ptr(i.Description),
		Amount:      i.Amount,
	}
}

//========================================== Helper functions ==========================================

func nullStr2Ptr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullTime2Ptr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func i32ToNullInt(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: true}
}
