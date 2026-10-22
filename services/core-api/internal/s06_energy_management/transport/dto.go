// Package transport provides request/response data transfer objects and validation
// for S06 Energy Management, load shedding rationing, and Safety Permit-to-Work (PTW) protocols.
package transport

import (
	"errors"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
)

var (
	// ErrMissingMandatoryFields indicates one or more required input fields are empty.
	ErrMissingMandatoryFields = errors.New("transport: mandatory fields missing")
	// ErrInvalidShiftWindow indicates an unapproved shift window enumeration.
	ErrInvalidShiftWindow = errors.New("transport: shift_window must be MORNING, AFTERNOON, or NIGHT")
	// ErrNegativeMW indicates power value cannot be negative.
	ErrNegativeMW = errors.New("transport: generation or TCN power cannot be negative")
	// ErrInvalidRole indicates approver role is not authorized for PTW signoff.
	ErrInvalidRole = errors.New("transport: approver_role must be HEAD_SERVICES or SAFETY_REP")
	// ErrNotesTooShort indicates safety isolation or clearance proof notes are too short.
	ErrNotesTooShort = errors.New("transport: confirmation notes must be at least 10 characters")
)

// RecordBulkIntakeDTO defines the payload for POST /api/v1/energy/bulk-intake.
type RecordBulkIntakeDTO struct {
	ShiftID            string             `json:"shift_id"`
	ShiftWindow        domain.ShiftWindow `json:"shift_window"`
	GenerationMWIn     float64            `json:"generation_mw_in"`
	TCNMWIn            float64            `json:"tcn_mw_in"`
	SubstationRollcall string             `json:"substation_rollcall"`
	HandoverNotes      string             `json:"handover_notes"`
}

// Validate verifies bulk intake submission parameters.
func (d RecordBulkIntakeDTO) Validate() error {
	if strings.TrimSpace(d.ShiftID) == "" {
		return ErrMissingMandatoryFields
	}
	if d.ShiftWindow != domain.ShiftMorning && d.ShiftWindow != domain.ShiftAfternoon && d.ShiftWindow != domain.ShiftNight {
		return ErrInvalidShiftWindow
	}
	if d.GenerationMWIn < 0.0 || d.TCNMWIn < 0.0 {
		return ErrNegativeMW
	}
	return nil
}

// RationingCalculateDTO defines the payload for POST /api/v1/energy/rationing/calculate.
type RationingCalculateDTO struct {
	AvailableBulkMW float64 `json:"available_bulk_mw"`
}

// Validate verifies rationing parameters.
func (d RationingCalculateDTO) Validate() error {
	if d.AvailableBulkMW < 0.0 {
		return ErrNegativeMW
	}
	return nil
}

// InitiatePTWDTO defines the payload for Step 1: POST /api/v1/energy/ptw/request.
type InitiatePTWDTO struct {
	AssetID         string  `json:"asset_id"`
	WorkOrderID     *string `json:"work_order_id,omitempty"`
	WorkDescription string  `json:"work_description"`
}

// Validate verifies PTW initiation fields.
func (d InitiatePTWDTO) Validate() error {
	if strings.TrimSpace(d.AssetID) == "" || strings.TrimSpace(d.WorkDescription) == "" {
		return ErrMissingMandatoryFields
	}
	return nil
}

// ApprovePTWDTO defines the payload for Step 2: POST /api/v1/energy/ptw/{id}/approve.
type ApprovePTWDTO struct {
	ApproverRole string `json:"approver_role"`
}

// Validate verifies dual approval role.
func (d ApprovePTWDTO) Validate() error {
	r := strings.ToUpper(strings.TrimSpace(d.ApproverRole))
	if r != "HEAD_SERVICES" && r != "DEPT_HEAD" && r != "SAFETY_REP" && r != "SAFETY_REPRESENTATIVE_HSE" {
		return ErrInvalidRole
	}
	return nil
}

// IsolatePTWDTO defines the payload for Step 3: POST /api/v1/energy/ptw/{id}/isolate.
type IsolatePTWDTO struct {
	IsolationNotes string `json:"isolation_notes"`
}

// Validate verifies substation operator isolation declaration.
func (d IsolatePTWDTO) Validate() error {
	if len(strings.TrimSpace(d.IsolationNotes)) < 10 {
		return ErrNotesTooShort
	}
	return nil
}

// ClearSafetyPTWDTO defines the payload for Step 5a: POST /api/v1/energy/ptw/{id}/clear-safety.
type ClearSafetyPTWDTO struct {
	Notes string `json:"notes"`
}

// Validate verifies safety rep clearance declaration.
func (d ClearSafetyPTWDTO) Validate() error {
	if len(strings.TrimSpace(d.Notes)) < 10 {
		return ErrNotesTooShort
	}
	return nil
}

