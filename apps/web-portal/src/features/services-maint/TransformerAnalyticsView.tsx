/**
 * Transformer Fleet Health & Loading Analytics Dashboard.
 * Visualizes 33kV/11kV and 11kV/0.415kV DT fleet health, oil condition, and loading.
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

const DT_HEALTH_STATES = [
  { name: "Optimal Health (< 75% Load, Good Oil)", value: 118, color: "#10b981" },
  { name: "Warning Heavy Load (75 - 85% Load)", value: 18, color: "#06b6d4" },
  { name: "Critical Overloaded (> 85% Load)", value: 3, color: "#f59e0b" },
  { name: "Low Dielectric BDV (< 30kV)", value: 2, color: "#ef4444" },
  { name: "Workshop Major Overhaul", value: 1, color: "#8b5cf6" },
];

const DISTRICT_LOADING = [
  { district: "Anglo-Jos Ind.", loadingPct: 82, count: 28 },
  { district: "Bukuru", loadingPct: 74, count: 34 },
  { district: "Rayfield", loadingPct: 68, count: 30 },
  { district: "Vom Research", loadingPct: 58, count: 18 },
  { district: "Makeri", loadingPct: 71, count: 20 },
  { district: "Bokkos Rural", loadingPct: 48, count: 12 },
];

export function TransformerAnalyticsView() {
  const totalDts = DT_HEALTH_STATES.reduce((acc, s) => acc + s.value, 0);

  return (
    <div className="space-y-6">
      {/* Top Fleet KPI Ribbon */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Total Installed DT Fleet</div>
          <div className="text-2xl font-black text-amber-400 mt-1">{totalDts} Units</div>
          <div className="text-[10px] text-slate-500 mt-0.5">82.5 MVA Total Nameplate Capacity</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Average Fleet Loading</div>
          <div className="text-2xl font-black text-emerald-400 mt-1">64.5%</div>
          <div className="text-[10px] text-emerald-400 mt-0.5">35.5% Dynamic Load Headroom</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Oil Dielectric Pass Rate</div>
          <div className="text-2xl font-black text-cyan-400 mt-1">97.2%</div>
          <div className="text-[10px] text-slate-400 mt-0.5">BDV &gt; 40 kV (IEEE C57 Standard)</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Overhaul MTTR</div>
          <div className="text-2xl font-black text-white mt-1">4.2 Days</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Turnaround Time in Workshop</div>
        </div>
      </div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Health States Pie Chart */}
        <div className="lg:col-span-5 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-white">Transformer Fleet Health Distribution</h3>
            <p className="text-xs text-slate-400 mt-0.5">
              Operating state and dielectric oil risk classification.
            </p>
          </div>

          <div className="h-60 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={DT_HEALTH_STATES}
                  cx="50%"
                  cy="50%"
                  innerRadius={55}
                  outerRadius={85}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {DT_HEALTH_STATES.map((entry, idx) => (
                    <Cell key={`cell-${idx}`} fill={entry.color} />
                  ))}
                </Pie>
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      const pct = ((data.value / totalDts) * 100).toFixed(1);
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.name}</div>
                          <div className="text-amber-400 font-mono">
                            {data.value} DTs ({pct}%)
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

          <div className="space-y-1.5 text-[11px] border-t border-slate-800/80 pt-3">
            {DT_HEALTH_STATES.map((s) => (
              <div key={s.name} className="flex items-center justify-between">
                <span className="flex items-center space-x-2 truncate">
                  <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ backgroundColor: s.color }} />
                  <span className="text-slate-300 truncate">{s.name}</span>
                </span>
                <span className="font-mono font-bold text-white ml-2">{s.value} DTs</span>
              </div>
            ))}
          </div>
        </div>

        {/* District Loading Bar Chart */}
        <div className="lg:col-span-7 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-white">District Average DT Loading (% Capacity)</h3>
            <p className="text-xs text-slate-400 mt-0.5">Transformer loading across operational corridors.</p>
          </div>

          <div className="h-60 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={DISTRICT_LOADING} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                <XAxis dataKey="district" stroke="#64748b" tick={{ fontSize: 10 }} />
                <YAxis stroke="#64748b" tick={{ fontSize: 11 }} />
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.district}</div>
                          <div className="text-amber-400 font-mono">Average Loading: {data.loadingPct}%</div>
                          <div className="text-slate-400 font-mono">Installed Substation DTs: {data.count}</div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
                <Bar dataKey="loadingPct" fill="#f59e0b" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="border-t border-slate-800/80 pt-3 flex items-center justify-between text-xs text-slate-400">
            <span>Critical Load Redistribution Target:</span>
            <span className="text-amber-400 font-semibold font-mono">Anglo-Jos Industrial Substation #4</span>
          </div>
        </div>
      </div>
    </div>
  );
}

