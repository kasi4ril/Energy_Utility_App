// Package transport_test exercises REST endpoints for adjustments and revenue protection AI for S08.
package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	s03domain "disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/transport"
)

func TestHTTPHandler_Adjustments(t *testing.T) {
	mux, _, _, ticketRepo := setupTestServer(t)
	token := generateTestToken(t, "BILL-01", "BILLING_OFFICER")
	headToken := generateTestToken(t, "HEAD-COMM-01", "HEAD_COMMERCIAL_OPS")

	// Prepopulate Customer Care ticket
	custID := "CUST-001"
	ticket := s03domain.ServiceTicket{
		ID:             "TICK-001",
		TicketNo:       "TCK-001",
		CustomerID:     &custID,
		Category:       s03domain.CategoryWrongBilling,
		Description:    "Disputed multiplier",
		AssignedUnitID: "UNIT-COMM-01",
		Status:         s03domain.StatusLogged,
	}
	ctx := context.Background()
	_ = ticketRepo.CreateTicket(ctx, ticket)

	var adjID string

	t.Run("propose adjustment success", func(t *testing.T) {
		dto := transport.ProposeAdjustmentDTO{
			TicketID:          "TICK-001",
			CustomerID:        "CUST-001",
			BillingOfficerID:  "STAFF-01",
			OldBalance:        20000.0,
			AdjustedBalance:   15000.0,
			AdjustmentAmount:  5000.0,
			RootCauseCategory: domain.RootCauseMultiplierError,
			RootCauseNote:     "Fixed incorrect multiplier value",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/adjustments", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
		}
		var created domain.BillingAdjustment
		_ = json.NewDecoder(rec.Body).Decode(&created)
		adjID = created.ID
		if created.AdjustmentAmount != 5000.0 {
			t.Fatalf("expected 5000, got %f", created.AdjustmentAmount)
		}
	})

	t.Run("propose adjustment validation error", func(t *testing.T) {
		dto := transport.ProposeAdjustmentDTO{
			TicketID:   "TICK-001",
			CustomerID: "CUST-001",
			OldBalance: 20000.0,
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/adjustments", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
		}
	})

	t.Run("get adjustment by id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/billing/adjustments/"+adjID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}

		// Not found
		nfReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/adjustments/non-existent", nil)
		nfReq.Header.Set("Authorization", "Bearer "+token)
		nfRec := httptest.NewRecorder()
		mux.ServeHTTP(nfRec, nfReq)
		if nfRec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got %d", nfRec.Code)
		}

		// Empty ID
		emptyReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/adjustments/%20", nil)
		emptyReq.Header.Set("Authorization", "Bearer "+token)
		emptyRec := httptest.NewRecorder()
		mux.ServeHTTP(emptyRec, emptyReq)
		if emptyRec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", emptyRec.Code)
		}
	})

	t.Run("list adjustments", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/billing/adjustments?customer_id=CUST-001", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}
	})

	t.Run("decide adjustment approval", func(t *testing.T) {
		dto := transport.DecideAdjustmentDTO{
			HeadCommercialOpsID: "HEAD-COMM-01",
			Approve:             true,
			Notes:               "Verified against substation log sheets",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/adjustments/"+adjID+"/decide", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+headToken)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
		var approved domain.BillingAdjustment
		_ = json.NewDecoder(rec.Body).Decode(&approved)
		if approved.Status != domain.AdjustmentAppliedRevertedToCare {
			t.Fatalf("expected APPLIED_REVERTED_TO_CARE, got %s", approved.Status)
		}

		// Verify ticket reverted to ATTENDED_TO
		tkt, _ := ticketRepo.GetTicketByID(context.Background(), "TICK-001")
		if tkt.Status != s03domain.StatusAttendedTo {
			t.Fatalf("expected ticket status ATTENDED_TO, got %s", tkt.Status)
		}

		// Trying to decide again on already approved adjustment should fail with 409
		reqRepeat := httptest.NewRequest(http.MethodPost, "/api/v1/billing/adjustments/"+adjID+"/decide", bytes.NewBuffer(body))
		reqRepeat.Header.Set("Authorization", "Bearer "+headToken)
		recRepeat := httptest.NewRecorder()
		mux.ServeHTTP(recRepeat, reqRepeat)
		if recRepeat.Code != http.StatusConflict {
			t.Fatalf("expected 409 Conflict on repeated decision, got %d", recRepeat.Code)
		}
	})

	t.Run("decide adjustment not found", func(t *testing.T) {
		dto := transport.DecideAdjustmentDTO{
			HeadCommercialOpsID: "HEAD-COMM-01",
			Approve:             false,
			Notes:               "Rejected",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/adjustments/non-existent/decide", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+headToken)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got %d", rec.Code)
		}
	})

	t.Run("decide adjustment bad json and validation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/adjustments/adj-123/decide", bytes.NewBufferString("{bad-json}"))
		req.Header.Set("Authorization", "Bearer "+headToken)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
		}

		valReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/adjustments/adj-123/decide", bytes.NewBufferString(`{"head_commercial_ops_id":""}`))
		valReq.Header.Set("Authorization", "Bearer "+headToken)
		valRec := httptest.NewRecorder()
		mux.ServeHTTP(valRec, valReq)
		if valRec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request on validation error, got %d", valRec.Code)
		}
	})
}

func TestHTTPHandler_RevenueProtectionAI(t *testing.T) {
	mux, _, _, _ := setupTestServer(t)
	token := generateTestToken(t, "REV-PROT-01", "REVENUE_PROTECTION_OFFICER")

	var alertID string

	t.Run("scan triggers bypass alert", func(t *testing.T) {
		dto := transport.RevenueScanDTO{
			CustomerOrDTID:  "DT-RAYFIELD-01",
			EntityType:      "DISTRIBUTION_TRANSFORMER",
			DTBulkKWh:       5000.0,
			SumMeteredKWh:   3000.0, // 40% loss
			FeederEnergized: true,
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/revenue-protection/scan", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
		}
		var alert domain.RevenueProtectionAlert
		_ = json.NewDecoder(rec.Body).Decode(&alert)
		alertID = alert.ID
		if alert.AlertType != domain.AlertEnergyTheftBypass {
			t.Fatalf("expected AlertEnergyTheftBypass, got %s", alert.AlertType)
		}
	})

	t.Run("scan clean grid returns 200 with clean status", func(t *testing.T) {
		dto := transport.RevenueScanDTO{
			CustomerOrDTID:  "DT-CLEAN-01",
			EntityType:      "DISTRIBUTION_TRANSFORMER",
			DTBulkKWh:       5000.0,
			SumMeteredKWh:   4800.0, // 4% technical loss
			FeederEnergized: true,
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/revenue-protection/scan", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("review AI alert", func(t *testing.T) {
		unit := "UNIT-SURV-01"
		notes := "Investigated meter seals: intact"
		dto := transport.ReviewAIAlertDTO{
			AssignedUnitID:     &unit,
			InvestigationNotes: &notes,
			Status:             domain.AlertResolvedRecovered,
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/revenue-protection/alerts/"+alertID+"/review", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}

		// GET Alert
		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/revenue-protection/alerts/"+alertID, nil)
		getReq.Header.Set("Authorization", "Bearer "+token)
		getRec := httptest.NewRecorder()
		mux.ServeHTTP(getRec, getReq)
		if getRec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", getRec.Code)
		}

		// List Alerts
		listReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/revenue-protection/alerts?status=RESOLVED_REVENUE_RECOVERED", nil)
		listReq.Header.Set("Authorization", "Bearer "+token)
		listRec := httptest.NewRecorder()
		mux.ServeHTTP(listRec, listReq)
		if listRec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", listRec.Code)
		}

		// GET Alert not found
		nfReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/revenue-protection/alerts/unknown-id", nil)
		nfReq.Header.Set("Authorization", "Bearer "+token)
		nfRec := httptest.NewRecorder()
		mux.ServeHTTP(nfRec, nfReq)
		if nfRec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got %d", nfRec.Code)
		}

		// GET Alert empty ID
		emptyReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/revenue-protection/alerts/%20", nil)
		emptyReq.Header.Set("Authorization", "Bearer "+token)
		emptyRec := httptest.NewRecorder()
		mux.ServeHTTP(emptyRec, emptyReq)
		if emptyRec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request on empty id, got %d", emptyRec.Code)
		}
	})

	t.Run("scan bad json and validation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/revenue-protection/scan", bytes.NewBufferString("{bad-json}"))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
		}

		valReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/revenue-protection/scan", bytes.NewBufferString(`{"customer_or_dt_id":""}`))
		valReq.Header.Set("Authorization", "Bearer "+token)
		valRec := httptest.NewRecorder()
		mux.ServeHTTP(valRec, valReq)
		if valRec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", valRec.Code)
		}
	})

	t.Run("review AI alert bad json and not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/revenue-protection/alerts/alr-1/review", bytes.NewBufferString("{bad-json}"))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
		}

		nfReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/revenue-protection/alerts/unknown-id/review", bytes.NewBufferString(`{"status":"OPEN"}`))
		nfReq.Header.Set("Authorization", "Bearer "+token)
		nfRec := httptest.NewRecorder()
		mux.ServeHTTP(nfRec, nfReq)
		if nfRec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got %d", nfRec.Code)
		}
	})
}
