// Package transport implements REST API handlers for Foundation (S01).
package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler encapsulates HTTP endpoints for the foundation bounded context.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTP handler instance.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers all Foundation endpoints onto the provided ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// Public Auth Route
	mux.HandleFunc("POST /api/v1/auth/login", h.HandleLogin)

	// Org Units
	mux.Handle("GET /api/v1/foundation/org-units", authMW(http.HandlerFunc(h.HandleListOrgUnits)))
	mux.Handle("POST /api/v1/foundation/org-units", authMW(http.HandlerFunc(h.HandleRegisterOrgUnit)))

	// Staff
	mux.Handle("GET /api/v1/foundation/staff", authMW(http.HandlerFunc(h.HandleListStaff)))
	mux.Handle("POST /api/v1/foundation/staff", authMW(http.HandlerFunc(h.HandleRegisterStaff)))
	mux.Handle("POST /api/v1/foundation/staff/transfer", authMW(http.HandlerFunc(h.HandleTransferStaff)))

	// Grid Assets
	mux.Handle("GET /api/v1/foundation/assets", authMW(http.HandlerFunc(h.HandleListAssets)))
	mux.Handle("POST /api/v1/foundation/assets", authMW(http.HandlerFunc(h.HandleRegisterAsset)))

	// WorkNotes Ledger
	mux.Handle("POST /api/v1/notes/{entityType}/{entityId}", authMW(http.HandlerFunc(h.HandleAppendWorkNote)))
	mux.Handle("GET /api/v1/notes/{entityType}/{entityId}", authMW(http.HandlerFunc(h.HandleListWorkNotes)))
}

// HandleLogin authenticates staff and issues JWT tokens.
func (h *HTTPHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := req.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	claims, token, err := h.svc.AuthenticateStaff(r.Context(), req.Email, req.Password)
	if err != nil {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, LoginResponse{Token: token, Claims: *claims})
}

// HandleListOrgUnits retrieves all registered administrative units.
func (h *HTTPHandler) HandleListOrgUnits(w http.ResponseWriter, r *http.Request) {
	units, err := h.repo.ListOrgUnits(r.Context())
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, units)
}

// HandleRegisterOrgUnit registers a new administrative unit.
func (h *HTTPHandler) HandleRegisterOrgUnit(w http.ResponseWriter, r *http.Request) {
	var req RegisterOrgUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := req.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	var gps *domain.GPSPoint
	if req.Latitude != nil && req.Longitude != nil {
		gps = &domain.GPSPoint{Latitude: *req.Latitude, Longitude: *req.Longitude}
	}

	unit, err := h.svc.RegisterOrgUnit(r.Context(), req.UnitCode, req.Name, req.UnitType, req.ParentID, gps)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "unprocessable_entity", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, unit)
}

// HandleListStaff returns staff profiles, optionally filtered by org_unit_id query param.
func (h *HTTPHandler) HandleListStaff(w http.ResponseWriter, r *http.Request) {
	orgUnitID := r.URL.Query().Get("org_unit_id")
	staff, err := h.repo.ListStaffProfiles(r.Context(), orgUnitID)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, staff)
}

// HandleRegisterStaff adds a new employee to the federated HR registry.
func (h *HTTPHandler) HandleRegisterStaff(w http.ResponseWriter, r *http.Request) {
	var req RegisterStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := req.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	staff, err := h.svc.RegisterStaffProfile(r.Context(), req.StaffNumber, req.FullName, req.Email, req.Password, req.Role, req.OrgUnitID)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "unprocessable_entity", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, staff)
}

// HandleTransferStaff re-deploys a staff member to a new department or role.
func (h *HTTPHandler) HandleTransferStaff(w http.ResponseWriter, r *http.Request) {
	var req TransferStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := req.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	staff, err := h.svc.TransferStaff(r.Context(), req.StaffID, req.NewOrgUnitID, req.NewRole)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "unprocessable_entity", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, staff)
}

// HandleListAssets retrieves grid assets.
func (h *HTTPHandler) HandleListAssets(w http.ResponseWriter, r *http.Request) {
	parentID := r.URL.Query().Get("parent_id")
	assetType := domain.AssetType(r.URL.Query().Get("asset_type"))
	assets, err := h.repo.ListGridAssets(r.Context(), parentID, assetType)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, assets)
}

// HandleRegisterAsset registers a physical transmission or distribution asset.
func (h *HTTPHandler) HandleRegisterAsset(w http.ResponseWriter, r *http.Request) {
	var req RegisterGridAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := req.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	gps := domain.GPSPoint{Latitude: req.Latitude, Longitude: req.Longitude}
	asset, err := h.svc.RegisterGridAsset(r.Context(), req.AssetCode, req.AssetType, req.ParentAssetID, req.OrgUnitID, req.RatedCapacityKVA, req.TariffBand, gps)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "unprocessable_entity", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, asset)
}

// HandleAppendWorkNote appends a note to the hash chain.
func (h *HTTPHandler) HandleAppendWorkNote(w http.ResponseWriter, r *http.Request) {
	entityType := r.PathValue("entityType")
	entityID := r.PathValue("entityId")
	if entityType == "" || entityID == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "entityType and entityId path parameters required")
		return
	}

	var req AppendWorkNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := req.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	note, err := h.svc.AppendWorkNote(r.Context(), entityType, entityID, req.AuthorStaffID, req.TaggedUnitIDs, req.NoteText)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "unprocessable_entity", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, note)
}

// HandleListWorkNotes lists notes attached to an entity.
func (h *HTTPHandler) HandleListWorkNotes(w http.ResponseWriter, r *http.Request) {
	entityType := r.PathValue("entityType")
	entityID := r.PathValue("entityId")
	notes, err := h.repo.ListWorkNotesByEntity(r.Context(), entityType, entityID)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, notes)
}

func (h *HTTPHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *HTTPHandler) jsonError(w http.ResponseWriter, status int, errCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: strings.ToLower(errCode), Message: msg})
}

