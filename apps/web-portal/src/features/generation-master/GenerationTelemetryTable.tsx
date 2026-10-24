/**
 * Real-Time Database Feed & Hourly Telemetry Log Table.
 * Live audit trail of hourly readings submitted into generation_master.master_station_telemetry_logs.
 */

import React from "react";
import { HydroPlantSummary } from "./types";
import { Clock } from "lucide-react";

interface GenerationTelemetryTableProps {
  plants: HydroPlantSummary[];
}

export const GenerationTelemetryTable: React.FC<GenerationTelemetryTableProps> = ({
  plants,
}) => {
  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-3">
      <div className="flex items-center justify-between pb-2 border-b border-slate-800">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <Clock className="w-4 h-4 text-cyan-400" />
            <span>Real-Time Database Feed &bull; Hourly Telemetry Log</span>
          </h3>
          <p className="text-[11px] text-slate-400 mt-0.5">
            Live audit trail of hourly readings submitted into{" "}
            <code className="text-cyan-300 font-mono">
              generation_master.master_station_telemetry_logs
            </code>.
          </p>
        </div>
        <span className="text-[10px] text-slate-400 font-mono">
          {plants.length} Active Feeds
        </span>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs">
          <thead className="bg-slate-950/80 text-slate-400 uppercase text-[10px] font-semibold border-b border-slate-800">
            <tr>
              <th className="p-2.5">Plant Station</th>
              <th className="p-2.5">Status</th>
              <th className="p-2.5">Active Power</th>
              <th className="p-2.5">Water Level</th>
              <th className="p-2.5">Stator Temp</th>
              <th className="p-2.5">Vibration</th>
              <th className="p-2.5">Lube Oil</th>
              <th className="p-2.5">Last Logged</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/60 font-mono">
            {plants.map((p) => (
              <tr key={p.id} className="hover:bg-slate-800/40 transition">
                <td className="p-2.5 font-sans font-medium text-slate-200">{p.name}</td>
                <td className="p-2.5">
                  <span
                    className={`px-1.5 py-0.5 rounded text-[10px] font-bold ${
                      p.state === "IN_SERVICE"
                        ? "bg-emerald-500/20 text-emerald-300"
                        : p.state === "STANDBY"
                        ? "bg-slate-700 text-slate-300"
                        : "bg-rose-500/20 text-rose-300"
                    }`}
                  >
                    {p.state}
                  </span>
                </td>
                <td className="p-2.5 font-bold text-cyan-300">{p.current_mw.toFixed(1)} MW</td>
                <td className="p-2.5 text-blue-300">{p.water_level_m} m</td>
                <td className="p-2.5 text-slate-300">{p.stator_temp_c}°C</td>
                <td
                  className={`p-2.5 font-bold ${
                    p.bearing_vibration_mm_s >= 4.5 ? "text-amber-400" : "text-emerald-400"
                  }`}
                >
                  {p.bearing_vibration_mm_s} mm/s
                </td>
                <td className="p-2.5 text-slate-300">{p.lube_oil_pressure_bar} bar</td>
                <td className="p-2.5 text-slate-400 text-[10px]">
                  {new Date(p.last_collated_at).toLocaleTimeString("en-GB", {
                    hour: "2-digit",
                    minute: "2-digit",
                  })}{" "}
                  WAT
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

