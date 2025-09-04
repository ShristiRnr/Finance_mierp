-- =====================================================
-- Accounts
-- =====================================================

-- name: CreateAccount :one
INSERT INTO accounts (code, name, type, parent_id, status, allow_manual_journal, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id = $1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY code LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET code = $2, name = $3, type = $4, parent_id = $5,
    status = $6, allow_manual_journal = $7,
    updated_by = $8, updated_at = now(), revision = revision + 1
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;

-- =====================================================
-- Journal Entries
-- =====================================================

-- name: CreateJournalEntry :one
INSERT INTO journal_entries (journal_date, reference, memo, source_type, source_id, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetJournalEntry :one
SELECT * FROM journal_entries WHERE id = $1;

-- name: ListJournalEntries :many
SELECT * FROM journal_entries ORDER BY journal_date DESC LIMIT $1 OFFSET $2;

-- name: UpdateJournalEntry :one
UPDATE journal_entries
SET journal_date = $2, reference = $3, memo = $4,
    source_type = $5, source_id = $6,
    updated_by = $7, updated_at = now(), revision = revision + 1
WHERE id = $1
RETURNING *;

-- name: DeleteJournalEntry :exec
DELETE FROM journal_entries WHERE id = $1;

-- =====================================================
-- Ledger Entries
-- =====================================================

-- name: ListLedgerEntries :many
SELECT * FROM ledger_entries ORDER BY transaction_date DESC LIMIT $1 OFFSET $2;