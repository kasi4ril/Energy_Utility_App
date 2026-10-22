-- ==============================================================================
-- DATABASE MIGRATION: 005_s05_generation_master.up.sql
-- Bounded Context: S05 Generation Department & Master Station Hub
-- Units: Electrical Wing, Mechanical Wing, Workshop Team, 6 Generating Plants
-- Target Database: PostgreSQL 16+
-- Compliance: Call Collation Telemetry, Google Forms HMAC Bridge & GM Gate
-- ==============================================================================

-- 1. Create Isolated Generation Master Schema
CREATE SCHEMA IF NOT EXISTS generation_master;

-- ------------------------------------------------------------------------------
-- Table 1: generation_master.master_station_telemetry_logs
-- Purpose: Hourly & shift hydrological/generator telemetry collated from Plants 1–6
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS generation_master.master_station_telemetry_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shift_id VARCHAR(50) NOT NULL,
    shift_window VARCHAR(30) NOT NULL CHECK (
        shift_window IN ('MORNING', 'AFTERNOON', 'NIGHT')
    ),
    plant_id UUID NOT NULL REFERENCES foundation.grid_assets(id) ON DELETE RESTRICT,
    plant_state VARCHAR(30) NOT NULL CHECK (
        plant_state IN ('IN_SERVICE', 'STANDBY', 'OUT_OF_SERVICE')
    ),
    outage_reason VARCHAR(50) CHECK (
        outage_reason IN (
            'HYDRO_LOW_WATER', 'TURBINE_MECH_FAULT', 'GENERATOR_ELEC_FAULT',
            'PROTECTION_RELAY_TRIP', 'PLANNED_MAINTENANCE'
        )
    ),
    outage_details TEXT,
    water_level_m NUMERIC(8, 2) NOT NULL,
    headwater_elevation_m NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    tailwater_elevation_m NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    inflow_rate_m3_s NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    mw_generated NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    mvar_reactive NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    stator_winding_temp_c NUMERIC(6, 2) NOT NULL DEFAULT 0.00,
    bearing_vibration_mm_s NUMERIC(6, 2) NOT NULL DEFAULT 0.00,
    lube_oil_pressure_bar NUMERIC(6, 2) NOT NULL DEFAULT 0.00,
    logged_by_operator_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_telemetry_shift ON generation_master.master_station_telemetry_logs(shift_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_plant ON generation_master.master_station_telemetry_logs(plant_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_recorded_at ON generation_master.master_station_telemetry_logs(recorded_at);
CREATE INDEX IF NOT EXISTS idx_telemetry_state ON generation_master.master_station_telemetry_logs(plant_state);

-- ------------------------------------------------------------------------------
-- Table 2: generation_master.remote_plant_shift_rollcall
-- Purpose: Remote physical roll-call of operators, engineers, and techs across Plants 1–6
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS generation_master.remote_plant_shift_rollcall (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shift_id VARCHAR(50) NOT NULL,
    shift_window VARCHAR(30) NOT NULL CHECK (
        shift_window IN ('MORNING', 'AFTERNOON', 'NIGHT')
    ),
    plant_id UUID NOT NULL REFERENCES foundation.grid_assets(id) ON DELETE RESTRICT,
    operator_names TEXT[] NOT NULL DEFAULT '{}',
    protection_engineers TEXT[] NOT NULL DEFAULT '{}',
    technicians TEXT[] NOT NULL DEFAULT '{}',
    handover_notes TEXT NOT NULL,
    synced_to_hq_hr BOOLEAN NOT NULL DEFAULT FALSE,
    master_operator_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_rollcall_shift ON generation_master.remote_plant_shift_rollcall(shift_id);
CREATE INDEX IF NOT EXISTS idx_rollcall_plant ON generation_master.remote_plant_shift_rollcall(plant_id);

-- ------------------------------------------------------------------------------
-- Table 3: generation_master.plant_fault_and_workshop_jobs
-- Purpose: Plant fault reporting, GM material approval gate, and Workshop fabrication tracking
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS generation_master.plant_fault_and_workshop_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_no VARCHAR(60) NOT NULL UNIQUE,
    plant_id UUID NOT NULL REFERENCES foundation.grid_assets(id) ON DELETE RESTRICT,
    reported_by_operator_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    wing VARCHAR(30) NOT NULL CHECK (
        wing IN ('ELECTRICAL', 'MECHANICAL')
    ),
    fault_description TEXT NOT NULL,
    severity VARCHAR(20) NOT NULL DEFAULT 'MEDIUM' CHECK (
        severity IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')
    ),
    unit_head_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    gm_approved BOOLEAN NOT NULL DEFAULT FALSE,
    gm_approval_notes TEXT,
    requires_workshop BOOLEAN NOT NULL DEFAULT FALSE,
    workshop_status VARCHAR(30) NOT NULL DEFAULT 'NOT_APPLICABLE' CHECK (
        workshop_status IN (
            'NOT_APPLICABLE', 'PENDING_DISPATCH', 'MACHINING_IN_PROGRESS',
            'WINDING_FABRICATION', 'QUALITY_CHECK', 'COMPLETED_DELIVERED', 'REJECTED'
        )
    ),
    raw_materials_summary TEXT,
    store_requisition_ref VARCHAR(100),
    status VARCHAR(30) NOT NULL DEFAULT 'REPORTED' CHECK (
        status IN (
            'REPORTED', 'UNIT_HEAD_REVIEWED', 'GM_APPROVED',
            'WORKSHOP_IN_PROGRESS', 'RESOLVED_REINSTALLED', 'CANCELLED'
        )
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_plant_jobs_plant ON generation_master.plant_fault_and_workshop_jobs(plant_id);
CREATE INDEX IF NOT EXISTS idx_plant_jobs_status ON generation_master.plant_fault_and_workshop_jobs(status);
CREATE INDEX IF NOT EXISTS idx_plant_jobs_workshop ON generation_master.plant_fault_and_workshop_jobs(workshop_status);

