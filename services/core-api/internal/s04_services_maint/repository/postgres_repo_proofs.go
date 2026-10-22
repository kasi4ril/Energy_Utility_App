// Package repository implements the PostgreSQL + PostGIS storage engine for S04 Services & Maintenance.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
)

// CreateProofOfWork inserts a proof of work log with spatial PostGIS point.
func (r *PostgresRepository) CreateProofOfWork(ctx context.Context, proof domain.ProofOfWorkLog) error {
	query := `
		INSERT INTO services_maint.proof_of_work_logs
		(id, work_order_id, team_lead_id, captured_gps, distance_from_asset_meters, before_photo_url, after_photo_url, photo_hash_before, photo_hash_after, meter_seal_number, work_details, ollama_gps_verified, ai_risk_score, ai_reasoning, unit_head_verified, unit_head_notes, submitted_at, verified_at)
		VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($5, $4), 4326), $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	_, err := r.db.ExecContext(ctx, query,
		proof.ID, proof.WorkOrderID, proof.TeamLeadID,
		proof.CapturedGPS.Latitude, proof.CapturedGPS.Longitude,
		proof.DistanceFromAssetMeters, proof.BeforePhotoURL, proof.AfterPhotoURL,
		proof.PhotoHashBefore, proof.PhotoHashAfter, proof.MeterSealNumber,
		proof.WorkDetails, proof.OllamaGPSVerified, proof.AIRiskScore, proof.AIReasoning,
		proof.UnitHeadVerified, proof.UnitHeadNotes, proof.SubmittedAt, proof.VerifiedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateProofOfWork failed: %w", err)
	}
	return nil
}

// GetProofByWorkOrderID retrieves the proof log for a given work order.
func (r *PostgresRepository) GetProofByWorkOrderID(ctx context.Context, workOrderID string) (*domain.ProofOfWorkLog, error) {
	query := `
		SELECT id, work_order_id, team_lead_id, ST_Y(captured_gps), ST_X(captured_gps), distance_from_asset_meters, before_photo_url, after_photo_url, photo_hash_before, photo_hash_after, meter_seal_number, work_details, ollama_gps_verified, ai_risk_score, ai_reasoning, unit_head_verified, unit_head_notes, submitted_at, verified_at
		FROM services_maint.proof_of_work_logs WHERE work_order_id = $1`

	var p domain.ProofOfWorkLog
	var seal, headNotes sql.NullString
	var verifiedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, workOrderID).Scan(
		&p.ID, &p.WorkOrderID, &p.TeamLeadID,
		&p.CapturedGPS.Latitude, &p.CapturedGPS.Longitude,
		&p.DistanceFromAssetMeters, &p.BeforePhotoURL, &p.AfterPhotoURL,
		&p.PhotoHashBefore, &p.PhotoHashAfter, &seal,
		&p.WorkDetails, &p.OllamaGPSVerified, &p.AIRiskScore, &p.AIReasoning,
		&p.UnitHeadVerified, &headNotes, &p.SubmittedAt, &verifiedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetProofByWorkOrderID failed: %w", err)
	}
	if seal.Valid {
		p.MeterSealNumber = &seal.String
	}
	if headNotes.Valid {
		p.UnitHeadNotes = &headNotes.String
	}
	if verifiedAt.Valid {
		p.VerifiedAt = &verifiedAt.Time
	}
	return &p, nil
}

// UpdateProofVerification updates the Unit Head verification flag, notes, and timestamp.
func (r *PostgresRepository) UpdateProofVerification(ctx context.Context, proofID string, verified bool, notes string) error {
	query := `
		UPDATE services_maint.proof_of_work_logs
		SET unit_head_verified = $2, unit_head_notes = $3, verified_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, proofID, verified, notes)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdateProofVerification failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListProofsByTeamLead returns proofs submitted by a team lead.
func (r *PostgresRepository) ListProofsByTeamLead(ctx context.Context, teamLeadID string, limit int) ([]domain.ProofOfWorkLog, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, work_order_id, team_lead_id, ST_Y(captured_gps), ST_X(captured_gps), distance_from_asset_meters, before_photo_url, after_photo_url, photo_hash_before, photo_hash_after, meter_seal_number, work_details, ollama_gps_verified, ai_risk_score, ai_reasoning, unit_head_verified, unit_head_notes, submitted_at, verified_at
		FROM services_maint.proof_of_work_logs
		WHERE ($1 = '' OR team_lead_id = $1::uuid)
		ORDER BY submitted_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, teamLeadID, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListProofsByTeamLead failed: %w", err)
	}
	defer rows.Close()

	var proofs []domain.ProofOfWorkLog
	for rows.Next() {
		var p domain.ProofOfWorkLog
		var seal, headNotes sql.NullString
		var verifiedAt sql.NullTime
		if err := rows.Scan(
			&p.ID, &p.WorkOrderID, &p.TeamLeadID,
			&p.CapturedGPS.Latitude, &p.CapturedGPS.Longitude,
			&p.DistanceFromAssetMeters, &p.BeforePhotoURL, &p.AfterPhotoURL,
			&p.PhotoHashBefore, &p.PhotoHashAfter, &seal,
			&p.WorkDetails, &p.OllamaGPSVerified, &p.AIRiskScore, &p.AIReasoning,
			&p.UnitHeadVerified, &headNotes, &p.SubmittedAt, &verifiedAt); err != nil {
			return nil, err
		}
		if seal.Valid {
			p.MeterSealNumber = &seal.String
		}
		if headNotes.Valid {
			p.UnitHeadNotes = &headNotes.String
		}
		if verifiedAt.Valid {
			p.VerifiedAt = &verifiedAt.Time
		}
		proofs = append(proofs, p)
	}
	return proofs, rows.Err()
}

