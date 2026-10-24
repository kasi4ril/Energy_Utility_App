/**
 * Activity 3: Change of Shift Desk & Handover Register.
 * Compliant with NERC Grid Code Handover Regulations:
 * - Shift Window, Incoming & Outgoing Operator Muster
 * - Remote 6-Plant Physical Roll-Call (Operators, Protection Engineers, Techs)
 * - Short column for Handover Notes which, once written, CANNOT BE EDITED.
 * - Permanent immutable audit trail of past shift handovers.
 */

import React, { useState, useEffect } from "react";
import {
  HydroPlantSummary,
  PlantID,
  ShiftWindow,
  ShiftHandoverRecord,
} from "./types";
import { fetchShiftHandovers, submitShiftHandover } from "./api";
import {
  Users,
  Lock,
  ShieldCheck,
  CheckCircle,
  AlertTriangle,
  Clock,
  Send,
  FileCheck,
  Hash,
  AlertCircle,
} from "lucide-react";

interface ChangeOfShiftProps {
  plants: HydroPlantSummary[];
  staffId?: string;
  staffName?: string;
  onSuccess?: () => void;
}

export const ChangeOfShiftForm: React.FC<ChangeOfShiftProps> = ({
  plants,
  staffId = "STAFF-GEN-OP-01",
  staffName = "Ibrahim Pam (Operator)",
  onSuccess,
}) => {
  const [handovers, setHandovers] = useState<ShiftHandoverRecord[]>([]);
  const [loadingHistory, setLoadingHistory] = useState(true);

  // Form State
  const [shiftWindow, setShiftWindow] = useState<ShiftWindow>("AFTERNOON");
  const [plantId, setPlantId] = useState<PlantID>("PLANT_01_KURRA_FALLS");
  const [outgoingOperatorId, setOutgoingOperatorId] = useState<string>(staffId);
  const [outgoingOperatorName, setOutgoingOperatorName] = useState<string>(staffName);
  const [incomingOperatorId, setIncomingOperatorId] = useState<string>("STAFF-GEN-OP-04");
  const [incomingOperatorName, setIncomingOperatorName] = useState<string>("Yusuf Chuwang (Incoming Lead)");
  const [operatorsList, setOperatorsList] = useState<string>("Ibrahim Pam, Godiya Bitrus");
  const [protectionEng, setProtectionEng] = useState<string>("Engr. D. Gyang, Engr. S. Lar");
  const [technicians, setTechnicians] = useState<string>("Bulus Mark, Sunday Azi, Danjuma Moses");
  
  // Handover Notes (Short Column - Immutable once submitted)
  const [handoverNotes, setHandoverNotes] = useState<string>("");
  const [submitting, setSubmitting] = useState(false);
  const [toastMessage, setToastMessage] = useState<string | null>(null);
  const [blockedEditAlert, setBlockedEditAlert] = useState<string | null>(null);

  const loadData = async () => {
    setLoadingHistory(true);
    try {
      const records = await fetchShiftHandovers();
      setHandovers(records);
    } finally {
      setLoadingHistory(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!handoverNotes.trim()) {
      alert("Please provide the short handover notes before submitting.");
      return;
    }

    setSubmitting(true);
    setToastMessage(null);

    try {
      const shiftId = `SH-${shiftWindow}-${new Date().toISOString().slice(0, 10).replace(/-/g, "")}`;
      const newRecord = await submitShiftHandover({
        shift_id: shiftId,
        shift_window: shiftWindow,
        plant_id: plantId,
        outgoing_operator_id: outgoingOperatorId,
        outgoing_operator_name: outgoingOperatorName,
        incoming_operator_id: incomingOperatorId,
        incoming_operator_name: incomingOperatorName,
        operators_on_duty: operatorsList.split(",").map((s) => s.trim()),
        protection_engineers: protectionEng.split(",").map((s) => s.trim()),
        technicians: technicians.split(",").map((s) => s.trim()),
        handover_notes: handoverNotes.trim(),
      });

      setHandovers([newRecord, ...handovers]);
      setHandoverNotes("");
      setToastMessage("Shift Handover successfully completed! The handover notes have been cryptographically locked and recorded into the permanent log.");
      if (onSuccess) onSuccess();
    } catch (err: any) {
      alert(err.message || "Failed to commit shift handover");
    } finally {
      setSubmitting(false);
    }
  };

  const handleAttemptEdit = (recordId: string) => {
    setBlockedEditAlert(`Shift record [${recordId}] is permanently locked. Under NERC Grid Code Regulations and NESCO Master Station Standard Operating Procedures, shift handover notes cannot be altered or deleted once signed.`);
  };

  return (
    <div className="space-y-6 text-xs text-slate-100 max-w-5xl">
      {/* Header Banner */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-xl">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div className="flex items-center space-x-3">
            <div className="p-2.5 rounded-lg bg-amber-500/20 text-amber-400 border border-amber-500/30">
              <Users className="w-6 h-6" />
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <span className="text-[10px] font-bold uppercase tracking-wider text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded border border-amber-500/20">
                  NERC Grid Code Section 5.3 &bull; Change of Shift Protocol
                </span>
                <span className="text-[10px] font-mono text-slate-400">Tamper-Proof Log</span>
              </div>
              <h1 className="text-lg sm:text-xl font-black text-white mt-1">
                Change of Shift Desk &amp; Muster Register
              </h1>
              <p className="text-slate-400 text-xs mt-0.5">
                Incoming and outgoing operator muster with mandatory physical personnel verification across Plants 1–6.
              </p>
            </div>
          </div>

          <div className="bg-amber-500/10 border border-amber-500/30 px-3 py-2 rounded-lg text-right shrink-0">
            <div className="flex items-center space-x-1.5 text-amber-300 font-bold text-xs">
              <Lock className="w-3.5 h-3.5" />
              <span>Immutable Notes Rule</span>
            </div>
            <div className="text-[10px] text-slate-400 mt-0.5">Once written, cannot be edited</div>
          </div>
        </div>
      </div>

      {/* Success Notification */}
      {toastMessage && (
        <div className="p-4 rounded-xl bg-emerald-950/40 border border-emerald-500/60 shadow-lg flex items-start justify-between">
          <div className="flex items-center space-x-3">
            <CheckCircle className="w-5 h-5 text-emerald-400 shrink-0" />
            <div className="text-xs text-emerald-200 font-semibold">{toastMessage}</div>
          </div>
          <button
            onClick={() => setToastMessage(null)}
            className="text-xs text-emerald-400 hover:text-white font-bold ml-4"
          >
            Dismiss
          </button>
        </div>
      )}

      {/* Blocked Edit Alert Notification */}
      {blockedEditAlert && (
        <div className="p-4 rounded-xl bg-amber-950/40 border border-amber-500/60 shadow-lg flex items-start justify-between">
          <div className="flex items-center space-x-3">
            <Lock className="w-5 h-5 text-amber-400 shrink-0" />
            <div className="text-xs text-amber-200 font-semibold leading-relaxed">
              {blockedEditAlert}
            </div>
          </div>
          <button
            onClick={() => setBlockedEditAlert(null)}
            className="text-xs text-amber-400 hover:text-white font-bold ml-4"
          >
            Dismiss
          </button>
        </div>
      )}

      {/* Handover Form */}
      <form onSubmit={handleSubmit} className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-5">
        <div className="flex items-center justify-between pb-2 border-b border-slate-800">
          <h2 className="text-xs font-bold text-white uppercase tracking-wider flex items-center space-x-2">
            <ShieldCheck className="w-4 h-4 text-emerald-400" />
            <span>Shift Handover Authentication</span>
          </h2>
          <span className="text-[10px] text-slate-400 font-mono">Form S05-HANDOVER</span>
        </div>

        {/* Shift Window & Cascade Station */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              Concluded Shift Window *
            </label>
            <select
              value={shiftWindow}
              onChange={(e) => setShiftWindow(e.target.value as ShiftWindow)}
              className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-white focus:outline-none focus:border-cyan-500"
            >
              <option value="MORNING">Morning Shift (08:00 - 16:00)</option>
              <option value="AFTERNOON">Afternoon Shift (16:00 - 00:00)</option>
              <option value="NIGHT">Night Shift (00:00 - 08:00)</option>
            </select>
          </div>

          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              Target Generating Station / Cascade *
            </label>
            <select
              value={plantId}
              onChange={(e) => setPlantId(e.target.value as PlantID)}
              className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-white focus:outline-none focus:border-cyan-500"
            >
              {plants.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              Handover Time (WAT)
            </label>
            <div className="bg-slate-950 border border-slate-800 rounded-lg p-2.5 font-mono text-cyan-400 font-semibold text-xs">
              {new Date().toLocaleTimeString("en-GB", { timeZone: "Africa/Lagos" })} WAT (Live)
            </div>
          </div>
        </div>

        {/* Outgoing & Incoming Operators */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-2">
          <div className="p-3 bg-slate-950/70 border border-slate-800 rounded-lg space-y-2">
            <span className="text-[10px] font-bold text-amber-400 uppercase tracking-wider">
              Outgoing Master Operator (Handing Over)
            </span>
            <div>
              <label className="block text-[10px] text-slate-400 mb-0.5">Staff Name</label>
              <input
                type="text"
                value={outgoingOperatorName}
                onChange={(e) => setOutgoingOperatorName(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-white"
              />
            </div>
            <div>
              <label className="block text-[10px] text-slate-400 mb-0.5">Staff ID</label>
              <input
                type="text"
                value={outgoingOperatorId}
                onChange={(e) => setOutgoingOperatorId(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-cyan-300 font-mono"
              />
            </div>
          </div>

          <div className="p-3 bg-slate-950/70 border border-slate-800 rounded-lg space-y-2">
            <span className="text-[10px] font-bold text-emerald-400 uppercase tracking-wider">
              Incoming Master Operator (Taking Over)
            </span>
            <div>
              <label className="block text-[10px] text-slate-400 mb-0.5">Staff Name</label>
              <input
                type="text"
                value={incomingOperatorName}
                onChange={(e) => setIncomingOperatorName(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-white"
              />
            </div>
            <div>
              <label className="block text-[10px] text-slate-400 mb-0.5">Staff ID</label>
              <input
                type="text"
                value={incomingOperatorId}
                onChange={(e) => setIncomingOperatorId(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-cyan-300 font-mono"
              />
            </div>
          </div>
        </div>

        {/* Remote Physical Muster (Plants 1 to 6) */}
        <div className="p-3.5 bg-slate-950/70 border border-slate-800 rounded-lg space-y-3">
          <span className="text-[10px] font-bold text-cyan-400 uppercase tracking-wider block">
            Radio / Call-Collated Remote Physical Muster (Plants 1–6)
          </span>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label className="block text-[10px] text-slate-400 mb-1">
                Operators on Duty (Comma-separated)
              </label>
              <input
                type="text"
                value={operatorsList}
                onChange={(e) => setOperatorsList(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-white"
              />
            </div>

            <div>
              <label className="block text-[10px] text-slate-400 mb-1">
                Protection Engineers
              </label>
              <input
                type="text"
                value={protectionEng}
                onChange={(e) => setProtectionEng(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-white"
              />
            </div>

            <div>
              <label className="block text-[10px] text-slate-400 mb-1">
                Mechanical &amp; Electrical Techs
              </label>
              <input
                type="text"
                value={technicians}
                onChange={(e) => setTechnicians(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-white"
              />
            </div>
          </div>
        </div>

        {/* Short Column for Handover Notes (IMMUTABLE ONCE WRITTEN) */}
        <div className="p-4 bg-amber-950/20 border border-amber-500/50 rounded-xl space-y-2">
          <div className="flex items-center justify-between">
            <label className="text-xs font-bold text-amber-300 flex items-center space-x-2">
              <Lock className="w-4 h-4 text-amber-400" />
              <span>Short Handover Notes (IMMUTABLE &bull; CANNOT BE EDITED ONCE WRITTEN) *</span>
            </label>
            <span className="text-[10px] font-mono text-amber-400">
              {handoverNotes.length} / 500 chars
            </span>
          </div>

          <p className="text-[11px] text-amber-200/80 leading-normal">
            Enter concise, essential inter-shift directives (e.g. water level trends, bearing temperature watch, workshop jobs dispatched, or isolation safety tags).
            <strong> Note: Under NERC Grid Regulations, this note cannot be edited, altered, or deleted after sign-off.</strong>
          </p>

          <textarea
            rows={3}
            maxLength={500}
            required
            value={handoverNotes}
            onChange={(e) => setHandoverNotes(e.target.value)}
            placeholder="e.g. Kurra Dam water level +0.3m. Station 1 Units 1-3 running normal. Penstock 4 valve packing seep inspected hourly. Plant 6 runner at workshop."
            className="w-full bg-slate-950 border border-amber-500/60 rounded-lg p-3 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-amber-400 font-sans"
          />
        </div>

        {/* Submit Handover */}
        <div className="flex items-center justify-between pt-2">
          <div className="text-[11px] text-slate-400 font-mono flex items-center space-x-1.5">
            <Lock className="w-3.5 h-3.5 text-amber-400" />
            <span>Digital SHA-256 Stamp Applied on Submission</span>
          </div>

          <button
            type="submit"
            disabled={submitting || !handoverNotes.trim()}
            className="px-6 py-2.5 rounded-xl bg-gradient-to-r from-amber-600 to-amber-700 hover:from-amber-500 hover:to-amber-600 text-white font-bold text-xs flex items-center space-x-2 shadow-lg transition-all disabled:opacity-50"
          >
            <Send className="w-4 h-4" />
            <span>Sign &amp; Lock Shift Handover</span>
          </button>
        </div>
      </form>

      {/* Permanent Shift Handover Register (Immutable Historical Records) */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
        <div className="flex items-center justify-between pb-3 border-b border-slate-800">
          <div>
            <h3 className="text-sm font-bold text-white flex items-center space-x-2">
              <Clock className="w-4 h-4 text-cyan-400" />
              <span>Permanent Shift Handover Register (Tamper-Proof Audit Trail)</span>
            </h3>
            <p className="text-[11px] text-slate-400 mt-0.5">
              Historical record of signed shifts. Handover notes are strictly read-only and sealed.
            </p>
          </div>
          <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-cyan-500/10 text-cyan-300 border border-cyan-500/20">
            {handovers.length} Locked Shifts
          </span>
        </div>

        <div className="space-y-3">
          {handovers.map((record) => (
            <div
              key={record.id}
              className="p-4 rounded-xl bg-slate-950/80 border border-slate-800 space-y-2.5 transition"
            >
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                <div className="flex items-center space-x-2">
                  <span className="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-amber-500/20 text-amber-300 border border-amber-500/30">
                    {record.shift_window} SHIFT
                  </span>
                  <span className="font-bold text-white text-xs">{record.plant_id.replace(/_/g, " ")}</span>
                  <span className="text-[10px] text-slate-400 font-mono">({record.shift_id})</span>
                </div>

                <div className="flex items-center space-x-2">
                  <span className="px-2 py-0.5 rounded-full text-[9px] font-bold bg-slate-800 text-emerald-400 border border-emerald-500/30 flex items-center space-x-1">
                    <Lock className="w-3 h-3 text-amber-400" />
                    <span>LOCKED &bull; CANNOT BE EDITED</span>
                  </span>
                  <button
                    onClick={() => handleAttemptEdit(record.id)}
                    className="text-[10px] text-slate-500 hover:text-amber-400 underline"
                  >
                    Edit Note?
                  </button>
                </div>
              </div>

              {/* Personnel Flow */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-[10px] text-slate-400 bg-slate-900/60 p-2.5 rounded-lg">
                <div>
                  Outgoing Operator: <strong className="text-slate-200">{record.outgoing_operator_name}</strong> ({record.outgoing_operator_id})
                </div>
                <div>
                  Incoming Operator: <strong className="text-slate-200">{record.incoming_operator_name}</strong> ({record.incoming_operator_id})
                </div>
              </div>

              {/* Immutable Handover Notes Column */}
              <div className="bg-amber-950/10 border border-amber-500/30 rounded-lg p-3 space-y-1">
                <div className="flex items-center justify-between text-[10px] text-amber-400 font-semibold uppercase tracking-wider">
                  <span className="flex items-center space-x-1">
                    <Lock className="w-3 h-3" />
                    <span>Handover Notes (Immutable)</span>
                  </span>
                  <span className="font-mono text-slate-500 text-[9px]">{record.tamper_hash?.slice(0, 24)}...</span>
                </div>
                <p className="text-xs text-slate-200 font-sans leading-relaxed select-all">
                  &ldquo;{record.handover_notes}&rdquo;
                </p>
              </div>

              {/* Footer Timestamp */}
              <div className="flex items-center justify-between text-[10px] text-slate-500 font-mono pt-1">
                <span>Signed at: {new Date(record.created_at).toLocaleString("en-GB", { timeZone: "Africa/Lagos" })} WAT</span>
                <span className="text-emerald-400 font-semibold">Verified NERC Log</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
