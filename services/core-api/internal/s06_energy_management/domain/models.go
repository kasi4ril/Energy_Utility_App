// Package domain encapsulates core entities, state machine enumerations,
// and pure business logic invariants for the S06 Energy Management bounded context.
package domain

import (
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

// FeederState represents the operational condition of a 33kV or 11kV transmission/distribution line.
type FeederState string

const (
	// FeederInService indicates the feeder breaker is closed and carrying live electrical load.
	FeederInService FeederState = "IN_SERVICE"
	// FeederOutOnFault indicates the feeder tripped automatically on protective relay actuation.
	FeederOutOnFault FeederState = "OUT_ON_FAULT"
	// FeederOutOnPTW indicates the line is physically isolated, locked out, and earthed under a Safety PTW.
	FeederOutOnPTW FeederState = "OUT_ON_PERMIT_PTW"
	// FeederOutOnLoadShedding indicates the line was opened by Central Dispatch to balance demand with bulk intake.
	FeederOutOnLoadShedding FeederState = "OUT_ON_LOAD_SHEDDING"
)

// PTWStatus governs the sequential progression of the 5-Step Safety Permit-to-Work protocol.
type PTWStatus string

const (
	// PTWRequested indicates Step 1: Unit Head raised the permit request.
	PTWRequested PTWStatus = "REQUESTED"
	// PTWApprovedPendingIsolation indicates Step 2: Head of Services AND HSE Safety Rep both approved.
	PTWApprovedPendingIsolation PTWStatus = "APPROVED_PENDING_ISOLATION"
	// PTWIsolatedEarthed indicates Step 3: Substation Operator isolated breaker, racked out, and earthed.
	PTWIsolatedEarthed PTWStatus = "ISOLATED_EARTHED"
	// PTWSafetyCleared indicates Step 5a: Safety Rep signed "Crew Clear & Safe to Energize".
	PTWSafetyCleared PTWStatus = "SAFETY_CLEARED"
	// PTWRestoredClosed indicates Step 5b: Substation Operator removed earths, closed breaker, and restored line.
	PTWRestoredClosed PTWStatus = "RESTORED_CLOSED"
	// PTWRejectedCancelled indicates the permit request was rejected or voided.
	PTWRejectedCancelled PTWStatus = "REJECTED_CANCELLED"
)

// ShiftWindow identifies the active dispatch shift window.
type ShiftWindow string

const (
	ShiftMorning   ShiftWindow = "MORNING"
	ShiftAfternoon ShiftWindow = "AFTERNOON"
	ShiftNight     ShiftWindow = "NIGHT"
)

// DispatchTelemetryLog records hourly bulk energy intake and remote substation roll-call logs.
type DispatchTelemetryLog struct {
	ID                          string              `json:"id"`
	ShiftID                     string              `json:"shift_id"`
	ShiftWindow                 ShiftWindow         `json:"shift_window"`
	GenerationMWIn              float64             `json:"generation_mw_in"`
	TCNMWIn                     float64             `json:"tcn_mw_in"`
	TotalBulkMW                 float64             `json:"total_bulk_mw"`
	SubstationStaffRollcallJSON string              `json:"substation_staff_rollcall_json"`
	HandoverNotes               string              `json:"handover_notes"`
	OperatorID                  string              `json:"operator_id"`
	LoggedAt                    time.Time           `json:"logged_at"`
}

// FeederStatus represents real-time telemetry, load, and priority dispatch attributes of a feeder.
type FeederStatus struct {
	ID                       string               `json:"id"`
	FeederAssetID            string               `json:"feeder_asset_id"`
	FeederCode               string               `json:"feeder_code"`
	CurrentLoadMW            float64              `json:"current_load_mw"`
	CurrentAmps              float64              `json:"current_amps"`
	State                    FeederState          `json:"state"`
	OutageReasonCode         *string              `json:"outage_reason_code,omitempty"`
	OutageNarrative          *string              `json:"outage_narrative,omitempty"`
	PriorityRank             int                  `json:"priority_rank"`
	TariffBand               s01domain.TariffBand `json:"tariff_band"`
	ActivePTWID              *string              `json:"active_ptw_id,omitempty"`
	ScheduledRestorationTime *time.Time           `json:"scheduled_restoration_time,omitempty"`
	UpdatedAt                time.Time            `json:"updated_at"`
}

// SafetyPermitPTW encapsulates the 5-step life cycle of high-voltage line isolation and restoration.
type SafetyPermitPTW struct {
	ID                      string     `json:"id"`
	PermitNo                string     `json:"permit_no"`
	WorkOrderID             *string    `json:"work_order_id,omitempty"`
	AssetID                 string     `json:"asset_id"`
	UnitHeadID              string     `json:"unit_head_id"`
	DeptHeadID              *string    `json:"dept_head_id,omitempty"`
	DeptHeadApproved        bool       `json:"dept_head_approved"`
	DeptHeadApprovedAt      *time.Time `json:"dept_head_approved_at,omitempty"`
	SafetyRepID             *string    `json:"safety_rep_id,omitempty"`
	SafetyRepApproved       bool       `json:"safety_rep_approved"`
	SafetyRepApprovedAt     *time.Time `json:"safety_rep_approved_at,omitempty"`
	SubstationOperatorID    *string    `json:"substation_operator_id,omitempty"`
	IsolationGrantedAt      *time.Time `json:"isolation_granted_at,omitempty"`
	IsolationNotes          *string    `json:"isolation_notes,omitempty"`
	SafetyClearanceSignedAt *time.Time `json:"safety_clearance_signed_at,omitempty"`
	SafetyClearanceNotes    *string    `json:"safety_clearance_notes,omitempty"`
	RestoredAt              *time.Time `json:"restored_at,omitempty"`
	Status                  PTWStatus  `json:"status"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// RationingDecision details the allocation outcome determined by the priority rationing algorithm.
type RationingDecision struct {
	FeederCode    string      `json:"feeder_code"`
	AllocatedMW   float64     `json:"allocated_mw"`
	TargetState   FeederState `json:"target_state"`
	PriorityRank  int         `json:"priority_rank"`
	Justification string      `json:"justification"`
}
