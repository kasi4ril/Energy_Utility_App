// Package transport_test validates HTTP REST API endpoints and authorization for S02.
package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/repository"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/service"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/transport"
	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const (
	testJWTSecret = "s02-transport-test-jwt-secret-32b!"
)

func setupTestServer() (http.Handler, *service.WorkflowService, *repository.MemoryRepository) {
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, bus, "")
	handler := transport.NewHTTPHandler(svc, repo, testJWTSecret)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, svc, repo
}

func makeAuthHeader(staffID, role, fullName string) string {
	claims := auth.Claims{
		StaffID:     staffID,
		PrimaryRole: role,
		FullName:    fullName,
	}
	token, _ := auth.GenerateToken(claims, testJWTSecret)
	return "Bearer " + token
}

func TestHTTP_MakerCheckerWorkflow(t *testing.T) {
	mux, _, repo := setupTestServer()

	makerAuth := makeAuthHeader("maker-1", "CUSTOMER_CARE_OFFICER", "Aisha Bello")
	checkerMDAuth := makeAuthHeader("checker-md", "MD", "Dr. Usman")

	// 1. Submit Request
	submitBody, _ := json.Marshal(transport.SubmitMakerCheckerRequestDTO{
		Module:         "customercare",
		ChangeCategory: domain.ChangeTariffBandChange,
		TargetRecordID: "cust-001",
		PayloadJSON:    `{"tariff_band":"A"}`,
		Justification:  "Commercial verification completed",
	})
	req := httptest.NewRequest("POST", "/api/v1/maker-checker/request", bytes.NewReader(submitBody))
	req.Header.Set("Authorization", makerAuth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on request submission, got %d: %s", rec.Code, rec.Body.String())
	}

	var createdReq domain.MakerCheckerRequest
	_ = json.NewDecoder(rec.Body).Decode(&createdReq)

	// Negative submit: empty body -> 400
	badReq := httptest.NewRequest("POST", "/api/v1/maker-checker/request", bytes.NewReader([]byte("{}")))
	badReq.Header.Set("Authorization", makerAuth)
	badRec := httptest.NewRecorder()
	mux.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on empty submit body")
	}

	// Negative submit: malformed json -> 400
	malformedReq := httptest.NewRequest("POST", "/api/v1/maker-checker/request", bytes.NewReader([]byte("{bad-json")))
	malformedReq.Header.Set("Authorization", makerAuth)
	malformedRec := httptest.NewRecorder()
	mux.ServeHTTP(malformedRec, malformedReq)
	if malformedRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on malformed json")
	}

	// 2. List Requests
	listReq := httptest.NewRequest("GET", "/api/v1/maker-checker/requests?status=PENDING", nil)
	listReq.Header.Set("Authorization", checkerMDAuth)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing requests, got %d", listRec.Code)
	}

	// 3. Negative: Maker self-approves -> 422
	selfApproveReq := httptest.NewRequest("POST", "/api/v1/maker-checker/"+createdReq.ID+"/approve", bytes.NewReader([]byte(`{"checker_note":"self"}`)))
	selfApproveReq.Header.Set("Authorization", makerAuth)
	selfApproveRec := httptest.NewRecorder()
	mux.ServeHTTP(selfApproveRec, selfApproveReq)
	if selfApproveRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for self-approval, got %d", selfApproveRec.Code)
	}

	// 4. Executive approves -> 200 OK
	approveBody, _ := json.Marshal(transport.ApproveRequestDTO{CheckerNote: "Verified by MD"})
	approveReq := httptest.NewRequest("POST", "/api/v1/maker-checker/"+createdReq.ID+"/approve", bytes.NewReader(approveBody))
	approveReq.Header.Set("Authorization", checkerMDAuth)
	approveRec := httptest.NewRecorder()
	mux.ServeHTTP(approveRec, approveReq)
	if approveRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on approval, got %d: %s", approveRec.Code, approveRec.Body.String())
	}

	// Negative: Approve non-existent request -> 422
	missingApproveReq := httptest.NewRequest("POST", "/api/v1/maker-checker/missing-id/approve", bytes.NewReader(approveBody))
	missingApproveReq.Header.Set("Authorization", checkerMDAuth)
	missingApproveRec := httptest.NewRecorder()
	mux.ServeHTTP(missingApproveRec, missingApproveReq)
	if missingApproveRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for missing request on approval")
	}

	// 5. Query Audit Chain
	auditReq := httptest.NewRequest("GET", "/api/v1/audit/chain/customercare/cust-001", nil)
	auditReq.Header.Set("Authorization", checkerMDAuth)
	auditRec := httptest.NewRecorder()
	mux.ServeHTTP(auditRec, auditReq)
	if auditRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK querying audit chain, got %d", auditRec.Code)
	}

	var auditResp transport.AuditChainVerificationDTO
	_ = json.NewDecoder(auditRec.Body).Decode(&auditResp)
	if !auditResp.IsValidChain || auditResp.ChainCount != 1 {
		t.Fatalf("expected valid audit chain with 1 record, got %+v", auditResp)
	}

	// Broken chain detection -> 409 Conflict
	tamperedRec := domain.AuditRecord{
		ID:          "tampered-id-1",
		TableName:   "customercare",
		RecordID:    "cust-001",
		ActionType:  domain.ActionUpdate,
		NewJSON:     `{"tariff_band":"tampered"}`,
		MakerID:     "hacker-1",
		PrevHash:    "genesis-hash",
		CurrentHash: "invalid-hash",
		Timestamp:   time.Now().Add(5 * time.Minute),
	}
	_ = repo.CreateAuditRecord(context.Background(), tamperedRec)

	brokenReq := httptest.NewRequest("GET", "/api/v1/audit/chain/customercare/cust-001", nil)
	brokenReq.Header.Set("Authorization", checkerMDAuth)
	brokenRec := httptest.NewRecorder()
	mux.ServeHTTP(brokenRec, brokenReq)
	if brokenRec.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict on broken audit chain, got %d", brokenRec.Code)
	}
}

func TestHTTP_RejectRequestAndChatFlow(t *testing.T) {
	mux, svc, _ := setupTestServer()
	ctx := context.Background()

	makerAuth := makeAuthHeader("maker-2", "STORE_OFFICER", "Musa Ali")
	checkerAuth := makeAuthHeader("checker-head", "HEAD_SERVICES", "Engr. Paul")

	// Submit request
	reqObj, _ := svc.SubmitMakerCheckerRequest(
		ctx, "stores_audit", domain.ChangeStoreStockEdit, "item-100",
		`{"delta":-5}`, "Spares count variance", "maker-2",
	)

	// Reject request
	rejectBody, _ := json.Marshal(transport.RejectRequestDTO{Reason: "Physical recount required first"})
	rejectReq := httptest.NewRequest("POST", "/api/v1/maker-checker/"+reqObj.ID+"/reject", bytes.NewReader(rejectBody))
	rejectReq.Header.Set("Authorization", checkerAuth)
	rejectRec := httptest.NewRecorder()
	mux.ServeHTTP(rejectRec, rejectReq)

	if rejectRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on rejection, got %d", rejectRec.Code)
	}

	// Negative: Reject non-existent request -> 422
	missingRejectReq := httptest.NewRequest("POST", "/api/v1/maker-checker/missing-id/reject", bytes.NewReader(rejectBody))
	missingRejectReq.Header.Set("Authorization", checkerAuth)
	missingRejectRec := httptest.NewRecorder()
	mux.ServeHTTP(missingRejectRec, missingRejectReq)
	if missingRejectRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 rejecting missing request")
	}

	// Negative reject: empty reason -> 400
	badRejectReq := httptest.NewRequest("POST", "/api/v1/maker-checker/"+reqObj.ID+"/reject", bytes.NewReader([]byte(`{"reason":""}`)))
	badRejectReq.Header.Set("Authorization", checkerAuth)
	badRejectRec := httptest.NewRecorder()
	mux.ServeHTTP(badRejectRec, badRejectReq)
	if badRejectRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty rejection reason")
	}

	// Negative reject: malformed JSON -> 400
	malformedReject := httptest.NewRequest("POST", "/api/v1/maker-checker/"+reqObj.ID+"/reject", bytes.NewReader([]byte(`{invalid-json`)))
	malformedReject.Header.Set("Authorization", checkerAuth)
	malformedRec := httptest.NewRecorder()
	mux.ServeHTTP(malformedRec, malformedReject)
	if malformedRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed rejection payload")
	}

	// Chat Message: POST
	chatBody, _ := json.Marshal(transport.ChatMessageDTO{
		ChannelID:   "general",
		MessageText: "@GridAssistant check meter stock",
	})
	chatReq := httptest.NewRequest("POST", "/api/v1/chat/message", bytes.NewReader(chatBody))
	chatReq.Header.Set("Authorization", makerAuth)
	chatRec := httptest.NewRecorder()
	mux.ServeHTTP(chatRec, chatReq)

	if chatRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on chat message, got %d", chatRec.Code)
	}

	// Negative chat: empty message -> 400
	badChatReq := httptest.NewRequest("POST", "/api/v1/chat/message", bytes.NewReader([]byte(`{"message_text":""}`)))
	badChatReq.Header.Set("Authorization", makerAuth)
	badChatRec := httptest.NewRecorder()
	mux.ServeHTTP(badChatRec, badChatReq)
	if badChatRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty chat message")
	}

	// Negative chat: malformed json -> 400
	malformedChat := httptest.NewRequest("POST", "/api/v1/chat/message", bytes.NewReader([]byte(`{bad`)))
	malformedChat.Header.Set("Authorization", makerAuth)
	malformedChatRec := httptest.NewRecorder()
	mux.ServeHTTP(malformedChatRec, malformedChat)
	if malformedChatRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed chat json")
	}

	// Chat Messages: GET without limit param
	listChatReqNoLimit := httptest.NewRequest("GET", "/api/v1/chat/messages", nil)
	listChatReqNoLimit.Header.Set("Authorization", makerAuth)
	listChatRecNoLimit := httptest.NewRecorder()
	mux.ServeHTTP(listChatRecNoLimit, listChatReqNoLimit)
	if listChatRecNoLimit.Code != http.StatusOK {
		t.Fatalf("expected 200 OK listing messages without limit")
	}
}
