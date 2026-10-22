// Package domain_test verifies domain rules, hierarchy validation, and boundary conditions.
package domain_test

import (
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
	"disco-enterprise-platform/services/core-api/pkg/crypto"
)

func TestValidateGPSPoint_ValidAndInvalid(t *testing.T) {
	validGPS := domain.GPSPoint{Latitude: 9.9281, Longitude: 8.8921}
	if err := domain.ValidateGPSPoint(validGPS); err != nil {
		t.Fatalf("expected valid GPS, got error: %v", err)
	}

	// Boundary test: invalid latitude
	invalidLat := domain.GPSPoint{Latitude: 95.0, Longitude: 8.8921}
	if err := domain.ValidateGPSPoint(invalidLat); err == nil {
		t.Fatalf("expected error for latitude > 90, got nil")
	}

	invalidLatNeg := domain.GPSPoint{Latitude: -95.0, Longitude: 8.8921}
	if err := domain.ValidateGPSPoint(invalidLatNeg); err == nil {
		t.Fatalf("expected error for latitude < -90, got nil")
	}

	// Boundary test: invalid longitude
	invalidLon := domain.GPSPoint{Latitude: 9.9281, Longitude: 185.0}
	if err := domain.ValidateGPSPoint(invalidLon); err == nil {
		t.Fatalf("expected error for longitude > 180, got nil")
	}

	invalidLonNeg := domain.GPSPoint{Latitude: 9.9281, Longitude: -185.0}
	if err := domain.ValidateGPSPoint(invalidLonNeg); err == nil {
		t.Fatalf("expected error for longitude < -180, got nil")
	}
}

func TestValidateOrgUnitHierarchy(t *testing.T) {
	hq := domain.OrgUnitHQ
	district := domain.OrgUnitDistrict
	masterStation := domain.OrgUnitMasterStation
	serviceUnit := domain.OrgUnitServiceUnit

	// Valid tree relationships
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitHQ, nil); err != nil {
		t.Fatalf("HQ with nil parent should be valid, got: %v", err)
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitDistrict, &hq); err != nil {
		t.Fatalf("DISTRICT under HQ should be valid: %v", err)
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitServiceUnit, &district); err != nil {
		t.Fatalf("SERVICE_UNIT under DISTRICT should be valid: %v", err)
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitMasterStation, &hq); err != nil {
		t.Fatalf("MASTER_STATION under HQ should be valid: %v", err)
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitGeneratingPlant, &masterStation); err != nil {
		t.Fatalf("GENERATING_PLANT under MASTER_STATION should be valid: %v", err)
	}

	// Negative tests: invalid hierarchies
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitHQ, &district); err == nil {
		t.Fatalf("HQ must not have a parent")
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitDistrict, nil); err == nil {
		t.Fatalf("DISTRICT with nil parent must fail")
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitServiceUnit, &hq); err == nil {
		t.Fatalf("SERVICE_UNIT directly under HQ must fail")
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitMasterStation, nil); err == nil {
		t.Fatalf("MASTER_STATION with nil parent must fail")
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitGeneratingPlant, &district); err == nil {
		t.Fatalf("GENERATING_PLANT under DISTRICT must fail")
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitGeneratingPlant, nil); err == nil {
		t.Fatalf("GENERATING_PLANT with nil parent must fail")
	}
	if err := domain.ValidateOrgUnitHierarchy(domain.OrgUnitType("UNKNOWN"), &serviceUnit); err == nil {
		t.Fatalf("UNKNOWN unit type must fail")
	}
}

func TestValidateGridTopology(t *testing.T) {
	plant := domain.AssetPlant
	substation := domain.AssetSubstation33KV
	feeder := domain.AssetFeeder11KV
	dt := domain.AssetTransformerDT
	pole := domain.AssetServicePole

	// Valid relationships
	if err := domain.ValidateGridTopology(domain.AssetPlant, nil); err != nil {
		t.Fatalf("Root Plant with nil parent should be valid: %v", err)
	}
	if err := domain.ValidateGridTopology(domain.AssetPlant, &plant); err != nil {
		t.Fatalf("Plant connected to bulk Plant intake should be valid: %v", err)
	}
	if err := domain.ValidateGridTopology(domain.AssetSubstation33KV, &plant); err != nil {
		t.Fatalf("Substation under Plant should be valid: %v", err)
	}
	if err := domain.ValidateGridTopology(domain.AssetFeeder11KV, &substation); err != nil {
		t.Fatalf("Feeder under Substation should be valid: %v", err)
	}
	if err := domain.ValidateGridTopology(domain.AssetTransformerDT, &feeder); err != nil {
		t.Fatalf("DT under Feeder should be valid: %v", err)
	}
	if err := domain.ValidateGridTopology(domain.AssetServicePole, &dt); err != nil {
		t.Fatalf("Pole under DT should be valid: %v", err)
	}
	if err := domain.ValidateGridTopology(domain.AssetCustomerMeter, &pole); err != nil {
		t.Fatalf("Meter under Pole should be valid: %v", err)
	}
	if err := domain.ValidateGridTopology(domain.AssetCustomerMeter, &dt); err != nil {
		t.Fatalf("Meter directly under DT should be valid: %v", err)
	}

	// Negative tests: illegal grid skips
	if err := domain.ValidateGridTopology(domain.AssetPlant, &substation); err == nil {
		t.Fatalf("Plant connected to Substation must fail")
	}
	if err := domain.ValidateGridTopology(domain.AssetSubstation33KV, nil); err == nil {
		t.Fatalf("Substation without parent must fail")
	}
	if err := domain.ValidateGridTopology(domain.AssetFeeder11KV, nil); err == nil {
		t.Fatalf("Feeder without parent must fail")
	}
	if err := domain.ValidateGridTopology(domain.AssetTransformerDT, &substation); err == nil {
		t.Fatalf("DT directly under Substation without Feeder must fail")
	}
	if err := domain.ValidateGridTopology(domain.AssetServicePole, nil); err == nil {
		t.Fatalf("Pole without parent must fail")
	}
	if err := domain.ValidateGridTopology(domain.AssetCustomerMeter, &plant); err == nil {
		t.Fatalf("Customer Meter directly under Plant must fail")
	}
	if err := domain.ValidateGridTopology(domain.AssetType("INVALID"), &plant); err == nil {
		t.Fatalf("Invalid asset type must fail")
	}
}

func TestValidateStaffProfile(t *testing.T) {
	valid := domain.StaffProfile{
		StaffNumber:         "NSC-100",
		FullName:            "Ibrahim Bello",
		Email:               "ibrahim.bello@nesco.ng",
		Role:                domain.RoleTeamLeadFeeder,
		OrgUnitID:           "unit-123",
		CertificationStatus: domain.CertActive,
		ActiveShiftState:    domain.ShiftOnDuty,
	}

	if err := domain.ValidateStaffProfile(valid); err != nil {
		t.Fatalf("valid profile failed validation: %v", err)
	}

	// Negative: missing staff number
	missingNum := valid
	missingNum.StaffNumber = ""
	if err := domain.ValidateStaffProfile(missingNum); err == nil {
		t.Fatalf("expected error on empty staff_number")
	}

	// Negative: missing full name
	missingName := valid
	missingName.FullName = ""
	if err := domain.ValidateStaffProfile(missingName); err == nil {
		t.Fatalf("expected error on empty full_name")
	}

	// Negative: missing org unit id
	missingUnit := valid
	missingUnit.OrgUnitID = ""
	if err := domain.ValidateStaffProfile(missingUnit); err == nil {
		t.Fatalf("expected error on empty org_unit_id")
	}

	// Negative: invalid email
	invalidEmail := valid
	invalidEmail.Email = "not-an-email"
	if err := domain.ValidateStaffProfile(invalidEmail); err == nil {
		t.Fatalf("expected error on bad email")
	}

	// Negative: suspended staff on duty
	suspended := valid
	suspended.CertificationStatus = domain.CertSuspended
	suspended.ActiveShiftState = domain.ShiftOnDuty
	if err := domain.ValidateStaffProfile(suspended); err != domain.ErrSuspendedStaffOnDuty {
		t.Fatalf("expected ErrSuspendedStaffOnDuty, got: %v", err)
	}
}

func TestValidateWorkNoteIntegrity(t *testing.T) {
	now := time.Now()
	prevHash := crypto.GenesisHash
	hash := crypto.ComputeWorkNoteHash(prevHash, "TICKET", "T-100", "staff-1", "Power restored on Feeder 2", now.Unix())

	note := domain.WorkNote{
		EntityType:    "TICKET",
		EntityID:      "T-100",
		AuthorStaffID: "staff-1",
		NoteText:      "Power restored on Feeder 2",
		PrevHash:      prevHash,
		CurrentHash:   hash,
		CreatedAt:     now,
	}

	if err := domain.ValidateWorkNoteIntegrity(note, prevHash); err != nil {
		t.Fatalf("valid work note integrity check failed: %v", err)
	}

	// Negative: empty entity type
	emptyType := note
	emptyType.EntityType = ""
	if err := domain.ValidateWorkNoteIntegrity(emptyType, prevHash); err == nil {
		t.Fatalf("expected error on empty entity_type")
	}

	// Negative: empty author
	emptyAuthor := note
	emptyAuthor.AuthorStaffID = ""
	if err := domain.ValidateWorkNoteIntegrity(emptyAuthor, prevHash); err == nil {
		t.Fatalf("expected error on empty author_staff_id")
	}

	// Negative: empty note text
	emptyText := note
	emptyText.NoteText = ""
	if err := domain.ValidateWorkNoteIntegrity(emptyText, prevHash); err == nil {
		t.Fatalf("expected error on empty note_text")
	}

	// Negative: tampered hash
	note.CurrentHash = "tampered-hash-64-characters-0000000000000000000000000000000000000000"
	if err := domain.ValidateWorkNoteIntegrity(note, prevHash); err == nil {
		t.Fatalf("expected error on tampered hash")
	}

	// Negative: broken chain prev_hash
	note.CurrentHash = hash
	if err := domain.ValidateWorkNoteIntegrity(note, "wrong-prev-hash"); err == nil {
		t.Fatalf("expected error on wrong prev_hash")
	}
}

