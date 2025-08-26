-- Create a new exchange rate
-- name: CreateExchangeRate :one
INSERT INTO exchange_rates (
    base_currency, quote_currency, rate, as_of,
    created_by, updated_by, revision
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7
) RETURNING *;

-- Fetch single exchange rate by ID
-- name: GetExchangeRate :one
SELECT * FROM exchange_rates WHERE id = $1;

-- Update exchange rate
-- name: UpdateExchangeRate :one
UPDATE exchange_rates
SET
    base_currency = $2,
    quote_currency = $3,
    rate = $4,
    as_of = $5,
    updated_by = $6,
    revision = $7,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- Delete exchange rate
-- name: DeleteExchangeRate :exec
DELETE FROM exchange_rates WHERE id = $1;

-- List exchange rates (with optional filters)
-- name: ListExchangeRates :many
SELECT *
FROM exchange_rates
WHERE
    ($1::TEXT IS NULL OR base_currency = $1)
    AND ($2::TEXT IS NULL OR quote_currency = $2)
ORDER BY as_of DESC
LIMIT $3 OFFSET $4;

-- Fetch the latest rate for a currency pair as of timestamp
-- name: GetLatestRate :one
SELECT *
FROM exchange_rates
WHERE base_currency = $1
  AND quote_currency = $2
  AND as_of <= $3
ORDER BY as_of DESC
LIMIT 1;