// Package domain encapsulates core entities, enumerations, and mathematical invariants
// for the S09 Executive Command, Corporate Secretarial & Predictive Maintenance bounded context.
package domain

import (
	"time"
)

// ExecutiveRole defines executive C-Suite roles authorized to access the command center.
type ExecutiveRole string

const (
	RoleMD                ExecutiveRole = "MD"
	RoleHeadServices      ExecutiveRole = "HEAD_SERVICES"
	RoleHeadInternalAudit ExecutiveRole = "HEAD_INTERNAL_AUDIT"
	RoleHeadCommercialOps ExecutiveRole = "HEAD_COMMERCIAL_OPS"
	RoleGMGeneration      ExecutiveRole = "GM_GENERATION"
	RoleHeadEnergyMgmt    ExecutiveRole = "HEAD_ENERGY_MGMT"
)

// AssetCategory categorizes grid and generation assets monitored by predictive models.
type AssetCategory string

const (
	AssetGenerator      AssetCategory = "GENERATOR"
	AssetTransformer    AssetCategory = "TRANSFORMER"
	AssetFeeder         AssetCategory = "FEEDER"
	AssetHydroReservoir AssetCategory = "HYDRO_RESERVOIR"
)

// PredictiveModelName specifies the machine learning or physics-based prognostic pipeline.
type PredictiveModelName string

const (
	ModelHydroWaterAvailability PredictiveModelName = "HYDRO_WATER_AVAILABILITY"
	ModelGeneratorPreventive    PredictiveModelName = "GENERATOR_PREVENTIVE"
	ModelTransformerCBM         PredictiveModelName = "TRANSFORMER_CBM"
	ModelFeederPatrol           PredictiveModelName = "FEEDER_PATROL"
)

// RiskLevel categorizes the severity of equipment failure probability or water depletion.
type RiskLevel string

const (
	RiskLow      RiskLevel = "LOW"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskHigh     RiskLevel = "HIGH"
	RiskCritical RiskLevel = "CRITICAL"
)

// ActionItemType categorizes pending C-Suite operational approval requests.
type ActionItemType string

const (
	ActionMDProcurementApproval     ActionItemType = "MD_PROCUREMENT_APPROVAL"
	ActionHeadCommBillingAdjustment ActionItemType = "HEAD_COMM_BILLING_ADJUSTMENT"
	ActionHeadAuditClaimsConsent    ActionItemType = "HEAD_AUDIT_CLAIMS_CONSENT"
	ActionHeadServicesPTW           ActionItemType = "HEAD_SERVICES_PTW"
	ActionGMGenMaterialReq          ActionItemType = "GM_GEN_MATERIAL_REQ"
)

// ActionItemStatus tracks executive queue progression.
type ActionItemStatus string

const (
	ActionStatusPending  ActionItemStatus = "PENDING"
	ActionStatusApproved ActionItemStatus = "APPROVED"
	ActionStatusRejected ActionItemStatus = "REJECTED"
	ActionStatusExpired  ActionItemStatus = "EXPIRED"
)

// ExecutiveSnapshot models periodic aggregated system telemetry and business metrics.
type ExecutiveSnapshot struct {
	ID                            string    `json:"id"`
	SnapshotDate                  string    `json:"snapshot_date"`
	ShiftWindow                   string    `json:"shift_window"`
	TotalGenerationMW             float64   `json:"total_generation_mw"`
	TCNGridIntakeMW               float64   `json:"tcn_grid_intake_mw"`
	TotalSystemLoadMW             float64   `json:"total_system_load_mw"`
	FeedersInServiceCount         int       `json:"feeders_in_service_count"`
	FeedersOutFaultCount          int       `json:"feeders_out_fault_count"`
	FeedersOutPTWCount            int       `json:"feeders_out_ptw_count"`
	FeedersOutRationingCount      int       `json:"feeders_out_rationing_count"`
	TotalRevenueCollected         float64   `json:"total_revenue_collected"`
	PrepaidRevenueSharePct        float64   `json:"prepaid_revenue_share_pct"`
	UnresolvedTheftAlertsCount    int       `json:"unresolved_theft_alerts_count"`
	UnresolvedChurnAlertsCount    int       `json:"unresolved_churn_alerts_count"`
	AuditPendingProcurementsCount int       `json:"audit_pending_procurements_count"`
	UnreturnedScrapCount          int       `json:"unreturned_scrap_count"`
	Notes                         string    `json:"notes,omitempty"`
	CapturedByStaffID             string    `json:"captured_by_staff_id"`
	CreatedAt                     time.Time `json:"created_at"`
}

// PredictiveMaintenanceTrigger models automated prognostic failure alerts and work order recommendations.
type PredictiveMaintenanceTrigger struct {
	ID                    string              `json:"id"`
	TriggerNo             string              `json:"trigger_no"`
	AssetID               string              `json:"asset_id"`
	AssetType             AssetCategory       `json:"asset_type"`
	ModelName             PredictiveModelName `json:"model_name"`
	RiskScore             float64             `json:"risk_score"`
	RiskLevel             RiskLevel           `json:"risk_level"`
	TelemetrySummaryJSON  string              `json:"telemetry_summary_json,omitempty"`
	RecommendedAction     string              `json:"recommended_action"`
	AutoWorkOrderCreated  bool                `json:"auto_work_order_created"`
	LinkedWorkOrderID     *string             `json:"linked_work_order_id,omitempty"`
	AcknowledgedByStaffID *string             `json:"acknowledged_by_staff_id,omitempty"`
	AcknowledgedAt        *time.Time          `json:"acknowledged_at,omitempty"`
	Status                string              `json:"status"`
	CreatedAt             time.Time           `json:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at"`
}

// ExecutiveActionItem models role-specific pending approvals across departments for the C-Suite console.
type ExecutiveActionItem struct {
	ID               string           `json:"id"`
	ActionType       ActionItemType   `json:"action_type"`
	TargetRole       ExecutiveRole    `json:"target_role"`
	SourceEntityID   string           `json:"source_entity_id"`
	SourceEntityRef  string           `json:"source_entity_ref"`
	SourceDepartment string           `json:"source_department"`
	Summary          string           `json:"summary"`
	Amount           float64          `json:"amount"`
	Priority         string           `json:"priority"`
	Status           ActionItemStatus `json:"status"`
	DecidedByStaffID *string          `json:"decided_by_staff_id,omitempty"`
	DecidedAt        *time.Time       `json:"decided_at,omitempty"`
	DecisionNotes    *string          `json:"decision_notes,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

// CSuiteDashboardSummary encapsulates the aggregated live command telemetry for frontend views.
type CSuiteDashboardSummary struct {
	Timestamp                   time.Time                      `json:"timestamp"`
	TotalGenerationMW           float64                        `json:"total_generation_mw"`
	TCNGridIntakeMW             float64                        `json:"tcn_grid_intake_mw"`
	TotalAvailablePowerMW       float64                        `json:"total_available_power_mw"`
	FeedersSummary              FeedersMetricSummary           `json:"feeders_summary"`
	RevenueSummary              RevenuePulseSummary            `json:"revenue_summary"`
	IntegritySummary            IntegrityRadarSummary          `json:"integrity_summary"`
	ActivePredictiveTriggers    []PredictiveMaintenanceTrigger `json:"active_predictive_triggers"`
	PendingExecutiveActionItems []ExecutiveActionItem          `json:"pending_executive_action_items"`
}

// FeedersMetricSummary tracks feeder availability metrics.
type FeedersMetricSummary struct {
	TotalFeeders         int `json:"total_feeders"`
	InServiceCount       int `json:"in_service_count"`
	OutOnFaultCount      int `json:"out_on_fault_count"`
	OutOnPTWCount        int `json:"out_on_ptw_count"`
	OutOnRationingCount  int `json:"out_on_rationing_count"`
	SubstationStaffCount int `json:"substation_staff_count"`
}

// RevenuePulseSummary aggregates financial collection channels.
type RevenuePulseSummary struct {
	TotalCollectedAmount float64            `json:"total_collected_amount"`
	PrepaidAmount        float64            `json:"prepaid_amount"`
	PostpaidAmount       float64            `json:"postpaid_amount"`
	ChannelBreakdown     map[string]float64 `json:"channel_breakdown"`
}

// IntegrityRadarSummary tracks utility loss risks and audit pending items.
type IntegrityRadarSummary struct {
	UnresolvedTheftAlertsCount int `json:"unresolved_theft_alerts_count"`
	HighChurnRiskAccountsCount int `json:"high_churn_risk_accounts_count"`
	AuditPendingItemsCount     int `json:"audit_pending_items_count"`
	UnreturnedScrapCount       int `json:"unreturned_scrap_count"`
}
