// Package transport implements REST routing and JSON handlers for staff allowance claims and audit review.
package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

func (h *HTTPHandler) registerClaimsRoutes(mux *http.ServeMux, authMW func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/audit/claims", authMW(http.HandlerFunc(h.HandleSubmitClaim)))
	mux.Handle("GET /api/v1/audit/claims", authMW(http.HandlerFunc(h.HandleListClaims)))
	mux.Handle("GET /api/v1/audit/claims/{id}", authMW(http.HandlerFunc(h.HandleGetClaim)))
	mux.Handle("POST /api/v1/audit/claims/{id}/review", authMW(http.HandlerFunc(h.HandleReviewClaim)))
}

// HandleSubmitClaim records an allowance claim and initiates 3-way verification.
func (h *HTTPHandler) HandleSubmitClaim(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto SubmitClaimDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	claim := domain.StaffAllowanceClaim{
		StaffID:             claims.StaffID,
		ClaimType:           dto.ClaimType,
		Amount:              dto.Amount,
		ClaimDate:           dto.ClaimDate,
		ShiftWindow:         dto.ShiftWindow,
		LinkedPTWOrTicketID: dto.LinkedPTWOrTicketID,
		RosterMatchVerified: dto.RosterMatchVerified,
		PTWCrewVerified:     dto.PTWCrewVerified,
		GPSMatchVerified:    dto.GPSMatchVerified,
	}

	created, err := h.svc.SubmitAllowanceClaim(r.Context(), claim)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "claim_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, created)
}

// HandleReviewClaim records an Internal Auditor's formal approval or rejection on a claim.
func (h *HTTPHandler) HandleReviewClaim(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	claimID := r.PathValue("id")
	if strings.TrimSpace(claimID) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "claim id is required")
		return
	}

	var dto ReviewClaimDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	decided, err := h.svc.ReviewClaimByAuditor(r.Context(), claimID, claims.StaffID, dto.Approve, dto.Notes)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "review_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, decided)
}

// HandleGetClaim retrieves an allowance claim by UUID.
func (h *HTTPHandler) HandleGetClaim(w http.ResponseWriter, r *http.Request) {
	claimID := r.PathValue("id")
	claim, err := h.svc.GetClaimByID(r.Context(), claimID)
	if err != nil {
		if err == repository.ErrNotFound {
			h.jsonError(w, http.StatusNotFound, "not_found", "claim not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, claim)
}

// HandleListClaims queries allowance claims filtered by staff ID or status.
func (h *HTTPHandler) HandleListClaims(w http.ResponseWriter, r *http.Request) {
	staffID := r.URL.Query().Get("staff_id")
	status := domain.ClaimStatus(r.URL.Query().Get("status"))

	claimsList, err := h.svc.ListClaims(r.Context(), staffID, status)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if claimsList == nil {
		claimsList = []domain.StaffAllowanceClaim{}
	}
	h.jsonResponse(w, http.StatusOK, claimsList)
}

