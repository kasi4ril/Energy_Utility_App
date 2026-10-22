/**
 * Task Navigation Pane for Generation Department & Master Station Command Hub (S05).
 */

import React from "react";
import { GenerationTask } from "./types";

interface NavProps {
  activeTask: GenerationTask;
  onSelectTask: (task: GenerationTask) => void;
  activePlantsCount: number;
  openFaultsCount: number;
  activeWorkshopJobsCount: number;
}

export const GenerationMasterNavigation: React.FC<NavProps> = ({
  activeTask,
  onSelectTask,
  activePlantsCount,
  openFaultsCount,
  activeWorkshopJobsCount,
}) => {
  const tasks: {
    id: GenerationTask;
    label: string;
    description: string;
    icon: string;
    badge?: { count: number | string; variant: "amber" | "emerald" | "cyan" | "rose" };
  }[] = [
    {
      id: "fleet-board",
      label: "6-Plant Cascade Fleet Board",
      description: "Live MW, Head & Vibration",
      icon: "🌐",
      badge: { count: `${activePlantsCount}/6 In Service`, variant: "cyan" },
    },
    {
      id: "call-collation",
      label: "Hourly Call Collation Desk",
      description: "Shift muster & Google Form A",
      icon: "📻",
      badge: { count: "Log", variant: "emerald" },
    },
    {
      id: "plant-faults",
      label: "Plant Electrical & Mech Faults",
      description: "Root causes & GM approval gate",
      icon: "🚨",
      badge:
        openFaultsCount > 0
          ? { count: `${openFaultsCount} Outage`, variant: "rose" }
          : undefined,
    },
    {
      id: "workshop",
      label: "Central Machine Workshop",
      description: "Turbine lathe & coil winding",
      icon: "⚙️",
      badge:
        activeWorkshopJobsCount > 0
          ? { count: `${activeWorkshopJobsCount} Jobs`, variant: "amber" }
          : undefined,
    },
    {
      id: "analytics",
      label: "Hydrology & Gen Analytics",
      description: "Reservoir head & MWh curves",
      icon: "📊",
      badge: { count: "94.2% PAF", variant: "emerald" },
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
                  ? "bg-slate-800 border-cyan-500/70 shadow-lg ring-1 ring-cyan-500/30"
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
                        : task.badge.variant === "cyan"
                        ? "bg-cyan-500/20 text-cyan-300 border border-cyan-500/40"
                        : "bg-emerald-500/20 text-emerald-300 border border-emerald-500/40"
                    }`}
                  >
                    {task.badge.count}
                  </span>
                )}
              </div>
              <p
                className={`text-[11px] truncate mt-0.5 ${
                  isActive ? "text-cyan-400" : "text-slate-500"
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
