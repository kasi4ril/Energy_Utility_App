// Package repository_test validates the in-memory repository implementation for S02.
package repository_test

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/repository"
)

func TestMemoryRepository_AuditChainOperations(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	// 1. Initial latest record check (empty)
	latest, err := repo.GetLatestAuditRecord(ctx, "customers", "cust-1")
	if err != nil || latest != nil {
		t.Fatalf("expected nil for non-existent audit record")
	}

	// 2. Append Audit Record
	rec1 := domain.AuditRecord{
		ID:          "rec-1",
		TableName:   "customers",
		RecordID:    "cust-1",
		ActionType:  domain.ActionInsert,
		NewJSON:     `{"name":"Danladi Musa"}`,
		MakerID:     "maker-uuid-1",
		PrevHash:    "genesis-hash",
		CurrentHash: "hash-001",
		Timestamp:   time.Now(),
	}
	if err := repo.CreateAuditRecord(ctx, rec1); err != nil {
		t.Fatalf("failed to create audit record: %v", err)
	}

	// 3. Verify latest record
	latest, err = repo.GetLatestAuditRecord(ctx, "customers", "cust-1")
	if err != nil || latest == nil || latest.ID != rec1.ID {
		t.Fatalf("latest record does not match created record: %+v", latest)
	}

	// 4. List records
	list, err := repo.ListAuditRecords(ctx, "customers", "cust-1")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 audit record, got: %d", len(list))
	}
}

func TestMemoryRepository_MakerCheckerOperations(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	req := domain.MakerCheckerRequest{
		ID:              "req-uuid-1",
		Module:          "s08_billing",
		ChangeCategory:  domain.ChangeTariffBandChange,
		TargetRecordID:  "cust-100",
		PayloadJSON:     `{"old_band":"A","new_band":"B"}`,
		Justification:   "Reclassified feeder supply hours",
		OllamaRiskScore: 0.45,
		OllamaReasoning: "Moderate tariff change",
		Status:          domain.StatusPending,
		MakerID:         "staff-1",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := repo.CreateMakerCheckerRequest(ctx, req); err != nil {
		t.Fatalf("failed to create maker-checker request: %v", err)
	}

	// Fetch by ID
	fetched, err := repo.GetMakerCheckerRequestByID(ctx, req.ID)
	if err != nil || fetched.TargetRecordID != req.TargetRecordID {
		t.Fatalf("failed to fetch created request: %v", err)
	}

	// List by status and module
	pendingBilling, err := repo.ListMakerCheckerRequests(ctx, domain.StatusPending, "s08_billing")
	if err != nil || len(pendingBilling) != 1 {
		t.Fatalf("expected 1 pending billing request, got: %d", len(pendingBilling))
	}

	// Update status
	err = repo.UpdateMakerCheckerStatus(ctx, req.ID, domain.StatusApproved, "checker-uuid-2", "Approved by Head of Commercial")
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	updated, _ := repo.GetMakerCheckerRequestByID(ctx, req.ID)
	if updated.Status != domain.StatusApproved || *updated.CheckerID != "checker-uuid-2" {
		t.Fatalf("updated status fields mismatched: %+v", updated)
	}

	// Missing ID checks
	if _, err := repo.GetMakerCheckerRequestByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing request ID")
	}
	if err := repo.UpdateMakerCheckerStatus(ctx, "unknown", domain.StatusRejected, "chk", "note"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound updating missing request")
	}
}

func TestMemoryRepository_ChatMessages(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	senderID := "staff-uuid-1"
	msg := domain.ChatMessage{
		ID:          "msg-1",
		ChannelID:   "general",
		SenderID:    &senderID,
		SenderName:  "Engr. Danladi",
		MessageText: "@GridAssistant do we have spare transformers?",
		CreatedAt:   time.Now(),
	}

	if err := repo.CreateChatMessage(ctx, msg); err != nil {
		t.Fatalf("failed to create chat message: %v", err)
	}

	msgs, err := repo.ListChatMessages(ctx, "general", 10)
	if err != nil || len(msgs) != 1 {
		t.Fatalf("expected 1 message in general channel, got: %d", len(msgs))
	}
}

