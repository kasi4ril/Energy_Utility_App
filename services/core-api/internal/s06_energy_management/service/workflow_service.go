// Package service orchestrates bulk energy dispatch telemetry, Google Forms HMAC webhook ingestion,
// priority feeder load shedding algorithms, and the 5-step Safety Permit-to-Work (PTW) lifecycle.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

var (
	// ErrInvalidWebhookSignature indicates Google Forms HMAC-SHA256 signature verification failed.
	ErrInvalidWebhookSignature = errors.New("service: invalid HMAC signature on incoming Google Forms webhook")
	// ErrNoDispatchData indicates no intake telemetry has been logged yet for rationing.
	ErrNoDispatchData = errors.New("service: no bulk intake telemetry available for rationing calculation")
	// ErrAssetNotFound indicates the referenced grid asset does not exist.
	ErrAssetNotFound = errors.New("service: referenced grid asset not found")
	// ErrUnauthorizedApprover indicates the role provided is unauthorized for PTW dual approval.
	ErrUnauthorizedApprover = errors.New("service: unauthorized role for PTW approval (must be HEAD_SERVICES or SAFETY_REP)")
)

// WorkflowService coordinates energy dispatch telemetry, priority rationing, and PTWs.
type WorkflowService struct {
	repo       repository.Repository
	gridRepo   s01repo.Repository
	bus        events.EventBus
	hmacSecret string
}

// NewWorkflowService constructs an energy management service orchestrator.
func NewWorkflowService(
	repo repository.Repository,
	gridRepo s01repo.Repository,
	bus events.EventBus,
	hmacSecret string,
) *WorkflowService {
	return &WorkflowService{
		repo:       repo,
		gridRepo:   gridRepo,
		bus:        bus,
		hmacSecret: hmacSecret,
	}
}

// RecordBulkIntake validates generation and TCN intake and persists hourly dispatch telemetry.
func (s *WorkflowService) RecordBulkIntake(ctx context.Context, log domain.DispatchTelemetryLog) (*domain.DispatchTelemetryLog, error) {
	totalMW, err := domain.CalculateBulkMW(log.GenerationMWIn, log.TCNMWIn)
	if err != nil {
		return nil, err
	}
	log.TotalBulkMW = totalMW

	if log.ID == "" {
		log.ID = crypto.MustUUID()
	}
	if log.LoggedAt.IsZero() {
		log.LoggedAt = time.Now().UTC()
	}

	if err := s.repo.CreateDispatchLog(ctx, log); err != nil {
		return nil, fmt.Errorf("service: failed to persist dispatch log: %w", err)
	}

	_ = s.bus.Publish(ctx, "energy.bulk_intake_logged", "s06_energy_management", map[string]interface{}{
		"log_id":        log.ID,
		"shift_id":      log.ShiftID,
		"generation_mw": log.GenerationMWIn,
		"tcn_mw":        log.TCNMWIn,
		"total_bulk_mw": log.TotalBulkMW,
		"logged_at":     log.LoggedAt,
	})

	return &log, nil
}

// IngestGoogleFormsWebhook validates HMAC-SHA256 signature and records Form B dispatch telemetry.
func (s *WorkflowService) IngestGoogleFormsWebhook(ctx context.Context, rawPayload []byte, signature string) (*domain.DispatchTelemetryLog, error) {
	if s.hmacSecret != "" && !crypto.VerifyRowHMAC(s.hmacSecret, string(rawPayload), signature) {
		return nil, ErrInvalidWebhookSignature
	}

	var formPayload struct {
		ShiftID            string             `json:"shift_id"`
		ShiftWindow        domain.ShiftWindow `json:"shift_window"`
		OperatorID         string             `json:"operator_id"`
		GenerationMWIn     float64            `json:"generation_mw_in"`
		TCNMWIn            float64            `json:"tcn_mw_in"`
		SubstationRollcall string             `json:"substation_rollcall"`
		HandoverNotes      string             `json:"handover_notes"`
		Feeders            []struct {
			FeederAssetID string             `json:"feeder_asset_id"`
			FeederCode    string             `json:"feeder_code"`
			CurrentLoadMW float64            `json:"current_load_mw"`
			CurrentAmps   float64            `json:"current_amps"`
			State         domain.FeederState `json:"state"`
			PriorityRank  int                `json:"priority_rank"`
			TariffBand    string             `json:"tariff_band"`
		} `json:"feeders"`
	}

	if err := json.Unmarshal(rawPayload, &formPayload); err != nil {
		return nil, fmt.Errorf("service: malformed Google Forms payload: %w", err)
	}

	totalMW, err := domain.CalculateBulkMW(formPayload.GenerationMWIn, formPayload.TCNMWIn)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	log := domain.DispatchTelemetryLog{
		ID:                          crypto.MustUUID(),
		ShiftID:                     formPayload.ShiftID,
		ShiftWindow:                 formPayload.ShiftWindow,
		GenerationMWIn:              formPayload.GenerationMWIn,
		TCNMWIn:                     formPayload.TCNMWIn,
		TotalBulkMW:                 totalMW,
		SubstationStaffRollcallJSON: formPayload.SubstationRollcall,
		HandoverNotes:               formPayload.HandoverNotes,
		OperatorID:                  formPayload.OperatorID,
		LoggedAt:                    now,
	}

	if err := s.repo.CreateDispatchLog(ctx, log); err != nil {
		return nil, fmt.Errorf("service: failed to persist webhook dispatch log: %w", err)
	}

	for _, f := range formPayload.Feeders {
		state := f.State
		if state == "" {
			state = domain.FeederInService
		}
		status := domain.FeederStatus{
			ID:            crypto.MustUUID(),
			FeederAssetID: f.FeederAssetID,
			FeederCode:    f.FeederCode,
			CurrentLoadMW: f.CurrentLoadMW,
			CurrentAmps:   f.CurrentAmps,
			State:         state,
			PriorityRank:  f.PriorityRank,
			TariffBand:    s01domain.TariffBand(f.TariffBand),
			UpdatedAt:     now,
		}
		_ = s.repo.UpsertFeederStatus(ctx, status)
	}

	_ = s.bus.Publish(ctx, "energy.google_forms_ingested", "s06_energy_management", map[string]interface{}{
		"log_id":        log.ID,
		"shift_id":      log.ShiftID,
		"total_bulk_mw": log.TotalBulkMW,
		"feeder_count":  len(formPayload.Feeders),
	})

	return &log, nil
}

// RunFeederRationing executes the priority rationing algorithm across all registered feeders.
func (s *WorkflowService) RunFeederRationing(ctx context.Context, bulkMW float64) ([]domain.RationingDecision, error) {
	if bulkMW <= 0.0 {
		latest, err := s.repo.GetLatestDispatchLog(ctx)
		if err != nil {
			return nil, ErrNoDispatchData
		}
		bulkMW = latest.TotalBulkMW
	}

	feeders, err := s.repo.ListFeederStatuses(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list feeders: %w", err)
	}

	decisions := domain.ComputePriorityFeederRationing(bulkMW, feeders)

	// Update feeder states according to decision
	feederMap := make(map[string]*domain.FeederStatus)
	for i := range feeders {
		feederMap[feeders[i].FeederCode] = &feeders[i]
	}

	for _, d := range decisions {
		if f, exists := feederMap[d.FeederCode]; exists {
			if f.State != domain.FeederOutOnPTW && f.State != domain.FeederOutOnFault {
				f.State = d.TargetState
				f.CurrentLoadMW = d.AllocatedMW
				narrative := d.Justification
				f.OutageNarrative = &narrative
				_ = s.repo.UpsertFeederStatus(ctx, *f)
			}
		}
	}

	_ = s.bus.Publish(ctx, "energy.feeder_rationing_applied", "s06_energy_management", map[string]interface{}{
		"available_bulk_mw": bulkMW,
		"decisions_count":   len(decisions),
	})

	return decisions, nil
}

// GetFeederStatus retrieves telemetry and state for a feeder by asset UUID.
func (s *WorkflowService) GetFeederStatus(ctx context.Context, assetID string) (*domain.FeederStatus, error) {
	return s.repo.GetFeederStatusByAssetID(ctx, assetID)
}

// ListFeederStatuses retrieves all known feeders.
func (s *WorkflowService) ListFeederStatuses(ctx context.Context) ([]domain.FeederStatus, error) {
	return s.repo.ListFeederStatuses(ctx)
}

// ListDispatchLogs retrieves intake logs filtered by shift.
func (s *WorkflowService) ListDispatchLogs(ctx context.Context, shiftID string, limit int) ([]domain.DispatchTelemetryLog, error) {
	return s.repo.ListDispatchLogs(ctx, shiftID, limit)
}
