-- =====================================================
-- Accrual Queries
-- =====================================================

-- name: CreateAccrual :one
INSERT INTO accruals (
    description, amount, accrual_date, account_id,
    created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAccrualById :one
SELECT * FROM accruals WHERE id = $1;

-- name: ListAccruals :many
SELECT * FROM accruals
ORDER BY accrual_date DESC
LIMIT $1 OFFSET $2;

-- name: UpdateAccrual :one
UPDATE accruals
SET description = $2,
    amount = $3,
    accrual_date = $4,
    account_id = $5,
    updated_by = $6,
    revision = revision + 1,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAccrual :exec
DELETE FROM accruals WHERE id = $1;

-- External refs for accrual
-- name: AddAccrualExternalRef :one
INSERT INTO accrual_external_refs (accrual_id, system, ref_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: ListAccrualExternalRefs :many
SELECT * FROM accrual_external_refs WHERE accrual_id = $1;

-- =====================================================
-- Allocation Rule Queries
-- =====================================================

-- name: CreateAllocationRule :one
INSERT INTO allocation_rules (
    name, basis, source_account_id, target_cost_center_ids,
    formula, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetAllocationRule :one
SELECT * FROM allocation_rules WHERE id = $1;

-- name: ListAllocationRules :many
SELECT * FROM allocation_rules
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateAllocationRule :one
UPDATE allocation_rules
SET name = $2,
    basis = $3,
    source_account_id = $4,
    target_cost_center_ids = $5,
    formula = $6,
    updated_by = $7,
    revision = revision + 1,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAllocationRule :exec
DELETE FROM allocation_rules WHERE id = $1;
