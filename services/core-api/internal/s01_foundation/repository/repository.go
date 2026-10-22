// Package repository provides persistence abstractions and storage implementations
// for the Foundation bounded context (S01).
package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

var (
	// ErrNotFound is returned when an entity does not exist in storage.
	ErrNotFound = errors.New("repository: entity not found")
	// ErrDuplicateKey is returned on unique constraint collision.
	ErrDuplicateKey = errors.New("repository: duplicate unique identifier")
)

// Repository specifies the foundational data access contract across schemas.
type Repository interface {
	// OrgUnit operations
	CreateOrgUnit(ctx context.Context, unit domain.OrgUnit) error
	GetOrgUnitByID(ctx context.Context, id string) (*domain.OrgUnit, error)
	ListOrgUnits(ctx context.Context) ([]domain.OrgUnit, error)

	// StaffProfile operations
	CreateStaffProfile(ctx context.Context, staff domain.StaffProfile) error
	GetStaffProfileByID(ctx context.Context, id string) (*domain.StaffProfile, error)
	GetStaffProfileByEmail(ctx context.Context, email string) (*domain.StaffProfile, error)
	GetStaffProfileByStaffNumber(ctx context.Context, staffNumber string) (*domain.StaffProfile, error)
	UpdateStaffTransfer(ctx context.Context, staffID, newOrgUnitID string, newRole domain.UserRole, newHMAC string) error
	ListStaffProfiles(ctx context.Context, orgUnitID string) ([]domain.StaffProfile, error)

	// GridAsset operations
	CreateGridAsset(ctx context.Context, asset domain.GridAsset) error
	GetGridAssetByID(ctx context.Context, id string) (*domain.GridAsset, error)
	GetGridAssetByCode(ctx context.Context, code string) (*domain.GridAsset, error)
	ListGridAssets(ctx context.Context, parentID string, assetType domain.AssetType) ([]domain.GridAsset, error)

	// WorkNotes operations
	CreateWorkNote(ctx context.Context, note domain.WorkNote) error
	GetLatestWorkNote(ctx context.Context, entityType, entityID string) (*domain.WorkNote, error)
	ListWorkNotesByEntity(ctx context.Context, entityType, entityID string) ([]domain.WorkNote, error)
}

// MemoryRepository provides a concurrent in-memory repository for unit testing and edge operation.
type MemoryRepository struct {
	mu         sync.RWMutex
	orgUnits   map[string]domain.OrgUnit
	staff      map[string]domain.StaffProfile
	assets     map[string]domain.GridAsset
	notes      []domain.WorkNote
}

// NewMemoryRepository constructs a new thread-safe in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		orgUnits: make(map[string]domain.OrgUnit),
		staff:    make(map[string]domain.StaffProfile),
		assets:   make(map[string]domain.GridAsset),
		notes:    make([]domain.WorkNote, 0),
	}
}

// CreateOrgUnit saves a new org unit.
func (m *MemoryRepository) CreateOrgUnit(_ context.Context, unit domain.OrgUnit) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.orgUnits {
		if existing.UnitCode == unit.UnitCode {
			return ErrDuplicateKey
		}
	}
	m.orgUnits[unit.ID] = unit
	return nil
}

// GetOrgUnitByID retrieves an org unit by primary key.
func (m *MemoryRepository) GetOrgUnitByID(_ context.Context, id string) (*domain.OrgUnit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	unit, ok := m.orgUnits[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &unit, nil
}

// ListOrgUnits retrieves all registered org units.
func (m *MemoryRepository) ListOrgUnits(_ context.Context) ([]domain.OrgUnit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	units := make([]domain.OrgUnit, 0, len(m.orgUnits))
	for _, u := range m.orgUnits {
		units = append(units, u)
	}
	return units, nil
}

// CreateStaffProfile persists a new staff profile.
func (m *MemoryRepository) CreateStaffProfile(_ context.Context, staff domain.StaffProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.staff {
		if existing.StaffNumber == staff.StaffNumber || existing.Email == staff.Email {
			return ErrDuplicateKey
		}
	}
	m.staff[staff.ID] = staff
	return nil
}

// GetStaffProfileByID fetches staff by ID.
func (m *MemoryRepository) GetStaffProfileByID(_ context.Context, id string) (*domain.StaffProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.staff[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &s, nil
}

// GetStaffProfileByEmail fetches staff by email address.
func (m *MemoryRepository) GetStaffProfileByEmail(_ context.Context, email string) (*domain.StaffProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.staff {
		if s.Email == email {
			return &s, nil
		}
	}
	return nil, ErrNotFound
}

// GetStaffProfileByStaffNumber fetches staff by company staff code.
func (m *MemoryRepository) GetStaffProfileByStaffNumber(_ context.Context, staffNumber string) (*domain.StaffProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.staff {
		if s.StaffNumber == staffNumber {
			return &s, nil
		}
	}
	return nil, ErrNotFound
}

// UpdateStaffTransfer executes administrative re-assignment of staff.
func (m *MemoryRepository) UpdateStaffTransfer(_ context.Context, staffID, newOrgUnitID string, newRole domain.UserRole, newHMAC string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.staff[staffID]
	if !ok {
		return ErrNotFound
	}
	s.OrgUnitID = newOrgUnitID
	s.Role = newRole
	s.RowHMACSignature = newHMAC
	s.UpdatedAt = time.Now()
	m.staff[staffID] = s
	return nil
}

// ListStaffProfiles lists staff deployed in a given unit or all staff if unit is empty.
func (m *MemoryRepository) ListStaffProfiles(_ context.Context, orgUnitID string) ([]domain.StaffProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.StaffProfile, 0)
	for _, s := range m.staff {
		if orgUnitID == "" || s.OrgUnitID == orgUnitID {
			res = append(res, s)
		}
	}
	return res, nil
}

// CreateGridAsset registers an electrical asset.
func (m *MemoryRepository) CreateGridAsset(_ context.Context, asset domain.GridAsset) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.assets {
		if a.AssetCode == asset.AssetCode {
			return ErrDuplicateKey
		}
	}
	m.assets[asset.ID] = asset
	return nil
}

// GetGridAssetByID fetches an asset by UUID.
func (m *MemoryRepository) GetGridAssetByID(_ context.Context, id string) (*domain.GridAsset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.assets[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &a, nil
}

// GetGridAssetByCode fetches an asset by unique grid asset code.
func (m *MemoryRepository) GetGridAssetByCode(_ context.Context, code string) (*domain.GridAsset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, a := range m.assets {
		if a.AssetCode == code {
			return &a, nil
		}
	}
	return nil, ErrNotFound
}

// ListGridAssets lists assets filtered by parent asset or asset type.
func (m *MemoryRepository) ListGridAssets(_ context.Context, parentID string, assetType domain.AssetType) ([]domain.GridAsset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.GridAsset, 0)
	for _, a := range m.assets {
		if parentID != "" && (a.ParentAssetID == nil || *a.ParentAssetID != parentID) {
			continue
		}
		if assetType != "" && a.AssetType != assetType {
			continue
		}
		res = append(res, a)
	}
	return res, nil
}

// CreateWorkNote appends a note to the immutable ledger.
func (m *MemoryRepository) CreateWorkNote(_ context.Context, note domain.WorkNote) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notes = append(m.notes, note)
	return nil
}

// GetLatestWorkNote returns the most recent note in the hash chain for the given entity.
func (m *MemoryRepository) GetLatestWorkNote(_ context.Context, entityType, entityID string) (*domain.WorkNote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for i := len(m.notes) - 1; i >= 0; i-- {
		n := m.notes[i]
		if n.EntityType == entityType && n.EntityID == entityID {
			return &n, nil
		}
	}
	return nil, nil // No notes yet
}

// ListWorkNotesByEntity returns the chronologically ordered note stream for an entity.
func (m *MemoryRepository) ListWorkNotesByEntity(_ context.Context, entityType, entityID string) ([]domain.WorkNote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.WorkNote, 0)
	for _, n := range m.notes {
		if n.EntityType == entityType && n.EntityID == entityID {
			res = append(res, n)
		}
	}
	return res, nil
}

