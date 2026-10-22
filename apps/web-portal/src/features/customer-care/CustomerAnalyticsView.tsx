/**
 * Task View 5: Customer Care Departmental Analytics & SLA Scorecard.
 * Comprehensive intelligence on complaints, MTTR/MTTA, channel volume, and feeder hotspots.
 */

import React from "react";
import { CustomerCareAnalytics } from "./types";
import { Card } from "../../components/ui/Card";

interface CustomerAnalyticsViewProps {
  analytics: CustomerCareAnalytics;
  onRefresh?: () => void;
}

export const CustomerAnalyticsView: React.FC<CustomerAnalyticsViewProps> = ({
  analytics,
  onRefresh,
}) => {
  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-800 gap-3">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-bold text-white tracking-wide">
              📊 Operational Analytics &bull; Customer Care Performance Scorecard
            </h2>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-purple-500/10 text-purple-300 border border-purple-500/20">
              Executive Audit
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Real-time NERC SLA metrics, resolution times (MTTA &amp; MTTR), intake channel distribution, and feeder fault hotspots.
          </p>
        </div>

        {onRefresh && (
          <button
            onClick={onRefresh}
            className="self-start sm:self-auto px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 border border-slate-700 text-xs font-semibold flex items-center gap-1.5 transition"
          >
            <span>🔄 Recalculate KPIs</span>
          </button>
        )}
      </div>

      {/* Top 6 KPI Scorecard Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Total Grievances Logged</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-white font-mono">{analytics.total_tickets_logged}</span>
            <span className="text-[11px] font-bold text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded">
              {analytics.closed_tickets_count} Resolved
            </span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Active in triage: {analytics.open_tickets_count} tickets</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">NERC SLA Compliance Rate</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-emerald-400 font-mono">{analytics.sla_compliance_pct}%</span>
            <span className="text-[11px] font-bold text-slate-300 bg-slate-800 px-2 py-0.5 rounded">
              Target &ge;95%
            </span>
          </div>
          <p className="text-[11px] text-emerald-400 mt-2 font-medium">✓ Exceeds regulatory threshold</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">First-Contact Resolution (FCR)</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-sky-400 font-mono">{analytics.first_contact_resolution_pct}%</span>
            <span className="text-[11px] font-bold text-sky-400 bg-sky-500/10 px-2 py-0.5 rounded">
              High Efficiency
            </span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Resolved without secondary field re-dispatch</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Mean Time to Acknowledge (MTTA)</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-amber-300 font-mono">{analytics.average_mtta_minutes} min</span>
            <span className="text-[11px] text-slate-400 font-mono">From logging</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Immediate engineering dispatch to assigned unit</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Mean Time to Repair (MTTR)</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-amber-300 font-mono">{analytics.average_mttr_hours} hrs</span>
            <span className="text-[11px] text-slate-400 font-mono">Average</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Inclusive of fault inspection, PTW, and repair</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Consumer Satisfaction Score (CSAT)</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-purple-300 font-mono">
              ★ {analytics.customer_satisfaction_score} <span className="text-xs text-slate-400 font-normal">/ 5.0</span>
            </span>
            <span className="text-[11px] font-bold text-purple-400 bg-purple-500/10 px-2 py-0.5 rounded">
              Verified Calls
            </span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Derived from 2-step closure feedback logs</p>
        </Card>
      </div>

      {/* Breakdown Section: Category Breakdown + Channel Breakdown */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Fault Category Breakdown */}
        <div className="lg:col-span-7 bg-slate-950/70 border border-slate-800 rounded-xl p-5 space-y-4">
          <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
            Fault Category Distribution &bull; Root Cause Analysis
          </h3>

          <div className="space-y-3 pt-1">
            {analytics.category_distribution.map((cat) => (
              <div key={cat.category} className="space-y-1">
                <div className="flex justify-between text-xs">
                  <span className="text-slate-300 font-medium">{cat.label}</span>
                  <span className="font-mono text-slate-400">
                    {cat.count} cases ({cat.pct}%)
                  </span>
                </div>
                <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                  <div
                    className={`h-full rounded-full transition-all duration-500 ${
                      cat.pct > 25
                        ? "bg-amber-500"
                        : cat.pct > 15
                        ? "bg-sky-500"
                        : "bg-emerald-500"
                    }`}
                    style={{ width: `${cat.pct}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Omnichannel Intake Split */}
        <div className="lg:col-span-5 bg-slate-950/70 border border-slate-800 rounded-xl p-5 space-y-4">
          <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
            Omnichannel Grievance Intake Volume
          </h3>

          <div className="space-y-3 pt-1">
            {analytics.channel_distribution.map((ch) => (
              <div key={ch.channel} className="p-3 rounded-lg bg-slate-900 border border-slate-800/80 flex items-center justify-between">
                <div>
                  <div className="text-xs font-bold text-white">
                    {ch.channel === "CALL_CENTER" && "📞 24/7 Call Center Hotline"}
                    {ch.channel === "WALK_IN" && "🏢 Walk-In Customer Care Desk"}
                    {ch.channel === "ONLINE_PORTAL" && "🌐 Online Web Portal & App"}
                    {ch.channel === "FIELD_AGENT" && "🚶 Field Agents & Revenue Officers"}
                  </div>
                  <div className="text-[11px] text-slate-500 mt-0.5">{ch.count} logged complaints</div>
                </div>
                <div className="text-right">
                  <span className="text-sm font-bold text-emerald-400 font-mono">{ch.pct}%</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Feeder Hotspots */}
      <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-5 space-y-3">
        <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
          Distribution Feeder Fault Hotspots (Priority Inspection Targets)
        </h3>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 pt-1">
          {analytics.hotspot_feeders.map((hf) => (
            <div key={hf.feeder} className="p-3 rounded-lg bg-slate-900 border border-slate-800">
              <div className="text-xs font-bold text-white truncate">{hf.feeder}</div>
              <div className="mt-2 flex justify-between items-center text-[11px]">
                <span className="text-slate-400">Total Faults:</span>
                <span className="font-mono font-bold text-amber-300">{hf.fault_count}</span>
              </div>
              <div className="mt-1 flex justify-between items-center text-[11px]">
                <span className="text-slate-400">Open in Triage:</span>
                <span className="font-mono font-bold text-rose-400">{hf.open_count} active</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
