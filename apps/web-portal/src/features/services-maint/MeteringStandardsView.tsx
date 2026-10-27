/**
 * NERC Metering Code & IEC 62053 Standards Reference.
 * Guidelines for Meter Accuracy, Seal Integrity, and Testing Protocols.
 */

import React from "react";

export function MeteringStandardsView() {
  const standards = [
    {
      code: "NERC Metering Code §4.2",
      parameter: "Accuracy Class by Customer Category",
      limit: "Class 0.5S (Maximum Demand Industrial) | Class 1.0 (Residential & Commercial)",
      description: "All Maximum Demand meters must maintain CT/PT Class 0.5S precision. Residential meters must remain within ±1.0% error across full load.",
    },
    {
      code: "IEC 62055-41",
      parameter: "STS Token Encryption & TID Rollover",
      limit: "20-Digit STS Token | Base Date: 2014 Rollover Compliant",
      description: "Standard Transfer Specification encryption prevents token forgery. TID rollover must be completed on all legacy meters before expiration.",
    },
    {
      code: "NERC Connection Code",
      parameter: "Tamper Holographic Seals",
      limit: "Numbered Security Polycarbonate Seals (Terminal & Meter Body)",
      description: "Any seal breach or missing serial tag automatically triggers an on-site investigation and energy theft assessment.",
    },
    {
      code: "NEMSA Meter Recertification",
      parameter: "Periodic Recalibration Interval",
      limit: "Maximum 10-Year In-Service Recertification Cycle",
      description: "Meters reaching 10 years in-service must undergo test-bench recalibration or fleet replacement to guarantee fair billing.",
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>📘</span>
            <span>NERC Metering Code &amp; IEC 62053 Standards Reference</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Statutory metrology standards enforcing fair, accurate, and tamper-proof energy measurement.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-slate-800 text-slate-300 border border-slate-700">
          NEMSA Benchmark Standard
        </span>
      </div>

      {/* Standards Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {standards.map((s) => (
          <div key={s.code} className="p-5 rounded-xl bg-slate-900/70 border border-slate-800 space-y-2.5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-mono font-bold text-cyan-400">{s.code}</span>
              <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-cyan-950 text-cyan-300 border border-cyan-800">
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

