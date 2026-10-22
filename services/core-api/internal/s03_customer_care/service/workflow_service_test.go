// Package service_test verifies Customer Care workflows, Anti-Ghost locks, and closure gates.
package service_test

import (
	"context"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/service"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const (
	testHMACSecret = "s03-service-test-hmac-secret-32b!"
)

func setupTestService() (*service.WorkflowService, *repository.MemoryRepository, *s01repo.MemoryRepository) {
	repo := repository.NewMemoryRepository()
	gridRepo := s01repo.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, gridRepo, bus, testHMACSecret)
	return svc, repo, gridRepo
}

func TestWorkflowService_CustomerOnboarding_SuccessAndAntiGhost(t *testing.T) {
	svc, _, gridRepo := setupTestService()
	ctx := context.Background()

	// 1. Create a Distribution Transformer in grid repo at Jos (9.9281, 8.8921)
	dt := s01domain.GridAsset{
		ID:               "tx-dt-jos-01",
		AssetCode:        "DT-JOS-CENTRAL-01",
		AssetType:        s01domain.AssetTransformerDT,
		OrgUnitID:        "org-unit-1",
		RatedCapacityKVA: 500.0,
		GPSPoint:         s01domain.GPSPoint{Latitude: 9.9281, Longitude: 8.8921},
		OperationalState: s01domain.StateInService,
	}
	_ = gridRepo.CreateGridAsset(ctx, dt)

	// 2. Onboard customer within 100m radius -> Success
	customerGPS := s01domain.GPSPoint{Latitude: 9.9285, Longitude: 8.8923}
	cust, err := svc.OnboardCustomer(
		ctx, "Audu Bako", "+2348021112233", "NIN-100200",
		domain.BuildingResidential, customerGPS, dt.ID, "MTR-SMART-001", s01domain.TariffBandA,
	)
	if err != nil {
		t.Fatalf("failed to onboard valid customer: %v", err)
	}
	if cust.AccountNo == "" || cust.CIN == "" || cust.RowHMACSignature == "" {
		t.Fatalf("customer attributes missing: %+v", cust)
	}

	// 3. Negative: Anti-Ghost Lock triggers when premises GPS is 1.5km away
	farGPS := s01domain.GPSPoint{Latitude: 9.9450, Longitude: 8.9100}
	_, err = svc.OnboardCustomer(
		ctx, "Ghost Customer", "+2348029998877", "NIN-999999",
		domain.BuildingResidential, farGPS, dt.ID, "MTR-SMART-002", s01domain.TariffBandA,
	)
	if err != domain.ErrAntiGhostToleranceExceeded && err == nil {
		t.Fatalf("expected ErrAntiGhostToleranceExceeded for customer 1.5km away")
	}

	// 4. Negative: Missing Transformer -> ErrTransformerNotFound
	_, err = svc.OnboardCustomer(
		ctx, "No TX Customer", "+2348024445566", "NIN-444555",
		domain.BuildingResidential, customerGPS, "non-existent-tx", "MTR-003", s01domain.TariffBandB,
	)
	if err != service.ErrTransformerNotFound {
		t.Fatalf("expected ErrTransformerNotFound, got: %v", err)
	}
}

func TestWorkflowService_TicketDirectRoutingAndClosureGate(t *testing.T) {
	svc, repo, _ := setupTestService()
	ctx := context.Background()

	custID := "cust-100"

	// 1. Log ticket: Fallen Pole -> Routed to Feeder Unit Head with 4-hour SLA
	ticket, err := svc.LogServiceTicket(ctx, &custID, domain.ChannelCallCenter, domain.CategoryFallenPole, "Pole fell across road", "unit-feeder-1")
	if err != nil {
		t.Fatalf("failed to log ticket: %v", err)
	}
	if ticket.Status != domain.StatusRouted || ticket.TicketNo == "" {
		t.Fatalf("invalid ticket status or ticket number: %+v", ticket)
	}

	// 2. Technical crew repairs line and marks ATTENDED_TO with proof
	attended, err := svc.AttendToTicket(ctx, ticket.ID, "Replaced 11kV crossarm and strung AAC conductor", "staff-team-lead-1")
	if err != nil {
		t.Fatalf("failed to attend ticket: %v", err)
	}
	if attended.Status != domain.StatusAttendedTo || attended.AttendedProofNote == "" {
		t.Fatalf("expected StatusAttendedTo, got: %s", attended.Status)
	}

	// 3. Negative: Technical staff attempts to close ticket directly -> Forbidden!
	_, err = svc.ConfirmAndCloseTicket(ctx, ticket.ID, "staff-team-lead-1", "TEAM_LEAD_FEEDER", true, "I finished the work")
	if err != domain.ErrUnauthorizedClosureRole {
		t.Fatalf("expected ErrUnauthorizedClosureRole for technical staff, got: %v", err)
	}

	// 4. Positive: Customer Care calls customer, customer confirms power is restored -> CLOSED
	closed, err := svc.ConfirmAndCloseTicket(ctx, ticket.ID, "officer-cc-1", "CUSTOMER_CARE_OFFICER", true, "Customer confirmed resolution via call")
	if err != nil {
		t.Fatalf("customer care confirmation closure failed: %v", err)
	}
	if closed.Status != domain.StatusClosed || !closed.CustomerCallVerified {
		t.Fatalf("expected StatusClosed with call verified, got: %+v", closed)
	}

	// 5. Test Customer Rejection Branch -> REOPENED_ESCALATED
	ticket2, _ := svc.LogServiceTicket(ctx, &custID, domain.ChannelWalkIn, domain.CategoryTransformerFault, "Transformer smoking", "unit-tx-1")
	_, _ = svc.AttendToTicket(ctx, ticket2.ID, "Replaced blown HV fuse", "staff-tx-1")

	escalated, err := svc.ConfirmAndCloseTicket(ctx, ticket2.ID, "officer-cc-1", "CUSTOMER_CARE_OFFICER", false, "Customer reports voltage is still zero")
	if err != nil {
		t.Fatalf("confirmation call processing failed: %v", err)
	}
	if escalated.Status != domain.StatusReopenedEscalated {
		t.Fatalf("expected StatusReopenedEscalated when customer denies resolution, got: %s", escalated.Status)
	}

	_ = repo
}

func TestWorkflowService_ScheduleMeterActivity(t *testing.T) {
	svc, repo, _ := setupTestService()
	ctx := context.Background()

	cust := domain.Customer{
		ID:        "cust-meter-1",
		AccountNo: "ACC-MTR-1",
		FullName:  "Grace Danladi",
		Phone:     "+2348035556677",
	}
	_ = repo.CreateCustomer(ctx, cust)

	scheduledDate := time.Now().Add(48 * time.Hour)
	sched, err := svc.ScheduleMeterActivity(ctx, cust.ID, nil, domain.ScheduleNewInstallation, scheduledDate)
	if err != nil {
		t.Fatalf("failed to schedule meter activity: %v", err)
	}
	if sched.Status != domain.ScheduleStatusScheduled || sched.CustomerID != cust.ID {
		t.Fatalf("unexpected meter schedule attributes: %+v", sched)
	}

	// Negative: Missing customer
	_, err = svc.ScheduleMeterActivity(ctx, "non-existent-customer", nil, domain.ScheduleRepair, scheduledDate)
	if err == nil {
		t.Fatalf("expected error scheduling for missing customer")
	}
}

