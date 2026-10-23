/**
 * Task Navigation Pane for Services & Maintenance Department.
 */

import React from "react";

export type ServicesTask =
  | "work-orders"
  | "create-order"
  | "proof-verification"
  | "unit-analytics"
  | "customer-complaints";

interface NavProps {
  activeTask: ServicesTask;
  onSelectTask: (task: ServicesTask) => void;
  orderCount: number;
}

export const ServicesMaintNavigation: React.FC<NavProps> = ({
  activeTask,
  onSelectTask,
  orderCount,
}) => {
  const tasks: {
    id: ServicesTask;
    label: string;
    description: string;
    icon: string;
    badge?: { count: number | string; variant: "amber" | "emerald" | "sky" | "rose" };
  }[] = [
    {
      id: "work-orders",
      label: "Active Work Orders",
      description: "Field crew dispatch queue",
      icon: "📋",
      badge:
        orderCount > 0
          ? { count: `${orderCount} Active`, variant: "amber" }
          : undefined,
    },
    {
      id: "create-order",
      label: "Create Work Order",
      description: "Asset repair & PTW requisition",
      icon: "📝",
      badge: { count: "Form", variant: "emerald" },
    },
    {
      id: "proof-verification",
      label: "Proof of Work Desk",
      description: "Photo audit & GPS verification",
      icon: "🛠️",
      badge: { count: "Verify", variant: "sky" },
    },
    {
      id: "unit-analytics",
      label: "Engineering Analytics",
      description: "MTTR & reliability KPIs",
      icon: "📊",
      badge: { count: "MTTR", variant: "emerald" },
    },
    {
      id: "customer-complaints",
      label: "CC Complaints Feed",
      description: "Feeder, Transformer & Meter tickets",
      icon: "📢",
      badge: { count: "Live", variant: "sky" },
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

