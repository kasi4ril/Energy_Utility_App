/**
 * Live Grid Metrics Ticker for NESCO Sovereign Landing Page.
 * Displays real-time frequency, bulk intake, feeder health, and collections.
 */

import React from "react";

interface LandingGridTickerProps {
  frequencyHz?: number;
  hydroMw?: number;
  tcnMw?: number;
  activeFeeders?: string;
  collectionsNgn?: string;
  compliancePct?: number;
}

export function LandingGridTicker({
  frequencyHz = 50.02,
  hydroMw = 48.6,
  tcnMw = 15.0,
  activeFeeders = "28 / 32",
  collectionsNgn = "₦4.85M",
  compliancePct = 96.8,
}: LandingGridTickerProps) {
  const totalBulk = hydroMw + tcnMw;

  return (
    <section className="max-w-7xl w-full mx-auto px-6 py-6 select-none">
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 sm:gap-4">
        {/* Hydro Generation */}
        <div className="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800 hover:border-cyan-500/30 transition shadow-lg">
          <div className="flex items-center justify-between">
            <span className="text-[10px] font-bold uppercase text-slate-400 font-mono">Hydro Generation</span>
            <span className="w-2 h-2 rounded-full bg-cyan-400 animate-pulse" />
          </div>
          <div className="text-xl sm:text-2xl font-black text-cyan-400 mt-1 font-mono">
            {hydroMw.toFixed(1)} MW
          </div>
          <div className="text-[10px] text-slate-500 mt-0.5">Kurra Falls Plants 1–6</div>
        </div>

        {/* TCN Bulk Hook */}
        <div className="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800 hover:border-purple-500/30 transition shadow-lg">
          <div className="flex items-center justify-between">
            <span className="text-[10px] font-bold uppercase text-slate-400 font-mono">TCN Bulk Intake</span>
            <span className="w-2 h-2 rounded-full bg-purple-400" />
          </div>
          <div className="text-xl sm:text-2xl font-black text-purple-400 mt-1 font-mono">
            {tcnMw.toFixed(1)} MW
          </div>
          <div className="text-[10px] text-slate-500 mt-0.5">132/33kV Interconnection</div>
        </div>

        {/* Total Bulk Power */}
        <div className="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800 hover:border-emerald-500/30 transition shadow-lg">
          <div className="flex items-center justify-between">
            <span className="text-[10px] font-bold uppercase text-slate-400 font-mono">Total Bulk Load</span>
            <span className="w-2 h-2 rounded-full bg-emerald-400" />
          </div>
          <div className="text-xl sm:text-2xl font-black text-white mt-1 font-mono">
            {totalBulk.toFixed(1)} MW
          </div>
          <div className="text-[10px] text-emerald-400 mt-0.5">Frequency: {frequencyHz.toFixed(2)} Hz</div>
        </div>

        {/* Feeder Corridors */}
        <div className="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800 hover:border-emerald-500/30 transition shadow-lg">
          <div className="flex items-center justify-between">
            <span className="text-[10px] font-bold uppercase text-slate-400 font-mono">Feeders Online</span>
            <span className="w-2 h-2 rounded-full bg-emerald-400" />
          </div>
          <div className="text-xl sm:text-2xl font-black text-emerald-400 mt-1 font-mono">
            {activeFeeders}
          </div>
          <div className="text-[10px] text-slate-500 mt-0.5">1 on PTW &bull; 1 Rationing</div>
        </div>

        {/* Shift Collections */}
        <div className="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800 hover:border-teal-500/30 transition shadow-lg">
          <div className="flex items-center justify-between">
            <span className="text-[10px] font-bold uppercase text-slate-400 font-mono">Today's Revenue</span>
            <span className="w-2 h-2 rounded-full bg-teal-400" />
          </div>
          <div className="text-xl sm:text-2xl font-black text-teal-400 mt-1 font-mono">
            {collectionsNgn}
          </div>
          <div className="text-[10px] text-slate-500 mt-0.5">91.2% Prepaid STS Share</div>
        </div>

        {/* Audit Compliance */}
        <div className="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800 hover:border-amber-500/30 transition shadow-lg">
          <div className="flex items-center justify-between">
            <span className="text-[10px] font-bold uppercase text-slate-400 font-mono">Audit Index</span>
            <span className="w-2 h-2 rounded-full bg-amber-400" />
          </div>
          <div className="text-xl sm:text-2xl font-black text-amber-400 mt-1 font-mono">
            {compliancePct}%
          </div>
          <div className="text-[10px] text-slate-500 mt-0.5">Zero Tamper Breaches</div>
        </div>
      </div>
    </section>
  );
}

