/**
 * Commercial Revenue Analytics View.
 * Displays revenue reconciliation, vending channels, and theft recovery audits.
 */

import React from "react";
import { Card } from "../../components/ui/Card";

export const CommercialRevenueAnalyticsView: React.FC = () => {
  const channelSplit = [
    { name: "OPay Instant Webhook", amount: "₦28,450,000", pct: 45, color: "bg-emerald-500" },
    { name: "Moniepoint / Interswitch POS", amount: "₦17,680,000", pct: 28, color: "bg-sky-500" },
    { name: "Commercial Banking Channels", amount: "₦12,620,000", pct: 20, color: "bg-purple-500" },
    { name: "Over-the-Counter / Walk-In", amount: "₦4,420,000", pct: 7, color: "bg-amber-500" },
  ];

  const tariffPerformance = [
    { band: "Band A (20+ hrs)", revenue: "₦39,170,000", pct: 62 },
    { band: "Band B (16-20 hrs)", revenue: "₦15,160,000", pct: 24 },
    { band: "Band C (12-16 hrs)", revenue: "₦8,840,000", pct: 14 },
  ];

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-800 gap-3">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-bold text-white tracking-wide">
              📊 Commercial Revenue Performance &amp; Settlement Scorecard
            </h2>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              Zero-Estimation Audit
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Real-time collection reconciliation across digital channels, tariff bands, and revenue protection recoveries.
          </p>
        </div>
      </div>

      {/* Top Scorecard Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Total Billed Revenue (MTD)</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-white font-mono">₦63,170,000</span>
          </div>
          <p className="text-[11px] text-emerald-400 mt-2 font-medium">+14.2% vs previous billing cycle</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Collection Efficiency</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-emerald-400 font-mono">98.1%</span>
            <span className="text-[11px] font-bold text-slate-300 bg-slate-800 px-2 py-0.5 rounded">Target &ge;95%</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Zero unmetered billing estimates</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Prepaid STS Share</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-sky-400 font-mono">87.4%</span>
            <span className="text-[11px] font-bold text-sky-400 bg-sky-500/10 px-2 py-0.5 rounded">High STS</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-2">Guaranteed upfront energy collection</p>
        </Card>

        <Card className="bg-slate-950/70 border border-slate-800 p-4">
          <span className="text-[11px] font-semibold text-slate-400">Theft &amp; Tamper Penalties Recovered</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-2xl font-bold text-purple-300 font-mono">₦4,250,000</span>
          </div>
          <p className="text-[11px] text-purple-400 mt-2">Recovered via secondary DT energy audits</p>
        </Card>
      </div>

      {/* Breakdown Grids */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-6 bg-slate-950/70 border border-slate-800 rounded-xl p-5 space-y-4">
          <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
            Payment &amp; Vending Channel Breakdown
          </h3>
          <div className="space-y-3 pt-1">
            {channelSplit.map((ch) => (
              <div key={ch.name} className="space-y-1">
                <div className="flex justify-between text-xs">
                  <span className="text-slate-300 font-medium">{ch.name}</span>
                  <span className="font-mono text-slate-400">{ch.amount} ({ch.pct}%)</span>
                </div>
                <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                  <div className={`h-full ${ch.color} rounded-full`} style={{ width: `${ch.pct}%` }} />
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="lg:col-span-6 bg-slate-950/70 border border-slate-800 rounded-xl p-5 space-y-4">
          <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
            Revenue by NERC Tariff Band
          </h3>
          <div className="space-y-3 pt-1">
            {tariffPerformance.map((tp) => (
              <div key={tp.band} className="space-y-1">
                <div className="flex justify-between text-xs">
                  <span className="text-slate-300 font-medium">{tp.band}</span>
                  <span className="font-mono text-slate-400">{tp.revenue} ({tp.pct}%)</span>
                </div>
                <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                  <div className="h-full bg-amber-500 rounded-full" style={{ width: `${tp.pct}%` }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
