/**
 * Task Navigation Pane for Customer Care Department.
 * Allows officers to switch between dedicated operational workspaces.
 */

import React from "react";
import { CustomerCareTask } from "./types";

interface NavProps {
  activeTask: CustomerCareTask;
  onSelectTask: (task: CustomerCareTask) => void;
  openTicketsCount: number;
  awaitingConfirmationCount: number;
  totalCustomersCount: number;
}

export const CustomerCareNavigation: React.FC<NavProps> = ({
  activeTask,
  onSelectTask,
  openTicketsCount,
  awaitingConfirmationCount,
  totalCustomersCount,
}) => {
  const tasks: {
    id: CustomerCareTask;
    label: string;
    description: string;
    icon: string;
    badge?: { count: number | string; variant: "amber" | "emerald" | "sky" | "rose" };
  }[] = [
    {
      id: "triage-queue",
      label: "Active Triage & SLA Queue",
      description: "Live fault queue & 2-step closure",
      icon: "📋",
      badge:
        openTicketsCount > 0
          ? { count: `${openTicketsCount} Open`, variant: "amber" }
          : undefined,
    },
    {
      id: "log-complaint",
      label: "Log Customer Complaint",
      description: "Direct fault & outage intake",
      icon: "📝",
      badge: { count: "Form", variant: "emerald" },
    },
    {
      id: "kyc-onboarding",
      label: "Register New Customer",
      description: "NIN & PostGIS DT GPS Anchor",
      icon: "👤",
      badge: { count: "KYC", variant: "sky" },
    },
    {
      id: "crm-directory",
      label: "CRM & Customer 360",
      description: "Accounts, meters, & debt audit",
      icon: "🔍",
      badge: { count: `${totalCustomersCount} Active`, variant: "emerald" },
    },
    {
      id: "analytics",
      label: "Department Analytics",
      description: "SLA compliance, MTTR & channels",
      icon: "📊",
      badge: { count: "KPIs", variant: "emerald" },
    },
  ];

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-2 shadow-md">
      <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-5 gap-2">
        {tasks.map((task) => {
          const isActive = activeTask === task.id;
          return (
            <button
              key={task.id}
              onClick={() => onSelectTask(task.id)}
              className={`flex flex-col text-left p-3 rounded-lg transition-all duration-150 border ${
                isActive
                  ? "bg-slate-800 border-emerald-500/70 shadow-lg ring-1 ring-emerald-500/30"
                  : "bg-slate-900/50 border-slate-800/80 hover:bg-slate-800/60 hover:border-slate-700 text-slate-400 hover:text-slate-200"
              }`}
            >
              <div className="flex items-center justify-between w-full mb-1">
                <div className="flex items-center gap-2">
                  <span className="text-lg">{task.icon}</span>
                  <span
                    className={`text-xs font-bold ${
                      isActive ? "text-white" : "text-slate-300"
                    }`}
                  >
                    {task.label}
                  </span>
                </div>
                {task.badge && (
                  <span
                    className={`text-[10px] font-bold px-1.5 py-0.5 rounded-full ${
                      task.badge.variant === "amber"
                        ? "bg-amber-500/20 text-amber-300 border border-amber-500/40"
                        : task.badge.variant === "rose"
                        ? "bg-rose-500/20 text-rose-300 border border-rose-500/40"
                        : task.badge.variant === "sky"
                        ? "bg-sky-500/20 text-sky-300 border border-sky-500/40"
                        : "bg-emerald-500/20 text-emerald-300 border border-emerald-500/40"
                    }`}
                  >
                    {task.badge.count}
                  </span>
                )}
              </div>
              <p
                className={`text-[11px] truncate mt-0.5 ${
                  isActive ? "text-emerald-400" : "text-slate-500"
                }`}
              >
                {task.description}
              </p>
            </button>
          );
        })}
      </div>
    </div>
  );
};
