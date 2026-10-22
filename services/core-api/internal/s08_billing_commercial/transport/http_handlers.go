// Package transport implements REST routing and JSON handlers for S08 Billing & Commercial Operations.
package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler exposes REST endpoints for billing, vending, refunds, and revenue protection.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTPHandler instance for S08.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers all S08 HTTP routes onto the ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// Vending & Payment Transactions (STS Token Generator)
	mux.Handle("POST /api/v1/billing/vending/vend", authMW(http.HandlerFunc(h.HandleVendPayment)))
	mux.Handle("GET /api/v1/billing/vending/{id}", authMW(http.HandlerFunc(h.HandleGetVendingTransaction)))
	mux.Handle("GET /api/v1/billing/vending", authMW(http.HandlerFunc(h.HandleListVendingTransactions)))

	// Customer Meter-Purchase Refund Ledgers
	mux.Handle("POST /api/v1/billing/refunds", authMW(http.HandlerFunc(h.HandleRegisterRefundLedger)))
	mux.Handle("GET /api/v1/billing/refunds/{customer_id}", authMW(http.HandlerFunc(h.HandleGetRefundLedgerByCustomer)))

	// Wrong Billing Adjustments & Commercial Head Gate
	mux.Handle("POST /api/v1/billing/adjustments", authMW(http.HandlerFunc(h.HandleProposeAdjustment)))
	mux.Handle("POST /api/v1/billing/adjustments/{id}/decide", authMW(http.HandlerFunc(h.HandleDecideAdjustment)))
	mux.Handle("GET /api/v1/billing/adjustments/{id}", authMW(http.HandlerFunc(h.HandleGetAdjustment)))
	mux.Handle("GET /api/v1/billing/adjustments", authMW(http.HandlerFunc(h.HandleListAdjustments)))

	// Revenue Protection AI Watchdog & Anomaly Scans
	mux.Handle("POST /api/v1/billing/revenue-protection/scan", authMW(http.HandlerFunc(h.HandleRevenueScan)))
	mux.Handle("POST /api/v1/billing/revenue-protection/alerts/{id}/review", authMW(http.HandlerFunc(h.HandleReviewAIAlert)))
	mux.Handle("GET /api/v1/billing/revenue-protection/alerts/{id}", authMW(http.HandlerFunc(h.HandleGetAIAlert)))
	mux.Handle("GET /api/v1/billing/revenue-protection/alerts", authMW(http.HandlerFunc(h.HandleListAIAlerts)))
}

// HandleVendPayment processes a prepaid or postpaid vending payment.
func (h *HTTPHandler) HandleVendPayment(w http.ResponseWriter, r *http.Request) {
	var dto VendPaymentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	tx, err := h.svc.RecordVendingPayment(r.Context(), service.RecordPaymentRequest{
		TxReference:            dto.TxReference,
		CustomerID:             dto.CustomerID,
		AccountNumber:          dto.AccountNumber,
		MeterNumber:            dto.MeterNumber,
		FeederID:               dto.FeederID,
		DTID:                   dto.DTID,
		PaymentChannelCategory: dto.PaymentChannelCategory,
		PaymentPlatformSource:  dto.PaymentPlatformSource,
		AmountPaid:             dto.AmountPaid,
		TariffRatePerKWh:       dto.TariffRatePerKWh,
		SettlementVerified:     dto.SettlementVerified,
		GatewayPayloadJSON:     dto.GatewayPayloadJSON,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidVendingParams) {
			h.jsonError(w, http.StatusBadRequest, "invalid_parameters", err.Error())
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, tx)
}

// HandleGetVendingTransaction retrieves a transaction by ID.
func (h *HTTPHandler) HandleGetVendingTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "missing_param", "Transaction id is required")
		return
	}

	tx, err := h.svc.GetVendingTransaction(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "Transaction not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, tx)
}

// HandleListVendingTransactions queries transactions matching filters.
func (h *HTTPHandler) HandleListVendingTransactions(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	meterNumber := r.URL.Query().Get("meter_number")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	txs, err := h.svc.ListVendingTransactions(r.Context(), customerID, meterNumber, limit)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if txs == nil {
		txs = []domain.VendingTransaction{}
	}
	h.jsonResponse(w, http.StatusOK, txs)
}

// HandleRegisterRefundLedger registers a customer meter cost refund ledger.
func (h *HTTPHandler) HandleRegisterRefundLedger(w http.ResponseWriter, r *http.Request) {
	var dto RegisterRefundLedgerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	ledger, err := h.svc.RegisterMeterRefundLedger(r.Context(), domain.MeterRefundLedger{
		CustomerID:            dto.CustomerID,
		MeterSerial:           dto.MeterSerial,
		InitialMeterCost:      dto.InitialMeterCost,
		RefundPercentagePerTx: dto.RefundPercentagePerTx,
	})
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			h.jsonError(w, http.StatusConflict, "conflict", "Meter refund ledger already exists")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, ledger)
}

// HandleGetRefundLedgerByCustomer retrieves the refund ledger of a customer.
func (h *HTTPHandler) HandleGetRefundLedgerByCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")
	if strings.TrimSpace(customerID) == "" {
		h.jsonError(w, http.StatusBadRequest, "missing_param", "customer_id is required")
		return
	}

	ledger, err := h.svc.GetRefundLedgerByCustomer(r.Context(), customerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "Refund ledger not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, ledger)
}

func (h *HTTPHandler) jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *HTTPHandler) jsonError(w http.ResponseWriter, code int, errCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error_code": errCode,
		"message":    message,
	})
}

