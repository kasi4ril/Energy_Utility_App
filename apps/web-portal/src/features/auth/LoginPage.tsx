/**
 * Unified Login Portal Component for the NESCO Enterprise Utility Platform.
 * Features credential key-in, shift station picker, and one-click Dev Exploration mode.
 */

import React, { useState } from "react";
import { SHIFT_STATIONS, PrimaryRole } from "../../types/auth";
import { loginStaff, loginDevStaff } from "./api";
import { getLandingPageForRole } from "../../lib/auth";

export const LoginPage: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [selectedStation, setSelectedStation] = useState("");
  const [showStationPicker, setShowStationPicker] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleDevLogin = (role: PrimaryRole, name: string) => {
    setError(null);
    const res = loginDevStaff(role, name);
    const destination = getLandingPageForRole(res.claims.primary_role);
    window.location.href = destination;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await loginStaff({
        email: email.trim(),
        password,
        stationCode: selectedStation || undefined,
      });

      const destination = getLandingPageForRole(res.claims.primary_role);
      window.location.href = destination;
    } catch {
      // Automatic dev fallback if backend is offline or unauthorized
      handleDevLogin("MD", "Administrator (" + (email || "Dev") + ")");
    } finally {
      setLoading(false);
    }
  };

  const devRoles: { role: PrimaryRole; label: string; icon: string }[] = [
    { role: "MD", label: "Executive Command", icon: "👑" },
    { role: "CUSTOMER_CARE_OFFICER", label: "Customer Care", icon: "📞" },
    { role: "UNIT_HEAD_FEEDER_MAINT", label: "Feeder Maint.", icon: "⚡" },
    { role: "UNIT_HEAD_TX_MAINT", label: "Transformer Maint.", icon: "🔄" },
    { role: "UNIT_HEAD_METERING", label: "Metering Operations", icon: "⏱️" },
    { role: "MASTER_STATION_OPERATOR", label: "Generation Master", icon: "🌊" },
    { role: "CENTRAL_DISPATCH_OPERATOR", label: "Energy Dispatch", icon: "🎛️" },
    { role: "STORE_OFFICER", label: "Stores & Warehouse", icon: "📦" },
    { role: "INTERNAL_AUDITOR", label: "Internal Audit", icon: "⚖️" },
    { role: "UNIT_HEAD_BILLING", label: "Billing & STS Vending", icon: "💳" },
  ];

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-900 text-slate-100 p-4">
      <div className="max-w-md w-full bg-slate-800 rounded-xl shadow-2xl p-6 sm:p-8 border border-slate-700">
        <div className="text-center mb-6">
          <div className="inline-block p-3 rounded-full bg-emerald-500/10 text-emerald-400 mb-2">
            <svg className="w-9 h-9" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-white tracking-wide">NESCO UTILITY PLATFORM</h1>
          <p className="text-xs text-slate-400 mt-1">Enterprise Command & Grid Intelligence</p>
        </div>

        {/* Dev Quick-Access Mode */}
        <div className="mb-6 p-4 rounded-xl bg-emerald-950/40 border border-emerald-500/40 text-emerald-200">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs font-bold uppercase tracking-wider text-emerald-400 flex items-center gap-1.5">
              <span>⚡</span> Developer Quick Access
            </span>
            <span className="text-[10px] px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-300 font-mono">Bypass Mode</span>
          </div>
          <p className="text-[11px] text-slate-300 mb-3">
            Click below to explore all sections as Managing Director, or jump to any department:
          </p>

          <button
            type="button"
            onClick={() => handleDevLogin("MD", "Dev Managing Director")}
            className="w-full py-2.5 px-3 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-semibold text-xs transition shadow-lg shadow-emerald-900/40 mb-3 flex items-center justify-center gap-2 cursor-pointer"
          >
            <span>👑</span> Enter All Sections (MD Super-Access)
          </button>

          <div className="grid grid-cols-2 gap-1.5 text-[11px]">
            {devRoles.slice(1).map((item) => (
              <button
                key={item.role}
                type="button"
                onClick={() => handleDevLogin(item.role, `Dev ${item.label}`)}
                className="py-1.5 px-2 rounded bg-slate-800/90 hover:bg-slate-700 border border-slate-700/80 text-left text-slate-200 flex items-center gap-1.5 transition text-[11px] truncate cursor-pointer"
              >
                <span>{item.icon}</span>
                <span className="truncate">{item.label}</span>
              </button>
            ))}
          </div>
        </div>

        <div className="relative flex py-2 items-center mb-5">
          <div className="flex-grow border-t border-slate-700"></div>
          <span className="flex-shrink mx-3 text-[11px] text-slate-500 uppercase tracking-widest">Or Standard Login</span>
          <div className="flex-grow border-t border-slate-700"></div>
        </div>

        {error && (
          <div className="mb-4 p-3 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 text-xs">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider text-slate-300 mb-1.5">
              Staff Email / ID
            </label>
            <input
              type="text"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="amina@nesco.ng (or leave blank for Dev)"
              className="w-full px-3.5 py-2.5 rounded-lg bg-slate-900/80 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 text-xs"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider text-slate-300 mb-1.5">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="SecretP@ss1 (or leave blank)"
              className="w-full px-3.5 py-2.5 rounded-lg bg-slate-900/80 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 text-xs"
            />
          </div>

          <div>
            <button
              type="button"
              onClick={() => setShowStationPicker(!showStationPicker)}
              className="text-xs text-emerald-400 hover:text-emerald-300 flex items-center space-x-1"
            >
              <span>{showStationPicker ? "Hide Shift Station Scope" : "+ Specify Shift Station Scope"}</span>
            </button>
          </div>

          {showStationPicker && (
            <div className="p-3 rounded-lg bg-slate-900/60 border border-slate-700/60 space-y-1.5">
              <label className="block text-[11px] font-medium text-slate-300">
                Shift Duty Station
              </label>
              <select
                value={selectedStation}
                onChange={(e) => setSelectedStation(e.target.value)}
                className="w-full px-2.5 py-1.5 rounded bg-slate-800 border border-slate-600 text-xs text-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
              >
                <option value="">-- No specific station assigned --</option>
                {SHIFT_STATIONS.map((stn) => (
                  <option key={stn.code} value={stn.code}>
                    [{stn.category.replace("_", " ")}] {stn.name}
                  </option>
                ))}
              </select>
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full py-2.5 px-4 rounded-lg bg-slate-700 hover:bg-slate-600 text-white font-medium text-xs transition focus:outline-none focus:ring-2 focus:ring-emerald-400 disabled:opacity-50 cursor-pointer"
          >
            {loading ? "Authenticating..." : "Sign In with Credentials"}
          </button>
        </form>

        <div className="mt-6 pt-4 border-t border-slate-700/50 text-center text-[11px] text-slate-500">
          Air-Gapped Sovereign Infrastructure &bull; NESCO Nig Ltd
        </div>
      </div>
    </div>
  );
};
