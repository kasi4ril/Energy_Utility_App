// Package transport implements REST API endpoints for Maker-Checker, audit chains, and chat.
package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/repository"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/service"
	"disco-enterprise-platform/services/core-api/pkg/auth"
)

// HTTPHandler encapsulates REST endpoints for S02 Ollama Guardian.
type HTTPHandler struct {
	svc       *service.WorkflowService
	repo      repository.Repository
	jwtSecret string
}

// NewHTTPHandler constructs a new HTTPHandler instance for S02.
func NewHTTPHandler(svc *service.WorkflowService, repo repository.Repository, jwtSecret string) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo, jwtSecret: jwtSecret}
}

// RegisterRoutes registers S02 routes onto the ServeMux.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authMW := auth.RequireAuth(h.jwtSecret)

	// Maker-Checker Workflows
	mux.Handle("POST /api/v1/maker-checker/request", authMW(http.HandlerFunc(h.HandleSubmitRequest)))
	mux.Handle("POST /api/v1/maker-checker/{id}/approve", authMW(http.HandlerFunc(h.HandleApproveRequest)))
	mux.Handle("POST /api/v1/maker-checker/{id}/reject", authMW(http.HandlerFunc(h.HandleRejectRequest)))
	mux.Handle("GET /api/v1/maker-checker/requests", authMW(http.HandlerFunc(h.HandleListRequests)))

	// Immutable Cryptographic Audit Chain
	mux.Handle("GET /api/v1/audit/chain/{tableName}/{recordId}", authMW(http.HandlerFunc(h.HandleGetAuditChain)))

	// Community Chat & @GridAssistant Co-Pilot
	mux.Handle("POST /api/v1/chat/message", authMW(http.HandlerFunc(h.HandleSendMessage)))
	mux.Handle("GET /api/v1/chat/messages", authMW(http.HandlerFunc(h.HandleListMessages)))
}

// HandleSubmitRequest submits a high-risk operational change for AI fraud scoring and checker review.
func (h *HTTPHandler) HandleSubmitRequest(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing staff identity claims")
		return
	}

	var dto SubmitMakerCheckerRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	req, err := h.svc.SubmitMakerCheckerRequest(
		r.Context(), dto.Module, dto.ChangeCategory, dto.TargetRecordID, dto.PayloadJSON, dto.Justification, claims.StaffID,
	)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "unprocessable_entity", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusCreated, req)
}

// HandleApproveRequest verifies segregation of duties and records the cryptographic audit row.
func (h *HTTPHandler) HandleApproveRequest(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing staff identity claims")
		return
	}

	reqID := r.PathValue("id")
	if reqID == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "request id is required")
		return
	}

	var dto ApproveRequestDTO
	_ = json.NewDecoder(r.Body).Decode(&dto)

	updatedReq, auditRec, err := h.svc.ApproveMakerCheckerRequest(r.Context(), reqID, claims.StaffID, claims.PrimaryRole, dto.CheckerNote)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "approval_rejected", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"request":      updatedReq,
		"audit_record": auditRec,
	})
}

// HandleRejectRequest marks a pending request as rejected with checker justification.
func (h *HTTPHandler) HandleRejectRequest(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing staff identity claims")
		return
	}

	reqID := r.PathValue("id")
	if reqID == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "request id is required")
		return
	}

	var dto RejectRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	rejectedReq, err := h.svc.RejectMakerCheckerRequest(r.Context(), reqID, claims.StaffID, claims.PrimaryRole, dto.Reason)
	if err != nil {
		h.jsonError(w, http.StatusUnprocessableEntity, "rejection_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, rejectedReq)
}

// HandleListRequests lists change requests filtered by status or module.
func (h *HTTPHandler) HandleListRequests(w http.ResponseWriter, r *http.Request) {
	status := domain.RequestStatus(r.URL.Query().Get("status"))
	module := r.URL.Query().Get("module")

	list, err := h.repo.ListMakerCheckerRequests(r.Context(), status, module)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, list)
}

// HandleGetAuditChain validates and returns the complete SHA-256 chain for an entity record.
func (h *HTTPHandler) HandleGetAuditChain(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	recordID := r.PathValue("recordId")
	if tableName == "" || recordID == "" {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "tableName and recordId path params are required")
		return
	}

	isValid, err := h.svc.VerifyAuditChainContinuity(r.Context(), tableName, recordID)
	if err != nil {
		h.jsonError(w, http.StatusConflict, "audit_chain_broken", err.Error())
		return
	}

	records, _ := h.repo.ListAuditRecords(r.Context(), tableName, recordID)

	h.jsonResponse(w, http.StatusOK, AuditChainVerificationDTO{
		TableName:    tableName,
		RecordID:     recordID,
		ChainCount:   len(records),
		IsValidChain: isValid,
		Records:      records,
	})
}

// HandleSendMessage processes in-app chat and queries the @GridAssistant co-pilot.
func (h *HTTPHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "unauthorized", "missing staff identity claims")
		return
	}

	var dto ChatMessageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	channel := dto.ChannelID
	if channel == "" {
		channel = "general"
	}

	userMsg, aiMsg, err := h.svc.ProcessChatMessage(r.Context(), channel, &claims.StaffID, claims.FullName, claims.PrimaryRole, dto.MessageText)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"user_message": userMsg,
		"ai_message":   aiMsg,
	})
}

// HandleListMessages lists messages for a channel.
func (h *HTTPHandler) HandleListMessages(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel_id")
	if channel == "" {
		channel = "general"
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	msgs, err := h.repo.ListChatMessages(r.Context(), channel, limit)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	h.jsonResponse(w, http.StatusOK, msgs)
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

