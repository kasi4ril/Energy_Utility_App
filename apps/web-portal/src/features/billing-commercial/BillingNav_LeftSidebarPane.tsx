/**
 * Persistent Left-Hand Navigation Sidebar for Commercial Operations & Billing.
 * Outlines core responsibilities as endpoints matching the Customer Care standard.
 */

import React from "react";
import Link from "next/link";
import { BillingWorkspace } from "./types";

interface BillingNav_LeftSidebarPaneProps {
  activeWorkspace: BillingWorkspace;
  onSelectWorkspace: (workspace: BillingWorkspace) => void;
  pendingDisputesCount?: number;
  criticalTheftAlertsCount?: number;
  activeRefundsCount?: number;
}

export function BillingNav_LeftSidebarPane({
  activeWorkspace,
  onSelectWorkspace,
  pendingDisputesCount = 1,
  criticalTheftAlertsCount = 2,
  activeRefundsCount = 2,
}: BillingNav_LeftSidebarPaneProps) {
  const navItems = [
    {
      id: "analytics" as BillingWorkspace,
      label: "Revenue & Tariff Analytics",
      sublabel: "Prepaid vs Postpaid & Bands A–E",
      icon: "📊",
      badge: "₦4.85M",
      badgeColor: "bg-emerald-950/80 text-emerald-400 border-emerald-800/40",
    },
    {
      id: "vending" as BillingWorkspace,
      label: "STS Vending Terminal",
      sublabel: "20-Digit Token & Fintech Hooks",
      icon: "💳",
      badge: "STS-20",
      badgeColor: "bg-teal-950/80 text-teal-400 border-teal-800/40",
    },
    {
      id: "refunds" as BillingWorkspace,
      label: "Meter Refund Amortization",
      sublabel: "NERC Customer Rebate Schedule",
      icon: "📑",
      badge: `${activeRefundsCount} Active`,
      badgeColor: "bg-slate-800 text-slate-300 border-slate-700",
    },
    {
      id: "disputes" as BillingWorkspace,
      label: "Billing Dispute Desk",
      sublabel: "CT Multipliers & Customer Care Link",
      icon: "⚖️",
      badge: `${pendingDisputesCount} Pending`,
      badgeColor: "bg-amber-950/80 text-amber-400 border-amber-800/40",
    },
    {
      id: "theft-ai" as BillingWorkspace,
      label: "Revenue Protection & Theft AI",
      sublabel: "Local Anomaly Scans & Field Checks",
      icon: "🛡️",
      badge: `${criticalTheftAlertsCount} Alert`,
      badgeColor: "bg-rose-950/80 text-rose-400 border-rose-800/40",
    },
  ];

  return (
    <aside className="w-72 bg-slate-950/95 border-r border-slate-800/80 flex flex-col flex-shrink-0 select-none">
      {/* Top Section Header */}
      <div className="p-3.5 border-b border-slate-800/70 bg-slate-900/30">
        <div className="text-[10px] font-bold uppercase tracking-wider text-teal-400 font-mono">
          Commercial Division
        </div>
        <div className="text-xs font-semibold text-slate-200 mt-0.5">
          Revenue &amp; STS Vending
        </div>
      </div>

      {/* Primary Endpoints Navigation */}
      <nav className="flex-1 p-2 space-y-1.5 overflow-y-auto">
        {navItems.map((item) => {
          const isActive = activeWorkspace === item.id;
          return (
            <button
              key={item.id}
              onClick={() => onSelectWorkspace(item.id)}
              className={`w-full text-left p-2.5 rounded-lg transition-all flex items-start space-x-2.5 group relative ${
                isActive
                  ? "bg-teal-950/70 border border-teal-500/50 text-white shadow-lg shadow-teal-950/40"
                  : "hover:bg-slate-900/80 border border-transparent text-slate-400 hover:text-slate-200"
              }`}
            >
              <span className="text-lg mt-0.5 group-hover:scale-110 transition-transform flex-shrink-0">
                {item.icon}
              </span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between">
                  <span
                    className={`text-xs font-bold truncate ${
                      isActive ? "text-teal-300" : "text-slate-200 group-hover:text-white"
                    }`}
                  >
                    {item.label}
                  </span>
                </div>
                <div className="text-[10px] text-slate-500 truncate mt-0.5">
                  {item.sublabel}
                </div>
              </div>
              {item.badge && (
                <span
                  className={`text-[9px] font-mono px-1.5 py-0.5 rounded border font-semibold flex-shrink-0 self-center ${item.badgeColor}`}
                >
                  {item.badge}
                </span>
              )}
            </button>
          );
        })}

        {/* Live Commercial Pulse Box */}
        <div className="mt-4 p-3 rounded-lg bg-slate-900/70 border border-slate-800 text-[11px] space-y-1.5">
          <div className="flex items-center justify-between font-bold text-slate-300 text-[10px] uppercase tracking-wider">
            <span>Commercial Efficiency</span>
            <span className="w-2 h-2 rounded-full bg-teal-400 animate-pulse" />
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Prepaid STS Ratio:</span>
            <span className="font-mono font-bold text-emerald-400">91.2% Share</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Average Tariff Realized:</span>
            <span className="font-mono font-bold text-teal-400">&#8358;184.50 / kWh</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Refund Amortizations:</span>
            <span className="font-mono font-bold text-slate-300">&#8358;125K / Month</span>
          </div>
        </div>

        {/* Sister Portals Quick Access */}
        <div className="pt-3 border-t border-slate-800/80 space-y-1">
          <div className="text-[9px] font-bold uppercase tracking-wider text-slate-500 px-2 font-mono">
            Inter-Department Gateways
          </div>
          <Link
            href="/portal/customer-care"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-blue-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>📞</span>
              <span>Customer Care (CRM)</span>
            </span>
            <span className="text-slate-600">&rarr;</span>
          </Link>
          <Link
            href="/portal/energy-management"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-cyan-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>🎛️</span>
              <span>Central Dispatch Hub</span>
            </span>
            <span className="text-slate-600">&rarr;</span>
          </Link>
          <Link
            href="/portal/executive-command"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-amber-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>👑</span>
              <span>Executive Command Desk</span>
            </span>
            <span className="text-slate-600">&rarr;</span>
          </Link>
        </div>
      </nav>

      {/* Regulatory Standard Footer */}
      <div className="p-3 border-t border-slate-800/70 bg-slate-950 text-[10px] text-slate-500 text-center font-mono">
        NERC Multi-Year Tariff Order (MYTO)
      </div>
    </aside>
  );
}

