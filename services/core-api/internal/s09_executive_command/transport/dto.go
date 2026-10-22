// Package transport provides request/response DTOs and input validations
// for S09 Executive Command, C-Suite Decision Console & Predictive Maintenance.
package transport

import (
	"errors"
	"fmt"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

var (
	// ErrMissingMandatoryFields indicates one or more required DTO fields are empty.
	ErrMissingMandatoryFields = errors.New("transport: mandatory fields missing")
	// ErrInvalidRole indicates an unapproved C-Suite executive role enumeration.
	ErrInvalidRole = errors.New("transport: invalid executive role")
	// ErrInvalidActionType indicates an unrecognized C-Suite action category.
	ErrInvalidActionType = errors.New("transport: invalid action type")
)

// CaptureSnapshotDTO defines payload for POST /api/v1/executive/snapshots.
type CaptureSnapshotDTO struct {
	SnapshotDate                  string  `json:"snapshot_date"`
	ShiftWindow                   string  `json:"shift_window"`
	TotalGenerationMW             float64 `json:"total_generation_mw"`
	TCNGridIntakeMW               float64 `json:"tcn_grid_intake_mw"`
	TotalSystemLoadMW             float64 `json:"total_system_load_mw"`
	FeedersInServiceCount         int     `json:"feeders_in_service_count"`
	FeedersOutFaultCount          int     `json:"feeders_out_fault_count"`
	FeedersOutPTWCount            int     `json:"feeders_out_ptw_count"`
	FeedersOutRationingCount      int     `json:"feeders_out_rationing_count"`
	TotalRevenueCollected         float64 `json:"total_revenue_collected"`
	PrepaidRevenueSharePct        float64 `json:"prepaid_revenue_share_pct"`
	UnresolvedTheftAlertsCount    int     `json:"unresolved_theft_alerts_count"`
	UnresolvedChurnAlertsCount    int     `json:"unresolved_churn_alerts_count"`
	AuditPendingProcurementsCount int     `json:"audit_pending_procurements_count"`
	UnreturnedScrapCount          int     `json:"unreturned_scrap_count"`
	Notes                         string  `json:"notes,omitempty"`
	CapturedByStaffID             string  `json:"captured_by_staff_id"`
}

// Validate verifies mandatory snapshot parameters.
func (d CaptureSnapshotDTO) Validate() error {
	if strings.TrimSpace(d.SnapshotDate) == "" || strings.TrimSpace(d.ShiftWindow) == "" || strings.TrimSpace(d.CapturedByStaffID) == "" {
		return fmt.Errorf("%w: snapshot_date, shift_window, and captured_by_staff_id are required", ErrMissingMandatoryFields)
	}
	return nil
}

// EnqueueActionDTO defines payload for POST /api/v1/executive/actions.
type EnqueueActionDTO struct {
	ActionType       domain.ActionItemType `json:"action_type"`
	TargetRole       domain.ExecutiveRole  `json:"target_role"`
	SourceEntityID   string                `json:"source_entity_id"`
	SourceEntityRef  string                `json:"source_entity_ref"`
	SourceDepartment string                `json:"source_department"`
	Summary          string                `json:"summary"`
	Amount           float64               `json:"amount"`
	Priority         string                `json:"priority"`
}

// Validate verifies executive action console enqueuing.
func (d EnqueueActionDTO) Validate() error {
	if strings.TrimSpace(d.SourceEntityID) == "" || strings.TrimSpace(d.SourceEntityRef) == "" || strings.TrimSpace(d.Summary) == "" {
		return fmt.Errorf("%w: source_entity_id, source_entity_ref, and summary are required", ErrMissingMandatoryFields)
	}
	if d.ActionType == "" {
		return ErrInvalidActionType
	}
	if d.TargetRole == "" {
		return ErrInvalidRole
	}
	return nil
}

// DecideActionDTO defines payload for POST /api/v1/executive/actions/{id}/decide.
type DecideActionDTO struct {
	DecidedByStaffID string `json:"decided_by_staff_id"`
	Approve          bool   `json:"approve"`
	Notes            string `json:"notes"`
}

// Validate verifies executive digital decision.
func (d DecideActionDTO) Validate() error {
	if strings.TrimSpace(d.DecidedByStaffID) == "" || strings.TrimSpace(d.Notes) == "" {
		return fmt.Errorf("%w: decided_by_staff_id and notes are required", ErrMissingMandatoryFields)
	}
	return nil
}

// HydroPrognosisDTO defines payload for POST /api/v1/executive/predictive/hydro.
type HydroPrognosisDTO struct {
	AssetID           string  `json:"asset_id"`
	InflowM3s         float64 `json:"inflow_m3s"`
	DischargeM3s      float64 `json:"discharge_m3s"`
	CurrentHeadM      float64 `json:"current_head_m"`
	TurbineEfficiency float64 `json:"turbine_efficiency"`
	ForecastDays      int     `json:"forecast_days"`
}

// Validate verifies hydro input parameters.
func (d HydroPrognosisDTO) Validate() error {
	if strings.TrimSpace(d.AssetID) == "" {
		return fmt.Errorf("%w: asset_id is required", ErrMissingMandatoryFields)
	}
	return nil
}

// GeneratorPrognosisDTO defines payload for POST /api/v1/executive/predictive/generator.
type GeneratorPrognosisDTO struct {
	AssetID               string  `json:"asset_id"`
	BearingVibrationMmS   float64 `json:"bearing_vibration_mm_s"`
	WindingTempC          float64 `json:"winding_temp_c"`
	LubeOilPressureBar    float64 `json:"lube_oil_pressure_bar"`
	RunHoursSinceOverhaul float64 `json:"run_hours_since_overhaul"`
	UnitHeadID            string  `json:"unit_head_id"`
	TeamLeadID            string  `json:"team_lead_id"`
}

// Validate verifies generator input parameters.
func (d GeneratorPrognosisDTO) Validate() error {
	if strings.TrimSpace(d.AssetID) == "" {
		return fmt.Errorf("%w: asset_id is required", ErrMissingMandatoryFields)
	}
	return nil
}

// TransformerPrognosisDTO defines payload for POST /api/v1/executive/predictive/transformer.
type TransformerPrognosisDTO struct {
	AssetID         string  `json:"asset_id"`
	LoadPct         float64 `json:"load_pct"`
	TopOilTempC     float64 `json:"top_oil_temp_c"`
	FuseBlows30Days int     `json:"fuse_blows_30_days"`
	UnitHeadID      string  `json:"unit_head_id"`
	TeamLeadID      string  `json:"team_lead_id"`
}

// Validate verifies transformer CBM input parameters.
func (d TransformerPrognosisDTO) Validate() error {
	if strings.TrimSpace(d.AssetID) == "" {
		return fmt.Errorf("%w: asset_id is required", ErrMissingMandatoryFields)
	}
	return nil
}

// FeederPrognosisDTO defines payload for POST /api/v1/executive/predictive/feeder.
type FeederPrognosisDTO struct {
	AssetID             string  `json:"asset_id"`
	TripCount30Days     int     `json:"trip_count_30_days"`
	FaultDensityPerKm   float64 `json:"fault_density_per_km"`
	DaysSinceLastPatrol float64 `json:"days_since_last_patrol"`
	StormForecast       bool    `json:"storm_forecast"`
	UnitHeadID          string  `json:"unit_head_id"`
	TeamLeadID          string  `json:"team_lead_id"`
}

// Validate verifies feeder prognostic input parameters.
func (d FeederPrognosisDTO) Validate() error {
	if strings.TrimSpace(d.AssetID) == "" {
		return fmt.Errorf("%w: asset_id is required", ErrMissingMandatoryFields)
	}
	return nil
}

// AcknowledgeTriggerDTO defines payload for POST /api/v1/executive/predictive/triggers/{id}/ack.
type AcknowledgeTriggerDTO struct {
	StaffID string `json:"staff_id"`
	Status  string `json:"status"`
}

// Validate verifies trigger acknowledgement.
func (d AcknowledgeTriggerDTO) Validate() error {
	if strings.TrimSpace(d.StaffID) == "" {
		return fmt.Errorf("%w: staff_id is required", ErrMissingMandatoryFields)
	}
	return nil
}

