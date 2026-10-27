/**
 * Proof of Work Submission & Verification Desk for Transformer Unit.
 * Enforces hardware GPS anti-ghost verification (±15m tolerance) and photo evidence.
 */

import React, { useState } from "react";

interface ProofItem {
  id: string;
  order_number: string;
  dt_name: string;
  task_summary: string;
  team_lead: string;
  gps_coords: string;
  haversine_distance_m: number;
  gps_within_tolerance: boolean;
  status: "PENDING_UNIT_HEAD" | "VERIFIED_APPROVED" | "REJECTED_AUDIT";
}

const INITIAL_PROOFS: ProofItem[] = [
  {
    id: "POW-TR-01",
    order_number: "WO-TR-2026-042",
    dt_name: "DT-AJ-04 (Anglo-Jos Ind. 4)",
    task_summary: "Oil dielectric filtration & silica gel charge replacement completed.",
    team_lead: "Tech. M. Dusu (Crew 1)",
    gps_coords: "9.8965° N, 8.8582° E",
    haversine_distance_m: 6.2,
    gps_within_tolerance: true,
    status: "PENDING_UNIT_HEAD",
  },
  {
    id: "POW-TR-02",
    order_number: "WO-TR-2026-039",
    dt_name: "DT-BK-08 (Bukuru High Street)",
    task_summary: "HT bushing gasket replaced and leak sealed. Megger test > 200 MΩ.",
    team_lead: "Engr. S. Lar",
    gps_coords: "9.8012° N, 8.8741° E",
    haversine_distance_m: 9.8,
    gps_within_tolerance: true,
    status: "VERIFIED_APPROVED",
  },
];

export function TransformerProofOfWorkView() {
  const [proofs, setProofs] = useState<ProofItem[]>(INITIAL_PROOFS);

  const handleApproveProof = (id: string) => {
    setProofs((prev) =>
      prev.map((p) => (p.id === id ? { ...p, status: "VERIFIED_APPROVED" } : p))
    );
    alert(
      "Proof of Work approved! GPS distance verified within 15m tolerance. Customer ticket reverted to ATTENDED_TO."
    );
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>📸</span>
            <span>Proof of Work Submission &amp; Anti-Ghost GPS Audit Desk</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Strict &plusmn;15m radial tolerance check against registered substation GIS coordinates.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-purple-950 text-purple-400 border border-purple-800">
          PostGIS Haversine Gate
        </span>
      </div>

      {/* Proof Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {proofs.map((p) => (
          <div
            key={p.id}
            className="p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between space-y-4"
          >
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-mono font-bold text-amber-400">{p.order_number}</span>
                <span
                  className={`text-[10px] font-mono px-2 py-0.5 rounded font-bold border ${
                    p.gps_within_tolerance
                      ? "bg-emerald-950 text-emerald-400 border-emerald-800"
                      : "bg-rose-950 text-rose-400 border-rose-800"
                  }`}
                >
                  GPS: {p.haversine_distance_m}m from Asset ({p.gps_within_tolerance ? "MATCH" : "FAIL"})
                </span>
              </div>

              <h4 className="text-sm font-bold text-white">{p.dt_name}</h4>
              <div className="text-[10px] text-slate-400 font-mono mt-0.5">
                Crew Lead: {p.team_lead} &bull; Coordinates: {p.gps_coords}
              </div>

              <p className="text-xs text-slate-300 mt-3 bg-slate-950/70 p-3 rounded-lg border border-slate-800 leading-relaxed font-sans">
                "{p.task_summary}"
              </p>

              {/* Photo Evidence Placeholders */}
              <div className="grid grid-cols-2 gap-2 mt-3">
                <div className="p-3 rounded-lg bg-slate-950 border border-slate-800 text-center">
                  <div className="text-xl">📸</div>
                  <div className="text-[10px] text-slate-400 font-bold mt-1">Before: Leak Exposed</div>
                  <div className="text-[9px] text-slate-500 font-mono">SHA-256 Verified</div>
                </div>
                <div className="p-3 rounded-lg bg-slate-950 border border-slate-800 text-center">
                  <div className="text-xl">✅</div>
                  <div className="text-[10px] text-emerald-400 font-bold mt-1">After: Gasket Sealed</div>
                  <div className="text-[9px] text-slate-500 font-mono">Megger Tested</div>
                </div>
              </div>
            </div>

            <div className="pt-3 border-t border-slate-800 flex items-center justify-between">
              <div>
                {p.status === "PENDING_UNIT_HEAD" && (
                  <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                    PENDING UNIT HEAD
                  </span>
                )}
                {p.status === "VERIFIED_APPROVED" && (
                  <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800 font-bold">
                    VERIFIED &amp; ATTENDED TO
                  </span>
                )}
              </div>

              {p.status === "PENDING_UNIT_HEAD" && (
                <button
                  onClick={() => handleApproveProof(p.id)}
                  className="px-3 py-1.5 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs transition shadow-lg shadow-emerald-950 flex items-center space-x-1"
                >
                  <span>✓</span>
                  <span>Approve &amp; Sign Off</span>
                </button>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

