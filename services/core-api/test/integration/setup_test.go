// Package integration exercises end-to-end multi-sprint workflows across NESCO bounded contexts.
package integration

import (
	"context"
	"net/http"
	"testing"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	s01repo "disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	s01svc "disco-enterprise-platform/services/core-api/internal/s01_foundation/service"
	s01trans "disco-enterprise-platform/services/core-api/internal/s01_foundation/transport"
	s03repo "disco-enterprise-platform/services/core-api/internal/s03_customer_care/repository"
	s03svc "disco-enterprise-platform/services/core-api/internal/s03_customer_care/service"
	s03trans "disco-enterprise-platform/services/core-api/internal/s03_customer_care/transport"
	s04repo "disco-enterprise-platform/services/core-api/internal/s04_services_maint/repository"
	s04svc "disco-enterprise-platform/services/core-api/internal/s04_services_maint/service"
	s04trans "disco-enterprise-platform/services/core-api/internal/s04_services_maint/transport"
	s05repo "disco-enterprise-platform/services/core-api/internal/s05_generation_master/repository"
	s05svc "disco-enterprise-platform/services/core-api/internal/s05_generation_master/service"
	s05trans "disco-enterprise-platform/services/core-api/internal/s05_generation_master/transport"
	s06repo "disco-enterprise-platform/services/core-api/internal/s06_energy_management/repository"
	s06svc "disco-enterprise-platform/services/core-api/internal/s06_energy_management/service"
	s06trans "disco-enterprise-platform/services/core-api/internal/s06_energy_management/transport"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const (
	testHMACKey = "nesco-master-hmac-integrity-key-32b!"
	testJWTKey  = "nesco-master-jwt-security-key-32b!!"
)

func generateHMACHeader(payload []byte, secret string) string {
	sig, _ := crypto.ComputeRowHMAC(secret, string(payload))
	return sig
}

func setupE2EEnvironment(t *testing.T) (*http.ServeMux, *s01repo.MemoryRepository, *s03repo.MemoryRepository, *s04repo.MemoryRepository, *s05repo.MemoryRepository, *s06repo.MemoryRepository) {
	t.Helper()
	bus := events.NewMemoryEventBus()

	s01Repo := s01repo.NewMemoryRepository()
	s01Service := s01svc.NewWorkflowService(s01Repo, bus, testHMACKey, testJWTKey)

	// Seed staff profile for login test
	ctx := context.Background()
	hq, _ := s01Service.RegisterOrgUnit(ctx, "HQ", "HQ Main", s01domain.OrgUnitHQ, nil, nil)
	_, _ = s01Service.RegisterStaffProfile(ctx, "NSC-01", "Amina Bello", "amina@nesco.ng", "SecretP@ss1", s01domain.RoleMD, hq.ID)

	// Seed Grid Asset for Generation Plant 1
	_ = s01Repo.CreateGridAsset(ctx, s01domain.GridAsset{
		ID:        "PLANT_01_KURRA_FALLS",
		AssetCode: "GEN-01",
		AssetType: s01domain.AssetPlant,
		OrgUnitID: hq.ID,
		GPSPoint:  s01domain.GPSPoint{Latitude: 9.42, Longitude: 8.85},
	})

	// Seed Grid Asset for 33kV Vom Feeder
	_ = s01Repo.CreateGridAsset(ctx, s01domain.GridAsset{
		ID:        "FDR-33KV-VOM",
		AssetCode: "FDR-VOM-01",
		AssetType: s01domain.AssetFeeder11KV,
		OrgUnitID: hq.ID,
		GPSPoint:  s01domain.GPSPoint{Latitude: 9.7523, Longitude: 8.8541},
	})

	s03Repo := s03repo.NewMemoryRepository()
	s03Service := s03svc.NewWorkflowService(s03Repo, s01Repo, bus, testHMACKey)

	s04Repo := s04repo.NewMemoryRepository()
	s04Service := s04svc.NewWorkflowService(s04Repo, s01Repo, s03Repo, bus, "")

	s05Repo := s05repo.NewMemoryRepository()
	s05Service := s05svc.NewWorkflowService(s05Repo, s01Repo, bus, testHMACKey)

	s06Repo := s06repo.NewMemoryRepository()
	s06Service := s06svc.NewWorkflowService(s06Repo, s01Repo, bus, testHMACKey)

	mux := http.NewServeMux()
	s01trans.NewHTTPHandler(s01Service, s01Repo, testJWTKey).RegisterRoutes(mux)
	s03trans.NewHTTPHandler(s03Service, s03Repo, testJWTKey).RegisterRoutes(mux)
	s04trans.NewHTTPHandler(s04Service, s04Repo, testJWTKey).RegisterRoutes(mux)
	s05trans.NewHTTPHandler(s05Service, s05Repo, testJWTKey).RegisterRoutes(mux)
	s06trans.NewHTTPHandler(s06Service, s06Repo, testJWTKey).RegisterRoutes(mux)

	return mux, s01Repo, s03Repo, s04Repo, s05Repo, s06Repo
}

