/**
 * Persistent Left-Hand Navigation Sidebar for Internal Audit & Compliance.
 * Outlines core responsibilities as endpoints matching the Customer Care standard.
 */

import React from "react";
import Link from "next/link";
import { AuditWorkspace } from "./auditTypes";

interface AuditNav_LeftSidebarPaneProps {
  activeWorkspace: AuditWorkspace;
  onSelectWorkspace: (workspace: AuditWorkspace) => void;
  openFindingsCount?: number;
  pendingClaimsCount?: number;
  unauthorizedScrapCount?: number;
}

export function AuditNav_LeftSidebarPane({
  activeWorkspace,
  onSelectWorkspace,
  openFindingsCount = 1,
  pendingClaimsCount = 1,
  unauthorizedScrapCount = 0,
}: AuditNav_LeftSidebarPaneProps) {
  const navItems = [
    {
      id: "analytics" as AuditWorkspace,
      label: "Audit & Compliance Analytics",
      sublabel: "Risk Heatmap & Variance Tracking",
      icon: "📊",
      badge: "96.8%",
      badgeColor: "bg-emerald-950/80 text-emerald-400 border-emerald-800/40",
    },
    {
      id: "ledger" as AuditWorkspace,
      label: "Cryptographic Ledger Explorer",
      sublabel: "SHA-256 Chained WorkNotes Audit",
      icon: "⛓️",
      badge: "Verified",
      badgeColor: "bg-cyan-950/80 text-cyan-400 border-cyan-800/40",
    },
    {
      id: "claims-audit" as AuditWorkspace,
      label: "3-Way Claims Final Approval",
      sublabel: "Overtime & Travel Authorization",
      icon: "📍",
      badge: `${pendingClaimsCount} Pending`,
      badgeColor: "bg-purple-950/80 text-purple-400 border-purple-800/40",
    },
    {
      id: "scrap-audit" as AuditWorkspace,
      label: "Scrap Disposal Auction Clearance",
      sublabel: "Condemned Asset Disposal Sign-off",
      icon: "⚖️",
      badge: "Cleared",
      badgeColor: "bg-teal-950/80 text-teal-400 border-teal-800/40",
    },
    {
      id: "compliance" as AuditWorkspace,
      label: "NERC Regulatory Monitor",
      sublabel: "Statutory SLA & Safety Standards",
      icon: "🛡️",
      badge: "100% SLA",
      badgeColor: "bg-slate-800 text-slate-300 border-slate-700",
    },
  ];

  return (
    <aside className="w-72 bg-slate-950/95 border-r border-slate-800/80 flex flex-col flex-shrink-0 select-none">
      {/* Top Section Header */}
      <div className="p-3.5 border-b border-slate-800/70 bg-slate-900/30">
        <div className="text-[10px] font-bold uppercase tracking-wider text-rose-400 font-mono">
          Internal Audit Tier
        </div>
        <div className="text-xs font-semibold text-slate-200 mt-0.5">
          Compliance &amp; Risk Governance
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
                  ? "bg-rose-950/70 border border-rose-500/50 text-white shadow-lg shadow-rose-950/40"
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
                      isActive ? "text-rose-300" : "text-slate-200 group-hover:text-white"
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

        {/* Live Audit Pulse Box */}
        <div className="mt-4 p-3 rounded-lg bg-slate-900/70 border border-slate-800 text-[11px] space-y-1.5">
          <div className="flex items-center justify-between font-bold text-slate-300 text-[10px] uppercase tracking-wider">
            <span>Sovereign Risk Monitor</span>
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Open Observations:</span>
            <span className="font-mono font-bold text-amber-400">{openFindingsCount} Item</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Hash Ledger State:</span>
            <span className="font-mono font-bold text-cyan-400">0 Tamper Breaches</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Claims Triangulation:</span>
            <span className="font-mono font-bold text-purple-400">1 Flagged Hearing</span>
          </div>
        </div>

        {/* Sister Portals Quick Access */}
        <div className="pt-3 border-t border-slate-800/80 space-y-1">
          <div className="text-[9px] font-bold uppercase tracking-wider text-slate-500 px-2 font-mono">
            Inter-Department Gateways
          </div>
          <Link
            href="/portal/stores"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-orange-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>📦</span>
              <span>Central Stores &amp; Stockpile</span>
            </span>
            <span className="text-slate-600">&rarr;</span>
          </Link>
          <Link
            href="/portal/billing-commercial"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-teal-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>💳</span>
              <span>Billing &amp; STS Vending</span>
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
        ISA 240 Fraud Standard &bull; Independence
      </div>
    </aside>
  );
}

