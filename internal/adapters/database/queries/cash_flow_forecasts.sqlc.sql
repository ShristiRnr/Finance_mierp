-- name: GenerateCashFlowForecast :one
INSERT INTO cash_flow_forecasts (
    organization_id, period_start, period_end, forecast_details
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCashFlowForecast :one
SELECT * 
FROM cash_flow_forecasts
WHERE id = $1;

-- name: ListCashFlowForecasts :many
SELECT *
FROM cash_flow_forecasts
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;