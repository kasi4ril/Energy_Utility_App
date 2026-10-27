/**
 * Metering Fleet Analytics Dashboard for Metering Services Unit.
 * Visualizes 14,850 meter population distribution, meter types, and tariff bands.
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

const METER_CATEGORIES = [
  { name: "Single-Phase STS Prepaid (Residential)", value: 10200, color: "#06b6d4" },
  { name: "Three-Phase STS Prepaid (Commercial)", value: 3850, color: "#3b82f6" },
  { name: "Maximum Demand (MD CT/PT Industrial)", value: 800, color: "#8b5cf6" },
];

const DISTRICT_METERS = [
  { district: "Anglo-Jos", count: 4200 },
  { district: "Bukuru", count: 3800 },
  { district: "Rayfield", count: 3400 },
  { district: "Vom", count: 1450 },
  { district: "Makeri", count: 1200 },
  { district: "Bokkos", count: 800 },
];

export function MeteringAnalyticsView() {
  const totalFleet = METER_CATEGORIES.reduce((acc, m) => acc + m.value, 0);

  return (
    <div className="space-y-6">
      {/* Top Fleet KPI Ribbon */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Total Metered Population</div>
          <div className="text-2xl font-black text-cyan-400 mt-1">
            {totalFleet.toLocaleString()} Meters
          </div>
          <div className="text-[10px] text-slate-500 mt-0.5">100% STS IEC 62055-41 Standard</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">TID Rollover Progress</div>
          <div className="text-2xl font-black text-emerald-400 mt-1">94.8%</div>
          <div className="text-[10px] text-emerald-400 mt-0.5">14,078 Meters Successfully Migrated</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Class 0.5S MD Fleet</div>
          <div className="text-2xl font-black text-purple-400 mt-1">800 Units</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Precision Industrial Tariff Billing</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Bench Recalibration Rate</div>
          <div className="text-2xl font-black text-white mt-1">99.4%</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Within &plusmn;0.5% NEMSA Error Band</div>
        </div>
      </div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Meter Type Donut Pie Chart */}
        <div className="lg:col-span-5 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-white">Meter Fleet by Operational Class</h3>
            <p className="text-xs text-slate-400 mt-0.5">
              Single-phase residential vs three-phase commercial vs industrial MD.
            </p>
          </div>

          <div className="h-60 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={METER_CATEGORIES}
                  cx="50%"
                  cy="50%"
                  innerRadius={55}
                  outerRadius={85}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {METER_CATEGORIES.map((entry, idx) => (
                    <Cell key={`cell-${idx}`} fill={entry.color} />
                  ))}
                </Pie>
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      const pct = ((data.value / totalFleet) * 100).toFixed(1);
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.name}</div>
                          <div className="text-cyan-400 font-mono">
                            {data.value.toLocaleString()} Meters ({pct}%)
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
            {METER_CATEGORIES.map((m) => (
              <div key={m.name} className="flex items-center justify-between">
                <span className="flex items-center space-x-2 truncate">
                  <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ backgroundColor: m.color }} />
                  <span className="text-slate-300 truncate">{m.name}</span>
                </span>
                <span className="font-mono font-bold text-white ml-2">
                  {m.value.toLocaleString()}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Meters by District Bar Chart */}
        <div className="lg:col-span-7 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-white">Meter Population by Operational District</h3>
            <p className="text-xs text-slate-400 mt-0.5">Active metered connection points across service areas.</p>
          </div>

          <div className="h-60 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={DISTRICT_METERS} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                <XAxis dataKey="district" stroke="#64748b" tick={{ fontSize: 11 }} />
                <YAxis stroke="#64748b" tick={{ fontSize: 11 }} />
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.district} District</div>
                          <div className="text-cyan-400 font-mono">{data.count.toLocaleString()} Meters</div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
                <Bar dataKey="count" fill="#06b6d4" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="border-t border-slate-800/80 pt-3 flex items-center justify-between text-xs text-slate-400">
            <span>Universal Prepayment Meter Penetration:</span>
            <span className="text-emerald-400 font-semibold font-mono">100% Zero-Estimation Achieved</span>
          </div>
        </div>
      </div>
    </div>
  );
}

