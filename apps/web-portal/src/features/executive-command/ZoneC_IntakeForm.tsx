/**
 * Zone C: Main Operational Workspace (Right Pane - 35% Intake Form - Executive Command).
 * Condition-Based Maintenance (CBM) Predictive Trigger Center & Auto Work Orders.
 */

import React from "react";
import { Badge } from "../../components/ui/Badge";
import { PredictiveTriggerDTO } from "./api";

interface IntakeFormProps {
  triggers: PredictiveTriggerDTO[];
  onSelectTrigger?: (trigger: PredictiveTriggerDTO) => void;
}

export const ZoneC_IntakeForm: React.FC<IntakeFormProps> = ({ triggers, onSelectTrigger }) => {
  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-3 border-b border-slate-700/60 mb-4 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white tracking-wide">Predictive CBM Radar</h3>
          <p className="text-xs text-slate-400 mt-0.5">Automated Work Order Dispatch into S04</p>
        </div>
        <span className="text-[10px] text-emerald-400 font-mono">4 Models Active</span>
      </div>

      <div className="flex-1 overflow-y-auto space-y-3 pr-1 text-xs">
        {triggers.length === 0 ? (
          <div className="text-center py-12 text-xs text-slate-500">
            All generation, transformer, and feeder assets operating within safe envelopes.
          </div>
        ) : (
          triggers.map((t) => (
            <div
              key={t.id}
              onClick={() => onSelectTrigger && onSelectTrigger(t)}
              className="p-3.5 rounded-lg bg-slate-900 border border-slate-700/80 hover:border-slate-600 transition-colors cursor-pointer space-y-2"
            >
              <div className="flex items-center justify-between">
                <span className="font-mono text-emerald-400 font-bold text-[11px]">{t.trigger_no}</span>
                <Badge variant={t.risk_level === "CRITICAL" ? "danger" : "warning"}>
                  {t.risk_level} ({t.risk_score}%)
                </Badge>
              </div>

              <div>
                <p className="text-[11px] text-slate-400">
                  Asset: <span className="text-white font-medium">{t.asset_id}</span> ({t.asset_type})
                </p>
                <p className="text-[11px] text-slate-300 mt-1">{t.recommended_action}</p>
              </div>

              {t.auto_work_order_created && (
                <div className="pt-2 border-t border-slate-800 flex items-center justify-between text-[10px]">
                  <span className="text-emerald-400 font-semibold">⚡ Auto Work Order Dispatched:</span>
                  <span className="font-mono text-slate-300">{t.linked_work_order_id || "Active"}</span>
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

