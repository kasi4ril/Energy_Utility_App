/**
 * View 4: Prepaid Token, TID Rollover (KCT 1 & 2) & Meter Scheduler Desk.
 * Compliant with Section 12.3 Page 4: STS Token History lookup, TID Rollover
 * Key Change Token generator (20-digit STS tokens), and Metering Unit scheduler.
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import { fetchRecentTokens, generateTIDRolloverKCT, scheduleMeterActivity } from "./api";
import { STSTokenRecord, TIDRolloverRecord } from "./types";
import { KeyRound, Calendar, History, Copy, Check, ShieldCheck, Zap } from "lucide-react";

interface View4Props {
  staffId: string;
  staffName?: string;
}

export function View4_TokenAndKCTDesk({
  staffId,
  staffName = "Amina Yusuf (CC Officer)",
}: View4Props) {
  const [meterSerial, setMeterSerial] = useState("MTR-2026-0921");
  const [tokens, setTokens] = useState<STSTokenRecord[]>([]);
  const [kctRecord, setKctRecord] = useState<TIDRolloverRecord | null>(null);
  const [copiedToken, setCopiedToken] = useState<string | null>(null);
  const [isLoadingTokens, setIsLoadingTokens] = useState(false);
  const [isGeneratingKCT, setIsGeneratingKCT] = useState(false);

  // Meter Scheduler State
  const [schedCustomerId, setSchedCustomerId] = useState("CUST-KURA-2026-001");
  const [scheduleType, setScheduleType] = useState("NEW_INSTALLATION");
  const [scheduledDate, setScheduledDate] = useState("2026-10-28");
  const [schedSuccess, setSchedSuccess] = useState<string | null>(null);

  const handleFetchTokens = async () => {
    setIsLoadingTokens(true);
    try {
      const res = await fetchRecentTokens(meterSerial);
      setTokens(res);
    } finally {
      setIsLoadingTokens(false);
    }
  };

  const handleGenerateKCT = async () => {
    setIsGeneratingKCT(true);
    try {
      const res = await generateTIDRolloverKCT(meterSerial, staffId);
      setKctRecord(res);
    } finally {
      setIsGeneratingKCT(false);
    }
  };

  const handleCopy = (text: string, label: string) => {
    navigator.clipboard?.writeText(text);
    setCopiedToken(label);
    setTimeout(() => setCopiedToken(null), 2000);
  };

  const handleScheduleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await scheduleMeterActivity({
        customer_id: schedCustomerId,
        schedule_type: scheduleType,
        scheduled_date: scheduledDate,
        staff_id: staffId,
        staff_name: staffName,
      });
      setSchedSuccess(`Meter appointment booked for ${scheduledDate}. Routed to Metering Field Queue.`);
    } catch (err: any) {
      alert("Failed to schedule meter appointment");
    }
  };

  return (
    <div className="space-y-6 max-w-4xl text-xs">
      {/* 1. Header & Meter Lookup Bar */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5">
        <div className="flex items-center space-x-3 mb-4">
          <div className="p-2.5 rounded-lg bg-blue-600/20 text-blue-400 border border-blue-500/30">
            <KeyRound className="w-6 h-6" />
          </div>
          <div>
            <h2 className="text-base font-bold text-slate-100 uppercase tracking-wide">
              Prepaid STS Token, TID Rollover &amp; Meter Scheduler
            </h2>
            <p className="text-xs text-slate-400">
              Lookup token transactions &bull; Generate TID Key Change Tokens (KCT 1 &amp; 2) &bull; Meter Work Orders
            </p>
          </div>
        </div>

        <div className="flex flex-col sm:flex-row items-center gap-2.5">
          <input
            type="text"
            placeholder="Enter 11-digit STS Meter Serial..."
            value={meterSerial}
            onChange={(e) => setMeterSerial(e.target.value)}
            className="flex-1 bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 font-mono text-xs focus:border-blue-500 focus:outline-none"
          />
          <button
            onClick={handleFetchTokens}
            disabled={isLoadingTokens}
            className="bg-blue-600 hover:bg-blue-500 text-white px-3.5 py-2 rounded-lg font-semibold flex items-center space-x-1.5 transition-colors"
          >
            <History className="w-3.5 h-3.5" />
            <span>Lookup Last 5 Tokens</span>
          </button>
          <button
            onClick={handleGenerateKCT}
            disabled={isGeneratingKCT}
            className="bg-indigo-600 hover:bg-indigo-500 text-white px-3.5 py-2 rounded-lg font-semibold flex items-center space-x-1.5 transition-colors"
          >
            <Zap className="w-3.5 h-3.5" />
            <span>Generate KCT 1 &amp; 2</span>
          </button>
        </div>
      </div>

      {/* 2. TID Rollover Key Change Token (KCT 1 & 2) Display Panel */}
      {kctRecord && (
        <div className="bg-gradient-to-r from-indigo-950/70 to-slate-900 border border-indigo-500/40 rounded-xl p-4 sm:p-5 shadow-lg">
          <div className="flex items-center justify-between mb-3 border-b border-indigo-500/20 pb-2">
            <div className="flex items-center space-x-2">
              <ShieldCheck className="w-4 h-4 text-emerald-400" />
              <span className="font-bold text-slate-100 uppercase text-xs">
                STS Key Change Tokens Generated (TID Rollover 2024+)
              </span>
            </div>
            <span className="text-[10px] bg-indigo-900/60 text-indigo-300 px-2 py-0.5 rounded font-mono">
              Meter: {kctRecord.meter_serial}
            </span>
          </div>

          <p className="text-slate-300 text-xs mb-3">
            Instruct the customer to punch <strong>KCT 1</strong> first into the meter keypad, press Enter, wait for "ACCEPT", and then punch <strong>KCT 2</strong>.
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {/* KCT 1 */}
            <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3">
              <div className="text-[10px] text-slate-400 uppercase font-mono mb-1">Key Change Token 1 (KCT 1)</div>
              <div className="text-sm font-bold font-mono text-emerald-400 tracking-wider flex items-center justify-between">
                <span>{kctRecord.kct_1}</span>
                <button
                  onClick={() => handleCopy(kctRecord.kct_1, "kct1")}
                  className="p-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 transition-colors"
                  title="Copy KCT 1"
                >
                  {copiedToken === "kct1" ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
                </button>
              </div>
            </div>

            {/* KCT 2 */}
            <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3">
              <div className="text-[10px] text-slate-400 uppercase font-mono mb-1">Key Change Token 2 (KCT 2)</div>
              <div className="text-sm font-bold font-mono text-blue-400 tracking-wider flex items-center justify-between">
                <span>{kctRecord.kct_2}</span>
                <button
                  onClick={() => handleCopy(kctRecord.kct_2, "kct2")}
                  className="p-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 transition-colors"
                  title="Copy KCT 2"
                >
                  {copiedToken === "kct2" ? <Check className="w-3.5 h-3.5 text-blue-400" /> : <Copy className="w-3.5 h-3.5" />}
                </button>
              </div>
            </div>
          </div>

          <div className="mt-3 text-[10px] text-slate-400 font-mono flex items-center justify-between">
            <span>Operator: {staffName} ({staffId})</span>
            <span>Generated: {new Date(kctRecord.generated_at).toLocaleString()}</span>
          </div>
        </div>
      )}

      {/* 3. Recent STS Token Transaction Ledger */}
      {tokens.length > 0 && (
        <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5">
          <div className="flex items-center space-x-2 mb-3">
            <History className="w-4 h-4 text-blue-400" />
            <h3 className="text-xs font-bold text-slate-100 uppercase tracking-wide">
              Last 5 STS Vending Records for {meterSerial}
            </h3>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead>
                <tr className="border-b border-slate-800 text-slate-400 text-[10px] uppercase">
                  <th className="pb-2">20-Digit STS Token Number</th>
                  <th className="pb-2 text-right">Units (kWh)</th>
                  <th className="pb-2 text-right">Amount (₦)</th>
                  <th className="pb-2 text-right">Tariff Rate</th>
                  <th className="pb-2 text-right">Purchased At</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800/60 font-mono">
                {tokens.map((t, idx) => (
                  <tr key={idx} className="hover:bg-slate-800/40">
                    <td className="py-2.5 font-bold text-slate-200 tracking-wider">{t.token_number}</td>
                    <td className="py-2.5 text-right text-emerald-400">{t.units_kwh}</td>
                    <td className="py-2.5 text-right text-slate-100">₦{t.amount_paid.toLocaleString()}</td>
                    <td className="py-2.5 text-right text-slate-400">₦{t.tariff_rate}/kWh</td>
                    <td className="py-2.5 text-right text-slate-400">
                      {new Date(t.purchased_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* 4. Meter Installation & Repair Capacity Scheduler */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5">
        <div className="flex items-center space-x-2 mb-4">
          <Calendar className="w-4 h-4 text-blue-400" />
          <h3 className="text-xs font-bold text-slate-100 uppercase tracking-wide">
            Meter Installation / Repair Slot Booking (Metering Unit Capacity)
          </h3>
        </div>

        {schedSuccess && (
          <div className="p-2.5 mb-3 rounded bg-emerald-950/60 border border-emerald-500/40 text-emerald-300 text-xs">
            {schedSuccess}
          </div>
        )}

        <form onSubmit={handleScheduleSubmit} className="grid grid-cols-1 md:grid-cols-3 gap-3">
          <div>
            <label className="block text-slate-400 mb-1">Customer Account / ID *</label>
            <input
              type="text"
              required
              value={schedCustomerId}
              onChange={(e) => setSchedCustomerId(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 font-mono"
            />
          </div>

          <div>
            <label className="block text-slate-400 mb-1">Activity Schedule Type *</label>
            <select
              value={scheduleType}
              onChange={(e) => setScheduleType(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100"
            >
              <option value="NEW_INSTALLATION">New MAP Meter Installation</option>
              <option value="REPAIR">Field Recertification &amp; Repair</option>
              <option value="CALIBRATION">Meter Accuracy Calibration</option>
            </select>
          </div>

          <div>
            <label className="block text-slate-400 mb-1">Appointment Date *</label>
            <input
              type="date"
              required
              value={scheduledDate}
              onChange={(e) => setScheduledDate(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 font-mono"
            />
          </div>

          <div className="md:col-span-3 mt-2 flex justify-end">
            <button
              type="submit"
              className="bg-blue-600 hover:bg-blue-500 text-white font-semibold py-2 px-5 rounded-lg transition-colors flex items-center space-x-1.5"
            >
              <Calendar className="w-3.5 h-3.5" />
              <span>Book Appointment &amp; Dispatch to Metering Head</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

