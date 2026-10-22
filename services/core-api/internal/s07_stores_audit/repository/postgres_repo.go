// Package repository implements PostgreSQL persistence for store inventory and material requisitions.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
)

// PostgresRepository provides a PostgreSQL-backed implementation of Repository.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new PostgreSQL repository for S07.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateStoreItem inserts a new inventory SKU.
func (p *PostgresRepository) CreateStoreItem(ctx context.Context, item domain.StoreItem) error {
	query := `
		INSERT INTO stores_audit.store_inventory (
			id, item_code, item_name, category, quantity_available,
			reorder_threshold, unit_of_measure, unit_cost, bin_location,
			last_updated_by_store_officer_id, row_hmac_signature, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	now := time.Now().UTC()
	_, err := p.db.ExecContext(ctx, query,
		item.ID, item.ItemCode, item.ItemName, string(item.Category), item.QuantityAvailable,
		item.ReorderThreshold, item.UnitOfMeasure, item.UnitCost, item.BinLocation,
		item.LastUpdatedBy, item.RowHMACSignature, now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to create store item: %w", err)
	}
	return nil
}

// GetStoreItemByID queries a SKU by UUID.
func (p *PostgresRepository) GetStoreItemByID(ctx context.Context, id string) (*domain.StoreItem, error) {
	query := `
		SELECT id, item_code, item_name, category, quantity_available,
		       reorder_threshold, unit_of_measure, unit_cost, bin_location,
		       last_updated_by_store_officer_id, row_hmac_signature, created_at, updated_at
		FROM stores_audit.store_inventory WHERE id = $1`

	var it domain.StoreItem
	var cat string
	var updatedBy sql.NullString

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&it.ID, &it.ItemCode, &it.ItemName, &cat, &it.QuantityAvailable,
		&it.ReorderThreshold, &it.UnitOfMeasure, &it.UnitCost, &it.BinLocation,
		&updatedBy, &it.RowHMACSignature, &it.CreatedAt, &it.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres: failed to get store item by id: %w", err)
	}
	it.Category = domain.ItemCategory(cat)
	if updatedBy.Valid {
		it.LastUpdatedBy = updatedBy.String
	}
	return &it, nil
}

// GetStoreItemByCode queries a SKU by item code.
func (p *PostgresRepository) GetStoreItemByCode(ctx context.Context, code string) (*domain.StoreItem, error) {
	query := `
		SELECT id, item_code, item_name, category, quantity_available,
		       reorder_threshold, unit_of_measure, unit_cost, bin_location,
		       last_updated_by_store_officer_id, row_hmac_signature, created_at, updated_at
		FROM stores_audit.store_inventory WHERE item_code = $1`

	var it domain.StoreItem
	var cat string
	var updatedBy sql.NullString

	err := p.db.QueryRowContext(ctx, query, code).Scan(
		&it.ID, &it.ItemCode, &it.ItemName, &cat, &it.QuantityAvailable,
		&it.ReorderThreshold, &it.UnitOfMeasure, &it.UnitCost, &it.BinLocation,
		&updatedBy, &it.RowHMACSignature, &it.CreatedAt, &it.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres: failed to get store item by code: %w", err)
	}
	it.Category = domain.ItemCategory(cat)
	if updatedBy.Valid {
		it.LastUpdatedBy = updatedBy.String
	}
	return &it, nil
}

// UpdateStoreItem updates SKU balances and HMAC signatures.
func (p *PostgresRepository) UpdateStoreItem(ctx context.Context, item domain.StoreItem) error {
	query := `
		UPDATE stores_audit.store_inventory
		SET quantity_available = $1, unit_cost = $2, bin_location = $3,
		    last_updated_by_store_officer_id = $4, row_hmac_signature = $5, updated_at = $6
		WHERE id = $7`

	res, err := p.db.ExecContext(ctx, query,
		item.QuantityAvailable, item.UnitCost, item.BinLocation,
		item.LastUpdatedBy, item.RowHMACSignature, time.Now().UTC(), item.ID,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to update store item: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListStoreItems returns SKUs filtered by category and threshold.
func (p *PostgresRepository) ListStoreItems(ctx context.Context, category domain.ItemCategory, lowStockOnly bool) ([]domain.StoreItem, error) {
	query := `
		SELECT id, item_code, item_name, category, quantity_available,
		       reorder_threshold, unit_of_measure, unit_cost, bin_location,
		       last_updated_by_store_officer_id, row_hmac_signature, created_at, updated_at
		FROM stores_audit.store_inventory
		WHERE ($1 = '' OR category = $1)
		  AND (NOT $2 OR quantity_available <= reorder_threshold)
		ORDER BY item_code ASC`

	rows, err := p.db.QueryContext(ctx, query, string(category), lowStockOnly)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to list store items: %w", err)
	}
	defer rows.Close()

	var items []domain.StoreItem
	for rows.Next() {
		var it domain.StoreItem
		var cat string
		var updatedBy sql.NullString
		if err := rows.Scan(
			&it.ID, &it.ItemCode, &it.ItemName, &cat, &it.QuantityAvailable,
			&it.ReorderThreshold, &it.UnitOfMeasure, &it.UnitCost, &it.BinLocation,
			&updatedBy, &it.RowHMACSignature, &it.CreatedAt, &it.UpdatedAt,
		); err != nil {
			return nil, err
		}
		it.Category = domain.ItemCategory(cat)
		if updatedBy.Valid {
			it.LastUpdatedBy = updatedBy.String
		}
		items = append(items, it)
	}
	return items, nil
}

// CreateRequisition persists a material request.
func (p *PostgresRepository) CreateRequisition(ctx context.Context, req domain.MaterialRequisition) error {
	query := `
		INSERT INTO stores_audit.material_requisitions (
			id, req_no, origin_dept_id, linked_ticket_or_job_id, items_json,
			total_cost, in_stock_flag, unit_head_id, audit_consent_status, status,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return fmt.Errorf("postgres: failed to marshal requisition items: %w", err)
	}

	now := time.Now().UTC()
	_, err = p.db.ExecContext(ctx, query,
		req.ID, req.ReqNo, req.OriginDeptID, req.LinkedTicketOrJobID, itemsJSON,
		req.TotalCost, req.InStockFlag, req.UnitHeadID, string(req.AuditConsentStatus),
		string(req.Status), now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to create requisition: %w", err)
	}
	return nil
}

// GetRequisitionByID retrieves a requisition by UUID.
func (p *PostgresRepository) GetRequisitionByID(ctx context.Context, id string) (*domain.MaterialRequisition, error) {
	query := `
		SELECT id, req_no, origin_dept_id, linked_ticket_or_job_id, items_json,
		       total_cost, in_stock_flag, unit_head_id, dept_head_id, dept_or_gm_approved,
		       dept_approved_at, audit_officer_id, audit_consent_status, audit_consent_at,
		       audit_notes, md_procurement_approved, md_approved_at, md_notes,
		       store_officer_id, issued_at, waybill_no, scrap_returned_flag,
		       scrap_return_id, status, created_at, updated_at
		FROM stores_audit.material_requisitions WHERE id = $1`

	return p.scanRequisition(p.db.QueryRowContext(ctx, query, id))
}

// GetRequisitionByReqNo retrieves a requisition by req_no.
func (p *PostgresRepository) GetRequisitionByReqNo(ctx context.Context, reqNo string) (*domain.MaterialRequisition, error) {
	query := `
		SELECT id, req_no, origin_dept_id, linked_ticket_or_job_id, items_json,
		       total_cost, in_stock_flag, unit_head_id, dept_head_id, dept_or_gm_approved,
		       dept_approved_at, audit_officer_id, audit_consent_status, audit_consent_at,
		       audit_notes, md_procurement_approved, md_approved_at, md_notes,
		       store_officer_id, issued_at, waybill_no, scrap_returned_flag,
		       scrap_return_id, status, created_at, updated_at
		FROM stores_audit.material_requisitions WHERE req_no = $1`

	return p.scanRequisition(p.db.QueryRowContext(ctx, query, reqNo))
}

// UpdateRequisition updates requisition state.
func (p *PostgresRepository) UpdateRequisition(ctx context.Context, req domain.MaterialRequisition) error {
	query := `
		UPDATE stores_audit.material_requisitions
		SET dept_head_id = $1, dept_or_gm_approved = $2, dept_approved_at = $3,
		    audit_officer_id = $4, audit_consent_status = $5, audit_consent_at = $6,
		    audit_notes = $7, md_procurement_approved = $8, md_approved_at = $9,
		    md_notes = $10, store_officer_id = $11, issued_at = $12, waybill_no = $13,
		    scrap_returned_flag = $14, scrap_return_id = $15, status = $16, updated_at = $17
		WHERE id = $18`

	res, err := p.db.ExecContext(ctx, query,
		req.DeptHeadID, req.DeptOrGMApproved, req.DeptApprovedAt,
		req.AuditOfficerID, string(req.AuditConsentStatus), req.AuditConsentAt,
		req.AuditNotes, req.MDProcurementApproved, req.MDApprovedAt,
		req.MDNotes, req.StoreOfficerID, req.IssuedAt, req.WaybillNo,
		req.ScrapReturnedFlag, req.ScrapReturnID, string(req.Status), time.Now().UTC(), req.ID,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to update requisition: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListRequisitions retrieves filtered requisitions.
func (p *PostgresRepository) ListRequisitions(ctx context.Context, status domain.RequisitionStatus, unitHeadID string) ([]domain.MaterialRequisition, error) {
	query := `
		SELECT id, req_no, origin_dept_id, linked_ticket_or_job_id, items_json,
		       total_cost, in_stock_flag, unit_head_id, dept_head_id, dept_or_gm_approved,
		       dept_approved_at, audit_officer_id, audit_consent_status, audit_consent_at,
		       audit_notes, md_procurement_approved, md_approved_at, md_notes,
		       store_officer_id, issued_at, waybill_no, scrap_returned_flag,
		       scrap_return_id, status, created_at, updated_at
		FROM stores_audit.material_requisitions
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR unit_head_id = $2)
		ORDER BY created_at DESC`

	rows, err := p.db.QueryContext(ctx, query, string(status), unitHeadID)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to list requisitions: %w", err)
	}
	defer rows.Close()

	var results []domain.MaterialRequisition
	for rows.Next() {
		r, err := p.scanRequisitionRow(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, *r)
	}
	return results, nil
}

