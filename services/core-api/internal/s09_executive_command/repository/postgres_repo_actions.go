// Package repository implements PostgreSQL persistence for executive action console items.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

// CreateActionItem inserts a new pending C-Suite action item.
func (p *PostgresRepository) CreateActionItem(ctx context.Context, item domain.ExecutiveActionItem) error {
	query := `
		INSERT INTO executive_command.executive_action_items (
			id, action_type, target_role, source_entity_id, source_entity_ref,
			source_department, summary, amount, priority, status,
			decided_by_staff_id, decided_at, decision_notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	now := time.Now().UTC()
	_, err := p.db.ExecContext(ctx, query,
		item.ID, string(item.ActionType), string(item.TargetRole), item.SourceEntityID,
		item.SourceEntityRef, item.SourceDepartment, item.Summary, item.Amount,
		item.Priority, string(item.Status), item.DecidedByStaffID, item.DecidedAt,
		item.DecisionNotes, now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to insert executive action item: %w", err)
	}
	return nil
}

// GetActionItemByID retrieves an action item by UUID.
func (p *PostgresRepository) GetActionItemByID(ctx context.Context, id string) (*domain.ExecutiveActionItem, error) {
	query := `
		SELECT id, action_type, target_role, source_entity_id, source_entity_ref,
		       source_department, summary, amount, priority, status,
		       decided_by_staff_id, decided_at, decision_notes, created_at, updated_at
		FROM executive_command.executive_action_items
		WHERE id = $1`

	return p.scanActionItemRow(p.db.QueryRowContext(ctx, query, id))
}

// GetActionItemByRef retrieves an action item by source entity reference.
func (p *PostgresRepository) GetActionItemByRef(ctx context.Context, ref string) (*domain.ExecutiveActionItem, error) {
	query := `
		SELECT id, action_type, target_role, source_entity_id, source_entity_ref,
		       source_department, summary, amount, priority, status,
		       decided_by_staff_id, decided_at, decision_notes, created_at, updated_at
		FROM executive_command.executive_action_items
		WHERE source_entity_ref = $1`

	return p.scanActionItemRow(p.db.QueryRowContext(ctx, query, ref))
}

func (p *PostgresRepository) scanActionItemRow(row *sql.Row) (*domain.ExecutiveActionItem, error) {
	var item domain.ExecutiveActionItem
	var aType, tRole, st string
	var staffID, notes sql.NullString
	var decidedAt sql.NullTime

	err := row.Scan(
		&item.ID, &aType, &tRole, &item.SourceEntityID, &item.SourceEntityRef,
		&item.SourceDepartment, &item.Summary, &item.Amount, &item.Priority,
		&st, &staffID, &decidedAt, &notes, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: scan action item error: %w", err)
	}

	item.ActionType = domain.ActionItemType(aType)
	item.TargetRole = domain.ExecutiveRole(tRole)
	item.Status = domain.ActionItemStatus(st)
	if staffID.Valid {
		item.DecidedByStaffID = &staffID.String
	}
	if decidedAt.Valid {
		item.DecidedAt = &decidedAt.Time
	}
	if notes.Valid {
		item.DecisionNotes = &notes.String
	}
	return &item, nil
}

// UpdateActionItem updates decision status and notes.
func (p *PostgresRepository) UpdateActionItem(ctx context.Context, item domain.ExecutiveActionItem) error {
	query := `
		UPDATE executive_command.executive_action_items
		SET status = $2,
		    decided_by_staff_id = $3,
		    decided_at = $4,
		    decision_notes = $5,
		    updated_at = $6
		WHERE id = $1`

	now := time.Now().UTC()
	res, err := p.db.ExecContext(ctx, query,
		item.ID, string(item.Status), item.DecidedByStaffID,
		item.DecidedAt, item.DecisionNotes, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: update action item error: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListActionItems filters action items by target role and status.
func (p *PostgresRepository) ListActionItems(ctx context.Context, targetRole domain.ExecutiveRole, status domain.ActionItemStatus) ([]domain.ExecutiveActionItem, error) {
	query := `
		SELECT id, action_type, target_role, source_entity_id, source_entity_ref,
		       source_department, summary, amount, priority, status,
		       decided_by_staff_id, decided_at, decision_notes, created_at, updated_at
		FROM executive_command.executive_action_items
		WHERE ($1 = '' OR target_role = $1)
		  AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC`

	rows, err := p.db.QueryContext(ctx, query, string(targetRole), string(status))
	if err != nil {
		return nil, fmt.Errorf("postgres: list action items error: %w", err)
	}
	defer rows.Close()

	res := make([]domain.ExecutiveActionItem, 0)
	for rows.Next() {
		var item domain.ExecutiveActionItem
		var aType, tRole, st string
		var staffID, notes sql.NullString
		var decidedAt sql.NullTime

		if err := rows.Scan(
			&item.ID, &aType, &tRole, &item.SourceEntityID, &item.SourceEntityRef,
			&item.SourceDepartment, &item.Summary, &item.Amount, &item.Priority,
			&st, &staffID, &decidedAt, &notes, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: scan action item row: %w", err)
		}
		item.ActionType = domain.ActionItemType(aType)
		item.TargetRole = domain.ExecutiveRole(tRole)
		item.Status = domain.ActionItemStatus(st)
		if staffID.Valid {
			item.DecidedByStaffID = &staffID.String
		}
		if decidedAt.Valid {
			item.DecidedAt = &decidedAt.Time
		}
		if notes.Valid {
			item.DecisionNotes = &notes.String
		}
		res = append(res, item)
	}
	return res, nil
}

