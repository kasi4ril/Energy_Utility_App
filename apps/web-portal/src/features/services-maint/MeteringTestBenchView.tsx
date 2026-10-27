/**
 * Meter Testing & Recalibration Bench for Metering Services Unit.
 * Logs laboratory accuracy verification and NEMSA calibration certificates.
 */

import React, { useState } from "react";

interface MeterTestCertificate {
  id: string;
  cert_number: string;
  meter_serial: string;
  meter_class: "CLASS_0_5S" | "CLASS_1_0" | "CLASS_2_0";
  creep_test: "PASSED_NO_ROTATION" | "FAILED_CREEP";
  starting_current_test: "PASSED_SENSITIVE" | "FAILED";
  full_load_error_pct: number;
  nemsa_pass: boolean;
  tested_by_staff_id: string;
  cert_date_wat: string;
}

const INITIAL_CERTS: MeterTestCertificate[] = [
  {
    id: "CERT-01",
    cert_number: "NEMSA-CAL-2026-441",
    meter_serial: "0419-3321-8842",
    meter_class: "CLASS_1_0",
    creep_test: "PASSED_NO_ROTATION",
    starting_current_test: "PASSED_SENSITIVE",
    full_load_error_pct: 0.22,
    nemsa_pass: true,
    tested_by_staff_id: "STAFF-MET-03",
    cert_date_wat: "2026-10-27 10:00 WAT",
  },
  {
    id: "CERT-02",
    cert_number: "NEMSA-CAL-2026-440",
    meter_serial: "6214-9981-2201",
    meter_class: "CLASS_0_5S",
    creep_test: "PASSED_NO_ROTATION",
    starting_current_test: "PASSED_SENSITIVE",
    full_load_error_pct: 0.08,
    nemsa_pass: true,
    tested_by_staff_id: "STAFF-MET-03",
    cert_date_wat: "2026-10-26 14:30 WAT",
  },
];

export function MeteringTestBenchView() {
  const [certs] = useState<MeterTestCertificate[]>(INITIAL_CERTS);

  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>🔬</span>
            <span>NEMSA Certified Meter Test &amp; Recalibration Bench</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            IEC 62053-21/22 precision testing: Full Load, Light Load, Creep, and Starting Current.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-emerald-950 text-emerald-400 border border-emerald-800">
          NEMSA Accredited Lab
        </span>
      </div>

      {/* Certificates Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
          <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
            Issued Calibration &amp; Accuracy Certificates
          </span>
          <span className="text-xs text-slate-400 font-mono">Precision Ledger</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Certificate No</th>
                <th className="py-3 px-3">Meter Serial</th>
                <th className="py-3 px-3">Accuracy Class</th>
                <th className="py-3 px-3">Creep &bull; Start Current</th>
                <th className="py-3 px-3">Full Load Error</th>
                <th className="py-3 px-4 text-right">NEMSA Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {certs.map((c) => (
                <tr key={c.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-cyan-400">
                    {c.cert_number}
                    <span className="block text-[10px] text-slate-500 font-normal">{c.cert_date_wat}</span>
                  </td>
                  <td className="py-3 px-3 font-mono font-bold text-white">{c.meter_serial}</td>
                  <td className="py-3 px-3">
                    <span className="font-mono text-cyan-300 font-bold">
                      {c.meter_class.replace(/_/g, " ")}
                    </span>
                  </td>
                  <td className="py-3 px-3 text-[11px] text-slate-300">
                    <div>✓ No Creep Rotation</div>
                    <div className="text-[10px] text-emerald-400 font-mono">Starting Current &lt; 20mA</div>
                  </td>
                  <td className="py-3 px-3 font-mono font-bold text-emerald-400">
                    &plusmn;{c.full_load_error_pct}%
                  </td>
                  <td className="py-3 px-4 text-right">
                    <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800 font-bold">
                      CERTIFIED PASS
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

