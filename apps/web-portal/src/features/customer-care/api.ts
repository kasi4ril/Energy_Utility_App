/**
 * Customer Care & Consumer Relations API Client with Offline Resilience.
 */

import { apiClient } from "../../lib/api-client";
import {
  ServiceTicket,
  CustomerRecord,
  NewComplaintPayload,
  KYCOnboardingPayload,
  CustomerCareAnalytics,
  STSTokenRecord,
  TIDRolloverRecord,
  AnalyticsPeriod,
  PeriodMetrics,
  BillingVendingStats,
  MeterCategoryStats,
  NINVerificationResult,
} from "./types";
import {
  SEEDED_CUSTOMERS,
  SEEDED_TICKETS,
  MOCK_TOKENS,
  MOCK_ANALYTICS,
} from "./mockData";
import {
  PERIOD_METRICS_MAP,
  MONTHLY_BILLING_VENDING_DATA,
  METER_INSTALLATION_STATS,
  lookupNIMCName,
} from "./analyticsData";

export type ServiceTicketDTO = ServiceTicket;
export type CustomerOnboardDTO = KYCOnboardingPayload;

export async function fetchCustomerTickets(): Promise<ServiceTicket[]> {
  try {
    const data = await apiClient.get<ServiceTicket[]>("/tickets");
    if (Array.isArray(data) && data.length > 0) return data;
  } catch {
    // fallback to mock
  }
  return SEEDED_TICKETS;
}

export async function fetchCustomerDirectory(): Promise<CustomerRecord[]> {
  try {
    const data = await apiClient.get<CustomerRecord[]>("/customers");
    if (Array.isArray(data) && data.length > 0) return data;
  } catch {
    // fallback to mock
  }
  return SEEDED_CUSTOMERS;
}

export async function createServiceTicket(payload: any): Promise<ServiceTicket> {
  const isEmergency = payload.is_safety_emergency || payload.category === "LINE_SNAP";
  const newTicket: ServiceTicket = {
    id: `TKT-${Date.now().toString().slice(-4)}`,
    ticket_number: `TKT-2026-${Date.now().toString().slice(-4)}`,
    customer_id: payload.customer_id || `CUST-WALK-${Date.now().toString().slice(-4)}`,
    customer_name: payload.customer_name,
    customer_phone: payload.customer_phone,
    service_address: payload.service_address,
    channel: payload.channel,
    category: payload.category,
    sub_category: payload.sub_category,
    is_safety_emergency: isEmergency,
    nerc_escalation_tier: isEmergency ? "TIER_3" : "TIER_1",
    priority: payload.priority || (isEmergency ? "CRITICAL" : "MEDIUM"),
    unit_type: payload.assigned_unit_id,
    assigned_unit_name: payload.assigned_unit_name,
    feeder_line_code: payload.feeder_line_code,
    dt_station_code: payload.dt_station_code,
    landmark: payload.landmark,
    description: payload.description,
    logged_by_cc_staff_name: payload.logged_by_cc_staff_name || "Amina Yusuf (CC Officer)",
    routed_by_cc_staff_name: payload.routed_by_cc_staff_name || "Amina Yusuf (CC Officer)",
    status: "ROUTED",
    sla_due_at: new Date(Date.now() + (isEmergency ? 2 : 4) * 3600 * 1000).toISOString(),
    created_at: new Date().toISOString(),
  };

  try {
    const res = await apiClient.post<ServiceTicket>("/tickets", payload);
    return res && res.ticket_number ? res : newTicket;
  } catch {
    return newTicket;
  }
}

export async function onboardCustomer(payload: any): Promise<CustomerRecord> {
  const newCust: CustomerRecord = {
    id: `CUST-2026-${Date.now().toString().slice(-4)}`,
    account_no: `ACC-${Date.now().toString().slice(-6)}`,
    cin: `CIN-2026-${Date.now().toString().slice(-5)}`,
    full_name: payload.full_name,
    phone: payload.phone,
    email: payload.email,
    nin: payload.nin,
    building_class: payload.building_class,
    tariff_class: payload.tariff_class || "R2 Residential",
    estimated_load_kva: payload.estimated_load_kva || 5.0,
    landmark: payload.landmark,
    service_address: payload.service_address,
    transformer_id: payload.transformer_id,
    feeder_id: payload.feeder_id,
    feeder_line_code: payload.feeder_line_code,
    meter_serial: payload.meter_serial,
    map_scheme_type: payload.map_scheme_type || "Single-Phase MAP",
    meter_type: payload.meter_type || "PREPAID_STS",
    tariff_band: payload.tariff_band,
    ledger_balance: 0.0,
    kyc_status: "VERIFIED",
    latitude: payload.latitude,
    longitude: payload.longitude,
    registered_by_cc_staff_id: payload.created_by_staff_id || "STAFF-CC-01",
    registered_by_cc_staff_name: payload.registered_by_cc_staff_name || "Sani Bello (CC Desk)",
    created_at: new Date().toISOString(),
  };

  try {
    const res = await apiClient.post<CustomerRecord>("/customers/onboard", payload);
    return res && res.cin ? res : newCust;
  } catch {
    return newCust;
  }
}

export async function confirmTicketClosure(
  ticketId: string,
  notes: string,
  staffId: string,
  officerName?: string,
  confirmedByCustomer: boolean = true
): Promise<any> {
  try {
    return await apiClient.post(`/tickets/${ticketId}/confirm-close`, {
      confirmed_by_customer: confirmedByCustomer,
      feedback_notes: notes,
      officer_name: officerName || "Customer Care Representative",
      staff_id: staffId,
    });
  } catch {
    return { status: confirmedByCustomer ? "CLOSED" : "REOPENED_ESCALATED" };
  }
}

export async function fetchCustomerAnalytics(): Promise<CustomerCareAnalytics> {
  return MOCK_ANALYTICS;
}

export async function fetchRecentTokens(_meterSerial: string): Promise<STSTokenRecord[]> {
  return MOCK_TOKENS;
}

export async function generateTIDRolloverKCT(
  meterSerial: string,
  staffId: string = "STAFF-CC-01"
): Promise<TIDRolloverRecord> {
  try {
    const res = await apiClient.post<TIDRolloverRecord>("/customer-care/tokens/kct", {
      meter_serial: meterSerial,
      staff_id: staffId,
    });
    if (res && res.kct_1) return res;
  } catch {
    // fallback
  }

  return {
    meter_serial: meterSerial,
    kct_1: "5821-9304-1842-6092-3317",
    kct_2: "4190-2831-7490-1152-9084",
    tid_status: "GENERATED_PENDING_INPUT",
    generated_at: new Date().toISOString(),
    operator_staff_id: staffId,
  };
}

export async function scheduleMeterActivity(payload: {
  customer_id: string;
  ticket_id?: string;
  schedule_type: string;
  scheduled_date: string;
  staff_id: string;
  staff_name: string;
}): Promise<any> {
  try {
    return await apiClient.post("/meters/schedule", payload);
  } catch {
    return { status: "SCHEDULED", scheduled_date: payload.scheduled_date };
  }
}

export async function verifyNIN(nin: string): Promise<NINVerificationResult> {
  try {
    const res = await apiClient.post<NINVerificationResult>("/customers/verify-nin", { nin });
    if (res && res.verified_name) return res;
  } catch {
    // fallback
  }
  return lookupNIMCName(nin);
}

export async function fetchPeriodMetrics(period: AnalyticsPeriod): Promise<PeriodMetrics> {
  return PERIOD_METRICS_MAP[period] || PERIOD_METRICS_MAP.MONTHLY;
}

export async function fetchBillingVendingStats(monthKey: string = "2026-10"): Promise<BillingVendingStats> {
  return MONTHLY_BILLING_VENDING_DATA[monthKey] || MONTHLY_BILLING_VENDING_DATA["2026-10"];
}

export async function fetchMeterInstallationStats(): Promise<MeterCategoryStats> {
  return METER_INSTALLATION_STATS;
}
