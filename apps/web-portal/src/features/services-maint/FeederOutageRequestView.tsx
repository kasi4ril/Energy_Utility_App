import React, { useState } from "react";
import { FeederCorridor, OutageType, OutageStatus, FeederClosureOutcome, FeederOutageRequest } from "./types";
import { MOCK_FEEDER_OUTAGE_REQUESTS } from "./mockData";
import { FeederOutageRegisterTable } from "./FeederOutageRegisterTable";

export const FeederOutageRequestView: React.FC = () => {
  const [requests, setRequests] = useState<FeederOutageRequest[]>(MOCK_FEEDER_OUTAGE_REQUESTS);
  const [corridor, setCorridor] = useState<FeederCorridor>("FEEDER_33KV_ANGLO_JOS");
  const [outageType, setOutageType] = useState<OutageType>("PLANNED_PREVENTIVE_MAINTENANCE");
  const [urgency, setUrgency] = useState<"EMERGENCY" | "URGENT" | "ROUTINE">("ROUTINE");
  const [reqDate, setReqDate] = useState("2026-10-25");
  const [startTime, setStartTime] = useState("09:00 WAT");
  const [durationHours, setDurationHours] = useState<number>(3);
  const [openingTime, setOpeningTime] = useState("");
  const [surrenderingTime, setSurrenderingTime] = useState("");
  const [actionTaken, setActionTaken] = useState("");
  const [closureOutcome, setClosureOutcome] = useState<FeederClosureOutcome>("PENDING_SURRENDER");
  const [loadMw, setLoadMw] = useState<number>(4.5);
  const [affectedDts, setAffectedDts] = useState<number>(28);
  const [isolator, setIsolator] = useState("ISO-AJ-04 (Makeri Gantry)");
  const [earthSwitch, setEarthSwitch] = useState("ES-AJ-04A");
  const [justification, setJustification] = useState("");
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const isExhausted = closureOutcome !== "PENDING_SURRENDER";
    const newReq: FeederOutageRequest = {
      id: `OUT-${Date.now()}`,
      request_no: `PTW-REQ-2026-${Math.floor(100 + Math.random() * 900)}`,
      corridor,
      outage_type: outageType,
      urgency,
      requested_date: reqDate,
      start_time: startTime,
      duration_hours: Number(durationHours),
      opening_time_wat: openingTime || undefined,
      surrendering_time_wat: surrenderingTime || undefined,
      time_returned_wat: surrenderingTime || undefined,
      action_taken_on_return: actionTaken || undefined,
      is_permit_exhausted: isExhausted,
      closure_outcome: closureOutcome,
      estimated_load_mw: Number(loadMw),
      affected_dts_count: Number(affectedDts),
      isolation_isolator: isolator,
      earth_switch: earthSwitch,
      justification: justification || "Scheduled preventive overhaul & right-of-way clearance.",
      requested_by_staff: "Engr. D. Gyang (Unit Head, Feeder Maint)",
      status: urgency === "EMERGENCY" ? "APPROVED_PTW_ISSUED" : "PENDING_DISPATCH",
      ptw_issued_no: urgency === "EMERGENCY" ? `PTW-2026-0${Math.floor(50 + Math.random() * 40)}` : undefined,
      created_at: new Date().toISOString(),
    };
    setRequests([newReq, ...requests]);
    setToastMsg(`Request ${newReq.request_no} submitted to Central Dispatch & Safety.`);
    setJustification("");
    setOpeningTime("");
    setSurrenderingTime("");
    setActionTaken("");
    setTimeout(() => setToastMsg(null), 5000);
  };

  const handleRecordReturn = (reqId: string, returnTime: string, actionText: string, outcome: FeederClosureOutcome) => {
    setRequests((prev) =>
      prev.map((r) =>
        r.id === reqId
          ? {
              ...r,
              time_returned_wat: returnTime,
              surrendering_time_wat: returnTime,
              action_taken_on_return: actionText,
              closure_outcome: outcome,
              is_permit_exhausted: true,
              status: outcome === "CLOSED_ENERGIZED_IMMEDIATELY" || outcome === "LEFT_ON_LOAD_SHEDDING"
                ? "COMPLETED_LINE_ENERGIZED"
                : r.status,
            }
          : r
      )
    );
    setToastMsg("Permit return time, action taken, and closure state logged in Return Register.");
    setTimeout(() => setToastMsg(null), 5000);
  };

  const isApproved = (s: OutageStatus) => s === "APPROVED_PTW_ISSUED" || s === "ISOLATED_IN_PROGRESS" || s === "COMPLETED_LINE_ENERGIZED";
  const isDeclined = (s: OutageStatus) => s === "REJECTED_BY_DISPATCH";
  const isPending = (s: OutageStatus) => s === "PENDING_DISPATCH";

  const approvedCount = requests.filter((r) => isApproved(r.status)).length;
  const pendingCount = requests.filter((r) => isPending(r.status)).length;
  const declinedCount = requests.filter((r) => isDeclined(r.status)).length;

  return (
    <div className="space-y-6">
      {/* Header Banner */}
      <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div className="inline-flex items-center gap-2 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-300 border border-amber-500/30 mb-2">
            <span className="w-2 h-2 rounded-full bg-amber-400 animate-pulse" />
            NERC Grid Code Part 3 • Safety Protocol PTW-09
          </div>
          <h2 className="text-xl font-bold text-white tracking-tight">33kV / 11kV Feeder Outage Requisition &amp; Return</h2>
          <p className="text-sm text-slate-400 mt-1">Formal isolation requisition, operational grant/return timing, and post-closure outcome logging.</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="bg-slate-950/70 border border-slate-800 rounded-lg px-3 py-2 text-center">
            <div className="text-[10px] text-emerald-400 uppercase font-semibold">Approved</div>
            <div className="text-lg font-bold text-emerald-400">{approvedCount}</div>
          </div>
          <div className="bg-slate-950/70 border border-slate-800 rounded-lg px-3 py-2 text-center">
            <div className="text-[10px] text-amber-400 uppercase font-semibold">Pending</div>
            <div className="text-lg font-bold text-amber-400">{pendingCount}</div>
          </div>
          <div className="bg-slate-950/70 border border-slate-800 rounded-lg px-3 py-2 text-center">
            <div className="text-[10px] text-rose-400 uppercase font-semibold">Declined</div>
            <div className="text-lg font-bold text-rose-400">{declinedCount}</div>
          </div>
        </div>
      </div>

      {toastMsg && (
        <div className="bg-emerald-950/80 border border-emerald-600/50 text-emerald-200 px-4 py-3 rounded-lg text-sm flex items-center justify-between shadow-md">
          <span>✓ {toastMsg}</span>
          <button onClick={() => setToastMsg(null)} className="text-emerald-400 hover:text-emerald-100 font-bold ml-4">✕</button>
        </div>
      )}

      {/* Grid: Form (5 cols) & Register Table (7 cols) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Outage Request Form */}
        <div className="lg:col-span-5 bg-slate-900/70 border border-slate-800 rounded-xl p-5 shadow-sm space-y-4">
          <div className="border-b border-slate-800 pb-3">
            <h3 className="text-base font-semibold text-white">Outage Requisition Form</h3>
            <p className="text-xs text-slate-400 mt-0.5">Permit granting, return on exhaustion, and post-closure state.</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-3 text-xs">
            <div>
              <label className="block text-slate-300 font-semibold mb-1">Target Feeder Corridor</label>
              <select
                value={corridor}
                onChange={(e) => setCorridor(e.target.value as FeederCorridor)}
                style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                className="w-full px-3 py-2 rounded-lg border border-slate-700 text-xs focus:outline-none"
              >
                <option value="FEEDER_33KV_ANGLO_JOS" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Anglo-Jos 33kV (Makeri - Kwall)</option>
                <option value="FEEDER_33KV_BUKURU" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Bukuru 33kV (Bukuru Substation)</option>
                <option value="FEEDER_33KV_RAYFIELD" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Rayfield 33kV (Residential/Gov)</option>
                <option value="FEEDER_33KV_MAKERI" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Makeri 33kV Industrial Tie Line</option>
                <option value="FEEDER_33KV_BOKKOS" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Bokkos 33kV Intertie</option>
                <option value="FEEDER_11KV_VOM" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Vom 11kV Radial (NVRI / NIPSS)</option>
                <option value="FEEDER_11KV_ZAWAN" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Zawan 11kV Feeder Corridor</option>
              </select>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Outage Category</label>
                <select
                  value={outageType}
                  onChange={(e) => setOutageType(e.target.value as OutageType)}
                  style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                  className="w-full px-2.5 py-1.5 rounded-lg border border-slate-700 text-xs"
                >
                  <option value="PLANNED_PREVENTIVE_MAINTENANCE" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Preventive Maint.</option>
                  <option value="EMERGENCY_FAULT_REPAIR" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Emergency Fault</option>
                  <option value="VEGETATION_BUSH_CLEARING" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Vegetation Clearing</option>
                  <option value="RECONDUCTORING_PROJECT" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Reconductoring</option>
                </select>
              </div>
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Priority Urgency</label>
                <select
                  value={urgency}
                  onChange={(e) => setUrgency(e.target.value as any)}
                  style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                  className="w-full px-2.5 py-1.5 rounded-lg border border-slate-700 text-xs"
                >
                  <option value="ROUTINE" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Routine (24h)</option>
                  <option value="URGENT" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Urgent (4h)</option>
                  <option value="EMERGENCY" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Emergency (Now)</option>
                </select>
              </div>
            </div>

            <div className="grid grid-cols-3 gap-2">
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Date</label>
                <input type="date" value={reqDate} onChange={(e) => setReqDate(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs" required />
              </div>
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Start Time</label>
                <input type="text" value={startTime} onChange={(e) => setStartTime(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs" placeholder="09:00 WAT" required />
              </div>
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Hours</label>
                <input type="number" step="0.5" value={durationHours} onChange={(e) => setDurationHours(parseFloat(e.target.value) || 1)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs" required />
              </div>
            </div>

            {/* Granting & Return Timings */}
            <div className="grid grid-cols-2 gap-3 bg-slate-950/60 p-2.5 rounded-lg border border-slate-800">
              <div>
                <label className="block text-emerald-400 font-semibold mb-1">Granting Time (Opened)</label>
                <input
                  type="text"
                  value={openingTime}
                  onChange={(e) => setOpeningTime(e.target.value)}
                  placeholder="e.g. 13:15 WAT"
                  style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                  className="w-full px-2.5 py-1.5 rounded border border-slate-700 text-xs focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-amber-400 font-semibold mb-1">Time Returned (Exhausted)</label>
                <input
                  type="text"
                  value={surrenderingTime}
                  onChange={(e) => setSurrenderingTime(e.target.value)}
                  placeholder="e.g. 15:45 WAT"
                  style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                  className="w-full px-2.5 py-1.5 rounded border border-slate-700 text-xs focus:outline-none"
                />
              </div>
            </div>

            {/* Action Taken on Return */}
            <div>
              <label className="block text-slate-300 font-semibold mb-1">Action Taken on Return / Permit Exhaustion</label>
              <input
                type="text"
                value={actionTaken}
                onChange={(e) => setActionTaken(e.target.value)}
                placeholder="e.g. Insulator replaced, earths removed, line reclosed."
                style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                className="w-full px-2.5 py-1.5 rounded border border-slate-700 text-xs focus:outline-none"
              />
            </div>

            {/* Feeder Closure Outcome */}
            <div>
              <label className="block text-slate-300 font-semibold mb-1">Feeder Closure State (Post-Return)</label>
              <select
                value={closureOutcome}
                onChange={(e) => setClosureOutcome(e.target.value as FeederClosureOutcome)}
                style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
                className="w-full px-2.5 py-1.5 rounded-lg border border-slate-700 text-xs font-semibold focus:outline-none"
              >
                <option value="CLOSED_ENERGIZED_IMMEDIATELY" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                  ✓ Closed &amp; Energized Immediately
                </option>
                <option value="LEFT_ON_LOAD_SHEDDING" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                  ⚠ Left on Load Shedding (Rationing)
                </option>
                <option value="TRIPPED_ON_CLOSURE" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                  ⚡ Tripped on Closure (Fault / Inrush)
                </option>
                <option value="PENDING_SURRENDER" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                  ⏳ Pending Return / Work Ongoing
                </option>
              </select>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Est. Load (MW)</label>
                <input type="number" step="0.1" value={loadMw} onChange={(e) => setLoadMw(parseFloat(e.target.value) || 0)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs" />
              </div>
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Affected DTs</label>
                <input type="number" value={affectedDts} onChange={(e) => setAffectedDts(parseInt(e.target.value, 10) || 0)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2.5 py-1.5 rounded border border-slate-700 text-xs" />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Isolator Point</label>
                <input type="text" value={isolator} onChange={(e) => setIsolator(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs" required />
              </div>
              <div>
                <label className="block text-slate-300 font-semibold mb-1">Earth Switch</label>
                <input type="text" value={earthSwitch} onChange={(e) => setEarthSwitch(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs" required />
              </div>
            </div>

            <div>
              <label className="block text-slate-300 font-semibold mb-1">Work Scope &amp; Notes</label>
              <textarea rows={2} value={justification} onChange={(e) => setJustification(e.target.value)} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }} className="w-full px-3 py-1.5 rounded border border-slate-700 text-xs" placeholder="Specify spans, clearances..." />
            </div>

            <button type="submit" className="w-full py-2 bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 text-slate-950 font-bold rounded-lg text-xs shadow transition-all">
              Dispatch PTW Outage Request to Central Control
            </button>
          </form>
        </div>

        {/* Outage Request and Return Register Table */}
        <div className="lg:col-span-7">
          <FeederOutageRegisterTable requests={requests} onRecordReturn={handleRecordReturn} />
        </div>
      </div>
    </div>
  );
};
export default FeederOutageRequestView;
