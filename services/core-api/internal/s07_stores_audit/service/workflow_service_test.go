// Package service_test exercises store inventory ledgers, dual-path material requisitions,
// and scrap return reverse logistics workflows for S07.
package service_test

import (
	"context"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/service"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const testHMACKey = "nesco-stores-hmac-32b-secret-test!"

func setupService() (*service.WorkflowService, *repository.MemoryRepository, *events.MemoryEventBus) {
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, bus, testHMACKey, "http://127.0.0.1:8001")
	return svc, repo, bus
}

func TestWorkflowService_StoreInventory(t *testing.T) {
	svc, _, _ := setupService()
	ctx := context.Background()

	item := domain.StoreItem{
		ItemCode:          "TX-OIL-DRUM",
		ItemName:          "Transformer Insulating Oil 200L",
		Category:          domain.CategoryConsumablesOils,
		QuantityAvailable: 25.0,
		ReorderThreshold:  5.0,
		UnitCost:          185000.0,
		BinLocation:       "YARD-OIL-01",
	}

	created, err := svc.RegisterStoreItem(ctx, item)
	if err != nil {
		t.Fatalf("unexpected error registering item: %v", err)
	}
	if created.RowHMACSignature == "" {
		t.Error("expected HMAC signature to be computed on item registration")
	}

	// Stock update (addition)
	updated, err := svc.UpdateInventoryStock(ctx, created.ID, 10.0, "store-op-1")
	if err != nil || updated.QuantityAvailable != 35.0 {
		t.Fatalf("stock addition failed: %+v, err: %v", updated, err)
	}

	// Stock reduction
	updated2, err := svc.UpdateInventoryStock(ctx, created.ID, -32.0, "store-op-1")
	if err != nil || updated2.QuantityAvailable != 3.0 {
		t.Fatalf("stock reduction failed: %+v, err: %v", updated2, err)
	}

	// Negative reduction exceeding available balance -> ErrInsufficientStock
	if _, err := svc.UpdateInventoryStock(ctx, created.ID, -10.0, "store-op-1"); err != service.ErrInsufficientStock {
		t.Fatalf("expected ErrInsufficientStock, got: %v", err)
	}

	// Query stock by code
	stockItem, err := svc.GetStockByItemCode(ctx, "TX-OIL-DRUM")
	if err != nil || stockItem.QuantityAvailable != 3.0 {
		t.Fatalf("GetStockByItemCode failed: %+v, err: %v", stockItem, err)
	}

	// List low stock (threshold is 5.0, current is 3.0)
	lowList, err := svc.ListInventory(ctx, domain.CategoryConsumablesOils, true)
	if err != nil || len(lowList) != 1 {
		t.Fatalf("expected 1 low stock item, got count: %d", len(lowList))
	}
}

func TestWorkflowService_RequisitionPathA_InStock(t *testing.T) {
	svc, _, _ := setupService()
	ctx := context.Background()

	// Seed In-Stock Item
	_, _ = svc.RegisterStoreItem(ctx, domain.StoreItem{
		ItemCode:          "FUSE-100A",
		ItemName:          "100A Drop-out Fuse Element",
		Category:          domain.CategorySubstationHardware,
		QuantityAvailable: 50.0,
		ReorderThreshold:  10.0,
		UnitCost:          3500.0,
		BinLocation:       "BIN-F10",
	})

	// Initiate Requisition
	req := domain.MaterialRequisition{
		OriginDeptID: "dept-services-feeder",
		UnitHeadID:   "unit-head-feeder",
		Items: []domain.RequisitionItem{
			{ItemCode: "FUSE-100A", QuantityRequested: 5.0},
		},
	}

	createdReq, err := svc.InitiateRequisition(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error initiating requisition: %v", err)
	}
	if !createdReq.InStockFlag {
		t.Fatal("expected InStockFlag to be true")
	}
	if createdReq.TotalCost != 17500.0 {
		t.Fatalf("expected total cost 17500.0, got: %f", createdReq.TotalCost)
	}

	// Step 2: Department Head Approval
	deptApproved, err := svc.ApproveRequisitionByDept(ctx, createdReq.ID, "dept-head-services")
	if err != nil || deptApproved.Status != domain.StatusRequisitionDeptApproved {
		t.Fatalf("dept approval failed: %+v, err: %v", deptApproved, err)
	}

	// Step 3: Internal Audit Consent
	auditConsented, err := svc.ConsentRequisitionByAudit(ctx, createdReq.ID, "auditor-01", true, "Checked ticket and prior scrap")
	if err != nil || auditConsented.Status != domain.StatusRequisitionAuditConsented {
		t.Fatalf("audit consent failed: %+v, err: %v", auditConsented, err)
	}

	// Step 4: Store Officer Issuance & Waybill
	issued, err := svc.IssueMaterialWaybill(ctx, createdReq.ID, "store-op-danladi")
	if err != nil || issued.Status != domain.StatusRequisitionIssued || issued.WaybillNo == nil {
		t.Fatalf("material issuance failed: %+v, err: %v", issued, err)
	}

	// Verify inventory was deducted from 50.0 to 45.0
	itemAfter, _ := svc.GetStockByItemCode(ctx, "FUSE-100A")
	if itemAfter.QuantityAvailable != 45.0 {
		t.Fatalf("expected 45.0 available in stock after issuance, got: %f", itemAfter.QuantityAvailable)
	}

	// Step 5: Reverse Logistics Scrap Return
	scrap := domain.ScrapReturn{
		RequisitionID:       createdReq.ID,
		TeamLeadID:          "lead-maintenance",
		StoreOfficerID:      "store-op-danladi",
		ItemsReturnedJSON:   `[{"item":"FUSE-100A-BLOWN","qty":5}]`,
		ScrapConditionNotes: "Blown fuse links collected from Bukuru feeder maintenance",
	}

	loggedScrap, err := svc.LogScrapReturn(ctx, scrap)
	if err != nil {
		t.Fatalf("unexpected error logging scrap return: %v", err)
	}
	if loggedScrap.ReturnNo == "" {
		t.Fatal("expected generated return number")
	}

	// Verify requisition scrap_returned_flag is now true
	reqFinal, _ := svc.GetRequisition(ctx, createdReq.ID)
	if !reqFinal.ScrapReturnedFlag || reqFinal.ScrapReturnID == nil {
		t.Fatalf("expected scrap_returned_flag to be true: %+v", reqFinal)
	}
}

func TestWorkflowService_RequisitionPathB_MDProcurement(t *testing.T) {
	svc, _, _ := setupService()
	ctx := context.Background()

	// Seed Item with 0 stock (Out of stock)
	_, _ = svc.RegisterStoreItem(ctx, domain.StoreItem{
		ItemCode:          "TX-BREAKER-33KV",
		ItemName:          "33kV Vacuum Circuit Breaker",
		Category:          domain.CategorySubstationHardware,
		QuantityAvailable: 0.0,
		ReorderThreshold:  2.0,
		UnitCost:          8500000.0,
		BinLocation:       "YARD-HV-01",
	})

	// Initiate Requisition
	req := domain.MaterialRequisition{
		OriginDeptID: "dept-generation-p1",
		UnitHeadID:   "unit-head-mech",
		Items: []domain.RequisitionItem{
			{ItemCode: "TX-BREAKER-33KV", QuantityRequested: 1.0, UnitCost: 8500000.0},
		},
	}

	createdReq, err := svc.InitiateRequisition(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if createdReq.InStockFlag {
		t.Fatal("expected InStockFlag to be false for zero inventory")
	}

	// Step 2: Generation GM Approval -> transitions to MD_PROCUREMENT_PENDING
	gmApproved, err := svc.ApproveRequisitionByDept(ctx, createdReq.ID, "gm-generation")
	if err != nil || gmApproved.Status != domain.StatusRequisitionMDProcurementPending {
		t.Fatalf("GM approval failed: %+v, err: %v", gmApproved, err)
	}

	// Negative: MD cannot approve before Audit Consent
	if _, err := svc.ApproveRequisitionByMD(ctx, createdReq.ID, "MD Approved"); err != domain.ErrAuditConsentRequired {
		t.Fatalf("expected ErrAuditConsentRequired, got: %v", err)
	}

	// Step 3a: Internal Audit Consents supplier quotes
	auditConsented, err := svc.ConsentRequisitionByAudit(ctx, createdReq.ID, "auditor-01", true, "Vendor pricing audited and compliant")
	if err != nil || auditConsented.AuditConsentStatus != domain.AuditConsentConsented {
		t.Fatalf("audit consent failed: %+v, err: %v", auditConsented, err)
	}

	// Step 3b: Managing Director Grants Final Procurement Approval
	mdApproved, err := svc.ApproveRequisitionByMD(ctx, createdReq.ID, "Procurement authorized for critical plant spares")
	if err != nil || mdApproved.Status != domain.StatusRequisitionMDApproved {
		t.Fatalf("MD approval failed: %+v, err: %v", mdApproved, err)
	}

	// Stock arrives from vendor, Store Officer adds stock
	item, _ := svc.GetStockByItemCode(ctx, "TX-BREAKER-33KV")
	_, _ = svc.UpdateInventoryStock(ctx, item.ID, 1.0, "store-op-1")

	// Step 4: Material Issuance & Waybill
	issued, err := svc.IssueMaterialWaybill(ctx, createdReq.ID, "store-op-1")
	if err != nil || issued.Status != domain.StatusRequisitionIssued {
		t.Fatalf("material issuance failed: %+v, err: %v", issued, err)
	}
}
