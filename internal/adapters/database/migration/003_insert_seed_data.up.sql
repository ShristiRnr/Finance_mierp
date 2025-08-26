-- =====================================================
-- Seed Data for Core Lookups & Defaults
-- =====================================================

-- Chart of Accounts (basic structure)
INSERT INTO accounts (code, name, type, status, allow_manual_journal, created_by)
VALUES 
('1000', 'Cash', 'ASSET', 'ACTIVE', true, 'system'),
('2000', 'Accounts Payable', 'LIABILITY', 'ACTIVE', true, 'system'),
('3000', 'Revenue', 'INCOME', 'ACTIVE', false, 'system'),
('4000', 'Operating Expenses', 'EXPENSE', 'ACTIVE', true, 'system')
ON CONFLICT (code) DO NOTHING;

-- Party Refs: Seed one default Customer & Vendor
INSERT INTO party_refs (kind, display_name)
VALUES
(1, 'Default Customer'),
(2, 'Default Vendor')
ON CONFLICT DO NOTHING;

-- Cost Centers
INSERT INTO cost_centers (name, description, created_by)
VALUES
('General', 'Default Cost Center', 'system'),
('Sales', 'Sales Department', 'system'),
('Operations', 'Operations Department', 'system')
ON CONFLICT DO NOTHING;

-- GST Regimes (example)
INSERT INTO gst_regimes (invoice_id, gstin, place_of_supply, reverse_charge, created_by)
SELECT id, '22AAAAA0000A1Z5', '22-Chhattisgarh', false, 'system'
FROM invoices LIMIT 1
ON CONFLICT DO NOTHING;

-- Exchange Rates (sample)
INSERT INTO exchange_rates (base_currency, quote_currency, rate, as_of, created_by)
VALUES
('USD', 'INR', 83.5000, now(), 'system'),
('EUR', 'INR', 91.2000, now(), 'system'),
('GBP', 'INR', 105.0000, now(), 'system')
ON CONFLICT DO NOTHING;

-- Allocation Rules
INSERT INTO allocation_rules (name, basis, source_account_id, target_cost_center_ids, formula, created_by)
VALUES
('Revenue Split by Sales', 'REVENUE', '3000', ARRAY['Sales','Operations'], 'Split by revenue contribution', 'system')
ON CONFLICT DO NOTHING;

-- Default Budgets
INSERT INTO budgets (name, total_amount, status, created_by)
VALUES
('Annual Budget 2025', 1000000.00, 'ACTIVE', 'system')
ON CONFLICT DO NOTHING;

-- Compliance Reports (placeholder)
INSERT INTO compliance_reports (organization_id, period_start, period_end, jurisdiction, details, created_by)
VALUES
('ORG001', date_trunc('year', now()), date_trunc('year', now()) + interval '1 year' - interval '1 day', 
 'IN-GST', '{"note":"Initial compliance placeholder"}', 'system')
ON CONFLICT DO NOTHING;

-- Audit Events (marker event)
INSERT INTO audit_events (user_id, action, details, resource_type, resource_id)
VALUES
('system', 'SEED_DATA_INSERTED', 'Initial seed data loaded', 'SYSTEM', 'seed-v3')
ON CONFLICT DO NOTHING;
