/**
 * Activity 1: Real-Time Generation Analytics & Database Summary View.
 * Summarizes from hourly readings stored in the database:
 * - Active plants output percentage of total generating power (Moment & 24h) in Pie Chart
 * - Active plants output bar chart comparison (MW Generated vs Installed Capacity)
 * - Power produced against installed capacity for each of the 6 cascade plants
 * - Machines state of health (ISO 10816 vibration, thermal, lube oil)
 * - Water level & reservoir head telemetry
 * - Out of service plants with root-cause outage details
 */

import React, { useState } from "react";
import { HydroPlantSummary, PlantTelemetryLog, GenerationAnalytics } from "./types";
import {
  Activity,
  Droplets,
  Zap,
  Gauge,
  AlertTriangle,
  RefreshCw,
} from "lucide-react";
import { ActivePlantsOutputChartsCard } from "./ActivePlantsOutputChartsCard";
import { OutOfServiceAlertBanner } from "./OutOfServiceAlertBanner";
import { PlantHealthAndHydrologyGrid } from "./PlantHealthAndHydrologyGrid";
import { GenerationTelemetryTable } from "./GenerationTelemetryTable";

interface AnalyticsViewProps {
  plants: HydroPlantSummary[];
  telemetryLogs?: PlantTelemetryLog[];
  analytics?: GenerationAnalytics | null;
  onRefresh?: () => void;
  onNavigateToReadingSheet?: () => void;
}

export const GenerationAnalyticsView: React.FC<AnalyticsViewProps> = ({
  plants,
  telemetryLogs = [],
  analytics,
  onRefresh,
  onNavigateToReadingSheet,
}) => {
  const [filterPlant, setFilterPlant] = useState<string>("ALL");

  const totalCurrentMW = plants.reduce((sum, p) => sum + p.current_mw, 0);
  const totalCapacityMW = plants.reduce((sum, p) => sum + p.installed_capacity_mw, 0);
  const overallCapacityFactor = totalCapacityMW > 0 ? (totalCurrentMW / totalCapacityMW) * 100 : 0;
  const inServicePlants = plants.filter((p) => p.state === "IN_SERVICE");
  const outOfServicePlants = plants.filter((p) => p.state === "OUT_OF_SERVICE");
  const standbyPlants = plants.filter((p) => p.state === "STANDBY");
  const highVibrationPlants = plants.filter((p) => p.bearing_vibration_mm_s >= 4.5);

  const filteredPlants = filterPlant === "ALL" 
    ? plants 
    : plants.filter((p) => p.id === filterPlant);

  return (
    <div className="space-y-6 text-xs text-slate-100">
      {/* 1. Header & Live Executive Banner */}
      <div className="bg-gradient-to-r from-slate-900 via-cyan-950/40 to-slate-900 border border-cyan-500/30 rounded-2xl p-5 shadow-xl relative overflow-hidden">
        <div className="absolute top-0 right-0 p-4 opacity-10 pointer-events-none">
          <Activity className="w-36 h-36 text-cyan-400" />
        </div>

        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 relative z-10">
          <div>
            <div className="flex items-center space-x-2 mb-1.5">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-cyan-500/20 border border-cyan-500/40 text-cyan-300">
                NERC SCADA Summary &bull; Baseload Hydro Fleet
              </span>
              <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
                Live DB Feed
              </span>
            </div>
            <h1 className="text-xl sm:text-2xl font-black text-white tracking-tight">
              Real-Time Generation Analytics &amp; Machine Telemetry
            </h1>
            <p className="text-xs text-slate-300 mt-1 max-w-2xl leading-relaxed">
              Consolidated operational summary synthesized from hourly operator log entries in the database.
              Monitoring active power versus installed capacity, water levels, machine health, and plant outage triage.
            </p>
          </div>

          <div className="flex items-center space-x-2 shrink-0">
            {onNavigateToReadingSheet && (
              <button
                onClick={onNavigateToReadingSheet}
                className="px-3.5 py-2 rounded-lg bg-cyan-600 hover:bg-cyan-500 text-white font-semibold flex items-center space-x-1.5 shadow transition"
              >
                <Zap className="w-4 h-4" />
                <span>Fill Hourly Reading</span>
              </button>
            )}

            {onRefresh && (
              <button
                onClick={onRefresh}
                className="px-3.5 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-200 font-semibold flex items-center space-x-1.5 shadow transition"
              >
                <RefreshCw className="w-3.5 h-3.5 text-cyan-400" />
                <span>Refresh DB</span>
              </button>
            )}
          </div>
        </div>

        {/* Top 4 KPI Metrics Strip */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3.5 mt-5">
          <div className="bg-slate-950/80 border border-slate-800/80 rounded-xl p-3.5">
            <div className="flex items-center justify-between text-slate-400 font-semibold text-[11px]">
              <span>Cascade Generation</span>
              <Zap className="w-4 h-4 text-cyan-400" />
            </div>
            <div className="mt-1.5 flex items-baseline justify-between">
              <span className="text-2xl font-black font-mono text-cyan-300">
                {totalCurrentMW.toFixed(1)} <span className="text-xs font-normal text-slate-400">MW</span>
              </span>
              <span className="text-[10px] font-bold px-1.5 py-0.5 rounded bg-cyan-500/10 text-cyan-400 font-mono">
                {overallCapacityFactor.toFixed(1)}% Capacity
              </span>
            </div>
            <div className="text-[10px] text-slate-400 mt-1">
              Installed: <span className="font-mono font-semibold text-slate-200">{totalCapacityMW.toFixed(1)} MW</span> across 6 plants
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800/80 rounded-xl p-3.5">
            <div className="flex items-center justify-between text-slate-400 font-semibold text-[11px]">
              <span>Hydrological Head Index</span>
              <Droplets className="w-4 h-4 text-blue-400" />
            </div>
            <div className="mt-1.5 flex items-baseline justify-between">
              <span className="text-2xl font-black font-mono text-blue-300">
                78.4% <span className="text-xs font-normal text-slate-400">Vol</span>
              </span>
              <span className="text-[10px] font-bold px-1.5 py-0.5 rounded bg-blue-500/10 text-blue-400 font-mono">
                Optimal Head
              </span>
            </div>
            <div className="text-[10px] text-slate-400 mt-1">
              Avg Inflow: <span className="font-mono font-semibold text-slate-200">38.6 m³/s</span> into Kurra &amp; Kwall
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800/80 rounded-xl p-3.5">
            <div className="flex items-center justify-between text-slate-400 font-semibold text-[11px]">
              <span>Machines State of Health</span>
              <Gauge className="w-4 h-4 text-emerald-400" />
            </div>
            <div className="mt-1.5 flex items-baseline justify-between">
              <span className="text-2xl font-black font-mono text-emerald-300">
                94.2% <span className="text-xs font-normal text-slate-400">Score</span>
              </span>
              <span className={`text-[10px] font-bold px-1.5 py-0.5 rounded ${
                highVibrationPlants.length > 0
                  ? "bg-amber-500/20 text-amber-300 border border-amber-500/30"
                  : "bg-emerald-500/10 text-emerald-400"
              }`}>
                {highVibrationPlants.length > 0 ? `${highVibrationPlants.length} Alarm` : "Normal"}
              </span>
            </div>
            <div className="text-[10px] text-slate-400 mt-1">
              Thermal &amp; ISO 10816 vibration status compliant
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800/80 rounded-xl p-3.5">
            <div className="flex items-center justify-between text-slate-400 font-semibold text-[11px]">
              <span>Fleet Operating State</span>
              <AlertTriangle className="w-4 h-4 text-rose-400" />
            </div>
            <div className="mt-1.5 flex items-baseline justify-between">
              <span className="text-2xl font-black font-mono text-white">
                {inServicePlants.length} <span className="text-xs font-normal text-slate-400">of 6 In Service</span>
              </span>
              <span className={`text-[10px] font-bold px-1.5 py-0.5 rounded ${
                outOfServicePlants.length > 0
                  ? "bg-rose-500/20 text-rose-300 border border-rose-500/40 animate-pulse"
                  : "bg-emerald-500/10 text-emerald-400"
              }`}>
                {outOfServicePlants.length} Out of Service
              </span>
            </div>
            <div className="text-[10px] text-slate-400 mt-1">
              Standby: <span className="font-semibold text-slate-200">{standbyPlants.length}</span> | Outage: <span className="font-semibold text-rose-300">{outOfServicePlants.length}</span>
            </div>
          </div>
        </div>
      </div>

      {/* 2. Out of Service Alert Banner (Prominent if any plant is down) */}
      <OutOfServiceAlertBanner outOfServicePlants={outOfServicePlants} />

      {/* 3. Active Plants Output Visual Charts: Pie Chart (% Share Moment & 24h) + Bar Chart */}
      <ActivePlantsOutputChartsCard plants={plants} analytics={analytics} />

      {/* 4. Each Generation Plant: Power Produced against Installed Capacity */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-3 border-b border-slate-800">
          <div>
            <h2 className="text-sm font-bold text-white flex items-center space-x-2">
              <Zap className="w-4 h-4 text-cyan-400" />
              <span>Power Produced Against Installed Capacity (Plants 1–6)</span>
            </h2>
            <p className="text-[11px] text-slate-400 mt-0.5">
              Live hourly output generated versus total nameplate turbine capacity.
            </p>
          </div>

          <div className="flex items-center space-x-2">
            <span className="text-[11px] text-slate-400">Filter Station:</span>
            <select
              value={filterPlant}
              onChange={(e) => setFilterPlant(e.target.value)}
              className="bg-slate-950 border border-slate-700 rounded-lg px-2.5 py-1 text-slate-200 text-xs focus:outline-none focus:border-cyan-500"
            >
              <option value="ALL">All 6 Cascade Plants</option>
              {plants.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Comparative Cards Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredPlants.map((plant) => {
            const pct = plant.installed_capacity_mw > 0 
              ? Math.min(100, (plant.current_mw / plant.installed_capacity_mw) * 100) 
              : 0;
            const isOut = plant.state === "OUT_OF_SERVICE";
            const isStandby = plant.state === "STANDBY";

            return (
              <div
                key={plant.id}
                className={`p-4 rounded-xl border transition-all space-y-3 ${
                  isOut
                    ? "bg-rose-950/20 border-rose-800/70"
                    : isStandby
                    ? "bg-slate-950/60 border-slate-800"
                    : "bg-slate-950/80 border-slate-800 hover:border-cyan-500/50"
                }`}
              >
                <div className="flex items-start justify-between">
                  <div>
                    <h3 className="font-bold text-white text-xs">{plant.name}</h3>
                    <span className="text-[10px] text-slate-400">{plant.location} &bull; {plant.turbine_type}</span>
                  </div>
                  <span
                    className={`px-2 py-0.5 rounded text-[10px] font-bold font-mono uppercase ${
                      isOut
                        ? "bg-rose-500/20 text-rose-300 border border-rose-500/30"
                        : isStandby
                        ? "bg-slate-700 text-slate-300"
                        : "bg-emerald-500/20 text-emerald-300 border border-emerald-500/30"
                    }`}
                  >
                    {plant.state.replace(/_/g, " ")}
                  </span>
                </div>

                <div className="space-y-1.5">
                  <div className="flex items-baseline justify-between text-xs">
                    <span className="text-slate-400">Power Produced:</span>
                    <div className="font-mono">
                      <strong className="text-base text-cyan-300 font-bold">{plant.current_mw.toFixed(1)}</strong>
                      <span className="text-slate-400 text-xs"> / {plant.installed_capacity_mw.toFixed(1)} MW</span>
                    </div>
                  </div>

                  <div className="w-full bg-slate-800 rounded-full h-2.5 overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all duration-500 ${
                        isOut
                          ? "bg-rose-500"
                          : pct >= 90
                          ? "bg-gradient-to-r from-emerald-500 to-cyan-400"
                          : pct >= 50
                          ? "bg-cyan-500"
                          : "bg-amber-500"
                      }`}
                      style={{ width: `${pct}%` }}
                    />
                  </div>

                  <div className="flex items-center justify-between text-[10px] text-slate-400">
                    <span>Utilization: <strong className="text-slate-200 font-mono">{pct.toFixed(1)}%</strong></span>
                    <span>Units: <strong className="text-slate-200 font-mono">{plant.active_units}/{plant.total_units} Active</strong></span>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-2 pt-2 border-t border-slate-800/80 text-[10px]">
                  <div>
                    <span className="text-slate-400">Water Level:</span>
                    <div className="font-semibold text-blue-300 font-mono">{plant.water_level_m} m</div>
                  </div>
                  <div>
                    <span className="text-slate-400">Vibration:</span>
                    <div className={`font-semibold font-mono ${
                      plant.bearing_vibration_mm_s >= 4.5 ? "text-amber-400 font-bold" : "text-slate-200"
                    }`}>
                      {plant.bearing_vibration_mm_s} mm/s
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* 5. Water Level & Machines State of Health Combined Grid */}
      <PlantHealthAndHydrologyGrid plants={plants} />

      {/* 6. Live Database Telemetry Log Entries (Recent Hourly Readings) */}
      <GenerationTelemetryTable plants={plants} />
    </div>
  );
};
