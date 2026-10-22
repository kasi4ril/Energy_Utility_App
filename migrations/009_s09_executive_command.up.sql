-- Migration: 009_s09_executive_command.up.sql
-- Description: Establishes schema executive_command for C-Suite aggregated intelligence,
--              predictive maintenance triggers (Hydro, Generator, Transformer, Feeder),
--              and unified executive action console queue.

CREATE SCHEMA IF NOT EXISTS executive_command;

-- 1. Executive Intelligence Snapshots Table
CREATE TABLE IF NOT EXISTS executive_command.executive_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snapshot_date DATE NOT NULL,
    shift_window VARCHAR(32) NOT NULL,
    total_generation_mw NUMERIC(10,2) NOT NULL DEFAULT 0.00 CHECK (total_generation_mw >= 0.00),
    tcn_grid_intake_mw NUMERIC(10,2) NOT NULL DEFAULT 0.00 CHECK (tcn_grid_intake_mw >= 0.00),
    total_system_load_mw NUMERIC(10,2) NOT NULL DEFAULT 0.00 CHECK (total_system_load_mw >= 0.00),
    feeders_in_service_count INT NOT NULL DEFAULT 0 CHECK (feeders_in_service_count >= 0),
    feeders_out_fault_count INT NOT NULL DEFAULT 0 CHECK (feeders_out_fault_count >= 0),
    feeders_out_ptw_count INT NOT NULL DEFAULT 0 CHECK (feeders_out_ptw_count >= 0),
    feeders_out_rationing_count INT NOT NULL DEFAULT 0 CHECK (feeders_out_rationing_count >= 0),
    total_revenue_collected NUMERIC(14,2) NOT NULL DEFAULT 0.00 CHECK (total_revenue_collected >= 0.00),
    prepaid_revenue_share_pct NUMERIC(5,2) NOT NULL DEFAULT 0.00 CHECK (prepaid_revenue_share_pct >= 0.00 AND prepaid_revenue_share_pct <= 100.00),
    unresolved_theft_alerts_count INT NOT NULL DEFAULT 0 CHECK (unresolved_theft_alerts_count >= 0),
    unresolved_churn_alerts_count INT NOT NULL DEFAULT 0 CHECK (unresolved_churn_alerts_count >= 0),
    audit_pending_procurements_count INT NOT NULL DEFAULT 0 CHECK (audit_pending_procurements_count >= 0),
    unreturned_scrap_count INT NOT NULL DEFAULT 0 CHECK (unreturned_scrap_count >= 0),
    notes TEXT,
    captured_by_staff_id VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_exec_snapshot_date ON executive_command.executive_snapshots(snapshot_date);
CREATE INDEX IF NOT EXISTS idx_exec_snapshot_shift ON executive_command.executive_snapshots(shift_window);
CREATE INDEX IF NOT EXISTS idx_exec_snapshot_created_at ON executive_command.executive_snapshots(created_at);

-- 2. Predictive Maintenance Triggers Table
CREATE TABLE IF NOT EXISTS executive_command.predictive_maintenance_triggers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trigger_no VARCHAR(64) UNIQUE NOT NULL,
    asset_id VARCHAR(64) NOT NULL,
    asset_type VARCHAR(32) NOT NULL,
    model_name VARCHAR(64) NOT NULL,
    risk_score NUMERIC(4,3) NOT NULL CHECK (risk_score >= 0.000 AND risk_score <= 1.000),
    risk_level VARCHAR(16) NOT NULL DEFAULT 'LOW',
    telemetry_summary_json JSONB,
    recommended_action TEXT NOT NULL,
    auto_work_order_created BOOLEAN NOT NULL DEFAULT FALSE,
    linked_work_order_id VARCHAR(64),
    acknowledged_by_staff_id VARCHAR(64),
    acknowledged_at TIMESTAMPTZ,
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pmt_no ON executive_command.predictive_maintenance_triggers(trigger_no);
CREATE INDEX IF NOT EXISTS idx_pmt_asset ON executive_command.predictive_maintenance_triggers(asset_id);
CREATE INDEX IF NOT EXISTS idx_pmt_model ON executive_command.predictive_maintenance_triggers(model_name);
CREATE INDEX IF NOT EXISTS idx_pmt_risk_level ON executive_command.predictive_maintenance_triggers(risk_level);
CREATE INDEX IF NOT EXISTS idx_pmt_status ON executive_command.predictive_maintenance_triggers(status);

-- 3. Executive Action Console Items Table
CREATE TABLE IF NOT EXISTS executive_command.executive_action_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action_type VARCHAR(32) NOT NULL,
    target_role VARCHAR(32) NOT NULL,
    source_entity_id VARCHAR(64) NOT NULL,
    source_entity_ref VARCHAR(64) NOT NULL,
    source_department VARCHAR(64) NOT NULL,
    summary TEXT NOT NULL,
    amount NUMERIC(14,2) DEFAULT 0.00,
    priority VARCHAR(16) NOT NULL DEFAULT 'ROUTINE',
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    decided_by_staff_id VARCHAR(64),
    decided_at TIMESTAMPTZ,
    decision_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_eai_target_role ON executive_command.executive_action_items(target_role);
CREATE INDEX IF NOT EXISTS idx_eai_action_type ON executive_command.executive_action_items(action_type);
CREATE INDEX IF NOT EXISTS idx_eai_status ON executive_command.executive_action_items(status);
CREATE INDEX IF NOT EXISTS idx_eai_priority ON executive_command.executive_action_items(priority);
CREATE INDEX IF NOT EXISTS idx_eai_source ON executive_command.executive_action_items(source_entity_id);

