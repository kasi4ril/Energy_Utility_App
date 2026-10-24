import React, { useState } from "react";
import { FeederPatrolRecord, FeederCorridor } from "./types";
import { MOCK_FEEDER_PATROLS } from "./mockData";

export const FeederStandardsView: React.FC = () => {
  const [patrols, setPatrols] = useState<FeederPatrolRecord[]>(MOCK_FEEDER_PATROLS);
  const [showAddModal, setShowAddModal] = useState(false);

  // Form State
  const [corridor, setCorridor] = useState<FeederCorridor>("FEEDER_33KV_ANGLO_JOS");
  const [poleSpan, setPoleSpan] = useState("Poles #121 to #190");
  const [leadPatroller, setLeadPatroller] = useState("Engr. D. Gyang");
  const [vegRisk, setVegRisk] = useState<"SAFE" | "ATTENTION_REQUIRED" | "IMMINENT_FLASHOVER">("SAFE");
  const [insulatorCond, setInsulatorCond] = useState<"ALL_INTACT" | "CRACKED_PIN_FOUND" | "FLASHED_DISC">("ALL_INTACT");
  const [hotspotDelta, setHotspotDelta] = useState("3.4");
  const [earthOhms, setEarthOhms] = useState("2.9");

  const handleAddPatrol = (e: React.FormEvent) => {
    e.preventDefault();
    const delta = parseFloat(hotspotDelta) || 0;
    const ohms = parseFloat(earthOhms) || 0;
    const isCompliant = ohms < 5.0 && delta < 10.0 && vegRisk === "SAFE";

    const newRec: FeederPatrolRecord = {
      id: `PAT-0${patrols.length + 1}`,
      patrol_code: `PAT-${corridor.split("_")[2]}-2026-W43`,
      corridor,
      pole_span: poleSpan,
      patrol_date: new Date().toISOString().split("T")[0],
      lead_patroller: leadPatroller,
      vegetation_risk_level: vegRisk,
      insulators_condition: insulatorCond,
      hotspot_delta_c: delta,
      earth_resistance_ohms: ohms,
      compliance_status: isCompliant ? "PASSED" : "ACTION_ORDERED",
    };

    setPatrols([newRec, ...patrols]);
    setShowAddModal(false);
  };

  return (
    <div className="space-y-6">
      {/* Banner */}
      <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="inline-flex items-center gap-2 px-2.5 py-1 rounded-full text-xs font-semibold bg-sky-500/10 text-sky-300 border border-sky-500/30 mb-2">
              <span className="w-2 h-2 rounded-full bg-sky-400" />
              NERC Grid Code Standards • Section 12 Engineering Mandate
            </div>
            <h2 className="text-xl font-bold text-white tracking-tight">
              Feeder Engineering Standards & Line Patrol Quality Register
            </h2>
            <p className="text-sm text-slate-400 mt-1">
              Operational parameters: Right-of-Way clearance, Thermographic IR scans, Earth Resistance testing (&lt; 5.0 Ω), and 5-Step Safety Isolation.
            </p>
          </div>
          <button
            onClick={() => setShowAddModal(true)}
            className="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-slate-950 font-bold text-xs rounded-lg shadow transition-all flex items-center gap-2 self-start md:self-auto"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            Record Line Patrol Inspection
          </button>
        </div>
      </div>

      {/* 3 Engineering Standards Reference Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
        {/* Card 1: Earthing */}
        <div className="bg-slate-900/70 border border-slate-800 rounded-xl p-4 shadow-sm space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold text-emerald-400 uppercase tracking-wider">Earth Resistance Standard</span>
            <span className="text-xs bg-emerald-950 text-emerald-300 px-2 py-0.5 rounded font-mono border border-emerald-800/40">&lt; 5.0 Ω Mandate</span>
          </div>
          <div className="text-2xl font-black text-white">IEEE Std 80 / NERC</div>
          <p className="text-xs text-slate-400 leading-relaxed">
            All 33kV switchyard surge arresters, transformer neutrals, and distribution steel structure down-leads must maintain ground resistance &lt; 5.0 Ω. Quarterly copper sulphate ground treatment required on rocky terrain.
          </p>
        </div>

        {/* Card 2: Thermography */}
        <div className="bg-slate-900/70 border border-slate-800 rounded-xl p-4 shadow-sm space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold text-amber-400 uppercase tracking-wider">IR Thermography Standards</span>
            <span className="text-xs bg-amber-950 text-amber-300 px-2 py-0.5 rounded font-mono border border-amber-800/40">FLIR E8-XT Criteria</span>
          </div>
          <div className="text-2xl font-black text-white">ΔT &lt; 10.0°C Normal</div>
          <p className="text-xs text-slate-400 leading-relaxed">
            Temperature difference between phase conductors and GOABS jaw contacts:
            <br />• ΔT &lt; 10°C: Safe operating range.
            <br />• ΔT 10-25°C: Maintenance order within 48 hours.
            <br />• ΔT &gt; 25°C: Critical fault risk, immediate emergency isolation.
          </p>
        </div>

        {/* Card 3: Right of Way */}
        <div className="bg-slate-900/70 border border-slate-800 rounded-xl p-4 shadow-sm space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold text-sky-400 uppercase tracking-wider">Right of Way (ROW)</span>
            <span className="text-xs bg-sky-950 text-sky-300 px-2 py-0.5 rounded font-mono border border-sky-800/40">15m Corridor</span>
          </div>
          <div className="text-2xl font-black text-white">33kV Clearance Zone</div>
          <p className="text-xs text-slate-400 leading-relaxed">
            7.5 meters horizontal clearance on either side of the centerline must be free from eucalyptus, mango branches, or unauthorized structures to prevent high-impedance vegetation flashovers during rain.
          </p>
        </div>
      </div>

      {/* 5-Step Safety Isolation Protocol Bar */}
      <div className="bg-slate-950 border border-slate-800 rounded-xl p-5 shadow-lg space-y-3">
        <h3 className="text-sm font-bold text-amber-400 flex items-center gap-2">
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
          Mandatory NERC 5-Step Electrical Safety Isolation Protocol
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-5 gap-3">
          {[
            { step: "1", title: "Disconnect", desc: "Open isolator switches & rack out circuit breaker." },
            { step: "2", title: "Lock & Tag", desc: "Affix padlock and danger warning DO NOT OPERATE tags." },
            { step: "3", title: "Test Zero V", desc: "Use calibrated detector stick to prove dead on all phases." },
            { step: "4", title: "Earth & Short", desc: "Apply heavy-duty portable earth clamps to conductors." },
            { step: "5", title: "Barricade", desc: "Cordon off live adjacent zones with high-vis safety tape." },
          ].map((s) => (
            <div key={s.step} className="bg-slate-900 border border-slate-800 rounded-lg p-3 space-y-1">
              <div className="flex items-center gap-2">
                <span className="w-5 h-5 rounded-full bg-amber-500/20 text-amber-400 font-bold text-xs flex items-center justify-center border border-amber-500/40">
                  {s.step}
                </span>
                <span className="text-xs font-bold text-white">{s.title}</span>
              </div>
              <p className="text-[11px] text-slate-400 leading-snug">{s.desc}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Line Patrol Register Table */}
      <div className="bg-slate-900/70 border border-slate-800 rounded-xl p-5 shadow-sm space-y-4">
        <div className="flex items-center justify-between border-b border-slate-800 pb-3">
          <div>
            <h3 className="text-base font-semibold text-white">Feeder Line Patrol & Diagnostics Register</h3>
            <p className="text-xs text-slate-400 mt-0.5">Bi-weekly preventive patrol logs across all 14 corridors.</p>
          </div>
          <span className="text-xs font-mono px-2.5 py-1 rounded bg-slate-800 text-slate-300">
            {patrols.length} Patrols Logged
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs border-collapse">
            <thead>
              <tr className="border-b border-slate-800 text-slate-400 uppercase tracking-wider font-semibold">
                <th className="py-2.5 px-3">Patrol ID / Corridor</th>
                <th className="py-2.5 px-3">Pole Span & Patroller</th>
                <th className="py-2.5 px-3">Vegetation ROW</th>
                <th className="py-2.5 px-3">Insulators</th>
                <th className="py-2.5 px-3">Hotspot ΔT</th>
                <th className="py-2.5 px-3">Earth Pit (Ω)</th>
                <th className="py-2.5 px-3">Audit Result</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60 text-slate-300">
              {patrols.map((p) => (
                <tr key={p.id} className="hover:bg-slate-800/40 transition-colors">
                  <td className="py-3 px-3">
                    <div className="font-semibold text-white">{p.patrol_code}</div>
                    <div className="text-[11px] text-amber-400 font-mono">{p.corridor.replace("FEEDER_", "")}</div>
                  </td>
                  <td className="py-3 px-3">
                    <div className="text-slate-200">{p.pole_span}</div>
                    <div className="text-slate-400 text-[11px]">{p.lead_patroller} • {p.patrol_date}</div>
                  </td>
                  <td className="py-3 px-3">
                    <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                      p.vegetation_risk_level === "SAFE"
                        ? "bg-emerald-950 text-emerald-300 border border-emerald-800/40"
                        : p.vegetation_risk_level === "ATTENTION_REQUIRED"
                        ? "bg-amber-950 text-amber-300 border border-amber-800/40"
                        : "bg-rose-950 text-rose-300 border border-rose-800/40"
                    }`}>
                      {p.vegetation_risk_level.replace(/_/g, " ")}
                    </span>
                  </td>
                  <td className="py-3 px-3">
                    <span className="text-slate-200">{p.insulators_condition.replace(/_/g, " ")}</span>
                  </td>
                  <td className="py-3 px-3">
                    <span className={`font-mono font-bold ${p.hotspot_delta_c < 10 ? "text-emerald-400" : "text-amber-400"}`}>
                      +{p.hotspot_delta_c}°C
                    </span>
                  </td>
                  <td className="py-3 px-3">
                    <span className={`font-mono font-bold ${p.earth_resistance_ohms < 5.0 ? "text-emerald-400" : "text-rose-400"}`}>
                      {p.earth_resistance_ohms} Ω
                    </span>
                  </td>
                  <td className="py-3 px-3">
                    <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                      p.compliance_status === "PASSED"
                        ? "bg-emerald-900/80 text-emerald-200 border border-emerald-700/50"
                        : "bg-rose-900/80 text-rose-200 border border-rose-700/50"
                    }`}>
                      {p.compliance_status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Add Patrol Modal */}
      {showAddModal && (
        <div className="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-700 rounded-xl p-6 max-w-lg w-full space-y-4 shadow-2xl">
            <div className="flex items-center justify-between border-b border-slate-800 pb-3">
              <h3 className="text-base font-bold text-white">Record Line Patrol Inspection</h3>
              <button onClick={() => setShowAddModal(false)} className="text-slate-400 hover:text-white font-bold text-sm">✕</button>
            </div>
            <form onSubmit={handleAddPatrol} className="space-y-3 text-xs">
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Feeder Corridor</label>
                <select
                  value={corridor}
                  onChange={(e) => setCorridor(e.target.value as FeederCorridor)}
                  style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                  className="w-full px-3 py-2 rounded border border-slate-700"
                >
                  <option value="FEEDER_33KV_ANGLO_JOS" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Anglo-Jos 33kV</option>
                  <option value="FEEDER_33KV_BUKURU" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Bukuru 33kV</option>
                  <option value="FEEDER_33KV_RAYFIELD" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Rayfield 33kV</option>
                  <option value="FEEDER_33KV_MAKERI" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Makeri 33kV</option>
                  <option value="FEEDER_33KV_BOKKOS" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Bokkos 33kV</option>
                  <option value="FEEDER_11KV_VOM" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Vom 11kV</option>
                  <option value="FEEDER_11KV_ZAWAN" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Zawan 11kV</option>
                </select>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Pole Span</label>
                  <input
                    type="text"
                    value={poleSpan}
                    onChange={(e) => setPoleSpan(e.target.value)}
                    style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                    className="w-full px-2.5 py-1.5 rounded border border-slate-700"
                    required
                  />
                </div>
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Lead Patroller</label>
                  <input
                    type="text"
                    value={leadPatroller}
                    onChange={(e) => setLeadPatroller(e.target.value)}
                    style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                    className="w-full px-2.5 py-1.5 rounded border border-slate-700"
                    required
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Vegetation Risk</label>
                  <select
                    value={vegRisk}
                    onChange={(e) => setVegRisk(e.target.value as any)}
                    style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                    className="w-full px-2.5 py-1.5 rounded border border-slate-700"
                  >
                    <option value="SAFE" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Safe (&gt;15m clear)</option>
                    <option value="ATTENTION_REQUIRED" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Attention Required</option>
                    <option value="IMMINENT_FLASHOVER" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Imminent Flashover</option>
                  </select>
                </div>
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Insulator Condition</label>
                  <select
                    value={insulatorCond}
                    onChange={(e) => setInsulatorCond(e.target.value as any)}
                    style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                    className="w-full px-2.5 py-1.5 rounded border border-slate-700"
                  >
                    <option value="ALL_INTACT" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>All Intact</option>
                    <option value="CRACKED_PIN_FOUND" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Cracked Pin Found</option>
                    <option value="FLASHED_DISC" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Flashed Disc</option>
                  </select>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Hotspot ΔT (°C)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={hotspotDelta}
                    onChange={(e) => setHotspotDelta(e.target.value)}
                    style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                    className="w-full px-2.5 py-1.5 rounded border border-slate-700"
                    required
                  />
                </div>
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Earth Resistance (Ohms)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={earthOhms}
                    onChange={(e) => setEarthOhms(e.target.value)}
                    style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                    className="w-full px-2.5 py-1.5 rounded border border-slate-700"
                    required
                  />
                </div>
              </div>
              <div className="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowAddModal(false)}
                  className="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded font-semibold"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-1.5 bg-amber-500 hover:bg-amber-600 text-slate-950 font-bold rounded"
                >
                  Save Patrol Record
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
export default FeederStandardsView;

