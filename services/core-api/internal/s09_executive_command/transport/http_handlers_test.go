// Package transport_test exercises REST endpoints for S09 Executive Command.
package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	s04repo "disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	s05repo "disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	s06repo "disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	s07repo "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	s08repo "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/service"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/transport"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const testJWTKey = "nesco-exec-jwt-32-byte-secret-key-12345!"

func generateTestToken(t *testing.T, staffID, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(auth.Claims{
		StaffID:                   staffID,
		StaffNumber:               "NES-" + staffID,
		FullName:                  "Officer " + staffID,
		PrimaryRole:               role,
		OrgUnitID:                 "unit-exec-hq",
		CanViewExecutiveDashboard: true,
		IssuedAt:                  time.Now().Unix(),
		ExpiresAt:                 time.Now().Add(15 * time.Minute).Unix(),
	}, testJWTKey)
	if err != nil {
		t.Fatalf("failed to generate auth token: %v", err)
	}
	return token
}

func setupTestServer(t *testing.T) (*http.ServeMux, repository.Repository, *service.WorkflowService, *s04repo.MemoryRepository) {
	t.Helper()
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	maintRepo := s04repo.NewMemoryRepository()
	genRepo := s05repo.NewMemoryRepository()
	energyRepo := s06repo.NewMemoryRepository()
	auditRepo := s07repo.NewMemoryRepository()
	billRepo := s08repo.NewMemoryRepository()

	svc := service.NewWorkflowService(repo, bus, maintRepo, genRepo, energyRepo, auditRepo, billRepo)
	handler := transport.NewHTTPHandler(svc, repo, testJWTKey)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, repo, svc, maintRepo
}

func TestGetDashboardSummary(t *testing.T) {
	mux, _, _, _ := setupTestServer(t)

	// Negative: No auth token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/executive/dashboard", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized without token, got %d", rec.Code)
	}

	// Positive: With valid executive token
	token := generateTestToken(t, "STAFF-MD-01", "MD")
	req = httptest.NewRequest(http.MethodGet, "/api/v1/executive/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	var summary domain.CSuiteDashboardSummary
	if err := json.Unmarshal(rec.Body.Bytes(), &summary); err != nil {
		t.Fatalf("failed to unmarshal summary: %v", err)
	}
	if summary.RevenueSummary.ChannelBreakdown == nil {
		t.Fatalf("expected channel breakdown map initialized")
	}
}

func TestSnapshotsEndpoints(t *testing.T) {
	mux, _, _, _ := setupTestServer(t)
	token := generateTestToken(t, "STAFF-EXEC-01", "MD")

	// 1. Get latest when none exists -> 404
	req := httptest.NewRequest(http.MethodGet, "/api/v1/executive/snapshots/latest", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found for initial snapshot, got %d", rec.Code)
	}

	// 2. Negative: Invalid JSON payload
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/snapshots", bytes.NewBufferString("{invalid_json"))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid json, got %d", rec.Code)
	}

	// 3. Negative: Missing mandatory fields
	invalidDTO := transport.CaptureSnapshotDTO{
		SnapshotDate: "",
	}
	raw, _ := json.Marshal(invalidDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/snapshots", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing mandatory fields, got %d", rec.Code)
	}

	// 4. Positive: Capture snapshot
	validDTO := transport.CaptureSnapshotDTO{
		SnapshotDate:                  "2026-10-22",
		ShiftWindow:                   "MORNING_08_16",
		TotalGenerationMW:             42.5,
		TCNGridIntakeMW:               15.0,
		TotalSystemLoadMW:             56.2,
		FeedersInServiceCount:         28,
		FeedersOutFaultCount:          2,
		FeedersOutPTWCount:            1,
		FeedersOutRationingCount:      1,
		TotalRevenueCollected:         2500000.0,
		PrepaidRevenueSharePct:        88.5,
		UnresolvedTheftAlertsCount:    3,
		UnresolvedChurnAlertsCount:    1,
		AuditPendingProcurementsCount: 2,
		UnreturnedScrapCount:          1,
		Notes:                         "Steady morning peak demand",
		CapturedByStaffID:             "STAFF-EXEC-01",
	}
	raw, _ = json.Marshal(validDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/snapshots", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var created domain.ExecutiveSnapshot
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode created snapshot: %v", err)
	}
	if created.ID == "" || created.SnapshotDate != "2026-10-22" {
		t.Fatalf("unexpected snapshot data: %+v", created)
	}

	// 5. Get latest -> 200
	req = httptest.NewRequest(http.MethodGet, "/api/v1/executive/snapshots/latest", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for latest snapshot, got %d", rec.Code)
	}

	// 6. List snapshots -> 200
	req = httptest.NewRequest(http.MethodGet, "/api/v1/executive/snapshots?limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for list snapshots, got %d", rec.Code)
	}
	var list []domain.ExecutiveSnapshot
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil || len(list) != 1 {
		t.Fatalf("expected 1 snapshot in list, got %d", len(list))
	}
}

func TestActionItemLifecycleEndpoints(t *testing.T) {
	mux, _, _, _ := setupTestServer(t)
	mdToken := generateTestToken(t, "STAFF-MD-01", "MD")
	feederToken := generateTestToken(t, "STAFF-ENG-01", "HEAD_ENERGY_MGMT")

	// 1. Negative Enqueue: Invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions", bytes.NewBufferString("{bad_json"))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad json, got %d", rec.Code)
	}

	// 2. Negative Enqueue: Missing fields
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions", bytes.NewBufferString(`{"action_type":""}`))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing fields, got %d", rec.Code)
	}

	// 3. Positive Enqueue: MD Procurement Requisition
	enqueueDTO := transport.EnqueueActionDTO{
		ActionType:       domain.ActionMDProcurementApproval,
		TargetRole:       domain.RoleMD,
		SourceEntityID:   "REQ-2026-999",
		SourceEntityRef:  "PO-TX-TRANSFORMER-OIL-DRUMS",
		SourceDepartment: "STORES_AUDIT",
		Summary:          "Emergency procurement of 50 drums dielectric oil",
		Amount:           12500000.0,
		Priority:         "CRITICAL",
	}
	raw, _ := json.Marshal(enqueueDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}
	var item domain.ExecutiveActionItem
	if err := json.Unmarshal(rec.Body.Bytes(), &item); err != nil {
		t.Fatalf("failed to decode enqueued action: %v", err)
	}

	// 4. Negative Enqueue: Conflict duplicate source entity ref
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions", bytes.NewBuffer(raw))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict for duplicate entity ref, got %d", rec.Code)
	}

	// 5. List action items
	req = httptest.NewRequest(http.MethodGet, "/api/v1/executive/actions?target_role=MD&status=PENDING", nil)
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for list action items, got %d", rec.Code)
	}

	// 6. Negative Decide: Invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions/"+item.ID+"/decide", bytes.NewBufferString("{bad}"))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad json, got %d", rec.Code)
	}

	// 7. Negative Decide: Missing decision fields
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions/"+item.ID+"/decide", bytes.NewBufferString(`{"decided_by_staff_id":""}`))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing decision fields, got %d", rec.Code)
	}

	// 8. Negative Decide: Unauthorized role (Segregation of Duties breach: HEAD_ENERGY_MGMT cannot sign MD procurement)
	decideDTO := transport.DecideActionDTO{
		DecidedByStaffID: "STAFF-ENG-01",
		Approve:          true,
		Notes:            "Attempting to approve as energy mgmt head",
	}
	decideRaw, _ := json.Marshal(decideDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions/"+item.ID+"/decide", bytes.NewBuffer(decideRaw))
	req.Header.Set("Authorization", "Bearer "+feederToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for unauthorized role, got %d: %s", rec.Code, rec.Body.String())
	}

	// 9. Negative Decide: Non-existent item
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions/non-existent-id/decide", bytes.NewBuffer(decideRaw))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found for bogus action item, got %d", rec.Code)
	}

	// 10. Positive Decide: MD digital approval
	mdDecideDTO := transport.DecideActionDTO{
		DecidedByStaffID: "STAFF-MD-01",
		Approve:          true,
		Notes:            "Approved for immediate emergency procurement. Validated with stores.",
	}
	mdDecideRaw, _ := json.Marshal(mdDecideDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions/"+item.ID+"/decide", bytes.NewBuffer(mdDecideRaw))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for MD digital sign-off, got %d: %s", rec.Code, rec.Body.String())
	}
	var decidedItem domain.ExecutiveActionItem
	if err := json.Unmarshal(rec.Body.Bytes(), &decidedItem); err != nil {
		t.Fatalf("failed to decode decided item: %v", err)
	}
	if decidedItem.Status != domain.ActionStatusApproved {
		t.Fatalf("expected status APPROVED, got %s", decidedItem.Status)
	}

	// 11. Negative Decide: Double decision on already decided action item -> 409 Conflict
	req = httptest.NewRequest(http.MethodPost, "/api/v1/executive/actions/"+item.ID+"/decide", bytes.NewBuffer(mdDecideRaw))
	req.Header.Set("Authorization", "Bearer "+mdToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict for double decision, got %d: %s", rec.Code, rec.Body.String())
	}
}
