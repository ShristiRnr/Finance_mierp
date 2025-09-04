-- Organizations / Tenants ka part alag schema mein hoga
-- फिलहाल सिर्फ common entities

-- name: CreateRequestMetadata :one
INSERT INTO request_metadata (request_id, organization_id, tenant_id, auth_subject, source_system, trace_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRequestMetadata :one
SELECT * FROM request_metadata WHERE id = $1;

-- name: CreateAuditFields :one
INSERT INTO audit_fields (created_by, updated_by, revision)
VALUES ($1, $2, $3)
RETURNING *;
