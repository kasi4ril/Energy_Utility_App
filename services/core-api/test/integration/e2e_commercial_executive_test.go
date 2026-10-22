// Package integration exercises commercial vending, refund amortization, audit and executive command.
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	s04domain "disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
	s04repo "disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	s05repo "disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	s06repo "disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	s07domain "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	s07repo "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	s07svc "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/service"
	s07trans "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/transport"
	s08domain "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	s08repo "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	s08svc "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/service"
	s08trans "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/transport"
	s09domain "disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	s09repo "disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
	s09svc "disco-enterprise-platform/services/core-api/internal/s09_executive_command/service"
	s09trans "disco-enterprise-platform/services/core-api/internal/s09_executive_command/transport"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

func setupCommercialExecEnvironment(t *testing.T) (*http.ServeMux, *s07repo.MemoryRepository, *s08repo.MemoryRepository, *s09repo.MemoryRepository, *s04repo.MemoryRepository) {
	t.Helper()
	bus := events.NewMemoryEventBus()

	s03Repo := s03repo.NewMemoryRepository()
	s04Repo := s04repo.NewMemoryRepository()
	s05Repo := s05repo.NewMemoryRepository()
	s06Repo := s06repo.NewMemoryRepository()
	s07Repo := s07repo.NewMemoryRepository()
	s08Repo := s08repo.NewMemoryRepository()
	s09Repo := s09repo.NewMemoryRepository()

	s07Service := s07svc.NewWorkflowService(s07Repo, bus, testHMACKey, "")
	s08Service := s08svc.NewWorkflowService(s08Repo, bus, s03Repo)
	s09Service := s09svc.NewWorkflowService(s09Repo, bus, s04Repo, s05Repo, s06Repo, s07Repo, s08Repo)

	mux := http.NewServeMux()
	s07trans.NewHTTPHandler(s07Service, s07Repo, testJWTKey).RegisterRoutes(mux)
	s08trans.NewHTTPHandler(s08Service, s08Repo, testJWTKey).RegisterRoutes(mux)
	s09trans.NewHTTPHandler(s09Service, s09Repo, testJWTKey).RegisterRoutes(mux)

	return mux, s07Repo, s08Repo, s09Repo, s04Repo
}

func TestE2E_MeteredVending_And_RefundAmortization_Lifecycle(t *testing.T) {
	mux, _, s08Repo, _, _ := setupCommercialExecEnvironment(t)
	token := auth.Claims{
		StaffID:     "STAFF-BILLING-01",
		StaffNumber: "NES-BILL-01",
		FullName:    "Commercial Officer",
		PrimaryRole: "UNIT_HEAD_BILLING",
		OrgUnitID:   "unit-commercial-hq",
	}
	jwtStr, _ := auth.GenerateToken(token, testJWTKey)
	authHeader := "Bearer " + jwtStr

	// Seed customer meter refund amortization ledger (NERC MAP Reimbursement)
	customerID := "CUST-KURA-2026-001"
	meterNo := "MTR-98214"
	err := s08Repo.CreateRefundLedger(context.Background(), s08domain.MeterRefundLedger{
		ID:                     "REF-LEDGER-01",
		CustomerID:             customerID,
		MeterSerial:            meterNo,
		InitialMeterCost:       85000.0,
		RefundPercentagePerTx:  10.0,
		TotalRefundedToDate:    0.0,
		RemainingRefundBalance: 85000.0,
		Status:                 s08domain.RefundActive,
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to seed meter refund ledger: %v", err)
	}

	// Ingest Metered Vending Transaction with 10% refund credit
	vendPayload, _ := json.Marshal(map[string]interface{}{
		"tx_reference":             "OPAY-TX-99881122",
		"customer_id":              customerID,
		"account_number":           "ACC-00142",
		"meter_number":             meterNo,
		"tariff_rate_per_kwh":      68.5,
		"amount_paid":              20000.0,
		"payment_platform_source":  "OPAY",
		"payment_channel_category": s08domain.ChannelFintechAgency,
		"settlement_verified":      true,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/vending/vend", bytes.NewBuffer(vendPayload))
	req.Header.Set("Authorization", authHeader)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for vending, got %d: %s", rec.Code, rec.Body.String())
	}
	var vendTx s08domain.VendingTransaction
	_ = json.Unmarshal(rec.Body.Bytes(), &vendTx)

	// Verify zero-estimation kWh derivation and 20-digit STS token
	if len(vendTx.STSToken) != 24 { // 20 digits + 4 hyphens
		t.Fatalf("expected 24-char formatted STS token, got %s", vendTx.STSToken)
	}
	if vendTx.MeterRefundCredited != 2000.0 { // 10% of 20,000
		t.Fatalf("expected 2000.0 meter refund credited, got %f", vendTx.MeterRefundCredited)
	}

	// Verify refund ledger balance decremented
	updatedLedger, err := s08Repo.GetRefundLedgerByMeter(context.Background(), meterNo)
	if err != nil || updatedLedger.RemainingRefundBalance != 83000.0 {
		t.Fatalf("expected remaining balance 83000.0, got %f", updatedLedger.RemainingRefundBalance)
	}
}

func TestE2E_ExecutiveCommand_And_PredictiveMaintenance_To_S04(t *testing.T) {
	mux, _, _, s09Repo, s04Repo := setupCommercialExecEnvironment(t)
	mdToken := auth.Claims{
		StaffID:                    "STAFF-MD-01",
		StaffNumber:                "NES-MD-01",
		FullName:                   "Managing Director",
		PrimaryRole:                "MD",
		OrgUnitID:                  "unit-exec-hq",
		CanViewExecutiveDashboard: true,
	}
	mdJwt, _ := auth.GenerateToken(mdToken, testJWTKey)
	mdAuth := "Bearer " + mdJwt

	// Step 1: C-Suite Live Intelligence Dashboard Query (S09)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/executive/dashboard", nil)
	req.Header.Set("Authorization", mdAuth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for C-Suite dashboard, got %d: %s", rec.Code, rec.Body.String())
	}
	var summary s09domain.CSuiteDashboardSummary
	_ = json.Unmarshal(rec.Body.Bytes(), &summary)

	// Step 2: Enqueue and Decide High-Value MD Action Item (S09)
	actionPayload, _ := json.Marshal(map[string]interface{}{
		"action_type":       s09domain.ActionMDProcurementApproval,
		"target_role":       s09domain.RoleMD,
		"source_entity_id":  "REQ-STORES-2026-991",
		"source_entity_ref": "PO-TX-OIL-DRUMS-50",
		"source_department": "STORES_AUDIT",
		"summary":           "Out-of-stock emergency procurement of 50 drums transformer dielectric oil",
		"amount":            12500000.0,
		"priority":          "CRITICAL",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions", bytes.NewBuffer(actionPayload))
	req.Header.Set("Authorization", mdAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for action item enqueue, got %d: %s", rec.Code, rec.Body.String())
	}
	var createdAction s09domain.ExecutiveActionItem
	_ = json.Unmarshal(rec.Body.Bytes(), &createdAction)

	// MD signs off action item
	decidePayload, _ := json.Marshal(map[string]interface{}{
		"decided_by_staff_id": "STAFF-MD-01",
		"approve":             true,
		"notes":               "Approved for immediate bank letter of credit issuance.",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions/"+createdAction.ID+"/decide", bytes.NewBuffer(decidePayload))
	req.Header.Set("Authorization", mdAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for MD digital sign-off, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify action item approved in repository
	decidedAction, _ := s09Repo.GetActionItemByID(context.Background(), createdAction.ID)
	if decidedAction.Status != s09domain.ActionStatusApproved {
		t.Fatalf("expected status APPROVED, got %s", decidedAction.Status)
	}

	// Step 3: Run AI Condition-Based Maintenance (CBM) Generator Prognosis (S09)
	// Critical mechanical vibration > 7.1 mm/s triggers automated S04 work order creation
	genPrognosisPayload, _ := json.Marshal(map[string]interface{}{
		"asset_id":                 "GEN-UNIT-04-TURBINE",
		"bearing_vibration_mm_s":   8.5, // Zone D Danger (> 7.1)
		"winding_temp_c":          96.5,
		"lube_oil_pressure_bar":   1.4,
		"run_hours_since_overhaul": 8500.0,
		"unit_head_id":            "HEAD-GEN-01",
		"team_lead_id":            "LEAD-MECH-01",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/generator", bytes.NewBuffer(genPrognosisPayload))
	req.Header.Set("Authorization", mdAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for generator prognosis, got %d: %s", rec.Code, rec.Body.String())
	}
	var trigger s09domain.PredictiveMaintenanceTrigger
	_ = json.Unmarshal(rec.Body.Bytes(), &trigger)

	if !trigger.AutoWorkOrderCreated || trigger.LinkedWorkOrderID == nil {
		t.Fatalf("expected auto work order created for critical generator vibration")
	}

	// Step 4: Verify Work Order was autonomously created inside S04 Services & Maintenance repository
	wos, err := s04Repo.ListWorkOrders(context.Background(), "", "", "")
	if err != nil || len(wos) != 1 {
		t.Fatalf("expected 1 work order created in S04 from S09 predictive trigger, got %d (err: %v)", len(wos), err)
	}
	if wos[0].WorkOrderType != s04domain.WOPreventiveMaintenance || wos[0].Priority != s04domain.PriorityCritical {
		t.Fatalf("unexpected work order parameters: %+v", wos[0])
	}
}

func TestE2E_StoreStockLedger_And_DualPathRequisitions(t *testing.T) {
	mux, s07Repo, _, _, _ := setupCommercialExecEnvironment(t)
	token := auth.Claims{
		StaffID:     "STAFF-STORE-01",
		StaffNumber: "NES-STR-01",
		FullName:    "Store Officer",
		PrimaryRole: "STORE_OFFICER",
		OrgUnitID:   "unit-stores-bukuru",
	}
	jwtStr, _ := auth.GenerateToken(token, testJWTKey)
	authHeader := "Bearer " + jwtStr

	// Seed store inventory
	itemCode := "SKU-DIELECTRIC-OIL-DRUM"
	err := s07Repo.CreateStoreItem(context.Background(), s07domain.StoreItem{
		ID:                "ITEM-01",
		ItemCode:          itemCode,
		ItemName:          "Transformer Dielectric Mineral Oil (200L Drum)",
		Category:          s07domain.CategoryConsumablesOils,
		QuantityAvailable: 25,
		ReorderThreshold:  10,
		UnitCost:          250000.0,
		UnitOfMeasure:     "DRUM",
		BinLocation:       "Bukuru Yard Bay A4",
		RowHMACSignature:  "valid-item-hmac",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to seed store inventory: %v", err)
	}

	// Query store inventory (@GridAssistant stock search)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stores/items", nil)
	req.Header.Set("Authorization", authHeader)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for store inventory search, got %d: %s", rec.Code, rec.Body.String())
	}
	var items []s07domain.StoreItem
	_ = json.Unmarshal(rec.Body.Bytes(), &items)
	if len(items) != 1 || items[0].ItemCode != itemCode {
		t.Fatalf("expected item returned from inventory search, got %+v", items)
	}

	// Submit Path A In-Stock Material Requisition
	reqPayload, _ := json.Marshal(map[string]interface{}{
		"origin_dept_id": "DEPT-MAINT-01",
		"items": []map[string]interface{}{
			{
				"item_code":          itemCode,
				"quantity_requested": 4.0,
				"unit_cost":          250000.0,
			},
		},
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/stores/requisitions", bytes.NewBuffer(reqPayload))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for requisition, got %d: %s", rec.Code, rec.Body.String())
	}
	var createdReq s07domain.MaterialRequisition
	_ = json.Unmarshal(rec.Body.Bytes(), &createdReq)

	if !createdReq.InStockFlag {
		t.Fatalf("expected Path A in-stock requisition since stock (25) >= requested (4)")
	}
}
