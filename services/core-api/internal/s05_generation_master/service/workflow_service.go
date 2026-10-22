// Package service orchestrates generation telemetry collection, remote plant call roll-calls,
// Google Forms HMAC webhook ingestion, GM material approvals, and workshop fabrication loops.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

var (
	// ErrPlantNotFound indicates the generating station was not located in grid assets.
	ErrPlantNotFound = errors.New("service: generating plant asset not found")
	// ErrInvalidWebhookSignature indicates the Google Forms webhook HMAC signature failed cryptographic verification.
	ErrInvalidWebhookSignature = errors.New("service: invalid HMAC signature on incoming Google Forms webhook")
)

// WorkflowService coordinates generation telemetry, remote roll-call, and workshop workflows.
type WorkflowService struct {
	repo       repository.Repository
	gridRepo   s01repo.Repository
	bus        events.EventBus
	hmacSecret string
}

// NewWorkflowService constructs a new Generation workflow orchestrator.
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

// RecordTelemetry validates physics, checks early warning thresholds, and logs generation telemetry.
func (s *WorkflowService) RecordTelemetry(ctx context.Context, log domain.MasterStationTelemetryLog) ([]string, error) {
	if _, err := s.gridRepo.GetGridAssetByID(ctx, log.PlantID); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPlantNotFound, err)
	}

	alerts, err := domain.EvaluateTelemetryPhysics(log)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if log.ID == "" {
		log.ID = crypto.MustUUID()
	}
	if log.RecordedAt.IsZero() {
		log.RecordedAt = now
	}

	if err := s.repo.CreateTelemetryLog(ctx, log); err != nil {
		return nil, err
	}

	// Publish generation telemetry event
	_ = s.bus.Publish(ctx, "generation.telemetry.recorded", "s05_generation_master", map[string]interface{}{
		"log_id":       log.ID,
		"plant_id":     log.PlantID,
		"mw_generated": log.MWGenerated,
		"state":        log.PlantState,
		"alerts_count": len(alerts),
	})

	if len(alerts) > 0 {
		_ = s.bus.Publish(ctx, "generation.telemetry.alert", "s05_generation_master", map[string]interface{}{
			"plant_id": log.PlantID,
			"alerts":   alerts,
		})
	}

	return alerts, nil
}

// RecordRollcall captures physical personnel presence on duty at a remote plant.
func (s *WorkflowService) RecordRollcall(ctx context.Context, rc domain.RemotePlantShiftRollcall) (*domain.RemotePlantShiftRollcall, error) {
	if _, err := s.gridRepo.GetGridAssetByID(ctx, rc.PlantID); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPlantNotFound, err)
	}

	if err := domain.ValidateRollcallSubmission(rc); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if rc.ID == "" {
		rc.ID = crypto.MustUUID()
	}
	if rc.CreatedAt.IsZero() {
		rc.CreatedAt = now
	}

	if err := s.repo.CreateRollcall(ctx, rc); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "generation.rollcall.recorded", "s05_generation_master", map[string]interface{}{
		"rollcall_id": rc.ID,
		"shift_id":    rc.ShiftID,
		"plant_id":    rc.PlantID,
		"staff_count": len(rc.OperatorNames) + len(rc.ProtectionEngineers) + len(rc.Technicians),
	})

	return &rc, nil
}

// IngestGoogleFormsWebhook validates HMAC signature and ingests collated shift telemetry and roll-calls.
func (s *WorkflowService) IngestGoogleFormsWebhook(ctx context.Context, rawPayload []byte, signature string) ([]string, error) {
	// 1. Verify HMAC SHA-256 signature
	if !crypto.VerifyRowHMAC(s.hmacSecret, string(rawPayload), signature) {
		return nil, ErrInvalidWebhookSignature
	}

	// 2. Decode standard Google Forms JSON payload
	var formPayload struct {
		ShiftID          string `json:"shift_id"`
		ShiftWindow      string `json:"shift_window"`
		MasterOperatorID string `json:"master_operator_id"`
		Plants           []struct {
			PlantID             string   `json:"plant_id"`
			PlantState          string   `json:"plant_state"`
			WaterLevelM         float64  `json:"water_level_m"`
			MWGenerated         float64  `json:"mw_generated"`
			StatorTempC         float64  `json:"stator_temp_c"`
			VibrationMMS        float64  `json:"vibration_mm_s"`
			OilPressureBar      float64  `json:"oil_pressure_bar"`
			OperatorNames       []string `json:"operator_names"`
			ProtectionEngineers []string `json:"protection_engineers"`
			Technicians         []string `json:"technicians"`
			HandoverNotes       string   `json:"handover_notes"`
		} `json:"plants"`
	}

	if err := json.Unmarshal(rawPayload, &formPayload); err != nil {
		return nil, fmt.Errorf("service: malformed Google Forms payload: %w", err)
	}

	var allAlerts []string
	now := time.Now().UTC()

	for _, p := range formPayload.Plants {
		// Log Telemetry
		tLog := domain.MasterStationTelemetryLog{
			ID:                  crypto.MustUUID(),
			ShiftID:             formPayload.ShiftID,
			ShiftWindow:         domain.ShiftWindow(formPayload.ShiftWindow),
			PlantID:             p.PlantID,
			PlantState:          domain.PlantState(p.PlantState),
			WaterLevelM:         p.WaterLevelM,
			MWGenerated:         p.MWGenerated,
			StatorWindingTempC:  p.StatorTempC,
			BearingVibrationMMS: p.VibrationMMS,
			LubeOilPressureBar:  p.OilPressureBar,
			LoggedByOperatorID:  formPayload.MasterOperatorID,
			RecordedAt:          now,
		}

		alerts, err := s.RecordTelemetry(ctx, tLog)
		if err == nil && len(alerts) > 0 {
			allAlerts = append(allAlerts, alerts...)
		}

		// Log Roll-call if operators provided
		if len(p.OperatorNames) > 0 {
			rc := domain.RemotePlantShiftRollcall{
				ID:                  crypto.MustUUID(),
				ShiftID:             formPayload.ShiftID,
				ShiftWindow:         domain.ShiftWindow(formPayload.ShiftWindow),
				PlantID:             p.PlantID,
				OperatorNames:       p.OperatorNames,
				ProtectionEngineers: p.ProtectionEngineers,
				Technicians:         p.Technicians,
				HandoverNotes:       p.HandoverNotes,
				MasterOperatorID:    formPayload.MasterOperatorID,
				CreatedAt:           now,
			}
			_, _ = s.RecordRollcall(ctx, rc)
		}
	}

	return allAlerts, nil
}

// ReportPlantFault logs a new fault and initiates GM material approval or workshop dispatch.
func (s *WorkflowService) ReportPlantFault(
	ctx context.Context,
	plantID, operatorID string,
	wing domain.TechnicalWing,
	faultDesc, severity string,
	requiresWorkshop bool,
) (*domain.PlantFaultAndWorkshopJob, error) {
	if _, err := s.gridRepo.GetGridAssetByID(ctx, plantID); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPlantNotFound, err)
	}

	now := time.Now().UTC()
	uuidPart := crypto.MustUUID()
	jobNo := fmt.Sprintf("JOB-%s-%s", now.Format("20060102"), strings.ToUpper(uuidPart[:6]))

	wStatus := domain.WorkshopNotApplicable
	if requiresWorkshop {
		wStatus = domain.WorkshopPendingDispatch
	}

	job := domain.PlantFaultAndWorkshopJob{
		ID:                   crypto.MustUUID(),
		JobNo:                jobNo,
		PlantID:              plantID,
		ReportedByOperatorID: operatorID,
		Wing:                 wing,
		FaultDescription:     faultDesc,
		Severity:             severity,
		RequiresWorkshop:     requiresWorkshop,
		WorkshopStatus:       wStatus,
		Status:               domain.JobStatusReported,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := domain.ValidateFaultJobCreation(job); err != nil {
		return nil, err
	}

	if err := s.repo.CreatePlantJob(ctx, job); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "generation.fault.reported", "s05_generation_master", map[string]interface{}{
		"job_id":            job.ID,
		"job_no":            job.JobNo,
		"plant_id":          job.PlantID,
		"wing":              job.Wing,
		"requires_workshop": job.RequiresWorkshop,
	})

	return &job, nil
}

// ApproveGMJob applies the General Manager's approval or rejection decision.
func (s *WorkflowService) ApproveGMJob(ctx context.Context, jobID, gmStaffID string, approved bool, notes string) (*domain.PlantFaultAndWorkshopJob, error) {
	job, err := s.repo.GetPlantJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	if _, err := domain.ValidateGMApproval(job.Status, approved, notes); err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePlantJobGMApproval(ctx, jobID, approved, notes); err != nil {
		return nil, err
	}

	job.GMApproved = approved
	job.GMApprovalNotes = &notes
	if approved {
		job.Status = domain.JobStatusGMApproved
	} else {
		job.Status = domain.JobStatusCancelled
	}

	_ = s.bus.Publish(ctx, "generation.gm.decision", "s05_generation_master", map[string]interface{}{
		"job_id":    job.ID,
		"approved":  approved,
		"gm_staff":  gmStaffID,
		"status":    job.Status,
	})

	return job, nil
}

// UpdateWorkshopFabrication advances machining or winding reconditioning in the central workshop.
func (s *WorkflowService) UpdateWorkshopFabrication(ctx context.Context, jobID string, targetStatus domain.WorkshopStatus, materialsRef *string) (*domain.PlantFaultAndWorkshopJob, error) {
	job, err := s.repo.GetPlantJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateWorkshopTransition(job.WorkshopStatus, targetStatus); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateWorkshopStatus(ctx, jobID, targetStatus, materialsRef); err != nil {
		return nil, err
	}

	job.WorkshopStatus = targetStatus
	if materialsRef != nil {
		job.StoreRequisitionRef = materialsRef
	}

	_ = s.bus.Publish(ctx, "generation.workshop.status_updated", "s05_generation_master", map[string]interface{}{
		"job_id":          job.ID,
		"workshop_status": targetStatus,
	})

	return job, nil
}

