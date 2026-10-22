-- ==============================================================================
-- DATABASE MIGRATION ROLLBACK: 005_s05_generation_master.down.sql
-- Bounded Context: S05 Generation Department & Master Station Hub
-- Target Database: PostgreSQL 16+
-- ==============================================================================

-- Drop tables in reverse order of foreign key dependency
DROP TABLE IF EXISTS generation_master.plant_fault_and_workshop_jobs CASCADE;
DROP TABLE IF EXISTS generation_master.remote_plant_shift_rollcall CASCADE;
DROP TABLE IF EXISTS generation_master.master_station_telemetry_logs CASCADE;

-- Drop schema
DROP SCHEMA IF EXISTS generation_master CASCADE;

