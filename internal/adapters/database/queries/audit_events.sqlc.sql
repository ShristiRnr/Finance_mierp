-- name: RecordAuditEvent :one
INSERT INTO audit_events (
    user_id, action, timestamp, details, resource_type, resource_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: ListAuditEvents :many
SELECT * FROM audit_events
ORDER BY timestamp DESC
LIMIT $1 OFFSET $2;

-- name: GetAuditEventById :one
SELECT * FROM audit_events WHERE id = $1;

-- name: FilterAuditEvents :many
SELECT *
FROM audit_events
WHERE 
    ($1::text IS NULL OR user_id = $1)
    AND ($2::text IS NULL OR action = $2)
    AND ($3::text IS NULL OR resource_type = $3)
    AND ($4::text IS NULL OR resource_id = $4)
    AND ($5::timestamptz IS NULL OR timestamp >= $5)
    AND ($6::timestamptz IS NULL OR timestamp <= $6)
ORDER BY timestamp DESC
LIMIT $7 OFFSET $8;