-- =====================================================
-- Finance Invoice Created Events
-- =====================================================
-- name: InsertInvoiceCreatedEvent :one
INSERT INTO finance_invoice_created_events (
    invoice_id, invoice_number, invoice_date, party_ref_id, total, organization_id
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListInvoiceCreatedEvents :many
SELECT * 
FROM finance_invoice_created_events
WHERE organization_id = $1
ORDER BY invoice_date DESC
LIMIT $2 OFFSET $3;

-- =====================================================
-- Finance Payment Received Events
-- =====================================================
-- name: InsertPaymentReceivedEvent :one
INSERT INTO finance_payment_received_events (
    payment_due_id, invoice_id, amount_paid, paid_at, reference, organization_id
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListPaymentReceivedEvents :many
SELECT * 
FROM finance_payment_received_events
WHERE organization_id = $1
ORDER BY paid_at DESC
LIMIT $2 OFFSET $3;

-- =====================================================
-- Inventory Cost Posted Events
-- =====================================================
-- name: InsertInventoryCostPostedEvent :one
INSERT INTO inventory_cost_posted_events (
    reference_type, reference_id, amount, cost_center_id, organization_id
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListInventoryCostPostedEvents :many
SELECT * 
FROM inventory_cost_posted_events
WHERE organization_id = $1
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- =====================================================
-- Payroll Posted Events
-- =====================================================
-- name: InsertPayrollPostedEvent :one
INSERT INTO payroll_posted_events (
    payroll_run_id, total_gross, total_net, run_date, organization_id
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListPayrollPostedEvents :many
SELECT *
FROM payroll_posted_events
WHERE organization_id = $1
ORDER BY run_date DESC
LIMIT $2 OFFSET $3;

-- =====================================================
-- Vendor Bill Approved Events
-- =====================================================
-- name: InsertVendorBillApprovedEvent :one
INSERT INTO vendor_bill_approved_events (
    vendor_bill_id, amount, approved_at, organization_id
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListVendorBillApprovedEvents :many
SELECT *
FROM vendor_bill_approved_events
WHERE organization_id = $1
ORDER BY approved_at DESC
LIMIT $2 OFFSET $3;