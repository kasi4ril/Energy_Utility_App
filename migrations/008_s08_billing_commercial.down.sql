-- Migration: 008_s08_billing_commercial.down.sql
-- Description: Reversibly tears down billing_commercial schema and associated tables.

DROP TABLE IF EXISTS billing_commercial.revenue_protection_ai_alerts CASCADE;
DROP TABLE IF EXISTS billing_commercial.billing_adjustments CASCADE;
DROP TABLE IF EXISTS billing_commercial.meter_refund_ledgers CASCADE;
DROP TABLE IF EXISTS billing_commercial.vending_and_payment_transactions CASCADE;

DROP SCHEMA IF EXISTS billing_commercial CASCADE;

