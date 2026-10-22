package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
)

func TestDTO_RecordTelemetryValidation(t *testing.T) {
	// Happy path
	valid := RecordTelemetryDTO{
		ShiftID:     "s1",
		PlantID:     "p1",
		PlantState:  domain.StateInService,
		WaterLevelM: 80.0,
		MWGenerated: 10.0,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	// 1. Missing shift_id or plant_id
	dto := RecordTelemetryDTO{PlantID: "p1", PlantState: domain.StateInService}
	if err := dto.Validate(); err == nil {
		t.Error("expected error on empty shift_id")
	}

	// 2. Invalid plant_state
	dto = RecordTelemetryDTO{ShiftID: "s1", PlantID: "p1", PlantState: "UNKNOWN"}
	if err := dto.Validate(); err == nil {
		t.Error("expected error on invalid plant_state")
	}

	// 3. Water level out of range
	dto = RecordTelemetryDTO{ShiftID: "s1", PlantID: "p1", PlantState: domain.StateInService, WaterLevelM: -10}
	if err := dto.Validate(); err == nil {
		t.Error("expected error on negative water level")
	}
	dto.WaterLevelM = 300.0
	if err := dto.Validate(); err == nil {
		t.Error("expected error on excessive water level")
	}

	// 4. Negative MW
	dto.WaterLevelM = 50.0
	dto.MWGenerated = -2.0
	if err := dto.Validate(); err == nil {
		t.Error("expected error on negative MWGenerated")
	}
}

func TestDTO_RecordRollcallValidation(t *testing.T) {
	// Happy path
	valid := RecordRollcallDTO{
		ShiftID:       "s1",
		PlantID:       "p1",
		OperatorNames: []string{"Musa"},
		HandoverNotes: "Valid 10+ characters handover notes",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid rollcall DTO, got: %v", err)
	}

	// Missing shift or plant
	dto := RecordRollcallDTO{PlantID: "p1"}
	if err := dto.Validate(); err == nil {
		t.Error("expected error on empty shift_id")
	}

	// Missing operators
	dto = RecordRollcallDTO{ShiftID: "s1", PlantID: "p1", OperatorNames: []string{""}}
	if err := dto.Validate(); err == nil {
		t.Error("expected error on empty operators")
	}

	// Short handover notes
	dto = RecordRollcallDTO{ShiftID: "s1", PlantID: "p1", OperatorNames: []string{"John"}, HandoverNotes: "short"}
	if err := dto.Validate(); err == nil {
		t.Error("expected error on short handover notes")
	}
}

func TestDTO_FaultJobAndGMValidation(t *testing.T) {
	// Happy paths
	fValid := ReportFaultJobDTO{
		PlantID:          "p1",
		Wing:             domain.WingElectrical,
		FaultDescription: "Valid 15+ characters fault description",
	}
	if err := fValid.Validate(); err != nil {
		t.Fatalf("expected valid fault DTO, got: %v", err)
	}

	gmValid := GMApprovalDTO{Approved: true}
	if err := gmValid.Validate(); err != nil {
		t.Fatalf("expected valid GM approval DTO, got: %v", err)
	}

	wValid := UpdateWorkshopDTO{WorkshopStatus: domain.WorkshopMachiningInProgress}
	if err := wValid.Validate(); err != nil {
		t.Fatalf("expected valid workshop DTO, got: %v", err)
	}

	// Missing plant_id
	fDTO := ReportFaultJobDTO{Wing: domain.WingElectrical, FaultDescription: "Valid fault description"}
	if err := fDTO.Validate(); err == nil {
		t.Error("expected error on empty plant_id")
	}

	// Invalid wing
	fDTO = ReportFaultJobDTO{PlantID: "p1", Wing: "CIVIL", FaultDescription: "Valid fault description"}
	if err := fDTO.Validate(); err == nil {
		t.Error("expected error on invalid wing")
	}

	// Short description
	fDTO = ReportFaultJobDTO{PlantID: "p1", Wing: domain.WingMechanical, FaultDescription: "short"}
	if err := fDTO.Validate(); err == nil {
		t.Error("expected error on short description")
	}

	// GM rejection without notes
	gmDTO := GMApprovalDTO{Approved: false, Notes: "short"}
	if err := gmDTO.Validate(); err == nil {
		t.Error("expected error on short rejection notes")
	}

	// Empty workshop status
	wDTO := UpdateWorkshopDTO{WorkshopStatus: ""}
	if err := wDTO.Validate(); err == nil {
		t.Error("expected error on empty workshop status")
	}
}

func TestHTTP_HandlerBranchCoverage(t *testing.T) {
	handler, mux, _ := setupHTTPTestEnv(t)
	token := generateToken(t, "op-1", "SYSTEM_OPERATOR")

	// 1. Rollcall validation error via HTTP
	badRC := RecordRollcallDTO{ShiftID: "s1", PlantID: "p1", OperatorNames: []string{""}}
	b, _ := json.Marshal(badRC)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/generation/rollcall", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on invalid rollcall DTO, got %d", rec.Code)
	}

	// 2. Fault reporting validation error via HTTP
	badF := ReportFaultJobDTO{PlantID: "plant-1-uuid", Wing: "CIVIL", FaultDescription: "short"}
	b, _ = json.Marshal(badF)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on invalid fault DTO, got %d", rec.Code)
	}

	// 3. GM Approval bad JSON & validation failure
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults/job-1/gm-approval", bytes.NewReader([]byte("{bad")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on bad JSON in GM approval, got %d", rec.Code)
	}

	badGM := GMApprovalDTO{Approved: false, Notes: "short"}
	b, _ = json.Marshal(badGM)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults/job-1/gm-approval", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on short GM notes, got %d", rec.Code)
	}

	// 4. Update Workshop bad JSON & validation failure
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults/job-1/workshop", bytes.NewReader([]byte("{bad")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on bad JSON in workshop, got %d", rec.Code)
	}

	badW := UpdateWorkshopDTO{WorkshopStatus: ""}
	b, _ = json.Marshal(badW)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/generation/faults/job-1/workshop", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on empty workshop status, got %d", rec.Code)
	}

	// 5. Direct invocation of handlers without auth claims
	noAuthReq := httptest.NewRequest(http.MethodPost, "/any", nil)
	w := httptest.NewRecorder()
	handler.HandleRecordRollcall(w, noAuthReq)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without claims on rollcall, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	handler.HandleGMApproval(w, noAuthReq)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without claims on gm approval, got %d", w.Code)
	}
}

