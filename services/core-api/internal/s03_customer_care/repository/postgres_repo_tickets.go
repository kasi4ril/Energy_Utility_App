// Package repository implements the PostgreSQL + PostGIS storage engine for S03 Customer Care.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
)

// CreateTicket records a service ticket in customercare.service_tickets.
func (r *PostgresRepository) CreateTicket(ctx context.Context, t domain.ServiceTicket) error {
	query := `
		INSERT INTO customercare.service_tickets
		(id, ticket_no, customer_id, channel, category, description, assigned_unit_id, unit_head_id, dept_head_supervisor_id, status, sla_due_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.TicketNo, t.CustomerID, t.Channel, t.Category, t.Description,
		t.AssignedUnitID, t.UnitHeadID, t.DeptHeadSupervisorID, t.Status, t.SLADueAt, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateTicket failed: %w", err)
	}
	return nil
}

// GetTicketByID fetches ticket by UUID.
func (r *PostgresRepository) GetTicketByID(ctx context.Context, id string) (*domain.ServiceTicket, error) {
	query := `
		SELECT id, ticket_no, customer_id, channel, category, description, assigned_unit_id, unit_head_id, dept_head_supervisor_id, status, attended_proof_note, attended_by_staff_id, customer_call_verified, verified_by_cc_officer_id, customer_feedback_notes, sla_due_at, created_at, updated_at
		FROM customercare.service_tickets WHERE id = $1`

	var t domain.ServiceTicket
	var custID, unitHead, deptHead, attendedNote, attendedStaff, ccOfficer, feedback sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.TicketNo, &custID, &t.Channel, &t.Category, &t.Description,
		&t.AssignedUnitID, &unitHead, &deptHead, &t.Status, &attendedNote, &attendedStaff,
		&t.CustomerCallVerified, &ccOfficer, &feedback, &t.SLADueAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetTicketByID failed: %w", err)
	}
	if custID.Valid {
		t.CustomerID = &custID.String
	}
	if unitHead.Valid {
		t.UnitHeadID = &unitHead.String
	}
	if deptHead.Valid {
		t.DeptHeadSupervisorID = &deptHead.String
	}
	if attendedNote.Valid {
		t.AttendedProofNote = attendedNote.String
	}
	if attendedStaff.Valid {
		t.AttendedByStaffID = &attendedStaff.String
	}
	if ccOfficer.Valid {
		t.VerifiedByCCOfficerID = &ccOfficer.String
	}
	if feedback.Valid {
		t.CustomerFeedbackNotes = feedback.String
	}
	return &t, nil
}

// GetTicketByNo fetches ticket by ticket number.
func (r *PostgresRepository) GetTicketByNo(ctx context.Context, ticketNo string) (*domain.ServiceTicket, error) {
	query := `
		SELECT id, ticket_no, customer_id, channel, category, description, assigned_unit_id, unit_head_id, dept_head_supervisor_id, status, attended_proof_note, attended_by_staff_id, customer_call_verified, verified_by_cc_officer_id, customer_feedback_notes, sla_due_at, created_at, updated_at
		FROM customercare.service_tickets WHERE ticket_no = $1`

	var t domain.ServiceTicket
	var custID, unitHead, deptHead, attendedNote, attendedStaff, ccOfficer, feedback sql.NullString
	err := r.db.QueryRowContext(ctx, query, ticketNo).Scan(
		&t.ID, &t.TicketNo, &custID, &t.Channel, &t.Category, &t.Description,
		&t.AssignedUnitID, &unitHead, &deptHead, &t.Status, &attendedNote, &attendedStaff,
		&t.CustomerCallVerified, &ccOfficer, &feedback, &t.SLADueAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetTicketByNo failed: %w", err)
	}
	if custID.Valid {
		t.CustomerID = &custID.String
	}
	if unitHead.Valid {
		t.UnitHeadID = &unitHead.String
	}
	if deptHead.Valid {
		t.DeptHeadSupervisorID = &deptHead.String
	}
	if attendedNote.Valid {
		t.AttendedProofNote = attendedNote.String
	}
	if attendedStaff.Valid {
		t.AttendedByStaffID = &attendedStaff.String
	}
	if ccOfficer.Valid {
		t.VerifiedByCCOfficerID = &ccOfficer.String
	}
	if feedback.Valid {
		t.CustomerFeedbackNotes = feedback.String
	}
	return &t, nil
}

// UpdateTicketAttended updates technical repair proof and marks ticket ATTENDED_TO.
func (r *PostgresRepository) UpdateTicketAttended(ctx context.Context, id string, proofNote string, staffID string) error {
	query := `
		UPDATE customercare.service_tickets
		SET status = 'ATTENDED_TO', attended_proof_note = $2, attended_by_staff_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id, proofNote, staffID)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdateTicketAttended failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateTicketClosure executes customer confirmation verification and final status transition.
func (r *PostgresRepository) UpdateTicketClosure(ctx context.Context, id string, status domain.TicketStatus, callVerified bool, officerID string, feedback string) error {
	query := `
		UPDATE customercare.service_tickets
		SET status = $2, customer_call_verified = $3, verified_by_cc_officer_id = $4, customer_feedback_notes = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id, string(status), callVerified, officerID, feedback)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdateTicketClosure failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListTickets returns tickets with optional status, unit, and category filters.
func (r *PostgresRepository) ListTickets(ctx context.Context, status domain.TicketStatus, unitID string, cat domain.TicketCategory) ([]domain.ServiceTicket, error) {
	query := `
		SELECT id, ticket_no, customer_id, channel, category, description, assigned_unit_id, unit_head_id, dept_head_supervisor_id, status, attended_proof_note, attended_by_staff_id, customer_call_verified, verified_by_cc_officer_id, customer_feedback_notes, sla_due_at, created_at, updated_at
		FROM customercare.service_tickets
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR assigned_unit_id = $2::uuid)
		  AND ($3 = '' OR category = $3)
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, string(status), unitID, string(cat))
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListTickets failed: %w", err)
	}
	defer rows.Close()

	var tickets []domain.ServiceTicket
	for rows.Next() {
		var t domain.ServiceTicket
		var custID, unitHead, deptHead, attendedNote, attendedStaff, ccOfficer, feedback sql.NullString
		if err := rows.Scan(&t.ID, &t.TicketNo, &custID, &t.Channel, &t.Category, &t.Description,
			&t.AssignedUnitID, &unitHead, &deptHead, &t.Status, &attendedNote, &attendedStaff,
			&t.CustomerCallVerified, &ccOfficer, &feedback, &t.SLADueAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if custID.Valid {
			t.CustomerID = &custID.String
		}
		if unitHead.Valid {
			t.UnitHeadID = &unitHead.String
		}
		if deptHead.Valid {
			t.DeptHeadSupervisorID = &deptHead.String
		}
		if attendedNote.Valid {
			t.AttendedProofNote = attendedNote.String
		}
		if attendedStaff.Valid {
			t.AttendedByStaffID = &attendedStaff.String
		}
		if ccOfficer.Valid {
			t.VerifiedByCCOfficerID = &ccOfficer.String
		}
		if feedback.Valid {
			t.CustomerFeedbackNotes = feedback.String
		}
		tickets = append(tickets, t)
	}
	return tickets, rows.Err()
}

// CreateMeterSchedule records a new meter appointment.
func (r *PostgresRepository) CreateMeterSchedule(ctx context.Context, s domain.MeterSchedule) error {
	query := `
		INSERT INTO customercare.meter_schedules
		(id, customer_id, ticket_id, schedule_type, scheduled_date, assigned_metering_team_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.CustomerID, s.TicketID, s.ScheduleType, s.ScheduledDate, s.AssignedMeteringTeamID, s.Status, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateMeterSchedule failed: %w", err)
	}
	return nil
}

// ListMeterSchedules lists meter appointments with optional status filter.
func (r *PostgresRepository) ListMeterSchedules(ctx context.Context, status domain.ScheduleStatus) ([]domain.MeterSchedule, error) {
	query := `
		SELECT id, customer_id, ticket_id, schedule_type, scheduled_date, assigned_metering_team_id, status, created_at, updated_at
		FROM customercare.meter_schedules
		WHERE ($1 = '' OR status = $1)
		ORDER BY scheduled_date ASC`

	rows, err := r.db.QueryContext(ctx, query, string(status))
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListMeterSchedules failed: %w", err)
	}
	defer rows.Close()

	var schedules []domain.MeterSchedule
	for rows.Next() {
		var s domain.MeterSchedule
		var ticketID, teamID sql.NullString
		if err := rows.Scan(&s.ID, &s.CustomerID, &ticketID, &s.ScheduleType, &s.ScheduledDate, &teamID, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		if ticketID.Valid {
			s.TicketID = &ticketID.String
		}
		if teamID.Valid {
			s.AssignedMeteringTeamID = &teamID.String
		}
		schedules = append(schedules, s)
	}
	return schedules, rows.Err()
}

