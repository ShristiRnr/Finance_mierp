package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/adapters/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite" // pure-Go SQLite driver
)

// setupTestDB connects to an in-memory SQLite DB for testing.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	sqlDB, err := sql.Open("sqlite", "file::memory:?cache=shared")
	require.NoError(t, err)

	// Create tables
	_, err = sqlDB.Exec(`
	CREATE TABLE invoices (
		id TEXT PRIMARY KEY,
		invoice_number TEXT,
		type TEXT,
		invoice_date DATETIME,
		due_date DATETIME,
		delivery_date DATETIME,
		organization_id TEXT,
		po_number TEXT,
		eway_number_legacy TEXT,
		status_note TEXT,
		status TEXT,
		payment_reference TEXT,
		challan_number TEXT,
		challan_date DATETIME,
		lr_number TEXT,
		transporter_name TEXT,
		transporter_id TEXT,
		vehicle_number TEXT,
		against_invoice_number TEXT,
		against_invoice_date DATETIME,
		subtotal TEXT,
		grand_total TEXT,
		gst_rate TEXT,
		gst_cgst TEXT,
		gst_sgst TEXT,
		gst_igst TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT,
		revision INTEGER
	);
	CREATE TABLE invoice_items (
		id TEXT PRIMARY KEY,
		invoice_id TEXT,
		name TEXT,
		description TEXT,
		hsn TEXT,
		quantity INTEGER,
		unit_price TEXT,
		line_subtotal TEXT,
		line_total TEXT,
		cost_center_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT,
		revision INTEGER
	);
	CREATE TABLE invoice_taxes (
		id TEXT PRIMARY KEY,
		invoice_id TEXT,
		name TEXT,
		rate TEXT,
		amount TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		revision INTEGER
	);
	CREATE TABLE invoice_discounts (
		id TEXT PRIMARY KEY,
		invoice_id TEXT,
		description TEXT,
		amount TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		revision INTEGER
	);
	`)
	require.NoError(t, err)

	cleanup := func() { sqlDB.Close() }
	return sqlDB, cleanup
}

// helper function to create a test invoice
func createTestInvoice(t *testing.T, q *db.Queries) db.Invoice {
	arg := db.CreateInvoiceParams{
		InvoiceNumber:  "INV-" + uuid.New().String(),
		Type:           "SALES",
		InvoiceDate:    time.Now(),
		DueDate:        sql.NullTime{Time: time.Now().AddDate(0, 0, 15), Valid: true},
		OrganizationID: uuid.New().String(),
		Status:         "PENDING",
		Subtotal:       "1000",
		GrandTotal:     "1180",
		CreatedBy:      sql.NullString{String: "tester", Valid: true},
		UpdatedBy:      sql.NullString{String: "tester", Valid: true},
		Revision:       sql.NullInt32{Int32: 1, Valid: true},
	}
	inv, err := q.CreateInvoice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, inv)
	return inv
}

func TestInvoiceRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	q := db.New(sqlDB)
	repo := repository.NewInvoiceRepo(q)

	// Create
	inv := createTestInvoice(t, q)
	require.NotEmpty(t, inv.ID)

	// Get
	fetched, err := repo.GetInvoice(ctx, inv.ID)
	require.NoError(t, err)
	require.Equal(t, inv.ID, fetched.ID)

	// List
	list, err := repo.ListInvoices(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)

	// Update
	inv.Status = "APPROVED"
	inv.UpdatedBy = sql.NullString{String: "manager", Valid: true}
	updated, err := repo.UpdateInvoice(ctx, inv)
	require.NoError(t, err)
	require.Equal(t, "APPROVED", updated.Status)

	// Delete
	err = repo.DeleteInvoice(ctx, inv.ID)
	require.NoError(t, err)

	_, err = repo.GetInvoice(ctx, inv.ID)
	require.Error(t, err)
}

func TestInvoiceItemsAndTaxes(t *testing.T) {
	ctx := context.Background()
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	q := db.New(sqlDB)
	repo := repository.NewInvoiceRepo(q)

	// Create invoice
	inv := createTestInvoice(t, q)

	// Add Item
	item := db.CreateInvoiceItemParams{
		InvoiceID:    inv.ID,
		Name:         "Product A",
		Description:  sql.NullString{String: "Test product", Valid: true},
		Hsn:          sql.NullString{String: "1001", Valid: true},
		Quantity:     2,
		UnitPrice:    "200",
		LineSubtotal: "400",
		LineTotal:    "400",
		CostCenterID: sql.NullString{Valid: false},
	}
	createdItem, err := repo.CreateInvoiceItem(ctx, item)
	require.NoError(t, err)
	require.Equal(t, "Product A", createdItem.Name)

	items, err := repo.ListInvoiceItems(ctx, inv.ID)
	require.NoError(t, err)
	require.Len(t, items, 1)

	// Add Tax
	tax := db.AddInvoiceTaxParams{
		InvoiceID: inv.ID,
		Name:      "GST",
		Rate:      "18",
		Amount:    "72",
	}
	createdTax, err := repo.AddInvoiceTax(ctx, tax)
	require.NoError(t, err)
	require.Equal(t, "GST", createdTax.Name)

	// Add Discount
	disc := db.AddInvoiceDiscountParams{
		InvoiceID:   inv.ID,
		Description: sql.NullString{String: "Promo", Valid: true},
		Amount:      "50",
	}
	createdDisc, err := repo.AddInvoiceDiscount(ctx, disc)
	require.NoError(t, err)
	require.Equal(t, "50", createdDisc.Amount)
}