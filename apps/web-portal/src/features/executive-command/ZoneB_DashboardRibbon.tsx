/**
 * Zone B: Upper Intelligence Dashboard Ribbon (Executive Command Center).
 * 4 Strategic KPI Cards + Real-time Power & Commercial Pulse.
 */

import React from "react";
import { Card } from "../../components/ui/Card";
import { MetricCard } from "../../types/domain";
import { DashboardSummaryDTO } from "./api";

interface ZoneBProps {
  summary?: DashboardSummaryDTO | null;
  onCaptureSnapshot?: () => void;
}

export const ZoneB_DashboardRibbon: React.FC<ZoneBProps> = ({ summary, onCaptureSnapshot }) => {
  const cards: MetricCard[] = [
    {
      id: "power",
      title: "Total Available Power",
      value: summary ? `${summary.total_available_power_mw.toFixed(1)} MW` : "62.5 MW",
      changePct: 2.1,
      trend: "UP",
      status: "NORMAL",
    },
    {
      id: "feeders",
      title: "Feeders in Service",
      value: summary
        ? `${summary.feeders_summary.in_service_count} / ${summary.feeders_summary.total_feeders}`
        : "28 / 32",
      status: "NORMAL",
    },
    {
      id: "revenue",
      title: "Total Shift Collections",
      value: summary
        ? `₦${(summary.revenue_summary.total_collected_amount / 1000000).toFixed(2)}M`
        : "₦4.85M",
      changePct: 6.4,
      trend: "UP",
      status: "NORMAL",
    },
    {
      id: "theft",
      title: "Open Integrity Alerts",
      value: summary ? `${summary.theft_alerts_open_count} Alerts` : "3 Alerts",
      status: "WARNING",
    },
  ];

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400">
          Zone B: Operational Intelligence Pulse &bull; Executive Command Tier
        </h2>
        {onCaptureSnapshot && (
          <button
            onClick={onCaptureSnapshot}
            className="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-xs font-semibold flex items-center space-x-1.5 transition-colors shadow-lg"
          >
            <span>📸 Capture Operational Snapshot</span>
          </button>
        )}
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
                style={{ width: "85%" }}
              />
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

