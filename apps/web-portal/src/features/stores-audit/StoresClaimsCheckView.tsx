/**
 * Staff Allowance Claims Initial 3-Way Check for Stores & Logistics.
 * Cross-references Duty Roster, Safety PTW isolation, and field GPS before forwarding to Internal Audit.
 */

import React, { useState } from "react";
import { MOCK_STAFF_CLAIMS } from "./storesMockData";
import { StaffAllowanceClaimItem } from "./storesTypes";

export function StoresClaimsCheckView() {
  const [claims, setClaims] = useState<StaffAllowanceClaimItem[]>(MOCK_STAFF_CLAIMS);

  const handleForwardToAudit = (id: string) => {
    alert("3-Way verified claim forwarded to Head Internal Audit for final payment approval.");
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>📍</span>
            <span>Overtime &amp; Field Allowance 3-Way Verification Desk</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Automated triangulated check: Duty Roster &bull; Substation Safety PTW &bull; Hardware GPS Distance.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-purple-950 text-purple-400 border border-purple-800">
          3-Way Triangulation Gate
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
                <th className="py-3 px-3">Roster &bull; PTW Check</th>
                <th className="py-3 px-3">GPS Proximity Check</th>
                <th className="py-3 px-3">Amount</th>
                <th className="py-3 px-4 text-right">Status / Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {claims.map((c) => (
                <tr key={c.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-purple-400">
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
                    <div className="flex items-center space-x-1.5 text-[11px]">
                      <span className={c.duty_roster_verified ? "text-emerald-400 font-bold" : "text-rose-400"}>
                        {c.duty_roster_verified ? "✓ Roster" : "✗ Roster"}
                      </span>
                      <span className="text-slate-600">&bull;</span>
                      <span className={c.safety_ptw_verified ? "text-emerald-400 font-bold" : "text-rose-400"}>
                        {c.safety_ptw_verified ? "✓ PTW" : "✗ PTW"}
                      </span>
                    </div>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-mono text-[11px]">
                      {c.gps_within_tolerance ? (
                        <span className="text-emerald-400 font-bold">
                          {c.field_gps_distance_m}m (Matched)
                        </span>
                      ) : (
                        <span className="text-rose-400 font-bold">
                          {c.field_gps_distance_m}m (Out of Bounds!)
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="py-3 px-3 font-mono font-bold text-white">
                    &#8358;{c.amount_claimed_ngn.toLocaleString()}
                  </td>
                  <td className="py-3 px-4 text-right">
                    {c.stores_verification_status === "VERIFIED_PASSED" && (
                      <button
                        onClick={() => handleForwardToAudit(c.id)}
                        className="px-2.5 py-1 rounded bg-purple-950 hover:bg-purple-900 text-purple-300 border border-purple-800 text-[11px] font-semibold transition"
                      >
                        Forward to Audit
                      </button>
                    )}
                    {c.stores_verification_status === "FLAGGED_DISCREPANCY" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-rose-950 text-rose-400 border border-rose-800 font-bold">
                        FLAGGED FRAUD RISK
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

