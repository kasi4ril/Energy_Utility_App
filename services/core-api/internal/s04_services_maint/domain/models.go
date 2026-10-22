// Package domain provides core entity representations, enumeration types,
// and pure business logic invariants for the S04 Services & Maintenance bounded context.
package domain

import (
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

// UnitType identifies the specific engineering branch under the Head of Services.
type UnitType string

const (
	// UnitFeederMaint handles 33kV and 11kV overhead lines, crossarms, insulators, and poles.
	UnitFeederMaint UnitType = "FEEDER_MAINT"
	// UnitTransformerMaint handles 33/11kV substations, distribution transformers, fuses, and oil filtering.
	UnitTransformerMaint UnitType = "TRANSFORMER_MAINT"
	// UnitMetering handles maximum demand (MD) and non-MD metering, seal enforcement, and calibration.
	UnitMetering UnitType = "METERING"
)

// WorkOrderType specifies the nature of the field maintenance engagement.
type WorkOrderType string

const (
	// WOCorrectiveRepair represents fault clearance arising from customer complaints or grid trips.
	WOCorrectiveRepair WorkOrderType = "CORRECTIVE_REPAIR"
	// WOPreventiveMaintenance represents routine scheduled asset inspection and servicing.
	WOPreventiveMaintenance WorkOrderType = "PREVENTIVE_MAINTENANCE"
	// WOEmergencyOutage represents immediate safety hazard isolation or catastrophic failure response.
	WOEmergencyOutage WorkOrderType = "EMERGENCY_OUTAGE"
	// WOMeterInstallation represents customer meter replacement, recertification, or new deployment.
	WOMeterInstallation WorkOrderType = "METER_INSTALLATION"
	// WODTInspection represents load balancing, oil level checks, and thermographic scans on transformers.
	WODTInspection WorkOrderType = "DT_INSPECTION"
)

// WorkOrderPriority dictates the dispatch urgency and response SLA.
type WorkOrderPriority string

const (
	PriorityLow      WorkOrderPriority = "LOW"
	PriorityMedium   WorkOrderPriority = "MEDIUM"
	PriorityHigh     WorkOrderPriority = "HIGH"
	PriorityCritical WorkOrderPriority = "CRITICAL"
)

// WorkOrderStatus tracks the formal state transitions of a maintenance work order.
type WorkOrderStatus string

const (
	// StatusDraft indicates the work order is drafted by the Unit Head.
	StatusDraft WorkOrderStatus = "DRAFT"
	// StatusDispatched indicates work is officially issued to the Field Team Lead.
	StatusDispatched WorkOrderStatus = "DISPATCHED_TO_FIELD"
	// StatusInProgress indicates the field crew is on-site and actively carrying out work.
	StatusInProgress WorkOrderStatus = "WORK_IN_PROGRESS"
	// StatusProofSubmitted indicates the Team Lead uploaded photos, GPS, and technical proof.
	StatusProofSubmitted WorkOrderStatus = "PROOF_SUBMITTED"
	// StatusAttendedVerified indicates the Unit Head verified the proof; ticket reverts to Customer Care.
	StatusAttendedVerified WorkOrderStatus = "ATTENDED_VERIFIED"
	// StatusRejectedRework indicates proof was rejected by the Unit Head; crew must remediate.
	StatusRejectedRework WorkOrderStatus = "REJECTED_REWORK"
	// StatusCancelled indicates the work order was voided.
	StatusCancelled WorkOrderStatus = "CANCELLED"
)

// MaintenanceWorkOrder encapsulates an operational directive issued to a field team.
type MaintenanceWorkOrder struct {
	ID                       string            `json:"id"`
	OrderNo                  string            `json:"order_no"`
	TicketID                 *string           `json:"ticket_id,omitempty"`
	UnitType                 UnitType          `json:"unit_type"`
	WorkOrderType            WorkOrderType     `json:"work_order_type"`
	Priority                 WorkOrderPriority `json:"priority"`
	UnitHeadID               string            `json:"unit_head_id"`
	TeamLeadID               string            `json:"team_lead_id"`
	AssetID                  string            `json:"asset_id"`
	PTWRequired              bool              `json:"ptw_required"`
	PTWPermitID              *string           `json:"ptw_permit_id,omitempty"`
	Status                   WorkOrderStatus   `json:"status"`
	MaterialsRequisitionRef  *string           `json:"materials_requisition_ref,omitempty"`
	Instructions             string            `json:"instructions"`
	ReworkNotes              *string           `json:"rework_notes,omitempty"`
	CreatedAt                time.Time         `json:"created_at"`
	UpdatedAt                time.Time         `json:"updated_at"`
}

// ProofOfWorkLog represents the tamper-evident evidence submitted from the mobile PWA.
type ProofOfWorkLog struct {
	ID                      string            `json:"id"`
	WorkOrderID             string            `json:"work_order_id"`
	TeamLeadID              string            `json:"team_lead_id"`
	CapturedGPS             s01domain.GPSPoint `json:"captured_gps"`
	DistanceFromAssetMeters float64           `json:"distance_from_asset_meters"`
	BeforePhotoURL          string            `json:"before_photo_url"`
	AfterPhotoURL           string            `json:"after_photo_url"`
	PhotoHashBefore         string            `json:"photo_hash_before"`
	PhotoHashAfter          string            `json:"photo_hash_after"`
	MeterSealNumber         *string           `json:"meter_seal_number,omitempty"`
	WorkDetails             string            `json:"work_details"`
	OllamaGPSVerified       bool              `json:"ollama_gps_verified"`
	AIRiskScore             float64           `json:"ai_risk_score"`
	AIReasoning             string            `json:"ai_reasoning"`
	UnitHeadVerified        bool              `json:"unit_head_verified"`
	UnitHeadNotes           *string           `json:"unit_head_notes,omitempty"`
	SubmittedAt             time.Time         `json:"submitted_at"`
	VerifiedAt              *time.Time        `json:"verified_at,omitempty"`
}

