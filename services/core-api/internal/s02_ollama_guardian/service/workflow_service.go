// Package service coordinates Maker-Checker risk scoring, cryptographic audit chain
// hashing, AI sidecar communication, and community co-pilot interactions.
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/repository"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

var (
	// ErrRequestNotPending indicates the change request is not in a reviewable state.
	ErrRequestNotPending = errors.New("service: request is not in pending status")
)

// WorkflowService orchestrates all operations for S02 Ollama Local AI Guardian.
type WorkflowService struct {
	repo         repository.Repository
	bus          events.EventBus
	aiSidecarURL string
	httpClient   *http.Client
}

// NewWorkflowService constructs a new Security & AI Workflow service.
func NewWorkflowService(repo repository.Repository, bus events.EventBus, aiSidecarURL string) *WorkflowService {
	if aiSidecarURL == "" {
		aiSidecarURL = "http://127.0.0.1:8001"
	}
	return &WorkflowService{
		repo:         repo,
		bus:          bus,
		aiSidecarURL: aiSidecarURL,
		httpClient:   &http.Client{Timeout: 5 * time.Second},
	}
}

// SubmitMakerCheckerRequest evaluates AI fraud risk, records request, and alerts supervisors.
func (s *WorkflowService) SubmitMakerCheckerRequest(ctx context.Context, module string, category domain.ChangeCategory, targetID, payloadJSON, justification, makerID string) (*domain.MakerCheckerRequest, error) {
	risk := s.evaluateRiskWithSidecar(ctx, category, payloadJSON, justification)

	reqID := crypto.MustUUID()
	now := time.Now()

	req := domain.MakerCheckerRequest{
		ID:              reqID,
		Module:          module,
		ChangeCategory:  category,
		TargetRecordID:  targetID,
		PayloadJSON:     payloadJSON,
		Justification:   justification,
		OllamaRiskScore: risk.RiskScore,
		OllamaReasoning: risk.Reasoning,
		Status:          domain.StatusPending,
		MakerID:         makerID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.CreateMakerCheckerRequest(ctx, req); err != nil {
		return nil, fmt.Errorf("service: failed to save request: %w", err)
	}

	_ = s.bus.Publish(ctx, "security.maker_checker.submitted", "s02_ollama_guardian", req)

	if req.OllamaRiskScore >= domain.HighRiskThreshold {
		_ = s.bus.Publish(ctx, "security.fraud_alert.high_risk", "s02_ollama_guardian", map[string]interface{}{
			"request_id": req.ID,
			"category":   category,
			"risk_score": req.OllamaRiskScore,
			"reasoning":  req.OllamaReasoning,
		})
	}

	return &req, nil
}

// ApproveMakerCheckerRequest enforces maker-checker segregation, updates state, and appends to audit chain.
func (s *WorkflowService) ApproveMakerCheckerRequest(ctx context.Context, requestID, checkerID, checkerRole, note string) (*domain.MakerCheckerRequest, *domain.AuditRecord, error) {
	req, err := s.repo.GetMakerCheckerRequestByID(ctx, requestID)
	if err != nil {
		return nil, nil, err
	}

	if err := domain.ValidateMakerCheckerApproval(*req, checkerID, checkerRole); err != nil {
		return nil, nil, err
	}

	if err := s.repo.UpdateMakerCheckerStatus(ctx, requestID, domain.StatusApproved, checkerID, note); err != nil {
		return nil, nil, fmt.Errorf("service: failed to update status: %w", err)
	}

	// Append to cryptographic audit chain
	latestAudit, err := s.repo.GetLatestAuditRecord(ctx, req.Module, req.TargetRecordID)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to query audit chain: %w", err)
	}

	prevHash := crypto.GenesisHash
	if latestAudit != nil {
		prevHash = latestAudit.CurrentHash
	}

	now := time.Now()
	auditHash := domain.ComputeAuditHash(
		prevHash, req.Module, req.TargetRecordID, domain.ActionMakerCheckerApproval, req.PayloadJSON, checkerID, now.Unix(),
	)

	auditRec := domain.AuditRecord{
		ID:          crypto.MustUUID(),
		TableName:   req.Module,
		RecordID:    req.TargetRecordID,
		ActionType:  domain.ActionMakerCheckerApproval,
		OldJSON:     "",
		NewJSON:     req.PayloadJSON,
		MakerID:     req.MakerID,
		CheckerID:   &checkerID,
		PrevHash:    prevHash,
		CurrentHash: auditHash,
		Timestamp:   now,
	}

	if err := s.repo.CreateAuditRecord(ctx, auditRec); err != nil {
		return nil, nil, fmt.Errorf("service: failed to append audit record: %w", err)
	}

	_ = s.bus.Publish(ctx, "security.maker_checker.approved", "s02_ollama_guardian", map[string]interface{}{
		"request_id":       req.ID,
		"module":           req.Module,
		"target_record_id": req.TargetRecordID,
		"payload_json":     req.PayloadJSON,
		"audit_hash":       auditHash,
	})

	updatedReq, _ := s.repo.GetMakerCheckerRequestByID(ctx, requestID)
	return updatedReq, &auditRec, nil
}

// RejectMakerCheckerRequest rejects a pending change request with checker justification.
func (s *WorkflowService) RejectMakerCheckerRequest(ctx context.Context, requestID, checkerID, checkerRole, note string) (*domain.MakerCheckerRequest, error) {
	req, err := s.repo.GetMakerCheckerRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	if req.Status != domain.StatusPending {
		return nil, ErrRequestNotPending
	}
	if req.MakerID == checkerID {
		return nil, domain.ErrSelfApprovalForbidden
	}

	if err := s.repo.UpdateMakerCheckerStatus(ctx, requestID, domain.StatusRejected, checkerID, note); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "security.maker_checker.rejected", "s02_ollama_guardian", map[string]interface{}{
		"request_id": req.ID,
		"module":     req.Module,
		"reason":     note,
	})

	return s.repo.GetMakerCheckerRequestByID(ctx, requestID)
}

// ProcessChatMessage handles in-app messaging and queries the @GridAssistant co-pilot.
func (s *WorkflowService) ProcessChatMessage(ctx context.Context, channelID string, senderID *string, senderName, senderRole, text string) (*domain.ChatMessage, *domain.ChatMessage, error) {
	userMsg := domain.ChatMessage{
		ID:           crypto.MustUUID(),
		ChannelID:    channelID,
		SenderID:     senderID,
		SenderName:   senderName,
		MessageText:  text,
		IsAIResponse: false,
		CreatedAt:    time.Now(),
	}
	if err := s.repo.CreateChatMessage(ctx, userMsg); err != nil {
		return nil, nil, err
	}

	// Query AI Sidecar if message addresses @GridAssistant
	var aiMsg *domain.ChatMessage
	if isAssistantQuery(text) {
		aiRespText := s.queryAssistantWithSidecar(ctx, senderRole, text)
		aiRecord := domain.ChatMessage{
			ID:           crypto.MustUUID(),
			ChannelID:    channelID,
			SenderName:   "@GridAssistant",
			MessageText:  aiRespText,
			IsAIResponse: true,
			CreatedAt:    time.Now(),
		}
		_ = s.repo.CreateChatMessage(ctx, aiRecord)
		aiMsg = &aiRecord
	}

	_ = s.bus.Publish(ctx, "chat.message.received", "s02_ollama_guardian", userMsg)
	return &userMsg, aiMsg, nil
}

// VerifyAuditChainContinuity verifies the complete cryptographic SHA-256 chain for an entity.
func (s *WorkflowService) VerifyAuditChainContinuity(ctx context.Context, tableName, recordID string) (bool, error) {
	records, err := s.repo.ListAuditRecords(ctx, tableName, recordID)
	if err != nil {
		return false, err
	}
	if len(records) == 0 {
		return true, nil
	}

	expectedPrev := crypto.GenesisHash
	for _, rec := range records {
		if err := domain.ValidateAuditChainRecord(rec, expectedPrev); err != nil {
			return false, err
		}
		expectedPrev = rec.CurrentHash
	}
	return true, nil
}

func (s *WorkflowService) evaluateRiskWithSidecar(ctx context.Context, category domain.ChangeCategory, payloadJSON, justification string) domain.RiskEvaluation {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"change_category": category,
		"payload":         json.RawMessage(payloadJSON),
		"justification":   justification,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", s.aiSidecarURL+"/internal/ai/risk-score", bytes.NewReader(reqBody))
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var res struct {
				RiskScore float64 `json:"risk_score"`
				Reasoning string  `json:"reasoning"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
				return domain.RiskEvaluation{RiskScore: res.RiskScore, Reasoning: res.Reasoning}
			}
		}
	}

	// Deterministic local fallback if AI Sidecar container is starting up or offline
	score := 0.45
	if category == domain.ChangeTariffBandChange || category == domain.ChangeDebtAdjustment {
		score = 0.85
	}
	return domain.RiskEvaluation{
		RiskScore: score,
		Reasoning: "Evaluated by local security heuristic engine (Air-Gapped Fallback).",
	}
}

func (s *WorkflowService) queryAssistantWithSidecar(ctx context.Context, role, text string) string {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"sender_role":  role,
		"message_text": text,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", s.aiSidecarURL+"/internal/ai/chat", bytes.NewReader(reqBody))
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var res struct {
				Response string `json:"response"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&res); err == nil && res.Response != "" {
				return res.Response
			}
		}
	}

	return "[@GridAssistant]: Live store database online. Please verify catalog code before raising requisition."
}

func isAssistantQuery(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "@gridassistant") || strings.Contains(lower, "@assistant")
}
