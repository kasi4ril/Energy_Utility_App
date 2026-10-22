// Package transport implements REST routing and JSON handlers for S04 Services & Maintenance.
package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler exposes REST endpoints for field maintenance and PWA proof-of-work workflows.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTPHandler instance for S04.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers S04 maintenance routes onto the ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// Work Order Management
	mux.Handle("POST /api/v1/services/work-orders", authMW(http.HandlerFunc(h.HandleCreateWorkOrder)))
	mux.Handle("GET /api/v1/services/work-orders", authMW(http.HandlerFunc(h.HandleListWorkOrders)))
	mux.Handle("GET /api/v1/services/work-orders/{id}", authMW(http.HandlerFunc(h.HandleGetWorkOrder)))

	// PWA Field Proof of Work
	mux.Handle("POST /api/v1/services/work-orders/{id}/proof", authMW(http.HandlerFunc(h.HandleSubmitProof)))
	mux.Handle("GET /api/v1/services/work-orders/{id}/proof", authMW(http.HandlerFunc(h.HandleGetProof)))

	// Unit Head Verification & Customer Care Reversion
	mux.Handle("POST /api/v1/services/work-orders/{id}/verify", authMW(http.HandlerFunc(h.HandleVerifyProof)))
}

// HandleCreateWorkOrder handles work order dispatch from the Unit Head.
func (h *HTTPHandler) HandleCreateWorkOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto CreateWorkOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	wo, err := h.svc.CreateWorkOrder(
		r.Context(), dto.TicketID, dto.UnitType, dto.WorkOrderType, dto.Priority,
		claims.StaffID, dto.TeamLeadID, dto.AssetID,
		dto.PTWRequired, dto.PTWPermitID, dto.MaterialsRequisitionRef, dto.Instructions,
	)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "dispatch_failed", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, wo)
}

// HandleListWorkOrders retrieves work orders filtered by unit, status, and team lead.
func (h *HTTPHandler) HandleListWorkOrders(w http.ResponseWriter, r *http.Request) {
	unit := domain.UnitType(r.URL.Query().Get("unit"))
	status := domain.WorkOrderStatus(r.URL.Query().Get("status"))
	teamLeadID := r.URL.Query().Get("team_lead_id")

	orders, err := h.repo.ListWorkOrders(r.Context(), unit, status, teamLeadID)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if orders == nil {
		orders = []domain.MaintenanceWorkOrder{}
	}
	h.jsonResponse(w, http.StatusOK, orders)
}

// HandleGetWorkOrder retrieves a single work order by UUID.
func (h *HTTPHandler) HandleGetWorkOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "work order id is required")
		return
	}

	wo, err := h.repo.GetWorkOrderByID(r.Context(), id)
	if err != nil {
		h.jsonError(w, http.StatusNotFound, "not_found", "work order not found")
		return
	}
	h.jsonResponse(w, http.StatusOK, wo)
}

// HandleSubmitProof processes mobile PWA photo and GPS evidence upload from the field crew.
func (h *HTTPHandler) HandleSubmitProof(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "work order id is required")
		return
	}

	var dto SubmitProofDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	capturedGPS := s01domain.GPSPoint{Latitude: dto.Latitude, Longitude: dto.Longitude}
	proof, err := h.svc.SubmitProofOfWork(
		r.Context(), id, claims.StaffID,
		capturedGPS, dto.BeforePhotoURL, dto.AfterPhotoURL,
		dto.PhotoHashBefore, dto.PhotoHashAfter, dto.MeterSealNumber,
		dto.WorkDetails,
	)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "proof_rejected", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, proof)
}

// HandleGetProof retrieves the proof-of-work record for a given work order.
func (h *HTTPHandler) HandleGetProof(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "work order id is required")
		return
	}

	proof, err := h.repo.GetProofByWorkOrderID(r.Context(), id)
	if err != nil {
		h.jsonError(w, http.StatusNotFound, "not_found", "proof not found for this work order")
		return
	}
	h.jsonResponse(w, http.StatusOK, proof)
}

// HandleVerifyProof executes Unit Head approval or rework assignment.
func (h *HTTPHandler) HandleVerifyProof(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "work order id is required")
		return
	}

	var dto VerifyProofDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	wo, err := h.svc.VerifyProofAndRevert(r.Context(), id, claims.StaffID, dto.Approved, dto.Notes)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "verification_failed", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, wo)
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

