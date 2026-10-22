/**
 * Generation Department & Master Station API Client with Offline Resilience (S05).
 */

import { apiClient } from "../../lib/api-client";
import {
  HydroPlantSummary,
  PlantTelemetryLog,
  ShiftRollcallRecord,
  PlantFaultJob,
  GenerationAnalytics,
} from "./types";

// Re-export PlantTelemetryDTO for backwards compatibility
export type PlantTelemetryDTO = PlantTelemetryLog;

const SEEDED_PLANTS: HydroPlantSummary[] = [
  {
    id: "PLANT_01_KURRA_FALLS",
    name: "Kurra Falls Station 1",
    location: "Kurra River Gorge, Plateau",
    turbine_type: "FRANCIS",
    installed_capacity_mw: 16.0,
    active_units: 3,
    total_units: 4,
    current_mw: 14.8,
    water_level_m: 24.5,
    headwater_elevation_m: 1240.2,
    tailwater_elevation_m: 980.1,
    stator_temp_c: 68.2,
    bearing_vibration_mm_s: 2.1,
    lube_oil_pressure_bar: 3.2,
    state: "IN_SERVICE",
    last_collated_at: new Date(Date.now() - 25 * 60 * 1000).toISOString(),
  },
  {
    id: "PLANT_02_KURRA_FALLS",
    name: "Kurra Falls Station 2",
    location: "Lower Cascade Falls, Plateau",
    turbine_type: "PELTON",
    installed_capacity_mw: 12.0,
    active_units: 3,
    total_units: 3,
    current_mw: 11.2,
    water_level_m: 22.8,
    headwater_elevation_m: 975.0,
    tailwater_elevation_m: 780.4,
    stator_temp_c: 64.5,
    bearing_vibration_mm_s: 1.8,
    lube_oil_pressure_bar: 3.4,
    state: "IN_SERVICE",
    last_collated_at: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
  },
  {
    id: "PLANT_03_KWALL",
    name: "Kwall Hydro Station 1",
    location: "Miango River Basin, Bassa",
    turbine_type: "PELTON",
    installed_capacity_mw: 8.0,
    active_units: 2,
    total_units: 2,
    current_mw: 7.4,
    water_level_m: 18.2,
    headwater_elevation_m: 1120.0,
    tailwater_elevation_m: 950.0,
    stator_temp_c: 62.0,
    bearing_vibration_mm_s: 2.3,
    lube_oil_pressure_bar: 3.1,
    state: "IN_SERVICE",
    last_collated_at: new Date(Date.now() - 40 * 60 * 1000).toISOString(),
  },
  {
    id: "PLANT_04_KWALL",
    name: "Kwall Hydro Station 2",
    location: "Lower Kwall Penstock Line",
    turbine_type: "FRANCIS",
    installed_capacity_mw: 6.0,
    active_units: 1,
    total_units: 2,
    current_mw: 4.8,
    water_level_m: 17.5,
    headwater_elevation_m: 945.0,
    tailwater_elevation_m: 810.0,
    stator_temp_c: 72.1,
    bearing_vibration_mm_s: 4.7, // High vibration alarm!
    lube_oil_pressure_bar: 2.8,
    state: "IN_SERVICE",
    last_collated_at: new Date(Date.now() - 15 * 60 * 1000).toISOString(),
  },
  {
    id: "PLANT_05_ANOMA",
    name: "Anoma Hydro Station",
    location: "Anoma Weir Headworks",
    turbine_type: "KAPLAN",
    installed_capacity_mw: 5.0,
    active_units: 1,
    total_units: 2,
    current_mw: 4.2,
    water_level_m: 14.2,
    headwater_elevation_m: 890.0,
    tailwater_elevation_m: 820.0,
    stator_temp_c: 59.4,
    bearing_vibration_mm_s: 1.6,
    lube_oil_pressure_bar: 3.5,
    state: "IN_SERVICE",
    last_collated_at: new Date(Date.now() - 50 * 60 * 1000).toISOString(),
  },
  {
    id: "PLANT_06_JERE",
    name: "Jere Hydro Station",
    location: "Jere Reservoir Outlet",
    turbine_type: "CROSSFLOW",
    installed_capacity_mw: 4.0,
    active_units: 0,
    total_units: 1,
    current_mw: 0.0,
    water_level_m: 9.8,
    headwater_elevation_m: 760.0,
    tailwater_elevation_m: 720.0,
    stator_temp_c: 24.0,
    bearing_vibration_mm_s: 0.0,
    lube_oil_pressure_bar: 0.0,
    state: "OUT_OF_SERVICE",
    outage_reason: "TURBINE_MECH_FAULT",
    last_collated_at: new Date(Date.now() - 10 * 60 * 1000).toISOString(),
  },
];

const SEEDED_ROLLCALLS: ShiftRollcallRecord[] = [
  {
    id: "RC-01",
    shift_id: "SH-MORNING-2026-10",
    shift_window: "MORNING",
    plant_id: "PLANT_01_KURRA_FALLS",
    operator_names: ["Ibrahim Pam", "Yusuf Chuwang"],
    protection_engineers: ["Engr. D. Gyang"],
    technicians: ["Bulus Mark", "Sunday Azi"],
    handover_notes: "Units 1, 2, 3 in full commercial generation. Penstock 4 valve packed for minor leak.",
    synced_to_hq_hr: true,
    master_operator_id: "STAFF-GEN-OP-01",
    created_at: new Date(Date.now() - 120 * 60 * 1000).toISOString(),
  },
  {
    id: "RC-02",
    shift_id: "SH-MORNING-2026-10",
    shift_window: "MORNING",
    plant_id: "PLANT_06_JERE",
    operator_names: ["Godiya Bitrus"],
    protection_engineers: ["Engr. S. Lar"],
    technicians: ["Danjuma Moses"],
    handover_notes: "Plant 6 runner dismantled for workshop lathe machining and shaft sleeve replacement.",
    synced_to_hq_hr: true,
    master_operator_id: "STAFF-GEN-OP-01",
    created_at: new Date(Date.now() - 110 * 60 * 1000).toISOString(),
  },
];

const SEEDED_FAULTS: PlantFaultJob[] = [
  {
    id: "FLT-GEN-01",
    job_no: "JOB-2026-GEN-001",
    plant_id: "PLANT_06_JERE",
    reported_by_operator_id: "STAFF-GEN-OP-01",
    wing: "MECHANICAL",
    fault_description: "Turbine runner excessive vibration (6.2 mm/s). Runner hub cavitation pitting and shaft sleeve groove.",
    severity: "CRITICAL",
    unit_head_id: "STAFF-HEAD-MECH-01",
    gm_approved: true,
    gm_approval_notes: "Approved emergency bronze alloy requisition and workshop overhaul.",
    requires_workshop: true,
    workshop_status: "MACHINING_IN_PROGRESS",
    raw_materials_summary: "Phosphor Bronze Bushing (120mm OD), 400kg Stainless Welding Rods",
    store_requisition_ref: "REQ-2026-ST-042",
    status: "WORKSHOP_IN_PROGRESS",
    created_at: new Date(Date.now() - 36 * 3600 * 1000).toISOString(),
  },
  {
    id: "FLT-GEN-02",
    job_no: "JOB-2026-GEN-002",
    plant_id: "PLANT_04_KWALL",
    reported_by_operator_id: "STAFF-GEN-OP-01",
    wing: "MECHANICAL",
    fault_description: "Unit 2 guide bearing vibration trending upwards to 4.7 mm/s. High lube oil temperature.",
    severity: "HIGH",
    unit_head_id: "STAFF-HEAD-MECH-01",
    gm_approved: false,
    requires_workshop: true,
    workshop_status: "PENDING_DISPATCH",
    status: "UNIT_HEAD_REVIEWED",
    created_at: new Date(Date.now() - 4 * 3600 * 1000).toISOString(),
  },
  {
    id: "FLT-GEN-03",
    job_no: "JOB-2026-GEN-003",
    plant_id: "PLANT_01_KURRA_FALLS",
    reported_by_operator_id: "STAFF-GEN-OP-01",
    wing: "ELECTRICAL",
    fault_description: "Unit 4 static excitation automatic voltage regulator (AVR) secondary thyristor overheating.",
    severity: "HIGH",
    unit_head_id: "STAFF-HEAD-ELEC-01",
    gm_approved: true,
    gm_approval_notes: "Spares released from Central Stores; replacement thyristor stack issued.",
    requires_workshop: false,
    workshop_status: "NOT_APPLICABLE",
    status: "GM_APPROVED",
    created_at: new Date(Date.now() - 18 * 3600 * 1000).toISOString(),
  },
];

export async function fetchPlantFleet(): Promise<HydroPlantSummary[]> {
  return SEEDED_PLANTS;
}

export async function fetchPlantTelemetry(): Promise<PlantTelemetryLog[]> {
  try {
    const data = await apiClient.get<PlantTelemetryLog[]>("/generation/telemetry");
    if (Array.isArray(data) && data.length > 0) return data;
  } catch {
    // fallback
  }
  return SEEDED_PLANTS.map((p) => ({
    id: `LOG-${p.id}`,
    shift_id: "SH-MORNING-2026",
    shift_window: "MORNING",
    plant_id: p.id,
    plant_state: p.state,
    outage_reason: p.outage_reason,
    water_level_m: p.water_level_m,
    headwater_elevation_m: p.headwater_elevation_m,
    tailwater_elevation_m: p.tailwater_elevation_m,
    inflow_rate_m3_s: 14.5,
    mw_generated: p.current_mw,
    mvar_reactive: p.current_mw * 0.28,
    stator_winding_temp_c: p.stator_temp_c,
    bearing_vibration_mm_s: p.bearing_vibration_mm_s,
    lube_oil_pressure_bar: p.lube_oil_pressure_bar,
    logged_by_operator_id: "STAFF-GEN-OP-01",
    recorded_at: p.last_collated_at,
  }));
}

export async function submitPlantTelemetry(payload: {
  plant_id: string;
  mw_generated: number;
  water_level_meters: number;
  stator_temp_celsius: number;
  bearing_vibration_mm_s: number;
  lube_oil_pressure_bar: number;
  shift_roster_id: string;
}): Promise<any> {
  return apiClient.post("/generation/telemetry", payload);
}

export async function fetchShiftRollcalls(): Promise<ShiftRollcallRecord[]> {
  try {
    const data = await apiClient.get<ShiftRollcallRecord[]>("/generation/rollcall");
    if (Array.isArray(data) && data.length > 0) return data;
  } catch {
    // fallback
  }
  return SEEDED_ROLLCALLS;
}

export async function submitShiftRollcall(payload: Partial<ShiftRollcallRecord>): Promise<any> {
  return apiClient.post("/generation/rollcall", payload);
}

export async function fetchPlantFaults(): Promise<PlantFaultJob[]> {
  try {
    const data = await apiClient.get<PlantFaultJob[]>("/generation/faults");
    if (Array.isArray(data) && data.length > 0) return data;
  } catch {
    // fallback
  }
  return SEEDED_FAULTS;
}

export async function reportPlantFault(payload: Partial<PlantFaultJob>): Promise<any> {
  return apiClient.post("/generation/faults", payload);
}

export async function approveGMFault(
  faultId: string,
  approvalNotes: string,
  gmStaffId: string
): Promise<any> {
  return apiClient.post(`/generation/faults/${faultId}/gm-approval`, {
    gm_approved: true,
    approval_notes: approvalNotes,
    gm_staff_id: gmStaffId,
  });
}

export async function updateWorkshopJob(
  faultId: string,
  workshopStatus: string,
  rawMaterials: string
): Promise<any> {
  return apiClient.post(`/generation/faults/${faultId}/workshop`, {
    workshop_status: workshopStatus,
    raw_materials_summary: rawMaterials,
  });
}

export async function fetchGenerationAnalytics(): Promise<GenerationAnalytics> {
  return {
    total_generation_mw: 42.4,
    daily_mwh_delivered: 986.5,
    reservoir_capacity_pct: 78.4,
    inflow_total_m3_s: 38.6,
    plant_availability_factor_pct: 94.2,
    equivalent_forced_outage_rate_pct: 1.8,
    vibration_alerts_count: 1,
    active_workshop_jobs_count: 1,
    plants_generation_share: [
      { plant_id: "PLANT_01_KURRA_FALLS", name: "Kurra Falls Station 1", mw: 14.8, pct: 35 },
      { plant_id: "PLANT_02_KURRA_FALLS", name: "Kurra Falls Station 2", mw: 11.2, pct: 26 },
      { plant_id: "PLANT_03_KWALL", name: "Kwall Hydro Station 1", mw: 7.4, pct: 17 },
      { plant_id: "PLANT_04_KWALL", name: "Kwall Hydro Station 2", mw: 4.8, pct: 11 },
      { plant_id: "PLANT_05_ANOMA", name: "Anoma Hydro Station", mw: 4.2, pct: 10 },
      { plant_id: "PLANT_06_JERE", name: "Jere Hydro Station", mw: 0.0, pct: 0 },
    ],
    weekly_generation_mwh: [
      { day: "Mon", hydro_mwh: 960 },
      { day: "Tue", hydro_mwh: 985 },
      { day: "Wed", hydro_mwh: 1010 },
      { day: "Thu", hydro_mwh: 975 },
      { day: "Fri", hydro_mwh: 990 },
      { day: "Sat", hydro_mwh: 980 },
      { day: "Sun", hydro_mwh: 986 },
    ],
  };
}
