-- Migration: 007_s07_stores_audit.down.sql
-- Description: Reversibly tears down stores_audit schema and associated tables.

DROP TABLE IF EXISTS stores_audit.staff_allowance_claims CASCADE;
DROP TABLE IF EXISTS stores_audit.scrap_returns CASCADE;
DROP TABLE IF EXISTS stores_audit.material_requisitions CASCADE;
DROP TABLE IF EXISTS stores_audit.store_inventory CASCADE;

DROP SCHEMA IF EXISTS stores_audit CASCADE;

