// Package repository_test verifies the persistence operations for billing adjustments and AI alerts.
package repository_test

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
)

func TestMemoryRepository_BillingAdjustments(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()

	adj := domain.BillingAdjustment{
		ID:                "adj-001",
		AdjustmentNo:      "ADJ-2026-0001",
		TicketID:          "TICK-001",
		CustomerID:        "CUST-001",
		BillingOfficerID:  "STAFF-BILL-01",
		OldBalance:        20000.0,
		AdjustedBalance:   12000.0,
		AdjustmentAmount:  8000.0,
		RootCauseCategory: domain.RootCauseMeterReadingTypo,
		RootCauseNote:     "Typo in initial reading",
		Status:            domain.AdjustmentProposed,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	// Create
	if err := repo.CreateBillingAdjustment(ctx, adj); err != nil {
		t.Fatalf("failed to create billing adjustment: %v", err)
	}

	// Conflict ID
	if err := repo.CreateBillingAdjustment(ctx, adj); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate ID, got %v", err)
	}

	// Conflict No
	a2 := adj
	a2.ID = "adj-002"
	if err := repo.CreateBillingAdjustment(ctx, a2); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate adjustment_no, got %v", err)
	}

	// Get by ID
	fetched, err := repo.GetBillingAdjustmentByID(ctx, "adj-001")
	if err != nil {
		t.Fatalf("failed to get adjustment by ID: %v", err)
	}
	if fetched.AdjustmentNo != "ADJ-2026-0001" {
		t.Fatalf("expected ADJ-2026-0001, got %s", fetched.AdjustmentNo)
	}

	// Get by ID - Not Found
	if _, err := repo.GetBillingAdjustmentByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Get by No
	byNo, err := repo.GetBillingAdjustmentByNo(ctx, "ADJ-2026-0001")
	if err != nil {
		t.Fatalf("failed to get by No: %v", err)
	}
	if byNo.ID != "adj-001" {
		t.Fatalf("expected adj-001, got %s", byNo.ID)
	}

	// Get by No - Not Found
	if _, err := repo.GetBillingAdjustmentByNo(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Update
	headID := "HEAD-COMM-01"
	adj.HeadCommercialOpsID = &headID
	adj.HeadCommercialOpsApproved = true
	adj.Status = domain.AdjustmentAppliedRevertedToCare
	if err := repo.UpdateBillingAdjustment(ctx, adj); err != nil {
		t.Fatalf("failed to update adjustment: %v", err)
	}

	// Update - Not Found
	nonExistent := adj
	nonExistent.ID = "non-existent"
	if err := repo.UpdateBillingAdjustment(ctx, nonExistent); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// List
	list, err := repo.ListBillingAdjustments(ctx, domain.AdjustmentAppliedRevertedToCare, "CUST-001")
	if err != nil {
		t.Fatalf("failed to list adjustments: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 adjustment, got %d", len(list))
	}

	empty, err := repo.ListBillingAdjustments(ctx, domain.AdjustmentProposed, "")
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected 0, got %d", len(empty))
	}
}

func TestMemoryRepository_AIAlerts(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()

	alert := domain.RevenueProtectionAlert{
		ID:              "alert-001",
		AlertNo:         "ALR-2026-0001",
		CustomerOrDTID:  "CUST-001",
		EntityType:      "CUSTOMER",
		AlertType:       domain.AlertEnergyTheftBypass,
		ConfidenceScore: 0.88,
		EvidenceSummary: "Feeder energized but vending dropped > 50%",
		Status:          domain.AlertOpen,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	// Create
	if err := repo.CreateAIAlert(ctx, alert); err != nil {
		t.Fatalf("failed to create AI alert: %v", err)
	}

	// Conflict ID
	if err := repo.CreateAIAlert(ctx, alert); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate ID, got %v", err)
	}

	// Conflict No
	al2 := alert
	al2.ID = "alert-002"
	if err := repo.CreateAIAlert(ctx, al2); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate alert_no, got %v", err)
	}

	// Get by ID
	fetched, err := repo.GetAIAlertByID(ctx, "alert-001")
	if err != nil {
		t.Fatalf("failed to get AI alert by ID: %v", err)
	}
	if fetched.AlertNo != "ALR-2026-0001" {
		t.Fatalf("expected ALR-2026-0001, got %s", fetched.AlertNo)
	}

	// Get by ID - Not Found
	if _, err := repo.GetAIAlertByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Get by No
	byNo, err := repo.GetAIAlertByNo(ctx, "ALR-2026-0001")
	if err != nil {
		t.Fatalf("failed to get by No: %v", err)
	}
	if byNo.ID != "alert-001" {
		t.Fatalf("expected alert-001, got %s", byNo.ID)
	}

	// Get by No - Not Found
	if _, err := repo.GetAIAlertByNo(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Update
	alert.Status = domain.AlertInvestigationDispatched
	unit := "UNIT-REV-PROT-01"
	alert.AssignedUnitID = &unit
	if err := repo.UpdateAIAlert(ctx, alert); err != nil {
		t.Fatalf("failed to update AI alert: %v", err)
	}

	// Update - Not Found
	nonExistent := alert
	nonExistent.ID = "non-existent"
	if err := repo.UpdateAIAlert(ctx, nonExistent); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// List
	list, err := repo.ListAIAlerts(ctx, domain.AlertInvestigationDispatched, domain.AlertEnergyTheftBypass)
	if err != nil {
		t.Fatalf("failed to list alerts: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(list))
	}

	empty, err := repo.ListAIAlerts(ctx, domain.AlertOpen, "")
	if err != nil {
		t.Fatalf("failed to list empty: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected 0, got %d", len(empty))
	}
}

