/**
 * Customer Care & KYC API client (S03).
 */

import { apiClient } from "../../lib/api-client";

export interface ServiceTicketDTO {
  id: string;
  ticket_number: string;
  customer_id: string;
  category: string;
  priority: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  status: string;
  unit_type: string;
  description: string;
  created_at: string;
}

export interface CustomerOnboardDTO {
  national_id_nin: string;
  full_name: string;
  phone_number: string;
  account_number: string;
  meter_number: string;
  feeder_line_code: string;
  dt_station_code: string;
  tariff_band: string;
  service_address: string;
  latitude: number;
  longitude: number;
  created_by_staff_id: string;
}

export async function fetchCustomerTickets(): Promise<ServiceTicketDTO[]> {
  try {
    return await apiClient.get<ServiceTicketDTO[]>("/tickets");
  } catch {
    return [];
  }
}

export async function onboardCustomer(payload: CustomerOnboardDTO): Promise<any> {
  return apiClient.post("/customers/onboard", payload);
}

export async function confirmTicketClosure(ticketId: string, notes: string, staffId: string): Promise<any> {
  return apiClient.post(`/tickets/${ticketId}/confirm-close`, {
    customer_contacted: true,
    resolution_satisfactory: true,
    officer_notes: notes,
    staff_id: staffId,
  });
}

