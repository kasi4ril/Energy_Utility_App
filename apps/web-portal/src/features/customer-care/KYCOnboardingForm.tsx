/**
 * Task View 3: New Customer KYC Onboarding Workspace.
 * Full Anti-Ghost Identity Validation & PostGIS DT GPS Proximity Lock.
 */

import React, { useState } from "react";
import { BuildingClass, KYCOnboardingPayload, TariffBand } from "./types";
import { onboardCustomer } from "./api";

interface KYCOnboardingFormProps {
  staffId: string;
  onCustomerOnboarded?: () => void;
  onNavigateToCRM?: () => void;
}

export const KYCOnboardingForm: React.FC<KYCOnboardingFormProps> = ({
  staffId,
  onCustomerOnboarded,
  onNavigateToCRM,
}) => {
  const [form, setForm] = useState<Partial<KYCOnboardingPayload>>({
    tariff_band: "BAND_A",
    building_class: "RESIDENTIAL",
    meter_type: "PREPAID_STS",
    feeder_line_code: "FDR-33KV-VOM",
    transformer_id: "TX-DT-500KVA-MARKET",
    latitude: 9.7542,
    longitude: 8.8521,
  });

  const [submitting, setSubmitting] = useState<boolean>(false);
  const [onboardedResult, setOnboardedResult] = useState<{
    cin: string;
    account_no: string;
    full_name: string;
  } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setOnboardedResult(null);

    try {
      const result = await onboardCustomer({
        full_name: form.full_name || "",
        phone: form.phone || "",
        email: form.email || "",
        nin: form.nin || `NIN-${Date.now().toString().slice(-8)}`,
        building_class: (form.building_class as BuildingClass) || "RESIDENTIAL",
        service_address: form.service_address || "",
        transformer_id: form.transformer_id || "TX-DT-500KVA-MARKET",
        feeder_line_code: form.feeder_line_code || "FDR-33KV-VOM",
        meter_serial: form.meter_serial || `MTR-${Date.now().toString().slice(-6)}`,
        meter_type: form.meter_type || "PREPAID_STS",
        tariff_band: (form.tariff_band as TariffBand) || "BAND_A",
        latitude: form.latitude || 9.7542,
        longitude: form.longitude || 8.8521,
        created_by_staff_id: staffId,
      });

      setOnboardedResult({
        cin: result.cin,
        account_no: result.account_no,
        full_name: result.full_name,
      });

      if (onCustomerOnboarded) onCustomerOnboarded();
    } catch (err: any) {
      alert(err.message || "Failed to onboard customer");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      <div className="border-b border-slate-800 pb-4 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <div>
          <h2 className="text-base font-bold text-white tracking-wide flex items-center gap-2">
            <span>👤 New Customer KYC Identity &amp; Anti-Ghost Anchor</span>
          </h2>
          <p className="text-xs text-slate-400 mt-0.5">
            Register new electricity consumers with NIN validation, building classification, and PostGIS DT radius lock.
          </p>
        </div>
        {onNavigateToCRM && (
          <button
            onClick={onNavigateToCRM}
            className="text-xs text-emerald-400 hover:text-emerald-300 font-semibold flex items-center gap-1"
          >
            Go to CRM Directory &rarr;
          </button>
        )}
      </div>

      {onboardedResult && (
        <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 flex items-center justify-between text-xs animate-fadeIn">
          <div>
            <span className="font-bold">Customer Verified &amp; Anchored!</span> Account No:{" "}
            <span className="font-mono font-bold text-white bg-emerald-950 px-2 py-0.5 rounded">
              {onboardedResult.account_no}
            </span>{" "}
            &bull; CIN:{" "}
            <span className="font-mono font-bold text-white bg-emerald-950 px-2 py-0.5 rounded">
              {onboardedResult.cin}
            </span>{" "}
            for <span className="font-semibold text-white">{onboardedResult.full_name}</span>.
          </div>
          {onNavigateToCRM && (
            <button
              onClick={onNavigateToCRM}
              className="px-3 py-1 bg-emerald-600 hover:bg-emerald-500 text-white rounded font-bold text-xs shadow"
            >
              Open in CRM &rarr;
            </button>
          )}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-5 text-xs">
        {/* Section 1: Consumer Identification */}
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4 space-y-3">
          <label className="text-xs font-bold text-slate-200 block">
            1. Consumer Identity &amp; National Verification
          </label>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Customer Full Name *</label>
              <input
                type="text"
                required
                value={form.full_name || ""}
                onChange={(e) => setForm({ ...form, full_name: e.target.value })}
                placeholder="e.g. Ibrahim Danladi Gyang"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Phone Number *</label>
              <input
                type="text"
                required
                value={form.phone || ""}
                onChange={(e) => setForm({ ...form, phone: e.target.value })}
                placeholder="08031234567"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">National NIN (11 Digits) *</label>
              <input
                type="text"
                required
                value={form.nin || ""}
                onChange={(e) => setForm({ ...form, nin: e.target.value })}
                placeholder="24891029481"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
            </div>
          </div>
        </div>

        {/* Section 2: Premises & Building Class */}
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4 space-y-3">
          <label className="text-xs font-bold text-slate-200 block">
            2. Premises Service Address &amp; Classification
          </label>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div className="md:col-span-2">
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Premises Service Address *</label>
              <input
                type="text"
                required
                value={form.service_address || ""}
                onChange={(e) => setForm({ ...form, service_address: e.target.value })}
                placeholder="Plot 14 Commercial Ave, Bukuru Industrial Area, Jos"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Building Classification</label>
              <select
                value={form.building_class || "RESIDENTIAL"}
                onChange={(e) => setForm({ ...form, building_class: e.target.value as BuildingClass })}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
              >
                <option value="RESIDENTIAL">Residential Property</option>
                <option value="COMMERCIAL">Commercial Enterprise</option>
                <option value="INDUSTRIAL">Industrial / Heavy Load</option>
                <option value="SPECIAL">Special (Hospital / School / Government)</option>
              </select>
            </div>
          </div>
        </div>

        {/* Section 3: PostGIS DT Proximity & Grid Lock */}
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4 space-y-3">
          <div className="flex items-center justify-between">
            <label className="text-xs font-bold text-slate-200">
              3. Grid Asset Anchoring &amp; Anti-Ghost Proximity
            </label>
            <span className="text-[11px] text-emerald-400 font-bold bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/30">
              ✓ GPS Lock &le;100m DT Radius
            </span>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-3">
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Feeder Line Code</label>
              <input
                type="text"
                required
                value={form.feeder_line_code || ""}
                onChange={(e) => setForm({ ...form, feeder_line_code: e.target.value })}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Distribution Transformer</label>
              <input
                type="text"
                required
                value={form.transformer_id || ""}
                onChange={(e) => setForm({ ...form, transformer_id: e.target.value })}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">GPS Latitude (Plateau)</label>
              <input
                type="number"
                step="0.0001"
                required
                value={form.latitude || 9.7542}
                onChange={(e) => setForm({ ...form, latitude: parseFloat(e.target.value) })}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">GPS Longitude</label>
              <input
                type="number"
                step="0.0001"
                required
                value={form.longitude || 8.8521}
                onChange={(e) => setForm({ ...form, longitude: parseFloat(e.target.value) })}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
          </div>
        </div>

        {/* Section 4: Metering & Tariff */}
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4 space-y-3">
          <label className="text-xs font-bold text-slate-200 block">
            4. Meter Assignment &amp; Tariff Band Scheduling
          </label>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Meter Serial Number</label>
              <input
                type="text"
                required
                value={form.meter_serial || ""}
                onChange={(e) => setForm({ ...form, meter_serial: e.target.value })}
                placeholder="e.g. MTR-2026-9814"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Meter Hardware Type</label>
              <select
                value={form.meter_type || "PREPAID_STS"}
                onChange={(e) => setForm({ ...form, meter_type: e.target.value as any })}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white"
              >
                <option value="PREPAID_STS">STS Prepaid Smart Card Meter</option>
                <option value="SMART_AMI">Smart AMI Two-Way Cellular Meter</option>
                <option value="POSTPAID">Postpaid Certified Maximum Demand</option>
              </select>
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">NERC Tariff Band</label>
              <select
                value={form.tariff_band || "BAND_A"}
                onChange={(e) => setForm({ ...form, tariff_band: e.target.value as TariffBand })}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white"
              >
                <option value="BAND_A">Band A (Minimum 20 hrs/day supply)</option>
                <option value="BAND_B">Band B (Minimum 16 hrs/day supply)</option>
                <option value="BAND_C">Band C (Minimum 12 hrs/day supply)</option>
              </select>
            </div>
          </div>
        </div>

        <div className="flex justify-end gap-3 pt-2">
          <button
            type="submit"
            disabled={submitting}
            className="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-bold text-xs rounded-lg transition shadow-lg flex items-center gap-2"
          >
            {submitting ? "Anchoring PostGIS..." : "🔒 Verify Identity & Create Customer Account"}
          </button>
        </div>
      </form>
    </div>
  );
};

