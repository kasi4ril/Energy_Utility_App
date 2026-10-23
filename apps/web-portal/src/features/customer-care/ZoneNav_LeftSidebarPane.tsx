/**
 * Persistent Left-Hand Navigation Sidebar Pane for Customer Care Department.
 * Positioned on the left side facing the screen, providing one-click access
 * across the 5 integrated sub-pages per Section 12.1. Strictly <= 350 lines.
 */

import React from "react";
import { CustomerCareTask } from "./types";
import {
  BarChart3,
  UserPlus,
  AlertTriangle,
  KeyRound,
  PhoneCall,
  ChevronRight,
  ShieldAlert,
} from "lucide-react";

interface LeftSidebarProps {
  activeTask: CustomerCareTask;
  onSelectTask: (task: CustomerCareTask) => void;
  openTicketsCount: number;
  awaitingClosureCount: number;
  emergencyCount?: number;
}

interface NavItem {
  id: CustomerCareTask;
  label: string;
  sublabel: string;
  icon: React.ElementType;
  badge?: number;
  badgeColor?: string;
  path: string;
}

export function ZoneNav_LeftSidebarPane({
  activeTask,
  onSelectTask,
  openTicketsCount,
  awaitingClosureCount,
  emergencyCount = 1,
}: LeftSidebarProps) {
  const navItems: NavItem[] = [
    {
      id: "analytics",
      label: "Analytics & Storytelling",
      sublabel: "KPI velocity & Ollama AI summary",
      icon: BarChart3,
      path: "/portal/customer-care/analytics",
    },
    {
      id: "register",
      label: "Customer KYC & Metering",
      sublabel: "Premises GPS, DT & Tariff MAP",
      icon: UserPlus,
      path: "/portal/customer-care/register",
    },
    {
      id: "complain",
      label: "Smart Complaint Desk",
      sublabel: "Auto-Routing & Safety Hazard",
      icon: AlertTriangle,
      badge: openTicketsCount,
      badgeColor: "bg-blue-600 text-white",
      path: "/portal/customer-care/complain",
    },
    {
      id: "metering-kct",
      label: "Token & TID Rollover",
      sublabel: "STS History, KCT 1 & 2 Generator",
      icon: KeyRound,
      path: "/portal/customer-care/metering-kct",
    },
    {
      id: "closure-desk",
      label: "Call Confirmation Gate",
      sublabel: "Mandatory Customer Phone Verification",
      icon: PhoneCall,
      badge: awaitingClosureCount,
      badgeColor: "bg-emerald-600 text-white animate-pulse",
      path: "/portal/customer-care/closure-desk",
    },
  ];

  return (
    <aside className="w-full md:w-64 bg-slate-900 border-r border-slate-800 flex flex-col shrink-0 h-full md:sticky md:top-14 md:h-[calc(100vh-3.5rem)] overflow-y-auto">
      {/* Sidebar Header */}
      <div className="p-3.5 border-b border-slate-800 bg-slate-950/60">
        <div className="text-[11px] font-bold tracking-wider text-slate-400 uppercase">
          Operating Workspaces
        </div>
        <div className="text-xs text-slate-200 font-semibold mt-0.5">
          5-Point Service Delivery Mesh
        </div>
      </div>

      {/* Navigation Links */}
      <nav className="p-2 space-y-1.5 flex-1">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive =
            activeTask === item.id ||
            (item.id === "complain" && activeTask === "triage-queue") ||
            (item.id === "register" && activeTask === "kyc-onboarding") ||
            (item.id === "closure-desk" && activeTask === "crm-directory");

          return (
            <button
              key={item.id}
              onClick={() => onSelectTask(item.id)}
              className={`w-full flex items-center justify-between p-2.5 rounded-lg text-left transition-all group ${
                isActive
                  ? "bg-blue-600/15 border border-blue-500/40 text-blue-300 shadow-sm"
                  : "hover:bg-slate-800/80 text-slate-300 border border-transparent"
              }`}
            >
              <div className="flex items-start space-x-2.5 min-w-0">
                <div
                  className={`mt-0.5 p-1.5 rounded-md ${
                    isActive
                      ? "bg-blue-600 text-white"
                      : "bg-slate-800 text-slate-400 group-hover:text-slate-200"
                  }`}
                >
                  <Icon className="w-4 h-4" />
                </div>
                <div className="truncate">
                  <div
                    className={`text-xs font-semibold truncate ${
                      isActive ? "text-blue-300" : "text-slate-200"
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
                {item.badge !== undefined && item.badge > 0 && (
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
                    isActive ? "text-blue-400 translate-x-0.5" : "text-slate-600"
                  }`}
                />
              </div>
            </button>
          );
        })}
      </nav>

      {/* Emergency Radar Pill at bottom of sidebar */}
      {emergencyCount > 0 && (
        <div className="p-3 m-2 rounded-lg bg-red-950/40 border border-red-500/40 text-red-300 text-xs">
          <div className="flex items-center space-x-2 font-bold text-red-400 mb-1">
            <ShieldAlert className="w-4 h-4 animate-bounce text-red-400" />
            <span>Hazard Alert Active</span>
          </div>
          <p className="text-[10px] text-red-200/80 leading-relaxed">
            1 live wire safety ticket routed under &lt;2-Hour SLA to Feeder Maintenance &amp; HSE.
          </p>
        </div>
      )}

      {/* Regulatory Badge */}
      <div className="p-3 border-t border-slate-800 text-[10px] text-slate-400 text-center font-mono bg-slate-950/40">
        NERC Customer Protection Order &bull; v2026.4
      </div>
    </aside>
  );
}

