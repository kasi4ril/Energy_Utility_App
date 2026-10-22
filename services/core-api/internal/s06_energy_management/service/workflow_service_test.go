// Package service_test exercises bulk intake logging, Google Forms HMAC webhook verification,
// and priority feeder load shedding in the S06 Energy Management bounded context.
package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/service"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const testHMACSecret = "nesco-dispatch-hmac-integrity-key-32b!"

func setupService() (*service.WorkflowService, *repository.MemoryRepository, *s01repo.MemoryRepository, *events.MemoryEventBus) {
	repo := repository.NewMemoryRepository()
	gridRepo := s01repo.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, gridRepo, bus, testHMACSecret)
	return svc, repo, gridRepo, bus
}

func TestWorkflowService_RecordBulkIntake(t *testing.T) {
	svc, _, _, bus := setupService()
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	bus.Subscribe("energy.bulk_intake_logged", func(_ context.Context, e events.Event) error {
		wg.Done()
		return nil
	})

	log := domain.DispatchTelemetryLog{
		ShiftID:        "SHIFT-20261022-M",
		ShiftWindow:    domain.ShiftMorning,
		GenerationMWIn: 12.5,
		TCNMWIn:        5.0,
		OperatorID:     "op-001",
	}

	created, err := svc.RecordBulkIntake(ctx, log)
	if err != nil {
		t.Fatalf("unexpected error recording bulk intake: %v", err)
	}
	if created.TotalBulkMW != 17.5 {
		t.Fatalf("expected TotalBulkMW 17.5, got %f", created.TotalBulkMW)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty UUID for dispatch log")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for energy.bulk_intake_logged event")
	}

	// Negative bulk intake test
	negLog := domain.DispatchTelemetryLog{
		GenerationMWIn: -1.0,
		TCNMWIn:        5.0,
	}
	if _, err := svc.RecordBulkIntake(ctx, negLog); err != domain.ErrNegativeBulkMW {
		t.Fatalf("expected ErrNegativeBulkMW, got %v", err)
	}
}

func TestWorkflowService_IngestGoogleFormsWebhook(t *testing.T) {
	svc, repo, _, _ := setupService()
	ctx := context.Background()

	payloadJSON := `{
		"shift_id": "SHIFT-20261022-M",
		"shift_window": "MORNING",
		"operator_id": "op-007",
		"generation_mw_in": 14.0,
		"tcn_mw_in": 0.0,
		"substation_rollcall": "[{\"sub\":\"Sub-Bukuru\",\"op\":\"Danladi\"}]",
		"handover_notes": "All lines energized",
		"feeders": [
			{
				"feeder_asset_id": "feeder-ind-1",
				"feeder_code": "FDR-IND-01",
				"current_load_mw": 5.0,
				"current_amps": 260.0,
				"state": "IN_SERVICE",
				"priority_rank": 1,
				"tariff_band": "A"
			},
			{
				"feeder_asset_id": "feeder-res-1",
				"feeder_code": "FDR-RES-01",
				"current_load_mw": 4.0,
				"current_amps": 210.0,
				"state": "IN_SERVICE",
				"priority_rank": 4,
				"tariff_band": "D"
			}
		]
	}`

	// 1. Invalid HMAC
	_, err := svc.IngestGoogleFormsWebhook(ctx, []byte(payloadJSON), "invalid-hex-sig")
	if err != service.ErrInvalidWebhookSignature {
		t.Fatalf("expected ErrInvalidWebhookSignature, got %v", err)
	}

	// 2. Valid HMAC
	validSig, err := crypto.ComputeRowHMAC(testHMACSecret, payloadJSON)
	if err != nil {
		t.Fatalf("failed to compute HMAC: %v", err)
	}

	log, err := svc.IngestGoogleFormsWebhook(ctx, []byte(payloadJSON), validSig)
	if err != nil {
		t.Fatalf("unexpected error ingesting valid webhook: %v", err)
	}
	if log.TotalBulkMW != 14.0 {
		t.Fatalf("expected TotalBulkMW 14.0, got %f", log.TotalBulkMW)
	}

	// Verify feeders were upserted
	f1, err := repo.GetFeederStatusByAssetID(ctx, "feeder-ind-1")
	if err != nil || f1.CurrentLoadMW != 5.0 || f1.TariffBand != s01domain.TariffBandA {
		t.Fatalf("feeder-ind-1 was not persisted properly: %+v, err: %v", f1, err)
	}

	// 3. Malformed JSON
	badSig, _ := crypto.ComputeRowHMAC(testHMACSecret, "not-json")
	if _, err := svc.IngestGoogleFormsWebhook(ctx, []byte("not-json"), badSig); err == nil {
		t.Fatal("expected error on malformed JSON payload")
	}

	// 4. Negative generation in payload
	negPayload := `{"generation_mw_in": -2.0, "tcn_mw_in": 0.0}`
	negSig, _ := crypto.ComputeRowHMAC(testHMACSecret, negPayload)
	if _, err := svc.IngestGoogleFormsWebhook(ctx, []byte(negPayload), negSig); err != domain.ErrNegativeBulkMW {
		t.Fatalf("expected ErrNegativeBulkMW, got %v", err)
	}
}

func TestWorkflowService_RunFeederRationing(t *testing.T) {
	svc, repo, _, _ := setupService()
	ctx := context.Background()

	// 1. No dispatch data error when bulkMW <= 0 and repository is empty
	if _, err := svc.RunFeederRationing(ctx, 0.0); err != service.ErrNoDispatchData {
		t.Fatalf("expected ErrNoDispatchData, got %v", err)
	}

	// Seed Feeders: Band A (5MW, rank 1), Band B (4MW, rank 2), Band E (3MW, rank 5)
	_ = repo.UpsertFeederStatus(ctx, domain.FeederStatus{
		FeederAssetID: "fdr-1",
		FeederCode:    "FDR-A",
		PriorityRank:  1,
		TariffBand:    s01domain.TariffBandA,
		CurrentLoadMW: 5.0,
		State:         domain.FeederInService,
		UpdatedAt:     time.Now().UTC(),
	})
	_ = repo.UpsertFeederStatus(ctx, domain.FeederStatus{
		FeederAssetID: "fdr-2",
		FeederCode:    "FDR-B",
		PriorityRank:  2,
		TariffBand:    s01domain.TariffBandB,
		CurrentLoadMW: 4.0,
		State:         domain.FeederInService,
		UpdatedAt:     time.Now().UTC(),
	})
	_ = repo.UpsertFeederStatus(ctx, domain.FeederStatus{
		FeederAssetID: "fdr-3",
		FeederCode:    "FDR-E",
		PriorityRank:  5,
		TariffBand:    s01domain.TariffBandE,
		CurrentLoadMW: 3.0,
		State:         domain.FeederInService,
		UpdatedAt:     time.Now().UTC(),
	})

	// Total demand = 12MW. Available = 7MW.
	// Rank 1 (5MW) fits (5MW <= 7MW).
	// Rank 2 (4MW) requires 9MW > 7MW -> Shed.
	// Rank 5 (3MW) requires 8MW > 7MW -> Shed.
	decisions, err := svc.RunFeederRationing(ctx, 7.0)
	if err != nil {
		t.Fatalf("unexpected error running rationing: %v", err)
	}
	if len(decisions) != 3 {
		t.Fatalf("expected 3 decisions, got %d", len(decisions))
	}

	fdr1, _ := repo.GetFeederStatusByAssetID(ctx, "fdr-1")
	fdr2, _ := repo.GetFeederStatusByAssetID(ctx, "fdr-2")
	fdr3, _ := repo.GetFeederStatusByAssetID(ctx, "fdr-3")

	if fdr1.State != domain.FeederInService {
		t.Errorf("expected FDR-A to be IN_SERVICE, got %s", fdr1.State)
	}
	if fdr2.State != domain.FeederOutOnLoadShedding {
		t.Errorf("expected FDR-B to be OUT_ON_LOAD_SHEDDING, got %s", fdr2.State)
	}
	if fdr3.State != domain.FeederOutOnLoadShedding {
		t.Errorf("expected FDR-E to be OUT_ON_LOAD_SHEDDING, got %s", fdr3.State)
	}

	// Test rationing with automatic latest dispatch intake
	_ = repo.CreateDispatchLog(ctx, domain.DispatchTelemetryLog{
		ID:          "intake-latest",
		TotalBulkMW: 15.0,
		LoggedAt:    time.Now().UTC(),
	})
	decisionsAuto, err := svc.RunFeederRationing(ctx, 0.0)
	if err != nil {
		t.Fatalf("unexpected error running auto rationing: %v", err)
	}
	// With 15MW, all 3 feeders should be IN_SERVICE
	for _, d := range decisionsAuto {
		if d.TargetState != domain.FeederInService {
			t.Errorf("expected %s to be IN_SERVICE under 15MW bulk, got %s", d.FeederCode, d.TargetState)
		}
	}
}
