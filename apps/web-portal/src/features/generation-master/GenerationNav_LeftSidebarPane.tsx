/**
 * Persistent Left-Hand Navigation Sidebar Pane for Generation Department & Master Station Command.
 * Positions direct endpoints for core generation activities:
 * 1. Analytics & Telemetry Summary
 * 2. Hourly Reading Sheet Form
 * 3. Change of Shift Form (Immutable Notes)
 * 4. Supervisor Requisition & Notes
 */

import React from "react";
import Link from "next/link";
import { GenerationTask } from "./types";
import {
  Activity,
  ClipboardList,
  Users,
  FileText,
  Radio,
  Gauge,
  Wrench,
  AlertOctagon,
  ChevronRight,
  ShieldCheck,
  ArrowLeft,
  BarChart3,
} from "lucide-react";

interface LeftSidebarProps {
  activeTask: GenerationTask;
  onSelectTask: (task: GenerationTask) => void;
  activePlantsCount: number;
  openFaultsCount: number;
  pendingRequisitionsCount?: number;
}

interface NavItem {
  id: GenerationTask;
  label: string;
  sublabel: string;
  icon: React.ElementType;
  badge?: string | number;
  badgeColor?: string;
  path: string;
}

export function GenerationNav_LeftSidebarPane({
  activeTask,
  onSelectTask,
  activePlantsCount,
  openFaultsCount,
  pendingRequisitionsCount = 2,
}: LeftSidebarProps) {
  // 4 Primary Activities requested by NESCO Operations
  const primaryActivities: NavItem[] = [
    {
      id: "analytics",
      label: "Real-Time Analytics",
      sublabel: "Water head, MW vs capacity & health",
      icon: Activity,
      badge: `${activePlantsCount}/6 Online`,
      badgeColor: "bg-cyan-500/20 text-cyan-300 border border-cyan-500/30",
      path: "/portal/generation-master/analytics",
    },
    {
      id: "reading-sheet",
      label: "Hourly Reading Sheet",
      sublabel: "Operator form straight to DB",
      icon: ClipboardList,
      badge: "Log Due",
      badgeColor: "bg-emerald-500/20 text-emerald-300 border border-emerald-500/30",
      path: "/portal/generation-master/reading-sheet",
    },
    {
      id: "change-of-shift",
      label: "Change of Shift Desk",
      sublabel: "Muster & locked handover notes",
      icon: Users,
      badge: "Locked",
      badgeColor: "bg-amber-500/20 text-amber-300 border border-amber-500/30",
      path: "/portal/generation-master/change-of-shift",
    },
    {
      id: "requisition",
      label: "Supervisor Requisition",
      sublabel: "Direct notes to Unit Head & GM",
      icon: FileText,
      badge: pendingRequisitionsCount > 0 ? `${pendingRequisitionsCount} Active` : undefined,
      badgeColor: "bg-purple-500/20 text-purple-300 border border-purple-500/30",
      path: "/portal/generation-master/requisition",
    },
  ];

  // Secondary Auxiliary Station Views
  const auxiliaryViews: NavItem[] = [
    {
      id: "fleet-board",
      label: "6-Plant Fleet SCADA",
      sublabel: "Live synoptic telemetry board",
      icon: Gauge,
      path: "/portal/generation-master?view=fleet-board",
    },
    {
      id: "workshop",
      label: "Central Machine Workshop",
      sublabel: "Turbine lathe & stator winding",
      icon: Wrench,
      path: "/portal/generation-master?view=workshop",
    },
  ];

  return (
    <aside className="w-full md:w-64 bg-slate-900 border-r border-slate-800 flex flex-col shrink-0 h-full md:sticky md:top-14 md:h-[calc(100vh-3.5rem)] overflow-y-auto">
      {/* Sidebar Header */}
      <div className="p-3.5 border-b border-slate-800 bg-slate-950/60">
        <div className="flex items-center space-x-2 text-[11px] font-bold tracking-wider text-cyan-400 uppercase">
          <Radio className="w-3.5 h-3.5 animate-pulse" />
          <span>Master Station Operations</span>
        </div>
        <div className="text-xs text-slate-200 font-semibold mt-0.5">
          Generation Department Hub
        </div>
      </div>

      {/* Primary Department Activities (Requested 4 Endpoints) */}
      <div className="px-3 pt-3 pb-1 text-[10px] font-bold uppercase tracking-wider text-slate-500">
        Core Activities
      </div>
      <nav className="p-2 space-y-1.5 flex-1">
        {primaryActivities.map((item) => {
          const isActive = activeTask === item.id;
          const Icon = item.icon;

          return (
            <button
              key={item.id}
              onClick={() => onSelectTask(item.id)}
              className={`w-full text-left p-2.5 rounded-lg transition-all flex items-start justify-between group border ${
                isActive
                  ? "bg-slate-800/90 border-cyan-500/70 shadow-md ring-1 ring-cyan-500/30 text-white"
                  : "bg-slate-900/40 border-slate-800/80 hover:bg-slate-800/50 hover:border-slate-700 text-slate-300"
              }`}
            >
              <div className="flex items-start space-x-2.5 min-w-0">
                <div
                  className={`p-1.5 rounded-md mt-0.5 shrink-0 ${
                    isActive
                      ? "bg-cyan-500/20 text-cyan-400 border border-cyan-500/40"
                      : "bg-slate-800 text-slate-400 group-hover:text-cyan-300"
                  }`}
                >
                  <Icon className="w-4 h-4" />
                </div>
                <div className="min-w-0">
                  <div className="flex items-center space-x-1.5">
                    <span className="font-semibold text-xs truncate">
                      {item.label}
                    </span>
                  </div>
                  <div className="text-[10px] text-slate-400 truncate mt-0.5">
                    {item.sublabel}
                  </div>
                </div>
              </div>

              <div className="flex flex-col items-end space-y-1 shrink-0 ml-1">
                {item.badge && (
                  <span
                    className={`text-[9px] font-bold px-1.5 py-0.5 rounded-full ${item.badgeColor}`}
                  >
                    {item.badge}
                  </span>
                )}
                <ChevronRight
                  className={`w-3.5 h-3.5 transition-transform ${
                    isActive ? "text-cyan-400 translate-x-0.5" : "text-slate-600 opacity-0 group-hover:opacity-100"
                  }`}
                />
              </div>
            </button>
          );
        })}

        {/* Secondary Auxiliary Views */}
        <div className="pt-3 pb-1 px-1 text-[10px] font-bold uppercase tracking-wider text-slate-500">
          Cascade Facilities
        </div>
        {auxiliaryViews.map((item) => {
          const isActive = activeTask === item.id;
          const Icon = item.icon;

          return (
            <button
              key={item.id}
              onClick={() => onSelectTask(item.id)}
              className={`w-full text-left p-2 rounded-lg transition-all flex items-center justify-between group border ${
                isActive
                  ? "bg-slate-800/90 border-cyan-500/70 shadow-md text-white"
                  : "bg-slate-900/40 border-slate-800/80 hover:bg-slate-800/50 hover:border-slate-700 text-slate-400 hover:text-slate-200"
              }`}
            >
              <div className="flex items-center space-x-2 min-w-0">
                <Icon className="w-3.5 h-3.5 text-slate-400 group-hover:text-cyan-400 shrink-0" />
                <span className="text-xs font-medium truncate">{item.label}</span>
              </div>
              <ChevronRight className="w-3 h-3 text-slate-600 opacity-0 group-hover:opacity-100" />
            </button>
          );
        })}

        {/* Global Utility Navigation */}
        <div className="pt-3 pb-1 px-1 text-[10px] font-bold uppercase tracking-wider text-slate-500">
          Enterprise Portals
        </div>
        <Link
          href="/"
          className="w-full text-left p-2 rounded-lg transition-all flex items-center justify-between group border bg-slate-900/40 border-slate-800/80 hover:bg-slate-800/50 hover:border-slate-700 text-slate-300 hover:text-white"
          title="Back to General Landing Page"
        >
          <div className="flex items-center space-x-2 min-w-0">
            <ArrowLeft className="w-3.5 h-3.5 text-cyan-400 shrink-0" />
            <span className="text-xs font-medium truncate">General Landing Page</span>
          </div>
          <ChevronRight className="w-3 h-3 text-slate-600 opacity-0 group-hover:opacity-100" />
        </Link>
        <Link
          href="/portal/central-analytics"
          className="w-full text-left p-2 rounded-lg transition-all flex items-center justify-between group border bg-slate-900/40 border-slate-800/80 hover:bg-slate-800/50 hover:border-slate-700 text-emerald-300 hover:text-emerald-200 mt-1"
          title="General Company Analytics & Storyboard for the Day"
        >
          <div className="flex items-center space-x-2 min-w-0">
            <BarChart3 className="w-3.5 h-3.5 text-emerald-400 shrink-0" />
            <span className="text-xs font-medium truncate">General Analytics Page</span>
          </div>
          <ChevronRight className="w-3 h-3 text-slate-600 opacity-0 group-hover:opacity-100" />
        </Link>
      </nav>

      {/* Station Footer Attribution */}
      <div className="p-3 border-t border-slate-800 bg-slate-950/70 space-y-2 text-[11px]">
        <div className="flex items-center space-x-2 text-slate-300">
          <ShieldCheck className="w-4 h-4 text-emerald-400 shrink-0" />
          <div className="leading-tight">
            <div className="font-semibold text-slate-200">NERC Grid Code Sync</div>
            <div className="text-[10px] text-slate-400 font-mono">50.00 Hz Baseload Active</div>
          </div>
        </div>
        <div className="pt-1 text-[10px] text-slate-500 font-mono">
          Cascade Plants 1 to 6 &bull; Sovereign Utility
        </div>
      </div>
    </aside>
  );
}
