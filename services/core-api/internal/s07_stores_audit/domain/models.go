// Package domain encapsulates core entities, state machine enumerations,
// and pure business logic invariants for the S07 Stores, Inventory & Internal Audit bounded context.
package domain

import (
	"time"
)

// ItemCategory classifies materials stocked in the central warehouse and district stores.
type ItemCategory string

const (
	CategoryCablesConductors    ItemCategory = "CABLES_CONDUCTORS"
	CategoryTransformerParts    ItemCategory = "TRANSFORMER_PARTS"
	CategoryMetersAccessories   ItemCategory = "METERS_ACCESSORIES"
	CategorySafetyGear          ItemCategory = "SAFETY_GEAR"
	CategorySubstationHardware  ItemCategory = "SUBSTATION_HARDWARE"
	CategoryConsumablesOils     ItemCategory = "CONSUMABLES_OILS"
)

// RequisitionStatus defines the progression of material requisition and procurement requests.
type RequisitionStatus string

const (
	StatusRequisitionSubmitted          RequisitionStatus = "SUBMITTED"
	StatusRequisitionDeptApproved       RequisitionStatus = "DEPT_APPROVED"
	StatusRequisitionAuditConsented     RequisitionStatus = "AUDIT_CONSENTED"
	StatusRequisitionMDProcurementPending RequisitionStatus = "MD_PROCUREMENT_PENDING"
	StatusRequisitionMDApproved         RequisitionStatus = "MD_APPROVED"
	StatusRequisitionIssued             RequisitionStatus = "ISSUED"
	StatusRequisitionRejected           RequisitionStatus = "REJECTED_QUERIED"
	StatusRequisitionCancelled          RequisitionStatus = "CANCELLED"
)

// AuditConsentStatus indicates compliance review outcomes by Internal Audit.
type AuditConsentStatus string

const (
	AuditConsentPending         AuditConsentStatus = "PENDING"
	AuditConsentConsented       AuditConsentStatus = "CONSENTED"
	AuditConsentQueriedRejected AuditConsentStatus = "QUERIED_REJECTED"
)

// ClaimType specifies staff financial reimbursement allowances.
type ClaimType string

const (
	ClaimOvertime            ClaimType = "OVERTIME"
	ClaimTravelOutOfStation  ClaimType = "TRAVEL_OUT_OF_STATION"
	ClaimHazard              ClaimType = "HAZARD"
	ClaimMealSubsistence     ClaimType = "MEAL_SUBSISTENCE"
)

// ClaimStatus tracks the 3-way audit verification and payment status of staff allowances.
type ClaimStatus string

const (
	StatusClaimSubmitted     ClaimStatus = "SUBMITTED"
	StatusClaimAIVerified    ClaimStatus = "AI_VERIFIED"
	StatusClaimAIFlagged     ClaimStatus = "AI_FLAGGED"
	StatusClaimAuditApproved ClaimStatus = "AUDIT_APPROVED"
	StatusClaimAuditRejected ClaimStatus = "AUDIT_REJECTED"
	StatusClaimPaid          ClaimStatus = "PAID"
)

// AIAuditFlag marks whether a claim passed or failed the 3-way automated cross-check.
type AIAuditFlag string

const (
	AIFlagVerified        AIAuditFlag = "VERIFIED"
	AIFlagFlaggedMismatch AIAuditFlag = "FLAGGED_MISMATCH"
	AIFlagFlaggedFraud    AIAuditFlag = "FLAGGED_FRAUD"
)

// StoreItem represents a managed warehouse or district store inventory SKU.
type StoreItem struct {
	ID                string       `json:"id"`
	ItemCode          string       `json:"item_code"`
	ItemName          string       `json:"item_name"`
	Category          ItemCategory `json:"category"`
	QuantityAvailable float64      `json:"quantity_available"`
	ReorderThreshold  float64      `json:"reorder_threshold"`
	UnitOfMeasure     string       `json:"unit_of_measure"`
	UnitCost          float64      `json:"unit_cost"`
	BinLocation       string       `json:"bin_location"`
	LastUpdatedBy     string       `json:"last_updated_by"`
	RowHMACSignature  string       `json:"row_hmac_signature"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// RequisitionItem details a specific SKU and quantity inside a material request.
type RequisitionItem struct {
	ItemID            string  `json:"item_id"`
	ItemCode          string  `json:"item_code"`
	ItemName          string  `json:"item_name"`
	QuantityRequested float64 `json:"quantity_requested"`
	UnitCost          float64 `json:"unit_cost"`
}

// MaterialRequisition models dual-path store issuing and MD procurement approvals.
type MaterialRequisition struct {
	ID                    string             `json:"id"`
	ReqNo                 string             `json:"req_no"`
	OriginDeptID          string             `json:"origin_dept_id"`
	LinkedTicketOrJobID   *string            `json:"linked_ticket_or_job_id,omitempty"`
	Items                 []RequisitionItem  `json:"items"`
	TotalCost             float64            `json:"total_cost"`
	InStockFlag           bool               `json:"in_stock_flag"`
	UnitHeadID            string             `json:"unit_head_id"`
	DeptHeadID            *string            `json:"dept_head_id,omitempty"`
	DeptOrGMApproved      bool               `json:"dept_or_gm_approved"`
	DeptApprovedAt        *time.Time         `json:"dept_approved_at,omitempty"`
	AuditOfficerID        *string            `json:"audit_officer_id,omitempty"`
	AuditConsentStatus    AuditConsentStatus `json:"audit_consent_status"`
	AuditConsentAt        *time.Time         `json:"audit_consent_at,omitempty"`
	AuditNotes            *string            `json:"audit_notes,omitempty"`
	MDProcurementApproved bool               `json:"md_procurement_approved"`
	MDApprovedAt          *time.Time         `json:"md_approved_at,omitempty"`
	MDNotes               *string            `json:"md_notes,omitempty"`
	StoreOfficerID        *string            `json:"store_officer_id,omitempty"`
	IssuedAt              *time.Time         `json:"issued_at,omitempty"`
	WaybillNo             *string            `json:"waybill_no,omitempty"`
	ScrapReturnedFlag     bool               `json:"scrap_returned_flag"`
	ScrapReturnID         *string            `json:"scrap_return_id,omitempty"`
	Status                RequisitionStatus  `json:"status"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
}

// ScrapReturn tracks damaged or replaced assets returned to the Store Scrap Yard.
type ScrapReturn struct {
	ID                  string    `json:"id"`
	ReturnNo            string    `json:"return_no"`
	RequisitionID       string    `json:"requisition_id"`
	TeamLeadID          string    `json:"team_lead_id"`
	StoreOfficerID      string    `json:"store_officer_id"`
	ItemsReturnedJSON   string    `json:"items_returned_json"`
	ScrapConditionNotes string    `json:"scrap_condition_notes"`
	ReturnedAt          time.Time `json:"returned_at"`
	CreatedAt           time.Time `json:"created_at"`
}

// StaffAllowanceClaim represents an employee overtime, travel, or hazard financial claim.
type StaffAllowanceClaim struct {
	ID                   string      `json:"id"`
	ClaimNo              string      `json:"claim_no"`
	StaffID              string      `json:"staff_id"`
	ClaimType            ClaimType   `json:"claim_type"`
	Amount               float64     `json:"amount"`
	ClaimDate            string      `json:"claim_date"`
	ShiftWindow          *string     `json:"shift_window,omitempty"`
	LinkedPTWOrTicketID  *string     `json:"linked_ptw_or_ticket_id,omitempty"`
	RosterMatchVerified  bool        `json:"roster_match_verified"`
	PTWCrewVerified      bool        `json:"ptw_crew_verified"`
	GPSMatchVerified     bool        `json:"gps_match_verified"`
	OllamaAuditFlag      AIAuditFlag `json:"ollama_audit_flag"`
	OllamaFraudScore     float64     `json:"ollama_fraud_score"`
	AuditorID            *string     `json:"auditor_id,omitempty"`
	AuditorNotes         *string     `json:"auditor_notes,omitempty"`
	AuditDecisionAt      *time.Time  `json:"audit_decision_at,omitempty"`
	Status               ClaimStatus `json:"status"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

