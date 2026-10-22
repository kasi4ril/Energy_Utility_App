// Package transport_test tests request/response DTO validations for S09.
package transport_test

import (
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/transport"
)

func TestCaptureSnapshotDTO_Validation(t *testing.T) {
	valid := transport.CaptureSnapshotDTO{
		SnapshotDate:      "2026-10-22",
		ShiftWindow:       "MORNING",
		CapturedByStaffID: "STAFF-01",
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO to pass, got: %v", err)
	}

	invalid := valid
	invalid.SnapshotDate = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty date, got nil")
	}

	invalid = valid
	invalid.ShiftWindow = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty shift window, got nil")
	}

	invalid = valid
	invalid.CapturedByStaffID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty staff ID, got nil")
	}
}

func TestEnqueueActionDTO_Validation(t *testing.T) {
	valid := transport.EnqueueActionDTO{
		ActionType:       domain.ActionMDProcurementApproval,
		TargetRole:       domain.RoleMD,
		SourceEntityID:   "REQ-001",
		SourceEntityRef:  "PO-001",
		SourceDepartment: "STORES_AUDIT",
		Summary:          "Critical Out-of-Stock Transformer Procurement",
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	invalid := valid
	invalid.SourceEntityID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty source entity ID, got nil")
	}

	invalid = valid
	invalid.ActionType = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty action type, got nil")
	}

	invalid = valid
	invalid.TargetRole = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty target role, got nil")
	}
}

func TestDecideActionDTO_Validation(t *testing.T) {
	valid := transport.DecideActionDTO{
		DecidedByStaffID: "STAFF-MD-01",
		Approve:          true,
		Notes:            "Approved",
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	invalid := valid
	invalid.DecidedByStaffID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty staff ID, got nil")
	}

	invalid = valid
	invalid.Notes = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty notes, got nil")
	}
}

func TestPredictiveDTO_Validation(t *testing.T) {
	// Hydro
	hydro := transport.HydroPrognosisDTO{AssetID: "RES-01"}
	if err := hydro.Validate(); err != nil {
		t.Fatalf("expected valid hydro DTO, got %v", err)
	}
	hydro.AssetID = ""
	if err := hydro.Validate(); err == nil {
		t.Fatal("expected error on empty hydro asset ID")
	}

	// Generator
	gen := transport.GeneratorPrognosisDTO{AssetID: "GEN-01"}
	if err := gen.Validate(); err != nil {
		t.Fatalf("expected valid generator DTO, got %v", err)
	}
	gen.AssetID = ""
	if err := gen.Validate(); err == nil {
		t.Fatal("expected error on empty gen asset ID")
	}

	// Transformer
	tx := transport.TransformerPrognosisDTO{AssetID: "TX-01"}
	if err := tx.Validate(); err != nil {
		t.Fatalf("expected valid tx DTO, got %v", err)
	}
	tx.AssetID = ""
	if err := tx.Validate(); err == nil {
		t.Fatal("expected error on empty tx asset ID")
	}

	// Feeder
	fdr := transport.FeederPrognosisDTO{AssetID: "FDR-01"}
	if err := fdr.Validate(); err != nil {
		t.Fatalf("expected valid fdr DTO, got %v", err)
	}
	fdr.AssetID = ""
	if err := fdr.Validate(); err == nil {
		t.Fatal("expected error on empty fdr asset ID")
	}

	// Acknowledge Trigger
	ack := transport.AcknowledgeTriggerDTO{StaffID: "STAFF-01"}
	if err := ack.Validate(); err != nil {
		t.Fatalf("expected valid ack DTO, got %v", err)
	}
	ack.StaffID = ""
	if err := ack.Validate(); err == nil {
		t.Fatal("expected error on empty ack staff ID")
	}
}

