-- name: CreateInvoice :one
INSERT INTO invoices (
    invoice_number, type, invoice_date, due_date, delivery_date, 
    party_ref_id, organization_id, po_number, eway_number_legacy, 
    status_note, status, payment_reference, challan_number, challan_date,
    lr_number, transporter_name, transporter_id, vehicle_number,
    against_invoice_number, against_invoice_date,
    subtotal, gst_cgst, gst_sgst, gst_igst, gst_rate, grand_total,
    created_by, updated_by, revision
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12, $13, $14,
    $15, $16, $17, $18,
    $19, $20,
    $21, $22, $23, $24, $25, $26,
    $27, $28, $29
) RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices WHERE id = $1;

-- name: ListInvoices :many
SELECT * FROM invoices ORDER BY invoice_date DESC LIMIT $1 OFFSET $2;

-- name: DeleteInvoice :exec
DELETE FROM invoices WHERE id = $1;

-- name: CreateInvoiceItem :one
INSERT INTO invoice_items (
    invoice_id, name, description, hsn, quantity, unit_price, line_subtotal, line_total, cost_center_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: ListInvoiceItems :many
SELECT * FROM invoice_items WHERE invoice_id = $1;

-- name: AddInvoiceTax :one
INSERT INTO invoice_taxes (invoice_id, name, rate, amount)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: AddInvoiceDiscount :one
INSERT INTO invoice_discounts (invoice_id, description, amount)
VALUES ($1, $2, $3) RETURNING *;

-- name: SearchInvoices :many
SELECT *
FROM invoices
WHERE 
    (invoice_number ILIKE '%' || $1 || '%'
     OR po_number ILIKE '%' || $1 || '%'
     OR status_note ILIKE '%' || $1 || '%'
     OR payment_reference ILIKE '%' || $1 || '%'
     OR transporter_name ILIKE '%' || $1 || '%'
     OR transporter_id ILIKE '%' || $1 || '%'
     OR vehicle_number ILIKE '%' || $1 || '%')
ORDER BY invoice_date DESC
LIMIT $2 OFFSET $3;

-- name: UpdateInvoice :one
UPDATE invoices
SET 
    invoice_number = $2,
    type = $3,
    invoice_date = $4,
    due_date = $5,
    delivery_date = $6,
    party_ref_id = $7,
    organization_id = $8,
    po_number = $9,
    eway_number_legacy = $10,
    status_note = $11,
    status = $12,
    payment_reference = $13,
    challan_number = $14,
    challan_date = $15,
    lr_number = $16,
    transporter_name = $17,
    transporter_id = $18,
    vehicle_number = $19,
    against_invoice_number = $20,
    against_invoice_date = $21,
    subtotal = $22,
    grand_total = $23,
    gst_rate = $24,
    gst_cgst = $25,
    gst_sgst = $26,
    gst_igst = $27,
    updated_by = $28,
    revision = $29,
    updated_at = now()
WHERE id = $1
RETURNING *;