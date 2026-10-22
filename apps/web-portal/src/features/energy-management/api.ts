/**
 * Energy Management, Central Dispatch & Safety PTW API client (S06).
 */

import { apiClient } from "../../lib/api-client";

export interface FeederStatusDTO {
  id: string;
  feeder_code: string;
  feeder_name: string;
  voltage_level_kv: number;
  current_load_mw: number;
  state: "IN_SERVICE" | "OUT_ON_FAULT" | "OUT_ON_PERMIT_PTW" | "OUT_ON_LOAD_SHEDDING";
  priority_tier: number;
  tariff_band: string;
}

export interface SafetyPermitDTO {
  id: string;
  permit_number: string;
  feeder_id: string;
  substation_code: string;
  isolation_point: string;
  step: "REQUESTED" | "SAFETY_APPROVED" | "ISOLATED" | "WORK_IN_PROGRESS" | "CLEARANCE_DECLARED" | "RESTORED";
  requested_by_staff_id: string;
  approved_by_safety_id?: string;
  isolated_by_operator_id?: string;
  created_at: string;
}

export async function fetchFeeders(): Promise<FeederStatusDTO[]> {
  try {
    return await apiClient.get<FeederStatusDTO[]>("/energy/feeders");
  } catch {
    return [];
  }
}

export async function fetchPTWPermits(): Promise<SafetyPermitDTO[]> {
  try {
    return await apiClient.get<SafetyPermitDTO[]>("/energy/ptw");
  } catch {
    return [];
  }
}

export async function executeLoadShedding(shedMW: number): Promise<any> {
  return apiClient.post("/energy/load-shedding", { required_shed_mw: shedMW });
}

