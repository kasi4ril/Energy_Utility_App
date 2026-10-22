// Package transport_test validates HTTP REST API endpoints and authorization for S03 Customer Care.
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
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/service"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/transport"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const (
	testJWTSecret  = "s03-transport-test-jwt-secret-32b!"
	testHMACSecret = "s03-transport-test-hmac-secret-32b"
)

func setupTestServer() (http.Handler, *service.WorkflowService, *repository.MemoryRepository, *s01repo.MemoryRepository) {
	repo := repository.NewMemoryRepository()
	gridRepo := s01repo.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, gridRepo, bus, testHMACSecret)
	handler := transport.NewHTTPHandler(svc, repo, testJWTSecret)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, svc, repo, gridRepo
}

func makeAuthHeader(staffID, role, fullName string) string {
	claims := auth.Claims{
		StaffID:     staffID,
		PrimaryRole: role,
		FullName:    fullName,
	}
	token, _ := auth.GenerateToken(claims, testJWTSecret)
	return "Bearer " + token
}

func TestHTTP_CustomerOnboardingFlow(t *testing.T) {
	mux, _, _, gridRepo := setupTestServer()
	ctx := context.Background()

	// Create Distribution Transformer in grid assets
	dt := s01domain.GridAsset{
		ID:               "tx-dt-001",
		AssetCode:        "DT-RAYFIELD-01",
		AssetType:        s01domain.AssetTransformerDT,
		OrgUnitID:        "org-unit-1",
		GPSPoint:         s01domain.GPSPoint{Latitude: 9.892, Longitude: 8.875},
		OperationalState: s01domain.StateInService,
	}
	_ = gridRepo.CreateGridAsset(ctx, dt)

	ccAuth := makeAuthHeader("officer-1", "CUSTOMER_CARE_OFFICER", "Zainab Ahmed")

	// 1. Onboard Customer
	onboardBody, _ := json.Marshal(transport.OnboardCustomerDTO{
		FullName:      "Kabiru Mohammed",
		Phone:         "+2348039991122",
		NIN:           "NIN-123456789",
		BuildingClass: domain.BuildingResidential,
		Latitude:      9.8923,
		Longitude:     8.8752,
		TransformerID: dt.ID,
		MeterSerial:   "MTR-PRE-099",
		TariffBand:    s01domain.TariffBandA,
	})
	req := httptest.NewRequest("POST", "/api/v1/customers/onboard", bytes.NewReader(onboardBody))
	req.Header.Set("Authorization", ccAuth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on customer onboarding, got %d: %s", rec.Code, rec.Body.String())
	}

	var createdCust domain.Customer
	_ = json.NewDecoder(rec.Body).Decode(&createdCust)

	// Negative: Empty onboarding payload -> 400
	badReq := httptest.NewRequest("POST", "/api/v1/customers/onboard", bytes.NewReader([]byte("{}")))
	badReq.Header.Set("Authorization", ccAuth)
	badRec := httptest.NewRecorder()
	mux.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on empty onboarding payload")
	}

	// Negative: Malformed JSON -> 400
	malformedReq := httptest.NewRequest("POST", "/api/v1/customers/onboard", bytes.NewReader([]byte("{bad-json")))
	malformedReq.Header.Set("Authorization", ccAuth)
	malformedRec := httptest.NewRecorder()
	mux.ServeHTTP(malformedRec, malformedReq)
	if malformedRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on malformed json")
	}

	// Negative: Anti-Ghost Lock triggers on distant GPS -> 422
	ghostBody, _ := json.Marshal(transport.OnboardCustomerDTO{
		FullName:      "Ghost Customer",
		Phone:         "+2348039994455",
		BuildingClass: domain.BuildingResidential,
		Latitude:      9.9500,
		Longitude:     8.9500, // Miles away from DT
		TransformerID: dt.ID,
		TariffBand:    s01domain.TariffBandA,
	})
	ghostReq := httptest.NewRequest("POST", "/api/v1/customers/onboard", bytes.NewReader(ghostBody))
	ghostReq.Header.Set("Authorization", ccAuth)
	ghostRec := httptest.NewRecorder()
	mux.ServeHTTP(ghostRec, ghostReq)
	if ghostRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 Anti-Ghost rejection, got %d", ghostRec.Code)
	}

	// 2. List Customers
	listReq := httptest.NewRequest("GET", "/api/v1/customers?transformer_id="+dt.ID, nil)
	listReq.Header.Set("Authorization", ccAuth)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing customers, got %d", listRec.Code)
	}

	// 3. Get Customer by ID
	getReq := httptest.NewRequest("GET", "/api/v1/customers/"+createdCust.ID, nil)
	getReq.Header.Set("Authorization", ccAuth)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK fetching customer by ID, got %d", getRec.Code)
	}

	// 4. Get non-existent customer -> 404
	missingReq := httptest.NewRequest("GET", "/api/v1/customers/missing-uuid", nil)
	missingReq.Header.Set("Authorization", ccAuth)
	missingRec := httptest.NewRecorder()
	mux.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing customer, got %d", missingRec.Code)
	}
}

func TestHTTP_TicketDirectRoutingAndClosureFlow(t *testing.T) {
	mux, _, _, _ := setupTestServer()

	ccAuth := makeAuthHeader("cc-officer-1", "CUSTOMER_CARE_OFFICER", "Fatima Aliyu")
	feederAuth := makeAuthHeader("feeder-lead-1", "TEAM_LEAD_FEEDER", "Engr. Bitrus")

	// 1. Customer Care logs ticket (Fallen Pole)
	logBody, _ := json.Marshal(transport.LogTicketDTO{
		Channel:        domain.ChannelCallCenter,
		Category:       domain.CategoryFallenPole,
		Description:    "Pole tilted and touching trees",
		AssignedUnitID: "unit-feeder-1",
	})
	req := httptest.NewRequest("POST", "/api/v1/tickets", bytes.NewReader(logBody))
	req.Header.Set("Authorization", ccAuth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on ticket log, got %d: %s", rec.Code, rec.Body.String())
	}

	var ticket domain.ServiceTicket
	_ = json.NewDecoder(rec.Body).Decode(&ticket)

	// Negative: Log ticket with empty payload -> 400
	badLog := httptest.NewRequest("POST", "/api/v1/tickets", bytes.NewReader([]byte("{}")))
	badLog.Header.Set("Authorization", ccAuth)
	badRec := httptest.NewRecorder()
	mux.ServeHTTP(badRec, badLog)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on empty ticket payload")
	}

	// 2. Technical staff attends ticket with proof
	attendBody, _ := json.Marshal(transport.AttendTicketDTO{ProofNote: "Trees trimmed and pole re-aligned"})
	attendReq := httptest.NewRequest("POST", "/api/v1/tickets/"+ticket.ID+"/attend", bytes.NewReader(attendBody))
	attendReq.Header.Set("Authorization", feederAuth)
	attendRec := httptest.NewRecorder()
	mux.ServeHTTP(attendRec, attendReq)
	if attendRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on attend ticket, got %d: %s", attendRec.Code, attendRec.Body.String())
	}

	// Negative: Attend with empty proof -> 400
	badAttend := httptest.NewRequest("POST", "/api/v1/tickets/"+ticket.ID+"/attend", bytes.NewReader([]byte(`{"proof_note":""}`)))
	badAttend.Header.Set("Authorization", feederAuth)
	badAttendRec := httptest.NewRecorder()
	mux.ServeHTTP(badAttendRec, badAttend)
	if badAttendRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty proof note")
	}

	// 3. Negative: Technical staff attempts to execute confirmation closure -> 422
	closeBody, _ := json.Marshal(transport.ConfirmClosureDTO{ConfirmedByCustomer: true})
	badCloseReq := httptest.NewRequest("POST", "/api/v1/tickets/"+ticket.ID+"/confirm-close", bytes.NewReader(closeBody))
	badCloseReq.Header.Set("Authorization", feederAuth) // Unauthorized role!
	badCloseRec := httptest.NewRecorder()
	mux.ServeHTTP(badCloseRec, badCloseReq)
	if badCloseRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for non-customer care closure attempt, got %d", badCloseRec.Code)
	}

	// 4. Positive: Customer Care calls customer, confirms fix -> 200 CLOSED
	closeReq := httptest.NewRequest("POST", "/api/v1/tickets/"+ticket.ID+"/confirm-close", bytes.NewReader(closeBody))
	closeReq.Header.Set("Authorization", ccAuth)
	closeRec := httptest.NewRecorder()
	mux.ServeHTTP(closeRec, closeReq)
	if closeRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on customer care closure confirmation, got %d: %s", closeRec.Code, closeRec.Body.String())
	}

	// 5. Query Ticket Details
	getTicketReq := httptest.NewRequest("GET", "/api/v1/tickets/"+ticket.ID, nil)
	getTicketReq.Header.Set("Authorization", ccAuth)
	getTicketRec := httptest.NewRecorder()
	mux.ServeHTTP(getTicketRec, getTicketReq)
	if getTicketRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK fetching ticket details, got %d", getTicketRec.Code)
	}

	// Negative: Missing ticket -> 404
	missingTicket := httptest.NewRequest("GET", "/api/v1/tickets/missing-id", nil)
	missingTicket.Header.Set("Authorization", ccAuth)
	missingRec := httptest.NewRecorder()
	mux.ServeHTTP(missingRec, missingTicket)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing ticket, got %d", missingRec.Code)
	}

	// 6. List Tickets
	listTickets := httptest.NewRequest("GET", "/api/v1/tickets?status=CLOSED&category=FALLEN_POLE", nil)
	listTickets.Header.Set("Authorization", ccAuth)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listTickets)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing tickets, got %d", listRec.Code)
	}
}

func TestHTTP_MeterScheduling(t *testing.T) {
	mux, _, repo, _ := setupTestServer()
	ctx := context.Background()

	ccAuth := makeAuthHeader("cc-officer-1", "CUSTOMER_CARE_OFFICER", "Fatima Aliyu")

	cust := domain.Customer{
		ID:        "cust-sch-1",
		AccountNo: "ACC-SCH-01",
		FullName:  "Joshua Nengel",
		Phone:     "+2348028889900",
	}
	_ = repo.CreateCustomer(ctx, cust)

	// 1. Schedule meter installation
	schedBody, _ := json.Marshal(transport.ScheduleMeterDTO{
		CustomerID:    cust.ID,
		ScheduleType:  domain.ScheduleNewInstallation,
		ScheduledDate: time.Now().Add(24 * time.Hour),
	})
	req := httptest.NewRequest("POST", "/api/v1/meters/schedule", bytes.NewReader(schedBody))
	req.Header.Set("Authorization", ccAuth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on meter schedule, got %d: %s", rec.Code, rec.Body.String())
	}

	// Negative: Empty schedule payload -> 400
	badSched := httptest.NewRequest("POST", "/api/v1/meters/schedule", bytes.NewReader([]byte("{}")))
	badSched.Header.Set("Authorization", ccAuth)
	badRec := httptest.NewRecorder()
	mux.ServeHTTP(badRec, badSched)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on empty schedule payload")
	}

	// Negative: Schedule for missing customer -> 422
	missingCustSched, _ := json.Marshal(transport.ScheduleMeterDTO{
		CustomerID:    "missing-cust-id",
		ScheduleType:  domain.ScheduleRepair,
		ScheduledDate: time.Now().Add(24 * time.Hour),
	})
	missingReq := httptest.NewRequest("POST", "/api/v1/meters/schedule", bytes.NewReader(missingCustSched))
	missingReq.Header.Set("Authorization", ccAuth)
	missingRec := httptest.NewRecorder()
	mux.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 scheduling for missing customer, got %d", missingRec.Code)
	}

	// 2. List meter schedules
	listReq := httptest.NewRequest("GET", "/api/v1/meters/schedules?status=SCHEDULED", nil)
	listReq.Header.Set("Authorization", ccAuth)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing meter schedules, got %d", listRec.Code)
	}
}

