/**
 * IEEE C57 / IEC 60076 Transformer Engineering Standards Reference.
 * Guidelines for Dielectric Oil, Temperature Limits, and Neutral Earthing.
 */

import React from "react";

export function TransformerStandardsView() {
  const standards = [
    {
      code: "IEEE C57.106",
      parameter: "Dielectric Breakdown Voltage (BDV)",
      limit: "> 40.0 kV (New/Treated) | > 30.0 kV (In-Service Min)",
      description: "Tested per ASTM D877. If BDV falls below 30 kV, immediate vacuum dehydration and oil filtration is mandatory.",
    },
    {
      code: "IEEE C57.91",
      parameter: "Top-Oil & Winding Temperature Rise",
      limit: "Top-Oil Rise: < 55°C (Max 65°C) | Winding Hotspot: < 75°C",
      description: "Sustained temperature above limits rapidly degrades cellulose insulation paper, halving transformer life per 6°C rise.",
    },
    {
      code: "IEEE 80 / NERC",
      parameter: "Transformer Neutral & Arrester Grounding",
      limit: "< 5.0 Ω (Target < 2.0 Ω)",
      description: "Substation transformer neutral ground pit and surge arresters must maintain solid grounding to dissipate lightning surges safely.",
    },
    {
      code: "IEC 60076-7",
      parameter: "Silica Gel Breather Saturation",
      limit: "Blue (Active Desiccant) | Pink (Replace Immediately)",
      description: "Moisture ingress via saturated breather breaks down oil insulation and increases combustible gas formation.",
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>📘</span>
            <span>IEEE C57 &amp; IEC 60076 Power Transformer Engineering Standards</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Prescribed operating tolerances protecting NESCO substation distribution transformers.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-slate-800 text-slate-300 border border-slate-700">
          Utility Benchmark Standard
        </span>
      </div>

      {/* Standards Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {standards.map((s) => (
          <div key={s.code} className="p-5 rounded-xl bg-slate-900/70 border border-slate-800 space-y-2.5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-mono font-bold text-amber-400">{s.code}</span>
              <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-amber-950 text-amber-300 border border-amber-800">
                Mandatory
              </span>
            </div>
            <h4 className="text-sm font-bold text-white">{s.parameter}</h4>
            <div className="text-xs font-mono font-bold text-emerald-400 bg-slate-950/80 p-2 rounded border border-slate-800">
              Limit: {s.limit}
            </div>
            <p className="text-xs text-slate-400 leading-relaxed font-sans">{s.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

