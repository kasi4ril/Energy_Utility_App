// Package transport_test validates DTO field validations and error handling.
package transport_test

import (
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/transport"
)

func TestRecordBulkIntakeDTO_Validate(t *testing.T) {
	valid := transport.RecordBulkIntakeDTO{
		ShiftID:        "SHIFT-001",
		ShiftWindow:    domain.ShiftMorning,
		GenerationMWIn: 10.0,
		TCNMWIn:        2.0,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO to pass, got: %v", err)
	}

	// Empty Shift ID
	invalidShift := valid
	invalidShift.ShiftID = "   "
	if err := invalidShift.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields, got: %v", err)
	}

	// Invalid Window
	invalidWin := valid
	invalidWin.ShiftWindow = "MIDDAY"
	if err := invalidWin.Validate(); err != transport.ErrInvalidShiftWindow {
		t.Errorf("expected ErrInvalidShiftWindow, got: %v", err)
	}

	// Negative MW
	negMW := valid
	negMW.GenerationMWIn = -1.0
	if err := negMW.Validate(); err != transport.ErrNegativeMW {
		t.Errorf("expected ErrNegativeMW, got: %v", err)
	}
}

func TestRationingCalculateDTO_Validate(t *testing.T) {
	valid := transport.RationingCalculateDTO{AvailableBulkMW: 12.0}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	zeroVal := transport.RationingCalculateDTO{AvailableBulkMW: 0.0}
	if err := zeroVal.Validate(); err != nil {
		t.Fatalf("expected zero MW to be valid, got: %v", err)
	}

	negVal := transport.RationingCalculateDTO{AvailableBulkMW: -5.0}
	if err := negVal.Validate(); err != transport.ErrNegativeMW {
		t.Errorf("expected ErrNegativeMW, got: %v", err)
	}
}

func TestInitiatePTWDTO_Validate(t *testing.T) {
	valid := transport.InitiatePTWDTO{
		AssetID:         "asset-001",
		WorkDescription: "Overhaul line isolator",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	missingAsset := valid
	missingAsset.AssetID = ""
	if err := missingAsset.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields, got: %v", err)
	}

	missingDesc := valid
	missingDesc.WorkDescription = "  "
	if err := missingDesc.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields, got: %v", err)
	}
}

func TestApprovePTWDTO_Validate(t *testing.T) {
	validDept := transport.ApprovePTWDTO{ApproverRole: "HEAD_SERVICES"}
	if err := validDept.Validate(); err != nil {
		t.Fatalf("expected HEAD_SERVICES to pass, got: %v", err)
	}

	validSafety := transport.ApprovePTWDTO{ApproverRole: "SAFETY_REP"}
	if err := validSafety.Validate(); err != nil {
		t.Fatalf("expected SAFETY_REP to pass, got: %v", err)
	}

	invalidRole := transport.ApprovePTWDTO{ApproverRole: "CLERK"}
	if err := invalidRole.Validate(); err != transport.ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got: %v", err)
	}
}

func TestIsolateAndClearSafetyPTWDTO_Validate(t *testing.T) {
	// Isolate
	validIso := transport.IsolatePTWDTO{IsolationNotes: "Breaker racked out and earthing attached"}
	if err := validIso.Validate(); err != nil {
		t.Fatalf("expected valid isolation, got: %v", err)
	}

	shortIso := transport.IsolatePTWDTO{IsolationNotes: "done"}
	if err := shortIso.Validate(); err != transport.ErrNotesTooShort {
		t.Errorf("expected ErrNotesTooShort, got: %v", err)
	}

	// Safety Clearance
	validClear := transport.ClearSafetyPTWDTO{Notes: "Crews accounted for and line confirmed safe"}
	if err := validClear.Validate(); err != nil {
		t.Fatalf("expected valid clearance, got: %v", err)
	}

	shortClear := transport.ClearSafetyPTWDTO{Notes: "ok"}
	if err := shortClear.Validate(); err != transport.ErrNotesTooShort {
		t.Errorf("expected ErrNotesTooShort, got: %v", err)
	}
}

