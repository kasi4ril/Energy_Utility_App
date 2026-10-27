/**
 * Scrap Disposal & Salvage Auction Audit Clearance for Internal Audit.
 * Prevents unauthorized salvage sales and certifies public auction revenue.
 */

import React, { useState } from "react";
import { MOCK_SCRAP_AUCTION } from "./auditMockData";
import { ScrapAuctionItem } from "./auditTypes";

export function AuditScrapClearanceView() {
  const [lots, setLots] = useState<ScrapAuctionItem[]>(MOCK_SCRAP_AUCTION);

  const handleApproveAuction = (id: string) => {
    setLots((prev) =>
      prev.map((l) =>
        l.id === id ? { ...l, audit_clearance: "AUCTION_COMPLETED" } : l
      )
    );
    alert("Auction lot audited and disposal proceeds reconciled with Finance.");
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>⚖️</span>
            <span>Salvage Scrap &amp; Condemned Asset Disposal Audit Gate</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Mandatory internal audit verification before any decommissioned copper or transformers leave the yard.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-teal-950 text-teal-400 border border-teal-800">
          Asset Protection Control
        </span>
      </div>

      {/* Auction Lots Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Lot Number</th>
                <th className="py-3 px-3">Asset Description</th>
                <th className="py-3 px-3">Weight</th>
                <th className="py-3 px-3">Reserve Price</th>
                <th className="py-3 px-3">Highest Bid</th>
                <th className="py-3 px-3">Clearance Status</th>
                <th className="py-3 px-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {lots.map((lot) => (
                <tr key={lot.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-teal-400">{lot.lot_number}</td>
                  <td className="py-3 px-3 text-white font-semibold">{lot.asset_type}</td>
                  <td className="py-3 px-3 font-mono text-slate-300">{lot.weighed_weight_kg} kg</td>
                  <td className="py-3 px-3 font-mono text-slate-300">
                    &#8358;{lot.reserve_price_ngn.toLocaleString()}
                  </td>
                  <td className="py-3 px-3 font-mono text-emerald-400 font-bold">
                    &#8358;{lot.highest_bid_ngn.toLocaleString()}
                  </td>
                  <td className="py-3 px-3">
                    {lot.audit_clearance === "CLEARED_FOR_AUCTION" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-teal-950 text-teal-400 border border-teal-800">
                        CLEARED FOR AUCTION
                      </span>
                    )}
                    {lot.audit_clearance === "AUCTION_COMPLETED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                        PROCEEDS RECONCILED
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-4 text-right">
                    {lot.audit_clearance === "CLEARED_FOR_AUCTION" && (
                      <button
                        onClick={() => handleApproveAuction(lot.id)}
                        className="px-2.5 py-1 rounded bg-teal-950 hover:bg-teal-900 text-teal-300 border border-teal-800 text-[11px] font-semibold transition"
                      >
                        Reconcile Sale
                      </button>
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

