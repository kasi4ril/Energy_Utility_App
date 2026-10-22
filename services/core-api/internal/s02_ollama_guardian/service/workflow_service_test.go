// Package service_test validates end-to-end Maker-Checker orchestration,
// cryptographic audit chain hashing, and @GridAssistant chat handling.
package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/repository"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/service"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

func setupTestService(sidecarURL string) (*service.WorkflowService, *repository.MemoryRepository) {
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, bus, sidecarURL)
	return svc, repo
}

func TestWorkflowService_SubmitAndApproveMakerChecker(t *testing.T) {
	// Setup mock AI sidecar server
	mockSidecar := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/internal/ai/risk-score" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"risk_score": 0.88,
				"reasoning":  "Mock sidecar: High risk debt adjustment",
			})
		} else if r.URL.Path == "/internal/ai/chat" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"response": "[@GridAssistant]: In-stock items verified via sidecar.",
			})
		}
	}))
	defer mockSidecar.Close()

	svc, repo := setupTestService(mockSidecar.URL)
	ctx := context.Background()

	// 1. Submit Request
	req, err := svc.SubmitMakerCheckerRequest(
		ctx, "customercare", domain.ChangeDebtAdjustment, "cust-uuid-10",
		`{"adjusted_balance":0.0,"old_balance":65000.0}`,
		"Management debt forgiveness settlement",
		"maker-uuid-1",
	)
	if err != nil {
		t.Fatalf("failed to submit maker-checker request: %v", err)
	}

	if req.Status != domain.StatusPending || req.OllamaRiskScore != 0.88 {
		t.Fatalf("expected sidecar score of 0.88, got: %f", req.OllamaRiskScore)
	}

	// 2. Negative: Maker self-approves -> Must fail
	_, _, err = svc.ApproveMakerCheckerRequest(ctx, req.ID, "maker-uuid-1", "MD", "Self approval")
	if err != domain.ErrSelfApprovalForbidden {
		t.Fatalf("expected ErrSelfApprovalForbidden, got: %v", err)
	}

	// 3. Negative: Non-executive approves high risk (score 0.88) -> Must fail
	_, _, err = svc.ApproveMakerCheckerRequest(ctx, req.ID, "checker-unit-head", "UNIT_HEAD_BILLING", "Approved")
	if err == nil {
		t.Fatalf("expected ErrExecutiveApprovalRequired for high-risk change")
	}

	// 4. Positive: Executive approves -> Success & Audit Chain row created
	updatedReq, auditRec, err := svc.ApproveMakerCheckerRequest(ctx, req.ID, "checker-md", "MD", "Executive approved")
	if err != nil {
		t.Fatalf("executive approval failed: %v", err)
	}

	if updatedReq.Status != domain.StatusApproved || auditRec == nil {
		t.Fatalf("expected approved status and valid audit record")
	}

	// 5. Verify audit chain continuity
	validChain, err := svc.VerifyAuditChainContinuity(ctx, "customercare", "cust-uuid-10")
	if err != nil || !validChain {
		t.Fatalf("expected audit chain continuity to pass, got err: %v", err)
	}

	// Empty audit record continuity check
	emptyChain, err := svc.VerifyAuditChainContinuity(ctx, "unknown_table", "rec-0")
	if err != nil || !emptyChain {
		t.Fatalf("expected empty chain to be valid")
	}

	// Negative: Tampered audit record in database
	tamperedRec := *auditRec
	tamperedRec.ID = "audit-tampered-id"
	tamperedRec.CurrentHash = "tampered-hash-value-0000000000000000000000000000000000000000000"
	tamperedRec.Timestamp = time.Now().Add(1 * time.Minute)
	_ = repo.CreateAuditRecord(ctx, tamperedRec)

	validChain, err = svc.VerifyAuditChainContinuity(ctx, "customercare", "cust-uuid-10")
	if validChain || err == nil {
		t.Fatalf("expected audit chain tampering to be detected")
	}

	// Negative: Approve non-existent request
	_, _, err = svc.ApproveMakerCheckerRequest(ctx, "non-existent-id", "chk", "MD", "note")
	if err == nil {
		t.Fatalf("expected error approving non-existent request")
	}
}

func TestWorkflowService_RejectMakerChecker(t *testing.T) {
	svc, _ := setupTestService("")
	ctx := context.Background()

	req, _ := svc.SubmitMakerCheckerRequest(
		ctx, "stores_audit", domain.ChangeStoreStockEdit, "item-10",
		`{"delta":-50}`, "Stock discrepancy", "maker-uuid-1",
	)

	// Negative: Maker self-rejects
	_, err := svc.RejectMakerCheckerRequest(ctx, req.ID, "maker-uuid-1", "STORE_OFFICER", "Reject myself")
	if err != domain.ErrSelfApprovalForbidden {
		t.Fatalf("expected ErrSelfApprovalForbidden for self-rejection, got: %v", err)
	}

	// Negative: Reject non-existent request
	_, err = svc.RejectMakerCheckerRequest(ctx, "non-existent", "checker-2", "UNIT_HEAD_STORE", "note")
	if err == nil {
		t.Fatalf("expected error rejecting missing request")
	}

	// Positive reject
	rejected, err := svc.RejectMakerCheckerRequest(ctx, req.ID, "checker-2", "UNIT_HEAD_STORE", "Discrepancy not justified")
	if err != nil {
		t.Fatalf("failed to reject request: %v", err)
	}
	if rejected.Status != domain.StatusRejected {
		t.Fatalf("expected status REJECTED, got: %s", rejected.Status)
	}

	// Attempting to reject again -> ErrRequestNotPending
	_, err = svc.RejectMakerCheckerRequest(ctx, req.ID, "checker-2", "UNIT_HEAD_STORE", "Second reject")
	if err != service.ErrRequestNotPending {
		t.Fatalf("expected ErrRequestNotPending on finalized request")
	}
}

func TestWorkflowService_ProcessChatMessage(t *testing.T) {
	mockSidecar := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"response": "[@GridAssistant]: In-stock items verified via sidecar.",
		})
	}))
	defer mockSidecar.Close()

	svc, repo := setupTestService(mockSidecar.URL)
	ctx := context.Background()

	senderID := "staff-uuid-1"

	// 1. Standard chat message without AI mention
	userMsg, aiMsg, err := svc.ProcessChatMessage(ctx, "general", &senderID, "Engr. Paul", "TEAM_LEAD_FEEDER", "Good morning team")
	if err != nil || userMsg == nil || aiMsg != nil {
		t.Fatalf("expected standard message without AI reply: %+v", userMsg)
	}

	// 2. Chat message with @GridAssistant mention -> AI response generated
	userMsg2, aiMsg2, err := svc.ProcessChatMessage(ctx, "general", &senderID, "Engr. Paul", "TEAM_LEAD_FEEDER", "@GridAssistant what is the conductor stock?")
	if err != nil || userMsg2 == nil || aiMsg2 == nil {
		t.Fatalf("expected AI co-pilot reply: %+v, %+v", userMsg2, aiMsg2)
	}
	if !aiMsg2.IsAIResponse || aiMsg2.SenderName != "@GridAssistant" {
		t.Fatalf("invalid AI response fields: %+v", aiMsg2)
	}

	msgs, _ := repo.ListChatMessages(ctx, "general", 10)
	if len(msgs) != 3 {
		t.Fatalf("expected 3 total messages recorded, got %d", len(msgs))
	}
}

