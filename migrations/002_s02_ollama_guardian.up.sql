-- ==============================================================================
-- DATABASE MIGRATION: 002_s02_ollama_guardian.up.sql
-- Bounded Context: S02 Ollama Local AI Guardian, Cryptographic Ledger & Maker-Checker
-- Target Database: PostgreSQL 16+
-- Compliance: Non-Negotiable SHA-256 Chaining & Dual-Brain Integrity Oversight
-- ==============================================================================

-- 1. Create Isolated Security and Local AI Governance Schema
CREATE SCHEMA IF NOT EXISTS security_ai;

-- ------------------------------------------------------------------------------
-- Table 1: security_ai.immutable_audit_chain
-- Purpose: Tamper-evident, cryptographically chained mutation ledger for all sensitive records.
--          Formula: current_hash = SHA256(prev_hash + table + id + action + payload + actor + ts).
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS security_ai.immutable_audit_chain (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_name VARCHAR(60) NOT NULL,
    record_id VARCHAR(64) NOT NULL,
    action_type VARCHAR(30) NOT NULL CHECK (
        action_type IN ('INSERT', 'UPDATE', 'DELETE', 'MAKER_CHECKER_APPROVAL')
    ),
    old_json JSONB,
    new_json JSONB NOT NULL,
    maker_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    checker_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    prev_hash VARCHAR(64) NOT NULL,
    current_hash VARCHAR(64) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_chain_entity ON security_ai.immutable_audit_chain(table_name, record_id);
CREATE INDEX IF NOT EXISTS idx_audit_chain_maker ON security_ai.immutable_audit_chain(maker_id);
CREATE INDEX IF NOT EXISTS idx_audit_chain_checker ON security_ai.immutable_audit_chain(checker_id);
CREATE INDEX IF NOT EXISTS idx_audit_chain_ts ON security_ai.immutable_audit_chain(timestamp);

-- ------------------------------------------------------------------------------
-- Table 2: security_ai.maker_checker_requests
-- Purpose: Intercepts high-risk mutations (Tariff Band changes, Debt adjustments,
--          Meter swaps, Store stock edits), computes an AI Fraud Risk Score (0.00-1.00),
--          and blocks execution until approved by an authorized Checker.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS security_ai.maker_checker_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module VARCHAR(50) NOT NULL,
    change_category VARCHAR(50) NOT NULL CHECK (
        change_category IN (
            'TARIFF_BAND_CHANGE', 'DEBT_ADJUSTMENT', 'METER_SWAP',
            'STORE_STOCK_EDIT', 'ALLOWANCE_OVERRIDE'
        )
    ),
    target_record_id VARCHAR(64) NOT NULL,
    payload_json JSONB NOT NULL,
    justification TEXT NOT NULL,
    ollama_risk_score NUMERIC(4, 2) NOT NULL DEFAULT 0.00 CHECK (
        ollama_risk_score >= 0.00 AND ollama_risk_score <= 1.00
    ),
    ollama_reasoning TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING' CHECK (
        status IN ('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED')
    ),
    maker_id UUID NOT NULL REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    checker_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE RESTRICT,
    checker_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_maker_checker_status ON security_ai.maker_checker_requests(status);
CREATE INDEX IF NOT EXISTS idx_maker_checker_module ON security_ai.maker_checker_requests(module);
CREATE INDEX IF NOT EXISTS idx_maker_checker_maker ON security_ai.maker_checker_requests(maker_id);
CREATE INDEX IF NOT EXISTS idx_maker_checker_checker ON security_ai.maker_checker_requests(checker_id);

-- ------------------------------------------------------------------------------
-- Table 3: security_ai.community_chat_messages
-- Purpose: Enterprise community channels and persistent @GridAssistant chat queries.
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS security_ai.community_chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id VARCHAR(50) NOT NULL DEFAULT 'general',
    sender_id UUID REFERENCES foundation.staff_profiles(id) ON DELETE SET NULL,
    sender_name VARCHAR(150) NOT NULL,
    message_text TEXT NOT NULL,
    is_ai_response BOOLEAN NOT NULL DEFAULT FALSE,
    referenced_entity_id VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chat_messages_channel ON security_ai.community_chat_messages(channel_id, created_at);
CREATE INDEX IF NOT EXISTS idx_chat_messages_sender ON security_ai.community_chat_messages(sender_id);

-- ------------------------------------------------------------------------------
-- Trigger: Immutability Guarantee on security_ai.immutable_audit_chain
-- Core Directive: Prevent raw UPDATE or DELETE on cryptographic audit trail.
-- ------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION security_ai.prevent_audit_chain_tampering()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Cryptographic Violation: security_ai.immutable_audit_chain is an append-only ledger. UPDATE and DELETE operations are forbidden.';
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_immutable_audit_chain ON security_ai.immutable_audit_chain;
CREATE TRIGGER trg_immutable_audit_chain
    BEFORE UPDATE OR DELETE ON security_ai.immutable_audit_chain
    FOR EACH ROW
    EXECUTE FUNCTION security_ai.prevent_audit_chain_tampering();

