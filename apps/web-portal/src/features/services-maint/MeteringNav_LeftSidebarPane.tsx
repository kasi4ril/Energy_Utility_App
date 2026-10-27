/**
 * Persistent Left-Hand Navigation Sidebar for Metering Services Engineering Unit.
 * Outlines core responsibilities as endpoints matching the Customer Care standard.
 */

import React from "react";
import Link from "next/link";

export type MeteringWorkspace =
  | "analytics"
  | "kct-terminal"
  | "testing-bench"
  | "tamper-desk"
  | "standards";

interface MeteringNav_LeftSidebarPaneProps {
  activeWorkspace: MeteringWorkspace;
  onSelectWorkspace: (workspace: MeteringWorkspace) => void;
  totalMetersFleet?: number;
  pendingKctsCount?: number;
  activeTampersCount?: number;
}

export function MeteringNav_LeftSidebarPane({
  activeWorkspace,
  onSelectWorkspace,
  totalMetersFleet = 14850,
  pendingKctsCount = 3,
  activeTampersCount = 2,
}: MeteringNav_LeftSidebarPaneProps) {
  const navItems = [
    {
      id: "analytics" as MeteringWorkspace,
      label: "Metering Fleet Analytics",
      sublabel: "Population, Types & Band Share",
      icon: "📊",
      badge: "14.8K Fleet",
      badgeColor: "bg-emerald-950/80 text-emerald-400 border-emerald-800/40",
    },
    {
      id: "kct-terminal" as MeteringWorkspace,
      label: "STS Key Change Tokens (KCT)",
      sublabel: "TID Rollover & Tariff Group Code",
      icon: "🔑",
      badge: `${pendingKctsCount} Pending`,
      badgeColor: "bg-cyan-950/80 text-cyan-400 border-cyan-800/40",
    },
    {
      id: "testing-bench" as MeteringWorkspace,
      label: "Meter Recalibration Bench",
      sublabel: "Class 0.5S / 1.0 Accuracy Test",
      icon: "🔬",
      badge: "Certified",
      badgeColor: "bg-slate-800 text-slate-300 border-slate-700",
    },
    {
      id: "tamper-desk" as MeteringWorkspace,
      label: "Bypass & Tamper Investigation",
      sublabel: "Physical Seals & Shunt Inspections",
      icon: "🔍",
      badge: `${activeTampersCount} Active`,
      badgeColor: "bg-rose-950/80 text-rose-400 border-rose-800/40",
    },
    {
      id: "standards" as MeteringWorkspace,
      label: "NERC Metering Code Standards",
      sublabel: "Accuracy Tolerances & Seal Regulations",
      icon: "📘",
      badge: "NERC Code",
      badgeColor: "bg-slate-800 text-slate-300 border-slate-700",
    },
  ];

  return (
    <aside className="w-72 bg-slate-950/95 border-r border-slate-800/80 flex flex-col flex-shrink-0 select-none">
      {/* Top Section Header */}
      <div className="p-3.5 border-b border-slate-800/70 bg-slate-900/30">
        <div className="text-[10px] font-bold uppercase tracking-wider text-cyan-400 font-mono">
          Engineering Services
        </div>
        <div className="text-xs font-semibold text-slate-200 mt-0.5">
          Metering Services Unit
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
                  ? "bg-cyan-950/70 border border-cyan-500/50 text-white shadow-lg shadow-cyan-950/40"
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
                      isActive ? "text-cyan-300" : "text-slate-200 group-hover:text-white"
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

        {/* Live Metering Pulse Box */}
        <div className="mt-4 p-3 rounded-lg bg-slate-900/70 border border-slate-800 text-[11px] space-y-1.5">
          <div className="flex items-center justify-between font-bold text-slate-300 text-[10px] uppercase tracking-wider">
            <span>TID Rollover Pulse</span>
            <span className="w-2 h-2 rounded-full bg-cyan-400 animate-pulse" />
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>TID Rollover Progress:</span>
            <span className="font-mono font-bold text-emerald-400">94.8% Migrated</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Bench Accuracy Pass:</span>
            <span className="font-mono font-bold text-cyan-400">99.4% (&plusmn;0.5%)</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Tamper Seal Breaches:</span>
            <span className="font-mono font-bold text-rose-400">{activeTampersCount} Under Investigation</span>
          </div>
        </div>

        {/* Sister Portals Quick Access */}
        <div className="pt-3 border-t border-slate-800/80 space-y-1">
          <div className="text-[9px] font-bold uppercase tracking-wider text-slate-500 px-2 font-mono">
            Sister Engineering Units
          </div>
          <Link
            href="/portal/services/feeder"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-emerald-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>⚡</span>
              <span>Feeder Maintenance Unit</span>
            </span>
            <span className="text-slate-600">&rarr;</span>
          </Link>
          <Link
            href="/portal/services/transformer"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-amber-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>🛢️</span>
              <span>Transformer Maintenance</span>
            </span>
            <span className="text-slate-600">&rarr;</span>
          </Link>
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
        </div>
      </nav>

      {/* Regulatory Standard Footer */}
      <div className="p-3 border-t border-slate-800/70 bg-slate-950 text-[10px] text-slate-500 text-center font-mono">
        NERC Metering Code &bull; NEMSA Certified
      </div>
    </aside>
  );
}

