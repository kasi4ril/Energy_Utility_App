// Package service implements the 5-step Safety Permit-to-Work (PTW) isolation protocol
// to protect human lives and prevent re-energization during maintenance.
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
)

// InitiatePTWRequest contains input data for Step 1 of the safety permit.
type InitiatePTWRequest struct {
	AssetID         string  `json:"asset_id"`
	WorkOrderID     *string `json:"work_order_id,omitempty"`
	UnitHeadID      string  `json:"unit_head_id"`
	WorkDescription string  `json:"work_description"`
}

// InitiatePTW (Step 1) registers a high-voltage isolation request by a Unit Head.
func (s *WorkflowService) InitiatePTW(ctx context.Context, req InitiatePTWRequest) (*domain.SafetyPermitPTW, error) {
	if err := domain.ValidatePTWRequest(req.AssetID, req.UnitHeadID); err != nil {
		return nil, err
	}

	if s.gridRepo != nil {
		if _, err := s.gridRepo.GetGridAssetByID(ctx, req.AssetID); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrAssetNotFound, err)
		}
	}

	now := time.Now().UTC()
	permitNo := fmt.Sprintf("PTW-%d-%s", now.Year(), strings.ToUpper(crypto.MustUUID()[:8]))

	ptw := domain.SafetyPermitPTW{
		ID:          crypto.MustUUID(),
		PermitNo:    permitNo,
		WorkOrderID: req.WorkOrderID,
		AssetID:     req.AssetID,
		UnitHeadID:  req.UnitHeadID,
		Status:      domain.PTWRequested,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreatePTW(ctx, ptw); err != nil {
		return nil, fmt.Errorf("service: failed to persist PTW: %w", err)
	}

	_ = s.bus.Publish(ctx, "energy.ptw_requested", "s06_energy_management", map[string]interface{}{
		"permit_id":    ptw.ID,
		"permit_no":    ptw.PermitNo,
		"asset_id":     ptw.AssetID,
		"unit_head_id": ptw.UnitHeadID,
	})

	return &ptw, nil
}

// DualApprovePTW (Step 2) records digital signatures from Head of Services or HSE Safety Rep.
// Transitions to APPROVED_PENDING_ISOLATION once both parties have signed.
func (s *WorkflowService) DualApprovePTW(ctx context.Context, permitID string, approverRole string, approverID string) (*domain.SafetyPermitPTW, error) {
	ptw, err := s.repo.GetPTWByID(ctx, permitID)
	if err != nil {
		return nil, err
	}

	roleUpper := strings.ToUpper(strings.TrimSpace(approverRole))
	isDeptHead := roleUpper == "HEAD_SERVICES" || roleUpper == "DEPT_HEAD" || roleUpper == "HEAD_OF_SERVICES"
	isSafetyRep := roleUpper == "SAFETY_REPRESENTATIVE_HSE" || roleUpper == "SAFETY_REP" || roleUpper == "HSE_OFFICER"

	if !isDeptHead && !isSafetyRep {
		return nil, ErrUnauthorizedApprover
	}

	newStatus, err := domain.ValidateDualApproval(ptw, isDeptHead, isSafetyRep, approverID)
	if err != nil {
		return nil, err
	}
	ptw.Status = newStatus

	if err := s.repo.UpdatePTWStatus(ctx, *ptw); err != nil {
		return nil, fmt.Errorf("service: failed to update PTW status: %w", err)
	}

	eventName := "energy.ptw_partial_approved"
	if ptw.Status == domain.PTWApprovedPendingIsolation {
		eventName = "energy.ptw_dual_approved"
	}

	_ = s.bus.Publish(ctx, eventName, "s06_energy_management", map[string]interface{}{
		"permit_id":  ptw.ID,
		"permit_no":  ptw.PermitNo,
		"status":     ptw.Status,
		"dept_head":  ptw.DeptHeadApproved,
		"safety_rep": ptw.SafetyRepApproved,
	})

	return ptw, nil
}

// GrantSubstationIsolation (Step 3) records physical breaker racking out and line earthing by Substation Operator.
func (s *WorkflowService) GrantSubstationIsolation(ctx context.Context, permitID, operatorID, notes string) (*domain.SafetyPermitPTW, error) {
	ptw, err := s.repo.GetPTWByID(ctx, permitID)
	if err != nil {
		return nil, err
	}

	newStatus, err := domain.ValidateSubstationIsolation(ptw.Status, operatorID, notes)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	ptw.Status = newStatus
	ptw.SubstationOperatorID = &operatorID
	ptw.IsolationNotes = &notes
	ptw.IsolationGrantedAt = &now

	if err := s.repo.UpdatePTWStatus(ctx, *ptw); err != nil {
		return nil, fmt.Errorf("service: failed to persist isolation: %w", err)
	}

	// Update Feeder state to OUT_ON_PERMIT_PTW
	if feeder, err := s.repo.GetFeederStatusByAssetID(ctx, ptw.AssetID); err == nil && feeder != nil {
		feeder.State = domain.FeederOutOnPTW
		feeder.ActivePTWID = &ptw.ID
		narrative := fmt.Sprintf("Isolated under PTW %s by Operator %s: %s", ptw.PermitNo, operatorID, notes)
		feeder.OutageNarrative = &narrative
		_ = s.repo.UpsertFeederStatus(ctx, *feeder)
	}

	_ = s.bus.Publish(ctx, "energy.ptw_isolated_earthed", "s06_energy_management", map[string]interface{}{
		"permit_id":   ptw.ID,
		"permit_no":   ptw.PermitNo,
		"asset_id":    ptw.AssetID,
		"operator_id": operatorID,
	})

	return ptw, nil
}

// SignSafetyClearance (Step 5a) records Safety Rep certification that all work crews are clear and earthing removed.
func (s *WorkflowService) SignSafetyClearance(ctx context.Context, permitID, safetyRepID, notes string) (*domain.SafetyPermitPTW, error) {
	ptw, err := s.repo.GetPTWByID(ctx, permitID)
	if err != nil {
		return nil, err
	}

	newStatus, err := domain.ValidateSafetyClearance(ptw.Status, safetyRepID, notes)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	ptw.Status = newStatus
	ptw.SafetyRepID = &safetyRepID
	ptw.SafetyClearanceNotes = &notes
	ptw.SafetyClearanceSignedAt = &now

	if err := s.repo.UpdatePTWStatus(ctx, *ptw); err != nil {
		return nil, fmt.Errorf("service: failed to persist safety clearance: %w", err)
	}

	_ = s.bus.Publish(ctx, "energy.ptw_safety_cleared", "s06_energy_management", map[string]interface{}{
		"permit_id":     ptw.ID,
		"permit_no":     ptw.PermitNo,
		"safety_rep_id": safetyRepID,
	})

	return ptw, nil
}

// RestoreAndClosePTW (Step 5b) verifies safe clearance, closes breaker, and restores line to service.
func (s *WorkflowService) RestoreAndClosePTW(ctx context.Context, permitID, operatorID string) (*domain.SafetyPermitPTW, error) {
	ptw, err := s.repo.GetPTWByID(ctx, permitID)
	if err != nil {
		return nil, err
	}

	newStatus, err := domain.ValidateRestoration(ptw.Status, operatorID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	ptw.Status = newStatus
	ptw.RestoredAt = &now

	if err := s.repo.UpdatePTWStatus(ctx, *ptw); err != nil {
		return nil, fmt.Errorf("service: failed to restore PTW: %w", err)
	}

	// Revert Feeder state to IN_SERVICE
	if feeder, err := s.repo.GetFeederStatusByAssetID(ctx, ptw.AssetID); err == nil && feeder != nil {
		feeder.State = domain.FeederInService
		feeder.ActivePTWID = nil
		feeder.OutageNarrative = nil
		_ = s.repo.UpsertFeederStatus(ctx, *feeder)
	}

	_ = s.bus.Publish(ctx, "energy.ptw_restored_closed", "s06_energy_management", map[string]interface{}{
		"permit_id":   ptw.ID,
		"permit_no":   ptw.PermitNo,
		"asset_id":    ptw.AssetID,
		"operator_id": operatorID,
	})

	return ptw, nil
}

// GetPTWByID retrieves a permit by UUID.
func (s *WorkflowService) GetPTWByID(ctx context.Context, id string) (*domain.SafetyPermitPTW, error) {
	return s.repo.GetPTWByID(ctx, id)
}

// GetPTWByPermitNo retrieves a permit by human-readable permit number.
func (s *WorkflowService) GetPTWByPermitNo(ctx context.Context, permitNo string) (*domain.SafetyPermitPTW, error) {
	return s.repo.GetPTWByPermitNo(ctx, permitNo)
}

// ListPTWs returns safety permits filtered by status and asset ID.
func (s *WorkflowService) ListPTWs(ctx context.Context, status domain.PTWStatus, assetID string) ([]domain.SafetyPermitPTW, error) {
	return s.repo.ListPTWs(ctx, status, assetID)
}

