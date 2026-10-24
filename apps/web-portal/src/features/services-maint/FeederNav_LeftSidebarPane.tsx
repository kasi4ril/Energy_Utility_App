/**
 * Persistent Left-Hand Navigation Sidebar Pane for Feeder Maintenance Unit.
 * Positioned on the left side facing the screen (Customer Care 5-point layout standard).
 * Strictly <= 350 lines.
 */

import React from "react";
import { FeederWorkspace } from "./types";
import {
  BarChart3,
  FileText,
  Camera,
  ShieldCheck,
  Wrench,
  ChevronRight,
  ShieldAlert,
  ArrowLeft,
  Activity,
  Layers,
} from "lucide-react";
import Link from "next/link";

interface LeftSidebarProps {
  activeWorkspace: FeederWorkspace;
  onSelectWorkspace: (ws: FeederWorkspace) => void;
  activeOrdersCount: number;
  pendingProofCount: number;
  activeOutageCount?: number;
}

interface NavItem {
  id: FeederWorkspace;
  label: string;
  sublabel: string;
  icon: React.ElementType;
  badge?: number | string;
  badgeColor?: string;
  path: string;
}

export const FeederNav_LeftSidebarPane: React.FC<LeftSidebarProps> = ({
  activeWorkspace,
  onSelectWorkspace,
  activeOrdersCount,
  pendingProofCount,
  activeOutageCount = 1,
}) => {
  const navItems: NavItem[] = [
    {
      id: "analytics",
      label: "Feeder Analytics Dashboard",
      sublabel: "SAIDI, SAIFI, MTTR & Corridor States",
      icon: BarChart3,
      path: "/portal/services/feeder/analytics",
    },
    {
      id: "outage-request",
      label: "Outage Request & Return Register",
      sublabel: "Permit Granting, Return & Closure Log",
      icon: FileText,
      badge: activeOutageCount > 0 ? "1 Active" : undefined,
      badgeColor: "bg-orange-600 text-white animate-pulse",
      path: "/portal/services/feeder/outage-request",
    },
    {
      id: "proof-of-work",
      label: "Proof of Work Submission Desk",
      sublabel: "Anti-Ghost GPS & Before/After Photos",
      icon: Camera,
      badge: pendingProofCount > 0 ? `${pendingProofCount} Due` : undefined,
      badgeColor: "bg-cyan-600 text-white",
      path: "/portal/services/feeder/proof-of-work",
    },
    {
      id: "standards",
      label: "Feeder Technical Standards",
      sublabel: "Patrols, IR Scans & 5-Step Safety Isolation",
      icon: ShieldCheck,
      path: "/portal/services/feeder/standards",
    },
    {
      id: "work-orders",
      label: "Active Feeder Work Orders",
      sublabel: "Field Crew Dispatch & Reversion",
      icon: Wrench,
      badge: activeOrdersCount > 0 ? activeOrdersCount : undefined,
      badgeColor: "bg-slate-700 text-slate-200",
      path: "/portal/services/feeder/work-orders",
    },
  ];

  return (
    <aside className="w-full md:w-64 bg-slate-900 border-r border-slate-800 flex flex-col shrink-0 h-full md:sticky md:top-14 md:h-[calc(100vh-3.5rem)] overflow-y-auto">
      {/* Sidebar Header */}
      <div className="p-3.5 border-b border-slate-800 bg-slate-950/60">
        <div className="text-[11px] font-bold tracking-wider text-orange-400 uppercase">
          Feeder Engineering Ops
        </div>
        <div className="text-xs text-slate-200 font-semibold mt-0.5">
          5-Point Service Delivery Mesh
        </div>
      </div>

      {/* Navigation Links */}
      <nav className="p-2 space-y-1.5 flex-1">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = activeWorkspace === item.id;

          return (
            <button
              key={item.id}
              onClick={() => onSelectWorkspace(item.id)}
              className={`w-full flex items-center justify-between p-2.5 rounded-lg text-left transition-all group ${
                isActive
                  ? "bg-orange-500/15 border border-orange-500/40 text-orange-300 shadow-sm"
                  : "hover:bg-slate-800/80 text-slate-300 border border-transparent"
              }`}
            >
              <div className="flex items-start space-x-2.5 min-w-0">
                <div
                  className={`mt-0.5 p-1.5 rounded-md ${
                    isActive
                      ? "bg-orange-600 text-white"
                      : "bg-slate-800 text-slate-400 group-hover:text-slate-200"
                  }`}
                >
                  <Icon className="w-4 h-4" />
                </div>
                <div className="truncate">
                  <div
                    className={`text-xs font-semibold truncate ${
                      isActive ? "text-orange-300" : "text-slate-200"
                    }`}
                  >
                    {item.label}
                  </div>
                  <div className="text-[10px] text-slate-400 truncate">
                    {item.sublabel}
                  </div>
                </div>
              </div>

              <div className="flex items-center space-x-1 shrink-0 ml-2">
                {item.badge !== undefined && (
                  <span
                    className={`text-[10px] font-bold px-1.5 py-0.5 rounded-full ${
                      item.badgeColor || "bg-slate-700 text-slate-200"
                    }`}
                  >
                    {item.badge}
                  </span>
                )}
                <ChevronRight
                  className={`w-3.5 h-3.5 transition-transform ${
                    isActive ? "text-orange-400 translate-x-0.5" : "text-slate-600"
                  }`}
                />
              </div>
            </button>
          );
        })}
      </nav>

      {/* Sister Units & Quick Links */}
      <div className="p-2 border-t border-slate-800/80 space-y-1">
        <div className="px-2 py-1 text-[10px] font-bold uppercase tracking-wider text-slate-400">
          Services Portals
        </div>
        <Link
          href="/"
          className="flex items-center justify-between p-2 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition"
        >
          <span className="flex items-center space-x-2">
            <ArrowLeft className="w-3.5 h-3.5 text-orange-400" />
            <span>General Landing Page</span>
          </span>
          <span className="text-[10px] text-slate-400 font-mono">/</span>
        </Link>

        <Link
          href="/portal/central-analytics"
          className="flex items-center justify-between p-2 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition"
        >
          <span className="flex items-center space-x-2">
            <BarChart3 className="w-3.5 h-3.5 text-cyan-400" />
            <span>Central Analytics</span>
          </span>
          <span className="text-[10px] text-slate-400 font-mono">/portal</span>
        </Link>

        <Link
          href="/portal/services/transformer"
          className="flex items-center justify-between p-2 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition"
        >
          <span className="flex items-center space-x-2">
            <Layers className="w-3.5 h-3.5 text-blue-400" />
            <span>Transformer Maint Unit</span>
          </span>
          <span className="text-[10px] text-slate-400 font-mono">&rarr;</span>
        </Link>

        <Link
          href="/portal/services/metering"
          className="flex items-center justify-between p-2 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition"
        >
          <span className="flex items-center space-x-2">
            <Activity className="w-3.5 h-3.5 text-purple-400" />
            <span>Metering Field Unit</span>
          </span>
          <span className="text-[10px] text-slate-400 font-mono">&rarr;</span>
        </Link>
      </div>

      {/* Emergency Active Radar Pill at bottom of sidebar */}
      {activeOutageCount > 0 && (
        <div className="p-3 m-2 rounded-lg bg-rose-950/40 border border-rose-500/40 text-rose-300 text-xs">
          <div className="flex items-center space-x-2 font-bold text-rose-400 mb-1">
            <ShieldAlert className="w-4 h-4 animate-bounce text-rose-400" />
            <span>Active Outage Alert</span>
          </div>
          <p className="text-[10px] text-rose-200/80 leading-relaxed">
            Anglo-Jos 33kV corridor isolated at Pole #AJ-142 for emergency pin insulator overhaul under PTW-2026-042.
          </p>
        </div>
      )}

      {/* Regulatory Badge */}
      <div className="p-3 border-t border-slate-800 text-[10px] text-slate-400 text-center font-mono bg-slate-950/40">
        NERC Grid Code Part 3 &bull; Feeder Standards v2026.4
      </div>
    </aside>
  );
};

