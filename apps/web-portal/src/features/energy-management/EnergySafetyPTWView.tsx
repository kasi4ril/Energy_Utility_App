/**
 * 5-Step Electrical Safety PTW Isolation Console for Central Dispatch.
 * Strict NERC & NEMSA electrical safety lock protecting high-voltage maintenance crews.
 */

import React, { useState } from "react";
import { MOCK_PTW_PERMITS } from "./mockData";
import { SafetyPermitRecord, PTWStep } from "./types";

export function EnergySafetyPTWView() {
  const [permits, setPermits] = useState<SafetyPermitRecord[]>(MOCK_PTW_PERMITS);
  const [selectedPermit, setSelectedPermit] = useState<SafetyPermitRecord | null>(null);

  const steps: { key: PTWStep; label: string; number: number }[] = [
    { key: "REQUESTED", label: "Requisition", number: 1 },
    { key: "SAFETY_APPROVED", label: "HSE Approved", number: 2 },
    { key: "ISOLATED", label: "Line Isolated", number: 3 },
    { key: "WORK_IN_PROGRESS", label: "Work In Progress", number: 4 },
    { key: "CLEARANCE_DECLARED", label: "Clearance Handover", number: 5 },
    { key: "RESTORED", label: "Re-Energized", number: 6 },
  ];

  const handleAdvanceStep = (id: string) => {
    setPermits((prev) =>
      prev.map((p) => {
        if (p.id !== id) return p;
        if (p.step === "WORK_IN_PROGRESS") {
          return { ...p, step: "CLEARANCE_DECLARED", updated_at_wat: "Just Now" };
        }
        if (p.step === "CLEARANCE_DECLARED") {
          return { ...p, step: "RESTORED", updated_at_wat: "Just Now", restored_by_operator_id: "STAFF-DSP-01" };
        }
        return p;
      })
    );
  };

  return (
    <div className="space-y-6">
      {/* 5-Step NERC Isolation Protocol Header Card */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
        <div className="flex items-center justify-between mb-3">
          <div>
            <h3 className="text-sm font-bold text-white flex items-center space-x-2">
              <span>🔒</span>
              <span>NERC 5-Step Substation Electrical Safety Isolation Protocol</span>
            </h3>
            <p className="text-xs text-slate-400 mt-0.5">
              Strict lockout-tagout (LOTO) &amp; earth grounding verification before circuit handback.
            </p>
          </div>
          <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-sky-950 text-sky-400 border border-sky-800">
            NEMSA Standard Compliant
          </span>
        </div>

        {/* Visual Stepper */}
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-2 pt-2">
          {steps.map((s, idx) => (
            <div
              key={s.key}
              className="p-2.5 rounded-lg bg-slate-950/60 border border-slate-800 text-center relative"
            >
              <div className="text-[10px] font-mono font-bold text-cyan-400">Step {s.number}</div>
              <div className="text-xs font-semibold text-slate-200 mt-0.5">{s.label}</div>
              {idx < steps.length - 1 && (
                <span className="hidden lg:block absolute -right-2 top-1/2 -translate-y-1/2 text-slate-600 z-10">
                  &rarr;
                </span>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Active PTW Register Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
          <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
            Active Safety PTW Permits
          </span>
          <span className="text-xs text-slate-400">Total: {permits.length} Records</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Permit No</th>
                <th className="py-3 px-3">Feeder &bull; Substation</th>
                <th className="py-3 px-3">Isolation &bull; Ground Point</th>
                <th className="py-3 px-3">Permit Step</th>
                <th className="py-3 px-3">Requested By</th>
                <th className="py-3 px-4 text-right">Dispatch Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {permits.map((p) => (
                <tr key={p.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-cyan-400">
                    {p.permit_number}
                    <span className="block text-[10px] text-slate-500 font-normal">{p.created_at_wat}</span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{p.feeder_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{p.substation_code}</div>
                  </td>
                  <td className="py-3 px-3">
                    <div className="text-slate-300 font-mono text-[11px]">{p.isolation_point}</div>
                    <div className="text-[10px] text-emerald-400 font-mono">{p.earth_switch_tag}</div>
                  </td>
                  <td className="py-3 px-3">
                    {p.step === "WORK_IN_PROGRESS" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-sky-950 text-sky-400 border border-sky-800">
                        WORK IN PROGRESS
                      </span>
                    )}
                    {p.step === "CLEARANCE_DECLARED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                        CLEARANCE DECLARED
                      </span>
                    )}
                    {p.step === "RESTORED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                        RESTORED (CLOSED)
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-3 text-slate-300">
                    <div>{p.requested_by_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{p.requested_by_staff_id}</div>
                  </td>
                  <td className="py-3 px-4 text-right">
                    {p.step === "WORK_IN_PROGRESS" && (
                      <button
                        onClick={() => handleAdvanceStep(p.id)}
                        className="px-2.5 py-1 rounded bg-amber-950 hover:bg-amber-900 text-amber-300 border border-amber-800 text-[11px] font-semibold transition"
                      >
                        Accept Surrender
                      </button>
                    )}
                    {p.step === "CLEARANCE_DECLARED" && (
                      <button
                        onClick={() => handleAdvanceStep(p.id)}
                        className="px-2.5 py-1 rounded bg-emerald-950 hover:bg-emerald-900 text-emerald-300 border border-emerald-800 text-[11px] font-semibold transition"
                      >
                        Re-Energize Line
                      </button>
                    )}
                    {p.step === "RESTORED" && (
                      <span className="text-[11px] text-emerald-400 font-mono">Restored</span>
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

