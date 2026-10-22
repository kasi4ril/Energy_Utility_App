/**
 * Authentication and RBAC type definitions for the NESCO Enterprise Utility Platform.
 * Mirrors Go pkg/auth Claims structure and Blueprint Section 10.1 role matrices.
 */

export type PrimaryRole =
  | "MD"
  | "HEAD_SERVICES"
  | "HEAD_INTERNAL_AUDIT"
  | "HEAD_COMMERCIAL_OPS"
  | "GM_GENERATION"
  | "HEAD_ENERGY_MGMT"
  | "CUSTOMER_CARE_OFFICER"
  | "HEAD_CUSTOMER_CARE"
  | "UNIT_HEAD_FEEDER_MAINT"
  | "UNIT_HEAD_TX_MAINT"
  | "UNIT_HEAD_METERING"
  | "TEAM_LEAD_FEEDER"
  | "TEAM_LEAD_TX"
  | "TEAM_LEAD_METERING"
  | "SAFETY_REP_HSE"
  | "MASTER_STATION_OPERATOR"
  | "UNIT_HEAD_GEN_ELECTRICAL"
  | "UNIT_HEAD_GEN_MECHANICAL"
  | "WORKSHOP_LEAD"
  | "HR_OFFICER_PLANT"
  | "CENTRAL_DISPATCH_OPERATOR"
  | "SUBSTATION_OPERATOR"
  | "STORE_OFFICER"
  | "INTERNAL_AUDITOR"
  | "UNIT_HEAD_BILLING";

export interface AuthClaims {
  staff_id: string;
  staff_number: string;
  full_name: string;
  primary_role: PrimaryRole;
  org_unit_id: string;
  can_view_executive_dashboard: boolean;
  assigned_asset_id?: string;
  iat?: number;
  exp?: number;
}

export interface LoginResponse {
  token: string;
  claims: AuthClaims;
}

export interface ShiftStationOption {
  code: string;
  name: string;
  category: "GENERATION_STATION" | "CENTRAL_DISPATCH" | "TRANSMISSION_SUBSTATION";
}

export const SHIFT_STATIONS: ShiftStationOption[] = [
  { code: "STN-GEN-KURA", name: "Plant 1 - Kurra Falls Hydro", category: "GENERATION_STATION" },
  { code: "STN-GEN-JUKAL", name: "Plant 2 - Jekko 1 Hydro", category: "GENERATION_STATION" },
  { code: "STN-GEN-JEKKO2", name: "Plant 3 - Jekko 2 Hydro", category: "GENERATION_STATION" },
  { code: "STN-GEN-ANP", name: "Plant 4 - Ankwil Hydro", category: "GENERATION_STATION" },
  { code: "STN-GEN-TRES", name: "Plant 5 - Tenti Hydro", category: "GENERATION_STATION" },
  { code: "STN-GEN-DIESEL", name: "Plant 6 - Bukuru Thermal Backup", category: "GENERATION_STATION" },
  { code: "STN-DISPATCH-HQ", name: "Bukuru Central Dispatch & Grid Control", category: "CENTRAL_DISPATCH" },
  { code: "STN-SUB-33KV-VOM", name: "33kV Vom Transmission Substation", category: "TRANSMISSION_SUBSTATION" },
  { code: "STN-SUB-33KV-BARK", name: "33kV Barkin Ladi Substation", category: "TRANSMISSION_SUBSTATION" },
];

