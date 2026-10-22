package repository

import (
	"context"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
)

func TestMemoryRepository_WorkOrderLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	wo := domain.MaintenanceWorkOrder{
		ID:            "wo-uuid-1",
		OrderNo:       "WO-20261022-0001",
		UnitType:      domain.UnitFeederMaint,
		WorkOrderType: domain.WOCorrectiveRepair,
		Priority:      domain.PriorityHigh,
		UnitHeadID:    "unit-head-1",
		TeamLeadID:    "team-lead-1",
		AssetID:       "asset-1",
		Status:        domain.StatusDispatched,
		Instructions:  "Clear fallen palm tree branch over 11kV line span.",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	// 1. Create
	if err := repo.CreateWorkOrder(ctx, wo); err != nil {
		t.Fatalf("CreateWorkOrder failed: %v", err)
	}

	// Duplicate create conflict
	if err := repo.CreateWorkOrder(ctx, wo); err != ErrConflict {
		t.Errorf("expected ErrConflict on duplicate, got: %v", err)
	}

	// 2. GetByID & GetByOrderNo
	fetched, err := repo.GetWorkOrderByID(ctx, "wo-uuid-1")
	if err != nil || fetched.OrderNo != wo.OrderNo {
		t.Fatalf("GetWorkOrderByID failed: %v, got: %+v", err, fetched)
	}

	byNo, err := repo.GetWorkOrderByOrderNo(ctx, "WO-20261022-0001")
	if err != nil || byNo.ID != wo.ID {
		t.Fatalf("GetWorkOrderByOrderNo failed: %v, got: %+v", err, byNo)
	}

	// Not found
	if _, err := repo.GetWorkOrderByID(ctx, "non-existent"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
	if _, err := repo.GetWorkOrderByOrderNo(ctx, "NO-MATCH"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}

	// 3. Update status & notes
	rework := "Inadequate photo clarity on pin insulator."
	if err := repo.UpdateWorkOrderStatus(ctx, "wo-uuid-1", domain.StatusRejectedRework, &rework); err != nil {
		t.Fatalf("UpdateWorkOrderStatus failed: %v", err)
	}
	updated, _ := repo.GetWorkOrderByID(ctx, "wo-uuid-1")
	if updated.Status != domain.StatusRejectedRework || updated.ReworkNotes == nil || *updated.ReworkNotes != rework {
		t.Fatalf("expected status %s and rework note, got %+v", domain.StatusRejectedRework, updated)
	}

	// Update non-existent
	if err := repo.UpdateWorkOrderStatus(ctx, "non-existent", domain.StatusCancelled, nil); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}

	// 4. List orders
	orders, err := repo.ListWorkOrders(ctx, domain.UnitFeederMaint, domain.StatusRejectedRework, "team-lead-1")
	if err != nil || len(orders) != 1 {
		t.Fatalf("expected 1 order, got %d, err: %v", len(orders), err)
	}

	// Filter with no match
	orders, err = repo.ListWorkOrders(ctx, domain.UnitTransformerMaint, "", "")
	if err != nil || len(orders) != 0 {
		t.Fatalf("expected 0 orders for transformer unit, got %d", len(orders))
	}
}

func TestMemoryRepository_ProofOfWorkLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	proof := domain.ProofOfWorkLog{
		ID:                      "proof-uuid-1",
		WorkOrderID:             "wo-uuid-1",
		TeamLeadID:              "team-lead-1",
		CapturedGPS:             s01domain.GPSPoint{Latitude: 9.9322, Longitude: 8.8921},
		DistanceFromAssetMeters: 23.5,
		BeforePhotoURL:          "/uploads/before_1.jpg",
		AfterPhotoURL:           "/uploads/after_1.jpg",
		PhotoHashBefore:         "1111111111111111111111111111111111111111111111111111111111111111",
		PhotoHashAfter:          "2222222222222222222222222222222222222222222222222222222222222222",
		WorkDetails:             "Re-insulated secondary jumper and tested line continuity.",
		OllamaGPSVerified:       true,
		AIRiskScore:             0.10,
		AIReasoning:             "GPS within 23.5m; photo hashes distinct.",
		SubmittedAt:             time.Now().UTC(),
	}

	// 1. Create
	if err := repo.CreateProofOfWork(ctx, proof); err != nil {
		t.Fatalf("CreateProofOfWork failed: %v", err)
	}

	// Duplicate create conflict
	if err := repo.CreateProofOfWork(ctx, proof); err != ErrConflict {
		t.Errorf("expected ErrConflict on duplicate proof, got: %v", err)
	}

	// 2. GetByWorkOrderID
	fetched, err := repo.GetProofByWorkOrderID(ctx, "wo-uuid-1")
	if err != nil || fetched.ID != proof.ID {
		t.Fatalf("GetProofByWorkOrderID failed: %v, got: %+v", err, fetched)
	}

	// Not found
	if _, err := repo.GetProofByWorkOrderID(ctx, "non-existent"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}

	// 3. Update Verification
	if err := repo.UpdateProofVerification(ctx, "proof-uuid-1", true, "Verified by Unit Head Engineer."); err != nil {
		t.Fatalf("UpdateProofVerification failed: %v", err)
	}
	verified, _ := repo.GetProofByWorkOrderID(ctx, "wo-uuid-1")
	if !verified.UnitHeadVerified || verified.UnitHeadNotes == nil || verified.VerifiedAt == nil {
		t.Fatalf("expected proof to be verified, got: %+v", verified)
	}

	// Update non-existent
	if err := repo.UpdateProofVerification(ctx, "non-existent", true, ""); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}

	// 4. ListProofsByTeamLead
	list, err := repo.ListProofsByTeamLead(ctx, "team-lead-1", 10)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 proof log, got %d, err: %v", len(list), err)
	}

	emptyList, err := repo.ListProofsByTeamLead(ctx, "other-team-lead", 10)
	if err != nil || len(emptyList) != 0 {
		t.Fatalf("expected 0 proofs, got %d", len(emptyList))
	}
}

