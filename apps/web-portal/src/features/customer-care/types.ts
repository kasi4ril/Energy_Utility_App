/**
 * Type definitions for Customer Care & Consumer Relations Department.
 */

export type CustomerCareTask =
  | "analytics"
  | "register"
  | "complain"
  | "metering-kct"
  | "closure-desk"
  | "triage-queue"
  | "log-complaint"
  | "kyc-onboarding"
  | "crm-directory";

export type TicketPriority = "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
export type TicketStatus = "LOGGED" | "ROUTED" | "ATTENDED_TO" | "CLOSED" | "REOPENED_ESCALATED";
export type TicketChannel = "WALK_IN" | "CALL_CENTER" | "ONLINE_PORTAL" | "WHATSAPP" | "FIELD_AGENT";

export type TicketCategory =
  | "FALLEN_POLE"
  | "LINE_SNAP"
  | "TRANSFORMER_FAULT"
  | "FUSE_BLOWN"
  | "OIL_LEAK"
  | "FAULTY_METER"
  | "NEW_METER_SCHEDULE"
  | "WRONG_BILLING"
  | "UNCREDITED_PAYMENT"
  | "GENERAL_INQUIRY"
  | "SAFETY_HAZARD_EMERGENCY";

export interface ServiceTicket {
  id: string;
  ticket_number: string;
  customer_id?: string;
  customer_name?: string;
  customer_phone?: string;
  channel: TicketChannel;
  category: TicketCategory;
  sub_category?: string;
  is_safety_emergency?: boolean;
  nerc_escalation_tier?: "TIER_1" | "TIER_2" | "TIER_3";
  priority: TicketPriority;
  unit_type: string;
  assigned_unit_name?: string;
  description: string;
  status: TicketStatus;
  feeder_line_code?: string;
  dt_station_code?: string;
  service_address?: string;
  landmark?: string;
  sla_due_at: string;
  logged_by_cc_staff_id?: string;
  logged_by_cc_staff_name?: string;
  routed_by_cc_staff_id?: string;
  routed_by_cc_staff_name?: string;
  attended_proof_note?: string;
  work_done_summary?: string;
  attended_by_staff_id?: string;
  attended_by_staff_name?: string;
  customer_call_verified?: boolean;
  verified_by_cc_officer_id?: string;
  closed_by_cc_staff_id?: string;
  closed_by_cc_staff_name?: string;
  customer_feedback_notes?: string;
  closed_at?: string;
  created_at: string;
}

export type BuildingClass = "RESIDENTIAL" | "COMMERCIAL" | "INDUSTRIAL" | "SPECIAL";
export type TariffClass = "R2 Residential" | "C1 Commercial" | "D1 Industrial" | "MD Maximum Demand";
export type TariffBand = "BAND_A" | "BAND_B" | "BAND_C" | "BAND_D" | "BAND_E";
export type KYCStatus = "VERIFIED" | "PENDING" | "SUSPENDED" | "DEBTOR_FLAGGED";

export interface CustomerRecord {
  id: string;
  account_no: string;
  cin: string;
  full_name: string;
  phone: string;
  email?: string;
  nin?: string;
  building_class: BuildingClass;
  tariff_class?: string;
  estimated_load_kva?: number;
  landmark?: string;
  service_address: string;
  transformer_id: string;
  feeder_id?: string;
  feeder_line_code: string;
  meter_serial: string;
  map_scheme_type?: string;
  meter_type: "PREPAID_STS" | "POSTPAID" | "SMART_AMI";
  tariff_band: TariffBand;
  ledger_balance: number;
  kyc_status: KYCStatus;
  latitude: number;
  longitude: number;
  registered_by_cc_staff_id?: string;
  registered_by_cc_staff_name?: string;
  created_at: string;
}

export interface KYCOnboardingPayload {
  nin?: string;
  full_name?: string;
  phone?: string;
  email?: string;
  meter_serial?: string;
  meter_type?: "PREPAID_STS" | "POSTPAID" | "SMART_AMI";
  feeder_id?: string;
  feeder_line_code?: string;
  transformer_id?: string;
  tariff_band?: TariffBand;
  tariff_class?: string;
  building_class?: BuildingClass;
  estimated_load_kva?: number;
  landmark?: string;
  service_address?: string;
  latitude?: number;
  longitude?: number;
  created_by_staff_id?: string;
  registered_by_cc_staff_name?: string;
  map_scheme_type?: string;
}

export interface NewComplaintPayload {
  customer_id?: string;
  customer_name?: string;
  customer_phone?: string;
  service_address?: string;
  landmark?: string;
  channel?: TicketChannel;
  category?: TicketCategory;
  sub_category?: string;
  is_safety_emergency?: boolean;
  priority?: TicketPriority;
  assigned_unit_id?: string;
  assigned_unit_name?: string;
  feeder_line_code?: string;
  dt_station_code?: string;
  description?: string;
  logged_by_cc_staff_id?: string;
  logged_by_cc_staff_name?: string;
  routed_by_cc_staff_id?: string;
  routed_by_cc_staff_name?: string;
}

export interface TIDRolloverRecord {
  meter_serial: string;
  kct_1: string;
  kct_2: string;
  tid_status: string;
  generated_at: string;
  operator_staff_id: string;
}

export interface STSTokenRecord {
  token_number: string;
  units_kwh: number;
  amount_paid: number;
  purchased_at: string;
  tariff_rate: number;
}

export interface StaffLeaderboardEntry {
  staff_id: string;
  staff_name: string;
  role: string;
  logged_count: number;
  routed_count: number;
  attended_count: number;
  closed_count: number;
}

export interface DepartmentMTTREntry {
  department: string;
  avg_mttr_hours: number;
  target_hours: number;
  open_count: number;
  closed_count: number;
}

export interface CustomerCareAnalytics {
  total_tickets_logged: number;
  open_tickets_count: number;
  attended_pending_close_count: number;
  closed_tickets_count: number;
  first_contact_resolution_pct: number;
  sla_compliance_pct: number;
  average_mtta_minutes: number;
  average_mttr_hours: number;
  customer_satisfaction_score: number;
  narrative_story?: string;
  channel_distribution: { channel: TicketChannel; count: number; pct: number }[];
  category_distribution: { category: TicketCategory; label: string; count: number; pct: number }[];
  hotspot_feeders: { feeder: string; fault_count: number; open_count: number }[];
  staff_leaderboard?: StaffLeaderboardEntry[];
  department_mttr?: DepartmentMTTREntry[];
}

export type AnalyticsPeriod = "DAILY" | "WEEKLY" | "MONTHLY" | "QUARTERLY" | "YEARLY";

export interface PeriodMetrics {
  period: AnalyticsPeriod;
  label: string;
  lodged_count: number;
  resolved_count: number;
  unresolved_count: number;
  new_customers_count: number;
  resolved_pct: number;
  unresolved_pct: number;
}

export interface BillingVendingStats {
  month: string;
  postpaid_paid_count: number;
  postpaid_paid_amount: number;
  postpaid_unpaid_count: number;
  postpaid_unpaid_amount: number;
  prepaid_vended_count: number;
  prepaid_vended_amount: number;
  prepaid_unvended_count: number;
}

export interface MeterCategoryStats {
  total_installed: number;
  residential_count: number;
  md_count: number;
  industrial_count: number;
  period_view: "WEEKLY" | "MONTHLY";
  period_installs: {
    period_label: string;
    residential: number;
    md: number;
    industrial: number;
    total: number;
  }[];
}

export interface NINVerificationResult {
  nin: string;
  verified_name: string;
  gender: string;
  date_of_birth: string;
  status: "VERIFIED" | "INVALID";
}

