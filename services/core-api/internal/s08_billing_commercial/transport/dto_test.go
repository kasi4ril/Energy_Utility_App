// Package transport_test tests request/response DTO validations for S08.
package transport_test

import (
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/transport"
)

func TestVendPaymentDTO_Validation(t *testing.T) {
	valid := transport.VendPaymentDTO{
		CustomerID:             "CUST-001",
		AccountNumber:          "ACC-001",
		MeterNumber:            "01234567890",
		PaymentChannelCategory: domain.ChannelOnlinePortal,
		PaymentPlatformSource:  "PAYSTACK",
		AmountPaid:             5000.0,
		TariffRatePerKWh:       50.0,
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO to pass, got: %v", err)
	}

	// Empty customer ID
	invalid := valid
	invalid.CustomerID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty customer ID, got nil")
	}

	// Empty platform source
	invalid = valid
	invalid.PaymentPlatformSource = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty platform source, got nil")
	}

	// Empty channel category
	invalid = valid
	invalid.PaymentChannelCategory = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty channel, got nil")
	}

	// Non-positive amount
	invalid = valid
	invalid.AmountPaid = -10.0
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on negative amount, got nil")
	}
}

func TestRegisterRefundLedgerDTO_Validation(t *testing.T) {
	valid := transport.RegisterRefundLedgerDTO{
		CustomerID:            "CUST-001",
		MeterSerial:           "MTR-001",
		InitialMeterCost:      50000.0,
		RefundPercentagePerTx: 10.0,
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	invalid := valid
	invalid.MeterSerial = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty meter serial, got nil")
	}

	invalid = valid
	invalid.InitialMeterCost = 0.0
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on zero cost, got nil")
	}
}

func TestProposeAdjustmentDTO_Validation(t *testing.T) {
	valid := transport.ProposeAdjustmentDTO{
		TicketID:          "TICK-001",
		CustomerID:        "CUST-001",
		BillingOfficerID:  "STAFF-01",
		OldBalance:        15000.0,
		AdjustedBalance:   10000.0,
		AdjustmentAmount:  5000.0,
		RootCauseCategory: domain.RootCauseMultiplierError,
		RootCauseNote:     "Fixed incorrect multiplier",
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	invalid := valid
	invalid.TicketID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on missing ticket ID, got nil")
	}

	invalid = valid
	invalid.RootCauseCategory = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty category, got nil")
	}

	invalid = valid
	invalid.AdjustmentAmount = 4000.0
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on amount mismatch, got nil")
	}
}

func TestDecideAdjustmentDTO_Validation(t *testing.T) {
	valid := transport.DecideAdjustmentDTO{
		HeadCommercialOpsID: "HEAD-01",
		Approve:             true,
		Notes:               "Audited and approved",
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	invalid := valid
	invalid.HeadCommercialOpsID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty head ID, got nil")
	}

	invalid = valid
	invalid.Notes = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty notes, got nil")
	}
}

func TestRevenueScanDTO_Validation(t *testing.T) {
	valid := transport.RevenueScanDTO{
		CustomerOrDTID: "DT-001",
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	invalid := valid
	invalid.CustomerOrDTID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty customer or DT ID, got nil")
	}
}

func TestReviewAIAlertDTO_Validation(t *testing.T) {
	valid := transport.ReviewAIAlertDTO{
		Status: domain.AlertInvestigationDispatched,
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO, got: %v", err)
	}

	invalid := valid
	invalid.Status = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected error on empty status, got nil")
	}
}

