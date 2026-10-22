// Package repository implements PostgreSQL queries for staff allowance claims.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
)

// CreateClaim persists a new staff financial allowance claim.
func (p *PostgresRepository) CreateClaim(ctx context.Context, claim domain.StaffAllowanceClaim) error {
	query := `
		INSERT INTO stores_audit.staff_allowance_claims (
			id, claim_no, staff_id, claim_type, amount, claim_date,
			shift_window, linked_ptw_or_ticket_id, roster_match_verified,
			ptw_crew_verified, gps_match_verified, ollama_audit_flag,
			ollama_fraud_score, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	now := time.Now().UTC()
	_, err := p.db.ExecContext(ctx, query,
		claim.ID, claim.ClaimNo, claim.StaffID, string(claim.ClaimType), claim.Amount,
		claim.ClaimDate, claim.ShiftWindow, claim.LinkedPTWOrTicketID, claim.RosterMatchVerified,
		claim.PTWCrewVerified, claim.GPSMatchVerified, string(claim.OllamaAuditFlag),
		claim.OllamaFraudScore, string(claim.Status), now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to create claim: %w", err)
	}
	return nil
}

// GetClaimByID retrieves a financial claim by UUID.
func (p *PostgresRepository) GetClaimByID(ctx context.Context, id string) (*domain.StaffAllowanceClaim, error) {
	query := `
		SELECT id, claim_no, staff_id, claim_type, amount, claim_date,
		       shift_window, linked_ptw_or_ticket_id, roster_match_verified,
		       ptw_crew_verified, gps_match_verified, ollama_audit_flag,
		       ollama_fraud_score, auditor_id, auditor_notes, audit_decision_at,
		       status, created_at, updated_at
		FROM stores_audit.staff_allowance_claims WHERE id = $1`

	var c domain.StaffAllowanceClaim
	var cType, flag, stat string
	var shift, linked, auditorID, auditorNotes sql.NullString
	var decisionAt sql.NullTime

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.ClaimNo, &c.StaffID, &cType, &c.Amount, &c.ClaimDate,
		&shift, &linked, &c.RosterMatchVerified,
		&c.PTWCrewVerified, &c.GPSMatchVerified, &flag,
		&c.OllamaFraudScore, &auditorID, &auditorNotes, &decisionAt,
		&stat, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres: failed to get claim by id: %w", err)
	}

	c.ClaimType = domain.ClaimType(cType)
	c.OllamaAuditFlag = domain.AIAuditFlag(flag)
	c.Status = domain.ClaimStatus(stat)

	if shift.Valid {
		c.ShiftWindow = &shift.String
	}
	if linked.Valid {
		c.LinkedPTWOrTicketID = &linked.String
	}
	if auditorID.Valid {
		c.AuditorID = &auditorID.String
	}
	if auditorNotes.Valid {
		c.AuditorNotes = &auditorNotes.String
	}
	if decisionAt.Valid {
		c.AuditDecisionAt = &decisionAt.Time
	}
	return &c, nil
}

// UpdateClaim updates audit flags and decision notes.
func (p *PostgresRepository) UpdateClaim(ctx context.Context, claim domain.StaffAllowanceClaim) error {
	query := `
		UPDATE stores_audit.staff_allowance_claims
		SET roster_match_verified = $1, ptw_crew_verified = $2, gps_match_verified = $3,
		    ollama_audit_flag = $4, ollama_fraud_score = $5, auditor_id = $6,
		    auditor_notes = $7, audit_decision_at = $8, status = $9, updated_at = $10
		WHERE id = $11`

	res, err := p.db.ExecContext(ctx, query,
		claim.RosterMatchVerified, claim.PTWCrewVerified, claim.GPSMatchVerified,
		string(claim.OllamaAuditFlag), claim.OllamaFraudScore, claim.AuditorID,
		claim.AuditorNotes, claim.AuditDecisionAt, string(claim.Status), time.Now().UTC(), claim.ID,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to update claim: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListClaims returns allowance claims filtered by staff and status.
func (p *PostgresRepository) ListClaims(ctx context.Context, staffID string, status domain.ClaimStatus) ([]domain.StaffAllowanceClaim, error) {
	query := `
		SELECT id, claim_no, staff_id, claim_type, amount, claim_date,
		       shift_window, linked_ptw_or_ticket_id, roster_match_verified,
		       ptw_crew_verified, gps_match_verified, ollama_audit_flag,
		       ollama_fraud_score, auditor_id, auditor_notes, audit_decision_at,
		       status, created_at, updated_at
		FROM stores_audit.staff_allowance_claims
		WHERE ($1 = '' OR staff_id = $1)
		  AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC`

	rows, err := p.db.QueryContext(ctx, query, staffID, string(status))
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to list claims: %w", err)
	}
	defer rows.Close()

	var results []domain.StaffAllowanceClaim
	for rows.Next() {
		var c domain.StaffAllowanceClaim
		var cType, flag, stat string
		var shift, linked, auditorID, auditorNotes sql.NullString
		var decisionAt sql.NullTime

		if err := rows.Scan(
			&c.ID, &c.ClaimNo, &c.StaffID, &cType, &c.Amount, &c.ClaimDate,
			&shift, &linked, &c.RosterMatchVerified,
			&c.PTWCrewVerified, &c.GPSMatchVerified, &flag,
			&c.OllamaFraudScore, &auditorID, &auditorNotes, &decisionAt,
			&stat, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		c.ClaimType = domain.ClaimType(cType)
		c.OllamaAuditFlag = domain.AIAuditFlag(flag)
		c.Status = domain.ClaimStatus(stat)
		if shift.Valid {
			c.ShiftWindow = &shift.String
		}
		if linked.Valid {
			c.LinkedPTWOrTicketID = &linked.String
		}
		if auditorID.Valid {
			c.AuditorID = &auditorID.String
		}
		if auditorNotes.Valid {
			c.AuditorNotes = &auditorNotes.String
		}
		if decisionAt.Valid {
			c.AuditDecisionAt = &decisionAt.Time
		}
		results = append(results, c)
	}
	return results, nil
}

