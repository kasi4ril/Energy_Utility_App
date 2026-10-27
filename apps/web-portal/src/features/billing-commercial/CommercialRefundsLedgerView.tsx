/**
 * Meter-Purchase Refund Amortization Ledger for Commercial Operations.
 * NERC MAP/NMMP regulatory compliance ledger tracking energy rebates to customers.
 */

import React, { useState } from "react";
import { MOCK_REFUND_LEDGER } from "./mockData";
import { MeterRefundLedgerEntry } from "./types";

export function CommercialRefundsLedgerView() {
  const [ledgers, setLedgers] = useState<MeterRefundLedgerEntry[]>(MOCK_REFUND_LEDGER);

  const handleApplyRebate = (id: string) => {
    setLedgers((prev) =>
      prev.map((entry) => {
        if (entry.id !== id) return entry;
        const rebate = Math.min(2500, entry.remaining_balance_ngn);
        const newBal = entry.remaining_balance_ngn - rebate;
        return {
          ...entry,
          total_refunded_to_date_ngn: entry.total_refunded_to_date_ngn + rebate,
          remaining_balance_ngn: newBal,
          status: newBal <= 0 ? "FULLY_REFUNDED" : "ACTIVE_AMORTIZING",
          last_rebate_wat: "Today",
        };
      })
    );
    alert("Rebate applied to customer's account successfully!");
  };

  return (
    <div className="space-y-6">
      {/* Overview Card */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>📑</span>
            <span>NERC Meter-Purchase Refund Amortization Ledger</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Automatic energy credit amortization for customer-financed meter installations.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-emerald-950 text-emerald-400 border border-emerald-800">
          NERC MAP/NMMP Regulated
        </span>
      </div>

      {/* Refund Ledger Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Customer &bull; Meter</th>
                <th className="py-3 px-3">Original Cost</th>
                <th className="py-3 px-3">Refunded to Date</th>
                <th className="py-3 px-3">Outstanding Balance</th>
                <th className="py-3 px-3">Amortization Progress</th>
                <th className="py-3 px-3">Status</th>
                <th className="py-3 px-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {ledgers.map((l) => {
                const pct = Math.min(
                  100,
                  Math.round((l.total_refunded_to_date_ngn / l.original_meter_cost_ngn) * 100)
                );
                return (
                  <tr key={l.id} className="hover:bg-slate-800/40 transition">
                    <td className="py-3 px-4">
                      <div className="font-bold text-white">{l.customer_name}</div>
                      <div className="text-[10px] text-slate-500 font-mono">{l.meter_number}</div>
                    </td>
                    <td className="py-3 px-3 font-mono font-bold text-white">
                      &#8358;{l.original_meter_cost_ngn.toLocaleString()}
                    </td>
                    <td className="py-3 px-3 font-mono text-emerald-400 font-bold">
                      &#8358;{l.total_refunded_to_date_ngn.toLocaleString()}
                    </td>
                    <td className="py-3 px-3 font-mono text-teal-300">
                      &#8358;{l.remaining_balance_ngn.toLocaleString()}
                    </td>
                    <td className="py-3 px-3">
                      <div className="w-28 bg-slate-800 h-2 rounded-full overflow-hidden">
                        <div
                          className="bg-teal-400 h-full rounded-full transition-all"
                          style={{ width: `${pct}%` }}
                        />
                      </div>
                      <span className="text-[10px] text-slate-400 font-mono mt-0.5 block">{pct}% Complete</span>
                    </td>
                    <td className="py-3 px-3">
                      {l.status === "ACTIVE_AMORTIZING" ? (
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-teal-950 text-teal-400 border border-teal-800">
                          AMORTIZING
                        </span>
                      ) : (
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                          SETTLED
                        </span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-right">
                      {l.status === "ACTIVE_AMORTIZING" && (
                        <button
                          onClick={() => handleApplyRebate(l.id)}
                          className="px-2.5 py-1 rounded bg-teal-950 hover:bg-teal-900 text-teal-300 border border-teal-800 text-[11px] font-semibold transition"
                        >
                          Credit Vend Rebate
                        </button>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

