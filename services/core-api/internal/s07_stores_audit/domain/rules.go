// Package domain implements business validation rules, dual-path procurement state transitions,
// and automated 3-way staff claims verification invariants for S07 Stores & Internal Audit.
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrInvalidItemData indicates mandatory inventory fields or quantities are invalid.
	ErrInvalidItemData = errors.New("domain: invalid store item parameters")
	// ErrEmptyRequisitionItems indicates a material request contains no item entries.
	ErrEmptyRequisitionItems = errors.New("domain: material requisition must contain at least one item")
	// ErrInvalidQuantity indicates requested or available quantity is invalid.
	ErrInvalidQuantity = errors.New("domain: quantity must be greater than zero")
	// ErrIllegalRequisitionTransition indicates an unauthorized or invalid requisition status progression.
	ErrIllegalRequisitionTransition = errors.New("domain: illegal material requisition state transition")
	// ErrAuditConsentRequired indicates procurement or issuance cannot proceed without prior Audit Consent.
	ErrAuditConsentRequired = errors.New("domain: internal audit consent is mandatory before MD procurement or store issuance")
	// ErrUnreturnedScrapFound indicates pending unreturned scrap prevents new material issuance.
	ErrUnreturnedScrapFound = errors.New("domain: unreturned scrap from prior work order blocks new issuance")
	// ErrInvalidClaimData indicates allowance claim fields are invalid.
	ErrInvalidClaimData = errors.New("domain: invalid staff allowance claim parameters")
)

// ValidateStoreItem verifies SKU master data sanity.
func ValidateStoreItem(item StoreItem) error {
	if strings.TrimSpace(item.ItemCode) == "" || strings.TrimSpace(item.ItemName) == "" {
		return ErrInvalidItemData
	}
	if item.Category == "" {
		return ErrInvalidItemData
	}
	if item.QuantityAvailable < 0.0 || item.ReorderThreshold < 0.0 || item.UnitCost < 0.0 {
		return ErrInvalidQuantity
	}
	return nil
}

// EvaluateReorderStatus determines if an inventory item has breached its reorder threshold.
func EvaluateReorderStatus(quantity, threshold float64) bool {
	return quantity <= threshold
}

// CalculateRequisitionTotal computes line item costs and verifies item requests.
func CalculateRequisitionTotal(items []RequisitionItem) (float64, error) {
	if len(items) == 0 {
		return 0.0, ErrEmptyRequisitionItems
	}
	total := 0.0
	for _, it := range items {
		if it.QuantityRequested <= 0.0 {
			return 0.0, fmt.Errorf("%w: item %s quantity must be > 0", ErrInvalidQuantity, it.ItemCode)
		}
		if it.UnitCost < 0.0 {
			return 0.0, fmt.Errorf("%w: item %s unit cost cannot be negative", ErrInvalidItemData, it.ItemCode)
		}
		total += it.QuantityRequested * it.UnitCost
	}
	return total, nil
}

// ValidateDeptApproval (Step 1 -> 2) advances requisition upon Unit/Dept Head or Generation GM sign-off.
func ValidateDeptApproval(req *MaterialRequisition, deptHeadID string) error {
	if req.Status != StatusRequisitionSubmitted {
		return fmt.Errorf("%w: cannot approve requisition in status %s", ErrIllegalRequisitionTransition, req.Status)
	}
	if strings.TrimSpace(deptHeadID) == "" {
		return fmt.Errorf("%w: departmental head ID is required", ErrIllegalRequisitionTransition)
	}

	now := time.Now().UTC()
	req.DeptHeadID = &deptHeadID
	req.DeptOrGMApproved = true
	req.DeptApprovedAt = &now

	if req.InStockFlag {
		req.Status = StatusRequisitionDeptApproved
	} else {
		req.Status = StatusRequisitionMDProcurementPending
	}
	return nil
}

// ValidateAuditConsent (Step 2 -> 3) records compliance check and Audit Consent from Internal Audit.
func ValidateAuditConsent(req *MaterialRequisition, auditorID string, consent bool, notes string) error {
	if req.Status != StatusRequisitionDeptApproved && req.Status != StatusRequisitionMDProcurementPending {
		return fmt.Errorf("%w: cannot perform audit on requisition in status %s", ErrIllegalRequisitionTransition, req.Status)
	}
	if strings.TrimSpace(auditorID) == "" {
		return fmt.Errorf("%w: auditor ID is required", ErrIllegalRequisitionTransition)
	}

	now := time.Now().UTC()
	req.AuditOfficerID = &auditorID
	req.AuditConsentAt = &now
	req.AuditNotes = &notes

	if !consent {
		req.AuditConsentStatus = AuditConsentQueriedRejected
		req.Status = StatusRequisitionRejected
		return nil
	}

	req.AuditConsentStatus = AuditConsentConsented
	if req.InStockFlag {
		req.Status = StatusRequisitionAuditConsented
	}
	// If out-of-stock, status remains MD_PROCUREMENT_PENDING awaiting MD approval
	return nil
}

// ValidateMDProcurementApproval (Step 3b) records final approval by Managing Director for supplier purchase.
func ValidateMDProcurementApproval(req *MaterialRequisition, mdNotes string) error {
	if req.InStockFlag {
		return fmt.Errorf("%w: in-stock items do not require MD procurement approval", ErrIllegalRequisitionTransition)
	}
	if req.AuditConsentStatus != AuditConsentConsented {
		return ErrAuditConsentRequired
	}
	if req.Status != StatusRequisitionMDProcurementPending {
		return fmt.Errorf("%w: cannot grant MD approval on status %s", ErrIllegalRequisitionTransition, req.Status)
	}

	now := time.Now().UTC()
	req.MDProcurementApproved = true
	req.MDApprovedAt = &now
	req.MDNotes = &mdNotes
	req.Status = StatusRequisitionMDApproved
	return nil
}

// ValidateMaterialIssuance records store dispatch, updates status to ISSUED, and assigns a waybill.
func ValidateMaterialIssuance(req *MaterialRequisition, storeOfficerID, waybillNo string) error {
	if req.InStockFlag {
		if req.Status != StatusRequisitionAuditConsented {
			return fmt.Errorf("%w: in-stock item requires AUDIT_CONSENTED status prior to issuance (current: %s)", ErrIllegalRequisitionTransition, req.Status)
		}
	} else {
		if req.Status != StatusRequisitionMDApproved {
			return fmt.Errorf("%w: out-of-stock item requires MD_APPROVED status prior to issuance (current: %s)", ErrIllegalRequisitionTransition, req.Status)
		}
	}

	if strings.TrimSpace(storeOfficerID) == "" || strings.TrimSpace(waybillNo) == "" {
		return fmt.Errorf("%w: store_officer_id and waybill_no are mandatory", ErrIllegalRequisitionTransition)
	}

	now := time.Now().UTC()
	req.StoreOfficerID = &storeOfficerID
	req.WaybillNo = &waybillNo
	req.IssuedAt = &now
	req.Status = StatusRequisitionIssued
	return nil
}

// Evaluate3WayClaimVerification conducts automated cross-checking across Attendance, PTW, and GPS.
func Evaluate3WayClaimVerification(claim StaffAllowanceClaim, rosterMatch, ptwMatch, gpsMatch bool) (AIAuditFlag, float64, ClaimStatus) {
	// If staff member was not logged in the shift roster/call sheet, flag immediately as high-risk fraud
	if !rosterMatch {
		return AIFlagFlaggedFraud, 0.950, StatusClaimAIFlagged
	}

	// For Overtime and Hazard, both PTW and GPS verification are required if linked to a PTW
	if claim.LinkedPTWOrTicketID != nil && *claim.LinkedPTWOrTicketID != "" {
		if !ptwMatch || !gpsMatch {
			return AIFlagFlaggedMismatch, 0.720, StatusClaimAIFlagged
		}
	} else {
		// Non-PTW general field claim (e.g. Travel/Meal) requires GPS alignment
		if !gpsMatch && claim.ClaimType == ClaimTravelOutOfStation {
			return AIFlagFlaggedMismatch, 0.680, StatusClaimAIFlagged
		}
	}

	// All relevant sources verified cleanly
	return AIFlagVerified, 0.050, StatusClaimAIVerified
}

// ValidateAuditorClaimDecision records the internal auditor's formal approval or rejection of an allowance claim.
func ValidateAuditorClaimDecision(claim *StaffAllowanceClaim, auditorID string, approve bool, notes string) error {
	if claim.Status != StatusClaimSubmitted && claim.Status != StatusClaimAIVerified && claim.Status != StatusClaimAIFlagged {
		return fmt.Errorf("domain: cannot decide on claim in status %s", claim.Status)
	}
	if strings.TrimSpace(auditorID) == "" {
		return fmt.Errorf("domain: auditor ID is required")
	}

	now := time.Now().UTC()
	claim.AuditorID = &auditorID
	claim.AuditorNotes = &notes
	claim.AuditDecisionAt = &now

	if approve {
		claim.Status = StatusClaimAuditApproved
	} else {
		claim.Status = StatusClaimAuditRejected
	}
	return nil
}

