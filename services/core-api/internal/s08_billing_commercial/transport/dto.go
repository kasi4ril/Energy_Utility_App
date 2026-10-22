// Package transport provides request/response DTOs and input validations
// for S08 Billing, Commercial Operations & Revenue Protection AI.
package transport

import (
	"errors"
	"fmt"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
)

var (
	// ErrMissingMandatoryFields indicates one or more required DTO fields are empty.
	ErrMissingMandatoryFields = errors.New("transport: mandatory fields missing")
	// ErrInvalidAmount indicates financial amount or tariff rate is non-positive.
	ErrInvalidAmount = errors.New("transport: amount or tariff rate must be greater than zero")
	// ErrInvalidCategory indicates an unrecognized enumeration value.
	ErrInvalidCategory = errors.New("transport: invalid category or status enumeration")
	// ErrAmountMismatch indicates calculated adjustment doesn't match balance diff.
	ErrAmountMismatch = errors.New("transport: adjustment amount must equal old_balance minus adjusted_balance")
)

// VendPaymentDTO defines the payload for POST /api/v1/billing/vending/vend.
type VendPaymentDTO struct {
	TxReference            string                        `json:"tx_reference"`
	CustomerID             string                        `json:"customer_id"`
	AccountNumber          string                        `json:"account_number"`
	MeterNumber            string                        `json:"meter_number"`
	FeederID               *string                       `json:"feeder_id,omitempty"`
	DTID                   *string                       `json:"dt_id,omitempty"`
	PaymentChannelCategory domain.PaymentChannelCategory `json:"payment_channel_category"`
	PaymentPlatformSource  string                        `json:"payment_platform_source"`
	AmountPaid             float64                       `json:"amount_paid"`
	TariffRatePerKWh       float64                       `json:"tariff_rate_per_kwh"`
	SettlementVerified     bool                          `json:"settlement_verified"`
	GatewayPayloadJSON     string                        `json:"gateway_payload_json,omitempty"`
}

// Validate verifies mandatory payment parameters.
func (d VendPaymentDTO) Validate() error {
	if strings.TrimSpace(d.CustomerID) == "" || strings.TrimSpace(d.AccountNumber) == "" || strings.TrimSpace(d.MeterNumber) == "" {
		return fmt.Errorf("%w: customer_id, account_number, and meter_number are required", ErrMissingMandatoryFields)
	}
	if strings.TrimSpace(d.PaymentPlatformSource) == "" {
		return fmt.Errorf("%w: payment_platform_source is required", ErrMissingMandatoryFields)
	}
	if d.PaymentChannelCategory == "" {
		return fmt.Errorf("%w: payment_channel_category is required", ErrInvalidCategory)
	}
	if d.AmountPaid <= 0.0 {
		return fmt.Errorf("%w: amount_paid must be greater than zero", ErrInvalidAmount)
	}
	return nil
}

// RegisterRefundLedgerDTO defines the payload for POST /api/v1/billing/refunds.
type RegisterRefundLedgerDTO struct {
	CustomerID            string  `json:"customer_id"`
	MeterSerial           string  `json:"meter_serial"`
	InitialMeterCost      float64 `json:"initial_meter_cost"`
	RefundPercentagePerTx float64 `json:"refund_percentage_per_tx"`
}

// Validate verifies customer refund ledger registration.
func (d RegisterRefundLedgerDTO) Validate() error {
	if strings.TrimSpace(d.CustomerID) == "" || strings.TrimSpace(d.MeterSerial) == "" {
		return fmt.Errorf("%w: customer_id and meter_serial are required", ErrMissingMandatoryFields)
	}
	if d.InitialMeterCost <= 0.0 {
		return fmt.Errorf("%w: initial_meter_cost must be greater than zero", ErrInvalidAmount)
	}
	return nil
}

// ProposeAdjustmentDTO defines the payload for POST /api/v1/billing/adjustments.
type ProposeAdjustmentDTO struct {
	TicketID          string                   `json:"ticket_id"`
	CustomerID        string                   `json:"customer_id"`
	BillingOfficerID  string                   `json:"billing_officer_id"`
	OldBalance        float64                  `json:"old_balance"`
	AdjustedBalance   float64                  `json:"adjusted_balance"`
	AdjustmentAmount  float64                  `json:"adjustment_amount"`
	RootCauseCategory domain.RootCauseCategory `json:"root_cause_category"`
	RootCauseNote     string                   `json:"root_cause_note"`
}

// Validate verifies wrong billing adjustment proposal parameters.
func (d ProposeAdjustmentDTO) Validate() error {
	if strings.TrimSpace(d.TicketID) == "" || strings.TrimSpace(d.CustomerID) == "" || strings.TrimSpace(d.BillingOfficerID) == "" {
		return fmt.Errorf("%w: ticket_id, customer_id, and billing_officer_id are required", ErrMissingMandatoryFields)
	}
	if strings.TrimSpace(d.RootCauseNote) == "" || d.RootCauseCategory == "" {
		return fmt.Errorf("%w: root_cause_category and root_cause_note are required", ErrMissingMandatoryFields)
	}
	expectedDiff := d.OldBalance - d.AdjustedBalance
	if (d.AdjustmentAmount - expectedDiff) > 0.01 || (expectedDiff - d.AdjustmentAmount) > 0.01 {
		return fmt.Errorf("%w: expected %.2f, got %.2f", ErrAmountMismatch, expectedDiff, d.AdjustmentAmount)
	}
	return nil
}

// DecideAdjustmentDTO defines the payload for POST /api/v1/billing/adjustments/{id}/decide.
type DecideAdjustmentDTO struct {
	HeadCommercialOpsID string `json:"head_commercial_ops_id"`
	Approve             bool   `json:"approve"`
	Notes               string `json:"notes"`
}

// Validate verifies Head of Commercial Operations approval decision.
func (d DecideAdjustmentDTO) Validate() error {
	if strings.TrimSpace(d.HeadCommercialOpsID) == "" {
		return fmt.Errorf("%w: head_commercial_ops_id is required", ErrMissingMandatoryFields)
	}
	if strings.TrimSpace(d.Notes) == "" {
		return fmt.Errorf("%w: review notes are required", ErrMissingMandatoryFields)
	}
	return nil
}

// RevenueScanDTO defines payload for POST /api/v1/billing/revenue-protection/scan.
type RevenueScanDTO struct {
	CustomerOrDTID  string  `json:"customer_or_dt_id"`
	EntityType      string  `json:"entity_type"`
	DTBulkKWh       float64 `json:"dt_bulk_kwh"`
	SumMeteredKWh   float64 `json:"sum_metered_kwh"`
	VendingDropPct  float64 `json:"vending_drop_pct"`
	FeederEnergized bool    `json:"feeder_energized"`
	OutageCount     int     `json:"outage_count"`
	PurchaseDropPct float64 `json:"purchase_drop_pct"`
}

// Validate verifies scan input parameters.
func (d RevenueScanDTO) Validate() error {
	if strings.TrimSpace(d.CustomerOrDTID) == "" {
		return fmt.Errorf("%w: customer_or_dt_id is required", ErrMissingMandatoryFields)
	}
	return nil
}

// ReviewAIAlertDTO defines payload for POST /api/v1/billing/revenue-protection/alerts/{id}/review.
type ReviewAIAlertDTO struct {
	AssignedUnitID     *string              `json:"assigned_unit_id,omitempty"`
	InvestigationNotes *string              `json:"investigation_notes,omitempty"`
	Status             domain.AIAlertStatus `json:"status"`
}

// Validate verifies review parameters.
func (d ReviewAIAlertDTO) Validate() error {
	if d.Status == "" {
		return fmt.Errorf("%w: status is required", ErrMissingMandatoryFields)
	}
	return nil
}

