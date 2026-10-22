// Package service orchestrates store stock inventory management, dual-path material requisitions,
// reverse logistics scrap tracking, and internal audit compliance workflows for S07.
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

var (
	// ErrItemNotFound indicates SKU was not found in warehouse inventory.
	ErrItemNotFound = errorsNotFound()
	// ErrInsufficientStock indicates requested quantity exceeds warehouse balance for in-stock delivery.
	ErrInsufficientStock = fmt.Errorf("service: insufficient warehouse stock available for item")
)

func errorsNotFound() error {
	return repository.ErrNotFound
}

// WorkflowService coordinates inventory ledgers, material requisitions, and audit gates.
type WorkflowService struct {
	repo         repository.Repository
	bus          events.EventBus
	hmacSecret   string
	aiSidecarURL string
}

// NewWorkflowService constructs an initialized S07 service orchestrator.
func NewWorkflowService(
	repo repository.Repository,
	bus events.EventBus,
	hmacSecret string,
	aiSidecarURL string,
) *WorkflowService {
	return &WorkflowService{
		repo:         repo,
		bus:          bus,
		hmacSecret:   hmacSecret,
		aiSidecarURL: aiSidecarURL,
	}
}

// RegisterStoreItem adds a new SKU to inventory and computes its cryptographic row HMAC.
func (s *WorkflowService) RegisterStoreItem(ctx context.Context, item domain.StoreItem) (*domain.StoreItem, error) {
	if err := domain.ValidateStoreItem(item); err != nil {
		return nil, err
	}

	if item.ID == "" {
		item.ID = crypto.MustUUID()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now

	hmacPayload := fmt.Sprintf("%s:%s:%.2f:%.2f:%s", item.ItemCode, item.ItemName, item.QuantityAvailable, item.UnitCost, item.BinLocation)
	sig, _ := crypto.ComputeRowHMAC(s.hmacSecret, hmacPayload)
	item.RowHMACSignature = sig

	if err := s.repo.CreateStoreItem(ctx, item); err != nil {
		return nil, fmt.Errorf("service: failed to create store item: %w", err)
	}

	_ = s.bus.Publish(ctx, "stores.item_registered", "s07_stores_audit", map[string]interface{}{
		"item_id":   item.ID,
		"item_code": item.ItemCode,
		"category":  item.Category,
	})

	return &item, nil
}

// UpdateInventoryStock adjusts quantity, updates HMAC, and checks reorder threshold.
func (s *WorkflowService) UpdateInventoryStock(ctx context.Context, itemID string, delta float64, officerID string) (*domain.StoreItem, error) {
	item, err := s.repo.GetStoreItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	newQty := item.QuantityAvailable + delta
	if newQty < 0.0 {
		return nil, ErrInsufficientStock
	}
	item.QuantityAvailable = newQty
	item.LastUpdatedBy = officerID

	hmacPayload := fmt.Sprintf("%s:%s:%.2f:%.2f:%s", item.ItemCode, item.ItemName, item.QuantityAvailable, item.UnitCost, item.BinLocation)
	sig, _ := crypto.ComputeRowHMAC(s.hmacSecret, hmacPayload)
	item.RowHMACSignature = sig

	if err := s.repo.UpdateStoreItem(ctx, *item); err != nil {
		return nil, fmt.Errorf("service: failed to update store item: %w", err)
	}

	if domain.EvaluateReorderStatus(item.QuantityAvailable, item.ReorderThreshold) {
		_ = s.bus.Publish(ctx, "stores.reorder_threshold_breached", "s07_stores_audit", map[string]interface{}{
			"item_id":            item.ID,
			"item_code":          item.ItemCode,
			"quantity_available": item.QuantityAvailable,
			"reorder_threshold":  item.ReorderThreshold,
		})
	}

	return item, nil
}

// GetStockByItemCode retrieves real-time stock balances (used by @GridAssistant bot and API).
func (s *WorkflowService) GetStockByItemCode(ctx context.Context, itemCode string) (*domain.StoreItem, error) {
	return s.repo.GetStoreItemByCode(ctx, itemCode)
}

// ListInventory returns warehouse inventory.
func (s *WorkflowService) ListInventory(ctx context.Context, category domain.ItemCategory, lowStockOnly bool) ([]domain.StoreItem, error) {
	return s.repo.ListStoreItems(ctx, category, lowStockOnly)
}

// InitiateRequisition initiates Path A or Path B material requisition depending on warehouse balances.
func (s *WorkflowService) InitiateRequisition(ctx context.Context, req domain.MaterialRequisition) (*domain.MaterialRequisition, error) {
	// Inspect warehouse availability to determine InStockFlag and populate UnitCost
	inStock := true
	for i := range req.Items {
		item, err := s.repo.GetStoreItemByCode(ctx, req.Items[i].ItemCode)
		if err != nil || item.QuantityAvailable < req.Items[i].QuantityRequested {
			inStock = false
		} else {
			req.Items[i].ItemID = item.ID
			if req.Items[i].UnitCost == 0.0 {
				req.Items[i].UnitCost = item.UnitCost
			}
		}
	}
	req.InStockFlag = inStock

	totalCost, err := domain.CalculateRequisitionTotal(req.Items)
	if err != nil {
		return nil, err
	}
	req.TotalCost = totalCost

	now := time.Now().UTC()
	if req.ID == "" {
		req.ID = crypto.MustUUID()
	}
	if req.ReqNo == "" {
		req.ReqNo = fmt.Sprintf("REQ-%d-%s", now.Year(), strings.ToUpper(crypto.MustUUID()[:8]))
	}
	req.Status = domain.StatusRequisitionSubmitted
	req.AuditConsentStatus = domain.AuditConsentPending
	req.CreatedAt = now
	req.UpdatedAt = now

	if err := s.repo.CreateRequisition(ctx, req); err != nil {
		return nil, fmt.Errorf("service: failed to create requisition: %w", err)
	}

	_ = s.bus.Publish(ctx, "stores.requisition_submitted", "s07_stores_audit", map[string]interface{}{
		"req_id":        req.ID,
		"req_no":        req.ReqNo,
		"in_stock_flag": req.InStockFlag,
		"total_cost":    req.TotalCost,
	})

	return &req, nil
}

// ApproveRequisitionByDept records Departmental Head / Generation GM approval.
func (s *WorkflowService) ApproveRequisitionByDept(ctx context.Context, reqID, deptHeadID string) (*domain.MaterialRequisition, error) {
	req, err := s.repo.GetRequisitionByID(ctx, reqID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateDeptApproval(req, deptHeadID); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateRequisition(ctx, *req); err != nil {
		return nil, fmt.Errorf("service: failed to update requisition: %w", err)
	}

	_ = s.bus.Publish(ctx, "stores.requisition_dept_approved", "s07_stores_audit", map[string]interface{}{
		"req_id":        req.ID,
		"req_no":        req.ReqNo,
		"dept_head_id":  deptHeadID,
		"in_stock_flag": req.InStockFlag,
	})

	return req, nil
}

// ConsentRequisitionByAudit records Internal Audit consent after verifying tickets and prior scrap returns.
func (s *WorkflowService) ConsentRequisitionByAudit(ctx context.Context, reqID, auditorID string, consent bool, notes string) (*domain.MaterialRequisition, error) {
	req, err := s.repo.GetRequisitionByID(ctx, reqID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateAuditConsent(req, auditorID, consent, notes); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateRequisition(ctx, *req); err != nil {
		return nil, fmt.Errorf("service: failed to update requisition: %w", err)
	}

	_ = s.bus.Publish(ctx, "stores.requisition_audit_decided", "s07_stores_audit", map[string]interface{}{
		"req_id":         req.ID,
		"req_no":         req.ReqNo,
		"consent_status": req.AuditConsentStatus,
		"auditor_id":     auditorID,
	})

	return req, nil
}

// ApproveRequisitionByMD records final Managing Director procurement approval for out-of-stock purchases.
func (s *WorkflowService) ApproveRequisitionByMD(ctx context.Context, reqID, mdNotes string) (*domain.MaterialRequisition, error) {
	req, err := s.repo.GetRequisitionByID(ctx, reqID)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateMDProcurementApproval(req, mdNotes); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateRequisition(ctx, *req); err != nil {
		return nil, fmt.Errorf("service: failed to update requisition: %w", err)
	}

	_ = s.bus.Publish(ctx, "stores.requisition_md_approved", "s07_stores_audit", map[string]interface{}{
		"req_id":     req.ID,
		"req_no":     req.ReqNo,
		"total_cost": req.TotalCost,
	})

	return req, nil
}

// IssueMaterialWaybill issues materials from store, deducts quantities, and assigns waybill.
func (s *WorkflowService) IssueMaterialWaybill(ctx context.Context, reqID, storeOfficerID string) (*domain.MaterialRequisition, error) {
	req, err := s.repo.GetRequisitionByID(ctx, reqID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	waybillNo := fmt.Sprintf("WAYBILL-%d-%s", now.Year(), strings.ToUpper(crypto.MustUUID()[:8]))

	if err := domain.ValidateMaterialIssuance(req, storeOfficerID, waybillNo); err != nil {
		return nil, err
	}

	// Deduct stock for each line item
	for _, it := range req.Items {
		item, err := s.repo.GetStoreItemByCode(ctx, it.ItemCode)
		if err == nil && item != nil {
			_, _ = s.UpdateInventoryStock(ctx, item.ID, -it.QuantityRequested, storeOfficerID)
		}
	}

	if err := s.repo.UpdateRequisition(ctx, *req); err != nil {
		return nil, fmt.Errorf("service: failed to update requisition: %w", err)
	}

	_ = s.bus.Publish(ctx, "stores.material_issued", "s07_stores_audit", map[string]interface{}{
		"req_id":           req.ID,
		"waybill_no":       req.WaybillNo,
		"store_officer_id": storeOfficerID,
	})

	return req, nil
}

// LogScrapReturn records returned scrap, linking to requisition and clearing reverse logistics holds.
func (s *WorkflowService) LogScrapReturn(ctx context.Context, ret domain.ScrapReturn) (*domain.ScrapReturn, error) {
	req, err := s.repo.GetRequisitionByID(ctx, ret.RequisitionID)
	if err != nil {
		return nil, fmt.Errorf("service: linked requisition not found: %w", err)
	}

	now := time.Now().UTC()
	if ret.ID == "" {
		ret.ID = crypto.MustUUID()
	}
	if ret.ReturnNo == "" {
		ret.ReturnNo = fmt.Sprintf("SCRAP-%d-%s", now.Year(), strings.ToUpper(crypto.MustUUID()[:8]))
	}
	ret.ReturnedAt = now
	ret.CreatedAt = now

	if err := s.repo.CreateScrapReturn(ctx, ret); err != nil {
		return nil, fmt.Errorf("service: failed to create scrap return: %w", err)
	}

	req.ScrapReturnedFlag = true
	req.ScrapReturnID = &ret.ID
	_ = s.repo.UpdateRequisition(ctx, *req)

	_ = s.bus.Publish(ctx, "stores.scrap_returned", "s07_stores_audit", map[string]interface{}{
		"scrap_id":       ret.ID,
		"return_no":      ret.ReturnNo,
		"requisition_id": ret.RequisitionID,
		"team_lead_id":   ret.TeamLeadID,
	})

	return &ret, nil
}

// GetRequisition retrieves requisition by ID.
func (s *WorkflowService) GetRequisition(ctx context.Context, id string) (*domain.MaterialRequisition, error) {
	return s.repo.GetRequisitionByID(ctx, id)
}

// ListRequisitions retrieves material requisitions.
func (s *WorkflowService) ListRequisitions(ctx context.Context, status domain.RequisitionStatus, unitHeadID string) ([]domain.MaterialRequisition, error) {
	return s.repo.ListRequisitions(ctx, status, unitHeadID)
}
