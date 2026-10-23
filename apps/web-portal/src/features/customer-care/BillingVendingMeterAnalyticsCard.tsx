/**
 * Observations 2 & 3: Postpaid / Prepaid Vending Analytics & Meter Installation Milestones.
 * - Month selector for Postpaid paid/unpaid & Prepaid vended/non-vended
 * - Meter Installation Query Engine (Year & Month dropdowns) with Stacked Category Bar Chart
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import {
  MONTHLY_BILLING_VENDING_DATA,
  METER_INSTALLATION_STATS,
  getMeterInstallationData,
} from "./analyticsData";
import { CreditCard, Gauge, ShieldAlert, CheckCircle2, XCircle, ChevronDown } from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip as RechartsTooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

const METER_MONTHS = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December"
].map((name, i) => ({
  key: (i + 1).toString().padStart(2, "0"),
  name: `${(i + 1).toString().padStart(2, "0")} - ${name}`,
}));

export function BillingVendingMeterAnalyticsCard() {
  const [mounted, setMounted] = useState(false);
  const [selectedMonth, setSelectedMonth] = useState<string>("2026-10");
  const [meterYear, setMeterYear] = useState<string>("2026");
  const [meterMonth, setMeterMonth] = useState<string>("10");

  useEffect(() => {
    setMounted(true);
  }, []);

  const billingData = MONTHLY_BILLING_VENDING_DATA[selectedMonth] || MONTHLY_BILLING_VENDING_DATA["2026-10"];
  const totalFleetInstalled = METER_INSTALLATION_STATS.total_installed;

  // Query meter telemetry based on Year and Month selection
  const meterData = getMeterInstallationData(meterYear, meterMonth);

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
                    {billingData.postpaid_paid_count.toLocaleString()}
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
                    {billingData.postpaid_unpaid_count.toLocaleString()}
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
                    {billingData.prepaid_unvended_count.toLocaleString()}
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

      {/* 2. Observation 3: Total Meters Installed & Year/Month Query with Stacked Bar */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5 shadow-md flex flex-col justify-between">
        <div>
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between pb-3 border-b border-slate-800 mb-3 gap-2">
            <div className="flex items-center space-x-2">
              <Gauge className="w-4 h-4 text-indigo-400" />
              <div>
                <h3 className="font-bold text-slate-100 uppercase tracking-wide">
                  Meter Installation Fleet Intelligence
                </h3>
                <p className="text-[10px] text-indigo-400 font-mono">{meterData.period_label}</p>
              </div>
            </div>

            {/* Total Fleet Badge */}
            <div className="bg-indigo-950/80 border border-indigo-500/40 px-2.5 py-1 rounded-lg text-right font-mono">
              <span className="text-[9px] text-indigo-300 block uppercase">Total Fleet Installed</span>
              <span className="text-sm font-bold text-indigo-200">
                {totalFleetInstalled.toLocaleString()} Meters
              </span>
            </div>
          </div>

          {/* Interactive Year & Month Query Dropdowns */}
          <div className="grid grid-cols-2 gap-2 bg-slate-950/70 p-2 rounded-lg border border-slate-800 mb-3">
            <div className="relative">
              <select
                value={meterYear}
                onChange={(e) => setMeterYear(e.target.value)}
                className="w-full appearance-none bg-slate-900 border border-slate-700 text-slate-200 text-[11px] font-mono rounded px-2.5 py-1.5 pr-6 focus:outline-none focus:border-indigo-500 cursor-pointer"
              >
                <option value="2026">2026 (Current Operating Year)</option>
                <option value="2025">2025 (Historical Audit)</option>
                <option value="2024">2024 (Historical Audit)</option>
                <option value="2023">2023 (Centenary Audit)</option>
              </select>
              <ChevronDown className="w-3 h-3 text-slate-400 absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>

            <div className="relative">
              <select
                value={meterMonth}
                onChange={(e) => setMeterMonth(e.target.value)}
                className="w-full appearance-none bg-slate-900 border border-slate-700 text-slate-200 text-[11px] font-mono rounded px-2.5 py-1.5 pr-6 focus:outline-none focus:border-indigo-500 cursor-pointer"
              >
                <option value="ALL">All Months (Annual Rollout)</option>
                {METER_MONTHS.map((m) => (
                  <option key={m.key} value={m.key}>{m.name}</option>
                ))}
              </select>
              <ChevronDown className="w-3 h-3 text-slate-400 absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none" />
            </div>
          </div>

          {/* 3 Categories Breakdown: Residential, MD, Industrial */}
          <div className="grid grid-cols-3 gap-2 mb-3">
            <div className="bg-slate-950/70 border border-slate-800 p-2 rounded-lg text-center">
              <div className="text-[10px] text-slate-400 uppercase font-mono">Residential (R2)</div>
              <div className="text-base font-bold text-emerald-400 font-mono mt-0.5">
                {meterData.residential_count.toLocaleString()}
              </div>
              <div className="text-[9px] text-slate-500">
                {((meterData.residential_count / meterData.total_installed) * 100).toFixed(1)}% of period
              </div>
            </div>

            <div className="bg-slate-950/70 border border-slate-800 p-2 rounded-lg text-center">
              <div className="text-[10px] text-slate-400 uppercase font-mono">MD Customers</div>
              <div className="text-base font-bold text-indigo-300 font-mono mt-0.5">
                {meterData.md_count.toLocaleString()}
              </div>
              <div className="text-[9px] text-slate-500">
                {((meterData.md_count / meterData.total_installed) * 100).toFixed(1)}% of period
              </div>
            </div>

            <div className="bg-slate-950/70 border border-slate-800 p-2 rounded-lg text-center">
              <div className="text-[10px] text-slate-400 uppercase font-mono">Industrial (D1/C1)</div>
              <div className="text-base font-bold text-amber-300 font-mono mt-0.5">
                {meterData.industrial_count.toLocaleString()}
              </div>
              <div className="text-[9px] text-slate-500">
                {((meterData.industrial_count / meterData.total_installed) * 100).toFixed(1)}% of period
              </div>
            </div>
          </div>

          {/* Stacked Category Bar Chart: Residential, MD, Industrial */}
          {mounted && (
            <div className="bg-slate-950/70 border border-slate-800 p-2.5 rounded-lg">
              <div className="flex justify-between items-center text-[10px] font-mono text-slate-400 mb-1">
                <span>Category Deployment Composition</span>
                <span className="text-emerald-400 font-semibold">+{meterData.total_installed} Energized</span>
              </div>
              <div className="h-32 w-full">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={meterData.chart_data} margin={{ top: 4, right: 10, left: -25, bottom: 0 }}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" opacity={0.6} />
                    <XAxis dataKey="label" stroke="#64748b" fontSize={9} tickLine={false} />
                    <YAxis stroke="#64748b" fontSize={9} tickLine={false} />
                    <RechartsTooltip
                      contentStyle={{
                        backgroundColor: "#090d16",
                        borderColor: "#1e293b",
                        borderRadius: "8px",
                        fontSize: "11px",
                        color: "#f8fafc",
                      }}
                    />
                    <Legend wrapperStyle={{ fontSize: "9px", paddingTop: "2px" }} />
                    <Bar dataKey="residential" stackId="a" fill="#10b981" name="Residential (R2)" />
                    <Bar dataKey="md" stackId="a" fill="#6366f1" name="MD Customers" />
                    <Bar dataKey="industrial" stackId="a" fill="#f59e0b" name="Industrial" radius={[2, 2, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>
          )}
        </div>

        <div className="text-[10px] text-slate-500 font-mono pt-2.5 border-t border-slate-800/80 mt-2 flex justify-between">
          <span>Field Crew: Metering Installation Unit</span>
          <span>100% Sealed &amp; PostGIS Geotagged</span>
        </div>
      </div>
    </div>
  );
}
