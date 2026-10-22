// Package repository_test exercises persistence operations for the C-Suite executive action console.
package repository_test

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
)

func TestMemoryRepository_ActionItems(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()

	item := domain.ExecutiveActionItem{
		ID:               "act-001",
		ActionType:       domain.ActionMDProcurementApproval,
		TargetRole:       domain.RoleMD,
		SourceEntityID:   "REQ-001",
		SourceEntityRef:  "PO-2026-0001",
		SourceDepartment: "STORES_AUDIT",
		Summary:          "Critical Out-of-Stock Transformer Procurement (₦4,500,000)",
		Amount:           4500000.0,
		Priority:         "CRITICAL",
		Status:           domain.ActionStatusPending,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	// Create
	if err := repo.CreateActionItem(ctx, item); err != nil {
		t.Fatalf("failed to create action item: %v", err)
	}

	// Conflict ID
	if err := repo.CreateActionItem(ctx, item); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate ID, got %v", err)
	}

	// Conflict Ref
	i2 := item
	i2.ID = "act-002"
	if err := repo.CreateActionItem(ctx, i2); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate sourceEntityRef, got %v", err)
	}

	// Get by ID
	fetched, err := repo.GetActionItemByID(ctx, "act-001")
	if err != nil {
		t.Fatalf("failed to get action item by ID: %v", err)
	}
	if fetched.SourceEntityRef != "PO-2026-0001" {
		t.Fatalf("expected PO-2026-0001, got %s", fetched.SourceEntityRef)
	}

	// Get by ID - Not Found
	if _, err := repo.GetActionItemByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Get by Ref
	byRef, err := repo.GetActionItemByRef(ctx, "PO-2026-0001")
	if err != nil {
		t.Fatalf("failed to get by ref: %v", err)
	}
	if byRef.ID != "act-001" {
		t.Fatalf("expected act-001, got %s", byRef.ID)
	}

	// Get by Ref - Not Found
	if _, err := repo.GetActionItemByRef(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Update
	staffID := "STAFF-MD-01"
	now := time.Now().UTC()
	notes := "Approved and signed"
	item.DecidedByStaffID = &staffID
	item.DecidedAt = &now
	item.DecisionNotes = &notes
	item.Status = domain.ActionStatusApproved

	if err := repo.UpdateActionItem(ctx, item); err != nil {
		t.Fatalf("failed to update action item: %v", err)
	}

	// Update - Not Found
	nonExistent := item
	nonExistent.ID = "non-existent"
	if err := repo.UpdateActionItem(ctx, nonExistent); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound on non-existent item, got %v", err)
	}

	// List
	list, err := repo.ListActionItems(ctx, domain.RoleMD, domain.ActionStatusApproved)
	if err != nil {
		t.Fatalf("failed to list action items: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 action item in list, got %d", len(list))
	}

	empty, err := repo.ListActionItems(ctx, domain.RoleHeadServices, domain.ActionStatusPending)
	if err != nil {
		t.Fatalf("failed to list empty: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected 0 items, got %d", len(empty))
	}
}

