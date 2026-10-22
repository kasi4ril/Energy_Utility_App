// Package domain_test verifies direct ticket routing, customer closure gates, and anti-ghost locks.
package domain_test

import (
	"testing"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
)

func TestResolveDirectRouting(t *testing.T) {
	// Fallen pole -> Feeder Unit Head, supervised by Head of Services
	targetFeeder := domain.ResolveDirectRouting(domain.CategoryFallenPole)
	if targetFeeder.TargetUnitHeadRole != string(s01domain.RoleUnitHeadFeeder) ||
		targetFeeder.SupervisingDeptRole != string(s01domain.RoleHeadServices) ||
		targetFeeder.SLADurationHours != 4 {
		t.Fatalf("unexpected routing for fallen pole: %+v", targetFeeder)
	}

	// Bad transformer -> Transformer Unit Head, supervised by Head of Services
	targetTx := domain.ResolveDirectRouting(domain.CategoryTransformerFault)
	if targetTx.TargetUnitHeadRole != string(s01domain.RoleUnitHeadTx) ||
		targetTx.SupervisingDeptRole != string(s01domain.RoleHeadServices) ||
		targetTx.SLADurationHours != 6 {
		t.Fatalf("unexpected routing for transformer fault: %+v", targetTx)
	}

	// Meter schedule -> Metering Unit Head
	targetMeter := domain.ResolveDirectRouting(domain.CategoryNewMeterSchedule)
	if targetMeter.TargetUnitHeadRole != string(s01domain.RoleUnitHeadMetering) ||
		targetMeter.SupervisingDeptRole != string(s01domain.RoleHeadServices) ||
		targetMeter.SLADurationHours != 24 {
		t.Fatalf("unexpected routing for new meter: %+v", targetMeter)
	}

	// Wrong billing -> Billing Unit Head, supervised by Head of Commercial Operations
	targetBilling := domain.ResolveDirectRouting(domain.CategoryWrongBilling)
	if targetBilling.TargetUnitHeadRole != string(s01domain.RoleUnitHeadBilling) ||
		targetBilling.SupervisingDeptRole != string(s01domain.RoleHeadCommercial) ||
		targetBilling.SLADurationHours != 24 {
		t.Fatalf("unexpected routing for wrong billing: %+v", targetBilling)
	}
}

func TestValidateClosureGate_Scenarios(t *testing.T) {
	ticket := domain.ServiceTicket{
		ID:     "tck-1",
		Status: domain.StatusAttendedTo,
	}

	// 1. Technical staff (e.g. Feeder Unit Head) attempts to close ticket directly -> Forbidden!
	_, err := domain.ValidateClosureGate(ticket, string(s01domain.RoleUnitHeadFeeder), true)
	if err != domain.ErrUnauthorizedClosureRole {
		t.Fatalf("expected ErrUnauthorizedClosureRole for technical unit, got: %v", err)
	}

	// 2. Ticket not in ATTENDED_TO state (e.g. ROUTED) -> Forbidden
	unattendedTicket := ticket
	unattendedTicket.Status = domain.StatusRouted
	_, err = domain.ValidateClosureGate(unattendedTicket, string(s01domain.RoleCustomerCareOp), true)
	if err == nil {
		t.Fatalf("expected error closing unattended ticket")
	}

	// 3. Customer Care calls customer and customer confirms fix -> CLOSED
	finalStatus, err := domain.ValidateClosureGate(ticket, string(s01domain.RoleCustomerCareOp), true)
	if err != nil || finalStatus != domain.StatusClosed {
		t.Fatalf("expected StatusClosed on confirmed resolution, got: %s, err: %v", finalStatus, err)
	}

	// 4. Customer Care calls customer and customer DENIES fix -> REOPENED_ESCALATED
	escalatedStatus, err := domain.ValidateClosureGate(ticket, string(s01domain.RoleCustomerCareOp), false)
	if err != nil || escalatedStatus != domain.StatusReopenedEscalated {
		t.Fatalf("expected StatusReopenedEscalated on denied resolution, got: %s, err: %v", escalatedStatus, err)
	}
}

func TestAntiGhostDTDistance(t *testing.T) {
	// Distribution Transformer at Bukuru (approx 9.795, 8.865)
	dtGPS := s01domain.GPSPoint{Latitude: 9.7950, Longitude: 8.8650}

	// 1. Customer within 150m radius -> Valid
	customerNear := s01domain.GPSPoint{Latitude: 9.7955, Longitude: 8.8655}
	dist, err := domain.ValidateAntiGhostDTDistance(customerNear, dtGPS, 500.0)
	if err != nil || dist > 500.0 {
		t.Fatalf("expected nearby customer to validate, got dist %.1f, err %v", dist, err)
	}

	// 2. Customer 1.2km away -> Anti-Ghost Lock Triggered
	customerFar := s01domain.GPSPoint{Latitude: 9.8050, Longitude: 8.8750}
	distFar, err := domain.ValidateAntiGhostDTDistance(customerFar, dtGPS, 500.0)
	if err == nil {
		t.Fatalf("expected ErrAntiGhostToleranceExceeded for customer 1.2km away, got dist: %.1f", distFar)
	}
}

func TestValidateCustomerOnboarding(t *testing.T) {
	validCust := domain.Customer{
		FullName:      "Yakubu Pam",
		Phone:         "+2348030001122",
		TransformerID: "tx-uuid-1",
		KYCStatus:     domain.KYCVerified,
		LedgerBalance: 0.0,
	}

	if err := domain.ValidateCustomerOnboarding(validCust); err != nil {
		t.Fatalf("valid customer failed onboarding validation: %v", err)
	}

	// Flagged debtor account
	debtorCust := validCust
	debtorCust.KYCStatus = domain.KYCDebtorFlagged
	if err := domain.ValidateCustomerOnboarding(debtorCust); err != domain.ErrDebtorAccountBlocked {
		t.Fatalf("expected ErrDebtorAccountBlocked, got: %v", err)
	}
}

