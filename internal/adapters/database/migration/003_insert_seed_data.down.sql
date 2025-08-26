-- Rollback Seed Data

DELETE FROM audit_events WHERE resource_id = 'seed-v3';

DELETE FROM compliance_reports WHERE jurisdiction = 'IN-GST';
DELETE FROM budgets WHERE name = 'Annual Budget 2025';
DELETE FROM allocation_rules WHERE name = 'Revenue Split by Sales';

DELETE FROM exchange_rates WHERE base_currency IN ('USD','EUR','GBP') AND quote_currency = 'INR';

DELETE FROM gst_regimes WHERE gstin = '22AAAAA0000A1Z5';

DELETE FROM cost_centers WHERE name IN ('General','Sales','Operations');
DELETE FROM party_refs WHERE display_name IN ('Default Customer','Default Vendor');

DELETE FROM accounts WHERE code IN ('1000','2000','3000','4000');
