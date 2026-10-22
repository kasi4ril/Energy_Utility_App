// Package transport implements REST routing and JSON handlers for S09 Executive Command.
package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler exposes REST endpoints for C-Suite dashboard aggregation and decision gates.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTPHandler instance for S09.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers all S09 HTTP routes onto the ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// C-Suite Live Intelligence Dashboard
	mux.Handle("GET /api/v1/executive/dashboard", authMW(http.HandlerFunc(h.HandleGetDashboardSummary)))

	// Operational Telemetry Snapshots
	mux.Handle("POST /api/v1/executive/snapshots", authMW(http.HandlerFunc(h.HandleCaptureSnapshot)))
	mux.Handle("GET /api/v1/executive/snapshots/latest", authMW(http.HandlerFunc(h.HandleGetLatestSnapshot)))
	mux.Handle("GET /api/v1/executive/snapshots", authMW(http.HandlerFunc(h.HandleListSnapshots)))

	// Executive Action Console (One-Click Digital Sign-Off Queue)
	mux.Handle("POST /api/v1/executive/actions", authMW(http.HandlerFunc(h.HandleEnqueueActionItem)))
	mux.Handle("POST /api/v1/executive/actions/{id}/decide", authMW(http.HandlerFunc(h.HandleDecideActionItem)))
	mux.Handle("GET /api/v1/executive/actions", authMW(http.HandlerFunc(h.HandleListActionItems)))

	// Predictive Maintenance Prognostics Center
	mux.Handle("POST /api/v1/executive/predictive/hydro", authMW(http.HandlerFunc(h.HandleHydroPrognosis)))
	mux.Handle("POST /api/v1/executive/predictive/generator", authMW(http.HandlerFunc(h.HandleGeneratorPrognosis)))
	mux.Handle("POST /api/v1/executive/predictive/transformer", authMW(http.HandlerFunc(h.HandleTransformerPrognosis)))
	mux.Handle("POST /api/v1/executive/predictive/feeder", authMW(http.HandlerFunc(h.HandleFeederPrognosis)))
	mux.Handle("POST /api/v1/executive/predictive/triggers/{id}/ack", authMW(http.HandlerFunc(h.HandleAcknowledgeTrigger)))
	mux.Handle("GET /api/v1/executive/predictive/triggers", authMW(http.HandlerFunc(h.HandleListPredictiveTriggers)))
}

// HandleGetDashboardSummary returns the aggregated live telemetry for the C-Suite view.
func (h *HTTPHandler) HandleGetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.svc.GetLiveDashboardSummary(r.Context())
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, summary)
}

// HandleCaptureSnapshot records an operational telemetry snapshot.
func (h *HTTPHandler) HandleCaptureSnapshot(w http.ResponseWriter, r *http.Request) {
	var dto CaptureSnapshotDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	snap, err := h.svc.CaptureExecutiveSnapshot(r.Context(), domain.ExecutiveSnapshot{
		SnapshotDate:                  dto.SnapshotDate,
		ShiftWindow:                   dto.ShiftWindow,
		TotalGenerationMW:             dto.TotalGenerationMW,
		TCNGridIntakeMW:               dto.TCNGridIntakeMW,
		TotalSystemLoadMW:             dto.TotalSystemLoadMW,
		FeedersInServiceCount:         dto.FeedersInServiceCount,
		FeedersOutFaultCount:          dto.FeedersOutFaultCount,
		FeedersOutPTWCount:            dto.FeedersOutPTWCount,
		FeedersOutRationingCount:      dto.FeedersOutRationingCount,
		TotalRevenueCollected:         dto.TotalRevenueCollected,
		PrepaidRevenueSharePct:        dto.PrepaidRevenueSharePct,
		UnresolvedTheftAlertsCount:    dto.UnresolvedTheftAlertsCount,
		UnresolvedChurnAlertsCount:    dto.UnresolvedChurnAlertsCount,
		AuditPendingProcurementsCount: dto.AuditPendingProcurementsCount,
		UnreturnedScrapCount:          dto.UnreturnedScrapCount,
		Notes:                         dto.Notes,
		CapturedByStaffID:             dto.CapturedByStaffID,
	})
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, snap)
}

// HandleGetLatestSnapshot retrieves the latest operational snapshot.
func (h *HTTPHandler) HandleGetLatestSnapshot(w http.ResponseWriter, r *http.Request) {
	snap, err := h.svc.GetLatestSnapshot(r.Context())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "No executive snapshot recorded yet")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, snap)
}

// HandleListSnapshots lists historical snapshots.
func (h *HTTPHandler) HandleListSnapshots(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	snaps, err := h.svc.ListSnapshots(r.Context(), limit)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if snaps == nil {
		snaps = []domain.ExecutiveSnapshot{}
	}
	h.jsonResponse(w, http.StatusOK, snaps)
}

// HandleEnqueueActionItem queues an operational approval for C-Suite review.
func (h *HTTPHandler) HandleEnqueueActionItem(w http.ResponseWriter, r *http.Request) {
	var dto EnqueueActionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	item, err := h.svc.EnqueueActionItem(r.Context(), domain.ExecutiveActionItem{
		ActionType:       dto.ActionType,
		TargetRole:       dto.TargetRole,
		SourceEntityID:   dto.SourceEntityID,
		SourceEntityRef:  dto.SourceEntityRef,
		SourceDepartment: dto.SourceDepartment,
		Summary:          dto.Summary,
		Amount:           dto.Amount,
		Priority:         dto.Priority,
	})
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			h.jsonError(w, http.StatusConflict, "conflict", "Action item for this entity reference already exists")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, item)
}

// HandleDecideActionItem records an executive digital approval or rejection.
func (h *HTTPHandler) HandleDecideActionItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "missing_param", "Action item id is required")
		return
	}

	var dto DecideActionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok || claims == nil {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication token required")
		return
	}

	role := domain.ExecutiveRole(claims.PrimaryRole)
	decided, err := h.svc.DecideActionItem(r.Context(), id, dto.DecidedByStaffID, role, dto.Approve, dto.Notes)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "Action item not found")
			return
		}
		if errors.Is(err, domain.ErrUnauthorizedExecutiveRole) {
			h.jsonError(w, http.StatusForbidden, "forbidden", err.Error())
			return
		}
		if errors.Is(err, domain.ErrIllegalActionItemState) {
			h.jsonError(w, http.StatusConflict, "illegal_state", err.Error())
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, decided)
}

// HandleListActionItems lists action items filtered by role and status.
func (h *HTTPHandler) HandleListActionItems(w http.ResponseWriter, r *http.Request) {
	role := domain.ExecutiveRole(r.URL.Query().Get("target_role"))
	status := domain.ActionItemStatus(r.URL.Query().Get("status"))

	items, err := h.svc.ListActionItems(r.Context(), role, status)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if items == nil {
		items = []domain.ExecutiveActionItem{}
	}
	h.jsonResponse(w, http.StatusOK, items)
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
