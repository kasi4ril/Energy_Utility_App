// Package repository_test verifies the persistence operations for vending and meter refund ledgers.
package repository_test

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
)

func TestMemoryRepository_VendingTransactions(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()

	feederID := "FDR-01"
	dtID := "DT-01"
	tx := domain.VendingTransaction{
		ID:                     "tx-001",
		TxReference:            "VEND-2026-0001",
		CustomerID:             "CUST-001",
		AccountNumber:          "ACC-001",
		MeterNumber:            "01234567890",
		FeederID:               &feederID,
		DTID:                   &dtID,
		PaymentChannelCategory: domain.ChannelOnlinePortal,
		PaymentPlatformSource:  "PAYSTACK",
		AmountPaid:             5000.0,
		TariffRatePerKWh:       50.0,
		KWhVended:              100.0,
		MeterRefundCredited:    0.0,
		STSToken:               "1111-2222-3333-4444-5555",
		SettlementVerified:     true,
		GatewayPayloadJSON:     `{"status":"success"}`,
		CreatedAt:              time.Now().UTC(),
	}

	// Create
	if err := repo.CreateVendingTransaction(ctx, tx); err != nil {
		t.Fatalf("failed to create vending tx: %v", err)
	}

	// Conflict ID
	if err := repo.CreateVendingTransaction(ctx, tx); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate ID, got %v", err)
	}

	// Conflict Ref
	tx2 := tx
	tx2.ID = "tx-002"
	if err := repo.CreateVendingTransaction(ctx, tx2); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate tx_reference, got %v", err)
	}

	// Get by ID
	fetched, err := repo.GetVendingTransactionByID(ctx, "tx-001")
	if err != nil {
		t.Fatalf("failed to get tx by ID: %v", err)
	}
	if fetched.TxReference != tx.TxReference {
		t.Fatalf("expected ref %s, got %s", tx.TxReference, fetched.TxReference)
	}

	// Get by ID - Not Found
	if _, err := repo.GetVendingTransactionByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Get by Ref
	fetchedByRef, err := repo.GetVendingTransactionByRef(ctx, "VEND-2026-0001")
	if err != nil {
		t.Fatalf("failed to get tx by ref: %v", err)
	}
	if fetchedByRef.ID != tx.ID {
		t.Fatalf("expected ID %s, got %s", tx.ID, fetchedByRef.ID)
	}

	// Get by Ref - Not Found
	if _, err := repo.GetVendingTransactionByRef(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// List
	list, err := repo.ListVendingTransactions(ctx, "CUST-001", "01234567890", 10)
	if err != nil {
		t.Fatalf("failed to list transactions: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 tx in list, got %d", len(list))
	}

	// Filter mismatch
	emptyList, err := repo.ListVendingTransactions(ctx, "CUST-OTHER", "", 10)
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(emptyList) != 0 {
		t.Fatalf("expected 0 tx, got %d", len(emptyList))
	}
}

func TestMemoryRepository_MeterRefundLedgers(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()

	ledger := domain.MeterRefundLedger{
		ID:                     "ref-001",
		CustomerID:             "CUST-001",
		MeterSerial:            "MTR-9999",
		InitialMeterCost:       50000.0,
		RefundPercentagePerTx:  10.0,
		TotalRefundedToDate:    0.0,
		RemainingRefundBalance: 50000.0,
		Status:                 domain.RefundActive,
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	}

	// Create
	if err := repo.CreateRefundLedger(ctx, ledger); err != nil {
		t.Fatalf("failed to create refund ledger: %v", err)
	}

	// Conflict ID
	if err := repo.CreateRefundLedger(ctx, ledger); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict, got %v", err)
	}

	// Conflict Meter
	l2 := ledger
	l2.ID = "ref-002"
	if err := repo.CreateRefundLedger(ctx, l2); err != repository.ErrConflict {
		t.Fatalf("expected ErrConflict on duplicate meter serial, got %v", err)
	}

	// Get by ID
	fetched, err := repo.GetRefundLedgerByID(ctx, "ref-001")
	if err != nil {
		t.Fatalf("failed to get ledger by ID: %v", err)
	}
	if fetched.MeterSerial != "MTR-9999" {
		t.Fatalf("expected meter serial MTR-9999, got %s", fetched.MeterSerial)
	}

	// Get by ID - Not Found
	if _, err := repo.GetRefundLedgerByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Get by Customer
	byCust, err := repo.GetRefundLedgerByCustomer(ctx, "CUST-001")
	if err != nil {
		t.Fatalf("failed to get ledger by customer: %v", err)
	}
	if byCust.ID != "ref-001" {
		t.Fatalf("expected ID ref-001, got %s", byCust.ID)
	}

	// Get by Customer - Not Found
	if _, err := repo.GetRefundLedgerByCustomer(ctx, "UNKNOWN"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Get by Meter
	byMeter, err := repo.GetRefundLedgerByMeter(ctx, "MTR-9999")
	if err != nil {
		t.Fatalf("failed to get ledger by meter: %v", err)
	}
	if byMeter.CustomerID != "CUST-001" {
		t.Fatalf("expected customer CUST-001, got %s", byMeter.CustomerID)
	}

	// Get by Meter - Not Found
	if _, err := repo.GetRefundLedgerByMeter(ctx, "UNKNOWN"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Update
	ledger.TotalRefundedToDate = 5000.0
	ledger.RemainingRefundBalance = 45000.0
	if err := repo.UpdateRefundLedger(ctx, ledger); err != nil {
		t.Fatalf("failed to update ledger: %v", err)
	}

	updated, _ := repo.GetRefundLedgerByID(ctx, "ref-001")
	if updated.RemainingRefundBalance != 45000.0 {
		t.Fatalf("expected 45000, got %f", updated.RemainingRefundBalance)
	}

	// Update - Not Found
	nonExistent := ledger
	nonExistent.ID = "non-existent"
	if err := repo.UpdateRefundLedger(ctx, nonExistent); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound on updating non-existent, got %v", err)
	}
}

func TestNewPostgresRepository(t *testing.T) {
	pg := repository.NewPostgresRepository(nil)
	if pg == nil {
		t.Fatal("expected non-nil PostgresRepository")
	}
}
