// Package repository implements helper scanning routines for SQL query rows.
package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
)

func (p *PostgresRepository) scanRequisition(row *sql.Row) (*domain.MaterialRequisition, error) {
	var r domain.MaterialRequisition
	var itemsJSON []byte
	var auditStatus, status string
	var linked, deptHead, auditOfficer, auditNotes, mdNotes, storeOfficer, waybill, scrapID sql.NullString
	var deptApprovedAt, auditConsentAt, mdApprovedAt, issuedAt sql.NullTime

	err := row.Scan(
		&r.ID, &r.ReqNo, &r.OriginDeptID, &linked, &itemsJSON,
		&r.TotalCost, &r.InStockFlag, &r.UnitHeadID, &deptHead, &r.DeptOrGMApproved,
		&deptApprovedAt, &auditOfficer, &auditStatus, &auditConsentAt,
		&auditNotes, &r.MDProcurementApproved, &mdApprovedAt, &mdNotes,
		&storeOfficer, &issuedAt, &waybill, &r.ScrapReturnedFlag,
		&scrapID, &status, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres: scan requisition failed: %w", err)
	}

	_ = json.Unmarshal(itemsJSON, &r.Items)
	r.AuditConsentStatus = domain.AuditConsentStatus(auditStatus)
	r.Status = domain.RequisitionStatus(status)

	if linked.Valid {
		r.LinkedTicketOrJobID = &linked.String
	}
	if deptHead.Valid {
		r.DeptHeadID = &deptHead.String
	}
	if deptApprovedAt.Valid {
		r.DeptApprovedAt = &deptApprovedAt.Time
	}
	if auditOfficer.Valid {
		r.AuditOfficerID = &auditOfficer.String
	}
	if auditConsentAt.Valid {
		r.AuditConsentAt = &auditConsentAt.Time
	}
	if auditNotes.Valid {
		r.AuditNotes = &auditNotes.String
	}
	if mdApprovedAt.Valid {
		r.MDApprovedAt = &mdApprovedAt.Time
	}
	if mdNotes.Valid {
		r.MDNotes = &mdNotes.String
	}
	if storeOfficer.Valid {
		r.StoreOfficerID = &storeOfficer.String
	}
	if issuedAt.Valid {
		r.IssuedAt = &issuedAt.Time
	}
	if waybill.Valid {
		r.WaybillNo = &waybill.String
	}
	if scrapID.Valid {
		r.ScrapReturnID = &scrapID.String
	}
	return &r, nil
}

func (p *PostgresRepository) scanRequisitionRow(rows *sql.Rows) (*domain.MaterialRequisition, error) {
	var r domain.MaterialRequisition
	var itemsJSON []byte
	var auditStatus, status string
	var linked, deptHead, auditOfficer, auditNotes, mdNotes, storeOfficer, waybill, scrapID sql.NullString
	var deptApprovedAt, auditConsentAt, mdApprovedAt, issuedAt sql.NullTime

	err := rows.Scan(
		&r.ID, &r.ReqNo, &r.OriginDeptID, &linked, &itemsJSON,
		&r.TotalCost, &r.InStockFlag, &r.UnitHeadID, &deptHead, &r.DeptOrGMApproved,
		&deptApprovedAt, &auditOfficer, &auditStatus, &auditConsentAt,
		&auditNotes, &r.MDProcurementApproved, &mdApprovedAt, &mdNotes,
		&storeOfficer, &issuedAt, &waybill, &r.ScrapReturnedFlag,
		&scrapID, &status, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("postgres: scan requisition row failed: %w", err)
	}

	_ = json.Unmarshal(itemsJSON, &r.Items)
	r.AuditConsentStatus = domain.AuditConsentStatus(auditStatus)
	r.Status = domain.RequisitionStatus(status)

	if linked.Valid {
		r.LinkedTicketOrJobID = &linked.String
	}
	if deptHead.Valid {
		r.DeptHeadID = &deptHead.String
	}
	if deptApprovedAt.Valid {
		r.DeptApprovedAt = &deptApprovedAt.Time
	}
	if auditOfficer.Valid {
		r.AuditOfficerID = &auditOfficer.String
	}
	if auditConsentAt.Valid {
		r.AuditConsentAt = &auditConsentAt.Time
	}
	if auditNotes.Valid {
		r.AuditNotes = &auditNotes.String
	}
	if mdApprovedAt.Valid {
		r.MDApprovedAt = &mdApprovedAt.Time
	}
	if mdNotes.Valid {
		r.MDNotes = &mdNotes.String
	}
	if storeOfficer.Valid {
		r.StoreOfficerID = &storeOfficer.String
	}
	if issuedAt.Valid {
		r.IssuedAt = &issuedAt.Time
	}
	if waybill.Valid {
		r.WaybillNo = &waybill.String
	}
	if scrapID.Valid {
		r.ScrapReturnID = &scrapID.String
	}
	return &r, nil
}

