// Package domain_test verifies Maker-Checker segregation, risk threshold guards, and audit hash chaining.
package domain_test

import (
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
)

func TestValidateMakerCheckerApproval(t *testing.T) {
	req := domain.MakerCheckerRequest{
		ID:              "req-1",
		Module:          "s08_billing",
		ChangeCategory:  domain.ChangeTariffBandChange,
		TargetRecordID:  "cust-100",
		OllamaRiskScore: 0.85,
		Status:          domain.StatusPending,
		MakerID:         "staff-maker-1",
	}

	// 1. Maker attempting to self-approve -> Forbidden
	err := domain.ValidateMakerCheckerApproval(req, "staff-maker-1", "HEAD_COMMERCIAL_OPS")
	if err != domain.ErrSelfApprovalForbidden {
		t.Fatalf("expected ErrSelfApprovalForbidden, got: %v", err)
	}

	// 2. High-risk change (0.85) attempted by non-executive -> Forbidden
	err = domain.ValidateMakerCheckerApproval(req, "staff-checker-2", "UNIT_HEAD_BILLING")
	if err != domain.ErrExecutiveApprovalRequired && err == nil {
		t.Fatalf("expected ErrExecutiveApprovalRequired for score >= 0.80")
	}

	// 3. High-risk change (0.85) approved by Executive -> Success
	err = domain.ValidateMakerCheckerApproval(req, "staff-checker-md", "MD")
	if err != nil {
		t.Fatalf("expected executive approval to succeed, got: %v", err)
	}

	// 4. Low-risk change (0.35) approved by Unit Head -> Success
	lowRiskReq := req
	lowRiskReq.OllamaRiskScore = 0.35
	err = domain.ValidateMakerCheckerApproval(lowRiskReq, "staff-checker-2", "UNIT_HEAD_BILLING")
	if err != nil {
		t.Fatalf("expected low-risk unit head approval to succeed, got: %v", err)
	}

	// 5. Already approved request re-evaluated -> Error
	finalizedReq := req
	finalizedReq.Status = domain.StatusApproved
	err = domain.ValidateMakerCheckerApproval(finalizedReq, "staff-checker-md", "MD")
	if err == nil {
		t.Fatalf("expected error evaluating already finalized request")
	}
}

func TestAuditChain_ComputeAndValidate(t *testing.T) {
	now := time.Now()
	prevHash := crypto.GenesisHash
	newJSON := `{"tariff_band":"A","status":"ACTIVE"}`
	hash := domain.ComputeAuditHash(prevHash, "customers", "cust-1", domain.ActionInsert, newJSON, "maker-1", now.Unix())

	rec := domain.AuditRecord{
		ID:          "audit-1",
		TableName:   "customers",
		RecordID:    "cust-1",
		ActionType:  domain.ActionInsert,
		NewJSON:     newJSON,
		MakerID:     "maker-1",
		PrevHash:    prevHash,
		CurrentHash: hash,
		Timestamp:   now,
	}

	// Valid chain record
	if err := domain.ValidateAuditChainRecord(rec, prevHash); err != nil {
		t.Fatalf("expected valid audit record validation, got: %v", err)
	}

	// Negative: Tampered current hash
	tamperedRec := rec
	tamperedRec.CurrentHash = "tampered-hash-0000000000000000000000000000000000000000000000000000"
	if err := domain.ValidateAuditChainRecord(tamperedRec, prevHash); err == nil {
		t.Fatalf("expected error on tampered hash")
	}

	// Negative: Wrong predecessor hash
	if err := domain.ValidateAuditChainRecord(rec, "wrong-predecessor-hash"); err == nil {
		t.Fatalf("expected error on wrong prev_hash")
	}
}

