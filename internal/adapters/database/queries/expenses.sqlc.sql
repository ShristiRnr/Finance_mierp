-- =====================================================
-- Expense Queries
-- =====================================================

-- name: CreateExpense :one
INSERT INTO expenses (category, amount, expense_date, cost_center_id, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetExpense :one
SELECT * FROM expenses WHERE id = $1;

-- name: ListExpenses :many
SELECT * FROM expenses
ORDER BY expense_date DESC
LIMIT $1 OFFSET $2;

-- name: UpdateExpense :one
UPDATE expenses
SET category = $2,
    amount = $3,
    expense_date = $4,
    cost_center_id = $5,
    updated_by = $6,
    revision = revision + 1,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1;

-- name: AddExpenseExternalRef :one
INSERT INTO expense_external_refs (expense_id, system, ref_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListExpenseExternalRefs :many
SELECT * FROM expense_external_refs WHERE expense_id = $1;

-- =====================================================
-- Cost Center Queries
-- =====================================================

-- name: CreateCostCenter :one
INSERT INTO cost_centers (name, description, created_by, updated_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCostCenter :one
SELECT * FROM cost_centers WHERE id = $1;

-- name: ListCostCenters :many
SELECT * FROM cost_centers
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateCostCenter :one
UPDATE cost_centers
SET name = $2,
    description = $3,
    updated_by = $4,
    revision = revision + 1,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCostCenter :exec
DELETE FROM cost_centers WHERE id = $1;

-- =====================================================
-- Cost Allocation Queries
-- =====================================================

-- name: AllocateCost :one
INSERT INTO cost_allocations (cost_center_id, amount, reference_type, reference_id, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListCostAllocations :many
SELECT * FROM cost_allocations
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
