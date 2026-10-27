/**
 * Substation Hourly Dispatch Log Sheet for Central Dispatch.
 * Telemetry intake for Frequency, 33kV Bus Voltage, Bulk MW Intake, and Shift Remarks.
 */

import React, { useState } from "react";
import { MOCK_HOURLY_DISPATCH } from "./mockData";
import { HourlyDispatchEntry } from "./types";

export function EnergyDispatchLogView() {
  const [logs, setLogs] = useState<HourlyDispatchEntry[]>(MOCK_HOURLY_DISPATCH);
  const [frequency, setFrequency] = useState("50.02");
  const [voltage, setVoltage] = useState("33.2");
  const [hydroMw, setHydroMw] = useState("48.6");
  const [tcnMw, setTcnMw] = useState("15.0");
  const [remarks, setRemarks] = useState("");

  const handleAddLog = (e: React.FormEvent) => {
    e.preventDefault();
    const fVal = parseFloat(frequency) || 50.0;
    const vVal = parseFloat(voltage) || 33.0;
    const hVal = parseFloat(hydroMw) || 0;
    const tVal = parseFloat(tcnMw) || 0;
    const now = new Date();
    const hourSlot = `${now.getHours().toString().padStart(2, "0")}:00 WAT`;

    const newEntry: HourlyDispatchEntry = {
      id: `DISP-${Date.now()}`,
      hour_slot_wat: hourSlot,
      system_frequency_hz: fVal,
      grid_voltage_33kv: vVal,
      hydro_generation_mw: hVal,
      tcn_intake_mw: tVal,
      total_bulk_mw: Number((hVal + tVal).toFixed(1)),
      load_shed_mw: 6.0,
      dispatcher_name: "Engr. P. Bulus",
      dispatcher_staff_id: "STAFF-DSP-01",
      grid_status: "NORMAL_STABLE",
      shift_remarks: remarks || "Nominal stability maintained. Switching logs normal.",
    };

    setLogs([newEntry, ...logs]);
    setRemarks("");
    alert("Hourly dispatch telemetry reading logged successfully!");
  };

  return (
    <div className="space-y-6">
      {/* Hourly Intake Form */}
      <div className="p-5 rounded-xl bg-slate-900/80 border border-slate-800">
        <h3 className="text-sm font-bold text-white flex items-center space-x-2 mb-3">
          <span>📝</span>
          <span>Log Hourly Transmission &amp; Dispatch Telemetry</span>
        </h3>

        <form onSubmit={handleAddLog} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
          <div>
            <label className="text-slate-400 font-semibold block mb-1">System Frequency (Hz)</label>
            <input
              type="number"
              step="0.01"
              value={frequency}
              onChange={(e) => setFrequency(e.target.value)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold focus:border-cyan-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">33kV Bus Voltage (kV)</label>
            <input
              type="number"
              step="0.1"
              value={voltage}
              onChange={(e) => setVoltage(e.target.value)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold focus:border-cyan-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Hydro Generation (MW)</label>
            <input
              type="number"
              step="0.1"
              value={hydroMw}
              onChange={(e) => setHydroMw(e.target.value)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold text-cyan-400 focus:border-cyan-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">TCN Intake (MW)</label>
            <input
              type="number"
              step="0.1"
              value={tcnMw}
              onChange={(e) => setTcnMw(e.target.value)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold text-purple-400 focus:border-cyan-400 focus:outline-none"
              required
            />
          </div>

          <div className="sm:col-span-2 lg:col-span-3">
            <label className="text-slate-400 font-semibold block mb-1">Dispatcher Operational Remarks</label>
            <input
              type="text"
              value={remarks}
              onChange={(e) => setRemarks(e.target.value)}
              placeholder="e.g. Load transferred from Vom to Bukuru. No feeder trips reported."
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 focus:border-cyan-400 focus:outline-none text-xs"
            />
          </div>

          <div className="flex items-end">
            <button
              type="submit"
              className="w-full py-2 px-4 rounded-lg bg-cyan-600 hover:bg-cyan-500 text-white font-bold transition shadow-lg shadow-cyan-950 flex items-center justify-center space-x-1"
            >
              <span>💾</span>
              <span>Record Reading</span>
            </button>
          </div>
        </form>
      </div>

      {/* Hourly Log Register */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
          <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
            Substation Dispatch Shift Telemetry Log
          </span>
          <span className="text-xs text-slate-400 font-mono">24-Hour Rolling Register</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Time Slot</th>
                <th className="py-3 px-3">Frequency</th>
                <th className="py-3 px-3">33kV Voltage</th>
                <th className="py-3 px-3">Hydro / TCN</th>
                <th className="py-3 px-3">Total Bulk</th>
                <th className="py-3 px-4">Dispatcher Remarks</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {logs.map((log) => (
                <tr key={log.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-cyan-400">{log.hour_slot_wat}</td>
                  <td className="py-3 px-3 font-mono text-emerald-400 font-bold">{log.system_frequency_hz} Hz</td>
                  <td className="py-3 px-3 font-mono text-white">{log.grid_voltage_33kv} kV</td>
                  <td className="py-3 px-3 font-mono text-slate-300">
                    {log.hydro_generation_mw} / {log.tcn_intake_mw} MW
                  </td>
                  <td className="py-3 px-3 font-mono font-black text-white">{log.total_bulk_mw} MW</td>
                  <td className="py-3 px-4 text-slate-400 text-xs">
                    <div>{log.shift_remarks}</div>
                    <div className="text-[10px] text-slate-500 font-mono mt-0.5">
                      Logged by {log.dispatcher_name} ({log.dispatcher_staff_id})
                    </div>
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

