// Package transport implements wrong billing and revenue protection AI HTTP handlers for S08.
package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/service"
)

// HandleProposeAdjustment submits a proposed billing adjustment.
func (h *HTTPHandler) HandleProposeAdjustment(w http.ResponseWriter, r *http.Request) {
	var dto ProposeAdjustmentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	adj, err := h.svc.ProposeBillingAdjustment(r.Context(), domain.BillingAdjustment{
		TicketID:          dto.TicketID,
		CustomerID:        dto.CustomerID,
		BillingOfficerID:  dto.BillingOfficerID,
		OldBalance:        dto.OldBalance,
		AdjustedBalance:   dto.AdjustedBalance,
		AdjustmentAmount:  dto.AdjustmentAmount,
		RootCauseCategory: dto.RootCauseCategory,
		RootCauseNote:     dto.RootCauseNote,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidAdjustmentParams) {
			h.jsonError(w, http.StatusBadRequest, "invalid_parameters", err.Error())
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, adj)
}

// HandleDecideAdjustment applies Head of Commercial Operations approval or rejection.
func (h *HTTPHandler) HandleDecideAdjustment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "missing_param", "Adjustment id is required")
		return
	}

	var dto DecideAdjustmentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	adj, err := h.svc.ApproveBillingAdjustment(r.Context(), id, dto.HeadCommercialOpsID, dto.Approve, dto.Notes)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "Billing adjustment not found")
			return
		}
		if errors.Is(err, domain.ErrIllegalAdjustmentTransition) {
			h.jsonError(w, http.StatusConflict, "illegal_transition", err.Error())
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, adj)
}

// HandleGetAdjustment retrieves an adjustment by ID.
func (h *HTTPHandler) HandleGetAdjustment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "missing_param", "Adjustment id is required")
		return
	}

	adj, err := h.svc.GetBillingAdjustment(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "Adjustment not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, adj)
}

// HandleListAdjustments lists adjustments by status and customer ID.
func (h *HTTPHandler) HandleListAdjustments(w http.ResponseWriter, r *http.Request) {
	status := domain.AdjustmentStatus(r.URL.Query().Get("status"))
	customerID := r.URL.Query().Get("customer_id")

	adjustments, err := h.svc.ListBillingAdjustments(r.Context(), status, customerID)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if adjustments == nil {
		adjustments = []domain.BillingAdjustment{}
	}
	h.jsonResponse(w, http.StatusOK, adjustments)
}

// HandleRevenueScan triggers an AI scan for energy bypass or churn risks.
func (h *HTTPHandler) HandleRevenueScan(w http.ResponseWriter, r *http.Request) {
	var dto RevenueScanDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	alert, err := h.svc.RunRevenueProtectionScan(r.Context(), service.RevenueScanRequest{
		CustomerOrDTID:  dto.CustomerOrDTID,
		EntityType:      dto.EntityType,
		DTBulkKWh:       dto.DTBulkKWh,
		SumMeteredKWh:   dto.SumMeteredKWh,
		VendingDropPct:  dto.VendingDropPct,
		FeederEnergized: dto.FeederEnergized,
		OutageCount:     dto.OutageCount,
		PurchaseDropPct: dto.PurchaseDropPct,
	})
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	if alert == nil {
		h.jsonResponse(w, http.StatusOK, map[string]interface{}{
			"status":  "clean",
			"message": "No energy theft or churn anomalies detected across grid and customer telemetry.",
		})
		return
	}

	h.jsonResponse(w, http.StatusCreated, alert)
}

// HandleReviewAIAlert updates investigation status and notes on an AI alert.
func (h *HTTPHandler) HandleReviewAIAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "missing_param", "Alert id is required")
		return
	}

	var dto ReviewAIAlertDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	alert, err := h.svc.ReviewAIAlert(r.Context(), id, dto.AssignedUnitID, dto.InvestigationNotes, dto.Status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "Revenue protection alert not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, alert)
}

// HandleGetAIAlert retrieves an AI alert by ID.
func (h *HTTPHandler) HandleGetAIAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "missing_param", "Alert id is required")
		return
	}

	alert, err := h.svc.GetAIAlert(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "Revenue protection alert not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, alert)
}

// HandleListAIAlerts lists AI alerts by status and alert type.
func (h *HTTPHandler) HandleListAIAlerts(w http.ResponseWriter, r *http.Request) {
	status := domain.AIAlertStatus(r.URL.Query().Get("status"))
	alertType := domain.AIAlertType(r.URL.Query().Get("alert_type"))

	alerts, err := h.svc.ListAIAlerts(r.Context(), status, alertType)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if alerts == nil {
		alerts = []domain.RevenueProtectionAlert{}
	}
	h.jsonResponse(w, http.StatusOK, alerts)
}

