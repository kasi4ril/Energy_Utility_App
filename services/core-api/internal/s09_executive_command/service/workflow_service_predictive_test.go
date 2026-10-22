// Package service_test exercises predictive maintenance evaluations and C-Suite action items.
package service_test

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/service"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

func TestWorkflowService_PredictivePrognostics(t *testing.T) {
	svc, _, bus, maintRepo, _, _, _, _ := setupTestService()
	ctx := context.Background()

	eventCh := make(chan struct{}, 1)
	_, _ = bus.Subscribe("executive.predictive.trigger_raised", func(_ context.Context, ev events.Event) error {
		select {
		case eventCh <- struct{}{}:
		default:
		}
		return nil
	})

	// 1. Hydro Prognosis (Critical low head)
	hydroTrigger, err := svc.EvaluateHydroPrognosis(ctx, service.HydroPrognosisRequest{
		AssetID:           "RES-KURA-FALLS-01",
		InflowM3s:         5.0,
		DischargeM3s:      45.0,
		CurrentHeadM:      14.0,
		TurbineEfficiency: 0.88,
		ForecastDays:      7,
	})
	if err != nil {
		t.Fatalf("unexpected hydro prognosis error: %v", err)
	}
	if hydroTrigger.RiskLevel != domain.RiskCritical {
		t.Fatalf("expected RiskCritical, got %s", hydroTrigger.RiskLevel)
	}

	select {
	case <-eventCh:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for trigger_raised event")
	}

	// 2. Generator Prognosis (Critical vibration -> Auto Work Order dispatched!)
	genTrigger, err := svc.EvaluateGeneratorPrognosis(ctx, service.GeneratorPrognosisRequest{
		AssetID:               "GEN-PLANT1-UNIT01",
		BearingVibrationMmS:   12.8,
		WindingTempC:          125.0,
		LubeOilPressureBar:    1.0,
		RunHoursSinceOverhaul: 8500.0,
		UnitHeadID:            "HEAD-GEN-01",
		TeamLeadID:            "LEAD-MECH-01",
	})
	if err != nil {
		t.Fatalf("unexpected generator error: %v", err)
	}
	if !genTrigger.AutoWorkOrderCreated {
		t.Fatal("expected AutoWorkOrderCreated to be true")
	}
	if genTrigger.LinkedWorkOrderID == nil {
		t.Fatal("expected LinkedWorkOrderID to be populated")
	}

	// Verify work order exists in maintenance repo
	wo, err := maintRepo.GetWorkOrderByID(ctx, *genTrigger.LinkedWorkOrderID)
	if err != nil {
		t.Fatalf("failed to retrieve auto-dispatched work order: %v", err)
	}
	if wo.AssetID != "GEN-PLANT1-UNIT01" {
		t.Fatalf("expected asset GEN-PLANT1-UNIT01, got %s", wo.AssetID)
	}

	// 3. Transformer Prognosis (High load/temp -> Auto Work Order)
	txTrigger, err := svc.EvaluateTransformerPrognosis(ctx, service.TransformerPrognosisRequest{
		AssetID:         "TX-33KV-BUKURU-01",
		LoadPct:         110.0,
		TopOilTempC:     88.0,
		FuseBlows30Days: 2,
	})
	if err != nil {
		t.Fatalf("unexpected transformer error: %v", err)
	}
	if !txTrigger.AutoWorkOrderCreated {
		t.Fatal("expected auto work order created for transformer")
	}

	// 4. Feeder Prognosis (Storm + high trip density -> Auto Patrol Work Order)
	fdrTrigger, err := svc.EvaluateFeederPrognosis(ctx, service.FeederPrognosisRequest{
		AssetID:             "FDR-RAYFIELD-33KV",
		TripCount30Days:     6,
		FaultDensityPerKm:   1.8,
		DaysSinceLastPatrol: 100.0,
		StormForecast:       true,
	})
	if err != nil {
		t.Fatalf("unexpected feeder error: %v", err)
	}
	if !fdrTrigger.AutoWorkOrderCreated {
		t.Fatal("expected auto patrol work order for feeder")
	}

	// 5. Acknowledge Trigger
	ack, err := svc.AcknowledgePredictiveTrigger(ctx, genTrigger.ID, "STAFF-EXEC-01", "WORK_ORDER_DISPATCHED")
	if err != nil {
		t.Fatalf("failed to acknowledge trigger: %v", err)
	}
	if ack.Status != "WORK_ORDER_DISPATCHED" {
		t.Fatalf("expected status WORK_ORDER_DISPATCHED, got %s", ack.Status)
	}
	if ack.AcknowledgedByStaffID == nil || *ack.AcknowledgedByStaffID != "STAFF-EXEC-01" {
		t.Fatal("expected acknowledged_by_staff_id to match")
	}

	// 6. List Triggers
	triggers, err := svc.ListPredictiveTriggers(ctx, "ACTIVE", "", 10)
	if err != nil {
		t.Fatalf("failed to list triggers: %v", err)
	}
	if len(triggers) < 3 {
		t.Fatalf("expected at least 3 active triggers, got %d", len(triggers))
	}
}

func TestWorkflowService_ActionConsole(t *testing.T) {
	svc, _, bus, _, _, _, _, _ := setupTestService()
	ctx := context.Background()

	eventCh := make(chan struct{}, 1)
	_, _ = bus.Subscribe("executive.action_item.decided", func(_ context.Context, ev events.Event) error {
		select {
		case eventCh <- struct{}{}:
		default:
		}
		return nil
	})

	// 1. Enqueue MD Procurement Action
	item := domain.ExecutiveActionItem{
		ActionType:       domain.ActionMDProcurementApproval,
		TargetRole:       domain.RoleMD,
		SourceEntityID:   "REQ-2026-9999",
		SourceEntityRef:  "PO-CRITICAL-001",
		SourceDepartment: "STORES_AUDIT",
		Summary:          "Critical Out-of-Stock Transformer Procurement (₦4,500,000)",
		Amount:           4500000.0,
		Priority:         "CRITICAL",
	}

	enqueued, err := svc.EnqueueActionItem(ctx, item)
	if err != nil {
		t.Fatalf("failed to enqueue action item: %v", err)
	}
	if enqueued.Status != domain.ActionStatusPending {
		t.Fatalf("expected PENDING status, got %s", enqueued.Status)
	}

	// 2. Unauthorized Role Attempt (Head of Commercial Ops cannot approve MD item)
	_, err = svc.DecideActionItem(ctx, enqueued.ID, "STAFF-COMM-01", domain.RoleHeadCommercialOps, true, "Unauthorized approve")
	if err == nil {
		t.Fatal("expected segregation of duties error, got nil")
	}

	// 3. Authorized MD Approval
	decided, err := svc.DecideActionItem(ctx, enqueued.ID, "STAFF-MD-01", domain.RoleMD, true, "Approved for immediate vendor purchase")
	if err != nil {
		t.Fatalf("unexpected decide error: %v", err)
	}
	if decided.Status != domain.ActionStatusApproved {
		t.Fatalf("expected APPROVED status, got %s", decided.Status)
	}

	select {
	case <-eventCh:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for action_item.decided event")
	}

	// 4. List Action Items
	list, err := svc.ListActionItems(ctx, domain.RoleMD, domain.ActionStatusApproved)
	if err != nil {
		t.Fatalf("failed to list items: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 approved item, got %d", len(list))
	}
}

