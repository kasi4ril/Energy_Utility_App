// Package service orchestrates C-Suite operational intelligence aggregation,
// predictive maintenance prognostics, and the executive action console for S09.
package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	s04repo "disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	s05repo "disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	s06repo "disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	s07repo "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	s08repo "disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

// WorkflowService coordinates executive dashboard aggregation and predictive maintenance triggers.
type WorkflowService struct {
	repo       repository.Repository
	bus        events.EventBus
	maintRepo  s04repo.Repository
	genRepo    s05repo.Repository
	energyRepo s06repo.Repository
	auditRepo  s07repo.Repository
	billRepo   s08repo.Repository
}

// NewWorkflowService constructs an initialized WorkflowService instance.
func NewWorkflowService(
	repo repository.Repository,
	bus events.EventBus,
	maintRepo s04repo.Repository,
	genRepo s05repo.Repository,
	energyRepo s06repo.Repository,
	auditRepo s07repo.Repository,
	billRepo s08repo.Repository,
) *WorkflowService {
	return &WorkflowService{
		repo:       repo,
		bus:        bus,
		maintRepo:  maintRepo,
		genRepo:    genRepo,
		energyRepo: energyRepo,
		auditRepo:  auditRepo,
		billRepo:   billRepo,
	}
}

// GetLiveDashboardSummary aggregates real-time operational telemetry across all departments for the C-Suite view.
func (s *WorkflowService) GetLiveDashboardSummary(ctx context.Context) (*domain.CSuiteDashboardSummary, error) {
	now := time.Now().UTC()
	summary := &domain.CSuiteDashboardSummary{
		Timestamp: now,
		RevenueSummary: domain.RevenuePulseSummary{
			ChannelBreakdown: make(map[string]float64),
		},
	}

	// 1. Generation & Hydro Telemetry (S05)
	if s.genRepo != nil {
		logs, err := s.genRepo.ListTelemetryLogs(ctx, "", "", 10)
		if err == nil {
			for _, l := range logs {
				summary.TotalGenerationMW += l.MWGenerated
			}
		}
	}

	// 2. Dispatch, Intake & Feeder Statuses (S06)
	if s.energyRepo != nil {
		latestDispatch, err := s.energyRepo.GetLatestDispatchLog(ctx)
		if err == nil && latestDispatch != nil {
			summary.TCNGridIntakeMW = latestDispatch.TCNMWIn
		}
		summary.TotalAvailablePowerMW = summary.TotalGenerationMW + summary.TCNGridIntakeMW

		feeders, err := s.energyRepo.ListFeederStatuses(ctx)
		if err == nil {
			summary.FeedersSummary.TotalFeeders = len(feeders)
			for _, f := range feeders {
				switch string(f.State) {
				case "IN_SERVICE":
					summary.FeedersSummary.InServiceCount++
				case "OUT_ON_FAULT":
					summary.FeedersSummary.OutOnFaultCount++
				case "OUT_ON_PERMIT_PTW":
					summary.FeedersSummary.OutOnPTWCount++
				case "OUT_ON_LOAD_SHEDDING":
					summary.FeedersSummary.OutOnRationingCount++
				}
			}
		}
	}

	// 3. Billing Inflow & Integrity Radar (S08)
	if s.billRepo != nil {
		txs, err := s.billRepo.ListVendingTransactions(ctx, "", "", 100)
		if err == nil {
			for _, tx := range txs {
				summary.RevenueSummary.TotalCollectedAmount += tx.AmountPaid
				if tx.PaymentChannelCategory == "POSTPAID_DIRECT" {
					summary.RevenueSummary.PostpaidAmount += tx.AmountPaid
				} else {
					summary.RevenueSummary.PrepaidAmount += tx.AmountPaid
				}
				platform := tx.PaymentPlatformSource
				if platform == "" {
					platform = "OTHER"
				}
				summary.RevenueSummary.ChannelBreakdown[platform] += tx.AmountPaid
			}
		}

		theftAlerts, err := s.billRepo.ListAIAlerts(ctx, "OPEN", "ENERGY_THEFT_BYPASS")
		if err == nil {
			summary.IntegritySummary.UnresolvedTheftAlertsCount = len(theftAlerts)
		}
		churnAlerts, err := s.billRepo.ListAIAlerts(ctx, "OPEN", "CUSTOMER_CHURN_RISK")
		if err == nil {
			summary.IntegritySummary.HighChurnRiskAccountsCount = len(churnAlerts)
		}
	}

	// 4. Stores & Audit Pending Items (S07)
	if s.auditRepo != nil {
		pendingReqs, err := s.auditRepo.ListRequisitions(ctx, "MD_PROCUREMENT_PENDING", "")
		if err == nil {
			summary.IntegritySummary.AuditPendingItemsCount = len(pendingReqs)
		}
	}

	// 5. Active Predictive Maintenance Triggers (S09 Local)
	activeTriggers, err := s.repo.ListPredictiveTriggers(ctx, "ACTIVE", "", 10)
	if err == nil {
		summary.ActivePredictiveTriggers = activeTriggers
	}

	// 6. Pending C-Suite Action Items (S09 Local)
	pendingActions, err := s.repo.ListActionItems(ctx, "", domain.ActionStatusPending)
	if err == nil {
		summary.PendingExecutiveActionItems = pendingActions
	}

	return summary, nil
}

// CaptureExecutiveSnapshot persists an aggregated operational snapshot.
func (s *WorkflowService) CaptureExecutiveSnapshot(ctx context.Context, snapshot domain.ExecutiveSnapshot) (*domain.ExecutiveSnapshot, error) {
	if strings.TrimSpace(snapshot.SnapshotDate) == "" || strings.TrimSpace(snapshot.ShiftWindow) == "" || strings.TrimSpace(snapshot.CapturedByStaffID) == "" {
		return nil, fmt.Errorf("service: snapshot_date, shift_window, and captured_by_staff_id are required")
	}

	now := time.Now().UTC()
	snapshot.ID = generateUUID()
	snapshot.CreatedAt = now

	if err := s.repo.CreateExecutiveSnapshot(ctx, snapshot); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "executive.snapshot.captured", "s09_executive_command", map[string]interface{}{
			"snapshot_id":   snapshot.ID,
			"snapshot_date": snapshot.SnapshotDate,
			"shift_window":  snapshot.ShiftWindow,
			"total_gen_mw":  snapshot.TotalGenerationMW,
			"total_rev":     snapshot.TotalRevenueCollected,
			"captured_by":   snapshot.CapturedByStaffID,
		})
	}

	return &snapshot, nil
}

// GetLatestSnapshot retrieves the latest operational snapshot.
func (s *WorkflowService) GetLatestSnapshot(ctx context.Context) (*domain.ExecutiveSnapshot, error) {
	return s.repo.GetLatestExecutiveSnapshot(ctx)
}

// ListSnapshots lists historical snapshots up to the limit.
func (s *WorkflowService) ListSnapshots(ctx context.Context, limit int) ([]domain.ExecutiveSnapshot, error) {
	return s.repo.ListExecutiveSnapshots(ctx, limit)
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func randomHexString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
