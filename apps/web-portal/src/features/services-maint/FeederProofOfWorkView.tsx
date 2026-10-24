import React, { useState } from "react";
import { FeederProofOfWork } from "./types";
import { MOCK_FEEDER_PROOFS, MOCK_FEEDER_WORK_ORDERS } from "./mockData";

export const FeederProofOfWorkView: React.FC = () => {
  const [proofs, setProofs] = useState<FeederProofOfWork[]>(MOCK_FEEDER_PROOFS);
  const [selectedProof, setSelectedProof] = useState<FeederProofOfWork | null>(MOCK_FEEDER_PROOFS[0] || null);
  const [verifyNotes, setVerifyNotes] = useState("");
  const [showNewModal, setShowNewModal] = useState(false);

  // New Proof Form State
  const [newOrderNo, setNewOrderNo] = useState("WO-2026-FDR-089");
  const [newAssetId, setNewAssetId] = useState("POLE-AJ-142");
  const [newTeamLead, setNewTeamLead] = useState("Engr. S. Pam (Line Crew A)");
  const [newLat, setNewLat] = useState("9.89645");
  const [newLng, setNewLng] = useState("8.85831");
  const [newDetails, setNewDetails] = useState("Replaced cracked pin insulator on crossarm 2; cleaned jumper contacts.");
  const [newMaterials, setNewMaterials] = useState("2x 33kV Pin Insulators, 15m Binding Wire");

  const handleAction = (status: "VERIFIED_APPROVED" | "REJECTED_REVISIT") => {
    if (!selectedProof) return;
    const updated = proofs.map((p) =>
      p.id === selectedProof.id
        ? {
            ...p,
            verification_status: status,
            verified_by_staff_id: "STAFF-ENG-04 (Engr. D. Gyang, Feeder Unit Head)",
            verification_notes: verifyNotes || (status === "VERIFIED_APPROVED" ? "GPS and photographic evidence verified." : "Evidence insufficient. Crew dispatched for rework."),
          }
        : p
    );
    setProofs(updated);
    setSelectedProof(updated.find((p) => p.id === selectedProof.id) || null);
    setVerifyNotes("");
  };

  const handleCreateProof = (e: React.FormEvent) => {
    e.preventDefault();
    const lat = parseFloat(newLat) || 9.89645;
    const lng = parseFloat(newLng) || 8.85831;
    const assetLat = 9.89648;
    const assetLng = 8.85834;
    const distMeters = Math.round(Math.hypot((lat - assetLat) * 111000, (lng - assetLng) * 111000) * 10) / 10;
    const withinTol = distMeters <= 15.0;

    const newEntry: FeederProofOfWork = {
      id: `PRF-${Date.now()}`,
      work_order_id: "WO-FDR-001",
      order_no: newOrderNo,
      asset_id: newAssetId,
      team_lead_name: newTeamLead,
      team_lead_id: "STAFF-ENG-12",
      captured_gps: { lat, lng, accuracy_meters: 2.4 },
      asset_gps: { lat: assetLat, lng: assetLng },
      distance_from_asset_meters: distMeters,
      gps_within_tolerance: withinTol,
      before_photo_url: "https://images.unsplash.com/photo-1544717305-2782549b5136?w=600&auto=format&fit=crop&q=80",
      after_photo_url: "https://images.unsplash.com/photo-1473341304170-971dccb5ac1e?w=600&auto=format&fit=crop&q=80",
      work_details: newDetails,
      materials_consumed: newMaterials.split(",").map((s) => s.trim()),
      submitted_at: new Date().toISOString(),
      verification_status: "PENDING_UNIT_HEAD",
    };
    setProofs([newEntry, ...proofs]);
    setSelectedProof(newEntry);
    setShowNewModal(false);
  };

  return (
    <div className="space-y-6">
      {/* Top Banner */}
      <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div className="inline-flex items-center gap-2 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-300 border border-emerald-500/30 mb-2">
            <span className="w-2 h-2 rounded-full bg-emerald-400" />
            Anti-Ghosting Audit • Geofenced Quality Assurance
          </div>
          <h2 className="text-xl font-bold text-white tracking-tight">Field Proof-of-Work Verification Desk (GPS & Visual Inspection)</h2>
          <p className="text-sm text-slate-400 mt-1">Physical GPS verification (±15m tolerance) and before/after photo audits for all feeder work orders.</p>
        </div>
        <button onClick={() => setShowNewModal(true)} className="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-slate-950 font-bold text-xs rounded-lg shadow transition-all flex items-center gap-2 self-start md:self-auto">
          <span>+ Submit Field Proof Record</span>
        </button>
      </div>

      {/* Main Grid: Left Submissions (4 cols), Right Detail (8 cols) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-4 bg-slate-900/70 border border-slate-800 rounded-xl p-4 shadow-sm space-y-3">
          <div className="flex items-center justify-between border-b border-slate-800 pb-2">
            <h3 className="text-sm font-semibold text-white">Submissions Queue</h3>
            <span className="text-xs px-2 py-0.5 rounded-full bg-slate-800 text-slate-300 font-mono">{proofs.length} records</span>
          </div>

          <div className="space-y-2.5">
            {proofs.map((p) => {
              const isSelected = selectedProof?.id === p.id;
              return (
                <div key={p.id} onClick={() => setSelectedProof(p)} className={`p-3 rounded-lg border cursor-pointer transition-all ${isSelected ? "bg-amber-950/20 border-amber-500/60 shadow-md" : "bg-slate-950/60 border-slate-800/80 hover:bg-slate-800/40"}`}>
                  <div className="flex items-center justify-between">
                    <span className="font-mono font-bold text-xs text-white">{p.order_no}</span>
                    <span className={`text-[10px] font-bold px-2 py-0.5 rounded ${p.verification_status === "VERIFIED_APPROVED" ? "bg-emerald-950 text-emerald-300 border border-emerald-700/50" : p.verification_status === "REJECTED_REVISIT" ? "bg-rose-950 text-rose-300 border border-rose-700/50" : "bg-amber-950 text-amber-300 border border-amber-700/50"}`}>
                      {p.verification_status.replace(/_/g, " ")}
                    </span>
                  </div>
                  <div className="text-xs text-amber-400 font-medium mt-1">Asset: {p.asset_id}</div>
                  <div className="text-[11px] text-slate-400 flex items-center justify-between mt-1">
                    <span>{p.team_lead_name}</span>
                    <span className={`font-mono text-[10px] ${p.gps_within_tolerance ? "text-emerald-400" : "text-rose-400"}`}>
                      {p.distance_from_asset_meters}m {p.gps_within_tolerance ? "✓" : "⚠"}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Right Detail Inspection Desk */}
        <div className="lg:col-span-8 bg-slate-900/70 border border-slate-800 rounded-xl p-5 shadow-sm space-y-4">
          {selectedProof ? (
            <>
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-slate-800 pb-3">
                <div>
                  <div className="text-xs text-amber-400 font-mono font-bold">{selectedProof.order_no}</div>
                  <h3 className="text-lg font-bold text-white mt-0.5">Asset: {selectedProof.asset_id}</h3>
                  <p className="text-xs text-slate-400 mt-0.5">Submitted by {selectedProof.team_lead_name} ({selectedProof.team_lead_id})</p>
                </div>
                <div className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold border ${selectedProof.gps_within_tolerance ? "bg-emerald-950/80 text-emerald-300 border-emerald-500/40" : "bg-rose-950/80 text-rose-300 border-rose-500/40"}`}>
                  <span>{selectedProof.gps_within_tolerance ? "✓ GPS VALIDATED" : "⚠ GPS OUT OF BOUNDS"}</span>
                  <span className="font-mono text-[11px]">({selectedProof.distance_from_asset_meters}m &lt; 15m)</span>
                </div>
              </div>

              {/* Before & After Photos */}
              <div>
                <h4 className="text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">Photographic Evidence (Before & After)</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="bg-slate-950 rounded-lg p-2.5 border border-slate-800">
                    <div className="text-[11px] font-semibold text-slate-400 mb-1 flex items-center justify-between">
                      <span>BEFORE REPAIR</span>
                      <span className="text-rose-400 font-mono">Fault Exposed</span>
                    </div>
                    <div className="aspect-video bg-slate-900 rounded overflow-hidden">
                      <img src={selectedProof.before_photo_url} alt="Before repair" className="w-full h-full object-cover" />
                    </div>
                  </div>
                  <div className="bg-slate-950 rounded-lg p-2.5 border border-slate-800">
                    <div className="text-[11px] font-semibold text-slate-400 mb-1 flex items-center justify-between">
                      <span>AFTER REPAIR</span>
                      <span className="text-emerald-400 font-mono">Rectified & Inspected</span>
                    </div>
                    <div className="aspect-video bg-slate-900 rounded overflow-hidden">
                      <img src={selectedProof.after_photo_url} alt="After repair" className="w-full h-full object-cover" />
                    </div>
                  </div>
                </div>
              </div>

              {/* GPS Telemetry & Scope */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
                <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-3 space-y-1 font-mono">
                  <div className="text-slate-400 font-sans font-semibold text-xs mb-1">GPS Telemetry Check</div>
                  <div className="text-slate-300">Field: {selectedProof.captured_gps.lat.toFixed(5)}°N, {selectedProof.captured_gps.lng.toFixed(5)}°E</div>
                  <div className="text-slate-300">Asset: {selectedProof.asset_gps.lat.toFixed(5)}°N, {selectedProof.asset_gps.lng.toFixed(5)}°E</div>
                  <div className="text-amber-400">Radial Offset: {selectedProof.distance_from_asset_meters}m (Limit: 15m)</div>
                </div>
                <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-3 space-y-1">
                  <div className="text-slate-400 font-semibold text-xs mb-1">Materials Consumed</div>
                  <div className="flex flex-wrap gap-1">
                    {selectedProof.materials_consumed.map((m, idx) => (
                      <span key={idx} className="bg-slate-800 px-2 py-0.5 rounded text-[10px] text-slate-300">• {m}</span>
                    ))}
                  </div>
                  <p className="text-slate-300 text-xs mt-1 leading-snug">{selectedProof.work_details}</p>
                </div>
              </div>

              {/* Verification Desk */}
              <div className="bg-slate-950/90 border border-slate-800 rounded-lg p-3.5 space-y-2.5">
                <div className="text-xs font-semibold text-slate-300 flex items-center justify-between">
                  <span>Unit Head Engineering Sign-Off Desk</span>
                  {selectedProof.verified_by_staff_id && <span className="text-emerald-400 text-[11px] font-mono">{selectedProof.verified_by_staff_id}</span>}
                </div>
                {selectedProof.verification_notes && (
                  <p className="text-xs text-amber-300/90 bg-amber-950/30 p-2 rounded border border-amber-900/50"><strong>Note:</strong> {selectedProof.verification_notes}</p>
                )}
                {selectedProof.verification_status === "PENDING_UNIT_HEAD" && (
                  <div className="space-y-2">
                    <input type="text" placeholder="Audit notes..." value={verifyNotes} onChange={(e) => setVerifyNotes(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-3 py-1.5 text-xs rounded border border-slate-700 focus:outline-none" />
                    <div className="flex items-center gap-3">
                      <button onClick={() => handleAction("VERIFIED_APPROVED")} className="flex-1 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded shadow transition-all">✓ Approve & Mark Attended to</button>
                      <button onClick={() => handleAction("REJECTED_REVISIT")} className="flex-1 py-1.5 bg-rose-600 hover:bg-rose-700 text-white font-bold text-xs rounded shadow transition-all">✕ Reject & Order Field Revisit</button>
                    </div>
                  </div>
                )}
              </div>
            </>
          ) : (
            <div className="text-center py-16 text-slate-500 text-sm">Select a proof record to review.</div>
          )}
        </div>
      </div>

      {/* New Proof Modal */}
      {showNewModal && (
        <div className="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-700 rounded-xl p-5 max-w-lg w-full space-y-3 shadow-2xl">
            <div className="flex items-center justify-between border-b border-slate-800 pb-2">
              <h3 className="text-sm font-bold text-white">Record Field Proof of Work</h3>
              <button onClick={() => setShowNewModal(false)} className="text-slate-400 hover:text-white font-bold">✕</button>
            </div>
            <form onSubmit={handleCreateProof} className="space-y-2.5 text-xs">
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Work Order Reference</label>
                <select value={newOrderNo} onChange={(e) => setNewOrderNo(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2.5 py-1.5 rounded border border-slate-700">
                  {MOCK_FEEDER_WORK_ORDERS.map((wo) => (
                    <option key={wo.id} value={wo.order_no} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>{wo.order_no} - {wo.asset_id} ({wo.fault_type})</option>
                  ))}
                </select>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Asset ID</label>
                  <input type="text" value={newAssetId} onChange={(e) => setNewAssetId(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1 rounded border border-slate-700" required />
                </div>
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Team Lead</label>
                  <input type="text" value={newTeamLead} onChange={(e) => setNewTeamLead(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1 rounded border border-slate-700" required />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Captured Lat (°N)</label>
                  <input type="text" value={newLat} onChange={(e) => setNewLat(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1 rounded border border-slate-700" required />
                </div>
                <div>
                  <label className="block text-slate-300 font-semibold mb-1">Captured Lng (°E)</label>
                  <input type="text" value={newLng} onChange={(e) => setNewLng(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1 rounded border border-slate-700" required />
                </div>
              </div>
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Materials Consumed</label>
                <input type="text" value={newMaterials} onChange={(e) => setNewMaterials(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1 rounded border border-slate-700" />
              </div>
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Work Narrative</label>
                <textarea rows={2} value={newDetails} onChange={(e) => setNewDetails(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1 rounded border border-slate-700" />
              </div>
              <div className="flex justify-end gap-2 pt-1">
                <button type="button" onClick={() => setShowNewModal(false)} className="px-3 py-1 bg-slate-800 text-slate-300 rounded font-semibold">Cancel</button>
                <button type="submit" className="px-4 py-1 bg-amber-500 hover:bg-amber-600 text-slate-950 font-bold rounded">Submit</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
export default FeederProofOfWorkView;

