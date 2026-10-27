/**
 * Types & Domain Interfaces for Stores & Stockpile Inventory.
 */

export type StoresWorkspace =
  | "analytics"
  | "requisitions"
  | "waybills"
  | "scrap-returns"
  | "claims";

export interface StockpileItem {
  id: string;
  sku: string;
  item_name: string;
  category: "CONDUCTORS" | "INSULATORS" | "TRANSFORMER_OIL" | "SWITCHGEAR" | "SAFETY_PPE";
  current_quantity: number;
  reorder_threshold: number;
  unit_of_measure: string;
  unit_cost_ngn: number;
  bin_location: string;
  reorder_status: "HEALTHY" | "REORDER_TRIGGERED" | "OUT_OF_STOCK";
}

export interface MaterialRequisitionRecord {
  id: string;
  requisition_no: string;
  item_sku: string;
  item_name: string;
  quantity_requested: number;
  department: string;
  purpose: string;
  path_type: "PATH_A_IN_STOCK" | "PATH_B_MD_PROCUREMENT";
  estimated_procurement_cost_ngn?: number;
  status: "ALLOCATED_WAYBILLED" | "PENDING_STOCK_CHECK" | "AWAITING_MD_APPROVAL" | "PROCURED_RECEIVED";
  requester_name: string;
  created_at_wat: string;
}

export interface SupplierWaybillItem {
  id: string;
  waybill_number: string;
  supplier_name: string;
  item_description: string;
  quantity_delivered: number;
  po_reference: string;
  inspection_status: "PASSED_QC" | "PENDING_INSPECTION" | "REJECTED_DEFECTIVE";
  store_bin_assigned: string;
  received_by_staff_id: string;
  delivered_date_wat: string;
}

export interface ScrapReturnRecord {
  id: string;
  return_ticket_no: string;
  source_department: string;
  asset_type: "BURNT_DISTRIBUTION_TRANSFORMER" | "SCRAP_COPPER_COILS" | "CONDEMNED_ALUMINUM_CABLE" | "DIRTY_OIL_SLUDGE";
  weight_kg: number;
  estimated_salvage_value_ngn: number;
  salvage_yard_bin: string;
  status: "LOGGED_IN_YARD" | "AUDIT_CLEARED_FOR_AUCTION" | "AUCTIONED_DISPOSED";
  received_at_wat: string;
}

export interface StaffAllowanceClaimItem {
  id: string;
  claim_number: string;
  staff_name: string;
  staff_id: string;
  claim_type: "FIELD_OVERTIME" | "REMOTE_PLANT_TRAVEL" | "HIGH_VOLTAGE_EMERGENCY";
  duty_roster_verified: boolean;
  safety_ptw_verified: boolean;
  field_gps_distance_m: number;
  gps_within_tolerance: boolean;
  amount_claimed_ngn: number;
  stores_verification_status: "VERIFIED_PASSED" | "FLAGGED_DISCREPANCY" | "PENDING_CHECK";
  created_at_wat: string;
}

