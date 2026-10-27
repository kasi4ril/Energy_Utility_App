/**
 * DGA & Dielectric Oil Breakdown Test Desk for Transformer Maintenance.
 * Logs ASTM D877 / D1816 BDV tests, moisture content, and acidity.
 */

import React, { useState } from "react";

interface OilTestRecord {
  id: string;
  substation_dt: string;
  dt_capacity_kva: number;
  bdv_breakdown_kv: number;
  moisture_ppm: number;
  acidity_mg_koh: number;
  test_date_wat: string;
  tester_staff_id: string;
  status: "OPTIMAL_PASS" | "MODERATE_WARNING" | "CRITICAL_ACTION_REQUIRED";
}

const INITIAL_OIL_TESTS: OilTestRecord[] = [
  {
    id: "OIL-01",
    substation_dt: "DT-AJ-04 (Anglo-Jos Ind. 4)",
    dt_capacity_kva: 1000,
    bdv_breakdown_kv: 48.5,
    moisture_ppm: 14,
    acidity_mg_koh: 0.04,
    test_date_wat: "2026-10-27 11:30 WAT",
    tester_staff_id: "STAFF-ENG-08",
    status: "OPTIMAL_PASS",
  },
  {
    id: "OIL-02",
    substation_dt: "DT-BK-12 (Bukuru Market DT)",
    dt_capacity_kva: 500,
    bdv_breakdown_kv: 28.2,
    moisture_ppm: 38,
    acidity_mg_koh: 0.16,
    test_date_wat: "2026-10-26 14:00 WAT",
    tester_staff_id: "STAFF-ENG-08",
    status: "CRITICAL_ACTION_REQUIRED",
  },
  {
    id: "OIL-03",
    substation_dt: "DT-RF-02 (Rayfield Roundabout)",
    dt_capacity_kva: 500,
    bdv_breakdown_kv: 36.0,
    moisture_ppm: 22,
    acidity_mg_koh: 0.08,
    test_date_wat: "2026-10-25 09:45 WAT",
    tester_staff_id: "STAFF-ENG-08",
    status: "MODERATE_WARNING",
  },
];

export function TransformerDGAView() {
  const [tests, setTests] = useState<OilTestRecord[]>(INITIAL_OIL_TESTS);
  const [dtName, setDtName] = useState("");
  const [bdv, setBdv] = useState("45.0");
  const [moisture, setMoisture] = useState("16");
  const [acidity, setAcidity] = useState("0.05");

  const handleAddTest = (e: React.FormEvent) => {
    e.preventDefault();
    const b = parseFloat(bdv) || 0;
    const m = parseInt(moisture, 10) || 0;
    const a = parseFloat(acidity) || 0;
    if (!dtName) {
      alert("Please specify the transformer substation identifier.");
      return;
    }

    const stat = b < 30 ? "CRITICAL_ACTION_REQUIRED" : b < 40 ? "MODERATE_WARNING" : "OPTIMAL_PASS";

    const newTest: OilTestRecord = {
      id: `OIL-${Date.now()}`,
      substation_dt: dtName,
      dt_capacity_kva: 500,
      bdv_breakdown_kv: b,
      moisture_ppm: m,
      acidity_mg_koh: a,
      test_date_wat: "Just Now",
      tester_staff_id: "STAFF-ENG-08",
      status: stat,
    };

    setTests([newTest, ...tests]);
    setDtName("");
    alert("Transformer oil dielectric breakdown reading recorded successfully!");
  };

  return (
    <div className="space-y-6">
      {/* Intake Form */}
      <div className="p-5 rounded-xl bg-slate-900/80 border border-slate-800">
        <h3 className="text-sm font-bold text-white flex items-center space-x-2 mb-3">
          <span>🧪</span>
          <span>Record Transformer Oil Dielectric Breakdown Test (ASTM D877 / D1816)</span>
        </h3>

        <form onSubmit={handleAddTest} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
          <div>
            <label className="text-slate-400 font-semibold block mb-1">Transformer Substation</label>
            <input
              type="text"
              value={dtName}
              onChange={(e) => setDtName(e.target.value)}
              placeholder="e.g. DT-VM-03 (Vom Hospital)"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-bold focus:border-amber-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">BDV Dielectric Breakdown (kV)</label>
            <input
              type="number"
              step="0.1"
              value={bdv}
              onChange={(e) => setBdv(e.target.value)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold text-amber-400 focus:border-amber-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Moisture Content (ppm)</label>
            <input
              type="number"
              value={moisture}
              onChange={(e) => setMoisture(e.target.value)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono focus:border-amber-400 focus:outline-none"
              required
            />
          </div>

          <div className="flex items-end">
            <button
              type="submit"
              className="w-full py-2 px-4 rounded-lg bg-amber-600 hover:bg-amber-500 text-white font-bold transition shadow-lg shadow-amber-950 flex items-center justify-center space-x-1"
            >
              <span>🔬</span>
              <span>Log Lab Analysis</span>
            </button>
          </div>
        </form>
      </div>

      {/* Test Log Register */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
          <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
            Dielectric Breakdown Voltage (BDV) &amp; Oil Quality Ledger
          </span>
          <span className="text-xs text-slate-400 font-mono">ASTM Standards</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Substation DT</th>
                <th className="py-3 px-3">BDV Strength</th>
                <th className="py-3 px-3">Moisture</th>
                <th className="py-3 px-3">Acidity</th>
                <th className="py-3 px-3">Test Date</th>
                <th className="py-3 px-4 text-right">Condition</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {tests.map((t) => (
                <tr key={t.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4">
                    <div className="font-bold text-white">{t.substation_dt}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{t.dt_capacity_kva} kVA &bull; {t.tester_staff_id}</div>
                  </td>
                  <td className="py-3 px-3 font-mono font-bold text-white text-sm">
                    {t.bdv_breakdown_kv} kV
                  </td>
                  <td className="py-3 px-3 font-mono text-slate-300">{t.moisture_ppm} ppm</td>
                  <td className="py-3 px-3 font-mono text-slate-300">{t.acidity_mg_koh} mg KOH/g</td>
                  <td className="py-3 px-3 text-slate-400 font-mono text-[11px]">{t.test_date_wat}</td>
                  <td className="py-3 px-4 text-right">
                    {t.status === "OPTIMAL_PASS" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800 font-bold">
                        PASS (&gt; 40kV)
                      </span>
                    )}
                    {t.status === "MODERATE_WARNING" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                        MONITOR (30-40kV)
                      </span>
                    )}
                    {t.status === "CRITICAL_ACTION_REQUIRED" && (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-rose-950 text-rose-400 border border-rose-800 font-bold">
                        OIL PURGE DUE (&lt; 30kV)
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

