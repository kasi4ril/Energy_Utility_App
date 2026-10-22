// Package integration exercises end-to-end multi-sprint workflows across NESCO bounded contexts.
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s03domain "disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	s04domain "disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

func TestE2E_Auth_And_CustomerCare_To_Services_Lifecycle(t *testing.T) {
	mux, _, s03Repo, s04Repo, _, _ := setupE2EEnvironment(t)

	// Step 1: Staff Authentication Handshake (S01)
	loginPayload, _ := json.Marshal(map[string]string{
		"email":    "amina@nesco.ng",
		"password": "SecretP@ss1",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginPayload))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on auth login, got %d: %s", rec.Code, rec.Body.String())
	}
	var loginResp struct {
		Token  string      `json:"token"`
		Claims auth.Claims `json:"claims"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &loginResp)
	if loginResp.Token == "" {
		t.Fatalf("expected non-empty JWT token from login")
	}

	// Prepare role-specific JWT tokens for actor segregation
	leadToken, _ := auth.GenerateToken(auth.Claims{
		StaffID:     "STAFF-LEAD-01",
		StaffNumber: "NES-LEAD-01",
		FullName:    "Team Lead Feeder",
		PrimaryRole: "TEAM_LEAD_FEEDER",
		OrgUnitID:   "unit-feeder-maint",
	}, testJWTKey)
	leadAuth := "Bearer " + leadToken

	headToken, _ := auth.GenerateToken(auth.Claims{
		StaffID:     "STAFF-HEAD-01",
		StaffNumber: "NES-HEAD-01",
		FullName:    "Unit Head Feeder",
		PrimaryRole: "UNIT_HEAD_FEEDER_MAINT",
		OrgUnitID:   "unit-feeder-maint",
	}, testJWTKey)
	headAuth := "Bearer " + headToken

	ccToken, _ := auth.GenerateToken(auth.Claims{
		StaffID:     "STAFF-CC-01",
		StaffNumber: "NES-CC-01",
		FullName:    "Customer Care Officer",
		PrimaryRole: "CUSTOMER_CARE_OFFICER",
		OrgUnitID:   "unit-customer-care",
	}, testJWTKey)
	ccAuth := "Bearer " + ccToken

	// Step 2: KYC Customer Onboarding (S03)
	customerID := "CUST-KURA-2026-001"
	err := s03Repo.CreateCustomer(context.Background(), s03domain.Customer{
		ID:            customerID,
		AccountNo:     "ACC-00142",
		CIN:           "CIN-00142",
		FullName:      "Danladi Gyang",
		Phone:         "08031234567",
		NIN:           "12345678901",
		BuildingClass: s03domain.BuildingResidential,
		PremisesGPS: s01domain.GPSPoint{
			Latitude:  9.7523,
			Longitude: 8.8541,
		},
		TransformerID:    "TX-DT-500KVA-MARKET",
		MeterSerial:      "MTR-98214",
		TariffBand:       "BAND_A",
		KYCStatus:        s03domain.KYCVerified,
		RowHMACSignature: "valid-hmac-sig",
		CreatedAt:        time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to seed customer: %v", err)
	}

	// Step 3: Customer Care logs a Service Ticket (S03)
	ticketPayload, _ := json.Marshal(map[string]interface{}{
		"customer_id":      customerID,
		"channel":          s03domain.ChannelCallCenter,
		"category":         s03domain.CategoryLineSnap,
		"description":      "Feeder jumper burnt at transformer pole junction",
		"assigned_unit_id": "UNIT-FEEDER-MAINT-01",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewBuffer(ticketPayload))
	req.Header.Set("Authorization", ccAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for ticket, got %d: %s", rec.Code, rec.Body.String())
	}
	var createdTicket s03domain.ServiceTicket
	_ = json.Unmarshal(rec.Body.Bytes(), &createdTicket)

	// Step 4: Maintenance Work Order Dispatched to Field Crew (S04)
	woID := "WO-2026-FDR-009"
	now := time.Now().UTC()
	err = s04Repo.CreateWorkOrder(context.Background(), s04domain.MaintenanceWorkOrder{
		ID:            woID,
		OrderNo:       "WO-FDR-009",
		TicketID:      &createdTicket.ID,
		UnitType:      s04domain.UnitFeederMaint,
		WorkOrderType: s04domain.WOCorrectiveRepair,
		Priority:      s04domain.PriorityHigh,
		UnitHeadID:    "STAFF-HEAD-01",
		TeamLeadID:    "STAFF-LEAD-01",
		AssetID:       "FDR-33KV-VOM",
		Status:        s04domain.StatusDispatched,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("failed to create work order: %v", err)
	}

	// Step 5: Field PWA Proof-of-Work submitted with GPS by assigned Team Lead (S04)
	proofPayload, _ := json.Marshal(map[string]interface{}{
		"latitude":          9.7524,
		"longitude":         8.8542,
		"before_photo_url":  "https://storage.disco.local/photos/before_wo_009.jpg",
		"after_photo_url":   "https://storage.disco.local/photos/after_wo_009.jpg",
		"photo_hash_before": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"photo_hash_after":  "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		"work_details":      "Replaced cracked insulator and jumper cleanly on tower 42",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders/"+woID+"/proof", bytes.NewBuffer(proofPayload))
	req.Header.Set("Authorization", leadAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for proof-of-work, got %d: %s", rec.Code, rec.Body.String())
	}

	// Step 6: Unit Head verifies proof
	reviewPayload, _ := json.Marshal(map[string]interface{}{
		"approved": true,
		"notes":    "GPS proximity within tolerance and photos verified.",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/services/work-orders/"+woID+"/verify", bytes.NewBuffer(reviewPayload))
	req.Header.Set("Authorization", headAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for proof review, got %d: %s", rec.Code, rec.Body.String())
	}

	// Step 7: Customer Care marks ATTENDED_TO then customer phone confirmation (S03)
	attendPayload, _ := json.Marshal(map[string]interface{}{
		"proof_note": "Field crew completed tower jumper replacement successfully.",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/tickets/"+createdTicket.ID+"/attend", bytes.NewBuffer(attendPayload))
	req.Header.Set("Authorization", ccAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on ticket attend, got %d: %s", rec.Code, rec.Body.String())
	}

	confirmPayload, _ := json.Marshal(map[string]interface{}{
		"confirmed_by_customer": true,
		"feedback_notes":        "Customer Danladi confirmed power fully restored on his street.",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/tickets/"+createdTicket.ID+"/confirm-close", bytes.NewBuffer(confirmPayload))
	req.Header.Set("Authorization", ccAuth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on 2-step customer confirmation, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify final state
	finalTicket, err := s03Repo.GetTicketByID(context.Background(), createdTicket.ID)
	if err != nil || finalTicket.Status != s03domain.StatusClosed {
		t.Fatalf("expected ticket status CLOSED after customer confirmation gate, got %s", finalTicket.Status)
	}
}

