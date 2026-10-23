/**
 * Observation 1: Complaints Lodged & Resolved, Customer Growth Query Engine.
 * Supports:
 * - Yearly: historical recorded years (2026 YTD, 2025, 2024, 2023)
 * - Quarterly: Q1 (Jan-Mar), Q2 (Apr-Jun), Q3 (Jul-Sep), Q4 (Oct-Dec)
 * - Monthly: Full 12-Month spectrum (January through December)
 * - Daily: Any given month and day query with auto-telemetry
 * - Displays: Total Lodged, Resolved, Unresolved, %, and New Customers
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import { AnalyticsPeriod, PeriodMetrics } from "./types";
import {
  YEARLY_METRICS_DATA,
  QUARTERLY_METRICS_DATA,
  MONTHLY_METRICS_DATA,
  getDailyMetricsForMonth,
  PERIOD_METRICS_MAP,
} from "./analyticsData";
import { Filter, CheckCircle2, AlertCircle, UserPlus, Clock, Calendar } from "lucide-react";

const MONTHS = [
  { key: "01", name: "January", short: "Jan" },
  { key: "02", name: "February", short: "Feb" },
  { key: "03", name: "March", short: "Mar" },
  { key: "04", name: "April", short: "Apr" },
  { key: "05", name: "May", short: "May" },
  { key: "06", name: "June", short: "Jun" },
  { key: "07", name: "July", short: "Jul" },
  { key: "08", name: "August", short: "Aug" },
  { key: "09", name: "September", short: "Sep" },
  { key: "10", name: "October", short: "Oct" },
  { key: "11", name: "November", short: "Nov" },
  { key: "12", name: "December", short: "Dec" },
];

const QUARTERS = [
  { id: "Q1", label: "Q1", range: "Jan, Feb, Mar" },
  { id: "Q2", label: "Q2", range: "Apr, May, Jun" },
  { id: "Q3", label: "Q3", range: "Jul, Aug, Sep" },
  { id: "Q4", label: "Q4", range: "Oct, Nov, Dec" },
];

const RECORDED_YEARS = [
  { year: "2026", badge: "Current YTD" },
  { year: "2025", badge: "Historical Audit" },
  { year: "2024", badge: "Historical Audit" },
  { year: "2023", badge: "Centenary Audit" },
];

export function AnalyticsPeriodQueryCard() {
  const [activePeriod, setActivePeriod] = useState<AnalyticsPeriod>("MONTHLY");
  const [selectedYear, setSelectedYear] = useState<string>("2026");
  const [selectedQuarter, setSelectedQuarter] = useState<string>("Q4");
  const [selectedMonth, setSelectedMonth] = useState<string>("10");
  const [dailyMonth, setDailyMonth] = useState<string>("10");
  const [dailyDay, setDailyDay] = useState<number>(23);

  // Compute active telemetry metrics based on active primary tab & sub-selections
  let data: PeriodMetrics;
  if (activePeriod === "YEARLY") {
    data = YEARLY_METRICS_DATA[selectedYear] || YEARLY_METRICS_DATA["2026"];
  } else if (activePeriod === "QUARTERLY") {
    data = QUARTERLY_METRICS_DATA[selectedQuarter] || QUARTERLY_METRICS_DATA.Q4;
  } else if (activePeriod === "MONTHLY") {
    data = MONTHLY_METRICS_DATA[selectedMonth] || MONTHLY_METRICS_DATA["10"];
  } else if (activePeriod === "DAILY") {
    data = getDailyMetricsForMonth(dailyMonth, dailyDay, "2026");
  } else {
    data = PERIOD_METRICS_MAP.WEEKLY;
  }

  const primaryTabs: { id: AnalyticsPeriod; label: string }[] = [
    { id: "DAILY", label: "Daily Query" },
    { id: "WEEKLY", label: "Weekly View" },
    { id: "MONTHLY", label: "Monthly (Jan-Dec)" },
    { id: "QUARTERLY", label: "Quarterly (Q1-Q4)" },
    { id: "YEARLY", label: "Yearly Records" },
  ];

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5 text-xs shadow-md">
      {/* Top Header & Primary Period Tabs */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-3 pb-3 border-b border-slate-800">
        <div className="flex items-center space-x-2">
          <Filter className="w-4 h-4 text-blue-400" />
          <div>
            <h3 className="font-bold text-slate-100 uppercase tracking-wide">
              Complaints Lodged &amp; Resolved Query Engine
            </h3>
            <p className="text-[11px] text-blue-400/90 font-mono flex items-center gap-1.5 mt-0.5">
              <Calendar className="w-3 h-3 text-slate-400" />
              <span>Active Scope: {data.label}</span>
            </p>
          </div>
        </div>

        {/* 5-Period Primary Tabs */}
        <div className="flex flex-wrap gap-1 bg-slate-950 p-1 rounded-lg border border-slate-800">
          {primaryTabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActivePeriod(tab.id)}
              className={`px-2.5 py-1 rounded text-[11px] font-semibold transition-all ${
                activePeriod === tab.id
                  ? "bg-blue-600 text-white shadow-sm"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Dynamic Sub-Filter Options Bar */}
      <div className="bg-slate-950/80 border border-slate-800/80 rounded-lg p-2.5 mb-4">
        {/* YEARLY SUB-OPTIONS */}
        {activePeriod === "YEARLY" && (
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider mr-1">
              Select Recorded Year:
            </span>
            {RECORDED_YEARS.map((y) => (
              <button
                key={y.year}
                onClick={() => setSelectedYear(y.year)}
                className={`px-3 py-1 rounded text-[11px] font-mono font-bold transition-all border ${
                  selectedYear === y.year
                    ? "bg-blue-600 text-white border-blue-500 shadow-sm"
                    : "bg-slate-900 text-slate-400 border-slate-800 hover:border-slate-700 hover:text-slate-200"
                }`}
              >
                {y.year} <span className="text-[9px] font-sans font-normal opacity-75">({y.badge})</span>
              </button>
            ))}
          </div>
        )}

        {/* QUARTERLY SUB-OPTIONS (Q1-Q4 with months) */}
        {activePeriod === "QUARTERLY" && (
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider mr-1">
              Select Quarter:
            </span>
            {QUARTERS.map((q) => (
              <button
                key={q.id}
                onClick={() => setSelectedQuarter(q.id)}
                className={`px-3 py-1 rounded text-[11px] font-semibold transition-all border ${
                  selectedQuarter === q.id
                    ? "bg-blue-600 text-white border-blue-500 shadow-sm"
                    : "bg-slate-900 text-slate-400 border-slate-800 hover:border-slate-700 hover:text-slate-200"
                }`}
              >
                <span className="font-bold">{q.label}:</span>{" "}
                <span className="text-[10px] text-slate-300 font-normal">({q.range})</span>
              </button>
            ))}
          </div>
        )}

        {/* MONTHLY SUB-OPTIONS (Jan - Dec) */}
        {activePeriod === "MONTHLY" && (
          <div className="flex flex-col gap-1.5">
            <div className="flex items-center justify-between text-[11px]">
              <span className="font-semibold text-slate-400 uppercase tracking-wider">
                Select Month of Record (2026):
              </span>
              <span className="text-blue-400 font-mono font-medium">
                {MONTHS.find((m) => m.key === selectedMonth)?.name} 2026
              </span>
            </div>
            <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-12 gap-1">
              {MONTHS.map((m) => (
                <button
                  key={m.key}
                  onClick={() => setSelectedMonth(m.key)}
                  className={`px-2 py-1 rounded text-center text-[11px] font-semibold transition-all border ${
                    selectedMonth === m.key
                      ? "bg-blue-600 text-white border-blue-500 shadow-sm"
                      : "bg-slate-900 text-slate-400 border-slate-800 hover:border-slate-700 hover:text-slate-200"
                  }`}
                >
                  {m.short}
                </button>
              ))}
            </div>
          </div>
        )}

        {/* DAILY SUB-OPTIONS (Month + Day selector) */}
        {activePeriod === "DAILY" && (
          <div className="flex flex-wrap items-center gap-3">
            <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
              Daily Query Filter:
            </span>

            {/* Month Select */}
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] text-slate-400">Month:</span>
              <select
                value={dailyMonth}
                onChange={(e) => setDailyMonth(e.target.value)}
                className="bg-slate-900 border border-slate-700 text-slate-200 text-[11px] rounded px-2.5 py-1 focus:outline-none focus:border-blue-500"
              >
                {MONTHS.map((m) => (
                  <option key={m.key} value={m.key}>
                    {m.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Day Select */}
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] text-slate-400">Day:</span>
              <select
                value={dailyDay}
                onChange={(e) => setDailyDay(parseInt(e.target.value, 10))}
                className="bg-slate-900 border border-slate-700 text-slate-200 text-[11px] font-mono rounded px-2.5 py-1 focus:outline-none focus:border-blue-500"
              >
                {Array.from({ length: 31 }, (_, i) => i + 1).map((d) => (
                  <option key={d} value={d}>
                    Day {d < 10 ? `0${d}` : d}
                  </option>
                ))}
              </select>
            </div>

            {/* Quick Reset to Today Button */}
            <button
              onClick={() => {
                setDailyMonth("10");
                setDailyDay(23);
              }}
              className="text-[10px] font-semibold px-2 py-1 rounded bg-blue-900/40 text-blue-300 border border-blue-700/50 hover:bg-blue-800/50 transition-all ml-auto"
            >
              Today (23 Oct 2026)
            </button>
          </div>
        )}

        {/* WEEKLY SUB-INFO */}
        {activePeriod === "WEEKLY" && (
          <div className="flex items-center justify-between text-[11px]">
            <span className="text-slate-400">
              Active Operational Week: <strong className="text-slate-200">Week 43, 2026 (Oct 19 - Oct 25)</strong>
            </span>
            <span className="text-emerald-400 font-mono text-[10px]">● Live Real-Time Aggregation</span>
          </div>
        )}
      </div>

      {/* 4 Core Telemetry Metric Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-4">
        {/* Total Complaints Lodged */}
        <div className="bg-slate-950/70 border border-slate-800/80 p-3 rounded-lg">
          <div className="flex items-center justify-between text-slate-400 text-[11px] mb-1">
            <span>Total Lodged</span>
            <AlertCircle className="w-3.5 h-3.5 text-blue-400" />
          </div>
          <div className="text-2xl font-bold text-slate-100 font-mono">{data.lodged_count.toLocaleString()}</div>
          <div className="text-[10px] text-slate-500 mt-0.5">Tickets Inflow</div>
        </div>

        {/* Complaints Resolved */}
        <div className="bg-slate-950/70 border border-slate-800/80 p-3 rounded-lg">
          <div className="flex items-center justify-between text-slate-400 text-[11px] mb-1">
            <span>Tickets Resolved</span>
            <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400" />
          </div>
          <div className="text-2xl font-bold text-emerald-400 font-mono">{data.resolved_count.toLocaleString()}</div>
          <div className="text-[10px] text-emerald-500/80 mt-0.5">
            Confirmed Closes ({data.resolved_pct}%)
          </div>
        </div>

        {/* Active Unresolved */}
        <div className="bg-slate-950/70 border border-slate-800/80 p-3 rounded-lg">
          <div className="flex items-center justify-between text-slate-400 text-[11px] mb-1">
            <span>Active Unresolved</span>
            <Clock className="w-3.5 h-3.5 text-amber-400" />
          </div>
          <div className="text-2xl font-bold text-amber-400 font-mono">{data.unresolved_count.toLocaleString()}</div>
          <div className="text-[10px] text-amber-500/80 mt-0.5">
            Pending / Backlog ({data.unresolved_pct}%)
          </div>
        </div>

        {/* New Customers Registered */}
        <div className="bg-slate-950/70 border border-slate-800/80 p-3 rounded-lg">
          <div className="flex items-center justify-between text-slate-400 text-[11px] mb-1">
            <span>New Customers</span>
            <UserPlus className="w-3.5 h-3.5 text-indigo-400" />
          </div>
          <div className="text-2xl font-bold text-indigo-400 font-mono">{data.new_customers_count.toLocaleString()}</div>
          <div className="text-[10px] text-slate-500 mt-0.5">NIN Verified Onboardings</div>
        </div>
      </div>

      {/* Resolution Velocity & Percentage Breakdown Bar */}
      <div className="bg-slate-950/80 border border-slate-800 p-3.5 rounded-lg space-y-2">
        <div className="flex justify-between items-center text-[11px]">
          <span className="font-semibold text-slate-200">
            Resolution Efficiency Rate:{" "}
            <span className="text-emerald-400 font-bold">{data.resolved_pct}% Resolved</span>
          </span>
          <span className="text-amber-400 font-bold font-mono">
            {data.unresolved_pct}% Unresolved ({data.unresolved_count} pending)
          </span>
        </div>

        <div className="w-full bg-slate-800 h-3 rounded-full overflow-hidden flex">
          <div
            className="bg-emerald-500 h-full transition-all duration-500 flex items-center justify-center text-[9px] text-slate-950 font-bold"
            style={{ width: `${data.resolved_pct}%` }}
          >
            {data.resolved_pct}%
          </div>
          <div
            className="bg-amber-500 h-full transition-all duration-500"
            style={{ width: `${data.unresolved_pct}%` }}
          />
        </div>

        <div className="flex justify-between text-[10px] text-slate-500 font-mono pt-1">
          <span>Target Standard: &ge;80.0% NERC Closure Rate</span>
          <span className={data.resolved_pct >= 80 ? "text-emerald-400 font-semibold" : "text-amber-400 font-semibold"}>
            {data.resolved_pct >= 80 ? "✓ REGULATORY COMPLIANT" : "⚠ ESCALATION WATCH"}
          </span>
        </div>
      </div>
    </div>
  );
}
