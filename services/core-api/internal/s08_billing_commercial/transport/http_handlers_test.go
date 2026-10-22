// Package transport_test exercises REST endpoints for vending and refund ledgers for S08.
package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/service"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/transport"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const testJWTKey = "nesco-billing-jwt-key-32b-secret-test!"

func generateTestToken(t *testing.T, staffID, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(auth.Claims{
		StaffID:     staffID,
		StaffNumber: "NES-" + staffID,
		FullName:    "Officer " + staffID,
		PrimaryRole: role,
		OrgUnitID:   "unit-commercial-hq",
		IssuedAt:    time.Now().Unix(),
		ExpiresAt:   time.Now().Add(15 * time.Minute).Unix(),
	}, testJWTKey)
	if err != nil {
		t.Fatalf("failed to generate auth token: %v", err)
	}
	return token
}

func setupTestServer(t *testing.T) (*http.ServeMux, *repository.MemoryRepository, *service.WorkflowService, *s03repo.MemoryRepository) {
	t.Helper()
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	ticketRepo := s03repo.NewMemoryRepository()
	svc := service.NewWorkflowService(repo, bus, ticketRepo)
	handler := transport.NewHTTPHandler(svc, repo, testJWTKey)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, repo, svc, ticketRepo
}

func TestHTTPHandler_VendPayment(t *testing.T) {
	mux, _, _, _ := setupTestServer(t)
	token := generateTestToken(t, "BILL-01", "BILLING_OFFICER")

	t.Run("unauthorized without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/vending/vend", bytes.NewBufferString("{}"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("invalid json body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/vending/vend", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
		}
	})

	t.Run("validation failure", func(t *testing.T) {
		payload := `{"customer_id":"","amount_paid":-100}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/vending/vend", bytes.NewBufferString(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
		}
	})

	t.Run("successful vending payment", func(t *testing.T) {
		dto := transport.VendPaymentDTO{
			CustomerID:             "CUST-001",
			AccountNumber:          "ACC-001",
			MeterNumber:            "01234567890",
			PaymentChannelCategory: domain.ChannelOnlinePortal,
			PaymentPlatformSource:  "PAYSTACK",
			AmountPaid:             5000.0,
			TariffRatePerKWh:       50.0,
			SettlementVerified:     true,
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/vending/vend", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
		}

		var created domain.VendingTransaction
		if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if created.KWhVended != 100.0 {
			t.Fatalf("expected 100 kWh, got %f", created.KWhVended)
		}
		if len(created.STSToken) != 24 {
			t.Fatalf("expected STS token length 24, got %s", created.STSToken)
		}

		// GET by ID
		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/vending/"+created.ID, nil)
		getReq.Header.Set("Authorization", "Bearer "+token)
		getRec := httptest.NewRecorder()
		mux.ServeHTTP(getRec, getReq)
		if getRec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", getRec.Code)
		}

		// GET by ID - not found
		nfReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/vending/non-existent-id", nil)
		nfReq.Header.Set("Authorization", "Bearer "+token)
		nfRec := httptest.NewRecorder()
		mux.ServeHTTP(nfRec, nfReq)
		if nfRec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got %d", nfRec.Code)
		}

		// GET by ID - empty ID
		emptyReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/vending/%20", nil)
		emptyReq.Header.Set("Authorization", "Bearer "+token)
		emptyRec := httptest.NewRecorder()
		mux.ServeHTTP(emptyRec, emptyReq)
		if emptyRec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", emptyRec.Code)
		}

		// LIST
		listReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/vending?customer_id=CUST-001&limit=10", nil)
		listReq.Header.Set("Authorization", "Bearer "+token)
		listRec := httptest.NewRecorder()
		mux.ServeHTTP(listRec, listReq)
		if listRec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", listRec.Code)
		}
		var list []domain.VendingTransaction
		_ = json.NewDecoder(listRec.Body).Decode(&list)
		if len(list) != 1 {
			t.Fatalf("expected 1 tx in list, got %d", len(list))
		}
	})
}

func TestHTTPHandler_RefundLedger(t *testing.T) {
	mux, _, _, _ := setupTestServer(t)
	token := generateTestToken(t, "BILL-02", "BILLING_SUPERVISOR")

	t.Run("register refund ledger successfully", func(t *testing.T) {
		dto := transport.RegisterRefundLedgerDTO{
			CustomerID:            "CUST-002",
			MeterSerial:           "MTR-002",
			InitialMeterCost:      35000.0,
			RefundPercentagePerTx: 10.0,
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/refunds", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
		}

		// GET by Customer
		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/refunds/CUST-002", nil)
		getReq.Header.Set("Authorization", "Bearer "+token)
		getRec := httptest.NewRecorder()
		mux.ServeHTTP(getRec, getReq)
		if getRec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", getRec.Code)
		}

		// Conflict on duplicate register
		reqDup := httptest.NewRequest(http.MethodPost, "/api/v1/billing/refunds", bytes.NewBuffer(body))
		reqDup.Header.Set("Authorization", "Bearer "+token)
		recDup := httptest.NewRecorder()
		mux.ServeHTTP(recDup, reqDup)
		if recDup.Code != http.StatusConflict {
			t.Fatalf("expected 409 Conflict, got %d", recDup.Code)
		}

		// Not found for unknown customer
		nfReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/refunds/UNKNOWN", nil)
		nfReq.Header.Set("Authorization", "Bearer "+token)
		nfRec := httptest.NewRecorder()
		mux.ServeHTTP(nfRec, nfReq)
		if nfRec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got %d", nfRec.Code)
		}

		// Missing customer ID in path
		emptyReq := httptest.NewRequest(http.MethodGet, "/api/v1/billing/refunds/%20", nil)
		emptyReq.Header.Set("Authorization", "Bearer "+token)
		emptyRec := httptest.NewRecorder()
		mux.ServeHTTP(emptyRec, emptyReq)
		if emptyRec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", emptyRec.Code)
		}
	})

	t.Run("invalid json body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/refunds", bytes.NewBufferString("{bad-json}"))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		dto := transport.RegisterRefundLedgerDTO{
			CustomerID:       "",
			MeterSerial:      "MTR-003",
			InitialMeterCost: -500.0,
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/refunds", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
		}
	})
}
