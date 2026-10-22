// Package repository implements the PostgreSQL + PostGIS storage engine for S04 Services & Maintenance.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
)

// PostgresRepository provides production PostgreSQL operations for the services_maint schema.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new Postgres repository handler.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateWorkOrder inserts a maintenance work order into services_maint.maintenance_work_orders.
func (r *PostgresRepository) CreateWorkOrder(ctx context.Context, wo domain.MaintenanceWorkOrder) error {
	query := `
		INSERT INTO services_maint.maintenance_work_orders
		(id, order_no, ticket_id, unit_type, work_order_type, priority, unit_head_id, team_lead_id, asset_id, ptw_required, ptw_permit_id, status, materials_requisition_ref, instructions, rework_notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := r.db.ExecContext(ctx, query,
		wo.ID, wo.OrderNo, wo.TicketID, wo.UnitType, wo.WorkOrderType, wo.Priority,
		wo.UnitHeadID, wo.TeamLeadID, wo.AssetID, wo.PTWRequired, wo.PTWPermitID,
		wo.Status, wo.MaterialsRequisitionRef, wo.Instructions, wo.ReworkNotes,
		wo.CreatedAt, wo.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateWorkOrder failed: %w", err)
	}
	return nil
}

// GetWorkOrderByID retrieves a work order by UUID.
func (r *PostgresRepository) GetWorkOrderByID(ctx context.Context, id string) (*domain.MaintenanceWorkOrder, error) {
	query := `
		SELECT id, order_no, ticket_id, unit_type, work_order_type, priority, unit_head_id, team_lead_id, asset_id, ptw_required, ptw_permit_id, status, materials_requisition_ref, instructions, rework_notes, created_at, updated_at
		FROM services_maint.maintenance_work_orders WHERE id = $1`

	var wo domain.MaintenanceWorkOrder
	var ticketID, ptwPermitID, matRef, rework sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wo.ID, &wo.OrderNo, &ticketID, &wo.UnitType, &wo.WorkOrderType, &wo.Priority,
		&wo.UnitHeadID, &wo.TeamLeadID, &wo.AssetID, &wo.PTWRequired, &ptwPermitID,
		&wo.Status, &matRef, &wo.Instructions, &rework, &wo.CreatedAt, &wo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetWorkOrderByID failed: %w", err)
	}
	if ticketID.Valid {
		wo.TicketID = &ticketID.String
	}
	if ptwPermitID.Valid {
		wo.PTWPermitID = &ptwPermitID.String
	}
	if matRef.Valid {
		wo.MaterialsRequisitionRef = &matRef.String
	}
	if rework.Valid {
		wo.ReworkNotes = &rework.String
	}
	return &wo, nil
}

// GetWorkOrderByOrderNo retrieves a work order by human-readable order number.
func (r *PostgresRepository) GetWorkOrderByOrderNo(ctx context.Context, orderNo string) (*domain.MaintenanceWorkOrder, error) {
	query := `
		SELECT id, order_no, ticket_id, unit_type, work_order_type, priority, unit_head_id, team_lead_id, asset_id, ptw_required, ptw_permit_id, status, materials_requisition_ref, instructions, rework_notes, created_at, updated_at
		FROM services_maint.maintenance_work_orders WHERE order_no = $1`

	var wo domain.MaintenanceWorkOrder
	var ticketID, ptwPermitID, matRef, rework sql.NullString
	err := r.db.QueryRowContext(ctx, query, orderNo).Scan(
		&wo.ID, &wo.OrderNo, &ticketID, &wo.UnitType, &wo.WorkOrderType, &wo.Priority,
		&wo.UnitHeadID, &wo.TeamLeadID, &wo.AssetID, &wo.PTWRequired, &ptwPermitID,
		&wo.Status, &matRef, &wo.Instructions, &rework, &wo.CreatedAt, &wo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetWorkOrderByOrderNo failed: %w", err)
	}
	if ticketID.Valid {
		wo.TicketID = &ticketID.String
	}
	if ptwPermitID.Valid {
		wo.PTWPermitID = &ptwPermitID.String
	}
	if matRef.Valid {
		wo.MaterialsRequisitionRef = &matRef.String
	}
	if rework.Valid {
		wo.ReworkNotes = &rework.String
	}
	return &wo, nil
}

// UpdateWorkOrderStatus transitions the work order status and optional rework notes.
func (r *PostgresRepository) UpdateWorkOrderStatus(ctx context.Context, id string, status domain.WorkOrderStatus, reworkNotes *string) error {
	query := `
		UPDATE services_maint.maintenance_work_orders
		SET status = $2, rework_notes = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id, string(status), reworkNotes)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdateWorkOrderStatus failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListWorkOrders queries orders by optional unit type, status, and team lead filters.
func (r *PostgresRepository) ListWorkOrders(ctx context.Context, unit domain.UnitType, status domain.WorkOrderStatus, teamLeadID string) ([]domain.MaintenanceWorkOrder, error) {
	query := `
		SELECT id, order_no, ticket_id, unit_type, work_order_type, priority, unit_head_id, team_lead_id, asset_id, ptw_required, ptw_permit_id, status, materials_requisition_ref, instructions, rework_notes, created_at, updated_at
		FROM services_maint.maintenance_work_orders
		WHERE ($1 = '' OR unit_type = $1)
		  AND ($2 = '' OR status = $2)
		  AND ($3 = '' OR team_lead_id = $3::uuid)
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, string(unit), string(status), teamLeadID)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListWorkOrders failed: %w", err)
	}
	defer rows.Close()

	var orders []domain.MaintenanceWorkOrder
	for rows.Next() {
		var wo domain.MaintenanceWorkOrder
		var ticketID, ptwPermitID, matRef, rework sql.NullString
		if err := rows.Scan(
			&wo.ID, &wo.OrderNo, &ticketID, &wo.UnitType, &wo.WorkOrderType, &wo.Priority,
			&wo.UnitHeadID, &wo.TeamLeadID, &wo.AssetID, &wo.PTWRequired, &ptwPermitID,
			&wo.Status, &matRef, &wo.Instructions, &rework, &wo.CreatedAt, &wo.UpdatedAt); err != nil {
			return nil, err
		}
		if ticketID.Valid {
			wo.TicketID = &ticketID.String
		}
		if ptwPermitID.Valid {
			wo.PTWPermitID = &ptwPermitID.String
		}
		if matRef.Valid {
			wo.MaterialsRequisitionRef = &matRef.String
		}
		if rework.Valid {
			wo.ReworkNotes = &rework.String
		}
		orders = append(orders, wo)
	}
	return orders, rows.Err()
}

