/**
 * Type definitions for Generation Department & Master Station Command Hub (S05).
 */

export type GenerationTask =
  | "fleet-board"
  | "call-collation"
  | "plant-faults"
  | "workshop"
  | "analytics";

export type PlantID =
  | "PLANT_01_KURRA_FALLS"
  | "PLANT_02_KURRA_FALLS"
  | "PLANT_03_KWALL"
  | "PLANT_04_KWALL"
  | "PLANT_05_ANOMA"
  | "PLANT_06_JERE";

export type PlantState = "IN_SERVICE" | "STANDBY" | "OUT_OF_SERVICE";

export type ShiftWindow = "MORNING" | "AFTERNOON" | "NIGHT";

export type OutageReason =
  | "HYDRO_LOW_WATER"
  | "TURBINE_MECH_FAULT"
  | "GENERATOR_ELEC_FAULT"
  | "PROTECTION_RELAY_TRIP"
  | "PLANNED_MAINTENANCE";

export type TechnicalWing = "ELECTRICAL" | "MECHANICAL";

export type WorkshopStatus =
  | "NOT_APPLICABLE"
  | "PENDING_DISPATCH"
  | "MACHINING_IN_PROGRESS"
  | "WINDING_FABRICATION"
  | "QUALITY_CHECK"
  | "COMPLETED_DELIVERED"
  | "REJECTED";

export type PlantJobStatus =
  | "REPORTED"
  | "UNIT_HEAD_REVIEWED"
  | "GM_APPROVED"
  | "WORKSHOP_IN_PROGRESS"
  | "RESOLVED_REINSTALLED"
  | "CANCELLED";

export interface HydroPlantSummary {
  id: PlantID;
  name: string;
  location: string;
  turbine_type: "FRANCIS" | "PELTON" | "KAPLAN" | "CROSSFLOW";
  installed_capacity_mw: number;
  active_units: number;
  total_units: number;
  current_mw: number;
  water_level_m: number;
  headwater_elevation_m: number;
  tailwater_elevation_m: number;
  stator_temp_c: number;
  bearing_vibration_mm_s: number;
  lube_oil_pressure_bar: number;
  state: PlantState;
  outage_reason?: OutageReason;
  last_collated_at: string;
}

export interface PlantTelemetryLog {
  id: string;
  shift_id: string;
  shift_window: ShiftWindow;
  plant_id: PlantID;
  plant_state: PlantState;
  outage_reason?: OutageReason;
  outage_details?: string;
  water_level_m: number;
  headwater_elevation_m: number;
  tailwater_elevation_m: number;
  inflow_rate_m3_s: number;
  mw_generated: number;
  mvar_reactive: number;
  stator_winding_temp_c: number;
  bearing_vibration_mm_s: number;
  lube_oil_pressure_bar: number;
  logged_by_operator_id: string;
  recorded_at: string;
}

export interface ShiftRollcallRecord {
  id: string;
  shift_id: string;
  shift_window: ShiftWindow;
  plant_id: PlantID;
  operator_names: string[];
  protection_engineers: string[];
  technicians: string[];
  handover_notes: string;
  synced_to_hq_hr: boolean;
  master_operator_id: string;
  created_at: string;
}

export interface PlantFaultJob {
  id: string;
  job_no: string;
  plant_id: PlantID;
  reported_by_operator_id: string;
  wing: TechnicalWing;
  fault_description: string;
  severity: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  unit_head_id?: string;
  gm_approved: boolean;
  gm_approval_notes?: string;
  requires_workshop: boolean;
  workshop_status: WorkshopStatus;
  raw_materials_summary?: string;
  store_requisition_ref?: string;
  status: PlantJobStatus;
  created_at: string;
}

export interface GenerationAnalytics {
  total_generation_mw: number;
  daily_mwh_delivered: number;
  reservoir_capacity_pct: number;
  inflow_total_m3_s: number;
  plant_availability_factor_pct: number;
  equivalent_forced_outage_rate_pct: number;
  vibration_alerts_count: number;
  plants_generation_share: { plant_id: string; name: string; mw: number; pct: number }[];
  weekly_generation_mwh: { day: string; hydro_mwh: number }[];
  active_workshop_jobs_count: number;
}
