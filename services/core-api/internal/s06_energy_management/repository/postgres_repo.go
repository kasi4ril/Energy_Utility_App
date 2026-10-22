// Package repository implements the PostgreSQL storage engine for S06 Energy Management.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
)

// PostgresRepository provides production PostgreSQL operations for the energy_management schema.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new Postgres repository handler.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateDispatchLog inserts bulk energy intake and roll-call data into dispatch_telemetry_logs.
func (r *PostgresRepository) CreateDispatchLog(ctx context.Context, log domain.DispatchTelemetryLog) error {
	query := `
		INSERT INTO energy_management.dispatch_telemetry_logs
		(id, shift_id, shift_window, generation_mw_in, tcn_mw_in, total_bulk_mw, substation_staff_rollcall_json, handover_notes, operator_id, logged_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.ShiftID, log.ShiftWindow,
		log.GenerationMWIn, log.TCNMWIn, log.TotalBulkMW,
		log.SubstationStaffRollcallJSON, log.HandoverNotes,
		log.OperatorID, log.LoggedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateDispatchLog failed: %w", err)
	}
	return nil
}

// GetLatestDispatchLog retrieves the most recent bulk intake entry.
func (r *PostgresRepository) GetLatestDispatchLog(ctx context.Context) (*domain.DispatchTelemetryLog, error) {
	query := `
		SELECT id, shift_id, shift_window, generation_mw_in, tcn_mw_in, total_bulk_mw, substation_staff_rollcall_json, handover_notes, operator_id, logged_at
		FROM energy_management.dispatch_telemetry_logs
		ORDER BY logged_at DESC LIMIT 1`

	var l domain.DispatchTelemetryLog
	err := r.db.QueryRowContext(ctx, query).Scan(
		&l.ID, &l.ShiftID, &l.ShiftWindow,
		&l.GenerationMWIn, &l.TCNMWIn, &l.TotalBulkMW,
		&l.SubstationStaffRollcallJSON, &l.HandoverNotes,
		&l.OperatorID, &l.LoggedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetLatestDispatchLog failed: %w", err)
	}
	return &l, nil
}

// ListDispatchLogs returns historical dispatch logs by shift.
func (r *PostgresRepository) ListDispatchLogs(ctx context.Context, shiftID string, limit int) ([]domain.DispatchTelemetryLog, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, shift_id, shift_window, generation_mw_in, tcn_mw_in, total_bulk_mw, substation_staff_rollcall_json, handover_notes, operator_id, logged_at
		FROM energy_management.dispatch_telemetry_logs
		WHERE ($1 = '' OR shift_id = $1)
		ORDER BY logged_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, shiftID, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListDispatchLogs failed: %w", err)
	}
	defer rows.Close()

	var logs []domain.DispatchTelemetryLog
	for rows.Next() {
		var l domain.DispatchTelemetryLog
		if err := rows.Scan(
			&l.ID, &l.ShiftID, &l.ShiftWindow,
			&l.GenerationMWIn, &l.TCNMWIn, &l.TotalBulkMW,
			&l.SubstationStaffRollcallJSON, &l.HandoverNotes,
			&l.OperatorID, &l.LoggedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// UpsertFeederStatus updates or inserts feeder operational state and telemetry.
func (r *PostgresRepository) UpsertFeederStatus(ctx context.Context, status domain.FeederStatus) error {
	query := `
		INSERT INTO energy_management.feeder_status_ledger
		(id, feeder_asset_id, feeder_code, current_load_mw, current_amps, state, outage_reason_code, outage_narrative, priority_rank, tariff_band, active_ptw_id, scheduled_restoration_time, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP)
		ON CONFLICT (feeder_asset_id) DO UPDATE SET
			current_load_mw = EXCLUDED.current_load_mw,
			current_amps = EXCLUDED.current_amps,
			state = EXCLUDED.state,
			outage_reason_code = EXCLUDED.outage_reason_code,
			outage_narrative = EXCLUDED.outage_narrative,
			priority_rank = EXCLUDED.priority_rank,
			tariff_band = EXCLUDED.tariff_band,
			active_ptw_id = EXCLUDED.active_ptw_id,
			scheduled_restoration_time = EXCLUDED.scheduled_restoration_time,
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.ExecContext(ctx, query,
		status.ID, status.FeederAssetID, status.FeederCode,
		status.CurrentLoadMW, status.CurrentAmps, string(status.State),
		status.OutageReasonCode, status.OutageNarrative,
		status.PriorityRank, string(status.TariffBand),
		status.ActivePTWID, status.ScheduledRestorationTime)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpsertFeederStatus failed: %w", err)
	}
	return nil
}

// GetFeederStatusByAssetID retrieves current telemetry for a feeder asset.
func (r *PostgresRepository) GetFeederStatusByAssetID(ctx context.Context, assetID string) (*domain.FeederStatus, error) {
	query := `
		SELECT id, feeder_asset_id, feeder_code, current_load_mw, current_amps, state, outage_reason_code, outage_narrative, priority_rank, tariff_band, active_ptw_id, scheduled_restoration_time, updated_at
		FROM energy_management.feeder_status_ledger WHERE feeder_asset_id = $1::uuid`

	var f domain.FeederStatus
	var reasonCode, narrative, ptwID sql.NullString
	var restoTime sql.NullTime
	var band string
	err := r.db.QueryRowContext(ctx, query, assetID).Scan(
		&f.ID, &f.FeederAssetID, &f.FeederCode,
		&f.CurrentLoadMW, &f.CurrentAmps, &f.State,
		&reasonCode, &narrative, &f.PriorityRank, &band,
		&ptwID, &restoTime, &f.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetFeederStatusByAssetID failed: %w", err)
	}
	f.TariffBand = s01domain.TariffBand(band)
	if reasonCode.Valid {
		f.OutageReasonCode = &reasonCode.String
	}
	if narrative.Valid {
		f.OutageNarrative = &narrative.String
	}
	if ptwID.Valid {
		f.ActivePTWID = &ptwID.String
	}
	if restoTime.Valid {
		f.ScheduledRestorationTime = &restoTime.Time
	}
	return &f, nil
}

// ListFeederStatuses retrieves all feeder statuses ordered by priority.
func (r *PostgresRepository) ListFeederStatuses(ctx context.Context) ([]domain.FeederStatus, error) {
	query := `
		SELECT id, feeder_asset_id, feeder_code, current_load_mw, current_amps, state, outage_reason_code, outage_narrative, priority_rank, tariff_band, active_ptw_id, scheduled_restoration_time, updated_at
		FROM energy_management.feeder_status_ledger
		ORDER BY priority_rank ASC, feeder_code ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListFeederStatuses failed: %w", err)
	}
	defer rows.Close()

	var feeders []domain.FeederStatus
	for rows.Next() {
		var f domain.FeederStatus
		var reasonCode, narrative, ptwID sql.NullString
		var restoTime sql.NullTime
		var band string
		if err := rows.Scan(
			&f.ID, &f.FeederAssetID, &f.FeederCode,
			&f.CurrentLoadMW, &f.CurrentAmps, &f.State,
			&reasonCode, &narrative, &f.PriorityRank, &band,
			&ptwID, &restoTime, &f.UpdatedAt); err != nil {
			return nil, err
		}
		f.TariffBand = s01domain.TariffBand(band)
		if reasonCode.Valid {
			f.OutageReasonCode = &reasonCode.String
		}
		if narrative.Valid {
			f.OutageNarrative = &narrative.String
		}
		if ptwID.Valid {
			f.ActivePTWID = &ptwID.String
		}
		if restoTime.Valid {
			f.ScheduledRestorationTime = &restoTime.Time
		}
		feeders = append(feeders, f)
	}
	return feeders, rows.Err()
}

