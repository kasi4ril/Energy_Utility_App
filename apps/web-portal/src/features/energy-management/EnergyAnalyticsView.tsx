/**
 * Energy Analytics Dashboard for Central Transmission Dispatch.
 * Visualizes Dual-Source Bulk Intake (Hydro + TCN) and 33kV Substation Loading.
 */

import React from "react";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip as RechartsTooltip,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
} from "recharts";
import { MOCK_BULK_SOURCES, MOCK_FEEDERS } from "./mockData";

export function EnergyAnalyticsView() {
  const totalBulk = MOCK_BULK_SOURCES.reduce((acc, s) => acc + s.output_mw, 0);
  const hydroBulk = MOCK_BULK_SOURCES.filter((s) => s.source_type === "HYDRO").reduce(
    (acc, s) => acc + s.output_mw,
    0
  );
  const tcnBulk = MOCK_BULK_SOURCES.find((s) => s.source_type === "TCN_GRID")?.output_mw || 0;

  const pieData = MOCK_BULK_SOURCES.map((s) => ({
    name: s.name,
    value: s.output_mw,
    color: s.color,
  }));

  const barData = MOCK_FEEDERS.map((f) => ({
    name: f.feeder_code,
    load: f.current_load_mw,
    capacity: f.capacity_mw,
    fullName: f.feeder_name,
  }));

  return (
    <div className="space-y-6">
      {/* Top Grid KPI Ribbon */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Total Bulk Intake</div>
          <div className="text-2xl font-black text-cyan-400 mt-1">{totalBulk.toFixed(1)} MW</div>
          <div className="text-[10px] text-slate-500 mt-0.5">
            Hydro: {hydroBulk.toFixed(1)} MW &bull; TCN: {tcnBulk.toFixed(1)} MW
          </div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Grid Frequency</div>
          <div className="text-2xl font-black text-emerald-400 mt-1">50.02 Hz</div>
          <div className="text-[10px] text-emerald-400 mt-0.5">&plusmn;0.04% (Nominal Target 50.0Hz)</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">33kV Bus Voltage</div>
          <div className="text-2xl font-black text-white mt-1">33.2 kV</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Phase-to-Phase Balanced</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">System Power Factor</div>
          <div className="text-2xl font-black text-purple-400 mt-1">0.94 Lag</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Optimal Reactive Compensation</div>
        </div>
      </div>

      {/* Interactive Charts: Bulk Intake Donut + Feeder Loading Bar */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Bulk Intake Donut Pie Chart */}
        <div className="lg:col-span-5 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-bold text-white">Bulk Power Intake Distribution</h3>
              <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-cyan-950/60 text-cyan-400 border border-cyan-800/40">
                Live MW Share
              </span>
            </div>
            <p className="text-xs text-slate-400 mt-0.5">
              Hydrological fleet generation vs National Grid TCN bulk hook.
            </p>
          </div>

          <div className="h-64 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={88}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {pieData.map((entry, idx) => (
                    <Cell key={`cell-${idx}`} fill={entry.color} />
                  ))}
                </Pie>
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      const pct = ((data.value / totalBulk) * 100).toFixed(1);
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.name}</div>
                          <div className="text-cyan-400 font-mono">
                            {data.value} MW ({pct}%)
                          </div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>

          <div className="grid grid-cols-2 gap-2 text-[10px] border-t border-slate-800/80 pt-3">
            {pieData.map((s) => (
              <div key={s.name} className="flex items-center space-x-2">
                <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ backgroundColor: s.color }} />
                <span className="text-slate-300 truncate">{s.name}:</span>
                <span className="font-mono font-bold text-white ml-auto">{s.value} MW</span>
              </div>
            ))}
          </div>
        </div>

        {/* Feeder Loading Bar Chart */}
        <div className="lg:col-span-7 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-bold text-white">Feeder Corridor Loading (MW)</h3>
              <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-slate-800 text-slate-300 border border-slate-700">
                Peak vs Capacity
              </span>
            </div>
            <p className="text-xs text-slate-400 mt-0.5">
              Live power flow through 33kV &amp; 11kV distribution corridors.
            </p>
          </div>

          <div className="h-64 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={barData} margin={{ top: 10, right: 10, left: -20, bottom: 20 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                <XAxis dataKey="name" stroke="#64748b" tick={{ fontSize: 10 }} angle={-25} textAnchor="end" />
                <YAxis stroke="#64748b" tick={{ fontSize: 10 }} />
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.fullName}</div>
                          <div className="text-cyan-400 font-mono">Current Load: {data.load} MW</div>
                          <div className="text-slate-400 font-mono">Thermal Limit: {data.capacity} MW</div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
                <Bar dataKey="load" fill="#06b6d4" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="flex items-center justify-between text-[11px] text-slate-400 border-t border-slate-800/80 pt-3">
            <span>Thermal Line Capacity Margin:</span>
            <span className="font-mono font-bold text-emerald-400">42% Available Headroom</span>
          </div>
        </div>
      </div>
    </div>
  );
}

