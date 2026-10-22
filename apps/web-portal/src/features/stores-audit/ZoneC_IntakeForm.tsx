/**
 * Zone C: Main Operational Workspace (Right Pane - 35% Intake Form - Stores & Audit).
 * Dual-Path Material Requisition Submission Form.
 */

import React, { useState } from "react";
import { createRequisition } from "./api";

export const ZoneC_IntakeForm: React.FC<{ staffId: string; onSuccess?: () => void }> = ({
  staffId,
  onSuccess,
}) => {
  const [sku, setSku] = useState("SKU-DIELECTRIC-OIL-DRUM");
  const [quantity, setQuantity] = useState(5);
  const [workOrderId, setWorkOrderId] = useState("");
  const [purpose, setPurpose] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [msg, setMsg] = useState<{ text: string; type: "success" | "error" } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMsg(null);
    setSubmitting(true);

    try {
      await createRequisition({
        item_sku: sku,
        quantity,
        work_order_id: workOrderId || undefined,
        purpose,
        requester_staff_id: staffId,
      });

      setMsg({
        text: "Requisition Submitted & Routed through Internal Audit Clearance!",
        type: "success",
      });
      setPurpose("");
      if (onSuccess) onSuccess();
    } catch (err: any) {
      setMsg({ text: err.message || "Failed to submit requisition", type: "error" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-3 border-b border-slate-700/60 mb-4">
        <h3 className="text-sm font-bold text-white tracking-wide">Submit Material Requisition</h3>
        <p className="text-xs text-slate-400 mt-0.5">Dual-Path Inventory & Purchase Request</p>
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

      <form onSubmit={handleSubmit} className="space-y-3.5 text-xs flex-1 overflow-y-auto">
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Item SKU</label>
          <select
            value={sku}
            onChange={(e) => setSku(e.target.value)}
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
          >
            <option value="SKU-DIELECTRIC-OIL-DRUM">SKU-DIELECTRIC-OIL-DRUM (Transformer Oil 200L)</option>
            <option value="SKU-11KV-DROP-OUT-FUSE">SKU-11KV-DROP-OUT-FUSE (11kV Expulsion Fuses)</option>
            <option value="SKU-ALUMINUM-CONDUCTOR-100M">SKU-ALUMINUM-CONDUCTOR-100M (100mm ACSR Line)</option>
            <option value="SKU-PIN-INSULATOR-33KV">SKU-PIN-INSULATOR-33KV (33kV Porcelain Insulators)</option>
          </select>
        </div>

        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Requested Quantity</label>
          <input
            type="number"
            min="1"
            required
            value={quantity}
            onChange={(e) => setQuantity(parseInt(e.target.value, 10))}
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500 font-bold"
          />
        </div>

        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">
            Linked Work Order / Asset (Optional)
          </label>
          <input
            type="text"
            value={workOrderId}
            onChange={(e) => setWorkOrderId(e.target.value)}
            placeholder="e.g. WO-PRED-1729600123"
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>

        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Operational Purpose</label>
          <textarea
            required
            rows={3}
            value={purpose}
            onChange={(e) => setPurpose(e.target.value)}
            placeholder="Describe maintenance scope and justification..."
            className="w-full p-2.5 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500 resize-none"
          />
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-semibold rounded-lg text-xs transition-colors"
        >
          {submitting ? "Routing..." : "Submit Requisition"}
        </button>
      </form>
    </div>
  );
};

