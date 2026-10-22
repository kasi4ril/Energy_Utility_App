// Package repository_test verifies the persistence operations and invariants
// of the S07 Stores & Internal Audit MemoryRepository.
package repository_test

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
)

func TestMemoryRepository_StoreItems(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	// 1. Not found
	if _, err := repo.GetStoreItemByID(ctx, "non-existent"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for item by ID, got: %v", err)
	}
	if _, err := repo.GetStoreItemByCode(ctx, "NON-EXISTENT"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for item by Code, got: %v", err)
	}

	item := domain.StoreItem{
		ID:                "item-01",
		ItemCode:          "CBL-33KV-150",
		ItemName:          "33kV XLPE Cable 150mm",
		Category:          domain.CategoryCablesConductors,
		QuantityAvailable: 1000.0,
		ReorderThreshold:  150.0,
		UnitOfMeasure:     "METERS",
		UnitCost:          45000.0,
		BinLocation:       "YARD-A1",
		LastUpdatedBy:     "store-op-1",
		RowHMACSignature:  "dummy-hmac-sig",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	// 2. Create
	if err := repo.CreateStoreItem(ctx, item); err != nil {
		t.Fatalf("unexpected error creating item: %v", err)
	}

	// Conflict on duplicate ID
	if err := repo.CreateStoreItem(ctx, item); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate ID, got: %v", err)
	}

	// Conflict on duplicate Code
	dupCode := item
	dupCode.ID = "item-02"
	if err := repo.CreateStoreItem(ctx, dupCode); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate code, got: %v", err)
	}

	// 3. Get by ID and Code
	byID, err := repo.GetStoreItemByID(ctx, "item-01")
	if err != nil || byID.ItemCode != "CBL-33KV-150" {
		t.Fatalf("mismatched item by ID: %+v, err: %v", byID, err)
	}

	byCode, err := repo.GetStoreItemByCode(ctx, "CBL-33KV-150")
	if err != nil || byCode.ID != "item-01" {
		t.Fatalf("mismatched item by Code: %+v, err: %v", byCode, err)
	}

	// 4. Update
	item.QuantityAvailable = 100.0 // below threshold (150.0)
	if err := repo.UpdateStoreItem(ctx, item); err != nil {
		t.Fatalf("unexpected error updating item: %v", err)
	}

	// Update non-existent returns ErrNotFound
	nonEx := item
	nonEx.ID = "item-999"
	if err := repo.UpdateStoreItem(ctx, nonEx); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound on non-existent item update, got: %v", err)
	}

	// 5. List items with lowStockOnly filter
	lowStockList, err := repo.ListStoreItems(ctx, domain.CategoryCablesConductors, true)
	if err != nil || len(lowStockList) != 1 {
		t.Fatalf("expected 1 low stock item, got count: %d, err: %v", len(lowStockList), err)
	}

	emptyCatList, err := repo.ListStoreItems(ctx, domain.CategorySafetyGear, false)
	if err != nil || len(emptyCatList) != 0 {
		t.Fatalf("expected 0 items for empty category, got count: %d", len(emptyCatList))
	}
}

func TestMemoryRepository_RequisitionsAndScrap(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	req := domain.MaterialRequisition{
		ID:           "req-101",
		ReqNo:        "REQ-2026-0001",
		OriginDeptID: "dept-services",
		UnitHeadID:   "unithead-1",
		Status:       domain.StatusRequisitionSubmitted,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := repo.CreateRequisition(ctx, req); err != nil {
		t.Fatalf("unexpected error creating requisition: %v", err)
	}

	// Conflict
	if err := repo.CreateRequisition(ctx, req); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate requisition ID, got: %v", err)
	}

	// Fetch by ID and ReqNo
	byID, err := repo.GetRequisitionByID(ctx, "req-101")
	if err != nil || byID.ReqNo != "REQ-2026-0001" {
		t.Fatalf("GetRequisitionByID failed: %+v, err: %v", byID, err)
	}
	byNo, err := repo.GetRequisitionByReqNo(ctx, "REQ-2026-0001")
	if err != nil || byNo.ID != "req-101" {
		t.Fatalf("GetRequisitionByReqNo failed: %+v, err: %v", byNo, err)
	}

	// Not found
	if _, err := repo.GetRequisitionByID(ctx, "req-999"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for unknown ID, got: %v", err)
	}

	// Update
	req.Status = domain.StatusRequisitionAuditConsented
	if err := repo.UpdateRequisition(ctx, req); err != nil {
		t.Fatalf("unexpected error updating requisition: %v", err)
	}

	// List
	list, err := repo.ListRequisitions(ctx, domain.StatusRequisitionAuditConsented, "unithead-1")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 requisition matching filters, got: %d", len(list))
	}

	// Scrap Return
	scrap := domain.ScrapReturn{
		ID:                  "scrap-01",
		ReturnNo:            "SCRAP-2026-0001",
		RequisitionID:       "req-101",
		TeamLeadID:          "lead-danladi",
		StoreOfficerID:      "store-op-1",
		ItemsReturnedJSON:   `[{"item":"CBL-BURNT","qty":150}]`,
		ScrapConditionNotes: "Burnt cable sheath from line fault",
		ReturnedAt:          time.Now().UTC(),
		CreatedAt:           time.Now().UTC(),
	}

	if err := repo.CreateScrapReturn(ctx, scrap); err != nil {
		t.Fatalf("unexpected error creating scrap return: %v", err)
	}
	if err := repo.CreateScrapReturn(ctx, scrap); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate scrap ID, got: %v", err)
	}

	sByID, err := repo.GetScrapReturnByID(ctx, "scrap-01")
	if err != nil || sByID.ReturnNo != "SCRAP-2026-0001" {
		t.Fatalf("GetScrapReturnByID failed: %+v, err: %v", sByID, err)
	}

	scrapList, err := repo.ListScrapReturns(ctx, "lead-danladi")
	if err != nil || len(scrapList) != 1 {
		t.Fatalf("expected 1 scrap return for lead-danladi, got %d", len(scrapList))
	}
}

func TestMemoryRepository_StaffClaims(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	// Not found
	if _, err := repo.GetClaimByID(ctx, "claim-none"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for claim, got: %v", err)
	}

	claim := domain.StaffAllowanceClaim{
		ID:        "claim-01",
		ClaimNo:   "CLM-2026-0001",
		StaffID:   "staff-101",
		ClaimType: domain.ClaimOvertime,
		Amount:    15000.0,
		ClaimDate: "2026-10-22",
		Status:    domain.StatusClaimSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := repo.CreateClaim(ctx, claim); err != nil {
		t.Fatalf("unexpected error creating claim: %v", err)
	}
	if err := repo.CreateClaim(ctx, claim); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate claim ID, got: %v", err)
	}

	got, err := repo.GetClaimByID(ctx, "claim-01")
	if err != nil || got.ClaimNo != "CLM-2026-0001" {
		t.Fatalf("GetClaimByID failed: %+v, err: %v", got, err)
	}

	// Update
	claim.Status = domain.StatusClaimAuditApproved
	if err := repo.UpdateClaim(ctx, claim); err != nil {
		t.Fatalf("unexpected error updating claim: %v", err)
	}

	// Non-existent update
	nonEx := claim
	nonEx.ID = "claim-none"
	if err := repo.UpdateClaim(ctx, nonEx); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound updating unknown claim, got: %v", err)
	}

	list, err := repo.ListClaims(ctx, "staff-101", domain.StatusClaimAuditApproved)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 approved claim, got: %d", len(list))
	}
}

