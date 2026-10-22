/**
 * Zone C: Main Operational Workspace (Right Pane - 35% Data Intake Form).
 * Context-aware KYC Customer Key-In and Fault Dispatch Form.
 */

import React, { useState } from "react";
import { onboardCustomer, CustomerOnboardDTO } from "./api";

interface IntakeFormProps {
  staffId: string;
  onSuccess?: () => void;
}

export const ZoneC_IntakeForm: React.FC<IntakeFormProps> = ({ staffId, onSuccess }) => {
  const [form, setForm] = useState<Partial<CustomerOnboardDTO>>({
    tariff_band: "BAND_A",
    feeder_line_code: "FDR-33KV-VOM",
    dt_station_code: "TX-DT-500KVA-MARKET",
  });
  const [submitting, setSubmitting] = useState(false);
  const [msg, setMsg] = useState<{ text: string; type: "success" | "error" } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMsg(null);
    setSubmitting(true);

    try {
      await onboardCustomer({
        national_id_nin: form.national_id_nin || "NIN-2026-0001",
        full_name: form.full_name || "",
        phone_number: form.phone_number || "",
        account_number: form.account_number || `ACC-${Date.now().toString().slice(-6)}`,
        meter_number: form.meter_number || `MTR-${Date.now().toString().slice(-6)}`,
        feeder_line_code: form.feeder_line_code || "FDR-33KV-VOM",
        dt_station_code: form.dt_station_code || "TX-DT-500KVA",
        tariff_band: form.tariff_band || "BAND_A",
        service_address: form.service_address || "",
        latitude: form.latitude || 9.75,
        longitude: form.longitude || 8.85,
        created_by_staff_id: staffId,
      });

      setMsg({ text: "KYC Customer Onboarded & PostGIS Anchor Locked!", type: "success" });
      setForm({
        tariff_band: "BAND_A",
        feeder_line_code: "FDR-33KV-VOM",
        dt_station_code: "TX-DT-500KVA-MARKET",
      });
      if (onSuccess) onSuccess();
    } catch (err: any) {
      setMsg({ text: err.message || "Failed to onboard customer", type: "error" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-3 border-b border-slate-700/60 mb-4">
        <h3 className="text-sm font-bold text-white tracking-wide">KYC Customer Onboarding</h3>
        <p className="text-xs text-slate-400 mt-0.5">PostGIS DT Radius & Meter Anchor Form</p>
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

      <form onSubmit={handleSubmit} className="space-y-3.5 text-xs flex-1 overflow-y-auto pr-1">
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Customer Full Name</label>
          <input
            type="text"
            required
            value={form.full_name || ""}
            onChange={(e) => setForm({ ...form, full_name: e.target.value })}
            placeholder="e.g. Ibrahim Danladi"
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>

        <div className="grid grid-cols-2 gap-2">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">Phone Number</label>
            <input
              type="text"
              required
              value={form.phone_number || ""}
              onChange={(e) => setForm({ ...form, phone_number: e.target.value })}
              placeholder="08031234567"
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
            />
          </div>
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">National NIN</label>
            <input
              type="text"
              required
              value={form.national_id_nin || ""}
              onChange={(e) => setForm({ ...form, national_id_nin: e.target.value })}
              placeholder="12345678901"
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-2">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">Feeder Code</label>
            <input
              type="text"
              required
              value={form.feeder_line_code || ""}
              onChange={(e) => setForm({ ...form, feeder_line_code: e.target.value })}
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
            />
          </div>
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">Tariff Band</label>
            <select
              value={form.tariff_band || "BAND_A"}
              onChange={(e) => setForm({ ...form, tariff_band: e.target.value })}
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
            >
              <option value="BAND_A">Band A (20+ hrs)</option>
              <option value="BAND_B">Band B (16-20 hrs)</option>
              <option value="BAND_C">Band C (12-16 hrs)</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Service Address</label>
          <input
            type="text"
            required
            value={form.service_address || ""}
            onChange={(e) => setForm({ ...form, service_address: e.target.value })}
            placeholder="No 14 Commercial Ave, Bukuru, Jos"
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full mt-2 py-2.5 px-4 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-semibold rounded-lg text-xs transition-colors"
        >
          {submitting ? "Anchoring PostGIS..." : "Submit KYC Onboarding"}
        </button>
      </form>
    </div>
  );
};

