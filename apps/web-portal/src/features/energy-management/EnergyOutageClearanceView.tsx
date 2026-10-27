/**
 * Outage Coordination & Feeder Clearance Register for Central Dispatch.
 * Inter-departmental coordination terminal connecting Central Dispatch with Feeder Maintenance.
 */

import React, { useState } from "react";

interface FeederClearance {
  id: string;
  outage_code: string;
  feeder_name: string;
  voltage_level: string;
  crew_leader: string;
  requested_window: string;
  dispatch_status: "PENDING_DISPATCH_GRANT" | "GRANTED_ISOLATED" | "SURRENDER_RECEIVED" | "CLOSED_ENERGIZED";
  load_relief_mw: number;
  granting_time_wat?: string;
  closure_state?: "CLOSED_ENERGIZED_IMMEDIATELY" | "LEFT_ON_LOAD_SHEDDING" | "TRIPPED_ON_CLOSURE";
}

const INITIAL_CLEARANCES: FeederClearance[] = [
  {
    id: "CLR-01",
    outage_code: "OUT-001-33KV",
    feeder_name: "Anglo-Jos 33kV Industrial Corridor",
    voltage_level: "33kV",
    crew_leader: "Engr. D. Gyang (Feeder Unit Head)",
    requested_window: "2026-10-27 14:00 - 18:00 WAT",
    dispatch_status: "GRANTED_ISOLATED",
    load_relief_mw: 8.4,
    granting_time_wat: "14:10 WAT",
  },
  {
    id: "CLR-02",
    outage_code: "OUT-002-33KV",
    feeder_name: "Bukuru Express 33kV Feeder",
    voltage_level: "33kV",
    crew_leader: "Engr. S. Lar",
    requested_window: "2026-10-28 09:00 - 13:00 WAT",
    dispatch_status: "PENDING_DISPATCH_GRANT",
    load_relief_mw: 6.2,
  },
  {
    id: "CLR-03",
    outage_code: "OUT-003-11KV",
    feeder_name: "Vom Vet-Research 11kV Feeder",
    voltage_level: "11kV",
    crew_leader: "Engr. D. Gyang",
    requested_window: "2026-10-26 10:00 - 14:00 WAT",
    dispatch_status: "CLOSED_ENERGIZED",
    load_relief_mw: 3.1,
    granting_time_wat: "10:05 WAT",
    closure_state: "CLOSED_ENERGIZED_IMMEDIATELY",
  },
];

export function EnergyOutageClearanceView() {
  const [clearances, setClearances] = useState<FeederClearance[]>(INITIAL_CLEARANCES);

  const handleGrantOutage = (id: string) => {
    const now = new Date();
    const timeStr = `${now.getHours().toString().padStart(2, "0")}:${now.getMinutes().toString().padStart(2, "0")} WAT`;
    setClearances((prev) =>
      prev.map((c) =>
        c.id === id
          ? {
              ...c,
              dispatch_status: "GRANTED_ISOLATED",
              granting_time_wat: timeStr,
            }
          : c
      )
    );
    alert(`Outage clearance GRANTED & line tagged open at ${timeStr}. Feeder Unit notified.`);
  };

  const handleCloseFeeder = (id: string, outcome: "CLOSED_ENERGIZED_IMMEDIATELY" | "LEFT_ON_LOAD_SHEDDING" | "TRIPPED_ON_CLOSURE") => {
    setClearances((prev) =>
      prev.map((c) =>
        c.id === id
          ? {
              ...c,
              dispatch_status: "CLOSED_ENERGIZED",
              closure_state: outcome,
            }
          : c
      )
    );
    alert(`Feeder closed with status: ${outcome}. Outage lifecycle completed.`);
  };

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>📋</span>
            <span>Feeder Outage Requisition &amp; Clearance Dispatch Desk</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Real-time handoff with Services &amp; Maintenance Feeder Unit for 33kV/11kV line clearance.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-cyan-950 text-cyan-400 border border-cyan-800">
          Grid Handover Portal
        </span>
      </div>

      {/* Clearances Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Outage Code</th>
                <th className="py-3 px-3">Corridor &bull; Voltage</th>
                <th className="py-3 px-3">Crew Lead</th>
                <th className="py-3 px-3">Window</th>
                <th className="py-3 px-3">Dispatch Status</th>
                <th className="py-3 px-4 text-right">Dispatch Authority</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {clearances.map((c) => (
                <tr key={c.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-cyan-400">
                    {c.outage_code}
                    <span className="block text-[10px] text-slate-500 font-normal">
                      Load Drop: {c.load_relief_mw} MW
                    </span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{c.feeder_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{c.voltage_level}</div>
                  </td>
                  <td className="py-3 px-3 text-slate-300">{c.crew_leader}</td>
                  <td className="py-3 px-3 font-mono text-slate-300 text-[11px]">{c.requested_window}</td>
                  <td className="py-3 px-3">
                    {c.dispatch_status === "PENDING_DISPATCH_GRANT" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                        PENDING DISPATCH
                      </span>
                    )}
                    {c.dispatch_status === "GRANTED_ISOLATED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-sky-950 text-sky-400 border border-sky-800">
                        ISOLATED &bull; {c.granting_time_wat}
                      </span>
                    )}
                    {c.dispatch_status === "CLOSED_ENERGIZED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                        CLOSED &amp; ENERGIZED
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-4 text-right">
                    {c.dispatch_status === "PENDING_DISPATCH_GRANT" && (
                      <button
                        onClick={() => handleGrantOutage(c.id)}
                        className="px-2.5 py-1 rounded bg-cyan-950 hover:bg-cyan-900 text-cyan-300 border border-cyan-800 text-[11px] font-semibold transition"
                      >
                        Grant Clearance
                      </button>
                    )}
                    {c.dispatch_status === "GRANTED_ISOLATED" && (
                      <div className="flex items-center justify-end space-x-1.5">
                        <button
                          onClick={() => handleCloseFeeder(c.id, "CLOSED_ENERGIZED_IMMEDIATELY")}
                          className="px-2 py-1 rounded bg-emerald-950 hover:bg-emerald-900 text-emerald-300 border border-emerald-800 text-[10px] font-semibold transition"
                          title="Close & Energize Immediately"
                        >
                          Close Line
                        </button>
                        <button
                          onClick={() => handleCloseFeeder(c.id, "LEFT_ON_LOAD_SHEDDING")}
                          className="px-2 py-1 rounded bg-amber-950 hover:bg-amber-900 text-amber-300 border border-amber-800 text-[10px] font-semibold transition"
                          title="Keep on Load Shedding"
                        >
                          Ration
                        </button>
                      </div>
                    )}
                    {c.dispatch_status === "CLOSED_ENERGIZED" && (
                      <span className="text-[10px] font-mono text-slate-500">
                        {c.closure_state || "ENERGIZED"}
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

