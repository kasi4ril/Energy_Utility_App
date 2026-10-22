// Package transport implements REST routing and JSON handlers for S06 Energy Management.
package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler exposes REST endpoints for bulk intake, priority rationing, and safety PTWs.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTPHandler for S06.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers all S06 HTTP endpoints onto the ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// Bulk Intake Telemetry
	mux.Handle("POST /api/v1/energy/bulk-intake", authMW(http.HandlerFunc(h.HandleRecordBulkIntake)))
	mux.Handle("GET /api/v1/energy/bulk-intake/latest", authMW(http.HandlerFunc(h.HandleGetLatestBulkIntake)))
	mux.Handle("GET /api/v1/energy/bulk-intake", authMW(http.HandlerFunc(h.HandleListBulkIntake)))

	// Feeder Priority Rationing
	mux.Handle("POST /api/v1/energy/rationing/calculate", authMW(http.HandlerFunc(h.HandleCalculateRationing)))
	mux.Handle("GET /api/v1/energy/feeders", authMW(http.HandlerFunc(h.HandleListFeeders)))
	mux.Handle("GET /api/v1/energy/feeders/{assetId}", authMW(http.HandlerFunc(h.HandleGetFeeder)))

	// Google Forms B Webhook (HMAC-SHA256 Signed)
	mux.HandleFunc("POST /api/v1/energy/webhook/google-forms", h.HandleGoogleFormsWebhook)

	// PTW Protocol Endpoints (registered in http_handlers_ptw.go)
	h.registerPTWRoutes(mux, authMW)
}

// HandleRecordBulkIntake logs bulk power generation & TCN intake.
func (h *HTTPHandler) HandleRecordBulkIntake(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto RecordBulkIntakeDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	log := domain.DispatchTelemetryLog{
		ShiftID:                     dto.ShiftID,
		ShiftWindow:                 dto.ShiftWindow,
		GenerationMWIn:              dto.GenerationMWIn,
		TCNMWIn:                     dto.TCNMWIn,
		SubstationStaffRollcallJSON: dto.SubstationRollcall,
		HandoverNotes:               dto.HandoverNotes,
		OperatorID:                  claims.StaffID,
	}

	recorded, err := h.svc.RecordBulkIntake(r.Context(), log)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "intake_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, recorded)
}

// HandleGetLatestBulkIntake returns the most recently recorded intake record.
func (h *HTTPHandler) HandleGetLatestBulkIntake(w http.ResponseWriter, r *http.Request) {
	latest, err := h.repo.GetLatestDispatchLog(r.Context())
	if err != nil {
		if err == repository.ErrNotFound {
			h.jsonError(w, http.StatusNotFound, "not_found", "no dispatch intake records logged yet")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, latest)
}

// HandleListBulkIntake lists historical intake logs.
func (h *HTTPHandler) HandleListBulkIntake(w http.ResponseWriter, r *http.Request) {
	shiftID := r.URL.Query().Get("shift_id")
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := h.svc.ListDispatchLogs(r.Context(), shiftID, limit)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if logs == nil {
		logs = []domain.DispatchTelemetryLog{}
	}
	h.jsonResponse(w, http.StatusOK, logs)
}

// HandleCalculateRationing calculates priority load shedding decisions.
func (h *HTTPHandler) HandleCalculateRationing(w http.ResponseWriter, r *http.Request) {
	var dto RationingCalculateDTO
	_ = json.NewDecoder(r.Body).Decode(&dto)

	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	decisions, err := h.svc.RunFeederRationing(r.Context(), dto.AvailableBulkMW)
	if err != nil {
		if err == service.ErrNoDispatchData {
			h.jsonError(w, http.StatusPreconditionFailed, "precondition_failed", err.Error())
			return
		}
		h.jsonError(w, http.StatusUnprocessableEntity, "rationing_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, decisions)
}

// HandleListFeeders retrieves current telemetry for all feeders.
func (h *HTTPHandler) HandleListFeeders(w http.ResponseWriter, r *http.Request) {
	feeders, err := h.svc.ListFeederStatuses(r.Context())
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if feeders == nil {
		feeders = []domain.FeederStatus{}
	}
	h.jsonResponse(w, http.StatusOK, feeders)
}

// HandleGetFeeder retrieves telemetry for a single feeder asset.
func (h *HTTPHandler) HandleGetFeeder(w http.ResponseWriter, r *http.Request) {
	assetID := r.PathValue("assetId")
	if strings.TrimSpace(assetID) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "assetId is required")
		return
	}

	feeder, err := h.svc.GetFeederStatus(r.Context(), assetID)
	if err != nil {
		if err == repository.ErrNotFound {
			h.jsonError(w, http.StatusNotFound, "not_found", "feeder asset telemetry not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, feeder)
}

// HandleGoogleFormsWebhook ingests Form B hourly substation & dispatch submissions.
func (h *HTTPHandler) HandleGoogleFormsWebhook(w http.ResponseWriter, r *http.Request) {
	sig := r.Header.Get("X-Signature-SHA256")
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "failed to read webhook body")
		return
	}

	log, err := h.svc.IngestGoogleFormsWebhook(r.Context(), rawBody, sig)
	if err != nil {
		if err == service.ErrInvalidWebhookSignature {
			h.jsonError(w, http.StatusUnauthorized, "unauthorized", err.Error())
			return
		}
		h.jsonError(w, http.StatusUnprocessableEntity, "webhook_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"status":        "INGESTED",
		"log_id":        log.ID,
		"shift_id":      log.ShiftID,
		"total_bulk_mw": log.TotalBulkMW,
	})
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

