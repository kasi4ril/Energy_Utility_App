/**
 * Zone B: Upper Intelligence Dashboard Ribbon (Customer Care).
 * 4 Summary KPI Cards + Live Resolution Metrics.
 */

import React from "react";
import { Card } from "../../components/ui/Card";
import { MetricCard } from "../../types/domain";

interface ZoneBProps {
  metrics?: MetricCard[];
  onGenerateReport?: () => void;
}

export const ZoneB_DashboardRibbon: React.FC<ZoneBProps> = ({ onGenerateReport }) => {
  const cards: MetricCard[] = [
    { id: "kyc", title: "Total KYC Customers", value: "128,450", changePct: 4.2, trend: "UP", status: "NORMAL" },
    { id: "open", title: "Active Open Tickets", value: "42", changePct: -12.5, trend: "DOWN", status: "NORMAL" },
    { id: "att", title: "Attended-To Confirmation", value: "9", changePct: 15.0, trend: "UP", status: "WARNING" },
    { id: "sla", title: "SLA Resolution Rate", value: "97.4%", unit: "24h", trend: "UP", status: "NORMAL" },
  ];

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h2 className="text-xs font-bold uppercase tracking-wider text-slate-400">
          Zone B: Operational Intelligence Pulse &bull; Customer Care
        </h2>
        {onGenerateReport && (
          <button
            onClick={onGenerateReport}
            className="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-emerald-400 border border-emerald-500/30 rounded-lg text-xs font-semibold flex items-center space-x-1.5 transition-colors"
          >
            <span>✨ Generate Shift Report (Ollama)</span>
          </button>
        )}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map((c) => (
          <Card key={c.id} className="relative overflow-hidden">
            <div className="flex justify-between items-start">
              <div>
                <p className="text-[11px] font-medium text-slate-400">{c.title}</p>
                <div className="mt-1 flex items-baseline space-x-1.5">
                  <span className="text-2xl font-bold text-white tracking-tight">{c.value}</span>
                  {c.unit && <span className="text-xs text-slate-400">{c.unit}</span>}
                </div>
              </div>
              {c.changePct !== undefined && (
                <span
                  className={`text-[11px] font-semibold px-1.5 py-0.5 rounded ${
                    c.changePct >= 0
                      ? "text-emerald-400 bg-emerald-500/10"
                      : "text-rose-400 bg-rose-500/10"
                  }`}
                >
                  {c.changePct >= 0 ? `+${c.changePct}%` : `${c.changePct}%`}
                </span>
              )}
            </div>
            <div className="mt-3 w-full bg-slate-700/50 h-1 rounded-full overflow-hidden">
              <div
                className={`h-full ${c.status === "WARNING" ? "bg-amber-400" : "bg-emerald-400"}`}
                style={{ width: "70%" }}
              />
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

