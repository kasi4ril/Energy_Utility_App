// Package service coordinates predictive maintenance prognostics and automated work orders.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	s04domain "disco-enterprise-platform/services/core-api/internal/s04_services_maint/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
)

// HydroPrognosisRequest encapsulates water reservoir telemetry.
type HydroPrognosisRequest struct {
	AssetID           string  `json:"asset_id"`
	InflowM3s         float64 `json:"inflow_m3s"`
	DischargeM3s      float64 `json:"discharge_m3s"`
	CurrentHeadM      float64 `json:"current_head_m"`
	TurbineEfficiency float64 `json:"turbine_efficiency"`
	ForecastDays      int     `json:"forecast_days"`
}

// GeneratorPrognosisRequest encapsulates mechanical and thermal telemetry for a turbine generator.
type GeneratorPrognosisRequest struct {
	AssetID               string  `json:"asset_id"`
	BearingVibrationMmS   float64 `json:"bearing_vibration_mm_s"`
	WindingTempC          float64 `json:"winding_temp_c"`
	LubeOilPressureBar    float64 `json:"lube_oil_pressure_bar"`
	RunHoursSinceOverhaul float64 `json:"run_hours_since_overhaul"`
	UnitHeadID            string  `json:"unit_head_id"`
	TeamLeadID            string  `json:"team_lead_id"`
}

// TransformerPrognosisRequest encapsulates loading and top oil thermal telemetry.
type TransformerPrognosisRequest struct {
	AssetID         string  `json:"asset_id"`
	LoadPct         float64 `json:"load_pct"`
	TopOilTempC     float64 `json:"top_oil_temp_c"`
	FuseBlows30Days int     `json:"fuse_blows_30_days"`
	UnitHeadID      string  `json:"unit_head_id"`
	TeamLeadID      string  `json:"team_lead_id"`
}

// FeederPrognosisRequest encapsulates line fault density and bush patrol history.
type FeederPrognosisRequest struct {
	AssetID             string  `json:"asset_id"`
	TripCount30Days     int     `json:"trip_count_30_days"`
	FaultDensityPerKm   float64 `json:"fault_density_per_km"`
	DaysSinceLastPatrol float64 `json:"days_since_last_patrol"`
	StormForecast       bool    `json:"storm_forecast"`
	UnitHeadID          string  `json:"unit_head_id"`
	TeamLeadID          string  `json:"team_lead_id"`
}

// EvaluateHydroPrognosis evaluates water availability and generates a predictive trigger.
func (s *WorkflowService) EvaluateHydroPrognosis(ctx context.Context, req HydroPrognosisRequest) (*domain.PredictiveMaintenanceTrigger, error) {
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, fmt.Errorf("service: asset_id is required")
	}

	head, riskScore, riskLevel, mwh, summary := domain.EvaluateHydroWaterAvailability(
		req.InflowM3s, req.DischargeM3s, req.CurrentHeadM, req.TurbineEfficiency, req.ForecastDays,
	)

	now := time.Now().UTC()
	payload, _ := json.Marshal(map[string]interface{}{
		"forecasted_head_m": head,
		"forecasted_mwh":    mwh,
		"inflow_m3s":        req.InflowM3s,
		"discharge_m3s":     req.DischargeM3s,
	})

	trigger := domain.PredictiveMaintenanceTrigger{
		ID:                   generateUUID(),
		TriggerNo:            fmt.Sprintf("PMT-HYD-%d-%s", now.Unix(), randomHexString(3)),
		AssetID:              req.AssetID,
		AssetType:            domain.AssetHydroReservoir,
		ModelName:            domain.ModelHydroWaterAvailability,
		RiskScore:            riskScore,
		RiskLevel:            riskLevel,
		TelemetrySummaryJSON: string(payload),
		RecommendedAction:    summary,
		Status:               "ACTIVE",
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := s.repo.CreatePredictiveTrigger(ctx, trigger); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "executive.predictive.trigger_raised", "s09_executive_command", map[string]interface{}{
			"trigger_id": trigger.ID,
			"trigger_no": trigger.TriggerNo,
			"asset_id":   trigger.AssetID,
			"risk_level": trigger.RiskLevel,
		})
	}

	return &trigger, nil
}

// EvaluateGeneratorPrognosis evaluates generator physics and auto-creates a maintenance work order on critical risks.
func (s *WorkflowService) EvaluateGeneratorPrognosis(ctx context.Context, req GeneratorPrognosisRequest) (*domain.PredictiveMaintenanceTrigger, error) {
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, fmt.Errorf("service: asset_id is required")
	}

	riskScore, riskLevel, action, autoWO := domain.EvaluateGeneratorHealth(
		req.BearingVibrationMmS, req.WindingTempC, req.LubeOilPressureBar, req.RunHoursSinceOverhaul,
	)

	now := time.Now().UTC()
	payload, _ := json.Marshal(map[string]interface{}{
		"bearing_vibration_mm_s":   req.BearingVibrationMmS,
		"winding_temp_c":           req.WindingTempC,
		"lube_oil_pressure_bar":    req.LubeOilPressureBar,
		"run_hours_since_overhaul": req.RunHoursSinceOverhaul,
	})

	trigger := domain.PredictiveMaintenanceTrigger{
		ID:                   generateUUID(),
		TriggerNo:            fmt.Sprintf("PMT-GEN-%d-%s", now.Unix(), randomHexString(3)),
		AssetID:              req.AssetID,
		AssetType:            domain.AssetGenerator,
		ModelName:            domain.ModelGeneratorPreventive,
		RiskScore:            riskScore,
		RiskLevel:            riskLevel,
		TelemetrySummaryJSON: string(payload),
		RecommendedAction:    action,
		AutoWorkOrderCreated: autoWO,
		Status:               "ACTIVE",
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if autoWO && s.maintRepo != nil {
		unitHead := req.UnitHeadID
		if unitHead == "" {
			unitHead = "HEAD-GEN-ELEC-01"
		}
		teamLead := req.TeamLeadID
		if teamLead == "" {
			teamLead = "LEAD-GEN-TECH-01"
		}
		wo := s04domain.MaintenanceWorkOrder{
			ID:            generateUUID(),
			OrderNo:       fmt.Sprintf("WO-PRED-%d-%s", now.Unix(), randomHexString(3)),
			UnitType:      s04domain.UnitTransformerMaint,
			WorkOrderType: s04domain.WOPreventiveMaintenance,
			Priority:      s04domain.PriorityCritical,
			UnitHeadID:    unitHead,
			TeamLeadID:    teamLead,
			AssetID:       req.AssetID,
			PTWRequired:   true,
			Status:        s04domain.StatusDraft,
			Instructions:  fmt.Sprintf("[AI PREDICTIVE CBM TRIGGER %s] %s", trigger.TriggerNo, action),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.maintRepo.CreateWorkOrder(ctx, wo); err == nil {
			trigger.LinkedWorkOrderID = &wo.ID
		}
	}

	if err := s.repo.CreatePredictiveTrigger(ctx, trigger); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = s.bus.Publish(ctx, "executive.predictive.trigger_raised", "s09_executive_command", map[string]interface{}{
			"trigger_id": trigger.ID,
			"trigger_no": trigger.TriggerNo,
			"asset_id":   trigger.AssetID,
			"risk_level": trigger.RiskLevel,
		})
	}

	return &trigger, nil
}

// EvaluateTransformerPrognosis evaluates transformer CBM metrics and auto-creates work orders if needed.
func (s *WorkflowService) EvaluateTransformerPrognosis(ctx context.Context, req TransformerPrognosisRequest) (*domain.PredictiveMaintenanceTrigger, error) {
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, fmt.Errorf("service: asset_id is required")
	}

	riskScore, riskLevel, action, autoWO := domain.EvaluateTransformerCBM(
		req.LoadPct, req.TopOilTempC, req.FuseBlows30Days,
	)

	now := time.Now().UTC()
	payload, _ := json.Marshal(map[string]interface{}{
		"load_pct":           req.LoadPct,
		"top_oil_temp_c":     req.TopOilTempC,
		"fuse_blows_30_days": req.FuseBlows30Days,
	})

	trigger := domain.PredictiveMaintenanceTrigger{
		ID:                   generateUUID(),
		TriggerNo:            fmt.Sprintf("PMT-TX-%d-%s", now.Unix(), randomHexString(3)),
		AssetID:              req.AssetID,
		AssetType:            domain.AssetTransformer,
		ModelName:            domain.ModelTransformerCBM,
		RiskScore:            riskScore,
		RiskLevel:            riskLevel,
		TelemetrySummaryJSON: string(payload),
		RecommendedAction:    action,
		AutoWorkOrderCreated: autoWO,
		Status:               "ACTIVE",
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if autoWO && s.maintRepo != nil {
		unitHead := req.UnitHeadID
		if unitHead == "" {
			unitHead = "HEAD-TX-MAINT-01"
		}
		teamLead := req.TeamLeadID
		if teamLead == "" {
			teamLead = "LEAD-TX-CREW-01"
		}
		wo := s04domain.MaintenanceWorkOrder{
			ID:            generateUUID(),
			OrderNo:       fmt.Sprintf("WO-PRED-%d-%s", now.Unix(), randomHexString(3)),
			UnitType:      s04domain.UnitTransformerMaint,
			WorkOrderType: s04domain.WODTInspection,
			Priority:      s04domain.PriorityHigh,
			UnitHeadID:    unitHead,
			TeamLeadID:    teamLead,
			AssetID:       req.AssetID,
			Status:        s04domain.StatusDraft,
			Instructions:  fmt.Sprintf("[AI TRANSFORMER CBM TRIGGER %s] %s", trigger.TriggerNo, action),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.maintRepo.CreateWorkOrder(ctx, wo); err == nil {
			trigger.LinkedWorkOrderID = &wo.ID
		}
	}

	if err := s.repo.CreatePredictiveTrigger(ctx, trigger); err != nil {
		return nil, err
	}

	return &trigger, nil
}

// EvaluateFeederPrognosis evaluates feeder trip probability and schedules patrol work orders.
func (s *WorkflowService) EvaluateFeederPrognosis(ctx context.Context, req FeederPrognosisRequest) (*domain.PredictiveMaintenanceTrigger, error) {
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, fmt.Errorf("service: asset_id is required")
	}

	riskScore, riskLevel, action, autoWO := domain.EvaluateFeederPatrolRisk(
		req.TripCount30Days, req.FaultDensityPerKm, req.DaysSinceLastPatrol, req.StormForecast,
	)

	now := time.Now().UTC()
	payload, _ := json.Marshal(map[string]interface{}{
		"trip_count_30_days":     req.TripCount30Days,
		"fault_density_per_km":   req.FaultDensityPerKm,
		"days_since_last_patrol": req.DaysSinceLastPatrol,
		"storm_forecast":         req.StormForecast,
	})

	trigger := domain.PredictiveMaintenanceTrigger{
		ID:                   generateUUID(),
		TriggerNo:            fmt.Sprintf("PMT-FDR-%d-%s", now.Unix(), randomHexString(3)),
		AssetID:              req.AssetID,
		AssetType:            domain.AssetFeeder,
		ModelName:            domain.ModelFeederPatrol,
		RiskScore:            riskScore,
		RiskLevel:            riskLevel,
		TelemetrySummaryJSON: string(payload),
		RecommendedAction:    action,
		AutoWorkOrderCreated: autoWO,
		Status:               "ACTIVE",
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if autoWO && s.maintRepo != nil {
		unitHead := req.UnitHeadID
		if unitHead == "" {
			unitHead = "HEAD-FDR-MAINT-01"
		}
		teamLead := req.TeamLeadID
		if teamLead == "" {
			teamLead = "LEAD-FDR-PATROL-01"
		}
		wo := s04domain.MaintenanceWorkOrder{
			ID:            generateUUID(),
			OrderNo:       fmt.Sprintf("WO-PRED-%d-%s", now.Unix(), randomHexString(3)),
			UnitType:      s04domain.UnitFeederMaint,
			WorkOrderType: s04domain.WOPreventiveMaintenance,
			Priority:      s04domain.PriorityHigh,
			UnitHeadID:    unitHead,
			TeamLeadID:    teamLead,
			AssetID:       req.AssetID,
			Status:        s04domain.StatusDraft,
			Instructions:  fmt.Sprintf("[AI FEEDER PATROL TRIGGER %s] %s", trigger.TriggerNo, action),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.maintRepo.CreateWorkOrder(ctx, wo); err == nil {
			trigger.LinkedWorkOrderID = &wo.ID
		}
	}

	if err := s.repo.CreatePredictiveTrigger(ctx, trigger); err != nil {
		return nil, err
	}

	return &trigger, nil
}
