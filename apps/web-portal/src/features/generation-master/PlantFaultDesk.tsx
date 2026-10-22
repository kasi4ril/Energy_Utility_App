/**
 * Task View 3: Plant Electrical & Mechanical Fault Desk with GM Approval Gate.
 * Coordinates outage triage, engineering wing dispatch, and GM material authorizations.
 */

import React, { useState } from "react";
import { HydroPlantSummary, PlantFaultJob, PlantID, TechnicalWing } from "./types";
import { approveGMFault, reportPlantFault } from "./api";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";

interface FaultDeskProps {
  plants: HydroPlantSummary[];
  faults: PlantFaultJob[];
  onSuccess?: () => void;
}

export const PlantFaultDesk: React.FC<FaultDeskProps> = ({ plants, faults, onSuccess }) => {
  const [showReportForm, setShowReportForm] = useState<boolean>(false);
  const [selectedPlantId, setSelectedPlantId] = useState<PlantID>("PLANT_06_JERE");
  const [wing, setWing] = useState<TechnicalWing>("MECHANICAL");
  const [severity, setSeverity] = useState<"LOW" | "MEDIUM" | "HIGH" | "CRITICAL">("HIGH");
  const [description, setDescription] = useState<string>("");
  const [requiresWorkshop, setRequiresWorkshop] = useState<boolean>(true);
  const [approvingFault, setApprovingFault] = useState<PlantFaultJob | null>(null);
  const [gmNotes, setGmNotes] = useState<string>("");
  const [submitting, setSubmitting] = useState<boolean>(false);

  const handleReportFault = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await reportPlantFault({
        plant_id: selectedPlantId,
        wing,
        severity,
        fault_description: description,
        requires_workshop: requiresWorkshop,
        reported_by_operator_id: "STAFF-GEN-OP-01",
      });
      setShowReportForm(false);
      setDescription("");
      if (onSuccess) onSuccess();
    } catch (err: any) {
      alert(err.message || "Failed to log fault");
    } finally {
      setSubmitting(false);
    }
  };

  const handleExecuteGMApproval = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!approvingFault) return;
    setSubmitting(true);
    try {
      await approveGMFault(
        approvingFault.id,
        gmNotes || "GM Generation authorized spare parts and store requisition.",
        "STAFF-GM-GEN-01"
      );
      setApprovingFault(null);
      setGmNotes("");
      if (onSuccess) onSuccess();
    } catch (err: any) {
      alert(err.message || "Failed to approve fault");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-800 gap-3">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-bold text-white tracking-wide">
              🚨 Plant Electrical &amp; Mechanical Fault Desk
            </h2>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-rose-500/10 text-rose-300 border border-rose-500/20">
              {faults.length} Plant Incident Logs
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Outage dispatch to Electrical &amp; Mechanical Unit Heads. High-value materials require GM Generation sign-off.
          </p>
        </div>

        <button
          onClick={() => setShowReportForm(!showReportForm)}
          className="px-3.5 py-1.5 rounded-lg bg-rose-600 hover:bg-rose-500 text-white text-xs font-bold transition shadow"
        >
          {showReportForm ? "✕ Cancel" : "+ Log Plant Trip / Fault"}
        </button>
      </div>

      {/* New Fault Report Form */}
      {showReportForm && (
        <form onSubmit={handleReportFault} className="bg-slate-950/80 border border-slate-800 rounded-xl p-5 space-y-4 text-xs">
          <h3 className="text-xs font-bold text-slate-200 uppercase tracking-wider">
            Log New Hydro Generating Station Outage / Trip
          </h3>

          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-3">
            <div>
              <label className="block text-[11px] text-slate-400 mb-1">Generating Station</label>
              <select
                value={selectedPlantId}
                onChange={(e) => setSelectedPlantId(e.target.value as PlantID)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white"
              >
                {plants.map((p) => (
                  <option key={p.id} value={p.id}>{p.name}</option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-[11px] text-slate-400 mb-1">Engineering Wing</label>
              <select
                value={wing}
                onChange={(e) => setWing(e.target.value as TechnicalWing)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white"
              >
                <option value="MECHANICAL">Mechanical Wing (Turbine / Penstock / Bearings)</option>
                <option value="ELECTRICAL">Electrical Wing (Alternator / AVR / Relays)</option>
              </select>
            </div>

            <div>
              <label className="block text-[11px] text-slate-400 mb-1">Severity Rating</label>
              <select
                value={severity}
                onChange={(e) => setSeverity(e.target.value as any)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white"
              >
                <option value="CRITICAL">Critical (Total Unit Shutdown)</option>
                <option value="HIGH">High (Derated / Vibration Warning)</option>
                <option value="MEDIUM">Medium (Minor Leak / Temperature)</option>
                <option value="LOW">Low (Preventive Inspection)</option>
              </select>
            </div>

            <div className="flex items-center pt-5">
              <label className="flex items-center gap-2 cursor-pointer text-slate-300 font-semibold">
                <input
                  type="checkbox"
                  checked={requiresWorkshop}
                  onChange={(e) => setRequiresWorkshop(e.target.checked)}
                  className="rounded bg-slate-900 border-slate-700 text-cyan-500 focus:ring-cyan-500"
                />
                <span>Requires Central Workshop Lathe/Winding</span>
              </label>
            </div>
          </div>

          <div>
            <label className="block text-[11px] text-slate-400 mb-1">Detailed Fault Narrative *</label>
            <textarea
              required
              rows={3}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Describe physical damage, vibration level, oil pressure drops, or protective relay codes..."
              className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500"
            />
          </div>

          <div className="flex justify-end gap-2 pt-1">
            <button
              type="submit"
              disabled={submitting}
              className="px-5 py-2 rounded bg-rose-600 hover:bg-rose-500 text-white font-bold transition shadow"
            >
              {submitting ? "Dispatching..." : "🚨 Dispatch Fault to Wing Head"}
            </button>
          </div>
        </form>
      )}

      {/* Fault Queue Table */}
      <div className="overflow-x-auto">
        <Table>
          <TableHead>
            <tr>
              <TableHeaderCell>Job #</TableHeaderCell>
              <TableHeaderCell>Station</TableHeaderCell>
              <TableHeaderCell>Wing</TableHeaderCell>
              <TableHeaderCell>Severity</TableHeaderCell>
              <TableHeaderCell>Workshop</TableHeaderCell>
              <TableHeaderCell>GM Approval Gate</TableHeaderCell>
              <TableHeaderCell>Status</TableHeaderCell>
            </tr>
          </TableHead>
          <TableBody>
            {faults.map((f) => (
              <TableRow key={f.id}>
                <TableCell className="font-mono text-xs font-bold text-rose-400">
                  {f.job_no}
                </TableCell>
                <TableCell className="font-semibold text-white">
                  {f.plant_id.replace(/_/g, " ")}
                </TableCell>
                <TableCell>
                  <span className={`text-xs font-bold ${f.wing === "ELECTRICAL" ? "text-cyan-400" : "text-amber-400"}`}>
                    {f.wing}
                  </span>
                </TableCell>
                <TableCell>
                  <Badge variant={f.severity === "CRITICAL" ? "danger" : f.severity === "HIGH" ? "warning" : "neutral"}>
                    {f.severity}
                  </Badge>
                </TableCell>
                <TableCell>
                  {f.requires_workshop ? (
                    <span className="text-[11px] font-mono text-amber-300">
                      ⚙️ {f.workshop_status.replace(/_/g, " ")}
                    </span>
                  ) : (
                    <span className="text-[11px] text-slate-500">Field Repair</span>
                  )}
                </TableCell>
                <TableCell>
                  {f.gm_approved ? (
                    <span className="text-[11px] font-bold text-emerald-400">
                      ✓ GM Authorized
                    </span>
                  ) : (
                    <button
                      onClick={() => setApprovingFault(f)}
                      className="px-2.5 py-1 rounded bg-amber-600 hover:bg-amber-500 text-white text-[11px] font-bold shadow animate-pulse"
                    >
                      Sign GM Approval &rarr;
                    </button>
                  )}
                </TableCell>
                <TableCell>
                  <Badge variant={f.status === "RESOLVED_REINSTALLED" ? "success" : "info"}>
                    {f.status.replace(/_/g, " ")}
                  </Badge>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* GM Approval Modal */}
      {approvingFault && (
        <div className="fixed inset-0 z-50 bg-black/80 flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-700 rounded-xl max-w-lg w-full p-6 space-y-4 shadow-2xl">
            <div className="flex items-center justify-between border-b border-slate-800 pb-3">
              <h3 className="text-sm font-bold text-white flex items-center gap-2">
                <span>✍️ General Manager (Generation) Authorization Gate</span>
              </h3>
              <button onClick={() => setApprovingFault(null)} className="text-slate-400 hover:text-white text-xs font-bold">
                ✕ Close
              </button>
            </div>

            <div className="bg-amber-500/10 border border-amber-500/30 rounded-lg p-3 text-xs text-amber-200">
              <p className="font-semibold">Executive Procurement Authority:</p>
              <p className="text-[11px] text-amber-300 mt-1">
                Authorizes release of high-value raw materials (bronze alloys, stainless rods, thyristors) from Central Stores for{" "}
                <span className="font-bold underline">{approvingFault.job_no}</span> at {approvingFault.plant_id}.
              </p>
            </div>

            <form onSubmit={handleExecuteGMApproval} className="space-y-3 text-xs pt-1">
              <div>
                <label className="block text-[11px] font-semibold text-slate-300 mb-1">
                  GM Executive Approval Endorsement:
                </label>
                <textarea
                  required
                  rows={3}
                  value={gmNotes}
                  onChange={(e) => setGmNotes(e.target.value)}
                  placeholder="e.g. Approved emergency requisition for turbine runner overhaul. Stores & Internal Audit authorized."
                  className="w-full px-3 py-2 rounded bg-slate-950 border border-slate-700 text-white placeholder-slate-500"
                />
              </div>

              <div className="flex justify-end gap-2 pt-2">
                <button
                  type="button"
                  onClick={() => setApprovingFault(null)}
                  className="px-3 py-1.5 rounded bg-slate-800 text-slate-300"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="px-4 py-1.5 rounded bg-emerald-600 hover:bg-emerald-500 font-bold text-white shadow"
                >
                  {submitting ? "Signing..." : "Authorize Requisition (GM Sign-Off)"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
