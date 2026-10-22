// Package repository defines persistence contracts and memory repository for S03 Customer Care.
package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
)

var (
	// ErrNotFound is returned when a customer, ticket, or schedule is missing.
	ErrNotFound = errors.New("repository: entity not found")
	// ErrDuplicateKey is returned on unique constraint collision.
	ErrDuplicateKey = errors.New("repository: duplicate unique identifier")
)

// Repository specifies data operations for the customercare schema.
type Repository interface {
	CreateCustomer(ctx context.Context, customer domain.Customer) error
	GetCustomerByID(ctx context.Context, id string) (*domain.Customer, error)
	GetCustomerByAccountNo(ctx context.Context, accountNo string) (*domain.Customer, error)
	GetCustomerByPhone(ctx context.Context, phone string) (*domain.Customer, error)
	ListCustomers(ctx context.Context, transformerID string, limit int) ([]domain.Customer, error)

	CreateTicket(ctx context.Context, ticket domain.ServiceTicket) error
	GetTicketByID(ctx context.Context, id string) (*domain.ServiceTicket, error)
	GetTicketByNo(ctx context.Context, ticketNo string) (*domain.ServiceTicket, error)
	UpdateTicketAttended(ctx context.Context, id string, proofNote string, staffID string) error
	UpdateTicketClosure(ctx context.Context, id string, status domain.TicketStatus, callVerified bool, officerID string, feedback string) error
	ListTickets(ctx context.Context, status domain.TicketStatus, assignedUnitID string, category domain.TicketCategory) ([]domain.ServiceTicket, error)

	CreateMeterSchedule(ctx context.Context, sched domain.MeterSchedule) error
	ListMeterSchedules(ctx context.Context, status domain.ScheduleStatus) ([]domain.MeterSchedule, error)
}

// MemoryRepository provides a concurrent, in-memory repository for S03 testing.
type MemoryRepository struct {
	mu        sync.RWMutex
	customers map[string]domain.Customer
	tickets   map[string]domain.ServiceTicket
	schedules []domain.MeterSchedule
}

// NewMemoryRepository constructs a new in-memory customer care repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		customers: make(map[string]domain.Customer),
		tickets:   make(map[string]domain.ServiceTicket),
		schedules: make([]domain.MeterSchedule, 0),
	}
}

// CreateCustomer stores a new onboarded customer.
func (m *MemoryRepository) CreateCustomer(_ context.Context, c domain.Customer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.customers {
		if existing.AccountNo == c.AccountNo || existing.CIN == c.CIN {
			return ErrDuplicateKey
		}
	}
	m.customers[c.ID] = c
	return nil
}

// GetCustomerByID fetches customer by UUID.
func (m *MemoryRepository) GetCustomerByID(_ context.Context, id string) (*domain.Customer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.customers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &c, nil
}

// GetCustomerByAccountNo fetches customer by account number.
func (m *MemoryRepository) GetCustomerByAccountNo(_ context.Context, accNo string) (*domain.Customer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, c := range m.customers {
		if c.AccountNo == accNo {
			return &c, nil
		}
	}
	return nil, ErrNotFound
}

// GetCustomerByPhone fetches customer by telephone number.
func (m *MemoryRepository) GetCustomerByPhone(_ context.Context, phone string) (*domain.Customer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, c := range m.customers {
		if c.Phone == phone {
			return &c, nil
		}
	}
	return nil, ErrNotFound
}

// ListCustomers returns customers optionally filtered by Distribution Transformer.
func (m *MemoryRepository) ListCustomers(_ context.Context, txID string, limit int) ([]domain.Customer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.Customer, 0)
	for _, c := range m.customers {
		if txID != "" && c.TransformerID != txID {
			continue
		}
		res = append(res, c)
		if limit > 0 && len(res) >= limit {
			break
		}
	}
	return res, nil
}

// CreateTicket records a new omnichannel complaint.
func (m *MemoryRepository) CreateTicket(_ context.Context, t domain.ServiceTicket) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.tickets {
		if existing.TicketNo == t.TicketNo {
			return ErrDuplicateKey
		}
	}
	m.tickets[t.ID] = t
	return nil
}

// GetTicketByID fetches ticket by UUID.
func (m *MemoryRepository) GetTicketByID(_ context.Context, id string) (*domain.ServiceTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.tickets[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &t, nil
}

// GetTicketByNo fetches ticket by ticket number code.
func (m *MemoryRepository) GetTicketByNo(_ context.Context, ticketNo string) (*domain.ServiceTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, t := range m.tickets {
		if t.TicketNo == ticketNo {
			return &t, nil
		}
	}
	return nil, ErrNotFound
}

// UpdateTicketAttended updates technical repair proof and marks ticket ATTENDED_TO.
func (m *MemoryRepository) UpdateTicketAttended(_ context.Context, id string, proofNote string, staffID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tickets[id]
	if !ok {
		return ErrNotFound
	}
	t.Status = domain.StatusAttendedTo
	t.AttendedProofNote = proofNote
	t.AttendedByStaffID = &staffID
	t.UpdatedAt = time.Now()
	m.tickets[id] = t
	return nil
}

// UpdateTicketClosure executes customer confirmation verification and final status transition.
func (m *MemoryRepository) UpdateTicketClosure(_ context.Context, id string, status domain.TicketStatus, callVerified bool, officerID string, feedback string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tickets[id]
	if !ok {
		return ErrNotFound
	}
	t.Status = status
	t.CustomerCallVerified = callVerified
	t.VerifiedByCCOfficerID = &officerID
	t.CustomerFeedbackNotes = feedback
	t.UpdatedAt = time.Now()
	m.tickets[id] = t
	return nil
}

// ListTickets lists tickets with optional status, unit, and category filters.
func (m *MemoryRepository) ListTickets(_ context.Context, status domain.TicketStatus, unitID string, cat domain.TicketCategory) ([]domain.ServiceTicket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.ServiceTicket, 0)
	for _, t := range m.tickets {
		if status != "" && t.Status != status {
			continue
		}
		if unitID != "" && t.AssignedUnitID != unitID {
			continue
		}
		if cat != "" && t.Category != cat {
			continue
		}
		res = append(res, t)
	}
	return res, nil
}

// CreateMeterSchedule records a new meter appointment.
func (m *MemoryRepository) CreateMeterSchedule(_ context.Context, sched domain.MeterSchedule) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.schedules = append(m.schedules, sched)
	return nil
}

// ListMeterSchedules lists meter schedules with optional status filter.
func (m *MemoryRepository) ListMeterSchedules(_ context.Context, status domain.ScheduleStatus) ([]domain.MeterSchedule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.MeterSchedule, 0)
	for _, s := range m.schedules {
		if status != "" && s.Status != status {
			continue
		}
		res = append(res, s)
	}
	return res, nil
}

