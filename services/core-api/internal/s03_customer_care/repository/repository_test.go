// Package repository_test validates Customer Care data access operations in memory.
package repository_test

import (
	"context"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
)

func TestMemoryRepository_CustomerOperations(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	cust := domain.Customer{
		ID:               "cust-uuid-1",
		AccountNo:        "ACC-00100",
		CIN:              "CIN-99100",
		FullName:         "Maimuna Yusuf",
		Phone:            "+2348031234567",
		BuildingClass:    domain.BuildingResidential,
		PremisesGPS:      s01domain.GPSPoint{Latitude: 9.92, Longitude: 8.89},
		TransformerID:    "tx-uuid-1",
		MeterSerial:      "MTR-00129",
		TariffBand:       s01domain.TariffBandA,
		LedgerBalance:    0.0,
		KYCStatus:        domain.KYCVerified,
		RowHMACSignature: "hmac-sig-1",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := repo.CreateCustomer(ctx, cust); err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	// Duplicate account number rejection
	if err := repo.CreateCustomer(ctx, cust); err != repository.ErrDuplicateKey {
		t.Fatalf("expected ErrDuplicateKey on duplicate customer")
	}

	// Fetch by account number
	byAcc, err := repo.GetCustomerByAccountNo(ctx, "ACC-00100")
	if err != nil || byAcc.FullName != "Maimuna Yusuf" {
		t.Fatalf("failed to fetch customer by account no: %v", err)
	}

	// Fetch by phone
	byPhone, err := repo.GetCustomerByPhone(ctx, "+2348031234567")
	if err != nil || byPhone.ID != cust.ID {
		t.Fatalf("failed to fetch customer by phone: %v", err)
	}

	// List customers by transformer
	list, err := repo.ListCustomers(ctx, "tx-uuid-1", 10)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 customer under transformer, got: %d", len(list))
	}

	// Missing ID checks
	if _, err := repo.GetCustomerByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing customer ID")
	}
	if _, err := repo.GetCustomerByAccountNo(ctx, "UNKNOWN"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing account number")
	}
	if _, err := repo.GetCustomerByPhone(ctx, "0000000000"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing phone number")
	}
}

func TestMemoryRepository_TicketAndScheduleOperations(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	now := time.Now()
	ticket := domain.ServiceTicket{
		ID:             "tck-uuid-1",
		TicketNo:       "TCK-2026-0001",
		Channel:        domain.ChannelCallCenter,
		Category:       domain.CategoryFallenPole,
		Description:    "High tension line dropped across road",
		AssignedUnitID: "unit-feeder-uuid",
		Status:         domain.StatusLogged,
		SLADueAt:       now.Add(4 * time.Hour),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := repo.CreateTicket(ctx, ticket); err != nil {
		t.Fatalf("failed to create ticket: %v", err)
	}

	// Duplicate ticket number rejection
	if err := repo.CreateTicket(ctx, ticket); err != repository.ErrDuplicateKey {
		t.Fatalf("expected ErrDuplicateKey on duplicate ticket")
	}

	// Fetch by ticket number
	byNo, err := repo.GetTicketByNo(ctx, "TCK-2026-0001")
	if err != nil || byNo.ID != ticket.ID {
		t.Fatalf("failed to fetch ticket by ticket no: %v", err)
	}

	// Mark attended
	err = repo.UpdateTicketAttended(ctx, ticket.ID, "Pole replaced and line re-tensioned", "staff-feeder-1")
	if err != nil {
		t.Fatalf("failed to mark ticket attended: %v", err)
	}

	attended, _ := repo.GetTicketByID(ctx, ticket.ID)
	if attended.Status != domain.StatusAttendedTo || attended.AttendedProofNote == "" {
		t.Fatalf("attended fields not updated: %+v", attended)
	}

	// Customer confirmation & closure
	err = repo.UpdateTicketClosure(ctx, ticket.ID, domain.StatusClosed, true, "officer-cc-1", "Customer confirms line energized")
	if err != nil {
		t.Fatalf("failed to close ticket: %v", err)
	}

	closed, _ := repo.GetTicketByID(ctx, ticket.ID)
	if closed.Status != domain.StatusClosed || !closed.CustomerCallVerified {
		t.Fatalf("closed ticket fields mismatched: %+v", closed)
	}

	// Meter Schedule Operations
	sched := domain.MeterSchedule{
		ID:            "sched-1",
		CustomerID:    "cust-uuid-1",
		ScheduleType:  domain.ScheduleNewInstallation,
		ScheduledDate: now,
		Status:        domain.ScheduleStatusScheduled,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := repo.CreateMeterSchedule(ctx, sched); err != nil {
		t.Fatalf("failed to create meter schedule: %v", err)
	}

	scheds, err := repo.ListMeterSchedules(ctx, domain.ScheduleStatusScheduled)
	if err != nil || len(scheds) != 1 {
		t.Fatalf("expected 1 scheduled appointment, got: %d", len(scheds))
	}
}

