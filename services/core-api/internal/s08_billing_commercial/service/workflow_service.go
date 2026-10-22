// Package service coordinates zero-estimation metered billing, customer meter-purchase refund ledgers,
// wrong billing commercial approval gates, and revenue protection AI anomaly investigations for S08.
package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/domain"
	"disco-enterprise-platform/services/core-api/internal/s08_billing_commercial/repository"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

// WorkflowService coordinates billing workflows, vending integrations, and revenue protection.
type WorkflowService struct {
	repo       repository.Repository
	bus        events.EventBus
	ticketRepo s03repo.Repository
}

// NewWorkflowService creates a new initialized WorkflowService instance.
func NewWorkflowService(
	repo repository.Repository,
	bus events.EventBus,
	ticketRepo s03repo.Repository,
) *WorkflowService {
	return &WorkflowService{
		repo:       repo,
		bus:        bus,
		ticketRepo: ticketRepo,
	}
}

// RecordPaymentRequest parameters for recording a vending payment.
type RecordPaymentRequest struct {
	TxReference            string                        `json:"tx_reference"`
	CustomerID             string                        `json:"customer_id"`
	AccountNumber          string                        `json:"account_number"`
	MeterNumber            string                        `json:"meter_number"`
	FeederID               *string                       `json:"feeder_id,omitempty"`
	DTID                   *string                       `json:"dt_id,omitempty"`
	PaymentChannelCategory domain.PaymentChannelCategory `json:"payment_channel_category"`
	PaymentPlatformSource  string                        `json:"payment_platform_source"`
	AmountPaid             float64                       `json:"amount_paid"`
	TariffRatePerKWh       float64                       `json:"tariff_rate_per_kwh"`
	SettlementVerified     bool                          `json:"settlement_verified"`
	GatewayPayloadJSON     string                        `json:"gateway_payload_json,omitempty"`
}

// RecordVendingPayment processes a vending payment, amortizes any meter refund, generates an STS token,
// and records the transaction without estimated billing.
func (s *WorkflowService) RecordVendingPayment(ctx context.Context, req RecordPaymentRequest) (*domain.VendingTransaction, error) {
	if strings.TrimSpace(req.CustomerID) == "" || strings.TrimSpace(req.AccountNumber) == "" || strings.TrimSpace(req.MeterNumber) == "" {
		return nil, fmt.Errorf("%w: customer_id, account_number, and meter_number are required", domain.ErrInvalidVendingParams)
	}
	if req.TariffRatePerKWh <= 0.0 {
		req.TariffRatePerKWh = 68.00 // Default NERC approved band rate
	}

	txRef := req.TxReference
	if strings.TrimSpace(txRef) == "" {
		txRef = fmt.Sprintf("VEND-%d-%s", time.Now().Unix(), randomHexString(4))
	}

	now := time.Now().UTC()
	var refundCredited float64

	// Check if customer has an active meter refund ledger
	ledger, err := s.repo.GetRefundLedgerByCustomer(ctx, req.CustomerID)
	if err == nil && ledger != nil && ledger.Status == domain.RefundActive {
		refundCredited = domain.ComputeRefundInstallment(ledger, req.AmountPaid)
		_ = s.repo.UpdateRefundLedger(ctx, *ledger)
	}

	kwhVended, err := domain.CalculateVendedEnergy(req.AmountPaid, req.TariffRatePerKWh, refundCredited)
	if err != nil {
		return nil, err
	}

	stsToken := domain.GenerateSTSToken(req.MeterNumber, req.AmountPaid, now.Unix())

	tx := domain.VendingTransaction{
		ID:                     generateUUID(),
		TxReference:            txRef,
		CustomerID:             req.CustomerID,
		AccountNumber:          req.AccountNumber,
		MeterNumber:            req.MeterNumber,
		FeederID:               req.FeederID,
		DTID:                   req.DTID,
		PaymentChannelCategory: req.PaymentChannelCategory,
		PaymentPlatformSource:  req.PaymentPlatformSource,
		AmountPaid:             req.AmountPaid,
		TariffRatePerKWh:       req.TariffRatePerKWh,
		KWhVended:              kwhVended,
		MeterRefundCredited:    refundCredited,
		STSToken:               stsToken,
		SettlementVerified:     req.SettlementVerified,
		GatewayPayloadJSON:     req.GatewayPayloadJSON,
		CreatedAt:              now,
	}

	if err := s.repo.CreateVendingTransaction(ctx, tx); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "billing.payment_received", "s08_billing_commercial", map[string]interface{}{
			"tx_reference":          tx.TxReference,
			"customer_id":           tx.CustomerID,
			"meter_number":          tx.MeterNumber,
			"amount_paid":           tx.AmountPaid,
			"kwh_vended":            tx.KWhVended,
			"meter_refund_credited": tx.MeterRefundCredited,
			"sts_token":             tx.STSToken,
		})
	}

	return &tx, nil
}

// RegisterMeterRefundLedger creates an active refund ledger for a customer who financed their meter.
func (s *WorkflowService) RegisterMeterRefundLedger(ctx context.Context, ledger domain.MeterRefundLedger) (*domain.MeterRefundLedger, error) {
	if strings.TrimSpace(ledger.CustomerID) == "" || strings.TrimSpace(ledger.MeterSerial) == "" {
		return nil, fmt.Errorf("service: customer_id and meter_serial are required")
	}
	if ledger.InitialMeterCost <= 0.0 {
		return nil, fmt.Errorf("service: initial_meter_cost must be greater than zero")
	}
	if ledger.RefundPercentagePerTx <= 0.0 || ledger.RefundPercentagePerTx > 100.0 {
		ledger.RefundPercentagePerTx = 10.00
	}

	now := time.Now().UTC()
	ledger.ID = generateUUID()
	ledger.RemainingRefundBalance = ledger.InitialMeterCost
	ledger.TotalRefundedToDate = 0.0
	ledger.Status = domain.RefundActive
	ledger.CreatedAt = now
	ledger.UpdatedAt = now

	if err := s.repo.CreateRefundLedger(ctx, ledger); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "billing.refund_ledger.registered", "s08_billing_commercial", map[string]interface{}{
			"ledger_id":          ledger.ID,
			"customer_id":        ledger.CustomerID,
			"initial_meter_cost": ledger.InitialMeterCost,
		})
	}

	return &ledger, nil
}

// GetRefundLedgerByCustomer retrieves the refund ledger associated with a customer.
func (s *WorkflowService) GetRefundLedgerByCustomer(ctx context.Context, customerID string) (*domain.MeterRefundLedger, error) {
	return s.repo.GetRefundLedgerByCustomer(ctx, customerID)
}

// GetVendingTransaction retrieves a transaction by ID.
func (s *WorkflowService) GetVendingTransaction(ctx context.Context, id string) (*domain.VendingTransaction, error) {
	return s.repo.GetVendingTransactionByID(ctx, id)
}

// ListVendingTransactions queries vending transactions matching filters.
func (s *WorkflowService) ListVendingTransactions(ctx context.Context, customerID, meterNumber string, limit int) ([]domain.VendingTransaction, error) {
	return s.repo.ListVendingTransactions(ctx, customerID, meterNumber, limit)
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func randomHexString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

