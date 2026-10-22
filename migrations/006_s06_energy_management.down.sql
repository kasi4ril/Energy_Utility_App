-- ==============================================================================
-- DATABASE MIGRATION ROLLBACK: 006_s06_energy_management.down.sql
-- Bounded Context: S06 Energy Management, Priority Dispatch & 5-Step Safety PTW
-- Target Database: PostgreSQL 16+
-- ==============================================================================

-- Drop tables in reverse order of foreign key dependency
DROP TABLE IF EXISTS energy_management.safety_permits_ptw CASCADE;
DROP TABLE IF EXISTS energy_management.feeder_status_ledger CASCADE;
DROP TABLE IF EXISTS energy_management.dispatch_telemetry_logs CASCADE;

-- Drop schema
DROP SCHEMA IF EXISTS energy_management CASCADE;

