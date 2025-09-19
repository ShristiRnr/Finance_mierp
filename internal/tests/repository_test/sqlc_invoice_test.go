package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
)

// helpers to convert Null types for sqlmock
func nt(t sql.NullTime) interface{} {
	if t.Valid {
		return t.Time
	}
	return nil
}

func ns(s sql.NullString) interface{} {
	if s.Valid {
		return s.String
	}
	return nil
}

func ni(i sql.NullInt32) interface{} {
	if i.Valid {
		return i.Int32
	}
	return nil
}

func TestUpdateInvoice(t *testing.T) {
	ctx := context.Background()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	queries := db.New(dbConn)
	repo := repository.NewInvoiceRepo(queries)

	invoiceID := uuid.New()
	now := time.Now()

	inv := db.Invoice{
		ID:                   invoiceID,
		InvoiceNumber:        "INV-001",
		Type:                 "SALE",
		InvoiceDate:          now,
		DueDate:              sql.NullTime{Time: now.AddDate(0, 0, 30), Valid: true},
		DeliveryDate:         sql.NullTime{Time: now.AddDate(0, 0, 2), Valid: true},
		OrganizationID:       uuid.New().String(),
		PoNumber:             sql.NullString{String: "PO-1001", Valid: true},
		EwayNumberLegacy:     sql.NullString{String: "EWAY123", Valid: true},
		StatusNote:           sql.NullString{String: "Pending approval", Valid: true},
		Status:               "DRAFT",
		PaymentReference:     sql.NullString{String: "PAY123", Valid: true},
		ChallanNumber:        sql.NullString{String: "CH123", Valid: true},
		ChallanDate:          sql.NullTime{Time: now, Valid: true},
		LrNumber:             sql.NullString{String: "LR987", Valid: true},
		TransporterName:      sql.NullString{String: "BlueDart", Valid: true},
		TransporterID:        sql.NullString{String: uuid.New().String(), Valid: true},
		VehicleNumber:        sql.NullString{String: "UP32AB1234", Valid: true},
		AgainstInvoiceNumber: sql.NullString{String: "INV-REF", Valid: true},
		AgainstInvoiceDate:   sql.NullTime{Time: now, Valid: true},
		Subtotal:             "1000.0",
		GrandTotal:           "1180.0",
		GstRate:              sql.NullString{String: fmt.Sprintf("%f", 18.0), Valid: true},
		GstCgst:              sql.NullString{String: fmt.Sprintf("%f", 90.0), Valid: true},
		GstSgst:              sql.NullString{String: fmt.Sprintf("%f", 90.0), Valid: true},
		GstIgst:              sql.NullString{String: fmt.Sprintf("%f", 0.0), Valid: true},
		CreatedBy:            sql.NullString{String: "creator", Valid: true},
		UpdatedBy:            sql.NullString{String: "tester", Valid: true},
		Revision:             sql.NullInt32{Int32: 2, Valid: true},
		CreatedAt:            sql.NullTime{Time: now, Valid: true},
		UpdatedAt:            sql.NullTime{Time: now, Valid: true},
	}

	columns := []string{
		"id", "invoice_number", "type", "invoice_date", "due_date", "delivery_date",
		"organization_id", "po_number", "eway_number_legacy", "status_note", "status",
		"payment_reference", "challan_number", "challan_date", "lr_number",
		"transporter_name", "transporter_id", "vehicle_number",
		"against_invoice_number", "against_invoice_date",
		"subtotal", "grand_total", "gst_rate", "gst_cgst", "gst_sgst", "gst_igst",
		"created_by", "updated_by", "revision",
		"created_at", "updated_at",
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows(columns).AddRow(
			inv.ID, inv.InvoiceNumber, inv.Type, inv.InvoiceDate,
			nt(inv.DueDate), nt(inv.DeliveryDate), inv.OrganizationID,
			ns(inv.PoNumber), ns(inv.EwayNumberLegacy), ns(inv.StatusNote), inv.Status,
			ns(inv.PaymentReference), ns(inv.ChallanNumber), nt(inv.ChallanDate),
			ns(inv.LrNumber), ns(inv.TransporterName), ns(inv.TransporterID),
			ns(inv.VehicleNumber), ns(inv.AgainstInvoiceNumber), nt(inv.AgainstInvoiceDate),
			inv.Subtotal, inv.GrandTotal, ns(inv.GstRate), ns(inv.GstCgst), ns(inv.GstSgst), ns(inv.GstIgst),
			ns(inv.CreatedBy), ns(inv.UpdatedBy), ni(inv.Revision),
			nt(inv.CreatedAt), nt(inv.UpdatedAt),
		)

		mock.ExpectQuery(`UPDATE invoices SET`).
			WithArgs(
				inv.ID,
				inv.InvoiceNumber, inv.Type, inv.InvoiceDate,
				inv.DueDate, inv.DeliveryDate, inv.OrganizationID,
				inv.PoNumber, inv.EwayNumberLegacy, inv.StatusNote, inv.Status,
				inv.PaymentReference, inv.ChallanNumber, inv.ChallanDate,
				inv.LrNumber, inv.TransporterName, inv.TransporterID,
				inv.VehicleNumber, inv.AgainstInvoiceNumber, inv.AgainstInvoiceDate,
				inv.Subtotal, inv.GrandTotal, inv.GstRate, inv.GstCgst,
				inv.GstSgst, inv.GstIgst,
				inv.UpdatedBy, inv.Revision,
			).
			WillReturnRows(rows)

		got, err := repo.UpdateInvoice(ctx, inv)
		require.NoError(t, err)
		require.Equal(t, inv.ID, got.ID)
		require.Equal(t, inv.InvoiceNumber, got.InvoiceNumber)
		require.Equal(t, inv.UpdatedBy, got.UpdatedBy)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`UPDATE invoices SET`).
			WithArgs(
				inv.ID,
				inv.InvoiceNumber, inv.Type, inv.InvoiceDate,
				inv.DueDate, inv.DeliveryDate, inv.OrganizationID,
				inv.PoNumber, inv.EwayNumberLegacy, inv.StatusNote, inv.Status,
				inv.PaymentReference, inv.ChallanNumber, inv.ChallanDate,
				inv.LrNumber, inv.TransporterName, inv.TransporterID,
				inv.VehicleNumber, inv.AgainstInvoiceNumber, inv.AgainstInvoiceDate,
				inv.Subtotal, inv.GrandTotal, inv.GstRate, inv.GstCgst,
				inv.GstSgst, inv.GstIgst,
				inv.UpdatedBy, inv.Revision,
			).WillReturnError(sql.ErrNoRows)

		_, err := repo.UpdateInvoice(ctx, inv)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`UPDATE invoices SET`).
			WithArgs(
				inv.ID,
				inv.InvoiceNumber, inv.Type, inv.InvoiceDate,
				inv.DueDate, inv.DeliveryDate, inv.OrganizationID,
				inv.PoNumber, inv.EwayNumberLegacy, inv.StatusNote, inv.Status,
				inv.PaymentReference, inv.ChallanNumber, inv.ChallanDate,
				inv.LrNumber, inv.TransporterName, inv.TransporterID,
				inv.VehicleNumber, inv.AgainstInvoiceNumber, inv.AgainstInvoiceDate,
				inv.Subtotal, inv.GrandTotal, inv.GstRate, inv.GstCgst,
				inv.GstSgst, inv.GstIgst,
				inv.UpdatedBy, inv.Revision,
			).WillReturnError(sql.ErrConnDone)

		_, err := repo.UpdateInvoice(ctx, inv)
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
	})
}
