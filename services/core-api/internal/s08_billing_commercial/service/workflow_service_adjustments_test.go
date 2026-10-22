// Package service_test exercises wrong billing adjustments and revenue protection AI anomaly workflows.
package service_test

import (
	"context"
	"testing"
	"time"

	s03domain "disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/service"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

func TestWorkflowService_BillingAdjustment_ApprovalLoop(t *testing.T) {
	svc, _, bus, ticketRepo := setupTestService()
	ctx := context.Background()

	// 1. Create a customer ticket in Customer Care repo
	custID := "CUST-001"
	ticket := s03domain.ServiceTicket{
		ID:             "TICK-BILLING-001",
		TicketNo:       "TCK-2026-0001",
		CustomerID:     &custID,
		Category:       s03domain.CategoryWrongBilling,
		Description:    "Overbilled by 5,000 NGN on estimated tariff multiplier",
		AssignedUnitID: "UNIT-COMMERCIAL-01",
		Status:         s03domain.StatusLogged,
	}
	_ = ticketRepo.CreateTicket(ctx, ticket)

	// 2. Billing Officer proposes adjustment
	adjProposal := domain.BillingAdjustment{
		TicketID:          "TICK-BILLING-001",
		CustomerID:        "CUST-001",
		BillingOfficerID:  "STAFF-BILL-01",
		OldBalance:        15000.0,
		AdjustedBalance:   10000.0,
		AdjustmentAmount:  5000.0,
		RootCauseCategory: domain.RootCauseMultiplierError,
		RootCauseNote:     "Multiplier was set to 2.0 instead of 1.0",
	}

	proposedCh := make(chan struct{}, 1)
	_, _ = bus.Subscribe("billing.adjustment.proposed", func(_ context.Context, ev events.Event) error {
		select {
		case proposedCh <- struct{}{}:
		default:
		}
		return nil
	})

	proposed, err := svc.ProposeBillingAdjustment(ctx, adjProposal)
	if err != nil {
		t.Fatalf("failed to propose adjustment: %v", err)
	}
	if proposed.Status != domain.AdjustmentProposed {
		t.Fatalf("expected PROPOSED, got %s", proposed.Status)
	}

	select {
	case <-proposedCh:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for billing.adjustment.proposed event")
	}

	// 3. Head of Commercial Operations approves adjustment
	approvedCh := make(chan struct{}, 1)
	_, _ = bus.Subscribe("billing.adjustment.approved", func(_ context.Context, ev events.Event) error {
		select {
		case approvedCh <- struct{}{}:
		default:
		}
		return nil
	})

	approved, err := svc.ApproveBillingAdjustment(ctx, proposed.ID, "HEAD-COMM-01", true, "Audit trail verified against feeder logs")
	if err != nil {
		t.Fatalf("failed to approve adjustment: %v", err)
	}
	if approved.Status != domain.AdjustmentAppliedRevertedToCare {
		t.Fatalf("expected APPLIED_REVERTED_TO_CARE, got %s", approved.Status)
	}
	if !approved.HeadCommercialOpsApproved {
		t.Fatal("expected HeadCommercialOpsApproved to be true")
	}

	select {
	case <-approvedCh:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for billing.adjustment.approved event")
	}

	// 4. VERIFY: Linked Customer Care ticket reverted to ATTENDED_TO!
	updatedTicket, err := ticketRepo.GetTicketByID(ctx, "TICK-BILLING-001")
	if err != nil {
		t.Fatalf("failed to get updated ticket: %v", err)
	}
	if updatedTicket.Status != s03domain.StatusAttendedTo {
		t.Fatalf("expected ticket status to be ATTENDED_TO, got %s", updatedTicket.Status)
	}

	// 5. Fetch and List adjustments
	fetched, err := svc.GetBillingAdjustment(ctx, proposed.ID)
	if err != nil {
		t.Fatalf("failed to get adjustment: %v", err)
	}
	if fetched.AdjustmentAmount != 5000.0 {
		t.Fatalf("expected 5000, got %f", fetched.AdjustmentAmount)
	}

	list, err := svc.ListBillingAdjustments(ctx, domain.AdjustmentAppliedRevertedToCare, "CUST-001")
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 adjustment, got %d", len(list))
	}
}

func TestWorkflowService_BillingAdjustment_RejectionLoop(t *testing.T) {
	svc, _, bus, _ := setupTestService()
	ctx := context.Background()

	adjProposal := domain.BillingAdjustment{
		TicketID:          "TICK-BILLING-002",
		CustomerID:        "CUST-002",
		BillingOfficerID:  "STAFF-BILL-02",
		OldBalance:        20000.0,
		AdjustedBalance:   15000.0,
		AdjustmentAmount:  5000.0,
		RootCauseCategory: domain.RootCauseTariffBandMismatch,
		RootCauseNote:     "Proposed band reclassification",
	}

	proposed, err := svc.ProposeBillingAdjustment(ctx, adjProposal)
	if err != nil {
		t.Fatalf("failed to propose adjustment: %v", err)
	}

	rejectedCh := make(chan struct{}, 1)
	_, _ = bus.Subscribe("billing.adjustment.rejected", func(_ context.Context, ev events.Event) error {
		select {
		case rejectedCh <- struct{}{}:
		default:
		}
		return nil
	})

	rejected, err := svc.ApproveBillingAdjustment(ctx, proposed.ID, "HEAD-COMM-01", false, "Disputed: Customer property verified as commercial welding shop")
	if err != nil {
		t.Fatalf("failed to reject adjustment: %v", err)
	}
	if rejected.Status != domain.AdjustmentRejected {
		t.Fatalf("expected REJECTED, got %s", rejected.Status)
	}
	if rejected.HeadCommercialOpsApproved {
		t.Fatal("expected HeadCommercialOpsApproved to be false")
	}

	select {
	case <-rejectedCh:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for billing.adjustment.rejected event")
	}
}

func TestWorkflowService_RevenueProtectionAI(t *testing.T) {
	svc, _, bus, _ := setupTestService()
	ctx := context.Background()

	alertCh := make(chan struct{}, 1)
	_, _ = bus.Subscribe("billing.revenue_protection.alert_raised", func(_ context.Context, ev events.Event) error {
		select {
		case alertCh <- struct{}{}:
		default:
		}
		return nil
	})

	// 1. Energy Theft Bypass Scan
	bypassAlert, err := svc.RunRevenueProtectionScan(ctx, service.RevenueScanRequest{
		CustomerOrDTID:  "DT-JOS-CENTRAL-01",
		EntityType:      "DISTRIBUTION_TRANSFORMER",
		DTBulkKWh:       10000.0,
		SumMeteredKWh:   6500.0, // 35% loss (>30% threshold)
		VendingDropPct:  0.0,
		FeederEnergized: true,
	})
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	if bypassAlert == nil {
		t.Fatal("expected AI bypass alert to be generated, got nil")
	}
	if bypassAlert.AlertType != domain.AlertEnergyTheftBypass {
		t.Fatalf("expected AlertEnergyTheftBypass, got %s", bypassAlert.AlertType)
	}

	select {
	case <-alertCh:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for alert_raised event")
	}

	// 2. Churn Risk Scan
	churnAlert, err := svc.RunRevenueProtectionScan(ctx, service.RevenueScanRequest{
		CustomerOrDTID:  "CUST-LARGE-COMM-01",
		EntityType:      "CUSTOMER",
		OutageCount:     14,
		PurchaseDropPct: 45.0,
	})
	if err != nil {
		t.Fatalf("unexpected churn scan error: %v", err)
	}
	if churnAlert == nil {
		t.Fatal("expected churn alert, got nil")
	}
	if churnAlert.AlertType != domain.AlertCustomerChurnRisk {
		t.Fatalf("expected AlertCustomerChurnRisk, got %s", churnAlert.AlertType)
	}

	// 3. Normal / Clean Grid Scan (No Anomaly)
	cleanAlert, err := svc.RunRevenueProtectionScan(ctx, service.RevenueScanRequest{
		CustomerOrDTID:  "DT-CLEAN-01",
		EntityType:      "DISTRIBUTION_TRANSFORMER",
		DTBulkKWh:       1000.0,
		SumMeteredKWh:   900.0, // 10% technical loss (< 30%)
		FeederEnergized: true,
		VendingDropPct:  5.0,
		OutageCount:     2,
		PurchaseDropPct: 5.0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cleanAlert != nil {
		t.Fatalf("expected nil alert for clean grid, got: %+v", cleanAlert)
	}

	// 4. Review AI Alert
	unitID := "UNIT-SURVEILLANCE-01"
	notes := "Crew dispatched with clamp meters"
	reviewed, err := svc.ReviewAIAlert(ctx, bypassAlert.ID, &unitID, &notes, domain.AlertInvestigationDispatched)
	if err != nil {
		t.Fatalf("failed to review AI alert: %v", err)
	}
	if reviewed.Status != domain.AlertInvestigationDispatched {
		t.Fatalf("expected status INVESTIGATION_DISPATCHED, got %s", reviewed.Status)
	}
	if *reviewed.AssignedUnitID != unitID {
		t.Fatalf("expected unit %s, got %s", unitID, *reviewed.AssignedUnitID)
	}

	// 5. Fetch and List
	fetchedAlert, err := svc.GetAIAlert(ctx, bypassAlert.ID)
	if err != nil {
		t.Fatalf("failed to get AI alert: %v", err)
	}
	if fetchedAlert.AlertNo != bypassAlert.AlertNo {
		t.Fatalf("expected alert no %s, got %s", bypassAlert.AlertNo, fetchedAlert.AlertNo)
	}

	alertsList, err := svc.ListAIAlerts(ctx, domain.AlertInvestigationDispatched, domain.AlertEnergyTheftBypass)
	if err != nil {
		t.Fatalf("failed to list AI alerts: %v", err)
	}
	if len(alertsList) != 1 {
		t.Fatalf("expected 1 alert in list, got %d", len(alertsList))
	}
}

