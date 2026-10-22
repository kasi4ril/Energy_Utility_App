/**
 * Zone C: Main Operational Workspace (Right Pane - 35% Intake Form - Energy Management).
 * Priority Feeder Load Shedding & Safety PTW Permit Ingestion.
 */

import React, { useState } from "react";
import { executeLoadShedding } from "./api";

export const ZoneC_IntakeForm: React.FC<{ onSuccess?: () => void }> = ({ onSuccess }) => {
  const [shedMW, setShedMW] = useState(5.0);
  const [submitting, setSubmitting] = useState(false);
  const [msg, setMsg] = useState<{ text: string; type: "success" | "error" } | null>(null);

  const handleShed = async (e: React.FormEvent) => {
    e.preventDefault();
    setMsg(null);
    setSubmitting(true);

    try {
      await executeLoadShedding(shedMW);
      setMsg({
        text: `Priority Rationing Executed! ${shedMW} MW shed across lowest priority tiers while preserving PTW safety locks.`,
        type: "success",
      });
      if (onSuccess) onSuccess();
    } catch (err: any) {
      setMsg({ text: err.message || "Failed to execute load shedding", type: "error" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-3 border-b border-slate-700/60 mb-4">
        <h3 className="text-sm font-bold text-white tracking-wide">Priority Grid Rationing</h3>
        <p className="text-xs text-slate-400 mt-0.5">Automated Tariff Band Load Shedding Engine</p>
      </div>

      {msg && (
        <div
          className={`p-3 rounded-lg text-xs mb-4 ${
            msg.type === "success"
              ? "bg-emerald-500/10 border border-emerald-500/30 text-emerald-400"
              : "bg-rose-500/10 border border-rose-500/30 text-rose-400"
          }`}
        >
          {msg.text}
        </div>
      )}

      <form onSubmit={handleShed} className="space-y-4 text-xs flex-1 overflow-y-auto">
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">
            Required Load Drop (MW)
          </label>
          <input
            type="number"
            step="0.5"
            min="0.5"
            required
            value={shedMW}
            onChange={(e) => setShedMW(parseFloat(e.target.value))}
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500 text-sm font-bold"
          />
        </div>

        <div className="p-3 rounded-lg bg-slate-900/60 border border-slate-700/60 space-y-2 text-[11px] text-slate-400">
          <p className="font-semibold text-slate-300">Rationing Sequence:</p>
          <p>1. Low Priority Bands (D & E non-essential industrial)</p>
          <p>2. Mid Priority Bands (C & B residential feeders)</p>
          <p className="text-amber-400 font-medium">
            Protected: Band A premium, Hospital feeders & Lines with Active PTW.
          </p>
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full py-2.5 bg-rose-600 hover:bg-rose-500 disabled:opacity-50 text-white font-semibold rounded-lg text-xs transition-colors shadow-lg"
        >
          {submitting ? "Calculating Cut Sequence..." : "Execute Automated Load Shedding"}
        </button>
      </form>
    </div>
  );
};

