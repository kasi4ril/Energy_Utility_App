// Package service coordinates field work orders, mobile PWA proof-of-work ingestion,
// local AI GPS watchdog verification, and Unit Head reversion to Customer Care.
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

var (
	// ErrAssetNotFound indicates the designated electrical grid asset does not exist.
	ErrAssetNotFound = errors.New("service: target electrical grid asset not found")
	// ErrUnauthorizedAction indicates an actor attempted a mutation outside their assignment.
	ErrUnauthorizedAction = errors.New("service: actor not authorized for this work order action")
)

// WorkflowService orchestrates S04 field maintenance operations.
type WorkflowService struct {
	repo         repository.Repository
	gridRepo     s01repo.Repository
	ticketRepo   s03repo.Repository
	bus          events.EventBus
	aiSidecarURL string
	httpClient   *http.Client
}

// NewWorkflowService instantiates a new Services & Maintenance workflow orchestrator.
func NewWorkflowService(
	repo repository.Repository,
	gridRepo s01repo.Repository,
	ticketRepo s03repo.Repository,
	bus events.EventBus,
	aiSidecarURL string,
) *WorkflowService {
	return &WorkflowService{
		repo:         repo,
		gridRepo:     gridRepo,
		ticketRepo:   ticketRepo,
		bus:          bus,
		aiSidecarURL: strings.TrimRight(aiSidecarURL, "/"),
		httpClient:   &http.Client{Timeout: 3 * time.Second},
	}
}

// CreateWorkOrder dispatches a new field job to a team lead.
func (s *WorkflowService) CreateWorkOrder(
	ctx context.Context,
	ticketID *string,
	unitType domain.UnitType,
	woType domain.WorkOrderType,
	priority domain.WorkOrderPriority,
	unitHeadID, teamLeadID, assetID string,
	ptwRequired bool,
	ptwPermitID, matRef *string,
	instructions string,
) (*domain.MaintenanceWorkOrder, error) {
	// 1. Verify Grid Asset existence
	if _, err := s.gridRepo.GetGridAssetByID(ctx, assetID); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAssetNotFound, err)
	}

	// 2. Generate unique human-readable Order Number: WO-YYYYMMDD-XXXX
	now := time.Now().UTC()
	uuidPart := crypto.MustUUID()
	orderNo := fmt.Sprintf("WO-%s-%s", now.Format("20060102"), strings.ToUpper(uuidPart[:6]))

	wo := domain.MaintenanceWorkOrder{
		ID:                      crypto.MustUUID(),
		OrderNo:                 orderNo,
		TicketID:                ticketID,
		UnitType:                unitType,
		WorkOrderType:           woType,
		Priority:                priority,
		UnitHeadID:              unitHeadID,
		TeamLeadID:              teamLeadID,
		AssetID:                 assetID,
		PTWRequired:             ptwRequired,
		PTWPermitID:             ptwPermitID,
		Status:                  domain.StatusDispatched,
		MaterialsRequisitionRef: matRef,
		Instructions:            instructions,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	if err := domain.ValidateWorkOrderCreation(wo); err != nil {
		return nil, err
	}

	if err := s.repo.CreateWorkOrder(ctx, wo); err != nil {
		return nil, err
	}

	// Emit event for real-time mobile push to Team Lead
	_ = s.bus.Publish(ctx, "services.work_order.dispatched", "s04_services_maint", map[string]interface{}{
		"work_order_id": wo.ID,
		"order_no":      wo.OrderNo,
		"unit_type":     wo.UnitType,
		"team_lead_id":  wo.TeamLeadID,
	})

	return &wo, nil
}

// SubmitProofOfWork ingests field photos, GPS, and technical proof from the Team Lead PWA.
func (s *WorkflowService) SubmitProofOfWork(
	ctx context.Context,
	workOrderID, teamLeadID string,
	capturedGPS s01domain.GPSPoint,
	beforeURL, afterURL, beforeHash, afterHash string,
	meterSeal *string,
	workDetails string,
) (*domain.ProofOfWorkLog, error) {
	// 1. Retrieve work order
	wo, err := s.repo.GetWorkOrderByID(ctx, workOrderID)
	if err != nil {
		return nil, err
	}
	if wo.TeamLeadID != teamLeadID {
		return nil, fmt.Errorf("%w: only assigned team lead can submit field proof", ErrUnauthorizedAction)
	}

	// 2. Retrieve asset coordinates for Haversine validation
	asset, err := s.gridRepo.GetGridAssetByID(ctx, wo.AssetID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAssetNotFound, err)
	}

	// 3. Pure domain verification (photo hashes, seal, PTW, distance)
	dist, gpsValid, err := domain.ValidateProofSubmission(
		wo.UnitType, asset.GPSPoint, capturedGPS,
		beforeHash, afterHash, meterSeal, workDetails,
		wo.PTWRequired, wo.PTWPermitID,
	)
	if err != nil {
		return nil, err
	}

	// 4. Ollama AI Watchdog Verification (synchronous sidecar call with graceful fallback)
	aiVerified, riskScore, reasoning := s.callAIWatchdog(ctx, capturedGPS, asset.GPSPoint, domain.GetGPSToleranceForUnit(wo.UnitType))
	if !gpsValid {
		aiVerified = false
		riskScore = 0.85
		reasoning = fmt.Sprintf("Anti-Ghost Alert: Distance of %.1fm exceeds maximum unit tolerance", dist)
	}

	now := time.Now().UTC()
	proof := domain.ProofOfWorkLog{
		ID:                      crypto.MustUUID(),
		WorkOrderID:             wo.ID,
		TeamLeadID:              teamLeadID,
		CapturedGPS:             capturedGPS,
		DistanceFromAssetMeters: dist,
		BeforePhotoURL:          beforeURL,
		AfterPhotoURL:           afterURL,
		PhotoHashBefore:         beforeHash,
		PhotoHashAfter:          afterHash,
		MeterSealNumber:         meterSeal,
		WorkDetails:             workDetails,
		OllamaGPSVerified:       aiVerified,
		AIRiskScore:             riskScore,
		AIReasoning:             reasoning,
		SubmittedAt:             now,
	}

	if err := s.repo.CreateProofOfWork(ctx, proof); err != nil {
		return nil, err
	}

	// Transition work order status to PROOF_SUBMITTED
	if err := s.repo.UpdateWorkOrderStatus(ctx, wo.ID, domain.StatusProofSubmitted, nil); err != nil {
		return nil, err
	}

	_ = s.bus.Publish(ctx, "services.proof.submitted", "s04_services_maint", map[string]interface{}{
		"work_order_id": wo.ID,
		"team_lead_id":  teamLeadID,
		"gps_verified":  aiVerified,
		"risk_score":    riskScore,
	})

	return &proof, nil
}

// VerifyProofAndRevert allows the Unit Head to inspect proof, approve, or reject for rework.
func (s *WorkflowService) VerifyProofAndRevert(
	ctx context.Context,
	workOrderID, unitHeadID string,
	approved bool,
	notes string,
) (*domain.MaintenanceWorkOrder, error) {
	wo, err := s.repo.GetWorkOrderByID(ctx, workOrderID)
	if err != nil {
		return nil, err
	}

	newStatus, err := domain.ValidateUnitHeadVerification(wo.Status, approved, notes)
	if err != nil {
		return nil, err
	}

	proof, err := s.repo.GetProofByWorkOrderID(ctx, wo.ID)
	if err != nil {
		return nil, err
	}

	// Update proof verification record
	if err := s.repo.UpdateProofVerification(ctx, proof.ID, approved, notes); err != nil {
		return nil, err
	}

	// Update work order
	var reworkNotes *string
	if !approved {
		reworkNotes = &notes
	}
	if err := s.repo.UpdateWorkOrderStatus(ctx, wo.ID, newStatus, reworkNotes); err != nil {
		return nil, err
	}
	wo.Status = newStatus
	wo.ReworkNotes = reworkNotes

	if approved {
		// If linked to a customer service ticket, update ticket to ATTENDED_TO
		if wo.TicketID != nil && s.ticketRepo != nil {
			proofSummary := fmt.Sprintf("Field repair verified by Unit Head: %s", notes)
			_ = s.ticketRepo.UpdateTicketAttended(ctx, *wo.TicketID, proofSummary, wo.TeamLeadID)
		}

		_ = s.bus.Publish(ctx, "services.work_order.attended_verified", "s04_services_maint", map[string]interface{}{
			"work_order_id": wo.ID,
			"ticket_id":     wo.TicketID,
			"unit_head_id":  unitHeadID,
		})
	} else {
		_ = s.bus.Publish(ctx, "services.work_order.rejected_rework", "s04_services_maint", map[string]interface{}{
			"work_order_id": wo.ID,
			"team_lead_id":  wo.TeamLeadID,
			"rework_notes":  notes,
		})
	}

	return wo, nil
}

// callAIWatchdog queries the local Python AI sidecar or falls back to internal rules.
func (s *WorkflowService) callAIWatchdog(
	ctx context.Context,
	reported, asset s01domain.GPSPoint,
	tolerance float64,
) (bool, float64, string) {
	if s.aiSidecarURL == "" {
		dist := domain.CalculateHaversineDistance(reported, asset)
		valid := dist <= tolerance
		return valid, 0.15, "Local heuristic: GPS within tolerance."
	}

	payload := map[string]interface{}{
		"reported_lat":     reported.Latitude,
		"reported_lon":     reported.Longitude,
		"asset_lat":        asset.Latitude,
		"asset_lon":        asset.Longitude,
		"tolerance_meters": tolerance,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.aiSidecarURL+"/internal/ai/verify-gps", bytes.NewReader(data))
	if err != nil {
		return true, 0.20, "AI sidecar unavailable; local validation fallback applied."
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return true, 0.20, "AI sidecar timeout/error; local validation fallback applied."
	}
	defer resp.Body.Close()

	var result struct {
		IsValid   bool    `json:"is_valid"`
		Distance  float64 `json:"distance_meters"`
		Reasoning string  `json:"reasoning"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return true, 0.20, "Failed to decode AI sidecar response; local fallback applied."
	}
	return result.IsValid, 0.10, result.Reasoning
}

