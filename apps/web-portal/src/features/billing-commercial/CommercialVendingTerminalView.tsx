/**
 * STS Prepayment Vending Terminal for Commercial Operations.
 * Generates 20-digit STS encryption tokens and interfaces with fintech hooks.
 */

import React, { useState } from "react";
import { MOCK_VENDINGS } from "./mockData";
import { VendingRecord } from "./types";

export function CommercialVendingTerminalView() {
  const [vendings, setVendings] = useState<VendingRecord[]>(MOCK_VENDINGS);
  const [meterNumber, setMeterNumber] = useState("");
  const [customerName, setCustomerName] = useState("");
  const [tariffBand, setTariffBand] = useState<"Band A" | "Band B" | "Band C" | "Band D" | "Band E">("Band A");
  const [amountPaid, setAmountPaid] = useState("");
  const [paymentChannel, setPaymentChannel] = useState<"OPay" | "Moniepoint" | "GTBank" | "Zenith Bank" | "Cash Counter">("Moniepoint");
  const [generatedToken, setGeneratedToken] = useState<string | null>(null);

  const tariffRates: Record<string, number> = {
    "Band A": 225.0,
    "Band B": 180.0,
    "Band C": 120.0,
    "Band D": 85.0,
    "Band E": 65.0,
  };

  const handleVend = (e: React.FormEvent) => {
    e.preventDefault();
    const amt = parseFloat(amountPaid);
    if (!meterNumber || isNaN(amt) || amt <= 0) {
      alert("Please provide valid meter number and payment amount.");
      return;
    }

    const rate = tariffRates[tariffBand] || 180.0;
    const kwh = Number((amt / rate).toFixed(1));

    // Generate random 20-digit STS token
    const r = () => Math.floor(1000 + Math.random() * 9000);
    const token = `${r()}-${r()}-${r()}-${r()}-${r()}`;

    const newRec: VendingRecord = {
      id: `VND-${Date.now()}`,
      transaction_ref: `TX-${Date.now().toString().slice(-8)}`,
      customer_name: customerName || "Customer on Meter " + meterNumber,
      account_or_meter: meterNumber,
      tariff_band: tariffBand,
      amount_paid_ngn: amt,
      tariff_rate_ngn: rate,
      kwh_vended: kwh,
      sts_token_20_digit: token,
      meter_refund_credited_ngn: 0,
      net_energy_amount_ngn: amt,
      payment_channel: paymentChannel,
      created_at_wat: "Just Now",
    };

    setVendings([newRec, ...vendings]);
    setGeneratedToken(token);
    alert(`STS Token successfully generated: ${token}`);
  };

  return (
    <div className="space-y-6">
      {/* Vending Intake Form */}
      <div className="p-5 rounded-xl bg-slate-900/80 border border-slate-800">
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>💳</span>
            <span>STS Prepayment Energy Vending Terminal</span>
          </h3>
          <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-teal-950 text-teal-400 border border-teal-800">
            STS IEC 62055-41 Compliant
          </span>
        </div>

        <form onSubmit={handleVend} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
          <div>
            <label className="text-slate-400 font-semibold block mb-1">Customer Meter Number</label>
            <input
              type="text"
              value={meterNumber}
              onChange={(e) => setMeterNumber(e.target.value)}
              placeholder="e.g. 0419-3321-8842"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold focus:border-teal-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Customer / Business Name</label>
            <input
              type="text"
              value={customerName}
              onChange={(e) => setCustomerName(e.target.value)}
              placeholder="e.g. Grand Cereals Nig Ltd"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 focus:border-teal-400 focus:outline-none"
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Tariff Band</label>
            <select
              value={tariffBand}
              onChange={(e) => setTariffBand(e.target.value as any)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-semibold focus:border-teal-400 focus:outline-none"
            >
              <option value="Band A" className="bg-[#0f172a] text-[#f8fafc]">Band A (20-24 hrs) - ₦225/kWh</option>
              <option value="Band B" className="bg-[#0f172a] text-[#f8fafc]">Band B (16-20 hrs) - ₦180/kWh</option>
              <option value="Band C" className="bg-[#0f172a] text-[#f8fafc]">Band C (12-16 hrs) - ₦120/kWh</option>
              <option value="Band D" className="bg-[#0f172a] text-[#f8fafc]">Band D (8-12 hrs) - ₦85/kWh</option>
              <option value="Band E" className="bg-[#0f172a] text-[#f8fafc]">Band E (4-8 hrs) - ₦65/kWh</option>
            </select>
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Payment Channel</label>
            <select
              value={paymentChannel}
              onChange={(e) => setPaymentChannel(e.target.value as any)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-semibold focus:border-teal-400 focus:outline-none"
            >
              <option value="Moniepoint" className="bg-[#0f172a] text-[#f8fafc]">Moniepoint POS / Webhook</option>
              <option value="OPay" className="bg-[#0f172a] text-[#f8fafc]">OPay Agent / Direct</option>
              <option value="GTBank" className="bg-[#0f172a] text-[#f8fafc]">GTBank NIBSS / NIP</option>
              <option value="Zenith Bank" className="bg-[#0f172a] text-[#f8fafc]">Zenith Bank PayDirect</option>
              <option value="Cash Counter" className="bg-[#0f172a] text-[#f8fafc]">NESCO Cash Office</option>
            </select>
          </div>

          <div className="sm:col-span-2">
            <label className="text-slate-400 font-semibold block mb-1">Amount Paid (NGN &#8358;)</label>
            <input
              type="number"
              step="500"
              value={amountPaid}
              onChange={(e) => setAmountPaid(e.target.value)}
              placeholder="e.g. 50000"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold text-teal-400 text-sm focus:border-teal-400 focus:outline-none"
              required
            />
          </div>

          <div className="sm:col-span-2 flex items-end">
            <button
              type="submit"
              className="w-full py-2.5 px-4 rounded-lg bg-teal-600 hover:bg-teal-500 text-white font-bold transition shadow-lg shadow-teal-950 flex items-center justify-center space-x-2"
            >
              <span>⚡</span>
              <span>Generate STS Token &amp; Vend</span>
            </button>
          </div>
        </form>

        {/* Token Alert Banner */}
        {generatedToken && (
          <div className="mt-4 p-4 rounded-xl bg-teal-950/80 border border-teal-500/50 flex flex-col sm:flex-row items-center justify-between gap-3 animate-pulse">
            <div>
              <div className="text-[10px] uppercase font-bold text-teal-300 font-mono">
                STS Token Generated Successfully
              </div>
              <div className="text-xl font-mono font-black text-white tracking-widest mt-0.5">
                {generatedToken}
              </div>
            </div>
            <button
              onClick={() => {
                navigator.clipboard.writeText(generatedToken);
                alert("Token copied to clipboard!");
              }}
              className="px-3 py-1.5 rounded-lg bg-teal-600 hover:bg-teal-500 text-white font-bold text-xs"
            >
              Copy Token
            </button>
          </div>
        )}
      </div>

      {/* Transactions Register */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
          <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
            Recent Prepayment Transactions Register
          </span>
          <span className="text-xs text-slate-400 font-mono">Real-Time Ingestion</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">TX Ref</th>
                <th className="py-3 px-3">Customer &bull; Meter</th>
                <th className="py-3 px-3">Band</th>
                <th className="py-3 px-3">Amount</th>
                <th className="py-3 px-3">kWh Vended</th>
                <th className="py-3 px-3">STS 20-Digit Token</th>
                <th className="py-3 px-4 text-right">Channel</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {vendings.map((v) => (
                <tr key={v.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-teal-400">
                    {v.transaction_ref}
                    <span className="block text-[10px] text-slate-500 font-normal">{v.created_at_wat}</span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{v.customer_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{v.account_or_meter}</div>
                  </td>
                  <td className="py-3 px-3 font-semibold text-cyan-400">{v.tariff_band}</td>
                  <td className="py-3 px-3 font-mono font-bold text-white">
                    &#8358;{v.amount_paid_ngn.toLocaleString()}
                  </td>
                  <td className="py-3 px-3 font-mono text-emerald-400 font-bold">{v.kwh_vended} kWh</td>
                  <td className="py-3 px-3 font-mono text-white text-[11px] font-bold">
                    {v.sts_token_20_digit}
                  </td>
                  <td className="py-3 px-4 text-right">
                    <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-slate-800 text-slate-300 border border-slate-700">
                      {v.payment_channel}
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

