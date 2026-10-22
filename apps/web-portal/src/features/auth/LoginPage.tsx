/**
 * Unified Login Portal Component for the NESCO Enterprise Utility Platform.
 * Features credential key-in and dynamic shift station selection.
 */

import React, { useState } from "react";
import { SHIFT_STATIONS } from "../../types/auth";
import { loginStaff } from "./api";
import { getLandingPageForRole } from "../../lib/auth";

export const LoginPage: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [selectedStation, setSelectedStation] = useState("");
  const [showStationPicker, setShowStationPicker] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

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
    } catch (err: any) {
      setError(err.message || "Authentication failed. Check credentials.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-900 text-slate-100 p-4">
      <div className="max-w-md w-full bg-slate-800 rounded-xl shadow-2xl p-8 border border-slate-700">
        <div className="text-center mb-8">
          <div className="inline-block p-3 rounded-full bg-emerald-500/10 text-emerald-400 mb-3">
            <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-white tracking-wide">NESCO UTILITY PLATFORM</h1>
          <p className="text-sm text-slate-400 mt-1">Enterprise Command & Grid Intelligence</p>
        </div>

        {error && (
          <div className="mb-6 p-4 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider text-slate-300 mb-2">
              Staff Number / Official Email
            </label>
            <input
              type="text"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="e.g. amina@nesco.ng or NES-STF-001"
              className="w-full px-4 py-3 rounded-lg bg-slate-900/80 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 text-sm"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold uppercase tracking-wider text-slate-300 mb-2">
              Password
            </label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••••••"
              className="w-full px-4 py-3 rounded-lg bg-slate-900/80 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 text-sm"
            />
          </div>

          <div className="pt-1">
            <button
              type="button"
              onClick={() => setShowStationPicker(!showStationPicker)}
              className="text-xs text-emerald-400 hover:text-emerald-300 flex items-center space-x-1"
            >
              <span>{showStationPicker ? "Hide Shift Station Scope" : "+ Specify Shift Station Scope"}</span>
            </button>
          </div>

          {showStationPicker && (
            <div className="p-3 rounded-lg bg-slate-900/60 border border-slate-700/60 space-y-2">
              <label className="block text-xs font-medium text-slate-300">
                Shift Duty Station (Plants 1–6 / Central Dispatch / 33kV Substations)
              </label>
              <select
                value={selectedStation}
                onChange={(e) => setSelectedStation(e.target.value)}
                className="w-full px-3 py-2 rounded bg-slate-800 border border-slate-600 text-xs text-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
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
            className="w-full py-3 px-4 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-medium text-sm transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-emerald-400 focus:ring-offset-2 focus:ring-offset-slate-800 disabled:opacity-50"
          >
            {loading ? "Verifying Credentials..." : "Access Enterprise Portal"}
          </button>
        </form>

        <div className="mt-8 pt-6 border-t border-slate-700/50 text-center text-xs text-slate-500">
          Air-Gapped Sovereign Infrastructure &bull; NESCO Nig Ltd
        </div>
      </div>
    </div>
  );
};

