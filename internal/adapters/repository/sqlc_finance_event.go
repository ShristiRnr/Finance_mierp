package repository

import (
	"context"

	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

type FinanceEventRepo struct {
	q *db.Queries
}

func NewFinanceEventRepo(q *db.Queries) ports.FinanceEventRepository {
	return &FinanceEventRepo{q: q}
}

// ======================================================= INVOICE CREATED ===========================================================

func (r *FinanceEventRepo) InsertInvoiceCreated(ctx context.Context, e db.FinanceInvoiceCreatedEvent) (db.FinanceInvoiceCreatedEvent, error) {
	params := db.InsertInvoiceCreatedEventParams{
		InvoiceID:      e.InvoiceID,
		InvoiceNumber:  e.InvoiceNumber,
		InvoiceDate:    e.InvoiceDate,
		Total:          e.Total,
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertInvoiceCreatedEvent(ctx, params)
	if err != nil {
		return db.FinanceInvoiceCreatedEvent{}, err
	}
	return mapInvoiceCreatedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListInvoiceCreated(ctx context.Context, orgID string, limit, offset int32) ([]db.FinanceInvoiceCreatedEvent, error) {
	rows, err := r.q.ListInvoiceCreatedEvents(ctx, db.ListInvoiceCreatedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]db.FinanceInvoiceCreatedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapInvoiceCreatedRowToDomain(row)
	}
	return events, nil
}

// ======================================================= PAYMENT RECEIVED ==============================================================

func (r *FinanceEventRepo) InsertPaymentReceived(ctx context.Context, e db.FinancePaymentReceivedEvent) (db.FinancePaymentReceivedEvent, error) {
	params := db.InsertPaymentReceivedEventParams{
		PaymentDueID:   e.PaymentDueID,
		InvoiceID:      e.InvoiceID,
		AmountPaid:     e.AmountPaid,
		PaidAt:         e.PaidAt,
		Reference:      e.Reference,
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertPaymentReceivedEvent(ctx, params)
	if err != nil {
		return db.FinancePaymentReceivedEvent{}, err
	}
	return mapPaymentReceivedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListPaymentReceived(ctx context.Context, orgID string, limit, offset int32) ([]db.FinancePaymentReceivedEvent, error) {
	rows, err := r.q.ListPaymentReceivedEvents(ctx, db.ListPaymentReceivedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]db.FinancePaymentReceivedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapPaymentReceivedRowToDomain(row)
	}
	return events, nil
}

// ========================================================== INVENTORY COST POSTED ===============================================================

func (r *FinanceEventRepo) InsertInventoryCostPosted(ctx context.Context, e db.InventoryCostPostedEvent) (db.InventoryCostPostedEvent, error) {
	params := db.InsertInventoryCostPostedEventParams{
		ReferenceType:  e.ReferenceType,
		ReferenceID:    e.ReferenceID,
		Amount:         e.Amount,
		CostCenterID:   e.CostCenterID,
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertInventoryCostPostedEvent(ctx, params)
	if err != nil {
		return db.InventoryCostPostedEvent{}, err
	}
	return mapInventoryCostPostedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListInventoryCostPosted(ctx context.Context, orgID string, limit, offset int32) ([]db.InventoryCostPostedEvent, error) {
	rows, err := r.q.ListInventoryCostPostedEvents(ctx, db.ListInventoryCostPostedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]db.InventoryCostPostedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapInventoryCostPostedRowToDomain(row)
	}
	return events, nil
}

// =========================================================== PAYROLL POSTED ==============================================================

func (r *FinanceEventRepo) InsertPayrollPosted(ctx context.Context, e db.PayrollPostedEvent) (db.PayrollPostedEvent, error) {
	params := db.InsertPayrollPostedEventParams{
		PayrollRunID:   e.PayrollRunID,
		TotalGross:     e.TotalGross,
		TotalNet:       e.TotalNet,
		RunDate:        e.RunDate,
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertPayrollPostedEvent(ctx, params)
	if err != nil {
		return db.PayrollPostedEvent{}, err
	}
	return mapPayrollPostedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListPayrollPosted(ctx context.Context, orgID string, limit, offset int32) ([]db.PayrollPostedEvent, error) {
	rows, err := r.q.ListPayrollPostedEvents(ctx, db.ListPayrollPostedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]db.PayrollPostedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapPayrollPostedRowToDomain(row)
	}
	return events, nil
}

// =================================================== VENDOR BILL APPROVED ====================================================

func (r *FinanceEventRepo) InsertVendorBillApproved(ctx context.Context, e db.VendorBillApprovedEvent) (db.VendorBillApprovedEvent, error) {
	params := db.InsertVendorBillApprovedEventParams{
		VendorBillID:   e.VendorBillID,
		Amount:         e.Amount,
		ApprovedAt:     e.ApprovedAt,
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertVendorBillApprovedEvent(ctx, params)
	if err != nil {
		return db.VendorBillApprovedEvent{}, err
	}
	return mapVendorBillApprovedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListVendorBillApproved(ctx context.Context, orgID string, limit, offset int32) ([]db.VendorBillApprovedEvent, error) {
	rows, err := r.q.ListVendorBillApprovedEvents(ctx, db.ListVendorBillApprovedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]db.VendorBillApprovedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapVendorBillApprovedRowToDomain(row)
	}
	return events, nil
}

// =========================================================== MAPPERS ================================================================

func mapInvoiceCreatedRowToDomain(row db.FinanceInvoiceCreatedEvent) db.FinanceInvoiceCreatedEvent {
	return db.FinanceInvoiceCreatedEvent{
		ID:             row.ID,
		InvoiceID:      row.InvoiceID,
		InvoiceNumber:  row.InvoiceNumber,
		InvoiceDate:    row.InvoiceDate,
		Total:          row.Total,
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt,
	}
}

func mapPaymentReceivedRowToDomain(row db.FinancePaymentReceivedEvent) db.FinancePaymentReceivedEvent {
	return db.FinancePaymentReceivedEvent{
		ID:             row.ID,
		PaymentDueID:   row.PaymentDueID,
		InvoiceID:      row.InvoiceID,
		AmountPaid:     row.AmountPaid,
		PaidAt:         row.PaidAt,
		Reference:      row.Reference,
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt,
	}
}

func mapInventoryCostPostedRowToDomain(row db.InventoryCostPostedEvent) db.InventoryCostPostedEvent {
	return db.InventoryCostPostedEvent{
		ID:             row.ID,
		ReferenceType:  row.ReferenceType,
		ReferenceID:    row.ReferenceID,
		Amount:         row.Amount,
		CostCenterID:   row.CostCenterID,
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt,
	}
}

func mapPayrollPostedRowToDomain(row db.PayrollPostedEvent) db.PayrollPostedEvent {
	return db.PayrollPostedEvent{
		ID:             row.ID,
		PayrollRunID:   row.PayrollRunID,
		TotalGross:     row.TotalGross,
		TotalNet:       row.TotalNet,
		RunDate:        row.RunDate,
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt,
	}
}

func mapVendorBillApprovedRowToDomain(row db.VendorBillApprovedEvent) db.VendorBillApprovedEvent {
	return db.VendorBillApprovedEvent{
		ID:             row.ID,
		VendorBillID:   row.VendorBillID,
		Amount:         row.Amount,
		ApprovedAt:     row.ApprovedAt,
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt,
	}
}
