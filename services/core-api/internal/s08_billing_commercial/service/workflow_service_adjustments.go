// Package service coordinates wrong billing resolution gates and revenue protection AI anomaly workflows.
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
)

// ProposeBillingAdjustment creates a proposed adjustment from a customer service dispute.
func (s *WorkflowService) ProposeBillingAdjustment(ctx context.Context, adj domain.BillingAdjustment) (*domain.BillingAdjustment, error) {
	if err := domain.ValidateAdjustmentProposal(adj); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	adj.ID = generateUUID()
	if strings.TrimSpace(adj.AdjustmentNo) == "" {
		adj.AdjustmentNo = fmt.Sprintf("ADJ-%d-%s", now.Unix(), randomHexString(3))
	}
	adj.Status = domain.AdjustmentProposed
	adj.CreatedAt = now
	adj.UpdatedAt = now

	if err := s.repo.CreateBillingAdjustment(ctx, adj); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "billing.adjustment.proposed", "s08_billing_commercial", map[string]interface{}{
			"adjustment_id":     adj.ID,
			"adjustment_no":     adj.AdjustmentNo,
			"ticket_id":         adj.TicketID,
			"customer_id":       adj.CustomerID,
			"adjustment_amount": adj.AdjustmentAmount,
		})
	}

	return &adj, nil
}

// ApproveBillingAdjustment applies Head of Commercial Operations review and reverts ticket to ATTENDED_TO on approval.
func (s *WorkflowService) ApproveBillingAdjustment(ctx context.Context, adjID, headID string, approve bool, notes string) (*domain.BillingAdjustment, error) {
	adj, err := s.repo.GetBillingAdjustmentByID(ctx, adjID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateCommercialHeadApproval(adj, headID, approve, notes); err != nil {
		return nil, err
	}

	adj.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateBillingAdjustment(ctx, *adj); err != nil {
		return nil, err
	}

	if approve {
		if s.ticketRepo != nil && strings.TrimSpace(adj.TicketID) != "" {
			proofSummary := fmt.Sprintf("Wrong billing dispute resolved and approved by Commercial Head: %s (credited ₦%.2f)", notes, adj.AdjustmentAmount)
			_ = s.ticketRepo.UpdateTicketAttended(ctx, adj.TicketID, proofSummary, headID)
		}

		if s.bus != nil {
			_ = s.bus.Publish(ctx, "billing.adjustment.approved", "s08_billing_commercial", map[string]interface{}{
				"adjustment_id": adj.ID,
				"adjustment_no": adj.AdjustmentNo,
				"ticket_id":     adj.TicketID,
				"head_id":       headID,
			})
		}
	} else {
		if s.bus != nil {
			_ = s.bus.Publish(ctx, "billing.adjustment.rejected", "s08_billing_commercial", map[string]interface{}{
				"adjustment_id": adj.ID,
				"adjustment_no": adj.AdjustmentNo,
				"ticket_id":     adj.TicketID,
				"head_id":       headID,
				"reason":        notes,
			})
		}
	}

	return adj, nil
}

// GetBillingAdjustment retrieves an adjustment by ID.
func (s *WorkflowService) GetBillingAdjustment(ctx context.Context, id string) (*domain.BillingAdjustment, error) {
	return s.repo.GetBillingAdjustmentByID(ctx, id)
}

// ListBillingAdjustments lists adjustments by status and customer ID.
func (s *WorkflowService) ListBillingAdjustments(ctx context.Context, status domain.AdjustmentStatus, customerID string) ([]domain.BillingAdjustment, error) {
	return s.repo.ListBillingAdjustments(ctx, status, customerID)
}

// RevenueScanRequest parameters for evaluating physical grid energy balance or churn patterns.
type RevenueScanRequest struct {
	CustomerOrDTID  string  `json:"customer_or_dt_id"`
	EntityType      string  `json:"entity_type"`
	DTBulkKWh       float64 `json:"dt_bulk_kwh"`
	SumMeteredKWh   float64 `json:"sum_metered_kwh"`
	VendingDropPct  float64 `json:"vending_drop_pct"`
	FeederEnergized bool    `json:"feeder_energized"`
	OutageCount     int     `json:"outage_count"`
	PurchaseDropPct float64 `json:"purchase_drop_pct"`
}

// RunRevenueProtectionScan evaluates bypass and churn indicators and raises an AI alert if detected.
func (s *WorkflowService) RunRevenueProtectionScan(ctx context.Context, req RevenueScanRequest) (*domain.RevenueProtectionAlert, error) {
	if strings.TrimSpace(req.CustomerOrDTID) == "" {
		return nil, fmt.Errorf("service: customer_or_dt_id is required")
	}

	now := time.Now().UTC()

	// 1. Check Energy Theft Bypass
	bypassDetected, bypassConf, bypassEvidence := domain.DetectEnergyTheftBypass(
		req.DTBulkKWh, req.SumMeteredKWh, req.VendingDropPct, req.FeederEnergized,
	)
	if bypassDetected {
		alert := domain.RevenueProtectionAlert{
			ID:              generateUUID(),
			AlertNo:         fmt.Sprintf("ALR-%d-%s", now.Unix(), randomHexString(3)),
			CustomerOrDTID:  req.CustomerOrDTID,
			EntityType:      req.EntityType,
			AlertType:       domain.AlertEnergyTheftBypass,
			ConfidenceScore: bypassConf,
			EvidenceSummary: bypassEvidence,
			Status:          domain.AlertOpen,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		if err := s.repo.CreateAIAlert(ctx, alert); err != nil {
			return nil, err
		}
		if s.bus != nil {
			_ = s.bus.Publish(ctx, "billing.revenue_protection.alert_raised", "s08_billing_commercial", map[string]interface{}{
				"alert_id":   alert.ID,
				"alert_no":   alert.AlertNo,
				"alert_type": alert.AlertType,
				"entity_id":  alert.CustomerOrDTID,
			})
		}
		return &alert, nil
	}

	// 2. Check Customer Churn Risk
	churnDetected, churnConf, churnEvidence := domain.EvaluateCustomerChurnRisk(
		req.OutageCount, req.PurchaseDropPct,
	)
	if churnDetected {
		alert := domain.RevenueProtectionAlert{
			ID:              generateUUID(),
			AlertNo:         fmt.Sprintf("ALR-%d-%s", now.Unix(), randomHexString(3)),
			CustomerOrDTID:  req.CustomerOrDTID,
			EntityType:      req.EntityType,
			AlertType:       domain.AlertCustomerChurnRisk,
			ConfidenceScore: churnConf,
			EvidenceSummary: churnEvidence,
			Status:          domain.AlertOpen,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		if err := s.repo.CreateAIAlert(ctx, alert); err != nil {
			return nil, err
		}
		if s.bus != nil {
			_ = s.bus.Publish(ctx, "billing.revenue_protection.alert_raised", "s08_billing_commercial", map[string]interface{}{
				"alert_id":   alert.ID,
				"alert_no":   alert.AlertNo,
				"alert_type": alert.AlertType,
				"entity_id":  alert.CustomerOrDTID,
			})
		}
		return &alert, nil
	}

	return nil, nil
}

// ReviewAIAlert updates alert investigation notes and assignment.
func (s *WorkflowService) ReviewAIAlert(ctx context.Context, alertID string, unitID *string, notes *string, status domain.AIAlertStatus) (*domain.RevenueProtectionAlert, error) {
	alert, err := s.repo.GetAIAlertByID(ctx, alertID)
	if err != nil {
		return nil, err
	}

	if unitID != nil {
		alert.AssignedUnitID = unitID
	}
	if notes != nil {
		alert.InvestigationNotes = notes
	}
	if status != "" {
		alert.Status = status
	}
	alert.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateAIAlert(ctx, *alert); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "billing.revenue_protection.alert_updated", "s08_billing_commercial", map[string]interface{}{
			"alert_id": alert.ID,
			"alert_no": alert.AlertNo,
			"status":   alert.Status,
		})
	}

	return alert, nil
}

// GetAIAlert retrieves alert by ID.
func (s *WorkflowService) GetAIAlert(ctx context.Context, id string) (*domain.RevenueProtectionAlert, error) {
	return s.repo.GetAIAlertByID(ctx, id)
}

// ListAIAlerts filters alerts by status and alert type.
func (s *WorkflowService) ListAIAlerts(ctx context.Context, status domain.AIAlertStatus, alertType domain.AIAlertType) ([]domain.RevenueProtectionAlert, error) {
	return s.repo.ListAIAlerts(ctx, status, alertType)
}
