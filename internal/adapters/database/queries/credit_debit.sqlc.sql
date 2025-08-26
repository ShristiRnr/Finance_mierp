-- name: CreateCreditDebitNote :one
INSERT INTO credit_debit_notes (
    invoice_id, type, amount, reason, created_by, updated_by, revision
) VALUES (
    $1, $2, $3, $4, $5, $6, 1
) RETURNING *;

-- name: GetCreditDebitNote :one
SELECT * FROM credit_debit_notes WHERE id = $1;

-- name: ListCreditDebitNotes :many
SELECT * 
FROM credit_debit_notes
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCreditDebitNote :one
UPDATE credit_debit_notes
SET 
    invoice_id = $2,
    type = $3,
    amount = $4,
    reason = $5,
    updated_by = $6,
    revision = revision + 1,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCreditDebitNote :exec
DELETE FROM credit_debit_notes WHERE id = $1;

-- External Refs --------------------------------------

-- name: AddCreditDebitNoteExternalRef :one
INSERT INTO credit_debit_note_external_refs (note_id, system, ref_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: ListCreditDebitNoteExternalRefs :many
SELECT * FROM credit_debit_note_external_refs WHERE note_id = $1;