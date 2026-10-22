// Package transport_test validates DTO field validations and error handling for S07.
package transport_test

import (
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	"disco-enterprise-platform/services/core-api/internal/s07_stores_audit/transport"
)

func TestCreateStoreItemDTO_Validate(t *testing.T) {
	valid := transport.CreateStoreItemDTO{
		ItemCode:          "CBL-33KV-150",
		ItemName:          "33kV XLPE Cable",
		Category:          domain.CategoryCablesConductors,
		QuantityAvailable: 100.0,
		ReorderThreshold:  20.0,
		UnitOfMeasure:     "METERS",
		UnitCost:          45000.0,
		BinLocation:       "YARD-A1",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO to pass, got: %v", err)
	}

	// Missing code
	missingCode := valid
	missingCode.ItemCode = "  "
	if err := missingCode.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields, got: %v", err)
	}

	// Invalid category
	badCat := valid
	badCat.Category = ""
	if err := badCat.Validate(); err != transport.ErrInvalidItemCategory {
		t.Errorf("expected ErrInvalidItemCategory, got: %v", err)
	}

	// Negative cost
	negCost := valid
	negCost.UnitCost = -10.0
	if err := negCost.Validate(); err != transport.ErrNegativeAmount {
		t.Errorf("expected ErrNegativeAmount, got: %v", err)
	}
}

func TestInitiateRequisitionDTO_Validate(t *testing.T) {
	valid := transport.InitiateRequisitionDTO{
		OriginDeptID: "dept-services",
		Items: []transport.RequisitionItemDTO{
			{ItemCode: "FUSE-100A", QuantityRequested: 5.0, UnitCost: 2000.0},
		},
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid DTO to pass, got: %v", err)
	}

	// Missing origin dept
	noDept := valid
	noDept.OriginDeptID = ""
	if err := noDept.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields, got: %v", err)
	}

	// Empty items
	noItems := valid
	noItems.Items = []transport.RequisitionItemDTO{}
	if err := noItems.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields, got: %v", err)
	}

	// Zero quantity
	zeroQty := valid
	zeroQty.Items = []transport.RequisitionItemDTO{{ItemCode: "FUSE-100A", QuantityRequested: 0.0}}
	if err := zeroQty.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields on 0 quantity, got: %v", err)
	}
}

func TestLogScrapReturnDTO_Validate(t *testing.T) {
	valid := transport.LogScrapReturnDTO{
		RequisitionID:       "req-001",
		ItemsReturnedJSON:   `[{"item":"CBL-SCRAP","qty":50}]`,
		ScrapConditionNotes: "Damaged during storm",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid scrap return DTO, got: %v", err)
	}

	missingReq := valid
	missingReq.RequisitionID = ""
	if err := missingReq.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields, got: %v", err)
	}
}

func TestSubmitClaimDTO_Validate(t *testing.T) {
	valid := transport.SubmitClaimDTO{
		ClaimType: domain.ClaimOvertime,
		Amount:    25000.0,
		ClaimDate: "2026-10-22",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid claim DTO, got: %v", err)
	}

	// Missing claim type
	missingType := valid
	missingType.ClaimType = ""
	if err := missingType.Validate(); err != transport.ErrInvalidClaimType {
		t.Errorf("expected ErrInvalidClaimType, got: %v", err)
	}

	// Negative amount
	negAmount := valid
	negAmount.Amount = -500.0
	if err := negAmount.Validate(); err != transport.ErrNegativeAmount {
		t.Errorf("expected ErrNegativeAmount, got: %v", err)
	}

	// Missing date
	missingDate := valid
	missingDate.ClaimDate = ""
	if err := missingDate.Validate(); err != transport.ErrMissingMandatoryFields {
		t.Errorf("expected ErrMissingMandatoryFields, got: %v", err)
	}
}

