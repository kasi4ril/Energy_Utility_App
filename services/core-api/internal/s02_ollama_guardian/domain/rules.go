// Package domain contains pure business rules and state validation for S02 Ollama Guardian.
package domain

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"disco-enterprise-platform/services/core-api/pkg/auth"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
)

var (
	// ErrSelfApprovalForbidden indicates a maker cannot approve their own request.
	ErrSelfApprovalForbidden = errors.New("domain: maker is strictly forbidden from checking/approving their own change request")
	// ErrUnauthorizedChecker indicates staff member does not have permission to check this request.
	ErrUnauthorizedChecker = errors.New("domain: staff role lacks checker authorization for this request")
	// ErrExecutiveApprovalRequired indicates high-risk requests require executive sign-off.
	ErrExecutiveApprovalRequired = errors.New("domain: high-risk change (score >= 0.80) strictly requires an Executive Tier Checker")
	// ErrInvalidStateTransition indicates an illegal transition from a finalized request state.
	ErrInvalidStateTransition = errors.New("domain: request is already finalized and cannot be modified")
	// ErrBrokenAuditChain indicates a hash chain continuity violation in the audit ledger.
	ErrBrokenAuditChain = errors.New("domain: cryptographic audit chain hash verification failed")
)

// HighRiskThreshold defines the boundary where executive tier sign-off is mandatory.
const HighRiskThreshold = 0.80

// ValidateMakerCheckerApproval enforces maker-checker segregation of duties and risk tier authorization.
func ValidateMakerCheckerApproval(req MakerCheckerRequest, checkerID, checkerRole string) error {
	if req.Status != StatusPending {
		return fmt.Errorf("%w: current status is %s", ErrInvalidStateTransition, req.Status)
	}

	// 1. Strict Segregation of Duties: Maker cannot be Checker
	if req.MakerID == checkerID {
		return ErrSelfApprovalForbidden
	}

	// 2. High-Risk Threshold Enforcement
	if req.OllamaRiskScore >= HighRiskThreshold {
		if !auth.IsExecutiveRole(checkerRole) {
			return fmt.Errorf("%w: risk score is %.2f, requiring Executive sign-off", ErrExecutiveApprovalRequired, req.OllamaRiskScore)
		}
	}

	return nil
}

// ComputeAuditHash calculates an unalterable SHA-256 hash chaining an audit record.
// Formula: SHA256(prev_hash + ":" + table_name + ":" + record_id + ":" + action + ":" + new_json + ":" + maker_id + ":" + timestamp)
func ComputeAuditHash(prevHash, tableName, recordID string, action ActionType, newJSON, makerID string, timestamp int64) string {
	if strings.TrimSpace(prevHash) == "" {
		prevHash = crypto.GenesisHash
	}
	payload := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%d",
		prevHash, tableName, recordID, action, newJSON, makerID, timestamp)

	h := sha256.New()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// ValidateAuditChainRecord verifies that an audit record's current hash matches its predecessor and payload.
func ValidateAuditChainRecord(rec AuditRecord, expectedPrevHash string) error {
	if expectedPrevHash == "" {
		expectedPrevHash = crypto.GenesisHash
	}
	if rec.PrevHash != expectedPrevHash {
		return fmt.Errorf("%w: expected prev_hash %s but got %s", ErrBrokenAuditChain, expectedPrevHash, rec.PrevHash)
	}

	actorID := rec.MakerID
	if rec.CheckerID != nil && *rec.CheckerID != "" {
		actorID = *rec.CheckerID
	}

	expectedHash := ComputeAuditHash(
		rec.PrevHash, rec.TableName, rec.RecordID, rec.ActionType, rec.NewJSON, actorID, rec.Timestamp.Unix(),
	)

	if subtle.ConstantTimeCompare([]byte(rec.CurrentHash), []byte(expectedHash)) != 1 {
		return fmt.Errorf("%w: hash mismatch, recomputed %s vs stored %s", ErrBrokenAuditChain, expectedHash, rec.CurrentHash)
	}

	return nil
}
