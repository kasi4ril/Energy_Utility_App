// Package transport_test exercises Path B out-of-stock requisitions, MD approval, and error branches for S07.
package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/transport"
)

func TestHTTPHandler_RequisitionsPathB_MDProcurement(t *testing.T) {
	mux, _ := setupTestServer(t)

	tokenUnitHead := generateTestToken(t, "unithead-02", "UNIT_HEAD")
	tokenDeptHead := generateTestToken(t, "depthead-02", "DEPT_HEAD")
	tokenAuditor := generateTestToken(t, "auditor-02", "INTERNAL_AUDITOR")
	tokenMD := generateTestToken(t, "md-executive-01", "MANAGING_DIRECTOR")

	// 1. Initiate Out-of-Stock Requisition
	reqDTO := transport.InitiateRequisitionDTO{
		OriginDeptID: "dept-generation-plant1",
		Items: []transport.RequisitionItemDTO{
			{ItemCode: "TURBINE-SEAL-KIT", QuantityRequested: 2.0, UnitCost: 1500000.0},
		},
	}
	bodyReq, _ := json.Marshal(reqDTO)
	reqInit := httptest.NewRequest("POST", "/api/v1/stores/requisitions", bytes.NewReader(bodyReq))
	reqInit.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wInit := httptest.NewRecorder()
	mux.ServeHTTP(wInit, reqInit)

	if wInit.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got: %d", wInit.Code)
	}

	var req domain.MaterialRequisition
	_ = json.Unmarshal(wInit.Body.Bytes(), &req)

	// 2. Query Requisition by ID -> 200 OK
	reqGet := httptest.NewRequest("GET", "/api/v1/stores/requisitions/"+req.ID, nil)
	reqGet.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, reqGet)
	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GetRequisition, got: %d", wGet.Code)
	}

	// 3. Query Non-existent Requisition -> 404 Not Found
	reqGetNonEx := httptest.NewRequest("GET", "/api/v1/stores/requisitions/non-existent-req-id", nil)
	reqGetNonEx.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wGetNonEx := httptest.NewRecorder()
	mux.ServeHTTP(wGetNonEx, reqGetNonEx)
	if wGetNonEx.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent requisition, got: %d", wGetNonEx.Code)
	}

	// 4. Dept Head Approves -> 200 OK
	reqDept := httptest.NewRequest("POST", "/api/v1/stores/requisitions/"+req.ID+"/dept-approve", nil)
	reqDept.Header.Set("Authorization", "Bearer "+tokenDeptHead)
	wDept := httptest.NewRecorder()
	mux.ServeHTTP(wDept, reqDept)
	if wDept.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on dept approval, got: %d", wDept.Code)
	}

	// 5. Internal Audit Grants Consent -> 200 OK
	auditDTO := transport.AuditConsentDTO{Consent: true, Notes: "Procurement price verified"}
	bodyAudit, _ := json.Marshal(auditDTO)
	reqAudit := httptest.NewRequest("POST", "/api/v1/stores/requisitions/"+req.ID+"/audit-consent", bytes.NewReader(bodyAudit))
	reqAudit.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wAudit := httptest.NewRecorder()
	mux.ServeHTTP(wAudit, reqAudit)
	if wAudit.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on audit consent, got: %d", wAudit.Code)
	}

	// 6. Managing Director Grants Procurement Approval -> 200 OK
	mdDTO := transport.MDApprovalDTO{Notes: "Approved for supplier purchase order"}
	bodyMD, _ := json.Marshal(mdDTO)
	reqMD := httptest.NewRequest("POST", "/api/v1/stores/requisitions/"+req.ID+"/md-approve", bytes.NewReader(bodyMD))
	reqMD.Header.Set("Authorization", "Bearer "+tokenMD)
	wMD := httptest.NewRecorder()
	mux.ServeHTTP(wMD, reqMD)
	if wMD.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on MD approval, got: %d", wMD.Code)
	}

	// 7. Negative validations
	// Bad JSON on initiate
	reqBadInit := httptest.NewRequest("POST", "/api/v1/stores/requisitions", bytes.NewBufferString("not-json"))
	reqBadInit.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wBadInit := httptest.NewRecorder()
	mux.ServeHTTP(wBadInit, reqBadInit)
	if wBadInit.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on bad JSON, got: %d", wBadInit.Code)
	}

	// Validation error on initiate (empty origin dept)
	reqValErr := httptest.NewRequest("POST", "/api/v1/stores/requisitions", bytes.NewBufferString(`{"origin_dept_id":""}`))
	reqValErr.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wValErr := httptest.NewRecorder()
	mux.ServeHTTP(wValErr, reqValErr)
	if wValErr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on validation error, got: %d", wValErr.Code)
	}

	// Bad JSON on scrap return
	reqBadScrap := httptest.NewRequest("POST", "/api/v1/stores/scrap-returns", bytes.NewBufferString("not-json"))
	reqBadScrap.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wBadScrap := httptest.NewRecorder()
	mux.ServeHTTP(wBadScrap, reqBadScrap)
	if wBadScrap.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on bad JSON scrap, got: %d", wBadScrap.Code)
	}

	// Bad JSON on item register
	reqBadItem := httptest.NewRequest("POST", "/api/v1/stores/items", bytes.NewBufferString("not-json"))
	reqBadItem.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wBadItem := httptest.NewRecorder()
	mux.ServeHTTP(wBadItem, reqBadItem)
	if wBadItem.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON item, got: %d", wBadItem.Code)
	}

	// Validation error on item register
	reqValItem := httptest.NewRequest("POST", "/api/v1/stores/items", bytes.NewBufferString(`{"item_code":""}`))
	reqValItem.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wValItem := httptest.NewRecorder()
	mux.ServeHTTP(wValItem, reqValItem)
	if wValItem.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for item validation error, got: %d", wValItem.Code)
	}

	// Bad JSON on stock update
	reqBadStock := httptest.NewRequest("POST", "/api/v1/stores/items/item-01/stock", bytes.NewBufferString("not-json"))
	reqBadStock.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wBadStock := httptest.NewRecorder()
	mux.ServeHTTP(wBadStock, reqBadStock)
	if wBadStock.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON stock update, got: %d", wBadStock.Code)
	}

	// Bad JSON on audit consent
	reqBadConsent := httptest.NewRequest("POST", "/api/v1/stores/requisitions/"+req.ID+"/audit-consent", bytes.NewBufferString("not-json"))
	reqBadConsent.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wBadConsent := httptest.NewRecorder()
	mux.ServeHTTP(wBadConsent, reqBadConsent)
	if wBadConsent.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON audit consent, got: %d", wBadConsent.Code)
	}

	// Illegal Dept approval on non-existent req
	reqDeptNonEx := httptest.NewRequest("POST", "/api/v1/stores/requisitions/non-existent/dept-approve", nil)
	reqDeptNonEx.Header.Set("Authorization", "Bearer "+tokenDeptHead)
	wDeptNonEx := httptest.NewRecorder()
	mux.ServeHTTP(wDeptNonEx, reqDeptNonEx)
	if wDeptNonEx.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for non-existent dept approve, got: %d", wDeptNonEx.Code)
	}

	// Illegal Audit consent on non-existent req
	reqAuditNonEx := httptest.NewRequest("POST", "/api/v1/stores/requisitions/non-existent/audit-consent", bytes.NewReader(bodyAudit))
	reqAuditNonEx.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wAuditNonEx := httptest.NewRecorder()
	mux.ServeHTTP(wAuditNonEx, reqAuditNonEx)
	if wAuditNonEx.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for non-existent audit consent, got: %d", wAuditNonEx.Code)
	}

	// Illegal MD approval on non-existent req
	reqMDNonEx := httptest.NewRequest("POST", "/api/v1/stores/requisitions/non-existent/md-approve", bytes.NewReader(bodyMD))
	reqMDNonEx.Header.Set("Authorization", "Bearer "+tokenMD)
	wMDNonEx := httptest.NewRecorder()
	mux.ServeHTTP(wMDNonEx, reqMDNonEx)
	if wMDNonEx.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for non-existent MD approval, got: %d", wMDNonEx.Code)
	}

	// Illegal issuance on non-existent req
	reqIssueNonEx := httptest.NewRequest("POST", "/api/v1/stores/requisitions/non-existent/issue", nil)
	reqIssueNonEx.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wIssueNonEx := httptest.NewRecorder()
	mux.ServeHTTP(wIssueNonEx, reqIssueNonEx)
	if wIssueNonEx.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for non-existent issuance, got: %d", wIssueNonEx.Code)
	}

	// Scrap return validation error (empty notes)
	reqValScrap := httptest.NewRequest("POST", "/api/v1/stores/scrap-returns", bytes.NewBufferString(`{"requisition_id":"req-01"}`))
	reqValScrap.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wValScrap := httptest.NewRecorder()
	mux.ServeHTTP(wValScrap, reqValScrap)
	if wValScrap.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for scrap validation error, got: %d", wValScrap.Code)
	}
}
