/**
 * Task View 1: 6-Plant Cascade Fleet Telemetry & Live Synoptic Board.
 * Displays real-time operational status, MW generation, water head, and vibration telemetry.
 */

import React from "react";
import { HydroPlantSummary } from "./types";
import { Badge } from "../../components/ui/Badge";
import { Card } from "../../components/ui/Card";

interface FleetBoardProps {
  plants: HydroPlantSummary[];
  onSelectPlant: (plant: HydroPlantSummary) => void;
  onNavigateToCollation: (plantId: string) => void;
  onNavigateToFaultReport: (plantId: string) => void;
}

export const PlantFleetLiveBoard: React.FC<FleetBoardProps> = ({
  plants,
  onSelectPlant,
  onNavigateToCollation,
  onNavigateToFaultReport,
}) => {
  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-800 gap-3">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-bold text-white tracking-wide">
              🌐 6-Plant Hydroelectric Cascade Fleet &bull; Live Synoptic Telemetry
            </h2>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-cyan-500/10 text-cyan-300 border border-cyan-500/20">
              Real-Time SCADA
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Live generation status across Kurra Falls, Kwall, Anoma, and Jere cascade stations.
          </p>
        </div>
      </div>

      {/* 6-Plant Visual Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        {plants.map((p) => {
          const isVibrationAlarm = p.bearing_vibration_mm_s >= 4.5;
          const isOut = p.state === "OUT_OF_SERVICE";

          return (
            <div
              key={p.id}
              className={`p-5 rounded-xl border transition-all duration-200 cursor-pointer ${
                isOut
                  ? "bg-rose-950/20 border-rose-800/60 hover:border-rose-600"
                  : isVibrationAlarm
                  ? "bg-amber-950/20 border-amber-500/60 hover:border-amber-400"
                  : "bg-slate-950/70 border-slate-800 hover:border-cyan-500/50"
              }`}
              onClick={() => onSelectPlant(p)}
            >
              {/* Plant Header */}
              <div className="flex items-start justify-between pb-3 border-b border-slate-800/80">
                <div>
                  <h3 className="text-sm font-bold text-white tracking-wide">{p.name}</h3>
                  <p className="text-[11px] text-slate-400">{p.location}</p>
                </div>
                <Badge
                  variant={
                    p.state === "IN_SERVICE"
                      ? "success"
                      : p.state === "STANDBY"
                      ? "warning"
                      : "danger"
                  }
                >
                  {p.state.replace(/_/g, " ")}
                </Badge>
              </div>

              {/* Turbine Specs */}
              <div className="py-2.5 flex items-center justify-between text-[11px] border-b border-slate-900">
                <span className="text-slate-400">Turbine Type:</span>
                <span className="font-semibold text-cyan-400 font-mono">
                  {p.turbine_type} ({p.active_units}/{p.total_units} Active)
                </span>
              </div>

              {/* Core Telemetry Grid */}
              <div className="grid grid-cols-2 gap-3 py-3 text-xs">
                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Active Generation</span>
                  <div className="text-lg font-extrabold text-white font-mono mt-0.5">
                    {p.current_mw} <span className="text-xs text-slate-400 font-normal">/ {p.installed_capacity_mw} MW</span>
                  </div>
                </div>

                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Water Level / Head</span>
                  <div className="text-lg font-bold text-cyan-300 font-mono mt-0.5">
                    {p.water_level_m} m
                  </div>
                </div>

                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Bearing Vibration</span>
                  <div
                    className={`text-sm font-bold font-mono mt-0.5 ${
                      isVibrationAlarm ? "text-rose-400 animate-pulse font-extrabold" : "text-slate-300"
                    }`}
                  >
                    {p.bearing_vibration_mm_s} mm/s {isVibrationAlarm && "⚠️ ALARM"}
                  </div>
                </div>

                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Stator Temperature</span>
                  <div className="text-sm font-bold font-mono text-slate-300 mt-0.5">
                    {p.stator_temp_c}°C
                  </div>
                </div>
              </div>

              {p.outage_reason && (
                <div className="p-2 rounded bg-rose-500/10 border border-rose-500/30 text-[11px] text-rose-300 mb-2">
                  <span className="font-bold">Outage Reason:</span> {p.outage_reason.replace(/_/g, " ")}
                </div>
              )}

              {/* Quick Actions */}
              <div className="pt-2 border-t border-slate-900 flex gap-2">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onNavigateToCollation(p.id);
                  }}
                  className="flex-1 py-1.5 rounded bg-slate-800 hover:bg-slate-700 text-cyan-300 text-[11px] font-semibold border border-slate-700 transition"
                >
                  📻 Call Collate
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onNavigateToFaultReport(p.id);
                  }}
                  className="flex-1 py-1.5 rounded bg-slate-800 hover:bg-slate-700 text-rose-300 text-[11px] font-semibold border border-slate-700 transition"
                >
                  🚨 Report Fault
                </button>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
