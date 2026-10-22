// Package transport provides request/response data transfer objects and validation
// for S05 Generation Department and Master Station Hub HTTP endpoints.
package transport

import (
	"errors"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
)

// RecordTelemetryDTO defines payload for POST /api/v1/generation/telemetry.
type RecordTelemetryDTO struct {
	ShiftID             string              `json:"shift_id"`
	ShiftWindow         domain.ShiftWindow  `json:"shift_window"`
	PlantID             string              `json:"plant_id"`
	PlantState          domain.PlantState   `json:"plant_state"`
	OutageReason        *domain.OutageReason `json:"outage_reason,omitempty"`
	OutageDetails       *string             `json:"outage_details,omitempty"`
	WaterLevelM         float64             `json:"water_level_m"`
	HeadwaterElevationM float64             `json:"headwater_elevation_m"`
	TailwaterElevationM float64             `json:"tailwater_elevation_m"`
	InflowRateM3S       float64             `json:"inflow_rate_m3_s"`
	MWGenerated         float64             `json:"mw_generated"`
	MVArReactive        float64             `json:"mvar_reactive"`
	StatorWindingTempC  float64             `json:"stator_winding_temp_c"`
	BearingVibrationMMS float64             `json:"bearing_vibration_mm_s"`
	LubeOilPressureBar  float64             `json:"lube_oil_pressure_bar"`
}

// Validate checks telemetry payload sanity.
func (d RecordTelemetryDTO) Validate() error {
	if strings.TrimSpace(d.ShiftID) == "" || strings.TrimSpace(d.PlantID) == "" {
		return errors.New("transport: shift_id and plant_id are mandatory")
	}
	if d.PlantState != domain.StateInService && d.PlantState != domain.StateStandby && d.PlantState != domain.StateOutOfService {
		return errors.New("transport: invalid plant_state")
	}
	if d.WaterLevelM < 0.0 || d.WaterLevelM > 250.0 {
		return errors.New("transport: water_level_m must be between 0.0 and 250.0 meters")
	}
	if d.MWGenerated < 0.0 {
		return errors.New("transport: mw_generated cannot be negative")
	}
	return nil
}

// RecordRollcallDTO defines payload for POST /api/v1/generation/rollcall.
type RecordRollcallDTO struct {
	ShiftID             string             `json:"shift_id"`
	ShiftWindow         domain.ShiftWindow `json:"shift_window"`
	PlantID             string             `json:"plant_id"`
	OperatorNames       []string           `json:"operator_names"`
	ProtectionEngineers []string           `json:"protection_engineers"`
	Technicians         []string           `json:"technicians"`
	HandoverNotes       string             `json:"handover_notes"`
}

// Validate checks roll-call payload fields.
func (d RecordRollcallDTO) Validate() error {
	if strings.TrimSpace(d.ShiftID) == "" || strings.TrimSpace(d.PlantID) == "" {
		return errors.New("transport: shift_id and plant_id are mandatory")
	}
	if len(d.OperatorNames) == 0 || strings.TrimSpace(d.OperatorNames[0]) == "" {
		return errors.New("transport: at least one active operator name is mandatory")
	}
	if len(strings.TrimSpace(d.HandoverNotes)) < 10 {
		return errors.New("transport: handover_notes must be at least 10 characters")
	}
	return nil
}

// ReportFaultJobDTO defines payload for POST /api/v1/generation/faults.
type ReportFaultJobDTO struct {
	PlantID          string               `json:"plant_id"`
	Wing             domain.TechnicalWing `json:"wing"`
	FaultDescription string               `json:"fault_description"`
	Severity         string               `json:"severity"`
	RequiresWorkshop bool                 `json:"requires_workshop"`
}

// Validate checks fault reporting parameters.
func (d ReportFaultJobDTO) Validate() error {
	if strings.TrimSpace(d.PlantID) == "" {
		return errors.New("transport: plant_id is mandatory")
	}
	if d.Wing != domain.WingElectrical && d.Wing != domain.WingMechanical {
		return errors.New("transport: wing must be ELECTRICAL or MECHANICAL")
	}
	if len(strings.TrimSpace(d.FaultDescription)) < 15 {
		return errors.New("transport: fault_description must be at least 15 characters")
	}
	return nil
}

// GMApprovalDTO defines payload for POST /api/v1/generation/faults/{id}/gm-approval.
type GMApprovalDTO struct {
	Approved bool   `json:"approved"`
	Notes    string `json:"notes"`
}

// Validate checks GM approval/rejection notes.
func (d GMApprovalDTO) Validate() error {
	if !d.Approved && len(strings.TrimSpace(d.Notes)) < 10 {
		return errors.New("transport: rejection explanation notes must be at least 10 characters")
	}
	return nil
}

// UpdateWorkshopDTO defines payload for POST /api/v1/generation/faults/{id}/workshop.
type UpdateWorkshopDTO struct {
	WorkshopStatus      domain.WorkshopStatus `json:"workshop_status"`
	StoreRequisitionRef *string               `json:"store_requisition_ref,omitempty"`
}

// Validate checks workshop transition status.
func (d UpdateWorkshopDTO) Validate() error {
	if strings.TrimSpace(string(d.WorkshopStatus)) == "" {
		return errors.New("transport: workshop_status is mandatory")
	}
	return nil
}

