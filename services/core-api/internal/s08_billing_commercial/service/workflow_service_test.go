// Package service_test exercises vending payments, STS token generation, and customer meter refund amortization.
package service_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/service"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

func setupTestService() (*service.WorkflowService, *repository.MemoryRepository, *events.MemoryEventBus, *s03repo.MemoryRepository) {
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	ticketRepo := s03repo.NewMemoryRepository()
	svc := service.NewWorkflowService(repo, bus, ticketRepo)
	return svc, repo, bus, ticketRepo
}

func TestWorkflowService_RecordVendingPayment_Standard(t *testing.T) {
	svc, _, bus, _ := setupTestService()
	ctx := context.Background()

	eventCh := make(chan struct{}, 1)
	_, _ = bus.Subscribe("billing.payment_received", func(_ context.Context, ev events.Event) error {
		select {
		case eventCh <- struct{}{}:
		default:
		}
		return nil
	})

	req := service.RecordPaymentRequest{
		CustomerID:             "CUST-001",
		AccountNumber:          "ACC-001",
		MeterNumber:            "01445566778",
		PaymentChannelCategory: domain.ChannelOnlinePortal,
		PaymentPlatformSource:  "PAYSTACK",
		AmountPaid:             10000.0,
		TariffRatePerKWh:       50.0,
		SettlementVerified:     true,
	}

	tx, err := svc.RecordVendingPayment(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error recording payment: %v", err)
	}

	if tx.KWhVended != 200.0 {
		t.Fatalf("expected 200 kWh, got %f", tx.KWhVended)
	}
	if tx.MeterRefundCredited != 0.0 {
		t.Fatalf("expected 0 refund, got %f", tx.MeterRefundCredited)
	}

	stsRegex := regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{4}-\d{4}$`)
	if !stsRegex.MatchString(tx.STSToken) {
		t.Fatalf("invalid STS token format: %s", tx.STSToken)
	}

	select {
	case <-eventCh:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for billing.payment_received event")
	}

	// Fetch by ID
	fetched, err := svc.GetVendingTransaction(ctx, tx.ID)
	if err != nil {
		t.Fatalf("failed to get tx: %v", err)
	}
	if fetched.TxReference != tx.TxReference {
		t.Fatalf("expected ref %s, got %s", tx.TxReference, fetched.TxReference)
	}

	// List
	list, err := svc.ListVendingTransactions(ctx, "CUST-001", "01445566778", 10)
	if err != nil {
		t.Fatalf("failed to list tx: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 tx, got %d", len(list))
	}
}

func TestWorkflowService_RecordVendingPayment_WithRefundAmortization(t *testing.T) {
	svc, _, _, _ := setupTestService()
	ctx := context.Background()

	// 1. Register customer meter refund ledger
	ledgerReq := domain.MeterRefundLedger{
		CustomerID:            "CUST-REFUND-01",
		MeterSerial:           "MTR-PREPAID-001",
		InitialMeterCost:      1500.0,
		RefundPercentagePerTx: 10.0, // 10% refund on each vend
	}
	regLedger, err := svc.RegisterMeterRefundLedger(ctx, ledgerReq)
	if err != nil {
		t.Fatalf("failed to register refund ledger: %v", err)
	}
	if regLedger.Status != domain.RefundActive {
		t.Fatalf("expected status ACTIVE, got %s", regLedger.Status)
	}

	// 2. Vend 1: 10,000 @ 50 NGN/kWh -> 10% is 1,000 refund credit
	tx1, err := svc.RecordVendingPayment(ctx, service.RecordPaymentRequest{
		CustomerID:             "CUST-REFUND-01",
		AccountNumber:          "ACC-REFUND-01",
		MeterNumber:            "MTR-PREPAID-001",
		PaymentChannelCategory: domain.ChannelFintechAgency,
		PaymentPlatformSource:  "OPAY",
		AmountPaid:             10000.0,
		TariffRatePerKWh:       50.0,
	})
	if err != nil {
		t.Fatalf("failed vend 1: %v", err)
	}
	if tx1.MeterRefundCredited != 1000.0 {
		t.Fatalf("expected 1000 refund credit, got %f", tx1.MeterRefundCredited)
	}
	// Total value = 10,000 + 1,000 = 11,000 / 50 = 220 kWh
	if tx1.KWhVended != 220.0 {
		t.Fatalf("expected 220 kWh, got %f", tx1.KWhVended)
	}

	// Verify remaining ledger balance: 1500 - 1000 = 500
	l1, err := svc.GetRefundLedgerByCustomer(ctx, "CUST-REFUND-01")
	if err != nil {
		t.Fatalf("failed to get ledger: %v", err)
	}
	if l1.RemainingRefundBalance != 500.0 {
		t.Fatalf("expected remaining 500, got %f", l1.RemainingRefundBalance)
	}
	if l1.Status != domain.RefundActive {
		t.Fatalf("expected active status, got %s", l1.Status)
	}

	// 3. Vend 2: 10,000 @ 50 NGN/kWh -> 10% would be 1,000, but only 500 remaining!
	tx2, err := svc.RecordVendingPayment(ctx, service.RecordPaymentRequest{
		CustomerID:             "CUST-REFUND-01",
		AccountNumber:          "ACC-REFUND-01",
		MeterNumber:            "MTR-PREPAID-001",
		PaymentChannelCategory: domain.ChannelCommercialBank,
		PaymentPlatformSource:  "ZENITH_BANK",
		AmountPaid:             10000.0,
		TariffRatePerKWh:       50.0,
	})
	if err != nil {
		t.Fatalf("failed vend 2: %v", err)
	}
	if tx2.MeterRefundCredited != 500.0 {
		t.Fatalf("expected 500 refund credit (capped), got %f", tx2.MeterRefundCredited)
	}
	// Total value = 10,000 + 500 = 10,500 / 50 = 210 kWh
	if tx2.KWhVended != 210.0 {
		t.Fatalf("expected 210 kWh, got %f", tx2.KWhVended)
	}

	// Verify ledger is now AMORTIZED_COMPLETED
	l2, err := svc.GetRefundLedgerByCustomer(ctx, "CUST-REFUND-01")
	if err != nil {
		t.Fatalf("failed to get ledger: %v", err)
	}
	if l2.RemainingRefundBalance != 0.0 {
		t.Fatalf("expected 0 remaining, got %f", l2.RemainingRefundBalance)
	}
	if l2.Status != domain.RefundAmortizedCompleted {
		t.Fatalf("expected AMORTIZED_COMPLETED, got %s", l2.Status)
	}

	// 4. Vend 3: Subsequent vend after completion gets 0 refund
	tx3, err := svc.RecordVendingPayment(ctx, service.RecordPaymentRequest{
		CustomerID:             "CUST-REFUND-01",
		AccountNumber:          "ACC-REFUND-01",
		MeterNumber:            "MTR-PREPAID-001",
		PaymentChannelCategory: domain.ChannelCommercialBank,
		PaymentPlatformSource:  "FIRST_BANK",
		AmountPaid:             10000.0,
		TariffRatePerKWh:       50.0,
	})
	if err != nil {
		t.Fatalf("failed vend 3: %v", err)
	}
	if tx3.MeterRefundCredited != 0.0 {
		t.Fatalf("expected 0 refund credit after completion, got %f", tx3.MeterRefundCredited)
	}
	if tx3.KWhVended != 200.0 {
		t.Fatalf("expected 200 kWh, got %f", tx3.KWhVended)
	}
}

func TestWorkflowService_ValidationErrors(t *testing.T) {
	svc, _, _, _ := setupTestService()
	ctx := context.Background()

	// Missing customer/account/meter
	_, err := svc.RecordVendingPayment(ctx, service.RecordPaymentRequest{
		CustomerID: "",
		AmountPaid: 5000.0,
	})
	if err == nil {
		t.Fatal("expected error on empty customer ID, got nil")
	}

	// Invalid amount paid
	_, err = svc.RecordVendingPayment(ctx, service.RecordPaymentRequest{
		CustomerID:    "CUST-001",
		AccountNumber: "ACC-001",
		MeterNumber:   "MTR-001",
		AmountPaid:    -100.0,
	})
	if err == nil {
		t.Fatal("expected error on negative amount paid, got nil")
	}

	// Refund ledger validation
	_, err = svc.RegisterMeterRefundLedger(ctx, domain.MeterRefundLedger{
		CustomerID:  "",
		MeterSerial: "MTR-001",
	})
	if err == nil {
		t.Fatal("expected error on empty customer ID, got nil")
	}

	_, err = svc.RegisterMeterRefundLedger(ctx, domain.MeterRefundLedger{
		CustomerID:       "CUST-001",
		MeterSerial:      "MTR-001",
		InitialMeterCost: -50.0,
	})
	if err == nil {
		t.Fatal("expected error on negative initial meter cost, got nil")
	}
}

