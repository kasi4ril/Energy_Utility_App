// Package repository defines the persistence contract and thread-safe in-memory mock
// for the S07 Stores, Inventory & Internal Audit bounded context.
package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
)

var (
	// ErrNotFound indicates the requested entity was not found in storage.
	ErrNotFound = errors.New("repository: entity not found")
	// ErrConflict indicates an entity with the specified identifier or business key already exists.
	ErrConflict = errors.New("repository: entity conflict detected")
)

// Repository specifies the persistence operations for inventory, requisitions, scrap, and allowance claims.
type Repository interface {
	CreateStoreItem(ctx context.Context, item domain.StoreItem) error
	GetStoreItemByID(ctx context.Context, id string) (*domain.StoreItem, error)
	GetStoreItemByCode(ctx context.Context, code string) (*domain.StoreItem, error)
	UpdateStoreItem(ctx context.Context, item domain.StoreItem) error
	ListStoreItems(ctx context.Context, category domain.ItemCategory, lowStockOnly bool) ([]domain.StoreItem, error)

	CreateRequisition(ctx context.Context, req domain.MaterialRequisition) error
	GetRequisitionByID(ctx context.Context, id string) (*domain.MaterialRequisition, error)
	GetRequisitionByReqNo(ctx context.Context, reqNo string) (*domain.MaterialRequisition, error)
	UpdateRequisition(ctx context.Context, req domain.MaterialRequisition) error
	ListRequisitions(ctx context.Context, status domain.RequisitionStatus, unitHeadID string) ([]domain.MaterialRequisition, error)

	CreateScrapReturn(ctx context.Context, ret domain.ScrapReturn) error
	GetScrapReturnByID(ctx context.Context, id string) (*domain.ScrapReturn, error)
	ListScrapReturns(ctx context.Context, teamLeadID string) ([]domain.ScrapReturn, error)

	CreateClaim(ctx context.Context, claim domain.StaffAllowanceClaim) error
	GetClaimByID(ctx context.Context, id string) (*domain.StaffAllowanceClaim, error)
	UpdateClaim(ctx context.Context, claim domain.StaffAllowanceClaim) error
	ListClaims(ctx context.Context, staffID string, status domain.ClaimStatus) ([]domain.StaffAllowanceClaim, error)
}

// MemoryRepository provides a concurrent, thread-safe in-memory mock of Repository.
type MemoryRepository struct {
	mu           sync.RWMutex
	items        map[string]domain.StoreItem
	itemsByCode  map[string]string // itemCode -> id
	requisitions map[string]domain.MaterialRequisition
	reqByNo      map[string]string // reqNo -> id
	scrapReturns map[string]domain.ScrapReturn
	claims       map[string]domain.StaffAllowanceClaim
}

// NewMemoryRepository constructs an initialized in-memory store.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items:        make(map[string]domain.StoreItem),
		itemsByCode:  make(map[string]string),
		requisitions: make(map[string]domain.MaterialRequisition),
		reqByNo:      make(map[string]string),
		scrapReturns: make(map[string]domain.ScrapReturn),
		claims:       make(map[string]domain.StaffAllowanceClaim),
	}
}

// CreateStoreItem persists a new SKU master record.
func (m *MemoryRepository) CreateStoreItem(_ context.Context, item domain.StoreItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.items[item.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.itemsByCode[item.ItemCode]; exists {
		return ErrConflict
	}

	m.items[item.ID] = item
	m.itemsByCode[item.ItemCode] = item.ID
	return nil
}

// GetStoreItemByID retrieves a SKU by UUID.
func (m *MemoryRepository) GetStoreItemByID(_ context.Context, id string) (*domain.StoreItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &item, nil
}

// GetStoreItemByCode retrieves a SKU by unique item code.
func (m *MemoryRepository) GetStoreItemByCode(_ context.Context, code string) (*domain.StoreItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.itemsByCode[code]
	if !exists {
		return nil, ErrNotFound
	}
	item := m.items[id]
	return &item, nil
}

// UpdateStoreItem updates inventory balance, costs, or location.
func (m *MemoryRepository) UpdateStoreItem(_ context.Context, item domain.StoreItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.items[item.ID]; !exists {
		return ErrNotFound
	}
	item.UpdatedAt = time.Now().UTC()
	m.items[item.ID] = item
	return nil
}

// ListStoreItems lists SKUs optionally filtered by category or low stock reorder threshold.
func (m *MemoryRepository) ListStoreItems(_ context.Context, category domain.ItemCategory, lowStockOnly bool) ([]domain.StoreItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.StoreItem
	for _, it := range m.items {
		if category != "" && it.Category != category {
			continue
		}
		if lowStockOnly && it.QuantityAvailable > it.ReorderThreshold {
			continue
		}
		results = append(results, it)
	}
	return results, nil
}

// CreateRequisition registers a new material request.
func (m *MemoryRepository) CreateRequisition(_ context.Context, req domain.MaterialRequisition) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.requisitions[req.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.reqByNo[req.ReqNo]; exists {
		return ErrConflict
	}

	m.requisitions[req.ID] = req
	m.reqByNo[req.ReqNo] = req.ID
	return nil
}

// GetRequisitionByID retrieves a requisition by UUID.
func (m *MemoryRepository) GetRequisitionByID(_ context.Context, id string) (*domain.MaterialRequisition, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, exists := m.requisitions[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &r, nil
}

// GetRequisitionByReqNo retrieves a requisition by human-readable req_no.
func (m *MemoryRepository) GetRequisitionByReqNo(_ context.Context, reqNo string) (*domain.MaterialRequisition, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.reqByNo[reqNo]
	if !exists {
		return nil, ErrNotFound
	}
	r := m.requisitions[id]
	return &r, nil
}

// UpdateRequisition updates requisition approval, audit, or issuance state.
func (m *MemoryRepository) UpdateRequisition(_ context.Context, req domain.MaterialRequisition) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.requisitions[req.ID]; !exists {
		return ErrNotFound
	}
	req.UpdatedAt = time.Now().UTC()
	m.requisitions[req.ID] = req
	return nil
}

// ListRequisitions retrieves material requisitions filtered by status and unit head.
func (m *MemoryRepository) ListRequisitions(_ context.Context, status domain.RequisitionStatus, unitHeadID string) ([]domain.MaterialRequisition, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.MaterialRequisition
	for _, r := range m.requisitions {
		if status != "" && r.Status != status {
			continue
		}
		if unitHeadID != "" && r.UnitHeadID != unitHeadID {
			continue
		}
		results = append(results, r)
	}
	return results, nil
}

// CreateScrapReturn logs a damaged or replaced item back to the Store Scrap Yard.
func (m *MemoryRepository) CreateScrapReturn(_ context.Context, ret domain.ScrapReturn) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.scrapReturns[ret.ID]; exists {
		return ErrConflict
	}
	m.scrapReturns[ret.ID] = ret
	return nil
}

// GetScrapReturnByID retrieves a scrap return by UUID.
func (m *MemoryRepository) GetScrapReturnByID(_ context.Context, id string) (*domain.ScrapReturn, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, exists := m.scrapReturns[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &s, nil
}

// ListScrapReturns retrieves scrap returns filtered by team lead.
func (m *MemoryRepository) ListScrapReturns(_ context.Context, teamLeadID string) ([]domain.ScrapReturn, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.ScrapReturn
	for _, s := range m.scrapReturns {
		if teamLeadID != "" && s.TeamLeadID != teamLeadID {
			continue
		}
		results = append(results, s)
	}
	return results, nil
}

// CreateClaim registers a staff allowance claim.
func (m *MemoryRepository) CreateClaim(_ context.Context, claim domain.StaffAllowanceClaim) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.claims[claim.ID]; exists {
		return ErrConflict
	}
	m.claims[claim.ID] = claim
	return nil
}

// GetClaimByID retrieves a financial claim by UUID.
func (m *MemoryRepository) GetClaimByID(_ context.Context, id string) (*domain.StaffAllowanceClaim, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, exists := m.claims[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &c, nil
}

// UpdateClaim updates a claim's audit status and decision notes.
func (m *MemoryRepository) UpdateClaim(_ context.Context, claim domain.StaffAllowanceClaim) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.claims[claim.ID]; !exists {
		return ErrNotFound
	}
	claim.UpdatedAt = time.Now().UTC()
	m.claims[claim.ID] = claim
	return nil
}

// ListClaims returns allowance claims filtered by staff and status.
func (m *MemoryRepository) ListClaims(_ context.Context, staffID string, status domain.ClaimStatus) ([]domain.StaffAllowanceClaim, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.StaffAllowanceClaim
	for _, c := range m.claims {
		if staffID != "" && c.StaffID != staffID {
			continue
		}
		if status != "" && c.Status != status {
			continue
		}
		results = append(results, c)
	}
	return results, nil
}

