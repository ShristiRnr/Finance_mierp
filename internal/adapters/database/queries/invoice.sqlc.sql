-- name: CreateInvoice :one
INSERT INTO invoices (
    invoice_number, type, invoice_date, due_date, delivery_date, organization_id, po_number, eway_number_legacy, status_note, status, payment_reference, challan_number, challan_date,
    lr_number, transporter_name, transporter_id, vehicle_number, against_invoice_number, against_invoice_date, subtotal, gst_cgst, gst_sgst, gst_igst, gst_rate, grand_total,
    created_by, updated_by, revision
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12, $13, $14,
    $15, $16, $17, $18,
    $19, $20,
    $21, $22, $23, $24, $25, $26,
    $27, $28) RETURNING *;

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
-- name: UpdateInvoice :one
UPDATE invoices
SET 
    invoice_number       = $2,
    type                 = $3,
    invoice_date         = $4,
    due_date             = $5,
    delivery_date        = $6,
    organization_id      = $7,
    po_number            = $8,
    eway_number_legacy   = $9,
    status_note          = $10,
    status               = $11,
    payment_reference    = $12,
    challan_number       = $13,
    challan_date         = $14,
    lr_number            = $15,
    transporter_name     = $16,
    transporter_id       = $17,
    vehicle_number       = $18,
    against_invoice_number = $19,
    against_invoice_date = $20,
    subtotal             = $21,
    grand_total          = $22,
    gst_rate             = $23,
    gst_cgst             = $24,
    gst_sgst             = $25,
    gst_igst             = $26,
    updated_by           = $27,
    revision             = $28,
    updated_at           = now()
WHERE id = $1
RETURNING *;