/**
 * Billing, Commercial Operations & Revenue Protection AI API client (S08).
 */

import { apiClient } from "../../lib/api-client";

export interface VendingTransactionDTO {
  id: string;
  transaction_ref: string;
  customer_id: string;
  meter_number: string;
  amount_paid: number;
  tariff_rate_per_kwh: number;
  kwh_vended: number;
  sts_token_20_digit: string;
  meter_refund_credited: number;
  payment_platform_source: string;
  created_at: string;
}

export interface TheftAlertDTO {
  id: string;
  customer_or_dt_id: string;
  entity_type: string;
  alert_type: string;
  confidence_score: number;
  evidence_summary: string;
  status: string;
  created_at: string;
}

export async function fetchTransactions(): Promise<VendingTransactionDTO[]> {
  try {
    return await apiClient.get<VendingTransactionDTO[]>("/billing/vending");
  } catch {
    return [];
  }
}

export async function fetchTheftAlerts(): Promise<TheftAlertDTO[]> {
  try {
    return await apiClient.get<TheftAlertDTO[]>("/billing/revenue-protection/alerts");
  } catch {
    return [];
  }
}

export async function simulateVending(payload: {
  customer_id: string;
  meter_number: string;
  amount_paid: number;
  tariff_rate_per_kwh: number;
  payment_platform: string;
}): Promise<any> {
  return apiClient.post("/billing/vending", payload);
}

