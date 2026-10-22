-- ==============================================================================
-- DATABASE MIGRATION: 004_s04_services_maint.up.sql
-- Bounded Context: S04 Services & Maintenance Department
-- Units: Feeder Maintenance, Transformer Maintenance, Metering
-- Target Database: PostgreSQL 16+ with PostGIS
-- Compliance: Closed-Loop Field Proof, Anti-Photo-Reuse & Haversine GPS Locking
-- ==============================================================================

-- 1. Create Isolated Services Maintenance Schema
CREATE SCHEMA IF NOT EXISTS services_maint;

-- ------------------------------------------------------------------------------
-- Table 1: services_maint.maintenance_work_orders
-- Purpose: Closed-loop field work order dispatch from Unit Head to Team Lead,
--          linking tickets, assets, safety permits, and store requisitions.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS services_maint.maintenance_work_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_no VARCHAR(60) NOT NULL UNIQUE,
    ticket_id UUID REFERENCES customercare.service_tickets(id) ON DELETE SET NULL,
    unit_type VARCHAR(50) NOT NULL CHECK (
        unit_type IN ('FEEDER_MAINT', 'TRANSFORMER_MAINT', 'METERING')
    ),
    work_order_type VARCHAR(50) NOT NULL CHECK (
        work_order_type IN (
            'CORRECTIVE_REPAIR', 'PREVENTIVE_MAINTENANCE', 'EMERGENCY_OUTAGE',
            'METER_INSTALLATION', 'DT_INSPECTION'
        )
    ),
    priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM' CHECK (
        priority IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')
    ),
    unit_head_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    team_lead_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    asset_id UUID NOT NULL REFERENCES foundation.grid_assets(id) ON DELETE RESTRICT,
    ptw_required BOOLEAN NOT NULL DEFAULT FALSE,
    ptw_permit_id UUID,
    status VARCHAR(30) NOT NULL DEFAULT 'DISPATCHED_TO_FIELD' CHECK (
        status IN (
            'DRAFT', 'DISPATCHED_TO_FIELD', 'WORK_IN_PROGRESS',
            'PROOF_SUBMITTED', 'ATTENDED_VERIFIED', 'REJECTED_REWORK', 'CANCELLED'
        )
    ),
    materials_requisition_ref VARCHAR(100),
    instructions TEXT NOT NULL,
    rework_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_work_orders_ticket ON services_maint.maintenance_work_orders(ticket_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_unit ON services_maint.maintenance_work_orders(unit_type);
CREATE INDEX IF NOT EXISTS idx_work_orders_unit_head ON services_maint.maintenance_work_orders(unit_head_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_team_lead ON services_maint.maintenance_work_orders(team_lead_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_asset ON services_maint.maintenance_work_orders(asset_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_status ON services_maint.maintenance_work_orders(status);

-- ------------------------------------------------------------------------------
-- Table 2: services_maint.proof_of_work_logs
-- Purpose: Tamper-evident field proof submitted by Team Lead on Mobile PWA.
--          Captures live GPS, Haversine asset distance, photo hashes, and AI verification.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS services_maint.proof_of_work_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    work_order_id UUID NOT NULL REFERENCES services_maint.maintenance_work_orders(id) ON DELETE CASCADE,
    team_lead_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    captured_gps GEOMETRY(Point, 4326) NOT NULL,
    distance_from_asset_meters NUMERIC(10, 2) NOT NULL,
    before_photo_url TEXT NOT NULL,
    after_photo_url TEXT NOT NULL,
    photo_hash_before VARCHAR(64) NOT NULL,
    photo_hash_after VARCHAR(64) NOT NULL,
    meter_seal_number VARCHAR(60),
    work_details TEXT NOT NULL,
    ollama_gps_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ai_risk_score NUMERIC(4, 2) NOT NULL DEFAULT 0.00,
    ai_reasoning TEXT,
    unit_head_verified BOOLEAN NOT NULL DEFAULT FALSE,
    unit_head_notes TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_proof_work_order ON services_maint.proof_of_work_logs(work_order_id);
CREATE INDEX IF NOT EXISTS idx_proof_team_lead ON services_maint.proof_of_work_logs(team_lead_id);
CREATE INDEX IF NOT EXISTS idx_proof_gps ON services_maint.proof_of_work_logs USING GIST(captured_gps);
CREATE INDEX IF NOT EXISTS idx_proof_photo_before ON services_maint.proof_of_work_logs(photo_hash_before);
CREATE INDEX IF NOT EXISTS idx_proof_photo_after ON services_maint.proof_of_work_logs(photo_hash_after);

