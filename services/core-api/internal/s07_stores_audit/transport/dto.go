// Package transport provides request/response data transfer objects and validation
// for S07 Stores, Material Requisitions, Reverse Logistics, and Staff Allowance Claims.
package transport

import (
	"errors"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
)

var (
	// ErrMissingMandatoryFields indicates one or more required DTO fields are empty.
	ErrMissingMandatoryFields = errors.New("transport: mandatory fields missing")
	// ErrInvalidItemCategory indicates an unapproved inventory category.
	ErrInvalidItemCategory = errors.New("transport: invalid item category")
	// ErrInvalidClaimType indicates an unapproved allowance claim type.
	ErrInvalidClaimType = errors.New("transport: invalid allowance claim type")
	// ErrNegativeAmount indicates price, cost, or allowance amount is invalid.
	ErrNegativeAmount = errors.New("transport: amount or cost cannot be negative or zero")
)

// CreateStoreItemDTO defines payload for POST /api/v1/stores/items.
type CreateStoreItemDTO struct {
	ItemCode          string              `json:"item_code"`
	ItemName          string              `json:"item_name"`
	Category          domain.ItemCategory `json:"category"`
	QuantityAvailable float64             `json:"quantity_available"`
	ReorderThreshold  float64             `json:"reorder_threshold"`
	UnitOfMeasure     string              `json:"unit_of_measure"`
	UnitCost          float64             `json:"unit_cost"`
	BinLocation       string              `json:"bin_location"`
}

// Validate verifies store item parameters.
func (d CreateStoreItemDTO) Validate() error {
	if strings.TrimSpace(d.ItemCode) == "" || strings.TrimSpace(d.ItemName) == "" || strings.TrimSpace(d.BinLocation) == "" {
		return ErrMissingMandatoryFields
	}
	if d.Category == "" {
		return ErrInvalidItemCategory
	}
	if d.QuantityAvailable < 0.0 || d.ReorderThreshold < 0.0 || d.UnitCost < 0.0 {
		return ErrNegativeAmount
	}
	return nil
}

// UpdateStockDTO defines payload for POST /api/v1/stores/items/{id}/stock.
type UpdateStockDTO struct {
	Delta float64 `json:"delta"`
}

// RequisitionItemDTO represents a single requested SKU inside a requisition.
type RequisitionItemDTO struct {
	ItemCode          string  `json:"item_code"`
	QuantityRequested float64 `json:"quantity_requested"`
	UnitCost          float64 `json:"unit_cost"`
}

// InitiateRequisitionDTO defines payload for POST /api/v1/stores/requisitions.
type InitiateRequisitionDTO struct {
	OriginDeptID        string               `json:"origin_dept_id"`
	LinkedTicketOrJobID *string              `json:"linked_ticket_or_job_id,omitempty"`
	Items               []RequisitionItemDTO `json:"items"`
}

// Validate verifies requisition submission.
func (d InitiateRequisitionDTO) Validate() error {
	if strings.TrimSpace(d.OriginDeptID) == "" || len(d.Items) == 0 {
		return ErrMissingMandatoryFields
	}
	for _, it := range d.Items {
		if strings.TrimSpace(it.ItemCode) == "" || it.QuantityRequested <= 0.0 {
			return ErrMissingMandatoryFields
		}
	}
	return nil
}

// AuditConsentDTO defines payload for POST /api/v1/stores/requisitions/{id}/audit-consent.
type AuditConsentDTO struct {
	Consent bool   `json:"consent"`
	Notes   string `json:"notes"`
}

// MDApprovalDTO defines payload for POST /api/v1/stores/requisitions/{id}/md-approve.
type MDApprovalDTO struct {
	Notes string `json:"notes"`
}

// LogScrapReturnDTO defines payload for POST /api/v1/stores/scrap-returns.
type LogScrapReturnDTO struct {
	RequisitionID       string `json:"requisition_id"`
	ItemsReturnedJSON   string `json:"items_returned_json"`
	ScrapConditionNotes string `json:"scrap_condition_notes"`
}

// Validate verifies scrap return parameters.
func (d LogScrapReturnDTO) Validate() error {
	if strings.TrimSpace(d.RequisitionID) == "" || strings.TrimSpace(d.ScrapConditionNotes) == "" {
		return ErrMissingMandatoryFields
	}
	return nil
}

// SubmitClaimDTO defines payload for POST /api/v1/audit/claims.
type SubmitClaimDTO struct {
	ClaimType           domain.ClaimType `json:"claim_type"`
	Amount              float64          `json:"amount"`
	ClaimDate           string           `json:"claim_date"`
	ShiftWindow         *string          `json:"shift_window,omitempty"`
	LinkedPTWOrTicketID *string          `json:"linked_ptw_or_ticket_id,omitempty"`
	RosterMatchVerified bool             `json:"roster_match_verified"`
	PTWCrewVerified     bool             `json:"ptw_crew_verified"`
	GPSMatchVerified    bool             `json:"gps_match_verified"`
}

// Validate verifies staff financial allowance claim parameters.
func (d SubmitClaimDTO) Validate() error {
	if d.ClaimType == "" {
		return ErrInvalidClaimType
	}
	if d.Amount <= 0.0 {
		return ErrNegativeAmount
	}
	if strings.TrimSpace(d.ClaimDate) == "" {
		return ErrMissingMandatoryFields
	}
	return nil
}

// ReviewClaimDTO defines payload for POST /api/v1/audit/claims/{id}/review.
type ReviewClaimDTO struct {
	Approve bool   `json:"approve"`
	Notes   string `json:"notes"`
}

