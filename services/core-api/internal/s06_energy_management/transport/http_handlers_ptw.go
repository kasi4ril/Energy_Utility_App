// Package transport implements REST routing and handlers for the 5-step Safety Permit-to-Work protocol.
package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

func (h *HTTPHandler) registerPTWRoutes(mux *http.ServeMux, authMW func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/energy/ptw/request", authMW(http.HandlerFunc(h.HandleInitiatePTW)))
	mux.Handle("POST /api/v1/energy/ptw/{id}/approve", authMW(http.HandlerFunc(h.HandleApprovePTW)))
	mux.Handle("POST /api/v1/energy/ptw/{id}/isolate", authMW(http.HandlerFunc(h.HandleIsolatePTW)))
	mux.Handle("POST /api/v1/energy/ptw/{id}/clear-safety", authMW(http.HandlerFunc(h.HandleClearSafetyPTW)))
	mux.Handle("POST /api/v1/energy/ptw/{id}/restore", authMW(http.HandlerFunc(h.HandleRestorePTW)))
	mux.Handle("GET /api/v1/energy/ptw/{id}", authMW(http.HandlerFunc(h.HandleGetPTW)))
	mux.Handle("GET /api/v1/energy/ptw", authMW(http.HandlerFunc(h.HandleListPTWs)))
}

// HandleInitiatePTW (Step 1) receives isolation request from a Unit Head.
func (h *HTTPHandler) HandleInitiatePTW(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto InitiatePTWDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	req := service.InitiatePTWRequest{
		AssetID:         dto.AssetID,
		WorkOrderID:     dto.WorkOrderID,
		UnitHeadID:      claims.StaffID,
		WorkDescription: dto.WorkDescription,
	}

	ptw, err := h.svc.InitiatePTW(r.Context(), req)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "ptw_initiation_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, ptw)
}

// HandleApprovePTW (Step 2) records digital approval from Head of Services or HSE Safety Rep.
func (h *HTTPHandler) HandleApprovePTW(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	permitID := r.PathValue("id")
	if strings.TrimSpace(permitID) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "permit id is required")
		return
	}

	var dto ApprovePTWDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	ptw, err := h.svc.DualApprovePTW(r.Context(), permitID, dto.ApproverRole, claims.StaffID)
	if err != nil {
		if err == service.ErrUnauthorizedApprover {
			h.jsonError(w, http.StatusForbidden, "forbidden", err.Error())
			return
		}
		h.jsonError(w, http.StatusUnprocessableEntity, "approval_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, ptw)
}

// HandleIsolatePTW (Step 3) records Substation Operator breaker rack out and line earthing.
func (h *HTTPHandler) HandleIsolatePTW(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	permitID := r.PathValue("id")
	if strings.TrimSpace(permitID) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "permit id is required")
		return
	}

	var dto IsolatePTWDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	ptw, err := h.svc.GrantSubstationIsolation(r.Context(), permitID, claims.StaffID, dto.IsolationNotes)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "isolation_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, ptw)
}

// HandleClearSafetyPTW (Step 5a) records HSE Safety Rep declaration that crew is clear.
func (h *HTTPHandler) HandleClearSafetyPTW(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	permitID := r.PathValue("id")
	if strings.TrimSpace(permitID) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "permit id is required")
		return
	}

	var dto ClearSafetyPTWDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	ptw, err := h.svc.SignSafetyClearance(r.Context(), permitID, claims.StaffID, dto.Notes)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "clearance_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, ptw)
}

// HandleRestorePTW (Step 5b) records Substation Operator closing breaker and energizing line.
func (h *HTTPHandler) HandleRestorePTW(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	permitID := r.PathValue("id")
	if strings.TrimSpace(permitID) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "permit id is required")
		return
	}

	ptw, err := h.svc.RestoreAndClosePTW(r.Context(), permitID, claims.StaffID)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "restoration_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, ptw)
}

// HandleGetPTW retrieves a permit by UUID or PermitNo.
func (h *HTTPHandler) HandleGetPTW(w http.ResponseWriter, r *http.Request) {
	permitID := r.PathValue("id")
	ptw, err := h.svc.GetPTWByID(r.Context(), permitID)
	if err != nil {
		// Fallback to query by permit number
		ptw, err = h.svc.GetPTWByPermitNo(r.Context(), permitID)
		if err != nil {
			if err == repository.ErrNotFound {
				h.jsonError(w, http.StatusNotFound, "not_found", "permit not found")
				return
			}
			h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
			return
		}
	}
	h.jsonResponse(w, http.StatusOK, ptw)
}

// HandleListPTWs queries permits matching status or asset filters.
func (h *HTTPHandler) HandleListPTWs(w http.ResponseWriter, r *http.Request) {
	status := domain.PTWStatus(r.URL.Query().Get("status"))
	assetID := r.URL.Query().Get("asset_id")

	permits, err := h.svc.ListPTWs(r.Context(), status, assetID)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if permits == nil {
		permits = []domain.SafetyPermitPTW{}
	}
	h.jsonResponse(w, http.StatusOK, permits)
}

