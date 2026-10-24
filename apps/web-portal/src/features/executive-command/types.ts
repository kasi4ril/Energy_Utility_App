/**
 * Executive Command Tier (MD & C-Suite) Type Definitions.
 * Strictly <= 350 lines.
 */

export type ExecutiveWorkspace =
  | "approvals"
  | "analytics"
  | "inter-department"
  | "governance"
  | "ledger";

export type DepartmentSource =
  | "GENERATION"
  | "ENGINEERING_SERVICES"
  | "COMMERCIAL_OPS"
  | "INTERNAL_AUDIT"
  | "STORES_PROCUREMENT"
  | "HSE_SAFETY";

export type RequisitionCategory =
  | "CAPEX_OVERHAUL"
  | "GRID_REPAIR"
  | "COMMERCIAL_REFUND"
  | "BULK_PROCUREMENT"
  | "SCRAP_DISPOSAL"
  | "SAFETY_EMERGENCY";

export type RequisitionUrgency =
  | "EMERGENCY_OUTAGE"
  | "URGENT"
  | "ROUTINE";

export type RequisitionStatus =
  | "PENDING_APPROVAL"
  | "UNDER_REVIEW"
  | "APPROVED"
  | "REJECTED"
  | "ACKNOWLEDGED_MINUTED";

export interface RequisitionLineItem {
  id: string;
  item: string;
  quantity: number;
  unit: string;
  unit_price: number;
  total_price: number;
}

export interface ExecutiveRequisition {
  id: string;
  request_no: string;
  title: string;
  source_department: DepartmentSource;
  submitting_head_name: string;
  submitting_head_role: string;
  amount: number;
  category: RequisitionCategory;
  urgency: RequisitionUrgency;
  submission_date: string;
  status: RequisitionStatus;
  summary: string;
  technical_justification: string;
  line_items: RequisitionLineItem[];
  md_review_notes?: string;
  md_review_date?: string;
  minuted_to_departments?: DepartmentSource[];
  minuted_notes?: string;
  decision_date?: string;
  decision_by?: string;
  hash_chain_seal?: string;
}

export interface ExecutiveInterDeptMemo {
  id: string;
  memo_no: string;
  title: string;
  originating_departments: DepartmentSource[];
  date: string;
  priority: "CRITICAL" | "HIGH" | "NORMAL";
  executive_summary: string;
  technical_impact: string;
  status: "AWAITING_MD_REVIEW" | "NOTED_AND_MINUTED" | "ACTION_DISPATCHED";
  md_minute_remarks?: string;
  assigned_action_heads?: string[];
  target_completion_date?: string;
}

export interface BoardResolution {
  id: string;
  resolution_no: string;
  title: string;
  meeting_date: string;
  status: "ENACTED" | "IN_PROGRESS" | "PENDING_RATIFICATION";
  mandate_summary: string;
  allocated_budget_ngn: number;
  responsible_executive: string;
}

export interface RegulatoryFiling {
  id: string;
  filing_code: string;
  regulatory_body: "NERC" | "NEMSA" | "TCN_SO" | "FMWR";
  title: string;
  due_date: string;
  status: "COMPLIANT" | "SUBMITTED" | "HEARING_SCHEDULED" | "EXPIRING_SOON";
  license_benchmark: string;
}

export interface StrategicConcession {
  id: string;
  partner_name: string;
  concession_scope: string;
  active_period: string;
  capacity_commitment_mw: number;
  quarterly_revenue_est_ngn: number;
  status: "ACTIVE_PERFORMING" | "UNDER_RENEGOTIATION" | "EXPIRING";
}

export interface DepartmentHeadScorecard {
  department: DepartmentSource;
  head_name: string;
  role: string;
  kpi_score_pct: number;
  pending_tasks_count: number;
  sla_adherence_pct: number;
  budget_utilization_pct: number;
  status: "EXEMPLARY" | "ON_TRACK" | "ATTENTION_REQUIRED";
}

