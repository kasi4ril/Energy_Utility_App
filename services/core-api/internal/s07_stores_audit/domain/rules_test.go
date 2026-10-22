// Package domain_test verifies store inventory rules, dual-path procurement state machines,
// and automated 3-way staff claims verification invariants.
package domain_test

import (
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
)

func TestStoreItem_Validation(t *testing.T) {
	valid := domain.StoreItem{
		ItemCode:          "CBL-33KV-XLPE",
		ItemName:          "33kV XLPE Underground Cable 150sqmm",
		Category:          domain.CategoryCablesConductors,
		QuantityAvailable: 1500.0,
		ReorderThreshold:  200.0,
		UnitCost:          45000.0,
	}

	if err := domain.ValidateStoreItem(valid); err != nil {
		t.Fatalf("expected valid item to pass, got: %v", err)
	}

	// Empty Code
	invalidCode := valid
	invalidCode.ItemCode = ""
	if err := domain.ValidateStoreItem(invalidCode); err != domain.ErrInvalidItemData {
		t.Errorf("expected ErrInvalidItemData on empty code, got: %v", err)
	}

	// Negative Quantity
	negQty := valid
	negQty.QuantityAvailable = -5.0
	if err := domain.ValidateStoreItem(negQty); err != domain.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got: %v", err)
	}

	// Reorder evaluation
	if !domain.EvaluateReorderStatus(100.0, 200.0) {
		t.Error("expected 100 <= 200 to trigger reorder status")
	}
	if domain.EvaluateReorderStatus(500.0, 200.0) {
		t.Error("expected 500 > 200 not to trigger reorder status")
	}
}

func TestCalculateRequisitionTotal(t *testing.T) {
	items := []domain.RequisitionItem{
		{ItemCode: "FUSE-100A", QuantityRequested: 10, UnitCost: 2500.0},
		{ItemCode: "INS-PIN-33", QuantityRequested: 4, UnitCost: 15000.0},
	}

	total, err := domain.CalculateRequisitionTotal(items)
	if err != nil {
		t.Fatalf("unexpected error calculating total: %v", err)
	}
	if total != 85000.0 {
		t.Fatalf("expected total 85000.0, got: %f", total)
	}

	// Empty items
	if _, err := domain.CalculateRequisitionTotal([]domain.RequisitionItem{}); err != domain.ErrEmptyRequisitionItems {
		t.Errorf("expected ErrEmptyRequisitionItems, got: %v", err)
	}

	// Negative quantity
	negItems := []domain.RequisitionItem{{ItemCode: "A", QuantityRequested: -1, UnitCost: 100}}
	if _, err := domain.CalculateRequisitionTotal(negItems); err == nil {
		t.Error("expected error on negative item quantity")
	}
}

func TestRequisition_DualPathWorkflow(t *testing.T) {
	// --- PATH A: In-Stock Requisition Flow ---
	reqInStock := domain.MaterialRequisition{
		ID:          "req-01",
		ReqNo:       "REQ-2026-0001",
		InStockFlag: true,
		Status:      domain.StatusRequisitionSubmitted,
	}

	// Step 2: Dept Head Approves
	if err := domain.ValidateDeptApproval(&reqInStock, "dept-head-1"); err != nil {
		t.Fatalf("unexpected error on dept approval: %v", err)
	}
	if reqInStock.Status != domain.StatusRequisitionDeptApproved {
		t.Fatalf("expected DEPT_APPROVED, got: %s", reqInStock.Status)
	}

	// Step 3: Internal Audit Consents
	if err := domain.ValidateAuditConsent(&reqInStock, "auditor-1", true, "Verified with Work Order #42"); err != nil {
		t.Fatalf("unexpected error on audit consent: %v", err)
	}
	if reqInStock.Status != domain.StatusRequisitionAuditConsented {
		t.Fatalf("expected AUDIT_CONSENTED, got: %s", reqInStock.Status)
	}

	// Step 4: Store Issuance with Waybill
	if err := domain.ValidateMaterialIssuance(&reqInStock, "store-op-1", "WAYBILL-2026-0001"); err != nil {
		t.Fatalf("unexpected error on issuance: %v", err)
	}
	if reqInStock.Status != domain.StatusRequisitionIssued {
		t.Fatalf("expected ISSUED, got: %s", reqInStock.Status)
	}

	// --- PATH B: Out-of-Stock / MD Procurement Flow ---
	reqOutStock := domain.MaterialRequisition{
		ID:          "req-02",
		ReqNo:       "REQ-2026-0002",
		InStockFlag: false,
		Status:      domain.StatusRequisitionSubmitted,
	}

	// Step 2: Dept Head Approves -> Transitions to MD_PROCUREMENT_PENDING
	if err := domain.ValidateDeptApproval(&reqOutStock, "dept-head-1"); err != nil {
		t.Fatalf("unexpected error on out-of-stock dept approval: %v", err)
	}
	if reqOutStock.Status != domain.StatusRequisitionMDProcurementPending {
		t.Fatalf("expected MD_PROCUREMENT_PENDING, got: %s", reqOutStock.Status)
	}

	// Negative: MD cannot approve before Audit Consent
	if err := domain.ValidateMDProcurementApproval(&reqOutStock, "Approved"); err != domain.ErrAuditConsentRequired {
		t.Fatalf("expected ErrAuditConsentRequired, got: %v", err)
	}

	// Step 3a: Internal Audit Consents pricing
	if err := domain.ValidateAuditConsent(&reqOutStock, "auditor-1", true, "Quotes checked and valid"); err != nil {
		t.Fatalf("unexpected error on out-of-stock audit consent: %v", err)
	}
	if reqOutStock.AuditConsentStatus != domain.AuditConsentConsented {
		t.Fatalf("expected CONSENTED audit status, got: %s", reqOutStock.AuditConsentStatus)
	}

	// Step 3b: MD Grants Procurement Approval
	if err := domain.ValidateMDProcurementApproval(&reqOutStock, "MD Signed off for supplier PO"); err != nil {
		t.Fatalf("unexpected error on MD approval: %v", err)
	}
	if reqOutStock.Status != domain.StatusRequisitionMDApproved {
		t.Fatalf("expected MD_APPROVED, got: %s", reqOutStock.Status)
	}

	// Step 4: Store Officer Receives and Issues
	if err := domain.ValidateMaterialIssuance(&reqOutStock, "store-op-1", "WAYBILL-2026-0002"); err != nil {
		t.Fatalf("unexpected error on out-of-stock issuance: %v", err)
	}
	if reqOutStock.Status != domain.StatusRequisitionIssued {
		t.Fatalf("expected ISSUED, got: %s", reqOutStock.Status)
	}
}

func TestRequisition_AuditRejection(t *testing.T) {
	req := domain.MaterialRequisition{
		ID:          "req-reject",
		InStockFlag: true,
		Status:      domain.StatusRequisitionDeptApproved,
	}

	if err := domain.ValidateAuditConsent(&req, "auditor-1", false, "Missing maintenance proof"); err != nil {
		t.Fatalf("unexpected error rejecting: %v", err)
	}
	if req.Status != domain.StatusRequisitionRejected || req.AuditConsentStatus != domain.AuditConsentQueriedRejected {
		t.Fatalf("expected REJECTED_QUERIED, got status=%s, consent=%s", req.Status, req.AuditConsentStatus)
	}

	// Attempt issuance on rejected requisition should fail
	if err := domain.ValidateMaterialIssuance(&req, "store-1", "WAY-1"); err == nil {
		t.Error("expected error attempting to issue rejected requisition")
	}
}

func TestStaffAllowanceClaim_3WayVerification(t *testing.T) {
	ptwID := "ptw-2026-01"
	claim := domain.StaffAllowanceClaim{
		ClaimNo:             "CLM-001",
		StaffID:             "staff-101",
		ClaimType:           domain.ClaimOvertime,
		LinkedPTWOrTicketID: &ptwID,
	}

	// 1. Roster Mismatch -> High Risk Fraud (0.950)
	flag1, score1, stat1 := domain.Evaluate3WayClaimVerification(claim, false, true, true)
	if flag1 != domain.AIFlagFlaggedFraud || score1 < 0.90 || stat1 != domain.StatusClaimAIFlagged {
		t.Errorf("expected FLAGGED_FRAUD, got flag=%s, score=%f, status=%s", flag1, score1, stat1)
	}

	// 2. Roster matches, but PTW/GPS fails -> FLAGGED_MISMATCH (0.720)
	flag2, score2, stat2 := domain.Evaluate3WayClaimVerification(claim, true, false, true)
	if flag2 != domain.AIFlagFlaggedMismatch || score2 < 0.70 || stat2 != domain.StatusClaimAIFlagged {
		t.Errorf("expected FLAGGED_MISMATCH on ptw failure, got flag=%s, score=%f", flag2, score2)
	}

	// 3. All 3 match -> VERIFIED (0.050)
	flag3, score3, stat3 := domain.Evaluate3WayClaimVerification(claim, true, true, true)
	if flag3 != domain.AIFlagVerified || score3 > 0.10 || stat3 != domain.StatusClaimAIVerified {
		t.Errorf("expected VERIFIED, got flag=%s, score=%f, status=%s", flag3, score3, stat3)
	}

	// 4. Auditor Decision: Approve
	cApprove := claim
	cApprove.Status = domain.StatusClaimAIVerified
	if err := domain.ValidateAuditorClaimDecision(&cApprove, "auditor-1", true, "Verified and passed"); err != nil {
		t.Fatalf("unexpected error deciding claim: %v", err)
	}
	if cApprove.Status != domain.StatusClaimAuditApproved {
		t.Errorf("expected AUDIT_APPROVED, got: %s", cApprove.Status)
	}

	// 5. Auditor Decision: Reject
	cReject := claim
	cReject.Status = domain.StatusClaimAIFlagged
	if err := domain.ValidateAuditorClaimDecision(&cReject, "auditor-1", false, "Overtime not logged in shift roster"); err != nil {
		t.Fatalf("unexpected error deciding claim: %v", err)
	}
	if cReject.Status != domain.StatusClaimAuditRejected {
		t.Errorf("expected AUDIT_REJECTED, got: %s", cReject.Status)
	}
}

