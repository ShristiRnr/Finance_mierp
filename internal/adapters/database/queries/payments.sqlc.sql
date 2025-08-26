-- =====================================================
-- Payment Dues
-- =====================================================

-- name: CreatePaymentDue :one
INSERT INTO payment_dues (
    invoice_id, amount_due, due_date, status, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPaymentDue :one
SELECT * FROM payment_dues WHERE id = $1;

-- name: UpdatePaymentDue :one
UPDATE payment_dues
SET invoice_id = $2,
    amount_due = $3,
    due_date = $4,
    status = $5,
    updated_by = $6,
    updated_at = now(),
    revision = revision + 1
WHERE id = $1
RETURNING *;

-- name: DeletePaymentDue :exec
DELETE FROM payment_dues WHERE id = $1;

-- name: ListPaymentDues :many
SELECT * FROM payment_dues
ORDER BY due_date ASC
LIMIT $1 OFFSET $2;

-- name: MarkPaymentAsPaid :one
UPDATE payment_dues
SET status = 'PAID',
    updated_at = now(),
    updated_by = $2,
    revision = revision + 1
WHERE id = $1
RETURNING *;

-- =====================================================
-- Bank Accounts
-- =====================================================

-- name: CreateBankAccount :one
INSERT INTO bank_accounts (
    name, account_number, ifsc_or_swift, ledger_account_id, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetBankAccount :one
SELECT * FROM bank_accounts WHERE id = $1;

-- name: UpdateBankAccount :one
UPDATE bank_accounts
SET name = $2,
    account_number = $3,
    ifsc_or_swift = $4,
    ledger_account_id = $5,
    updated_by = $6,
    updated_at = now(),
    revision = revision + 1
WHERE id = $1
RETURNING *;

-- name: DeleteBankAccount :exec
DELETE FROM bank_accounts WHERE id = $1;

-- name: ListBankAccounts :many
SELECT * FROM bank_accounts
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- =====================================================
-- Bank Transactions
-- =====================================================

-- name: ImportBankTransaction :one
INSERT INTO bank_transactions (
    bank_account_id, amount, transaction_date, description, reference,
    reconciled, matched_reference_type, matched_reference_id,
    created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8,
    $9, $10
) RETURNING *;

-- name: ListBankTransactions :many
SELECT * FROM bank_transactions
WHERE bank_account_id = $1
ORDER BY transaction_date DESC
LIMIT $2 OFFSET $3;

-- name: ReconcileTransaction :one
UPDATE bank_transactions
SET reconciled = true,
    matched_reference_type = $2,
    matched_reference_id = $3,
    updated_at = now(),
    updated_by = $4,
    revision = revision + 1
WHERE id = $1
RETURNING *;