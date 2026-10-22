// Package repository defines the persistence contracts and thread-safe in-memory mock
// for S08 Billing, Commercial Operations & Revenue Protection AI.
package repository

import (
	"context"
	"errors"
	"sync"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
)

var (
	// ErrNotFound indicates the requested entity was not found in storage.
	ErrNotFound = errors.New("repository: entity not found")
	// ErrConflict indicates an entity with the specified identifier or reference already exists.
	ErrConflict = errors.New("repository: entity conflict detected")
)

// Repository specifies the persistence operations for vending, refunds, adjustments, and AI alerts.
type Repository interface {
	// Vending Transactions
	CreateVendingTransaction(ctx context.Context, tx domain.VendingTransaction) error
	GetVendingTransactionByID(ctx context.Context, id string) (*domain.VendingTransaction, error)
	GetVendingTransactionByRef(ctx context.Context, txRef string) (*domain.VendingTransaction, error)
	ListVendingTransactions(ctx context.Context, customerID, meterNumber string, limit int) ([]domain.VendingTransaction, error)

	// Meter Refund Ledgers
	CreateRefundLedger(ctx context.Context, ledger domain.MeterRefundLedger) error
	GetRefundLedgerByID(ctx context.Context, id string) (*domain.MeterRefundLedger, error)
	GetRefundLedgerByCustomer(ctx context.Context, customerID string) (*domain.MeterRefundLedger, error)
	GetRefundLedgerByMeter(ctx context.Context, meterSerial string) (*domain.MeterRefundLedger, error)
	UpdateRefundLedger(ctx context.Context, ledger domain.MeterRefundLedger) error

	// Billing Adjustments
	CreateBillingAdjustment(ctx context.Context, adj domain.BillingAdjustment) error
	GetBillingAdjustmentByID(ctx context.Context, id string) (*domain.BillingAdjustment, error)
	GetBillingAdjustmentByNo(ctx context.Context, adjNo string) (*domain.BillingAdjustment, error)
	UpdateBillingAdjustment(ctx context.Context, adj domain.BillingAdjustment) error
	ListBillingAdjustments(ctx context.Context, status domain.AdjustmentStatus, customerID string) ([]domain.BillingAdjustment, error)

	// Revenue Protection AI Alerts
	CreateAIAlert(ctx context.Context, alert domain.RevenueProtectionAlert) error
	GetAIAlertByID(ctx context.Context, id string) (*domain.RevenueProtectionAlert, error)
	GetAIAlertByNo(ctx context.Context, alertNo string) (*domain.RevenueProtectionAlert, error)
	UpdateAIAlert(ctx context.Context, alert domain.RevenueProtectionAlert) error
	ListAIAlerts(ctx context.Context, status domain.AIAlertStatus, alertType domain.AIAlertType) ([]domain.RevenueProtectionAlert, error)
}

// MemoryRepository provides a concurrent, thread-safe in-memory mock of Repository.
type MemoryRepository struct {
	mu              sync.RWMutex
	vending         map[string]domain.VendingTransaction
	vendingByRef    map[string]string // txReference -> id
	refunds         map[string]domain.MeterRefundLedger
	refundByMeter   map[string]string // meterSerial -> id
	refundByCust    map[string]string // customerID -> id
	adjustments     map[string]domain.BillingAdjustment
	adjustmentsByNo map[string]string // adjustmentNo -> id
	aiAlerts        map[string]domain.RevenueProtectionAlert
	aiAlertsByNo    map[string]string // alertNo -> id
}

// NewMemoryRepository initializes empty collections for in-memory persistence.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		vending:         make(map[string]domain.VendingTransaction),
		vendingByRef:    make(map[string]string),
		refunds:         make(map[string]domain.MeterRefundLedger),
		refundByMeter:   make(map[string]string),
		refundByCust:    make(map[string]string),
		adjustments:     make(map[string]domain.BillingAdjustment),
		adjustmentsByNo: make(map[string]string),
		aiAlerts:        make(map[string]domain.RevenueProtectionAlert),
		aiAlertsByNo:    make(map[string]string),
	}
}

// CreateVendingTransaction stores a vending record.
func (m *MemoryRepository) CreateVendingTransaction(_ context.Context, tx domain.VendingTransaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.vending[tx.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.vendingByRef[tx.TxReference]; exists {
		return ErrConflict
	}
	m.vending[tx.ID] = tx
	m.vendingByRef[tx.TxReference] = tx.ID
	return nil
}

// GetVendingTransactionByID retrieves transaction by primary UUID.
func (m *MemoryRepository) GetVendingTransactionByID(_ context.Context, id string) (*domain.VendingTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tx, exists := m.vending[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &tx, nil
}

// GetVendingTransactionByRef retrieves transaction by unique reference.
func (m *MemoryRepository) GetVendingTransactionByRef(_ context.Context, txRef string) (*domain.VendingTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.vendingByRef[txRef]
	if !exists {
		return nil, ErrNotFound
	}
	tx := m.vending[id]
	return &tx, nil
}

// ListVendingTransactions lists matching transactions up to limit.
func (m *MemoryRepository) ListVendingTransactions(_ context.Context, customerID, meterNumber string, limit int) ([]domain.VendingTransaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]domain.VendingTransaction, 0)
	for _, tx := range m.vending {
		if customerID != "" && tx.CustomerID != customerID {
			continue
		}
		if meterNumber != "" && tx.MeterNumber != meterNumber {
			continue
		}
		res = append(res, tx)
		if limit > 0 && len(res) >= limit {
			break
		}
	}
	return res, nil
}

// CreateRefundLedger inserts a new customer refund ledger.
func (m *MemoryRepository) CreateRefundLedger(_ context.Context, ledger domain.MeterRefundLedger) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.refunds[ledger.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.refundByMeter[ledger.MeterSerial]; exists {
		return ErrConflict
	}
	m.refunds[ledger.ID] = ledger
	m.refundByMeter[ledger.MeterSerial] = ledger.ID
	m.refundByCust[ledger.CustomerID] = ledger.ID
	return nil
}

// GetRefundLedgerByID retrieves ledger by ID.
func (m *MemoryRepository) GetRefundLedgerByID(_ context.Context, id string) (*domain.MeterRefundLedger, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ledger, exists := m.refunds[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &ledger, nil
}

// GetRefundLedgerByCustomer retrieves ledger by customer ID.
func (m *MemoryRepository) GetRefundLedgerByCustomer(_ context.Context, customerID string) (*domain.MeterRefundLedger, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.refundByCust[customerID]
	if !exists {
		return nil, ErrNotFound
	}
	ledger := m.refunds[id]
	return &ledger, nil
}

// GetRefundLedgerByMeter retrieves ledger by meter serial.
func (m *MemoryRepository) GetRefundLedgerByMeter(_ context.Context, meterSerial string) (*domain.MeterRefundLedger, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.refundByMeter[meterSerial]
	if !exists {
		return nil, ErrNotFound
	}
	ledger := m.refunds[id]
	return &ledger, nil
}

// UpdateRefundLedger updates ledger state.
func (m *MemoryRepository) UpdateRefundLedger(_ context.Context, ledger domain.MeterRefundLedger) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.refunds[ledger.ID]; !exists {
		return ErrNotFound
	}
	m.refunds[ledger.ID] = ledger
	return nil
}

