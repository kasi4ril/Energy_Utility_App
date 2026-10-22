package domain

import (
	"testing"
)

func TestEvaluateTelemetryPhysics_BoundariesAndAlerts(t *testing.T) {
	validTelemetry := MasterStationTelemetryLog{
		WaterLevelM:         85.4,
		MWGenerated:         12.5,
		PlantState:          StateInService,
		StatorWindingTempC:  78.0,
		BearingVibrationMMS: 1.2,
		LubeOilPressureBar:  3.2,
	}

	// 1. Happy path - normal operational conditions
	alerts, err := EvaluateTelemetryPhysics(validTelemetry)
	if err != nil {
		t.Fatalf("unexpected error on valid telemetry: %v", err)
	}
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got: %v", alerts)
	}

	// 2. Invalid water level (< 0 or > 250)
	badWater := validTelemetry
	badWater.WaterLevelM = -5.0
	if _, err := EvaluateTelemetryPhysics(badWater); err != ErrInvalidWaterLevel {
		t.Errorf("expected ErrInvalidWaterLevel for negative water level, got: %v", err)
	}

	badWater.WaterLevelM = 300.0
	if _, err := EvaluateTelemetryPhysics(badWater); err != ErrInvalidWaterLevel {
		t.Errorf("expected ErrInvalidWaterLevel for excessive water level, got: %v", err)
	}

	// 3. Negative power
	badPower := validTelemetry
	badPower.MWGenerated = -1.0
	if _, err := EvaluateTelemetryPhysics(badPower); err != ErrNegativePower {
		t.Errorf("expected ErrNegativePower, got: %v", err)
	}

	// 4. Power on out-of-service plant
	badOffline := validTelemetry
	badOffline.PlantState = StateOutOfService
	badOffline.MWGenerated = 5.0
	if _, err := EvaluateTelemetryPhysics(badOffline); err != ErrInvalidPowerForOfflinePlant {
		t.Errorf("expected ErrInvalidPowerForOfflinePlant, got: %v", err)
	}

	// 5. Critical Thermal Alert (> 110C)
	hotWinding := validTelemetry
	hotWinding.StatorWindingTempC = 115.0
	alerts, err = EvaluateTelemetryPhysics(hotWinding)
	if err != nil || len(alerts) != 1 {
		t.Fatalf("expected 1 critical thermal alert, got: %v", alerts)
	}

	// Warning Thermal Alert (> 95C)
	warmWinding := validTelemetry
	warmWinding.StatorWindingTempC = 98.0
	alerts, err = EvaluateTelemetryPhysics(warmWinding)
	if err != nil || len(alerts) != 1 {
		t.Fatalf("expected 1 warning thermal alert, got: %v", alerts)
	}

	// 6. Critical Vibration Alert (> 4.5 mm/s)
	highVib := validTelemetry
	highVib.BearingVibrationMMS = 5.2
	alerts, err = EvaluateTelemetryPhysics(highVib)
	if err != nil || len(alerts) != 1 {
		t.Fatalf("expected 1 critical vibration alert, got: %v", alerts)
	}

	// Warning Vibration Alert (> 3.0 mm/s)
	warnVib := validTelemetry
	warnVib.BearingVibrationMMS = 3.5
	alerts, err = EvaluateTelemetryPhysics(warnVib)
	if err != nil || len(alerts) != 1 {
		t.Fatalf("expected 1 warning vibration alert, got: %v", alerts)
	}

	// 7. Low Lube Oil Pressure Alert (< 1.5 bar) while IN_SERVICE
	lowOil := validTelemetry
	lowOil.LubeOilPressureBar = 1.1
	alerts, err = EvaluateTelemetryPhysics(lowOil)
	if err != nil || len(alerts) != 1 {
		t.Fatalf("expected 1 critical oil alert, got: %v", alerts)
	}
}

func TestValidateRollcallSubmission(t *testing.T) {
	validRC := RemotePlantShiftRollcall{
		PlantID:       "plant-1-uuid",
		ShiftID:       "SHIFT-20261022-M",
		OperatorNames: []string{"Operator John Yakubu", "Operator Grace Pam"},
		HandoverNotes: "Shift handed over normally; Turbine #2 bearing lubed.",
	}

	if err := ValidateRollcallSubmission(validRC); err != nil {
		t.Fatalf("expected valid rollcall, got: %v", err)
	}

	// Empty PlantID
	badRC := validRC
	badRC.PlantID = ""
	if err := ValidateRollcallSubmission(badRC); err != ErrInvalidRollcall {
		t.Errorf("expected ErrInvalidRollcall on empty plant, got: %v", err)
	}

	// Empty Operators
	badRC = validRC
	badRC.OperatorNames = []string{""}
	if err := ValidateRollcallSubmission(badRC); err == nil {
		t.Errorf("expected error on empty operator names")
	}

	// Short handover notes
	badRC = validRC
	badRC.HandoverNotes = "ok"
	if err := ValidateRollcallSubmission(badRC); err == nil {
		t.Errorf("expected error on short handover notes")
	}
}

func TestValidateFaultJobCreation(t *testing.T) {
	validJob := PlantFaultAndWorkshopJob{
		PlantID:              "plant-1-uuid",
		ReportedByOperatorID: "operator-uuid",
		Wing:                 WingMechanical,
		FaultDescription:     "Turbine guide bearing experiencing severe cavitation and oil weeping.",
	}

	if err := ValidateFaultJobCreation(validJob); err != nil {
		t.Fatalf("expected valid fault job, got: %v", err)
	}

	// Empty PlantID
	badJob := validJob
	badJob.PlantID = ""
	if err := ValidateFaultJobCreation(badJob); err == nil {
		t.Errorf("expected error on empty plant_id")
	}

	// Invalid wing
	badJob = validJob
	badJob.Wing = "CIVIL"
	if err := ValidateFaultJobCreation(badJob); err == nil {
		t.Errorf("expected error on invalid wing")
	}

	// Short description
	badJob = validJob
	badJob.FaultDescription = "leak"
	if err := ValidateFaultJobCreation(badJob); err == nil {
		t.Errorf("expected error on short description")
	}
}

func TestValidateGMApproval(t *testing.T) {
	// 1. Valid Approval from REPORTED
	newStatus, err := ValidateGMApproval(JobStatusReported, true, "Approved for emergency fabrication.")
	if err != nil || newStatus != JobStatusGMApproved {
		t.Errorf("expected JobStatusGMApproved, got %s, err: %v", newStatus, err)
	}

	// 2. Valid Rejection with notes
	newStatus, err = ValidateGMApproval(JobStatusUnitHeadReviewed, false, "Defer until dry season maintenance overhaul.")
	if err != nil || newStatus != JobStatusCancelled {
		t.Errorf("expected JobStatusCancelled on rejection, got %s, err: %v", newStatus, err)
	}

	// 3. Rejection without notes
	_, err = ValidateGMApproval(JobStatusReported, false, "bad")
	if err != ErrMissingGMRejectionNotes {
		t.Errorf("expected ErrMissingGMRejectionNotes, got: %v", err)
	}

	// 4. Invalid status for GM approval
	_, err = ValidateGMApproval(JobStatusGMApproved, true, "")
	if err == nil {
		t.Errorf("expected error approving already approved job")
	}
}

func TestValidateWorkshopTransition(t *testing.T) {
	// Permitted: NOT_APPLICABLE -> PENDING_DISPATCH
	if err := ValidateWorkshopTransition(WorkshopNotApplicable, WorkshopPendingDispatch); err != nil {
		t.Errorf("expected success for NotApplicable -> PendingDispatch: %v", err)
	}

	// Permitted: PENDING_DISPATCH -> MACHINING_IN_PROGRESS
	if err := ValidateWorkshopTransition(WorkshopPendingDispatch, WorkshopMachiningInProgress); err != nil {
		t.Errorf("expected success for PendingDispatch -> Machining: %v", err)
	}

	// Permitted: MACHINING -> QUALITY_CHECK
	if err := ValidateWorkshopTransition(WorkshopMachiningInProgress, WorkshopQualityCheck); err != nil {
		t.Errorf("expected success for Machining -> QualityCheck: %v", err)
	}

	// Permitted: QUALITY_CHECK -> COMPLETED_DELIVERED
	if err := ValidateWorkshopTransition(WorkshopQualityCheck, WorkshopCompletedDelivered); err != nil {
		t.Errorf("expected success for QualityCheck -> Delivered: %v", err)
	}

	// Illegal: NOT_APPLICABLE -> COMPLETED_DELIVERED
	if err := ValidateWorkshopTransition(WorkshopNotApplicable, WorkshopCompletedDelivered); err == nil {
		t.Errorf("expected error for illegal direct transition to delivered")
	}

	// Illegal: Unknown target
	if err := ValidateWorkshopTransition(WorkshopPendingDispatch, "UNKNOWN"); err == nil {
		t.Errorf("expected error for unknown workshop status")
	}
}

