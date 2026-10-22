// Package repository implements PostgreSQL persistence for vending and meter refund ledgers.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
)

// PostgresRepository provides a PostgreSQL-backed implementation of Repository.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a new PostgreSQL repository for S08.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateVendingTransaction inserts a new vending payment record.
func (p *PostgresRepository) CreateVendingTransaction(ctx context.Context, tx domain.VendingTransaction) error {
	query := `
		INSERT INTO billing_commercial.vending_and_payment_transactions (
			id, tx_reference, customer_id, account_number, meter_number,
			feeder_id, dt_id, payment_channel_category, payment_platform_source,
			amount_paid, tariff_rate_per_kwh, kwh_vended, meter_refund_credited,
			sts_token, settlement_verified, gateway_payload_json, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	now := tx.CreatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}

	var payload sql.NullString
	if tx.GatewayPayloadJSON != "" {
		payload = sql.NullString{String: tx.GatewayPayloadJSON, Valid: true}
	}

	_, err := p.db.ExecContext(ctx, query,
		tx.ID, tx.TxReference, tx.CustomerID, tx.AccountNumber, tx.MeterNumber,
		tx.FeederID, tx.DTID, string(tx.PaymentChannelCategory), tx.PaymentPlatformSource,
		tx.AmountPaid, tx.TariffRatePerKWh, tx.KWhVended, tx.MeterRefundCredited,
		tx.STSToken, tx.SettlementVerified, payload, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to create vending transaction: %w", err)
	}
	return nil
}

// GetVendingTransactionByID retrieves transaction by UUID.
func (p *PostgresRepository) GetVendingTransactionByID(ctx context.Context, id string) (*domain.VendingTransaction, error) {
	query := `
		SELECT id, tx_reference, customer_id, account_number, meter_number,
		       feeder_id, dt_id, payment_channel_category, payment_platform_source,
		       amount_paid, tariff_rate_per_kwh, kwh_vended, meter_refund_credited,
		       sts_token, settlement_verified, gateway_payload_json, created_at
		FROM billing_commercial.vending_and_payment_transactions
		WHERE id = $1`

	return p.scanVendingRow(p.db.QueryRowContext(ctx, query, id))
}

// GetVendingTransactionByRef retrieves transaction by reference.
func (p *PostgresRepository) GetVendingTransactionByRef(ctx context.Context, txRef string) (*domain.VendingTransaction, error) {
	query := `
		SELECT id, tx_reference, customer_id, account_number, meter_number,
		       feeder_id, dt_id, payment_channel_category, payment_platform_source,
		       amount_paid, tariff_rate_per_kwh, kwh_vended, meter_refund_credited,
		       sts_token, settlement_verified, gateway_payload_json, created_at
		FROM billing_commercial.vending_and_payment_transactions
		WHERE tx_reference = $1`

	return p.scanVendingRow(p.db.QueryRowContext(ctx, query, txRef))
}

func (p *PostgresRepository) scanVendingRow(row *sql.Row) (*domain.VendingTransaction, error) {
	var tx domain.VendingTransaction
	var channel string
	var feederID, dtID, stsToken, payload sql.NullString

	err := row.Scan(
		&tx.ID, &tx.TxReference, &tx.CustomerID, &tx.AccountNumber, &tx.MeterNumber,
		&feederID, &dtID, &channel, &tx.PaymentPlatformSource,
		&tx.AmountPaid, &tx.TariffRatePerKWh, &tx.KWhVended, &tx.MeterRefundCredited,
		&stsToken, &tx.SettlementVerified, &payload, &tx.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: scan vending error: %w", err)
	}

	tx.PaymentChannelCategory = domain.PaymentChannelCategory(channel)
	if feederID.Valid {
		tx.FeederID = &feederID.String
	}
	if dtID.Valid {
		tx.DTID = &dtID.String
	}
	if stsToken.Valid {
		tx.STSToken = stsToken.String
	}
	if payload.Valid {
		tx.GatewayPayloadJSON = payload.String
	}
	return &tx, nil
}

// ListVendingTransactions filters transactions by customer and meter number.
func (p *PostgresRepository) ListVendingTransactions(ctx context.Context, customerID, meterNumber string, limit int) ([]domain.VendingTransaction, error) {
	query := `
		SELECT id, tx_reference, customer_id, account_number, meter_number,
		       feeder_id, dt_id, payment_channel_category, payment_platform_source,
		       amount_paid, tariff_rate_per_kwh, kwh_vended, meter_refund_credited,
		       sts_token, settlement_verified, gateway_payload_json, created_at
		FROM billing_commercial.vending_and_payment_transactions
		WHERE ($1 = '' OR customer_id = $1)
		  AND ($2 = '' OR meter_number = $2)
		ORDER BY created_at DESC
		LIMIT $3`

	if limit <= 0 {
		limit = 50
	}

	rows, err := p.db.QueryContext(ctx, query, customerID, meterNumber, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres: list vending error: %w", err)
	}
	defer rows.Close()

	res := make([]domain.VendingTransaction, 0)
	for rows.Next() {
		var tx domain.VendingTransaction
		var channel string
		var feederID, dtID, stsToken, payload sql.NullString

		if err := rows.Scan(
			&tx.ID, &tx.TxReference, &tx.CustomerID, &tx.AccountNumber, &tx.MeterNumber,
			&feederID, &dtID, &channel, &tx.PaymentPlatformSource,
			&tx.AmountPaid, &tx.TariffRatePerKWh, &tx.KWhVended, &tx.MeterRefundCredited,
			&stsToken, &tx.SettlementVerified, &payload, &tx.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: scan vending row error: %w", err)
		}
		tx.PaymentChannelCategory = domain.PaymentChannelCategory(channel)
		if feederID.Valid {
			tx.FeederID = &feederID.String
		}
		if dtID.Valid {
			tx.DTID = &dtID.String
		}
		if stsToken.Valid {
			tx.STSToken = stsToken.String
		}
		if payload.Valid {
			tx.GatewayPayloadJSON = payload.String
		}
		res = append(res, tx)
	}
	return res, nil
}

// CreateRefundLedger inserts a new customer refund ledger.
func (p *PostgresRepository) CreateRefundLedger(ctx context.Context, ledger domain.MeterRefundLedger) error {
	query := `
		INSERT INTO billing_commercial.meter_refund_ledgers (
			id, customer_id, meter_serial, initial_meter_cost,
			refund_percentage_per_tx, total_refunded_to_date, remaining_refund_balance,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	now := time.Now().UTC()
	_, err := p.db.ExecContext(ctx, query,
		ledger.ID, ledger.CustomerID, ledger.MeterSerial, ledger.InitialMeterCost,
		ledger.RefundPercentagePerTx, ledger.TotalRefundedToDate, ledger.RemainingRefundBalance,
		string(ledger.Status), now, now,
	)
	if err != nil {
		return fmt.Errorf("postgres: failed to create refund ledger: %w", err)
	}
	return nil
}

// GetRefundLedgerByID retrieves ledger by primary ID.
func (p *PostgresRepository) GetRefundLedgerByID(ctx context.Context, id string) (*domain.MeterRefundLedger, error) {
	query := `
		SELECT id, customer_id, meter_serial, initial_meter_cost,
		       refund_percentage_per_tx, total_refunded_to_date, remaining_refund_balance,
		       status, created_at, updated_at
		FROM billing_commercial.meter_refund_ledgers
		WHERE id = $1`

	return p.scanRefundRow(p.db.QueryRowContext(ctx, query, id))
}

// GetRefundLedgerByCustomer retrieves ledger by customer ID.
func (p *PostgresRepository) GetRefundLedgerByCustomer(ctx context.Context, customerID string) (*domain.MeterRefundLedger, error) {
	query := `
		SELECT id, customer_id, meter_serial, initial_meter_cost,
		       refund_percentage_per_tx, total_refunded_to_date, remaining_refund_balance,
		       status, created_at, updated_at
		FROM billing_commercial.meter_refund_ledgers
		WHERE customer_id = $1`

	return p.scanRefundRow(p.db.QueryRowContext(ctx, query, customerID))
}

// GetRefundLedgerByMeter retrieves ledger by meter serial.
func (p *PostgresRepository) GetRefundLedgerByMeter(ctx context.Context, meterSerial string) (*domain.MeterRefundLedger, error) {
	query := `
		SELECT id, customer_id, meter_serial, initial_meter_cost,
		       refund_percentage_per_tx, total_refunded_to_date, remaining_refund_balance,
		       status, created_at, updated_at
		FROM billing_commercial.meter_refund_ledgers
		WHERE meter_serial = $1`

	return p.scanRefundRow(p.db.QueryRowContext(ctx, query, meterSerial))
}

func (p *PostgresRepository) scanRefundRow(row *sql.Row) (*domain.MeterRefundLedger, error) {
	var l domain.MeterRefundLedger
	var status string

	err := row.Scan(
		&l.ID, &l.CustomerID, &l.MeterSerial, &l.InitialMeterCost,
		&l.RefundPercentagePerTx, &l.TotalRefundedToDate, &l.RemainingRefundBalance,
		&status, &l.CreatedAt, &l.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres: scan refund error: %w", err)
	}
	l.Status = domain.RefundStatus(status)
	return &l, nil
}

// UpdateRefundLedger updates ledger balances and status.
func (p *PostgresRepository) UpdateRefundLedger(ctx context.Context, ledger domain.MeterRefundLedger) error {
	query := `
		UPDATE billing_commercial.meter_refund_ledgers
		SET total_refunded_to_date = $2,
		    remaining_refund_balance = $3,
		    status = $4,
		    updated_at = $5
		WHERE id = $1`

	now := time.Now().UTC()
	res, err := p.db.ExecContext(ctx, query,
		ledger.ID, ledger.TotalRefundedToDate, ledger.RemainingRefundBalance,
		string(ledger.Status), now,
	)
	if err != nil {
		return fmt.Errorf("postgres: update refund error: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
