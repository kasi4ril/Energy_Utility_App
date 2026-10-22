// Package transport_test exercises the REST endpoints for the 5-step Safety Permit-to-Work (PTW) lifecycle.
package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/transport"
)

func TestHTTPHandler_PTW_Lifecycle(t *testing.T) {
	mux, repo, gridRepo := setupTestServer(t)
	ctx := context.Background()

	tokenUnitHead := generateTestToken(t, "unit-head-01", "UNIT_HEAD")
	tokenDeptHead := generateTestToken(t, "dept-head-01", "HEAD_SERVICES")
	tokenSafetyRep := generateTestToken(t, "safety-rep-01", "SAFETY_REP")
	tokenOperator := generateTestToken(t, "operator-01", "SUBSTATION_OPERATOR")

	// Seed Grid Asset & Feeder
	gridAsset := s01domain.GridAsset{
		ID:               "feeder-33kv-bukuru",
		AssetCode:        "FDR-33KV-BUKURU",
		AssetType:        s01domain.AssetFeeder11KV,
		OrgUnitID:        "unit-energy-dispatch",
		OperationalState: s01domain.StateInService,
	}
	_ = gridRepo.CreateGridAsset(ctx, gridAsset)

	_ = repo.UpsertFeederStatus(ctx, domain.FeederStatus{
		FeederAssetID: gridAsset.ID,
		FeederCode:    gridAsset.AssetCode,
		PriorityRank:  1,
		TariffBand:    s01domain.TariffBandA,
		CurrentLoadMW: 5.5,
		State:         domain.FeederInService,
		UpdatedAt:     time.Now().UTC(),
	})

	// 1. Step 1: Request PTW
	// Negative: bad JSON
	reqBadJSON := httptest.NewRequest("POST", "/api/v1/energy/ptw/request", bytes.NewBufferString("not-json"))
	reqBadJSON.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wBadJSON := httptest.NewRecorder()
	mux.ServeHTTP(wBadJSON, reqBadJSON)
	if wBadJSON.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON on PTW request, got: %d", wBadJSON.Code)
	}

	// Negative: validation error (missing fields)
	reqValErr := httptest.NewRequest("POST", "/api/v1/energy/ptw/request", bytes.NewBufferString(`{"asset_id":""}`))
	reqValErr.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wValErr := httptest.NewRecorder()
	mux.ServeHTTP(wValErr, reqValErr)
	if wValErr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for validation error, got: %d", wValErr.Code)
	}

	// Negative: non-existent grid asset -> 422
	reqNonEx := httptest.NewRequest("POST", "/api/v1/energy/ptw/request", bytes.NewBufferString(`{"asset_id":"non-existent","work_description":"repair"}`))
	reqNonEx.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wNonEx := httptest.NewRecorder()
	mux.ServeHTTP(wNonEx, reqNonEx)
	if wNonEx.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for non-existent grid asset, got: %d", wNonEx.Code)
	}

	initDTO := transport.InitiatePTWDTO{
		AssetID:         gridAsset.ID,
		WorkDescription: "Urgent busbar replacement at Bukuru injection substation",
	}
	bodyInit, _ := json.Marshal(initDTO)

	reqInit := httptest.NewRequest("POST", "/api/v1/energy/ptw/request", bytes.NewReader(bodyInit))
	reqInit.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wInit := httptest.NewRecorder()
	mux.ServeHTTP(wInit, reqInit)

	if wInit.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for PTW request, got: %d, body: %s", wInit.Code, wInit.Body.String())
	}

	var ptw domain.SafetyPermitPTW
	_ = json.Unmarshal(wInit.Body.Bytes(), &ptw)
	if ptw.Status != domain.PTWRequested {
		t.Fatalf("expected status REQUESTED, got: %s", ptw.Status)
	}

	// 2. Step 2: Dual Approval
	// Negative: unauthorized role
	unauthDTO := transport.ApprovePTWDTO{ApproverRole: "ACCOUNTANT"}
	bodyUnauth, _ := json.Marshal(unauthDTO)
	reqUnauth := httptest.NewRequest("POST", "/api/v1/energy/ptw/"+ptw.ID+"/approve", bytes.NewReader(bodyUnauth))
	reqUnauth.Header.Set("Authorization", "Bearer "+tokenUnitHead)
	wUnauth := httptest.NewRecorder()
	mux.ServeHTTP(wUnauth, reqUnauth)
	if wUnauth.Code != http.StatusBadRequest { // validation error on DTO
		t.Errorf("expected 400 Bad Request on invalid approver role, got: %d", wUnauth.Code)
	}

	// Valid Sign-off 1: Head of Services
	approveDTO1 := transport.ApprovePTWDTO{ApproverRole: "HEAD_SERVICES"}
	bodyApprove1, _ := json.Marshal(approveDTO1)
	reqApprove1 := httptest.NewRequest("POST", "/api/v1/energy/ptw/"+ptw.ID+"/approve", bytes.NewReader(bodyApprove1))
	reqApprove1.Header.Set("Authorization", "Bearer "+tokenDeptHead)
	wApprove1 := httptest.NewRecorder()
	mux.ServeHTTP(wApprove1, reqApprove1)

	if wApprove1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Dept Head signoff, got: %d", wApprove1.Code)
	}

	// Valid Sign-off 2: HSE Safety Rep
	approveDTO2 := transport.ApprovePTWDTO{ApproverRole: "SAFETY_REP"}
	bodyApprove2, _ := json.Marshal(approveDTO2)
	reqApprove2 := httptest.NewRequest("POST", "/api/v1/energy/ptw/"+ptw.ID+"/approve", bytes.NewReader(bodyApprove2))
	reqApprove2.Header.Set("Authorization", "Bearer "+tokenSafetyRep)
	wApprove2 := httptest.NewRecorder()
	mux.ServeHTTP(wApprove2, reqApprove2)

	if wApprove2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Safety Rep signoff, got: %d", wApprove2.Code)
	}

	var ptwApproved domain.SafetyPermitPTW
	_ = json.Unmarshal(wApprove2.Body.Bytes(), &ptwApproved)
	if ptwApproved.Status != domain.PTWApprovedPendingIsolation {
		t.Fatalf("expected status APPROVED_PENDING_ISOLATION, got: %s", ptwApproved.Status)
	}

	// 3. Step 3: Substation Isolation
	isoDTO := transport.IsolatePTWDTO{IsolationNotes: "Main breaker isolated, locked with padlock #42, earthed."}
	bodyIso, _ := json.Marshal(isoDTO)
	reqIso := httptest.NewRequest("POST", "/api/v1/energy/ptw/"+ptw.ID+"/isolate", bytes.NewReader(bodyIso))
	reqIso.Header.Set("Authorization", "Bearer "+tokenOperator)
	wIso := httptest.NewRecorder()
	mux.ServeHTTP(wIso, reqIso)

	if wIso.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for isolation, got: %d, body: %s", wIso.Code, wIso.Body.String())
	}

	var ptwIso domain.SafetyPermitPTW
	_ = json.Unmarshal(wIso.Body.Bytes(), &ptwIso)
	if ptwIso.Status != domain.PTWIsolatedEarthed {
		t.Fatalf("expected ISOLATED_EARTHED, got: %s", ptwIso.Status)
	}

	// Verify Feeder state transitioned to OUT_ON_PERMIT_PTW
	fdrState, err := repo.GetFeederStatusByAssetID(ctx, gridAsset.ID)
	if err != nil || fdrState.State != domain.FeederOutOnPTW {
		t.Fatalf("feeder state must be OUT_ON_PERMIT_PTW, got: %+v", fdrState)
	}

	// Negative: premature restore before safety cleared -> 422
	reqPremature := httptest.NewRequest("POST", "/api/v1/energy/ptw/"+ptw.ID+"/restore", bytes.NewBufferString("{}"))
	reqPremature.Header.Set("Authorization", "Bearer "+tokenOperator)
	wPremature := httptest.NewRecorder()
	mux.ServeHTTP(wPremature, reqPremature)
	if wPremature.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for premature restore attempt, got: %d", wPremature.Code)
	}

	// 4. Step 5a: HSE Safety Clearance
	clearDTO := transport.ClearSafetyPTWDTO{Notes: "Work complete. All workers off tower and grounds removed."}
	bodyClear, _ := json.Marshal(clearDTO)
	reqClear := httptest.NewRequest("POST", "/api/v1/energy/ptw/"+ptw.ID+"/clear-safety", bytes.NewReader(bodyClear))
	reqClear.Header.Set("Authorization", "Bearer "+tokenSafetyRep)
	wClear := httptest.NewRecorder()
	mux.ServeHTTP(wClear, reqClear)

	if wClear.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for safety clearance, got: %d, body: %s", wClear.Code, wClear.Body.String())
	}

	// 5. Step 5b: Restoration & Closure
	reqRestore := httptest.NewRequest("POST", "/api/v1/energy/ptw/"+ptw.ID+"/restore", bytes.NewBufferString("{}"))
	reqRestore.Header.Set("Authorization", "Bearer "+tokenOperator)
	wRestore := httptest.NewRecorder()
	mux.ServeHTTP(wRestore, reqRestore)

	if wRestore.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for restore, got: %d, body: %s", wRestore.Code, wRestore.Body.String())
	}

	var ptwClosed domain.SafetyPermitPTW
	_ = json.Unmarshal(wRestore.Body.Bytes(), &ptwClosed)
	if ptwClosed.Status != domain.PTWRestoredClosed {
		t.Fatalf("expected RESTORED_CLOSED, got: %s", ptwClosed.Status)
	}

	// Verify Feeder state reverted to IN_SERVICE
	fdrReverted, err := repo.GetFeederStatusByAssetID(ctx, gridAsset.ID)
	if err != nil || fdrReverted.State != domain.FeederInService {
		t.Fatalf("feeder state must revert to IN_SERVICE, got: %+v", fdrReverted)
	}

	// 6. Query by ID and Query by PermitNo -> 200 OK
	wGetID := httptest.NewRecorder()
	reqGetID := httptest.NewRequest("GET", "/api/v1/energy/ptw/"+ptw.ID, nil)
	reqGetID.Header.Set("Authorization", "Bearer "+tokenOperator)
	mux.ServeHTTP(wGetID, reqGetID)
	if wGetID.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GetPTW by ID, got: %d", wGetID.Code)
	}

	wGetNo := httptest.NewRecorder()
	reqGetNo := httptest.NewRequest("GET", "/api/v1/energy/ptw/"+ptw.PermitNo, nil)
	reqGetNo.Header.Set("Authorization", "Bearer "+tokenOperator)
	mux.ServeHTTP(wGetNo, reqGetNo)
	if wGetNo.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GetPTW by PermitNo, got: %d", wGetNo.Code)
	}

	// Negative: query non-existent permit -> 404
	wGetNotFound := httptest.NewRecorder()
	reqGetNotFound := httptest.NewRequest("GET", "/api/v1/energy/ptw/non-existent-ptw", nil)
	reqGetNotFound.Header.Set("Authorization", "Bearer "+tokenOperator)
	mux.ServeHTTP(wGetNotFound, reqGetNotFound)
	if wGetNotFound.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent permit, got: %d", wGetNotFound.Code)
	}

	// 7. List PTWs -> 200 OK
	wList := httptest.NewRecorder()
	reqList := httptest.NewRequest("GET", "/api/v1/energy/ptw?status=RESTORED_CLOSED", nil)
	reqList.Header.Set("Authorization", "Bearer "+tokenOperator)
	mux.ServeHTTP(wList, reqList)
	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for ListPTWs, got: %d", wList.Code)
	}
}

func TestHTTPHandler_PTW_NegativePaths(t *testing.T) {
	mux, _, _ := setupTestServer(t)
	token := generateTestToken(t, "user-test", "HEAD_SERVICES")

	// 1. Approve with bad JSON -> 400
	reqAppBadJSON := httptest.NewRequest("POST", "/api/v1/energy/ptw/ptw-1/approve", bytes.NewBufferString("not-json"))
	reqAppBadJSON.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, reqAppBadJSON)
	if w1.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON on approve, got: %d", w1.Code)
	}

	// 2. Isolate with bad JSON -> 400
	reqIsoBadJSON := httptest.NewRequest("POST", "/api/v1/energy/ptw/ptw-1/isolate", bytes.NewBufferString("not-json"))
	reqIsoBadJSON.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, reqIsoBadJSON)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON on isolate, got: %d", w2.Code)
	}

	// 3. Isolate with short notes -> 400
	reqIsoShort := httptest.NewRequest("POST", "/api/v1/energy/ptw/ptw-1/isolate", bytes.NewBufferString(`{"isolation_notes":"short"}`))
	reqIsoShort.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, reqIsoShort)
	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for short isolation notes, got: %d", w3.Code)
	}

	// 4. Clear safety with bad JSON -> 400
	reqClearBadJSON := httptest.NewRequest("POST", "/api/v1/energy/ptw/ptw-1/clear-safety", bytes.NewBufferString("not-json"))
	reqClearBadJSON.Header.Set("Authorization", "Bearer "+token)
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, reqClearBadJSON)
	if w4.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON on clear safety, got: %d", w4.Code)
	}

	// 5. Clear safety with short notes -> 400
	reqClearShort := httptest.NewRequest("POST", "/api/v1/energy/ptw/ptw-1/clear-safety", bytes.NewBufferString(`{"notes":"short"}`))
	reqClearShort.Header.Set("Authorization", "Bearer "+token)
	w5 := httptest.NewRecorder()
	mux.ServeHTTP(w5, reqClearShort)
	if w5.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for short clearance notes, got: %d", w5.Code)
	}

	// 6. Approve on non-existent permit -> 422
	reqAppNotFound := httptest.NewRequest("POST", "/api/v1/energy/ptw/ptw-nonexistent/approve", bytes.NewBufferString(`{"approver_role":"HEAD_SERVICES"}`))
	reqAppNotFound.Header.Set("Authorization", "Bearer "+token)
	w6 := httptest.NewRecorder()
	mux.ServeHTTP(w6, reqAppNotFound)
	if w6.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for non-existent permit approve, got: %d", w6.Code)
	}
}

