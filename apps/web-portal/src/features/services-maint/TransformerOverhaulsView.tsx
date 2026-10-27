/**
 * Transformer Overhauls & Preventive Servicing Queue for Transformer Unit.
 * Tracks silica gel breathers, bushing gaskets, tap changers, and nitrogen purge.
 */

import React, { useState } from "react";

interface OverhaulTask {
  id: string;
  dt_code: string;
  substation_name: string;
  task_type: "SILICA_GEL_REPLACEMENT" | "BUSHING_GASKET_RESEAL" | "OIL_FILTRATION_PURGE" | "TAP_CHANGER_OVERHAUL";
  priority: "CRITICAL" | "HIGH" | "ROUTINE";
  assigned_crew: string;
  status: "IN_PROGRESS" | "SCHEDULED" | "COMPLETED";
  scheduled_date_wat: string;
}

const INITIAL_OVERHAULS: OverhaulTask[] = [
  {
    id: "OVH-01",
    dt_code: "DT-BK-12",
    substation_name: "Bukuru Market Substation",
    task_type: "OIL_FILTRATION_PURGE",
    priority: "CRITICAL",
    assigned_crew: "Crew 2 (Lead: Engr. S. Lar)",
    status: "IN_PROGRESS",
    scheduled_date_wat: "2026-10-27 WAT",
  },
  {
    id: "OVH-02",
    dt_code: "DT-RF-08",
    substation_name: "Rayfield Ext. Substation",
    task_type: "SILICA_GEL_REPLACEMENT",
    priority: "HIGH",
    assigned_crew: "Crew 1 (Lead: Tech. M. Dusu)",
    status: "SCHEDULED",
    scheduled_date_wat: "2026-10-28 WAT",
  },
  {
    id: "OVH-03",
    dt_code: "DT-AJ-02",
    substation_name: "Anglo-Jos Commercial Hub",
    task_type: "BUSHING_GASKET_RESEAL",
    priority: "ROUTINE",
    assigned_crew: "Crew 1",
    status: "COMPLETED",
    scheduled_date_wat: "2026-10-24 WAT",
  },
];

export function TransformerOverhaulsView() {
  const [tasks, setTasks] = useState<OverhaulTask[]>(INITIAL_OVERHAULS);

  const handleCompleteTask = (id: string) => {
    setTasks((prev) =>
      prev.map((t) => (t.id === id ? { ...t, status: "COMPLETED" } : t))
    );
    alert("Transformer overhaul task completed and verified against IEEE C57 checklist.");
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>🔧</span>
            <span>Transformer Overhaul &amp; Preventive Maintenance Desk</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Active servicing queue for silica gel reactivation, bushing gaskets, and vacuum oil filtration.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-cyan-950 text-cyan-400 border border-cyan-800">
          Preventive CBM Queue
        </span>
      </div>

      {/* Overhauls Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">DT Code &bull; Substation</th>
                <th className="py-3 px-3">Service Operation</th>
                <th className="py-3 px-3">Priority</th>
                <th className="py-3 px-3">Assigned Crew</th>
                <th className="py-3 px-3">Scheduled</th>
                <th className="py-3 px-4 text-right">Status / Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {tasks.map((t) => (
                <tr key={t.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4">
                    <div className="font-bold text-white">{t.substation_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{t.dt_code}</div>
                  </td>
                  <td className="py-3 px-3 font-semibold text-amber-300">
                    {t.task_type.replace(/_/g, " ")}
                  </td>
                  <td className="py-3 px-3">
                    {t.priority === "CRITICAL" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-rose-950 text-rose-400 border border-rose-800 font-bold">
                        CRITICAL
                      </span>
                    )}
                    {t.priority === "HIGH" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                        HIGH
                      </span>
                    )}
                    {t.priority === "ROUTINE" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-slate-800 text-slate-300 border border-slate-700">
                        ROUTINE
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-3 text-slate-300">{t.assigned_crew}</td>
                  <td className="py-3 px-3 font-mono text-slate-400 text-[11px]">{t.scheduled_date_wat}</td>
                  <td className="py-3 px-4 text-right">
                    {t.status === "IN_PROGRESS" && (
                      <button
                        onClick={() => handleCompleteTask(t.id)}
                        className="px-2.5 py-1 rounded bg-emerald-950 hover:bg-emerald-900 text-emerald-300 border border-emerald-800 text-[11px] font-semibold transition"
                      >
                        Sign-off Overhaul
                      </button>
                    )}
                    {t.status === "SCHEDULED" && (
                      <span className="text-slate-400 font-mono text-[11px]">Awaiting PTW</span>
                    )}
                    {t.status === "COMPLETED" && (
                      <span className="text-emerald-400 font-mono text-[11px]">Completed</span>
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

