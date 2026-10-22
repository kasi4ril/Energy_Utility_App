// Package domain_test verifies zero-estimation metered vending calculations
// and customer meter-purchase refund amortization logic for S08.
package domain_test

import (
	"regexp"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
)

func TestCalculateVendedEnergy(t *testing.T) {
	t.Run("valid vending calculation without refund", func(t *testing.T) {
		kwh, err := domain.CalculateVendedEnergy(10000.0, 50.0, 0.0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if kwh != 200.0 {
			t.Fatalf("expected 200.0 kWh, got %f", kwh)
		}
	})

	t.Run("valid vending calculation with refund credit", func(t *testing.T) {
		kwh, err := domain.CalculateVendedEnergy(10000.0, 50.0, 1000.0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if kwh != 220.0 {
			t.Fatalf("expected 220.0 kWh, got %f", kwh)
		}
	})

	t.Run("error on zero or negative amount paid", func(t *testing.T) {
		if _, err := domain.CalculateVendedEnergy(0.0, 50.0, 0.0); err == nil {
			t.Fatal("expected error on zero amount paid, got nil")
		}
		if _, err := domain.CalculateVendedEnergy(-500.0, 50.0, 0.0); err == nil {
			t.Fatal("expected error on negative amount paid, got nil")
		}
	})

	t.Run("error on zero or negative tariff rate", func(t *testing.T) {
		if _, err := domain.CalculateVendedEnergy(5000.0, 0.0, 0.0); err == nil {
			t.Fatal("expected error on zero tariff rate, got nil")
		}
		if _, err := domain.CalculateVendedEnergy(5000.0, -10.0, 0.0); err == nil {
			t.Fatal("expected error on negative tariff rate, got nil")
		}
	})

	t.Run("error on negative refund credited", func(t *testing.T) {
		if _, err := domain.CalculateVendedEnergy(5000.0, 50.0, -100.0); err == nil {
			t.Fatal("expected error on negative refund, got nil")
		}
	})
}

func TestGenerateSTSToken(t *testing.T) {
	token := domain.GenerateSTSToken("01234567890", 5000.0, 1698000000)
	pattern := regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{4}-\d{4}$`)
	if !pattern.MatchString(token) {
		t.Fatalf("Token %s does not match STS 20-digit pattern", token)
	}
	if len(token) != 24 {
		t.Fatalf("Expected token length 24, got %d", len(token))
	}

	token2 := domain.GenerateSTSToken("01234567890", 5000.0, 1698000000)
	if token != token2 {
		t.Fatalf("Expected identical tokens for same input, got %s and %s", token, token2)
	}

	token3 := domain.GenerateSTSToken("01234567890", 5000.0, 1698000001)
	if token == token3 {
		t.Fatalf("Expected different tokens for different timestamp, got same: %s", token)
	}
}

func TestComputeRefundInstallment(t *testing.T) {
	t.Run("nil ledger returns zero", func(t *testing.T) {
		if inst := domain.ComputeRefundInstallment(nil, 5000.0); inst != 0.0 {
			t.Fatalf("expected 0.0, got %f", inst)
		}
	})

	t.Run("non-active status returns zero", func(t *testing.T) {
		ledger := &domain.MeterRefundLedger{
			Status:                 domain.RefundPaused,
			RemainingRefundBalance: 10000.0,
			RefundPercentagePerTx:  10.0,
		}
		if inst := domain.ComputeRefundInstallment(ledger, 5000.0); inst != 0.0 {
			t.Fatalf("expected 0.0, got %f", inst)
		}
	})

	t.Run("zero remaining balance returns zero", func(t *testing.T) {
		ledger := &domain.MeterRefundLedger{
			Status:                 domain.RefundActive,
			RemainingRefundBalance: 0.0,
			RefundPercentagePerTx:  10.0,
		}
		if inst := domain.ComputeRefundInstallment(ledger, 5000.0); inst != 0.0 {
			t.Fatalf("expected 0.0, got %f", inst)
		}
	})

	t.Run("zero or negative vend amount returns zero", func(t *testing.T) {
		ledger := &domain.MeterRefundLedger{
			Status:                 domain.RefundActive,
			RemainingRefundBalance: 10000.0,
			RefundPercentagePerTx:  10.0,
		}
		if inst := domain.ComputeRefundInstallment(ledger, 0.0); inst != 0.0 {
			t.Fatalf("expected 0.0, got %f", inst)
		}
		if inst := domain.ComputeRefundInstallment(ledger, -100.0); inst != 0.0 {
			t.Fatalf("expected 0.0, got %f", inst)
		}
	})

	t.Run("standard installment amortization", func(t *testing.T) {
		ledger := &domain.MeterRefundLedger{
			ID:                     "ref-001",
			CustomerID:             "cust-001",
			RemainingRefundBalance: 50000.0,
			TotalRefundedToDate:    0.0,
			RefundPercentagePerTx:  10.0,
			Status:                 domain.RefundActive,
		}
		inst := domain.ComputeRefundInstallment(ledger, 10000.0)
		if inst != 1000.0 {
			t.Fatalf("expected 1000.0, got %f", inst)
		}
		if ledger.RemainingRefundBalance != 49000.0 {
			t.Fatalf("expected remaining 49000.0, got %f", ledger.RemainingRefundBalance)
		}
		if ledger.TotalRefundedToDate != 1000.0 {
			t.Fatalf("expected total 1000.0, got %f", ledger.TotalRefundedToDate)
		}
		if ledger.Status != domain.RefundActive {
			t.Fatalf("expected status ACTIVE, got %s", ledger.Status)
		}
	})

	t.Run("installment capped at remaining balance and transitions to completed", func(t *testing.T) {
		ledger := &domain.MeterRefundLedger{
			ID:                     "ref-002",
			CustomerID:             "cust-002",
			RemainingRefundBalance: 500.0,
			TotalRefundedToDate:    49500.0,
			RefundPercentagePerTx:  10.0,
			Status:                 domain.RefundActive,
		}
		inst := domain.ComputeRefundInstallment(ledger, 10000.0)
		if inst != 500.0 {
			t.Fatalf("expected 500.0, got %f", inst)
		}
		if ledger.RemainingRefundBalance != 0.0 {
			t.Fatalf("expected remaining 0.0, got %f", ledger.RemainingRefundBalance)
		}
		if ledger.TotalRefundedToDate != 50000.0 {
			t.Fatalf("expected total 50000.0, got %f", ledger.TotalRefundedToDate)
		}
		if ledger.Status != domain.RefundAmortizedCompleted {
			t.Fatalf("expected status AMORTIZED_COMPLETED, got %s", ledger.Status)
		}
	})
}

