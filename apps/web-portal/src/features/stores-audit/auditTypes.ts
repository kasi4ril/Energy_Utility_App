/**
 * Types & Domain Interfaces for Internal Audit & Compliance.
 */

export type AuditWorkspace =
  | "analytics"
  | "ledger"
  | "claims-audit"
  | "scrap-audit"
  | "compliance";

export interface AuditFindingItem {
  id: string;
  audit_ref: string;
  department: string;
  finding_title: string;
  risk_level: "HIGH" | "MEDIUM" | "LOW";
  variance_amount_ngn?: number;
  status: "OPEN_OBSERVATION" | "MANAGEMENT_RESPONSE_RECEIVED" | "RECTIFIED_CLOSED";
  auditor_staff_id: string;
  created_at_wat: string;
}

export interface ChainedLedgerRecord {
  id: string;
  entity_type: "work_order" | "safety_permit" | "requisition" | "claim";
  entity_id: string;
  sequence_number: number;
  author_staff_id: string;
  author_role: string;
  note_content: string;
  current_sha256_hash: string;
  prev_sha256_hash: string;
  integrity_valid: boolean;
  timestamp_wat: string;
}

export interface FinalClaimAuditItem {
  id: string;
  claim_number: string;
  staff_name: string;
  staff_id: string;
  claim_type: string;
  amount_claimed_ngn: number;
  triangulation_status: "PASSED_ROSTER_PTW_GPS" | "FLAGGED_DISTANCE_MISMATCH";
  audit_decision: "AUDIT_APPROVED" | "FLAGGED_FOR_HEARING" | "PENDING_AUDITOR";
  audit_remarks?: string;
  created_at_wat: string;
}

export interface ScrapAuctionItem {
  id: string;
  lot_number: string;
  asset_type: string;
  weighed_weight_kg: number;
  reserve_price_ngn: number;
  highest_bid_ngn: number;
  audit_clearance: "CLEARED_FOR_AUCTION" | "HELD_FOR_REVALUATION" | "AUCTION_COMPLETED";
  cleared_by_auditor_id?: string;
}

export interface NERCComplianceMetric {
  standard_code: string;
  title: string;
  target_sla: string;
  actual_performance: string;
  compliance_status: "COMPLIANT" | "WARNING" | "BREACH";
}

