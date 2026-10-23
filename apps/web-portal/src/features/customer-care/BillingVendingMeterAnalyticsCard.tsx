/**
 * Observations 2 & 3: Postpaid / Prepaid Vending Analytics & Meter Installation Milestones.
 * - Month selector for Postpaid paid/unpaid & Prepaid vended/non-vended
 * - Total meters installed so far (14,850) with query per month/week and category breakdown
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import { MONTHLY_BILLING_VENDING_DATA, METER_INSTALLATION_STATS } from "./analyticsData";
import { CreditCard, Gauge, Calendar, ShieldAlert, CheckCircle2, XCircle } from "lucide-react";

export function BillingVendingMeterAnalyticsCard() {
  const [selectedMonth, setSelectedMonth] = useState<string>("2026-10");
  const [meterQueryView, setMeterQueryView] = useState<"MONTHLY" | "WEEKLY">("MONTHLY");

  const billingData = MONTHLY_BILLING_VENDING_DATA[selectedMonth] || MONTHLY_BILLING_VENDING_DATA["2026-10"];
  const meterStats = METER_INSTALLATION_STATS;

  // Filter installs based on query view
  const visibleInstalls = meterStats.period_installs.filter((item) =>
    meterQueryView === "MONTHLY" ? item.period_label.includes("M") : item.period_label.includes("Week")
  );

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-5 text-xs">
      {/* 1. Observation 2: Postpaid & Prepaid Vending Analytics */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5 shadow-md flex flex-col justify-between">
        <div>
          <div className="flex items-center justify-between pb-3 border-b border-slate-800 mb-4">
            <div className="flex items-center space-x-2">
              <CreditCard className="w-4 h-4 text-emerald-400" />
              <div>
                <h3 className="font-bold text-slate-100 uppercase tracking-wide">
                  Postpaid &amp; Prepaid Revenue Analytics
                </h3>
                <p className="text-[10px] text-slate-400">Linked to Live Customer Billing Database</p>
              </div>
            </div>

            {/* Month Selector (Present / Past) */}
            <select
              value={selectedMonth}
              onChange={(e) => setSelectedMonth(e.target.value)}
              className="bg-slate-950 border border-slate-800 rounded-lg px-2.5 py-1 text-slate-200 text-[11px] font-mono focus:border-blue-500 focus:outline-none cursor-pointer"
            >
              <option value="2026-10">Oct 2026 (Present)</option>
              <option value="2026-09">Sep 2026 (Past)</option>
              <option value="2026-08">Aug 2026 (Past)</option>
              <option value="2026-07">Jul 2026 (Past)</option>
            </select>
          </div>

          <div className="space-y-3.5">
            {/* Postpaid Segment */}
            <div className="bg-slate-950/70 p-3 rounded-lg border border-slate-800/80">
              <div className="flex justify-between items-center text-[11px] font-semibold text-slate-300 mb-2">
                <span>Postpaid Customer Accounts</span>
                <span className="font-mono text-slate-400">
                  Total: {billingData.postpaid_paid_count + billingData.postpaid_unpaid_count}
                </span>
              </div>

              <div className="grid grid-cols-2 gap-2 text-[11px]">
                <div className="bg-emerald-950/40 border border-emerald-500/30 p-2 rounded">
                  <div className="flex items-center space-x-1 text-emerald-400 font-bold">
                    <CheckCircle2 className="w-3.5 h-3.5" />
                    <span>Paid in Full</span>
                  </div>
                  <div className="text-base font-bold text-emerald-300 font-mono mt-0.5">
                    {billingData.postpaid_paid_count}
                  </div>
                  <div className="text-[10px] text-emerald-400/80">
                    ₦{(billingData.postpaid_paid_amount / 1000000).toFixed(1)}M Collected
                  </div>
                </div>

                <div className="bg-red-950/40 border border-red-500/30 p-2 rounded">
                  <div className="flex items-center space-x-1 text-red-400 font-bold">
                    <XCircle className="w-3.5 h-3.5" />
                    <span>Unpaid / Arrears</span>
                  </div>
                  <div className="text-base font-bold text-red-300 font-mono mt-0.5">
                    {billingData.postpaid_unpaid_count}
                  </div>
                  <div className="text-[10px] text-red-400/80">
                    ₦{(billingData.postpaid_unpaid_amount / 1000000).toFixed(1)}M Outstanding
                  </div>
                </div>
              </div>
            </div>

            {/* Prepaid STS Vending Segment */}
            <div className="bg-slate-950/70 p-3 rounded-lg border border-slate-800/80">
              <div className="flex justify-between items-center text-[11px] font-semibold text-slate-300 mb-2">
                <span>Prepaid STS Metered Accounts</span>
                <span className="font-mono text-slate-400">
                  Total: {billingData.prepaid_vended_count + billingData.prepaid_unvended_count}
                </span>
              </div>

              <div className="grid grid-cols-2 gap-2 text-[11px]">
                <div className="bg-blue-950/40 border border-blue-500/30 p-2 rounded">
                  <div className="flex items-center space-x-1 text-blue-400 font-bold">
                    <CheckCircle2 className="w-3.5 h-3.5" />
                    <span>Vended This Month</span>
                  </div>
                  <div className="text-base font-bold text-blue-300 font-mono mt-0.5">
                    {billingData.prepaid_vended_count.toLocaleString()}
                  </div>
                  <div className="text-[10px] text-blue-400/80">
                    ₦{(billingData.prepaid_vended_amount / 1000000).toFixed(1)}M Generated
                  </div>
                </div>

                <div className="bg-amber-950/40 border border-amber-500/30 p-2 rounded">
                  <div className="flex items-center space-x-1 text-amber-400 font-bold">
                    <ShieldAlert className="w-3.5 h-3.5" />
                    <span>Zero Vend (&gt;30 Days)</span>
                  </div>
                  <div className="text-base font-bold text-amber-300 font-mono mt-0.5">
                    {billingData.prepaid_unvended_count}
                  </div>
                  <div className="text-[10px] text-amber-400/80">Flagged Anti-Bypass Audit</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="text-[10px] text-slate-500 font-mono pt-3 border-t border-slate-800/80 mt-3 flex justify-between">
          <span>Active Tariff: NERC Multi-Year MYTO Rate</span>
          <span>Automatic Ledger Reconciliation</span>
        </div>
      </div>

      {/* 2. Observation 3: Total Meters Installed & Category Query */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5 shadow-md flex flex-col justify-between">
        <div>
          <div className="flex items-center justify-between pb-3 border-b border-slate-800 mb-4">
            <div className="flex items-center space-x-2">
              <Gauge className="w-4 h-4 text-indigo-400" />
              <div>
                <h3 className="font-bold text-slate-100 uppercase tracking-wide">
                  Total Meters Installed So Far
                </h3>
                <p className="text-[10px] text-slate-400">MAP / NMMP Meter Rollout Milestones</p>
              </div>
            </div>

            {/* Total Badge */}
            <div className="bg-indigo-950/80 border border-indigo-500/40 px-2.5 py-1 rounded-lg text-right font-mono">
              <span className="text-[10px] text-indigo-300 block">Total Installed</span>
              <span className="text-sm font-bold text-indigo-200">
                {meterStats.total_installed.toLocaleString()}
              </span>
            </div>
          </div>

          {/* 3 Categories Breakdown: Residential, MD, Industrial */}
          <div className="grid grid-cols-3 gap-2.5 mb-4">
            <div className="bg-slate-950/70 border border-slate-800 p-2.5 rounded-lg text-center">
              <div className="text-[10px] text-slate-400 uppercase font-mono">Residential (R2)</div>
              <div className="text-base font-bold text-slate-100 font-mono mt-0.5">
                {meterStats.residential_count.toLocaleString()}
              </div>
              <div className="text-[9px] text-emerald-400">75.4% of fleet</div>
            </div>

            <div className="bg-slate-950/70 border border-slate-800 p-2.5 rounded-lg text-center">
              <div className="text-[10px] text-slate-400 uppercase font-mono">MD Customers</div>
              <div className="text-base font-bold text-indigo-300 font-mono mt-0.5">
                {meterStats.md_count.toLocaleString()}
              </div>
              <div className="text-[9px] text-indigo-400">High-yield (5.7%)</div>
            </div>

            <div className="bg-slate-950/70 border border-slate-800 p-2.5 rounded-lg text-center">
              <div className="text-[10px] text-slate-400 uppercase font-mono">Industrial (D1/C1)</div>
              <div className="text-base font-bold text-amber-300 font-mono mt-0.5">
                {meterStats.industrial_count.toLocaleString()}
              </div>
              <div className="text-[9px] text-amber-400">Commercial (18.9%)</div>
            </div>
          </div>

          {/* Query Per Month and Week */}
          <div className="bg-slate-950/70 border border-slate-800 p-3 rounded-lg space-y-2">
            <div className="flex items-center justify-between text-[11px]">
              <span className="font-semibold text-slate-300 flex items-center space-x-1.5">
                <Calendar className="w-3.5 h-3.5 text-blue-400" />
                <span>Installation Velocity Feed</span>
              </span>

              <div className="flex gap-1 bg-slate-900 p-0.5 rounded border border-slate-800">
                <button
                  onClick={() => setMeterQueryView("MONTHLY")}
                  className={`px-2 py-0.5 rounded text-[10px] font-semibold transition-all ${
                    meterQueryView === "MONTHLY" ? "bg-indigo-600 text-white" : "text-slate-400"
                  }`}
                >
                  Per Month
                </button>
                <button
                  onClick={() => setMeterQueryView("WEEKLY")}
                  className={`px-2 py-0.5 rounded text-[10px] font-semibold transition-all ${
                    meterQueryView === "WEEKLY" ? "bg-indigo-600 text-white" : "text-slate-400"
                  }`}
                >
                  Per Week
                </button>
              </div>
            </div>

            <div className="space-y-1.5 max-h-36 overflow-y-auto pr-1">
              {visibleInstalls.map((row, idx) => (
                <div
                  key={idx}
                  className="flex items-center justify-between p-2 rounded bg-slate-900/80 border border-slate-800/80 text-[11px] font-mono"
                >
                  <span className="text-slate-200 font-medium">{row.period_label}</span>
                  <div className="flex items-center space-x-3 text-[10px]">
                    <span className="text-slate-400">Res: <strong className="text-slate-200">{row.residential}</strong></span>
                    <span className="text-indigo-400">MD: <strong className="text-indigo-200">{row.md}</strong></span>
                    <span className="text-amber-400">Ind: <strong className="text-amber-200">{row.industrial}</strong></span>
                    <span className="text-emerald-400 font-bold bg-emerald-950 px-1.5 py-0.2 rounded">
                      +{row.total}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="text-[10px] text-slate-500 font-mono pt-3 border-t border-slate-800/80 mt-3 flex justify-between">
          <span>Field Crew: Metering Installation Unit</span>
          <span>100% Sealed &amp; PostGIS Geotagged</span>
        </div>
      </div>
    </div>
  );
}

