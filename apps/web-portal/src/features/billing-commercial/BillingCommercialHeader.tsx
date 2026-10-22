/**
 * Department Mission Header for Billing, Commercial Operations & Revenue Protection AI.
 */

import React from "react";

interface HeaderProps {
  txCount: number;
  alertCount: number;
  onRefresh?: () => void;
}

export const BillingCommercialHeader: React.FC<HeaderProps> = ({
  txCount,
  alertCount,
  onRefresh,
}) => {
  return (
    <div className="bg-gradient-to-r from-slate-900 via-slate-850 to-slate-900 border border-slate-700/80 rounded-2xl p-6 shadow-xl relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(#0ea5e9_1px,transparent_1px)] [background-size:24px_24px] opacity-10 pointer-events-none" />

      <div className="relative z-10 space-y-5">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center space-x-2.5 mb-2">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-sky-500/10 border border-sky-500/30 text-sky-400">
                NERC Regulated Utility &bull; Commercial Operations &amp; Retail
              </span>
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-amber-500/10 border border-amber-500/30 text-amber-400">
                Unit S08 &bull; Revenue Assurance AI
              </span>
            </div>
            <h1 className="text-2xl sm:text-3xl font-extrabold text-white tracking-tight flex items-center gap-3">
              <span>Billing, Commercial Operations &amp; Revenue Protection</span>
            </h1>
            <p className="text-xs sm:text-sm text-slate-300 max-w-3xl mt-1.5 leading-relaxed">
              Directs commercial energy metering, STS prepayment token generation, tariff band reconciliation,
              and machine learning-driven revenue protection. Deploys automated audits comparing distribution
              transformer bulk energy with consumer smart meters to detect unmetered load surges, meter bypasses,
              and commercial revenue leakages.
            </p>
          </div>

          {onRefresh && (
            <button
              onClick={onRefresh}
              className="self-start md:self-auto px-3.5 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-600 text-slate-200 text-xs font-semibold flex items-center gap-2 transition shadow"
            >
              <svg className="w-4 h-4 text-sky-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span>Refresh Commercial Telemetry</span>
            </button>
          )}
        </div>

        {/* 4 Core Mandates */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 pt-2">
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-emerald-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>🪙 Zero-Estimation Prepayment</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              100% token-based STS vending with transparent tariff-per-kWh conversion and debt deductions.
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-rose-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>🚨 AI Revenue Protection</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Continuous transformer-level energy discrepancy radar (&gt;15% loss flags automated inspection).
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-sky-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>⚡ Sub-Second STS Vending</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Multi-channel token generation for OPay, POS, and Banking webhooks with cryptographic proof.
            </p>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-3.5">
            <div className="flex items-center gap-2 text-amber-400 font-bold text-xs uppercase tracking-wider mb-1">
              <span>⚖️ NERC Capping Compliance</span>
            </div>
            <p className="text-[11px] text-slate-400 leading-normal">
              Automated capping guards to eliminate arbitrary estimation and ensure full regulatory compliance.
            </p>
          </div>
        </div>

        {/* Live Operational Strip */}
        <div className="flex flex-wrap items-center gap-6 pt-3 border-t border-slate-800/80 text-xs">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping" />
            <span className="text-slate-400">Vending Gateway:</span>
            <span className="font-bold text-emerald-400">ONLINE &bull; STS PROTOCOL ACTIVE</span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-slate-400">Processed Transactions:</span>
            <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-300 font-bold font-mono">
              {txCount} Tokens
            </span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-slate-400">Active AI Theft Flags:</span>
            <span className="px-2 py-0.5 rounded bg-rose-500/10 text-rose-300 font-bold font-mono">
              {alertCount} Alerts
            </span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-slate-400">Collection Efficiency:</span>
            <span className="font-bold text-sky-400">98.1% of Invoiced Energy</span>
          </div>
        </div>
      </div>
    </div>
  );
};
