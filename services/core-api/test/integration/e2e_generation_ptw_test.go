// Package integration exercises end-to-end multi-sprint workflows across NESCO bounded contexts.
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	s06domain "disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

func TestE2E_GoogleFormsWebhooks_And_SafetyPTW_Lifecycle(t *testing.T) {
	mux, _, _, _, s05Repo, s06Repo := setupE2EEnvironment(t)
	token := auth.Claims{
		StaffID:     "STAFF-HQ-01",
		StaffNumber: "NES-HQ-01",
		FullName:    "Engineer HQ",
		PrimaryRole: "HEAD_SERVICES",
		OrgUnitID:   "unit-grid-ctrl",
	}
	jwtStr, _ := auth.GenerateToken(token, testJWTKey)
	authHeader := "Bearer " + jwtStr

	// Step 1: Google Form A Webhook Ingestion (S05 Generation Master)
	formAPayload := []byte(`{
		"shift_id": "SHIFT-2026-MORN",
		"shift_window": "MORNING",
		"master_operator_id": "STAFF-OP-01",
		"plants": [
			{
				"plant_id": "PLANT_01_KURRA_FALLS",
				"plant_state": "IN_SERVICE",
				"water_level_m": 24.2,
				"mw_generated": 14.8,
				"stator_temp_c": 68.5,
				"vibration_mm_s": 2.2,
				"oil_pressure_bar": 3.1,
				"operator_names": ["Operator Bitrus"],
				"protection_engineers": ["Engr Danladi"],
				"technicians": ["Tech Gyang"],
				"handover_notes": "All generators operating normally"
			}
		]
	}`)
	formASig := generateHMACHeader(formAPayload, testHMACKey)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/generation/webhook/google-forms", bytes.NewBuffer(formAPayload))
	req.Header.Set("X-Signature-SHA256", formASig)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Google Form A webhook ingestion, got %d: %s", rec.Code, rec.Body.String())
	}

	// Step 2: Google Form B Webhook Ingestion (S06 Energy Management & Dispatch)
	formBPayload := []byte(`{
		"shift_id": "SHIFT-DISP-001",
		"shift_window": "MORNING",
		"operator_id": "STAFF-DISP-01",
		"generation_mw_in": 45.0,
		"tcn_mw_in": 15.0,
		"substation_rollcall": "Substations Vom, Barkin Ladi operational",
		"handover_notes": "Grid stable at nominal frequency 50Hz"
	}`)
	formBSig := generateHMACHeader(formBPayload, testHMACKey)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/energy/webhook/google-forms", bytes.NewBuffer(formBPayload))
	req.Header.Set("X-Signature-SHA256", formBSig)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for Google Form B webhook ingestion, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify telemetry recorded in repositories
	genLogs, _ := s05Repo.ListTelemetryLogs(context.Background(), "", "", 10)
	if len(genLogs) != 1 {
		t.Fatalf("expected 1 generation log recorded from Form A webhook, got %d", len(genLogs))
	}
	latestDispatch, _ := s06Repo.GetLatestDispatchLog(context.Background())
	if latestDispatch == nil || latestDispatch.TCNMWIn != 15.0 {
		t.Fatalf("expected dispatch log recorded from Form B webhook, got %+v", latestDispatch)
	}

	// Step 3: Execute 5-Step Safety Permit-to-Work (PTW) Protocol (S06)
	// Sub-step A: Request PTW
	ptwReqPayload, _ := json.Marshal(map[string]interface{}{
		"asset_id":         "FDR-33KV-VOM",
		"work_description": "High-voltage insulator replacement at tower 14",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/energy/ptw/request", bytes.NewBuffer(ptwReqPayload))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for PTW request, got %d: %s", rec.Code, rec.Body.String())
	}
	var permit s06domain.SafetyPermitPTW
	_ = json.Unmarshal(rec.Body.Bytes(), &permit)

	// Sub-step B: Safety Dual-Approval (HEAD_SERVICES then SAFETY_REP)
	headApprovePayload, _ := json.Marshal(map[string]string{
		"approver_role": "HEAD_SERVICES",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/energy/ptw/"+permit.ID+"/approve", bytes.NewBuffer(headApprovePayload))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for head services approval, got %d: %s", rec.Code, rec.Body.String())
	}

	safetyRepToken := auth.Claims{
		StaffID:     "STAFF-SAFETY-01",
		StaffNumber: "NES-HSE-01",
		FullName:    "Safety Representative",
		PrimaryRole: "SAFETY_REP",
		OrgUnitID:   "unit-safety-hse",
	}
	safetyJwt, _ := auth.GenerateToken(safetyRepToken, testJWTKey)
	safetyAuth := "Bearer " + safetyJwt

	safetyApprovePayload, _ := json.Marshal(map[string]string{
		"approver_role": "SAFETY_REP",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/energy/ptw/"+permit.ID+"/approve", bytes.NewBuffer(safetyApprovePayload))
	req.Header.Set("Authorization", safetyAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for safety rep approval, got %d: %s", rec.Code, rec.Body.String())
	}

	// Sub-step C: Substation Operator Isolation Declaration
	substationToken := auth.Claims{
		StaffID:     "STAFF-SUB-01",
		StaffNumber: "NES-SUB-01",
		FullName:    "Substation Operator",
		PrimaryRole: "SUBSTATION_OPERATOR",
		OrgUnitID:   "unit-sub-vom",
	}
	substationJwt, _ := auth.GenerateToken(substationToken, testJWTKey)
	substationAuth := "Bearer " + substationJwt

	isolatePayload, _ := json.Marshal(map[string]string{
		"isolation_notes": "CB-04 racked out and earthed cleanly.",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/energy/ptw/"+permit.ID+"/isolate", bytes.NewBuffer(isolatePayload))
	req.Header.Set("Authorization", substationAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for isolation, got %d: %s", rec.Code, rec.Body.String())
	}

	// Sub-step D: Safety Rep Clearance Declaration
	clearancePayload, _ := json.Marshal(map[string]string{
		"notes": "All personnel clear and grounds removed cleanly.",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/energy/ptw/"+permit.ID+"/clear-safety", bytes.NewBuffer(clearancePayload))
	req.Header.Set("Authorization", safetyAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for clearance, got %d: %s", rec.Code, rec.Body.String())
	}

	// Sub-step E: Substation Operator Restoration & Closure
	req = httptest.NewRequest(http.MethodPost, "/api/v1/energy/ptw/"+permit.ID+"/restore", nil)
	req.Header.Set("Authorization", substationAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for restore, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify completed PTW status
	finalPTW, err := s06Repo.GetPTWByID(context.Background(), permit.ID)
	if err != nil || finalPTW.Status != s06domain.PTWRestoredClosed {
		t.Fatalf("expected permit status RESTORED_CLOSED, got %s", finalPTW.Status)
	}
}

