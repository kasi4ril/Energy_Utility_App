// Package transport provides request/response data transfer objects and validation
// for S04 Services & Maintenance HTTP endpoints.
package transport

import (
	"errors"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
)

// CreateWorkOrderDTO defines the payload for POST /api/v1/services/work-orders.
type CreateWorkOrderDTO struct {
	TicketID                *string                  `json:"ticket_id,omitempty"`
	UnitType                domain.UnitType          `json:"unit_type"`
	WorkOrderType           domain.WorkOrderType     `json:"work_order_type"`
	Priority                domain.WorkOrderPriority `json:"priority"`
	TeamLeadID              string                   `json:"team_lead_id"`
	AssetID                 string                   `json:"asset_id"`
	PTWRequired             bool                     `json:"ptw_required"`
	PTWPermitID             *string                  `json:"ptw_permit_id,omitempty"`
	MaterialsRequisitionRef *string                  `json:"materials_requisition_ref,omitempty"`
	Instructions            string                   `json:"instructions"`
}

// Validate checks mandatory fields for work order dispatch.
func (d CreateWorkOrderDTO) Validate() error {
	if strings.TrimSpace(string(d.UnitType)) == "" || strings.TrimSpace(string(d.WorkOrderType)) == "" {
		return errors.New("transport: unit_type and work_order_type are mandatory")
	}
	if strings.TrimSpace(d.TeamLeadID) == "" || strings.TrimSpace(d.AssetID) == "" {
		return errors.New("transport: team_lead_id and asset_id are mandatory")
	}
	if len(strings.TrimSpace(d.Instructions)) < 10 {
		return errors.New("transport: instructions must be at least 10 characters")
	}
	if d.PTWRequired && (d.PTWPermitID == nil || strings.TrimSpace(*d.PTWPermitID) == "") {
		return errors.New("transport: ptw_permit_id is required when ptw_required is true")
	}
	return nil
}

// SubmitProofDTO defines the payload for POST /api/v1/services/work-orders/{id}/proof.
type SubmitProofDTO struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	BeforePhotoURL  string  `json:"before_photo_url"`
	AfterPhotoURL   string  `json:"after_photo_url"`
	PhotoHashBefore string  `json:"photo_hash_before"`
	PhotoHashAfter  string  `json:"photo_hash_after"`
	MeterSealNumber *string `json:"meter_seal_number,omitempty"`
	WorkDetails     string  `json:"work_details"`
}

// Validate checks the mobile PWA proof-of-work submission.
func (d SubmitProofDTO) Validate() error {
	if d.Latitude < -90 || d.Latitude > 90 || d.Longitude < -180 || d.Longitude > 180 {
		return errors.New("transport: valid GPS latitude (-90 to 90) and longitude (-180 to 180) are required")
	}
	if strings.TrimSpace(d.BeforePhotoURL) == "" || strings.TrimSpace(d.AfterPhotoURL) == "" {
		return errors.New("transport: before_photo_url and after_photo_url are mandatory")
	}
	if len(strings.TrimSpace(d.PhotoHashBefore)) < 32 || len(strings.TrimSpace(d.PhotoHashAfter)) < 32 {
		return errors.New("transport: before and after SHA-256 photo hashes are mandatory (min 32 characters)")
	}
	if len(strings.TrimSpace(d.WorkDetails)) < 15 {
		return errors.New("transport: work_details must be at least 15 characters of technical summary")
	}
	return nil
}

// VerifyProofDTO defines the payload for POST /api/v1/services/work-orders/{id}/verify.
type VerifyProofDTO struct {
	Approved bool   `json:"approved"`
	Notes    string `json:"notes"`
}

// Validate verifies Unit Head approval/rejection parameters.
func (d VerifyProofDTO) Validate() error {
	if !d.Approved && len(strings.TrimSpace(d.Notes)) < 10 {
		return errors.New("transport: rework notes must be at least 10 characters when rejecting proof")
	}
	return nil
}

