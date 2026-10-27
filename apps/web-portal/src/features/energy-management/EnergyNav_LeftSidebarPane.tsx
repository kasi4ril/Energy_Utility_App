/**
 * Persistent Left-Hand Navigation Sidebar for Central Energy Management & Transmission Dispatch.
 * Outlines core responsibilities as endpoints matching the Customer Care standard.
 */

import React from "react";
import Link from "next/link";
import { EnergyWorkspace } from "./types";

interface EnergyNav_LeftSidebarPaneProps {
  activeWorkspace: EnergyWorkspace;
  onSelectWorkspace: (workspace: EnergyWorkspace) => void;
  activeFeedersCount?: number;
  sheddedFeedersCount?: number;
  activePTWCount?: number;
  pendingClearancesCount?: number;
}

export function EnergyNav_LeftSidebarPane({
  activeWorkspace,
  onSelectWorkspace,
  activeFeedersCount = 6,
  sheddedFeedersCount = 1,
  activePTWCount = 1,
  pendingClearancesCount = 1,
}: EnergyNav_LeftSidebarPaneProps) {
  const navItems = [
    {
      id: "analytics" as EnergyWorkspace,
      label: "Dispatch & Bulk Intake Analytics",
      sublabel: "Dual-Source MW, Frequency & Voltage",
      icon: "📊",
      badge: "63.6 MW",
      badgeColor: "bg-emerald-950/80 text-emerald-400 border-emerald-800/40",
    },
    {
      id: "load-shedding" as EnergyWorkspace,
      label: "Feeder Rationing & Load Shedding",
      sublabel: "Bands A–E Priority Dispatch Console",
      icon: "⚡",
      badge: `${sheddedFeedersCount} Shedded`,
      badgeColor: "bg-amber-950/80 text-amber-400 border-amber-800/40",
    },
    {
      id: "safety-ptw" as EnergyWorkspace,
      label: "5-Step Safety PTW Isolation",
      sublabel: "NERC Grid Code Substation Safety Lock",
      icon: "🔒",
      badge: `${activePTWCount} Active`,
      badgeColor: "bg-sky-950/80 text-sky-400 border-sky-800/40",
    },
    {
      id: "dispatch-log" as EnergyWorkspace,
      label: "Hourly Dispatch Log Sheet",
      sublabel: "Direct Telemetry Ingestion & Remarks",
      icon: "📝",
      badge: "Hourly",
      badgeColor: "bg-slate-800 text-slate-300 border-slate-700",
    },
    {
      id: "outage-clearance" as EnergyWorkspace,
      label: "Outage Coordination & Clearances",
      sublabel: "Services Feeder Handover Register",
      icon: "📋",
      badge: `${pendingClearancesCount} Due`,
      badgeColor: "bg-purple-950/80 text-purple-400 border-purple-800/40",
    },
  ];

  return (
    <aside className="w-72 bg-slate-950/95 border-r border-slate-800/80 flex flex-col flex-shrink-0 select-none">
      {/* Top Section Header */}
      <div className="p-3.5 border-b border-slate-800/70 bg-slate-900/30">
        <div className="text-[10px] font-bold uppercase tracking-wider text-cyan-400 font-mono">
          Central Dispatch Hub
        </div>
        <div className="text-xs font-semibold text-slate-200 mt-0.5">
          Transmission &amp; Grid Control
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

        {/* Live Grid Alert Box */}
        <div className="mt-4 p-3 rounded-lg bg-slate-900/70 border border-slate-800 text-[11px] space-y-1.5">
          <div className="flex items-center justify-between font-bold text-slate-300 text-[10px] uppercase tracking-wider">
            <span>Grid Balance Pulse</span>
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Feeders In-Service:</span>
            <span className="font-mono font-bold text-emerald-400">{activeFeedersCount} Online</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Active Load Shed:</span>
            <span className="font-mono font-bold text-amber-400">{sheddedFeedersCount} Feeder</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>PTW Safety Isolations:</span>
            <span className="font-mono font-bold text-sky-400">{activePTWCount} Active</span>
          </div>
        </div>

        {/* Sister Portals Quick Access */}
        <div className="pt-3 border-t border-slate-800/80 space-y-1">
          <div className="text-[9px] font-bold uppercase tracking-wider text-slate-500 px-2 font-mono">
            Inter-Department Gateways
          </div>
          <Link
            href="/portal/generation-master"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-cyan-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>🌊</span>
              <span>Generation Fleet (Hydro)</span>
            </span>
            <span className="text-slate-600">&rarr;</span>
          </Link>
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
        NERC Grid Code Part 3 &bull; Transmission
      </div>
    </aside>
  );
}

