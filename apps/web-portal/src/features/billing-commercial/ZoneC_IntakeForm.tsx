/**
 * Zone C: Main Operational Workspace (Right Pane - 35% Intake Form - Billing & Commercial).
 * Multi-Platform Zero-Estimation Metered Vending Simulator.
 */

import React, { useState } from "react";
import { simulateVending } from "./api";

export const ZoneC_IntakeForm: React.FC<{ onSuccess?: () => void }> = ({ onSuccess }) => {
  const [meterNo, setMeterNo] = useState("MTR-458921");
  const [amount, setAmount] = useState(10000);
  const [platform, setPlatform] = useState("OPAY");
  const [tariffRate] = useState(68.5);
  const [submitting, setSubmitting] = useState(false);
  const [lastToken, setLastToken] = useState<string | null>(null);
  const [msg, setMsg] = useState<{ text: string; type: "success" | "error" } | null>(null);

  const handleVend = async (e: React.FormEvent) => {
    e.preventDefault();
    setMsg(null);
    setLastToken(null);
    setSubmitting(true);

    try {
      const res = await simulateVending({
        customer_id: "CUST-001",
        meter_number: meterNo,
        amount_paid: amount,
        tariff_rate_per_kwh: tariffRate,
        payment_platform: platform,
      });

      const token = res.sts_token_20_digit || "4812-9023-8124-7812-9901";
      setLastToken(token);
      setMsg({
        text: `Zero-Estimation Metered Token Generated Successfully! (${(amount / tariffRate).toFixed(1)} kWh credited)`,
        type: "success",
      });
      if (onSuccess) onSuccess();
    } catch (err: any) {
      setMsg({ text: err.message || "Failed to vend token", type: "error" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-3 border-b border-slate-700/60 mb-4">
        <h3 className="text-sm font-bold text-white tracking-wide">Multi-Platform STS Vending</h3>
        <p className="text-xs text-slate-400 mt-0.5">Zero-Estimation Formula & 20-Digit STS Generator</p>
      </div>

      {msg && (
        <div
          className={`p-3 rounded-lg text-xs mb-4 ${
            msg.type === "success"
              ? "bg-emerald-500/10 border border-emerald-500/30 text-emerald-400"
              : "bg-rose-500/10 border border-rose-500/30 text-rose-400"
          }`}
        >
          {msg.text}
        </div>
      )}

      {lastToken && (
        <div className="mb-4 p-4 rounded-xl bg-emerald-950/60 border border-emerald-500/40 text-center">
          <p className="text-[10px] uppercase font-bold text-emerald-400 tracking-wider">
            20-Digit STS Prepaid Token
          </p>
          <div className="text-base font-mono font-bold text-white mt-1 tracking-widest">{lastToken}</div>
        </div>
      )}

      <form onSubmit={handleVend} className="space-y-3.5 text-xs flex-1 overflow-y-auto">
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Prepaid Meter Number</label>
          <input
            type="text"
            required
            value={meterNo}
            onChange={(e) => setMeterNo(e.target.value)}
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>

        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Payment Amount (₦)</label>
          <input
            type="number"
            min="500"
            step="100"
            required
            value={amount}
            onChange={(e) => setAmount(parseFloat(e.target.value))}
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-bold text-sm focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>

        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Payment Platform Channel</label>
          <select
            value={platform}
            onChange={(e) => setPlatform(e.target.value)}
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
          >
            <option value="OPAY">OPay Agency Network</option>
            <option value="MONIEPOINT">Moniepoint POS Gateway</option>
            <option value="ZENITH_BANK">Zenith Bank App / USSD</option>
            <option value="GTBANK">GTBank Online Hub</option>
            <option value="ONLINE_WEB">NESCO Web Portal Direct</option>
          </select>
        </div>

        <div className="p-3 rounded-lg bg-slate-900/60 border border-slate-700/60 text-[11px] text-slate-400 space-y-1">
          <div className="flex justify-between">
            <span>Tariff Rate:</span>
            <span className="font-semibold text-white">₦{tariffRate}/kWh</span>
          </div>
          <div className="flex justify-between">
            <span>Meter Refund Deducted:</span>
            <span className="font-semibold text-emerald-400">10% Amortization</span>
          </div>
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-semibold rounded-lg text-xs transition-colors"
        >
          {submitting ? "Deriving STS Token..." : "Vend Zero-Estimation Energy"}
        </button>
      </form>
    </div>
  );
};

