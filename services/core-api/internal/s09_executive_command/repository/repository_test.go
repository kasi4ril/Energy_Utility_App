// Package repository_test exercises the persistence operations for snapshots and predictive triggers.
package repository_test

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
)

func TestMemoryRepository_Snapshots(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()

	snapshot := domain.ExecutiveSnapshot{
		ID:                            "snap-001",
		SnapshotDate:                  "2026-10-22",
		ShiftWindow:                   "MORNING",
		TotalGenerationMW:             18.5,
		TCNGridIntakeMW:               6.0,
		TotalSystemLoadMW:             24.0,
		FeedersInServiceCount:         12,
		FeedersOutFaultCount:          1,
		FeedersOutPTWCount:            1,
		FeedersOutRationingCount:      0,
		TotalRevenueCollected:         8500000.0,
		PrepaidRevenueSharePct:        75.0,
		UnresolvedTheftAlertsCount:    2,
		UnresolvedChurnAlertsCount:    1,
		AuditPendingProcurementsCount: 3,
		UnreturnedScrapCount:          0,
		Notes:                         "Morning peak load stable across 33kV corridor",
		CapturedByStaffID:             "STAFF-EXEC-01",
		CreatedAt:                     time.Now().UTC(),
	}

	// Create
	if err := repo.CreateExecutiveSnapshot(ctx, snapshot); err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}

	// Conflict
	if err := repo.CreateExecutiveSnapshot(ctx, snapshot); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict, got %v", err)
	}

	// Get by ID
	fetched, err := repo.GetExecutiveSnapshotByID(ctx, "snap-001")
	if err != nil {
		t.Fatalf("failed to get snapshot by ID: %v", err)
	}
	if fetched.TotalGenerationMW != 18.5 {
		t.Fatalf("expected 18.5 MW, got %f", fetched.TotalGenerationMW)
	}

	// Get by ID not found
	if _, err := repo.GetExecutiveSnapshotByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Get Latest
	latest, err := repo.GetLatestExecutiveSnapshot(ctx)
	if err != nil {
		t.Fatalf("failed to get latest snapshot: %v", err)
	}
	if latest.ID != "snap-001" {
		t.Fatalf("expected snap-001, got %s", latest.ID)
	}

	// List
	list, err := repo.ListExecutiveSnapshots(ctx, 10)
	if err != nil {
		t.Fatalf("failed to list snapshots: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(list))
	}
}

func TestMemoryRepository_PredictiveTriggers(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()

	trigger := domain.PredictiveMaintenanceTrigger{
		ID:                   "trg-001",
		TriggerNo:            "PMT-2026-0001",
		AssetID:              "GEN-PLANT1-UNIT01",
		AssetType:            domain.AssetGenerator,
		ModelName:            domain.ModelGeneratorPreventive,
		RiskScore:            0.82,
		RiskLevel:            domain.RiskCritical,
		TelemetrySummaryJSON: `{"vibration": 11.5}`,
		RecommendedAction:    "Immediate emergency shutdown and overhaul",
		AutoWorkOrderCreated: true,
		Status:               "ACTIVE",
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}

	// Create
	if err := repo.CreatePredictiveTrigger(ctx, trigger); err != nil {
		t.Fatalf("failed to create trigger: %v", err)
	}

	// Conflict ID
	if err := repo.CreatePredictiveTrigger(ctx, trigger); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate ID, got %v", err)
	}

	// Conflict No
	t2 := trigger
	t2.ID = "trg-002"
	if err := repo.CreatePredictiveTrigger(ctx, t2); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate triggerNo, got %v", err)
	}

	// Get by ID
	fetched, err := repo.GetPredictiveTriggerByID(ctx, "trg-001")
	if err != nil {
		t.Fatalf("failed to get trigger: %v", err)
	}
	if fetched.TriggerNo != "PMT-2026-0001" {
		t.Fatalf("expected PMT-2026-0001, got %s", fetched.TriggerNo)
	}

	// Get by ID - Not Found
	if _, err := repo.GetPredictiveTriggerByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Get by No
	byNo, err := repo.GetPredictiveTriggerByNo(ctx, "PMT-2026-0001")
	if err != nil {
		t.Fatalf("failed to get trigger by no: %v", err)
	}
	if byNo.ID != "trg-001" {
		t.Fatalf("expected trg-001, got %s", byNo.ID)
	}

	// Get by No - Not Found
	if _, err := repo.GetPredictiveTriggerByNo(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Update
	staffID := "STAFF-GEN-01"
	now := time.Now().UTC()
	trigger.AcknowledgedByStaffID = &staffID
	trigger.AcknowledgedAt = &now
	trigger.Status = "WORK_ORDER_DISPATCHED"
	if err := repo.UpdatePredictiveTrigger(ctx, trigger); err != nil {
		t.Fatalf("failed to update trigger: %v", err)
	}

	// Update - Not Found
	tNon := trigger
	tNon.ID = "non-existent"
	if err := repo.UpdatePredictiveTrigger(ctx, tNon); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound on non-existent trigger, got %v", err)
	}

	// List
	list, err := repo.ListPredictiveTriggers(ctx, "WORK_ORDER_DISPATCHED", domain.ModelGeneratorPreventive, 10)
	if err != nil {
		t.Fatalf("failed to list triggers: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 trigger, got %d", len(list))
	}
}

func TestNewPostgresRepository(t *testing.T) {
	pg := repository.NewPostgresRepository(nil)
	if pg == nil {
		t.Fatal("expected non-nil PostgresRepository")
	}
}

