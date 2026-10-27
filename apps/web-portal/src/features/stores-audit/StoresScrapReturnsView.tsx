/**
 * Scrap Return & Reverse Logistics Register for Stores & Salvage Yard.
 * Manages decommissioned transformers, copper windings, and aluminum scrap.
 */

import React, { useState } from "react";
import { MOCK_SCRAP_RETURNS } from "./storesMockData";
import { ScrapReturnRecord } from "./storesTypes";

export function StoresScrapReturnsView() {
  const [scraps, setScraps] = useState<ScrapReturnRecord[]>(MOCK_SCRAP_RETURNS);

  const handleClearForAudit = (id: string) => {
    setScraps((prev) =>
      prev.map((s) =>
        s.id === id ? { ...s, status: "AUDIT_CLEARED_FOR_AUCTION" } : s
      )
    );
    alert("Scrap asset weighed and cleared for Internal Audit Auction Verification.");
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>♻️</span>
            <span>Salvage Yard &bull; Reverse Logistics &amp; Scrap Asset Register</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Strict chain of custody for condemned distribution transformers and recovered copper/aluminum.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-teal-950 text-teal-400 border border-teal-800">
          Salvage Custody Standard
        </span>
      </div>

      {/* Scraps Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Ticket No</th>
                <th className="py-3 px-3">Asset Classification</th>
                <th className="py-3 px-3">Originating Dept</th>
                <th className="py-3 px-3">Weighed Weight</th>
                <th className="py-3 px-3">Est. Salvage Value</th>
                <th className="py-3 px-3">Yard Section</th>
                <th className="py-3 px-4 text-right">Status / Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {scraps.map((s) => (
                <tr key={s.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-teal-400">
                    {s.return_ticket_no}
                    <span className="block text-[10px] text-slate-500 font-normal">{s.received_at_wat}</span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{s.asset_type.replace(/_/g, " ")}</div>
                  </td>
                  <td className="py-3 px-3 text-slate-300">{s.source_department}</td>
                  <td className="py-3 px-3 font-mono font-bold text-white">{s.weight_kg} kg</td>
                  <td className="py-3 px-3 font-mono text-amber-400 font-bold">
                    &#8358;{s.estimated_salvage_value_ngn.toLocaleString()}
                  </td>
                  <td className="py-3 px-3 font-mono text-slate-400 text-[11px]">{s.salvage_yard_bin}</td>
                  <td className="py-3 px-4 text-right">
                    {s.status === "LOGGED_IN_YARD" && (
                      <button
                        onClick={() => handleClearForAudit(s.id)}
                        className="px-2.5 py-1 rounded bg-teal-950 hover:bg-teal-900 text-teal-300 border border-teal-800 text-[11px] font-semibold transition"
                      >
                        Send to Audit
                      </button>
                    )}
                    {s.status === "AUDIT_CLEARED_FOR_AUCTION" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                        AUDIT CLEARED
                      </span>
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

