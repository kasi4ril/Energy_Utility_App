// Package repository provides persistence interfaces and in-memory mock storage
// for the S05 Generation Department & Master Station Hub bounded context.
package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
)

var (
	// ErrNotFound indicates the requested entity was not found in storage.
	ErrNotFound = errors.New("repository: entity not found")
	// ErrConflict indicates an entity with the specified identifier or job number already exists.
	ErrConflict = errors.New("repository: entity conflict detected")
)

// Repository specifies the persistence contract for generation telemetry, roll-calls, and workshop jobs.
type Repository interface {
	CreateTelemetryLog(ctx context.Context, log domain.MasterStationTelemetryLog) error
	GetLatestTelemetryByPlant(ctx context.Context, plantID string) (*domain.MasterStationTelemetryLog, error)
	ListTelemetryLogs(ctx context.Context, shiftID, plantID string, limit int) ([]domain.MasterStationTelemetryLog, error)

	CreateRollcall(ctx context.Context, rc domain.RemotePlantShiftRollcall) error
	ListRollcallsByShift(ctx context.Context, shiftID string) ([]domain.RemotePlantShiftRollcall, error)

	CreatePlantJob(ctx context.Context, job domain.PlantFaultAndWorkshopJob) error
	GetPlantJobByID(ctx context.Context, id string) (*domain.PlantFaultAndWorkshopJob, error)
	UpdatePlantJobGMApproval(ctx context.Context, id string, approved bool, notes string) error
	UpdateWorkshopStatus(ctx context.Context, id string, status domain.WorkshopStatus, materialsRef *string) error
	ListPlantJobs(ctx context.Context, plantID string, wing domain.TechnicalWing, status domain.PlantJobStatus) ([]domain.PlantFaultAndWorkshopJob, error)
}

// MemoryRepository provides a concurrent, thread-safe in-memory mock of Repository.
type MemoryRepository struct {
	mu            sync.RWMutex
	telemetryLogs map[string]domain.MasterStationTelemetryLog
	rollcalls     map[string]domain.RemotePlantShiftRollcall
	plantJobs     map[string]domain.PlantFaultAndWorkshopJob
	jobsByNo      map[string]string // jobNo -> ID
}

// NewMemoryRepository constructs a new in-memory generation storage handler.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		telemetryLogs: make(map[string]domain.MasterStationTelemetryLog),
		rollcalls:     make(map[string]domain.RemotePlantShiftRollcall),
		plantJobs:     make(map[string]domain.PlantFaultAndWorkshopJob),
		jobsByNo:      make(map[string]string),
	}
}

// CreateTelemetryLog records a new hourly/shift telemetry reading.
func (m *MemoryRepository) CreateTelemetryLog(_ context.Context, log domain.MasterStationTelemetryLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.telemetryLogs[log.ID]; exists {
		return ErrConflict
	}
	m.telemetryLogs[log.ID] = log
	return nil
}

// GetLatestTelemetryByPlant retrieves the most recent telemetry reading for a given plant.
func (m *MemoryRepository) GetLatestTelemetryByPlant(_ context.Context, plantID string) (*domain.MasterStationTelemetryLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var latest *domain.MasterStationTelemetryLog
	for _, l := range m.telemetryLogs {
		if l.PlantID == plantID {
			if latest == nil || l.RecordedAt.After(latest.RecordedAt) {
				lCopy := l
				latest = &lCopy
			}
		}
	}
	if latest == nil {
		return nil, ErrNotFound
	}
	return latest, nil
}

// ListTelemetryLogs returns telemetry records filtered by shift or plant.
func (m *MemoryRepository) ListTelemetryLogs(_ context.Context, shiftID, plantID string, limit int) ([]domain.MasterStationTelemetryLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 {
		limit = 50
	}
	var results []domain.MasterStationTelemetryLog
	for _, l := range m.telemetryLogs {
		if shiftID != "" && l.ShiftID != shiftID {
			continue
		}
		if plantID != "" && l.PlantID != plantID {
			continue
		}
		results = append(results, l)
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

// CreateRollcall persists a remote plant shift roll-call record.
func (m *MemoryRepository) CreateRollcall(_ context.Context, rc domain.RemotePlantShiftRollcall) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rollcalls[rc.ID]; exists {
		return ErrConflict
	}
	m.rollcalls[rc.ID] = rc
	return nil
}

// ListRollcallsByShift returns roll-calls for a specific shift window.
func (m *MemoryRepository) ListRollcallsByShift(_ context.Context, shiftID string) ([]domain.RemotePlantShiftRollcall, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.RemotePlantShiftRollcall
	for _, rc := range m.rollcalls {
		if rc.ShiftID == shiftID {
			results = append(results, rc)
		}
	}
	return results, nil
}

// CreatePlantJob records a new plant fault or workshop job.
func (m *MemoryRepository) CreatePlantJob(_ context.Context, job domain.PlantFaultAndWorkshopJob) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plantJobs[job.ID]; exists {
		return ErrConflict
	}
	if _, exists := m.jobsByNo[job.JobNo]; exists {
		return ErrConflict
	}

	m.plantJobs[job.ID] = job
	m.jobsByNo[job.JobNo] = job.ID
	return nil
}

// GetPlantJobByID retrieves a fault/workshop job by UUID.
func (m *MemoryRepository) GetPlantJobByID(_ context.Context, id string) (*domain.PlantFaultAndWorkshopJob, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	job, exists := m.plantJobs[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &job, nil
}

// UpdatePlantJobGMApproval updates GM approval state, notes, and status.
func (m *MemoryRepository) UpdatePlantJobGMApproval(_ context.Context, id string, approved bool, notes string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.plantJobs[id]
	if !exists {
		return ErrNotFound
	}

	job.GMApproved = approved
	job.GMApprovalNotes = &notes
	if approved {
		job.Status = domain.JobStatusGMApproved
	} else {
		job.Status = domain.JobStatusCancelled
	}
	job.UpdatedAt = time.Now().UTC()
	m.plantJobs[id] = job
	return nil
}

// UpdateWorkshopStatus advances workshop fabrication status and attaches store requisition references.
func (m *MemoryRepository) UpdateWorkshopStatus(_ context.Context, id string, status domain.WorkshopStatus, materialsRef *string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.plantJobs[id]
	if !exists {
		return ErrNotFound
	}

	job.WorkshopStatus = status
	if materialsRef != nil {
		job.StoreRequisitionRef = materialsRef
	}
	if status == domain.WorkshopCompletedDelivered {
		job.Status = domain.JobStatusResolvedReinstalled
	} else if status == domain.WorkshopMachiningInProgress || status == domain.WorkshopWindingFabrication {
		job.Status = domain.JobStatusWorkshopInProgress
	}
	job.UpdatedAt = time.Now().UTC()
	m.plantJobs[id] = job
	return nil
}

// ListPlantJobs returns jobs matching plant, wing, and status filters.
func (m *MemoryRepository) ListPlantJobs(_ context.Context, plantID string, wing domain.TechnicalWing, status domain.PlantJobStatus) ([]domain.PlantFaultAndWorkshopJob, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.PlantFaultAndWorkshopJob
	for _, j := range m.plantJobs {
		if plantID != "" && j.PlantID != plantID {
			continue
		}
		if wing != "" && j.Wing != wing {
			continue
		}
		if status != "" && j.Status != status {
			continue
		}
		results = append(results, j)
	}
	return results, nil
}

