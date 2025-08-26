-- =====================================================
-- Budgets
-- =====================================================

-- name: CreateBudget :one
INSERT INTO budgets (name, total_amount, status, created_by, updated_by)
VALUES ($1, $2, COALESCE($3, 'DRAFT'), $4, $5)
RETURNING *;

-- name: GetBudget :one
SELECT * FROM budgets WHERE id = $1;

-- name: ListBudgets :many
SELECT * FROM budgets
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateBudget :one
UPDATE budgets
SET 
    name = $2,
    total_amount = $3,
    status = $4,
    updated_by = $5,
    updated_at = now(),
    revision = revision + 1
WHERE id = $1
RETURNING *;

-- name: DeleteBudget :exec
DELETE FROM budgets WHERE id = $1;

-- =====================================================
-- Budget Allocations
-- =====================================================

-- name: AllocateBudget :one
INSERT INTO budget_allocations (budget_id, department_id, allocated_amount, spent_amount, created_by, updated_by)
VALUES ($1, $2, $3, COALESCE($4, 0), $5, $6)
RETURNING *;

-- name: GetBudgetAllocation :one
SELECT * FROM budget_allocations WHERE id = $1;

-- name: ListBudgetAllocations :many
SELECT * FROM budget_allocations
WHERE budget_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateBudgetAllocation :one
UPDATE budget_allocations
SET 
    department_id = $2,
    allocated_amount = $3,
    spent_amount = $4,
    updated_by = $5,
    updated_at = now(),
    revision = revision + 1
WHERE id = $1
RETURNING *;

-- name: DeleteBudgetAllocation :exec
DELETE FROM budget_allocations WHERE id = $1;

-- =====================================================
-- Budget Comparison (Report)
-- =====================================================

-- name: GetBudgetComparisonReport :one
SELECT 
    b.id as budget_id,
    b.total_amount as total_budget,
    COALESCE(SUM(a.allocated_amount), 0) as total_allocated,
    COALESCE(SUM(a.spent_amount), 0) as total_spent,
    (b.total_amount - COALESCE(SUM(a.spent_amount), 0)) as remaining_budget
FROM budgets b
LEFT JOIN budget_allocations a ON b.id = a.budget_id
WHERE b.id = $1
GROUP BY b.id, b.total_amount;
