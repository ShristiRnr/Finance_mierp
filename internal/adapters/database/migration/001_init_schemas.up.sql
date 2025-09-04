-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ================================
-- Request Metadata
-- ================================
CREATE TABLE request_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id TEXT NOT NULL,             -- idempotency key
    organization_id TEXT,                 -- multi-tenant hints
    tenant_id TEXT,
    auth_subject TEXT,                    -- user/service principal
    source_system TEXT,                   -- CRM/VMS/HRMS etc.
    trace_id TEXT,                        -- distributed tracing id
    created_at TIMESTAMPTZ DEFAULT now()
);

-- ================================
-- Audit Fields
-- ================================
CREATE TABLE audit_fields (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision TEXT                         -- optimistic concurrency token
);

-- =====================================
-- Chart of Accounts (must exist early for FKs)
-- =====================================
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    type TEXT NOT NULL,                        -- AccountType (ASSET, LIABILITY, etc.)
    parent_id UUID REFERENCES accounts(id),    -- hierarchy
    status TEXT NOT NULL DEFAULT 'ACTIVE',     -- AccountStatus
    allow_manual_journal BOOLEAN DEFAULT true,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);

CREATE INDEX idx_accounts_code ON accounts(code);

-- =====================================================
-- Invoices
-- =====================================================
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,                       -- InvoiceType
    invoice_date TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ,
    delivery_date TIMESTAMPTZ,
    organization_id TEXT NOT NULL,
    po_number TEXT,
    eway_number_legacy TEXT,
    status_note TEXT,
    status TEXT NOT NULL,                     -- InvoiceStatus

    -- Payments/logistics references
    payment_reference TEXT,
    challan_number TEXT,
    challan_date TIMESTAMPTZ,
    lr_number TEXT,
    transporter_name TEXT,
    transporter_id TEXT,
    vehicle_number TEXT,
    against_invoice_number TEXT,
    against_invoice_date TIMESTAMPTZ,

    -- Totals
    subtotal NUMERIC(18,2) NOT NULL,
    grand_total NUMERIC(18,2) NOT NULL,

    -- GST breakup
    gst_rate NUMERIC(6,3) DEFAULT 0,
    gst_cgst NUMERIC(18,2) DEFAULT 0,
    gst_sgst NUMERIC(18,2) DEFAULT 0,
    gst_igst NUMERIC(18,2) DEFAULT 0,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Invoice Items
-- =====================================================
CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    hsn TEXT,                           -- may contain leading zeros
    quantity INT NOT NULL,
    unit_price NUMERIC(18,2) NOT NULL,
    line_subtotal NUMERIC(18,2) NOT NULL,
    line_total NUMERIC(18,2) NOT NULL,
    cost_center_id TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision INT DEFAULT 1
);

CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id); 

-- =====================================================
-- Item-level Discounts and Taxes (Normalized)
-- =====================================================
CREATE TABLE invoice_item_discounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES invoice_items(id) ON DELETE CASCADE,
    description TEXT,
    amount NUMERIC(18,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_invoice_item_discounts_item ON invoice_item_discounts(item_id);

CREATE TABLE invoice_item_taxes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES invoice_items(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    rate NUMERIC(6,3) NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_invoice_item_taxes_item ON invoice_item_taxes(item_id);

-- =====================================================
-- Invoice-level Discounts and Taxes
-- =====================================================
CREATE TABLE invoice_discounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    description TEXT,
    amount NUMERIC(18,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_invoice_discounts_invoice ON invoice_discounts(invoice_id);

CREATE TABLE invoice_taxes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    rate NUMERIC(6,3) NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_invoice_taxes_invoice ON invoice_taxes(invoice_id);

-- =====================================================
-- GST Breakup
-- =====================================================
CREATE TABLE gst_breakups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    taxable_amount NUMERIC(18,2) NOT NULL,
    cgst NUMERIC(18,2) DEFAULT 0,
    sgst NUMERIC(18,2) DEFAULT 0,
    igst NUMERIC(18,2) DEFAULT 0,
    total_gst NUMERIC(18,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- GST Regime
-- =====================================================
CREATE TABLE gst_regimes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    gstin TEXT NOT NULL,
    place_of_supply TEXT NOT NULL,
    reverse_charge BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- GST Document Status
-- =====================================================
CREATE TABLE gst_doc_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    einvoice_status TEXT DEFAULT 'EINV_STATUS_UNSPECIFIED',
    irn TEXT,
    ack_no TEXT,
    ack_date TIMESTAMPTZ,
    eway_status TEXT DEFAULT 'EWAY_STATUS_UNSPECIFIED',
    eway_bill_no TEXT,
    eway_valid_upto TIMESTAMPTZ,
    last_error TEXT,
    last_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    revision INT DEFAULT 1
);


-- =====================================================
-- Credit / Debit Notes
-- =====================================================
CREATE TABLE credit_debit_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    type TEXT NOT NULL,                     -- NoteType: CREDIT / DEBIT
    amount NUMERIC(18,2) NOT NULL,
    reason TEXT,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_credit_debit_notes_invoice_id ON credit_debit_notes(invoice_id);
-- =====================================================
-- Payment Dues
-- =====================================================
CREATE TABLE payment_dues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    amount_due NUMERIC(18,2) NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL, -- PaymentStatus

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_payment_dues_invoice ON payment_dues(invoice_id);

-- =====================================================
-- Bank Accounts
-- =====================================================
CREATE TABLE bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    account_number TEXT NOT NULL,
    ifsc_or_swift TEXT NOT NULL,
    ledger_account_id TEXT,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_bank_accounts_account_number ON bank_accounts(account_number);


-- =====================================================
-- Bank Transactions
-- =====================================================
CREATE TABLE bank_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bank_account_id UUID NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
    amount NUMERIC(18,2) NOT NULL,
    transaction_date TIMESTAMPTZ NOT NULL,
    description TEXT,
    reference TEXT,
    reconciled BOOLEAN DEFAULT false,
    matched_reference_type TEXT, -- "INVOICE","PAYMENT","JOURNAL"
    matched_reference_id UUID,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_bank_transactions_account_date ON bank_transactions(bank_account_id, transaction_date);

-- =====================================
-- Cost Centers
-- =====================================
CREATE TABLE cost_centers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Journal Entries (header)
-- =====================================================
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_date TIMESTAMPTZ NOT NULL,
    reference TEXT,
    memo TEXT,
    source_type TEXT,     -- "INVOICE"/"PAYMENT"/"ADJUSTMENT"
    source_id TEXT,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_journal_entries_journal_date ON journal_entries(journal_date);

-- Journal Lines (must balance DR = CR)
CREATE TABLE journal_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id),
    side TEXT NOT NULL,     -- "DR" / "CR"
    amount NUMERIC(18,2) NOT NULL,
    cost_center_id TEXT,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- =====================================================
-- Ledger Entries (derived)
-- =====================================================
CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    description TEXT,
    side TEXT NOT NULL,     -- "DR"/"CR"
    amount NUMERIC(18,2) NOT NULL,
    transaction_date TIMESTAMPTZ NOT NULL,
    cost_center_id TEXT,
    reference_type TEXT,    -- "INVOICE"/"PAYMENT"/"JOURNAL"
    reference_id TEXT,

    -- Audit
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_ledger_entries_account_date ON ledger_entries(account_id, transaction_date);

-- =====================================================
-- Budgets
-- =====================================================
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    total_amount NUMERIC(18,2) NOT NULL,
    status TEXT NOT NULL DEFAULT 'DRAFT',  -- (DRAFT/ACTIVE/CLOSED)

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Budget Allocations
-- =====================================================
CREATE TABLE budget_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    department_id TEXT NOT NULL,
    allocated_amount NUMERIC(18,2) NOT NULL,
    spent_amount NUMERIC(18,2) DEFAULT 0,
    remaining_amount NUMERIC(18,2) GENERATED ALWAYS AS (allocated_amount - spent_amount) STORED,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Expenses
-- =====================================================
CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category TEXT NOT NULL,                       -- "LABOR","MATERIAL","OPEX","CAPEX"
    amount NUMERIC(18,2) NOT NULL,
    expense_date TIMESTAMPTZ NOT NULL,
    cost_center_id UUID REFERENCES cost_centers(id),

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Cost Allocations
-- =====================================================
CREATE TABLE cost_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cost_center_id UUID NOT NULL REFERENCES cost_centers(id),
    amount NUMERIC(18,2) NOT NULL,
    reference_type TEXT NOT NULL, -- INVOICE / EXPENSE / JOURNAL
    reference_id TEXT NOT NULL,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Audit Events
-- =====================================================
CREATE TABLE audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
    details TEXT,
    resource_type TEXT,
    resource_id TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- indexes for filtering & search
CREATE INDEX idx_audit_events_user_id ON audit_events(user_id);
CREATE INDEX idx_audit_events_action ON audit_events(action);
CREATE INDEX idx_audit_events_resource ON audit_events(resource_type, resource_id);
CREATE INDEX idx_audit_events_timestamp ON audit_events(timestamp);

-- =====================================================
-- Accruals
-- =====================================================
CREATE TABLE accruals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT,
    amount NUMERIC(18,2) NOT NULL,
    accrual_date TIMESTAMPTZ NOT NULL,
    account_id TEXT NOT NULL,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Allocation Rules
-- =====================================================
CREATE TABLE allocation_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    basis TEXT NOT NULL, -- "ACTIVITY","HEADCOUNT","REVENUE","MACHINE_HOURS"
    source_account_id TEXT NOT NULL,
    target_cost_center_ids TEXT[] NOT NULL, -- array of cost centers
    formula TEXT,

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT,
    revision INT DEFAULT 1
);
CREATE INDEX idx_allocation_rules_source_account ON allocation_rules(source_account_id);

-- =====================================================
-- Profit & Loss Reports
-- =====================================================
CREATE TABLE profit_loss_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    total_revenue NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_expenses NUMERIC(18,2) NOT NULL DEFAULT 0,
    net_profit NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Balance Sheet Reports
-- =====================================================
CREATE TABLE balance_sheet_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    total_assets NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_liabilities NUMERIC(18,2) NOT NULL DEFAULT 0,
    net_worth NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Trial Balance Reports
-- =====================================================
CREATE TABLE trial_balance_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- Ledger Entries (Trial Balance child table)
CREATE TABLE trial_balance_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES trial_balance_reports(id) ON DELETE CASCADE,
    ledger_account TEXT NOT NULL,
    debit NUMERIC(18,2) NOT NULL DEFAULT 0,
    credit NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT
);    

-- =====================================================
-- Compliance Reports
-- =====================================================
CREATE TABLE compliance_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    jurisdiction TEXT NOT NULL,   -- "IN-GST", "US-GAAP", "IFRS"
    details TEXT NOT NULL,        -- JSON or text blob
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Consolidated Reports
-- =====================================================
CREATE TABLE consolidations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_ids TEXT[] NOT NULL,            -- list of entity IDs consolidated
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    report TEXT NOT NULL,                  -- consolidated report (JSON/text)
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
-- Exchange Rates
-- =====================================================
CREATE TABLE exchange_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_currency TEXT NOT NULL,   -- e.g., "USD"
    quote_currency TEXT NOT NULL,  -- e.g., "INR"
    rate NUMERIC(18,6) NOT NULL,   -- quote per base
    as_of TIMESTAMPTZ NOT NULL,    -- rate validity timestamp

    -- Audit fields
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision TEXT,

    UNIQUE(base_currency, quote_currency, as_of)
);

-- =====================================================
-- Cash Flow Forecasts
-- =====================================================
CREATE TABLE cash_flow_forecasts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id TEXT NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    forecast_details TEXT NOT NULL,    -- serialized forecast data (JSON, XML, etc.)
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by TEXT,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    revision INT DEFAULT 1
);

-- =====================================================
--  Events Tables
-- =====================================================
CREATE TABLE finance_invoice_created_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    invoice_number TEXT NOT NULL,
    invoice_date TIMESTAMPTZ NOT NULL,
    total NUMERIC(18,2) NOT NULL,
    organization_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_finance_invoice_created_events_invoice ON finance_invoice_created_events(invoice_id);


-- =====================================================
-- Finance Payment Received Events
-- =====================================================
CREATE TABLE finance_payment_received_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_due_id UUID NOT NULL,
    invoice_id UUID NOT NULL,
    amount_paid NUMERIC(18,2) NOT NULL,
    paid_at TIMESTAMPTZ NOT NULL,
    reference TEXT,
    organization_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_finance_payment_received_events_invoice ON finance_payment_received_events(invoice_id);

-- =====================================================
-- Inventory Cost Posted Events
-- =====================================================
CREATE TABLE inventory_cost_posted_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_type TEXT NOT NULL,       -- "PRODUCTION_ORDER", "GOODS_ISSUE"
    reference_id UUID NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    cost_center_id TEXT,
    organization_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- =====================================================
-- Payroll Posted Events
-- =====================================================
CREATE TABLE payroll_posted_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_run_id UUID NOT NULL,
    total_gross NUMERIC(18,2) NOT NULL,
    total_net NUMERIC(18,2) NOT NULL,
    run_date TIMESTAMPTZ NOT NULL,
    organization_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- =====================================================
-- Vendor Bill Approved Events
-- =====================================================
CREATE TABLE vendor_bill_approved_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_bill_id UUID NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    approved_at TIMESTAMPTZ NOT NULL,
    organization_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- =====================================
-- Helpful additional indexes
-- =====================================
CREATE INDEX idx_invoices_org_party_date ON invoices(organization_id, invoice_date);
CREATE INDEX idx_payment_dues_due_date ON payment_dues(due_date);
CREATE INDEX idx_accruals_account_date ON accruals(account_id, accrual_date);