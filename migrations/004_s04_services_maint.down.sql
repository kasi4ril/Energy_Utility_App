-- ==============================================================================
-- DATABASE MIGRATION ROLLBACK: 004_s04_services_maint.down.sql
-- Bounded Context: S04 Services & Maintenance Department
-- Target Database: PostgreSQL 16+ with PostGIS
-- ==============================================================================

-- Drop tables in reverse order of foreign key dependency
DROP TABLE IF EXISTS services_maint.proof_of_work_logs CASCADE;
DROP TABLE IF EXISTS services_maint.maintenance_work_orders CASCADE;

-- Drop schema
DROP SCHEMA IF EXISTS services_maint CASCADE;

