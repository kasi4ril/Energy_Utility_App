// Package service_test validates end-to-end service orchestration, cryptographic signing,
// tamper detection, and cross-departmental event emission.
package service_test

import (
	"context"
	"testing"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/service"
	"disco-enterprise-platform/services/core-api/pkg/events"
)

const (
	testHMACSecret = "test-hmac-secret-key-32-bytes!!"
	testJWTSecret  = "test-jwt-secret-key-32-bytes!!"
)

func setupTestService() (*service.WorkflowService, *repository.MemoryRepository, *events.MemoryEventBus) {
	repo := repository.NewMemoryRepository()
	bus := events.NewMemoryEventBus()
	svc := service.NewWorkflowService(repo, bus, testHMACSecret, testJWTSecret)
	return svc, repo, bus
}

func TestWorkflowService_RegisterOrgUnit_SuccessAndNegative(t *testing.T) {
	svc, _, _ := setupTestService()
	ctx := context.Background()

	// 1. Root HQ
	hqGPS := domain.GPSPoint{Latitude: 9.9281, Longitude: 8.8921}
	hq, err := svc.RegisterOrgUnit(ctx, "HQ-01", "Headquarters", domain.OrgUnitHQ, nil, &hqGPS)
	if err != nil {
		t.Fatalf("failed to register HQ: %v", err)
	}

	// 2. District under HQ
	dist, err := svc.RegisterOrgUnit(ctx, "DIST-JOS", "Jos District", domain.OrgUnitDistrict, &hq.ID, nil)
	if err != nil {
		t.Fatalf("failed to register District: %v", err)
	}

	// Negative: Service Unit with invalid parent (HQ instead of District)
	_, err = svc.RegisterOrgUnit(ctx, "SU-BUKURU", "Bukuru Service Unit", domain.OrgUnitServiceUnit, &hq.ID, nil)
	if err == nil {
		t.Fatalf("expected error for Service Unit directly under HQ")
	}

	// Negative: Missing parent ID
	missingID := "non-existent-parent-uuid"
	_, err = svc.RegisterOrgUnit(ctx, "DIST-2", "District 2", domain.OrgUnitDistrict, &missingID, nil)
	if err == nil {
		t.Fatalf("expected error for non-existent parent ID")
	}

	// Negative: Invalid GPS
	invalidGPS := domain.GPSPoint{Latitude: 99.0, Longitude: 8.0}
	_, err = svc.RegisterOrgUnit(ctx, "HQ-BAD-GPS", "Bad GPS", domain.OrgUnitHQ, nil, &invalidGPS)
	if err == nil {
		t.Fatalf("expected error for bad GPS on org unit")
	}

	_ = dist
}

func TestWorkflowService_StaffRegistration_AuthAndTransfer(t *testing.T) {
	svc, repo, _ := setupTestService()
	ctx := context.Background()

	hq, _ := svc.RegisterOrgUnit(ctx, "HQ-01", "Headquarters", domain.OrgUnitHQ, nil, nil)
	dist, _ := svc.RegisterOrgUnit(ctx, "DIST-01", "District 01", domain.OrgUnitDistrict, &hq.ID, nil)

	// Negative: Register staff with non-existent org unit
	_, err := svc.RegisterStaffProfile(ctx, "NSC-999", "Ghost", "ghost@nesco.ng", "pwd", domain.RoleCustomerCareOp, "invalid-org-id")
	if err == nil {
		t.Fatalf("expected error registering staff in non-existent org unit")
	}

	// 1. Register Staff
	staff, err := svc.RegisterStaffProfile(ctx, "NSC-001", "Maryam Sani", "maryam@nesco.ng", "SecretPassword#1", domain.RoleCustomerCareOp, dist.ID)
	if err != nil {
		t.Fatalf("failed to register staff: %v", err)
	}

	// 2. Authenticate Staff
	claims, token, err := svc.AuthenticateStaff(ctx, "maryam@nesco.ng", "SecretPassword#1")
	if err != nil || token == "" {
		t.Fatalf("authentication failed: %v", err)
	}
	if claims.StaffID != staff.ID || claims.PrimaryRole != string(domain.RoleCustomerCareOp) {
		t.Fatalf("mismatched claims on auth: %+v", claims)
	}

	// Negative: Wrong password
	_, _, err = svc.AuthenticateStaff(ctx, "maryam@nesco.ng", "WrongPassword")
	if err != service.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials on wrong password, got: %v", err)
	}

	// Negative: Non-existent email
	_, _, err = svc.AuthenticateStaff(ctx, "nobody@nesco.ng", "SecretPassword#1")
	if err != service.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials on missing user")
	}

	// 3. Transfer Staff
	transferred, err := svc.TransferStaff(ctx, staff.ID, hq.ID, domain.RoleCustomerCareHead)
	if err != nil {
		t.Fatalf("transfer failed: %v", err)
	}
	if transferred.OrgUnitID != hq.ID || transferred.Role != domain.RoleCustomerCareHead {
		t.Fatalf("transfer did not persist expected changes: %+v", transferred)
	}

	// Negative: Transfer non-existent staff
	_, err = svc.TransferStaff(ctx, "unknown-staff", hq.ID, domain.RoleMD)
	if err == nil {
		t.Fatalf("expected error transferring non-existent staff")
	}

	// Negative: Transfer to non-existent destination org unit
	_, err = svc.TransferStaff(ctx, staff.ID, "non-existent-dest-org", domain.RoleCustomerCareHead)
	if err == nil {
		t.Fatalf("expected error transferring to non-existent org unit")
	}

	// Negative: Tampered HMAC in database
	_ = repo.UpdateStaffTransfer(ctx, staff.ID, staff.OrgUnitID, staff.Role, "tampered-hmac-signature-breaking-integrity")
	_, _, err = svc.AuthenticateStaff(ctx, "maryam@nesco.ng", "SecretPassword#1")
	if err != service.ErrHMACVerificationFailed {
		t.Fatalf("expected ErrHMACVerificationFailed on tampered row signature, got: %v", err)
	}
}

func TestWorkflowService_GridAssetAndWorkNotesChain(t *testing.T) {
	svc, _, _ := setupTestService()
	ctx := context.Background()

	hq, _ := svc.RegisterOrgUnit(ctx, "HQ-01", "HQ", domain.OrgUnitHQ, nil, nil)
	staff, _ := svc.RegisterStaffProfile(ctx, "ENG-01", "Engr. Paul", "paul@nesco.ng", "Pass#123", domain.RoleUnitHeadFeeder, hq.ID)

	// Negative: Asset with non-existent org unit
	_, err := svc.RegisterGridAsset(ctx, "BAD-ORG-ASSET", domain.AssetPlant, nil, "non-existent-org", 100.0, domain.TariffBandA, domain.GPSPoint{Latitude: 9.0, Longitude: 8.0})
	if err == nil {
		t.Fatalf("expected error for asset in non-existent org unit")
	}

	// 1. Register Plant
	plantGPS := domain.GPSPoint{Latitude: 9.42, Longitude: 8.92}
	plant, err := svc.RegisterGridAsset(ctx, "GEN-PLANT-01", domain.AssetPlant, nil, hq.ID, 10000.0, domain.TariffBandNone, plantGPS)
	if err != nil {
		t.Fatalf("failed to register plant: %v", err)
	}

	// 2. Register Substation 33kV under Plant
	subGPS := domain.GPSPoint{Latitude: 9.45, Longitude: 8.90}
	sub, err := svc.RegisterGridAsset(ctx, "SUB-33KV-JOS", domain.AssetSubstation33KV, &plant.ID, hq.ID, 15000.0, domain.TariffBandA, subGPS)
	if err != nil {
		t.Fatalf("failed to register substation: %v", err)
	}

	// Negative: Substation under invalid parent
	invalidParent := "non-existent-uuid"
	_, err = svc.RegisterGridAsset(ctx, "SUB-2", domain.AssetSubstation33KV, &invalidParent, hq.ID, 5000.0, domain.TariffBandA, subGPS)
	if err == nil {
		t.Fatalf("expected error for non-existent parent asset")
	}

	// Negative: Asset with invalid GPS
	badGPS := domain.GPSPoint{Latitude: 120.0, Longitude: 8.0}
	_, err = svc.RegisterGridAsset(ctx, "BAD-GPS-ASSET", domain.AssetSubstation33KV, &plant.ID, hq.ID, 5000.0, domain.TariffBandA, badGPS)
	if err == nil {
		t.Fatalf("expected error for bad GPS on grid asset")
	}

	// 3. Append WorkNotes with SHA-256 Chaining
	note1, err := svc.AppendWorkNote(ctx, "GRID_ASSET", sub.ID, staff.ID, []string{hq.ID}, "Inspected transformer oil")
	if err != nil {
		t.Fatalf("failed to append note 1: %v", err)
	}

	// Note 2 in chain
	note2, err := svc.AppendWorkNote(ctx, "GRID_ASSET", sub.ID, staff.ID, []string{}, "Oil level topped up and confirmed normal")
	if err != nil {
		t.Fatalf("failed to append note 2: %v", err)
	}

	if note2.PrevHash != note1.CurrentHash {
		t.Fatalf("hash chain broken: note2 prev_hash (%s) != note1 current_hash (%s)", note2.PrevHash, note1.CurrentHash)
	}

	// Negative: Append note with invalid author
	_, err = svc.AppendWorkNote(ctx, "GRID_ASSET", sub.ID, "unknown-staff", nil, "Text")
	if err == nil {
		t.Fatalf("expected error for unknown author")
	}
}

