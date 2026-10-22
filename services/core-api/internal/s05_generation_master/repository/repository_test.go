package repository

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
)

func TestMemoryRepository_TelemetryAndRollcall(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	telemetry := domain.MasterStationTelemetryLog{
		ID:                  "telem-1",
		ShiftID:             "SHIFT-20261022-M",
		ShiftWindow:         domain.ShiftMorning,
		PlantID:             "plant-1",
		PlantState:          domain.StateInService,
		WaterLevelM:         88.2,
		MWGenerated:         15.4,
		LoggedByOperatorID:  "op-1",
		RecordedAt:          time.Now().UTC(),
	}

	// 1. Create Telemetry
	if err := repo.CreateTelemetryLog(ctx, telemetry); err != nil {
		t.Fatalf("CreateTelemetryLog failed: %v", err)
	}

	// Duplicate
	if err := repo.CreateTelemetryLog(ctx, telemetry); err != ErrConflict {
		t.Errorf("expected ErrConflict on duplicate telemetry, got: %v", err)
	}

	// Get latest
	latest, err := repo.GetLatestTelemetryByPlant(ctx, "plant-1")
	if err != nil || latest.MWGenerated != 15.4 {
		t.Fatalf("GetLatestTelemetryByPlant failed: %v, got: %+v", err, latest)
	}

	// Not found
	if _, err := repo.GetLatestTelemetryByPlant(ctx, "plant-99"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound for unknown plant, got: %v", err)
	}

	// List
	list, err := repo.ListTelemetryLogs(ctx, "SHIFT-20261022-M", "plant-1", 10)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 telemetry log, got %d, err: %v", len(list), err)
	}

	// 2. Rollcall
	rc := domain.RemotePlantShiftRollcall{
		ID:                  "rc-1",
		ShiftID:             "SHIFT-20261022-M",
		ShiftWindow:         domain.ShiftMorning,
		PlantID:             "plant-1",
		OperatorNames:       []string{"Operator Pam"},
		HandoverNotes:       "All turbines online.",
		MasterOperatorID:    "op-1",
		CreatedAt:           time.Now().UTC(),
	}

	if err := repo.CreateRollcall(ctx, rc); err != nil {
		t.Fatalf("CreateRollcall failed: %v", err)
	}
	if err := repo.CreateRollcall(ctx, rc); err != ErrConflict {
		t.Errorf("expected ErrConflict on duplicate rollcall, got: %v", err)
	}

	rcList, err := repo.ListRollcallsByShift(ctx, "SHIFT-20261022-M")
	if err != nil || len(rcList) != 1 {
		t.Fatalf("expected 1 rollcall, got %d, err: %v", len(rcList), err)
	}
}

func TestMemoryRepository_PlantJobsAndWorkshop(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	job := domain.PlantFaultAndWorkshopJob{
		ID:                   "job-1",
		JobNo:                "JOB-20261022-0001",
		PlantID:              "plant-1",
		ReportedByOperatorID: "op-1",
		Wing:                 domain.WingMechanical,
		FaultDescription:     "Turbine guide bearing cavitation.",
		Severity:             "HIGH",
		RequiresWorkshop:     true,
		WorkshopStatus:       domain.WorkshopPendingDispatch,
		Status:               domain.JobStatusReported,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	// 1. Create Job
	if err := repo.CreatePlantJob(ctx, job); err != nil {
		t.Fatalf("CreatePlantJob failed: %v", err)
	}
	if err := repo.CreatePlantJob(ctx, job); err != ErrConflict {
		t.Errorf("expected ErrConflict on duplicate job, got: %v", err)
	}

	// Get by ID
	fetched, err := repo.GetPlantJobByID(ctx, "job-1")
	if err != nil || fetched.JobNo != job.JobNo {
		t.Fatalf("GetPlantJobByID failed: %v, got: %+v", err, fetched)
	}
	if _, err := repo.GetPlantJobByID(ctx, "non-existent"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}

	// 2. GM Approval
	if err := repo.UpdatePlantJobGMApproval(ctx, "job-1", true, "Approved for workshop."); err != nil {
		t.Fatalf("UpdatePlantJobGMApproval failed: %v", err)
	}
	updated, _ := repo.GetPlantJobByID(ctx, "job-1")
	if !updated.GMApproved || updated.Status != domain.JobStatusGMApproved {
		t.Fatalf("expected GMApproved true and JobStatusGMApproved, got: %+v", updated)
	}

	// GM Rejection
	if err := repo.UpdatePlantJobGMApproval(ctx, "job-1", false, "Budget denied."); err != nil {
		t.Fatalf("UpdatePlantJobGMApproval failed: %v", err)
	}
	updated, _ = repo.GetPlantJobByID(ctx, "job-1")
	if updated.GMApproved || updated.Status != domain.JobStatusCancelled {
		t.Fatalf("expected JobStatusCancelled on rejection, got: %+v", updated)
	}

	// 3. Workshop Status Transitions
	matRef := "MAT-REQ-001"
	if err := repo.UpdateWorkshopStatus(ctx, "job-1", domain.WorkshopMachiningInProgress, &matRef); err != nil {
		t.Fatalf("UpdateWorkshopStatus failed: %v", err)
	}
	wJob, _ := repo.GetPlantJobByID(ctx, "job-1")
	if wJob.WorkshopStatus != domain.WorkshopMachiningInProgress || wJob.StoreRequisitionRef == nil {
		t.Fatalf("expected WorkshopMachiningInProgress and store ref, got: %+v", wJob)
	}

	// Delivered
	if err := repo.UpdateWorkshopStatus(ctx, "job-1", domain.WorkshopCompletedDelivered, nil); err != nil {
		t.Fatalf("UpdateWorkshopStatus delivery failed: %v", err)
	}
	wJob, _ = repo.GetPlantJobByID(ctx, "job-1")
	if wJob.Status != domain.JobStatusResolvedReinstalled {
		t.Fatalf("expected JobStatusResolvedReinstalled, got: %+v", wJob)
	}

	// 4. List jobs
	jobs, err := repo.ListPlantJobs(ctx, "plant-1", domain.WingMechanical, domain.JobStatusResolvedReinstalled)
	if err != nil || len(jobs) != 1 {
		t.Fatalf("expected 1 job matching filters, got %d, err: %v", len(jobs), err)
	}
}

