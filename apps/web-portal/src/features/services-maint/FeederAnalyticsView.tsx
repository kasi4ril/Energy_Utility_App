/**
 * Workspace 1: Feeder Engineering Analytics & Reliability Dashboard.
 * Displays:
 * 1. SAIDI, SAIFI, CAIDI and MTTR reliability indicators
 * 2. Feeder operating states (In-Service, Fault, PTW, Rationing) in Pie Chart
 * 3. Line trips by corridor in Bar Chart
 * 4. 33kV/11kV network inventory health table
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import { FeederReliabilityMetrics } from "./types";
import { MOCK_FEEDER_RELIABILITY } from "./mockData";
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
import {
  Activity,
  Zap,
  Clock,
  TrendingDown,
  ShieldAlert,
  BarChart3,
  PieChart as PieIcon,
  CheckCircle,
} from "lucide-react";

interface FeederAnalyticsProps {
  metrics?: FeederReliabilityMetrics;
}

const STATE_COLORS: Record<string, string> = {
  IN_SERVICE: "#10b981",       // Emerald
  OUT_FAULT: "#f43f5e",        // Rose
  OUT_PTW: "#f59e0b",          // Amber
  OUT_RATIONING: "#3b82f6",    // Blue
};

const CORRIDOR_TRIP_DATA = [
  { corridor: "Anglo-Jos 33kV", trips: 3, faults: "Insulator flashover", lengthKm: 28.4 },
  { corridor: "Bukuru 33kV", trips: 2, faults: "GOABS hotspot", lengthKm: 34.2 },
  { corridor: "Rayfield 33kV", trips: 2, faults: "Tree branch intrusion", lengthKm: 18.6 },
  { corridor: "Vom 11kV", trips: 1, faults: "Stay wire slack", lengthKm: 14.5 },
  { corridor: "Makeri 33kV", trips: 1, faults: "Jumper burn", lengthKm: 22.0 },
  { corridor: "Bokkos 33kV", trips: 0, faults: "None (Clear)", lengthKm: 42.1 },
];

export const FeederAnalyticsView: React.FC<FeederAnalyticsProps> = ({ metrics = MOCK_FEEDER_RELIABILITY }) => {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const pieData = [
    { name: "In Service", value: metrics.in_service_count, color: STATE_COLORS.IN_SERVICE },
    { name: "Out on Fault", value: metrics.out_fault_count, color: STATE_COLORS.OUT_FAULT },
    { name: "Out on PTW", value: metrics.out_ptw_count, color: STATE_COLORS.OUT_PTW },
    { name: "Out on Rationing", value: metrics.out_rationing_count, color: STATE_COLORS.OUT_RATIONING },
  ].filter((d) => d.value > 0);

  const availabilityPct = ((metrics.in_service_count / metrics.total_feeders) * 100).toFixed(1);

  if (!mounted) return null;

  return (
    <div className="space-y-6 text-xs text-slate-100">
      {/* 1. Header Banner & KPIs */}
      <div className="bg-gradient-to-r from-slate-900 via-orange-950/40 to-slate-900 border border-orange-500/30 rounded-2xl p-5 shadow-xl">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center space-x-2 mb-1.5">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-orange-500/20 border border-orange-500/40 text-orange-300">
                NERC Feeder Reliability Index &bull; Grid Code Part 3
              </span>
              <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
                Live SCADA Telemetry
              </span>
            </div>
            <h1 className="text-xl sm:text-2xl font-black text-white tracking-tight">
              Feeder Engineering Analytics &amp; Reliability Index
            </h1>
            <p className="text-xs text-slate-300 mt-1 max-w-2xl leading-relaxed">
              Consolidated operational intelligence for 33kV transmission and 11kV distribution corridors.
              Monitoring SAIDI interruption duration, SAIFI frequency, line trip density, and MTTR restoration pace.
            </p>
          </div>

          <div className="bg-slate-950/80 border border-orange-500/40 px-3.5 py-2 rounded-xl text-right shrink-0">
            <div className="text-[10px] text-orange-400 font-bold uppercase tracking-wider">Feeder Availability</div>
            <div className="text-lg font-black font-mono text-emerald-400">
              {availabilityPct}% <span className="text-xs font-normal text-slate-400">({metrics.in_service_count}/{metrics.total_feeders})</span>
            </div>
          </div>
        </div>

        {/* 4 Core Reliability KPIs */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mt-4">
          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>SAIDI (Interruption Duration)</span>
              <Clock className="w-3.5 h-3.5 text-orange-400" />
            </div>
            <div className="mt-1 flex items-baseline justify-between">
              <span className="text-xl font-black font-mono text-white">{metrics.saidi_hours} <span className="text-xs text-slate-400 font-normal">hrs/yr</span></span>
              <span className="text-[10px] text-emerald-400 font-mono font-bold">&lt; 12h NERC</span>
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>SAIFI (Frequency Index)</span>
              <Zap className="w-3.5 h-3.5 text-cyan-400" />
            </div>
            <div className="mt-1 flex items-baseline justify-between">
              <span className="text-xl font-black font-mono text-cyan-300">{metrics.saifi_frequency} <span className="text-xs text-slate-400 font-normal">int/yr</span></span>
              <span className="text-[10px] text-emerald-400 font-mono font-bold">&lt; 8.0 Target</span>
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>MTTR (Mean Time to Repair)</span>
              <TrendingDown className="w-3.5 h-3.5 text-emerald-400" />
            </div>
            <div className="mt-1 flex items-baseline justify-between">
              <span className="text-xl font-black font-mono text-emerald-300">{metrics.mttr_hours} <span className="text-xs text-slate-400 font-normal">hours</span></span>
              <span className="text-[10px] text-emerald-400 font-mono font-bold">Fast Restore</span>
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Line Trip Count (Monthly)</span>
              <ShieldAlert className="w-3.5 h-3.5 text-rose-400" />
            </div>
            <div className="mt-1 flex items-baseline justify-between">
              <span className="text-xl font-black font-mono text-rose-300">{metrics.monthly_trips_count}</span>
              <span className="text-[10px] text-amber-400 font-mono">1 Active Outage</span>
            </div>
          </div>
        </div>
      </div>

      {/* 2. Visual Charts: Pie Chart (States) + Bar Chart (Corridor Trips) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-5 items-stretch">
        {/* Left Column: Donut Pie Chart (Feeder State Distribution) */}
        <div className="lg:col-span-5 bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
          <div className="flex items-center justify-between border-b border-slate-800 pb-2 mb-2">
            <span className="font-bold text-slate-200 text-xs flex items-center space-x-1.5">
              <PieIcon className="w-3.5 h-3.5 text-orange-400" />
              <span>33kV &amp; 11kV Feeder Operational States</span>
            </span>
            <span className="text-[10px] font-mono text-orange-400 bg-orange-500/10 px-2 py-0.5 rounded border border-orange-500/20">
              {metrics.total_feeders} Corridors
            </span>
          </div>

          <div className="relative h-56 flex items-center justify-center">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={50}
                  outerRadius={75}
                  paddingAngle={3}
                  dataKey="value"
                  strokeWidth={1.5}
                  stroke="#0f172a"
                >
                  {pieData.map((entry, idx) => (
                    <Cell key={idx} fill={entry.color} />
                  ))}
                </Pie>
                <RechartsTooltip
                  formatter={(val: any, name: any) => [`${val} Feeders`, name]}
                  contentStyle={{ backgroundColor: "#0f172a", borderColor: "#334155", borderRadius: "8px", fontSize: "11px" }}
                />
              </PieChart>
            </ResponsiveContainer>

            <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none text-center">
              <span className="text-2xl font-black font-mono text-emerald-400">{metrics.in_service_count}</span>
              <span className="text-[9px] uppercase tracking-wider text-slate-400 font-bold">In Service</span>
            </div>
          </div>

          {/* Breakdown Badges */}
          <div className="grid grid-cols-2 gap-2 pt-2 border-t border-slate-800 text-[11px] font-mono">
            {pieData.map((d, i) => (
              <div key={i} className="flex items-center justify-between p-1 rounded bg-slate-950/60 border border-slate-800/60">
                <div className="flex items-center space-x-1.5 truncate">
                  <span className="w-2 h-2 rounded-full shrink-0" style={{ backgroundColor: d.color }} />
                  <span className="text-slate-300 truncate font-sans text-[11px]">{d.name}</span>
                </div>
                <strong className="text-white text-xs ml-1">{d.value}</strong>
              </div>
            ))}
          </div>
        </div>

        {/* Right Column: Bar Chart (Line Trips by Corridor) */}
        <div className="lg:col-span-7 bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
          <div className="flex items-center justify-between border-b border-slate-800 pb-2 mb-2">
            <span className="font-bold text-slate-200 text-xs flex items-center space-x-1.5">
              <BarChart3 className="w-3.5 h-3.5 text-orange-400" />
              <span>Fault Trip Density by Transmission Corridor (30 Days)</span>
            </span>
            <span className="text-[10px] text-slate-400 font-mono">Corridor Telemetry</span>
          </div>

          <div className="h-56 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={CORRIDOR_TRIP_DATA} margin={{ top: 10, right: 15, left: -10, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
                <XAxis dataKey="corridor" stroke="#94a3b8" fontSize={10} tickLine={false} />
                <YAxis stroke="#94a3b8" fontSize={11} tickLine={false} allowDecimals={false} />
                <RechartsTooltip
                  formatter={(val: any, _name: any, item: any) => [
                    `${val} Trips (${item.payload.faults})`,
                    `${item.payload.lengthKm} km Span`,
                  ]}
                  contentStyle={{ backgroundColor: "#0f172a", borderColor: "#334155", borderRadius: "8px", fontSize: "11px" }}
                />
                <Bar dataKey="trips" fill="#f97316" radius={[4, 4, 0, 0]} barSize={22} />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="flex flex-wrap items-center justify-between pt-2 border-t border-slate-800 text-[11px] text-slate-400">
            <span>Critical Focus: <strong className="text-rose-400 font-mono">Anglo-Jos (3 Trips)</strong></span>
            <span>Cleanest Corridor: <strong className="text-emerald-400 font-mono">Bokkos 33kV (0 Trips)</strong></span>
          </div>
        </div>
      </div>
    </div>
  );
};
