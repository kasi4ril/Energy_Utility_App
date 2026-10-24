/**
 * Plant Health and Hydrology Grid:
 * 1. Cascade Reservoir Water Levels & Elevation Telemetry
 * 2. Turbine & Generator Machine State of Health Scorecard (ISO 10816, Stator, Lube Oil)
 */

import React from "react";
import { HydroPlantSummary } from "./types";
import { Droplets, Gauge } from "lucide-react";

interface PlantHealthAndHydrologyGridProps {
  plants: HydroPlantSummary[];
}

export const PlantHealthAndHydrologyGrid: React.FC<PlantHealthAndHydrologyGridProps> = ({
  plants,
}) => {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* 1. Hydrological Water Level Details */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-3.5">
        <div className="flex items-center justify-between pb-2.5 border-b border-slate-800">
          <h2 className="text-sm font-bold text-white flex items-center space-x-2">
            <Droplets className="w-4 h-4 text-blue-400" />
            <span>Cascade Reservoir Water Levels &amp; Elevation</span>
          </h2>
          <span className="text-[10px] text-slate-400 font-mono">Headwater Telemetry</span>
        </div>

        <div className="space-y-3">
          {plants.map((plant) => (
            <div
              key={plant.id}
              className="p-3 bg-slate-950/70 border border-slate-800/80 rounded-lg flex items-center justify-between gap-3"
            >
              <div>
                <div className="font-semibold text-slate-200 text-xs">{plant.name}</div>
                <div className="text-[10px] text-slate-400 mt-0.5">
                  Headwater: <span className="font-mono text-slate-300">{plant.headwater_elevation_m}m</span> &bull; 
                  Tailwater: <span className="font-mono text-slate-300">{plant.tailwater_elevation_m}m</span>
                </div>
              </div>

              <div className="text-right">
                <div className="font-mono font-bold text-sm text-blue-300">
                  {plant.water_level_m.toFixed(1)} <span className="text-[10px] font-normal text-slate-400">meters</span>
                </div>
                <span className={`text-[9px] font-semibold px-1.5 py-0.2 rounded ${
                  plant.water_level_m >= 20
                    ? "bg-blue-500/20 text-blue-300"
                    : plant.water_level_m >= 15
                    ? "bg-emerald-500/20 text-emerald-300"
                    : "bg-amber-500/20 text-amber-300"
                }`}>
                  {plant.water_level_m >= 15 ? "Healthy Head" : "Low Water Gate"}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* 2. Machines State of Health Scorecard */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-3.5">
        <div className="flex items-center justify-between pb-2.5 border-b border-slate-800">
          <h2 className="text-sm font-bold text-white flex items-center space-x-2">
            <Gauge className="w-4 h-4 text-emerald-400" />
            <span>Turbine &amp; Generator Machine State of Health</span>
          </h2>
          <span className="text-[10px] text-slate-400 font-mono">ISO 10816 Vibration</span>
        </div>

        <div className="space-y-3">
          {plants.map((plant) => {
            const isVibAlert = plant.bearing_vibration_mm_s >= 4.5;
            const isTempAlert = plant.stator_temp_c >= 75.0;

            return (
              <div
                key={plant.id}
                className={`p-3 rounded-lg border flex items-center justify-between gap-3 ${
                  isVibAlert
                    ? "bg-amber-950/20 border-amber-500/50"
                    : "bg-slate-950/70 border-slate-800/80"
                }`}
              >
                <div>
                  <div className="flex items-center space-x-2">
                    <span className="font-semibold text-slate-200 text-xs">{plant.name}</span>
                    {isVibAlert && (
                      <span className="px-1.5 py-0.2 rounded text-[9px] font-bold bg-amber-500/20 text-amber-300 border border-amber-500/40">
                        Vibration Warning
                      </span>
                    )}
                  </div>
                  <div className="flex items-center space-x-3 text-[10px] text-slate-400 mt-1">
                    <span>Stator: <strong className={`font-mono ${isTempAlert ? "text-amber-300" : "text-slate-200"}`}>{plant.stator_temp_c}°C</strong></span>
                    <span>Lube Oil: <strong className="font-mono text-slate-200">{plant.lube_oil_pressure_bar} bar</strong></span>
                    <span>Vibration: <strong className={`font-mono ${isVibAlert ? "text-amber-400 font-bold" : "text-emerald-400"}`}>{plant.bearing_vibration_mm_s} mm/s</strong></span>
                  </div>
                </div>

                <div className="text-right">
                  <span
                    className={`px-2 py-0.5 rounded text-[10px] font-bold font-mono uppercase ${
                      plant.state === "OUT_OF_SERVICE"
                        ? "bg-rose-500/20 text-rose-300"
                        : isVibAlert
                        ? "bg-amber-500/20 text-amber-300"
                        : "bg-emerald-500/20 text-emerald-300"
                    }`}
                  >
                    {plant.state === "OUT_OF_SERVICE" ? "Trip / Outage" : isVibAlert ? "Attention" : "Healthy"}
                  </span>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

