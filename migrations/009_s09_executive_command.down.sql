-- Migration: 009_s09_executive_command.down.sql
-- Description: Reversible rollback for schema executive_command and associated tables.

DROP TABLE IF EXISTS executive_command.executive_action_items CASCADE;
DROP TABLE IF EXISTS executive_command.predictive_maintenance_triggers CASCADE;
DROP TABLE IF EXISTS executive_command.executive_snapshots CASCADE;

DROP SCHEMA IF EXISTS executive_command CASCADE;

