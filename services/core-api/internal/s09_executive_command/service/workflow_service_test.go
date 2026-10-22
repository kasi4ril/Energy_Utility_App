// Package service_test exercises C-Suite dashboard aggregation and snapshot capturing.
package service_test

import (
	"context"
	"testing"
	"time"

	s04repo "disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	s05domain "disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
	s05repo "disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	s06domain "disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	s06repo "disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	s07domain "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	s07repo "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	s08domain "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	s08repo "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/service"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

func setupTestService() (*service.WorkflowService, *repository.MemoryRepository, *events.MemoryEventBus, *s04repo.MemoryRepository, *s05repo.MemoryRepository, *s06repo.MemoryRepository, *s07repo.MemoryRepository, *s08repo.MemoryRepository) {
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	maintRepo := s04repo.NewMemoryRepository()
	genRepo := s05repo.NewMemoryRepository()
	energyRepo := s06repo.NewMemoryRepository()
	auditRepo := s07repo.NewMemoryRepository()
	billRepo := s08repo.NewMemoryRepository()

	svc := service.NewWorkflowService(repo, bus, maintRepo, genRepo, energyRepo, auditRepo, billRepo)
	return svc, repo, bus, maintRepo, genRepo, energyRepo, auditRepo, billRepo
}

func TestWorkflowService_GetLiveDashboardSummary(t *testing.T) {
	svc, _, _, _, genRepo, energyRepo, auditRepo, billRepo := setupTestService()
	ctx := context.Background()

	// 1. Populate Generation Telemetry in S05 (Plant 1 generating 21.0 MW)
	_ = genRepo.CreateTelemetryLog(ctx, s05domain.MasterStationTelemetryLog{
		ID:          "tel-001",
		PlantID:     "PLANT-1",
		MWGenerated: 21.0,
		RecordedAt:  time.Now().UTC(),
	})

	// 2. Populate Dispatch & Feeders in S06 (TCN Intake 5.0 MW, 2 Feeders)
	_ = energyRepo.CreateDispatchLog(ctx, s06domain.DispatchTelemetryLog{
		ID:       "dis-001",
		TCNMWIn:  5.0,
		LoggedAt: time.Now().UTC(),
	})
	_ = energyRepo.UpsertFeederStatus(ctx, s06domain.FeederStatus{
		ID:            "fdr-001",
		FeederAssetID: "ASSET-FDR-001",
		FeederCode:    "FDR-001",
		State:         s06domain.FeederInService,
	})
	_ = energyRepo.UpsertFeederStatus(ctx, s06domain.FeederStatus{
		ID:            "fdr-002",
		FeederAssetID: "ASSET-FDR-002",
		FeederCode:    "FDR-002",
		State:         s06domain.FeederOutOnFault,
	})

	// 3. Populate Billing & AI Alerts in S08 (Vending 10,000 NGN)
	_ = billRepo.CreateVendingTransaction(ctx, domainVendingTx("tx-001", 10000.0, "OPAY"))
	_ = billRepo.CreateAIAlert(ctx, s08domain.RevenueProtectionAlert{
		ID:        "alr-001",
		AlertNo:   "ALR-001",
		AlertType: s08domain.AlertEnergyTheftBypass,
		Status:    s08domain.AlertOpen,
	})

	// 4. Populate Out-of-Stock Requisition in S07
	_ = auditRepo.CreateRequisition(ctx, s07domain.MaterialRequisition{
		ID:     "req-001",
		ReqNo:  "REQ-001",
		Status: s07domain.StatusRequisitionMDProcurementPending,
	})

	// EXECUTE Dashboard Summary
	summary, err := svc.GetLiveDashboardSummary(ctx)
	if err != nil {
		t.Fatalf("unexpected summary error: %v", err)
	}

	if summary.TotalGenerationMW != 21.0 {
		t.Fatalf("expected generation 21.0 MW, got %f", summary.TotalGenerationMW)
	}
	if summary.TCNGridIntakeMW != 5.0 {
		t.Fatalf("expected TCN intake 5.0 MW, got %f", summary.TCNGridIntakeMW)
	}
	if summary.TotalAvailablePowerMW != 26.0 {
		t.Fatalf("expected total available power 26.0 MW, got %f", summary.TotalAvailablePowerMW)
	}

	if summary.FeedersSummary.TotalFeeders != 2 {
		t.Fatalf("expected 2 feeders, got %d", summary.FeedersSummary.TotalFeeders)
	}
	if summary.FeedersSummary.InServiceCount != 1 || summary.FeedersSummary.OutOnFaultCount != 1 {
		t.Fatalf("unexpected feeder counts: %+v", summary.FeedersSummary)
	}

	if summary.RevenueSummary.TotalCollectedAmount != 10000.0 {
		t.Fatalf("expected total revenue 10000, got %f", summary.RevenueSummary.TotalCollectedAmount)
	}
	if summary.RevenueSummary.ChannelBreakdown["OPAY"] != 10000.0 {
		t.Fatalf("expected OPay channel 10000, got %f", summary.RevenueSummary.ChannelBreakdown["OPAY"])
	}

	if summary.IntegritySummary.UnresolvedTheftAlertsCount != 1 {
		t.Fatalf("expected 1 theft alert, got %d", summary.IntegritySummary.UnresolvedTheftAlertsCount)
	}
	if summary.IntegritySummary.AuditPendingItemsCount != 1 {
		t.Fatalf("expected 1 audit pending requisition, got %d", summary.IntegritySummary.AuditPendingItemsCount)
	}
}

func TestWorkflowService_CaptureExecutiveSnapshot(t *testing.T) {
	svc, _, bus, _, _, _, _, _ := setupTestService()
	ctx := context.Background()

	eventCh := make(chan struct{}, 1)
	_, _ = bus.Subscribe("executive.snapshot.captured", func(_ context.Context, ev events.Event) error {
		select {
		case eventCh <- struct{}{}:
		default:
		}
		return nil
	})

	snapshot := domain.ExecutiveSnapshot{
		SnapshotDate:          "2026-10-22",
		ShiftWindow:           "DAILY_SUMMARY",
		TotalGenerationMW:     22.5,
		TCNGridIntakeMW:       4.5,
		TotalSystemLoadMW:     26.0,
		FeedersInServiceCount: 14,
		TotalRevenueCollected: 12500000.0,
		CapturedByStaffID:     "STAFF-MD-01",
	}

	captured, err := svc.CaptureExecutiveSnapshot(ctx, snapshot)
	if err != nil {
		t.Fatalf("unexpected snapshot error: %v", err)
	}

	if captured.ID == "" {
		t.Fatal("expected non-empty snapshot ID")
	}

	select {
	case <-eventCh:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for executive.snapshot.captured event")
	}

	// Fetch Latest
	latest, err := svc.GetLatestSnapshot(ctx)
	if err != nil {
		t.Fatalf("failed to get latest snapshot: %v", err)
	}
	if latest.ID != captured.ID {
		t.Fatalf("expected snapshot %s, got %s", captured.ID, latest.ID)
	}

	// List
	list, err := svc.ListSnapshots(ctx, 10)
	if err != nil {
		t.Fatalf("failed to list snapshots: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 snapshot in list, got %d", len(list))
	}
}

func domainVendingTx(id string, amount float64, platform string) s08domain.VendingTransaction {
	return s08domain.VendingTransaction{
		ID:                     id,
		TxReference:            "TX-" + id,
		CustomerID:             "CUST-001",
		AccountNumber:          "ACC-001",
		MeterNumber:            "01234567890",
		PaymentChannelCategory: s08domain.ChannelFintechAgency,
		PaymentPlatformSource:  platform,
		AmountPaid:             amount,
		TariffRatePerKWh:       50.0,
		KWhVended:              amount / 50.0,
		STSToken:               "1111-2222-3333-4444-5555",
		SettlementVerified:     true,
	}
}
