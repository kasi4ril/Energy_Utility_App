/**
 * Task View 4: Central Machine Workshop & Turbine Component Fabrication Desk.
 * Manages in-house lathe machining, runner dynamic balancing, and stator coil rewinding.
 */

import React, { useState } from "react";
import { PlantFaultJob, WorkshopStatus } from "./types";
import { updateWorkshopJob } from "./api";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";

interface WorkshopViewProps {
  faults: PlantFaultJob[];
  onSuccess?: () => void;
}

export const WorkshopFabricationView: React.FC<WorkshopViewProps> = ({ faults, onSuccess }) => {
  const workshopJobs = faults.filter((f) => f.requires_workshop);
  const [updatingJob, setUpdatingJob] = useState<PlantFaultJob | null>(null);
  const [newStatus, setNewStatus] = useState<WorkshopStatus>("QUALITY_CHECK");
  const [materials, setMaterials] = useState<string>("");
  const [submitting, setSubmitting] = useState<boolean>(false);

  const handleUpdateStatus = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!updatingJob) return;
    setSubmitting(true);
    try {
      await updateWorkshopJob(
        updatingJob.id,
        newStatus,
        materials || updatingJob.raw_materials_summary || "Machining completed"
      );
      setUpdatingJob(null);
      if (onSuccess) onSuccess();
    } catch (err: any) {
      alert(err.message || "Failed to update workshop job");
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
              ⚙️ Central Machine Workshop &amp; Turbine Component Fabrication
            </h2>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-amber-500/10 text-amber-300 border border-amber-500/20">
              {workshopJobs.length} Active Fabrication Orders
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            In-house lathe machining, Pelton/Francis runner dynamic balancing, bearing re-babbitting, and stator coil rewinding.
          </p>
        </div>
      </div>

      {/* Workshop Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4">
          <span className="text-[11px] font-semibold text-slate-400">Heavy Machine Lathes</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-xl font-bold text-white font-mono">2 Lathes Active</span>
            <span className="text-[11px] text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded">Shaft Sleeves</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-1">Precision tolerance &le;0.02mm</p>
        </div>

        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4">
          <span className="text-[11px] font-semibold text-slate-400">Stator Winding Bay</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-xl font-bold text-white font-mono">Vacuum Tank Ready</span>
            <span className="text-[11px] text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded">Class H Resin</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-1">11kV / 3.3kV coil insulation</p>
        </div>

        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4">
          <span className="text-[11px] font-semibold text-slate-400">Dynamic Balancing Rig</span>
          <div className="mt-1 flex items-baseline justify-between">
            <span className="text-xl font-bold text-white font-mono">Runner Rig Online</span>
            <span className="text-[11px] text-purple-400 bg-purple-500/10 px-2 py-0.5 rounded">ISO 1940 G2.5</span>
          </div>
          <p className="text-[11px] text-slate-500 mt-1">Vibration mitigation before installation</p>
        </div>
      </div>

      {/* Workshop Jobs Table */}
      <div className="overflow-x-auto">
        <Table>
          <TableHead>
            <tr>
              <TableHeaderCell>Job Order #</TableHeaderCell>
              <TableHeaderCell>Station Asset</TableHeaderCell>
              <TableHeaderCell>Fabrication Scope</TableHeaderCell>
              <TableHeaderCell>Raw Materials Consumed</TableHeaderCell>
              <TableHeaderCell>Workshop Lifecycle</TableHeaderCell>
              <TableHeaderCell>Action</TableHeaderCell>
            </tr>
          </TableHead>
          <TableBody>
            {workshopJobs.map((j) => (
              <TableRow key={j.id}>
                <TableCell className="font-mono text-xs font-bold text-amber-400">
                  {j.job_no}
                </TableCell>
                <TableCell className="font-semibold text-white">
                  {j.plant_id.replace(/_/g, " ")}
                </TableCell>
                <TableCell className="text-xs text-slate-300 max-w-xs">
                  {j.fault_description}
                </TableCell>
                <TableCell className="font-mono text-[11px] text-slate-400">
                  {j.raw_materials_summary || "Awaiting store issue"}
                </TableCell>
                <TableCell>
                  <Badge
                    variant={
                      j.workshop_status === "COMPLETED_DELIVERED"
                        ? "success"
                        : j.workshop_status === "QUALITY_CHECK"
                        ? "info"
                        : "warning"
                    }
                  >
                    {j.workshop_status.replace(/_/g, " ")}
                  </Badge>
                </TableCell>
                <TableCell>
                  <button
                    onClick={() => {
                      setUpdatingJob(j);
                      setMaterials(j.raw_materials_summary || "");
                    }}
                    className="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-amber-300 text-[11px] font-semibold border border-slate-700 transition"
                  >
                    Update Stage &rarr;
                  </button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Update Workshop Modal */}
      {updatingJob && (
        <div className="fixed inset-0 z-50 bg-black/80 flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-700 rounded-xl max-w-lg w-full p-6 space-y-4 shadow-2xl">
            <div className="flex items-center justify-between border-b border-slate-800 pb-3">
              <h3 className="text-sm font-bold text-white flex items-center gap-2">
                <span>⚙️ Update Workshop Fabrication Progress</span>
              </h3>
              <button onClick={() => setUpdatingJob(null)} className="text-slate-400 hover:text-white text-xs font-bold">
                ✕ Close
              </button>
            </div>

            <form onSubmit={handleUpdateStatus} className="space-y-3 text-xs pt-1">
              <div>
                <label className="block text-[11px] font-semibold text-slate-300 mb-1">
                  Fabrication Lifecycle Status:
                </label>
                <select
                  value={newStatus}
                  onChange={(e) => setNewStatus(e.target.value as WorkshopStatus)}
                  className="w-full px-3 py-2 rounded bg-slate-950 border border-slate-700 text-white"
                >
                  <option value="MACHINING_IN_PROGRESS">Machining in Progress (Lathe / Milling)</option>
                  <option value="WINDING_FABRICATION">Winding / Dynamic Balancing in Progress</option>
                  <option value="QUALITY_CHECK">Quality Inspection &amp; Pressure Testing</option>
                  <option value="COMPLETED_DELIVERED">Completed &amp; Dispatched to Plant Site</option>
                </select>
              </div>

              <div>
                <label className="block text-[11px] font-semibold text-slate-300 mb-1">
                  Raw Materials Consumed (Stores Reconciliation):
                </label>
                <input
                  type="text"
                  value={materials}
                  onChange={(e) => setMaterials(e.target.value)}
                  placeholder="e.g. 120mm Phosphor Bronze, 25kg Class H Resin"
                  className="w-full px-3 py-2 rounded bg-slate-950 border border-slate-700 text-white"
                />
              </div>

              <div className="flex justify-end gap-2 pt-2">
                <button type="button" onClick={() => setUpdatingJob(null)} className="px-3 py-1.5 rounded bg-slate-800 text-slate-300">
                  Cancel
                </button>
                <button type="submit" disabled={submitting} className="px-4 py-1.5 rounded bg-amber-600 hover:bg-amber-500 font-bold text-white shadow">
                  {submitting ? "Updating..." : "Save Workshop Milestone"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
