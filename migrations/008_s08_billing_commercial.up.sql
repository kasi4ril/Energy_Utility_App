-- Migration: 008_s08_billing_commercial.up.sql
-- Description: Establishes schema billing_commercial for zero-estimation metered billing,
--              customer meter refund ledgers, multi-platform vending transactions,
--              wrong billing resolution loop, and revenue protection AI alerts.

CREATE SCHEMA IF NOT EXISTS billing_commercial;

-- 1. Vending and Payment Transactions Ledger
CREATE TABLE IF NOT EXISTS billing_commercial.vending_and_payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tx_reference VARCHAR(64) UNIQUE NOT NULL,
    customer_id VARCHAR(64) NOT NULL,
    account_number VARCHAR(64) NOT NULL,
    meter_number VARCHAR(64) NOT NULL,
    feeder_id VARCHAR(64),
    dt_id VARCHAR(64),
    payment_channel_category VARCHAR(32) NOT NULL,
    payment_platform_source VARCHAR(64) NOT NULL,
    amount_paid NUMERIC(14,2) NOT NULL CHECK (amount_paid > 0.00),
    tariff_rate_per_kwh NUMERIC(10,2) NOT NULL DEFAULT 68.00,
    kwh_vended NUMERIC(12,2) NOT NULL CHECK (kwh_vended >= 0.00),
    meter_refund_credited NUMERIC(14,2) NOT NULL DEFAULT 0.00 CHECK (meter_refund_credited >= 0.00),
    sts_token VARCHAR(64),
    settlement_verified BOOLEAN NOT NULL DEFAULT FALSE,
    gateway_payload_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vending_tx_ref ON billing_commercial.vending_and_payment_transactions(tx_reference);
CREATE INDEX IF NOT EXISTS idx_vending_customer ON billing_commercial.vending_and_payment_transactions(customer_id);
CREATE INDEX IF NOT EXISTS idx_vending_meter ON billing_commercial.vending_and_payment_transactions(meter_number);
CREATE INDEX IF NOT EXISTS idx_vending_platform ON billing_commercial.vending_and_payment_transactions(payment_platform_source);
CREATE INDEX IF NOT EXISTS idx_vending_channel ON billing_commercial.vending_and_payment_transactions(payment_channel_category);
CREATE INDEX IF NOT EXISTS idx_vending_dt ON billing_commercial.vending_and_payment_transactions(dt_id);
CREATE INDEX IF NOT EXISTS idx_vending_created_at ON billing_commercial.vending_and_payment_transactions(created_at);

-- 2. Customer Meter-Purchase Refund Ledger
CREATE TABLE IF NOT EXISTS billing_commercial.meter_refund_ledgers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id VARCHAR(64) NOT NULL,
    meter_serial VARCHAR(64) UNIQUE NOT NULL,
    initial_meter_cost NUMERIC(14,2) NOT NULL CHECK (initial_meter_cost > 0.00),
    refund_percentage_per_tx NUMERIC(5,2) NOT NULL DEFAULT 10.00 CHECK (refund_percentage_per_tx > 0.00 AND refund_percentage_per_tx <= 100.00),
    total_refunded_to_date NUMERIC(14,2) NOT NULL DEFAULT 0.00 CHECK (total_refunded_to_date >= 0.00),
    remaining_refund_balance NUMERIC(14,2) NOT NULL CHECK (remaining_refund_balance >= 0.00),
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refund_customer ON billing_commercial.meter_refund_ledgers(customer_id);
CREATE INDEX IF NOT EXISTS idx_refund_meter ON billing_commercial.meter_refund_ledgers(meter_serial);
CREATE INDEX IF NOT EXISTS idx_refund_status ON billing_commercial.meter_refund_ledgers(status);

-- 3. Billing Adjustments Table (Head of Commercial Operations Gate)
CREATE TABLE IF NOT EXISTS billing_commercial.billing_adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    adjustment_no VARCHAR(64) UNIQUE NOT NULL,
    ticket_id VARCHAR(64) NOT NULL,
    customer_id VARCHAR(64) NOT NULL,
    billing_officer_id VARCHAR(64) NOT NULL,
    old_balance NUMERIC(14,2) NOT NULL,
    adjusted_balance NUMERIC(14,2) NOT NULL,
    adjustment_amount NUMERIC(14,2) NOT NULL,
    root_cause_category VARCHAR(64) NOT NULL,
    root_cause_note TEXT NOT NULL,
    head_commercial_ops_id VARCHAR(64),
    head_commercial_ops_approved BOOLEAN NOT NULL DEFAULT FALSE,
    head_commercial_ops_approved_at TIMESTAMPTZ,
    head_commercial_ops_notes TEXT,
    status VARCHAR(32) NOT NULL DEFAULT 'PROPOSED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_adj_no ON billing_commercial.billing_adjustments(adjustment_no);
CREATE INDEX IF NOT EXISTS idx_adj_ticket ON billing_commercial.billing_adjustments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_adj_customer ON billing_commercial.billing_adjustments(customer_id);
CREATE INDEX IF NOT EXISTS idx_adj_status ON billing_commercial.billing_adjustments(status);

-- 4. Revenue Protection AI Alerts Table
CREATE TABLE IF NOT EXISTS billing_commercial.revenue_protection_ai_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_no VARCHAR(64) UNIQUE NOT NULL,
    customer_or_dt_id VARCHAR(64) NOT NULL,
    entity_type VARCHAR(32) NOT NULL,
    alert_type VARCHAR(32) NOT NULL,
    confidence_score NUMERIC(4,3) NOT NULL CHECK (confidence_score >= 0.000 AND confidence_score <= 1.000),
    evidence_summary TEXT NOT NULL,
    assigned_investigation_unit_id VARCHAR(64),
    investigation_notes TEXT,
    status VARCHAR(32) NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_alert_no ON billing_commercial.revenue_protection_ai_alerts(alert_no);
CREATE INDEX IF NOT EXISTS idx_ai_alert_entity ON billing_commercial.revenue_protection_ai_alerts(customer_or_dt_id);
CREATE INDEX IF NOT EXISTS idx_ai_alert_type ON billing_commercial.revenue_protection_ai_alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_ai_alert_status ON billing_commercial.revenue_protection_ai_alerts(status);

