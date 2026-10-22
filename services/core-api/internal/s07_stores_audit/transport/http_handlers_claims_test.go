// Package transport_test exercises the staff allowance claims desk and internal auditor review endpoints.
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

func TestHTTPHandler_StaffClaims(t *testing.T) {
	mux, _ := setupTestServer(t)

	tokenStaff := generateTestToken(t, "staff-eng-01", "MAINTENANCE_TECH")
	tokenAuditor := generateTestToken(t, "auditor-01", "INTERNAL_AUDITOR")

	// 1. Submit Claim -> 201 Created (Verified with Roster + PTW + GPS)
	ptwID := "ptw-2026-bukuru-42"
	claimDTO := transport.SubmitClaimDTO{
		ClaimType:           domain.ClaimOvertime,
		Amount:              18000.0,
		ClaimDate:           "2026-10-22",
		LinkedPTWOrTicketID: &ptwID,
		RosterMatchVerified: true,
		PTWCrewVerified:     true,
		GPSMatchVerified:    true,
	}
	body, _ := json.Marshal(claimDTO)

	req := httptest.NewRequest("POST", "/api/v1/audit/claims", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+tokenStaff)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on claim submission, got: %d, body: %s", w.Code, w.Body.String())
	}

	var created domain.StaffAllowanceClaim
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	if created.Status != domain.StatusClaimAIVerified || created.OllamaAuditFlag != domain.AIFlagVerified {
		t.Fatalf("expected AI_VERIFIED claim, got status=%s, flag=%s", created.Status, created.OllamaAuditFlag)
	}

	// 2. Query Claim by ID -> 200 OK
	reqGet := httptest.NewRequest("GET", "/api/v1/audit/claims/"+created.ID, nil)
	reqGet.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK getting claim, got: %d", wGet.Code)
	}

	// 3. Auditor Reviews and Approves Claim -> 200 OK
	reviewDTO := transport.ReviewClaimDTO{
		Approve: true,
		Notes:   "All 3 sources verified cleanly",
	}
	bodyRev, _ := json.Marshal(reviewDTO)
	reqRev := httptest.NewRequest("POST", "/api/v1/audit/claims/"+created.ID+"/review", bytes.NewReader(bodyRev))
	reqRev.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wRev := httptest.NewRecorder()
	mux.ServeHTTP(wRev, reqRev)

	if wRev.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on auditor review, got: %d", wRev.Code)
	}

	var approved domain.StaffAllowanceClaim
	_ = json.Unmarshal(wRev.Body.Bytes(), &approved)
	if approved.Status != domain.StatusClaimAuditApproved {
		t.Fatalf("expected AUDIT_APPROVED, got: %s", approved.Status)
	}

	// 4. List Claims -> 200 OK
	reqList := httptest.NewRequest("GET", "/api/v1/audit/claims?status=AUDIT_APPROVED", nil)
	reqList.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wList := httptest.NewRecorder()
	mux.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing claims, got: %d", wList.Code)
	}

	// 5. Negative: Invalid JSON on claim submission -> 400 Bad Request
	reqBadJSON := httptest.NewRequest("POST", "/api/v1/audit/claims", bytes.NewBufferString("not-json"))
	reqBadJSON.Header.Set("Authorization", "Bearer "+tokenStaff)
	wBadJSON := httptest.NewRecorder()
	mux.ServeHTTP(wBadJSON, reqBadJSON)
	if wBadJSON.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request on malformed JSON, got: %d", wBadJSON.Code)
	}

	// 6. Negative: Non-existent claim review -> 422
	reqNotFoundRev := httptest.NewRequest("POST", "/api/v1/audit/claims/non-existent-claim/review", bytes.NewReader(bodyRev))
	reqNotFoundRev.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wNotFoundRev := httptest.NewRecorder()
	mux.ServeHTTP(wNotFoundRev, reqNotFoundRev)
	if wNotFoundRev.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for non-existent claim review, got: %d", wNotFoundRev.Code)
	}

	// 7. Negative: Query non-existent claim by ID -> 404 Not Found
	reqGetNonEx := httptest.NewRequest("GET", "/api/v1/audit/claims/non-existent-claim-id", nil)
	reqGetNonEx.Header.Set("Authorization", "Bearer "+tokenAuditor)
	wGetNonEx := httptest.NewRecorder()
	mux.ServeHTTP(wGetNonEx, reqGetNonEx)
	if wGetNonEx.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unknown claim, got: %d", wGetNonEx.Code)
	}

	// 8. Negative: Validation error on claim submit (negative amount) -> 400 Bad Request
	badDTO := transport.SubmitClaimDTO{
		ClaimType: domain.ClaimOvertime,
		Amount:    -100.0,
		ClaimDate: "2026-10-22",
	}
	bodyBad, _ := json.Marshal(badDTO)
	reqValErr := httptest.NewRequest("POST", "/api/v1/audit/claims", bytes.NewReader(bodyBad))
	reqValErr.Header.Set("Authorization", "Bearer "+tokenStaff)
	wValErr := httptest.NewRecorder()
	mux.ServeHTTP(wValErr, reqValErr)
	if wValErr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative amount, got: %d", wValErr.Code)
	}
}
