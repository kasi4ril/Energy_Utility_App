// Package transport implements predictive maintenance REST prognostic endpoints.
package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/domain"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/repository"
	"disco-enterprise-platform/services/core-api/internal/s09_executive_command/service"
)

// HandleHydroPrognosis evaluates water reservoir telemetry and returns predictive generation outlook.
func (h *HTTPHandler) HandleHydroPrognosis(w http.ResponseWriter, r *http.Request) {
	var dto HydroPrognosisDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	trigger, err := h.svc.EvaluateHydroPrognosis(r.Context(), service.HydroPrognosisRequest{
		AssetID:           dto.AssetID,
		InflowM3s:         dto.InflowM3s,
		DischargeM3s:      dto.DischargeM3s,
		CurrentHeadM:      dto.CurrentHeadM,
		TurbineEfficiency: dto.TurbineEfficiency,
		ForecastDays:      dto.ForecastDays,
	})
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, trigger)
}

// HandleGeneratorPrognosis evaluates mechanical/thermal health and raises automated work orders on critical risk.
func (h *HTTPHandler) HandleGeneratorPrognosis(w http.ResponseWriter, r *http.Request) {
	var dto GeneratorPrognosisDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	trigger, err := h.svc.EvaluateGeneratorPrognosis(r.Context(), service.GeneratorPrognosisRequest{
		AssetID:               dto.AssetID,
		BearingVibrationMmS:   dto.BearingVibrationMmS,
		WindingTempC:          dto.WindingTempC,
		LubeOilPressureBar:    dto.LubeOilPressureBar,
		RunHoursSinceOverhaul: dto.RunHoursSinceOverhaul,
		UnitHeadID:            dto.UnitHeadID,
		TeamLeadID:            dto.TeamLeadID,
	})
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, trigger)
}

// HandleTransformerPrognosis evaluates condition-based maintenance metrics for distribution transformers.
func (h *HTTPHandler) HandleTransformerPrognosis(w http.ResponseWriter, r *http.Request) {
	var dto TransformerPrognosisDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	trigger, err := h.svc.EvaluateTransformerPrognosis(r.Context(), service.TransformerPrognosisRequest{
		AssetID:         dto.AssetID,
		LoadPct:         dto.LoadPct,
		TopOilTempC:     dto.TopOilTempC,
		FuseBlows30Days: dto.FuseBlows30Days,
		UnitHeadID:      dto.UnitHeadID,
		TeamLeadID:      dto.TeamLeadID,
	})
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, trigger)
}

// HandleFeederPrognosis predicts feeder fault trip probability and triggers preventive vegetation patrols.
func (h *HTTPHandler) HandleFeederPrognosis(w http.ResponseWriter, r *http.Request) {
	var dto FeederPrognosisDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	trigger, err := h.svc.EvaluateFeederPrognosis(r.Context(), service.FeederPrognosisRequest{
		AssetID:             dto.AssetID,
		TripCount30Days:     dto.TripCount30Days,
		FaultDensityPerKm:   dto.FaultDensityPerKm,
		DaysSinceLastPatrol: dto.DaysSinceLastPatrol,
		StormForecast:       dto.StormForecast,
		UnitHeadID:          dto.UnitHeadID,
		TeamLeadID:          dto.TeamLeadID,
	})
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, trigger)
}

// HandleAcknowledgeTrigger records staff review of a predictive alarm.
func (h *HTTPHandler) HandleAcknowledgeTrigger(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.jsonError(w, http.StatusBadRequest, "missing_param", "Trigger id is required")
		return
	}

	var dto AcknowledgeTriggerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid_json", "Failed to decode request body")
		return
	}
	if err := dto.Validate(); err != nil {
		h.jsonError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	t, err := h.svc.AcknowledgePredictiveTrigger(r.Context(), id, dto.StaffID, dto.Status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.jsonError(w, http.StatusNotFound, "not_found", "Predictive trigger not found")
			return
		}
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, t)
}

// HandleListPredictiveTriggers queries alarms filtered by status and model name.
func (h *HTTPHandler) HandleListPredictiveTriggers(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	model := domain.PredictiveModelName(r.URL.Query().Get("model_name"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	triggers, err := h.svc.ListPredictiveTriggers(r.Context(), status, model, limit)
	if err != nil {
		h.jsonError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if triggers == nil {
		triggers = []domain.PredictiveMaintenanceTrigger{}
	}

	h.jsonResponse(w, http.StatusOK, triggers)
}

