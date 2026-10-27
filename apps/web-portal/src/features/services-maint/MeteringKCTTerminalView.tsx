/**
 * STS Key Change Token (KCT1/KCT2) Terminal for Metering Services.
 * Generates paired 20-digit Key Change Tokens for TID Rollover and SGC migration.
 */

import React, { useState } from "react";

interface KCTRecord {
  id: string;
  meter_number: string;
  customer_name: string;
  kct1_20_digit: string;
  kct2_20_digit: string;
  migration_reason: "TID_ROLLOVER_2024" | "SUPPLY_GROUP_CODE_TRANSFER" | "TARIFF_KEY_UPGRADE";
  status: "ENTERED_CONFIRMED" | "ISSUED_PENDING_INPUT";
  issued_at_wat: string;
}

const INITIAL_KCTS: KCTRecord[] = [
  {
    id: "KCT-01",
    meter_number: "0419-3321-8842",
    customer_name: "Dr. Joshua M. Pam",
    kct1_20_digit: "5821-0492-1193-8821-4491",
    kct2_20_digit: "9144-2201-7731-9042-5512",
    migration_reason: "TID_ROLLOVER_2024",
    status: "ENTERED_CONFIRMED",
    issued_at_wat: "2026-10-27 12:15 WAT",
  },
  {
    id: "KCT-02",
    meter_number: "6214-9981-2201",
    customer_name: "Grand Cereals Nigeria Ltd",
    kct1_20_digit: "4412-8891-0023-5591-3381",
    kct2_20_digit: "1190-3341-8892-7712-4401",
    migration_reason: "TARIFF_KEY_UPGRADE",
    status: "ISSUED_PENDING_INPUT",
    issued_at_wat: "2026-10-27 15:30 WAT",
  },
];

export function MeteringKCTTerminalView() {
  const [records, setRecords] = useState<KCTRecord[]>(INITIAL_KCTS);
  const [meterNo, setMeterNo] = useState("");
  const [custName, setCustName] = useState("");
  const [reason, setReason] = useState<"TID_ROLLOVER_2024" | "SUPPLY_GROUP_CODE_TRANSFER" | "TARIFF_KEY_UPGRADE">("TID_ROLLOVER_2024");
  const [activeKctPair, setActiveKctPair] = useState<{ kct1: string; kct2: string } | null>(null);

  const handleGenerateKCT = (e: React.FormEvent) => {
    e.preventDefault();
    if (!meterNo) {
      alert("Please provide the customer meter number.");
      return;
    }

    const r = () => Math.floor(1000 + Math.random() * 9000);
    const k1 = `${r()}-${r()}-${r()}-${r()}-${r()}`;
    const k2 = `${r()}-${r()}-${r()}-${r()}-${r()}`;

    const newRecord: KCTRecord = {
      id: `KCT-${Date.now()}`,
      meter_number: meterNo,
      customer_name: custName || "Premises Meter " + meterNo,
      kct1_20_digit: k1,
      kct2_20_digit: k2,
      migration_reason: reason,
      status: "ISSUED_PENDING_INPUT",
      issued_at_wat: "Just Now",
    };

    setRecords([newRecord, ...records]);
    setActiveKctPair({ kct1: k1, kct2: k2 });
    alert(`Paired Key Change Tokens successfully generated! KCT 1: ${k1} | KCT 2: ${k2}`);
  };

  return (
    <div className="space-y-6">
      {/* Intake Form */}
      <div className="p-5 rounded-xl bg-slate-900/80 border border-slate-800">
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>🔑</span>
            <span>STS Key Change Token (KCT) Issuance Terminal</span>
          </h3>
          <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-cyan-950 text-cyan-400 border border-cyan-800">
            STS Association Certified
          </span>
        </div>

        <form onSubmit={handleGenerateKCT} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
          <div>
            <label className="text-slate-400 font-semibold block mb-1">Meter Serial Number</label>
            <input
              type="text"
              value={meterNo}
              onChange={(e) => setMeterNo(e.target.value)}
              placeholder="e.g. 0419-3321-8842"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold focus:border-cyan-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Customer / Property Name</label>
            <input
              type="text"
              value={custName}
              onChange={(e) => setCustName(e.target.value)}
              placeholder="e.g. Dr. Joshua M. Pam"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 focus:border-cyan-400 focus:outline-none"
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Key Change Reason</label>
            <select
              value={reason}
              onChange={(e) => setReason(e.target.value as any)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-semibold focus:border-cyan-400 focus:outline-none"
            >
              <option value="TID_ROLLOVER_2024" className="bg-[#0f172a] text-[#f8fafc]">TID Rollover (2024 Base Date)</option>
              <option value="SUPPLY_GROUP_CODE_TRANSFER" className="bg-[#0f172a] text-[#f8fafc]">Supply Group Code (SGC) Transfer</option>
              <option value="TARIFF_KEY_UPGRADE" className="bg-[#0f172a] text-[#f8fafc]">Tariff Band Key Reclassification</option>
            </select>
          </div>

          <div className="flex items-end">
            <button
              type="submit"
              className="w-full py-2.5 px-4 rounded-lg bg-cyan-600 hover:bg-cyan-500 text-white font-bold transition shadow-lg shadow-cyan-950 flex items-center justify-center space-x-1"
            >
              <span>⚡</span>
              <span>Generate KCT 1 &amp; KCT 2</span>
            </button>
          </div>
        </form>

        {/* Active Pair Alert */}
        {activeKctPair && (
          <div className="mt-4 p-4 rounded-xl bg-cyan-950/80 border border-cyan-500/50 space-y-2">
            <div className="text-[10px] uppercase font-bold text-cyan-300 font-mono">
              Paired 20-Digit Key Change Tokens Issued
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div className="bg-slate-950 p-2.5 rounded border border-slate-800">
                <span className="text-[10px] text-slate-400 block font-mono">Step 1 &bull; KCT 1 (First Input):</span>
                <span className="text-base font-mono font-black text-cyan-300 tracking-wider">
                  {activeKctPair.kct1}
                </span>
              </div>
              <div className="bg-slate-950 p-2.5 rounded border border-slate-800">
                <span className="text-[10px] text-slate-400 block font-mono">Step 2 &bull; KCT 2 (Second Input):</span>
                <span className="text-base font-mono font-black text-emerald-300 tracking-wider">
                  {activeKctPair.kct2}
                </span>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Register */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
          <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
            KCT Key Change Issuance Ledger
          </span>
          <span className="text-xs text-slate-400 font-mono">STS Compliance</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Meter &bull; Customer</th>
                <th className="py-3 px-3">Reason</th>
                <th className="py-3 px-3">KCT 1 (First Token)</th>
                <th className="py-3 px-3">KCT 2 (Second Token)</th>
                <th className="py-3 px-4 text-right">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {records.map((r) => (
                <tr key={r.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4">
                    <div className="font-bold text-white">{r.customer_name}</div>
                    <div className="text-[10px] text-cyan-400 font-mono">{r.meter_number}</div>
                  </td>
                  <td className="py-3 px-3 text-slate-300 font-semibold">{r.migration_reason.replace(/_/g, " ")}</td>
                  <td className="py-3 px-3 font-mono font-bold text-slate-200 text-[11px]">{r.kct1_20_digit}</td>
                  <td className="py-3 px-3 font-mono font-bold text-emerald-300 text-[11px]">{r.kct2_20_digit}</td>
                  <td className="py-3 px-4 text-right">
                    {r.status === "ENTERED_CONFIRMED" ? (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                        CONFIRMED MIGRATED
                      </span>
                    ) : (
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-amber-950 text-amber-400 border border-amber-800">
                        PENDING ENTRY
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

