/**
 * Task View 2: Hourly Call Collation & Shift Handover Desk (Google Form A Bridge).
 * Radio/phone roll-call intake, remote personnel muster, and inter-shift handover notes.
 */

import React, { useState } from "react";
import { HydroPlantSummary, PlantID, ShiftRollcallRecord, ShiftWindow } from "./types";
import { submitPlantTelemetry, submitShiftRollcall } from "./api";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";

interface ShiftCollationDeskProps {
  plants: HydroPlantSummary[];
  rollcalls: ShiftRollcallRecord[];
  initialPlantId?: string;
  onSuccess?: () => void;
}

export const ShiftCollationDesk: React.FC<ShiftCollationDeskProps> = ({
  plants,
  rollcalls,
  initialPlantId,
  onSuccess,
}) => {
  const [selectedPlantId, setSelectedPlantId] = useState<PlantID>(
    (initialPlantId as PlantID) || "PLANT_01_KURRA_FALLS"
  );
  const [shiftWindow, setShiftWindow] = useState<ShiftWindow>("MORNING");
  const [mwGenerated, setMwGenerated] = useState<number>(14.5);
  const [waterLevelM, setWaterLevelM] = useState<number>(24.2);
  const [bearingVibration, setBearingVibration] = useState<number>(2.1);
  const [statorTempC, setStatorTempC] = useState<number>(68.0);
  const [lubeOilBar, setLubeOilBar] = useState<number>(3.2);
  const [operatorNames, setOperatorNames] = useState<string>("Ibrahim Pam, Yusuf Chuwang");
  const [protectionEng, setProtectionEng] = useState<string>("Engr. D. Gyang");
  const [technicians, setTechnicians] = useState<string>("Bulus Mark, Sunday Azi");
  const [handoverNotes, setHandoverNotes] = useState<string>("");
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [msg, setMsg] = useState<{ text: string; type: "success" | "error" } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setMsg(null);

    try {
      await Promise.all([
        submitPlantTelemetry({
          plant_id: selectedPlantId,
          mw_generated: mwGenerated,
          water_level_meters: waterLevelM,
          stator_temp_celsius: statorTempC,
          bearing_vibration_mm_s: bearingVibration,
          lube_oil_pressure_bar: lubeOilBar,
          shift_roster_id: `ROSTER-${shiftWindow}-01`,
        }),
        submitShiftRollcall({
          plant_id: selectedPlantId,
          shift_window: shiftWindow,
          operator_names: operatorNames.split(",").map((s) => s.trim()),
          protection_engineers: protectionEng.split(",").map((s) => s.trim()),
          technicians: technicians.split(",").map((s) => s.trim()),
          handover_notes: handoverNotes,
          synced_to_hq_hr: true,
          master_operator_id: "STAFF-GEN-OP-01",
        }),
      ]);

      setMsg({
        text: `Hourly Telemetry & Shift Muster Collated for ${selectedPlantId}! Synced to HQ HR & Core API.`,
        type: "success",
      });
      setHandoverNotes("");
      if (onSuccess) onSuccess();
    } catch (err: any) {
      setMsg({ text: err.message || "Failed to submit collation log", type: "error" });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      <div className="border-b border-slate-800 pb-4 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <div>
          <h2 className="text-base font-bold text-white tracking-wide flex items-center gap-2">
            <span>📻 Master Station Voice &amp; Radio Call Collation Desk</span>
          </h2>
          <p className="text-xs text-slate-400 mt-0.5">
            Hourly remote station check-in, hydrological/electrical telemetry ingestion, and shift personnel muster.
          </p>
        </div>
        <span className="px-2.5 py-1 rounded bg-slate-800 text-cyan-400 border border-cyan-500/30 text-xs font-semibold">
          Google Form A Webhook Bridge
        </span>
      </div>

      {msg && (
        <div
          className={`p-3 rounded-lg text-xs ${
            msg.type === "success"
              ? "bg-emerald-500/10 border border-emerald-500/30 text-emerald-400"
              : "bg-rose-500/10 border border-rose-500/30 text-rose-400"
          }`}
        >
          {msg.text}
        </div>
      )}

      {/* Main Collation Form */}
      <form onSubmit={handleSubmit} className="space-y-5 text-xs">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              1. Generating Station
            </label>
            <select
              value={selectedPlantId}
              onChange={(e) => setSelectedPlantId(e.target.value as PlantID)}
              className="w-full px-3 py-2 rounded bg-slate-950 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-cyan-500"
            >
              {plants.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name} ({p.location})
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              2. Active Shift Window
            </label>
            <select
              value={shiftWindow}
              onChange={(e) => setShiftWindow(e.target.value as ShiftWindow)}
              className="w-full px-3 py-2 rounded bg-slate-950 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-cyan-500"
            >
              <option value="MORNING">Morning Shift (08:00 – 16:00)</option>
              <option value="AFTERNOON">Afternoon Shift (16:00 – 00:00)</option>
              <option value="NIGHT">Night Shift (00:00 – 08:00)</option>
            </select>
          </div>

          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              3. Master Station Operator ID
            </label>
            <input
              type="text"
              readOnly
              value="STAFF-GEN-OP-01 (Master Desk)"
              className="w-full px-3 py-2 rounded bg-slate-950/60 border border-slate-800 text-slate-400 font-mono"
            />
          </div>
        </div>

        {/* Section 2: Telemetry Data */}
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4 space-y-3">
          <label className="text-xs font-bold text-slate-200 block">
            4. Remote Radio Telemetry (Hourly Plant Readings)
          </label>

          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3">
            <div>
              <label className="block text-[10px] text-slate-400 mb-1 font-semibold">Active MW</label>
              <input
                type="number"
                step="0.1"
                required
                value={mwGenerated}
                onChange={(e) => setMwGenerated(parseFloat(e.target.value))}
                className="w-full px-2.5 py-1.5 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[10px] text-slate-400 mb-1 font-semibold">Water Level (m)</label>
              <input
                type="number"
                step="0.1"
                required
                value={waterLevelM}
                onChange={(e) => setWaterLevelM(parseFloat(e.target.value))}
                className="w-full px-2.5 py-1.5 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[10px] text-slate-400 mb-1 font-semibold">Vibration (mm/s)</label>
              <input
                type="number"
                step="0.1"
                required
                value={bearingVibration}
                onChange={(e) => setBearingVibration(parseFloat(e.target.value))}
                className="w-full px-2.5 py-1.5 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[10px] text-slate-400 mb-1 font-semibold">Stator Temp (°C)</label>
              <input
                type="number"
                step="0.5"
                required
                value={statorTempC}
                onChange={(e) => setStatorTempC(parseFloat(e.target.value))}
                className="w-full px-2.5 py-1.5 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[10px] text-slate-400 mb-1 font-semibold">Lube Oil (bar)</label>
              <input
                type="number"
                step="0.1"
                required
                value={lubeOilBar}
                onChange={(e) => setLubeOilBar(parseFloat(e.target.value))}
                className="w-full px-2.5 py-1.5 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
          </div>
        </div>

        {/* Section 3: Remote Personnel Muster */}
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4 space-y-3">
          <label className="text-xs font-bold text-slate-200 block">
            5. Remote Station Personnel Shift Muster (Voice Verified)
          </label>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div>
              <label className="block text-[11px] text-slate-400 mb-1 font-semibold">Station Operators on Duty</label>
              <input
                type="text"
                required
                value={operatorNames}
                onChange={(e) => setOperatorNames(e.target.value)}
                placeholder="Ibrahim Pam, Yusuf Chuwang"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white"
              />
            </div>
            <div>
              <label className="block text-[11px] text-slate-400 mb-1 font-semibold">Protection Engineers</label>
              <input
                type="text"
                required
                value={protectionEng}
                onChange={(e) => setProtectionEng(e.target.value)}
                placeholder="Engr. D. Gyang"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white"
              />
            </div>
            <div>
              <label className="block text-[11px] text-slate-400 mb-1 font-semibold">Mechanical/Electrical Techs</label>
              <input
                type="text"
                required
                value={technicians}
                onChange={(e) => setTechnicians(e.target.value)}
                placeholder="Bulus Mark, Sunday Azi"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white"
              />
            </div>
          </div>
        </div>

        {/* Section 4: Handover Notes */}
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">
            6. Inter-Shift Operational Handover Notes &amp; Safety Advisories
          </label>
          <textarea
            rows={3}
            value={handoverNotes}
            onChange={(e) => setHandoverNotes(e.target.value)}
            placeholder="Log penstock valve adjustments, governor response, hydraulic gate status, weather conditions, or any pending workshop repairs..."
            className="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-cyan-500"
          />
        </div>

        <div className="flex justify-end pt-2">
          <button
            type="submit"
            disabled={submitting}
            className="px-6 py-2.5 bg-cyan-600 hover:bg-cyan-500 disabled:opacity-50 text-white font-bold text-xs rounded-lg transition shadow-lg flex items-center gap-2"
          >
            {submitting ? "Collation in Progress..." : "📻 Collate Shift Log & Sync to HQ HR"}
          </button>
        </div>
      </form>

      {/* Recent Muster History Table */}
      <div className="pt-4 border-t border-slate-800 space-y-3">
        <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
          Recent Shift Muster &amp; Handover Audit Trail
        </h3>

        <div className="overflow-x-auto">
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Shift Window</TableHeaderCell>
                <TableHeaderCell>Plant</TableHeaderCell>
                <TableHeaderCell>Operators on Duty</TableHeaderCell>
                <TableHeaderCell>Protection Eng</TableHeaderCell>
                <TableHeaderCell>Technicians</TableHeaderCell>
                <TableHeaderCell>HQ HR Sync</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {rollcalls.map((rc) => (
                <TableRow key={rc.id}>
                  <TableCell className="font-mono text-xs font-bold text-cyan-400">
                    {rc.shift_window}
                  </TableCell>
                  <TableCell className="font-semibold text-white">
                    {rc.plant_id.replace(/_/g, " ")}
                  </TableCell>
                  <TableCell className="text-slate-300">
                    {rc.operator_names.join(", ")}
                  </TableCell>
                  <TableCell className="text-slate-300">
                    {rc.protection_engineers.join(", ")}
                  </TableCell>
                  <TableCell className="text-slate-300">
                    {rc.technicians.join(", ")}
                  </TableCell>
                  <TableCell>
                    <span className="text-[11px] font-bold text-emerald-400">
                      ✓ Synced
                    </span>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  );
};
