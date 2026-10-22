// Package domain_test verifies C-Suite action console approval gates and segregation of duties.
package domain_test

import (
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

func TestValidateActionItemDecision(t *testing.T) {
	baseItem := domain.ExecutiveActionItem{
		ID:               "act-001",
		ActionType:       domain.ActionMDProcurementApproval,
		TargetRole:       domain.RoleMD,
		SourceEntityID:   "REQ-2026-0001",
		SourceEntityRef:  "PO-001",
		SourceDepartment: "STORES_AUDIT",
		Summary:          "Critical Out-of-Stock Transformer Procurement (₦4,500,000)",
		Amount:           4500000.0,
		Priority:         "CRITICAL",
		Status:           domain.ActionStatusPending,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	t.Run("successful approval by MD", func(t *testing.T) {
		item := baseItem
		err := domain.ValidateActionItemDecision(&item, "STAFF-MD-01", domain.RoleMD, true, "Approved for immediate vendor purchase")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.Status != domain.ActionStatusApproved {
			t.Fatalf("expected APPROVED, got %s", item.Status)
		}
		if item.DecidedByStaffID == nil || *item.DecidedByStaffID != "STAFF-MD-01" {
			t.Fatal("expected decided_by_staff_id to be populated")
		}
		if item.DecidedAt == nil {
			t.Fatal("expected decided_at to be populated")
		}
	})

	t.Run("successful rejection by MD", func(t *testing.T) {
		item := baseItem
		err := domain.ValidateActionItemDecision(&item, "STAFF-MD-01", domain.RoleMD, false, "Supplier quotation too high; renegotiate")
		if err != nil {
			t.Fatalf("unexpected rejection error: %v", err)
		}
		if item.Status != domain.ActionStatusRejected {
			t.Fatalf("expected REJECTED, got %s", item.Status)
		}
	})

	t.Run("role segregation enforcement (non-MD cannot approve MD action)", func(t *testing.T) {
		item := baseItem
		err := domain.ValidateActionItemDecision(&item, "STAFF-COMM-01", domain.RoleHeadCommercialOps, true, "Attempted override")
		if err == nil {
			t.Fatal("expected error on unauthorized role, got nil")
		}
	})

	t.Run("illegal state transition (cannot decide non-pending item)", func(t *testing.T) {
		item := baseItem
		item.Status = domain.ActionStatusApproved
		err := domain.ValidateActionItemDecision(&item, "STAFF-MD-01", domain.RoleMD, true, "Duplicate")
		if err == nil {
			t.Fatal("expected error on non-pending item, got nil")
		}
	})

	t.Run("missing mandatory decision parameters", func(t *testing.T) {
		item := baseItem
		err := domain.ValidateActionItemDecision(&item, "", domain.RoleMD, true, "Notes")
		if err == nil {
			t.Fatal("expected error on empty staff ID, got nil")
		}

		err = domain.ValidateActionItemDecision(&item, "STAFF-MD-01", domain.RoleMD, true, "   ")
		if err == nil {
			t.Fatal("expected error on empty notes, got nil")
		}
	})
}

