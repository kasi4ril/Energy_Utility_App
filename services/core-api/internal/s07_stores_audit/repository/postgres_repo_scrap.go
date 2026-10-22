// Package repository implements PostgreSQL queries for reverse logistics scrap returns.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
)

// CreateScrapReturn persists a returned scrap ledger entry.
func (p *PostgresRepository) CreateScrapReturn(ctx context.Context, ret domain.ScrapReturn) error {
	query := `
		INSERT INTO stores_audit.scrap_returns (
			id, return_no, requisition_id, team_lead_id, store_officer_id,
			items_returned_json, scrap_condition_notes, returned_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	now := time.Now().UTC()
	_, err := p.db.ExecContext(ctx, query,
		ret.ID, ret.ReturnNo, ret.RequisitionID, ret.TeamLeadID, ret.StoreOfficerID,
		ret.ItemsReturnedJSON, ret.ScrapConditionNotes, now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to create scrap return: %w", err)
	}
	return nil
}

// GetScrapReturnByID retrieves a scrap return record by UUID.
func (p *PostgresRepository) GetScrapReturnByID(ctx context.Context, id string) (*domain.ScrapReturn, error) {
	query := `
		SELECT id, return_no, requisition_id, team_lead_id, store_officer_id,
		       items_returned_json, scrap_condition_notes, returned_at, created_at
		FROM stores_audit.scrap_returns WHERE id = $1`

	var s domain.ScrapReturn
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.ReturnNo, &s.RequisitionID, &s.TeamLeadID, &s.StoreOfficerID,
		&s.ItemsReturnedJSON, &s.ScrapConditionNotes, &s.ReturnedAt, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres: failed to get scrap return by id: %w", err)
	}
	return &s, nil
}

// ListScrapReturns lists scrap records.
func (p *PostgresRepository) ListScrapReturns(ctx context.Context, teamLeadID string) ([]domain.ScrapReturn, error) {
	query := `
		SELECT id, return_no, requisition_id, team_lead_id, store_officer_id,
		       items_returned_json, scrap_condition_notes, returned_at, created_at
		FROM stores_audit.scrap_returns
		WHERE ($1 = '' OR team_lead_id = $1)
		ORDER BY returned_at DESC`

	rows, err := p.db.QueryContext(ctx, query, teamLeadID)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to list scrap returns: %w", err)
	}
	defer rows.Close()

	var results []domain.ScrapReturn
	for rows.Next() {
		var s domain.ScrapReturn
		if err := rows.Scan(
			&s.ID, &s.ReturnNo, &s.RequisitionID, &s.TeamLeadID, &s.StoreOfficerID,
			&s.ItemsReturnedJSON, &s.ScrapConditionNotes, &s.ReturnedAt, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}

