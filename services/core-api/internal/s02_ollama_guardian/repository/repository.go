// Package repository defines storage interfaces and memory repository for S02 Ollama Guardian.
package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s02_ollama_guardian/domain"
)

var (
	// ErrNotFound is returned when an audit record or request is missing.
	ErrNotFound = errors.New("repository: entity not found")
)

// Repository specifies data operations for the security_ai schema.
type Repository interface {
	CreateAuditRecord(ctx context.Context, record domain.AuditRecord) error
	GetLatestAuditRecord(ctx context.Context, tableName, recordID string) (*domain.AuditRecord, error)
	ListAuditRecords(ctx context.Context, tableName, recordID string) ([]domain.AuditRecord, error)

	CreateMakerCheckerRequest(ctx context.Context, req domain.MakerCheckerRequest) error
	GetMakerCheckerRequestByID(ctx context.Context, id string) (*domain.MakerCheckerRequest, error)
	UpdateMakerCheckerStatus(ctx context.Context, id string, status domain.RequestStatus, checkerID string, note string) error
	ListMakerCheckerRequests(ctx context.Context, status domain.RequestStatus, module string) ([]domain.MakerCheckerRequest, error)

	CreateChatMessage(ctx context.Context, msg domain.ChatMessage) error
	ListChatMessages(ctx context.Context, channelID string, limit int) ([]domain.ChatMessage, error)
}

// MemoryRepository provides a concurrent, in-memory repository for S02 testing and edge environments.
type MemoryRepository struct {
	mu           sync.RWMutex
	auditRecords []domain.AuditRecord
	requests     map[string]domain.MakerCheckerRequest
	chatMessages []domain.ChatMessage
}

// NewMemoryRepository constructs a new in-memory security repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		auditRecords: make([]domain.AuditRecord, 0),
		requests:     make(map[string]domain.MakerCheckerRequest),
		chatMessages: make([]domain.ChatMessage, 0),
	}
}

// CreateAuditRecord appends an audit record to the chain.
func (m *MemoryRepository) CreateAuditRecord(_ context.Context, rec domain.AuditRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auditRecords = append(m.auditRecords, rec)
	return nil
}

// GetLatestAuditRecord retrieves the most recent audit entry for a given table and record ID.
func (m *MemoryRepository) GetLatestAuditRecord(_ context.Context, tableName, recordID string) (*domain.AuditRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for i := len(m.auditRecords) - 1; i >= 0; i-- {
		r := m.auditRecords[i]
		if r.TableName == tableName && r.RecordID == recordID {
			return &r, nil
		}
	}
	return nil, nil // No previous audit entry
}

// ListAuditRecords returns all audit chain records for an entity in chronological sequence.
func (m *MemoryRepository) ListAuditRecords(_ context.Context, tableName, recordID string) ([]domain.AuditRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.AuditRecord, 0)
	for _, r := range m.auditRecords {
		if r.TableName == tableName && r.RecordID == recordID {
			res = append(res, r)
		}
	}
	return res, nil
}

// CreateMakerCheckerRequest registers a pending maker-checker change request.
func (m *MemoryRepository) CreateMakerCheckerRequest(_ context.Context, req domain.MakerCheckerRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests[req.ID] = req
	return nil
}

// GetMakerCheckerRequestByID fetches a request by primary key UUID.
func (m *MemoryRepository) GetMakerCheckerRequestByID(_ context.Context, id string) (*domain.MakerCheckerRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	req, ok := m.requests[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &req, nil
}

// UpdateMakerCheckerStatus updates approval/rejection state and checker remarks.
func (m *MemoryRepository) UpdateMakerCheckerStatus(_ context.Context, id string, status domain.RequestStatus, checkerID string, note string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	req, ok := m.requests[id]
	if !ok {
		return ErrNotFound
	}
	req.Status = status
	req.CheckerID = &checkerID
	req.CheckerNote = note
	req.UpdatedAt = time.Now()
	m.requests[id] = req
	return nil
}

// ListMakerCheckerRequests returns requests filtered by status and/or module.
func (m *MemoryRepository) ListMakerCheckerRequests(_ context.Context, status domain.RequestStatus, module string) ([]domain.MakerCheckerRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.MakerCheckerRequest, 0)
	for _, r := range m.requests {
		if status != "" && r.Status != status {
			continue
		}
		if module != "" && r.Module != module {
			continue
		}
		res = append(res, r)
	}
	return res, nil
}

// CreateChatMessage records an in-app or AI message.
func (m *MemoryRepository) CreateChatMessage(_ context.Context, msg domain.ChatMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.chatMessages = append(m.chatMessages, msg)
	return nil
}

// ListChatMessages returns messages for a given channel up to limit.
func (m *MemoryRepository) ListChatMessages(_ context.Context, channelID string, limit int) ([]domain.ChatMessage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.ChatMessage, 0)
	for _, msg := range m.chatMessages {
		if channelID == "" || msg.ChannelID == channelID {
			res = append(res, msg)
		}
	}
	if limit > 0 && len(res) > limit {
		res = res[len(res)-limit:]
	}
	return res, nil
}

