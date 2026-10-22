package domain

import (
	"testing"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

func TestCalculateBulkMW(t *testing.T) {
	// 1. Success
	total, err := CalculateBulkMW(25.4, 10.0)
	if err != nil || total != 35.4 {
		t.Fatalf("expected 35.4 MW, got %f, err: %v", total, err)
	}

	// 2. Negative Generation
	_, err = CalculateBulkMW(-1.0, 10.0)
	if err != ErrNegativeBulkMW {
		t.Errorf("expected ErrNegativeBulkMW on negative gen, got: %v", err)
	}

	// 3. Negative TCN
	_, err = CalculateBulkMW(10.0, -5.0)
	if err != ErrNegativeBulkMW {
		t.Errorf("expected ErrNegativeBulkMW on negative TCN, got: %v", err)
	}
}

func TestComputePriorityFeederRationing(t *testing.T) {
	feeders := []FeederStatus{
		{FeederCode: "FDR-RES-E", CurrentLoadMW: 5.0, PriorityRank: 5, TariffBand: s01domain.TariffBandE, State: FeederInService},
		{FeederCode: "FDR-IND-A", CurrentLoadMW: 10.0, PriorityRank: 1, TariffBand: s01domain.TariffBandA, State: FeederInService},
		{FeederCode: "FDR-COM-B", CurrentLoadMW: 8.0, PriorityRank: 2, TariffBand: s01domain.TariffBandB, State: FeederInService},
		{FeederCode: "FDR-PTW-C", CurrentLoadMW: 6.0, PriorityRank: 3, TariffBand: s01domain.TariffBandC, State: FeederOutOnPTW},
	}

	// Case 1: 20 MW available.
	// Rank 1 (10 MW) -> In Service (accum: 10)
	// Rank 2 (8 MW)  -> In Service (accum: 18)
	// Rank 3 (PTW)   -> Stays OutOnPTW
	// Rank 5 (5 MW)  -> Exceeds 20 MW (18 + 5 = 23 > 20) -> Load Shed!
	decisions := ComputePriorityFeederRationing(20.0, feeders)
	if len(decisions) != 4 {
		t.Fatalf("expected 4 decisions, got %d", len(decisions))
	}

	decisionMap := make(map[string]RationingDecision)
	for _, d := range decisions {
		decisionMap[d.FeederCode] = d
	}

	if decisionMap["FDR-IND-A"].TargetState != FeederInService {
		t.Errorf("expected Band A feeder to be InService, got %s", decisionMap["FDR-IND-A"].TargetState)
	}
	if decisionMap["FDR-COM-B"].TargetState != FeederInService {
		t.Errorf("expected Band B feeder to be InService, got %s", decisionMap["FDR-COM-B"].TargetState)
	}
	if decisionMap["FDR-PTW-C"].TargetState != FeederOutOnPTW {
		t.Errorf("expected PTW feeder to retain OutOnPTW, got %s", decisionMap["FDR-PTW-C"].TargetState)
	}
	if decisionMap["FDR-RES-E"].TargetState != FeederOutOnLoadShedding {
		t.Errorf("expected Band E feeder to be shed, got %s", decisionMap["FDR-RES-E"].TargetState)
	}
}

func TestSafetyPermitPTW_5StepProtocol(t *testing.T) {
	// Step 1: Validation
	if err := ValidatePTWRequest("asset-1", "unit-head-1"); err != nil {
		t.Fatalf("expected valid PTW request, got: %v", err)
	}
	if err := ValidatePTWRequest("", "unit-head-1"); err != ErrInvalidPTWParams {
		t.Errorf("expected ErrInvalidPTWParams on empty asset")
	}

	ptw := SafetyPermitPTW{
		ID:         "ptw-1",
		PermitNo:   "PTW-20261022-0001",
		AssetID:    "asset-1",
		UnitHeadID: "unit-head-1",
		Status:     PTWRequested,
	}

	// Step 2a: Dept Head approves
	status, err := ValidateDualApproval(&ptw, true, false, "dept-head-1")
	if err != nil || status != PTWRequested {
		t.Fatalf("expected status to stay PTWRequested after only Dept Head approves, got %s", status)
	}

	// Step 2b: Safety Rep approves -> Transitions to APPROVED_PENDING_ISOLATION
	status, err = ValidateDualApproval(&ptw, false, true, "safety-rep-1")
	if err != nil || status != PTWApprovedPendingIsolation {
		t.Fatalf("expected status PTWApprovedPendingIsolation after dual approval, got %s", status)
	}

	// Step 2 Negative: Approval from invalid status
	badPTW := ptw
	badPTW.Status = PTWIsolatedEarthed
	if _, err := ValidateDualApproval(&badPTW, true, true, "dh-1"); err == nil {
		t.Errorf("expected error approving already isolated permit")
	}

	// Step 3: Substation Isolation
	// Missing proof notes
	_, err = ValidateSubstationIsolation(ptw.Status, "op-1", "short")
	if err != ErrMissingIsolationProof {
		t.Errorf("expected ErrMissingIsolationProof on short notes, got: %v", err)
	}

	// Success isolation
	status, err = ValidateSubstationIsolation(ptw.Status, "op-1", "Opened 33kV circuit breaker, racked out, and closed busbar earth switches.")
	if err != nil || status != PTWIsolatedEarthed {
		t.Fatalf("expected PTWIsolatedEarthed, got %s, err: %v", status, err)
	}
	ptw.Status = status

	// Step 3 Negative: Isolation from wrong status
	if _, err := ValidateSubstationIsolation(PTWRequested, "op-1", "Valid notes here 12345"); err == nil {
		t.Errorf("expected error isolating unapproved permit")
	}

	// Step 5a: Safety Clearance
	// Missing clearance notes
	_, err = ValidateSafetyClearance(ptw.Status, "safety-rep-1", "short")
	if err != ErrMissingClearanceNotes {
		t.Errorf("expected ErrMissingClearanceNotes, got: %v", err)
	}

	// Success clearance
	status, err = ValidateSafetyClearance(ptw.Status, "safety-rep-1", "Crew clear, all earthing clusters removed, line safe to energize.")
	if err != nil || status != PTWSafetyCleared {
		t.Fatalf("expected PTWSafetyCleared, got %s, err: %v", status, err)
	}
	ptw.Status = status

	// Step 5b: Restoration & Energization
	status, err = ValidateRestoration(ptw.Status, "op-1")
	if err != nil || status != PTWRestoredClosed {
		t.Fatalf("expected PTWRestoredClosed, got %s, err: %v", status, err)
	}

	// Step 5b Negative: Restoration from un-cleared status
	if _, err := ValidateRestoration(PTWIsolatedEarthed, "op-1"); err == nil {
		t.Errorf("expected error restoring without safety clearance")
	}
}

