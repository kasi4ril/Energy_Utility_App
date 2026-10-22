// Package transport_test validates HTTP REST API endpoints, DTO validations,
// and authentication barriers for S01 Foundation.
package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/service"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/transport"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const (
	testHMACSecret = "transport-test-hmac-secret-key-32b"
	testJWTSecret  = "transport-test-jwt-secret-key-32b!"
)

func setupTestServer() (http.Handler, *service.WorkflowService, *repository.MemoryRepository) {
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, bus, testHMACSecret, testJWTSecret)
	handler := transport.NewHTTPHandler(svc, repo, testJWTSecret)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, svc, repo
}

func makeAuthHeader(staffID, role, orgUnitID string) string {
	claims := auth.Claims{
		StaffID:     staffID,
		PrimaryRole: role,
		OrgUnitID:   orgUnitID,
	}
	token, _ := auth.GenerateToken(claims, testJWTSecret)
	return "Bearer " + token
}

func TestHTTP_LoginEndpoint(t *testing.T) {
	mux, svc, _ := setupTestServer()
	ctx := context.Background()

	hq, _ := svc.RegisterOrgUnit(ctx, "HQ", "HQ Main", domain.OrgUnitHQ, nil, nil)
	_, _ = svc.RegisterStaffProfile(ctx, "NSC-01", "Amina Bello", "amina@nesco.ng", "SecretP@ss1", domain.RoleMD, hq.ID)

	// Valid login
	loginBody, _ := json.Marshal(transport.LoginRequest{Email: "amina@nesco.ng", Password: "SecretP@ss1"})
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(loginBody))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on login, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp transport.LoginResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Token == "" || resp.Claims.PrimaryRole != "MD" {
		t.Fatalf("unexpected login response: %+v", resp)
	}

	// Invalid password -> 401
	badLogin, _ := json.Marshal(transport.LoginRequest{Email: "amina@nesco.ng", Password: "WrongPassword"})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(badLogin))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 on bad password, got %d", rec.Code)
	}

	// Empty credentials -> 400
	emptyLogin, _ := json.Marshal(transport.LoginRequest{Email: "", Password: ""})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(emptyLogin))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on empty login, got %d", rec.Code)
	}

	// Malformed JSON -> 400
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader([]byte("{invalid-json")))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on malformed json, got %d", rec.Code)
	}
}

func TestHTTP_OrgUnitsAndStaffFlow(t *testing.T) {
	mux, _, _ := setupTestServer()
	authHeader := makeAuthHeader("admin-id", "MD", "hq-id")

	// 1. Create HQ
	lat, lon := 9.9281, 8.8921
	orgReq, _ := json.Marshal(transport.RegisterOrgUnitRequest{
		UnitCode:  "HQ-01",
		Name:      "Headquarters Plateau",
		UnitType:  domain.OrgUnitHQ,
		Latitude:  &lat,
		Longitude: &lon,
	})
	req := httptest.NewRequest("POST", "/api/v1/foundation/org-units", bytes.NewReader(orgReq))
	req.Header.Set("Authorization", authHeader)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for org unit, got %d: %s", rec.Code, rec.Body.String())
	}

	var createdOrg domain.OrgUnit
	_ = json.NewDecoder(rec.Body).Decode(&createdOrg)

	// Create District under HQ
	distReq, _ := json.Marshal(transport.RegisterOrgUnitRequest{
		UnitCode: "DIST-01",
		Name:     "District Central",
		UnitType: domain.OrgUnitDistrict,
		ParentID: &createdOrg.ID,
	})
	req = httptest.NewRequest("POST", "/api/v1/foundation/org-units", bytes.NewReader(distReq))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for district, got %d", rec.Code)
	}
	var createdDist domain.OrgUnit
	_ = json.NewDecoder(rec.Body).Decode(&createdDist)

	// Negative org unit validation -> 400
	badOrg, _ := json.Marshal(transport.RegisterOrgUnitRequest{})
	req = httptest.NewRequest("POST", "/api/v1/foundation/org-units", bytes.NewReader(badOrg))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty org unit request, got %d", rec.Code)
	}

	// Negative org unit domain rule -> 422
	ruleBadOrg, _ := json.Marshal(transport.RegisterOrgUnitRequest{
		UnitCode: "BAD-SU",
		Name:     "Bad Service Unit",
		UnitType: domain.OrgUnitServiceUnit,
		ParentID: &createdOrg.ID, // Service unit directly under HQ is illegal
	})
	req = httptest.NewRequest("POST", "/api/v1/foundation/org-units", bytes.NewReader(ruleBadOrg))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for illegal hierarchy, got %d", rec.Code)
	}

	// 2. List Org Units
	req = httptest.NewRequest("GET", "/api/v1/foundation/org-units", nil)
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on list org units, got %d", rec.Code)
	}

	// 3. Register Staff
	staffReq, _ := json.Marshal(transport.RegisterStaffRequest{
		StaffNumber: "NSC-100",
		FullName:    "Yakubu Gowon",
		Email:       "yakubu@nesco.ng",
		Password:    "SecurePass#1",
		Role:        domain.RoleGMGeneration,
		OrgUnitID:   createdOrg.ID,
	})
	req = httptest.NewRequest("POST", "/api/v1/foundation/staff", bytes.NewReader(staffReq))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for staff, got %d: %s", rec.Code, rec.Body.String())
	}
	var createdStaff domain.StaffProfile
	_ = json.NewDecoder(rec.Body).Decode(&createdStaff)

	// Negative staff validation -> 400
	badStaff, _ := json.Marshal(transport.RegisterStaffRequest{})
	req = httptest.NewRequest("POST", "/api/v1/foundation/staff", bytes.NewReader(badStaff))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on empty staff payload, got %d", rec.Code)
	}

	// 4. Transfer Staff
	transReq, _ := json.Marshal(transport.TransferStaffRequest{
		StaffID:      createdStaff.ID,
		NewOrgUnitID: createdDist.ID,
		NewRole:      domain.RoleHeadServices,
	})
	req = httptest.NewRequest("POST", "/api/v1/foundation/staff/transfer", bytes.NewReader(transReq))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on staff transfer, got %d: %s", rec.Code, rec.Body.String())
	}

	// Negative transfer validation -> 400
	badTrans, _ := json.Marshal(transport.TransferStaffRequest{})
	req = httptest.NewRequest("POST", "/api/v1/foundation/staff/transfer", bytes.NewReader(badTrans))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on empty transfer request, got %d", rec.Code)
	}

	// 5. List Staff
	req = httptest.NewRequest("GET", "/api/v1/foundation/staff?org_unit_id="+createdDist.ID, nil)
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing staff, got %d", rec.Code)
	}
}

func TestHTTP_GridAssetAndWorkNotes(t *testing.T) {
	mux, svc, _ := setupTestServer()
	ctx := context.Background()

	hq, _ := svc.RegisterOrgUnit(ctx, "HQ", "HQ Main", domain.OrgUnitHQ, nil, nil)
	staff, _ := svc.RegisterStaffProfile(ctx, "NSC-ENG", "Engr. Nnamdi", "nnamdi@nesco.ng", "Pass#1", domain.RoleHeadServices, hq.ID)
	authHeader := makeAuthHeader(staff.ID, "HEAD_SERVICES", hq.ID)

	// 1. Register Plant Asset
	assetReq, _ := json.Marshal(transport.RegisterGridAssetRequest{
		AssetCode:        "PLANT-KURRA-01",
		AssetType:        domain.AssetPlant,
		OrgUnitID:        hq.ID,
		RatedCapacityKVA: 12000.0,
		Latitude:         9.428,
		Longitude:        8.892,
	})
	req := httptest.NewRequest("POST", "/api/v1/foundation/assets", bytes.NewReader(assetReq))
	req.Header.Set("Authorization", authHeader)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for asset, got %d: %s", rec.Code, rec.Body.String())
	}

	var createdAsset domain.GridAsset
	_ = json.NewDecoder(rec.Body).Decode(&createdAsset)

	// Negative asset validation -> 400
	badAsset, _ := json.Marshal(transport.RegisterGridAssetRequest{})
	req = httptest.NewRequest("POST", "/api/v1/foundation/assets", bytes.NewReader(badAsset))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty asset request, got %d", rec.Code)
	}

	// 2. List Assets
	req = httptest.NewRequest("GET", "/api/v1/foundation/assets?asset_type=PLANT", nil)
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing assets, got %d", rec.Code)
	}

	// 3. Append Work Note
	noteReq, _ := json.Marshal(transport.AppendWorkNoteRequest{
		AuthorStaffID: staff.ID,
		TaggedUnitIDs: []string{hq.ID},
		NoteText:      "Inspected intake penstock and generator unit 1",
	})
	req = httptest.NewRequest("POST", "/api/v1/notes/GRID_ASSET/"+createdAsset.ID, bytes.NewReader(noteReq))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for work note, got %d: %s", rec.Code, rec.Body.String())
	}

	// Negative note validation -> 400
	badNote, _ := json.Marshal(transport.AppendWorkNoteRequest{})
	req = httptest.NewRequest("POST", "/api/v1/notes/GRID_ASSET/"+createdAsset.ID, bytes.NewReader(badNote))
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty note request, got %d", rec.Code)
	}

	// 4. List Work Notes
	req = httptest.NewRequest("GET", "/api/v1/notes/GRID_ASSET/"+createdAsset.ID, nil)
	req.Header.Set("Authorization", authHeader)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing notes, got %d", rec.Code)
	}

	var notes []domain.WorkNote
	_ = json.NewDecoder(rec.Body).Decode(&notes)
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
}

