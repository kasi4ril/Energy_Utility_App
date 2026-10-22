/**
 * Zone B: Upper Intelligence Dashboard Ribbon (Services & Maintenance).
 * 4 Summary KPI Cards + PWA Field Execution Metrics.
 */

import React from "react";
import { Card } from "../../components/ui/Card";
import { MetricCard } from "../../types/domain";

export const ZoneB_DashboardRibbon: React.FC<{ unitTitle?: string }> = ({ unitTitle = "Grid Maintenance" }) => {
  const cards: MetricCard[] = [
    { id: "active", title: "Active Work Orders", value: "24", changePct: 5.0, trend: "UP", status: "NORMAL" },
    { id: "proof", title: "PWA Proofs Under Review", value: "6", changePct: -25.0, trend: "DOWN", status: "WARNING" },
    { id: "ptw", title: "Active Safety PTWs", value: "4", changePct: 0.0, trend: "NEUTRAL", status: "NORMAL" },
    { id: "mttr", title: "MTTR (Mean Repair Time)", value: "2.4 hrs", changePct: -18.2, trend: "DOWN", status: "NORMAL" },
  ];

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400">
          Zone B: Operational Intelligence Pulse &bull; {unitTitle}
        </h2>
        <span className="text-[11px] text-emerald-400 font-mono">
          Field PWA Proximity Tolerance: Active (Haversine Lock)
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
                style={{ width: "65%" }}
              />
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

