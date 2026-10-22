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
	s03domain "disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const testSecret = "test-s04-super-secret-jwt-key-32b!"

func setupHTTPTestEnv(t *testing.T) (*HTTPHandler, *http.ServeMux, *repository.MemoryRepository) {
	t.Helper()

	repo := repository.NewMemoryRepository()
	gridRepo := s01repo.NewMemoryRepository()
	ticketRepo := s03repo.NewMemoryRepository()
	bus := events.NewMemoryEventBus()

	// Seed Grid Asset
	dtAsset := s01domain.GridAsset{
		ID:        "dt-asset-1",
		AssetCode: "DT-001",
		AssetType: s01domain.AssetTransformerDT,
		GPSPoint:  s01domain.GPSPoint{Latitude: 9.9320, Longitude: 8.8920},
	}
	_ = gridRepo.CreateGridAsset(context.Background(), dtAsset)

	// Seed Ticket
	ticket := s03domain.ServiceTicket{
		ID:             "ticket-1",
		TicketNo:       "TIC-1",
		Channel:        s03domain.ChannelCallCenter,
		Category:       s03domain.CategoryTransformerFault,
		Description:    "Transformer tripped.",
		AssignedUnitID: "unit-1",
		Status:         s03domain.StatusLogged,
		SLADueAt:       time.Now().Add(24 * time.Hour),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	_ = ticketRepo.CreateTicket(context.Background(), ticket)

	svc := service.NewWorkflowService(repo, gridRepo, ticketRepo, bus, "")
	handler := NewHTTPHandler(svc, repo, testSecret)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return handler, mux, repo
}

func generateTestToken(t *testing.T, staffID, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(auth.Claims{
		StaffID:     staffID,
		StaffNumber: "NES-" + staffID,
		FullName:    "Engineer " + staffID,
		PrimaryRole: role,
		OrgUnitID:   "unit-1",
		IssuedAt:    time.Now().Unix(),
		ExpiresAt:   time.Now().Add(15 * time.Minute).Unix(),
	}, testSecret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}

func TestHTTP_WorkOrderFlow(t *testing.T) {
	_, mux, _ := setupHTTPTestEnv(t)
	unitHeadToken := generateTestToken(t, "unit-head-1", "UNIT_HEAD_TX_MAINT")

	// 1. Create Work Order
	dto := CreateWorkOrderDTO{
		UnitType:      domain.UnitTransformerMaint,
		WorkOrderType: domain.WOCorrectiveRepair,
		Priority:      domain.PriorityHigh,
		TeamLeadID:    "team-lead-1",
		AssetID:       "dt-asset-1",
		Instructions:  "Inspect fuse box and top up cooling oil.",
	}
	body, _ := json.Marshal(dto)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+unitHeadToken)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var createdWO domain.MaintenanceWorkOrder
	_ = json.NewDecoder(rec.Body).Decode(&createdWO)

	// 2. List Work Orders
	req = httptest.NewRequest(http.MethodGet, "/api/v1/services/work-orders?unit=TRANSFORMER_MAINT", nil)
	req.Header.Set("Authorization", "Bearer "+unitHeadToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	// 3. Get Work Order by ID
	req = httptest.NewRequest(http.MethodGet, "/api/v1/services/work-orders/"+createdWO.ID, nil)
	req.Header.Set("Authorization", "Bearer "+unitHeadToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GetByID, got %d", rec.Code)
	}

	// 4. Get Non-Existent Work Order
	req = httptest.NewRequest(http.MethodGet, "/api/v1/services/work-orders/non-existent", nil)
	req.Header.Set("Authorization", "Bearer "+unitHeadToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", rec.Code)
	}
}

func TestHTTP_ProofSubmissionAndVerificationFlow(t *testing.T) {
	_, mux, _ := setupHTTPTestEnv(t)
	unitHeadToken := generateTestToken(t, "unit-head-1", "UNIT_HEAD_TX_MAINT")
	teamLeadToken := generateTestToken(t, "team-lead-1", "FIELD_TEAM_LEAD")

	ticketID := "ticket-1"
	// Create Work Order
	woDTO := CreateWorkOrderDTO{
		TicketID:      &ticketID,
		UnitType:      domain.UnitTransformerMaint,
		WorkOrderType: domain.WOCorrectiveRepair,
		Priority:      domain.PriorityCritical,
		TeamLeadID:    "team-lead-1",
		AssetID:       "dt-asset-1",
		Instructions:  "Emergency repair on burnt secondary fuse.",
	}
	body, _ := json.Marshal(woDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+unitHeadToken)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var wo domain.MaintenanceWorkOrder
	_ = json.NewDecoder(rec.Body).Decode(&wo)

	// 1. Submit Proof of Work
	proofDTO := SubmitProofDTO{
		Latitude:        9.9321,
		Longitude:       8.8920,
		BeforePhotoURL:  "/photos/before_wo.jpg",
		AfterPhotoURL:   "/photos/after_wo.jpg",
		PhotoHashBefore: "1111111111111111111111111111111111111111111111111111111111111111",
		PhotoHashAfter:  "2222222222222222222222222222222222222222222222222222222222222222",
		WorkDetails:     "Replaced burnt 200A fuse and cleaned transformer terminal lugs.",
	}
	proofBody, _ := json.Marshal(proofDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders/"+wo.ID+"/proof", bytes.NewReader(proofBody))
	req.Header.Set("Authorization", "Bearer "+teamLeadToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on proof submit, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. Get Proof by Work Order ID
	req = httptest.NewRequest(http.MethodGet, "/api/v1/services/work-orders/"+wo.ID+"/proof", nil)
	req.Header.Set("Authorization", "Bearer "+unitHeadToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on get proof, got %d", rec.Code)
	}

	// 3. Get Proof for non-existent work order
	req = httptest.NewRequest(http.MethodGet, "/api/v1/services/work-orders/missing/proof", nil)
	req.Header.Set("Authorization", "Bearer "+unitHeadToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing proof, got %d", rec.Code)
	}

	// 4. Verify Proof (Approval)
	verifyDTO := VerifyProofDTO{
		Approved: true,
		Notes:    "Verified fuse integrity and load parameters.",
	}
	vBody, _ := json.Marshal(verifyDTO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders/"+wo.ID+"/verify", bytes.NewReader(vBody))
	req.Header.Set("Authorization", "Bearer "+unitHeadToken)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on verify proof, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHTTP_NegativeValidations(t *testing.T) {
	_, mux, _ := setupHTTPTestEnv(t)
	token := generateTestToken(t, "unit-head-1", "UNIT_HEAD_TX_MAINT")

	// 1. Invalid JSON body on CreateWorkOrder
	req := httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request on invalid JSON, got %d", rec.Code)
	}

	// 2. Validation failure on CreateWorkOrder
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders", bytes.NewReader([]byte("{}")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request on empty CreateDTO, got %d", rec.Code)
	}

	// 3. Invalid JSON on SubmitProof
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders/wo-1/proof", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on invalid JSON in SubmitProof, got %d", rec.Code)
	}

	// 4. Validation failure on SubmitProof (out of range lat)
	badProof := SubmitProofDTO{
		Latitude: 100.0, // Invalid latitude
	}
	bpBody, _ := json.Marshal(badProof)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders/wo-1/proof", bytes.NewReader(bpBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on bad latitude, got %d", rec.Code)
	}

	// 5. Invalid JSON on VerifyProof
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders/wo-1/verify", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on invalid JSON in VerifyProof, got %d", rec.Code)
	}

	// 6. Validation failure on VerifyProof (rejecting without notes)
	badVerify := VerifyProofDTO{
		Approved: false,
		Notes:    "bad", // short notes
	}
	bvBody, _ := json.Marshal(badVerify)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders/wo-1/verify", bytes.NewReader(bvBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 on short rework notes, got %d", rec.Code)
	}

	// 7. Service error on CreateWorkOrder (non-existent asset)
	badAssetWO := CreateWorkOrderDTO{
		UnitType:      domain.UnitTransformerMaint,
		WorkOrderType: domain.WOCorrectiveRepair,
		Priority:      domain.PriorityHigh,
		TeamLeadID:    "team-lead-1",
		AssetID:       "missing-asset-id",
		Instructions:  "Inspect fuse box and top up cooling oil.",
	}
	baBody, _ := json.Marshal(badAssetWO)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders", bytes.NewReader(baBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 on missing asset, got %d", rec.Code)
	}
}

func TestDTO_ComprehensiveValidation(t *testing.T) {
	// CreateWorkOrderDTO branches
	woDTO := CreateWorkOrderDTO{UnitType: "FEEDER_MAINT", WorkOrderType: "CORRECTIVE_REPAIR", TeamLeadID: "", AssetID: "A1"}
	if err := woDTO.Validate(); err == nil {
		t.Error("expected error on empty TeamLeadID")
	}
	woDTO = CreateWorkOrderDTO{UnitType: "FEEDER_MAINT", WorkOrderType: "CORRECTIVE_REPAIR", TeamLeadID: "T1", AssetID: "A1", Instructions: "short"}
	if err := woDTO.Validate(); err == nil {
		t.Error("expected error on short instructions")
	}
	woDTO = CreateWorkOrderDTO{UnitType: "FEEDER_MAINT", WorkOrderType: "CORRECTIVE_REPAIR", TeamLeadID: "T1", AssetID: "A1", Instructions: "Valid instructions", PTWRequired: true}
	if err := woDTO.Validate(); err == nil {
		t.Error("expected error on missing PTW permit ID")
	}

	// SubmitProofDTO branches
	pDTO := SubmitProofDTO{Latitude: 9.0, Longitude: 8.0, BeforePhotoURL: "", AfterPhotoURL: "b"}
	if err := pDTO.Validate(); err == nil {
		t.Error("expected error on empty BeforePhotoURL")
	}
	pDTO = SubmitProofDTO{Latitude: 9.0, Longitude: 8.0, BeforePhotoURL: "a", AfterPhotoURL: "b", PhotoHashBefore: "short", PhotoHashAfter: "short"}
	if err := pDTO.Validate(); err == nil {
		t.Error("expected error on short photo hashes")
	}
	pDTO = SubmitProofDTO{Latitude: 9.0, Longitude: 8.0, BeforePhotoURL: "a", AfterPhotoURL: "b", PhotoHashBefore: "11111111111111111111111111111111", PhotoHashAfter: "22222222222222222222222222222222", WorkDetails: "short"}
	if err := pDTO.Validate(); err == nil {
		t.Error("expected error on short work details")
	}
}
