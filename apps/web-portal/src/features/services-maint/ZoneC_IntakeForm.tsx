/**
 * Zone C: Main Operational Workspace (Right Pane - 35% Intake Form - Services & Maintenance).
 * Emergency and Preventive Work Order Dispatch Form.
 */

import React, { useState } from "react";
import { apiClient } from "../../lib/api-client";

interface IntakeFormProps {
  staffId: string;
  unitType?: string;
  onSuccess?: () => void;
}

export const ZoneC_IntakeForm: React.FC<IntakeFormProps> = ({ staffId, unitType = "FEEDER_MAINT", onSuccess }) => {
  const [assetId, setAssetId] = useState("");
  const [teamLeadId, setTeamLeadId] = useState("");
  const [priority, setPriority] = useState("HIGH");
  const [ptwRequired, setPtwRequired] = useState(false);
  const [instructions, setInstructions] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [msg, setMsg] = useState<{ text: string; type: "success" | "error" } | null>(null);

  const handleDispatch = async (e: React.FormEvent) => {
    e.preventDefault();
    setMsg(null);
    setSubmitting(true);

    try {
      await apiClient.post("/maintenance/orders", {
        unit_type: unitType,
        work_order_type: "PREVENTIVE_MAINTENANCE",
        priority,
        unit_head_id: staffId,
        team_lead_id: teamLeadId,
        asset_id: assetId,
        ptw_required: ptwRequired,
        instructions,
      });

      setMsg({ text: "Work Order Dispatched to Field PWA Queue!", type: "success" });
      setAssetId("");
      setInstructions("");
      if (onSuccess) onSuccess();
    } catch (err: any) {
      setMsg({ text: err.message || "Failed to dispatch work order", type: "error" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-3 border-b border-slate-700/60 mb-4">
        <h3 className="text-sm font-bold text-white tracking-wide">Dispatch Maintenance Order</h3>
        <p className="text-xs text-slate-400 mt-0.5">Push Task to Field Crew Mobile PWA</p>
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

      <form onSubmit={handleDispatch} className="space-y-3 text-xs flex-1 overflow-y-auto">
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Grid Asset Code / Feeder</label>
          <input
            type="text"
            required
            value={assetId}
            onChange={(e) => setAssetId(e.target.value)}
            placeholder="e.g. FDR-33KV-VOM or TX-DT-500KVA"
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>

        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Assign Team Lead</label>
          <input
            type="text"
            required
            value={teamLeadId}
            onChange={(e) => setTeamLeadId(e.target.value)}
            placeholder="e.g. LEAD-FDR-01"
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>

        <div className="grid grid-cols-2 gap-2">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">Priority</label>
            <select
              value={priority}
              onChange={(e) => setPriority(e.target.value)}
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
            >
              <option value="CRITICAL">Critical</option>
              <option value="HIGH">High</option>
              <option value="MEDIUM">Medium</option>
            </select>
          </div>
          <div className="flex items-center pt-5">
            <label className="flex items-center space-x-2 text-slate-300 cursor-pointer">
              <input
                type="checkbox"
                checked={ptwRequired}
                onChange={(e) => setPtwRequired(e.target.checked)}
                className="rounded border-slate-700 text-emerald-500 focus:ring-emerald-500"
              />
              <span className="text-[11px] font-semibold">Safety PTW Lock</span>
            </label>
          </div>
        </div>

        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Work Instructions</label>
          <textarea
            required
            rows={3}
            value={instructions}
            onChange={(e) => setInstructions(e.target.value)}
            placeholder="Detail technical scope, feeder isolation requirement, and safety precautions..."
            className="w-full p-2.5 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500 resize-none"
          />
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-semibold rounded-lg text-xs transition-colors"
        >
          {submitting ? "Transmitting..." : "Dispatch Work Order"}
        </button>
      </form>
    </div>
  );
};

