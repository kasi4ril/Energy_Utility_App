// Package service coordinates C-Suite action items and predictive trigger lifecycle.
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

// EnqueueActionItem registers a new pending C-Suite operational decision in the console.
func (s *WorkflowService) EnqueueActionItem(ctx context.Context, item domain.ExecutiveActionItem) (*domain.ExecutiveActionItem, error) {
	if strings.TrimSpace(item.SourceEntityID) == "" || strings.TrimSpace(item.SourceEntityRef) == "" {
		return nil, fmt.Errorf("service: source_entity_id and source_entity_ref are required")
	}

	now := time.Now().UTC()
	item.ID = generateUUID()
	item.Status = domain.ActionStatusPending
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := s.repo.CreateActionItem(ctx, item); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "executive.action_item.enqueued", "s09_executive_command", map[string]interface{}{
			"action_id":   item.ID,
			"action_type": item.ActionType,
			"target_role": item.TargetRole,
			"source_ref":  item.SourceEntityRef,
		})
	}

	return &item, nil
}

// DecideActionItem records a C-Suite digital signature, enforcing role segregation of duties.
func (s *WorkflowService) DecideActionItem(ctx context.Context, itemID, decidedByStaffID string, role domain.ExecutiveRole, approve bool, notes string) (*domain.ExecutiveActionItem, error) {
	item, err := s.repo.GetActionItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateActionItemDecision(item, decidedByStaffID, role, approve, notes); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateActionItem(ctx, *item); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "executive.action_item.decided", "s09_executive_command", map[string]interface{}{
			"action_id":   item.ID,
			"action_type": item.ActionType,
			"target_role": item.TargetRole,
			"status":      item.Status,
			"decided_by":  decidedByStaffID,
		})
	}

	return item, nil
}

// ListActionItems queries action items matching role and status filters.
func (s *WorkflowService) ListActionItems(ctx context.Context, targetRole domain.ExecutiveRole, status domain.ActionItemStatus) ([]domain.ExecutiveActionItem, error) {
	return s.repo.ListActionItems(ctx, targetRole, status)
}

// AcknowledgePredictiveTrigger updates a trigger's acknowledgement and status.
func (s *WorkflowService) AcknowledgePredictiveTrigger(ctx context.Context, triggerID, staffID, status string) (*domain.PredictiveMaintenanceTrigger, error) {
	t, err := s.repo.GetPredictiveTriggerByID(ctx, triggerID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	t.AcknowledgedByStaffID = &staffID
	t.AcknowledgedAt = &now
	if status != "" {
		t.Status = status
	}
	t.UpdatedAt = now

	if err := s.repo.UpdatePredictiveTrigger(ctx, *t); err != nil {
		return nil, err
	}

	return t, nil
}

// ListPredictiveTriggers filters predictive maintenance alarms.
func (s *WorkflowService) ListPredictiveTriggers(ctx context.Context, status string, model domain.PredictiveModelName, limit int) ([]domain.PredictiveMaintenanceTrigger, error) {
	return s.repo.ListPredictiveTriggers(ctx, status, model, limit)
}

