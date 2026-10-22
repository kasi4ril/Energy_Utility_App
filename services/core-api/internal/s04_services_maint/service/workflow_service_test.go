package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	s03domain "disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

func setupTestService(t *testing.T, aiServerURL string) (*WorkflowService, *repository.MemoryRepository, *s01repo.MemoryRepository, *s03repo.MemoryRepository) {
	t.Helper()

	repo := repository.NewMemoryRepository()
	gridRepo := s01repo.NewMemoryRepository()
	ticketRepo := s03repo.NewMemoryRepository()
	bus := events.NewMemoryEventBus()

	// Seed Grid Asset: DT at Jos coordinates
	dtAsset := s01domain.GridAsset{
		ID:        "dt-asset-uuid-1",
		AssetCode: "DT-JOS-001",
		OrgUnitID: "unit-id-1",
		AssetType: s01domain.AssetTransformerDT,
		GPSPoint:  s01domain.GPSPoint{Latitude: 9.9320, Longitude: 8.8920},
	}
	if err := gridRepo.CreateGridAsset(context.Background(), dtAsset); err != nil {
		t.Fatalf("failed to seed grid asset: %v", err)
	}

	// Seed Ticket in S03
	ticket := s03domain.ServiceTicket{
		ID:             "ticket-uuid-1",
		TicketNo:       "TIC-20261022-0001",
		Channel:        s03domain.ChannelCallCenter,
		Category:       s03domain.CategoryTransformerFault,
		Description:    "500kVA DT phase 2 blown fuse.",
		AssignedUnitID: "unit-id-1",
		Status:         s03domain.StatusLogged,
		SLADueAt:       time.Now().Add(24 * time.Hour),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	if err := ticketRepo.CreateTicket(context.Background(), ticket); err != nil {
		t.Fatalf("failed to seed ticket: %v", err)
	}

	svc := NewWorkflowService(repo, gridRepo, ticketRepo, bus, aiServerURL)
	return svc, repo, gridRepo, ticketRepo
}

func TestWorkflowService_WorkOrderDispatchAndNegative(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _ := setupTestService(t, "")

	ticketID := "ticket-uuid-1"

	// 1. Success dispatch
	wo, err := svc.CreateWorkOrder(
		ctx, &ticketID, domain.UnitTransformerMaint,
		domain.WOCorrectiveRepair, domain.PriorityHigh,
		"unit-head-1", "team-lead-1", "dt-asset-uuid-1",
		false, nil, nil,
		"Replace blown 200A HRC fuse on phase 2.",
	)
	if err != nil {
		t.Fatalf("CreateWorkOrder failed: %v", err)
	}
	if wo.Status != domain.StatusDispatched || wo.OrderNo == "" {
		t.Fatalf("unexpected work order state: %+v", wo)
	}

	// 2. Negative: Missing Grid Asset
	_, err = svc.CreateWorkOrder(
		ctx, &ticketID, domain.UnitTransformerMaint,
		domain.WOCorrectiveRepair, domain.PriorityHigh,
		"unit-head-1", "team-lead-1", "non-existent-asset",
		false, nil, nil,
		"Replace blown 200A HRC fuse on phase 2.",
	)
	if err == nil {
		t.Errorf("expected error for non-existent asset")
	}

	// 3. Negative: Short instructions
	_, err = svc.CreateWorkOrder(
		ctx, &ticketID, domain.UnitTransformerMaint,
		domain.WOCorrectiveRepair, domain.PriorityHigh,
		"unit-head-1", "team-lead-1", "dt-asset-uuid-1",
		false, nil, nil,
		"fix",
	)
	if err == nil {
		t.Errorf("expected error on short instructions")
	}
}

func TestWorkflowService_ProofSubmissionAndAIWatchdog(t *testing.T) {
	ctx := context.Background()

	// Start a mock Python AI Sidecar server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"is_valid":true,"distance_meters":15.2,"reasoning":"Mock AI GPS verified."}`))
	}))
	defer mockServer.Close()

	svc, repo, _, _ := setupTestService(t, mockServer.URL)

	wo, _ := svc.CreateWorkOrder(
		ctx, nil, domain.UnitTransformerMaint,
		domain.WOCorrectiveRepair, domain.PriorityMedium,
		"unit-head-1", "team-lead-1", "dt-asset-uuid-1",
		false, nil, nil,
		"Check oil levels and top up silica gel.",
	)

	capturedGPS := s01domain.GPSPoint{Latitude: 9.9321, Longitude: 8.8920} // ~11m
	bHash := "1111111111111111111111111111111111111111111111111111111111111111"
	aHash := "2222222222222222222222222222222222222222222222222222222222222222"
	details := "Replaced desiccant breather and topped up 10L transformer oil."

	// 1. Unauthorized Team Lead
	_, err := svc.SubmitProofOfWork(
		ctx, wo.ID, "impostor-lead",
		capturedGPS, "/before.jpg", "/after.jpg",
		bHash, aHash, nil, details,
	)
	if err == nil {
		t.Errorf("expected unauthorized error for non-assigned team lead")
	}

	// 2. Successful proof submission
	proof, err := svc.SubmitProofOfWork(
		ctx, wo.ID, "team-lead-1",
		capturedGPS, "/before.jpg", "/after.jpg",
		bHash, aHash, nil, details,
	)
	if err != nil {
		t.Fatalf("SubmitProofOfWork failed: %v", err)
	}
	if !proof.OllamaGPSVerified || proof.AIRiskScore != 0.10 {
		t.Fatalf("unexpected proof state: %+v", proof)
	}

	// Verify work order status updated
	updatedWO, _ := repo.GetWorkOrderByID(ctx, wo.ID)
	if updatedWO.Status != domain.StatusProofSubmitted {
		t.Fatalf("expected status %s, got %s", domain.StatusProofSubmitted, updatedWO.Status)
	}

	// 3. Fallback AI Sidecar test (bad sidecar URL)
	badSvc, _, _, _ := setupTestService(t, "http://127.0.0.1:9999") // Closed port
	wo2, _ := badSvc.CreateWorkOrder(
		ctx, nil, domain.UnitTransformerMaint,
		domain.WOCorrectiveRepair, domain.PriorityMedium,
		"unit-head-1", "team-lead-1", "dt-asset-uuid-1",
		false, nil, nil,
		"Check oil levels and top up silica gel.",
	)
	proof2, err := badSvc.SubmitProofOfWork(
		ctx, wo2.ID, "team-lead-1",
		capturedGPS, "/before.jpg", "/after.jpg",
		bHash, aHash, nil, details,
	)
	if err != nil || !proof2.OllamaGPSVerified {
		t.Fatalf("expected graceful fallback on closed AI sidecar, got err: %v, proof: %+v", err, proof2)
	}
}

func TestWorkflowService_VerificationAndCustomerCareReversion(t *testing.T) {
	ctx := context.Background()
	svc, _, _, ticketRepo := setupTestService(t, "")

	ticketID := "ticket-uuid-1"
	wo, _ := svc.CreateWorkOrder(
		ctx, &ticketID, domain.UnitTransformerMaint,
		domain.WOCorrectiveRepair, domain.PriorityCritical,
		"unit-head-1", "team-lead-1", "dt-asset-uuid-1",
		false, nil, nil,
		"Restore phase 2 power on Tundun Wada DT.",
	)

	capturedGPS := s01domain.GPSPoint{Latitude: 9.9321, Longitude: 8.8920}
	bHash := "1111111111111111111111111111111111111111111111111111111111111111"
	aHash := "2222222222222222222222222222222222222222222222222222222222222222"
	details := "Replaced blown 200A HRC fuse and balanced phase load."

	// Submit proof
	_, err := svc.SubmitProofOfWork(
		ctx, wo.ID, "team-lead-1",
		capturedGPS, "/before.jpg", "/after.jpg",
		bHash, aHash, nil, details,
	)
	if err != nil {
		t.Fatalf("SubmitProofOfWork failed: %v", err)
	}

	// 1. Unit Head Rejects Proof for Rework
	rejectedWO, err := svc.VerifyProofAndRevert(
		ctx, wo.ID, "unit-head-1", false,
		"Secondary phase jumper is slack; please retorque clamps.",
	)
	if err != nil || rejectedWO.Status != domain.StatusRejectedRework {
		t.Fatalf("expected rejection, got: %+v, err: %v", rejectedWO, err)
	}

	// Resubmit proof after rework
	bHash2 := "3333333333333333333333333333333333333333333333333333333333333333"
	aHash2 := "4444444444444444444444444444444444444444444444444444444444444444"
	details2 := "Retorqued secondary clamps to 45Nm and re-measured contact resistance."
	_, err = svc.SubmitProofOfWork(
		ctx, wo.ID, "team-lead-1",
		capturedGPS, "/before2.jpg", "/after2.jpg",
		bHash2, aHash2, nil, details2,
	)
	if err != nil {
		t.Fatalf("re-submission failed: %v", err)
	}

	// 2. Unit Head Approves Proof
	approvedWO, err := svc.VerifyProofAndRevert(
		ctx, wo.ID, "unit-head-1", true,
		"Torque confirmed and phase voltage verified at 230V.",
	)
	if err != nil || approvedWO.Status != domain.StatusAttendedVerified {
		t.Fatalf("expected approval, got: %+v, err: %v", approvedWO, err)
	}

	// 3. Verify Customer Care ticket status was updated to ATTENDED_TO!
	ticket, err := ticketRepo.GetTicketByID(ctx, ticketID)
	if err != nil || ticket.Status != s03domain.StatusAttendedTo {
		t.Fatalf("expected ticket status %s, got %s, err: %v", s03domain.StatusAttendedTo, ticket.Status, err)
	}
	if ticket.AttendedProofNote == "" {
		t.Errorf("expected non-empty attended proof note on reverted ticket")
	}
}
