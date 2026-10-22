/**
 * Services & Maintenance API client (S04).
 */

import { apiClient } from "../../lib/api-client";

export interface WorkOrderDTO {
  id: string;
  order_no: string;
  unit_type: string;
  work_order_type: string;
  priority: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  status: string;
  unit_head_id: string;
  team_lead_id: string;
  asset_id: string;
  ptw_required: boolean;
  instructions: string;
  created_at: string;
}

export async function fetchWorkOrders(unitType?: string): Promise<WorkOrderDTO[]> {
  try {
    const url = unitType ? `/maintenance/orders?unit_type=${unitType}` : "/maintenance/orders";
    return await apiClient.get<WorkOrderDTO[]>(url);
  } catch {
    return [];
  }
}

export async function verifyProofOfWork(orderId: string, approved: boolean, notes: string, staffId: string): Promise<any> {
  return apiClient.post(`/maintenance/orders/${orderId}/verify-proof`, {
    approved,
    review_notes: notes,
    reviewed_by_staff_id: staffId,
  });
}

