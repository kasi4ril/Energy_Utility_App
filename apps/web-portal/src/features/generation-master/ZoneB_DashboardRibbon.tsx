/**
 * Zone B: Upper Intelligence Dashboard Ribbon (Generation Master Hub).
 * 4 Summary KPI Cards + 6-Plant Aggregated Generation MW.
 */

import React from "react";
import { Card } from "../../components/ui/Card";
import { MetricCard } from "../../types/domain";

export const ZoneB_DashboardRibbon: React.FC = () => {
  const cards: MetricCard[] = [
    { id: "gen", title: "Total Output (Plants 1–6)", value: "48.2 MW", changePct: 2.5, trend: "UP", status: "NORMAL" },
    { id: "hydro", title: "Kurra Reservoir Head", value: "24.5 m", changePct: 0.8, trend: "UP", status: "NORMAL" },
    { id: "stn", title: "Plants In Service", value: "5 / 6", unit: "Plants", status: "NORMAL" },
    { id: "roll", title: "Remote Shift Roll-Call", value: "100%", unit: "Synchronized", status: "NORMAL" },
  ];

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400">
          Zone B: Operational Intelligence Pulse &bull; Generation Master Station
        </h2>
        <span className="text-[11px] text-emerald-400 font-mono">
          Google Form A Ingestion Bridge: Connected (HMAC Validated)
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
              <div className="h-full bg-emerald-400" style={{ width: "80%" }} />
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

