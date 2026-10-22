-- ==============================================================================
-- DATABASE ROLLBACK: 001_s01_foundation.down.sql
-- Bounded Context: S01 Foundation, Grid Topology, IAM, Federated HR & WorkNotes
-- Purpose: Safely rolls back schema additions in reverse foreign-key order.
-- ==============================================================================

-- 1. Drop trigger and immutability function
DROP TRIGGER IF EXISTS trg_immutable_work_notes ON foundation.work_notes;
DROP FUNCTION IF EXISTS foundation.prevent_work_notes_tampering();

-- 2. Drop tables in reverse-dependency sequence
DROP TABLE IF EXISTS foundation.work_notes CASCADE;
DROP TABLE IF EXISTS foundation.grid_assets CASCADE;
DROP TABLE IF EXISTS foundation.staff_profiles CASCADE;
DROP TABLE IF EXISTS foundation.org_units CASCADE;

-- 3. Drop foundation schema
DROP SCHEMA IF EXISTS foundation CASCADE;

