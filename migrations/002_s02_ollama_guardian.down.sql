-- ==============================================================================
-- DATABASE ROLLBACK: 002_s02_ollama_guardian.down.sql
-- Bounded Context: S02 Ollama Local AI Guardian, Cryptographic Ledger & Maker-Checker
-- Purpose: Safely rolls back schema additions in reverse foreign-key order.
-- ==============================================================================

-- 1. Drop immutability trigger and function
DROP TRIGGER IF EXISTS trg_immutable_audit_chain ON security_ai.immutable_audit_chain;
DROP FUNCTION IF EXISTS security_ai.prevent_audit_chain_tampering();

-- 2. Drop tables in reverse-dependency sequence
DROP TABLE IF EXISTS security_ai.community_chat_messages CASCADE;
DROP TABLE IF EXISTS security_ai.maker_checker_requests CASCADE;
DROP TABLE IF EXISTS security_ai.immutable_audit_chain CASCADE;

-- 3. Drop schema
DROP SCHEMA IF EXISTS security_ai CASCADE;

