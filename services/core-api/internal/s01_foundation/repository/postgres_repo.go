// Package repository implements the PostgreSQL + PostGIS storage engine for S01 Foundation.
// This file covers organizational units and staff profiles.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

// PostgresRepository provides production PostgreSQL operations with spatial PostGIS querying.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new Postgres repository handler.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateOrgUnit inserts an administrative unit into foundation.org_units.
func (r *PostgresRepository) CreateOrgUnit(ctx context.Context, u domain.OrgUnit) error {
	query := `
		INSERT INTO foundation.org_units (id, parent_id, unit_code, name, unit_type, location_gps, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CASE WHEN $6::float8 IS NOT NULL THEN ST_SetSRID(ST_MakePoint($7, $6), 4326) ELSE NULL END, $8, $9, $10)`

	var lat, lon *float64
	if u.LocationGPS != nil {
		lat = &u.LocationGPS.Latitude
		lon = &u.LocationGPS.Longitude
	}

	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.ParentID, u.UnitCode, u.Name, u.UnitType, lat, lon, u.IsActive, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateOrgUnit failed: %w", err)
	}
	return nil
}

// GetOrgUnitByID retrieves an organizational unit by ID.
func (r *PostgresRepository) GetOrgUnitByID(ctx context.Context, id string) (*domain.OrgUnit, error) {
	query := `
		SELECT id, parent_id, unit_code, name, unit_type, ST_Y(location_gps), ST_X(location_gps), is_active, created_at, updated_at
		FROM foundation.org_units WHERE id = $1`

	var u domain.OrgUnit
	var lat, lon sql.NullFloat64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.ParentID, &u.UnitCode, &u.Name, &u.UnitType, &lat, &lon, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetOrgUnitByID failed: %w", err)
	}
	if lat.Valid && lon.Valid {
		u.LocationGPS = &domain.GPSPoint{Latitude: lat.Float64, Longitude: lon.Float64}
	}
	return &u, nil
}

// ListOrgUnits fetches all organizational units.
func (r *PostgresRepository) ListOrgUnits(ctx context.Context) ([]domain.OrgUnit, error) {
	query := `
		SELECT id, parent_id, unit_code, name, unit_type, ST_Y(location_gps), ST_X(location_gps), is_active, created_at, updated_at
		FROM foundation.org_units ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListOrgUnits failed: %w", err)
	}
	defer rows.Close()

	var units []domain.OrgUnit
	for rows.Next() {
		var u domain.OrgUnit
		var lat, lon sql.NullFloat64
		if err := rows.Scan(&u.ID, &u.ParentID, &u.UnitCode, &u.Name, &u.UnitType, &lat, &lon, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		if lat.Valid && lon.Valid {
			u.LocationGPS = &domain.GPSPoint{Latitude: lat.Float64, Longitude: lon.Float64}
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

// CreateStaffProfile inserts an employee record into foundation.staff_profiles.
func (r *PostgresRepository) CreateStaffProfile(ctx context.Context, s domain.StaffProfile) error {
	query := `
		INSERT INTO foundation.staff_profiles (id, staff_number, full_name, email, role, org_unit_id, certification_status, active_shift_state, password_hash, row_hmac_signature, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.StaffNumber, s.FullName, s.Email, s.Role, s.OrgUnitID, s.CertificationStatus, s.ActiveShiftState, s.PasswordHash, s.RowHMACSignature, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateStaffProfile failed: %w", err)
	}
	return nil
}

// GetStaffProfileByID retrieves staff by UUID.
func (r *PostgresRepository) GetStaffProfileByID(ctx context.Context, id string) (*domain.StaffProfile, error) {
	query := `
		SELECT id, staff_number, full_name, email, role, org_unit_id, certification_status, active_shift_state, password_hash, row_hmac_signature, created_at, updated_at
		FROM foundation.staff_profiles WHERE id = $1`

	var s domain.StaffProfile
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.StaffNumber, &s.FullName, &s.Email, &s.Role, &s.OrgUnitID, &s.CertificationStatus, &s.ActiveShiftState, &s.PasswordHash, &s.RowHMACSignature, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetStaffProfileByID failed: %w", err)
	}
	return &s, nil
}

// GetStaffProfileByEmail retrieves staff by official email.
func (r *PostgresRepository) GetStaffProfileByEmail(ctx context.Context, email string) (*domain.StaffProfile, error) {
	query := `
		SELECT id, staff_number, full_name, email, role, org_unit_id, certification_status, active_shift_state, password_hash, row_hmac_signature, created_at, updated_at
		FROM foundation.staff_profiles WHERE email = $1`

	var s domain.StaffProfile
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&s.ID, &s.StaffNumber, &s.FullName, &s.Email, &s.Role, &s.OrgUnitID, &s.CertificationStatus, &s.ActiveShiftState, &s.PasswordHash, &s.RowHMACSignature, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetStaffProfileByEmail failed: %w", err)
	}
	return &s, nil
}

// GetStaffProfileByStaffNumber retrieves staff by company staff code.
func (r *PostgresRepository) GetStaffProfileByStaffNumber(ctx context.Context, staffNumber string) (*domain.StaffProfile, error) {
	query := `
		SELECT id, staff_number, full_name, email, role, org_unit_id, certification_status, active_shift_state, password_hash, row_hmac_signature, created_at, updated_at
		FROM foundation.staff_profiles WHERE staff_number = $1`

	var s domain.StaffProfile
	err := r.db.QueryRowContext(ctx, query, staffNumber).Scan(
		&s.ID, &s.StaffNumber, &s.FullName, &s.Email, &s.Role, &s.OrgUnitID, &s.CertificationStatus, &s.ActiveShiftState, &s.PasswordHash, &s.RowHMACSignature, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetStaffProfileByStaffNumber failed: %w", err)
	}
	return &s, nil
}

// UpdateStaffTransfer executes administrative re-assignment of staff.
func (r *PostgresRepository) UpdateStaffTransfer(ctx context.Context, staffID, newOrgUnitID string, newRole domain.UserRole, newHMAC string) error {
	query := `
		UPDATE foundation.staff_profiles
		SET org_unit_id = $2, role = $3, row_hmac_signature = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, staffID, newOrgUnitID, newRole, newHMAC)
	if err != nil {
		return fmt.Errorf("postgres_repo: UpdateStaffTransfer failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListStaffProfiles lists all staff or filters by org unit.
func (r *PostgresRepository) ListStaffProfiles(ctx context.Context, orgUnitID string) ([]domain.StaffProfile, error) {
	query := `
		SELECT id, staff_number, full_name, email, role, org_unit_id, certification_status, active_shift_state, password_hash, row_hmac_signature, created_at, updated_at
		FROM foundation.staff_profiles
		WHERE ($1 = '' OR org_unit_id = $1::uuid)
		ORDER BY full_name ASC`

	rows, err := r.db.QueryContext(ctx, query, orgUnitID)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListStaffProfiles failed: %w", err)
	}
	defer rows.Close()

	var staffList []domain.StaffProfile
	for rows.Next() {
		var s domain.StaffProfile
		if err := rows.Scan(&s.ID, &s.StaffNumber, &s.FullName, &s.Email, &s.Role, &s.OrgUnitID, &s.CertificationStatus, &s.ActiveShiftState, &s.PasswordHash, &s.RowHMACSignature, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		staffList = append(staffList, s)
	}
	return staffList, rows.Err()
}

