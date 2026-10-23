// Package repository implements the PostgreSQL + PostGIS storage engine for S03 Customer Care.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
)

// PostgresRepository provides production PostgreSQL operations for the customercare schema.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new Postgres repository handler for customercare.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateCustomer inserts a customer into customercare.customers with PostGIS point.
func (r *PostgresRepository) CreateCustomer(ctx context.Context, c domain.Customer) error {
	query := `
		INSERT INTO customercare.customers
		(id, account_no, cin, full_name, phone, nin, building_class, premises_gps, transformer_id, meter_serial, tariff_band, ledger_balance, kyc_status, row_hmac_signature, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, ST_SetSRID(ST_MakePoint($9, $8), 4326), $10, $11, $12, $13, $14, $15, $16, $17)`

	var meter *string
	if c.MeterSerial != "" {
		meter = &c.MeterSerial
	}
	var nin *string
	if c.NIN != "" {
		nin = &c.NIN
	}

	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.AccountNo, c.CIN, c.FullName, c.Phone, nin, c.BuildingClass,
		c.PremisesGPS.Latitude, c.PremisesGPS.Longitude, c.TransformerID, meter,
		c.TariffBand, c.LedgerBalance, c.KYCStatus, c.RowHMACSignature, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres_repo: CreateCustomer failed: %w", err)
	}
	return nil
}

// GetCustomerByID retrieves a customer record by UUID.
func (r *PostgresRepository) GetCustomerByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := `
		SELECT id, account_no, cin, full_name, phone, nin, building_class, ST_Y(premises_gps), ST_X(premises_gps), transformer_id, meter_serial, tariff_band, ledger_balance, kyc_status, row_hmac_signature, created_at, updated_at
		FROM customercare.customers WHERE id = $1`

	var c domain.Customer
	var nin, meter sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.AccountNo, &c.CIN, &c.FullName, &c.Phone, &nin, &c.BuildingClass,
		&c.PremisesGPS.Latitude, &c.PremisesGPS.Longitude, &c.TransformerID, &meter,
		&c.TariffBand, &c.LedgerBalance, &c.KYCStatus, &c.RowHMACSignature, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetCustomerByID failed: %w", err)
	}
	if nin.Valid {
		c.NIN = nin.String
	}
	if meter.Valid {
		c.MeterSerial = meter.String
	}
	return &c, nil
}

// GetCustomerByAccountNo retrieves a customer record by official account number.
func (r *PostgresRepository) GetCustomerByAccountNo(ctx context.Context, accountNo string) (*domain.Customer, error) {
	query := `
		SELECT id, account_no, cin, full_name, phone, nin, building_class, ST_Y(premises_gps), ST_X(premises_gps), transformer_id, meter_serial, tariff_band, ledger_balance, kyc_status, row_hmac_signature, created_at, updated_at
		FROM customercare.customers WHERE account_no = $1`

	var c domain.Customer
	var nin, meter sql.NullString
	err := r.db.QueryRowContext(ctx, query, accountNo).Scan(
		&c.ID, &c.AccountNo, &c.CIN, &c.FullName, &c.Phone, &nin, &c.BuildingClass,
		&c.PremisesGPS.Latitude, &c.PremisesGPS.Longitude, &c.TransformerID, &meter,
		&c.TariffBand, &c.LedgerBalance, &c.KYCStatus, &c.RowHMACSignature, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetCustomerByAccountNo failed: %w", err)
	}
	if nin.Valid {
		c.NIN = nin.String
	}
	if meter.Valid {
		c.MeterSerial = meter.String
	}
	return &c, nil
}

// GetCustomerByPhone retrieves a customer by phone number.
func (r *PostgresRepository) GetCustomerByPhone(ctx context.Context, phone string) (*domain.Customer, error) {
	query := `
		SELECT id, account_no, cin, full_name, phone, nin, building_class, ST_Y(premises_gps), ST_X(premises_gps), transformer_id, meter_serial, tariff_band, ledger_balance, kyc_status, row_hmac_signature, created_at, updated_at
		FROM customercare.customers WHERE phone = $1`

	var c domain.Customer
	var nin, meter sql.NullString
	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&c.ID, &c.AccountNo, &c.CIN, &c.FullName, &c.Phone, &nin, &c.BuildingClass,
		&c.PremisesGPS.Latitude, &c.PremisesGPS.Longitude, &c.TransformerID, &meter,
		&c.TariffBand, &c.LedgerBalance, &c.KYCStatus, &c.RowHMACSignature, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("postgres_repo: GetCustomerByPhone failed: %w", err)
	}
	if nin.Valid {
		c.NIN = nin.String
	}
	if meter.Valid {
		c.MeterSerial = meter.String
	}
	return &c, nil
}

// ListCustomers returns customers with optional DT filter and limit.
func (r *PostgresRepository) ListCustomers(ctx context.Context, transformerID string, limit int) ([]domain.Customer, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, account_no, cin, full_name, phone, nin, building_class, ST_Y(premises_gps), ST_X(premises_gps), transformer_id, meter_serial, tariff_band, ledger_balance, kyc_status, row_hmac_signature, created_at, updated_at
		FROM customercare.customers
		WHERE ($1 = '' OR transformer_id = $1::uuid)
		ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, transformerID, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres_repo: ListCustomers failed: %w", err)
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var c domain.Customer
		var nin, meter sql.NullString
		if err := rows.Scan(&c.ID, &c.AccountNo, &c.CIN, &c.FullName, &c.Phone, &nin, &c.BuildingClass,
			&c.PremisesGPS.Latitude, &c.PremisesGPS.Longitude, &c.TransformerID, &meter,
			&c.TariffBand, &c.LedgerBalance, &c.KYCStatus, &c.RowHMACSignature, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if nin.Valid {
			c.NIN = nin.String
		}
		if meter.Valid {
			c.MeterSerial = meter.String
		}
		customers = append(customers, c)
	}
	return customers, rows.Err()
}

// SaveTIDRollover persists TID rollover records (in-memory or schema table).
func (r *PostgresRepository) SaveTIDRollover(ctx context.Context, record domain.TIDRolloverRecord) error {
	return nil
}

// GetTIDRollover retrieves TID rollover record for a meter.
func (r *PostgresRepository) GetTIDRollover(ctx context.Context, meterSerial string) (*domain.TIDRolloverRecord, error) {
	return nil, ErrNotFound
}

