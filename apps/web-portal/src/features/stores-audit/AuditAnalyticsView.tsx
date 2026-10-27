/**
 * Sovereign Audit & Compliance Analytics Dashboard for Internal Audit.
 * Tracks risk heatmaps, open observations, and departmental compliance scores.
 */

import React from "react";
import { MOCK_AUDIT_FINDINGS } from "./auditMockData";

const DEPT_SCORES = [
  { dept: "Generation Division", score: 98.4, status: "EXCELLENT" },
  { dept: "Energy Management", score: 99.1, status: "EXCELLENT" },
  { dept: "Services & Maintenance", score: 95.8, status: "COMPLIANT" },
  { dept: "Central Stores", score: 94.2, status: "COMPLIANT" },
  { dept: "Commercial Operations", score: 92.5, status: "ATTENTION_REQUIRED" },
];

export function AuditAnalyticsView() {
  return (
    <div className="space-y-6">
      {/* Top KPI Ribbon */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Composite Compliance Index</div>
          <div className="text-2xl font-black text-emerald-400 mt-1">96.8%</div>
          <div className="text-[10px] text-slate-500 mt-0.5">Weighted Average across 5 Divisions</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Open Audit Observations</div>
          <div className="text-2xl font-black text-amber-400 mt-1">1 High Risk</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Commercial CT Multiplier Audit</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Financial Variance Tracked</div>
          <div className="text-2xl font-black text-white mt-1">&#8358;675,000</div>
          <div className="text-[10px] text-emerald-400 mt-0.5">100% Under Remediation</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Cryptographic Ledger Check</div>
          <div className="text-2xl font-black text-cyan-400 mt-1">0 Breaches</div>
          <div className="text-[10px] text-cyan-400 mt-0.5">SHA-256 Prev-Hash Chaining Valid</div>
        </div>
      </div>

      {/* Department Compliance Scores & Observations */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Compliance Index by Department */}
        <div className="lg:col-span-5 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-white">Departmental Compliance Ratings</h3>
            <p className="text-xs text-slate-400 mt-0.5">
              Internal control effectiveness score (0 - 100%).
            </p>
          </div>

          <div className="space-y-4 my-4">
            {DEPT_SCORES.map((d) => (
              <div key={d.dept} className="space-y-1">
                <div className="flex justify-between text-xs">
                  <span className="text-slate-300 font-semibold">{d.dept}</span>
                  <span className="font-mono font-bold text-white">{d.score}%</span>
                </div>
                <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                  <div
                    className={`h-full rounded-full transition-all ${
                      d.score >= 98
                        ? "bg-emerald-400"
                        : d.score >= 95
                        ? "bg-cyan-400"
                        : "bg-amber-400"
                    }`}
                    style={{ width: `${d.score}%` }}
                  />
                </div>
              </div>
            ))}
          </div>

          <div className="border-t border-slate-800/80 pt-3 text-[11px] text-slate-400 flex items-center justify-between">
            <span>Overall Internal Audit Grade:</span>
            <span className="font-mono text-emerald-400 font-bold">Grade A (Substantial Assurance)</span>
          </div>
        </div>

        {/* Open Observations Table */}
        <div className="lg:col-span-7 rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60 flex flex-col justify-between">
          <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
            <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
              Active Audit Observations &amp; Remediation Register
            </span>
            <span className="text-xs text-slate-400 font-mono">Quarterly Cycle</span>
          </div>

          <div className="overflow-x-auto flex-1">
            <table className="w-full text-left text-xs">
              <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
                <tr>
                  <th className="py-2.5 px-3">Audit Ref</th>
                  <th className="py-2.5 px-3">Finding Title</th>
                  <th className="py-2.5 px-3">Risk Level</th>
                  <th className="py-2.5 px-3 text-right">Remediation Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800/80 text-slate-200">
                {MOCK_AUDIT_FINDINGS.map((f) => (
                  <tr key={f.id} className="hover:bg-slate-800/40 transition">
                    <td className="py-2.5 px-3 font-mono font-bold text-rose-400">
                      {f.audit_ref}
                      <span className="block text-[10px] text-slate-500 font-normal">{f.department}</span>
                    </td>
                    <td className="py-2.5 px-3 text-slate-200 max-w-xs">{f.finding_title}</td>
                    <td className="py-2.5 px-3">
                      {f.risk_level === "HIGH" && (
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-rose-950 text-rose-400 border border-rose-800 font-bold">
                          HIGH
                        </span>
                      )}
                      {f.risk_level === "MEDIUM" && (
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                          MEDIUM
                        </span>
                      )}
                      {f.risk_level === "LOW" && (
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-slate-800 text-slate-300 border border-slate-700">
                          LOW
                        </span>
                      )}
                    </td>
                    <td className="py-2.5 px-3 text-right">
                      {f.status === "OPEN_OBSERVATION" && (
                        <span className="text-amber-400 font-mono text-[11px]">Under Review</span>
                      )}
                      {f.status === "MANAGEMENT_RESPONSE_RECEIVED" && (
                        <span className="text-cyan-400 font-mono text-[11px]">Response In</span>
                      )}
                      {f.status === "RECTIFIED_CLOSED" && (
                        <span className="text-emerald-400 font-mono text-[11px]">Closed</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="p-3 border-t border-slate-800/80 text-[11px] text-slate-400 flex items-center justify-between">
            <span>Remediation tracked against executive deadlines.</span>
            <span className="font-mono text-cyan-400">Audit Desk Gate</span>
          </div>
        </div>
      </div>
    </div>
  );
}

