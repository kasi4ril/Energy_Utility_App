/**
 * Generation Master Station Hub API client (S05).
 */

import { apiClient } from "../../lib/api-client";

export interface PlantTelemetryDTO {
  id: string;
  plant_id: string;
  mw_generated: number;
  water_level_meters: number;
  stator_temp_celsius: number;
  bearing_vibration_mm_s: number;
  lube_oil_pressure_bar: number;
  in_service: boolean;
  logged_at: string;
}

export async function fetchPlantTelemetry(): Promise<PlantTelemetryDTO[]> {
  try {
    return await apiClient.get<PlantTelemetryDTO[]>("/generation/telemetry");
  } catch {
    return [];
  }
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

