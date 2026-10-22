// Package repository implements memory persistence for billing adjustments and AI alerts.
package repository

import (
	"context"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
)

// CreateBillingAdjustment stores a proposed adjustment.
func (m *MemoryRepository) CreateBillingAdjustment(_ context.Context, adj domain.BillingAdjustment) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.adjustments[adj.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.adjustmentsByNo[adj.AdjustmentNo]; exists {
		return ErrConflict
	}
	m.adjustments[adj.ID] = adj
	m.adjustmentsByNo[adj.AdjustmentNo] = adj.ID
	return nil
}

// GetBillingAdjustmentByID retrieves adjustment by ID.
func (m *MemoryRepository) GetBillingAdjustmentByID(_ context.Context, id string) (*domain.BillingAdjustment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	adj, exists := m.adjustments[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &adj, nil
}

// GetBillingAdjustmentByNo retrieves adjustment by unique adjustment number.
func (m *MemoryRepository) GetBillingAdjustmentByNo(_ context.Context, adjNo string) (*domain.BillingAdjustment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.adjustmentsByNo[adjNo]
	if !exists {
		return nil, ErrNotFound
	}
	adj := m.adjustments[id]
	return &adj, nil
}

// UpdateBillingAdjustment updates adjustment state.
func (m *MemoryRepository) UpdateBillingAdjustment(_ context.Context, adj domain.BillingAdjustment) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.adjustments[adj.ID]; !exists {
		return ErrNotFound
	}
	m.adjustments[adj.ID] = adj
	return nil
}

// ListBillingAdjustments filters adjustments by status and customer.
func (m *MemoryRepository) ListBillingAdjustments(_ context.Context, status domain.AdjustmentStatus, customerID string) ([]domain.BillingAdjustment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]domain.BillingAdjustment, 0)
	for _, adj := range m.adjustments {
		if status != "" && adj.Status != status {
			continue
		}
		if customerID != "" && adj.CustomerID != customerID {
			continue
		}
		res = append(res, adj)
	}
	return res, nil
}

// CreateAIAlert records a new revenue protection alert.
func (m *MemoryRepository) CreateAIAlert(_ context.Context, alert domain.RevenueProtectionAlert) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.aiAlerts[alert.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.aiAlertsByNo[alert.AlertNo]; exists {
		return ErrConflict
	}
	m.aiAlerts[alert.ID] = alert
	m.aiAlertsByNo[alert.AlertNo] = alert.ID
	return nil
}

// GetAIAlertByID retrieves alert by primary ID.
func (m *MemoryRepository) GetAIAlertByID(_ context.Context, id string) (*domain.RevenueProtectionAlert, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	alert, exists := m.aiAlerts[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &alert, nil
}

// GetAIAlertByNo retrieves alert by alert number.
func (m *MemoryRepository) GetAIAlertByNo(_ context.Context, alertNo string) (*domain.RevenueProtectionAlert, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.aiAlertsByNo[alertNo]
	if !exists {
		return nil, ErrNotFound
	}
	alert := m.aiAlerts[id]
	return &alert, nil
}

// UpdateAIAlert updates alert fields.
func (m *MemoryRepository) UpdateAIAlert(_ context.Context, alert domain.RevenueProtectionAlert) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.aiAlerts[alert.ID]; !exists {
		return ErrNotFound
	}
	m.aiAlerts[alert.ID] = alert
	return nil
}

// ListAIAlerts filters alerts by status and alert type.
func (m *MemoryRepository) ListAIAlerts(_ context.Context, status domain.AIAlertStatus, alertType domain.AIAlertType) ([]domain.RevenueProtectionAlert, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]domain.RevenueProtectionAlert, 0)
	for _, alert := range m.aiAlerts {
		if status != "" && alert.Status != status {
			continue
		}
		if alertType != "" && alert.AlertType != alertType {
			continue
		}
		res = append(res, alert)
	}
	return res, nil
}

