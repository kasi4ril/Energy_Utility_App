// Package repository implements memory persistence for the executive action console.
package repository

import (
	"context"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

// CreateActionItem stores a pending C-Suite action item.
func (m *MemoryRepository) CreateActionItem(_ context.Context, item domain.ExecutiveActionItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.actions[item.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.actionsByRef[item.SourceEntityRef]; exists {
		return ErrConflict
	}
	m.actions[item.ID] = item
	m.actionsByRef[item.SourceEntityRef] = item.ID
	return nil
}

// GetActionItemByID retrieves an action item by UUID.
func (m *MemoryRepository) GetActionItemByID(_ context.Context, id string) (*domain.ExecutiveActionItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.actions[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &item, nil
}

// GetActionItemByRef retrieves an action item by source entity reference.
func (m *MemoryRepository) GetActionItemByRef(_ context.Context, ref string) (*domain.ExecutiveActionItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.actionsByRef[ref]
	if !exists {
		return nil, ErrNotFound
	}
	item := m.actions[id]
	return &item, nil
}

// UpdateActionItem updates decision status and audit notes.
func (m *MemoryRepository) UpdateActionItem(_ context.Context, item domain.ExecutiveActionItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.actions[item.ID]; !exists {
		return ErrNotFound
	}
	m.actions[item.ID] = item
	return nil
}

// ListActionItems filters action items by target role and status.
func (m *MemoryRepository) ListActionItems(_ context.Context, targetRole domain.ExecutiveRole, status domain.ActionItemStatus) ([]domain.ExecutiveActionItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]domain.ExecutiveActionItem, 0)
	for _, item := range m.actions {
		if targetRole != "" && item.TargetRole != targetRole {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		res = append(res, item)
	}
	return res, nil
}

