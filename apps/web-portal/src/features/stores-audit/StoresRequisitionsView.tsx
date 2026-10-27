/**
 * Dual-Path Requisitions Desk for Stores & Inventory.
 * Manages Path A (In-Stock warehouse allocation) and Path B (MD Procurement escalation).
 */

import React, { useState } from "react";
import { MOCK_REQUISITIONS } from "./storesMockData";
import { MaterialRequisitionRecord } from "./storesTypes";

export function StoresRequisitionsView() {
  const [requisitions, setRequisitions] = useState<MaterialRequisitionRecord[]>(MOCK_REQUISITIONS);

  const handleAllocatePathA = (id: string) => {
    setRequisitions((prev) =>
      prev.map((r) =>
        r.id === id ? { ...r, status: "ALLOCATED_WAYBILLED" } : r
      )
    );
    alert("Path A In-Stock items allocated. Internal Material Store Issue Waybill printed.");
  };

  const handleEscalatePathB = (id: string) => {
    setRequisitions((prev) =>
      prev.map((r) =>
        r.id === id ? { ...r, status: "AWAITING_MD_APPROVAL" } : r
      )
    );
    alert("Path B Requisition escalated to Managing Director's Executive Approvals Desk.");
  };

  return (
    <div className="space-y-6">
      {/* Overview Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>📑</span>
            <span>Dual-Path Material Requisitions Triage Desk</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Path A (In-Stock): Immediate Store Binning &bull; Path B (Out-of-Stock): MD Capital Authorization.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-orange-950 text-orange-400 border border-orange-800">
          Dual-Path Governance
        </span>
      </div>

      {/* Requisitions Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Requisition No</th>
                <th className="py-3 px-3">Item Requested</th>
                <th className="py-3 px-3">Dept &bull; Purpose</th>
                <th className="py-3 px-3">Path Type</th>
                <th className="py-3 px-3">Status</th>
                <th className="py-3 px-4 text-right">Store Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {requisitions.map((r) => (
                <tr key={r.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-orange-400">
                    {r.requisition_no}
                    <span className="block text-[10px] text-slate-500 font-normal">{r.created_at_wat}</span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{r.item_name}</div>
                    <div className="text-[10px] text-slate-400 font-mono">
                      Qty: <strong className="text-amber-400">{r.quantity_requested}</strong> &bull; {r.item_sku}
                    </div>
                  </td>
                  <td className="py-3 px-3 max-w-xs">
                    <div className="font-semibold text-slate-200">{r.department}</div>
                    <div className="text-[10px] text-slate-400 truncate">{r.purpose}</div>
                  </td>
                  <td className="py-3 px-3">
                    {r.path_type === "PATH_A_IN_STOCK" ? (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-teal-950 text-teal-400 border border-teal-800">
                        PATH A (IN-STOCK)
                      </span>
                    ) : (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800 font-bold">
                        PATH B (MD PROCUREMENT)
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-3">
                    {r.status === "ALLOCATED_WAYBILLED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                        WAYBILLED
                      </span>
                    )}
                    {r.status === "AWAITING_MD_APPROVAL" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                        AWAITING MD
                      </span>
                    )}
                    {r.status === "PENDING_STOCK_CHECK" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-slate-800 text-slate-300 border border-slate-700">
                        PENDING CHECK
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-4 text-right">
                    {r.path_type === "PATH_A_IN_STOCK" && r.status !== "ALLOCATED_WAYBILLED" && (
                      <button
                        onClick={() => handleAllocatePathA(r.id)}
                        className="px-2.5 py-1 rounded bg-teal-950 hover:bg-teal-900 text-teal-300 border border-teal-800 text-[11px] font-semibold transition"
                      >
                        Allocate Stock
                      </button>
                    )}
                    {r.path_type === "PATH_B_MD_PROCUREMENT" && r.status !== "AWAITING_MD_APPROVAL" && (
                      <button
                        onClick={() => handleEscalatePathB(r.id)}
                        className="px-2.5 py-1 rounded bg-amber-950 hover:bg-amber-900 text-amber-300 border border-amber-800 text-[11px] font-semibold transition"
                      >
                        Escalate to MD
                      </button>
                    )}
                    {r.status === "ALLOCATED_WAYBILLED" && (
                      <span className="text-[11px] text-emerald-400 font-mono">Dispatched</span>
                    )}
                    {r.status === "AWAITING_MD_APPROVAL" && (
                      <span className="text-[11px] text-amber-400 font-mono">With MD</span>
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

