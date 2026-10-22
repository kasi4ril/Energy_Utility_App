/**
 * Executive Command & Predictive Maintenance Engine API client (S09).
 */

import { apiClient } from "../../lib/api-client";

export interface DashboardSummaryDTO {
  timestamp: string;
  total_generation_mw: number;
  tcn_grid_intake_mw: number;
  total_available_power_mw: number;
  feeders_summary: {
    total_feeders: number;
    in_service_count: number;
    out_on_fault_count: number;
    out_on_ptw_count: number;
    out_on_rationing_count: number;
  };
  revenue_summary: {
    total_collected_amount: number;
    prepaid_amount: number;
    postpaid_amount: number;
    channel_breakdown: Record<string, number>;
  };
  theft_alerts_open_count: number;
  churn_alerts_open_count: number;
  audit_pending_procurements_count: number;
  unreturned_scrap_items_count: number;
}

export interface ActionItemDTO {
  id: string;
  action_type: string;
  target_role: string;
  source_entity_id: string;
  source_entity_ref: string;
  source_department: string;
  summary: string;
  amount?: number;
  priority: string;
  status: string;
  created_at: string;
}

export interface PredictiveTriggerDTO {
  id: string;
  trigger_no: string;
  asset_id: string;
  asset_type: string;
  model_name: string;
  risk_score: number;
  risk_level: string;
  recommended_action: string;
  auto_work_order_created: boolean;
  linked_work_order_id?: string;
  status: string;
}

export async function fetchDashboardSummary(): Promise<DashboardSummaryDTO> {
  return apiClient.get<DashboardSummaryDTO>("/executive/dashboard");
}

export async function fetchActionItems(role?: string): Promise<ActionItemDTO[]> {
  try {
    const url = role ? `/executive/actions?target_role=${role}&status=PENDING` : "/executive/actions?status=PENDING";
    return await apiClient.get<ActionItemDTO[]>(url);
  } catch {
    return [];
  }
}

export async function decideActionItem(
  id: string,
  staffId: string,
  approve: boolean,
  notes: string
): Promise<any> {
  return apiClient.post(`/executive/actions/${id}/decide`, {
    decided_by_staff_id: staffId,
    approve,
    notes,
  });
}

export async function fetchPredictiveTriggers(): Promise<PredictiveTriggerDTO[]> {
  try {
    return await apiClient.get<PredictiveTriggerDTO[]>("/executive/predictive/triggers");
  } catch {
    return [];
  }
}

