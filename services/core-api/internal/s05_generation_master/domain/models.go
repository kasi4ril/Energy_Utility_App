// Package domain defines core entities, enumeration types, and physical invariants
// for the S05 Generation Department & Master Station Hub bounded context.
package domain

import (
	"time"
)

// ShiftWindow identifies the active 8-hour operational window at the Master Station.
type ShiftWindow string

const (
	ShiftMorning   ShiftWindow = "MORNING"   // 08:00 - 16:00
	ShiftAfternoon ShiftWindow = "AFTERNOON" // 16:00 - 00:00
	ShiftNight     ShiftWindow = "NIGHT"     // 00:00 - 08:00
)

// PlantState represents the operational status of a hydroelectric generating station.
type PlantState string

const (
	StateInService    PlantState = "IN_SERVICE"
	StateStandby      PlantState = "STANDBY"
	StateOutOfService PlantState = "OUT_OF_SERVICE"
)

// OutageReason categorizes the root cause when a plant or turbine trip occurs.
type OutageReason string

const (
	OutageHydroLowWater       OutageReason = "HYDRO_LOW_WATER"
	OutageTurbineMechFault    OutageReason = "TURBINE_MECH_FAULT"
	OutageGeneratorElecFault  OutageReason = "GENERATOR_ELEC_FAULT"
	OutageProtectionRelayTrip OutageReason = "PROTECTION_RELAY_TRIP"
	OutagePlannedMaintenance  OutageReason = "PLANNED_MAINTENANCE"
)

// TechnicalWing identifies the engineering discipline within Generation.
type TechnicalWing string

const (
	WingElectrical TechnicalWing = "ELECTRICAL"
	WingMechanical TechnicalWing = "MECHANICAL"
)

// WorkshopStatus tracks fabrication and reconditioning lifecycle in the central workshop.
type WorkshopStatus string

const (
	WorkshopNotApplicable       WorkshopStatus = "NOT_APPLICABLE"
	WorkshopPendingDispatch     WorkshopStatus = "PENDING_DISPATCH"
	WorkshopMachiningInProgress WorkshopStatus = "MACHINING_IN_PROGRESS"
	WorkshopWindingFabrication  WorkshopStatus = "WINDING_FABRICATION"
	WorkshopQualityCheck        WorkshopStatus = "QUALITY_CHECK"
	WorkshopCompletedDelivered  WorkshopStatus = "COMPLETED_DELIVERED"
	WorkshopRejected            WorkshopStatus = "REJECTED"
)

// PlantJobStatus tracks the workflow status of plant faults and repair authorizations.
type PlantJobStatus string

const (
	JobStatusReported             PlantJobStatus = "REPORTED"
	JobStatusUnitHeadReviewed     PlantJobStatus = "UNIT_HEAD_REVIEWED"
	JobStatusGMApproved           PlantJobStatus = "GM_APPROVED"
	JobStatusWorkshopInProgress   PlantJobStatus = "WORKSHOP_IN_PROGRESS"
	JobStatusResolvedReinstalled  PlantJobStatus = "RESOLVED_REINSTALLED"
	JobStatusCancelled            PlantJobStatus = "CANCELLED"
)

// MasterStationTelemetryLog records hourly hydrological and generator telemetry from Plants 1–6.
type MasterStationTelemetryLog struct {
	ID                   string        `json:"id"`
	ShiftID              string        `json:"shift_id"`
	ShiftWindow          ShiftWindow   `json:"shift_window"`
	PlantID              string        `json:"plant_id"`
	PlantState           PlantState    `json:"plant_state"`
	OutageReason         *OutageReason `json:"outage_reason,omitempty"`
	OutageDetails        *string       `json:"outage_details,omitempty"`
	WaterLevelM          float64       `json:"water_level_m"`
	HeadwaterElevationM  float64       `json:"headwater_elevation_m"`
	TailwaterElevationM  float64       `json:"tailwater_elevation_m"`
	InflowRateM3S        float64       `json:"inflow_rate_m3_s"`
	MWGenerated          float64       `json:"mw_generated"`
	MVArReactive         float64       `json:"mvar_reactive"`
	StatorWindingTempC   float64       `json:"stator_winding_temp_c"`
	BearingVibrationMMS  float64       `json:"bearing_vibration_mm_s"`
	LubeOilPressureBar   float64       `json:"lube_oil_pressure_bar"`
	LoggedByOperatorID   string        `json:"logged_by_operator_id"`
	RecordedAt           time.Time     `json:"recorded_at"`
}

// RemotePlantShiftRollcall captures on-duty personnel physically stationed at remote plants.
type RemotePlantShiftRollcall struct {
	ID                  string      `json:"id"`
	ShiftID             string      `json:"shift_id"`
	ShiftWindow         ShiftWindow `json:"shift_window"`
	PlantID             string      `json:"plant_id"`
	OperatorNames       []string    `json:"operator_names"`
	ProtectionEngineers []string    `json:"protection_engineers"`
	Technicians         []string    `json:"technicians"`
	HandoverNotes       string      `json:"handover_notes"`
	SyncedToHQHR        bool        `json:"synced_to_hq_hr"`
	MasterOperatorID    string      `json:"master_operator_id"`
	CreatedAt           time.Time   `json:"created_at"`
}

// PlantFaultAndWorkshopJob coordinates plant repair, GM material approval, and workshop fabrication.
type PlantFaultAndWorkshopJob struct {
	ID                   string         `json:"id"`
	JobNo                string         `json:"job_no"`
	PlantID              string         `json:"plant_id"`
	ReportedByOperatorID string         `json:"reported_by_operator_id"`
	Wing                 TechnicalWing  `json:"wing"`
	FaultDescription     string         `json:"fault_description"`
	Severity             string         `json:"severity"`
	UnitHeadID           *string        `json:"unit_head_id,omitempty"`
	GMApproved           bool           `json:"gm_approved"`
	GMApprovalNotes      *string        `json:"gm_approval_notes,omitempty"`
	RequiresWorkshop     bool           `json:"requires_workshop"`
	WorkshopStatus       WorkshopStatus `json:"workshop_status"`
	RawMaterialsSummary  *string        `json:"raw_materials_summary,omitempty"`
	StoreRequisitionRef  *string        `json:"store_requisition_ref,omitempty"`
	Status               PlantJobStatus `json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

