// Package domain implements the priority feeder rationing algorithm,
// bulk intake validation, and 5-step safety isolation state machine guards.
package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var (
	// ErrNegativeBulkMW indicates bulk energy input cannot be negative.
	ErrNegativeBulkMW = errors.New("domain: bulk generation or TCN power intake cannot be negative")
	// ErrInvalidPTWParams indicates mandatory PTW request parameters are missing.
	ErrInvalidPTWParams = errors.New("domain: PTW request requires valid asset_id and unit_head_id")
	// ErrIllegalPTWTransition indicates an unauthorized or out-of-sequence PTW state transition.
	ErrIllegalPTWTransition = errors.New("domain: illegal PTW state transition sequence")
	// ErrMissingIsolationProof indicates isolation requires notes confirming breaker open and line earthed.
	ErrMissingIsolationProof = errors.New("domain: isolation requires explicit confirmation notes (min 10 chars) confirming breaker open and earth applied")
	// ErrMissingClearanceNotes indicates safety clearance requires notes confirming crew is safe and clear.
	ErrMissingClearanceNotes = errors.New("domain: safety clearance requires explicit declaration (min 10 chars) confirming crew is clear")
)

// CalculateBulkMW sums available generation and TCN intake with physical sanity checks.
func CalculateBulkMW(generationMW, tcnMW float64) (float64, error) {
	if generationMW < 0.0 || tcnMW < 0.0 {
		return 0.0, ErrNegativeBulkMW
	}
	return generationMW + tcnMW, nil
}

// ComputePriorityFeederRationing allocates available bulk MW across 33kV and 11kV feeders
// ranked strictly by Priority Rank (Rank 1: Band A / High Revenue -> Rank 5: Band E).
// Feeders that cannot be supported within available intake are shed (OUT_ON_LOAD_SHEDDING).
func ComputePriorityFeederRationing(availableBulkMW float64, feeders []FeederStatus) []RationingDecision {
	// 1. Sort feeders by Priority Rank ascending (1 is highest priority)
	sorted := make([]FeederStatus, len(feeders))
	copy(sorted, feeders)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].PriorityRank < sorted[j].PriorityRank
	})

	var decisions []RationingDecision
	accumulatedMW := 0.0

	for _, f := range sorted {
		// Feeders currently under PTW or tripped on fault cannot be energized
		if f.State == FeederOutOnPTW || f.State == FeederOutOnFault {
			decisions = append(decisions, RationingDecision{
				FeederCode:    f.FeederCode,
				AllocatedMW:   0.0,
				TargetState:   f.State,
				PriorityRank:  f.PriorityRank,
				Justification: fmt.Sprintf("Retained in %s: cannot energize under safety permit or active fault.", f.State),
			})
			continue
		}

		load := f.CurrentLoadMW
		if load <= 0.0 {
			load = 1.0 // Nominal baseline allocation
		}

		if accumulatedMW+load <= availableBulkMW {
			accumulatedMW += load
			decisions = append(decisions, RationingDecision{
				FeederCode:    f.FeederCode,
				AllocatedMW:   load,
				TargetState:   FeederInService,
				PriorityRank:  f.PriorityRank,
				Justification: fmt.Sprintf("Energized under Priority Rank %d (Band %s). Cumulative demand: %.2f MW / %.2f MW available.", f.PriorityRank, f.TariffBand, accumulatedMW, availableBulkMW),
			})
		} else {
			decisions = append(decisions, RationingDecision{
				FeederCode:    f.FeederCode,
				AllocatedMW:   0.0,
				TargetState:   FeederOutOnLoadShedding,
				PriorityRank:  f.PriorityRank,
				Justification: fmt.Sprintf("Load shed: insufficient bulk MW (available: %.2f MW, cumulative required: %.2f MW).", availableBulkMW, accumulatedMW+load),
			})
		}
	}

	return decisions
}

// ValidatePTWRequest verifies mandatory fields when Step 1 of PTW is initiated.
func ValidatePTWRequest(assetID, unitHeadID string) error {
	if strings.TrimSpace(assetID) == "" || strings.TrimSpace(unitHeadID) == "" {
		return ErrInvalidPTWParams
	}
	return nil
}

// ValidateDualApproval (Step 2) updates digital sign-offs from Head of Services AND HSE Safety Rep.
// The status transitions to APPROVED_PENDING_ISOLATION only once BOTH parties have approved.
func ValidateDualApproval(ptw *SafetyPermitPTW, isDeptHead, isSafetyRep bool, approverID string) (PTWStatus, error) {
	if ptw.Status != PTWRequested && ptw.Status != PTWApprovedPendingIsolation {
		return ptw.Status, fmt.Errorf("%w: cannot approve permit in status %s", ErrIllegalPTWTransition, ptw.Status)
	}

	now := time.Now().UTC()
	if isDeptHead {
		ptw.DeptHeadApproved = true
		ptw.DeptHeadID = &approverID
		ptw.DeptHeadApprovedAt = &now
	}
	if isSafetyRep {
		ptw.SafetyRepApproved = true
		ptw.SafetyRepID = &approverID
		ptw.SafetyRepApprovedAt = &now
	}

	// Dual-approval threshold check
	if ptw.DeptHeadApproved && ptw.SafetyRepApproved {
		ptw.Status = PTWApprovedPendingIsolation
	}

	return ptw.Status, nil
}

// ValidateSubstationIsolation (Step 3) checks authorization before applying physical earth/isolation.
func ValidateSubstationIsolation(currentStatus PTWStatus, operatorID, isolationNotes string) (PTWStatus, error) {
	if currentStatus != PTWApprovedPendingIsolation {
		return currentStatus, fmt.Errorf("%w: can only grant isolation for permits in APPROVED_PENDING_ISOLATION (current: %s)", ErrIllegalPTWTransition, currentStatus)
	}
	if strings.TrimSpace(operatorID) == "" {
		return currentStatus, fmt.Errorf("%w: substation operator ID is mandatory", ErrIllegalPTWTransition)
	}
	if len(strings.TrimSpace(isolationNotes)) < 10 {
		return currentStatus, ErrMissingIsolationProof
	}
	return PTWIsolatedEarthed, nil
}

// ValidateSafetyClearance (Step 5a) verifies that HSE has declared the site safe and crews clear.
func ValidateSafetyClearance(currentStatus PTWStatus, safetyRepID, notes string) (PTWStatus, error) {
	if currentStatus != PTWIsolatedEarthed {
		return currentStatus, fmt.Errorf("%w: can only clear safety for permits in ISOLATED_EARTHED (current: %s)", ErrIllegalPTWTransition, currentStatus)
	}
	if strings.TrimSpace(safetyRepID) == "" {
		return currentStatus, fmt.Errorf("%w: safety representative ID is mandatory", ErrIllegalPTWTransition)
	}
	if len(strings.TrimSpace(notes)) < 10 {
		return currentStatus, ErrMissingClearanceNotes
	}
	return PTWSafetyCleared, nil
}

// ValidateRestoration (Step 5b) verifies that the substation operator can re-energize the breaker.
func ValidateRestoration(currentStatus PTWStatus, operatorID string) (PTWStatus, error) {
	if currentStatus != PTWSafetyCleared {
		return currentStatus, fmt.Errorf("%w: can only restore power after permit is in SAFETY_CLEARED status (current: %s)", ErrIllegalPTWTransition, currentStatus)
	}
	if strings.TrimSpace(operatorID) == "" {
		return currentStatus, fmt.Errorf("%w: substation operator ID is mandatory", ErrIllegalPTWTransition)
	}
	return PTWRestoredClosed, nil
}

