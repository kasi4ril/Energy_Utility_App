// Package repository implements PostgreSQL persistence for executive snapshots and predictive maintenance triggers.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

// PostgresRepository provides a PostgreSQL-backed implementation of Repository.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new PostgreSQL repository for S09.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateExecutiveSnapshot stores an aggregated telemetry snapshot.
func (p *PostgresRepository) CreateExecutiveSnapshot(ctx context.Context, s domain.ExecutiveSnapshot) error {
	query := `
		INSERT INTO executive_command.executive_snapshots (
			id, snapshot_date, shift_window, total_generation_mw, tcn_grid_intake_mw,
			total_system_load_mw, feeders_in_service_count, feeders_out_fault_count,
			feeders_out_ptw_count, feeders_out_rationing_count, total_revenue_collected,
			prepaid_revenue_share_pct, unresolved_theft_alerts_count, unresolved_churn_alerts_count,
			audit_pending_procurements_count, unreturned_scrap_count, notes, captured_by_staff_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	now := s.CreatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}

	_, err := p.db.ExecContext(ctx, query,
		s.ID, s.SnapshotDate, s.ShiftWindow, s.TotalGenerationMW, s.TCNGridIntakeMW,
		s.TotalSystemLoadMW, s.FeedersInServiceCount, s.FeedersOutFaultCount,
		s.FeedersOutPTWCount, s.FeedersOutRationingCount, s.TotalRevenueCollected,
		s.PrepaidRevenueSharePct, s.UnresolvedTheftAlertsCount, s.UnresolvedChurnAlertsCount,
		s.AuditPendingProcurementsCount, s.UnreturnedScrapCount, s.Notes, s.CapturedByStaffID, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to insert executive snapshot: %w", err)
	}
	return nil
}

// GetExecutiveSnapshotByID retrieves a snapshot by UUID.
func (p *PostgresRepository) GetExecutiveSnapshotByID(ctx context.Context, id string) (*domain.ExecutiveSnapshot, error) {
	query := `
		SELECT id, snapshot_date::text, shift_window, total_generation_mw, tcn_grid_intake_mw,
		       total_system_load_mw, feeders_in_service_count, feeders_out_fault_count,
		       feeders_out_ptw_count, feeders_out_rationing_count, total_revenue_collected,
		       prepaid_revenue_share_pct, unresolved_theft_alerts_count, unresolved_churn_alerts_count,
		       audit_pending_procurements_count, unreturned_scrap_count, COALESCE(notes, ''),
		       captured_by_staff_id, created_at
		FROM executive_command.executive_snapshots
		WHERE id = $1`

	return p.scanSnapshotRow(p.db.QueryRowContext(ctx, query, id))
}

// GetLatestExecutiveSnapshot queries the most recent snapshot.
func (p *PostgresRepository) GetLatestExecutiveSnapshot(ctx context.Context) (*domain.ExecutiveSnapshot, error) {
	query := `
		SELECT id, snapshot_date::text, shift_window, total_generation_mw, tcn_grid_intake_mw,
		       total_system_load_mw, feeders_in_service_count, feeders_out_fault_count,
		       feeders_out_ptw_count, feeders_out_rationing_count, total_revenue_collected,
		       prepaid_revenue_share_pct, unresolved_theft_alerts_count, unresolved_churn_alerts_count,
		       audit_pending_procurements_count, unreturned_scrap_count, COALESCE(notes, ''),
		       captured_by_staff_id, created_at
		FROM executive_command.executive_snapshots
		ORDER BY created_at DESC
		LIMIT 1`

	return p.scanSnapshotRow(p.db.QueryRowContext(ctx, query))
}

func (p *PostgresRepository) scanSnapshotRow(row *sql.Row) (*domain.ExecutiveSnapshot, error) {
	var s domain.ExecutiveSnapshot
	err := row.Scan(
		&s.ID, &s.SnapshotDate, &s.ShiftWindow, &s.TotalGenerationMW, &s.TCNGridIntakeMW,
		&s.TotalSystemLoadMW, &s.FeedersInServiceCount, &s.FeedersOutFaultCount,
		&s.FeedersOutPTWCount, &s.FeedersOutRationingCount, &s.TotalRevenueCollected,
		&s.PrepaidRevenueSharePct, &s.UnresolvedTheftAlertsCount, &s.UnresolvedChurnAlertsCount,
		&s.AuditPendingProcurementsCount, &s.UnreturnedScrapCount, &s.Notes,
		&s.CapturedByStaffID, &s.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: scan snapshot error: %w", err)
	}
	return &s, nil
}

// ListExecutiveSnapshots lists historical snapshots up to limit.
func (p *PostgresRepository) ListExecutiveSnapshots(ctx context.Context, limit int) ([]domain.ExecutiveSnapshot, error) {
	query := `
		SELECT id, snapshot_date::text, shift_window, total_generation_mw, tcn_grid_intake_mw,
		       total_system_load_mw, feeders_in_service_count, feeders_out_fault_count,
		       feeders_out_ptw_count, feeders_out_rationing_count, total_revenue_collected,
		       prepaid_revenue_share_pct, unresolved_theft_alerts_count, unresolved_churn_alerts_count,
		       audit_pending_procurements_count, unreturned_scrap_count, COALESCE(notes, ''),
		       captured_by_staff_id, created_at
		FROM executive_command.executive_snapshots
		ORDER BY created_at DESC
		LIMIT $1`

	if limit <= 0 {
		limit = 30
	}

	rows, err := p.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres: list snapshots error: %w", err)
	}
	defer rows.Close()

	res := make([]domain.ExecutiveSnapshot, 0)
	for rows.Next() {
		var s domain.ExecutiveSnapshot
		if err := rows.Scan(
			&s.ID, &s.SnapshotDate, &s.ShiftWindow, &s.TotalGenerationMW, &s.TCNGridIntakeMW,
			&s.TotalSystemLoadMW, &s.FeedersInServiceCount, &s.FeedersOutFaultCount,
			&s.FeedersOutPTWCount, &s.FeedersOutRationingCount, &s.TotalRevenueCollected,
			&s.PrepaidRevenueSharePct, &s.UnresolvedTheftAlertsCount, &s.UnresolvedChurnAlertsCount,
			&s.AuditPendingProcurementsCount, &s.UnreturnedScrapCount, &s.Notes,
			&s.CapturedByStaffID, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: scan snapshot row: %w", err)
		}
		res = append(res, s)
	}
	return res, nil
}

// CreatePredictiveTrigger stores an automated predictive maintenance alarm.
func (p *PostgresRepository) CreatePredictiveTrigger(ctx context.Context, t domain.PredictiveMaintenanceTrigger) error {
	query := `
		INSERT INTO executive_command.predictive_maintenance_triggers (
			id, trigger_no, asset_id, asset_type, model_name, risk_score,
			risk_level, telemetry_summary_json, recommended_action,
			auto_work_order_created, linked_work_order_id, acknowledged_by_staff_id,
			acknowledged_at, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	now := time.Now().UTC()
	var payload sql.NullString
	if t.TelemetrySummaryJSON != "" {
		payload = sql.NullString{String: t.TelemetrySummaryJSON, Valid: true}
	}

	_, err := p.db.ExecContext(ctx, query,
		t.ID, t.TriggerNo, t.AssetID, string(t.AssetType), string(t.ModelName),
		t.RiskScore, string(t.RiskLevel), payload, t.RecommendedAction,
		t.AutoWorkOrderCreated, t.LinkedWorkOrderID, t.AcknowledgedByStaffID,
		t.AcknowledgedAt, t.Status, now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to insert predictive trigger: %w", err)
	}
	return nil
}

// GetPredictiveTriggerByID retrieves a trigger by UUID.
func (p *PostgresRepository) GetPredictiveTriggerByID(ctx context.Context, id string) (*domain.PredictiveMaintenanceTrigger, error) {
	query := `
		SELECT id, trigger_no, asset_id, asset_type, model_name, risk_score,
		       risk_level, telemetry_summary_json, recommended_action,
		       auto_work_order_created, linked_work_order_id, acknowledged_by_staff_id,
		       acknowledged_at, status, created_at, updated_at
		FROM executive_command.predictive_maintenance_triggers
		WHERE id = $1`

	return p.scanTriggerRow(p.db.QueryRowContext(ctx, query, id))
}

// GetPredictiveTriggerByNo retrieves a trigger by trigger number.
func (p *PostgresRepository) GetPredictiveTriggerByNo(ctx context.Context, triggerNo string) (*domain.PredictiveMaintenanceTrigger, error) {
	query := `
		SELECT id, trigger_no, asset_id, asset_type, model_name, risk_score,
		       risk_level, telemetry_summary_json, recommended_action,
		       auto_work_order_created, linked_work_order_id, acknowledged_by_staff_id,
		       acknowledged_at, status, created_at, updated_at
		FROM executive_command.predictive_maintenance_triggers
		WHERE trigger_no = $1`

	return p.scanTriggerRow(p.db.QueryRowContext(ctx, query, triggerNo))
}

func (p *PostgresRepository) scanTriggerRow(row *sql.Row) (*domain.PredictiveMaintenanceTrigger, error) {
	var t domain.PredictiveMaintenanceTrigger
	var aType, mName, rLevel string
	var payload, linkedWO, ackStaff sql.NullString
	var ackAt sql.NullTime

	err := row.Scan(
		&t.ID, &t.TriggerNo, &t.AssetID, &aType, &mName, &t.RiskScore,
		&rLevel, &payload, &t.RecommendedAction,
		&t.AutoWorkOrderCreated, &linkedWO, &ackStaff,
		&ackAt, &t.Status, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: scan trigger error: %w", err)
	}

	t.AssetType = domain.AssetCategory(aType)
	t.ModelName = domain.PredictiveModelName(mName)
	t.RiskLevel = domain.RiskLevel(rLevel)
	if payload.Valid {
		t.TelemetrySummaryJSON = payload.String
	}
	if linkedWO.Valid {
		t.LinkedWorkOrderID = &linkedWO.String
	}
	if ackStaff.Valid {
		t.AcknowledgedByStaffID = &ackStaff.String
	}
	if ackAt.Valid {
		t.AcknowledgedAt = &ackAt.Time
	}
	return &t, nil
}

// UpdatePredictiveTrigger updates acknowledgment and work order linkage.
func (p *PostgresRepository) UpdatePredictiveTrigger(ctx context.Context, t domain.PredictiveMaintenanceTrigger) error {
	query := `
		UPDATE executive_command.predictive_maintenance_triggers
		SET auto_work_order_created = $2,
		    linked_work_order_id = $3,
		    acknowledged_by_staff_id = $4,
		    acknowledged_at = $5,
		    status = $6,
		    updated_at = $7
		WHERE id = $1`

	now := time.Now().UTC()
	res, err := p.db.ExecContext(ctx, query,
		t.ID, t.AutoWorkOrderCreated, t.LinkedWorkOrderID,
		t.AcknowledgedByStaffID, t.AcknowledgedAt, t.Status, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: update trigger error: %w", err)
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

// ListPredictiveTriggers filters triggers by status and model.
func (p *PostgresRepository) ListPredictiveTriggers(ctx context.Context, status string, model domain.PredictiveModelName, limit int) ([]domain.PredictiveMaintenanceTrigger, error) {
	query := `
		SELECT id, trigger_no, asset_id, asset_type, model_name, risk_score,
		       risk_level, telemetry_summary_json, recommended_action,
		       auto_work_order_created, linked_work_order_id, acknowledged_by_staff_id,
		       acknowledged_at, status, created_at, updated_at
		FROM executive_command.predictive_maintenance_triggers
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR model_name = $2)
		ORDER BY created_at DESC
		LIMIT $3`

	if limit <= 0 {
		limit = 50
	}

	rows, err := p.db.QueryContext(ctx, query, status, string(model), limit)
	if err != nil {
		return nil, fmt.Errorf("postgres: list triggers error: %w", err)
	}
	defer rows.Close()

	res := make([]domain.PredictiveMaintenanceTrigger, 0)
	for rows.Next() {
		var t domain.PredictiveMaintenanceTrigger
		var aType, mName, rLevel string
		var payload, linkedWO, ackStaff sql.NullString
		var ackAt sql.NullTime

		if err := rows.Scan(
			&t.ID, &t.TriggerNo, &t.AssetID, &aType, &mName, &t.RiskScore,
			&rLevel, &payload, &t.RecommendedAction,
			&t.AutoWorkOrderCreated, &linkedWO, &ackStaff,
			&ackAt, &t.Status, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: scan trigger row error: %w", err)
		}
		t.AssetType = domain.AssetCategory(aType)
		t.ModelName = domain.PredictiveModelName(mName)
		t.RiskLevel = domain.RiskLevel(rLevel)
		if payload.Valid {
			t.TelemetrySummaryJSON = payload.String
		}
		if linkedWO.Valid {
			t.LinkedWorkOrderID = &linkedWO.String
		}
		if ackStaff.Valid {
			t.AcknowledgedByStaffID = &ackStaff.String
		}
		if ackAt.Valid {
			t.AcknowledgedAt = &ackAt.Time
		}
		res = append(res, t)
	}
	return res, nil
}

