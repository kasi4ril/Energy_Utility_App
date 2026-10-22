// Package repository implements the PostgreSQL + PostGIS storage engine for S01 Foundation.
// This file covers electrical grid assets and immutable WorkNotes.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

// CreateGridAsset registers an electrical asset in foundation.grid_assets with PostGIS geometry.
func (r *PostgresRepository) CreateGridAsset(ctx context.Context, a domain.GridAsset) error {
	query := `
		INSERT INTO foundation.grid_assets (id, parent_asset_id, org_unit_id, asset_code, asset_type, rated_capacity_kva, tariff_band, gps_point, operational_state, row_hmac_signature, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, ST_SetSRID(ST_MakePoint($9, $8), 4326), $10, $11, $12, $13)`

	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.ParentAssetID, a.OrgUnitID, a.AssetCode, a.AssetType, a.RatedCapacityKVA, a.TariffBand, a.GPSPoint.Latitude, a.GPSPoint.Longitude, a.OperationalState, a.RowHMACSignature, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateGridAsset failed: %w", err)
	}
	return nil
}

// GetGridAssetByID fetches a grid asset by ID.
func (r *PostgresRepository) GetGridAssetByID(ctx context.Context, id string) (*domain.GridAsset, error) {
	query := `
		SELECT id, parent_asset_id, org_unit_id, asset_code, asset_type, rated_capacity_kva, tariff_band, ST_Y(gps_point), ST_X(gps_point), operational_state, row_hmac_signature, created_at, updated_at
		FROM foundation.grid_assets WHERE id = $1`

	var a domain.GridAsset
	var tariff sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.ParentAssetID, &a.OrgUnitID, &a.AssetCode, &a.AssetType, &a.RatedCapacityKVA, &tariff, &a.GPSPoint.Latitude, &a.GPSPoint.Longitude, &a.OperationalState, &a.RowHMACSignature, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetGridAssetByID failed: %w", err)
	}
	if tariff.Valid {
		a.TariffBand = domain.TariffBand(tariff.String)
	}
	return &a, nil
}

// GetGridAssetByCode fetches an asset by unique asset code.
func (r *PostgresRepository) GetGridAssetByCode(ctx context.Context, code string) (*domain.GridAsset, error) {
	query := `
		SELECT id, parent_asset_id, org_unit_id, asset_code, asset_type, rated_capacity_kva, tariff_band, ST_Y(gps_point), ST_X(gps_point), operational_state, row_hmac_signature, created_at, updated_at
		FROM foundation.grid_assets WHERE asset_code = $1`

	var a domain.GridAsset
	var tariff sql.NullString
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&a.ID, &a.ParentAssetID, &a.OrgUnitID, &a.AssetCode, &a.AssetType, &a.RatedCapacityKVA, &tariff, &a.GPSPoint.Latitude, &a.GPSPoint.Longitude, &a.OperationalState, &a.RowHMACSignature, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetGridAssetByCode failed: %w", err)
	}
	if tariff.Valid {
		a.TariffBand = domain.TariffBand(tariff.String)
	}
	return &a, nil
}

// ListGridAssets lists assets with optional filters.
func (r *PostgresRepository) ListGridAssets(ctx context.Context, parentID string, assetType domain.AssetType) ([]domain.GridAsset, error) {
	query := `
		SELECT id, parent_asset_id, org_unit_id, asset_code, asset_type, rated_capacity_kva, tariff_band, ST_Y(gps_point), ST_X(gps_point), operational_state, row_hmac_signature, created_at, updated_at
		FROM foundation.grid_assets
		WHERE ($1 = '' OR parent_asset_id = $1::uuid)
		  AND ($2 = '' OR asset_type = $2)
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, parentID, string(assetType))
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListGridAssets failed: %w", err)
	}
	defer rows.Close()

	var assets []domain.GridAsset
	for rows.Next() {
		var a domain.GridAsset
		var tariff sql.NullString
		if err := rows.Scan(&a.ID, &a.ParentAssetID, &a.OrgUnitID, &a.AssetCode, &a.AssetType, &a.RatedCapacityKVA, &tariff, &a.GPSPoint.Latitude, &a.GPSPoint.Longitude, &a.OperationalState, &a.RowHMACSignature, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		if tariff.Valid {
			a.TariffBand = domain.TariffBand(tariff.String)
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

// CreateWorkNote inserts an immutable note to foundation.work_notes.
func (r *PostgresRepository) CreateWorkNote(ctx context.Context, n domain.WorkNote) error {
	query := `
		INSERT INTO foundation.work_notes (id, entity_type, entity_id, author_staff_id, tagged_unit_ids, note_text, ollama_summary, prev_hash, current_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.EntityType, n.EntityID, n.AuthorStaffID, n.TaggedUnitIDs, n.NoteText, n.OllamaSummary, n.PrevHash, n.CurrentHash, n.CreatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateWorkNote failed: %w", err)
	}
	return nil
}

// GetLatestWorkNote returns the latest note in the chain for an entity.
func (r *PostgresRepository) GetLatestWorkNote(ctx context.Context, entityType, entityID string) (*domain.WorkNote, error) {
	query := `
		SELECT id, entity_type, entity_id, author_staff_id, tagged_unit_ids, note_text, ollama_summary, prev_hash, current_hash, created_at
		FROM foundation.work_notes
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at DESC LIMIT 1`

	var n domain.WorkNote
	var summary sql.NullString
	err := r.db.QueryRowContext(ctx, query, entityType, entityID).Scan(
		&n.ID, &n.EntityType, &n.EntityID, &n.AuthorStaffID, &n.TaggedUnitIDs, &n.NoteText, &summary, &n.PrevHash, &n.CurrentHash, &n.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No note yet
		}
		return nil, fmt.Errorf("postgres_repo: GetLatestWorkNote failed: %w", err)
	}
	if summary.Valid {
		n.OllamaSummary = summary.String
	}
	return &n, nil
}

// ListWorkNotesByEntity returns all notes for an entity ordered chronologically.
func (r *PostgresRepository) ListWorkNotesByEntity(ctx context.Context, entityType, entityID string) ([]domain.WorkNote, error) {
	query := `
		SELECT id, entity_type, entity_id, author_staff_id, tagged_unit_ids, note_text, ollama_summary, prev_hash, current_hash, created_at
		FROM foundation.work_notes
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListWorkNotesByEntity failed: %w", err)
	}
	defer rows.Close()

	var notes []domain.WorkNote
	for rows.Next() {
		var n domain.WorkNote
		var summary sql.NullString
		if err := rows.Scan(&n.ID, &n.EntityType, &n.EntityID, &n.AuthorStaffID, &n.TaggedUnitIDs, &n.NoteText, &summary, &n.PrevHash, &n.CurrentHash, &n.CreatedAt); err != nil {
			return nil, err
		}
		if summary.Valid {
			n.OllamaSummary = summary.String
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

