-- ==============================================================================
-- DATABASE MIGRATION: 001_s01_foundation.up.sql
-- Bounded Context: S01 Foundation, Grid Topology, IAM, Federated HR & WorkNotes
-- Target Database: PostgreSQL 16+ with PostGIS Extension
-- Compliance: NERC / NEMSA Standards, Non-Negotiable SHA-256 Chaining
-- ==============================================================================

-- 1. Enable Required Database Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

-- 2. Create Isolated Foundation Schema
CREATE SCHEMA IF NOT EXISTS foundation;

-- ------------------------------------------------------------------------------
-- Table 1: foundation.org_units
-- Purpose: Defines administrative and operational entity hierarchy:
--          HQ -> District -> Service Unit alongside Master Station -> Plants 1-6.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS foundation.org_units (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID REFERENCES foundation.org_units(id) ON DELETE RESTRICT,
    unit_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    unit_type VARCHAR(50) NOT NULL CHECK (
        unit_type IN ('HQ', 'DISTRICT', 'SERVICE_UNIT', 'MASTER_STATION', 'GENERATING_PLANT')
    ),
    location_gps GEOMETRY(Point, 4326),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_org_units_parent_id ON foundation.org_units(parent_id);
CREATE INDEX IF NOT EXISTS idx_org_units_unit_type ON foundation.org_units(unit_type);
CREATE INDEX IF NOT EXISTS idx_org_units_gps ON foundation.org_units USING GIST(location_gps);

-- ------------------------------------------------------------------------------
-- Table 2: foundation.staff_profiles
-- Purpose: Federated workforce directory across HQ, Districts, Generation & Dispatch.
--          Enforces cryptographic HMAC row signature and 3-Tier RBAC classification.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS foundation.staff_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    staff_number VARCHAR(50) NOT NULL UNIQUE,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL CHECK (
        role IN (
            -- Executive Command Tier (Shared Dashboard)
            'MD', 'HEAD_SERVICES', 'HEAD_INTERNAL_AUDIT', 'HEAD_COMMERCIAL_OPS', 'GM_GENERATION', 'HEAD_ENERGY_MGMT',
            -- Departmental & Unit Heads
            'CUSTOMER_CARE_HEAD', 'UNIT_HEAD_FEEDER_MAINT', 'UNIT_HEAD_TX_MAINT', 'UNIT_HEAD_METERING',
            'UNIT_HEAD_GEN_ELECTRICAL', 'UNIT_HEAD_GEN_MECHANICAL', 'WORKSHOP_LEAD', 'UNIT_HEAD_BILLING',
            'STORE_OFFICER', 'SAFETY_REP_HSE', 'PLANT_HR_OFFICER',
            -- Field & Operational Tier
            'MASTER_STATION_OPERATOR', 'CENTRAL_DISPATCH_OPERATOR', 'SUBSTATION_OPERATOR',
            'TEAM_LEAD_FEEDER', 'TEAM_LEAD_TX', 'TEAM_LEAD_METERING', 'CUSTOMER_CARE_OFFICER', 'INTERNAL_AUDITOR'
        )
    ),
    org_unit_id UUID NOT NULL REFERENCES foundation.org_units(id) ON DELETE RESTRICT,
    certification_status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE' CHECK (
        certification_status IN ('ACTIVE', 'SUSPENDED', 'EXPIRED', 'PROBATION')
    ),
    active_shift_state VARCHAR(50) NOT NULL DEFAULT 'OFF_DUTY' CHECK (
        active_shift_state IN ('ON_DUTY', 'OFF_DUTY', 'STANDBY', 'LEAVE')
    ),
    password_hash VARCHAR(255) NOT NULL,
    row_hmac_signature VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_staff_profiles_role ON foundation.staff_profiles(role);
CREATE INDEX IF NOT EXISTS idx_staff_profiles_org_unit ON foundation.staff_profiles(org_unit_id);
CREATE INDEX IF NOT EXISTS idx_staff_profiles_shift_state ON foundation.staff_profiles(active_shift_state);

-- ------------------------------------------------------------------------------
-- Table 3: foundation.grid_assets
-- Purpose: Physical and spatial grid topology:
--          Plants 1-6 / TCN -> 33kV Substations -> 11kV Feeders -> DTs -> Poles -> Meters.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS foundation.grid_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_asset_id UUID REFERENCES foundation.grid_assets(id) ON DELETE RESTRICT,
    org_unit_id UUID NOT NULL REFERENCES foundation.org_units(id) ON DELETE RESTRICT,
    asset_code VARCHAR(60) NOT NULL UNIQUE,
    asset_type VARCHAR(50) NOT NULL CHECK (
        asset_type IN ('PLANT', 'SUBSTATION_33KV', 'FEEDER_11KV', 'TRANSFORMER_DT', 'SERVICE_POLE', 'CUSTOMER_METER')
    ),
    rated_capacity_kva NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (rated_capacity_kva >= 0),
    tariff_band VARCHAR(5) CHECK (tariff_band IN ('A', 'B', 'C', 'D', 'E', NULL)),
    gps_point GEOMETRY(Point, 4326) NOT NULL,
    operational_state VARCHAR(50) NOT NULL DEFAULT 'IN_SERVICE' CHECK (
        operational_state IN ('IN_SERVICE', 'STANDBY', 'OUT_OF_SERVICE', 'MAINTENANCE', 'DECOMMISSIONED')
    ),
    row_hmac_signature VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_grid_assets_parent ON foundation.grid_assets(parent_asset_id);
CREATE INDEX IF NOT EXISTS idx_grid_assets_org_unit ON foundation.grid_assets(org_unit_id);
CREATE INDEX IF NOT EXISTS idx_grid_assets_type ON foundation.grid_assets(asset_type);
CREATE INDEX IF NOT EXISTS idx_grid_assets_tariff ON foundation.grid_assets(tariff_band);
CREATE INDEX IF NOT EXISTS idx_grid_assets_gps ON foundation.grid_assets USING GIST(gps_point);

-- ------------------------------------------------------------------------------
-- Table 4: foundation.work_notes
-- Purpose: Unified, append-only inter-departmental note stream with cryptographic
--          prev-hash chaining: current_hash = SHA256(prev_hash + payload + ts).
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS foundation.work_notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(60) NOT NULL,
    entity_id VARCHAR(64) NOT NULL,
    author_staff_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    tagged_unit_ids UUID[] DEFAULT '{}',
    note_text TEXT NOT NULL,
    ollama_summary TEXT,
    prev_hash VARCHAR(64) NOT NULL,
    current_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_work_notes_entity ON foundation.work_notes(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_work_notes_author ON foundation.work_notes(author_staff_id);
CREATE INDEX IF NOT EXISTS idx_work_notes_created_at ON foundation.work_notes(created_at);

-- ------------------------------------------------------------------------------
-- Trigger: Immutability Guarantee on foundation.work_notes
-- Core Directive: Strictly reject raw UPDATE or DELETE on notes ledger.
-- ------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION foundation.prevent_work_notes_tampering()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Cryptographic Violation: foundation.work_notes is an immutable ledger. UPDATE and DELETE operations are forbidden.';
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_immutable_work_notes ON foundation.work_notes;
CREATE TRIGGER trg_immutable_work_notes
    BEFORE UPDATE OR DELETE ON foundation.work_notes
    FOR EACH ROW
    EXECUTE FUNCTION foundation.prevent_work_notes_tampering();

