-- name: AddGstBreakup :one
INSERT INTO gst_breakups (invoice_id, taxable_amount, cgst, sgst, igst, total_gst)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetGstBreakup :one
SELECT * FROM gst_breakups WHERE invoice_id = $1;

-- name: AddGstRegime :one
INSERT INTO gst_regimes (invoice_id, gstin, place_of_supply, reverse_charge)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGstRegime :one
SELECT * FROM gst_regimes WHERE invoice_id = $1;

-- name: AddGstDocStatus :one
INSERT INTO gst_doc_statuses (
    invoice_id, einvoice_status, irn, ack_no, ack_date,
    eway_status, eway_bill_no, eway_valid_upto, last_error, last_synced_at
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING *;

-- name: GetGstDocStatus :one
SELECT * FROM gst_doc_statuses WHERE invoice_id = $1;