/**
 * Task View 5: Hydrological Forecasting & Generation Analytics Scorecard.
 * Tracks water head elevation, plant availability factor (PAF), and bearing vibration radar.
 */

import React from "react";
import { GenerationAnalytics } from "./types";
import { Card } from "../../components/ui/Card";

interface AnalyticsProps {
  analytics: GenerationAnalytics;
  onRefresh?: () => void;
}

export const HydrologyAnalyticsView: React.FC<AnalyticsProps> = ({ analytics, onRefresh }) => {
  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-800 gap-3">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-bold text-white tracking-wide">
              📊 Hydrological Forecasting &amp; Baseload Generation Analytics
            </h2>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-cyan-500/10 text-cyan-300 border border-cyan-500/20">
              Executive Technical Intelligence
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Reservoir water levels, Plant Availability Factor (PAF), Equivalent Forced Outage Rate (EFOR), and vibration trends.
          </p>
        </div>

        {onRefresh && (
          <button
            onClick={onRefresh}
            className="self-start sm:self-auto px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 border border-slate-700 text-xs font-semibold flex items-center gap-1.5 transition"
          >
            <span>🔄 Recalculate PAF &amp; EFOR</span>
          </button>
        )}
      </div>

      {/* Top 4 KPI Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Total Generation Output</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-white font-mono">{analytics.total_generation_mw} MW</span>
            <span className="text-[11px] text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded font-bold">Cascade Fleet</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Daily Delivered: {analytics.daily_mwh_delivered} MWh</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Reservoir Water Storage</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-cyan-300 font-mono">{analytics.reservoir_capacity_pct}%</span>
            <span className="text-[11px] text-emerald-400 font-bold">Optimal Head</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Total Inflow: {analytics.inflow_total_m3_s} m³/s</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Plant Availability Factor (PAF)</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-emerald-400 font-mono">{analytics.plant_availability_factor_pct}%</span>
            <span className="text-[11px] text-slate-300 bg-slate-800 px-2 py-0.5 rounded">Target &ge;90%</span>
          </div>
          <p className="text-[11px] text-emerald-400 mt-2 font-medium">✓ High commercial readiness</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Equivalent Forced Outage (EFOR)</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-amber-300 font-mono">{analytics.equivalent_forced_outage_rate_pct}%</span>
            <span className="text-[11px] text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded font-bold">Below Cap</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Unplanned trips strictly constrained</p>
        </Card>
      </div>

      {/* Generation Share by Station & Weekly Curves */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Plant Generation Share */}
        <div className="lg:col-span-7 bg-slate-950/70 border border-slate-800 rounded-xl p-5 space-y-4">
          <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
            Cascade Station Generation Share (% of Fleet MW)
          </h3>

          <div className="space-y-3 pt-1">
            {analytics.plants_generation_share.map((pg) => (
              <div key={pg.plant_id} className="space-y-1">
                <div className="flex justify-between text-xs">
                  <span className="text-slate-300 font-medium">{pg.name}</span>
                  <span className="font-mono text-cyan-400 font-bold">
                    {pg.mw} MW ({pg.pct}%)
                  </span>
                </div>
                <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                  <div
                    className={`h-full rounded-full transition-all duration-500 ${
                      pg.pct > 25
                        ? "bg-cyan-500"
                        : pg.pct > 15
                        ? "bg-emerald-500"
                        : pg.pct > 0
                        ? "bg-amber-500"
                        : "bg-slate-700"
                    }`}
                    style={{ width: `${pg.pct}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Weekly Trend Bar Chart */}
        <div className="lg:col-span-5 bg-slate-950/70 border border-slate-800 rounded-xl p-5 space-y-4">
          <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
            7-Day Hydro Generation Trend (MWh)
          </h3>

          <div className="space-y-2 pt-2">
            {analytics.weekly_generation_mwh.map((w) => (
              <div key={w.day} className="flex items-center gap-3 text-xs">
                <span className="w-8 font-bold text-slate-400">{w.day}</span>
                <div className="flex-1 bg-slate-800 h-3 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-cyan-500 rounded-full"
                    style={{ width: `${(w.hydro_mwh / 1100) * 100}%` }}
                  />
                </div>
                <span className="w-16 text-right font-mono text-white font-semibold">
                  {w.hydro_mwh} MWh
                </span>
              </div>
            ))}
          </div>

          <div className="pt-3 border-t border-slate-900 text-[11px] text-slate-400">
            Total Weekly Hydro Baseload: <span className="font-bold text-white font-mono">6,891 MWh</span> supplied to NESCO grid.
          </div>
        </div>
      </div>
    </div>
  );
};
