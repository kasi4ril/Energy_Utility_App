/**
 * Department Mission Header for Customer Care & Consumer Relations.
 * Clearly articulates departmental mission, regulatory mandate, and live SLA pulse.
 */

import React from "react";

interface HeaderProps {
  activeTicketsCount: number;
  attendedAwaitingCloseCount: number;
  onRefresh?: () => void;
}

export const CustomerCareHeader: React.FC<HeaderProps> = ({
  activeTicketsCount,
  attendedAwaitingCloseCount,
  onRefresh,
}) => {
  return (
    <div className="bg-gradient-to-r from-slate-900 via-slate-850 to-slate-900 border border-slate-700/80 rounded-2xl p-6 shadow-xl relative overflow-hidden">
      {/* Background Accent Grid */}
      <div className="absolute inset-0 bg-[radial-gradient(#10b981_1px,transparent_1px)] [background-size:24px_24px] opacity-10 pointer-events-none" />

      <div className="relative z-10 space-y-5">
        {/* Top Badges & Title */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center space-x-2.5 mb-2">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-emerald-500/10 border border-emerald-500/30 text-emerald-400">
                NERC Regulated Utility &bull; Distribution &amp; Retail Operations
              </span>
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-sky-500/10 border border-sky-500/30 text-sky-400">
                Unit S03 &bull; Frontline Command
              </span>
            </div>
            <h1 className="text-2xl sm:text-3xl font-extrabold text-white tracking-tight flex items-center gap-3">
              <span>Consumer Relations &amp; Customer Care Department</span>
            </h1>
            <p className="text-xs sm:text-sm text-slate-300 max-w-3xl mt-1.5 leading-relaxed">
              The primary operational gateway for electricity consumers across NESCO distribution zones.
              Mandated with rapid omnichannel fault logging, anti-ghost KYC identity verification,
              PostGIS transformer-meter anchoring, CRM 360 account management, and enforcing the strict
              2-step customer call confirmation gate before any service ticket is marked closed.
            </p>
          </div>

          {onRefresh && (
            <button
              onClick={onRefresh}
              className="self-start md:self-auto px-3.5 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-600 text-slate-200 text-xs font-semibold flex items-center gap-2 transition-all shadow"
            >
              <svg className="w-4 h-4 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span>Refresh Operational Telemetry</span>
            </button>
          )}
        </div>

        {/* 4 Core Departmental Mandates */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 pt-2">
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-emerald-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>🛡️ 2-Step Closure Gate</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Zero unconfirmed closures. Field maintenance fixes require customer verbal confirmation before closure.
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-amber-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>⏱️ NERC SLA Mandates</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Max 4h response on critical feeder faults. Automated countdown timers with real-time escalation.
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-sky-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>📍 Anti-Ghost KYC Anchoring</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Mandatory NIN validation and PostGIS GPS radius lock to distribution transformer stations.
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-purple-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>🔍 CRM 360 &amp; Debt Shield</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Unified customer directory with live ledger balance audit, meter serial validation, and vending history.
            </p>
          </div>
        </div>

        {/* Live Operational Quick Stats */}
        <div className="flex flex-wrap items-center gap-6 pt-3 border-t border-slate-800/80 text-xs">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping" />
            <span className="text-slate-400">Queue Status:</span>
            <span className="font-bold text-emerald-400">ONLINE &bull; DISPATCH ACTIVE</span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-slate-400">Active Triage Tickets:</span>
            <span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-300 font-bold font-mono">
              {activeTicketsCount}
            </span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-slate-400">Awaiting 2-Step Customer Call:</span>
            <span className="px-2 py-0.5 rounded bg-rose-500/10 text-rose-300 font-bold font-mono">
              {attendedAwaitingCloseCount}
            </span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-slate-400">SLA Compliance (30-day):</span>
            <span className="font-bold text-emerald-400">96.2% On-Time</span>
          </div>
        </div>
      </div>
    </div>
  );
};

