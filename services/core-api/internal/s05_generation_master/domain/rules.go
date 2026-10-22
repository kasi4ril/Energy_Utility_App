// Package domain implements physical telemetry boundary validation, early fault alerts,
// and state transition guards for Generation and Workshop workflows.
package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidWaterLevel indicates the reservoir water level reading is outside physical boundaries.
	ErrInvalidWaterLevel = errors.New("domain: reservoir water level reading is physically impossible (must be 0.0 - 250.0m)")
	// ErrNegativePower indicates generation active power cannot be negative.
	ErrNegativePower = errors.New("domain: active power generated (MW) cannot be negative")
	// ErrInvalidPowerForOfflinePlant indicates active power was recorded for an out-of-service plant.
	ErrInvalidPowerForOfflinePlant = errors.New("domain: out-of-service plant cannot report active power generation greater than zero")
	// ErrInvalidRollcall indicates mandatory shift roll-call fields are missing.
	ErrInvalidRollcall = errors.New("domain: remote plant roll-call requires valid plant_id, shift_id, and at least one on-duty operator")
	// ErrInvalidFaultJob indicates mandatory plant fault attributes are missing.
	ErrInvalidFaultJob = errors.New("domain: plant fault job requires plant_id, reported_by, valid wing, and description (min 15 chars)")
	// ErrInvalidGMApprovalTransition indicates an unauthorized status transition for GM approval.
	ErrInvalidGMApprovalTransition = errors.New("domain: GM approval can only be applied to jobs in REPORTED or UNIT_HEAD_REVIEWED status")
	// ErrMissingGMRejectionNotes indicates GM rejection requires an explanation.
	ErrMissingGMRejectionNotes = errors.New("domain: GM rejection requires explanation notes (min 10 characters)")
	// ErrIllegalWorkshopTransition indicates an invalid workshop progression sequence.
	ErrIllegalWorkshopTransition = errors.New("domain: illegal workshop status transition")
)

// EvaluateTelemetryPhysics evaluates telemetry against physical hydrological and generator constraints,
// returning anomalies and early predictive maintenance warning alerts.
func EvaluateTelemetryPhysics(t MasterStationTelemetryLog) ([]string, error) {
	// 1. Water Level Boundaries (Reservoir level between 0m and 250m)
	if t.WaterLevelM < 0.0 || t.WaterLevelM > 250.0 {
		return nil, ErrInvalidWaterLevel
	}

	// 2. Power Sanity
	if t.MWGenerated < 0.0 {
		return nil, ErrNegativePower
	}
	if t.PlantState == StateOutOfService && t.MWGenerated > 0.0 {
		return nil, ErrInvalidPowerForOfflinePlant
	}

	var alerts []string

	// 3. Early Warning Thresholds
	if t.StatorWindingTempC > 110.0 {
		alerts = append(alerts, fmt.Sprintf("CRITICAL_THERMAL_ALERT: Stator winding temperature (%.1f°C) exceeds 110°C limit", t.StatorWindingTempC))
	} else if t.StatorWindingTempC > 95.0 {
		alerts = append(alerts, fmt.Sprintf("WARNING: High stator winding temperature (%.1f°C)", t.StatorWindingTempC))
	}

	if t.BearingVibrationMMS > 4.5 {
		alerts = append(alerts, fmt.Sprintf("CRITICAL_VIBRATION_ALERT: Bearing vibration (%.2f mm/s) exceeds 4.5 mm/s ISO limit", t.BearingVibrationMMS))
	} else if t.BearingVibrationMMS > 3.0 {
		alerts = append(alerts, fmt.Sprintf("WARNING: Elevated bearing vibration (%.2f mm/s)", t.BearingVibrationMMS))
	}

	if t.PlantState == StateInService && t.LubeOilPressureBar < 1.5 {
		alerts = append(alerts, fmt.Sprintf("CRITICAL_LUBE_OIL_ALERT: Lube oil pressure (%.2f bar) below 1.5 bar operational limit", t.LubeOilPressureBar))
	}

	return alerts, nil
}

// ValidateRollcallSubmission checks that remote plant physical presence was captured.
func ValidateRollcallSubmission(rc RemotePlantShiftRollcall) error {
	if strings.TrimSpace(rc.PlantID) == "" || strings.TrimSpace(rc.ShiftID) == "" {
		return ErrInvalidRollcall
	}
	if len(rc.OperatorNames) == 0 || strings.TrimSpace(rc.OperatorNames[0]) == "" {
		return fmt.Errorf("%w: at least one active operator name is required", ErrInvalidRollcall)
	}
	if len(strings.TrimSpace(rc.HandoverNotes)) < 10 {
		return fmt.Errorf("%w: handover notes must be at least 10 characters", ErrInvalidRollcall)
	}
	return nil
}

// ValidateFaultJobCreation validates the initial fault intake from plant operators.
func ValidateFaultJobCreation(job PlantFaultAndWorkshopJob) error {
	if strings.TrimSpace(job.PlantID) == "" || strings.TrimSpace(job.ReportedByOperatorID) == "" {
		return ErrInvalidFaultJob
	}
	if job.Wing != WingElectrical && job.Wing != WingMechanical {
		return fmt.Errorf("%w: wing must be ELECTRICAL or MECHANICAL", ErrInvalidFaultJob)
	}
	if len(strings.TrimSpace(job.FaultDescription)) < 15 {
		return fmt.Errorf("%w: fault description must be at least 15 characters", ErrInvalidFaultJob)
	}
	return nil
}

// ValidateGMApproval evaluates General Manager approval or rejection of material/workshop work.
func ValidateGMApproval(currentStatus PlantJobStatus, approved bool, notes string) (PlantJobStatus, error) {
	if currentStatus != JobStatusReported && currentStatus != JobStatusUnitHeadReviewed {
		return currentStatus, fmt.Errorf("%w: current status is %s", ErrInvalidGMApprovalTransition, currentStatus)
	}

	if approved {
		return JobStatusGMApproved, nil
	}

	if len(strings.TrimSpace(notes)) < 10 {
		return currentStatus, ErrMissingGMRejectionNotes
	}
	return JobStatusCancelled, nil
}

// ValidateWorkshopTransition validates permitted sequential lifecycle steps in the central workshop.
func ValidateWorkshopTransition(current, target WorkshopStatus) error {
	switch target {
	case WorkshopPendingDispatch:
		if current != WorkshopNotApplicable && current != WorkshopRejected {
			return ErrIllegalWorkshopTransition
		}
	case WorkshopMachiningInProgress, WorkshopWindingFabrication:
		if current != WorkshopPendingDispatch && current != WorkshopMachiningInProgress && current != WorkshopWindingFabrication {
			return ErrIllegalWorkshopTransition
		}
	case WorkshopQualityCheck:
		if current != WorkshopMachiningInProgress && current != WorkshopWindingFabrication {
			return ErrIllegalWorkshopTransition
		}
	case WorkshopCompletedDelivered:
		if current != WorkshopQualityCheck {
			return ErrIllegalWorkshopTransition
		}
	case WorkshopRejected:
		if current != WorkshopQualityCheck {
			return ErrIllegalWorkshopTransition
		}
	default:
		return ErrIllegalWorkshopTransition
	}
	return nil
}

