// Package repository implements the PostgreSQL storage engine for S05 Generation Department.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
)

// CreatePlantJob inserts a plant fault or workshop job into plant_fault_and_workshop_jobs.
func (r *PostgresRepository) CreatePlantJob(ctx context.Context, job domain.PlantFaultAndWorkshopJob) error {
	query := `
		INSERT INTO generation_master.plant_fault_and_workshop_jobs
		(id, job_no, plant_id, reported_by_operator_id, wing, fault_description, severity, unit_head_id, gm_approved, gm_approval_notes, requires_workshop, workshop_status, raw_materials_summary, store_requisition_ref, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.JobNo, job.PlantID, job.ReportedByOperatorID,
		job.Wing, job.FaultDescription, job.Severity, job.UnitHeadID,
		job.GMApproved, job.GMApprovalNotes, job.RequiresWorkshop,
		job.WorkshopStatus, job.RawMaterialsSummary, job.StoreRequisitionRef,
		job.Status, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreatePlantJob failed: %w", err)
	}
	return nil
}

// GetPlantJobByID retrieves a fault/workshop job by UUID.
func (r *PostgresRepository) GetPlantJobByID(ctx context.Context, id string) (*domain.PlantFaultAndWorkshopJob, error) {
	query := `
		SELECT id, job_no, plant_id, reported_by_operator_id, wing, fault_description, severity, unit_head_id, gm_approved, gm_approval_notes, requires_workshop, workshop_status, raw_materials_summary, store_requisition_ref, status, created_at, updated_at
		FROM generation_master.plant_fault_and_workshop_jobs WHERE id = $1`

	var j domain.PlantFaultAndWorkshopJob
	var unitHead, gmNotes, rawMat, storeRef sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&j.ID, &j.JobNo, &j.PlantID, &j.ReportedByOperatorID,
		&j.Wing, &j.FaultDescription, &j.Severity, &unitHead,
		&j.GMApproved, &gmNotes, &j.RequiresWorkshop,
		&j.WorkshopStatus, &rawMat, &storeRef,
		&j.Status, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetPlantJobByID failed: %w", err)
	}
	if unitHead.Valid {
		j.UnitHeadID = &unitHead.String
	}
	if gmNotes.Valid {
		j.GMApprovalNotes = &gmNotes.String
	}
	if rawMat.Valid {
		j.RawMaterialsSummary = &rawMat.String
	}
	if storeRef.Valid {
		j.StoreRequisitionRef = &storeRef.String
	}
	return &j, nil
}

// UpdatePlantJobGMApproval updates GM approval decision, notes, and status.
func (r *PostgresRepository) UpdatePlantJobGMApproval(ctx context.Context, id string, approved bool, notes string) error {
	newStatus := string(domain.JobStatusGMApproved)
	if !approved {
		newStatus = string(domain.JobStatusCancelled)
	}

	query := `
		UPDATE generation_master.plant_fault_and_workshop_jobs
		SET gm_approved = $2, gm_approval_notes = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id, approved, notes, newStatus)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdatePlantJobGMApproval failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateWorkshopStatus advances workshop fabrication status.
func (r *PostgresRepository) UpdateWorkshopStatus(ctx context.Context, id string, status domain.WorkshopStatus, materialsRef *string) error {
	var jobStatus string
	if status == domain.WorkshopCompletedDelivered {
		jobStatus = string(domain.JobStatusResolvedReinstalled)
	} else if status == domain.WorkshopMachiningInProgress || status == domain.WorkshopWindingFabrication {
		jobStatus = string(domain.JobStatusWorkshopInProgress)
	}

	query := `
		UPDATE generation_master.plant_fault_and_workshop_jobs
		SET workshop_status = $2,
		    store_requisition_ref = COALESCE($3, store_requisition_ref),
		    status = CASE WHEN $4 != '' THEN $4 ELSE status END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, string(status), materialsRef, jobStatus)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdateWorkshopStatus failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListPlantJobs queries jobs by optional plant, wing, and status filters.
func (r *PostgresRepository) ListPlantJobs(ctx context.Context, plantID string, wing domain.TechnicalWing, status domain.PlantJobStatus) ([]domain.PlantFaultAndWorkshopJob, error) {
	query := `
		SELECT id, job_no, plant_id, reported_by_operator_id, wing, fault_description, severity, unit_head_id, gm_approved, gm_approval_notes, requires_workshop, workshop_status, raw_materials_summary, store_requisition_ref, status, created_at, updated_at
		FROM generation_master.plant_fault_and_workshop_jobs
		WHERE ($1 = '' OR plant_id = $1::uuid)
		  AND ($2 = '' OR wing = $2)
		  AND ($3 = '' OR status = $3)
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, plantID, string(wing), string(status))
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListPlantJobs failed: %w", err)
	}
	defer rows.Close()

	var jobs []domain.PlantFaultAndWorkshopJob
	for rows.Next() {
		var j domain.PlantFaultAndWorkshopJob
		var unitHead, gmNotes, rawMat, storeRef sql.NullString
		if err := rows.Scan(
			&j.ID, &j.JobNo, &j.PlantID, &j.ReportedByOperatorID,
			&j.Wing, &j.FaultDescription, &j.Severity, &unitHead,
			&j.GMApproved, &gmNotes, &j.RequiresWorkshop,
			&j.WorkshopStatus, &rawMat, &storeRef,
			&j.Status, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		if unitHead.Valid {
			j.UnitHeadID = &unitHead.String
		}
		if gmNotes.Valid {
			j.GMApprovalNotes = &gmNotes.String
		}
		if rawMat.Valid {
			j.RawMaterialsSummary = &rawMat.String
		}
		if storeRef.Valid {
			j.StoreRequisitionRef = &storeRef.String
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

