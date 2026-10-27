/**
 * Persistent Left-Hand Navigation Sidebar for Transformer Maintenance Engineering Unit.
 * Outlines core responsibilities as endpoints matching the Customer Care standard.
 */

import React from "react";
import Link from "next/link";

export type TransformerWorkspace =
  | "analytics"
  | "dga-testing"
  | "overhauls"
  | "proof-of-work"
  | "standards";

interface TransformerNav_LeftSidebarPaneProps {
  activeWorkspace: TransformerWorkspace;
  onSelectWorkspace: (workspace: TransformerWorkspace) => void;
  activeDtsCount?: number;
  criticalOilAlertsCount?: number;
  activeOverhaulsCount?: number;
}

export function TransformerNav_LeftSidebarPane({
  activeWorkspace,
  onSelectWorkspace,
  activeDtsCount = 142,
  criticalOilAlertsCount = 1,
  activeOverhaulsCount = 2,
}: TransformerNav_LeftSidebarPaneProps) {
  const navItems = [
    {
      id: "analytics" as TransformerWorkspace,
      label: "Transformer Fleet Analytics",
      sublabel: "Loading %, Temperature & Health",
      icon: "📊",
      badge: `${activeDtsCount} DTs`,
      badgeColor: "bg-emerald-950/80 text-emerald-400 border-emerald-800/40",
    },
    {
      id: "dga-testing" as TransformerWorkspace,
      label: "DGA & Oil Dielectric Lab",
      sublabel: "BDV Breakdown Testing & Moisture",
      icon: "🧪",
      badge: `${criticalOilAlertsCount} Alert`,
      badgeColor: "bg-amber-950/80 text-amber-400 border-amber-800/40",
    },
    {
      id: "overhauls" as TransformerWorkspace,
      label: "Overhaul & Servicing Queue",
      sublabel: "Silica Gel, Bushings & Purge",
      icon: "🔧",
      badge: `${activeOverhaulsCount} Active`,
      badgeColor: "bg-cyan-950/80 text-cyan-400 border-cyan-800/40",
    },
    {
      id: "proof-of-work" as TransformerWorkspace,
      label: "Proof of Work Desk (GPS/Photos)",
      sublabel: "Haversine Distance & Visual Audit",
      icon: "📸",
      badge: "GPS ±15m",
      badgeColor: "bg-purple-950/80 text-purple-400 border-purple-800/40",
    },
    {
      id: "standards" as TransformerWorkspace,
      label: "IEEE C57 Technical Standards",
      sublabel: "Top-Oil Limits & Winding Insulation",
      icon: "📘",
      badge: "IEEE C57",
      badgeColor: "bg-slate-800 text-slate-300 border-slate-700",
    },
  ];

  return (
    <aside className="w-72 bg-slate-950/95 border-r border-slate-800/80 flex flex-col flex-shrink-0 select-none">
      {/* Top Section Header */}
      <div className="p-3.5 border-b border-slate-800/70 bg-slate-900/30">
        <div className="text-[10px] font-bold uppercase tracking-wider text-amber-400 font-mono">
          Engineering Services
        </div>
        <div className="text-xs font-semibold text-slate-200 mt-0.5">
          Transformer Maintenance Unit
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
                  ? "bg-amber-950/70 border border-amber-500/50 text-white shadow-lg shadow-amber-950/40"
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
                      isActive ? "text-amber-300" : "text-slate-200 group-hover:text-white"
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

        {/* Live DT Fleet Pulse Box */}
        <div className="mt-4 p-3 rounded-lg bg-slate-900/70 border border-slate-800 text-[11px] space-y-1.5">
          <div className="flex items-center justify-between font-bold text-slate-300 text-[10px] uppercase tracking-wider">
            <span>Fleet Health Pulse</span>
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Average Fleet Load:</span>
            <span className="font-mono font-bold text-emerald-400">64.5% Capacity</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Overloaded DTs (&gt;85%):</span>
            <span className="font-mono font-bold text-amber-400">3 Substations</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Dielectric Test Pass Rate:</span>
            <span className="font-mono font-bold text-cyan-400">97.2% (&gt;40kV)</span>
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
            href="/portal/services/metering"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-cyan-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>⏱️</span>
              <span>Metering Services Unit</span>
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
        IEEE C57 / IEC 60076 &bull; Power Transformers
      </div>
    </aside>
  );
}

