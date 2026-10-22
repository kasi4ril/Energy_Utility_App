// Package transport implements REST routing and JSON handlers for reverse logistics scrap returns.
package transport

import (
	"encoding/json"
	"net/http"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HandleLogScrapReturn logs returned scrap into scrap yard.
func (h *HTTPHandler) HandleLogScrapReturn(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto LogScrapReturnDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	scrap := domain.ScrapReturn{
		RequisitionID:       dto.RequisitionID,
		TeamLeadID:          claims.StaffID,
		StoreOfficerID:      claims.StaffID,
		ItemsReturnedJSON:   dto.ItemsReturnedJSON,
		ScrapConditionNotes: dto.ScrapConditionNotes,
	}

	logged, err := h.svc.LogScrapReturn(r.Context(), scrap)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "scrap_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, logged)
}

// HandleListScrapReturns lists scrap records.
func (h *HTTPHandler) HandleListScrapReturns(w http.ResponseWriter, r *http.Request) {
	teamLeadID := r.URL.Query().Get("team_lead_id")
	returns, err := h.repo.ListScrapReturns(r.Context(), teamLeadID)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if returns == nil {
		returns = []domain.ScrapReturn{}
	}
	h.jsonResponse(w, http.StatusOK, returns)
}

