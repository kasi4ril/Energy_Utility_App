/**
 * Stores, Internal Audit & MD Procurement API client (S07).
 */

import { apiClient } from "../../lib/api-client";

export interface StoreItemDTO {
  id: string;
  sku: string;
  item_name: string;
  category: string;
  current_quantity: number;
  reorder_threshold: number;
  unit_cost: number;
}

export interface RequisitionDTO {
  id: string;
  requisition_no: string;
  item_sku: string;
  quantity_requested: number;
  is_out_of_stock_procurement: boolean;
  procurement_estimated_cost?: number;
  status: string;
  requester_staff_id: string;
  unit_head_approval: boolean;
  internal_audit_consent: boolean;
  md_approval?: boolean;
}

export async function searchInventory(query: string = ""): Promise<StoreItemDTO[]> {
  try {
    return await apiClient.get<StoreItemDTO[]>(`/stores/inventory?q=${query}`);
  } catch {
    return [];
  }
}

export async function fetchRequisitions(): Promise<RequisitionDTO[]> {
  try {
    return await apiClient.get<RequisitionDTO[]>("/stores/requisitions");
  } catch {
    return [];
  }
}

export async function createRequisition(payload: {
  item_sku: string;
  quantity: number;
  work_order_id?: string;
  purpose: string;
  requester_staff_id: string;
}): Promise<any> {
  return apiClient.post("/stores/requisitions", payload);
}

