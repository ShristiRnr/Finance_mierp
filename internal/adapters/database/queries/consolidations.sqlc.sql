-- name: CreateConsolidation :one
INSERT INTO consolidations (
    entity_ids, period_start, period_end, report
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetConsolidation :one
SELECT * FROM consolidations WHERE id = $1;

-- name: ListConsolidations :many
SELECT *
FROM consolidations
WHERE (entity_ids && $1::text[])  -- overlap with requested entities
  AND period_start >= $2
  AND period_end <= $3
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: DeleteConsolidation :exec
DELETE FROM consolidations WHERE id = $1;