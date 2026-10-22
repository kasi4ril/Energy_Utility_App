// Package repository_test verifies the persistence operations and invariants
// of the S06 Energy Management MemoryRepository.
package repository_test

import (
	"context"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
)

func TestMemoryRepository_DispatchLogs(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	// Initial get on empty repository should return ErrNotFound
	_, err := repo.GetLatestDispatchLog(ctx)
	if err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for latest dispatch log, got %v", err)
	}

	t1 := time.Now().UTC().Add(-1 * time.Hour)
	t2 := time.Now().UTC()

	log1 := domain.DispatchTelemetryLog{
		ID:                          "log-001",
		ShiftID:                     "SHIFT-MORNING",
		ShiftWindow:                 domain.ShiftMorning,
		GenerationMWIn:              7.5,
		TCNMWIn:                     0.0,
		TotalBulkMW:                 7.5,
		SubstationStaffRollcallJSON: `[{"substation_id":"sub-01","operator_name":"M. Garba"}]`,
		HandoverNotes:               "All lines normal",
		OperatorID:                  "op-101",
		LoggedAt:                    t1,
	}

	if err := repo.CreateDispatchLog(ctx, log1); err != nil {
		t.Fatalf("unexpected error creating log1: %v", err)
	}

	// Conflict detection on duplicate ID
	if err := repo.CreateDispatchLog(ctx, log1); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate log ID, got %v", err)
	}

	log2 := domain.DispatchTelemetryLog{
		ID:                          "log-002",
		ShiftID:                     "SHIFT-EVENING",
		ShiftWindow:                 domain.ShiftAfternoon,
		GenerationMWIn:              8.5,
		TCNMWIn:                     0.0,
		TotalBulkMW:                 8.5,
		SubstationStaffRollcallJSON: `[{"substation_id":"sub-01","operator_name":"J. Bello"}]`,
		HandoverNotes:               "Line 2 under monitoring",
		OperatorID:                  "op-102",
		LoggedAt:                    t2,
	}

	if err := repo.CreateDispatchLog(ctx, log2); err != nil {
		t.Fatalf("unexpected error creating log2: %v", err)
	}

	latest, err := repo.GetLatestDispatchLog(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting latest dispatch log: %v", err)
	}
	if latest.ID != "log-002" {
		t.Fatalf("expected latest log ID 'log-002', got '%s'", latest.ID)
	}

	// Filter by shift
	morningLogs, err := repo.ListDispatchLogs(ctx, "SHIFT-MORNING", 10)
	if err != nil {
		t.Fatalf("unexpected error listing logs: %v", err)
	}
	if len(morningLogs) != 1 || morningLogs[0].ID != "log-001" {
		t.Fatalf("expected 1 log with ID 'log-001', got %d", len(morningLogs))
	}

	allLogs, err := repo.ListDispatchLogs(ctx, "", 1)
	if err != nil {
		t.Fatalf("unexpected error listing all logs with limit: %v", err)
	}
	if len(allLogs) != 1 {
		t.Fatalf("expected 1 log due to limit, got %d", len(allLogs))
	}
}

func TestMemoryRepository_FeederStatus(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	// Not found
	_, err := repo.GetFeederStatusByAssetID(ctx, "feeder-unknown")
	if err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	status := domain.FeederStatus{
		ID:            "fs-01",
		FeederAssetID: "feeder-01",
		FeederCode:    "FDR-IND-01",
		PriorityRank:  1,
		TariffBand:    s01domain.TariffBandA,
		CurrentLoadMW: 3.2,
		CurrentAmps:   168.0,
		State:         domain.FeederInService,
		UpdatedAt:     time.Now().UTC(),
	}

	if err := repo.UpsertFeederStatus(ctx, status); err != nil {
		t.Fatalf("unexpected error upserting feeder status: %v", err)
	}

	got, err := repo.GetFeederStatusByAssetID(ctx, "feeder-01")
	if err != nil {
		t.Fatalf("unexpected error getting feeder status: %v", err)
	}
	if got.FeederCode != "FDR-IND-01" || got.State != domain.FeederInService {
		t.Fatalf("mismatched feeder status: %+v", got)
	}

	// Update load and state
	status.CurrentLoadMW = 3.8
	status.State = domain.FeederOutOnLoadShedding
	if err := repo.UpsertFeederStatus(ctx, status); err != nil {
		t.Fatalf("unexpected error updating feeder status: %v", err)
	}

	updated, err := repo.GetFeederStatusByAssetID(ctx, "feeder-01")
	if err != nil {
		t.Fatalf("unexpected error getting updated feeder: %v", err)
	}
	if updated.State != domain.FeederOutOnLoadShedding || updated.CurrentLoadMW != 3.8 {
		t.Fatalf("state update failed: %+v", updated)
	}

	list, err := repo.ListFeederStatuses(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing feeder statuses: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 feeder status, got %d", len(list))
	}
}

func TestMemoryRepository_SafetyPermits(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	// Not found
	_, err := repo.GetPTWByID(ctx, "ptw-nonexistent")
	if err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for ID, got %v", err)
	}
	_, err = repo.GetPTWByPermitNo(ctx, "PTW-NONEXISTENT")
	if err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for PermitNo, got %v", err)
	}

	ptw := domain.SafetyPermitPTW{
		ID:         "ptw-101",
		PermitNo:   "PTW-2026-0001",
		AssetID:    "feeder-02",
		UnitHeadID: "user-unithead-1",
		Status:     domain.PTWRequested,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	if err := repo.CreatePTW(ctx, ptw); err != nil {
		t.Fatalf("unexpected error creating PTW: %v", err)
	}

	// Conflict on duplicate ID
	if err := repo.CreatePTW(ctx, ptw); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate ID, got %v", err)
	}

	// Conflict on duplicate PermitNo
	ptwDupNo := ptw
	ptwDupNo.ID = "ptw-102"
	if err := repo.CreatePTW(ctx, ptwDupNo); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate PermitNo, got %v", err)
	}

	// Fetch by ID
	byID, err := repo.GetPTWByID(ctx, "ptw-101")
	if err != nil {
		t.Fatalf("unexpected error fetching PTW by ID: %v", err)
	}
	if byID.PermitNo != "PTW-2026-0001" {
		t.Fatalf("mismatched permit number: %s", byID.PermitNo)
	}

	// Fetch by PermitNo
	byNo, err := repo.GetPTWByPermitNo(ctx, "PTW-2026-0001")
	if err != nil {
		t.Fatalf("unexpected error fetching PTW by PermitNo: %v", err)
	}
	if byNo.ID != "ptw-101" {
		t.Fatalf("mismatched permit ID: %s", byNo.ID)
	}

	// Update non-existent returns ErrNotFound
	nonExistentPTW := ptw
	nonExistentPTW.ID = "ptw-999"
	if err := repo.UpdatePTWStatus(ctx, nonExistentPTW); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound updating non-existent PTW, got %v", err)
	}

	// Update existing
	ptw.Status = domain.PTWApprovedPendingIsolation
	ptw.DeptHeadApproved = true
	ptw.SafetyRepApproved = true
	if err := repo.UpdatePTWStatus(ctx, ptw); err != nil {
		t.Fatalf("unexpected error updating PTW: %v", err)
	}

	updated, err := repo.GetPTWByID(ctx, "ptw-101")
	if err != nil {
		t.Fatalf("unexpected error fetching updated PTW: %v", err)
	}
	if updated.Status != domain.PTWApprovedPendingIsolation {
		t.Fatalf("expected status APPROVED_PENDING_ISOLATION, got %s", updated.Status)
	}

	// List with filter
	list, err := repo.ListPTWs(ctx, domain.PTWApprovedPendingIsolation, "feeder-02")
	if err != nil {
		t.Fatalf("unexpected error listing PTWs: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 permit matching filters, got %d", len(list))
	}

	emptyList, err := repo.ListPTWs(ctx, domain.PTWIsolatedEarthed, "")
	if err != nil {
		t.Fatalf("unexpected error listing PTWs: %v", err)
	}
	if len(emptyList) != 0 {
		t.Fatalf("expected 0 permits matching ISOLATED_EARTHED, got %d", len(emptyList))
	}
}
