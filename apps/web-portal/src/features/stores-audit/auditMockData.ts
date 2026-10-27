/**
 * Mock data for Internal Audit & Compliance Desk.
 */

import {
  AuditFindingItem,
  ChainedLedgerRecord,
  FinalClaimAuditItem,
  ScrapAuctionItem,
  NERCComplianceMetric,
} from "./auditTypes";

export const MOCK_AUDIT_FINDINGS: AuditFindingItem[] = [
  {
    id: "FND-01",
    audit_ref: "AUD-2026-Q3-014",
    department: "Commercial Operations",
    finding_title: "CT Ratio Multiplier Discrepancy at Anglo-Jos Industrial Feeder",
    risk_level: "HIGH",
    variance_amount_ngn: 450000,
    status: "OPEN_OBSERVATION",
    auditor_staff_id: "STAFF-AUD-01",
    created_at_wat: "2026-10-27 10:15 WAT",
  },
  {
    id: "FND-02",
    audit_ref: "AUD-2026-Q3-012",
    department: "Central Stores",
    finding_title: "150mm² ACSR Conductor Physical Count Variance (50m difference)",
    risk_level: "MEDIUM",
    variance_amount_ngn: 225000,
    status: "MANAGEMENT_RESPONSE_RECEIVED",
    auditor_staff_id: "STAFF-AUD-02",
    created_at_wat: "2026-10-25 14:00 WAT",
  },
  {
    id: "FND-03",
    audit_ref: "AUD-2026-Q3-009",
    department: "Generation Master",
    finding_title: "Plant 2 Turbine Oil Consumption Log Reconciliation",
    risk_level: "LOW",
    status: "RECTIFIED_CLOSED",
    auditor_staff_id: "STAFF-AUD-01",
    created_at_wat: "2026-10-22 09:30 WAT",
  },
];

export const MOCK_CHAINED_LEDGER: ChainedLedgerRecord[] = [
  {
    id: "WN-8891",
    entity_type: "work_order",
    entity_id: "WO-2026-081",
    sequence_number: 1,
    author_staff_id: "STAFF-ENG-04",
    author_role: "UNIT_HEAD_FEEDER_MAINT",
    note_content: "Work order created for emergency insulator replacement along Anglo-Jos corridor.",
    current_sha256_hash: "a4f8e912c3b4a5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0",
    prev_sha256_hash: "0000000000000000000000000000000000000000000000000000000000000000",
    integrity_valid: true,
    timestamp_wat: "2026-10-27 14:15 WAT",
  },
  {
    id: "WN-8892",
    entity_type: "work_order",
    entity_id: "WO-2026-081",
    sequence_number: 2,
    author_staff_id: "STAFF-DSP-01",
    author_role: "CHIEF_DISPATCH_OFFICER",
    note_content: "Approved 33kV line isolation with PTW-041-33KV. Line tested dead.",
    current_sha256_hash: "b5c9f023d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1",
    prev_sha256_hash: "a4f8e912c3b4a5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0",
    integrity_valid: true,
    timestamp_wat: "2026-10-27 14:25 WAT",
  },
  {
    id: "WN-8893",
    entity_type: "work_order",
    entity_id: "WO-2026-081",
    sequence_number: 3,
    author_staff_id: "STAFF-ENG-04",
    author_role: "UNIT_HEAD_FEEDER_MAINT",
    note_content: "All 6 disc insulators replaced. Proof of work photos uploaded and GPS validated at 8.4m.",
    current_sha256_hash: "c6d0a134e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
    prev_sha256_hash: "b5c9f023d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1",
    integrity_valid: true,
    timestamp_wat: "2026-10-27 16:30 WAT",
  },
];

export const MOCK_FINAL_CLAIMS: FinalClaimAuditItem[] = [
  {
    id: "FCLM-01",
    claim_number: "CLM-2026-112",
    staff_name: "Iliya Pam (Field Lineman)",
    staff_id: "STAFF-ENG-12",
    claim_type: "HIGH_VOLTAGE_EMERGENCY",
    amount_claimed_ngn: 35000,
    triangulation_status: "PASSED_ROSTER_PTW_GPS",
    audit_decision: "PENDING_AUDITOR",
    created_at_wat: "2026-10-27 15:40 WAT",
  },
  {
    id: "FCLM-02",
    claim_number: "CLM-2026-111",
    staff_name: "Gyang Dung (Linesman Lead)",
    staff_id: "STAFF-ENG-14",
    claim_type: "FIELD_OVERTIME",
    amount_claimed_ngn: 45000,
    triangulation_status: "FLAGGED_DISTANCE_MISMATCH",
    audit_decision: "FLAGGED_FOR_HEARING",
    audit_remarks: "GPS captured at 620m away from asset coordinates. Roster was confirmed but location failed tolerance.",
    created_at_wat: "2026-10-26 18:20 WAT",
  },
];

export const MOCK_SCRAP_AUCTION: ScrapAuctionItem[] = [
  {
    id: "LOT-01",
    lot_number: "LOT-SCR-2026-04",
    asset_type: "Burnt 2.5MVA Transformer Tank & Scrap Copper",
    weighed_weight_kg: 1850,
    reserve_price_ngn: 2400000,
    highest_bid_ngn: 2750000,
    audit_clearance: "CLEARED_FOR_AUCTION",
    cleared_by_auditor_id: "STAFF-AUD-01",
  },
  {
    id: "LOT-02",
    lot_number: "LOT-SCR-2026-05",
    asset_type: "Condemned ACSR Aluminum Line Scrap",
    weighed_weight_kg: 420,
    reserve_price_ngn: 580000,
    highest_bid_ngn: 640000,
    audit_clearance: "CLEARED_FOR_AUCTION",
    cleared_by_auditor_id: "STAFF-AUD-01",
  },
];

export const MOCK_NERC_COMPLIANCE: NERCComplianceMetric[] = [
  {
    standard_code: "NERC-SLA-01",
    title: "Fault Outage Resolution SLA (Feeder Level)",
    target_sla: "< 4.0 Hours MTTR",
    actual_performance: "2.1 Hours (Average)",
    compliance_status: "COMPLIANT",
  },
  {
    standard_code: "NERC-MYTO-02",
    title: "Estimated Billing Cap Order Compliance",
    target_sla: "100% Zero-Estimation Vending",
    actual_performance: "91.2% Prepaid Metered",
    compliance_status: "COMPLIANT",
  },
  {
    standard_code: "NEMSA-SAF-04",
    title: "Permit-to-Work Safety Isolation Lock",
    target_sla: "100% 5-Step Lockout Compliance",
    actual_performance: "100% Compliant",
    compliance_status: "COMPLIANT",
  },
  {
    standard_code: "NERC-DIS-03",
    title: "Commercial Dispute Resolution SLA",
    target_sla: "< 72 Hours from Customer Complaint",
    actual_performance: "28 Hours",
    compliance_status: "COMPLIANT",
  },
];

