// Package domain encapsulates core entities, value objects, and business rules
// for the Foundation bounded context (S01) of the NESCO Enterprise Utility Platform.
// It enforces organizational hierarchy, electrical grid asset topologies, staff profiles,
// and append-only cryptographic WorkNotes.
package domain

import (
	"time"
)

// OrgUnitType identifies administrative classifications across the enterprise.
type OrgUnitType string

const (
	OrgUnitHQ              OrgUnitType = "HQ"
	OrgUnitDistrict        OrgUnitType = "DISTRICT"
	OrgUnitServiceUnit     OrgUnitType = "SERVICE_UNIT"
	OrgUnitMasterStation   OrgUnitType = "MASTER_STATION"
	OrgUnitGeneratingPlant OrgUnitType = "GENERATING_PLANT"
)

// UserRole identifies the role assigned to a staff profile across the 3-Tier RBAC matrix.
type UserRole string

const (
	// Executive Command Tier (Shared Dashboard)
	RoleMD              UserRole = "MD"
	RoleHeadServices    UserRole = "HEAD_SERVICES"
	RoleHeadAudit       UserRole = "HEAD_INTERNAL_AUDIT"
	RoleHeadCommercial  UserRole = "HEAD_COMMERCIAL_OPS"
	RoleGMGeneration    UserRole = "GM_GENERATION"
	RoleHeadEnergyMgmt  UserRole = "HEAD_ENERGY_MGMT"

	// Departmental & Unit Heads
	RoleCustomerCareHead  UserRole = "CUSTOMER_CARE_HEAD"
	RoleUnitHeadFeeder    UserRole = "UNIT_HEAD_FEEDER_MAINT"
	RoleUnitHeadTx        UserRole = "UNIT_HEAD_TX_MAINT"
	RoleUnitHeadMetering  UserRole = "UNIT_HEAD_METERING"
	RoleUnitHeadGenElec   UserRole = "UNIT_HEAD_GEN_ELECTRICAL"
	RoleUnitHeadGenMech   UserRole = "UNIT_HEAD_GEN_MECHANICAL"
	RoleWorkshopLead      UserRole = "WORKSHOP_LEAD"
	RoleUnitHeadBilling   UserRole = "UNIT_HEAD_BILLING"
	RoleStoreOfficer      UserRole = "STORE_OFFICER"
	RoleSafetyRepHSE      UserRole = "SAFETY_REP_HSE"
	RolePlantHROfficer    UserRole = "PLANT_HR_OFFICER"

	// Field & Operational Tier
	RoleMasterStationOp   UserRole = "MASTER_STATION_OPERATOR"
	RoleCentralDispatchOp UserRole = "CENTRAL_DISPATCH_OPERATOR"
	RoleSubstationOp      UserRole = "SUBSTATION_OPERATOR"
	RoleTeamLeadFeeder    UserRole = "TEAM_LEAD_FEEDER"
	RoleTeamLeadTx        UserRole = "TEAM_LEAD_TX"
	RoleTeamLeadMetering  UserRole = "TEAM_LEAD_METERING"
	RoleCustomerCareOp    UserRole = "CUSTOMER_CARE_OFFICER"
	RoleInternalAuditor   UserRole = "INTERNAL_AUDITOR"
)

// AssetType identifies physical assets in the transmission and distribution network.
type AssetType string

const (
	AssetPlant          AssetType = "PLANT"
	AssetSubstation33KV AssetType = "SUBSTATION_33KV"
	AssetFeeder11KV     AssetType = "FEEDER_11KV"
	AssetTransformerDT  AssetType = "TRANSFORMER_DT"
	AssetServicePole    AssetType = "SERVICE_POLE"
	AssetCustomerMeter  AssetType = "CUSTOMER_METER"
)

// TariffBand defines the NERC regulatory supply commitment bands (A to E).
type TariffBand string

const (
	TariffBandA TariffBand = "A" // >= 20 hrs/day
	TariffBandB TariffBand = "B" // >= 16 hrs/day
	TariffBandC TariffBand = "C" // >= 12 hrs/day
	TariffBandD TariffBand = "D" // >= 8 hrs/day
	TariffBandE TariffBand = "E" // >= 4 hrs/day
	TariffBandNone TariffBand = ""
)

// OperationalState defines the current service status of physical grid assets.
type OperationalState string

const (
	StateInService      OperationalState = "IN_SERVICE"
	StateStandby        OperationalState = "STANDBY"
	StateOutOfService   OperationalState = "OUT_OF_SERVICE"
	StateMaintenance    OperationalState = "MAINTENANCE"
	StateDecommissioned OperationalState = "DECOMMISSIONED"
)

// CertificationStatus defines technical/safety certifications of staff.
type CertificationStatus string

const (
	CertActive    CertificationStatus = "ACTIVE"
	CertSuspended CertificationStatus = "SUSPENDED"
	CertExpired   CertificationStatus = "EXPIRED"
	CertProbation CertificationStatus = "PROBATION"
)

// ShiftState defines current shift engagement of field/plant operators.
type ShiftState string

const (
	ShiftOnDuty  ShiftState = "ON_DUTY"
	ShiftOffDuty ShiftState = "OFF_DUTY"
	ShiftStandby ShiftState = "STANDBY"
	ShiftLeave   ShiftState = "LEAVE"
)

// GPSPoint represents geographic coordinates in WGS 84 (EPSG:4326).
type GPSPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// OrgUnit represents an administrative unit in the company directory.
type OrgUnit struct {
	ID          string      `json:"id"`
	ParentID    *string     `json:"parent_id,omitempty"`
	UnitCode    string      `json:"unit_code"`
	Name        string      `json:"name"`
	UnitType    OrgUnitType `json:"unit_type"`
	LocationGPS *GPSPoint   `json:"location_gps,omitempty"`
	IsActive    bool        `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// StaffProfile represents a federated employee record.
type StaffProfile struct {
	ID                  string              `json:"id"`
	StaffNumber         string              `json:"staff_number"`
	FullName            string              `json:"full_name"`
	Email               string              `json:"email"`
	Role                UserRole            `json:"role"`
	OrgUnitID           string              `json:"org_unit_id"`
	CertificationStatus CertificationStatus `json:"certification_status"`
	ActiveShiftState    ShiftState          `json:"active_shift_state"`
	PasswordHash        string              `json:"-"`
	RowHMACSignature    string              `json:"row_hmac_signature"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

// GridAsset represents a physical electrical asset in the grid hierarchy.
type GridAsset struct {
	ID               string           `json:"id"`
	ParentAssetID    *string          `json:"parent_asset_id,omitempty"`
	OrgUnitID        string           `json:"org_unit_id"`
	AssetCode        string           `json:"asset_code"`
	AssetType        AssetType        `json:"asset_type"`
	RatedCapacityKVA float64          `json:"rated_capacity_kva"`
	TariffBand       TariffBand       `json:"tariff_band,omitempty"`
	GPSPoint         GPSPoint         `json:"gps_point"`
	OperationalState OperationalState `json:"operational_state"`
	RowHMACSignature string           `json:"row_hmac_signature"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

// WorkNote represents an append-only inter-departmental note with cryptographic chaining.
type WorkNote struct {
	ID            string    `json:"id"`
	EntityType    string    `json:"entity_type"`
	EntityID      string    `json:"entity_id"`
	AuthorStaffID string    `json:"author_staff_id"`
	TaggedUnitIDs []string  `json:"tagged_unit_ids"`
	NoteText      string    `json:"note_text"`
	OllamaSummary string    `json:"ollama_summary,omitempty"`
	PrevHash      string    `json:"prev_hash"`
	CurrentHash   string    `json:"current_hash"`
	CreatedAt     time.Time `json:"created_at"`
}

