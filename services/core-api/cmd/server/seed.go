// Package main provides demonstration seeding for local runtime exploration.
package main

import (
	"context"
	"log"
	"time"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01svc "disco-enterprise-platform/services/core-api/internal/s01_foundation/service"
	s03domain "disco-enterprise-platform/services/core-api/internal/s03_customer_care/domain"
	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	s04domain "disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
	s04repo "disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	s05domain "disco-enterprise-platform/services/core-api/internal/s05_generation_master/domain"
	s05repo "disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	s06domain "disco-enterprise-platform/services/core-api/internal/s06_energy_management/domain"
	s06repo "disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	s07domain "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/domain"
	s07repo "disco-enterprise-platform/services/core-api/internal/s07_stores_audit/repository"
)

// SeedDemoData pre-populates in-memory stores with realistic demo staff, assets, and workflows.
func SeedDemoData(
	s01Svc *s01svc.WorkflowService,
	s03Repo *s03repo.MemoryRepository,
	s04Repo *s04repo.MemoryRepository,
	s05Repo *s05repo.MemoryRepository,
	s06Repo *s06repo.MemoryRepository,
	s07Repo *s07repo.MemoryRepository,
) {
	ctx := context.Background()

	// 1. Root OrgUnit
	hq, err := s01Svc.RegisterOrgUnit(ctx, "HQ-JOS", "NESCO Headquarters Jos", s01domain.OrgUnitHQ, nil, nil)
	if err != nil {
		log.Printf("[WARN] Seed HQ registration error: %v", err)
		return
	}

	// 2. Demo Staff Accounts across all Departments
	demoStaff := []struct {
		num, name, email, pass string
		role                   s01domain.UserRole
	}{
		{"NES-001", "Amina Bello", "amina@nesco.ng", "SecretP@ss1", s01domain.RoleMD},
		{"NES-002", "MD Executive", "md@nesco.ng", "SecretP@ss1", s01domain.RoleMD},
		{"NES-003", "Customer Care Rep", "care@nesco.ng", "SecretP@ss1", s01domain.RoleCustomerCareOp},
		{"NES-004", "Engr. Bitrus Feeder", "feeder@nesco.ng", "SecretP@ss1", s01domain.RoleUnitHeadFeeder},
		{"NES-005", "Engr. Gyang Transformer", "transformer@nesco.ng", "SecretP@ss1", s01domain.RoleUnitHeadTx},
		{"NES-006", "Engr. Danladi Metering", "metering@nesco.ng", "SecretP@ss1", s01domain.RoleUnitHeadMetering},
		{"NES-007", "Plant Master Operator", "generation@nesco.ng", "SecretP@ss1", s01domain.RoleMasterStationOp},
		{"NES-008", "Central Grid Dispatcher", "dispatch@nesco.ng", "SecretP@ss1", s01domain.RoleCentralDispatchOp},
		{"NES-009", "Chief Store Officer", "store@nesco.ng", "SecretP@ss1", s01domain.RoleStoreOfficer},
		{"NES-010", "Head Internal Auditor", "audit@nesco.ng", "SecretP@ss1", s01domain.RoleInternalAuditor},
		{"NES-011", "Billing & Vending Officer", "billing@nesco.ng", "SecretP@ss1", s01domain.RoleUnitHeadBilling},
	}

	for _, s := range demoStaff {
		_, _ = s01Svc.RegisterStaffProfile(ctx, s.num, s.name, s.email, s.pass, s.role, hq.ID)
	}

	// 3. Demo Customer & Service Ticket (S03)
	custID := "CUST-KURA-2026-001"
	_ = s03Repo.CreateCustomer(ctx, s03domain.Customer{
		ID:            custID,
		AccountNo:     "ACC-00142",
		CIN:           "CIN-00142",
		FullName:      "Danladi Gyang",
		Phone:         "08031234567",
		BuildingClass: s03domain.BuildingResidential,
		PremisesGPS:   s01domain.GPSPoint{Latitude: 9.7523, Longitude: 8.8541},
		TransformerID: "TX-DT-500KVA-MARKET",
		MeterSerial:   "MTR-98214",
		TariffBand:    "BAND_A",
		KYCStatus:     s03domain.KYCVerified,
		CreatedAt:     time.Now().UTC(),
	})
	_ = s03Repo.CreateTicket(ctx, s03domain.ServiceTicket{
		ID:             "TKT-2026-001",
		TicketNo:       "TKT-001",
		CustomerID:     &custID,
		Channel:        s03domain.ChannelCallCenter,
		Category:       s03domain.CategoryLineSnap,
		AssignedUnitID: "UNIT-FEEDER-01",
		Status:         s03domain.StatusLogged,
		SLADueAt:       time.Now().UTC().Add(24 * time.Hour),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})

	// 4. Demo Work Order (S04)
	_ = s04Repo.CreateWorkOrder(ctx, s04domain.MaintenanceWorkOrder{
		ID:            "WO-2026-FDR-001",
		OrderNo:       "WO-FDR-001",
		UnitType:      s04domain.UnitFeederMaint,
		WorkOrderType: s04domain.WOCorrectiveRepair,
		Priority:      s04domain.PriorityCritical,
		AssetID:       "FDR-33KV-VOM",
		Status:        s04domain.StatusDispatched,
		CreatedAt:     time.Now().UTC(),
	})

	// 5. Demo Generation Telemetry Log (S05)
	_ = s05Repo.CreateTelemetryLog(ctx, s05domain.MasterStationTelemetryLog{
		ID:                 "GEN-LOG-001",
		ShiftID:            "SHIFT-MORN-01",
		ShiftWindow:        s05domain.ShiftMorning,
		PlantID:            "PLANT_01_KURRA_FALLS",
		PlantState:         s05domain.StateInService,
		WaterLevelM:        24.5,
		MWGenerated:        14.8,
		StatorWindingTempC: 68.2,
		LoggedByOperatorID: "STAFF-OP-01",
		RecordedAt:         time.Now().UTC(),
	})

	// 6. Demo Dispatch Telemetry Log (S06)
	_ = s06Repo.CreateDispatchLog(ctx, s06domain.DispatchTelemetryLog{
		ID:             "DISP-LOG-001",
		ShiftID:        "SHIFT-DISP-01",
		ShiftWindow:    s06domain.ShiftMorning,
		GenerationMWIn: 45.0,
		TCNMWIn:        15.0,
		TotalBulkMW:    60.0,
		OperatorID:     "STAFF-DISP-01",
		LoggedAt:       time.Now().UTC(),
	})

	// 7. Demo Store Inventory Item (S07)
	_ = s07Repo.CreateStoreItem(ctx, s07domain.StoreItem{
		ID:                "ITEM-FUSE-11KV",
		ItemCode:          "FUSE-11KV-100A",
		ItemName:          "11kV Drop-Out Fuse 100A",
		Category:          s07domain.CategorySubstationHardware,
		UnitOfMeasure:     "PCS",
		QuantityAvailable: 45,
		ReorderThreshold:  10,
		UnitCost:          4500.0,
		BinLocation:       "BIN-A-04",
		LastUpdatedBy:     "STAFF-STORE-01",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	})

	log.Println("[INFO] Successfully seeded demo staff accounts and sample telemetry data.")
}
