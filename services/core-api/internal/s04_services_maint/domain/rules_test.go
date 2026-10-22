package domain

import (
	"testing"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

func TestGetGPSToleranceForUnit(t *testing.T) {
	tests := []struct {
		unit     UnitType
		expected float64
	}{
		{UnitFeederMaint, 150.0},
		{UnitTransformerMaint, 75.0},
		{UnitMetering, 50.0},
		{UnitType("UNKNOWN"), 100.0},
	}

	for _, tc := range tests {
		got := GetGPSToleranceForUnit(tc.unit)
		if got != tc.expected {
			t.Errorf("GetGPSToleranceForUnit(%s) = %f; want %f", tc.unit, got, tc.expected)
		}
	}
}

func TestCalculateHaversineDistance(t *testing.T) {
	// Jos Central Substation (approx 9.9320, 8.8920) and a point ~111 meters North
	p1 := s01domain.GPSPoint{Latitude: 9.9320, Longitude: 8.8920}
	p2 := s01domain.GPSPoint{Latitude: 9.9330, Longitude: 8.8920}

	dist := CalculateHaversineDistance(p1, p2)
	if dist < 100.0 || dist > 120.0 {
		t.Fatalf("expected distance around 111m, got %f", dist)
	}

	// Same point distance should be zero
	zeroDist := CalculateHaversineDistance(p1, p1)
	if zeroDist != 0.0 {
		t.Fatalf("expected distance 0, got %f", zeroDist)
	}
}

func TestValidateWorkOrderCreation(t *testing.T) {
	validWO := MaintenanceWorkOrder{
		OrderNo:       "WO-20261022-0001",
		UnitType:      UnitFeederMaint,
		WorkOrderType: WOCorrectiveRepair,
		Priority:      PriorityHigh,
		UnitHeadID:    "unit-head-uuid",
		TeamLeadID:    "team-lead-uuid",
		AssetID:       "asset-uuid",
		Instructions:  "Inspect 33kV corridor between pole 40 and 42 for broken pin insulators.",
	}

	if err := ValidateWorkOrderCreation(validWO); err != nil {
		t.Fatalf("expected valid work order, got: %v", err)
	}

	// Test missing OrderNo
	badWO := validWO
	badWO.OrderNo = ""
	if err := ValidateWorkOrderCreation(badWO); err == nil {
		t.Errorf("expected error on empty OrderNo")
	}

	// Test missing UnitHeadID
	badWO = validWO
	badWO.UnitHeadID = ""
	if err := ValidateWorkOrderCreation(badWO); err == nil {
		t.Errorf("expected error on empty UnitHeadID")
	}

	// Test missing TeamLeadID
	badWO = validWO
	badWO.TeamLeadID = ""
	if err := ValidateWorkOrderCreation(badWO); err == nil {
		t.Errorf("expected error on empty TeamLeadID")
	}

	// Test missing AssetID
	badWO = validWO
	badWO.AssetID = ""
	if err := ValidateWorkOrderCreation(badWO); err == nil {
		t.Errorf("expected error on empty AssetID")
	}

	// Test short instructions
	badWO = validWO
	badWO.Instructions = "fix line"
	if err := ValidateWorkOrderCreation(badWO); err == nil {
		t.Errorf("expected error on short instructions")
	}

	// Test unknown unit type
	badWO = validWO
	badWO.UnitType = "CIVIL_WORKS"
	if err := ValidateWorkOrderCreation(badWO); err == nil {
		t.Errorf("expected error on invalid unit type")
	}
}

func TestValidateProofSubmission_HappyAndNegativePaths(t *testing.T) {
	assetGPS := s01domain.GPSPoint{Latitude: 9.9320, Longitude: 8.8920}
	closeGPS := s01domain.GPSPoint{Latitude: 9.9321, Longitude: 8.8920} // ~11m
	farGPS := s01domain.GPSPoint{Latitude: 9.9400, Longitude: 8.8920}   // ~890m away

	validBeforeHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	validAfterHash := "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"
	seal := "SEAL-998811"
	workDetails := "Replaced damaged 33kV surge arrestor and verified ground resistance at 2.4 ohms."

	// 1. Happy path - Feeder Maintenance within 150m
	dist, valid, err := ValidateProofSubmission(
		UnitFeederMaint, assetGPS, closeGPS,
		validBeforeHash, validAfterHash, nil,
		workDetails, false, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error on valid feeder proof: %v", err)
	}
	if !valid || dist > 20.0 {
		t.Fatalf("expected valid GPS proximity, got valid=%v, dist=%f", valid, dist)
	}

	// 2. Anti-Photo-Reuse: Identical hashes
	_, _, err = ValidateProofSubmission(
		UnitFeederMaint, assetGPS, closeGPS,
		validBeforeHash, validBeforeHash, nil,
		workDetails, false, nil,
	)
	if err != ErrIdenticalPhotoHashes {
		t.Errorf("expected ErrIdenticalPhotoHashes, got: %v", err)
	}

	// 3. Empty / Short Photo Hash
	_, _, err = ValidateProofSubmission(
		UnitFeederMaint, assetGPS, closeGPS,
		"short-hash", validAfterHash, nil,
		workDetails, false, nil,
	)
	if err != ErrEmptyPhotoHash {
		t.Errorf("expected ErrEmptyPhotoHash, got: %v", err)
	}

	// 4. Metering Unit missing meter seal
	_, _, err = ValidateProofSubmission(
		UnitMetering, assetGPS, closeGPS,
		validBeforeHash, validAfterHash, nil,
		workDetails, false, nil,
	)
	if err != ErrMissingMeterSeal {
		t.Errorf("expected ErrMissingMeterSeal on Metering task, got: %v", err)
	}

	// 5. Short work details
	_, _, err = ValidateProofSubmission(
		UnitFeederMaint, assetGPS, closeGPS,
		validBeforeHash, validAfterHash, nil,
		"fixed fuse", false, nil,
	)
	if err != ErrShortWorkDetails {
		t.Errorf("expected ErrShortWorkDetails, got: %v", err)
	}

	// 6. PTW Required but missing
	_, _, err = ValidateProofSubmission(
		UnitFeederMaint, assetGPS, closeGPS,
		validBeforeHash, validAfterHash, nil,
		workDetails, true, nil,
	)
	if err != ErrPTWPermitRequired {
		t.Errorf("expected ErrPTWPermitRequired, got: %v", err)
	}

	// 7. PTW Required with valid permit
	permitID := "ptw-permit-uuid-123"
	_, valid, err = ValidateProofSubmission(
		UnitFeederMaint, assetGPS, closeGPS,
		validBeforeHash, validAfterHash, nil,
		workDetails, true, &permitID,
	)
	if err != nil || !valid {
		t.Errorf("expected success with permit, got err: %v, valid: %v", err, valid)
	}

	// 8. GPS Out-of-Tolerance (far away)
	dist, valid, err = ValidateProofSubmission(
		UnitFeederMaint, assetGPS, farGPS,
		validBeforeHash, validAfterHash, nil,
		workDetails, false, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error on far GPS: %v", err)
	}
	if valid {
		t.Errorf("expected valid=false for far GPS (%f meters), but was true", dist)
	}

	// 9. Metering Unit with Seal Success
	dist, valid, err = ValidateProofSubmission(
		UnitMetering, assetGPS, closeGPS,
		validBeforeHash, validAfterHash, &seal,
		workDetails, false, nil,
	)
	if err != nil || !valid {
		t.Errorf("expected valid metering proof with seal, got err: %v, valid: %v", err, valid)
	}
}

func TestValidateUnitHeadVerification(t *testing.T) {
	// 1. Verify from non-PROOF_SUBMITTED status should fail
	_, err := ValidateUnitHeadVerification(StatusDispatched, true, "")
	if err == nil {
		t.Errorf("expected error verifying non-submitted work order")
	}

	// 2. Valid Approval
	newStatus, err := ValidateUnitHeadVerification(StatusProofSubmitted, true, "Great job, approved.")
	if err != nil || newStatus != StatusAttendedVerified {
		t.Errorf("expected StatusAttendedVerified on approval, got %s, err: %v", newStatus, err)
	}

	// 3. Rejection with adequate notes
	newStatus, err = ValidateUnitHeadVerification(StatusProofSubmitted, false, "Phase 2 drop-out fuse was not tensioned properly; re-align.")
	if err != nil || newStatus != StatusRejectedRework {
		t.Errorf("expected StatusRejectedRework on rejection, got %s, err: %v", newStatus, err)
	}

	// 4. Rejection without adequate notes
	_, err = ValidateUnitHeadVerification(StatusProofSubmitted, false, "bad")
	if err != ErrMissingReworkNotes {
		t.Errorf("expected ErrMissingReworkNotes, got: %v", err)
	}
}

