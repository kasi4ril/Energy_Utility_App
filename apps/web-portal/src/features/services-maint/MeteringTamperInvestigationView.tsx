/**
 * Bypass & Tamper Investigation Desk for Metering Services Unit.
 * Tracks physical seal breaches, magnetic tamper, and energy bypass investigations.
 */

import React, { useState } from "react";

interface TamperCase {
  id: string;
  case_number: string;
  meter_number: string;
  customer_name: string;
  tamper_type: "PHYSICAL_SEAL_BROKEN" | "MAGNETIC_SENSOR_TRIGGERED" | "NEUTRAL_CURRENT_BYPASS" | "TERMINAL_COVER_OPEN";
  estimated_stolen_kwh: number;
  penalty_amount_ngn: number;
  investigation_status: "INSPECTION_PENDING" | "TAMPER_CONFIRMED" | "RECTIFIED_RESEALED";
  field_inspector: string;
}

const INITIAL_CASES: TamperCase[] = [
  {
    id: "TMP-01",
    case_number: "TMP-2026-031",
    meter_number: "0419-3321-8842",
    customer_name: "Anglo-Jos Commercial Compound",
    tamper_type: "NEUTRAL_CURRENT_BYPASS",
    estimated_stolen_kwh: 1450,
    penalty_amount_ngn: 326250,
    investigation_status: "TAMPER_CONFIRMED",
    field_inspector: "Tech. P. Pam (Metering Crew)",
  },
  {
    id: "TMP-02",
    case_number: "TMP-2026-030",
    meter_number: "5541-0021-9932",
    customer_name: "Bukuru Retail Stores",
    tamper_type: "PHYSICAL_SEAL_BROKEN",
    estimated_stolen_kwh: 450,
    penalty_amount_ngn: 101250,
    investigation_status: "INSPECTION_PENDING",
    field_inspector: "Tech. P. Pam",
  },
];

export function MeteringTamperInvestigationView() {
  const [cases, setCases] = useState<TamperCase[]>(INITIAL_CASES);

  const handleReseal = (id: string) => {
    setCases((prev) =>
      prev.map((c) =>
        c.id === id ? { ...c, investigation_status: "RECTIFIED_RESEALED" } : c
      )
    );
    alert("New holographic tamper seal applied. Commercial Billing notified for penalty surcharge.");
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>🔍</span>
            <span>Meter Bypass &amp; Tamper Investigation Desk</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Physical seal audits, optical cover-open sensors, and reverse current bypass remediation.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-rose-950 text-rose-400 border border-rose-800">
          Revenue Defense Unit
        </span>
      </div>

      {/* Cases Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Case &bull; Meter</th>
                <th className="py-3 px-3">Customer</th>
                <th className="py-3 px-3">Tamper Mode</th>
                <th className="py-3 px-3">Est. Stolen Energy</th>
                <th className="py-3 px-3">NERC Penalty</th>
                <th className="py-3 px-4 text-right">Status / Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {cases.map((c) => (
                <tr key={c.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-rose-400">
                    {c.case_number}
                    <span className="block text-[10px] text-slate-400 font-normal">{c.meter_number}</span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{c.customer_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">Inspector: {c.field_inspector}</div>
                  </td>
                  <td className="py-3 px-3 font-semibold text-amber-300">
                    {c.tamper_type.replace(/_/g, " ")}
                  </td>
                  <td className="py-3 px-3 font-mono font-bold text-white">{c.estimated_stolen_kwh} kWh</td>
                  <td className="py-3 px-3 font-mono text-emerald-400 font-bold">
                    &#8358;{c.penalty_amount_ngn.toLocaleString()}
                  </td>
                  <td className="py-3 px-4 text-right">
                    {c.investigation_status === "TAMPER_CONFIRMED" && (
                      <button
                        onClick={() => handleReseal(c.id)}
                        className="px-2.5 py-1 rounded bg-teal-950 hover:bg-teal-900 text-teal-300 border border-teal-800 text-[11px] font-semibold transition"
                      >
                        Reseal &amp; Fine
                      </button>
                    )}
                    {c.investigation_status === "INSPECTION_PENDING" && (
                      <span className="text-amber-400 font-mono text-[11px]">Inspection Out</span>
                    )}
                    {c.investigation_status === "RECTIFIED_RESEALED" && (
                      <span className="text-emerald-400 font-mono text-[11px]">Resealed</span>
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

