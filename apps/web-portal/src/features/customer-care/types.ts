/**
 * Type definitions for Customer Care & Consumer Relations Department.
 */

export type CustomerCareTask =
  | "triage-queue"
  | "log-complaint"
  | "kyc-onboarding"
  | "crm-directory"
  | "analytics";

export type TicketPriority = "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
export type TicketStatus = "LOGGED" | "ROUTED" | "ATTENDED_TO" | "CLOSED" | "REOPENED_ESCALATED";
export type TicketChannel = "WALK_IN" | "CALL_CENTER" | "ONLINE_PORTAL" | "FIELD_AGENT";

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
  | "GENERAL_INQUIRY";

export interface ServiceTicket {
  id: string;
  ticket_number: string;
  customer_id?: string;
  customer_name?: string;
  customer_phone?: string;
  channel: TicketChannel;
  category: TicketCategory;
  priority: TicketPriority;
  unit_type: string;
  description: string;
  status: TicketStatus;
  feeder_line_code?: string;
  dt_station_code?: string;
  service_address?: string;
  sla_due_at: string;
  attended_proof_note?: string;
  attended_by_staff_id?: string;
  customer_call_verified?: boolean;
  verified_by_cc_officer_id?: string;
  customer_feedback_notes?: string;
  created_at: string;
}

export type BuildingClass = "RESIDENTIAL" | "COMMERCIAL" | "INDUSTRIAL" | "SPECIAL";
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
  service_address: string;
  transformer_id: string;
  feeder_line_code: string;
  meter_serial: string;
  meter_type: "PREPAID_STS" | "POSTPAID" | "SMART_AMI";
  tariff_band: TariffBand;
  ledger_balance: number;
  kyc_status: KYCStatus;
  latitude: number;
  longitude: number;
  created_at: string;
}

export interface NewComplaintPayload {
  customer_id?: string;
  customer_name: string;
  customer_phone: string;
  service_address: string;
  channel: TicketChannel;
  category: TicketCategory;
  priority: TicketPriority;
  assigned_unit_id: string;
  feeder_line_code?: string;
  dt_station_code?: string;
  description: string;
}

export interface KYCOnboardingPayload {
  full_name: string;
  phone: string;
  email?: string;
  nin: string;
  building_class: BuildingClass;
  service_address: string;
  transformer_id: string;
  feeder_line_code: string;
  meter_serial: string;
  meter_type: "PREPAID_STS" | "POSTPAID" | "SMART_AMI";
  tariff_band: TariffBand;
  latitude: number;
  longitude: number;
  created_by_staff_id: string;
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
  channel_distribution: { channel: TicketChannel; count: number; pct: number }[];
  category_distribution: { category: TicketCategory; label: string; count: number; pct: number }[];
  hotspot_feeders: { feeder: string; fault_count: number; open_count: number }[];
}

