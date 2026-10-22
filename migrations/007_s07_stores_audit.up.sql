-- Migration: 007_s07_stores_audit.up.sql
-- Description: Establishes schema stores_audit for store stock ledgers, material requisitions,
--              reverse logistics scrap returns, and staff allowance claims (3-way verification).

CREATE SCHEMA IF NOT EXISTS stores_audit;

-- 1. Store Inventory Ledger Table
CREATE TABLE IF NOT EXISTS stores_audit.store_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_code VARCHAR(64) UNIQUE NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    category VARCHAR(64) NOT NULL,
    quantity_available NUMERIC(12,2) NOT NULL DEFAULT 0.00 CHECK (quantity_available >= 0.00),
    reorder_threshold NUMERIC(12,2) NOT NULL DEFAULT 10.00 CHECK (reorder_threshold >= 0.00),
    unit_of_measure VARCHAR(32) NOT NULL DEFAULT 'PCS',
    unit_cost NUMERIC(14,2) NOT NULL DEFAULT 0.00 CHECK (unit_cost >= 0.00),
    bin_location VARCHAR(64) NOT NULL,
    last_updated_by_store_officer_id VARCHAR(64),
    row_hmac_signature VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_store_inventory_code ON stores_audit.store_inventory(item_code);
CREATE INDEX IF NOT EXISTS idx_store_inventory_category ON stores_audit.store_inventory(category);
CREATE INDEX IF NOT EXISTS idx_store_inventory_reorder ON stores_audit.store_inventory(quantity_available, reorder_threshold);

-- 2. Material Requisitions Table (Dual-Path & MD Procurement)
CREATE TABLE IF NOT EXISTS stores_audit.material_requisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    req_no VARCHAR(64) UNIQUE NOT NULL,
    origin_dept_id VARCHAR(64) NOT NULL,
    linked_ticket_or_job_id VARCHAR(64),
    items_json JSONB NOT NULL,
    total_cost NUMERIC(14,2) NOT NULL DEFAULT 0.00 CHECK (total_cost >= 0.00),
    in_stock_flag BOOLEAN NOT NULL DEFAULT TRUE,
    unit_head_id VARCHAR(64) NOT NULL,
    dept_head_id VARCHAR(64),
    dept_or_gm_approved BOOLEAN NOT NULL DEFAULT FALSE,
    dept_approved_at TIMESTAMPTZ,
    audit_officer_id VARCHAR(64),
    audit_consent_status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    audit_consent_at TIMESTAMPTZ,
    audit_notes TEXT,
    md_procurement_approved BOOLEAN NOT NULL DEFAULT FALSE,
    md_approved_at TIMESTAMPTZ,
    md_notes TEXT,
    store_officer_id VARCHAR(64),
    issued_at TIMESTAMPTZ,
    waybill_no VARCHAR(64),
    scrap_returned_flag BOOLEAN NOT NULL DEFAULT FALSE,
    scrap_return_id VARCHAR(64),
    status VARCHAR(32) NOT NULL DEFAULT 'SUBMITTED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mat_req_no ON stores_audit.material_requisitions(req_no);
CREATE INDEX IF NOT EXISTS idx_mat_req_status ON stores_audit.material_requisitions(status);
CREATE INDEX IF NOT EXISTS idx_mat_req_unit_head ON stores_audit.material_requisitions(unit_head_id);
CREATE INDEX IF NOT EXISTS idx_mat_req_audit ON stores_audit.material_requisitions(audit_consent_status);
CREATE INDEX IF NOT EXISTS idx_mat_req_instock ON stores_audit.material_requisitions(in_stock_flag);

-- 3. Reverse Logistics Scrap Returns Table
CREATE TABLE IF NOT EXISTS stores_audit.scrap_returns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_no VARCHAR(64) UNIQUE NOT NULL,
    requisition_id UUID NOT NULL REFERENCES stores_audit.material_requisitions(id) ON DELETE CASCADE,
    team_lead_id VARCHAR(64) NOT NULL,
    store_officer_id VARCHAR(64) NOT NULL,
    items_returned_json JSONB NOT NULL,
    scrap_condition_notes TEXT NOT NULL,
    returned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_scrap_return_no ON stores_audit.scrap_returns(return_no);
CREATE INDEX IF NOT EXISTS idx_scrap_req_id ON stores_audit.scrap_returns(requisition_id);
CREATE INDEX IF NOT EXISTS idx_scrap_team_lead ON stores_audit.scrap_returns(team_lead_id);

-- 4. Staff Allowance & Financial Claims Table
CREATE TABLE IF NOT EXISTS stores_audit.staff_allowance_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_no VARCHAR(64) UNIQUE NOT NULL,
    staff_id VARCHAR(64) NOT NULL,
    claim_type VARCHAR(32) NOT NULL,
    amount NUMERIC(14,2) NOT NULL CHECK (amount > 0.00),
    claim_date DATE NOT NULL,
    shift_window VARCHAR(32),
    linked_ptw_or_ticket_id VARCHAR(64),
    roster_match_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ptw_crew_verified BOOLEAN NOT NULL DEFAULT FALSE,
    gps_match_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ollama_audit_flag VARCHAR(32) NOT NULL DEFAULT 'VERIFIED',
    ollama_fraud_score NUMERIC(4,3) NOT NULL DEFAULT 0.000 CHECK (ollama_fraud_score >= 0.000 AND ollama_fraud_score <= 1.000),
    auditor_id VARCHAR(64),
    auditor_notes TEXT,
    audit_decision_at TIMESTAMPTZ,
    status VARCHAR(32) NOT NULL DEFAULT 'SUBMITTED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_staff_claims_no ON stores_audit.staff_allowance_claims(claim_no);
CREATE INDEX IF NOT EXISTS idx_staff_claims_staff_id ON stores_audit.staff_allowance_claims(staff_id);
CREATE INDEX IF NOT EXISTS idx_staff_claims_status ON stores_audit.staff_allowance_claims(status);
CREATE INDEX IF NOT EXISTS idx_staff_claims_flag ON stores_audit.staff_allowance_claims(ollama_audit_flag);
CREATE INDEX IF NOT EXISTS idx_staff_claims_date ON stores_audit.staff_allowance_claims(claim_date);

