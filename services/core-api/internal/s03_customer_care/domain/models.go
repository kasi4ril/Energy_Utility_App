// Package domain defines core entities, value objects, and routing rules for
// the S03 Customer Care bounded context of the NESCO Enterprise Utility Platform.
package domain

import (
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

// TicketChannel identifies the communication channel through which a customer lodged a ticket.
type TicketChannel string

const (
	ChannelWalkIn       TicketChannel = "WALK_IN"
	ChannelOnlinePortal TicketChannel = "ONLINE_PORTAL"
	ChannelCallCenter   TicketChannel = "CALL_CENTER"
)

// TicketCategory defines standardized utility fault and service classifications.
type TicketCategory string

const (
	CategoryFallenPole         TicketCategory = "FALLEN_POLE"
	CategoryLineSnap           TicketCategory = "LINE_SNAP"
	CategoryTransformerFault   TicketCategory = "TRANSFORMER_FAULT"
	CategoryFuseBlown          TicketCategory = "FUSE_BLOWN"
	CategoryOilLeak            TicketCategory = "OIL_LEAK"
	CategoryFaultyMeter        TicketCategory = "FAULTY_METER"
	CategoryNewMeterSchedule   TicketCategory = "NEW_METER_SCHEDULE"
	CategoryWrongBilling       TicketCategory = "WRONG_BILLING"
	CategoryUncreditedPayment  TicketCategory = "UNCREDITED_PAYMENT"
	CategoryGeneralInquiry     TicketCategory = "GENERAL_INQUIRY"
)

// TicketStatus defines the strict lifecycle states of a customer service ticket.
type TicketStatus string

const (
	StatusLogged            TicketStatus = "LOGGED"
	StatusRouted            TicketStatus = "ROUTED"
	StatusAttendedTo        TicketStatus = "ATTENDED_TO"
	StatusClosed            TicketStatus = "CLOSED"
	StatusReopenedEscalated TicketStatus = "REOPENED_ESCALATED"
)

// ScheduleType defines the technical activity for metering field teams.
type ScheduleType string

const (
	ScheduleNewInstallation ScheduleType = "NEW_INSTALLATION"
	ScheduleRepair          ScheduleType = "REPAIR"
	ScheduleCalibration     ScheduleType = "CALIBRATION"
	ScheduleRecertification ScheduleType = "RECERTIFICATION"
)

// ScheduleStatus defines the progress of a meter installation schedule.
type ScheduleStatus string

const (
	ScheduleStatusScheduled   ScheduleStatus = "SCHEDULED"
	ScheduleStatusDispatched  ScheduleStatus = "DISPATCHED"
	ScheduleStatusCompleted   ScheduleStatus = "COMPLETED"
	ScheduleStatusRescheduled ScheduleStatus = "RESCHEDULED"
	ScheduleStatusCancelled   ScheduleStatus = "CANCELLED"
)

// KYCStatus defines verification and debtor standing of customer accounts.
type KYCStatus string

const (
	KYCPending        KYCStatus = "PENDING"
	KYCVerified       KYCStatus = "VERIFIED"
	KYCSuspended      KYCStatus = "SUSPENDED"
	KYCDebtorFlagged  KYCStatus = "DEBTOR_FLAGGED"
)

// BuildingClass defines premises categorization for billing tariff bands.
type BuildingClass string

const (
	BuildingResidential BuildingClass = "RESIDENTIAL"
	BuildingCommercial  BuildingClass = "COMMERCIAL"
	BuildingIndustrial  BuildingClass = "INDUSTRIAL"
	BuildingSpecial     BuildingClass = "SPECIAL"
)

// Customer encapsulates an onboarded electricity consumer.
type Customer struct {
	ID               string             `json:"id"`
	AccountNo        string             `json:"account_no"`
	CIN              string             `json:"cin"`
	FullName         string             `json:"full_name"`
	Phone            string             `json:"phone"`
	NIN              string             `json:"nin,omitempty"`
	BuildingClass    BuildingClass      `json:"building_class"`
	PremisesGPS      s01domain.GPSPoint `json:"premises_gps"`
	TransformerID    string             `json:"transformer_id"`
	MeterSerial      string             `json:"meter_serial,omitempty"`
	TariffBand       s01domain.TariffBand `json:"tariff_band"`
	LedgerBalance    float64            `json:"ledger_balance"`
	KYCStatus        KYCStatus          `json:"kyc_status"`
	RowHMACSignature string             `json:"row_hmac_signature"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// ServiceTicket encapsulates an omnichannel complaint or field service request.
type ServiceTicket struct {
	ID                     string         `json:"id"`
	TicketNo               string         `json:"ticket_no"`
	CustomerID             *string        `json:"customer_id,omitempty"`
	Channel                TicketChannel  `json:"channel"`
	Category               TicketCategory `json:"category"`
	Description            string         `json:"description"`
	AssignedUnitID         string         `json:"assigned_unit_id"`
	UnitHeadID             *string        `json:"unit_head_id,omitempty"`
	DeptHeadSupervisorID   *string        `json:"dept_head_supervisor_id,omitempty"`
	Status                 TicketStatus   `json:"status"`
	AttendedProofNote      string         `json:"attended_proof_note,omitempty"`
	AttendedByStaffID      *string        `json:"attended_by_staff_id,omitempty"`
	CustomerCallVerified   bool           `json:"customer_call_verified"`
	VerifiedByCCOfficerID  *string        `json:"verified_by_cc_officer_id,omitempty"`
	CustomerFeedbackNotes  string         `json:"customer_feedback_notes,omitempty"`
	SLADueAt               time.Time      `json:"sla_due_at"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

// MeterSchedule encapsulates an appointment for metering installations or recertifications.
type MeterSchedule struct {
	ID                     string         `json:"id"`
	CustomerID             string         `json:"customer_id"`
	TicketID               *string        `json:"ticket_id,omitempty"`
	ScheduleType           ScheduleType   `json:"schedule_type"`
	ScheduledDate          time.Time      `json:"scheduled_date"`
	AssignedMeteringTeamID *string        `json:"assigned_metering_team_id,omitempty"`
	Status                 ScheduleStatus `json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

// RoutingTarget defines the target unit head and supervising departmental head for ticket routing.
type RoutingTarget struct {
	TargetUnitHeadRole     string
	SupervisingDeptRole    string
	SLADurationHours       int
}

