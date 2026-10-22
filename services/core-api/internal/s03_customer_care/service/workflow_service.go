// Package service orchestrates Customer Care onboarding, Anti-Ghost GPS locking,
// direct ticket routing, and the exclusive confirmation closure gate.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

var (
	// ErrTransformerNotFound indicates the specified distribution transformer does not exist.
	ErrTransformerNotFound = errors.New("service: assigned distribution transformer not found in grid assets")
)

// WorkflowService coordinates all customer care, ticketing, and scheduling operations.
type WorkflowService struct {
	repo       repository.Repository
	gridRepo   s01repo.Repository
	bus        events.EventBus
	hmacSecret string
}

// NewWorkflowService constructs a new Customer Care workflow orchestrator.
func NewWorkflowService(repo repository.Repository, gridRepo s01repo.Repository, bus events.EventBus, hmacSecret string) *WorkflowService {
	return &WorkflowService{
		repo:       repo,
		gridRepo:   gridRepo,
		bus:        bus,
		hmacSecret: hmacSecret,
	}
}

// OnboardCustomer registers a new consumer after validating Anti-Ghost GPS proximity and debtor status.
func (s *WorkflowService) OnboardCustomer(ctx context.Context, name, phone, nin string, bClass domain.BuildingClass, premisesGPS s01domain.GPSPoint, txID, meterSerial string, tariff s01domain.TariffBand) (*domain.Customer, error) {
	// 1. Fetch Assigned Distribution Transformer to verify existence and coordinates
	txAsset, err := s.gridRepo.GetGridAssetByID(ctx, txID)
	if err != nil {
		return nil, ErrTransformerNotFound
	}

	// 2. Anti-Ghost GPS Lock: Verify customer premises is within 500m of assigned DT
	if _, err := domain.ValidateAntiGhostDTDistance(premisesGPS, txAsset.GPSPoint, domain.DefaultDTRadiusMeters); err != nil {
		return nil, err
	}

	// 3. Debtor & Duplicate Check
	if existing, err := s.repo.GetCustomerByPhone(ctx, phone); err == nil && existing != nil {
		if existing.KYCStatus == domain.KYCDebtorFlagged {
			return nil, domain.ErrDebtorAccountBlocked
		}
	}

	custID := crypto.MustUUID()
	now := time.Now()
	accNo := fmt.Sprintf("ACC-%s", custID[:8])
	cin := fmt.Sprintf("CIN-%s", custID[9:17])

	hmacPayload := fmt.Sprintf("%s:%s:%s:%s:%.2f", custID, accNo, phone, tariff, 0.0)
	sig, err := crypto.ComputeRowHMAC(s.hmacSecret, hmacPayload)
	if err != nil {
		return nil, err
	}

	cust := domain.Customer{
		ID:               custID,
		AccountNo:        accNo,
		CIN:              cin,
		FullName:         name,
		Phone:            phone,
		NIN:              nin,
		BuildingClass:    bClass,
		PremisesGPS:      premisesGPS,
		TransformerID:    txID,
		MeterSerial:      meterSerial,
		TariffBand:       tariff,
		LedgerBalance:    0.0,
		KYCStatus:        domain.KYCVerified,
		RowHMACSignature: sig,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := domain.ValidateCustomerOnboarding(cust); err != nil {
		return nil, err
	}
	if err := s.repo.CreateCustomer(ctx, cust); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "customercare.customer.onboarded", "s03_customer_care", cust)
	return &cust, nil
}

// LogServiceTicket routes a customer complaint directly to the responsible Unit Head with supervisory radar.
func (s *WorkflowService) LogServiceTicket(ctx context.Context, custID *string, channel domain.TicketChannel, cat domain.TicketCategory, desc, unitID string) (*domain.ServiceTicket, error) {
	routing := domain.ResolveDirectRouting(cat)
	now := time.Now()
	tckID := crypto.MustUUID()
	tckNo := fmt.Sprintf("TCK-%s-%s", now.Format("20060102"), tckID[:6])
	slaDue := now.Add(time.Duration(routing.SLADurationHours) * time.Hour)

	ticket := domain.ServiceTicket{
		ID:             tckID,
		TicketNo:       tckNo,
		CustomerID:     custID,
		Channel:        channel,
		Category:       cat,
		Description:    desc,
		AssignedUnitID: unitID,
		Status:         domain.StatusRouted,
		SLADueAt:       slaDue,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.CreateTicket(ctx, ticket); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "customercare.ticket.routed", "s03_customer_care", map[string]interface{}{
		"ticket_id":             ticket.ID,
		"ticket_no":             ticket.TicketNo,
		"category":              ticket.Category,
		"target_unit_head_role": routing.TargetUnitHeadRole,
		"supervising_dept_role": routing.SupervisingDeptRole,
		"sla_due_at":            ticket.SLADueAt,
	})

	return &ticket, nil
}

// AttendToTicket marks a ticket as attended by technical staff, attaching proof of resolution.
func (s *WorkflowService) AttendToTicket(ctx context.Context, ticketID, proofNote, staffID string) (*domain.ServiceTicket, error) {
	ticket, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if ticket.Status == domain.StatusClosed {
		return nil, errors.New("service: ticket is already closed")
	}

	if err := s.repo.UpdateTicketAttended(ctx, ticketID, proofNote, staffID); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "customercare.ticket.attended", "s03_customer_care", map[string]interface{}{
		"ticket_id":  ticketID,
		"ticket_no":  ticket.TicketNo,
		"proof_note": proofNote,
	})

	return s.repo.GetTicketByID(ctx, ticketID)
}

// ConfirmAndCloseTicket enforces the Customer Care exclusive confirmation call gate.
func (s *WorkflowService) ConfirmAndCloseTicket(ctx context.Context, ticketID, officerID, officerRole string, confirmedByCustomer bool, feedback string) (*domain.ServiceTicket, error) {
	ticket, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	newStatus, err := domain.ValidateClosureGate(*ticket, officerRole, confirmedByCustomer)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateTicketClosure(ctx, ticketID, newStatus, confirmedByCustomer, officerID, feedback); err != nil {
		return nil, err
	}

	eventSubject := "customercare.ticket.closed"
	if newStatus == domain.StatusReopenedEscalated {
		eventSubject = "customercare.ticket.reopened_escalated"
	}

	_ = s.bus.Publish(ctx, eventSubject, "s03_customer_care", map[string]interface{}{
		"ticket_id":       ticketID,
		"status":          newStatus,
		"customer_called": true,
		"confirmed":       confirmedByCustomer,
		"feedback":        feedback,
	})

	return s.repo.GetTicketByID(ctx, ticketID)
}

// ScheduleMeterActivity schedules a customer meter installation or repair.
func (s *WorkflowService) ScheduleMeterActivity(ctx context.Context, custID string, ticketID *string, sType domain.ScheduleType, date time.Time) (*domain.MeterSchedule, error) {
	if _, err := s.repo.GetCustomerByID(ctx, custID); err != nil {
		return nil, errors.New("service: customer not found for meter schedule")
	}

	sched := domain.MeterSchedule{
		ID:            crypto.MustUUID(),
		CustomerID:    custID,
		TicketID:      ticketID,
		ScheduleType:  sType,
		ScheduledDate: date,
		Status:        domain.ScheduleStatusScheduled,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateMeterSchedule(ctx, sched); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "customercare.meter.scheduled", "s03_customer_care", sched)
	return &sched, nil
}

