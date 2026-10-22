// Package repository defines the persistence interface and an in-memory mock engine
// for the S04 Services & Maintenance bounded context.
package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
)

var (
	// ErrNotFound indicates the requested entity was not located in storage.
	ErrNotFound = errors.New("repository: entity not found")
	// ErrConflict indicates an entity with the specified identifier or order number already exists.
	ErrConflict = errors.New("repository: entity conflict detected")
)

// Repository specifies the persistence contract for work orders and proof logs.
type Repository interface {
	CreateWorkOrder(ctx context.Context, wo domain.MaintenanceWorkOrder) error
	GetWorkOrderByID(ctx context.Context, id string) (*domain.MaintenanceWorkOrder, error)
	GetWorkOrderByOrderNo(ctx context.Context, orderNo string) (*domain.MaintenanceWorkOrder, error)
	UpdateWorkOrderStatus(ctx context.Context, id string, status domain.WorkOrderStatus, reworkNotes *string) error
	ListWorkOrders(ctx context.Context, unit domain.UnitType, status domain.WorkOrderStatus, teamLeadID string) ([]domain.MaintenanceWorkOrder, error)

	CreateProofOfWork(ctx context.Context, proof domain.ProofOfWorkLog) error
	GetProofByWorkOrderID(ctx context.Context, workOrderID string) (*domain.ProofOfWorkLog, error)
	UpdateProofVerification(ctx context.Context, proofID string, verified bool, notes string) error
	ListProofsByTeamLead(ctx context.Context, teamLeadID string, limit int) ([]domain.ProofOfWorkLog, error)
}

// MemoryRepository provides a concurrent, thread-safe in-memory implementation of Repository.
type MemoryRepository struct {
	mu          sync.RWMutex
	workOrders  map[string]domain.MaintenanceWorkOrder
	byOrderNo   map[string]string // orderNo -> ID
	proofLogs   map[string]domain.ProofOfWorkLog
	proofByWOID map[string]string // workOrderID -> proofID
}

// NewMemoryRepository initializes an empty in-memory store.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		workOrders:  make(map[string]domain.MaintenanceWorkOrder),
		byOrderNo:   make(map[string]string),
		proofLogs:   make(map[string]domain.ProofOfWorkLog),
		proofByWOID: make(map[string]string),
	}
}

// CreateWorkOrder persists a new work order.
func (m *MemoryRepository) CreateWorkOrder(_ context.Context, wo domain.MaintenanceWorkOrder) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.workOrders[wo.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.byOrderNo[wo.OrderNo]; exists {
		return ErrConflict
	}

	m.workOrders[wo.ID] = wo
	m.byOrderNo[wo.OrderNo] = wo.ID
	return nil
}

// GetWorkOrderByID retrieves a work order by UUID.
func (m *MemoryRepository) GetWorkOrderByID(_ context.Context, id string) (*domain.MaintenanceWorkOrder, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	wo, exists := m.workOrders[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &wo, nil
}

// GetWorkOrderByOrderNo retrieves a work order by order number.
func (m *MemoryRepository) GetWorkOrderByOrderNo(_ context.Context, orderNo string) (*domain.MaintenanceWorkOrder, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.byOrderNo[orderNo]
	if !exists {
		return nil, ErrNotFound
	}
	wo := m.workOrders[id]
	return &wo, nil
}

// UpdateWorkOrderStatus changes the status and optional rework notes.
func (m *MemoryRepository) UpdateWorkOrderStatus(_ context.Context, id string, status domain.WorkOrderStatus, reworkNotes *string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	wo, exists := m.workOrders[id]
	if !exists {
		return ErrNotFound
	}

	wo.Status = status
	wo.ReworkNotes = reworkNotes
	wo.UpdatedAt = time.Now().UTC()
	m.workOrders[id] = wo
	return nil
}

// ListWorkOrders returns orders matching unit, status, and team lead filters.
func (m *MemoryRepository) ListWorkOrders(_ context.Context, unit domain.UnitType, status domain.WorkOrderStatus, teamLeadID string) ([]domain.MaintenanceWorkOrder, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.MaintenanceWorkOrder
	for _, wo := range m.workOrders {
		if unit != "" && wo.UnitType != unit {
			continue
		}
		if status != "" && wo.Status != status {
			continue
		}
		if teamLeadID != "" && wo.TeamLeadID != teamLeadID {
			continue
		}
		results = append(results, wo)
	}
	return results, nil
}

// CreateProofOfWork stores a new proof log.
func (m *MemoryRepository) CreateProofOfWork(_ context.Context, proof domain.ProofOfWorkLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.proofLogs[proof.ID]; exists {
		return ErrConflict
	}
	m.proofLogs[proof.ID] = proof
	m.proofByWOID[proof.WorkOrderID] = proof.ID
	return nil
}

// GetProofByWorkOrderID returns the proof log associated with a work order.
func (m *MemoryRepository) GetProofByWorkOrderID(_ context.Context, workOrderID string) (*domain.ProofOfWorkLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	proofID, exists := m.proofByWOID[workOrderID]
	if !exists {
		return nil, ErrNotFound
	}
	proof := m.proofLogs[proofID]
	return &proof, nil
}

// UpdateProofVerification updates the verification outcome and notes.
func (m *MemoryRepository) UpdateProofVerification(_ context.Context, proofID string, verified bool, notes string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	proof, exists := m.proofLogs[proofID]
	if !exists {
		return ErrNotFound
	}

	proof.UnitHeadVerified = verified
	proof.UnitHeadNotes = &notes
	now := time.Now().UTC()
	proof.VerifiedAt = &now
	m.proofLogs[proofID] = proof
	return nil
}

// ListProofsByTeamLead returns proofs submitted by a team lead.
func (m *MemoryRepository) ListProofsByTeamLead(_ context.Context, teamLeadID string, limit int) ([]domain.ProofOfWorkLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 {
		limit = 50
	}
	var results []domain.ProofOfWorkLog
	for _, proof := range m.proofLogs {
		if teamLeadID != "" && proof.TeamLeadID != teamLeadID {
			continue
		}
		results = append(results, proof)
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

