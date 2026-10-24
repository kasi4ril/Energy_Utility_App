import React, { useState } from "react";
import { FeederOutageRequest, FeederClosureOutcome } from "./types";

interface RecordReturnModalProps {
  request: FeederOutageRequest;
  onClose: () => void;
  onSubmit: (reqId: string, returnTime: string, actionText: string, outcome: FeederClosureOutcome) => void;
}

export const FeederRecordReturnModal: React.FC<RecordReturnModalProps> = ({ request, onClose, onSubmit }) => {
  const [returnTime, setReturnTime] = useState(
    new Date().toLocaleTimeString("en-GB", { hour: "2-digit", minute: "2-digit" }) + " WAT"
  );
  const [actionTaken, setActionTaken] = useState("Work completed. Portable earths removed, line energized.");
  const [closureOutcome, setClosureOutcome] = useState<FeederClosureOutcome>("CLOSED_ENERGIZED_IMMEDIATELY");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(request.id, returnTime, actionTaken, closureOutcome);
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="bg-slate-900 border border-slate-700 rounded-xl p-5 max-w-md w-full space-y-3.5 shadow-2xl">
        <div className="flex items-center justify-between border-b border-slate-800 pb-2">
          <h3 className="text-sm font-bold text-white">Record Permit Return &amp; Action Taken</h3>
          <button onClick={onClose} className="text-slate-400 hover:text-white font-bold">✕</button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-3 text-xs">
          <div className="bg-slate-950 p-2.5 rounded border border-slate-800 text-[11px] space-y-1">
            <div className="text-white font-bold">{request.request_no} • {request.corridor.replace("FEEDER_", "")}</div>
            <div className="text-emerald-400 font-mono">Granted: {request.opening_time_wat || request.start_time}</div>
            <div className="text-slate-400">Scheduled Duration: {request.duration_hours} hours</div>
          </div>
          <div>
            <label className="block text-slate-300 font-semibold mb-1">Time Permit Returned / Exhausted</label>
            <input
              type="text"
              value={returnTime}
              onChange={(e) => setReturnTime(e.target.value)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="w-full px-2.5 py-1.5 rounded border border-slate-700 text-xs focus:outline-none"
              required
            />
          </div>
          <div>
            <label className="block text-slate-300 font-semibold mb-1">Action Taken on Return</label>
            <textarea
              rows={2}
              value={actionTaken}
              onChange={(e) => setActionTaken(e.target.value)}
              placeholder="Describe repairs completed, earths removed, line condition..."
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="w-full px-2.5 py-1.5 rounded border border-slate-700 text-xs focus:outline-none"
              required
            />
          </div>
          <div>
            <label className="block text-slate-300 font-semibold mb-1">Feeder Closure Outcome</label>
            <select
              value={closureOutcome}
              onChange={(e) => setClosureOutcome(e.target.value as FeederClosureOutcome)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="w-full px-2.5 py-1.5 rounded border border-slate-700 text-xs font-semibold focus:outline-none"
            >
              <option value="CLOSED_ENERGIZED_IMMEDIATELY" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                ✓ Closed &amp; Energized Immediately
              </option>
              <option value="LEFT_ON_LOAD_SHEDDING" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                ⚠ Left on Load Shedding (Rationing)
              </option>
              <option value="TRIPPED_ON_CLOSURE" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                ⚡ Tripped on Closure (Fault/Surge)
              </option>
            </select>
          </div>
          <div className="flex justify-end gap-2 pt-1 border-t border-slate-800">
            <button type="button" onClick={onClose} className="px-3 py-1 bg-slate-800 text-slate-300 rounded font-semibold text-xs">
              Cancel
            </button>
            <button type="submit" className="px-4 py-1 bg-amber-500 hover:bg-amber-600 text-slate-950 font-bold rounded text-xs shadow transition-all">
              Log Return in Register
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
export default FeederRecordReturnModal;

