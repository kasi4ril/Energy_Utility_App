/**
 * Zone B: Upper Intelligence Dashboard Ribbon (Energy Management & Central Dispatch).
 * 4 Summary KPI Cards + Grid Power Balance Metrics.
 */

import React from "react";
import { Card } from "../../components/ui/Card";
import { MetricCard } from "../../types/domain";

export const ZoneB_DashboardRibbon: React.FC = () => {
  const cards: MetricCard[] = [
    { id: "total_in", title: "Total Available Intake", value: "62.5 MW", changePct: 3.1, trend: "UP", status: "NORMAL" },
    { id: "load", title: "Total System Load", value: "58.1 MW", changePct: 1.2, trend: "UP", status: "NORMAL" },
    { id: "feeders", title: "Feeders In Service", value: "28 / 32", unit: "Feeders", status: "NORMAL" },
    { id: "ptw", title: "Safety PTW Isolations", value: "2 Active", unit: "Permits", status: "WARNING" },
  ];

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400">
          Zone B: Operational Intelligence Pulse &bull; Central Grid Dispatch
        </h2>
        <span className="text-[11px] text-emerald-400 font-mono">
          NEMSA 5-Step Safety PTW Protocol: Active Lock
        </span>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map((c) => (
          <Card key={c.id}>
            <p className="text-[11px] font-medium text-slate-400">{c.title}</p>
            <div className="mt-1 flex items-baseline space-x-1.5">
              <span className="text-2xl font-bold text-white tracking-tight">{c.value}</span>
            </div>
            <div className="mt-3 w-full bg-slate-700/50 h-1 rounded-full overflow-hidden">
              <div
                className={`h-full ${c.status === "WARNING" ? "bg-amber-400" : "bg-emerald-400"}`}
                style={{ width: "88%" }}
              />
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

