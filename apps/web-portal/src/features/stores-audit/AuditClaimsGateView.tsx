/**
 * 3-Way Claims Final Approval Gate for Internal Audit.
 * Authoritative financial sign-off for overtime and travel allowances.
 */

import React, { useState } from "react";
import { MOCK_FINAL_CLAIMS } from "./auditMockData";
import { FinalClaimAuditItem } from "./auditTypes";

export function AuditClaimsGateView() {
  const [claims, setClaims] = useState<FinalClaimAuditItem[]>(MOCK_FINAL_CLAIMS);

  const handleAuditAction = (id: string, action: "AUDIT_APPROVED" | "FLAGGED_FOR_HEARING") => {
    setClaims((prev) =>
      prev.map((c) =>
        c.id === id ? { ...c, audit_decision: action } : c
      )
    );
    alert(
      action === "AUDIT_APPROVED"
        ? "Claim approved for payroll disbursement."
        : "Claim flagged and referred to Internal Disciplinary Committee for distance falsification."
    );
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>📍</span>
            <span>Internal Audit 3-Way Claims Final Authorization Desk</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Sole authority for releasing payroll funds against field overtime and emergency allowances.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-purple-950 text-purple-400 border border-purple-800">
          Segregation of Duties
        </span>
      </div>

      {/* Claims Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Claim No</th>
                <th className="py-3 px-3">Staff &bull; Role</th>
                <th className="py-3 px-3">Claim Type</th>
                <th className="py-3 px-3">Triangulation Status</th>
                <th className="py-3 px-3">Amount Claimed</th>
                <th className="py-3 px-4 text-right">Audit Authority</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {claims.map((c) => (
                <tr key={c.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-rose-400">
                    {c.claim_number}
                    <span className="block text-[10px] text-slate-500 font-normal">{c.created_at_wat}</span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{c.staff_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{c.staff_id}</div>
                  </td>
                  <td className="py-3 px-3 font-semibold text-slate-300">
                    {c.claim_type.replace(/_/g, " ")}
                  </td>
                  <td className="py-3 px-3">
                    {c.triangulation_status === "PASSED_ROSTER_PTW_GPS" ? (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800 font-bold">
                        ✓ 3-WAY MATCHED
                      </span>
                    ) : (
                      <div>
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-rose-950 text-rose-400 border border-rose-800 font-bold">
                          ✗ GPS DISTANCE MISMATCH
                        </span>
                        <div className="text-[10px] text-slate-400 mt-1 max-w-xs">{c.audit_remarks}</div>
                      </div>
                    )}
                  </td>
                  <td className="py-3 px-3 font-mono font-bold text-white text-sm">
                    &#8358;{c.amount_claimed_ngn.toLocaleString()}
                  </td>
                  <td className="py-3 px-4 text-right">
                    {c.audit_decision === "PENDING_AUDITOR" ? (
                      <div className="flex items-center justify-end space-x-2">
                        <button
                          onClick={() => handleAuditAction(c.id, "AUDIT_APPROVED")}
                          className="px-2.5 py-1 rounded bg-emerald-950 hover:bg-emerald-900 text-emerald-300 border border-emerald-800 text-[11px] font-semibold transition"
                        >
                          Approve Payment
                        </button>
                        <button
                          onClick={() => handleAuditAction(c.id, "FLAGGED_FOR_HEARING")}
                          className="px-2.5 py-1 rounded bg-rose-950 hover:bg-rose-900 text-rose-300 border border-rose-800 text-[11px] font-semibold transition"
                        >
                          Flag Discrepancy
                        </button>
                      </div>
                    ) : c.audit_decision === "AUDIT_APPROVED" ? (
                      <span className="text-[11px] text-emerald-400 font-mono font-bold">Approved</span>
                    ) : (
                      <span className="text-[11px] text-rose-400 font-mono font-bold">Disciplinary Hearing</span>
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

