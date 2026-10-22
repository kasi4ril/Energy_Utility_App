/**
 * Zone B: Upper Intelligence Dashboard Ribbon (Billing & Commercial Operations).
 * 4 Summary KPI Cards + Revenue Pulse & Meter Refund Metrics.
 */

import React from "react";
import { Card } from "../../components/ui/Card";
import { MetricCard } from "../../types/domain";

export const ZoneB_DashboardRibbon: React.FC = () => {
  const cards: MetricCard[] = [
    { id: "rev", title: "Collections (Current Shift)", value: "₦4,850,000", changePct: 8.4, trend: "UP", status: "NORMAL" },
    { id: "prep", title: "Prepaid Revenue Share", value: "91.2%", changePct: 1.5, trend: "UP", status: "NORMAL" },
    { id: "refund", title: "Amortized Meter Refunds", value: "₦342,000", unit: "Credited", status: "NORMAL" },
    { id: "theft", title: "Open Theft / Bypass Flags", value: "5 Flags", changePct: -20.0, trend: "DOWN", status: "WARNING" },
  ];

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400">
          Zone B: Operational Intelligence Pulse &bull; Commercial Operations
        </h2>
        <span className="text-[11px] text-emerald-400 font-mono">
          Zero-Estimation Metered Engine &bull; NERC STS Active
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
                style={{ width: "82%" }}
              />
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

