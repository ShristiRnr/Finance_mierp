-- ============================
-- Profit & Loss
-- ============================

-- name: GenerateProfitLossReport :one
INSERT INTO profit_loss_reports (
    organization_id, period_start, period_end, total_revenue, total_expenses, net_profit
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProfitLossReport :one
SELECT * FROM profit_loss_reports WHERE id = $1;

-- name: ListProfitLossReports :many
SELECT * FROM profit_loss_reports
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- ============================
-- Balance Sheet
-- ============================

-- name: GenerateBalanceSheetReport :one
INSERT INTO balance_sheet_reports (
    organization_id, period_start, period_end, total_assets, total_liabilities, net_worth
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetBalanceSheetReport :one
SELECT * FROM balance_sheet_reports WHERE id = $1;

-- name: ListBalanceSheetReports :many
SELECT * FROM balance_sheet_reports
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- ============================
-- Trial Balance
-- ============================

-- name: CreateTrialBalanceReport :one
INSERT INTO trial_balance_reports (
    organization_id, period_start, period_end
) VALUES ($1, $2, $3)
RETURNING *;

-- name: AddTrialBalanceEntry :one
INSERT INTO trial_balance_entries (
    report_id, ledger_account, debit, credit
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListTrialBalanceEntries :many
SELECT * FROM trial_balance_entries WHERE report_id = $1;

-- name: GetTrialBalanceReport :one
SELECT * FROM trial_balance_reports WHERE id = $1;

-- name: ListTrialBalanceReports :many
SELECT * FROM trial_balance_reports
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- ============================
-- Compliance Reports
-- ============================

-- name: GenerateComplianceReport :one
INSERT INTO compliance_reports (
    organization_id, period_start, period_end, jurisdiction, details
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetComplianceReport :one
SELECT * FROM compliance_reports WHERE id = $1;

-- name: ListComplianceReports :many
SELECT *
FROM compliance_reports
WHERE organization_id = $1
  AND jurisdiction = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;