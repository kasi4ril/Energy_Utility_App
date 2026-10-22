package service

import (
	"context"
	"encoding/json"
	"testing"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const testHMACKey = "nesco-master-generation-hmac-32b!"

func setupServiceTest(t *testing.T) (*WorkflowService, *repository.MemoryRepository, *s01repo.MemoryRepository) {
	t.Helper()

	repo := repository.NewMemoryRepository()
	gridRepo := s01repo.NewMemoryRepository()
	bus := events.NewMemoryEventBus()

	// Seed Generating Plant 1 (Kurra Falls)
	plant := s01domain.GridAsset{
		ID:        "plant-1-uuid",
		AssetCode: "GEN-PLANT-01",
		AssetType: s01domain.AssetPlant,
		OrgUnitID: "org-unit-gen",
		GPSPoint:  s01domain.GPSPoint{Latitude: 9.4200, Longitude: 8.8500},
	}
	if err := gridRepo.CreateGridAsset(context.Background(), plant); err != nil {
		t.Fatalf("failed to seed plant: %v", err)
	}

	svc := NewWorkflowService(repo, gridRepo, bus, testHMACKey)
	return svc, repo, gridRepo
}

func TestWorkflowService_RecordTelemetryAndPhysicsAlerts(t *testing.T) {
	ctx := context.Background()
	svc, repo, _ := setupServiceTest(t)

	// 1. Success with High Stator Temp Warning Alert
	tLog := domain.MasterStationTelemetryLog{
		ShiftID:             "SHIFT-20261022-M",
		ShiftWindow:         domain.ShiftMorning,
		PlantID:             "plant-1-uuid",
		PlantState:          domain.StateInService,
		WaterLevelM:         92.5,
		MWGenerated:         18.2,
		StatorWindingTempC:  112.0, // Critical thermal alert
		BearingVibrationMMS: 1.5,
		LubeOilPressureBar:  3.5,
		LoggedByOperatorID:  "master-op-1",
	}

	alerts, err := svc.RecordTelemetry(ctx, tLog)
	if err != nil {
		t.Fatalf("RecordTelemetry failed: %v", err)
	}
	if len(alerts) != 1 {
		t.Fatalf("expected 1 critical alert, got: %v", alerts)
	}

	// Verify logged in repo
	latest, err := repo.GetLatestTelemetryByPlant(ctx, "plant-1-uuid")
	if err != nil || latest.MWGenerated != 18.2 {
		t.Fatalf("expected latest telemetry in repo, got err: %v, data: %+v", err, latest)
	}

	// 2. Negative: Missing Plant Asset
	tLog.PlantID = "missing-plant"
	_, err = svc.RecordTelemetry(ctx, tLog)
	if err == nil {
		t.Errorf("expected error for non-existent plant asset")
	}

	// 3. Negative: Physics Error (Negative Power)
	tLog.PlantID = "plant-1-uuid"
	tLog.MWGenerated = -5.0
	_, err = svc.RecordTelemetry(ctx, tLog)
	if err == nil {
		t.Errorf("expected error on negative power generation")
	}
}

func TestWorkflowService_RollcallAndGoogleFormsWebhook(t *testing.T) {
	ctx := context.Background()
	svc, repo, _ := setupServiceTest(t)

	// 1. Record Rollcall directly
	rc := domain.RemotePlantShiftRollcall{
		ShiftID:             "SHIFT-20261022-M",
		ShiftWindow:         domain.ShiftMorning,
		PlantID:             "plant-1-uuid",
		OperatorNames:       []string{"Operator Musa", "Operator Yakubu"},
		ProtectionEngineers: []string{"Engr. Pam"},
		HandoverNotes:       "Turbines 1 and 2 running smoothly.",
		MasterOperatorID:    "master-op-1",
	}
	recordedRC, err := svc.RecordRollcall(ctx, rc)
	if err != nil || recordedRC.ID == "" {
		t.Fatalf("RecordRollcall failed: %v", err)
	}

	// 2. Ingest Google Forms Webhook (Valid HMAC)
	webhookPayload := map[string]interface{}{
		"shift_id":           "SHIFT-20261022-A",
		"shift_window":       "AFTERNOON",
		"master_operator_id": "master-op-2",
		"plants": []map[string]interface{}{
			{
				"plant_id":             "plant-1-uuid",
				"plant_state":          "IN_SERVICE",
				"water_level_m":        91.0,
				"mw_generated":         17.5,
				"stator_temp_c":        75.0,
				"vibration_mm_s":       1.1,
				"oil_pressure_bar":     3.2,
				"operator_names":       []string{"Operator Grace"},
				"protection_engineers": []string{"Engr. Dung"},
				"handover_notes":       "Hydraulic governor tuned.",
			},
		},
	}
	rawBytes, _ := json.Marshal(webhookPayload)
	validHMAC, _ := crypto.ComputeRowHMAC(testHMACKey, string(rawBytes))

	alerts, err := svc.IngestGoogleFormsWebhook(ctx, rawBytes, validHMAC)
	if err != nil {
		t.Fatalf("IngestGoogleFormsWebhook failed: %v", err)
	}
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts for normal operation, got: %v", alerts)
	}

	// Verify telemetry was created from webhook
	latest, err := repo.GetLatestTelemetryByPlant(ctx, "plant-1-uuid")
	if err != nil || latest.MWGenerated != 17.5 {
		t.Fatalf("expected updated telemetry from webhook, got: %+v", latest)
	}

	// 3. Ingest with Invalid HMAC Signature
	_, err = svc.IngestGoogleFormsWebhook(ctx, rawBytes, "bad-signature")
	if err != ErrInvalidWebhookSignature {
		t.Errorf("expected ErrInvalidWebhookSignature, got: %v", err)
	}
}

func TestWorkflowService_PlantFaultAndGMWorkshopLoop(t *testing.T) {
	ctx := context.Background()
	svc, repo, _ := setupServiceTest(t)

	// 1. Report Plant Fault requiring workshop
	job, err := svc.ReportPlantFault(
		ctx, "plant-1-uuid", "master-op-1",
		domain.WingMechanical,
		"Turbine guide bearing experiencing severe cavitation and oil weeping.",
		"HIGH", true,
	)
	if err != nil {
		t.Fatalf("ReportPlantFault failed: %v", err)
	}
	if job.Status != domain.JobStatusReported || job.WorkshopStatus != domain.WorkshopPendingDispatch {
		t.Fatalf("unexpected job status: %+v", job)
	}

	// 2. GM Rejection first to test branch
	_, err = svc.ApproveGMJob(ctx, job.ID, "gm-uuid-1", false, "bad")
	if err != domain.ErrMissingGMRejectionNotes {
		t.Errorf("expected ErrMissingGMRejectionNotes, got: %v", err)
	}

	// 3. GM Approval
	approvedJob, err := svc.ApproveGMJob(ctx, job.ID, "gm-uuid-1", true, "Approved for immediate workshop machining.")
	if err != nil || approvedJob.Status != domain.JobStatusGMApproved {
		t.Fatalf("ApproveGMJob failed: %v, got: %+v", err, approvedJob)
	}

	// 4. Workshop Machining in progress
	matRef := "MAT-REQ-BRG-001"
	wJob, err := svc.UpdateWorkshopFabrication(ctx, job.ID, domain.WorkshopMachiningInProgress, &matRef)
	if err != nil || wJob.WorkshopStatus != domain.WorkshopMachiningInProgress {
		t.Fatalf("UpdateWorkshopFabrication machining failed: %v", err)
	}

	// 5. Workshop Quality Check
	wJob, err = svc.UpdateWorkshopFabrication(ctx, job.ID, domain.WorkshopQualityCheck, nil)
	if err != nil || wJob.WorkshopStatus != domain.WorkshopQualityCheck {
		t.Fatalf("UpdateWorkshopFabrication quality check failed: %v", err)
	}

	// 6. Workshop Delivered and Reinstalled
	wJob, err = svc.UpdateWorkshopFabrication(ctx, job.ID, domain.WorkshopCompletedDelivered, nil)
	if err != nil || wJob.WorkshopStatus != domain.WorkshopCompletedDelivered {
		t.Fatalf("UpdateWorkshopFabrication delivery failed: %v", err)
	}

	// Verify final job status in repo
	inRepo, _ := repo.GetPlantJobByID(ctx, job.ID)
	if inRepo.Status != domain.JobStatusResolvedReinstalled {
		t.Fatalf("expected JobStatusResolvedReinstalled, got: %s", inRepo.Status)
	}
}
