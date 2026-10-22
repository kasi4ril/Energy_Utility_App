// Package repository_test verifies data access methods and uniqueness guarantees.
package repository_test

import (
	"context"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/internal/s01_foundation/repository"
)

func TestMemoryRepository_OrgUnits(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	u := domain.OrgUnit{
		ID:        "org-hq-1",
		UnitCode:  "HQ-MAIN",
		Name:      "Headquarters Jos",
		UnitType:  domain.OrgUnitHQ,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.CreateOrgUnit(ctx, u); err != nil {
		t.Fatalf("failed to create org unit: %v", err)
	}

	// Duplicate code rejection
	if err := repo.CreateOrgUnit(ctx, u); err != repository.ErrDuplicateKey {
		t.Fatalf("expected ErrDuplicateKey, got: %v", err)
	}

	fetched, err := repo.GetOrgUnitByID(ctx, u.ID)
	if err != nil || fetched.UnitCode != u.UnitCode {
		t.Fatalf("failed to fetch created org unit: %v", err)
	}

	list, err := repo.ListOrgUnits(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 org unit, got: %v", err)
	}

	// Missing ID
	if _, err := repo.GetOrgUnitByID(ctx, "non-existent"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing unit")
	}
}

func TestMemoryRepository_StaffAndTransfer(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	staff := domain.StaffProfile{
		ID:                  "stf-1",
		StaffNumber:         "NSC-001",
		FullName:            "Musa Ibrahim",
		Email:               "musa@nesco.ng",
		Role:                domain.RoleTeamLeadFeeder,
		OrgUnitID:           "org-1",
		CertificationStatus: domain.CertActive,
		ActiveShiftState:    domain.ShiftOffDuty,
		RowHMACSignature:    "hmac-sig-1",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := repo.CreateStaffProfile(ctx, staff); err != nil {
		t.Fatalf("failed to create staff: %v", err)
	}

	// Duplicate staff number rejection
	if err := repo.CreateStaffProfile(ctx, staff); err != repository.ErrDuplicateKey {
		t.Fatalf("expected ErrDuplicateKey on duplicate staff number")
	}

	// Duplicate email rejection
	dupEmailStaff := staff
	dupEmailStaff.ID = "stf-2"
	dupEmailStaff.StaffNumber = "NSC-002"
	if err := repo.CreateStaffProfile(ctx, dupEmailStaff); err != repository.ErrDuplicateKey {
		t.Fatalf("expected ErrDuplicateKey on duplicate email")
	}

	// Query by staff number and email
	byNum, err := repo.GetStaffProfileByStaffNumber(ctx, "NSC-001")
	if err != nil || byNum.FullName != "Musa Ibrahim" {
		t.Fatalf("failed to fetch by staff number: %v", err)
	}

	byEmail, err := repo.GetStaffProfileByEmail(ctx, "musa@nesco.ng")
	if err != nil || byEmail.ID != staff.ID {
		t.Fatalf("failed to fetch by email: %v", err)
	}

	// List staff with filter and without filter
	staffList, err := repo.ListStaffProfiles(ctx, "org-1")
	if err != nil || len(staffList) != 1 {
		t.Fatalf("expected 1 staff in org-1")
	}
	staffListEmpty, err := repo.ListStaffProfiles(ctx, "non-existent-org")
	if err != nil || len(staffListEmpty) != 0 {
		t.Fatalf("expected 0 staff in non-existent org")
	}
	allStaff, err := repo.ListStaffProfiles(ctx, "")
	if err != nil || len(allStaff) != 1 {
		t.Fatalf("expected 1 staff overall")
	}

	// Transfer staff
	err = repo.UpdateStaffTransfer(ctx, staff.ID, "org-2", domain.RoleUnitHeadFeeder, "new-hmac-sig")
	if err != nil {
		t.Fatalf("transfer failed: %v", err)
	}

	updated, _ := repo.GetStaffProfileByID(ctx, staff.ID)
	if updated.OrgUnitID != "org-2" || updated.Role != domain.RoleUnitHeadFeeder {
		t.Fatalf("transfer fields did not update properly: %+v", updated)
	}

	// Negative tests: missing lookups
	if _, err := repo.GetStaffProfileByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for unknown staff ID")
	}
	if _, err := repo.GetStaffProfileByEmail(ctx, "unknown@nesco.ng"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for unknown email")
	}
	if _, err := repo.GetStaffProfileByStaffNumber(ctx, "UNKNOWN"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for unknown staff number")
	}
	if err := repo.UpdateStaffTransfer(ctx, "unknown", "org-2", domain.RoleMD, "hmac"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound updating unknown staff")
	}
}

func TestMemoryRepository_GridAssetAndNotes(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	parentPlantID := "plant-parent-1"
	asset := domain.GridAsset{
		ID:               "ast-1",
		ParentAssetID:    &parentPlantID,
		AssetCode:        "PLANT-KURRA-01",
		AssetType:        domain.AssetPlant,
		OrgUnitID:        "org-plant-1",
		RatedCapacityKVA: 12000.0,
		GPSPoint:         domain.GPSPoint{Latitude: 9.42, Longitude: 8.92},
		OperationalState: domain.StateInService,
		RowHMACSignature: "sig-ast-1",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := repo.CreateGridAsset(ctx, asset); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}

	// Duplicate code
	if err := repo.CreateGridAsset(ctx, asset); err != repository.ErrDuplicateKey {
		t.Fatalf("expected ErrDuplicateKey for asset")
	}

	byCode, err := repo.GetGridAssetByCode(ctx, "PLANT-KURRA-01")
	if err != nil || byCode.ID != asset.ID {
		t.Fatalf("failed to fetch asset by code: %v", err)
	}

	byID, err := repo.GetGridAssetByID(ctx, asset.ID)
	if err != nil || byID.AssetCode != asset.AssetCode {
		t.Fatalf("failed to fetch asset by ID: %v", err)
	}

	// List assets with filters
	assets, err := repo.ListGridAssets(ctx, parentPlantID, domain.AssetPlant)
	if err != nil || len(assets) != 1 {
		t.Fatalf("expected 1 plant asset in list, got: %d", len(assets))
	}
	assetsMismatch, err := repo.ListGridAssets(ctx, "other-parent", domain.AssetPlant)
	if err != nil || len(assetsMismatch) != 0 {
		t.Fatalf("expected 0 assets with mismatched parent")
	}
	assetsTypeMismatch, err := repo.ListGridAssets(ctx, parentPlantID, domain.AssetFeeder11KV)
	if err != nil || len(assetsTypeMismatch) != 0 {
		t.Fatalf("expected 0 assets with mismatched type")
	}

	// Missing asset checks
	if _, err := repo.GetGridAssetByID(ctx, "unknown"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for unknown asset ID")
	}
	if _, err := repo.GetGridAssetByCode(ctx, "UNKNOWN"); err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound for unknown asset code")
	}

	// WorkNotes tests
	latest, err := repo.GetLatestWorkNote(ctx, "ASSET", asset.ID)
	if err != nil || latest != nil {
		t.Fatalf("expected nil for empty note stream")
	}

	note1 := domain.WorkNote{
		ID:            "note-1",
		EntityType:    "ASSET",
		EntityID:      asset.ID,
		AuthorStaffID: "stf-1",
		NoteText:      "Inspected turbine bearings",
		PrevHash:      "0000000000000000000000000000000000000000000000000000000000000000",
		CurrentHash:   "hash-note-1",
		CreatedAt:     time.Now(),
	}

	if err := repo.CreateWorkNote(ctx, note1); err != nil {
		t.Fatalf("failed to create work note: %v", err)
	}

	latest, err = repo.GetLatestWorkNote(ctx, "ASSET", asset.ID)
	if err != nil || latest.ID != note1.ID {
		t.Fatalf("expected latest note to match note1: %v", err)
	}

	stream, err := repo.ListWorkNotesByEntity(ctx, "ASSET", asset.ID)
	if err != nil || len(stream) != 1 {
		t.Fatalf("expected 1 note in stream, got %d", len(stream))
	}
}

