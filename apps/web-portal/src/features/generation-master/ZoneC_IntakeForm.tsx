/**
 * Zone C: Main Operational Workspace (Right Pane - 35% Intake Form - Generation Master).
 * Phone Call Collation & Remote Plant Telemetry Entry Form.
 */

import React, { useState } from "react";
import { submitPlantTelemetry } from "./api";

export const ZoneC_IntakeForm: React.FC<{ onSuccess?: () => void }> = ({ onSuccess }) => {
  const [plantId, setPlantId] = useState("PLANT_01_KURRA_FALLS");
  const [mw, setMw] = useState(12.5);
  const [waterLevel, setWaterLevel] = useState(24.5);
  const [vibration, setVibration] = useState(2.1);
  const [statorTemp, setStatorTemp] = useState(65.0);
  const [lubeOil, setLubeOil] = useState(3.2);
  const [submitting, setSubmitting] = useState(false);
  const [msg, setMsg] = useState<{ text: string; type: "success" | "error" } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMsg(null);
    setSubmitting(true);

    try {
      await submitPlantTelemetry({
        plant_id: plantId,
        mw_generated: mw,
        water_level_meters: waterLevel,
        stator_temp_celsius: statorTemp,
        bearing_vibration_mm_s: vibration,
        lube_oil_pressure_bar: lubeOil,
        shift_roster_id: "ROSTER-MORNING-01",
      });

      setMsg({ text: "Plant Telemetry Collated & Verified!", type: "success" });
      if (onSuccess) onSuccess();
    } catch (err: any) {
      setMsg({ text: err.message || "Failed to submit telemetry", type: "error" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-3 border-b border-slate-700/60 mb-4">
        <h3 className="text-sm font-bold text-white tracking-wide">Hourly Call Collation Entry</h3>
        <p className="text-xs text-slate-400 mt-0.5">Direct Log from Remote Station Operator</p>
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

      <form onSubmit={handleSubmit} className="space-y-3.5 text-xs flex-1 overflow-y-auto">
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">Select Hydro / Thermal Plant</label>
          <select
            value={plantId}
            onChange={(e) => setPlantId(e.target.value)}
            className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
          >
            <option value="PLANT_01_KURRA_FALLS">Plant 1 - Kurra Falls Hydro</option>
            <option value="PLANT_02_JEKKO_1">Plant 2 - Jekko 1 Hydro</option>
            <option value="PLANT_03_JEKKO_2">Plant 3 - Jekko 2 Hydro</option>
            <option value="PLANT_04_ANKWIL">Plant 4 - Ankwil Hydro</option>
            <option value="PLANT_05_TENTI">Plant 5 - Tenti Hydro</option>
            <option value="PLANT_06_BUKURU_THERMAL">Plant 6 - Bukuru Thermal Backup</option>
          </select>
        </div>

        <div className="grid grid-cols-2 gap-2">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">Active MW Output</label>
            <input
              type="number"
              step="0.1"
              required
              value={mw}
              onChange={(e) => setMw(parseFloat(e.target.value))}
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
            />
          </div>
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">Water Level (m)</label>
            <input
              type="number"
              step="0.1"
              required
              value={waterLevel}
              onChange={(e) => setWaterLevel(parseFloat(e.target.value))}
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-2">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">Vibration (mm/s)</label>
            <input
              type="number"
              step="0.1"
              required
              value={vibration}
              onChange={(e) => setVibration(parseFloat(e.target.value))}
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
            />
          </div>
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">Stator Temp (°C)</label>
            <input
              type="number"
              step="0.1"
              required
              value={statorTemp}
              onChange={(e) => setStatorTemp(parseFloat(e.target.value))}
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
            />
          </div>
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-semibold rounded-lg text-xs transition-colors"
        >
          {submitting ? "Logging..." : "Log Station Telemetry"}
        </button>
      </form>
    </div>
  );
};

