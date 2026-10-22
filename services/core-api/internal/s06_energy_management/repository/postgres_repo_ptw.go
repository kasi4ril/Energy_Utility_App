// Package repository implements the PostgreSQL storage engine for S06 Energy Management.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
)

// CreatePTW inserts a new safety permit into safety_permits_ptw.
func (r *PostgresRepository) CreatePTW(ctx context.Context, ptw domain.SafetyPermitPTW) error {
	query := `
		INSERT INTO energy_management.safety_permits_ptw
		(id, permit_no, work_order_id, asset_id, unit_head_id, dept_head_id, dept_head_approved, dept_head_approved_at, safety_rep_id, safety_rep_approved, safety_rep_approved_at, substation_operator_id, isolation_granted_at, isolation_notes, safety_clearance_signed_at, safety_clearance_notes, restored_at, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	_, err := r.db.ExecContext(ctx, query,
		ptw.ID, ptw.PermitNo, ptw.WorkOrderID, ptw.AssetID, ptw.UnitHeadID,
		ptw.DeptHeadID, ptw.DeptHeadApproved, ptw.DeptHeadApprovedAt,
		ptw.SafetyRepID, ptw.SafetyRepApproved, ptw.SafetyRepApprovedAt,
		ptw.SubstationOperatorID, ptw.IsolationGrantedAt, ptw.IsolationNotes,
		ptw.SafetyClearanceSignedAt, ptw.SafetyClearanceNotes,
		ptw.RestoredAt, string(ptw.Status), ptw.CreatedAt, ptw.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreatePTW failed: %w", err)
	}
	return nil
}

// GetPTWByID retrieves a safety permit by UUID.
func (r *PostgresRepository) GetPTWByID(ctx context.Context, id string) (*domain.SafetyPermitPTW, error) {
	query := `
		SELECT id, permit_no, work_order_id, asset_id, unit_head_id, dept_head_id, dept_head_approved, dept_head_approved_at, safety_rep_id, safety_rep_approved, safety_rep_approved_at, substation_operator_id, isolation_granted_at, isolation_notes, safety_clearance_signed_at, safety_clearance_notes, restored_at, status, created_at, updated_at
		FROM energy_management.safety_permits_ptw WHERE id = $1`

	return r.scanPTWRow(r.db.QueryRowContext(ctx, query, id))
}

// GetPTWByPermitNo retrieves a safety permit by permit number.
func (r *PostgresRepository) GetPTWByPermitNo(ctx context.Context, permitNo string) (*domain.SafetyPermitPTW, error) {
	query := `
		SELECT id, permit_no, work_order_id, asset_id, unit_head_id, dept_head_id, dept_head_approved, dept_head_approved_at, safety_rep_id, safety_rep_approved, safety_rep_approved_at, substation_operator_id, isolation_granted_at, isolation_notes, safety_clearance_signed_at, safety_clearance_notes, restored_at, status, created_at, updated_at
		FROM energy_management.safety_permits_ptw WHERE permit_no = $1`

	return r.scanPTWRow(r.db.QueryRowContext(ctx, query, permitNo))
}

// UpdatePTWStatus updates all approval, isolation, clearance, and restoration fields of a permit.
func (r *PostgresRepository) UpdatePTWStatus(ctx context.Context, ptw domain.SafetyPermitPTW) error {
	query := `
		UPDATE energy_management.safety_permits_ptw
		SET dept_head_id = $2, dept_head_approved = $3, dept_head_approved_at = $4,
		    safety_rep_id = $5, safety_rep_approved = $6, safety_rep_approved_at = $7,
		    substation_operator_id = $8, isolation_granted_at = $9, isolation_notes = $10,
		    safety_clearance_signed_at = $11, safety_clearance_notes = $12, restored_at = $13,
		    status = $14, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query,
		ptw.ID, ptw.DeptHeadID, ptw.DeptHeadApproved, ptw.DeptHeadApprovedAt,
		ptw.SafetyRepID, ptw.SafetyRepApproved, ptw.SafetyRepApprovedAt,
		ptw.SubstationOperatorID, ptw.IsolationGrantedAt, ptw.IsolationNotes,
		ptw.SafetyClearanceSignedAt, ptw.SafetyClearanceNotes, ptw.RestoredAt,
		string(ptw.Status))
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdatePTWStatus failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListPTWs returns safety permits filtered by status and asset.
func (r *PostgresRepository) ListPTWs(ctx context.Context, status domain.PTWStatus, assetID string) ([]domain.SafetyPermitPTW, error) {
	query := `
		SELECT id, permit_no, work_order_id, asset_id, unit_head_id, dept_head_id, dept_head_approved, dept_head_approved_at, safety_rep_id, safety_rep_approved, safety_rep_approved_at, substation_operator_id, isolation_granted_at, isolation_notes, safety_clearance_signed_at, safety_clearance_notes, restored_at, status, created_at, updated_at
		FROM energy_management.safety_permits_ptw
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR asset_id = $2::uuid)
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, string(status), assetID)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListPTWs failed: %w", err)
	}
	defer rows.Close()

	var permits []domain.SafetyPermitPTW
	for rows.Next() {
		var p domain.SafetyPermitPTW
		var woID, deptHead, safetyRep, subOp, isoNotes, clearNotes sql.NullString
		var dhAppAt, srAppAt, isoAt, clearAt, restoAt sql.NullTime
		if err := rows.Scan(
			&p.ID, &p.PermitNo, &woID, &p.AssetID, &p.UnitHeadID,
			&deptHead, &p.DeptHeadApproved, &dhAppAt,
			&safetyRep, &p.SafetyRepApproved, &srAppAt,
			&subOp, &isoAt, &isoNotes,
			&clearAt, &clearNotes,
			&restoAt, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		r.hydratePTWNulls(&p, woID, deptHead, safetyRep, subOp, isoNotes, clearNotes, dhAppAt, srAppAt, isoAt, clearAt, restoAt)
		permits = append(permits, p)
	}
	return permits, rows.Err()
}

func (r *PostgresRepository) scanPTWRow(row *sql.Row) (*domain.SafetyPermitPTW, error) {
	var p domain.SafetyPermitPTW
	var woID, deptHead, safetyRep, subOp, isoNotes, clearNotes sql.NullString
	var dhAppAt, srAppAt, isoAt, clearAt, restoAt sql.NullTime
	err := row.Scan(
		&p.ID, &p.PermitNo, &woID, &p.AssetID, &p.UnitHeadID,
		&deptHead, &p.DeptHeadApproved, &dhAppAt,
		&safetyRep, &p.SafetyRepApproved, &srAppAt,
		&subOp, &isoAt, &isoNotes,
		&clearAt, &clearNotes,
		&restoAt, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: scanPTW failed: %w", err)
	}
	r.hydratePTWNulls(&p, woID, deptHead, safetyRep, subOp, isoNotes, clearNotes, dhAppAt, srAppAt, isoAt, clearAt, restoAt)
	return &p, nil
}

func (r *PostgresRepository) hydratePTWNulls(
	p *domain.SafetyPermitPTW,
	woID, deptHead, safetyRep, subOp, isoNotes, clearNotes sql.NullString,
	dhAppAt, srAppAt, isoAt, clearAt, restoAt sql.NullTime,
) {
	if woID.Valid {
		p.WorkOrderID = &woID.String
	}
	if deptHead.Valid {
		p.DeptHeadID = &deptHead.String
	}
	if dhAppAt.Valid {
		p.DeptHeadApprovedAt = &dhAppAt.Time
	}
	if safetyRep.Valid {
		p.SafetyRepID = &safetyRep.String
	}
	if srAppAt.Valid {
		p.SafetyRepApprovedAt = &srAppAt.Time
	}
	if subOp.Valid {
		p.SubstationOperatorID = &subOp.String
	}
	if isoAt.Valid {
		p.IsolationGrantedAt = &isoAt.Time
	}
	if isoNotes.Valid {
		p.IsolationNotes = &isoNotes.String
	}
	if clearAt.Valid {
		p.SafetyClearanceSignedAt = &clearAt.Time
	}
	if clearNotes.Valid {
		p.SafetyClearanceNotes = &clearNotes.String
	}
	if restoAt.Valid {
		p.RestoredAt = &restoAt.Time
	}
}

