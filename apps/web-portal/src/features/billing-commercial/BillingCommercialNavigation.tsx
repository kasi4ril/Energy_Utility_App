/**
 * Task Navigation Pane for Billing, Commercial Operations & Revenue Protection AI.
 */

import React from "react";

export type BillingTask = "vending-desk" | "theft-radar" | "vending-ledger" | "revenue-analytics";

interface NavProps {
  activeTask: BillingTask;
  onSelectTask: (task: BillingTask) => void;
  txCount: number;
  alertCount: number;
}

export const BillingCommercialNavigation: React.FC<NavProps> = ({
  activeTask,
  onSelectTask,
  txCount,
  alertCount,
}) => {
  const tasks: {
    id: BillingTask;
    label: string;
    description: string;
    icon: string;
    badge?: { count: number | string; variant: "amber" | "emerald" | "sky" | "rose" };
  }[] = [
    {
      id: "vending-desk",
      label: "STS Pre-Payment Vending Desk",
      description: "Generate 20-digit STS tokens",
      icon: "💳",
      badge: { count: "Vend", variant: "emerald" },
    },
    {
      id: "theft-radar",
      label: "Revenue Protection AI Radar",
      description: "Transformer unmetered load alerts",
      icon: "🚨",
      badge:
        alertCount > 0
          ? { count: `${alertCount} Anomaly`, variant: "rose" }
          : undefined,
    },
    {
      id: "vending-ledger",
      label: "Metered Vending Ledger",
      description: "Audit trail of tokens & channels",
      icon: "📜",
      badge: { count: `${txCount} Tokens`, variant: "sky" },
    },
    {
      id: "revenue-analytics",
      label: "Commercial Revenue Analytics",
      description: "Tariff collection & vending split",
      icon: "📊",
      badge: { count: "98.1%", variant: "emerald" },
    },
  ];

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-2 shadow-md">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-2">
        {tasks.map((task) => {
          const isActive = activeTask === task.id;
          return (
            <button
              key={task.id}
              onClick={() => onSelectTask(task.id)}
              className={`flex flex-col text-left p-3 rounded-lg transition-all duration-150 border ${
                isActive
                  ? "bg-slate-800 border-sky-500/70 shadow-lg ring-1 ring-sky-500/30"
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
                  isActive ? "text-sky-400" : "text-slate-500"
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
