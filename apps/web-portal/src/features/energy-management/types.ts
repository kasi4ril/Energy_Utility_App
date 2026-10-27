/**
 * Types & Domain Interfaces for Central Energy Management & Transmission Dispatch.
 * Adheres to NERC Grid Code & 5-Step Electrical Safety PTW Isolation Protocol.
 */

export type EnergyWorkspace =
  | "analytics"
  | "load-shedding"
  | "safety-ptw"
  | "dispatch-log"
  | "outage-clearance";

export interface FeederDispatchItem {
  id: string;
  feeder_code: string;
  feeder_name: string;
  substation: string;
  voltage_level_kv: number;
  current_load_mw: number;
  capacity_mw: number;
  state: "IN_SERVICE" | "OUT_ON_FAULT" | "OUT_ON_PERMIT_PTW" | "OUT_ON_LOAD_SHEDDING";
  priority_tier: number; // 1 (Highest) to 5 (Lowest)
  tariff_band: "Band A" | "Band B" | "Band C" | "Band D" | "Band E";
  min_supply_hours_day: number;
  actual_hours_today: number;
}

export type PTWStep =
  | "REQUESTED"
  | "SAFETY_APPROVED"
  | "ISOLATED"
  | "WORK_IN_PROGRESS"
  | "CLEARANCE_DECLARED"
  | "RESTORED";

export interface SafetyPermitRecord {
  id: string;
  permit_number: string;
  feeder_code: string;
  feeder_name: string;
  substation_code: string;
  isolation_point: string;
  earth_switch_tag: string;
  step: PTWStep;
  requested_by_staff_id: string;
  requested_by_name: string;
  approved_by_safety_id?: string;
  isolated_by_operator_id?: string;
  restored_by_operator_id?: string;
  created_at_wat: string;
  updated_at_wat: string;
  remarks: string;
}

export interface HourlyDispatchEntry {
  id: string;
  hour_slot_wat: string;
  system_frequency_hz: number;
  grid_voltage_33kv: number;
  hydro_generation_mw: number;
  tcn_intake_mw: number;
  total_bulk_mw: number;
  load_shed_mw: number;
  dispatcher_name: string;
  dispatcher_staff_id: string;
  grid_status: "NORMAL_STABLE" | "FREQUENCY_STRESS" | "UNDER_FREQUENCY" | "VOLTAGE_SAG";
  shift_remarks: string;
}

export interface BulkIntakeSource {
  name: string;
  source_type: "HYDRO" | "TCN_GRID";
  output_mw: number;
  capacity_mw: number;
  status: "ONLINE" | "PARTIAL" | "OFFLINE";
  color: string;
}

