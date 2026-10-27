/**
 * Revenue Protection & Energy Theft AI Alerts Console for Commercial Operations.
 * Air-gapped AI Guardian scans detecting unmetered load and tampering anomalies.
 */

import React, { useState } from "react";
import { MOCK_THEFT_ALERTS } from "./mockData";
import { TheftAIAlertItem } from "./types";

export function CommercialTheftAlertsView() {
  const [alerts, setAlerts] = useState<TheftAIAlertItem[]>(MOCK_THEFT_ALERTS);

  const handleDispatchInspection = (id: string) => {
    setAlerts((prev) =>
      prev.map((a) =>
        a.id === id ? { ...a, status: "INSPECTION_DISPATCHED" } : a
      )
    );
    alert(
      "Field Inspection Work Order created and dispatched to Services & Maintenance (Metering Unit) with GPS coordinates."
    );
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>🛡️</span>
            <span>Local AI Guardian &bull; Revenue Protection &amp; Anti-Theft Console</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Air-gapped ML scanning DT mass-energy balances, zero-vending anomalies, and CT phase shunts.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-rose-950 text-rose-400 border border-rose-800">
          Ollama Dual-Brain Watchdog
        </span>
      </div>

      {/* Alerts Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {alerts.map((a) => (
          <div
            key={a.id}
            className="p-5 rounded-xl bg-slate-900/70 border border-slate-800 hover:border-slate-700 transition flex flex-col justify-between space-y-4"
          >
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-mono font-bold text-teal-400">{a.target_id}</span>
                <span
                  className={`text-[10px] font-mono px-2 py-0.5 rounded font-bold border ${
                    a.confidence_score >= 0.9
                      ? "bg-rose-950 text-rose-400 border-rose-800"
                      : "bg-amber-950 text-amber-400 border-amber-800"
                  }`}
                >
                  {(a.confidence_score * 100).toFixed(0)}% AI CONFIDENCE
                </span>
              </div>

              <h4 className="text-sm font-bold text-white">{a.target_name}</h4>
              <div className="text-[10px] text-slate-400 font-mono mt-0.5">
                {a.entity_type} &bull; {a.alert_type.replace(/_/g, " ")}
              </div>

              <p className="text-xs text-slate-300 mt-3 bg-slate-950/70 p-3 rounded-lg border border-slate-800 leading-relaxed font-sans">
                {a.evidence_summary}
              </p>
            </div>

            <div className="pt-3 border-t border-slate-800 flex items-center justify-between">
              <div>
                {a.status === "UNRESOLVED_CRITICAL" && (
                  <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-rose-950 text-rose-400 border border-rose-800 font-bold">
                    UNRESOLVED CRITICAL
                  </span>
                )}
                {a.status === "INSPECTION_DISPATCHED" && (
                  <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-sky-950 text-sky-400 border border-sky-800 font-bold">
                    INSPECTION DISPATCHED
                  </span>
                )}
              </div>

              {a.status === "UNRESOLVED_CRITICAL" && (
                <button
                  onClick={() => handleDispatchInspection(a.id)}
                  className="px-3 py-1.5 rounded-lg bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs transition shadow-lg shadow-rose-950 flex items-center space-x-1"
                >
                  <span>🚨</span>
                  <span>Dispatch Inspection Crew</span>
                </button>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

