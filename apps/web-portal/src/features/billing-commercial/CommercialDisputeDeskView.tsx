/**
 * Billing Dispute & Wrong-Billing Resolution Desk for Commercial Operations.
 * Handles MD CT/PT multiplier error adjustments and ticket reversion to Customer Care.
 */

import React, { useState } from "react";
import { MOCK_DISPUTES } from "./mockData";
import { BillingDisputeItem } from "./types";

export function CommercialDisputeDeskView() {
  const [disputes, setDisputes] = useState<BillingDisputeItem[]>(MOCK_DISPUTES);

  const handleResolve = (id: string, action: "RESOLVED_CREDITED" | "REJECTED_VERIFIED") => {
    setDisputes((prev) =>
      prev.map((d) =>
        d.id === id
          ? {
              ...d,
              status: action,
              resolution_notes:
                action === "RESOLVED_CREDITED"
                  ? "Adjusted CT multiplier in billing engine; credit adjustment applied to customer account."
                  : "Meter audit confirmed CT multiplier and wiring correct; challenge rejected.",
            }
          : d
      )
    );
    alert(
      `Dispute marked ${action}. Ticket status reverted to ATTENDED_TO for Customer Care confirmation.`
    );
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>⚖️</span>
            <span>Billing Dispute &amp; Wrong-Billing Audit Desk</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Inter-departmental loop: Disputed tickets forwarded from Customer Care are reviewed and rectified here.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-amber-950 text-amber-400 border border-amber-800">
          Commercial Service Gate
        </span>
      </div>

      {/* Disputes Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Ticket No</th>
                <th className="py-3 px-3">Customer &bull; Account</th>
                <th className="py-3 px-3">Dispute Category</th>
                <th className="py-3 px-3">Disputed Amount</th>
                <th className="py-3 px-3">Audit Status</th>
                <th className="py-3 px-4 text-right">Commercial Authority</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {disputes.map((d) => (
                <tr key={d.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-teal-400">
                    {d.ticket_number}
                    <span className="block text-[10px] text-slate-500 font-normal">{d.created_at_wat}</span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{d.customer_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{d.account_number}</div>
                  </td>
                  <td className="py-3 px-3 font-semibold text-slate-300">
                    {d.category.replace(/_/g, " ")}
                  </td>
                  <td className="py-3 px-3 font-mono font-bold text-amber-400">
                    &#8358;{d.disputed_amount_ngn.toLocaleString()}
                  </td>
                  <td className="py-3 px-3">
                    {d.status === "PENDING_AUDIT" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                        PENDING AUDIT
                      </span>
                    )}
                    {d.status === "RESOLVED_CREDITED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                        CREDITED &amp; RESOLVED
                      </span>
                    )}
                    {d.status === "REJECTED_VERIFIED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-slate-800 text-slate-300 border border-slate-700">
                        REJECTED
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-4 text-right">
                    {d.status === "PENDING_AUDIT" ? (
                      <div className="flex items-center justify-end space-x-2">
                        <button
                          onClick={() => handleResolve(d.id, "RESOLVED_CREDITED")}
                          className="px-2.5 py-1 rounded bg-emerald-950 hover:bg-emerald-900 text-emerald-300 border border-emerald-800 text-[11px] font-semibold transition"
                        >
                          Approve Credit
                        </button>
                        <button
                          onClick={() => handleResolve(d.id, "REJECTED_VERIFIED")}
                          className="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 border border-slate-700 text-[11px] font-semibold transition"
                        >
                          Reject
                        </button>
                      </div>
                    ) : (
                      <div className="text-[10px] text-slate-400 italic max-w-xs text-right">
                        {d.resolution_notes}
                      </div>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

