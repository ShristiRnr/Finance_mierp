package repository

import (
	"context"
	"database/sql"

	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
)

type FinanceEventRepo struct {
	q *db.Queries
}

func NewFinanceEventRepo(q *db.Queries) ports.FinanceEventRepository {
	return &FinanceEventRepo{q: q}
}

//
// ─── INVOICE CREATED ─────────────────────────────────────────────────────────────
//

func (r *FinanceEventRepo) InsertInvoiceCreated(ctx context.Context, e domain.FinanceInvoiceCreatedEvent) (domain.FinanceInvoiceCreatedEvent, error) {
	params := db.InsertInvoiceCreatedEventParams{
		InvoiceID:      e.InvoiceID,
		InvoiceNumber:  e.InvoiceNumber,
		InvoiceDate:    e.InvoiceDate,
		PartyRefID:     e.PartyRefID,
		Total:          e.Total,
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertInvoiceCreatedEvent(ctx, params)
	if err != nil {
		return domain.FinanceInvoiceCreatedEvent{}, err
	}
	return mapInvoiceCreatedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListInvoiceCreated(ctx context.Context, orgID string, limit, offset int32) ([]domain.FinanceInvoiceCreatedEvent, error) {
	rows, err := r.q.ListInvoiceCreatedEvents(ctx, db.ListInvoiceCreatedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]domain.FinanceInvoiceCreatedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapInvoiceCreatedRowToDomain(row)
	}
	return events, nil
}

//
// ─── PAYMENT RECEIVED ───────────────────────────────────────────────────────────
//

func (r *FinanceEventRepo) InsertPaymentReceived(ctx context.Context, e domain.FinancePaymentReceivedEvent) (domain.FinancePaymentReceivedEvent, error) {
	params := db.InsertPaymentReceivedEventParams{
		PaymentDueID:   e.PaymentDueID,
		InvoiceID:      e.InvoiceID,
		AmountPaid:     e.AmountPaid,
		PaidAt:         e.PaidAt,
		Reference:      toNullString(e.Reference),
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertPaymentReceivedEvent(ctx, params)
	if err != nil {
		return domain.FinancePaymentReceivedEvent{}, err
	}
	return mapPaymentReceivedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListPaymentReceived(ctx context.Context, orgID string, limit, offset int32) ([]domain.FinancePaymentReceivedEvent, error) {
	rows, err := r.q.ListPaymentReceivedEvents(ctx, db.ListPaymentReceivedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]domain.FinancePaymentReceivedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapPaymentReceivedRowToDomain(row)
	}
	return events, nil
}

//
// ─── INVENTORY COST POSTED ──────────────────────────────────────────────────────
//

func (r *FinanceEventRepo) InsertInventoryCostPosted(ctx context.Context, e domain.InventoryCostPostedEvent) (domain.InventoryCostPostedEvent, error) {
	params := db.InsertInventoryCostPostedEventParams{
		ReferenceType:  e.ReferenceType,
		ReferenceID:    e.ReferenceID,
		Amount:         e.Amount,
		CostCenterID:   toNullString(e.CostCenterID),
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertInventoryCostPostedEvent(ctx, params)
	if err != nil {
		return domain.InventoryCostPostedEvent{}, err
	}
	return mapInventoryCostPostedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListInventoryCostPosted(ctx context.Context, orgID string, limit, offset int32) ([]domain.InventoryCostPostedEvent, error) {
	rows, err := r.q.ListInventoryCostPostedEvents(ctx, db.ListInventoryCostPostedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]domain.InventoryCostPostedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapInventoryCostPostedRowToDomain(row)
	}
	return events, nil
}

//
// ─── PAYROLL POSTED ─────────────────────────────────────────────────────────────
//

func (r *FinanceEventRepo) InsertPayrollPosted(ctx context.Context, e domain.PayrollPostedEvent) (domain.PayrollPostedEvent, error) {
	params := db.InsertPayrollPostedEventParams{
		PayrollRunID:   e.PayrollRunID,
		TotalGross:     e.TotalGross,
		TotalNet:       e.TotalNet,
		RunDate:        e.RunDate,
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertPayrollPostedEvent(ctx, params)
	if err != nil {
		return domain.PayrollPostedEvent{}, err
	}
	return mapPayrollPostedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListPayrollPosted(ctx context.Context, orgID string, limit, offset int32) ([]domain.PayrollPostedEvent, error) {
	rows, err := r.q.ListPayrollPostedEvents(ctx, db.ListPayrollPostedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]domain.PayrollPostedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapPayrollPostedRowToDomain(row)
	}
	return events, nil
}

//
// ─── VENDOR BILL APPROVED ───────────────────────────────────────────────────────
//

func (r *FinanceEventRepo) InsertVendorBillApproved(ctx context.Context, e domain.VendorBillApprovedEvent) (domain.VendorBillApprovedEvent, error) {
	params := db.InsertVendorBillApprovedEventParams{
		VendorBillID:   e.VendorBillID,
		Amount:         e.Amount,
		ApprovedAt:     e.ApprovedAt,
		OrganizationID: e.OrganizationID,
	}
	row, err := r.q.InsertVendorBillApprovedEvent(ctx, params)
	if err != nil {
		return domain.VendorBillApprovedEvent{}, err
	}
	return mapVendorBillApprovedRowToDomain(row), nil
}

func (r *FinanceEventRepo) ListVendorBillApproved(ctx context.Context, orgID string, limit, offset int32) ([]domain.VendorBillApprovedEvent, error) {
	rows, err := r.q.ListVendorBillApprovedEvents(ctx, db.ListVendorBillApprovedEventsParams{
		OrganizationID: orgID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}
	events := make([]domain.VendorBillApprovedEvent, len(rows))
	for i, row := range rows {
		events[i] = mapVendorBillApprovedRowToDomain(row)
	}
	return events, nil
}

//
// ─── MAPPERS ────────────────────────────────────────────────────────────────────
//

func mapInvoiceCreatedRowToDomain(row db.FinanceInvoiceCreatedEvent) domain.FinanceInvoiceCreatedEvent {
	return domain.FinanceInvoiceCreatedEvent{
		ID:             row.ID,
		InvoiceID:      row.InvoiceID,
		InvoiceNumber:  row.InvoiceNumber,
		InvoiceDate:    row.InvoiceDate,
		PartyRefID:     row.PartyRefID,
		Total:          row.Total,
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt.Time,
	}
}

func mapPaymentReceivedRowToDomain(row db.FinancePaymentReceivedEvent) domain.FinancePaymentReceivedEvent {
	return domain.FinancePaymentReceivedEvent{
		ID:             row.ID,
		PaymentDueID:   row.PaymentDueID,
		InvoiceID:      row.InvoiceID,
		AmountPaid:     row.AmountPaid,
		PaidAt:         row.PaidAt,
		Reference:      fromNullString(row.Reference),
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt.Time,
	}
}

func mapInventoryCostPostedRowToDomain(row db.InventoryCostPostedEvent) domain.InventoryCostPostedEvent {
	return domain.InventoryCostPostedEvent{
		ID:             row.ID,
		ReferenceType:  row.ReferenceType,
		ReferenceID:    row.ReferenceID,
		Amount:         row.Amount,
		CostCenterID:   fromNullString(row.CostCenterID),
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt.Time,
	}
}

func mapPayrollPostedRowToDomain(row db.PayrollPostedEvent) domain.PayrollPostedEvent {
	return domain.PayrollPostedEvent{
		ID:             row.ID,
		PayrollRunID:   row.PayrollRunID,
		TotalGross:     row.TotalGross,
		TotalNet:       row.TotalNet,
		RunDate:        row.RunDate,
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt.Time,
	}
}

func mapVendorBillApprovedRowToDomain(row db.VendorBillApprovedEvent) domain.VendorBillApprovedEvent {
	return domain.VendorBillApprovedEvent{
		ID:             row.ID,
		VendorBillID:   row.VendorBillID,
		Amount:         row.Amount,
		ApprovedAt:     row.ApprovedAt,
		OrganizationID: row.OrganizationID,
		CreatedAt:      row.CreatedAt.Time,
	}
}

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func fromNullString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}
