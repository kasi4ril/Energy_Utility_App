-- ==============================================================================
-- DATABASE MIGRATION: 003_s03_customer_care.up.sql
-- Bounded Context: S03 Customer Care, Anti-Ghost GPS Lock & Mandatory Closure Gate
-- Target Database: PostgreSQL 16+ with PostGIS
-- Compliance: NERC Customer Protection Standards, Direct Routing & Mandatory Call Confirmation
-- ==============================================================================

-- 1. Create Isolated Customer Care Schema
CREATE SCHEMA IF NOT EXISTS customercare;

-- ------------------------------------------------------------------------------
-- Table 1: customercare.customers
-- Purpose: KYC onboarded customer directory with mandatory premises GPS, assigned DT,
--          ledger balance, and cryptographic row HMAC signature.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS customercare.customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_no VARCHAR(50) NOT NULL UNIQUE,
    cin VARCHAR(50) NOT NULL UNIQUE,
    full_name VARCHAR(150) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    nin VARCHAR(30),
    building_class VARCHAR(50) NOT NULL DEFAULT 'RESIDENTIAL' CHECK (
        building_class IN ('RESIDENTIAL', 'COMMERCIAL', 'INDUSTRIAL', 'SPECIAL')
    ),
    premises_gps GEOMETRY(Point, 4326) NOT NULL,
    transformer_id UUID NOT NULL REFERENCES foundation.grid_assets(id) ON DELETE RESTRICT,
    meter_serial VARCHAR(60) UNIQUE,
    tariff_band VARCHAR(5) NOT NULL CHECK (tariff_band IN ('A', 'B', 'C', 'D', 'E')),
    ledger_balance NUMERIC(14, 2) NOT NULL DEFAULT 0.00,
    kyc_status VARCHAR(30) NOT NULL DEFAULT 'VERIFIED' CHECK (
        kyc_status IN ('PENDING', 'VERIFIED', 'SUSPENDED', 'DEBTOR_FLAGGED')
    ),
    row_hmac_signature VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_customers_phone ON customercare.customers(phone);
CREATE INDEX IF NOT EXISTS idx_customers_transformer ON customercare.customers(transformer_id);
CREATE INDEX IF NOT EXISTS idx_customers_gps ON customercare.customers USING GIST(premises_gps);
CREATE INDEX IF NOT EXISTS idx_customers_meter ON customercare.customers(meter_serial);

-- ------------------------------------------------------------------------------
-- Table 2: customercare.service_tickets
-- Purpose: Omnichannel ticketing ledger enforcing direct-to-unit-head routing and
--          the exclusive Customer Care confirmation closure gate.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS customercare.service_tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_no VARCHAR(50) NOT NULL UNIQUE,
    customer_id UUID REFERENCES customercare.customers(id) ON DELETE RESTRICT,
    channel VARCHAR(30) NOT NULL CHECK (
        channel IN ('WALK_IN', 'ONLINE_PORTAL', 'CALL_CENTER')
    ),
    category VARCHAR(50) NOT NULL CHECK (
        category IN (
            'FALLEN_POLE', 'LINE_SNAP', 'TRANSFORMER_FAULT', 'FUSE_BLOWN',
            'OIL_LEAK', 'FAULTY_METER', 'NEW_METER_SCHEDULE', 'WRONG_BILLING',
            'UNCREDITED_PAYMENT', 'GENERAL_INQUIRY'
        )
    ),
    description TEXT NOT NULL,
    assigned_unit_id UUID NOT NULL REFERENCES foundation.org_units(id) ON DELETE RESTRICT,
    unit_head_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    dept_head_supervisor_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    status VARCHAR(30) NOT NULL DEFAULT 'LOGGED' CHECK (
        status IN ('LOGGED', 'ROUTED', 'ATTENDED_TO', 'CLOSED', 'REOPENED_ESCALATED')
    ),
    attended_proof_note TEXT,
    attended_by_staff_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    customer_call_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by_cc_officer_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    customer_feedback_notes TEXT,
    sla_due_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tickets_customer ON customercare.service_tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON customercare.service_tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_assigned_unit ON customercare.service_tickets(assigned_unit_id);
CREATE INDEX IF NOT EXISTS idx_tickets_unit_head ON customercare.service_tickets(unit_head_id);
CREATE INDEX IF NOT EXISTS idx_tickets_dept_head ON customercare.service_tickets(dept_head_supervisor_id);
CREATE INDEX IF NOT EXISTS idx_tickets_sla ON customercare.service_tickets(sla_due_at);

-- ------------------------------------------------------------------------------
-- Table 3: customercare.meter_schedules
-- Purpose: Installation, repair, and recertification schedule ledger for Metering Unit.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS customercare.meter_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL REFERENCES customercare.customers(id) ON DELETE RESTRICT,
    ticket_id UUID REFERENCES customercare.service_tickets(id) ON DELETE SET NULL,
    schedule_type VARCHAR(50) NOT NULL CHECK (
        schedule_type IN ('NEW_INSTALLATION', 'REPAIR', 'CALIBRATION', 'RECERTIFICATION')
    ),
    scheduled_date DATE NOT NULL,
    assigned_metering_team_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE SET NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'SCHEDULED' CHECK (
        status IN ('SCHEDULED', 'DISPATCHED', 'COMPLETED', 'RESCHEDULED', 'CANCELLED')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_meter_schedules_customer ON customercare.meter_schedules(customer_id);
CREATE INDEX IF NOT EXISTS idx_meter_schedules_date ON customercare.meter_schedules(scheduled_date);
CREATE INDEX IF NOT EXISTS idx_meter_schedules_status ON customercare.meter_schedules(status);

