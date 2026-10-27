/**
 * NERC & Regulatory Compliance Monitor for Internal Audit.
 * Monitors statutory SLA breach risks, safety regulations, and tariff order adherence.
 */

import React from "react";
import { MOCK_NERC_COMPLIANCE } from "./auditMockData";

export function AuditComplianceView() {
  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>🛡️</span>
            <span>NERC &amp; NEMSA Regulatory Compliance Audit Monitor</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Continuous surveillance against statutory standards preventing regulatory fines and license sanctions.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-emerald-950 text-emerald-400 border border-emerald-800">
          NERC &amp; NEMSA Verified
        </span>
      </div>

      {/* Metrics Grid */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Standard Code</th>
                <th className="py-3 px-3">Regulatory Requirement</th>
                <th className="py-3 px-3">Statutory Target SLA</th>
                <th className="py-3 px-3">NESCO Performance</th>
                <th className="py-3 px-4 text-right">Compliance Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {MOCK_NERC_COMPLIANCE.map((m) => (
                <tr key={m.standard_code} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-cyan-400">{m.standard_code}</td>
                  <td className="py-3 px-3 font-semibold text-white">{m.title}</td>
                  <td className="py-3 px-3 font-mono text-slate-300">{m.target_sla}</td>
                  <td className="py-3 px-3 font-mono font-bold text-emerald-400">
                    {m.actual_performance}
                  </td>
                  <td className="py-3 px-4 text-right">
                    <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800 font-bold">
                      {m.compliance_status}
                    </span>
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

