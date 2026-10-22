-- ==============================================================================
-- DATABASE ROLLBACK: 003_s03_customer_care.down.sql
-- Bounded Context: S03 Customer Care, Anti-Ghost GPS Lock & Mandatory Closure Gate
-- Purpose: Safely rolls back schema additions in reverse foreign-key order.
-- ==============================================================================

-- 1. Drop tables in reverse-dependency sequence
DROP TABLE IF EXISTS customercare.meter_schedules CASCADE;
DROP TABLE IF EXISTS customercare.service_tickets CASCADE;
DROP TABLE IF EXISTS customercare.customers CASCADE;

-- 2. Drop customercare schema
DROP SCHEMA IF EXISTS customercare CASCADE;

