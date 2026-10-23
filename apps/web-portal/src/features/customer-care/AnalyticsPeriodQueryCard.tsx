/**
 * Observation 1: Complaints Lodged & Resolved, Customer Growth Query Engine (Option A).
 * Visual Dashboard with:
 * - Hierarchical Query Dropdowns: Year (2023-2026), Month (All/01-12), Day (All/01-31)
 * - Donut Chart Gauge: Resolved vs Active Unresolved (%) with Center NERC Indicator
 * - Dual Grouped Bar Histogram: Inflow vs Resolution Pace across queried intervals
 * - 4 Core KPI Summary Cards
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import { PeriodMetrics } from "./types";
import {
  YEARLY_METRICS_DATA,
  getMonthMetricsForYear,
  getDailyMetricsForMonth,
  getTicketHistogramData,
} from "./analyticsData";
import { Filter, CheckCircle2, AlertCircle, UserPlus, Clock, Calendar, ChevronDown } from "lucide-react";
import {
  PieChart,
  Pie,
  Cell,
  Tooltip as RechartsTooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Legend,
} from "recharts";

const MONTHS = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December"
].map((name, i) => ({
  key: (i + 1).toString().padStart(2, "0"),
  name: `${(i + 1).toString().padStart(2, "0")} - ${name}`,
}));

const RECORDED_YEARS = [
  { year: "2026", label: "2026 (Current Operating Year)" },
  { year: "2025", label: "2025 (Annual Historical Audit)" },
  { year: "2024", label: "2024 (Annual Historical Audit)" },
  { year: "2023", label: "2023 (Centenary Audit)" },
];

export function AnalyticsPeriodQueryCard() {
  const [mounted, setMounted] = useState(false);
  const [selectedYear, setSelectedYear] = useState<string>("2026");
  const [selectedMonth, setSelectedMonth] = useState<string>("10");
  const [selectedDay, setSelectedDay] = useState<string>("23");

  useEffect(() => {
    setMounted(true);
  }, []);

  // Compute active telemetry metrics based on hierarchical dropdown selections:
  let data: PeriodMetrics;
  if (selectedMonth === "ALL") {
    data = YEARLY_METRICS_DATA[selectedYear] || YEARLY_METRICS_DATA["2026"];
  } else if (selectedDay === "ALL") {
    data = getMonthMetricsForYear(selectedMonth, selectedYear);
  } else {
    data = getDailyMetricsForMonth(selectedMonth, parseInt(selectedDay, 10), selectedYear);
  }

  const donutData = [
    { name: "Tickets Resolved", value: data.resolved_count, color: "#10b981" },
    { name: "Active Unresolved", value: data.unresolved_count, color: "#f59e0b" },
  ];

  const histogramData = getTicketHistogramData(selectedYear, selectedMonth, selectedDay);
  const isCompliant = data.resolved_pct >= 80.0;

  const kpis = [
    { label: "Total Lodged", value: data.lodged_count.toLocaleString(), sub: "Tickets Inflow", icon: AlertCircle, color: "text-blue-400", valColor: "text-slate-100" },
    { label: "Tickets Resolved", value: data.resolved_count.toLocaleString(), sub: `Confirmed Closes (${data.resolved_pct}%)`, icon: CheckCircle2, color: "text-emerald-400", valColor: "text-emerald-400" },
    { label: "Active Unresolved", value: data.unresolved_count.toLocaleString(), sub: `Pending Backlog (${data.unresolved_pct}%)`, icon: Clock, color: "text-amber-400", valColor: "text-amber-400" },
    { label: "New Customers", value: data.new_customers_count.toLocaleString(), sub: "NIN Verified Onboardings", icon: UserPlus, color: "text-indigo-400", valColor: "text-indigo-400" },
  ];

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5 text-xs shadow-md">
      {/* Top Header & Active Scope Indicator */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-3 pb-3 border-b border-slate-800">
        <div className="flex items-center space-x-2">
          <Filter className="w-4 h-4 text-blue-400" />
          <div>
            <h3 className="font-bold text-slate-100 uppercase tracking-wide">
              Complaints &amp; SLA Resolution Intelligence
            </h3>
            <p className="text-[11px] text-blue-400/90 font-mono flex items-center gap-1.5 mt-0.5">
              <Calendar className="w-3 h-3 text-slate-400" />
              <span>Active Scope: {data.label}</span>
            </p>
          </div>
        </div>

        {/* Quick Shortcut Presets */}
        <div className="flex flex-wrap items-center gap-1.5">
          <button
            onClick={() => { setSelectedYear("2026"); setSelectedMonth("10"); setSelectedDay("23"); }}
            className="px-2.5 py-1 rounded text-[10px] font-semibold bg-emerald-950/70 text-emerald-300 border border-emerald-600/40 hover:bg-emerald-900/80 transition-all"
            title="Jump to Today's Live Operating Grid Cycle"
          >
            ⚡ Today (23 Oct 2026)
          </button>
          <button
            onClick={() => { setSelectedYear("2026"); setSelectedMonth("10"); setSelectedDay("ALL"); }}
            className="px-2.5 py-1 rounded text-[10px] font-semibold bg-blue-950/70 text-blue-300 border border-blue-600/40 hover:bg-blue-900/80 transition-all"
          >
            📅 Oct 2026 Month
          </button>
          <button
            onClick={() => { setSelectedYear("2026"); setSelectedMonth("ALL"); setSelectedDay("ALL"); }}
            className="px-2.5 py-1 rounded text-[10px] font-semibold bg-slate-800 text-slate-300 border border-slate-700 hover:bg-slate-700 transition-all"
          >
            📊 2026 Full Year
          </button>
        </div>
      </div>

      {/* Linked Hierarchical Dropdowns Bar: Year -> Month -> Day */}
      <div className="bg-slate-950/80 border border-slate-800/90 rounded-lg p-3 mb-4">
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {/* 1. Year Dropdown (2023 - 2026) */}
          <div className="space-y-1">
            <label className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Select Year:</label>
            <div className="relative">
              <select
                value={selectedYear}
                onChange={(e) => setSelectedYear(e.target.value)}
                className="w-full appearance-none bg-slate-900 border border-slate-700 hover:border-slate-600 text-slate-100 text-xs rounded-lg px-3 py-2 pr-8 focus:outline-none focus:border-blue-500 font-mono font-semibold cursor-pointer shadow-sm"
              >
                {RECORDED_YEARS.map((y) => (
                  <option key={y.year} value={y.year}>{y.label}</option>
                ))}
              </select>
              <ChevronDown className="w-3.5 h-3.5 text-slate-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>
          </div>

          {/* 2. Month Dropdown */}
          <div className="space-y-1">
            <label className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Select Month:</label>
            <div className="relative">
              <select
                value={selectedMonth}
                onChange={(e) => {
                  const m = e.target.value;
                  setSelectedMonth(m);
                  if (m === "ALL") setSelectedDay("ALL");
                }}
                className="w-full appearance-none bg-slate-900 border border-slate-700 hover:border-slate-600 text-slate-100 text-xs rounded-lg px-3 py-2 pr-8 focus:outline-none focus:border-blue-500 font-medium cursor-pointer shadow-sm"
              >
                <option value="ALL">Full Year Summary (All 12 Months)</option>
                {MONTHS.map((m) => (
                  <option key={m.key} value={m.key}>{m.name}</option>
                ))}
              </select>
              <ChevronDown className="w-3.5 h-3.5 text-slate-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>
          </div>

          {/* 3. Day Dropdown */}
          <div className="space-y-1">
            <label className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Select Day Report:</label>
            <div className="relative">
              <select
                value={selectedDay}
                disabled={selectedMonth === "ALL"}
                onChange={(e) => setSelectedDay(e.target.value)}
                className={`w-full appearance-none bg-slate-900 border border-slate-700 text-slate-100 text-xs rounded-lg px-3 py-2 pr-8 focus:outline-none focus:border-blue-500 font-mono font-medium cursor-pointer shadow-sm ${
                  selectedMonth === "ALL" ? "opacity-40 cursor-not-allowed" : "hover:border-slate-600"
                }`}
              >
                <option value="ALL">Entire Month Summary (All Days)</option>
                {Array.from({ length: 31 }, (_, i) => i + 1).map((d) => (
                  <option key={d} value={d}>Day {d < 10 ? `0${d}` : d}</option>
                ))}
              </select>
              <ChevronDown className="w-3.5 h-3.5 text-slate-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>
          </div>
        </div>
      </div>

      {/* 4 Core Telemetry Metric Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-4">
        {kpis.map((k, i) => {
          const Icon = k.icon;
          return (
            <div key={i} className="bg-slate-950/70 border border-slate-800/80 p-3 rounded-lg">
              <div className="flex items-center justify-between text-slate-400 text-[11px] mb-1">
                <span>{k.label}</span>
                <Icon className={`w-3.5 h-3.5 ${k.color}`} />
              </div>
              <div className={`text-2xl font-bold font-mono ${k.valColor}`}>{k.value}</div>
              <div className="text-[10px] text-slate-500 mt-0.5">{k.sub}</div>
            </div>
          );
        })}
      </div>

      {/* Option A Visual Dashboard: Donut Chart Gauge & Dual Grouped Bar Histogram */}
      {mounted && (
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-3.5">
          {/* 1. Left: Donut Chart Gauge (Resolved vs Active Unresolved) */}
          <div className="lg:col-span-4 bg-slate-950/80 border border-slate-800 rounded-lg p-3.5 flex flex-col justify-between">
            <div className="flex items-center justify-between border-b border-slate-800/70 pb-2 mb-1">
              <span className="font-semibold text-slate-200 uppercase tracking-wide text-[11px]">
                Resolution Efficiency
              </span>
              <span
                className={`text-[10px] px-2 py-0.5 rounded font-mono font-bold border ${
                  isCompliant
                    ? "bg-emerald-950/70 text-emerald-300 border-emerald-500/40"
                    : "bg-amber-950/70 text-amber-300 border-amber-500/40"
                }`}
              >
                {isCompliant ? "✓ NERC COMPLIANT" : "⚠ ESCALATION"}
              </span>
            </div>

            <div className="relative h-44 flex items-center justify-center">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={donutData}
                    cx="50%"
                    cy="50%"
                    innerRadius={50}
                    outerRadius={68}
                    startAngle={90}
                    endAngle={-270}
                    dataKey="value"
                    strokeWidth={2}
                    stroke="#090d16"
                  >
                    {donutData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                  <RechartsTooltip
                    contentStyle={{
                      backgroundColor: "#090d16",
                      borderColor: "#1e293b",
                      borderRadius: "8px",
                      fontSize: "11px",
                      color: "#f8fafc",
                    }}
                  />
                </PieChart>
              </ResponsiveContainer>

              {/* Center Donut Label */}
              <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none text-center">
                <span className="text-xl font-bold font-mono text-emerald-400">{data.resolved_pct}%</span>
                <span className="text-[9px] uppercase tracking-wider text-slate-400 font-semibold">Resolved</span>
              </div>
            </div>

            <div className="flex items-center justify-around pt-2 border-t border-slate-800/70 text-[10px] font-mono">
              <div className="flex items-center space-x-1.5">
                <span className="w-2.5 h-2.5 rounded-full bg-emerald-500 inline-block" />
                <span className="text-slate-300">Resolved ({data.resolved_count})</span>
              </div>
              <div className="flex items-center space-x-1.5">
                <span className="w-2.5 h-2.5 rounded-full bg-amber-500 inline-block" />
                <span className="text-slate-300">Backlog ({data.unresolved_count})</span>
              </div>
            </div>
          </div>

          {/* 2. Right: Dual Grouped Bar Histogram (Lodged vs Resolved Pace) */}
          <div className="lg:col-span-8 bg-slate-950/80 border border-slate-800 rounded-lg p-3.5 flex flex-col justify-between">
            <div className="flex items-center justify-between border-b border-slate-800/70 pb-2 mb-1">
              <div>
                <span className="font-semibold text-slate-200 uppercase tracking-wide text-[11px] block">
                  Tickets Inflow vs. Resolution Pace Histogram
                </span>
                <span className="text-[10px] text-slate-400 font-mono">
                  {selectedMonth === "ALL"
                    ? "Monthly Velocity (Jan - Dec)"
                    : selectedDay === "ALL"
                    ? "Weekly Operational Cadence (W1 - W4)"
                    : "Hourly Shift Dispatch Progression"}
                </span>
              </div>
              <span className="text-[10px] font-mono text-blue-400 bg-blue-950/60 px-2 py-0.5 rounded border border-blue-600/30">
                Omnichannel Intake
              </span>
            </div>

            <div className="h-44 w-full">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={histogramData} margin={{ top: 8, right: 10, left: -20, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" opacity={0.6} />
                  <XAxis dataKey="label" stroke="#64748b" fontSize={10} tickLine={false} />
                  <YAxis stroke="#64748b" fontSize={10} tickLine={false} />
                  <RechartsTooltip
                    contentStyle={{
                      backgroundColor: "#090d16",
                      borderColor: "#1e293b",
                      borderRadius: "8px",
                      fontSize: "11px",
                      color: "#f8fafc",
                    }}
                  />
                  <Legend wrapperStyle={{ fontSize: "10px", paddingTop: "4px" }} />
                  <Bar dataKey="lodged" name="Lodged" fill="#3b82f6" radius={[3, 3, 0, 0]} barSize={12} />
                  <Bar dataKey="resolved" name="Resolved" fill="#10b981" radius={[3, 3, 0, 0]} barSize={12} />
                </BarChart>
              </ResponsiveContainer>
            </div>

            <div className="flex justify-between items-center text-[10px] text-slate-500 font-mono pt-2 border-t border-slate-800/70">
              <span>Target Benchmark: &ge;80.0% Closure Within Regulatory SLA Window</span>
              <span className="text-slate-400">NERC Customer Protection Order &bull; v2026.4</span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
