// Package repository implements the PostgreSQL storage engine for S05 Generation Department.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
)

// PostgresRepository provides production PostgreSQL queries for the generation_master schema.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new Postgres repository handler.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateTelemetryLog records a generation telemetry snapshot into master_station_telemetry_logs.
func (r *PostgresRepository) CreateTelemetryLog(ctx context.Context, log domain.MasterStationTelemetryLog) error {
	query := `
		INSERT INTO generation_master.master_station_telemetry_logs
		(id, shift_id, shift_window, plant_id, plant_state, outage_reason, outage_details, water_level_m, headwater_elevation_m, tailwater_elevation_m, inflow_rate_m3_s, mw_generated, mvar_reactive, stator_winding_temp_c, bearing_vibration_mm_s, lube_oil_pressure_bar, logged_by_operator_id, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`

	var outageReason *string
	if log.OutageReason != nil {
		str := string(*log.OutageReason)
		outageReason = &str
	}

	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.ShiftID, log.ShiftWindow, log.PlantID, log.PlantState,
		outageReason, log.OutageDetails, log.WaterLevelM,
		log.HeadwaterElevationM, log.TailwaterElevationM, log.InflowRateM3S,
		log.MWGenerated, log.MVArReactive, log.StatorWindingTempC,
		log.BearingVibrationMMS, log.LubeOilPressureBar,
		log.LoggedByOperatorID, log.RecordedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateTelemetryLog failed: %w", err)
	}
	return nil
}

// GetLatestTelemetryByPlant retrieves the most recent telemetry entry for a given generating plant.
func (r *PostgresRepository) GetLatestTelemetryByPlant(ctx context.Context, plantID string) (*domain.MasterStationTelemetryLog, error) {
	query := `
		SELECT id, shift_id, shift_window, plant_id, plant_state, outage_reason, outage_details, water_level_m, headwater_elevation_m, tailwater_elevation_m, inflow_rate_m3_s, mw_generated, mvar_reactive, stator_winding_temp_c, bearing_vibration_mm_s, lube_oil_pressure_bar, logged_by_operator_id, recorded_at
		FROM generation_master.master_station_telemetry_logs
		WHERE plant_id = $1::uuid
		ORDER BY recorded_at DESC LIMIT 1`

	var l domain.MasterStationTelemetryLog
	var outageReason, outageDetails sql.NullString
	err := r.db.QueryRowContext(ctx, query, plantID).Scan(
		&l.ID, &l.ShiftID, &l.ShiftWindow, &l.PlantID, &l.PlantState,
		&outageReason, &outageDetails, &l.WaterLevelM,
		&l.HeadwaterElevationM, &l.TailwaterElevationM, &l.InflowRateM3S,
		&l.MWGenerated, &l.MVArReactive, &l.StatorWindingTempC,
		&l.BearingVibrationMMS, &l.LubeOilPressureBar,
		&l.LoggedByOperatorID, &l.RecordedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetLatestTelemetryByPlant failed: %w", err)
	}
	if outageReason.Valid {
		reason := domain.OutageReason(outageReason.String)
		l.OutageReason = &reason
	}
	if outageDetails.Valid {
		l.OutageDetails = &outageDetails.String
	}
	return &l, nil
}

// ListTelemetryLogs returns historical telemetry records by shift or plant filter.
func (r *PostgresRepository) ListTelemetryLogs(ctx context.Context, shiftID, plantID string, limit int) ([]domain.MasterStationTelemetryLog, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, shift_id, shift_window, plant_id, plant_state, outage_reason, outage_details, water_level_m, headwater_elevation_m, tailwater_elevation_m, inflow_rate_m3_s, mw_generated, mvar_reactive, stator_winding_temp_c, bearing_vibration_mm_s, lube_oil_pressure_bar, logged_by_operator_id, recorded_at
		FROM generation_master.master_station_telemetry_logs
		WHERE ($1 = '' OR shift_id = $1)
		  AND ($2 = '' OR plant_id = $2::uuid)
		ORDER BY recorded_at DESC LIMIT $3`

	rows, err := r.db.QueryContext(ctx, query, shiftID, plantID, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListTelemetryLogs failed: %w", err)
	}
	defer rows.Close()

	var logs []domain.MasterStationTelemetryLog
	for rows.Next() {
		var l domain.MasterStationTelemetryLog
		var outageReason, outageDetails sql.NullString
		if err := rows.Scan(
			&l.ID, &l.ShiftID, &l.ShiftWindow, &l.PlantID, &l.PlantState,
			&outageReason, &outageDetails, &l.WaterLevelM,
			&l.HeadwaterElevationM, &l.TailwaterElevationM, &l.InflowRateM3S,
			&l.MWGenerated, &l.MVArReactive, &l.StatorWindingTempC,
			&l.BearingVibrationMMS, &l.LubeOilPressureBar,
			&l.LoggedByOperatorID, &l.RecordedAt); err != nil {
			return nil, err
		}
		if outageReason.Valid {
			reason := domain.OutageReason(outageReason.String)
			l.OutageReason = &reason
		}
		if outageDetails.Valid {
			l.OutageDetails = &outageDetails.String
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// CreateRollcall records shift personnel on duty at a remote plant.
func (r *PostgresRepository) CreateRollcall(ctx context.Context, rc domain.RemotePlantShiftRollcall) error {
	query := `
		INSERT INTO generation_master.remote_plant_shift_rollcall
		(id, shift_id, shift_window, plant_id, operator_names, protection_engineers, technicians, handover_notes, synced_to_hq_hr, master_operator_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		rc.ID, rc.ShiftID, rc.ShiftWindow, rc.PlantID,
		rc.OperatorNames, rc.ProtectionEngineers, rc.Technicians,
		rc.HandoverNotes, rc.SyncedToHQHR, rc.MasterOperatorID, rc.CreatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateRollcall failed: %w", err)
	}
	return nil
}

// ListRollcallsByShift returns roll-call rosters across all 6 plants for a shift window.
func (r *PostgresRepository) ListRollcallsByShift(ctx context.Context, shiftID string) ([]domain.RemotePlantShiftRollcall, error) {
	query := `
		SELECT id, shift_id, shift_window, plant_id, operator_names, protection_engineers, technicians, handover_notes, synced_to_hq_hr, master_operator_id, created_at
		FROM generation_master.remote_plant_shift_rollcall
		WHERE shift_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, shiftID)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListRollcallsByShift failed: %w", err)
	}
	defer rows.Close()

	var results []domain.RemotePlantShiftRollcall
	for rows.Next() {
		var rc domain.RemotePlantShiftRollcall
		if err := rows.Scan(
			&rc.ID, &rc.ShiftID, &rc.ShiftWindow, &rc.PlantID,
			&rc.OperatorNames, &rc.ProtectionEngineers, &rc.Technicians,
			&rc.HandoverNotes, &rc.SyncedToHQHR, &rc.MasterOperatorID, &rc.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, rc)
	}
	return results, rows.Err()
}

