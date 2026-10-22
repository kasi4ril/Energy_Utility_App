// Package repository defines the persistence contracts and thread-safe in-memory mock
// for S09 Executive Command, C-Suite Decision Console & Predictive Maintenance.
package repository

import (
	"context"
	"errors"
	"sync"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

var (
	// ErrNotFound indicates the requested entity was not found in storage.
	ErrNotFound = errors.New("repository: entity not found")
	// ErrConflict indicates an entity with the specified identifier or reference already exists.
	ErrConflict = errors.New("repository: entity conflict detected")
)

// Repository specifies the persistence operations for executive snapshots, predictive triggers, and action items.
type Repository interface {
	// Executive Snapshots
	CreateExecutiveSnapshot(ctx context.Context, s domain.ExecutiveSnapshot) error
	GetExecutiveSnapshotByID(ctx context.Context, id string) (*domain.ExecutiveSnapshot, error)
	GetLatestExecutiveSnapshot(ctx context.Context) (*domain.ExecutiveSnapshot, error)
	ListExecutiveSnapshots(ctx context.Context, limit int) ([]domain.ExecutiveSnapshot, error)

	// Predictive Maintenance Triggers
	CreatePredictiveTrigger(ctx context.Context, t domain.PredictiveMaintenanceTrigger) error
	GetPredictiveTriggerByID(ctx context.Context, id string) (*domain.PredictiveMaintenanceTrigger, error)
	GetPredictiveTriggerByNo(ctx context.Context, triggerNo string) (*domain.PredictiveMaintenanceTrigger, error)
	UpdatePredictiveTrigger(ctx context.Context, t domain.PredictiveMaintenanceTrigger) error
	ListPredictiveTriggers(ctx context.Context, status string, model domain.PredictiveModelName, limit int) ([]domain.PredictiveMaintenanceTrigger, error)

	// Executive Action Items Console
	CreateActionItem(ctx context.Context, item domain.ExecutiveActionItem) error
	GetActionItemByID(ctx context.Context, id string) (*domain.ExecutiveActionItem, error)
	GetActionItemByRef(ctx context.Context, ref string) (*domain.ExecutiveActionItem, error)
	UpdateActionItem(ctx context.Context, item domain.ExecutiveActionItem) error
	ListActionItems(ctx context.Context, targetRole domain.ExecutiveRole, status domain.ActionItemStatus) ([]domain.ExecutiveActionItem, error)
}

// MemoryRepository provides a concurrent, thread-safe in-memory mock of Repository.
type MemoryRepository struct {
	mu           sync.RWMutex
	snapshots    map[string]domain.ExecutiveSnapshot
	latestSnapID string

	triggers     map[string]domain.PredictiveMaintenanceTrigger
	triggersByNo map[string]string // triggerNo -> id

	actions      map[string]domain.ExecutiveActionItem
	actionsByRef map[string]string // sourceEntityRef -> id
}

// NewMemoryRepository initializes empty collections for in-memory persistence.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		snapshots:    make(map[string]domain.ExecutiveSnapshot),
		triggers:     make(map[string]domain.PredictiveMaintenanceTrigger),
		triggersByNo: make(map[string]string),
		actions:      make(map[string]domain.ExecutiveActionItem),
		actionsByRef: make(map[string]string),
	}
}

// CreateExecutiveSnapshot stores an aggregated telemetry snapshot.
func (m *MemoryRepository) CreateExecutiveSnapshot(_ context.Context, s domain.ExecutiveSnapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.snapshots[s.ID]; exists {
		return ErrConflict
	}
	m.snapshots[s.ID] = s
	m.latestSnapID = s.ID
	return nil
}

// GetExecutiveSnapshotByID retrieves a snapshot by UUID.
func (m *MemoryRepository) GetExecutiveSnapshotByID(_ context.Context, id string) (*domain.ExecutiveSnapshot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, exists := m.snapshots[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &s, nil
}

// GetLatestExecutiveSnapshot retrieves the most recently captured operational snapshot.
func (m *MemoryRepository) GetLatestExecutiveSnapshot(_ context.Context) (*domain.ExecutiveSnapshot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.latestSnapID == "" {
		return nil, ErrNotFound
	}
	s, exists := m.snapshots[m.latestSnapID]
	if !exists {
		return nil, ErrNotFound
	}
	return &s, nil
}

// ListExecutiveSnapshots lists historical snapshots up to the limit.
func (m *MemoryRepository) ListExecutiveSnapshots(_ context.Context, limit int) ([]domain.ExecutiveSnapshot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]domain.ExecutiveSnapshot, 0, len(m.snapshots))
	for _, s := range m.snapshots {
		res = append(res, s)
		if limit > 0 && len(res) >= limit {
			break
		}
	}
	return res, nil
}

// CreatePredictiveTrigger stores an automated predictive maintenance alarm.
func (m *MemoryRepository) CreatePredictiveTrigger(_ context.Context, t domain.PredictiveMaintenanceTrigger) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.triggers[t.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.triggersByNo[t.TriggerNo]; exists {
		return ErrConflict
	}
	m.triggers[t.ID] = t
	m.triggersByNo[t.TriggerNo] = t.ID
	return nil
}

// GetPredictiveTriggerByID retrieves a trigger by UUID.
func (m *MemoryRepository) GetPredictiveTriggerByID(_ context.Context, id string) (*domain.PredictiveMaintenanceTrigger, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, exists := m.triggers[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &t, nil
}

// GetPredictiveTriggerByNo retrieves a trigger by trigger number.
func (m *MemoryRepository) GetPredictiveTriggerByNo(_ context.Context, triggerNo string) (*domain.PredictiveMaintenanceTrigger, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.triggersByNo[triggerNo]
	if !exists {
		return nil, ErrNotFound
	}
	t := m.triggers[id]
	return &t, nil
}

// UpdatePredictiveTrigger updates trigger status and acknowledgment.
func (m *MemoryRepository) UpdatePredictiveTrigger(_ context.Context, t domain.PredictiveMaintenanceTrigger) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.triggers[t.ID]; !exists {
		return ErrNotFound
	}
	m.triggers[t.ID] = t
	return nil
}

// ListPredictiveTriggers filters predictive maintenance alarms.
func (m *MemoryRepository) ListPredictiveTriggers(_ context.Context, status string, model domain.PredictiveModelName, limit int) ([]domain.PredictiveMaintenanceTrigger, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]domain.PredictiveMaintenanceTrigger, 0)
	for _, t := range m.triggers {
		if status != "" && t.Status != status {
			continue
		}
		if model != "" && t.ModelName != model {
			continue
		}
		res = append(res, t)
		if limit > 0 && len(res) >= limit {
			break
		}
	}
	return res, nil
}
