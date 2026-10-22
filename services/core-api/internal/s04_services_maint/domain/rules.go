// Package domain implements pure business logic validation, spatial proximity verification,
// and state transition guards for S04 Services & Maintenance operations.
package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"

	s01domain "disco-enterprise-platform/services/core-api/internal/s01_foundation/domain"
)

var (
	// ErrEmptyPhotoHash indicates a before or after photographic proof hash is missing or malformed.
	ErrEmptyPhotoHash = errors.New("domain: before and after photo hashes must be provided (min 32 characters)")
	// ErrIdenticalPhotoHashes indicates the before and after photo hashes match, signalling duplicate photo fraud.
	ErrIdenticalPhotoHashes = errors.New("domain: anti-fraud violation: before and after photos cannot be identical")
	// ErrMissingMeterSeal indicates a Metering work order is missing the required tamper-evident seal number.
	ErrMissingMeterSeal = errors.New("domain: metering work orders require a valid meter seal number (min 4 characters)")
	// ErrShortWorkDetails indicates the technician failed to provide an adequate technical summary.
	ErrShortWorkDetails = errors.New("domain: work details must contain at least 15 characters of technical description")
	// ErrPTWPermitRequired indicates an electrically isolated job was submitted without a validated PTW permit.
	ErrPTWPermitRequired = errors.New("domain: safety violation: PTW permit is mandatory before undertaking high-voltage line work")
	// ErrInvalidStatusTransition indicates an illegal work order state transition attempt.
	ErrInvalidStatusTransition = errors.New("domain: illegal work order status transition")
	// ErrMissingReworkNotes indicates rejection was attempted without providing remedial feedback.
	ErrMissingReworkNotes = errors.New("domain: unit head must provide clear rework instructions (min 10 characters) when rejecting proof")
	// ErrInvalidWorkOrderParams indicates missing mandatory fields on work order creation.
	ErrInvalidWorkOrderParams = errors.New("domain: invalid work order parameters")
)

// GetGPSToleranceForUnit returns the maximum permissible deviation (in meters) between
// the registered grid asset and the field technician's device GPS for each engineering branch.
func GetGPSToleranceForUnit(unit UnitType) float64 {
	switch unit {
	case UnitFeederMaint:
		return 150.0 // Line span tolerance for 33kV/11kV overhead corridors
	case UnitTransformerMaint:
		return 75.0 // Substation and distribution transformer perimeter
	case UnitMetering:
		return 50.0 // Customer premises meter installation tolerance
	default:
		return 100.0
	}
}

// CalculateHaversineDistance computes the great-circle surface distance in meters between two GPS coordinates.
func CalculateHaversineDistance(p1, p2 s01domain.GPSPoint) float64 {
	const earthRadiusMeters = 6371000.0

	phi1 := p1.Latitude * math.Pi / 180.0
	phi2 := p2.Latitude * math.Pi / 180.0
	deltaPhi := (p2.Latitude - p1.Latitude) * math.Pi / 180.0
	deltaLambda := (p2.Longitude - p1.Longitude) * math.Pi / 180.0

	a := math.Sin(deltaPhi/2.0)*math.Sin(deltaPhi/2.0) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2.0)*math.Sin(deltaLambda/2.0)

	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	return earthRadiusMeters * c
}

// ValidateWorkOrderCreation verifies mandatory fields for a new work order.
func ValidateWorkOrderCreation(wo MaintenanceWorkOrder) error {
	if strings.TrimSpace(wo.OrderNo) == "" {
		return fmt.Errorf("%w: order_no is mandatory", ErrInvalidWorkOrderParams)
	}
	if strings.TrimSpace(wo.UnitHeadID) == "" || strings.TrimSpace(wo.TeamLeadID) == "" {
		return fmt.Errorf("%w: unit_head_id and team_lead_id are mandatory", ErrInvalidWorkOrderParams)
	}
	if strings.TrimSpace(wo.AssetID) == "" {
		return fmt.Errorf("%w: target asset_id is mandatory", ErrInvalidWorkOrderParams)
	}
	if len(strings.TrimSpace(wo.Instructions)) < 10 {
		return fmt.Errorf("%w: technical instructions must be at least 10 characters", ErrInvalidWorkOrderParams)
	}
	switch wo.UnitType {
	case UnitFeederMaint, UnitTransformerMaint, UnitMetering:
	default:
		return fmt.Errorf("%w: unknown unit_type %s", ErrInvalidWorkOrderParams, wo.UnitType)
	}
	return nil
}

// ValidateProofSubmission performs comprehensive domain integrity checks on field proof uploads:
// - Verifies non-empty and non-identical photo hashes.
// - Enforces meter seal presence for Metering tasks.
// - Enforces PTW linkage if safety isolation was required.
// - Computes distance and evaluates GPS validity against unit tolerance.
func ValidateProofSubmission(
	unit UnitType,
	assetGPS, capturedGPS s01domain.GPSPoint,
	beforeHash, afterHash string,
	meterSeal *string,
	workDetails string,
	ptwRequired bool,
	ptwPermitID *string,
) (float64, bool, error) {
	// 1. Photo Hash Integrity & Anti-Reuse Check
	bHash := strings.TrimSpace(beforeHash)
	aHash := strings.TrimSpace(afterHash)
	if len(bHash) < 32 || len(aHash) < 32 {
		return 0, false, ErrEmptyPhotoHash
	}
	if strings.EqualFold(bHash, aHash) {
		return 0, false, ErrIdenticalPhotoHashes
	}

	// 2. Meter Seal Check for Metering Unit
	if unit == UnitMetering {
		if meterSeal == nil || len(strings.TrimSpace(*meterSeal)) < 4 {
			return 0, false, ErrMissingMeterSeal
		}
	}

	// 3. Technical Work Summary Depth
	if len(strings.TrimSpace(workDetails)) < 15 {
		return 0, false, ErrShortWorkDetails
	}

	// 4. PTW Safety Permit Lock
	if ptwRequired {
		if ptwPermitID == nil || strings.TrimSpace(*ptwPermitID) == "" {
			return 0, false, ErrPTWPermitRequired
		}
	}

	// 5. GPS Haversine Proximity Check
	dist := CalculateHaversineDistance(assetGPS, capturedGPS)
	tolerance := GetGPSToleranceForUnit(unit)
	gpsValid := dist <= tolerance

	return dist, gpsValid, nil
}

// ValidateUnitHeadVerification executes the approval or rework transition of a submitted proof.
func ValidateUnitHeadVerification(currentStatus WorkOrderStatus, approved bool, notes string) (WorkOrderStatus, error) {
	if currentStatus != StatusProofSubmitted {
		return currentStatus, fmt.Errorf("%w: can only verify work orders in PROOF_SUBMITTED status (current: %s)", ErrInvalidStatusTransition, currentStatus)
	}

	if approved {
		return StatusAttendedVerified, nil
	}

	// Rejection requires actionable rework feedback
	if len(strings.TrimSpace(notes)) < 10 {
		return currentStatus, ErrMissingReworkNotes
	}
	return StatusRejectedRework, nil
}

