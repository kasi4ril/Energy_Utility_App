// Package domain implements the predictive maintenance evaluation algorithms,
// hydro capacity forecasting, and C-Suite approval gate validations for S09.
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrUnauthorizedExecutiveRole indicates the deciding officer lacks the required C-Suite role.
	ErrUnauthorizedExecutiveRole = errors.New("domain: unauthorized executive role for requested action")
	// ErrIllegalActionItemState indicates an invalid state transition on an executive action item.
	ErrIllegalActionItemState = errors.New("domain: action item is not in pending state")
	// ErrInvalidDecisionParams indicates missing decision notes or deciding staff ID.
	ErrInvalidDecisionParams = errors.New("domain: deciding_staff_id and notes are required")
)

// EvaluateHydroWaterAvailability forecasts reservoir head, risk score, risk level, and available MWh energy capacity.
func EvaluateHydroWaterAvailability(inflowM3s, dischargeM3s, currentHeadM, turbineEfficiency float64, forecastDays int) (float64, float64, RiskLevel, float64, string) {
	if forecastDays <= 0 {
		forecastDays = 7
	}
	if turbineEfficiency <= 0.0 || turbineEfficiency > 1.0 {
		turbineEfficiency = 0.88 // Standard Francis turbine efficiency
	}

	// Approximate reservoir surface area: 10,000,000 m2 (10 sq km)
	const reservoirAreaM2 = 10000000.0
	netFlowM3s := inflowM3s - dischargeM3s
	dailyDeltaHeadM := (netFlowM3s * 86400.0) / reservoirAreaM2
	forecastedHeadM := currentHeadM + (dailyDeltaHeadM * float64(forecastDays))

	if forecastedHeadM < 0.0 {
		forecastedHeadM = 0.0
	}

	// Power P (MW) = (rho * g * Q * H * eta) / 10^6 where rho = 1000 kg/m3, g = 9.81 m/s2
	activeFlow := dischargeM3s
	if activeFlow > inflowM3s && forecastedHeadM < 10.0 {
		activeFlow = inflowM3s // Reduced dispatch due to depleted head
	}
	powerMW := (1000.0 * 9.81 * activeFlow * forecastedHeadM * turbineEfficiency) / 1000000.0
	forecastedMWh := powerMW * 24.0 * float64(forecastDays)

	riskScore := 0.10
	riskLevel := RiskLow
	summary := fmt.Sprintf("Stable hydro forecast: Head projected at %.1fm with %.1f MWh over %d days.", forecastedHeadM, forecastedMWh, forecastDays)

	if forecastedHeadM < 15.0 {
		riskScore = 0.85
		riskLevel = RiskCritical
		summary = fmt.Sprintf("Critical Low Water Alert: Reservoir head dropping to %.1fm within %d days! Imminent generation rationing required.", forecastedHeadM, forecastDays)
	} else if forecastedHeadM < 22.0 || netFlowM3s < -15.0 {
		riskScore = 0.65
		riskLevel = RiskHigh
		summary = fmt.Sprintf("Hydro Depletion Warning: Discharge exceeds inflow by %.1f m3/s. Head forecast at %.1fm.", -netFlowM3s, forecastedHeadM)
	} else if netFlowM3s < 0.0 {
		riskScore = 0.35
		riskLevel = RiskMedium
		summary = fmt.Sprintf("Moderate hydro drawdown: Inflow %.1f m3/s < Discharge %.1f m3/s.", inflowM3s, dischargeM3s)
	}

	return forecastedHeadM, riskScore, riskLevel, forecastedMWh, summary
}

// EvaluateGeneratorHealth calculates mechanical/electrical failure probability from sensor telemetry.
func EvaluateGeneratorHealth(bearingVibrationMmS, windingTempC, lubeOilPressureBar, runHoursSinceOverhaul float64) (float64, RiskLevel, string, bool) {
	var penalty float64

	// 1. Bearing Vibration (ISO 10816: standard <= 4.5 mm/s; alert > 7.1; critical trip > 11.2)
	if bearingVibrationMmS > 11.2 {
		penalty += 0.50
	} else if bearingVibrationMmS > 7.1 {
		penalty += 0.35
	} else if bearingVibrationMmS > 4.5 {
		penalty += 0.15
	}

	// 2. Stator Winding Temp (Normal <= 90C; warning > 105C; critical > 120C)
	if windingTempC > 120.0 {
		penalty += 0.40
	} else if windingTempC > 105.0 {
		penalty += 0.25
	} else if windingTempC > 90.0 {
		penalty += 0.10
	}

	// 3. Lube Oil Pressure (Normal >= 2.5 bar; low < 1.8; critical trip < 1.2)
	if lubeOilPressureBar < 1.2 && lubeOilPressureBar > 0.0 {
		penalty += 0.45
	} else if lubeOilPressureBar < 1.8 && lubeOilPressureBar > 0.0 {
		penalty += 0.25
	}

	// 4. Overhaul hours (Reconditioned threshold 8,000 hrs)
	if runHoursSinceOverhaul > 8000.0 {
		penalty += 0.20
	}

	riskScore := penalty
	if riskScore > 1.0 {
		riskScore = 1.0
	}

	riskLevel := RiskLow
	autoWorkOrder := false
	action := "Generator telemetry within acceptable operating parameters."

	if riskScore >= 0.75 {
		riskLevel = RiskCritical
		autoWorkOrder = true
		action = "Immediate Emergency Shutdown & Overhaul Recommended: Dangerous vibration/thermal exceedance detected."
	} else if riskScore >= 0.50 {
		riskLevel = RiskHigh
		autoWorkOrder = true
		action = "Preventive Mechanical Intervention Scheduled: Bearing lube flushing and vibration analysis required."
	} else if riskScore >= 0.25 {
		riskLevel = RiskMedium
		action = "Enhanced Shift Monitoring: Winding temperature and lube pressure trending upward."
	}

	return riskScore, riskLevel, action, autoWorkOrder
}

// EvaluateTransformerCBM predicts DT and Power Transformer failure from loading, top oil temp, and fuse blows.
func EvaluateTransformerCBM(loadPct, topOilTempC float64, fuseBlows30Days int) (float64, RiskLevel, string, bool) {
	var penalty float64

	// Load exceedance
	if loadPct > 120.0 {
		penalty += 0.50
	} else if loadPct > 100.0 {
		penalty += 0.30
	} else if loadPct > 85.0 {
		penalty += 0.15
	}

	// Top oil temperature
	if topOilTempC > 95.0 {
		penalty += 0.40
	} else if topOilTempC > 85.0 {
		penalty += 0.25
	}

	// Fuse blow frequency in last 30 days
	if fuseBlows30Days >= 4 {
		penalty += 0.35
	} else if fuseBlows30Days >= 2 {
		penalty += 0.15
	}

	riskScore := penalty
	if riskScore > 1.0 {
		riskScore = 1.0
	}

	riskLevel := RiskLow
	autoWorkOrder := false
	action := "Transformer operating within safe thermal and loading envelope."

	if riskScore >= 0.70 {
		riskLevel = RiskCritical
		autoWorkOrder = true
		action = "Immediate Transformer Load Relief & Oil Degassing Required: Severe thermal overload detected."
	} else if riskScore >= 0.45 {
		riskLevel = RiskHigh
		autoWorkOrder = true
		action = "Schedule Dielectric Breakdown Voltage (BDV) Test and Load Balancing across LV phases."
	} else if riskScore >= 0.20 {
		riskLevel = RiskMedium
		action = "Monitor peak evening loading: Transformer approaching rated capacity."
	}

	return riskScore, riskLevel, action, autoWorkOrder
}

// EvaluateFeederPatrolRisk predicts feeder line trips and pole collapses from trip history and vegetation patrol cadence.
func EvaluateFeederPatrolRisk(tripCount30Days int, faultDensityPerKm, daysSinceLastPatrol float64, stormForecast bool) (float64, RiskLevel, string, bool) {
	var penalty float64

	if tripCount30Days >= 8 {
		penalty += 0.40
	} else if tripCount30Days >= 4 {
		penalty += 0.20
	}

	if faultDensityPerKm >= 2.0 {
		penalty += 0.30
	} else if faultDensityPerKm >= 1.0 {
		penalty += 0.15
	}

	if daysSinceLastPatrol > 90.0 {
		penalty += 0.30
	} else if daysSinceLastPatrol > 60.0 {
		penalty += 0.15
	}

	if stormForecast {
		penalty += 0.25
	}

	riskScore := penalty
	if riskScore > 1.0 {
		riskScore = 1.0
	}

	riskLevel := RiskLow
	autoWorkOrder := false
	action := "Feeder corridor within normal reliability margins."

	if riskScore >= 0.70 {
		riskLevel = RiskCritical
		autoWorkOrder = true
		action = "Priority Bush Clearing & Insulator Wash Required: Imminent line trip risk under stormy conditions."
	} else if riskScore >= 0.45 {
		riskLevel = RiskHigh
		autoWorkOrder = true
		action = "Dispatch Route Foot-Patrol: High fault density and vegetation encroachment overdue."
	} else if riskScore >= 0.20 {
		riskLevel = RiskMedium
		action = "Schedule routine tree-trimming patrol for feeder right-of-way."
	}

	return riskScore, riskLevel, action, autoWorkOrder
}

// ValidateActionItemDecision enforces C-Suite role segregation of duties on one-click approvals.
func ValidateActionItemDecision(item *ExecutiveActionItem, decidedByStaffID string, role ExecutiveRole, approve bool, notes string) error {
	if item.Status != ActionStatusPending {
		return fmt.Errorf("%w: cannot decide item in status %s", ErrIllegalActionItemState, item.Status)
	}
	if strings.TrimSpace(decidedByStaffID) == "" || strings.TrimSpace(notes) == "" {
		return ErrInvalidDecisionParams
	}

	// Segregation of duties: verify officer role matches the target role
	if item.TargetRole != role {
		return fmt.Errorf("%w: action %s requires role %s, got %s",
			ErrUnauthorizedExecutiveRole, item.ActionType, item.TargetRole, role)
	}

	now := time.Now().UTC()
	item.DecidedByStaffID = &decidedByStaffID
	item.DecidedAt = &now
	item.DecisionNotes = &notes

	if approve {
		item.Status = ActionStatusApproved
	} else {
		item.Status = ActionStatusRejected
	}
	item.UpdatedAt = now
	return nil
}
