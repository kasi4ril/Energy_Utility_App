/**
 * View 2: New Customer Registration & Metering KYC (/portal/customer-care/register).
 * Captures Full Name, Phone (+234), NIN/CAC, GPS, Building & Tariff Class, Load kVA,
 * nearest Feeder & DT, NERC Band A-E, MAP Scheme, and stamps CC Staff Attribution.
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import { onboardCustomer, verifyNIN } from "./api";
import { TariffBand, BuildingClass, TariffClass, NINVerificationResult } from "./types";
import { UserCheck, MapPin, Shield, CheckCircle, AlertCircle, CheckCheck, Loader2 } from "lucide-react";

interface View2Props {
  staffId: string;
  staffName?: string;
  onCustomerOnboarded?: () => void;
}

const DEFAULT_FORM = {
  full_name: "", phone: "+234", email: "", nin: "",
  building_class: "RESIDENTIAL" as BuildingClass,
  tariff_class: "R2 Residential" as TariffClass,
  estimated_load_kva: 5.0, landmark: "", service_address: "",
  transformer_id: "TX-DT-500KVA-MARKET", feeder_id: "FDR-33KV-VOM", feeder_line_code: "FDR-33KV-VOM",
  map_scheme_type: "Single-Phase MAP", meter_serial: "", meter_type: "PREPAID_STS" as const,
  tariff_band: "BAND_A" as TariffBand, latitude: 9.7542, longitude: 8.8521,
};

export function View2_NewCustomerKYC({
  staffId,
  staffName = "Amina Yusuf (CC Officer)",
  onCustomerOnboarded,
}: View2Props) {
  const [formData, setFormData] = useState(DEFAULT_FORM);
  const [isVerifyingNIN, setIsVerifyingNIN] = useState(false);
  const [ninResult, setNinResult] = useState<NINVerificationResult | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const triggerNINLookup = async (ninValue: string) => {
    const trimmed = ninValue.trim();
    if (!trimmed || trimmed.length < 5) return;
    setIsVerifyingNIN(true);
    try {
      const res = await verifyNIN(trimmed);
      setNinResult(res);
      if (res?.verified_name) setFormData((prev) => ({ ...prev, full_name: res.verified_name }));
    } catch {} finally {
      setIsVerifyingNIN(false);
    }
  };

  const handleNINInputChange = (val: string) => {
    setFormData((prev) => ({ ...prev, nin: val }));
    if (val.trim().length === 11) triggerNINLookup(val.trim());
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSuccessMessage(null);
    setErrorMessage(null);
    try {
      const payload = { ...formData, created_by_staff_id: staffId, registered_by_cc_staff_name: staffName };
      const res = await onboardCustomer(payload);
      setSuccessMessage(
        `Customer successfully onboarded! Account: ${res.account_no} (CIN: ${res.cin}). Site survey ticket routed to Metering Unit Head.`
      );
      if (onCustomerOnboarded) onCustomerOnboarded();
      setFormData(DEFAULT_FORM);
      setNinResult(null);
    } catch (err: any) {
      setErrorMessage(err.message || "Failed to onboard customer");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-6 max-w-4xl">
      <div className="flex items-center space-x-3 mb-6 pb-4 border-b border-slate-800">
        <div className="p-2.5 rounded-lg bg-blue-600/20 text-blue-400 border border-blue-500/30">
          <UserCheck className="w-6 h-6" />
        </div>
        <div>
          <h2 className="text-base font-bold text-slate-100 uppercase tracking-wide">
            New Customer Registration &amp; Metering KYC Desk
          </h2>
          <p className="text-xs text-slate-400">
            Anti-Ghost GPS lock, NERC Tariff Band assignment &amp; Metering Scheme survey dispatch
          </p>
        </div>
      </div>

      {successMessage && (
        <div className="p-3.5 mb-5 rounded-lg bg-emerald-950/60 border border-emerald-500/40 text-emerald-300 text-xs flex items-center space-x-2.5">
          <CheckCircle className="w-4 h-4 shrink-0 text-emerald-400" />
          <span>{successMessage}</span>
        </div>
      )}

      {errorMessage && (
        <div className="p-3.5 mb-5 rounded-lg bg-red-950/60 border border-red-500/40 text-red-300 text-xs flex items-center space-x-2.5">
          <AlertCircle className="w-4 h-4 shrink-0 text-red-400" />
          <span>{errorMessage}</span>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-5 text-xs">
        {/* Row 1: Consumer Identity (NIN Verification & Name Auto-Population) */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <div className="flex items-center justify-between mb-1">
              <label className="text-slate-300 font-medium">National ID / NIN (11 Digits) *</label>
              {isVerifyingNIN && (
                <span className="text-[10px] text-blue-400 flex items-center space-x-1">
                  <Loader2 className="w-3 h-3 animate-spin" />
                  <span>NIMC Lookup...</span>
                </span>
              )}
            </div>
            <div className="relative">
              <input
                type="text"
                required
                maxLength={11}
                placeholder="11-digit NIN"
                value={formData.nin}
                onChange={(e) => handleNINInputChange(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 font-mono focus:border-blue-500 focus:outline-none pr-14"
              />
              <button
                type="button"
                onClick={() => triggerNINLookup(formData.nin)}
                className="absolute right-1 top-1 bottom-1 px-2 text-[10px] bg-blue-600/40 hover:bg-blue-600/60 text-blue-200 rounded border border-blue-500/40"
              >
                Verify
              </button>
            </div>
            {ninResult && (
              <div className="text-[10px] mt-1 text-emerald-400 flex items-center space-x-1 font-mono">
                <CheckCheck className="w-3 h-3 text-emerald-400" />
                <span>NIMC Verified: {ninResult.verified_name}</span>
              </div>
            )}
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">
              Full Customer Name *{" "}
              {ninResult && <span className="text-[10px] text-emerald-400 font-mono">(Verified)</span>}
            </label>
            <input
              type="text"
              required
              placeholder="e.g. Hajiya Zainab Mohammed"
              value={formData.full_name}
              onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
              className={`w-full bg-slate-950 border rounded-lg px-3 py-2 text-slate-100 focus:outline-none ${
                ninResult ? "border-emerald-500/50 bg-emerald-950/10" : "border-slate-800 focus:border-blue-500"
              }`}
            />
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">Phone Number (+234) *</label>
            <input
              type="text"
              required
              placeholder="+2348031234567"
              value={formData.phone}
              onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            />
          </div>
        </div>

        {/* Row 2: Tariff Classification & Load Assessment */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-slate-300 font-medium mb-1">Customer Tariff Class *</label>
            <select
              value={formData.tariff_class}
              onChange={(e) => setFormData({ ...formData, tariff_class: e.target.value as TariffClass })}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            >
              <option value="R2 Residential">R2 Residential (Domestic Single/Three-Phase)</option>
              <option value="C1 Commercial">C1 Commercial (Small Business/Shops)</option>
              <option value="D1 Industrial">D1 Industrial (Manufacturing/Workshops)</option>
              <option value="MD Maximum Demand">MD Maximum Demand (Large Commercial &gt;45kVA)</option>
            </select>
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">Estimated Load Rating (kVA) *</label>
            <input
              type="number"
              step="0.5"
              min="0.5"
              required
              value={formData.estimated_load_kva}
              onChange={(e) => setFormData({ ...formData, estimated_load_kva: parseFloat(e.target.value) || 5.0 })}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none font-mono"
            />
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">MAP / NMMP Scheme *</label>
            <select
              value={formData.map_scheme_type}
              onChange={(e) => setFormData({ ...formData, map_scheme_type: e.target.value })}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            >
              <option value="Single-Phase MAP">Single-Phase MAP (₦82,000 Upfront/DisCo Financed)</option>
              <option value="Three-Phase MAP">Three-Phase MAP (₦145,000 Industrial)</option>
              <option value="NMMP Phase-2 Grant">NMMP Phase-2 Federal Grant Rollout</option>
            </select>
          </div>
        </div>

        {/* Row 3: Grid Attachment & NERC Tariff Band */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-slate-300 font-medium mb-1">11kV / 33kV Feeder Line *</label>
            <select
              value={formData.feeder_id}
              onChange={(e) => setFormData({ ...formData, feeder_id: e.target.value, feeder_line_code: e.target.value })}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            >
              <option value="FDR-33KV-VOM">FDR-33KV-VOM (Vom Band A)</option>
              <option value="FDR-11KV-RAYFIELD">FDR-11KV-RAYFIELD (Rayfield Band A)</option>
              <option value="FDR-33KV-ANGLO">FDR-33KV-ANGLO (Anglo-Jos Band A)</option>
              <option value="FDR-11KV-KOWA">FDR-11KV-KOWA (Dadin Kowa Band B)</option>
            </select>
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">Nearest DT *</label>
            <select
              value={formData.transformer_id}
              onChange={(e) => setFormData({ ...formData, transformer_id: e.target.value })}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            >
              <option value="TX-DT-500KVA-MARKET">TX-DT-500KVA-MARKET (Bukuru Central)</option>
              <option value="TX-DT-300KVA-RAYFIELD">TX-DT-300KVA-RAYFIELD (Old Airport Rd)</option>
              <option value="TX-DT-1000KVA-ANGLO">TX-DT-1000KVA-ANGLO (Steel Mills)</option>
              <option value="TX-DT-500KVA-KOWA">TX-DT-500KVA-KOWA (Gindiri Sub)</option>
            </select>
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">Service Tariff Band (NERC)</label>
            <div className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-emerald-400 font-bold flex items-center justify-between">
              <span>{formData.tariff_band} (&ge; 20 Hrs)</span>
              <span className="text-[10px] text-slate-400 font-normal">Auto-Assigned</span>
            </div>
          </div>
        </div>

        {/* Row 4: Landmark, Physical Address & PostGIS GPS */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-slate-300 font-medium mb-1">Major Landmark *</label>
            <input
              type="text"
              required
              placeholder="e.g. Opposite Central Mosque / Market Gate"
              value={formData.landmark}
              onChange={(e) => setFormData({ ...formData, landmark: e.target.value })}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">Premises Physical Address *</label>
            <input
              type="text"
              required
              placeholder="e.g. Plot 18 Commercial Street, Vom Road"
              value={formData.service_address}
              onChange={(e) => setFormData({ ...formData, service_address: e.target.value })}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">Premises GPS (PostGIS Point) *</label>
            <div className="flex items-center space-x-2">
              <input
                type="number"
                step="0.0001"
                required
                value={formData.latitude}
                onChange={(e) => setFormData({ ...formData, latitude: parseFloat(e.target.value) || 9.7542 })}
                className="w-1/2 bg-slate-950 border border-slate-800 rounded-lg px-2 py-2 text-slate-100 font-mono text-[11px]"
                title="Latitude"
              />
              <input
                type="number"
                step="0.0001"
                required
                value={formData.longitude}
                onChange={(e) => setFormData({ ...formData, longitude: parseFloat(e.target.value) || 8.8521 })}
                className="w-1/2 bg-slate-950 border border-slate-800 rounded-lg px-2 py-2 text-slate-100 font-mono text-[11px]"
                title="Longitude"
              />
            </div>
          </div>
        </div>

        {/* Staff Attribution Footer Stamp */}
        <div className="p-3 rounded-lg bg-slate-950/70 border border-slate-800 flex items-center justify-between text-slate-400">
          <div className="flex items-center space-x-2">
            <Shield className="w-4 h-4 text-blue-400" />
            <span>
              Attribution Stamp: <strong className="text-slate-200">{staffName}</strong> ({staffId})
            </span>
          </div>
          <span className="text-[10px] font-mono text-emerald-400">Cryptographic HMAC Lock</span>
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          disabled={isSubmitting}
          className="w-full bg-blue-600 hover:bg-blue-500 text-white font-semibold py-2.5 px-4 rounded-lg transition-colors flex items-center justify-center space-x-2 shadow-md disabled:opacity-50"
        >
          {isSubmitting ? (
            <span>Registering &amp; Verifying Proximity...</span>
          ) : (
            <>
              <UserCheck className="w-4 h-4" />
              <span>Register Customer &amp; Route Site Survey to Metering Unit</span>
            </>
          )}
        </button>
      </form>
    </div>
  );
}

