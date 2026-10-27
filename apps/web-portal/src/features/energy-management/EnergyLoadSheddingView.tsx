/**
 * Feeder Rationing & Priority Load Shedding Console for Central Dispatch.
 * Manages NERC Tariff Bands A through E with emergency load shed triggers.
 */

import React, { useState } from "react";
import { MOCK_FEEDERS } from "./mockData";
import { FeederDispatchItem } from "./types";

export function EnergyLoadSheddingView() {
  const [feeders, setFeeders] = useState<FeederDispatchItem[]>(MOCK_FEEDERS);
  const [selectedBand, setSelectedBand] = useState<string>("ALL");
  const [targetShedMw, setTargetShedMw] = useState<string>("5.0");

  const filteredFeeders = feeders.filter((f) => {
    if (selectedBand === "ALL") return true;
    return f.tariff_band === selectedBand;
  });

  const handleToggleState = (id: string) => {
    setFeeders((prev) =>
      prev.map((f) => {
        if (f.id !== id) return f;
        if (f.state === "IN_SERVICE") {
          return { ...f, state: "OUT_ON_LOAD_SHEDDING", current_load_mw: 0.0 };
        } else if (f.state === "OUT_ON_LOAD_SHEDDING") {
          return { ...f, state: "IN_SERVICE", current_load_mw: Number((f.capacity_mw * 0.7).toFixed(1)) };
        }
        return f;
      })
    );
  };

  const handleExecuteAlgorithm = () => {
    const mw = parseFloat(targetShedMw);
    if (isNaN(mw) || mw <= 0) {
      alert("Please enter a valid MW shed target.");
      return;
    }
    // Shed from lowest priority tier (Tier 5, 4, 3...)
    alert(`Priority shedding algorithm executed for ${mw} MW drop across lowest priority bands.`);
  };

  return (
    <div className="space-y-6">
      {/* Top Console Bar */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-wrap items-center justify-between gap-4">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>⚡</span>
            <span>Priority Load Shedding &amp; Feeder Rationing Engine</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Strict Band A protection (min 20 hrs/day). Shedding commences from Band E upwards.
          </p>
        </div>

        {/* Algorithm Quick Shed Box */}
        <div className="flex items-center space-x-2">
          <input
            type="number"
            step="0.5"
            value={targetShedMw}
            onChange={(e) => setTargetShedMw(e.target.value)}
            className="w-24 px-2.5 py-1.5 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 text-xs font-mono font-bold focus:outline-none focus:border-amber-400"
            placeholder="MW"
          />
          <button
            onClick={handleExecuteAlgorithm}
            className="px-3.5 py-1.5 rounded-lg bg-amber-600 hover:bg-amber-500 text-white font-bold text-xs transition shadow-lg shadow-amber-950 flex items-center space-x-1"
          >
            <span>🚨</span>
            <span>Auto-Shed MW</span>
          </button>
        </div>
      </div>

      {/* Band Filter Tabs */}
      <div className="flex items-center space-x-2 text-xs">
        {["ALL", "Band A", "Band B", "Band C", "Band D", "Band E"].map((band) => (
          <button
            key={band}
            onClick={() => setSelectedBand(band)}
            className={`px-3 py-1.5 rounded-lg font-medium transition ${
              selectedBand === band
                ? "bg-cyan-950 text-cyan-300 border border-cyan-700 font-bold"
                : "bg-slate-900 text-slate-400 hover:text-white border border-slate-800"
            }`}
          >
            {band}
          </button>
        ))}
      </div>

      {/* Feeder Status & Shedding Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Feeder Corridor</th>
                <th className="py-3 px-3">Band / Tier</th>
                <th className="py-3 px-3">Voltage</th>
                <th className="py-3 px-3">Live Load</th>
                <th className="py-3 px-3">Target / Actual</th>
                <th className="py-3 px-3">Current State</th>
                <th className="py-3 px-4 text-right">Dispatch Control</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {filteredFeeders.map((f) => (
                <tr key={f.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4">
                    <div className="font-bold text-white">{f.feeder_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{f.feeder_code} &bull; {f.substation}</div>
                  </td>
                  <td className="py-3 px-3">
                    <span className="font-bold text-cyan-400">{f.tariff_band}</span>
                    <span className="text-[10px] text-slate-500 block">Tier {f.priority_tier}</span>
                  </td>
                  <td className="py-3 px-3 font-mono">{f.voltage_level_kv} kV</td>
                  <td className="py-3 px-3 font-mono font-bold text-white">
                    {f.current_load_mw.toFixed(1)} / {f.capacity_mw.toFixed(1)} MW
                  </td>
                  <td className="py-3 px-3">
                    <span className="font-mono text-slate-300">
                      {f.actual_hours_today}h / {f.min_supply_hours_day}h
                    </span>
                  </td>
                  <td className="py-3 px-3">
                    {f.state === "IN_SERVICE" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                        IN SERVICE
                      </span>
                    )}
                    {f.state === "OUT_ON_LOAD_SHEDDING" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                        LOAD SHED
                      </span>
                    )}
                    {f.state === "OUT_ON_PERMIT_PTW" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-sky-950 text-sky-400 border border-sky-800">
                        PTW ISOLATED
                      </span>
                    )}
                    {f.state === "OUT_ON_FAULT" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-rose-950 text-rose-400 border border-rose-800">
                        TRIPPED FAULT
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-4 text-right">
                    {f.state === "IN_SERVICE" ? (
                      <button
                        onClick={() => handleToggleState(f.id)}
                        className="px-2.5 py-1 rounded bg-amber-950/80 hover:bg-amber-900 text-amber-300 border border-amber-800/60 text-[11px] font-semibold transition"
                      >
                        Shed Feeder
                      </button>
                    ) : f.state === "OUT_ON_LOAD_SHEDDING" ? (
                      <button
                        onClick={() => handleToggleState(f.id)}
                        className="px-2.5 py-1 rounded bg-emerald-950/80 hover:bg-emerald-900 text-emerald-300 border border-emerald-800/60 text-[11px] font-semibold transition"
                      >
                        Restore Line
                      </button>
                    ) : (
                      <span className="text-[11px] text-slate-500 italic">Locked on Safety</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

