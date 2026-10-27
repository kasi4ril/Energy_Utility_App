/**
 * Persistent Left-Hand Navigation Sidebar for Stores & Stockpile Inventory.
 * Outlines core responsibilities as endpoints matching the Customer Care standard.
 */

import React from "react";
import Link from "next/link";
import { StoresWorkspace } from "./storesTypes";

interface StoresNav_LeftSidebarPaneProps {
  activeWorkspace: StoresWorkspace;
  onSelectWorkspace: (workspace: StoresWorkspace) => void;
  pendingRequisitionsCount?: number;
  lowStockItemsCount?: number;
  pendingClaimsCount?: number;
}

export function StoresNav_LeftSidebarPane({
  activeWorkspace,
  onSelectWorkspace,
  pendingRequisitionsCount = 1,
  lowStockItemsCount = 2,
  pendingClaimsCount = 1,
}: StoresNav_LeftSidebarPaneProps) {
  const navItems = [
    {
      id: "analytics" as StoresWorkspace,
      label: "Stockpile & Material Analytics",
      sublabel: "Valuation, Low Stock & Fast-Moving",
      icon: "📦",
      badge: "₦48.5M",
      badgeColor: "bg-amber-950/80 text-amber-400 border-amber-800/40",
    },
    {
      id: "requisitions" as StoresWorkspace,
      label: "Dual-Path Requisitions Desk",
      sublabel: "Path A (In-Stock) vs Path B (MD)",
      icon: "📑",
      badge: `${pendingRequisitionsCount} In-Review`,
      badgeColor: "bg-orange-950/80 text-orange-400 border-orange-800/40",
    },
    {
      id: "waybills" as StoresWorkspace,
      label: "Goods Received & Waybills",
      sublabel: "Supplier Delivery, QC & Binning",
      icon: "🚚",
      badge: "Intake",
      badgeColor: "bg-slate-800 text-slate-300 border-slate-700",
    },
    {
      id: "scrap-returns" as StoresWorkspace,
      label: "Scrap Return & Salvage Register",
      sublabel: "Reverse Logistics & Condemned Assets",
      icon: "♻️",
      badge: "Yard Log",
      badgeColor: "bg-teal-950/80 text-teal-400 border-teal-800/40",
    },
    {
      id: "claims" as StoresWorkspace,
      label: "Staff Claims 3-Way Check",
      sublabel: "Roster + PTW + Field GPS Check",
      icon: "📍",
      badge: `${pendingClaimsCount} To Audit`,
      badgeColor: "bg-purple-950/80 text-purple-400 border-purple-800/40",
    },
  ];

  return (
    <aside className="w-72 bg-slate-950/95 border-r border-slate-800/80 flex flex-col flex-shrink-0 select-none">
      {/* Top Section Header */}
      <div className="p-3.5 border-b border-slate-800/70 bg-slate-900/30">
        <div className="text-[10px] font-bold uppercase tracking-wider text-orange-400 font-mono">
          Warehouse &amp; Stores
        </div>
        <div className="text-xs font-semibold text-slate-200 mt-0.5">
          Inventory &amp; Material Logistics
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
                  ? "bg-orange-950/70 border border-orange-500/50 text-white shadow-lg shadow-orange-950/40"
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
                      isActive ? "text-orange-300" : "text-slate-200 group-hover:text-white"
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

        {/* Live Warehouse Pulse Box */}
        <div className="mt-4 p-3 rounded-lg bg-slate-900/70 border border-slate-800 text-[11px] space-y-1.5">
          <div className="flex items-center justify-between font-bold text-slate-300 text-[10px] uppercase tracking-wider">
            <span>Stock Status Radar</span>
            <span className="w-2 h-2 rounded-full bg-orange-400 animate-pulse" />
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Low Stock Items:</span>
            <span className="font-mono font-bold text-rose-400">{lowStockItemsCount} SKUs</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>Salvage Yard Value:</span>
            <span className="font-mono font-bold text-teal-400">&#8358;2.98M Est.</span>
          </div>
          <div className="flex justify-between text-slate-400 text-[10px]">
            <span>3-Way GPS Checks:</span>
            <span className="font-mono font-bold text-emerald-400">100% PostGIS Lock</span>
          </div>
        </div>

        {/* Sister Portals Quick Access */}
        <div className="pt-3 border-t border-slate-800/80 space-y-1">
          <div className="text-[9px] font-bold uppercase tracking-wider text-slate-500 px-2 font-mono">
            Inter-Department Gateways
          </div>
          <Link
            href="/portal/internal-audit"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-rose-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>⚖️</span>
              <span>Internal Audit Desk</span>
            </span>
            <span className="text-slate-600">&rarr;</span>
          </Link>
          <Link
            href="/portal/services/feeder"
            className="w-full block px-2.5 py-1.5 rounded text-[11px] text-slate-400 hover:text-emerald-300 hover:bg-slate-900/60 transition flex items-center justify-between"
          >
            <span className="flex items-center space-x-1.5">
              <span>⚡</span>
              <span>Services &amp; Maintenance</span>
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
        IAS 2 Inventory &bull; Segregation of Duties
      </div>
    </aside>
  );
}

