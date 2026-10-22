/**
 * Zone B: Upper Intelligence Dashboard Ribbon (Stores & Internal Audit).
 * 4 Summary KPI Cards + Store Stock Ledger Metrics.
 */

import React from "react";
import { Card } from "../../components/ui/Card";
import { MetricCard } from "../../types/domain";

export const ZoneB_DashboardRibbon: React.FC = () => {
  const cards: MetricCard[] = [
    { id: "skus", title: "Active Store SKUs", value: "1,240", status: "NORMAL" },
    { id: "reorder", title: "Low Stock Critical Items", value: "7 Items", status: "WARNING" },
    { id: "audit", title: "Pending Audit Consents", value: "3", status: "NORMAL" },
    { id: "claims", title: "3-Way Claims Verified", value: "98.2%", status: "NORMAL" },
  ];

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400">
          Zone B: Operational Intelligence Pulse &bull; Stores & Internal Audit
        </h2>
        <span className="text-[11px] text-emerald-400 font-mono">
          @GridAssistant Stock Query Bot: Live Sync
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
                style={{ width: "75%" }}
              />
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

