// Package service_test exercises the 5-step Safety Permit-to-Work (PTW) lifecycle,
// ensuring strict dual approval, physical isolation, safety clearance, and line restoration.
package service_test

import (
	"context"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/service"
)

func TestWorkflowService_PTW_5StepProtocol(t *testing.T) {
	svc, repo, gridRepo, _ := setupService()
	ctx := context.Background()

	// Seed Grid Asset for Feeder
	gridAsset := s01domain.GridAsset{
		ID:               "feeder-asset-33kv-01",
		AssetCode:        "FDR-33KV-BUKURU",
		AssetType:        s01domain.AssetFeeder11KV,
		OrgUnitID:        "org-unit-1",
		OperationalState: s01domain.StateInService,
	}
	_ = gridRepo.CreateGridAsset(ctx, gridAsset)

	// Seed Feeder Status in Repository
	_ = repo.UpsertFeederStatus(ctx, domain.FeederStatus{
		FeederAssetID: gridAsset.ID,
		FeederCode:    gridAsset.AssetCode,
		PriorityRank:  1,
		TariffBand:    s01domain.TariffBandA,
		CurrentLoadMW: 6.0,
		State:         domain.FeederInService,
		UpdatedAt:     time.Now().UTC(),
	})

	// Step 1: Initiate PTW Request
	initReq := service.InitiatePTWRequest{
		AssetID:         gridAsset.ID,
		UnitHeadID:      "user-unit-head-01",
		WorkDescription: "Replace cracked surge arrestor on Phase B",
	}

	ptw, err := svc.InitiatePTW(ctx, initReq)
	if err != nil {
		t.Fatalf("unexpected error initiating PTW: %v", err)
	}
	if ptw.Status != domain.PTWRequested {
		t.Fatalf("expected status REQUESTED, got %s", ptw.Status)
	}
	if ptw.PermitNo == "" {
		t.Fatal("expected generated permit number")
	}

	// Negative initiation: empty asset ID
	if _, err := svc.InitiatePTW(ctx, service.InitiatePTWRequest{UnitHeadID: "head-1"}); err != domain.ErrInvalidPTWParams {
		t.Fatalf("expected ErrInvalidPTWParams on empty asset, got %v", err)
	}

	// Negative initiation: non-existent asset
	if _, err := svc.InitiatePTW(ctx, service.InitiatePTWRequest{AssetID: "non-existent-asset", UnitHeadID: "head-1"}); err == nil {
		t.Fatal("expected error on non-existent grid asset")
	}

	// Step 2: Dual Approval
	// Negative: unauthorized role
	if _, err := svc.DualApprovePTW(ctx, ptw.ID, "ACCOUNTS_OFFICER", "user-acc"); err != service.ErrUnauthorizedApprover {
		t.Fatalf("expected ErrUnauthorizedApprover, got %v", err)
	}

	// Approval 1: Head of Services signs
	pPartial, err := svc.DualApprovePTW(ctx, ptw.ID, "HEAD_SERVICES", "user-dept-head")
	if err != nil {
		t.Fatalf("unexpected error on Dept Head approval: %v", err)
	}
	if !pPartial.DeptHeadApproved || pPartial.SafetyRepApproved {
		t.Fatalf("mismatched approval flags: dept=%v, safety=%v", pPartial.DeptHeadApproved, pPartial.SafetyRepApproved)
	}
	if pPartial.Status != domain.PTWRequested {
		t.Fatalf("permit should remain in REQUESTED until dual approval, got: %s", pPartial.Status)
	}

	// Approval 2: HSE Safety Representative signs
	pDual, err := svc.DualApprovePTW(ctx, ptw.ID, "SAFETY_REPRESENTATIVE_HSE", "user-safety-rep")
	if err != nil {
		t.Fatalf("unexpected error on Safety Rep approval: %v", err)
	}
	if !pDual.DeptHeadApproved || !pDual.SafetyRepApproved {
		t.Fatal("both approval flags must be true after dual signoff")
	}
	if pDual.Status != domain.PTWApprovedPendingIsolation {
		t.Fatalf("expected APPROVED_PENDING_ISOLATION, got %s", pDual.Status)
	}

	// Step 3: Substation Isolation
	// Negative: isolation notes too short (< 10 chars)
	if _, err := svc.GrantSubstationIsolation(ctx, ptw.ID, "op-danladi", "done"); err != domain.ErrMissingIsolationProof {
		t.Fatalf("expected ErrMissingIsolationProof for short notes, got %v", err)
	}

	// Valid Substation Isolation
	isolationNotes := "Breaker racked out, red caution tags applied, 3-phase line earthed."
	pIsolated, err := svc.GrantSubstationIsolation(ctx, ptw.ID, "op-danladi", isolationNotes)
	if err != nil {
		t.Fatalf("unexpected error granting isolation: %v", err)
	}
	if pIsolated.Status != domain.PTWIsolatedEarthed {
		t.Fatalf("expected ISOLATED_EARTHED, got %s", pIsolated.Status)
	}

	// Verify Feeder state automatically shifted to OUT_ON_PERMIT_PTW
	fdrCheck, err := repo.GetFeederStatusByAssetID(ctx, gridAsset.ID)
	if err != nil || fdrCheck.State != domain.FeederOutOnPTW || *fdrCheck.ActivePTWID != ptw.ID {
		t.Fatalf("feeder state was not properly transitioned to OUT_ON_PERMIT_PTW: %+v", fdrCheck)
	}

	// Out of order attempt: cannot restore before safety cleared
	if _, err := svc.RestoreAndClosePTW(ctx, ptw.ID, "op-danladi"); err == nil {
		t.Fatal("expected error attempting to restore permit before safety clearance")
	}

	// Step 4 / 5a: HSE Safety Clearance
	// Negative: clearance notes too short (< 10 chars)
	if _, err := svc.SignSafetyClearance(ctx, ptw.ID, "user-safety-rep", "ok"); err != domain.ErrMissingClearanceNotes {
		t.Fatalf("expected ErrMissingClearanceNotes, got %v", err)
	}

	clearanceNotes := "All maintenance crew accounted for. Earths removed and line declared safe."
	pCleared, err := svc.SignSafetyClearance(ctx, ptw.ID, "user-safety-rep", clearanceNotes)
	if err != nil {
		t.Fatalf("unexpected error signing safety clearance: %v", err)
	}
	if pCleared.Status != domain.PTWSafetyCleared {
		t.Fatalf("expected SAFETY_CLEARED, got %s", pCleared.Status)
	}

	// Step 5b: Restoration & Closure
	pRestored, err := svc.RestoreAndClosePTW(ctx, ptw.ID, "op-danladi")
	if err != nil {
		t.Fatalf("unexpected error restoring PTW: %v", err)
	}
	if pRestored.Status != domain.PTWRestoredClosed {
		t.Fatalf("expected RESTORED_CLOSED, got %s", pRestored.Status)
	}

	// Verify Feeder status reverted to IN_SERVICE and ActivePTWID cleared
	fdrRestored, err := repo.GetFeederStatusByAssetID(ctx, gridAsset.ID)
	if err != nil || fdrRestored.State != domain.FeederInService || fdrRestored.ActivePTWID != nil {
		t.Fatalf("feeder state failed to revert to IN_SERVICE: %+v", fdrRestored)
	}

	// Query helpers
	byID, err := svc.GetPTWByID(ctx, ptw.ID)
	if err != nil || byID.ID != ptw.ID {
		t.Fatalf("GetPTWByID failed: %v", err)
	}

	byNo, err := svc.GetPTWByPermitNo(ctx, ptw.PermitNo)
	if err != nil || byNo.PermitNo != ptw.PermitNo {
		t.Fatalf("GetPTWByPermitNo failed: %v", err)
	}

	list, err := svc.ListPTWs(ctx, domain.PTWRestoredClosed, gridAsset.ID)
	if err != nil || len(list) != 1 {
		t.Fatalf("ListPTWs failed: %v, count: %d", err, len(list))
	}
}
