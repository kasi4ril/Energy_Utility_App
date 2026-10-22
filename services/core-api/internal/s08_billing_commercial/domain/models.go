// Package domain encapsulates core entities, state machine enumerations,
// and pure business logic invariants for the S08 Billing, Commercial Operations & Revenue Protection AI bounded context.
package domain

import (
	"time"
)

// PaymentChannelCategory categorizes customer payment touchpoints.
type PaymentChannelCategory string

const (
	ChannelOnlinePortal   PaymentChannelCategory = "ONLINE_PORTAL"
	ChannelFintechAgency  PaymentChannelCategory = "FINTECH_AGENCY"
	ChannelCommercialBank PaymentChannelCategory = "COMMERCIAL_BANK"
	ChannelPostpaidDirect PaymentChannelCategory = "POSTPAID_DIRECT"
)

// RefundStatus tracks the amortization lifecycle of customer meter purchase refunds.
type RefundStatus string

const (
	RefundActive             RefundStatus = "ACTIVE"
	RefundAmortizedCompleted RefundStatus = "AMORTIZED_COMPLETED"
	RefundPaused             RefundStatus = "PAUSED"
)

// RootCauseCategory categorizes the origin of a wrong billing dispute.
type RootCauseCategory string

const (
	RootCauseMeterReadingTypo    RootCauseCategory = "METER_READING_TYPO"
	RootCauseTariffBandMismatch  RootCauseCategory = "TARIFF_BAND_MISMATCH"
	RootCauseMultiplierError     RootCauseCategory = "MULTIPLIER_ERROR"
	RootCausePaymentUncredited   RootCauseCategory = "PAYMENT_UNCREDITED"
)

// AdjustmentStatus governs the state machine of billing adjustment proposals.
type AdjustmentStatus string

const (
	AdjustmentProposed                AdjustmentStatus = "PROPOSED"
	AdjustmentHeadApproved            AdjustmentStatus = "HEAD_APPROVED"
	AdjustmentRejected                AdjustmentStatus = "REJECTED"
	AdjustmentAppliedRevertedToCare   AdjustmentStatus = "APPLIED_REVERTED_TO_CARE"
)

// AIAlertType specifies revenue protection risk categories.
type AIAlertType string

const (
	AlertEnergyTheftBypass       AIAlertType = "ENERGY_THEFT_BYPASS"
	AlertCustomerChurnRisk       AIAlertType = "CUSTOMER_CHURN_RISK"
	AlertUnremittedCollection    AIAlertType = "UNREMITTED_COLLECTION_LEAK"
)

// AIAlertStatus governs the investigation progression of revenue protection alarms.
type AIAlertStatus string

const (
	AlertOpen                     AIAlertStatus = "OPEN"
	AlertInvestigationDispatched  AIAlertStatus = "INVESTIGATION_DISPATCHED"
	AlertResolvedRecovered        AIAlertStatus = "RESOLVED_REVENUE_RECOVERED"
	AlertFalsePositive            AIAlertStatus = "FALSE_POSITIVE"
)

// VendingTransaction represents a prepaid STS token vending or postpaid settlement transaction.
type VendingTransaction struct {
	ID                     string                 `json:"id"`
	TxReference            string                 `json:"tx_reference"`
	CustomerID             string                 `json:"customer_id"`
	AccountNumber          string                 `json:"account_number"`
	MeterNumber            string                 `json:"meter_number"`
	FeederID               *string                `json:"feeder_id,omitempty"`
	DTID                   *string                `json:"dt_id,omitempty"`
	PaymentChannelCategory PaymentChannelCategory `json:"payment_channel_category"`
	PaymentPlatformSource  string                 `json:"payment_platform_source"`
	AmountPaid             float64                `json:"amount_paid"`
	TariffRatePerKWh       float64                `json:"tariff_rate_per_kwh"`
	KWhVended              float64                `json:"kwh_vended"`
	MeterRefundCredited    float64                `json:"meter_refund_credited"`
	STSToken               string                 `json:"sts_token"`
	SettlementVerified     bool                   `json:"settlement_verified"`
	GatewayPayloadJSON     string                 `json:"gateway_payload_json,omitempty"`
	CreatedAt              time.Time              `json:"created_at"`
}

// MeterRefundLedger models customer-financed meter cost amortization.
type MeterRefundLedger struct {
	ID                     string       `json:"id"`
	CustomerID             string       `json:"customer_id"`
	MeterSerial            string       `json:"meter_serial"`
	InitialMeterCost       float64      `json:"initial_meter_cost"`
	RefundPercentagePerTx  float64      `json:"refund_percentage_per_tx"`
	TotalRefundedToDate    float64      `json:"total_refunded_to_date"`
	RemainingRefundBalance float64      `json:"remaining_refund_balance"`
	Status                 RefundStatus `json:"status"`
	CreatedAt              time.Time    `json:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at"`
}

// BillingAdjustment models the wrong billing investigation and commercial head approval loop.
type BillingAdjustment struct {
	ID                          string            `json:"id"`
	AdjustmentNo                string            `json:"adjustment_no"`
	TicketID                    string            `json:"ticket_id"`
	CustomerID                  string            `json:"customer_id"`
	BillingOfficerID            string            `json:"billing_officer_id"`
	OldBalance                  float64           `json:"old_balance"`
	AdjustedBalance             float64           `json:"adjusted_balance"`
	AdjustmentAmount            float64           `json:"adjustment_amount"`
	RootCauseCategory           RootCauseCategory `json:"root_cause_category"`
	RootCauseNote               string            `json:"root_cause_note"`
	HeadCommercialOpsID         *string           `json:"head_commercial_ops_id,omitempty"`
	HeadCommercialOpsApproved   bool              `json:"head_commercial_ops_approved"`
	HeadCommercialOpsApprovedAt *time.Time        `json:"head_commercial_ops_approved_at,omitempty"`
	HeadCommercialOpsNotes      *string           `json:"head_commercial_ops_notes,omitempty"`
	Status                      AdjustmentStatus  `json:"status"`
	CreatedAt                   time.Time         `json:"created_at"`
	UpdatedAt                   time.Time         `json:"updated_at"`
}

// RevenueProtectionAlert represents AI-detected bypass, off-grid churn, or settlement leaks.
type RevenueProtectionAlert struct {
	ID                string        `json:"id"`
	AlertNo           string        `json:"alert_no"`
	CustomerOrDTID    string        `json:"customer_or_dt_id"`
	EntityType        string        `json:"entity_type"`
	AlertType         AIAlertType   `json:"alert_type"`
	ConfidenceScore   float64       `json:"confidence_score"`
	EvidenceSummary   string        `json:"evidence_summary"`
	AssignedUnitID    *string       `json:"assigned_unit_id,omitempty"`
	InvestigationNotes *string      `json:"investigation_notes,omitempty"`
	Status            AIAlertStatus `json:"status"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

