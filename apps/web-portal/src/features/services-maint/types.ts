/**
 * Services & Maintenance Feeder Unit Domain Types.
 * Strictly <= 350 lines.
 */

export type FeederWorkspace =
  | "analytics"
  | "outage-request"
  | "proof-of-work"
  | "standards"
  | "work-orders";

export type FeederCorridor =
  | "FEEDER_33KV_ANGLO_JOS"
  | "FEEDER_33KV_BUKURU"
  | "FEEDER_33KV_RAYFIELD"
  | "FEEDER_33KV_MAKERI"
  | "FEEDER_33KV_BOKKOS"
  | "FEEDER_11KV_VOM"
  | "FEEDER_11KV_ZAWAN";

export type OutageType =
  | "EMERGENCY_FAULT_REPAIR"
  | "PLANNED_PREVENTIVE_MAINTENANCE"
  | "VEGETATION_BUSH_CLEARING"
  | "RECONDUCTORING_PROJECT";

export type OutageStatus =
  | "PENDING_DISPATCH"
  | "APPROVED_PTW_ISSUED"
  | "ISOLATED_IN_PROGRESS"
  | "COMPLETED_LINE_ENERGIZED"
  | "REJECTED_BY_DISPATCH";

export type FeederClosureOutcome =
  | "CLOSED_ENERGIZED_IMMEDIATELY"
  | "LEFT_ON_LOAD_SHEDDING"
  | "TRIPPED_ON_CLOSURE"
  | "PENDING_SURRENDER";

export interface FeederWorkOrder {
  id: string;
  order_no: string;
  ticket_ref?: string;
  corridor: FeederCorridor;
  asset_id: string; // e.g. "POLE-AJ-142"
  title: string;
  fault_type: string;
  priority: "CRITICAL" | "HIGH" | "MEDIUM" | "LOW";
  status: "ASSIGNED" | "IN_PROGRESS" | "AWAITING_PROOF" | "ATTENDED_TO";
  assigned_team_lead: string;
  ptw_required: boolean;
  ptw_permit_no?: string;
  created_at: string;
  instructions: string;
}

export interface FeederOutageRequest {
  id: string;
  request_no: string;
  corridor: FeederCorridor;
  outage_type: OutageType;
  urgency: "EMERGENCY" | "URGENT" | "ROUTINE";
  requested_date: string;
  start_time: string;
  duration_hours: number;
  opening_time_wat?: string;         // Opening / Granting time by Central Dispatch
  surrendering_time_wat?: string;   // Time permit is formally surrendered
  time_returned_wat?: string;       // Time permit was returned upon exhaustion
  action_taken_on_return?: string;  // Detailed action taken on return / permit exhaustion
  is_permit_exhausted?: boolean;    // Whether permit time limit has been exhausted
  closure_outcome?: FeederClosureOutcome; // Closed immediately, left on load shedding, or tripped on closure
  estimated_load_mw: number;
  affected_dts_count: number;
  isolation_isolator: string; // e.g. "ISO-AJ-04"
  earth_switch: string;       // e.g. "ES-AJ-04A"
  justification: string;
  requested_by_staff: string;
  status: OutageStatus;
  ptw_issued_no?: string;
  created_at: string;
}

export interface FeederProofOfWork {
  id: string;
  work_order_id: string;
  order_no: string;
  asset_id: string;
  team_lead_name: string;
  team_lead_id: string;
  captured_gps: {
    lat: number;
    lng: number;
    accuracy_meters: number;
  };
  asset_gps: {
    lat: number;
    lng: number;
  };
  distance_from_asset_meters: number;
  gps_within_tolerance: boolean; // < 15 meters
  before_photo_url: string;
  after_photo_url: string;
  work_details: string;
  materials_consumed: string[];
  submitted_at: string;
  verification_status: "PENDING_UNIT_HEAD" | "VERIFIED_APPROVED" | "REJECTED_REVISIT";
  verified_by_staff_id?: string;
  verification_notes?: string;
}

export interface FeederPatrolRecord {
  id: string;
  patrol_code: string;
  corridor: FeederCorridor;
  pole_span: string; // e.g. "Poles #01 to #85"
  patrol_date: string;
  lead_patroller: string;
  vegetation_risk_level: "SAFE" | "ATTENTION_REQUIRED" | "IMMINENT_FLASHOVER";
  insulators_condition: "ALL_INTACT" | "CRACKED_PIN_FOUND" | "FLASHED_DISC";
  hotspot_delta_c: number; // e.g. 4.2 C
  earth_resistance_ohms: number; // < 5.0 ohms standard
  compliance_status: "PASSED" | "ACTION_ORDERED";
}

export interface FeederReliabilityMetrics {
  total_feeders: number;
  in_service_count: number;
  out_fault_count: number;
  out_ptw_count: number;
  out_rationing_count: number;
  saidi_hours: number;
  saifi_frequency: number;
  caidi_hours: number;
  mttr_hours: number;
  monthly_trips_count: number;
}

