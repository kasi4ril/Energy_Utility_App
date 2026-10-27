/**
 * Types & Domain Interfaces for Billing, Commercial Operations & Revenue Protection AI.
 */

export type BillingWorkspace =
  | "analytics"
  | "vending"
  | "refunds"
  | "disputes"
  | "theft-ai";

export interface VendingRecord {
  id: string;
  transaction_ref: string;
  customer_name: string;
  account_or_meter: string;
  tariff_band: "Band A" | "Band B" | "Band C" | "Band D" | "Band E";
  amount_paid_ngn: number;
  tariff_rate_ngn: number;
  kwh_vended: number;
  sts_token_20_digit: string;
  meter_refund_credited_ngn: number;
  net_energy_amount_ngn: number;
  payment_channel: "OPay" | "Moniepoint" | "GTBank" | "Zenith Bank" | "Cash Counter";
  created_at_wat: string;
}

export interface MeterRefundLedgerEntry {
  id: string;
  customer_name: string;
  meter_number: string;
  original_meter_cost_ngn: number;
  total_refunded_to_date_ngn: number;
  remaining_balance_ngn: number;
  monthly_amortization_rate_pct: number;
  last_rebate_wat: string;
  status: "ACTIVE_AMORTIZING" | "FULLY_REFUNDED";
}

export interface BillingDisputeItem {
  id: string;
  ticket_number: string;
  customer_name: string;
  account_number: string;
  category: "MD_CT_RATIO_MULTIPLIER" | "ESTIMATED_OVERBILLING" | "UNCREDITED_VEND" | "TARIFF_MIGRATION";
  disputed_amount_ngn: number;
  status: "PENDING_AUDIT" | "RESOLVED_CREDITED" | "REJECTED_VERIFIED";
  resolution_notes?: string;
  created_at_wat: string;
}

export interface TheftAIAlertItem {
  id: string;
  target_id: string;
  target_name: string;
  entity_type: "DISTRIBUTION_TRANSFORMER" | "MAXIMUM_DEMAND_CUSTOMER" | "RESIDENTIAL_PREMISES";
  alert_type: "ZERO_CONSUMPTION_HIGH_LOAD" | "CT_BYPASS_TAMPER" | "UNMETERED_TAP" | "REVERSE_ENERGY_FLOW";
  confidence_score: number; // 0.0 to 1.0
  evidence_summary: string;
  status: "UNRESOLVED_CRITICAL" | "INSPECTION_DISPATCHED" | "CONFIRMED_PENALIZED" | "FALSE_POSITIVE";
  created_at_wat: string;
}

