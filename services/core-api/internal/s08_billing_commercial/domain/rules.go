// Package domain implements zero-estimation metered vending calculations,
// customer meter-purchase refund amortization, wrong billing commercial approval gates,
// and revenue protection AI anomaly algorithms for S08.
package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrInvalidVendingParams indicates invalid payment amount or tariff rate.
	ErrInvalidVendingParams = errors.New("domain: amount_paid and tariff_rate_per_kwh must be greater than zero")
	// ErrIllegalAdjustmentTransition indicates an unauthorized or out-of-sequence billing adjustment state change.
	ErrIllegalAdjustmentTransition = errors.New("domain: illegal billing adjustment state transition")
	// ErrInvalidAdjustmentParams indicates mandatory adjustment fields or calculations are mismatched.
	ErrInvalidAdjustmentParams = errors.New("domain: invalid billing adjustment parameters")
)

// CalculateVendedEnergy calculates kilowatt-hours (kWh) delivered for a prepaid vend or postpaid billing receipt.
// Estimated billing is structurally impossible; energy is mathematically derived strictly from tariff rates and refund credits.
func CalculateVendedEnergy(amountPaid, tariffRatePerKWh, refundCredited float64) (float64, error) {
	if amountPaid <= 0.0 || tariffRatePerKWh <= 0.0 {
		return 0.0, ErrInvalidVendingParams
	}
	if refundCredited < 0.0 {
		return 0.0, fmt.Errorf("%w: refund_credited cannot be negative", ErrInvalidVendingParams)
	}

	totalEnergyValue := amountPaid + refundCredited
	kwh := totalEnergyValue / tariffRatePerKWh
	return kwh, nil
}

// GenerateSTSToken generates a standard 20-digit Standard Transfer Specification (STS) token.
// Formula: SHA-256 derivation grouped into 5 blocks of 4 digits: XXXX-XXXX-XXXX-XXXX-XXXX.
func GenerateSTSToken(meterNumber string, amountPaid float64, timestamp int64) string {
	raw := fmt.Sprintf("%s:%.2f:%d", meterNumber, amountPaid, timestamp)
	h := sha256.Sum256([]byte(raw))
	hexStr := hex.EncodeToString(h[:])

	// Convert hex characters to numeric digits modulo 10
	digits := make([]byte, 20)
	for i := 0; i < 20; i++ {
		digits[i] = '0' + (hexStr[i] % 10)
	}

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		string(digits[0:4]),
		string(digits[4:8]),
		string(digits[8:12]),
		string(digits[12:16]),
		string(digits[16:20]),
	)
}

// ComputeRefundInstallment applies regulatory amortization of customer-financed meter costs.
// When a refund ledger is active, a percentage of the transaction amount is credited until fully amortized.
func ComputeRefundInstallment(ledger *MeterRefundLedger, vendAmount float64) float64 {
	if ledger == nil || ledger.Status != RefundActive || ledger.RemainingRefundBalance <= 0.0 {
		return 0.0
	}
	if vendAmount <= 0.0 {
		return 0.0
	}

	installment := vendAmount * (ledger.RefundPercentagePerTx / 100.0)
	if installment > ledger.RemainingRefundBalance {
		installment = ledger.RemainingRefundBalance
	}

	ledger.TotalRefundedToDate += installment
	ledger.RemainingRefundBalance -= installment

	if ledger.RemainingRefundBalance <= 0.001 {
		ledger.RemainingRefundBalance = 0.0
		ledger.Status = RefundAmortizedCompleted
	}

	ledger.UpdatedAt = time.Now().UTC()
	return installment
}

// ValidateAdjustmentProposal verifies that the proposed adjustment amounts and root causes are sound.
func ValidateAdjustmentProposal(adj BillingAdjustment) error {
	if strings.TrimSpace(adj.TicketID) == "" || strings.TrimSpace(adj.CustomerID) == "" || strings.TrimSpace(adj.BillingOfficerID) == "" {
		return ErrInvalidAdjustmentParams
	}
	if strings.TrimSpace(adj.RootCauseNote) == "" || adj.RootCauseCategory == "" {
		return ErrInvalidAdjustmentParams
	}
	expectedAmount := adj.OldBalance - adj.AdjustedBalance
	// Tolerate minor floating point precision differences
	if (adj.AdjustmentAmount - expectedAmount) > 0.01 || (expectedAmount - adj.AdjustmentAmount) > 0.01 {
		return fmt.Errorf("%w: adjustment amount (%.2f) must match old_balance - adjusted_balance (%.2f)",
			ErrInvalidAdjustmentParams, adj.AdjustmentAmount, expectedAmount)
	}
	return nil
}

// ValidateCommercialHeadApproval records the Head of Commercial Operations digital decision.
// On approval, the status transitions to APPLIED_REVERTED_TO_CARE to trigger Customer Care's closure callback.
func ValidateCommercialHeadApproval(adj *BillingAdjustment, headID string, approve bool, notes string) error {
	if adj.Status != AdjustmentProposed {
		return fmt.Errorf("%w: cannot decide on adjustment in status %s", ErrIllegalAdjustmentTransition, adj.Status)
	}
	if strings.TrimSpace(headID) == "" {
		return fmt.Errorf("%w: head of commercial operations ID is mandatory", ErrIllegalAdjustmentTransition)
	}

	now := time.Now().UTC()
	adj.HeadCommercialOpsID = &headID
	adj.HeadCommercialOpsNotes = &notes
	adj.HeadCommercialOpsApprovedAt = &now

	if approve {
		adj.HeadCommercialOpsApproved = true
		adj.Status = AdjustmentAppliedRevertedToCare
	} else {
		adj.HeadCommercialOpsApproved = false
		adj.Status = AdjustmentRejected
	}
	return nil
}

// DetectEnergyTheftBypass evaluates physical meter vending cadence against feeder availability and DT delivery.
func DetectEnergyTheftBypass(dtBulkKWh, sumMeteredKWh, vendingDropPct float64, feederEnergized bool) (bool, float64, string) {
	// 1. Sudden vending drop > 50% while line is energized indicates direct line bypass or illegal tapping
	if feederEnergized && vendingDropPct >= 50.0 {
		confidence := 0.880
		summary := fmt.Sprintf("Energy Theft Bypass Flag: Meter vending cadence collapsed by %.1f%% while feeder line was fully energized.", vendingDropPct)
		return true, confidence, summary
	}

	// 2. Transformer level discrepancy: DT bulk input significantly exceeds cumulative metered sales (> 30% technical/commercial loss)
	if dtBulkKWh > 0.0 && sumMeteredKWh >= 0.0 {
		lossPct := ((dtBulkKWh - sumMeteredKWh) / dtBulkKWh) * 100.0
		if lossPct >= 30.0 {
			confidence := 0.790
			summary := fmt.Sprintf("Distribution Transformer Bypass Flag: DT bulk delivered (%.1f kWh) exceeds metered sales (%.1f kWh) by %.1f%% loss.", dtBulkKWh, sumMeteredKWh, lossPct)
			return true, confidence, summary
		}
	}

	return false, 0.0, ""
}

// EvaluateCustomerChurnRisk forecasts commercial/residential off-grid migration probability.
func EvaluateCustomerChurnRisk(outageCount int, purchaseDropPct float64) (bool, float64, string) {
	if outageCount > 10 && purchaseDropPct >= 40.0 {
		prob := 0.850
		summary := fmt.Sprintf("High Customer Churn Risk: Customer experienced %d outages and reduced vending by %.1f%%; high probability of solar/generator migration.", outageCount, purchaseDropPct)
		return true, prob, summary
	}
	return false, 0.10, ""
}

