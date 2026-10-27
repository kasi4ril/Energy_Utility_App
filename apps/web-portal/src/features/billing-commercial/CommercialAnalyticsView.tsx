/**
 * Revenue & Tariff Analytics Dashboard for Commercial Operations.
 * Visualizes collections by tariff bands, customer categories, and payment channels.
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

const REVENUE_SPLIT = [
  { name: "Prepaid STS (Residential & Commercial)", value: 3450000, color: "#14b8a6" },
  { name: "Maximum Demand (Industrial MD)", value: 1150000, color: "#3b82f6" },
  { name: "Postpaid Government & Institutions", value: 250000, color: "#8b5cf6" },
];

const WEEKLY_TREND = [
  { day: "Mon", collections: 4.1 },
  { day: "Tue", collections: 4.3 },
  { day: "Wed", collections: 4.85 },
  { day: "Thu", collections: 4.6 },
  { day: "Fri", collections: 5.2 },
  { day: "Sat", collections: 3.8 },
  { day: "Sun", collections: 2.9 },
];

const TARIFF_BANDS = [
  { band: "Band A", hours: "20 - 24 hrs", rate: "₦225.00 / kWh", share: "48%" },
  { band: "Band B", hours: "16 - 20 hrs", rate: "₦180.00 / kWh", share: "28%" },
  { band: "Band C", hours: "12 - 16 hrs", rate: "₦120.00 / kWh", share: "15%" },
  { band: "Band D", hours: "8 - 12 hrs", rate: "₦85.00 / kWh", share: "7%" },
  { band: "Band E", hours: "4 - 8 hrs", rate: "₦65.00 / kWh", share: "2%" },
];

export function CommercialAnalyticsView() {
  const totalCollections = REVENUE_SPLIT.reduce((acc, s) => acc + s.value, 0);

  return (
    <div className="space-y-6">
      {/* Revenue KPI Ribbon */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Total Shift Collections</div>
          <div className="text-2xl font-black text-teal-400 mt-1">
            &#8358;{(totalCollections / 1000000).toFixed(2)}M
          </div>
          <div className="text-[10px] text-emerald-400 mt-0.5">+4.8% vs Yesterday</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Prepaid STS Share</div>
          <div className="text-2xl font-black text-white mt-1">91.2%</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Zero-Estimation Compliant</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Average Tariff Yield</div>
          <div className="text-2xl font-black text-cyan-400 mt-1">&#8358;184.50</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Per kWh Delivered</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Theft Loss Recovery</div>
          <div className="text-2xl font-black text-emerald-400 mt-1">&#8358;840K</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Recovered via AI Audits</div>
        </div>
      </div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Revenue Pie Chart */}
        <div className="lg:col-span-5 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-white">Revenue by Customer Class</h3>
            <p className="text-xs text-slate-400 mt-0.5">
              STS Prepayment tokens vs Industrial Maximum Demand.
            </p>
          </div>

          <div className="h-60 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={REVENUE_SPLIT}
                  cx="50%"
                  cy="50%"
                  innerRadius={55}
                  outerRadius={85}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {REVENUE_SPLIT.map((entry, idx) => (
                    <Cell key={`cell-${idx}`} fill={entry.color} />
                  ))}
                </Pie>
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      const pct = ((data.value / totalCollections) * 100).toFixed(1);
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.name}</div>
                          <div className="text-teal-400 font-mono">
                            &#8358;{(data.value / 1000000).toFixed(2)}M ({pct}%)
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
            {REVENUE_SPLIT.map((s) => (
              <div key={s.name} className="flex items-center justify-between">
                <span className="flex items-center space-x-2 truncate">
                  <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ backgroundColor: s.color }} />
                  <span className="text-slate-300 truncate">{s.name}</span>
                </span>
                <span className="font-mono font-bold text-white ml-2">
                  &#8358;{(s.value / 1000000).toFixed(2)}M
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* 7-Day Collections Bar Chart */}
        <div className="lg:col-span-7 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-white">Daily Collections Trend (&#8358; Millions)</h3>
            <p className="text-xs text-slate-400 mt-0.5">Rolling 7-day cash flow performance.</p>
          </div>

          <div className="h-60 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={WEEKLY_TREND} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                <XAxis dataKey="day" stroke="#64748b" tick={{ fontSize: 11 }} />
                <YAxis stroke="#64748b" tick={{ fontSize: 11 }} />
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.day} Collections</div>
                          <div className="text-teal-400 font-mono">&#8358;{data.collections}M</div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
                <Bar dataKey="collections" fill="#14b8a6" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="border-t border-slate-800/80 pt-3 flex items-center justify-between text-xs text-slate-400">
            <span>Weekly Total: <strong className="text-white font-mono">&#8358;30.65M</strong></span>
            <span className="text-emerald-400 font-semibold">&uarr; 3.2% vs Previous Week</span>
          </div>
        </div>
      </div>

      {/* NERC Multi-Year Tariff Matrix */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
          <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
            Active NERC Tariff Order Rates (MYTO)
          </span>
          <span className="text-xs text-slate-400 font-mono">Approved NERC Matrix</span>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-2.5 px-4">Tariff Band</th>
                <th className="py-2.5 px-3">Mandatory Supply Hours</th>
                <th className="py-2.5 px-3">Approved Rate per kWh</th>
                <th className="py-2.5 px-4 text-right">Revenue Share</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {TARIFF_BANDS.map((t) => (
                <tr key={t.band} className="hover:bg-slate-800/40 transition">
                  <td className="py-2.5 px-4 font-bold text-teal-400">{t.band}</td>
                  <td className="py-2.5 px-3 font-mono text-slate-300">{t.hours}</td>
                  <td className="py-2.5 px-3 font-mono font-bold text-white">{t.rate}</td>
                  <td className="py-2.5 px-4 text-right font-mono text-slate-400">{t.share}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

