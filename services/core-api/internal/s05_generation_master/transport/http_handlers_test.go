package transport

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
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const testJWTKey = "nesco-gen-jwt-key-32b-secret-test!"
const testHMACKey = "nesco-gen-hmac-key-32b-secret-test"

func setupHTTPTestEnv(t *testing.T) (*HTTPHandler, *http.ServeMux, *repository.MemoryRepository) {
	t.Helper()

	repo := repository.NewMemoryRepository()
	gridRepo := s01repo.NewMemoryRepository()
	bus := events.NewMemoryEventBus()

	// Seed Plant 1
	plant := s01domain.GridAsset{
		ID:        "plant-1-uuid",
		AssetCode: "GEN-01",
		AssetType: s01domain.AssetPlant,
		OrgUnitID: "org-gen",
		GPSPoint:  s01domain.GPSPoint{Latitude: 9.42, Longitude: 8.85},
	}
	_ = gridRepo.CreateGridAsset(context.Background(), plant)

	svc := service.NewWorkflowService(repo, gridRepo, bus, testHMACKey)
	handler := NewHTTPHandler(svc, repo, testJWTKey)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return handler, mux, repo
}

func generateToken(t *testing.T, staffID, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(auth.Claims{
		StaffID:     staffID,
		StaffNumber: "NES-" + staffID,
		FullName:    "Engineer " + staffID,
		PrimaryRole: role,
		OrgUnitID:   "unit-gen",
		IssuedAt:    time.Now().Unix(),
		ExpiresAt:   time.Now().Add(15 * time.Minute).Unix(),
	}, testJWTKey)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}

func TestHTTP_TelemetryAndRollcallFlow(t *testing.T) {
	_, mux, _ := setupHTTPTestEnv(t)
	token := generateToken(t, "master-op-1", "SYSTEM_OPERATOR")

	// 1. Record Telemetry
	tDTO := RecordTelemetryDTO{
		ShiftID:             "SHIFT-20261022-M",
		ShiftWindow:         domain.ShiftMorning,
		PlantID:             "plant-1-uuid",
		PlantState:          domain.StateInService,
		WaterLevelM:         90.5,
		MWGenerated:         18.0,
		StatorWindingTempC:  72.0,
		BearingVibrationMMS: 1.2,
		LubeOilPressureBar:  3.0,
	}
	body, _ := json.Marshal(tDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/generation/telemetry", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on telemetry, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. List Telemetry
	req = httptest.NewRequest(http.MethodGet, "/api/v1/generation/telemetry?plant_id=plant-1-uuid", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on list telemetry, got %d", rec.Code)
	}

	// 3. Get Latest Telemetry
	req = httptest.NewRequest(http.MethodGet, "/api/v1/generation/telemetry/latest/plant-1-uuid", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on get latest telemetry, got %d", rec.Code)
	}

	// 4. Get Latest for non-existent plant
	req = httptest.NewRequest(http.MethodGet, "/api/v1/generation/telemetry/latest/unknown-plant", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", rec.Code)
	}

	// 5. Record Rollcall
	rcDTO := RecordRollcallDTO{
		ShiftID:             "SHIFT-20261022-M",
		ShiftWindow:         domain.ShiftMorning,
		PlantID:             "plant-1-uuid",
		OperatorNames:       []string{"Operator Yakubu"},
		ProtectionEngineers: []string{"Engr. Pam"},
		HandoverNotes:       "All turbines online and synchronised.",
	}
	rcBody, _ := json.Marshal(rcDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/rollcall", bytes.NewReader(rcBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on rollcall, got %d: %s", rec.Code, rec.Body.String())
	}

	// 6. List Rollcalls
	req = httptest.NewRequest(http.MethodGet, "/api/v1/generation/rollcall?shift_id=SHIFT-20261022-M", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on list rollcalls, got %d", rec.Code)
	}
}

func TestHTTP_FaultAndWorkshopFlow(t *testing.T) {
	_, mux, _ := setupHTTPTestEnv(t)
	opToken := generateToken(t, "op-1", "SYSTEM_OPERATOR")
	gmToken := generateToken(t, "gm-1", "GENERAL_MANAGER_GEN")

	// 1. Report Plant Fault
	fDTO := ReportFaultJobDTO{
		PlantID:          "plant-1-uuid",
		Wing:             domain.WingMechanical,
		FaultDescription: "Turbine runner cavitation damage observed during inspection.",
		Severity:         "HIGH",
		RequiresWorkshop: true,
	}
	body, _ := json.Marshal(fDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+opToken)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on fault report, got %d: %s", rec.Code, rec.Body.String())
	}

	var job domain.PlantFaultAndWorkshopJob
	_ = json.NewDecoder(rec.Body).Decode(&job)

	// 2. List Faults
	req = httptest.NewRequest(http.MethodGet, "/api/v1/generation/faults?plant_id=plant-1-uuid", nil)
	req.Header.Set("Authorization", "Bearer "+opToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on list faults, got %d", rec.Code)
	}

	// 3. GM Approval
	gmDTO := GMApprovalDTO{
		Approved: true,
		Notes:    "Approved for immediate runner welding and dynamic balancing.",
	}
	gmBody, _ := json.Marshal(gmDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults/"+job.ID+"/gm-approval", bytes.NewReader(gmBody))
	req.Header.Set("Authorization", "Bearer "+gmToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on GM approval, got %d: %s", rec.Code, rec.Body.String())
	}

	// 4. Update Workshop Status
	matRef := "MAT-RUNNER-001"
	wDTO := UpdateWorkshopDTO{
		WorkshopStatus:      domain.WorkshopMachiningInProgress,
		StoreRequisitionRef: &matRef,
	}
	wBody, _ := json.Marshal(wDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults/"+job.ID+"/workshop", bytes.NewReader(wBody))
	req.Header.Set("Authorization", "Bearer "+opToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on workshop update, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHTTP_GoogleFormsWebhook(t *testing.T) {
	_, mux, _ := setupHTTPTestEnv(t)

	payload := map[string]interface{}{
		"shift_id":           "SHIFT-20261022-N",
		"shift_window":       "NIGHT",
		"master_operator_id": "master-op-99",
		"plants": []map[string]interface{}{
			{
				"plant_id":       "plant-1-uuid",
				"plant_state":    "IN_SERVICE",
				"water_level_m":  92.0,
				"mw_generated":   19.0,
				"stator_temp_c":  74.0,
				"vibration_mm_s": 1.0,
				"oil_pressure_bar": 3.2,
				"operator_names": []string{"Operator John"},
				"handover_notes": "Night shift active.",
			},
		},
	}
	rawBytes, _ := json.Marshal(payload)
	validHMAC, _ := crypto.ComputeRowHMAC(testHMACKey, string(rawBytes))

	// 1. Valid Webhook
	req := httptest.NewRequest(http.MethodPost, "/api/v1/generation/webhook/google-forms", bytes.NewReader(rawBytes))
	req.Header.Set("X-Signature-SHA256", validHMAC)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on valid webhook, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. Missing HMAC Header
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/webhook/google-forms", bytes.NewReader(rawBytes))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing signature header, got %d", rec.Code)
	}

	// 3. Invalid HMAC Signature
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/webhook/google-forms", bytes.NewReader(rawBytes))
	req.Header.Set("X-Signature-SHA256", "invalid-signature")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for bad signature, got %d", rec.Code)
	}
}

func TestHTTP_NegativeValidations(t *testing.T) {
	_, mux, _ := setupHTTPTestEnv(t)
	token := generateToken(t, "op-1", "SYSTEM_OPERATOR")

	// Invalid JSON bodies -> 400
	req := httptest.NewRequest(http.MethodPost, "/api/v1/generation/telemetry", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on invalid JSON in telemetry, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/rollcall", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on invalid JSON in rollcall, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on invalid JSON in faults, got %d", rec.Code)
	}

	// Validation failure on DTOs
	badDTO := RecordTelemetryDTO{WaterLevelM: 500.0} // > 250m
	body, _ := json.Marshal(badDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/telemetry", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on bad water level, got %d", rec.Code)
	}

	// Service error on unknown plant
	unknownPlantDTO := RecordTelemetryDTO{ShiftID: "s1", PlantID: "unknown-plant", PlantState: domain.StateInService, WaterLevelM: 80}
	upBody, _ := json.Marshal(unknownPlantDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/telemetry", bytes.NewReader(upBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 on unknown plant telemetry, got %d", rec.Code)
	}

	// Direct invocation without claims
	handler, _, _ := setupHTTPTestEnv(t)
	noAuthReq := httptest.NewRequest(http.MethodPost, "/any", nil)
	w := httptest.NewRecorder()
	handler.HandleRecordTelemetry(w, noAuthReq)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without claims, got %d", w.Code)
	}
	w = httptest.NewRecorder()
	handler.HandleReportFault(w, noAuthReq)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without claims, got %d", w.Code)
	}
}
