// Package transport_test exercises REST endpoints for intake telemetry,
// priority feeder load shedding, and Google Forms B webhook ingestion.
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
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/service"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/transport"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const (
	testJWTKey  = "nesco-energy-jwt-key-32b-secret-test!"
	testHMACKey = "nesco-energy-hmac-key-32b-secret-test"
)

func generateTestToken(t *testing.T, staffID, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(auth.Claims{
		StaffID:     staffID,
		StaffNumber: "NES-" + staffID,
		FullName:    "Engineer " + staffID,
		PrimaryRole: role,
		OrgUnitID:   "unit-energy-dispatch",
		IssuedAt:    time.Now().Unix(),
		ExpiresAt:   time.Now().Add(15 * time.Minute).Unix(),
	}, testJWTKey)
	if err != nil {
		t.Fatalf("failed to generate auth token: %v", err)
	}
	return token
}

func setupTestServer(t *testing.T) (*http.ServeMux, *repository.MemoryRepository, *s01repo.MemoryRepository) {
	t.Helper()
	repo := repository.NewMemoryRepository()
	gridRepo := s01repo.NewMemoryRepository()
	bus := events.NewMemoryEventBus()

	svc := service.NewWorkflowService(repo, gridRepo, bus, testHMACKey)
	handler := transport.NewHTTPHandler(svc, repo, testJWTKey)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, repo, gridRepo
}

func TestHTTPHandler_BulkIntake(t *testing.T) {
	mux, _, _ := setupTestServer(t)
	token := generateTestToken(t, "op-101", "DISPATCH_OPERATOR")

	// 1. Unauthenticated request -> 401
	unauthReq := httptest.NewRequest("POST", "/api/v1/energy/bulk-intake", bytes.NewBufferString("{}"))
	wUnauth := httptest.NewRecorder()
	mux.ServeHTTP(wUnauth, unauthReq)
	if wUnauth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got: %d", wUnauth.Code)
	}

	// 2. Initial Get Latest -> 404 Not Found
	getReq := httptest.NewRequest("GET", "/api/v1/energy/bulk-intake/latest", nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, getReq)
	if wGet.Code != http.StatusNotFound {
		t.Errorf("expected 404 for empty latest intake, got: %d", wGet.Code)
	}

	// 3. Valid Bulk Intake Post -> 201 Created
	intakePayload := transport.RecordBulkIntakeDTO{
		ShiftID:            "SHIFT-20261022-M",
		ShiftWindow:        domain.ShiftMorning,
		GenerationMWIn:     15.0,
		TCNMWIn:            3.5,
		SubstationRollcall: "[{\"operator\":\"M. Danladi\"}]",
		HandoverNotes:      "All grid sectors stable",
	}
	body, _ := json.Marshal(intakePayload)

	postReq := httptest.NewRequest("POST", "/api/v1/energy/bulk-intake", bytes.NewReader(body))
	postReq.Header.Set("Authorization", "Bearer "+token)
	wPost := httptest.NewRecorder()
	mux.ServeHTTP(wPost, postReq)

	if wPost.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got: %d, body: %s", wPost.Code, wPost.Body.String())
	}

	var created domain.DispatchTelemetryLog
	_ = json.Unmarshal(wPost.Body.Bytes(), &created)
	if created.TotalBulkMW != 18.5 {
		t.Fatalf("expected TotalBulkMW 18.5, got: %f", created.TotalBulkMW)
	}

	// 4. Verify Latest Bulk Intake -> 200 OK
	wGetLatest := httptest.NewRecorder()
	getLatestReq := httptest.NewRequest("GET", "/api/v1/energy/bulk-intake/latest", nil)
	getLatestReq.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wGetLatest, getLatestReq)

	if wGetLatest.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got: %d", wGetLatest.Code)
	}

	// 5. List Bulk Intake with limit -> 200 OK
	wList := httptest.NewRecorder()
	listReq := httptest.NewRequest("GET", "/api/v1/energy/bulk-intake?limit=5", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wList, listReq)
	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got: %d", wList.Code)
	}

	// 6. Invalid JSON -> 400 Bad Request
	wBadJSON := httptest.NewRecorder()
	badJSONReq := httptest.NewRequest("POST", "/api/v1/energy/bulk-intake", bytes.NewBufferString("not-json"))
	badJSONReq.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wBadJSON, badJSONReq)
	if wBadJSON.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got: %d", wBadJSON.Code)
	}

	// 7. Validation Error (missing shift_id) -> 400 Bad Request
	wValErr := httptest.NewRecorder()
	valErrReq := httptest.NewRequest("POST", "/api/v1/energy/bulk-intake", bytes.NewBufferString(`{"shift_id":""}`))
	valErrReq.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wValErr, valErrReq)
	if wValErr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for validation error, got: %d", wValErr.Code)
	}

	// 8. Service unprocessable error (negative MW) -> 400 / 422
	wNegMW := httptest.NewRecorder()
	negMWReq := httptest.NewRequest("POST", "/api/v1/energy/bulk-intake", bytes.NewBufferString(`{"shift_id":"S1","shift_window":"MORNING","generation_mw_in":-5.0}`))
	negMWReq.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wNegMW, negMWReq)
	if wNegMW.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative generation, got: %d", wNegMW.Code)
	}
}

func TestHTTPHandler_RationingAndFeeders(t *testing.T) {
	mux, repo, _ := setupTestServer(t)
	token := generateTestToken(t, "op-101", "DISPATCH_OPERATOR")
	ctx := context.Background()

	// 1. Calculate rationing with no intake -> 412 Precondition Failed
	ratReq := httptest.NewRequest("POST", "/api/v1/energy/rationing/calculate", bytes.NewBufferString("{}"))
	ratReq.Header.Set("Authorization", "Bearer "+token)
	wRat := httptest.NewRecorder()
	mux.ServeHTTP(wRat, ratReq)
	if wRat.Code != http.StatusPreconditionFailed {
		t.Errorf("expected 412 Precondition Failed, got: %d", wRat.Code)
	}

	// 1b. List empty feeders -> 200 OK (empty array)
	wEmptyFeeders := httptest.NewRecorder()
	reqEmptyFeeders := httptest.NewRequest("GET", "/api/v1/energy/feeders", nil)
	reqEmptyFeeders.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wEmptyFeeders, reqEmptyFeeders)
	if wEmptyFeeders.Code != http.StatusOK {
		t.Errorf("expected 200 OK for empty feeders, got: %d", wEmptyFeeders.Code)
	}

	// 2. Seed a feeder
	_ = repo.UpsertFeederStatus(ctx, domain.FeederStatus{
		FeederAssetID: "feeder-bukuru-1",
		FeederCode:    "FDR-BUK-01",
		PriorityRank:  1,
		TariffBand:    s01domain.TariffBandA,
		CurrentLoadMW: 4.5,
		State:         domain.FeederInService,
		UpdatedAt:     time.Now().UTC(),
	})

	// 3. Calculate rationing with negative MW -> 400 Bad Request
	wRatNeg := httptest.NewRecorder()
	ratNegReq := httptest.NewRequest("POST", "/api/v1/energy/rationing/calculate", bytes.NewBufferString(`{"available_bulk_mw": -2.0}`))
	ratNegReq.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wRatNeg, ratNegReq)
	if wRatNeg.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for negative rationing MW, got: %d", wRatNeg.Code)
	}

	// 4. Calculate rationing with explicit AvailableBulkMW = 10.0 -> 200 OK
	calcPayload := transport.RationingCalculateDTO{AvailableBulkMW: 10.0}
	calcBody, _ := json.Marshal(calcPayload)
	ratReq2 := httptest.NewRequest("POST", "/api/v1/energy/rationing/calculate", bytes.NewReader(calcBody))
	ratReq2.Header.Set("Authorization", "Bearer "+token)
	wRat2 := httptest.NewRecorder()
	mux.ServeHTTP(wRat2, ratReq2)

	if wRat2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got: %d, body: %s", wRat2.Code, wRat2.Body.String())
	}

	// 4. List Feeders -> 200 OK
	wFeeders := httptest.NewRecorder()
	reqFeeders := httptest.NewRequest("GET", "/api/v1/energy/feeders", nil)
	reqFeeders.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wFeeders, reqFeeders)
	if wFeeders.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got: %d", wFeeders.Code)
	}

	// 5. Get Feeder by Asset ID -> 200 OK
	wFeederSingle := httptest.NewRecorder()
	reqSingle := httptest.NewRequest("GET", "/api/v1/energy/feeders/feeder-bukuru-1", nil)
	reqSingle.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wFeederSingle, reqSingle)
	if wFeederSingle.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got: %d", wFeederSingle.Code)
	}

	// 6. Get Non-existent Feeder -> 404 Not Found
	wNotFound := httptest.NewRecorder()
	reqNotFound := httptest.NewRequest("GET", "/api/v1/energy/feeders/non-existent-feeder", nil)
	reqNotFound.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(wNotFound, reqNotFound)
	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got: %d", wNotFound.Code)
	}
}

func TestHTTPHandler_GoogleFormsWebhook(t *testing.T) {
	mux, _, _ := setupTestServer(t)

	payloadJSON := `{
		"shift_id": "SHIFT-FORM-B",
		"shift_window": "NIGHT",
		"operator_id": "op-danladi",
		"generation_mw_in": 16.0,
		"tcn_mw_in": 0.0,
		"substation_rollcall": "[{\"sub\":\"Bukuru\",\"op\":\"Danladi\"}]",
		"handover_notes": "All lines active",
		"feeders": []
	}`

	// 1. Invalid HMAC signature -> 401 Unauthorized
	reqInvalid := httptest.NewRequest("POST", "/api/v1/energy/webhook/google-forms", bytes.NewBufferString(payloadJSON))
	reqInvalid.Header.Set("X-Signature-SHA256", "invalid-sig-hex")
	wInvalid := httptest.NewRecorder()
	mux.ServeHTTP(wInvalid, reqInvalid)
	if wInvalid.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for bad HMAC, got: %d", wInvalid.Code)
	}

	// 2. Valid HMAC signature -> 200 OK
	validSig, _ := crypto.ComputeRowHMAC(testHMACKey, payloadJSON)
	reqValid := httptest.NewRequest("POST", "/api/v1/energy/webhook/google-forms", bytes.NewBufferString(payloadJSON))
	reqValid.Header.Set("X-Signature-SHA256", validSig)
	wValid := httptest.NewRecorder()
	mux.ServeHTTP(wValid, reqValid)

	if wValid.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for valid webhook, got: %d, body: %s", wValid.Code, wValid.Body.String())
	}

	// 3. Webhook with invalid JSON payload -> 422 Unprocessable Entity
	badPayload := "not-valid-json-string"
	badSig, _ := crypto.ComputeRowHMAC(testHMACKey, badPayload)
	reqBad := httptest.NewRequest("POST", "/api/v1/energy/webhook/google-forms", bytes.NewBufferString(badPayload))
	reqBad.Header.Set("X-Signature-SHA256", badSig)
	wBad := httptest.NewRecorder()
	mux.ServeHTTP(wBad, reqBad)
	if wBad.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for malformed JSON webhook, got: %d", wBad.Code)
	}
}
