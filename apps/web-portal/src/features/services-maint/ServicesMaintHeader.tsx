/**
 * Department Mission Header for Engineering Services & Maintenance.
 */

import React from "react";

interface HeaderProps {
  unitTitle: string;
  orderCount: number;
  onRefresh?: () => void;
}

export const ServicesMaintHeader: React.FC<HeaderProps> = ({
  unitTitle,
  orderCount,
  onRefresh,
}) => {
  return (
    <div className="bg-gradient-to-r from-slate-900 via-slate-850 to-slate-900 border border-slate-700/80 rounded-2xl p-6 shadow-xl relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(#10b981_1px,transparent_1px)] [background-size:24px_24px] opacity-10 pointer-events-none" />

      <div className="relative z-10 space-y-5">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center space-x-2.5 mb-2">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-emerald-500/10 border border-emerald-500/30 text-emerald-400">
                NERC Regulated Utility &bull; Field Engineering &amp; Operations
              </span>
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-amber-500/10 border border-amber-500/30 text-amber-400">
                Unit S04 &bull; {unitTitle}
              </span>
            </div>
            <h1 className="text-2xl sm:text-3xl font-extrabold text-white tracking-tight flex items-center gap-3">
              <span>Services &amp; Maintenance: {unitTitle}</span>
            </h1>
            <p className="text-xs sm:text-sm text-slate-300 max-w-3xl mt-1.5 leading-relaxed">
              Operational field headquarters for high-voltage grid infrastructure maintenance, emergency feeder outage repairs,
              distribution transformer maintenance, and revenue metering services. Enforces rigorous Permit-To-Work (PTW) safety
              lockout/tagout rules and GPS-anchored proof-of-work before returning tickets to Customer Care.
            </p>
          </div>

          {onRefresh && (
            <button
              onClick={onRefresh}
              className="self-start md:self-auto px-3.5 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-600 text-slate-200 text-xs font-semibold flex items-center gap-2 transition shadow"
            >
              <svg className="w-4 h-4 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span>Refresh Field Queue</span>
            </button>
          )}
        </div>

        {/* 4 Core Mandates */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 pt-2">
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-emerald-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>📸 GPS Proof of Work</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Mandatory before/after photo verification and GPS proximity lock (&le;50m) prior to technical sign-off.
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-amber-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>⚡ PTW Safety Protocol</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Strict isolation and grounding confirmation with System Operators before work orders proceed.
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-sky-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>⏱️ NERC Rapid Response</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              30-minute crew mobilization on line snap hazards and major distribution transformer trips.
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-purple-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>🔄 Customer Revert Gate</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Repaired tickets automatically revert to Customer Care officers for 2-step verbal consumer closure.
            </p>
          </div>
        </div>

        {/* Live Operational Strip */}
        <div className="flex flex-wrap items-center gap-6 pt-3 border-t border-slate-800/80 text-xs">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping" />
            <span className="text-slate-400">Crew Dispatch Status:</span>
            <span className="font-bold text-emerald-400">ACTIVE &bull; FIELD TEAMS MOBILIZED</span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-slate-400">Active Work Orders:</span>
            <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-300 font-bold font-mono">
              {orderCount} Orders
            </span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-slate-400">Mean Time to Repair (MTTR):</span>
            <span className="font-bold text-sky-400">2.8 Hours Avg</span>
          </div>
        </div>
      </div>
    </div>
  );
};
