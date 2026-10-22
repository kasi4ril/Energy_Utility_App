// Package domain_test verifies wrong billing adjustment commercial approvals,
// and revenue protection AI anomaly algorithms for S08.
package domain_test

import (
	"strings"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
)

func TestValidateAdjustmentProposal(t *testing.T) {
	baseAdj := domain.BillingAdjustment{
		ID:                "adj-001",
		AdjustmentNo:      "ADJ-2026-0001",
		TicketID:          "TICK-001",
		CustomerID:        "CUST-001",
		BillingOfficerID:  "STAFF-BILL-01",
		OldBalance:        15000.0,
		AdjustedBalance:   10000.0,
		AdjustmentAmount:  5000.0,
		RootCauseCategory: domain.RootCauseMeterReadingTypo,
		RootCauseNote:     "Correction of inverted meter reading digits",
		Status:            domain.AdjustmentProposed,
	}

	t.Run("valid adjustment proposal", func(t *testing.T) {
		if err := domain.ValidateAdjustmentProposal(baseAdj); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		adj := baseAdj
		adj.TicketID = ""
		if err := domain.ValidateAdjustmentProposal(adj); err == nil {
			t.Fatal("expected error on missing TicketID, got nil")
		}

		adj = baseAdj
		adj.CustomerID = "   "
		if err := domain.ValidateAdjustmentProposal(adj); err == nil {
			t.Fatal("expected error on empty CustomerID, got nil")
		}

		adj = baseAdj
		adj.BillingOfficerID = ""
		if err := domain.ValidateAdjustmentProposal(adj); err == nil {
			t.Fatal("expected error on missing BillingOfficerID, got nil")
		}

		adj = baseAdj
		adj.RootCauseNote = ""
		if err := domain.ValidateAdjustmentProposal(adj); err == nil {
			t.Fatal("expected error on empty RootCauseNote, got nil")
		}

		adj = baseAdj
		adj.RootCauseCategory = ""
		if err := domain.ValidateAdjustmentProposal(adj); err == nil {
			t.Fatal("expected error on empty RootCauseCategory, got nil")
		}
	})

	t.Run("amount mismatch", func(t *testing.T) {
		adj := baseAdj
		adj.AdjustmentAmount = 4000.0 // Expected 5000
		err := domain.ValidateAdjustmentProposal(adj)
		if err == nil {
			t.Fatal("expected error on amount mismatch, got nil")
		}
		if !strings.Contains(err.Error(), "must match old_balance - adjusted_balance") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
}

func TestValidateCommercialHeadApproval(t *testing.T) {
	adj := domain.BillingAdjustment{
		ID:                "adj-001",
		Status:            domain.AdjustmentProposed,
		OldBalance:        20000.0,
		AdjustedBalance:   12000.0,
		AdjustmentAmount:  8000.0,
		RootCauseCategory: domain.RootCauseTariffBandMismatch,
		RootCauseNote:     "Commercial rate mistakenly applied to residential meter",
	}

	t.Run("invalid initial status", func(t *testing.T) {
		invalidAdj := adj
		invalidAdj.Status = domain.AdjustmentHeadApproved
		if err := domain.ValidateCommercialHeadApproval(&invalidAdj, "HEAD-COMM-01", true, "Looks good"); err == nil {
			t.Fatal("expected error for non-proposed status, got nil")
		}
	})

	t.Run("empty head of commercial operations ID", func(t *testing.T) {
		targetAdj := adj
		if err := domain.ValidateCommercialHeadApproval(&targetAdj, "  ", true, "Approved"); err == nil {
			t.Fatal("expected error on empty head ID, got nil")
		}
	})

	t.Run("approval path", func(t *testing.T) {
		targetAdj := adj
		err := domain.ValidateCommercialHeadApproval(&targetAdj, "HEAD-COMM-01", true, "Adjustment approved and audited")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !targetAdj.HeadCommercialOpsApproved {
			t.Fatal("expected HeadCommercialOpsApproved to be true")
		}
		if targetAdj.Status != domain.AdjustmentAppliedRevertedToCare {
			t.Fatalf("expected status APPLIED_REVERTED_TO_CARE, got %s", targetAdj.Status)
		}
		if targetAdj.HeadCommercialOpsApprovedAt == nil {
			t.Fatal("expected approved at to be populated")
		}
		if *targetAdj.HeadCommercialOpsID != "HEAD-COMM-01" {
			t.Fatalf("expected head ID HEAD-COMM-01, got %s", *targetAdj.HeadCommercialOpsID)
		}
	})

	t.Run("rejection path", func(t *testing.T) {
		targetAdj := adj
		err := domain.ValidateCommercialHeadApproval(&targetAdj, "HEAD-COMM-01", false, "Disputed root cause")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if targetAdj.HeadCommercialOpsApproved {
			t.Fatal("expected HeadCommercialOpsApproved to be false")
		}
		if targetAdj.Status != domain.AdjustmentRejected {
			t.Fatalf("expected status REJECTED, got %s", targetAdj.Status)
		}
		if targetAdj.HeadCommercialOpsApprovedAt == nil {
			t.Fatal("expected approved at to be populated")
		}
	})
}

func TestDetectEnergyTheftBypass(t *testing.T) {
	t.Run("feeder energized with vending drop >= 50%", func(t *testing.T) {
		detected, conf, summary := domain.DetectEnergyTheftBypass(0, 0, 65.0, true)
		if !detected {
			t.Fatal("expected theft detected to be true")
		}
		if conf != 0.880 {
			t.Fatalf("expected confidence 0.880, got %f", conf)
		}
		if !strings.Contains(summary, "Energy Theft Bypass Flag") {
			t.Fatalf("expected summary to contain flag, got: %s", summary)
		}
	})

	t.Run("feeder not energized with vending drop does not trigger theft flag", func(t *testing.T) {
		detected, conf, _ := domain.DetectEnergyTheftBypass(0, 0, 65.0, false)
		if detected {
			t.Fatal("expected theft detected to be false when feeder not energized")
		}
		if conf != 0.0 {
			t.Fatalf("expected confidence 0.0, got %f", conf)
		}
	})

	t.Run("transformer level loss >= 30%", func(t *testing.T) {
		detected, conf, summary := domain.DetectEnergyTheftBypass(1000.0, 600.0, 10.0, false)
		if !detected {
			t.Fatal("expected DT bypass detected to be true")
		}
		if conf != 0.790 {
			t.Fatalf("expected confidence 0.790, got %f", conf)
		}
		if !strings.Contains(summary, "Distribution Transformer Bypass Flag") {
			t.Fatalf("expected summary to contain flag, got: %s", summary)
		}
	})

	t.Run("transformer level loss < 30%", func(t *testing.T) {
		detected, conf, summary := domain.DetectEnergyTheftBypass(1000.0, 800.0, 20.0, false)
		if detected {
			t.Fatal("expected DT bypass detected to be false for < 30% loss")
		}
		if conf != 0.0 || summary != "" {
			t.Fatalf("expected empty summary and 0.0 confidence, got %f, %s", conf, summary)
		}
	})
}

func TestEvaluateCustomerChurnRisk(t *testing.T) {
	t.Run("high churn risk triggers with high outages and high purchase drop", func(t *testing.T) {
		detected, prob, summary := domain.EvaluateCustomerChurnRisk(15, 55.0)
		if !detected {
			t.Fatal("expected churn risk to be true")
		}
		if prob != 0.850 {
			t.Fatalf("expected prob 0.850, got %f", prob)
		}
		if !strings.Contains(summary, "High Customer Churn Risk") {
			t.Fatalf("expected summary to contain high churn risk, got %s", summary)
		}
	})

	t.Run("low churn risk under threshold", func(t *testing.T) {
		detected, prob, summary := domain.EvaluateCustomerChurnRisk(5, 20.0)
		if detected {
			t.Fatal("expected churn risk to be false")
		}
		if prob != 0.10 {
			t.Fatalf("expected prob 0.10, got %f", prob)
		}
		if summary != "" {
			t.Fatalf("expected empty summary, got %s", summary)
		}
	})
}

