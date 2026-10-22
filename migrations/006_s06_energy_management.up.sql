-- ==============================================================================
-- DATABASE MIGRATION: 006_s06_energy_management.up.sql
-- Bounded Context: S06 Energy Management, Priority Dispatch & 5-Step Safety PTW
-- Target Database: PostgreSQL 16+
-- Compliance: 5-Step Electrical Substation Isolation, Priority Feeder Rationing
-- ==============================================================================

-- 1. Create Isolated Energy Management Schema
CREATE SCHEMA IF NOT EXISTS energy_management;

-- ------------------------------------------------------------------------------
-- Table 1: energy_management.dispatch_telemetry_logs
-- Purpose: Bulk MW intake collation (6 Generating Plants + TCN) & Substation rollcall
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS energy_management.dispatch_telemetry_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shift_id VARCHAR(50) NOT NULL,
    shift_window VARCHAR(30) NOT NULL CHECK (
        shift_window IN ('MORNING', 'AFTERNOON', 'NIGHT')
    ),
    generation_mw_in NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    tcn_mw_in NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    total_bulk_mw NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    substation_staff_rollcall_json TEXT NOT NULL DEFAULT '{}',
    handover_notes TEXT NOT NULL,
    operator_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    logged_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dispatch_shift ON energy_management.dispatch_telemetry_logs(shift_id);
CREATE INDEX IF NOT EXISTS idx_dispatch_logged_at ON energy_management.dispatch_telemetry_logs(logged_at);

-- ------------------------------------------------------------------------------
-- Table 2: energy_management.feeder_status_ledger
-- Purpose: Real-time 33kV/11kV feeder operational state, load, priority & active PTW
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS energy_management.feeder_status_ledger (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    feeder_asset_id UUID NOT NULL REFERENCES foundation.grid_assets(id) ON DELETE RESTRICT UNIQUE,
    feeder_code VARCHAR(60) NOT NULL UNIQUE,
    current_load_mw NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    current_amps NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    state VARCHAR(30) NOT NULL DEFAULT 'IN_SERVICE' CHECK (
        state IN ('IN_SERVICE', 'OUT_ON_FAULT', 'OUT_ON_PERMIT_PTW', 'OUT_ON_LOAD_SHEDDING')
    ),
    outage_reason_code VARCHAR(60),
    outage_narrative TEXT,
    priority_rank INT NOT NULL DEFAULT 3,
    tariff_band VARCHAR(10) NOT NULL DEFAULT 'BAND_B',
    active_ptw_id UUID,
    scheduled_restoration_time TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_feeder_state ON energy_management.feeder_status_ledger(state);
CREATE INDEX IF NOT EXISTS idx_feeder_priority ON energy_management.feeder_status_ledger(priority_rank);

-- ------------------------------------------------------------------------------
-- Table 3: energy_management.safety_permits_ptw
-- Purpose: Mandatory 5-Step Safety Permit-to-Work (PTW) Substation Isolation Lock
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS energy_management.safety_permits_ptw (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    permit_no VARCHAR(60) NOT NULL UNIQUE,
    work_order_id UUID,
    asset_id UUID NOT NULL REFERENCES foundation.grid_assets(id) ON DELETE RESTRICT,
    unit_head_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    dept_head_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    dept_head_approved BOOLEAN NOT NULL DEFAULT FALSE,
    dept_head_approved_at TIMESTAMPTZ,
    safety_rep_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    safety_rep_approved BOOLEAN NOT NULL DEFAULT FALSE,
    safety_rep_approved_at TIMESTAMPTZ,
    substation_operator_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    isolation_granted_at TIMESTAMPTZ,
    isolation_notes TEXT,
    safety_clearance_signed_at TIMESTAMPTZ,
    safety_clearance_notes TEXT,
    restored_at TIMESTAMPTZ,
    status VARCHAR(40) NOT NULL DEFAULT 'REQUESTED' CHECK (
        status IN (
            'REQUESTED', 'APPROVED_PENDING_ISOLATION', 'ISOLATED_EARTHED',
            'SAFETY_CLEARED', 'RESTORED_CLOSED', 'REJECTED_CANCELLED'
        )
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ptw_asset ON energy_management.safety_permits_ptw(asset_id);
CREATE INDEX IF NOT EXISTS idx_ptw_status ON energy_management.safety_permits_ptw(status);
CREATE INDEX IF NOT EXISTS idx_ptw_work_order ON energy_management.safety_permits_ptw(work_order_id);

