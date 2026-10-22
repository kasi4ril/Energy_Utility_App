// Package repository defines the persistence contract and thread-safe in-memory mock
// for the S06 Energy Management bounded context.
package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
)

var (
	// ErrNotFound indicates the requested entity was not found in storage.
	ErrNotFound = errors.New("repository: entity not found")
	// ErrConflict indicates an entity with the specified identifier or permit number already exists.
	ErrConflict = errors.New("repository: entity conflict detected")
)

// Repository specifies the persistence contract for dispatch telemetry, feeder states, and safety PTWs.
type Repository interface {
	CreateDispatchLog(ctx context.Context, log domain.DispatchTelemetryLog) error
	GetLatestDispatchLog(ctx context.Context) (*domain.DispatchTelemetryLog, error)
	ListDispatchLogs(ctx context.Context, shiftID string, limit int) ([]domain.DispatchTelemetryLog, error)

	UpsertFeederStatus(ctx context.Context, status domain.FeederStatus) error
	GetFeederStatusByAssetID(ctx context.Context, assetID string) (*domain.FeederStatus, error)
	ListFeederStatuses(ctx context.Context) ([]domain.FeederStatus, error)

	CreatePTW(ctx context.Context, ptw domain.SafetyPermitPTW) error
	GetPTWByID(ctx context.Context, id string) (*domain.SafetyPermitPTW, error)
	GetPTWByPermitNo(ctx context.Context, permitNo string) (*domain.SafetyPermitPTW, error)
	UpdatePTWStatus(ctx context.Context, ptw domain.SafetyPermitPTW) error
	ListPTWs(ctx context.Context, status domain.PTWStatus, assetID string) ([]domain.SafetyPermitPTW, error)
}

// MemoryRepository provides a concurrent, thread-safe in-memory mock of Repository.
type MemoryRepository struct {
	mu           sync.RWMutex
	dispatchLogs map[string]domain.DispatchTelemetryLog
	feeders      map[string]domain.FeederStatus // assetID -> FeederStatus
	ptwPermits   map[string]domain.SafetyPermitPTW
	ptwByNo      map[string]string // permitNo -> ID
}

// NewMemoryRepository constructs a new in-memory store.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		dispatchLogs: make(map[string]domain.DispatchTelemetryLog),
		feeders:      make(map[string]domain.FeederStatus),
		ptwPermits:   make(map[string]domain.SafetyPermitPTW),
		ptwByNo:      make(map[string]string),
	}
}

// CreateDispatchLog persists a bulk energy intake and roll-call record.
func (m *MemoryRepository) CreateDispatchLog(_ context.Context, log domain.DispatchTelemetryLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.dispatchLogs[log.ID]; exists {
		return ErrConflict
	}
	m.dispatchLogs[log.ID] = log
	return nil
}

// GetLatestDispatchLog returns the most recently recorded dispatch telemetry.
func (m *MemoryRepository) GetLatestDispatchLog(_ context.Context) (*domain.DispatchTelemetryLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var latest *domain.DispatchTelemetryLog
	for _, l := range m.dispatchLogs {
		if latest == nil || l.LoggedAt.After(latest.LoggedAt) {
			lCopy := l
			latest = &lCopy
		}
	}
	if latest == nil {
		return nil, ErrNotFound
	}
	return latest, nil
}

// ListDispatchLogs returns dispatch logs filtered by shift.
func (m *MemoryRepository) ListDispatchLogs(_ context.Context, shiftID string, limit int) ([]domain.DispatchTelemetryLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 {
		limit = 50
	}
	var results []domain.DispatchTelemetryLog
	for _, l := range m.dispatchLogs {
		if shiftID != "" && l.ShiftID != shiftID {
			continue
		}
		results = append(results, l)
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

// UpsertFeederStatus updates or inserts a feeder's real-time state and load.
func (m *MemoryRepository) UpsertFeederStatus(_ context.Context, status domain.FeederStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	status.UpdatedAt = time.Now().UTC()
	m.feeders[status.FeederAssetID] = status
	return nil
}

// GetFeederStatusByAssetID retrieves real-time telemetry for a feeder asset.
func (m *MemoryRepository) GetFeederStatusByAssetID(_ context.Context, assetID string) (*domain.FeederStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	f, exists := m.feeders[assetID]
	if !exists {
		return nil, ErrNotFound
	}
	return &f, nil
}

// ListFeederStatuses returns all registered feeder states.
func (m *MemoryRepository) ListFeederStatuses(_ context.Context) ([]domain.FeederStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.FeederStatus
	for _, f := range m.feeders {
		results = append(results, f)
	}
	return results, nil
}

// CreatePTW persists a new safety permit.
func (m *MemoryRepository) CreatePTW(_ context.Context, ptw domain.SafetyPermitPTW) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.ptwPermits[ptw.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.ptwByNo[ptw.PermitNo]; exists {
		return ErrConflict
	}

	m.ptwPermits[ptw.ID] = ptw
	m.ptwByNo[ptw.PermitNo] = ptw.ID
	return nil
}

// GetPTWByID retrieves a permit by UUID.
func (m *MemoryRepository) GetPTWByID(_ context.Context, id string) (*domain.SafetyPermitPTW, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.ptwPermits[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &p, nil
}

// GetPTWByPermitNo retrieves a permit by human-readable permit number.
func (m *MemoryRepository) GetPTWByPermitNo(_ context.Context, permitNo string) (*domain.SafetyPermitPTW, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.ptwByNo[permitNo]
	if !exists {
		return nil, ErrNotFound
	}
	p := m.ptwPermits[id]
	return &p, nil
}

// UpdatePTWStatus updates a safety permit record.
func (m *MemoryRepository) UpdatePTWStatus(_ context.Context, ptw domain.SafetyPermitPTW) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.ptwPermits[ptw.ID]; !exists {
		return ErrNotFound
	}
	ptw.UpdatedAt = time.Now().UTC()
	m.ptwPermits[ptw.ID] = ptw
	return nil
}

// ListPTWs returns permits matching status and asset filters.
func (m *MemoryRepository) ListPTWs(_ context.Context, status domain.PTWStatus, assetID string) ([]domain.SafetyPermitPTW, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.SafetyPermitPTW
	for _, p := range m.ptwPermits {
		if status != "" && p.Status != status {
			continue
		}
		if assetID != "" && p.AssetID != assetID {
			continue
		}
		results = append(results, p)
	}
	return results, nil
}

