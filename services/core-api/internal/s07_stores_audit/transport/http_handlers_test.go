// Package transport_test exercises REST endpoints for stores inventory, requisitions,
// and reverse logistics scrap tracking for S07.
package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/service"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/transport"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const (
	testJWTKey  = "nesco-stores-jwt-key-32b-secret-test!"
	testHMACKey = "nesco-stores-hmac-key-32b-secret-test"
)

func generateTestToken(t *testing.T, staffID, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(auth.Claims{
		StaffID:     staffID,
		StaffNumber: "NES-" + staffID,
		FullName:    "Officer " + staffID,
		PrimaryRole: role,
		OrgUnitID:   "unit-stores-hq",
		IssuedAt:    time.Now().Unix(),
		ExpiresAt:   time.Now().Add(15 * time.Minute).Unix(),
	}, testJWTKey)
	if err != nil {
		t.Fatalf("failed to generate auth token: %v", err)
	}
	return token
}

func setupTestServer(t *testing.T) (*http.ServeMux, *repository.MemoryRepository) {
	t.Helper()
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, bus, testHMACKey, "http://127.0.0.1:8001")
	handler := transport.NewHTTPHandler(svc, repo, testJWTKey)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, repo
}

func TestHTTPHandler_StoreItems(t *testing.T) {
	mux, _ := setupTestServer(t)
	tokenStoreOp := generateTestToken(t, "store-op-01", "STORE_OFFICER")

	// 1. Unauthenticated -> 401
	unauthReq := httptest.NewRequest("POST", "/api/v1/stores/items", bytes.NewBufferString("{}"))
	wUnauth := httptest.NewRecorder()
	mux.ServeHTTP(wUnauth, unauthReq)
	if wUnauth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got: %d", wUnauth.Code)
	}

	// 2. Register Store Item -> 201 Created
	itemDTO := transport.CreateStoreItemDTO{
		ItemCode:          "CBL-33KV-150",
		ItemName:          "33kV XLPE 150mm Cable",
		Category:          domain.CategoryCablesConductors,
		QuantityAvailable: 100.0,
		ReorderThreshold:  20.0,
		UnitOfMeasure:     "METERS",
		UnitCost:          45000.0,
		BinLocation:       "YARD-A1",
	}
	body, _ := json.Marshal(itemDTO)

	req := httptest.NewRequest("POST", "/api/v1/stores/items", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got: %d, body: %s", w.Code, w.Body.String())
	}

	var created domain.StoreItem
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	if created.ItemCode != "CBL-33KV-150" || created.RowHMACSignature == "" {
		t.Fatalf("invalid registered item: %+v", created)
	}

	// 3. Live Stock Balance Check by Code -> 200 OK
	reqCode := httptest.NewRequest("GET", "/api/v1/stores/items/CBL-33KV-150", nil)
	reqCode.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	wCode := httptest.NewRecorder()
	mux.ServeHTTP(wCode, reqCode)

	if wCode.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for stock check, got: %d", wCode.Code)
	}

	// 4. Non-existent Item Check -> 404 Not Found
	reqNotFound := httptest.NewRequest("GET", "/api/v1/stores/items/NON-EXISTENT-SKU", nil)
	reqNotFound.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	wNotFound := httptest.NewRecorder()
	mux.ServeHTTP(wNotFound, reqNotFound)
	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unknown item, got: %d", wNotFound.Code)
	}

	// 5. Update Stock Balance -> 200 OK
	stockBody, _ := json.Marshal(transport.UpdateStockDTO{Delta: 50.0})
	reqStock := httptest.NewRequest("POST", "/api/v1/stores/items/"+created.ID+"/stock", bytes.NewReader(stockBody))
	reqStock.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	wStock := httptest.NewRecorder()
	mux.ServeHTTP(wStock, reqStock)

	if wStock.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on stock update, got: %d", wStock.Code)
	}

	// 6. Insufficient Stock Deduction -> 400 Bad Request
	stockExcess, _ := json.Marshal(transport.UpdateStockDTO{Delta: -500.0})
	reqExcess := httptest.NewRequest("POST", "/api/v1/stores/items/"+created.ID+"/stock", bytes.NewReader(stockExcess))
	reqExcess.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	wExcess := httptest.NewRecorder()
	mux.ServeHTTP(wExcess, reqExcess)
	if wExcess.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request on excessive stock deduction, got: %d", wExcess.Code)
	}

	// 7. List Items -> 200 OK
	reqList := httptest.NewRequest("GET", "/api/v1/stores/items?category=CABLES_CONDUCTORS", nil)
	reqList.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	wList := httptest.NewRecorder()
	mux.ServeHTTP(wList, reqList)
	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing items, got: %d", wList.Code)
	}
}

func TestHTTPHandler_RequisitionsAndScrap(t *testing.T) {
	mux, _ := setupTestServer(t)

	tokenUnitHead := generateTestToken(t, "unithead-01", "UNIT_HEAD")
	tokenDeptHead := generateTestToken(t, "depthead-01", "DEPT_HEAD")
	tokenAuditor := generateTestToken(t, "auditor-01", "INTERNAL_AUDITOR")
	tokenStoreOp := generateTestToken(t, "store-op-01", "STORE_OFFICER")

	// Seed Item in Stock
	itemDTO := transport.CreateStoreItemDTO{
		ItemCode:          "FUSE-100A",
		ItemName:          "100A Drop-out Fuse",
		Category:          domain.CategorySubstationHardware,
		QuantityAvailable: 40.0,
		ReorderThreshold:  10.0,
		UnitCost:          2500.0,
		BinLocation:       "BIN-F10",
	}
	bodyItem, _ := json.Marshal(itemDTO)
	reqSeed := httptest.NewRequest("POST", "/api/v1/stores/items", bytes.NewReader(bodyItem))
	reqSeed.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	wSeed := httptest.NewRecorder()
	mux.ServeHTTP(wSeed, reqSeed)

	// 1. Initiate Requisition -> 201 Created
	reqDTO := transport.InitiateRequisitionDTO{
		OriginDeptID: "dept-services",
		Items: []transport.RequisitionItemDTO{
			{ItemCode: "FUSE-100A", QuantityRequested: 4.0},
		},
	}
	bodyReq, _ := json.Marshal(reqDTO)

	reqInit := httptest.NewRequest("POST", "/api/v1/stores/requisitions", bytes.NewReader(bodyReq))
	reqInit.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wInit := httptest.NewRecorder()
	mux.ServeHTTP(wInit, reqInit)

	if wInit.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for requisition, got: %d, body: %s", wInit.Code, wInit.Body.String())
	}

	var createdReq domain.MaterialRequisition
	_ = json.Unmarshal(wInit.Body.Bytes(), &createdReq)
	if createdReq.Status != domain.StatusRequisitionSubmitted {
		t.Fatalf("expected SUBMITTED status, got: %s", createdReq.Status)
	}

	// 2. Department Head Approves -> 200 OK
	reqDept := httptest.NewRequest("POST", "/api/v1/stores/requisitions/"+createdReq.ID+"/dept-approve", nil)
	reqDept.Header.Set("Authorization", "Bearer "+tokenDeptHead)
	wDept := httptest.NewRecorder()
	mux.ServeHTTP(wDept, reqDept)

	if wDept.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on dept approval, got: %d", wDept.Code)
	}

	// 3. Internal Audit Grants Consent -> 200 OK
	auditDTO := transport.AuditConsentDTO{Consent: true, Notes: "Work Order #42 validated"}
	bodyAudit, _ := json.Marshal(auditDTO)
	reqAudit := httptest.NewRequest("POST", "/api/v1/stores/requisitions/"+createdReq.ID+"/audit-consent", bytes.NewReader(bodyAudit))
	reqAudit.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wAudit := httptest.NewRecorder()
	mux.ServeHTTP(wAudit, reqAudit)

	if wAudit.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on audit consent, got: %d", wAudit.Code)
	}

	// 4. Store Officer Issues Materials & Generates Waybill -> 200 OK
	reqIssue := httptest.NewRequest("POST", "/api/v1/stores/requisitions/"+createdReq.ID+"/issue", nil)
	reqIssue.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	wIssue := httptest.NewRecorder()
	mux.ServeHTTP(wIssue, reqIssue)

	if wIssue.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on material issuance, got: %d, body: %s", wIssue.Code, wIssue.Body.String())
	}

	var issued domain.MaterialRequisition
	_ = json.Unmarshal(wIssue.Body.Bytes(), &issued)
	if issued.Status != domain.StatusRequisitionIssued || issued.WaybillNo == nil {
		t.Fatalf("expected ISSUED status with waybill, got: %+v", issued)
	}

	// 5. Reverse Logistics Scrap Return -> 201 Created
	scrapDTO := transport.LogScrapReturnDTO{
		RequisitionID:       createdReq.ID,
		ItemsReturnedJSON:   `[{"item":"FUSE-100A-BLOWN","qty":4}]`,
		ScrapConditionNotes: "Damaged fuse barrels retrieved from transformer pole",
	}
	bodyScrap, _ := json.Marshal(scrapDTO)
	reqScrap := httptest.NewRequest("POST", "/api/v1/stores/scrap-returns", bytes.NewReader(bodyScrap))
	reqScrap.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	wScrap := httptest.NewRecorder()
	mux.ServeHTTP(wScrap, reqScrap)

	if wScrap.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on scrap return, got: %d, body: %s", wScrap.Code, wScrap.Body.String())
	}

	// 6. List Requisitions & Scrap -> 200 OK
	wListReq := httptest.NewRecorder()
	reqListReq := httptest.NewRequest("GET", "/api/v1/stores/requisitions?status=ISSUED", nil)
	reqListReq.Header.Set("Authorization", "Bearer "+tokenAuditor)
	mux.ServeHTTP(wListReq, reqListReq)
	if wListReq.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing requisitions, got: %d", wListReq.Code)
	}

	wListScrap := httptest.NewRecorder()
	reqListScrap := httptest.NewRequest("GET", "/api/v1/stores/scrap-returns", nil)
	reqListScrap.Header.Set("Authorization", "Bearer "+tokenStoreOp)
	mux.ServeHTTP(wListScrap, reqListScrap)
	if wListScrap.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing scrap returns, got: %d", wListScrap.Code)
	}
}

