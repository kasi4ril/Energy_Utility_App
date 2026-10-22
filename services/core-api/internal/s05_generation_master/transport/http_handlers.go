// Package transport implements REST routing and JSON handlers for S05 Generation Department.
package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler exposes REST endpoints for generation telemetry, roll-calls, and workshop workflows.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTPHandler instance for S05.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers S05 generation routes onto the ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// Telemetry
	mux.Handle("POST /api/v1/generation/telemetry", authMW(http.HandlerFunc(h.HandleRecordTelemetry)))
	mux.Handle("GET /api/v1/generation/telemetry", authMW(http.HandlerFunc(h.HandleListTelemetry)))
	mux.Handle("GET /api/v1/generation/telemetry/latest/{plantId}", authMW(http.HandlerFunc(h.HandleGetLatestTelemetry)))

	// Roll-Call
	mux.Handle("POST /api/v1/generation/rollcall", authMW(http.HandlerFunc(h.HandleRecordRollcall)))
	mux.Handle("GET /api/v1/generation/rollcall", authMW(http.HandlerFunc(h.HandleListRollcalls)))

	// Faults & Workshop
	mux.Handle("POST /api/v1/generation/faults", authMW(http.HandlerFunc(h.HandleReportFault)))
	mux.Handle("GET /api/v1/generation/faults", authMW(http.HandlerFunc(h.HandleListFaults)))
	mux.Handle("POST /api/v1/generation/faults/{id}/gm-approval", authMW(http.HandlerFunc(h.HandleGMApproval)))
	mux.Handle("POST /api/v1/generation/faults/{id}/workshop", authMW(http.HandlerFunc(h.HandleUpdateWorkshop)))

	// Low-Bandwidth Google Forms Webhook Ingestion (HMAC Signed)
	mux.HandleFunc("POST /api/v1/generation/webhook/google-forms", h.HandleGoogleFormsWebhook)
}

// HandleRecordTelemetry logs single plant hourly hydrological and electrical telemetry.
func (h *HTTPHandler) HandleRecordTelemetry(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto RecordTelemetryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	log := domain.MasterStationTelemetryLog{
		ShiftID:             dto.ShiftID,
		ShiftWindow:         dto.ShiftWindow,
		PlantID:             dto.PlantID,
		PlantState:          dto.PlantState,
		OutageReason:        dto.OutageReason,
		OutageDetails:       dto.OutageDetails,
		WaterLevelM:         dto.WaterLevelM,
		HeadwaterElevationM: dto.HeadwaterElevationM,
		TailwaterElevationM: dto.TailwaterElevationM,
		InflowRateM3S:       dto.InflowRateM3S,
		MWGenerated:         dto.MWGenerated,
		MVArReactive:        dto.MVArReactive,
		StatorWindingTempC:  dto.StatorWindingTempC,
		BearingVibrationMMS: dto.BearingVibrationMMS,
		LubeOilPressureBar:  dto.LubeOilPressureBar,
		LoggedByOperatorID:  claims.StaffID,
	}

	alerts, err := h.svc.RecordTelemetry(r.Context(), log)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "telemetry_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"status": "recorded",
		"alerts": alerts,
	})
}

// HandleListTelemetry retrieves telemetry records by shift and plant.
func (h *HTTPHandler) HandleListTelemetry(w http.ResponseWriter, r *http.Request) {
	shiftID := r.URL.Query().Get("shift_id")
	plantID := r.URL.Query().Get("plant_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	logs, err := h.repo.ListTelemetryLogs(r.Context(), shiftID, plantID, limit)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if logs == nil {
		logs = []domain.MasterStationTelemetryLog{}
	}
	h.jsonResponse(w, http.StatusOK, logs)
}

// HandleGetLatestTelemetry retrieves the most recent telemetry entry for a plant.
func (h *HTTPHandler) HandleGetLatestTelemetry(w http.ResponseWriter, r *http.Request) {
	plantID := r.PathValue("plantId")
	if strings.TrimSpace(plantID) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "plantId is required")
		return
	}

	latest, err := h.repo.GetLatestTelemetryByPlant(r.Context(), plantID)
	if err != nil {
		h.jsonError(w, http.StatusNotFound, "not_found", "no telemetry found for this plant")
		return
	}
	h.jsonResponse(w, http.StatusOK, latest)
}

// HandleRecordRollcall logs physical shift roster presence collated from a remote plant.
func (h *HTTPHandler) HandleRecordRollcall(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto RecordRollcallDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	rc := domain.RemotePlantShiftRollcall{
		ShiftID:             dto.ShiftID,
		ShiftWindow:         dto.ShiftWindow,
		PlantID:             dto.PlantID,
		OperatorNames:       dto.OperatorNames,
		ProtectionEngineers: dto.ProtectionEngineers,
		Technicians:         dto.Technicians,
		HandoverNotes:       dto.HandoverNotes,
		MasterOperatorID:    claims.StaffID,
	}

	recorded, err := h.svc.RecordRollcall(r.Context(), rc)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "rollcall_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, recorded)
}

// HandleListRollcalls retrieves shift roll-calls across all 6 plants.
func (h *HTTPHandler) HandleListRollcalls(w http.ResponseWriter, r *http.Request) {
	shiftID := r.URL.Query().Get("shift_id")
	records, err := h.repo.ListRollcallsByShift(r.Context(), shiftID)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if records == nil {
		records = []domain.RemotePlantShiftRollcall{}
	}
	h.jsonResponse(w, http.StatusOK, records)
}

// HandleReportFault logs a new fault and initiates GM material approval or workshop dispatch.
func (h *HTTPHandler) HandleReportFault(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto ReportFaultJobDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	job, err := h.svc.ReportPlantFault(
		r.Context(), dto.PlantID, claims.StaffID,
		dto.Wing, dto.FaultDescription, dto.Severity, dto.RequiresWorkshop,
	)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "fault_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, job)
}

// HandleListFaults returns plant jobs filtered by plant, wing, or status.
func (h *HTTPHandler) HandleListFaults(w http.ResponseWriter, r *http.Request) {
	plantID := r.URL.Query().Get("plant_id")
	wing := domain.TechnicalWing(r.URL.Query().Get("wing"))
	status := domain.PlantJobStatus(r.URL.Query().Get("status"))

	jobs, err := h.repo.ListPlantJobs(r.Context(), plantID, wing, status)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if jobs == nil {
		jobs = []domain.PlantFaultAndWorkshopJob{}
	}
	h.jsonResponse(w, http.StatusOK, jobs)
}

// HandleGMApproval applies the General Manager's approval or rejection decision.
func (h *HTTPHandler) HandleGMApproval(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	id := r.PathValue("id")
	var dto GMApprovalDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	job, err := h.svc.ApproveGMJob(r.Context(), id, claims.StaffID, dto.Approved, dto.Notes)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "approval_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, job)
}

// HandleUpdateWorkshop advances workshop fabrication status.
func (h *HTTPHandler) HandleUpdateWorkshop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var dto UpdateWorkshopDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	job, err := h.svc.UpdateWorkshopFabrication(r.Context(), id, dto.WorkshopStatus, dto.StoreRequisitionRef)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "workshop_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, job)
}

// HandleGoogleFormsWebhook ingests Google Forms webhook with cryptographic HMAC verification.
func (h *HTTPHandler) HandleGoogleFormsWebhook(w http.ResponseWriter, r *http.Request) {
	sig := r.Header.Get("X-Signature-SHA256")
	if sig == "" {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing X-Signature-SHA256 header")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "failed to read webhook body")
		return
	}

	alerts, err := h.svc.IngestGoogleFormsWebhook(r.Context(), body, sig)
	if err != nil {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"status": "ingested",
		"alerts": alerts,
	})
}

func (h *HTTPHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *HTTPHandler) jsonError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   code,
		"message": msg,
	})
}

