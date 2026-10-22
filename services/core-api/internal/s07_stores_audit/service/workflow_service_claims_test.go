// Package service_test exercises the 3-way staff allowance claims verification desk,
// automated fraud scoring, and internal auditor approval/rejection workflows for S07.
package service_test

import (
	"context"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
)

func TestWorkflowService_StaffAllowanceClaims(t *testing.T) {
	svc, _, _ := setupService()
	ctx := context.Background()

	ptwID := "ptw-2026-bukuru-42"

	// 1. Full 3-Way Match -> AI_VERIFIED
	claimClean := domain.StaffAllowanceClaim{
		StaffID:             "staff-101",
		ClaimType:           domain.ClaimOvertime,
		Amount:              25000.0,
		ClaimDate:           "2026-10-22",
		LinkedPTWOrTicketID: &ptwID,
		RosterMatchVerified: true,
		PTWCrewVerified:     true,
		GPSMatchVerified:    true,
	}

	createdClean, err := svc.SubmitAllowanceClaim(ctx, claimClean)
	if err != nil {
		t.Fatalf("unexpected error submitting clean claim: %v", err)
	}
	if createdClean.Status != domain.StatusClaimAIVerified || createdClean.OllamaAuditFlag != domain.AIFlagVerified {
		t.Fatalf("expected AI_VERIFIED claim, got status=%s, flag=%s", createdClean.Status, createdClean.OllamaAuditFlag)
	}
	if createdClean.OllamaFraudScore > 0.10 {
		t.Errorf("expected fraud score <= 0.10, got: %f", createdClean.OllamaFraudScore)
	}

	// Auditor approves clean claim
	approved, err := svc.ReviewClaimByAuditor(ctx, createdClean.ID, "auditor-01", true, "Roster, PTW, and GPS all verified")
	if err != nil || approved.Status != domain.StatusClaimAuditApproved {
		t.Fatalf("auditor approval failed: %+v, err: %v", approved, err)
	}

	// 2. Roster Mismatch -> AI_FLAGGED / FLAGGED_FRAUD (0.950)
	claimFraud := domain.StaffAllowanceClaim{
		StaffID:             "staff-ghost-01",
		ClaimType:           domain.ClaimOvertime,
		Amount:              35000.0,
		ClaimDate:           "2026-10-22",
		RosterMatchVerified: false, // staff was never on shift
		PTWCrewVerified:     true,
		GPSMatchVerified:    true,
	}

	createdFraud, err := svc.SubmitAllowanceClaim(ctx, claimFraud)
	if err != nil {
		t.Fatalf("unexpected error submitting claim: %v", err)
	}
	if createdFraud.Status != domain.StatusClaimAIFlagged || createdFraud.OllamaAuditFlag != domain.AIFlagFlaggedFraud {
		t.Fatalf("expected AI_FLAGGED / FLAGGED_FRAUD, got status=%s, flag=%s", createdFraud.Status, createdFraud.OllamaAuditFlag)
	}
	if createdFraud.OllamaFraudScore < 0.90 {
		t.Errorf("expected high fraud score >= 0.90, got: %f", createdFraud.OllamaFraudScore)
	}

	// Auditor rejects fraudulent claim
	rejected, err := svc.ReviewClaimByAuditor(ctx, createdFraud.ID, "auditor-01", false, "Staff was absent from shift roll-call on claim date")
	if err != nil || rejected.Status != domain.StatusClaimAuditRejected {
		t.Fatalf("auditor rejection failed: %+v, err: %v", rejected, err)
	}

	// 3. Negative: Invalid Claim Submission (negative amount or empty staff id)
	badClaim := domain.StaffAllowanceClaim{
		StaffID:   "",
		Amount:    -1000.0,
		ClaimDate: "2026-10-22",
	}
	if _, err := svc.SubmitAllowanceClaim(ctx, badClaim); err != domain.ErrInvalidClaimData {
		t.Fatalf("expected ErrInvalidClaimData, got: %v", err)
	}

	// 4. Query by ID and List
	byID, err := svc.GetClaimByID(ctx, createdClean.ID)
	if err != nil || byID.ClaimNo != createdClean.ClaimNo {
		t.Fatalf("GetClaimByID failed: %+v, err: %v", byID, err)
	}

	list, err := svc.ListClaims(ctx, "staff-101", domain.StatusClaimAuditApproved)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 approved claim, got: %d", len(list))
	}
}

