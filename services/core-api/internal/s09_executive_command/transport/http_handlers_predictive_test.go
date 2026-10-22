// Package transport_test exercises predictive maintenance prognostic REST endpoints for S09.
package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/transport"
)

func TestHydroPrognosisEndpoint(t *testing.T) {
	mux, _, _, _ := setupTestServer(t)
	token := generateTestToken(t, "STAFF-GEN-01", "GM_GENERATION")

	// 1. Negative: Bad JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/hydro", bytes.NewBufferString("{bad"))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for bad json, got %d", rec.Code)
	}

	// 2. Negative: Missing AssetID
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/hydro", bytes.NewBufferString(`{"asset_id":""}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for empty asset_id, got %d", rec.Code)
	}

	// 3. Positive: Valid Hydro Prognosis
	dto := transport.HydroPrognosisDTO{
		AssetID:           "RES-KURA-FALLS-01",
		InflowM3s:         14.5,
		DischargeM3s:      12.0,
		CurrentHeadM:      22.0,
		TurbineEfficiency: 0.88,
		ForecastDays:      14,
	}
	raw, _ := json.Marshal(dto)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/hydro", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for hydro prognosis, got %d: %s", rec.Code, rec.Body.String())
	}

	var trigger domain.PredictiveMaintenanceTrigger
	if err := json.Unmarshal(rec.Body.Bytes(), &trigger); err != nil {
		t.Fatalf("failed to decode trigger: %v", err)
	}
	if trigger.AssetType != domain.AssetHydroReservoir || trigger.ModelName != domain.ModelHydroWaterAvailability {
		t.Fatalf("unexpected hydro trigger data: %+v", trigger)
	}
}

func TestGeneratorPrognosisEndpoint(t *testing.T) {
	mux, _, _, maintRepo := setupTestServer(t)
	token := generateTestToken(t, "STAFF-GEN-01", "GM_GENERATION")

	// 1. Negative: Missing AssetID
	req := httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/generator", bytes.NewBufferString(`{"asset_id":""}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for missing asset_id, got %d", rec.Code)
	}

	// 2. Negative: Bad JSON
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/generator", bytes.NewBufferString("{bad"))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for bad json, got %d", rec.Code)
	}

	// 3. Positive Critical: High Vibration and Temperature -> Triggers auto work order
	dto := transport.GeneratorPrognosisDTO{
		AssetID:               "GEN-UNIT-04-TURBINE",
		BearingVibrationMmS:   8.2, // severe vibration > 7.1 mm/s
		WindingTempC:          98.0,
		LubeOilPressureBar:    1.8,
		RunHoursSinceOverhaul: 9200.0,
		UnitHeadID:            "HEAD-GEN-01",
		TeamLeadID:            "LEAD-MECH-01",
	}
	raw, _ := json.Marshal(dto)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/generator", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for generator prognosis, got %d: %s", rec.Code, rec.Body.String())
	}

	var trigger domain.PredictiveMaintenanceTrigger
	if err := json.Unmarshal(rec.Body.Bytes(), &trigger); err != nil {
		t.Fatalf("failed to decode trigger: %v", err)
	}
	if !trigger.AutoWorkOrderCreated || trigger.LinkedWorkOrderID == nil {
		t.Fatalf("expected auto work order created for critical generator risk")
	}

	// Verify work order created in maintRepo
	wos, _ := maintRepo.ListWorkOrders(req.Context(), "", "", "")
	if len(wos) != 1 {
		t.Fatalf("expected 1 work order created in maintRepo, got %d", len(wos))
	}
}

func TestTransformerPrognosisEndpoint(t *testing.T) {
	mux, _, _, maintRepo := setupTestServer(t)
	token := generateTestToken(t, "STAFF-ENG-01", "HEAD_SERVICES")

	// 1. Negative: Missing AssetID
	req := httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/transformer", bytes.NewBufferString(`{"asset_id":""}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing asset_id, got %d", rec.Code)
	}

	// 2. Negative: Bad JSON
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/transformer", bytes.NewBufferString("{bad"))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad json, got %d", rec.Code)
	}

	// 3. Positive Critical: Overloaded transformer with high top oil temp
	dto := transport.TransformerPrognosisDTO{
		AssetID:         "TX-DT-500KVA-MARKET-SQ",
		LoadPct:         125.0,
		TopOilTempC:     92.0,
		FuseBlows30Days: 4,
		UnitHeadID:      "HEAD-TX-01",
		TeamLeadID:      "LEAD-CREW-01",
	}
	raw, _ := json.Marshal(dto)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/transformer", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for transformer prognosis, got %d: %s", rec.Code, rec.Body.String())
	}

	var trigger domain.PredictiveMaintenanceTrigger
	if err := json.Unmarshal(rec.Body.Bytes(), &trigger); err != nil {
		t.Fatalf("failed to decode trigger: %v", err)
	}
	if !trigger.AutoWorkOrderCreated || trigger.LinkedWorkOrderID == nil {
		t.Fatalf("expected auto work order created for critical transformer CBM")
	}

	wos, _ := maintRepo.ListWorkOrders(req.Context(), "", "", "")
	if len(wos) != 1 {
		t.Fatalf("expected 1 work order in maintRepo, got %d", len(wos))
	}
}

func TestFeederPrognosisEndpoint(t *testing.T) {
	mux, _, _, maintRepo := setupTestServer(t)
	token := generateTestToken(t, "STAFF-ENG-01", "HEAD_ENERGY_MGMT")

	// 1. Negative: Missing AssetID
	req := httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/feeder", bytes.NewBufferString(`{"asset_id":""}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing asset_id, got %d", rec.Code)
	}

	// 2. Negative: Bad JSON
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/feeder", bytes.NewBufferString("{bad"))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad json, got %d", rec.Code)
	}

	// 3. Positive Critical: Heavy trips and storm forecast -> Vegetation work order
	dto := transport.FeederPrognosisDTO{
		AssetID:             "FDR-33KV-VOM-EXPRESS",
		TripCount30Days:     12,
		FaultDensityPerKm:   1.8,
		DaysSinceLastPatrol: 45.0,
		StormForecast:       true,
		UnitHeadID:          "HEAD-FDR-01",
		TeamLeadID:          "LEAD-PATROL-01",
	}
	raw, _ := json.Marshal(dto)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/feeder", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for feeder prognosis, got %d: %s", rec.Code, rec.Body.String())
	}

	var trigger domain.PredictiveMaintenanceTrigger
	if err := json.Unmarshal(rec.Body.Bytes(), &trigger); err != nil {
		t.Fatalf("failed to decode trigger: %v", err)
	}
	if !trigger.AutoWorkOrderCreated || trigger.LinkedWorkOrderID == nil {
		t.Fatalf("expected auto work order created for critical feeder risk")
	}

	wos, _ := maintRepo.ListWorkOrders(req.Context(), "", "", "")
	if len(wos) != 1 {
		t.Fatalf("expected 1 work order in maintRepo, got %d", len(wos))
	}
}

func TestPredictiveTriggersAckAndListEndpoints(t *testing.T) {
	mux, _, _, _ := setupTestServer(t)
	token := generateTestToken(t, "STAFF-EXEC-01", "MD")

	// Create a hydro prognosis to get a trigger
	hydroDTO := transport.HydroPrognosisDTO{
		AssetID:           "RES-KURA-01",
		InflowM3s:         10.0,
		DischargeM3s:      10.0,
		CurrentHeadM:      20.0,
		TurbineEfficiency: 0.85,
		ForecastDays:      7,
	}
	raw, _ := json.Marshal(hydroDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/hydro", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("failed to create initial hydro trigger: %s", rec.Body.String())
	}
	var created domain.PredictiveMaintenanceTrigger
	_ = json.Unmarshal(rec.Body.Bytes(), &created)

	// 1. List predictive triggers -> 200 OK
	req = httptest.NewRequest(http.MethodGet, "/api/v1/executive/predictive/triggers?status=ACTIVE&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for list triggers, got %d", rec.Code)
	}
	var list []domain.PredictiveMaintenanceTrigger
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil || len(list) != 1 {
		t.Fatalf("expected 1 trigger in list, got %d", len(list))
	}

	// 2. Negative Ack: Bad JSON
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/triggers/"+created.ID+"/ack", bytes.NewBufferString("{bad"))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad json ack, got %d", rec.Code)
	}

	// 3. Negative Ack: Missing staff_id
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/triggers/"+created.ID+"/ack", bytes.NewBufferString(`{"staff_id":""}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing staff_id, got %d", rec.Code)
	}

	// 4. Negative Ack: Non-existent trigger
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/triggers/non-existent-trigger/ack", bytes.NewBufferString(`{"staff_id":"STAFF-01"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for non-existent trigger, got %d", rec.Code)
	}

	// 5. Positive Ack: Staff acknowledges trigger
	ackDTO := transport.AcknowledgeTriggerDTO{
		StaffID: "STAFF-EXEC-01",
		Status:  "ACKNOWLEDGED",
	}
	raw, _ = json.Marshal(ackDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/predictive/triggers/"+created.ID+"/ack", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for ack trigger, got %d: %s", rec.Code, rec.Body.String())
	}
	var acked domain.PredictiveMaintenanceTrigger
	if err := json.Unmarshal(rec.Body.Bytes(), &acked); err != nil {
		t.Fatalf("failed to decode acked trigger: %v", err)
	}
	if acked.AcknowledgedByStaffID == nil || *acked.AcknowledgedByStaffID != "STAFF-EXEC-01" {
		t.Fatalf("expected staff_id set on acked trigger, got %+v", acked)
	}
	if acked.Status != "ACKNOWLEDGED" {
		t.Fatalf("expected status ACKNOWLEDGED, got %s", acked.Status)
	}
}
