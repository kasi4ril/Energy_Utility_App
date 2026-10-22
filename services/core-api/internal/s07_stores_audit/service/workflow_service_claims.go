// Package service implements the 3-way automated staff claims verification desk and auditor sign-off.
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
)

// SubmitAllowanceClaim registers an allowance claim and executes automated 3-way verification.
func (s *WorkflowService) SubmitAllowanceClaim(ctx context.Context, claim domain.StaffAllowanceClaim) (*domain.StaffAllowanceClaim, error) {
	if strings.TrimSpace(claim.StaffID) == "" || claim.Amount <= 0.0 || strings.TrimSpace(claim.ClaimDate) == "" {
		return nil, domain.ErrInvalidClaimData
	}

	now := time.Now().UTC()
	if claim.ID == "" {
		claim.ID = crypto.MustUUID()
	}
	if claim.ClaimNo == "" {
		claim.ClaimNo = fmt.Sprintf("CLM-%d-%s", now.Year(), strings.ToUpper(crypto.MustUUID()[:8]))
	}

	// Automated 3-way cross-check (Roster, PTW, GPS)
	flag, fraudScore, initialStatus := domain.Evaluate3WayClaimVerification(
		claim,
		claim.RosterMatchVerified,
		claim.PTWCrewVerified,
		claim.GPSMatchVerified,
	)

	claim.OllamaAuditFlag = flag
	claim.OllamaFraudScore = fraudScore
	claim.Status = initialStatus
	claim.CreatedAt = now
	claim.UpdatedAt = now

	if err := s.repo.CreateClaim(ctx, claim); err != nil {
		return nil, fmt.Errorf("service: failed to create allowance claim: %w", err)
	}

	_ = s.bus.Publish(ctx, "audit.claim_submitted", "s07_stores_audit", map[string]interface{}{
		"claim_id":    claim.ID,
		"claim_no":    claim.ClaimNo,
		"staff_id":    claim.StaffID,
		"amount":      claim.Amount,
		"audit_flag":  claim.OllamaAuditFlag,
		"fraud_score": claim.OllamaFraudScore,
		"status":      claim.Status,
	})

	return &claim, nil
}

// ReviewClaimByAuditor records an Internal Auditor's formal approval or rejection on a claim.
func (s *WorkflowService) ReviewClaimByAuditor(ctx context.Context, claimID, auditorID string, approve bool, notes string) (*domain.StaffAllowanceClaim, error) {
	claim, err := s.repo.GetClaimByID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateAuditorClaimDecision(claim, auditorID, approve, notes); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateClaim(ctx, *claim); err != nil {
		return nil, fmt.Errorf("service: failed to update claim: %w", err)
	}

	_ = s.bus.Publish(ctx, "audit.claim_decided", "s07_stores_audit", map[string]interface{}{
		"claim_id":   claim.ID,
		"claim_no":   claim.ClaimNo,
		"auditor_id": auditorID,
		"status":     claim.Status,
	})

	return claim, nil
}

// GetClaimByID retrieves a financial claim by UUID.
func (s *WorkflowService) GetClaimByID(ctx context.Context, id string) (*domain.StaffAllowanceClaim, error) {
	return s.repo.GetClaimByID(ctx, id)
}

// ListClaims returns allowance claims matching optional filters.
func (s *WorkflowService) ListClaims(ctx context.Context, staffID string, status domain.ClaimStatus) ([]domain.StaffAllowanceClaim, error) {
	return s.repo.ListClaims(ctx, staffID, status)
}

