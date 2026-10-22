// Package domain contains pure business rules and validation logic for the
// Foundation bounded context. It contains no database or HTTP dependencies.
package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"disco-enterprise-platform/services/core-api/pkg/crypto"
)

var (
	// ErrInvalidOrgHierarchy is returned when an organizational unit's parent-child relationship is invalid.
	ErrInvalidOrgHierarchy = errors.New("domain: invalid organizational unit hierarchy")
	// ErrInvalidGridTopology is returned when an electrical asset violates transmission/distribution order.
	ErrInvalidGridTopology = errors.New("domain: invalid electrical grid asset topology")
	// ErrInvalidGPSCoordinates indicates coordinates fall outside valid latitude/longitude bounds.
	ErrInvalidGPSCoordinates = errors.New("domain: invalid GPS coordinates")
	// ErrInvalidStaffData indicates missing or ill-formed staff profile data.
	ErrInvalidStaffData = errors.New("domain: invalid staff profile data")
	// ErrSuspendedStaffOnDuty indicates suspended staff cannot be put on duty.
	ErrSuspendedStaffOnDuty = errors.New("domain: staff with non-active certification cannot be placed on duty")
	// ErrBrokenWorkNoteChain indicates a hash chain continuity failure in WorkNotes.
	ErrBrokenWorkNoteChain = errors.New("domain: cryptographic work note hash chain mismatch")
)

// ValidateGPSPoint validates that latitude and longitude conform to standard spherical bounds
// and fall within reasonable operational proximity.
func ValidateGPSPoint(gps GPSPoint) error {
	if gps.Latitude < -90.0 || gps.Latitude > 90.0 {
		return fmt.Errorf("%w: latitude %f must be between -90 and +90", ErrInvalidGPSCoordinates, gps.Latitude)
	}
	if gps.Longitude < -180.0 || gps.Longitude > 180.0 {
		return fmt.Errorf("%w: longitude %f must be between -180 and +180", ErrInvalidGPSCoordinates, gps.Longitude)
	}
	return nil
}

// ValidateOrgUnitHierarchy enforces the NESCO administrative tree:
// HQ (root) -> DISTRICT -> SERVICE_UNIT; HQ -> MASTER_STATION -> GENERATING_PLANT.
func ValidateOrgUnitHierarchy(unitType OrgUnitType, parentType *OrgUnitType) error {
	switch unitType {
	case OrgUnitHQ:
		if parentType != nil {
			return fmt.Errorf("%w: HQ must be the root organizational unit without a parent", ErrInvalidOrgHierarchy)
		}
	case OrgUnitDistrict:
		if parentType == nil || *parentType != OrgUnitHQ {
			return fmt.Errorf("%w: DISTRICT parent must be HQ", ErrInvalidOrgHierarchy)
		}
	case OrgUnitServiceUnit:
		if parentType == nil || *parentType != OrgUnitDistrict {
			return fmt.Errorf("%w: SERVICE_UNIT parent must be a DISTRICT", ErrInvalidOrgHierarchy)
		}
	case OrgUnitMasterStation:
		if parentType == nil || *parentType != OrgUnitHQ {
			return fmt.Errorf("%w: MASTER_STATION parent must be HQ", ErrInvalidOrgHierarchy)
		}
	case OrgUnitGeneratingPlant:
		if parentType == nil || *parentType != OrgUnitMasterStation {
			return fmt.Errorf("%w: GENERATING_PLANT parent must be MASTER_STATION", ErrInvalidOrgHierarchy)
		}
	default:
		return fmt.Errorf("%w: unrecognized unit type %s", ErrInvalidOrgHierarchy, unitType)
	}
	return nil
}

// ValidateGridTopology verifies parent-child hierarchy across the electrical network:
// PLANT -> SUBSTATION_33KV -> FEEDER_11KV -> TRANSFORMER_DT -> SERVICE_POLE -> CUSTOMER_METER.
func ValidateGridTopology(childType AssetType, parentType *AssetType) error {
	switch childType {
	case AssetPlant:
		// Plants are root generation nodes (or future TCN bulk intake)
		if parentType != nil && *parentType != AssetPlant {
			return fmt.Errorf("%w: PLANT can only attach to another bulk supply node", ErrInvalidGridTopology)
		}
	case AssetSubstation33KV:
		if parentType == nil || *parentType != AssetPlant {
			return fmt.Errorf("%w: SUBSTATION_33KV must feed from a PLANT or bulk intake", ErrInvalidGridTopology)
		}
	case AssetFeeder11KV:
		if parentType == nil || *parentType != AssetSubstation33KV {
			return fmt.Errorf("%w: FEEDER_11KV must radiate from a SUBSTATION_33KV", ErrInvalidGridTopology)
		}
	case AssetTransformerDT:
		if parentType == nil || *parentType != AssetFeeder11KV {
			return fmt.Errorf("%w: TRANSFORMER_DT must connect to a FEEDER_11KV", ErrInvalidGridTopology)
		}
	case AssetServicePole:
		if parentType == nil || *parentType != AssetTransformerDT {
			return fmt.Errorf("%w: SERVICE_POLE must connect to a TRANSFORMER_DT", ErrInvalidGridTopology)
		}
	case AssetCustomerMeter:
		if parentType == nil || (*parentType != AssetServicePole && *parentType != AssetTransformerDT) {
			return fmt.Errorf("%w: CUSTOMER_METER must connect to a SERVICE_POLE or DT", ErrInvalidGridTopology)
		}
	default:
		return fmt.Errorf("%w: unrecognized asset type %s", ErrInvalidGridTopology, childType)
	}
	return nil
}

// ValidateStaffProfile verifies personal attributes, email syntax, and role validity.
func ValidateStaffProfile(profile StaffProfile) error {
	if strings.TrimSpace(profile.StaffNumber) == "" {
		return fmt.Errorf("%w: staff_number cannot be empty", ErrInvalidStaffData)
	}
	if strings.TrimSpace(profile.FullName) == "" {
		return fmt.Errorf("%w: full_name cannot be empty", ErrInvalidStaffData)
	}
	if _, err := mail.ParseAddress(profile.Email); err != nil {
		return fmt.Errorf("%w: invalid email address: %v", ErrInvalidStaffData, err)
	}
	if profile.OrgUnitID == "" {
		return fmt.Errorf("%w: org_unit_id is required", ErrInvalidStaffData)
	}

	// Certification guard: non-active cert cannot be placed on duty
	if profile.ActiveShiftState == ShiftOnDuty && profile.CertificationStatus != CertActive {
		return ErrSuspendedStaffOnDuty
	}

	return nil
}

// ValidateWorkNoteIntegrity checks the SHA-256 cryptographic chain of a work note.
func ValidateWorkNoteIntegrity(note WorkNote, expectedPrevHash string) error {
	if strings.TrimSpace(note.EntityType) == "" || strings.TrimSpace(note.EntityID) == "" {
		return errors.New("domain: work note entity_type and entity_id are required")
	}
	if strings.TrimSpace(note.AuthorStaffID) == "" {
		return errors.New("domain: author_staff_id is required")
	}
	if strings.TrimSpace(note.NoteText) == "" {
		return errors.New("domain: note_text cannot be empty")
	}

	if expectedPrevHash == "" {
		expectedPrevHash = crypto.GenesisHash
	}
	if note.PrevHash != expectedPrevHash {
		return fmt.Errorf("%w: expected prev_hash %s but got %s", ErrBrokenWorkNoteChain, expectedPrevHash, note.PrevHash)
	}

	recomputed := crypto.ComputeWorkNoteHash(note.PrevHash, note.EntityType, note.EntityID, note.AuthorStaffID, note.NoteText, note.CreatedAt.Unix())
	if recomputed != note.CurrentHash {
		return fmt.Errorf("%w: current_hash mismatch, recomputed %s vs stored %s", ErrBrokenWorkNoteChain, recomputed, note.CurrentHash)
	}

	return nil
}

