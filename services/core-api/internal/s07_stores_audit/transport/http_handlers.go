// Package transport implements REST routing and JSON handlers for S07 Stores & Internal Audit.
package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler exposes REST endpoints for stores inventory, requisitions, and audit.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTPHandler instance for S07.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers all S07 HTTP routes onto the ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// Inventory Ledger (Live @GridAssistant Stock Checks)
	mux.Handle("POST /api/v1/stores/items", authMW(http.HandlerFunc(h.HandleRegisterItem)))
	mux.Handle("GET /api/v1/stores/items", authMW(http.HandlerFunc(h.HandleListItems)))
	mux.Handle("GET /api/v1/stores/items/{code}", authMW(http.HandlerFunc(h.HandleGetItemByCode)))
	mux.Handle("POST /api/v1/stores/items/{id}/stock", authMW(http.HandlerFunc(h.HandleUpdateStock)))

	// Dual-Path Requisitions
	mux.Handle("POST /api/v1/stores/requisitions", authMW(http.HandlerFunc(h.HandleInitiateRequisition)))
	mux.Handle("GET /api/v1/stores/requisitions", authMW(http.HandlerFunc(h.HandleListRequisitions)))
	mux.Handle("GET /api/v1/stores/requisitions/{id}", authMW(http.HandlerFunc(h.HandleGetRequisition)))
	mux.Handle("POST /api/v1/stores/requisitions/{id}/dept-approve", authMW(http.HandlerFunc(h.HandleDeptApproveRequisition)))
	mux.Handle("POST /api/v1/stores/requisitions/{id}/audit-consent", authMW(http.HandlerFunc(h.HandleAuditConsentRequisition)))
	mux.Handle("POST /api/v1/stores/requisitions/{id}/md-approve", authMW(http.HandlerFunc(h.HandleMDApproveRequisition)))
	mux.Handle("POST /api/v1/stores/requisitions/{id}/issue", authMW(http.HandlerFunc(h.HandleIssueMaterialWaybill)))

	// Reverse Logistics Scrap Returns
	mux.Handle("POST /api/v1/stores/scrap-returns", authMW(http.HandlerFunc(h.HandleLogScrapReturn)))
	mux.Handle("GET /api/v1/stores/scrap-returns", authMW(http.HandlerFunc(h.HandleListScrapReturns)))

	// Allowance Claims Desk (Registered in http_handlers_claims.go)
	h.registerClaimsRoutes(mux, authMW)
}

// HandleRegisterItem registers a new SKU in warehouse stores.
func (h *HTTPHandler) HandleRegisterItem(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto CreateStoreItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	item := domain.StoreItem{
		ItemCode:          dto.ItemCode,
		ItemName:          dto.ItemName,
		Category:          dto.Category,
		QuantityAvailable: dto.QuantityAvailable,
		ReorderThreshold:  dto.ReorderThreshold,
		UnitOfMeasure:     dto.UnitOfMeasure,
		UnitCost:          dto.UnitCost,
		BinLocation:       dto.BinLocation,
		LastUpdatedBy:     claims.StaffID,
	}

	created, err := h.svc.RegisterStoreItem(r.Context(), item)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "registration_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, created)
}

// HandleGetItemByCode returns stock balance by item code.
func (h *HTTPHandler) HandleGetItemByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if strings.TrimSpace(code) == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "item code is required")
		return
	}

	item, err := h.svc.GetStockByItemCode(r.Context(), code)
	if err != nil {
		if err == repository.ErrNotFound {
			h.jsonError(w, http.StatusNotFound, "not_found", "item SKU not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, item)
}

// HandleListItems returns items filtered by category or low stock flag.
func (h *HTTPHandler) HandleListItems(w http.ResponseWriter, r *http.Request) {
	category := domain.ItemCategory(r.URL.Query().Get("category"))
	lowStock := r.URL.Query().Get("low_stock") == "true"

	items, err := h.svc.ListInventory(r.Context(), category, lowStock)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if items == nil {
		items = []domain.StoreItem{}
	}
	h.jsonResponse(w, http.StatusOK, items)
}

// HandleUpdateStock adjusts inventory quantity balance.
func (h *HTTPHandler) HandleUpdateStock(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	itemID := r.PathValue("id")
	var dto UpdateStockDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	updated, err := h.svc.UpdateInventoryStock(r.Context(), itemID, dto.Delta, claims.StaffID)
	if err != nil {
		if err == service.ErrInsufficientStock {
			h.jsonError(w, http.StatusBadRequest, "insufficient_stock", err.Error())
			return
		}
		h.jsonError(w, http.StatusUnprocessableEntity, "update_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, updated)
}

// HandleInitiateRequisition initiates a new material requisition.
func (h *HTTPHandler) HandleInitiateRequisition(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	var dto InitiateRequisitionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	var items []domain.RequisitionItem
	for _, it := range dto.Items {
		items = append(items, domain.RequisitionItem{
			ItemCode:          it.ItemCode,
			QuantityRequested: it.QuantityRequested,
			UnitCost:          it.UnitCost,
		})
	}

	req := domain.MaterialRequisition{
		OriginDeptID:        dto.OriginDeptID,
		LinkedTicketOrJobID: dto.LinkedTicketOrJobID,
		UnitHeadID:          claims.StaffID,
		Items:               items,
	}

	created, err := h.svc.InitiateRequisition(r.Context(), req)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "requisition_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, created)
}

// HandleDeptApproveRequisition records Departmental Head sign-off.
func (h *HTTPHandler) HandleDeptApproveRequisition(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	reqID := r.PathValue("id")
	approved, err := h.svc.ApproveRequisitionByDept(r.Context(), reqID, claims.StaffID)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "approval_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, approved)
}

// HandleAuditConsentRequisition records Internal Audit compliance consent.
func (h *HTTPHandler) HandleAuditConsentRequisition(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	reqID := r.PathValue("id")
	var dto AuditConsentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	decided, err := h.svc.ConsentRequisitionByAudit(r.Context(), reqID, claims.StaffID, dto.Consent, dto.Notes)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "audit_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, decided)
}

// HandleMDApproveRequisition records Managing Director approval for supplier purchase.
func (h *HTTPHandler) HandleMDApproveRequisition(w http.ResponseWriter, r *http.Request) {
	reqID := r.PathValue("id")
	var dto MDApprovalDTO
	_ = json.NewDecoder(r.Body).Decode(&dto)

	approved, err := h.svc.ApproveRequisitionByMD(r.Context(), reqID, dto.Notes)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "md_approval_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, approved)
}

// HandleIssueMaterialWaybill issues materials and creates waybill.
func (h *HTTPHandler) HandleIssueMaterialWaybill(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing authentication claims")
		return
	}

	reqID := r.PathValue("id")
	issued, err := h.svc.IssueMaterialWaybill(r.Context(), reqID, claims.StaffID)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "issuance_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, issued)
}

// HandleGetRequisition retrieves requisition details.
func (h *HTTPHandler) HandleGetRequisition(w http.ResponseWriter, r *http.Request) {
	reqID := r.PathValue("id")
	req, err := h.svc.GetRequisition(r.Context(), reqID)
	if err != nil {
		if err == repository.ErrNotFound {
			h.jsonError(w, http.StatusNotFound, "not_found", "requisition not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, req)
}

// HandleListRequisitions lists requisitions.
func (h *HTTPHandler) HandleListRequisitions(w http.ResponseWriter, r *http.Request) {
	status := domain.RequisitionStatus(r.URL.Query().Get("status"))
	unitHeadID := r.URL.Query().Get("unit_head_id")

	requisitions, err := h.svc.ListRequisitions(r.Context(), status, unitHeadID)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if requisitions == nil {
		requisitions = []domain.MaterialRequisition{}
	}
	h.jsonResponse(w, http.StatusOK, requisitions)
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
