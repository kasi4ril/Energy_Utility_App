// Package repository implements PostgreSQL persistence for billing adjustments and AI alerts.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
)

// CreateBillingAdjustment stores a proposed adjustment in PostgreSQL.
func (p *PostgresRepository) CreateBillingAdjustment(ctx context.Context, adj domain.BillingAdjustment) error {
	query := `
		INSERT INTO billing_commercial.billing_adjustments (
			id, adjustment_no, ticket_id, customer_id, billing_officer_id,
			old_balance, adjusted_balance, adjustment_amount,
			root_cause_category, root_cause_note,
			head_commercial_ops_id, head_commercial_ops_approved,
			head_commercial_ops_approved_at, head_commercial_ops_notes,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	now := time.Now().UTC()
	_, err := p.db.ExecContext(ctx, query,
		adj.ID, adj.AdjustmentNo, adj.TicketID, adj.CustomerID, adj.BillingOfficerID,
		adj.OldBalance, adj.AdjustedBalance, adj.AdjustmentAmount,
		string(adj.RootCauseCategory), adj.RootCauseNote,
		adj.HeadCommercialOpsID, adj.HeadCommercialOpsApproved,
		adj.HeadCommercialOpsApprovedAt, adj.HeadCommercialOpsNotes,
		string(adj.Status), now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to create billing adjustment: %w", err)
	}
	return nil
}

// GetBillingAdjustmentByID retrieves adjustment by primary ID.
func (p *PostgresRepository) GetBillingAdjustmentByID(ctx context.Context, id string) (*domain.BillingAdjustment, error) {
	query := `
		SELECT id, adjustment_no, ticket_id, customer_id, billing_officer_id,
		       old_balance, adjusted_balance, adjustment_amount,
		       root_cause_category, root_cause_note,
		       head_commercial_ops_id, head_commercial_ops_approved,
		       head_commercial_ops_approved_at, head_commercial_ops_notes,
		       status, created_at, updated_at
		FROM billing_commercial.billing_adjustments
		WHERE id = $1`

	return p.scanAdjustmentRow(p.db.QueryRowContext(ctx, query, id))
}

// GetBillingAdjustmentByNo retrieves adjustment by adjustment number.
func (p *PostgresRepository) GetBillingAdjustmentByNo(ctx context.Context, adjNo string) (*domain.BillingAdjustment, error) {
	query := `
		SELECT id, adjustment_no, ticket_id, customer_id, billing_officer_id,
		       old_balance, adjusted_balance, adjustment_amount,
		       root_cause_category, root_cause_note,
		       head_commercial_ops_id, head_commercial_ops_approved,
		       head_commercial_ops_approved_at, head_commercial_ops_notes,
		       status, created_at, updated_at
		FROM billing_commercial.billing_adjustments
		WHERE adjustment_no = $1`

	return p.scanAdjustmentRow(p.db.QueryRowContext(ctx, query, adjNo))
}

func (p *PostgresRepository) scanAdjustmentRow(row *sql.Row) (*domain.BillingAdjustment, error) {
	var adj domain.BillingAdjustment
	var rootCause, status string
	var headID, notes sql.NullString
	var approvedAt sql.NullTime

	err := row.Scan(
		&adj.ID, &adj.AdjustmentNo, &adj.TicketID, &adj.CustomerID, &adj.BillingOfficerID,
		&adj.OldBalance, &adj.AdjustedBalance, &adj.AdjustmentAmount,
		&rootCause, &adj.RootCauseNote,
		&headID, &adj.HeadCommercialOpsApproved,
		&approvedAt, &notes,
		&status, &adj.CreatedAt, &adj.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: scan adjustment error: %w", err)
	}

	adj.RootCauseCategory = domain.RootCauseCategory(rootCause)
	adj.Status = domain.AdjustmentStatus(status)
	if headID.Valid {
		adj.HeadCommercialOpsID = &headID.String
	}
	if notes.Valid {
		adj.HeadCommercialOpsNotes = &notes.String
	}
	if approvedAt.Valid {
		adj.HeadCommercialOpsApprovedAt = &approvedAt.Time
	}
	return &adj, nil
}

// UpdateBillingAdjustment updates adjustment approval status and notes.
func (p *PostgresRepository) UpdateBillingAdjustment(ctx context.Context, adj domain.BillingAdjustment) error {
	query := `
		UPDATE billing_commercial.billing_adjustments
		SET head_commercial_ops_id = $2,
		    head_commercial_ops_approved = $3,
		    head_commercial_ops_approved_at = $4,
		    head_commercial_ops_notes = $5,
		    status = $6,
		    updated_at = $7
		WHERE id = $1`

	now := time.Now().UTC()
	res, err := p.db.ExecContext(ctx, query,
		adj.ID, adj.HeadCommercialOpsID, adj.HeadCommercialOpsApproved,
		adj.HeadCommercialOpsApprovedAt, adj.HeadCommercialOpsNotes,
		string(adj.Status), now,
	)
	if err != nil {
		return fmt.Errorf("postgres: update adjustment error: %w", err)
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

// ListBillingAdjustments filters adjustments by status and customer ID.
func (p *PostgresRepository) ListBillingAdjustments(ctx context.Context, status domain.AdjustmentStatus, customerID string) ([]domain.BillingAdjustment, error) {
	query := `
		SELECT id, adjustment_no, ticket_id, customer_id, billing_officer_id,
		       old_balance, adjusted_balance, adjustment_amount,
		       root_cause_category, root_cause_note,
		       head_commercial_ops_id, head_commercial_ops_approved,
		       head_commercial_ops_approved_at, head_commercial_ops_notes,
		       status, created_at, updated_at
		FROM billing_commercial.billing_adjustments
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR customer_id = $2)
		ORDER BY created_at DESC`

	rows, err := p.db.QueryContext(ctx, query, string(status), customerID)
	if err != nil {
		return nil, fmt.Errorf("postgres: list adjustments error: %w", err)
	}
	defer rows.Close()

	res := make([]domain.BillingAdjustment, 0)
	for rows.Next() {
		var adj domain.BillingAdjustment
		var rootCause, st string
		var headID, notes sql.NullString
		var approvedAt sql.NullTime

		if err := rows.Scan(
			&adj.ID, &adj.AdjustmentNo, &adj.TicketID, &adj.CustomerID, &adj.BillingOfficerID,
			&adj.OldBalance, &adj.AdjustedBalance, &adj.AdjustmentAmount,
			&rootCause, &adj.RootCauseNote,
			&headID, &adj.HeadCommercialOpsApproved,
			&approvedAt, &notes,
			&st, &adj.CreatedAt, &adj.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: scan adjustment row error: %w", err)
		}
		adj.RootCauseCategory = domain.RootCauseCategory(rootCause)
		adj.Status = domain.AdjustmentStatus(st)
		if headID.Valid {
			adj.HeadCommercialOpsID = &headID.String
		}
		if notes.Valid {
			adj.HeadCommercialOpsNotes = &notes.String
		}
		if approvedAt.Valid {
			adj.HeadCommercialOpsApprovedAt = &approvedAt.Time
		}
		res = append(res, adj)
	}
	return res, nil
}

// CreateAIAlert persists an AI revenue protection alert.
func (p *PostgresRepository) CreateAIAlert(ctx context.Context, alert domain.RevenueProtectionAlert) error {
	query := `
		INSERT INTO billing_commercial.revenue_protection_ai_alerts (
			id, alert_no, customer_or_dt_id, entity_type, alert_type,
			confidence_score, evidence_summary, assigned_investigation_unit_id,
			investigation_notes, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	now := time.Now().UTC()
	_, err := p.db.ExecContext(ctx, query,
		alert.ID, alert.AlertNo, alert.CustomerOrDTID, alert.EntityType, string(alert.AlertType),
		alert.ConfidenceScore, alert.EvidenceSummary, alert.AssignedUnitID,
		alert.InvestigationNotes, string(alert.Status), now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to create AI alert: %w", err)
	}
	return nil
}

// GetAIAlertByID retrieves alert by UUID.
func (p *PostgresRepository) GetAIAlertByID(ctx context.Context, id string) (*domain.RevenueProtectionAlert, error) {
	query := `
		SELECT id, alert_no, customer_or_dt_id, entity_type, alert_type,
		       confidence_score, evidence_summary, assigned_investigation_unit_id,
		       investigation_notes, status, created_at, updated_at
		FROM billing_commercial.revenue_protection_ai_alerts
		WHERE id = $1`

	return p.scanAlertRow(p.db.QueryRowContext(ctx, query, id))
}

// GetAIAlertByNo retrieves alert by unique alert number.
func (p *PostgresRepository) GetAIAlertByNo(ctx context.Context, alertNo string) (*domain.RevenueProtectionAlert, error) {
	query := `
		SELECT id, alert_no, customer_or_dt_id, entity_type, alert_type,
		       confidence_score, evidence_summary, assigned_investigation_unit_id,
		       investigation_notes, status, created_at, updated_at
		FROM billing_commercial.revenue_protection_ai_alerts
		WHERE alert_no = $1`

	return p.scanAlertRow(p.db.QueryRowContext(ctx, query, alertNo))
}

func (p *PostgresRepository) scanAlertRow(row *sql.Row) (*domain.RevenueProtectionAlert, error) {
	var a domain.RevenueProtectionAlert
	var alertType, status string
	var unitID, notes sql.NullString

	err := row.Scan(
		&a.ID, &a.AlertNo, &a.CustomerOrDTID, &a.EntityType, &alertType,
		&a.ConfidenceScore, &a.EvidenceSummary, &unitID,
		&notes, &status, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: scan alert error: %w", err)
	}
	a.AlertType = domain.AIAlertType(alertType)
	a.Status = domain.AIAlertStatus(status)
	if unitID.Valid {
		a.AssignedUnitID = &unitID.String
	}
	if notes.Valid {
		a.InvestigationNotes = &notes.String
	}
	return &a, nil
}

// UpdateAIAlert updates alert status and notes.
func (p *PostgresRepository) UpdateAIAlert(ctx context.Context, alert domain.RevenueProtectionAlert) error {
	query := `
		UPDATE billing_commercial.revenue_protection_ai_alerts
		SET assigned_investigation_unit_id = $2,
		    investigation_notes = $3,
		    status = $4,
		    updated_at = $5
		WHERE id = $1`

	now := time.Now().UTC()
	res, err := p.db.ExecContext(ctx, query,
		alert.ID, alert.AssignedUnitID, alert.InvestigationNotes,
		string(alert.Status), now,
	)
	if err != nil {
		return fmt.Errorf("postgres: update AI alert error: %w", err)
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

// ListAIAlerts filters alerts by status and alert type.
func (p *PostgresRepository) ListAIAlerts(ctx context.Context, status domain.AIAlertStatus, alertType domain.AIAlertType) ([]domain.RevenueProtectionAlert, error) {
	query := `
		SELECT id, alert_no, customer_or_dt_id, entity_type, alert_type,
		       confidence_score, evidence_summary, assigned_investigation_unit_id,
		       investigation_notes, status, created_at, updated_at
		FROM billing_commercial.revenue_protection_ai_alerts
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR alert_type = $2)
		ORDER BY created_at DESC`

	rows, err := p.db.QueryContext(ctx, query, string(status), string(alertType))
	if err != nil {
		return nil, fmt.Errorf("postgres: list AI alerts error: %w", err)
	}
	defer rows.Close()

	res := make([]domain.RevenueProtectionAlert, 0)
	for rows.Next() {
		var a domain.RevenueProtectionAlert
		var aType, st string
		var unitID, notes sql.NullString

		if err := rows.Scan(
			&a.ID, &a.AlertNo, &a.CustomerOrDTID, &a.EntityType, &aType,
			&a.ConfidenceScore, &a.EvidenceSummary, &unitID,
			&notes, &st, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: scan alert row error: %w", err)
		}
		a.AlertType = domain.AIAlertType(aType)
		a.Status = domain.AIAlertStatus(st)
		if unitID.Valid {
			a.AssignedUnitID = &unitID.String
		}
		if notes.Valid {
			a.InvestigationNotes = &notes.String
		}
		res = append(res, a)
	}
	return res, nil
}
